// このファイルが属するパッケージ名を宣言する。
package handler

// ここから使用する外部パッケージの列挙を始める。
import (
	// この行は処理の一部として必要な操作を実行する。
	"encoding/json"
	// この行は処理の一部として必要な操作を実行する。
	"net/http"

	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/service"
	// このブロックや引数リストを閉じる。
)

// この行で構造体の型定義を始める。
type GenreHandler struct {
	// この行は処理の一部として必要な操作を実行する。
	service *service.HotpepperService
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func NewGenreHandler(service *service.HotpepperService) *GenreHandler {

	// この時点で関数の結果を返して処理を終了する。
	return &GenreHandler{service: service}
	// ここで処理ブロックを終了する。
}

// この行でメソッド定義を始める。
func (h *GenreHandler) Handle(w http.ResponseWriter, r *http.Request) {

	// 条件を評価し、真のときだけ次の処理を実行する。
	if r.Method != http.MethodGet {

		// この行は処理の一部として必要な操作を実行する。
		writeErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")

		// この時点で関数処理を終了する。
		return
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	response, err := h.service.GetGenreMaster()

	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この行は処理の一部として必要な操作を実行する。
		handleServiceError(w, err)
		// この時点で関数処理を終了する。
		return
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	body, err := json.Marshal(response)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この行は処理の一部として必要な操作を実行する。
		writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
		// この時点で関数処理を終了する。
		return
		// ここで処理ブロックを終了する。
	}

	// HTTPレスポンスへ値を書き込む処理を行う。
	w.Header().Set("Content-Type", "application/json")

	// HTTPレスポンスへ値を書き込む処理を行う。
	w.WriteHeader(http.StatusOK)

	// HTTPレスポンスへ値を書き込む処理を行う。
	w.Write(body)
	// ここで処理ブロックを終了する。
}
