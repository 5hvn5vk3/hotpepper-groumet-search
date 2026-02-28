package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"backend/internal/types"
)

// handleServiceError はservice層のエラーを判別し、適切なHTTPステータスと
// ユーザー向けメッセージ（自前定義）をJSONで返す。
// ホットペッパーAPIの生のエラーメッセージはログのみに出力し、フロントには返さない。
func handleServiceError(w http.ResponseWriter, err error) {
	var apiErr *types.HotpepperAPIError
	if errors.As(err, &apiErr) {
		log.Printf("hotpepper API error handled: code=%d", apiErr.Code)
		switch apiErr.Code {
		case 3000:
			writeErrorJSON(w, http.StatusBadRequest, "検索条件が正しくありません")
			return
		case 2000:
			writeErrorJSON(w, http.StatusBadGateway, "サービスに接続できませんでした")
			return
		}
	}
	log.Printf("internal server error: %v", err)
	writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
}

func writeErrorJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"message": message},
	})
}
