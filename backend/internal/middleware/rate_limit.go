package middleware

import (
	"encoding/json"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const rateLimitMessage = "リクエストが多すぎます。しばらく待ってから再試行してください。"

type clientBucket struct {
	tokens     float64
	lastRefill time.Time
	lastSeen   time.Time
}

type LimiterStore struct {
	mu sync.Mutex

	clients map[string]*clientBucket

	now             func() time.Time
	ratePerSecond   float64
	burst           float64
	retryAfterSec   int
	entryTTL        time.Duration
	cleanupInterval time.Duration
	lastCleanup     time.Time
}

type LimiterStoreOption func(*LimiterStore)

func WithNowFunc(now func() time.Time) LimiterStoreOption {
	return func(s *LimiterStore) {
		if now != nil {
			s.now = now
		}
	}
}

func WithCleanup(entryTTL, cleanupInterval time.Duration) LimiterStoreOption {
	return func(s *LimiterStore) {
		if entryTTL > 0 {
			s.entryTTL = entryTTL
		}
		if cleanupInterval > 0 {
			s.cleanupInterval = cleanupInterval
		}
	}
}

// WithBurst はバースト上限（トークンバケットの最大容量）を設定する。
// 未指定時は requests と同じ値が使われる。
func WithBurst(burst int) LimiterStoreOption {
	return func(s *LimiterStore) {
		if burst > 0 {
			s.burst = float64(burst)
		}
	}
}

func NewLimiterStore(requests int, window time.Duration, opts ...LimiterStoreOption) *LimiterStore {
	if requests <= 0 {
		requests = 60
	}
	if window <= 0 {
		window = 60 * time.Second
	}

	store := &LimiterStore{
		clients:         map[string]*clientBucket{},
		now:             time.Now,
		ratePerSecond:   float64(requests) / window.Seconds(),
		burst:           float64(requests),
		entryTTL:        10 * window,
		cleanupInterval: window,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(store)
		}
	}

	if math.IsNaN(store.ratePerSecond) || math.IsInf(store.ratePerSecond, 0) || store.ratePerSecond <= 0 {
		store.ratePerSecond = 1
	}
	if store.burst < 1 {
		store.burst = 1
	}
	if store.entryTTL <= 0 {
		store.entryTTL = 10 * window
	}
	if store.cleanupInterval <= 0 {
		store.cleanupInterval = window
	}

	// 1 トークン補充に要する最小秒数を Retry-After のデフォルト値として事前計算する
	store.retryAfterSec = int(math.Ceil(1.0 / store.ratePerSecond))
	if store.retryAfterSec < 1 {
		store.retryAfterSec = 1
	}

	return store
}

func RateLimit(store *LimiterStore, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next(w, r)
			return
		}

		if store == nil || store.Allow(clientKeyFromRequest(r)) {
			next(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", strconv.Itoa(store.retryAfterSec))
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{
				"message": rateLimitMessage,
			},
		})
	}
}

func (s *LimiterStore) Allow(clientKey string) bool {
	if s == nil {
		return true
	}

	key := strings.TrimSpace(clientKey)
	if key == "" {
		key = "unknown"
	}

	now := s.now()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpiredLocked(now)

	entry, exists := s.clients[key]
	if !exists {
		s.clients[key] = &clientBucket{
			tokens:     s.burst - 1,
			lastRefill: now,
			lastSeen:   now,
		}
		return true
	}

	entry.tokens = refillTokens(entry.tokens, now.Sub(entry.lastRefill), s.ratePerSecond, s.burst)
	entry.lastRefill = now
	entry.lastSeen = now

	if entry.tokens < 1 {
		return false
	}

	entry.tokens--
	return true
}

func refillTokens(tokens float64, elapsed time.Duration, ratePerSecond, burst float64) float64 {
	if elapsed > 0 {
		tokens += elapsed.Seconds() * ratePerSecond
		if tokens > burst {
			tokens = burst
		}
	}

	return tokens
}

func (s *LimiterStore) cleanupExpiredLocked(now time.Time) {
	if s.lastCleanup.IsZero() {
		s.lastCleanup = now
		return
	}

	if now.Sub(s.lastCleanup) < s.cleanupInterval {
		return
	}

	for key, entry := range s.clients {
		if now.Sub(entry.lastSeen) > s.entryTTL {
			delete(s.clients, key)
		}
	}

	s.lastCleanup = now
}

func clientKeyFromRequest(r *http.Request) string {
	if r == nil {
		return "unknown"
	}

	// Render 等のリバースプロキシ環境では r.RemoteAddr がプロキシ内部 IP になる。
	// X-Forwarded-For の末尾エントリがプロキシ層の付与した本物のクライアント IP のため優先して使用する。
	// 先頭エントリはクライアントが偽装できるため使用しない。
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		for i := len(parts) - 1; i >= 0; i-- {
			if ip := strings.TrimSpace(parts[i]); net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	remoteAddr := strings.TrimSpace(r.RemoteAddr)
	if remoteAddr == "" {
		return "unknown"
	}

	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil && strings.TrimSpace(host) != "" {
		return host
	}

	return remoteAddr
}
