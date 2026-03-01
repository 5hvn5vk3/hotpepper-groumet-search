package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/service"
)

// GenreHandler はジャンルマスター取得エンドポイントのハンドラー構造体。
// ジャンル一覧はホットペッパーAPIのマスターデータをそのまま返すシンプルな構成。
type GenreHandler struct {
	service *service.HotpepperService
}

// NewGenreHandler は GenreHandler のコンストラクタ関数。
// service（ビジネスロジック層）を外から注入することで、テスト時にモックへ差し替えが可能。
func NewGenreHandler(service *service.HotpepperService) *GenreHandler {

	return &GenreHandler{service: service}
}

// Handle は HTTP リクエストを受け取り、ジャンルマスターを JSON で返すメソッド。
// ジャンルマスターはパラメータ不要のシンプルな読み取り専用エンドポイント。
func (h *GenreHandler) Handle(w http.ResponseWriter, r *http.Request) {

	// GETメソッド以外は拒否する。マスターデータの参照なので GET のみ許可すれば十分。
	// 405 Method Not Allowed はHTTPの仕様（RFC）に沿ったステータスコード。
	if r.Method != http.MethodGet {

		writeErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")

		return
	}

	// service層に実際の外部API呼び出しを委譲する。
	// handler はHTTPの入出力だけを担当し、APIとの通信はservice層が責任を持つ。
	response, err := h.service.GetGenreMaster()

	if err != nil {
		// service層エラーは共通の handleServiceError で処理する（errors.go 参照）。
		// APIエラーコードに応じた 400/500/502 の使い分けも handleServiceError が行う。
		handleServiceError(w, err)
		return
	}

	// json.Marshal は response 構造体を JSON バイト列（[]byte）に変換する。
	// ここでのエラーは通常発生しないが、チャネル型など変換不能な型が混入した場合に失敗する。
	body, err := json.Marshal(response)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
		return
	}

	// Content-Type ヘッダーは WriteHeader より前に設定しなければならない。
	// 順序を守らないとヘッダーが正しく送信されないGoの仕様に注意。
	w.Header().Set("Content-Type", "application/json")

	// 200 OK: リクエストが正常に処理され、レスポンスボディにデータが含まれる。
	w.WriteHeader(http.StatusOK)

	w.Write(body)
}
