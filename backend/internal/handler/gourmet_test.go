// このファイルが属するパッケージ名を宣言する。
package handler

// ここから使用する外部パッケージの列挙を始める。
import (
	// この行は処理の一部として必要な操作を実行する。
	"encoding/json"
	// この行は処理の一部として必要な操作を実行する。
	"net/http"
	// この行は処理の一部として必要な操作を実行する。
	"net/http/httptest"
	// この行は処理の一部として必要な操作を実行する。
	"net/url"
	// この行は処理の一部として必要な操作を実行する。
	"strings"
	// この行は処理の一部として必要な操作を実行する。
	"testing"

	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/types"
	// このブロックや引数リストを閉じる。
)

// この行で関数定義を始める。
func TestParseParams_ValidationErrors(t *testing.T) {
	// この行は補足のためのコメントを記述する。
	// service はバリデーション失敗パスしか通らないため nil のままで安全
	// この行で新しい変数を宣言しつつ値を代入する。
	h := &GourmetHandler{}

	// この行で新しい変数を宣言しつつ値を代入する。
	testCases := []struct {
		// この行は処理の一部として必要な操作を実行する。
		name string
		// この行は処理の一部として必要な操作を実行する。
		query url.Values
		// この行は処理の一部として必要な操作を実行する。
		wantContains string
		// この行で新しいブロックまたはリテラルを開始する。
	}{
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "all params missing",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "either lat/lng, address, or keyword must be provided",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "lat only",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"lat": []string{"35.0"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "lat and lng must be provided together",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "lng only",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"lng": []string{"139.0"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "lat and lng must be provided together",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "lat not number",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"lat": []string{"abc"},
				// 構造体やマップのフィールドへ値を設定する。
				"lng": []string{"139.0"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "lat must be a number",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "lat NaN",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"lat": []string{"NaN"},
				// 構造体やマップのフィールドへ値を設定する。
				"lng": []string{"139.0"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "lat must be a finite number",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "lat Inf",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"lat": []string{"Inf"},
				// 構造体やマップのフィールドへ値を設定する。
				"lng": []string{"139.0"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "lat must be a finite number",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "lng NaN",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"lat": []string{"35.0"},
				// 構造体やマップのフィールドへ値を設定する。
				"lng": []string{"NaN"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "lng must be a finite number",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "lng not number",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"lat": []string{"35.0"},
				// 構造体やマップのフィールドへ値を設定する。
				"lng": []string{"abc"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "lng must be a number",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "range out of bounds",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"lat": []string{"35.0"},
				// 構造体やマップのフィールドへ値を設定する。
				"lng": []string{"139.0"},
				// 構造体やマップのフィールドへ値を設定する。
				"range": []string{"6"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "range must be 1, 2, 3, 4, or 5",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "range without location",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"keyword": []string{"寿司"},
				// 構造体やマップのフィールドへ値を設定する。
				"range": []string{"1"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "range requires lat and lng",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "start not integer",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"keyword": []string{"寿司"},
				// 構造体やマップのフィールドへ値を設定する。
				"start": []string{"abc"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "start must be an integer",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "count not integer",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"keyword": []string{"寿司"},
				// 構造体やマップのフィールドへ値を設定する。
				"count": []string{"abc"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// この行は処理の一部として必要な操作を実行する。
			wantContains: "count must be an integer",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここで処理ブロックを終了する。
	}

	// 繰り返し処理を開始する。
	for _, tc := range testCases {
		// サブテストを作成してケースごとに検証する。
		t.Run(tc.name, func(t *testing.T) {
			// この行で新しい変数を宣言しつつ値を代入する。
			path := "/"
			// 条件を評価し、真のときだけ次の処理を実行する。
			if encoded := tc.query.Encode(); encoded != "" {
				// この行で変数やフィールドへ値を代入する。
				path += "?" + encoded
				// ここで処理ブロックを終了する。
			}
			// この行で新しい変数を宣言しつつ値を代入する。
			r := httptest.NewRequest(http.MethodGet, path, nil)

			// 戻り値の一部を捨て、エラーだけを受け取る。
			_, err := h.parseParams(r)
			// 条件を評価し、真のときだけ次の処理を実行する。
			if err == nil {
				// この行は処理の一部として必要な操作を実行する。
				t.Fatal("expected error, got nil")
				// ここで処理ブロックを終了する。
			}
			// 条件を評価し、真のときだけ次の処理を実行する。
			if !strings.Contains(err.Error(), tc.wantContains) {
				// この行で変数やフィールドへ値を代入する。
				t.Fatalf("error = %q, want contains %q", err.Error(), tc.wantContains)
				// ここで処理ブロックを終了する。
			}
			// この行は処理の一部として必要な操作を実行する。
		})
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func TestHandle_BadRequestReturnsJSON(t *testing.T) {
	// この行は補足のためのコメントを記述する。
	// service はバリデーション失敗パスしか通らないため nil のままで安全
	// この行で新しい変数を宣言しつつ値を代入する。
	h := &GourmetHandler{}
	// この行で新しい変数を宣言しつつ値を代入する。
	wantMessage := "either lat/lng, address, or keyword must be provided"

	// この行で新しい変数を宣言しつつ値を代入する。
	rec := httptest.NewRecorder()
	// この行で新しい変数を宣言しつつ値を代入する。
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	// この行は処理の一部として必要な操作を実行する。
	h.Handle(rec, r)

	// 条件を評価し、真のときだけ次の処理を実行する。
	if rec.Code != http.StatusBadRequest {
		// この行で変数やフィールドへ値を代入する。
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	contentType := rec.Header().Get("Content-Type")
	// 条件を評価し、真のときだけ次の処理を実行する。
	if !strings.Contains(contentType, "application/json") {
		// この行で変数やフィールドへ値を代入する。
		t.Fatalf("content-type = %q, want contains application/json", contentType)
		// ここで処理ブロックを終了する。
	}

	// この行で変数を定義する。
	var body map[string]any
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		// この行は処理の一部として必要な操作を実行する。
		t.Fatalf("failed to decode response JSON: %v", err)
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	errorObj, ok := body["error"].(map[string]any)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if !ok {
		// この行は処理の一部として必要な操作を実行する。
		t.Fatalf("missing error object in response: %#v", body)
		// ここで処理ブロックを終了する。
	}
	// この行で新しい変数を宣言しつつ値を代入する。
	gotMessage, ok := errorObj["message"].(string)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if !ok {
		// この行は処理の一部として必要な操作を実行する。
		t.Fatalf("missing error.message in response: %#v", body)
		// ここで処理ブロックを終了する。
	}
	// 条件を評価し、真のときだけ次の処理を実行する。
	if gotMessage != wantMessage {
		// この行で変数やフィールドへ値を代入する。
		t.Fatalf("error.message = %q, want %q", gotMessage, wantMessage)
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func TestParseParams_ValidCases(t *testing.T) {
	// この行で新しい変数を宣言しつつ値を代入する。
	h := &GourmetHandler{}

	// この行で新しい変数を宣言しつつ値を代入する。
	testCases := []struct {
		// この行は処理の一部として必要な操作を実行する。
		name string
		// この行は処理の一部として必要な操作を実行する。
		query url.Values
		// この行は処理の一部として必要な操作を実行する。
		want types.GourmetSearchParams
		// この行で新しいブロックまたはリテラルを開始する。
	}{
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "lat/lng only",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"lat": []string{"35.6895"},
				// 構造体やマップのフィールドへ値を設定する。
				"lng": []string{"139.6917"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// 構造体やマップのフィールドへ値を設定する。
			want: types.GourmetSearchParams{Lat: 35.6895, Lng: 139.6917, Start: 1, Count: 20},
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "keyword only",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"keyword": []string{"寿司"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// 構造体やマップのフィールドへ値を設定する。
			want: types.GourmetSearchParams{Keyword: "寿司", Start: 1, Count: 20},
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "address only",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"address": []string{"  渋谷  "},
				// この行は処理の一部として必要な操作を実行する。
			},
			// 構造体やマップのフィールドへ値を設定する。
			want: types.GourmetSearchParams{Address: "渋谷", Start: 1, Count: 20},
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "keyword is trimmed",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"keyword": []string{"  ランチ  "},
				// この行は処理の一部として必要な操作を実行する。
			},
			// 構造体やマップのフィールドへ値を設定する。
			want: types.GourmetSearchParams{Keyword: "ランチ", Start: 1, Count: 20},
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "lat/lng with range",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"lat": []string{"35.0"},
				// 構造体やマップのフィールドへ値を設定する。
				"lng": []string{"139.0"},
				// 構造体やマップのフィールドへ値を設定する。
				"range": []string{"3"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// 構造体やマップのフィールドへ値を設定する。
			want: types.GourmetSearchParams{Lat: 35.0, Lng: 139.0, Range: 3, Start: 1, Count: 20},
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "keyword with start and count",
			// 構造体やマップのフィールドへ値を設定する。
			query: url.Values{
				// 構造体やマップのフィールドへ値を設定する。
				"keyword": []string{"ラーメン"},
				// 構造体やマップのフィールドへ値を設定する。
				"start": []string{"11"},
				// 構造体やマップのフィールドへ値を設定する。
				"count": []string{"5"},
				// この行は処理の一部として必要な操作を実行する。
			},
			// 構造体やマップのフィールドへ値を設定する。
			want: types.GourmetSearchParams{Keyword: "ラーメン", Start: 11, Count: 5},
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここで処理ブロックを終了する。
	}

	// 繰り返し処理を開始する。
	for _, tc := range testCases {
		// サブテストを作成してケースごとに検証する。
		t.Run(tc.name, func(t *testing.T) {
			// この行で新しい変数を宣言しつつ値を代入する。
			r := httptest.NewRequest(http.MethodGet, "/?"+tc.query.Encode(), nil)

			// この行で新しい変数を宣言しつつ値を代入する。
			got, err := h.parseParams(r)
			// 条件を評価し、真のときだけ次の処理を実行する。
			if err != nil {
				// この行は処理の一部として必要な操作を実行する。
				t.Fatalf("unexpected error: %v", err)
				// ここで処理ブロックを終了する。
			}
			// 条件を評価し、真のときだけ次の処理を実行する。
			if got != tc.want {
				// この行で変数やフィールドへ値を代入する。
				t.Fatalf("parseParams() = %+v, want %+v", got, tc.want)
				// ここで処理ブロックを終了する。
			}
			// この行は処理の一部として必要な操作を実行する。
		})
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}
