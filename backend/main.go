package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/service"
)

type appConfig struct {
	Port                   string
	APIKey                 string
	AllowedOrigin          string
	RateLimitRequests      int
	RateLimitWindowSeconds int
	RateLimitBurst         int
}

type appServices struct {
	hotpepper *service.HotpepperService
}

// main はアプリ起動処理を run に委譲し、終了時のエラーハンドリングを行う。
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run は設定読み込み、依存構築、ルーター構築、サーバー起動を統括する。
func run() error {
	cfg, err := loadConfigFromEnv()
	if err != nil {
		return err
	}

	svcs := buildServices(cfg)
	newLimiterStore := newLimiterFactory(cfg)
	mux := buildRouter(cfg, svcs, newLimiterStore)

	return startServer(cfg.Port, mux)
}

// loadConfigFromEnv は環境変数の読み込み・デフォルト適用・必須値検証を行う。
func loadConfigFromEnv() (appConfig, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("HOTPEPPER_API_KEY")
	if apiKey == "" {
		return appConfig{}, errors.New("HOTPEPPER_API_KEY environment variable is required")
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	rateLimitRequests := config.ParsePositiveIntEnv("RATE_LIMIT_REQUESTS", 60)
	rateLimitWindowSeconds := config.ParsePositiveIntEnv("RATE_LIMIT_WINDOW_SECONDS", 60)
	// RATE_LIMIT_BURST 未設定時は RATE_LIMIT_REQUESTS と同値（バースト=通常レートで現状と同じ挙動）
	rateLimitBurst := config.ParsePositiveIntEnv("RATE_LIMIT_BURST", rateLimitRequests)

	return appConfig{
		Port:                   port,
		APIKey:                 apiKey,
		AllowedOrigin:          allowedOrigin,
		RateLimitRequests:      rateLimitRequests,
		RateLimitWindowSeconds: rateLimitWindowSeconds,
		RateLimitBurst:         rateLimitBurst,
	}, nil
}

// buildServices はアプリケーションで利用するサービス層を構築する。
func buildServices(cfg appConfig) appServices {
	return appServices{
		hotpepper: service.NewHotpepperService(cfg.APIKey),
	}
}

// newLimiterFactory はレートリミット設定済みの LimiterStore 生成関数を返す。
func newLimiterFactory(cfg appConfig) func() *middleware.LimiterStore {
	// エンドポイントごとに独立したストアを持つことで、あるエンドポイントへの
	// 過剰アクセスが他のエンドポイントのレートリミットに影響しないようにする
	return func() *middleware.LimiterStore {
		return middleware.NewLimiterStore(
			cfg.RateLimitRequests,
			time.Duration(cfg.RateLimitWindowSeconds)*time.Second,
			middleware.WithBurst(cfg.RateLimitBurst),
		)
	}
}

// buildRouter はハンドラ生成と API エンドポイントのルーティング登録を行う。
func buildRouter(cfg appConfig, svcs appServices, newLimiterStore func() *middleware.LimiterStore) *http.ServeMux {
	gourmetHandler := handler.NewGourmetHandler(svcs.hotpepper)
	gourmetDetailHandler := handler.NewGourmetDetailHandler(svcs.hotpepper)
	genreHandler := handler.NewGenreHandler(svcs.hotpepper)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/gourmet", withMiddlewares(cfg.AllowedOrigin, newLimiterStore, gourmetHandler.Handle))
	mux.HandleFunc("/api/gourmet/detail", withMiddlewares(cfg.AllowedOrigin, newLimiterStore, gourmetDetailHandler.Handle))
	mux.HandleFunc("/api/genre", withMiddlewares(cfg.AllowedOrigin, newLimiterStore, genreHandler.Handle))

	return mux
}

// withMiddlewares は CORS と RateLimit を合成してハンドラへ適用する。
func withMiddlewares(allowedOrigin string, newLimiterStore func() *middleware.LimiterStore, next http.HandlerFunc) http.HandlerFunc {
	return middleware.CORS(allowedOrigin, middleware.RateLimit(newLimiterStore(), next))
}

// startServer は起動ログを出力し、HTTP サーバーを開始する。
func startServer(port string, h http.Handler) error {
	log.Printf("Server starting on port %s", port)
	return http.ListenAndServe(":"+port, h)
}
