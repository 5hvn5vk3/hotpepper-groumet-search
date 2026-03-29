package middleware

import (
	"math"
	"strings"
	"sync"
	"time"
)

// このファイルはトークンバケットのドメインロジックと LimiterStore の
// ライフサイクル管理のみを含む。HTTP には依存しない。

type clientBucket struct {
	tokens     float64
	lastRefill time.Time
	lastSeen   time.Time
}

// rateLimitResult は allowWithInfo が返す判定結果と付随情報をまとめた型。
type rateLimitResult struct {
	allowed   bool
	remaining int       // 消費後の残りトークン数（floor）
	limit     int       // バースト上限
	resetAt   time.Time // 次に 1 トークン分の余裕ができる推定時刻
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

	// 1 トークン補充に要する最小秒数を Retry-After のデフォルト値として事前計算する
	store.retryAfterSec = int(math.Ceil(1.0 / store.ratePerSecond))
	if store.retryAfterSec < 1 {
		store.retryAfterSec = 1
	}

	return store
}

// Allow は clientKey のリクエストを許可するかどうかを返す。
func (s *LimiterStore) Allow(clientKey string) bool {
	return s.allowWithInfo(clientKey).allowed
}

// allowWithInfo はレートリミット判定を行い、ヘッダ付与に必要な情報を返す。
func (s *LimiterStore) allowWithInfo(clientKey string) rateLimitResult {
	if s == nil {
		return rateLimitResult{allowed: true}
	}

	key := strings.TrimSpace(clientKey)
	if key == "" {
		key = "unknown"
	}

	now := s.now()
	limit := int(s.burst)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpiredLocked(now)

	entry, exists := s.clients[key]
	if !exists {
		tokens := s.burst - 1
		s.clients[key] = &clientBucket{
			tokens:     tokens,
			lastRefill: now,
			lastSeen:   now,
		}
		return rateLimitResult{
			allowed:   true,
			remaining: int(tokens),
			limit:     limit,
			resetAt:   resetTime(now, tokens, s.ratePerSecond),
		}
	}

	entry.tokens = refillTokens(entry.tokens, now.Sub(entry.lastRefill), s.ratePerSecond, s.burst)
	entry.lastRefill = now
	entry.lastSeen = now

	if entry.tokens < 1 {
		return rateLimitResult{
			allowed:   false,
			remaining: 0,
			limit:     limit,
			resetAt:   resetTime(now, entry.tokens, s.ratePerSecond),
		}
	}

	entry.tokens--
	return rateLimitResult{
		allowed:   true,
		remaining: int(entry.tokens),
		limit:     limit,
		resetAt:   resetTime(now, entry.tokens, s.ratePerSecond),
	}
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

// resetTime は tokens 個のトークンがある状態で、次に 1 リクエストが通過できるようになる時刻を返す。
// tokens >= 1 の場合は今すぐ通過できるため now をそのまま返す。
func resetTime(now time.Time, tokens float64, ratePerSecond float64) time.Time {
	if tokens >= 1 {
		return now
	}
	// tokens < 1: あと (1 - tokens) トークン必要。補充に要する秒数を切り上げで計算する。
	waitSec := int64(math.Ceil((1.0 - tokens) / ratePerSecond))
	return now.Add(time.Duration(waitSec) * time.Second)
}

// ClientCount は現在追跡中のクライアント数を返す。主にテストで使用する。
func (s *LimiterStore) ClientCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.clients)
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
