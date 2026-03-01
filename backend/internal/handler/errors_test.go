// このファイルが属するパッケージ名を宣言する。
package handler

// ここから使用する外部パッケージの列挙を始める。
import (
	// この行は処理の一部として必要な操作を実行する。
	"bytes"
	// この行は処理の一部として必要な操作を実行する。
	"encoding/json"
	// この行は処理の一部として必要な操作を実行する。
	"fmt"
	// この行は処理の一部として必要な操作を実行する。
	"log"
	// この行は処理の一部として必要な操作を実行する。
	"net/http"
	// この行は処理の一部として必要な操作を実行する。
	"net/http/httptest"
	// この行は処理の一部として必要な操作を実行する。
	"strings"
	// この行は処理の一部として必要な操作を実行する。
	"testing"

	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/types"
	// このブロックや引数リストを閉じる。
)

// この行で関数定義を始める。
func TestMaskAPIKeyInLog(t *testing.T) {
	// この行で新しい変数を宣言しつつ値を代入する。
	testCases := []struct {
		// この行は処理の一部として必要な操作を実行する。
		name string
		// この行は処理の一部として必要な操作を実行する。
		err error
		// この行は処理の一部として必要な操作を実行する。
		want string
		// この行で新しいブロックまたはリテラルを開始する。
	}{
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "? separator",
			// この行で変数やフィールドへ値を代入する。
			err: fmt.Errorf("url?key=SECRET&lat=35"),
			// この行で変数やフィールドへ値を代入する。
			want: "url?key=[REDACTED]&lat=35",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "& separator",
			// この行で変数やフィールドへ値を代入する。
			err: fmt.Errorf("url?lat=35&key=SECRET"),
			// この行で変数やフィールドへ値を代入する。
			want: "url?lat=35&key=[REDACTED]",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "no key",
			// この行は処理の一部として必要な操作を実行する。
			err: fmt.Errorf("network timeout"),
			// この行は処理の一部として必要な操作を実行する。
			want: "network timeout",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "nil",
			// この行は処理の一部として必要な操作を実行する。
			err: nil,
			// この行は処理の一部として必要な操作を実行する。
			want: "",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "symbol in key",
			// この行で変数やフィールドへ値を代入する。
			err: fmt.Errorf("url?key=A+B/C&lat=35"),
			// この行で変数やフィールドへ値を代入する。
			want: "url?key=[REDACTED]&lat=35",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "key is only param",
			// この行で変数やフィールドへ値を代入する。
			err: fmt.Errorf("url?key=SECRET"),
			// この行で変数やフィールドへ値を代入する。
			want: "url?key=[REDACTED]",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここで処理ブロックを終了する。
	}

	// 繰り返し処理を開始する。
	for _, tc := range testCases {
		// サブテストを作成してケースごとに検証する。
		t.Run(tc.name, func(t *testing.T) {
			// この行で新しい変数を宣言しつつ値を代入する。
			got := maskAPIKeyInLog(tc.err)
			// 条件を評価し、真のときだけ次の処理を実行する。
			if got != tc.want {
				// この行で変数やフィールドへ値を代入する。
				t.Fatalf("maskAPIKeyInLog() = %q, want %q", got, tc.want)
				// ここで処理ブロックを終了する。
			}
			// この行は処理の一部として必要な操作を実行する。
		})
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func TestHandleServiceError_MasksAPIKeyInLog(t *testing.T) {
	// この行で変数を定義する。
	var buf bytes.Buffer

	// この行で新しい変数を宣言しつつ値を代入する。
	original := log.Writer()
	// この行で新しいブロックまたはリテラルを開始する。
	t.Cleanup(func() {
		// デバッグや障害調査のためにログを出力する。
		log.SetOutput(original)
		// この行は処理の一部として必要な操作を実行する。
	})

	// デバッグや障害調査のためにログを出力する。
	log.SetOutput(&buf)

	// この行で新しい変数を宣言しつつ値を代入する。
	rec := httptest.NewRecorder()
	// この行で変数やフィールドへ値を代入する。
	handleServiceError(rec, fmt.Errorf("url?key=SECRETKEY&lat=35.0"))

	// 条件を評価し、真のときだけ次の処理を実行する。
	if strings.Contains(buf.String(), "SECRETKEY") {
		// この行は処理の一部として必要な操作を実行する。
		t.Fatalf("log leaked raw key: %s", buf.String())
		// ここで処理ブロックを終了する。
	}
	// 条件を評価し、真のときだけ次の処理を実行する。
	if !strings.Contains(buf.String(), "[REDACTED]") {
		// この行は処理の一部として必要な操作を実行する。
		t.Fatalf("masked key not found in log: %s", buf.String())
		// ここで処理ブロックを終了する。
	}
	// 条件を評価し、真のときだけ次の処理を実行する。
	if rec.Code != http.StatusInternalServerError {
		// この行で変数やフィールドへ値を代入する。
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		// ここで処理ブロックを終了する。
	}

	// この行で変数を定義する。
	var body map[string]map[string]string
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		// この行は処理の一部として必要な操作を実行する。
		t.Fatalf("failed to decode response JSON: %v", err)
		// ここで処理ブロックを終了する。
	}
	// 条件を評価し、真のときだけ次の処理を実行する。
	if got := body["error"]["message"]; got != "サービスが一時的に利用できません" {
		// この行で変数やフィールドへ値を代入する。
		t.Fatalf("message = %q, want %q", got, "サービスが一時的に利用できません")
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func TestHandleServiceError_HotpepperCodeMapping(t *testing.T) {
	// この行で新しい変数を宣言しつつ値を代入する。
	testCases := []struct {
		// この行は処理の一部として必要な操作を実行する。
		name string
		// この行は処理の一部として必要な操作を実行する。
		code int
		// この行は処理の一部として必要な操作を実行する。
		wantStatus int
		// この行は処理の一部として必要な操作を実行する。
		wantMessage string
		// この行で新しいブロックまたはリテラルを開始する。
	}{
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "3000 -> 400",
			// この行は処理の一部として必要な操作を実行する。
			code: 3000,
			// この行は処理の一部として必要な操作を実行する。
			wantStatus: http.StatusBadRequest,
			// この行は処理の一部として必要な操作を実行する。
			wantMessage: "検索条件が正しくありません",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "1000 -> 502",
			// この行は処理の一部として必要な操作を実行する。
			code: 1000,
			// この行は処理の一部として必要な操作を実行する。
			wantStatus: http.StatusBadGateway,
			// この行は処理の一部として必要な操作を実行する。
			wantMessage: "サービスが一時的に利用できません",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "2000 -> 500",
			// この行は処理の一部として必要な操作を実行する。
			code: 2000,
			// この行は処理の一部として必要な操作を実行する。
			wantStatus: http.StatusInternalServerError,
			// この行は処理の一部として必要な操作を実行する。
			wantMessage: "サービスが一時的に利用できません",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここから処理ブロックを開始する。
		{
			// この行は処理の一部として必要な操作を実行する。
			name: "unknown -> 500",
			// この行は処理の一部として必要な操作を実行する。
			code: 9999,
			// この行は処理の一部として必要な操作を実行する。
			wantStatus: http.StatusInternalServerError,
			// この行は処理の一部として必要な操作を実行する。
			wantMessage: "サービスが一時的に利用できません",
			// この行は処理の一部として必要な操作を実行する。
		},
		// ここで処理ブロックを終了する。
	}

	// 繰り返し処理を開始する。
	for _, tc := range testCases {
		// サブテストを作成してケースごとに検証する。
		t.Run(tc.name, func(t *testing.T) {
			// この行で新しい変数を宣言しつつ値を代入する。
			rec := httptest.NewRecorder()
			// 構造体やマップのフィールドへ値を設定する。
			handleServiceError(rec, &types.HotpepperAPIError{Code: tc.code})

			// 条件を評価し、真のときだけ次の処理を実行する。
			if rec.Code != tc.wantStatus {
				// この行で変数やフィールドへ値を代入する。
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
				// ここで処理ブロックを終了する。
			}

			// この行で変数を定義する。
			var body map[string]map[string]string
			// 条件を評価し、真のときだけ次の処理を実行する。
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				// この行は処理の一部として必要な操作を実行する。
				t.Fatalf("failed to decode response JSON: %v", err)
				// ここで処理ブロックを終了する。
			}

			// この行で新しい変数を宣言しつつ値を代入する。
			gotMessage := body["error"]["message"]
			// 条件を評価し、真のときだけ次の処理を実行する。
			if gotMessage != tc.wantMessage {
				// この行で変数やフィールドへ値を代入する。
				t.Fatalf("message = %q, want %q", gotMessage, tc.wantMessage)
				// ここで処理ブロックを終了する。
			}
			// この行は処理の一部として必要な操作を実行する。
		})
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}
