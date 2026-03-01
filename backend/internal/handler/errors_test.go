package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/types"
)

// TestMaskAPIKeyInLog は 2-a のテスト。
// maskAPIKeyInLog がAPIキーを [REDACTED] に置換することをテーブル駆動で検証する。
func TestMaskAPIKeyInLog(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "? 区切り",
			err:  fmt.Errorf("url?key=SECRET&lat=35"),
			want: "url?key=[REDACTED]&lat=35",
		},
		{
			name: "& 区切り",
			err:  fmt.Errorf("url?lat=35&key=SECRET"),
			want: "url?lat=35&key=[REDACTED]",
		},
		{
			name: "キーなし",
			err:  fmt.Errorf("network timeout"),
			want: "network timeout",
		},
		{
			name: "nil",
			err:  nil,
			want: "",
		},
		{
			name: "記号入りキー",
			err:  fmt.Errorf("url?key=A+B/C&lat=35"),
			want: "url?key=[REDACTED]&lat=35",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := maskAPIKeyInLog(tc.err)
			if got != tc.want {
				t.Errorf("maskAPIKeyInLog(%v) = %q, want %q", tc.err, got, tc.want)
			}
		})
	}
}

// TestHandleServiceError_MaskAPIKey は 2-b のテスト。
// handleServiceError がログに APIキーを出力しないことを検証する。
// log がグローバル状態を持つため t.Parallel() は付けない。
func TestHandleServiceError_MaskAPIKey(t *testing.T) {
	// ① 元の出力先を退避し、テスト終了時に確実に戻す
	original := log.Writer()
	t.Cleanup(func() { log.SetOutput(original) })

	// ② ログをメモリバッファに横取りする
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// ③ レスポンスを httptest.Recorder で捕捉する
	rec := httptest.NewRecorder()

	// ④ APIキーを含むエラーで handleServiceError を呼ぶ
	handleServiceError(rec, fmt.Errorf("url?key=SECRETKEY&lat=35.0"))

	// ⑤ 検証
	logOutput := buf.String()
	if strings.Contains(logOutput, "SECRETKEY") {
		t.Errorf("ログに生のAPIキー 'SECRETKEY' が含まれている: %s", logOutput)
	}
	if !strings.Contains(logOutput, "[REDACTED]") {
		t.Errorf("ログに '[REDACTED]' が含まれていない: %s", logOutput)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("HTTPステータス %d を期待したが %d だった", http.StatusInternalServerError, rec.Code)
	}
}

// TestHandleServiceError_ErrorCodeMapping は 2-c のテスト。
// HotpepperAPIError のコードが正しい HTTP ステータスとメッセージにマッピングされることをテーブル駆動で検証する。
func TestHandleServiceError_ErrorCodeMapping(t *testing.T) {
	tests := []struct {
		name           string
		code           int
		wantStatus     int
		wantMessage    string
	}{
		{
			name:        "code=3000 → 400",
			code:        3000,
			wantStatus:  http.StatusBadRequest,
			wantMessage: "検索条件が正しくありません",
		},
		{
			name:        "code=1000 → 502",
			code:        1000,
			wantStatus:  http.StatusBadGateway,
			wantMessage: "サービスが一時的に利用できません",
		},
		{
			name:        "code=2000 → 500",
			code:        2000,
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "サービスが一時的に利用できません",
		},
		{
			name:        "code=9999（未定義） → 500",
			code:        9999,
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "サービスが一時的に利用できません",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handleServiceError(rec, &types.HotpepperAPIError{Code: tc.code})

			if rec.Code != tc.wantStatus {
				t.Errorf("HTTPステータス %d を期待したが %d だった", tc.wantStatus, rec.Code)
			}

			var body map[string]map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("レスポンスボディの JSON デコードに失敗: %v", err)
			}
			gotMessage := body["error"]["message"]
			if gotMessage != tc.wantMessage {
				t.Errorf("error.message = %q, want %q", gotMessage, tc.wantMessage)
			}
		})
	}
}
