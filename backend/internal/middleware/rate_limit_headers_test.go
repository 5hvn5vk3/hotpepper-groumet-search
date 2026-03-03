package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// --- X-RateLimit-* ヘッダのテスト ---

// TestRateLimit_XRateLimitHeadersPresentOnOK は 200 レスポンス時に
// X-RateLimit-Limit / Remaining / Reset の 3 ヘッダが付与されることを確認する。
func TestRateLimit_XRateLimitHeadersPresentOnOK(t *testing.T) {
	clock := newFixedClock(time.Unix(1000, 0))
	store := NewLimiterStore(10, time.Minute, WithNowFunc(clock.Now))

	handler := RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "192.0.2.1:1234"
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	for _, h := range []string{"X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"} {
		if got := rec.Header().Get(h); got == "" {
			t.Fatalf("header %s is missing on 200", h)
		}
	}
}

// TestRateLimit_XRateLimitHeadersPresentOn429 は 429 レスポンス時に
// X-RateLimit-Limit / Remaining / Reset / Retry-After の 4 ヘッダが付与されることを確認する。
func TestRateLimit_XRateLimitHeadersPresentOn429(t *testing.T) {
	clock := newFixedClock(time.Unix(1000, 0))
	store := NewLimiterStore(1, time.Hour, WithNowFunc(clock.Now))

	handler := RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	makeReq := func() *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.RemoteAddr = "192.0.2.2:5678"
		handler(rec, req)
		return rec
	}

	makeReq() // バースト 1 を消費

	rec := makeReq() // 429
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}

	for _, h := range []string{"X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset", "Retry-After"} {
		if got := rec.Header().Get(h); got == "" {
			t.Fatalf("header %s is missing on 429", h)
		}
	}
}

// TestRateLimit_XRateLimitRemainingDecreases はリクエストごとに
// X-RateLimit-Remaining が 1 ずつ減少することを確認する。
func TestRateLimit_XRateLimitRemainingDecreases(t *testing.T) {
	clock := newFixedClock(time.Unix(1000, 0))
	store := NewLimiterStore(3, time.Hour, WithNowFunc(clock.Now))

	handler := RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	remaining := func() int {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.RemoteAddr = "192.0.2.3:9999"
		handler(rec, req)
		v, _ := strconv.Atoi(rec.Header().Get("X-RateLimit-Remaining"))
		return v
	}

	// burst=3: 1 回目消費後 tokens=2, 2 回目後 tokens=1, 3 回目後 tokens=0
	r1, r2, r3 := remaining(), remaining(), remaining()
	if r1 != 2 || r2 != 1 || r3 != 0 {
		t.Fatalf("remaining sequence = %d,%d,%d, want 2,1,0", r1, r2, r3)
	}
}

// TestRateLimit_XRateLimitResetOnBlocked は 429 時の X-RateLimit-Reset が
// 現在時刻より後の Unix 秒であることを確認する。
func TestRateLimit_XRateLimitResetOnBlocked(t *testing.T) {
	// 1 req / 60 sec → ratePerSecond = 1/60
	// tokens=0 でブロック時: resetAt = now + ceil(1/(1/60)) = now + 60s
	start := time.Unix(1000, 0)
	clock := newFixedClock(start)
	store := NewLimiterStore(1, 60*time.Second, WithNowFunc(clock.Now))

	handler := RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	makeReq := func() *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.RemoteAddr = "192.0.2.4:7777"
		handler(rec, req)
		return rec
	}

	makeReq() // バースト 1 を消費

	rec := makeReq() // 429
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", rec.Code)
	}

	got, err := strconv.ParseInt(rec.Header().Get("X-RateLimit-Reset"), 10, 64)
	if err != nil {
		t.Fatalf("X-RateLimit-Reset is not a valid integer: %v", err)
	}

	want := start.Add(60 * time.Second).Unix() // 1000 + 60 = 1060
	if got != want {
		t.Fatalf("X-RateLimit-Reset = %d, want %d", got, want)
	}
}
