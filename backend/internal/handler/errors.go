// このファイルが属するパッケージ名を宣言する。
package handler

// ここから使用する外部パッケージの列挙を始める。
import (
	// この行は処理の一部として必要な操作を実行する。
	"encoding/json"
	// この行は処理の一部として必要な操作を実行する。
	"errors"
	// この行は処理の一部として必要な操作を実行する。
	"log"
	// この行は処理の一部として必要な操作を実行する。
	"net/http"
	// この行は処理の一部として必要な操作を実行する。
	"regexp"

	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/types"
	// このブロックや引数リストを閉じる。
)

// この行で変数を定義する。
var apiKeyQueryPattern = regexp.MustCompile(`([?&]key=)[^&\s]+`)

// この行は補足のためのコメントを記述する。
// handleServiceError はservice層のエラーを判別し、適切なHTTPステータスと
// この行は補足のためのコメントを記述する。
// ユーザー向けメッセージ（自前定義）をJSONで返す。
// この行は補足のためのコメントを記述する。
// ホットペッパーAPIの生のエラーメッセージはログのみに出力し、フロントには返さない。
// この行で関数定義を始める。
func handleServiceError(w http.ResponseWriter, err error) {
	// この行で変数を定義する。
	var apiErr *types.HotpepperAPIError
	// 条件を評価し、真のときだけ次の処理を実行する。
	if errors.As(err, &apiErr) {
		// デバッグや障害調査のためにログを出力する。
		log.Printf("hotpepper API error handled: code=%d", apiErr.Code)
		// 値に応じて分岐する処理を開始する。
		switch apiErr.Code {
		// この条件に一致した場合の処理を記述する。
		case 3000:
			// この行は補足のためのコメントを記述する。
			// パラメータ不正: ユーザーが修正可能
			// この行は処理の一部として必要な操作を実行する。
			writeErrorJSON(w, http.StatusBadRequest, "検索条件が正しくありません")
			// この時点で関数処理を終了する。
			return
		// この条件に一致した場合の処理を記述する。
		case 1000:
			// この行は補足のためのコメントを記述する。
			// ホットペッパー側のサーバ障害: 上流障害を示す502
			// この行は処理の一部として必要な操作を実行する。
			writeErrorJSON(w, http.StatusBadGateway, "サービスが一時的に利用できません")
			// この時点で関数処理を終了する。
			return
		// この条件に一致した場合の処理を記述する。
		case 2000:
			// この行は補足のためのコメントを記述する。
			// APIキー/IP認証エラー: バックエンド設定ミスのため500
			// この行は処理の一部として必要な操作を実行する。
			writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
			// この時点で関数処理を終了する。
			return
			// ここで処理ブロックを終了する。
		}
		// ここで処理ブロックを終了する。
	}
	// デバッグや障害調査のためにログを出力する。
	log.Printf("internal server error: %s", maskAPIKeyInLog(err))
	// この行は処理の一部として必要な操作を実行する。
	writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func maskAPIKeyInLog(err error) string {
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err == nil {
		// この時点で関数の結果を返して処理を終了する。
		return ""
		// ここで処理ブロックを終了する。
	}

	// この時点で関数の結果を返して処理を終了する。
	return apiKeyQueryPattern.ReplaceAllString(err.Error(), "${1}[REDACTED]")
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func writeErrorJSON(w http.ResponseWriter, status int, message string) {
	// HTTPレスポンスへ値を書き込む処理を行う。
	w.Header().Set("Content-Type", "application/json")
	// HTTPレスポンスへ値を書き込む処理を行う。
	w.WriteHeader(status)
	// JSON形式への変換または解析を行う。
	json.NewEncoder(w).Encode(map[string]any{
		// 構造体やマップのフィールドへ値を設定する。
		"error": map[string]string{"message": message},
		// この行は処理の一部として必要な操作を実行する。
	})
	// ここで処理ブロックを終了する。
}
