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

// 【http.RoundTripper インターフェースをモックに使う手法】
// http.RoundTripper は net/http の通信レイヤーを抽象化したインターフェース。
// 定義は：type RoundTripper interface { RoundTrip(*Request) (*Response, error) }
// http.Client の Transport フィールドにセットすることで、
// 実際の TCP 通信をテスト用関数に差し替えられる。
//
// roundTripFunc は「関数」を RoundTripper として使えるようにするアダプター型。
// 各テスト内でインラインの関数リテラルを渡すだけでモック Transport を定義できる。
//
// 【httptest.NewServer との違い】
//   - httptest.NewServer：実際の TCP リスナーを起動するローカルサーバー。
//     統合テスト向きだが、ポート確保・サーバー起動のオーバーヘッドがある。
//   - RoundTripper モック：TCP 通信を一切行わず関数で直接レスポンスを返す。
//     高速でネットワーク環境に依存しない単体テスト向き。
type roundTripFunc func(*http.Request) (*http.Response, error)

// roundTripFunc が http.RoundTripper インターフェースを満たすようにメソッドを実装する。
// この1メソッドだけで、関数リテラルを http.Client の Transport として渡せるようになる。
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
		// ステータス 200 だが閉じていない不正 JSON を返す Transport を定義する。
		// これによりデコード失敗のコードパスをテストできる。
		//
		// 【io.NopCloser の役割】
		// http.Response.Body は io.ReadCloser インターフェース（Read + Close が必要）。
		// strings.NewReader は io.Reader は満たすが io.ReadCloser は満たさない。
		// io.NopCloser は io.Reader を受け取り、Close を何もしない（Nop = No Operation）実装で
		// io.ReadCloser を満たすラッパーを返す。テストで Body を手軽に作る定番パターン。
		svc := newTestService(roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"broken":`)), // 閉じていない不正 JSON
				Header:     make(http.Header),
			}, nil
		}), 5*time.Second)
		_, err := svc.SearchGourmet(types.GourmetSearchParams{Keyword: "test"})

		// エラーが nil の場合は後続の strings.Contains 確認が意味をなさないため t.Fatal で終了する。
		if err == nil {
			t.Fatal("エラーが返るべきなのに nil だった")
		}
		// 【strings.Contains によるエラーメッセージの部分一致検証】
		// fmt.Errorf("...: %w", err) でラップされたエラーは前後に追加情報を持つ。
		// 完全一致より部分一致の方が実装の細部に依存せず、より堅牢なテストになる。
		if !strings.Contains(err.Error(), "failed to decode response") {
			t.Errorf("期待するエラー文字列 'failed to decode response' が含まれない: %v", err)
		}
	})

	t.Run("タイムアウト_failedToFetchAPI", func(t *testing.T) {
		// タイムアウトのテスト：クライアントのタイムアウト（50ms）より長い時間待つ Transport を使う。
		// req.Context().Done() を監視することで、クライアントがキャンセルしたことを検出できる。
		svc := newTestService(roundTripFunc(func(req *http.Request) (*http.Response, error) {
			select {
			case <-time.After(500 * time.Millisecond): // クライアントのタイムアウトより長く待つ
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"results":{}}`)),
					Header:     make(http.Header),
				}, nil
			case <-req.Context().Done():
				// Context がキャンセル（タイムアウト）されたときにここに入る
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
		// errorTransport を使って即座にエラーを返す。
		// "connection refused" のようなネットワーク障害をシミュレートする。
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
		// ホットペッパー API はエラーでも HTTP 200 を返し、
		// レスポンスボディの "error" フィールドにエラー情報を含める仕様。
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

	// 【errors.As によるエラー型の検査】
	// errors.As(err, &target) はエラーチェーン（fmt.Errorf %w でラップされた連鎖）を
	// 再帰的にたどり、target の型に一致するエラーを見つけたら target にセットして true を返す。
	// 単純な型アサーション（err.(*types.HotpepperAPIError)）と違い、
	// エラーがラップされていても内側の型を取り出せる点が優れている。
	// 例：fmt.Errorf("outer: %w", &HotpepperAPIError{}) のようなラップ構造でも検出できる。
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
//
// 【テーブル駆動テストと t.Run の組み合わせ】
// テーブル駆動テストと t.Run を組み合わせることで複数ケースを独立したサブテストとして実行できる。
// 失敗時に「TestSearchGourmet_ClampStartAndCount/start=0 は 1 にクランプ」のように
// どのケースが失敗したか名前付きで出力されるため、問題箇所の特定が容易になる。
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
			// Transport 内でリクエストの URL クエリパラメータを取得して、
			// clampInt による補正が正しく適用されているか検証する。
			// 実際のネットワーク通信は行わず、すぐに空のレスポンスを返す。
			var gotStart, gotCount string
			svc := newTestService(roundTripFunc(func(req *http.Request) (*http.Response, error) {
				// リクエスト送信直前のクエリパラメータを取得する
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
