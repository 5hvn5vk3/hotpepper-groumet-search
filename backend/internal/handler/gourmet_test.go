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

	h := &GourmetHandler{} // service は nil のまま

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q := url.Values{}
			for k, v := range tc.params {
				q.Set(k, v)
			}
			r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)

			_, err := h.parseParams(r)
			if err == nil {
				t.Fatalf("エラーが返るべきなのに nil だった（期待: %q）", tc.wantErrMsg)
			}
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
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディの JSON デコードに失敗: %v", err)
	}
	errObj, ok := body["error"]
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
			if err != nil {
				t.Fatalf("予期しないエラーが返った: %v", err)
			}
			if got != tc.want {
				t.Errorf("parseParams() = %+v, want %+v", got, tc.want)
			}
		})
	}
}
