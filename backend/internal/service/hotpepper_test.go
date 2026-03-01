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

// newTestService はテスト用の HotpepperService を生成するヘルパー。
// 実 HTTP サーバーを立てず、RoundTripper でレスポンスを差し替える。
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
		svc := newTestService(roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"broken":`)), // 閉じていない不正 JSON
				Header:     make(http.Header),
			}, nil
		}), 5*time.Second)
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "test"})

		if err == nil {
			t.Fatal("エラーが返るべきなのに nil だった")
		}
		if !strings.Contains(err.Error(), "failed to decode response") {
			t.Errorf("期待するエラー文字列 'failed to decode response' が含まれない: %v", err)
		}
	})

	t.Run("タイムアウト_failedToFetchAPI", func(t *testing.T) {
		svc := newTestService(roundTripFunc(func(req *http.Request) (*http.Response, error) {
			select {
			case <-time.After(500 * time.Millisecond): // クライアントのタイムアウトより長く待つ
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"results":{}}`)),
					Header:     make(http.Header),
				}, nil
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
		}), 50*time.Millisecond) // 50ms で即タイムアウト
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "test"})

		if err == nil {
			t.Fatal("エラーが返るべきなのに nil だった")
		}
		if !strings.Contains(err.Error(), "failed to fetch API") {
			t.Errorf("期待するエラー文字列 'failed to fetch API' が含まれない: %v", err)
		}
	})

	t.Run("通信失敗_failedToFetchAPI", func(t *testing.T) {
		svc := newTestService(errorTransport{err: errors.New("connection refused")}, 5*time.Second)
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
	svc := newTestService(roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(
				`{"results":{"error":[{"code":3000,"message":"パラメータ不正"}]}}`,
			)),
			Header: make(http.Header),
		}, nil
	}), 5*time.Second)
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

// TestSearchGourmet_ClampStartAndCount は clampInt による start/count の補正を検証する。
// start=0 → 1、count=200 → 100 に補正されることを、実際に送信される HTTP クエリパラメータで確認する。
func TestSearchGourmet_ClampStartAndCount(t *testing.T) {
	testCases := []struct {
		name      string
		start     int
		count     int
		wantStart string
		wantCount string
	}{
		{name: "start=0 は 1 にクランプ", start: 0, count: 10, wantStart: "1", wantCount: "10"},
		{name: "count=200 は 100 にクランプ", start: 1, count: 200, wantStart: "1", wantCount: "100"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotStart, gotCount string
			svc := newTestService(roundTripFunc(func(req *http.Request) (*http.Response, error) {
				gotStart = req.URL.Query().Get("start")
				gotCount = req.URL.Query().Get("count")
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"results":{"shop":[]}}`)),
					Header:     make(http.Header),
				}, nil
			}), 5*time.Second)

			_, _ = svc.SearchGourmet(types.GourmetSearchParams{
				Keyword: "test",
				Start:   tc.start,
				Count:   tc.count,
			})

			if gotStart != tc.wantStart {
				t.Errorf("start = %q, want %q", gotStart, tc.wantStart)
			}
			if gotCount != tc.wantCount {
				t.Errorf("count = %q, want %q", gotCount, tc.wantCount)
			}
		})
	}
}
