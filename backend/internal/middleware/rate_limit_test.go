package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type fixedClock struct {
	current time.Time
}

func newFixedClock(start time.Time) *fixedClock {
	return &fixedClock{current: start}
}

func (c *fixedClock) Now() time.Time {
	return c.current
}

func (c *fixedClock) Advance(d time.Duration) {
	c.current = c.current.Add(d)
}

func TestRateLimit_AllowsWithinLimit(t *testing.T) {
	clock := newFixedClock(time.Unix(0, 0))
	store := NewLimiterStore(2, time.Second, WithNowFunc(clock.Now))

	nextCalls := 0
	handler := RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		nextCalls++
		w.WriteHeader(http.StatusOK)
	})

	for i := 0; i < 2; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.RemoteAddr = "192.0.2.10:12345"

		handler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	}

	if nextCalls != 2 {
		t.Fatalf("next calls = %d, want %d", nextCalls, 2)
	}
}

func TestRateLimit_BlocksOnLimitExceeded(t *testing.T) {
	clock := newFixedClock(time.Unix(0, 0))
	store := NewLimiterStore(1, time.Hour, WithNowFunc(clock.Now))

	nextCalls := 0
	handler := RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		nextCalls++
		w.WriteHeader(http.StatusOK)
	})

	firstRec := httptest.NewRecorder()
	firstReq := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	firstReq.RemoteAddr = "192.0.2.20:11111"
	handler(firstRec, firstReq)

	if firstRec.Code != http.StatusOK {
		t.Fatalf("first status = %d, want %d", firstRec.Code, http.StatusOK)
	}

	secondRec := httptest.NewRecorder()
	secondReq := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	secondReq.RemoteAddr = "192.0.2.20:22222"
	handler(secondRec, secondReq)

	if secondRec.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want %d", secondRec.Code, http.StatusTooManyRequests)
	}

	if nextCalls != 1 {
		t.Fatalf("next calls = %d, want %d", nextCalls, 1)
	}

	contentType := secondRec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("content-type = %q, want contains application/json", contentType)
	}

	var body map[string]map[string]string
	if err := json.NewDecoder(secondRec.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	got := body["error"]["message"]
	if got != rateLimitMessage {
		t.Fatalf("error.message = %q, want %q", got, rateLimitMessage)
	}
}

func TestRateLimit_OptionsBypassesAndDoesNotConsumeTokens(t *testing.T) {
	clock := newFixedClock(time.Unix(0, 0))
	store := NewLimiterStore(1, time.Hour, WithNowFunc(clock.Now))

	nextCalls := 0
	handler := RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		nextCalls++
		w.WriteHeader(http.StatusOK)
	})

	optionsRec := httptest.NewRecorder()
	optionsReq := httptest.NewRequest(http.MethodOptions, "/api/test", nil)
	optionsReq.RemoteAddr = "192.0.2.30:33333"
	handler(optionsRec, optionsReq)

	if optionsRec.Code != http.StatusOK {
		t.Fatalf("OPTIONS status = %d, want %d", optionsRec.Code, http.StatusOK)
	}

	firstGetRec := httptest.NewRecorder()
	firstGetReq := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	firstGetReq.RemoteAddr = "192.0.2.30:33333"
	handler(firstGetRec, firstGetReq)

	if firstGetRec.Code != http.StatusOK {
		t.Fatalf("first GET status = %d, want %d", firstGetRec.Code, http.StatusOK)
	}

	secondGetRec := httptest.NewRecorder()
	secondGetReq := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	secondGetReq.RemoteAddr = "192.0.2.30:44444"
	handler(secondGetRec, secondGetReq)

	if secondGetRec.Code != http.StatusTooManyRequests {
		t.Fatalf("second GET status = %d, want %d", secondGetRec.Code, http.StatusTooManyRequests)
	}

	if nextCalls != 2 {
		t.Fatalf("next calls = %d, want %d", nextCalls, 2)
	}
}

func TestRateLimit_DifferentIPsAreIndependent(t *testing.T) {
	clock := newFixedClock(time.Unix(0, 0))
	store := NewLimiterStore(1, time.Hour, WithNowFunc(clock.Now))

	handler := RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	reqA1 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	reqA1.RemoteAddr = "192.0.2.40:10001"
	recA1 := httptest.NewRecorder()
	handler(recA1, reqA1)

	reqA2 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	reqA2.RemoteAddr = "192.0.2.40:10002"
	recA2 := httptest.NewRecorder()
	handler(recA2, reqA2)

	reqB1 := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	reqB1.RemoteAddr = "198.51.100.10:20001"
	recB1 := httptest.NewRecorder()
	handler(recB1, reqB1)

	if recA1.Code != http.StatusOK {
		t.Fatalf("ipA first status = %d, want %d", recA1.Code, http.StatusOK)
	}
	if recA2.Code != http.StatusTooManyRequests {
		t.Fatalf("ipA second status = %d, want %d", recA2.Code, http.StatusTooManyRequests)
	}
	if recB1.Code != http.StatusOK {
		t.Fatalf("ipB first status = %d, want %d", recB1.Code, http.StatusOK)
	}
}

