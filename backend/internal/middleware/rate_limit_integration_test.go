package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimitIntegration_CorsHeaderRemainsOn429(t *testing.T) {
	allowedOrigin := "http://localhost:3000"
	store := NewLimiterStore(1, time.Hour)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/gourmet", CORS(allowedOrigin, RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	firstReq := httptest.NewRequest(http.MethodGet, "/api/gourmet", nil)
	firstReq.RemoteAddr = "203.0.113.200:41001"
	firstRec := httptest.NewRecorder()
	mux.ServeHTTP(firstRec, firstReq)

	if firstRec.Code != http.StatusOK {
		t.Fatalf("first status = %d, want %d", firstRec.Code, http.StatusOK)
	}

	secondReq := httptest.NewRequest(http.MethodGet, "/api/gourmet", nil)
	secondReq.RemoteAddr = "203.0.113.200:41002"
	secondRec := httptest.NewRecorder()
	mux.ServeHTTP(secondRec, secondReq)

	if secondRec.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want %d", secondRec.Code, http.StatusTooManyRequests)
	}

	if got := secondRec.Header().Get("Access-Control-Allow-Origin"); got != allowedOrigin {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, allowedOrigin)
	}

	var body map[string]map[string]string
	if err := json.NewDecoder(secondRec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if got := body["error"]["message"]; got != rateLimitMessage {
		t.Fatalf("error.message = %q, want %q", got, rateLimitMessage)
	}
}
