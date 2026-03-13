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

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

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

func buildServices(cfg appConfig) appServices {
	return appServices{
		hotpepper: service.NewHotpepperService(cfg.APIKey),
	}
}

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

func withMiddlewares(allowedOrigin string, newLimiterStore func() *middleware.LimiterStore, next http.HandlerFunc) http.HandlerFunc {
	return middleware.CORS(allowedOrigin, middleware.RateLimit(newLimiterStore(), next))
}

func startServer(port string, h http.Handler) error {
	log.Printf("Server starting on port %s", port)
	return http.ListenAndServe(":"+port, h)
}