func TestLimiterStore_RecoversAfterWindow(t *testing.T) {
	clock := newFixedClock(time.Unix(0, 0))
	store := NewLimiterStore(1, time.Second, WithNowFunc(clock.Now))

	if !store.Allow("192.0.2.50") {
		t.Fatal("first allow should be true")
	}
	if store.Allow("192.0.2.50") {
		t.Fatal("second allow should be false")
	}

	clock.Advance(time.Second)

	if !store.Allow("192.0.2.50") {
		t.Fatal("allow after window should be true")
	}
}

func TestLimiterStore_CleansUpStaleEntries(t *testing.T) {
	clock := newFixedClock(time.Unix(0, 0))
	store := NewLimiterStore(
		1,
		time.Second,
		WithNowFunc(clock.Now),
		WithCleanup(2*time.Second, time.Second),
	)

	if !store.Allow("192.0.2.60") {
		t.Fatal("first allow should be true")
	}

	clock.Advance(3 * time.Second)

	if !store.Allow("192.0.2.61") {
		t.Fatal("allow for second client should be true")
	}

	// クリーンアップ後は 192.0.2.60 のエントリが削除され、残存クライアントは 1 件のみ
	if got := store.ClientCount(); got != 1 {
		t.Fatalf("client count after cleanup = %d, want 1", got)
	}
}

func TestRateLimit_ConcurrentRequestsAreBounded(t *testing.T) {
	store := NewLimiterStore(20, time.Hour)

	handler := RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var okCount int32
	var tooManyCount int32

	const total = 100

	var wg sync.WaitGroup
	wg.Add(total)

	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
			req.RemoteAddr = "203.0.113.1:5555"

			handler(rec, req)

			switch rec.Code {
			case http.StatusOK:
				atomic.AddInt32(&okCount, 1)
			case http.StatusTooManyRequests:
				atomic.AddInt32(&tooManyCount, 1)
			default:
				t.Errorf("unexpected status: %d", rec.Code)
			}
		}()
	}

	wg.Wait()

	if okCount > 20 {
		t.Fatalf("okCount = %d, want <= 20", okCount)
	}

	if okCount+tooManyCount != total {
		t.Fatalf("handled=%d, want=%d", okCount+tooManyCount, total)
	}
}

func TestWithBurst_LimitsToBurstNotRequests(t *testing.T) {
	clock := newFixedClock(time.Unix(0, 0))
	// requests=10 だがバーストを 2 に絞る
	store := NewLimiterStore(10, time.Hour, WithNowFunc(clock.Now), WithBurst(2))

	handler := RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	makeReq := func() int {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.RemoteAddr = "203.0.113.100:1234"
		handler(rec, req)
		return rec.Code
	}

	if got := makeReq(); got != http.StatusOK {
		t.Fatalf("1st status = %d, want %d", got, http.StatusOK)
	}
	if got := makeReq(); got != http.StatusOK {
		t.Fatalf("2nd status = %d, want %d", got, http.StatusOK)
	}
	if got := makeReq(); got != http.StatusTooManyRequests {
		t.Fatalf("3rd status = %d, want %d (burst=2 should block)", got, http.StatusTooManyRequests)
	}
}

func TestRateLimit_Returns429WithRetryAfterHeader(t *testing.T) {
	clock := newFixedClock(time.Unix(0, 0))
	// 10 req / 60 sec → 1 トークン補充に 6 秒 → Retry-After: 6
	store := NewLimiterStore(10, 60*time.Second, WithNowFunc(clock.Now))

	handler := RateLimit(store, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	makeReq := func() *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.RemoteAddr = "203.0.113.110:9999"
		handler(rec, req)
		return rec
	}

	// 1 回目: 通過
	makeReq()

	// 2 回目以降: ブロック
	for i := 0; i < 10; i++ {
		makeReq()
	}

	rec := makeReq()
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}

	got := rec.Header().Get("Retry-After")
	if got != "6" {
		t.Fatalf("Retry-After = %q, want %q", got, "6")
	}
}

func TestClientKeyFromRequest_XForwardedFor(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		xff        string
		wantKey    string
	}{
		{
			name:       "XFF単一IP: プロキシが付与したIPを返す",
			remoteAddr: "10.0.0.1:9999",
			xff:        "203.0.113.10",
			wantKey:    "203.0.113.10",
		},
		{
			name:       "XFF複数IP: 末尾（プロキシ付与）のIPを返す",
			remoteAddr: "10.0.0.1:9999",
			xff:        "1.2.3.4, 5.6.7.8, 203.0.113.20",
			wantKey:    "203.0.113.20",
		},
		{
			name:       "XFF末尾が不正な場合: 手前の有効IPを返す",
			remoteAddr: "10.0.0.1:9999",
			xff:        "203.0.113.30, invalid",
			wantKey:    "203.0.113.30",
		},
		{
			name:       "XFFなし: RemoteAddrのホスト部を返す",
			remoteAddr: "192.0.2.50:12345",
			xff:        "",
			wantKey:    "192.0.2.50",
		},
		{
			name:       "XFF全エントリ不正: RemoteAddrのホスト部を返す",
			remoteAddr: "192.0.2.60:12345",
			xff:        "invalid, garbage",
			wantKey:    "192.0.2.60",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xff != "" {
				req.Header.Set("X-Forwarded-For", tt.xff)
			}

			got := clientKeyFromRequest(req)
			if got != tt.wantKey {
				t.Fatalf("clientKeyFromRequest = %q, want %q", got, tt.wantKey)
			}
		})
	}
}
