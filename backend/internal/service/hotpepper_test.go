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

// newTestService はテスト用の HotpepperService を生成するヘルパー。
// 本物の API サーバーの代わりに偽サーバーの URL を注入できる。
func newTestService(serverURL string, timeout time.Duration) *HotpepperService {
	return &HotpepperService{
		apiKey:  "test-key",
		baseURL: serverURL,
		client:  &http.Client{Timeout: timeout},
	}
}

// errorTransport は通信を行わず即座にエラーを返す偽 Transport。
// ネットワーク障害のシミュレーションに使う。
type errorTransport struct{ err error }

func (t errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, t.err
}

// TestSearchGourmet_NetworkAndFormat は 1-a のテスト群。
// 不正 JSON・タイムアウト・通信失敗の 3 パターンを検証する。
func TestSearchGourmet_NetworkAndFormat(t *testing.T) {
	t.Run("不正JSON_failedToDecodeResponse", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"broken":`)) // 閉じていない不正 JSON
		}))
		defer ts.Close()

		svc := newTestService(ts.URL, 5*time.Second)
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "test"})

		if err == nil {
			t.Fatal("エラーが返るべきなのに nil だった")
		}
		if !strings.Contains(err.Error(), "failed to decode response") {
			t.Errorf("期待するエラー文字列 'failed to decode response' が含まれない: %v", err)
		}
	})

	t.Run("タイムアウト_failedToFetchAPI", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(200 * time.Millisecond) // クライアントのタイムアウトより長く待つ
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"results":{}}`))
		}))
		defer ts.Close()

		svc := newTestService(ts.URL, 5*time.Millisecond) // 5ms で即タイムアウト
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "test"})

		if err == nil {
			t.Fatal("エラーが返るべきなのに nil だった")
		}
		if !strings.Contains(err.Error(), "failed to fetch API") {
			t.Errorf("期待するエラー文字列 'failed to fetch API' が含まれない: %v", err)
		}
	})

	t.Run("通信失敗_failedToFetchAPI", func(t *testing.T) {
		svc := &HotpepperService{
			apiKey:  "test-key",
			baseURL: "http://dummy",
			client: &http.Client{
				Transport: errorTransport{err: errors.New("connection refused")},
			},
		}
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "test"})

		if err == nil {
			t.Fatal("エラーが返るべきなのに nil だった")
		}
		if !strings.Contains(err.Error(), "failed to fetch API") {
			t.Errorf("期待するエラー文字列 'failed to fetch API' が含まれない: %v", err)
		}
	})
}

// TestSearchGourmet_HotpepperAPIError は 1-b のテスト。
// API が HTTP 200 で error フィールドを返した場合に *types.HotpepperAPIError が返ることを確認する。
func TestSearchGourmet_HotpepperAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results":{"error":[{"code":3000,"message":"パラメータ不正"}]}}`))
	}))
	defer ts.Close()

	svc := newTestService(ts.URL, 5*time.Second)
	_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "test"})

	if err == nil {
		t.Fatal("エラーが返るべきなのに nil だった")
	}

	var apiErr *types.HotpepperAPIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("*types.HotpepperAPIError が返るべきなのに別の型だった: %T (%v)", err, err)
	}
	if apiErr.Code != 3000 {
		t.Errorf("Code=3000 を期待したが Code=%d だった", apiErr.Code)
	}
}
