package service

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"backend/internal/types"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestService(transport http.RoundTripper, timeout time.Duration) *HotpepperService {
	if transport == nil {
		transport = http.DefaultTransport
	}

	return &HotpepperService{
		apiKey:  "test-key",
		baseURL: "http://example.test",
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

type errorTransport struct{ err error }

func (t errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, t.err
}

type failingReadCloser struct{}

func (failingReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func (failingReadCloser) Close() error {
	return nil
}

func TestSearchGourmet_NetworkAndFormatErrors(t *testing.T) {
	t.Run("invalid JSON", func(t *testing.T) {
		svc := newTestService(roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"broken":`)),
				Header:     make(http.Header),
			}, nil
		}), time.Second)

		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to decode response") {
			t.Fatalf("expected decode error, got: %v", err)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		svc := newTestService(roundTripFunc(func(req *http.Request) (*http.Response, error) {
			select {
			case <-time.After(500 * time.Millisecond):
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"results":{"shop":[]}}`)),
					Header:     make(http.Header),
				}, nil
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
		}), 50*time.Millisecond)

		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to fetch API") {
			t.Fatalf("expected fetch error, got: %v", err)
		}
	})

	t.Run("transport failure", func(t *testing.T) {
		svc := newTestService(errorTransport{err: errors.New("dial tcp: no route to host")}, time.Second)

		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to fetch API") {
			t.Fatalf("expected fetch error, got: %v", err)
		}
	})

	t.Run("non-200 status", func(t *testing.T) {
		svc := newTestService(roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Body:       io.NopCloser(strings.NewReader(`upstream failure`)),
				Header:     make(http.Header),
			}, nil
		}), time.Second)

		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "API returned status code 502") {
			t.Fatalf("expected status code error, got: %v", err)
		}
	})

	t.Run("failed to read response body", func(t *testing.T) {
		svc := newTestService(roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       failingReadCloser{},
				Header:     make(http.Header),
			}, nil
		}), time.Second)

		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to read response") {
			t.Fatalf("expected read error, got: %v", err)
		}
	})
}

func TestSearchGourmet_ClampStartAndCount(t *testing.T) {
	testCases := []struct {
		name      string
		start     int
		count     int
		wantStart string
		wantCount string
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
			svc := newTestService(roundTripFunc(func(r *http.Request) (*http.Response, error) {
				gotStart = r.URL.Query().Get("start")
				gotCount = r.URL.Query().Get("count")
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"results":{"shop":[]}}`)),
					Header:     make(http.Header),
				}, nil
			}), time.Second)

			_, err := svc.SearchGourmet(types.GourmetSearchParams{
				Keyword: "sushi",
				Start:   tc.start,
				Count:   tc.count,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

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
	svc := newTestService(roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(
				`{"results":{"error":[{"code":3000,"message":"パラメータ不正"}]}}`,
			)),
			Header: make(http.Header),
		}, nil
	}), time.Second)

	_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *HotpepperAPIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *HotpepperAPIError, got: %T (%v)", err, err)
	}
	if apiErr.Code != 3000 {
		t.Fatalf("expected code=3000, got: %d", apiErr.Code)
	}
}
