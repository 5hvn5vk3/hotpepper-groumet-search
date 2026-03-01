// このファイルが属するパッケージ名を宣言する。
package handler

// ここから使用する外部パッケージの列挙を始める。
import (
	// この行は処理の一部として必要な操作を実行する。
	"encoding/json"
	// この行は処理の一部として必要な操作を実行する。
	"net/http"
	// この行は処理の一部として必要な操作を実行する。
	"strings"

	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/service"
	// この行は処理の一部として必要な操作を実行する。
	"backend/internal/types"
	// このブロックや引数リストを閉じる。
)

// この行で構造体の型定義を始める。
type GourmetDetailHandler struct {
	// この行は処理の一部として必要な操作を実行する。
	service *service.HotpepperService
	// ここで処理ブロックを終了する。
}

// この行で関数定義を始める。
func NewGourmetDetailHandler(service *service.HotpepperService) *GourmetDetailHandler {
	// この時点で関数の結果を返して処理を終了する。
	return &GourmetDetailHandler{service: service}
	// ここで処理ブロックを終了する。
}

// この行でメソッド定義を始める。
func (h *GourmetDetailHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// 条件を評価し、真のときだけ次の処理を実行する。
	if r.Method != http.MethodGet {
		// この行は処理の一部として必要な操作を実行する。
		writeErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		// この時点で関数処理を終了する。
		return
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	shopID := strings.TrimSpace(r.URL.Query().Get("id"))
	// 条件を評価し、真のときだけ次の処理を実行する。
	if shopID == "" {
		// この行は処理の一部として必要な操作を実行する。
		writeErrorJSON(w, http.StatusBadRequest, "id is required")
		// この時点で関数処理を終了する。
		return
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	response, err := h.service.GetGourmetDetail(shopID)
	// 条件を評価し、真のときだけ次の処理を実行する。
	if err != nil {
		// この行は処理の一部として必要な操作を実行する。
		handleServiceError(w, err)
		// この時点で関数処理を終了する。
		return
		// ここで処理ブロックを終了する。
	}

	// 条件を評価し、真のときだけ次の処理を実行する。
	if len(response.Results.Shop) == 0 {
		// この行は処理の一部として必要な操作を実行する。
		writeErrorJSON(w, http.StatusNotFound, "店舗が見つかりません")
		// この時点で関数処理を終了する。
		return
		// ここで処理ブロックを終了する。
	}

	// この行で新しい変数を宣言しつつ値を代入する。
	detail := toShopDetailSupplement(response.Results.Shop[0])
	// この行で新しい変数を宣言しつつ値を代入する。
	body, err := json.Marshal(detail)
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

// この行で関数定義を始める。
func toShopDetailSupplement(shop types.Shop) types.ShopDetailSupplement {
	// この時点で関数の結果を返して処理を終了する。
	return types.ShopDetailSupplement{
		// この行は処理の一部として必要な操作を実行する。
		Open: strings.TrimSpace(shop.Open),
		// この行は処理の一部として必要な操作を実行する。
		Close: strings.TrimSpace(shop.Close),
		// この行は処理の一部として必要な操作を実行する。
		BudgetMemo: strings.TrimSpace(shop.BudgetMemo),
		// この行は処理の一部として必要な操作を実行する。
		Wifi: strings.TrimSpace(shop.Wifi),
		// この行は処理の一部として必要な操作を実行する。
		PrivateRoom: strings.TrimSpace(shop.PrivateRoom),
		// この行は処理の一部として必要な操作を実行する。
		NonSmoking: strings.TrimSpace(shop.NonSmoking),
		// この行は処理の一部として必要な操作を実行する。
		Parking: strings.TrimSpace(shop.Parking),
		// この行は処理の一部として必要な操作を実行する。
		Lunch: strings.TrimSpace(shop.Lunch),
		// この行は処理の一部として必要な操作を実行する。
		Midnight: strings.TrimSpace(shop.Midnight),
		// この行は処理の一部として必要な操作を実行する。
		ShopDetailMemo: strings.TrimSpace(shop.ShopDetailMemo),
		// ここで処理ブロックを終了する。
	}
	// ここで処理ブロックを終了する。
}
