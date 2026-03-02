package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"backend/internal/types"
)

var apiKeyQueryPattern = regexp.MustCompile(`([?&]key=)[^&\s]+`)

// handleServiceError はservice層のエラーを判別し、適切なHTTPステータスと
// ユーザー向けメッセージ（自前定義）をJSONで返す。
// ホットペッパーAPIの生のエラーメッセージはログのみに出力し、フロントには返さない。
func handleServiceError(w http.ResponseWriter, err error) {
	var apiErr *types.HotpepperAPIError
	if errors.As(err, &apiErr) {
		log.Printf("hotpepper API error handled: code=%d", apiErr.Code)
		switch apiErr.Code {
		case 3000:
			// パラメータ不正: ユーザーが修正可能
			writeErrorJSON(w, http.StatusBadRequest, "検索条件が正しくありません")
			return
		case 1000:
			// ホットペッパー側のサーバ障害: 上流障害を示す502
			writeErrorJSON(w, http.StatusBadGateway, "サービスが一時的に利用できません")
			return
		case 2000:
			// APIキー/IP認証エラー: バックエンド設定ミスのため500
			writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
			return
		}
	}
	log.Printf("internal server error: %s", maskAPIKeyInLog(err))
	writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
}

func maskAPIKeyInLog(err error) string {
	if err == nil {
		return ""
	}

	return apiKeyQueryPattern.ReplaceAllString(err.Error(), "${1}[REDACTED]")
}

func writeErrorJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"message": message},
	})
}
