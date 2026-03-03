package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestClientKeyFromRequest_XRealIPFallback は XFF がない場合に
// X-Real-IP ヘッダからクライアント IP を取得することを確認する。
func TestClientKeyFromRequest_XRealIPFallback(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "10.0.0.1:9999"
	req.Header.Set("X-Real-IP", "203.0.113.55")

	got := clientKeyFromRequest(req)
	want := "203.0.113.55"
	if got != want {
		t.Fatalf("clientKeyFromRequest = %q, want %q", got, want)
	}
}

// TestClientKeyFromRequest_XRealIPInvalidFallsThrough は X-Real-IP が不正な値の場合に
// RemoteAddr にフォールバックすることを確認する。
func TestClientKeyFromRequest_XRealIPInvalidFallsThrough(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "192.0.2.99:1234"
	req.Header.Set("X-Real-IP", "not-an-ip")

	got := clientKeyFromRequest(req)
	want := "192.0.2.99"
	if got != want {
		t.Fatalf("clientKeyFromRequest = %q, want %q", got, want)
	}
}

// TestClientKeyFromRequest_XFFTakesPrecedenceOverXRealIP は XFF が存在する場合に
// X-Real-IP より XFF が優先されることを確認する。
func TestClientKeyFromRequest_XFFTakesPrecedenceOverXRealIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "10.0.0.1:9999"
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	req.Header.Set("X-Real-IP", "203.0.113.99")

	got := clientKeyFromRequest(req)
	want := "203.0.113.10" // XFF が優先
	if got != want {
		t.Fatalf("clientKeyFromRequest = %q, want %q", got, want)
	}
}
