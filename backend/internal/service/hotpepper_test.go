// このファイルが属するパッケージ名を宣言する。
package service

// ここから使用する外部パッケージの列挙を始める。
import (
	// この行は処理の一部として必要な操作を実行する。
	"errors"
	// この行は処理の一部として必要な操作を実行する。
	"io"
	// この行は処理の一部として必要な操作を実行する。
	"net/http"
	// この行は処理の一部として必要な操作を実行する。
	"strings"
	// この行は処理の一部として必要な操作を実行する。
	"testing"
	// この行は処理の一部として必要な操作を実行する。
	"time"

	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/types"
	// このブロックや引数リストを閉じる。
)

// この行で独自の型を定義する。
type roundTripFunc func(*http.Request) (*http.Response, error)

// この行でメソッド定義を始める。
func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	// この時点で関数の結果を返して処理を終了する。
	return f(req)
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func newTestService(transport http.RoundTripper, timeout time.Duration) *HotpepperService {
	// 条件を評価し、真のときだけ次の処理を実行する。
	if transport == nil {
		// この行で変数やフィールドへ値を代入する。
		transport = http.DefaultTransport
		// ここで処理ブロックを終了する。
	}

	// この時点で関数の結果を返して処理を終了する。
	return &HotpepperService{
		// この行は処理の一部として必要な操作を実行する。
		apiKey: "test-key",
		// この行は処理の一部として必要な操作を実行する。
		baseURL: "http://example.test",
		// 構造体やマップのフィールドへ値を設定する。
		client: &http.Client{
			// この行は処理の一部として必要な操作を実行する。
			Timeout: timeout,
			// この行は処理の一部として必要な操作を実行する。
			Transport: transport,
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type errorTransport struct{ err error }

// この行でメソッド定義を始める。
func (t errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	// この時点で関数の結果を返して処理を終了する。
	return nil, t.err
	// ここで処理ブロックを終了する。
}

// この行で構造体の型定義を始める。
type failingReadCloser struct{}

// この行でメソッド定義を始める。
func (failingReadCloser) Read([]byte) (int, error) {
	// この時点で関数の結果を返して処理を終了する。
	return 0, errors.New("read failed")
	// ここで処理ブロックを終了する。
}

// この行でメソッド定義を始める。
func (failingReadCloser) Close() error {
	// この時点で関数の結果を返して処理を終了する。
	return nil
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func TestSearchGourmet_NetworkAndFormatErrors(t *testing.T) {
	// サブテストを作成してケースごとに検証する。
	t.Run("invalid JSON", func(t *testing.T) {
		// この行で新しい変数を宣言しつつ値を代入する。
		svc := newTestService(roundTripFunc(func(*http.Request) (*http.Response, error) {
			// この時点で関数の結果を返して処理を終了する。
			return &http.Response{
				// この行は処理の一部として必要な操作を実行する。
				StatusCode: http.StatusOK,
				// 構造体やマップのフィールドへ値を設定する。
				Body: io.NopCloser(strings.NewReader(`{"broken":`)),
				// この行は処理の一部として必要な操作を実行する。
				Header: make(http.Header),
				// この行は処理の一部として必要な操作を実行する。
			}, nil
			// この行は処理の一部として必要な操作を実行する。
		}), time.Second)

		// 戻り値の一部を捨て、エラーだけを受け取る。
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		// 条件を評価し、真のときだけ次の処理を実行する。
		if err == nil {
			// この行は処理の一部として必要な操作を実行する。
			t.Fatal("expected error, got nil")
			// ここで処理ブロックを終了する。
		}
		// 条件を評価し、真のときだけ次の処理を実行する。
		if !strings.Contains(err.Error(), "failed to decode response") {
			// この行は処理の一部として必要な操作を実行する。
			t.Fatalf("expected decode error, got: %v", err)
			// ここで処理ブロックを終了する。
		}
		// この行は処理の一部として必要な操作を実行する。
	})

	// サブテストを作成してケースごとに検証する。
	t.Run("timeout", func(t *testing.T) {
		// この行で新しい変数を宣言しつつ値を代入する。
		svc := newTestService(roundTripFunc(func(req *http.Request) (*http.Response, error) {
			// この行で新しいブロックまたはリテラルを開始する。
			select {
			// この条件に一致した場合の処理を記述する。
			case <-time.After(500 * time.Millisecond):
				// この時点で関数の結果を返して処理を終了する。
				return &http.Response{
					// この行は処理の一部として必要な操作を実行する。
					StatusCode: http.StatusOK,
					// 構造体やマップのフィールドへ値を設定する。
					Body: io.NopCloser(strings.NewReader(`{"results":{"shop":[]}}`)),
					// この行は処理の一部として必要な操作を実行する。
					Header: make(http.Header),
					// この行は処理の一部として必要な操作を実行する。
				}, nil
			// この条件に一致した場合の処理を記述する。
			case <-req.Context().Done():
				// この時点で関数の結果を返して処理を終了する。
				return nil, req.Context().Err()
				// ここで処理ブロックを終了する。
			}
			// この行は処理の一部として必要な操作を実行する。
		}), 50*time.Millisecond)

		// 戻り値の一部を捨て、エラーだけを受け取る。
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		// 条件を評価し、真のときだけ次の処理を実行する。
		if err == nil {
			// この行は処理の一部として必要な操作を実行する。
			t.Fatal("expected error, got nil")
			// ここで処理ブロックを終了する。
		}
		// 条件を評価し、真のときだけ次の処理を実行する。
		if !strings.Contains(err.Error(), "failed to fetch API") {
			// この行は処理の一部として必要な操作を実行する。
			t.Fatalf("expected fetch error, got: %v", err)
			// ここで処理ブロックを終了する。
		}
		// この行は処理の一部として必要な操作を実行する。
	})

	// サブテストを作成してケースごとに検証する。
	t.Run("transport failure", func(t *testing.T) {
		// この行で新しい変数を宣言しつつ値を代入する。
		svc := newTestService(errorTransport{err: errors.New("dial tcp: no route to host")}, time.Second)

		// 戻り値の一部を捨て、エラーだけを受け取る。
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		// 条件を評価し、真のときだけ次の処理を実行する。
		if err == nil {
			// この行は処理の一部として必要な操作を実行する。
			t.Fatal("expected error, got nil")
			// ここで処理ブロックを終了する。
		}
		// 条件を評価し、真のときだけ次の処理を実行する。
		if !strings.Contains(err.Error(), "failed to fetch API") {
			// この行は処理の一部として必要な操作を実行する。
			t.Fatalf("expected fetch error, got: %v", err)
			// ここで処理ブロックを終了する。
		}
		// この行は処理の一部として必要な操作を実行する。
	})

	// サブテストを作成してケースごとに検証する。
	t.Run("non-200 status", func(t *testing.T) {
		// この行で新しい変数を宣言しつつ値を代入する。
		svc := newTestService(roundTripFunc(func(*http.Request) (*http.Response, error) {
			// この時点で関数の結果を返して処理を終了する。
			return &http.Response{
				// この行は処理の一部として必要な操作を実行する。
				StatusCode: http.StatusBadGateway,
				// この行は処理の一部として必要な操作を実行する。
				Body: io.NopCloser(strings.NewReader(`upstream failure`)),
				// この行は処理の一部として必要な操作を実行する。
				Header: make(http.Header),
				// この行は処理の一部として必要な操作を実行する。
			}, nil
			// この行は処理の一部として必要な操作を実行する。
		}), time.Second)

		// 戻り値の一部を捨て、エラーだけを受け取る。
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		// 条件を評価し、真のときだけ次の処理を実行する。
		if err == nil {
			// この行は処理の一部として必要な操作を実行する。
			t.Fatal("expected error, got nil")
			// ここで処理ブロックを終了する。
		}
		// 条件を評価し、真のときだけ次の処理を実行する。
		if !strings.Contains(err.Error(), "API returned status code 502") {
			// この行は処理の一部として必要な操作を実行する。
			t.Fatalf("expected status code error, got: %v", err)
			// ここで処理ブロックを終了する。
		}
		// この行は処理の一部として必要な操作を実行する。
	})

	// サブテストを作成してケースごとに検証する。
	t.Run("failed to read response body", func(t *testing.T) {
		// この行で新しい変数を宣言しつつ値を代入する。
		svc := newTestService(roundTripFunc(func(*http.Request) (*http.Response, error) {
			// この時点で関数の結果を返して処理を終了する。
			return &http.Response{
				// この行は処理の一部として必要な操作を実行する。
				StatusCode: http.StatusOK,
				// 構造体やマップのフィールドへ値を設定する。
				Body: failingReadCloser{},
				// この行は処理の一部として必要な操作を実行する。
				Header: make(http.Header),
				// この行は処理の一部として必要な操作を実行する。
			}, nil
			// この行は処理の一部として必要な操作を実行する。
		}), time.Second)

		// 戻り値の一部を捨て、エラーだけを受け取る。
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
		// 条件を評価し、真のときだけ次の処理を実行する。
		if err == nil {
			// この行は処理の一部として必要な操作を実行する。
			t.Fatal("expected error, got nil")
			// ここで処理ブロックを終了する。
		}
		// 条件を評価し、真のときだけ次の処理を実行する。
		if !strings.Contains(err.Error(), "failed to read response") {
			// この行は処理の一部として必要な操作を実行する。
			t.Fatalf("expected read error, got: %v", err)
			// ここで処理ブロックを終了する。
		}
		// この行は処理の一部として必要な操作を実行する。
	})
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func TestSearchGourmet_ClampStartAndCount(t *testing.T) {
	// この行で新しい変数を宣言しつつ値を代入する。
	testCases := []struct {
		// この行は処理の一部として必要な操作を実行する。
		name string
		// この行は処理の一部として必要な操作を実行する。
		start int
		// この行は処理の一部として必要な操作を実行する。
		count int
		// この行は処理の一部として必要な操作を実行する。
		wantStart string
		// この行は処理の一部として必要な操作を実行する。
		wantCount string
		// この行で新しいブロックまたはリテラルを開始する。
	}{
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "start below min is clamped to 1",
			// この行は処理の一部として必要な操作を実行する。
			start: 0,
			// この行は処理の一部として必要な操作を実行する。
			count: 10,
			// この行は処理の一部として必要な操作を実行する。
			wantStart: "1",
			// この行は処理の一部として必要な操作を実行する。
			wantCount: "10",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "count above max is clamped to 100",
			// この行は処理の一部として必要な操作を実行する。
			start: 1,
			// この行は処理の一部として必要な操作を実行する。
			count: 200,
			// この行は処理の一部として必要な操作を実行する。
			wantStart: "1",
			// この行は処理の一部として必要な操作を実行する。
			wantCount: "100",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここで処理ブロックを終了する。
	}

	// 繰り返し処理を開始する。
	for _, tc := range testCases {
		// サブテストを作成してケースごとに検証する。
		t.Run(tc.name, func(t *testing.T) {
			// この行で変数を定義する。
			var gotStart, gotCount string
			// この行で新しい変数を宣言しつつ値を代入する。
			svc := newTestService(roundTripFunc(func(r *http.Request) (*http.Response, error) {
				// この行で変数やフィールドへ値を代入する。
				gotStart = r.URL.Query().Get("start")
				// この行で変数やフィールドへ値を代入する。
				gotCount = r.URL.Query().Get("count")
				// この時点で関数の結果を返して処理を終了する。
				return &http.Response{
					// この行は処理の一部として必要な操作を実行する。
					StatusCode: http.StatusOK,
					// 構造体やマップのフィールドへ値を設定する。
					Body: io.NopCloser(strings.NewReader(`{"results":{"shop":[]}}`)),
					// この行は処理の一部として必要な操作を実行する。
					Header: make(http.Header),
					// この行は処理の一部として必要な操作を実行する。
				}, nil
				// この行は処理の一部として必要な操作を実行する。
			}), time.Second)

			// 戻り値の一部を捨て、エラーだけを受け取る。
			_, err := svc.SearchGourmet(types.GourmetSearchParams{
				// この行は処理の一部として必要な操作を実行する。
				Keyword: "sushi",
				// この行は処理の一部として必要な操作を実行する。
				Start: tc.start,
				// この行は処理の一部として必要な操作を実行する。
				Count: tc.count,
				// この行は処理の一部として必要な操作を実行する。
			})
			// 条件を評価し、真のときだけ次の処理を実行する。
			if err != nil {
				// この行は処理の一部として必要な操作を実行する。
				t.Fatalf("unexpected error: %v", err)
				// ここで処理ブロックを終了する。
			}

			// 条件を評価し、真のときだけ次の処理を実行する。
			if gotStart != tc.wantStart {
				// この行で変数やフィールドへ値を代入する。
				t.Fatalf("start = %q, want %q", gotStart, tc.wantStart)
				// ここで処理ブロックを終了する。
			}
			// 条件を評価し、真のときだけ次の処理を実行する。
			if gotCount != tc.wantCount {
				// この行で変数やフィールドへ値を代入する。
				t.Fatalf("count = %q, want %q", gotCount, tc.wantCount)
				// ここで処理ブロックを終了する。
			}
			// この行は処理の一部として必要な操作を実行する。
		})
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func TestSearchGourmet_HotpepperAPIError(t *testing.T) {
	// この行で新しい変数を宣言しつつ値を代入する。
	svc := newTestService(roundTripFunc(func(*http.Request) (*http.Response, error) {
		// この時点で関数の結果を返して処理を終了する。
		return &http.Response{
			// この行は処理の一部として必要な操作を実行する。
			StatusCode: http.StatusOK,
			// この行は処理の一部として必要な操作を実行する。
			Body: io.NopCloser(strings.NewReader(
				// 構造体やマップのフィールドへ値を設定する。
				`{"results":{"error":[{"code":3000,"message":"パラメータ不正"}]}}`,
			// この行は処理の一部として必要な操作を実行する。
			)),
			// この行は処理の一部として必要な操作を実行する。
			Header: make(http.Header),
			// この行は処理の一部として必要な操作を実行する。
		}, nil
		// この行は処理の一部として必要な操作を実行する。
	}), time.Second)

	// 戻り値の一部を捨て、エラーだけを受け取る。
	_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "sushi"})
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err == nil {
		// この行は処理の一部として必要な操作を実行する。
		t.Fatal("expected error, got nil")
		// ここで処理ブロックを終了する。
	}

	// この行で変数を定義する。
	var apiErr *types.HotpepperAPIError
	// 条件を評価し、真のときだけ次の処理を実行する。
	if !errors.As(err, &apiErr) {
		// この行は処理の一部として必要な操作を実行する。
		t.Fatalf("expected *types.HotpepperAPIError, got: %T (%v)", err, err)
		// ここで処理ブロックを終了する。
	}
	// 条件を評価し、真のときだけ次の処理を実行する。
	if apiErr.Code != 3000 {
		// この行で変数やフィールドへ値を代入する。
		t.Fatalf("expected code=3000, got: %d", apiErr.Code)
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}
