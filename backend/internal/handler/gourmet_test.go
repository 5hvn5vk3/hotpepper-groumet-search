package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"backend/internal/types"
)

// TestParseParams は 3-a のテスト。
// parseParams のバリデーションロジックをテーブル駆動で検証する。
// service が nil のまま GourmetHandler を生成できる。
//
// 【テーブル駆動テスト（Table-Driven Tests）とは？】
// 複数の入力パターンを構造体スライスにまとめ、同じ検証ロジックを繰り返し適用する
// Go 標準のテストパターン。標準ライブラリ内部でも多用されている。
// 追加が容易：新しいケースはスライスへの追記だけで完結し、検証ロジックの重複がない。
func TestParseParams(t *testing.T) {
	tests := []struct {
		name       string
		params     map[string]string
		wantErrMsg string
	}{
		{
			name:       "全パラメータ未指定",
			params:     map[string]string{},
			wantErrMsg: "either lat/lng, address, or keyword must be provided",
		},
		{
			name:       "lat のみ指定（lng なし）",
			params:     map[string]string{"lat": "35.0"},
			wantErrMsg: "lat and lng must be provided together",
		},
		{
			name:       "lng のみ指定（lat なし）",
			params:     map[string]string{"lng": "139.0"},
			wantErrMsg: "lat and lng must be provided together",
		},
		{
			name:       "lat に非数値",
			params:     map[string]string{"lat": "abc", "lng": "139.0"},
			wantErrMsg: "lat must be a number",
		},
		{
			name:       "lat に NaN",
			params:     map[string]string{"lat": "NaN", "lng": "139.0"},
			wantErrMsg: "lat must be a finite number",
		},
		{
			name:       "lat に Inf",
			params:     map[string]string{"lat": "Inf", "lng": "139.0"},
			wantErrMsg: "lat must be a finite number",
		},
		{
			name:       "lng に NaN",
			params:     map[string]string{"lat": "35.0", "lng": "NaN"},
			wantErrMsg: "lng must be a finite number",
		},
		{
			name:       "range=6（範囲外）",
			params:     map[string]string{"lat": "35.0", "lng": "139.0", "range": "6"},
			wantErrMsg: "range must be 1, 2, 3, 4, or 5",
		},
		{
			name:       "range のみ（位置なし）",
			params:     map[string]string{"keyword": "寿司", "range": "1"},
			wantErrMsg: "range requires lat and lng",
		},
		{
			name:       "start に非整数",
			params:     map[string]string{"keyword": "寿司", "start": "abc"},
			wantErrMsg: "start must be an integer",
		},
		{
			name:       "count に非整数",
			params:     map[string]string{"keyword": "寿司", "count": "abc"},
			wantErrMsg: "count must be an integer",
		},
	}

	// service が nil でもバリデーション段階では service に到達しないためテスト可能。
	h := &GourmetHandler{} // service は nil のまま

	for _, tc := range tests {
		// 【t.Run によるサブテスト】
		// 各ケースを独立したサブテストとして実行する。
		// 失敗時に「TestParseParams/lat のみ指定（lng なし）」のように
		// どのケースが落ちたか名前付きで出力されるため、デバッグが容易になる。
		t.Run(tc.name, func(t *testing.T) {
			// 【httptest.NewRequest とは？】
			// テスト用の *http.Request を生成するヘルパー関数。
			// 実際のネットワーク接続を行わず、指定したメソッド・URL・ボディから
			// リクエストオブジェクトを作成する。本番では net/http サーバーが
			// 受け取ったリクエストを渡すが、テストでは httptest.NewRequest で代替できる。
			q := url.Values{}
			for k, v := range tc.params {
				q.Set(k, v)
			}
			r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)

			_, err := h.parseParams(r)
			// err が nil だと後続の err.Error() 呼び出しでパニックするため t.Fatalf で即終了する。
			// 「後続処理が意味をなさない・危険になる場合」は t.Fatal を使う鉄則。
			if err == nil {
				t.Fatalf("エラーが返るべきなのに nil だった（期待: %q）", tc.wantErrMsg)
			}
			// 【strings.Contains による部分一致検証】
			// エラーは fmt.Errorf("...: %w", err) でラップされ前後に情報が付くことがある。
			// 完全一致（==）より部分一致の方が実装変更に対して堅牢なテストになる。
			if !strings.Contains(err.Error(), tc.wantErrMsg) {
				t.Errorf("エラーメッセージ %q が含まれない: %q", tc.wantErrMsg, err.Error())
			}
		})
	}
}

