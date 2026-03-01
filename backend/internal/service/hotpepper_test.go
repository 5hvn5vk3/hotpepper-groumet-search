package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/internal/types"
)

func newTestService(serverURL string, timeout time.Duration) *HotpepperService {
	return &HotpepperService{
		apiKey:  "test-key",
		baseURL: serverURL,
		client:  &http.Client{Timeout: timeout},
	}
}

type errorTransport struct{ err error }

func (t errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, t.err
}

func TestSearchGourmet_NetworkAndFormatErrors(t *testing.T) {
	t.Run("invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"broken":`))
		}))
		defer server.Close()

		svc := newTestService(server.URL, time.Second)
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to decode response") {
			t.Fatalf("expected decode error, got: %v", err)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"results":{"shop":[]}}`))
		}))
		defer server.Close()

		svc := newTestService(server.URL, 5*time.Millisecond)
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to fetch API") {
			t.Fatalf("expected fetch error, got: %v", err)
		}
	})

	t.Run("transport failure", func(t *testing.T) {
		svc := newTestService("http://example.com", time.Second)
		svc.client = &http.Client{
			Timeout:   time.Second,
			Transport: errorTransport{err: errors.New("dial tcp: no route to host")},
		}

		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to fetch API") {
			t.Fatalf("expected fetch error, got: %v", err)
		}
	})
}

func TestSearchGourmet_HotpepperAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"results":{"error":[{"code":3000,"message":"パラメータ不正"}]}}`))
	}))
	defer server.Close()

	svc := newTestService(server.URL, time.Second)
	_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *types.HotpepperAPIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *types.HotpepperAPIError, got: %T (%v)", err, err)
	}
	if apiErr.Code != 3000 {
		t.Fatalf("expected code=3000, got: %d", apiErr.Code)
	}
}
