package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

	rateLimitRequests := parsePositiveIntEnv("RATE_LIMIT_REQUESTS", 60)
	rateLimitWindowSeconds := parsePositiveIntEnv("RATE_LIMIT_WINDOW_SECONDS", 60)
	limiterStore := middleware.NewLimiterStore(
		rateLimitRequests,
		time.Duration(rateLimitWindowSeconds)*time.Second,
	)

	gourmetHandler := handler.NewGourmetHandler(hotpepperService)
	gourmetDetailHandler := handler.NewGourmetDetailHandler(hotpepperService)

	genreHandler := handler.NewGenreHandler(hotpepperService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/gourmet", middleware.CORS(allowedOrigin, middleware.RateLimit(limiterStore, gourmetHandler.Handle)))
	mux.HandleFunc("/api/gourmet/detail", middleware.CORS(allowedOrigin, middleware.RateLimit(limiterStore, gourmetDetailHandler.Handle)))

	mux.HandleFunc("/api/genre", middleware.CORS(allowedOrigin, middleware.RateLimit(limiterStore, genreHandler.Handle)))

	log.Printf("Server starting on port %s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func parsePositiveIntEnv(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		log.Printf("invalid %s=%q; using default %d", name, value, fallback)
		return fallback
	}

	return parsed
}