// TestHandle_ValidationError は 3-b のテスト。
// Handle が全パラメータ未指定のリクエストに対して 400 と JSON エラーを返すことを検証する。
func TestHandle_ValidationError(t *testing.T) {
	wantMessage := "either lat/lng, address, or keyword must be provided"

	// ① recorder と全パラメータ未指定リクエストを準備
	//
	// 【httptest.NewRecorder と httptest.NewRequest の役割まとめ】
	// - httptest.NewRecorder：HTTP レスポンスをメモリ上に記録する ResponseWriter の実装。
	//   rec.Code でステータスコード、rec.Body でレスポンスボディを取得できる。
	//   実際のサーバーを起動せずにハンドラーの出力を検査するために使う。
	// - httptest.NewRequest：テスト用の *http.Request を生成するヘルパー。
	//   実際のネットワーク接続を行わず、指定したメソッド・URL・ボディからリクエストを作る。
	// この2つを組み合わせることで、HTTP ハンドラーをサーバーなしで単体テストできる。
	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	h := &GourmetHandler{} // service は nil のまま

	// ② Handle を呼ぶ（バリデーションで弾かれるため service は呼ばれない）
	h.Handle(rec, r)

	// ③ 検証
	if rec.Code != http.StatusBadRequest {
		t.Errorf("HTTPステータス 400 を期待したが %d だった", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Content-Type に 'application/json' が含まれない: %q", contentType)
	}

	var body map[string]map[string]string
	// JSON デコード失敗時は後続の body 参照が意味をなさないため t.Fatalf で即終了する。
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディの JSON デコードに失敗: %v", err)
	}
	errObj, ok := body["error"]
	// "error" キーが存在しない場合は errObj へのアクセスが意味をなさないため t.Fatalf で終了する。
	if !ok {
		t.Fatalf("レスポンスボディに 'error' キーが存在しない: %v", body)
	}
	if _, ok := errObj["message"]; !ok {
		t.Errorf("error オブジェクトに 'message' キーが存在しない: %v", errObj)
	}
	if got := errObj["message"]; got != wantMessage {
		t.Errorf("error.message = %q, want %q", got, wantMessage)
	}
}

// TestParseParams_ValidCases は parseParams の正常系テスト。
// 正しい入力が正しく GourmetSearchParams に変換されることを検証する。
func TestParseParams_ValidCases(t *testing.T) {
	h := &GourmetHandler{} // service は nil のまま

	testCases := []struct {
		name   string
		params map[string]string
		want   types.GourmetSearchParams
	}{
		{
			name:   "lat/lng のみ",
			params: map[string]string{"lat": "35.6895", "lng": "139.6917"},
			want:   types.GourmetSearchParams{Lat: 35.6895, Lng: 139.6917, Start: 1, Count: 20},
		},
		{
			name:   "keyword のみ",
			params: map[string]string{"keyword": "寿司"},
			want:   types.GourmetSearchParams{Keyword: "寿司", Start: 1, Count: 20},
		},
		{
			name:   "lat/lng と range",
			params: map[string]string{"lat": "35.0", "lng": "139.0", "range": "3"},
			want:   types.GourmetSearchParams{Lat: 35.0, Lng: 139.0, Range: 3, Start: 1, Count: 20},
		},
		{
			name:   "keyword と start/count",
			params: map[string]string{"keyword": "ラーメン", "start": "11", "count": "5"},
			want:   types.GourmetSearchParams{Keyword: "ラーメン", Start: 11, Count: 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := url.Values{}
			for k, v := range tc.params {
				q.Set(k, v)
			}
			r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)

			got, err := h.parseParams(r)
			// 正常系なのにエラーが返った場合、後続の got 比較が意味をなさないため t.Fatalf で終了する。
			if err != nil {
				t.Fatalf("予期しないエラーが返った: %v", err)
			}
			if got != tc.want {
				t.Errorf("parseParams() = %+v, want %+v", got, tc.want)
			}
		})
	}
}
