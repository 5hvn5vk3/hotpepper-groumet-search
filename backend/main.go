package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/service"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("HOTPEPPER_API_KEY")

	if apiKey == "" {
		log.Fatal("HOTPEPPER_API_KEY environment variable is required")
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}

	hotpepperService := service.NewHotpepperService(apiKey)

	rateLimitRequests := config.ParsePositiveIntEnv("RATE_LIMIT_REQUESTS", 60)
	rateLimitWindowSeconds := config.ParsePositiveIntEnv("RATE_LIMIT_WINDOW_SECONDS", 60)
	// RATE_LIMIT_BURST 未設定時は RATE_LIMIT_REQUESTS と同値（バースト=通常レートで現状と同じ挙動）
	rateLimitBurst := config.ParsePositiveIntEnv("RATE_LIMIT_BURST", rateLimitRequests)
	// エンドポイントごとに独立したストアを持つことで、あるエンドポイントへの
	// 過剰アクセスが他のエンドポイントのレートリミットに影響しないようにする
	newLimiterStore := func() *middleware.LimiterStore {
		return middleware.NewLimiterStore(
			rateLimitRequests,
			time.Duration(rateLimitWindowSeconds)*time.Second,
			middleware.WithBurst(rateLimitBurst),
		)
	}

	gourmetHandler := handler.NewGourmetHandler(hotpepperService)
	gourmetDetailHandler := handler.NewGourmetDetailHandler(hotpepperService)

	genreHandler := handler.NewGenreHandler(hotpepperService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/gourmet", middleware.CORS(allowedOrigin, middleware.RateLimit(newLimiterStore(), gourmetHandler.Handle)))
	mux.HandleFunc("/api/gourmet/detail", middleware.CORS(allowedOrigin, middleware.RateLimit(newLimiterStore(), gourmetDetailHandler.Handle)))

	mux.HandleFunc("/api/genre", middleware.CORS(allowedOrigin, middleware.RateLimit(newLimiterStore(), genreHandler.Handle)))

	log.Printf("Server starting on port %s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
