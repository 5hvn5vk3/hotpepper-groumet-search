package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// このファイルは HTTP ミドルウェアのみを含む。
// トークンバケットのドメインロジックとストア管理は rate_limit_store.go、
// クライアント IP 解決は rate_limit_client_ip.go を参照。
const rateLimitMessage = "リクエストが多すぎます。しばらく待ってから再試行してください。"

// RateLimit は GET リクエストに対して store のレートリミットを適用するミドルウェアを返す。
// store が nil の場合はレートリミットを無効化する（テスト・開発環境での一時的な無効化に利用できる）。
// 通過・拒否いずれの場合も X-RateLimit-Limit / Remaining / Reset ヘッダを付与する。
func RateLimit(store *LimiterStore, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next(w, r)
			return
		}

		if store == nil {
			next(w, r)
			return
		}

		result := store.allowWithInfo(clientKeyFromRequest(r))
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(result.limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.resetAt.Unix(), 10))

		if result.allowed {
			next(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", strconv.Itoa(store.retryAfterSec))
		w.WriteHeader(http.StatusTooManyRequests)
		body, _ := json.Marshal(map[string]any{
			"error": map[string]string{
				"message": rateLimitMessage,
			},
		})
		_, _ = w.Write(body)
	}
}
