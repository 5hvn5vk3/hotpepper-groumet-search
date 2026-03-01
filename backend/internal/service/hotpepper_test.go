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
			time.Sleep(500 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"results":{"shop":[]}}`))
		}))
		defer server.Close()

		svc := newTestService(server.URL, 50*time.Millisecond)
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

func TestSearchGourmet_ClampStartAndCount(t *testing.T) {
	testCases := []struct {
		name       string
		start      int
		count      int
		wantStart  string
		wantCount  string
	}{
		{
			name:      "start below min is clamped to 1",
			start:     0,
			count:     10,
			wantStart: "1",
			wantCount: "10",
		},
		{
			name:      "count above max is clamped to 100",
			start:     1,
			count:     200,
			wantStart: "1",
			wantCount: "100",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotStart, gotCount string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotStart = r.URL.Query().Get("start")
				gotCount = r.URL.Query().Get("count")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"results":{"shop":[]}}`))
			}))
			defer server.Close()

			svc := newTestService(server.URL, time.Second)
			_, _ = svc.SearchGourmet(types.GourmetSearchParams{
				Keyword: "sushi",
				Start:   tc.start,
				Count:   tc.count,
			})

			if gotStart != tc.wantStart {
				t.Fatalf("start = %q, want %q", gotStart, tc.wantStart)
			}
			if gotCount != tc.wantCount {
				t.Fatalf("count = %q, want %q", gotCount, tc.wantCount)
			}
		})
	}
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
