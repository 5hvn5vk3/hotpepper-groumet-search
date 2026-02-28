package main

import (
	"log"
	"net/http"
	"os"

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

	gourmetHandler := handler.NewGourmetHandler(hotpepperService)
	gourmetDetailHandler := handler.NewGourmetDetailHandler(hotpepperService)

	genreHandler := handler.NewGenreHandler(hotpepperService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/gourmet", middleware.CORS(allowedOrigin, gourmetHandler.Handle))
	mux.HandleFunc("/api/gourmet/detail", middleware.CORS(allowedOrigin, gourmetDetailHandler.Handle))

	mux.HandleFunc("/api/genre", middleware.CORS(allowedOrigin, genreHandler.Handle))

	log.Printf("Server starting on port %s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
