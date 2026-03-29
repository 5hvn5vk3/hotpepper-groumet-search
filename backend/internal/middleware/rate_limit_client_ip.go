package middleware

import (
	"net"
	"net/http"
	"strings"
)

// このファイルはリバースプロキシ環境を考慮したクライアント IP 解決のみを含む。
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

	// X-Real-IP は XFF を付与しない一部のプロキシ（Nginx 等）向けのフォールバック。
	if xri := strings.TrimSpace(r.Header.Get("X-Real-IP")); xri != "" {
		if net.ParseIP(xri) != nil {
			return xri
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
