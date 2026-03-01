package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend/internal/service"
	"backend/internal/types"
)

// GourmetDetailHandler は店舗詳細取得エンドポイントのハンドラー構造体。
// 店舗ID（id パラメータ）を受け取り、該当店舗の詳細情報を返す。
type GourmetDetailHandler struct {
	service *service.HotpepperService
}

// NewGourmetDetailHandler は GourmetDetailHandler のコンストラクタ関数。
func NewGourmetDetailHandler(service *service.HotpepperService) *GourmetDetailHandler {
	return &GourmetDetailHandler{service: service}
}

// Handle は HTTP リクエストを受け取り、店舗詳細を JSON で返すメソッド。
// 店舗IDがクエリパラメータに存在しない・店舗が見つからない場合は適切なエラーを返す。
func (h *GourmetDetailHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// 詳細取得も読み取り専用操作なので GET のみ許可する。
	if r.Method != http.MethodGet {
		writeErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// strings.TrimSpace でクエリパラメータの前後空白を除去してから空チェックを行う。
	// "id=" のように値が空の場合も "" として扱われ、400 を返す。
	shopID := strings.TrimSpace(r.URL.Query().Get("id"))
	if shopID == "" {
		// 必須パラメータが欠如しているため 400 Bad Request を返す。
		// "id is required" はクライアント開発者向けの明確なエラーメッセージ。
		writeErrorJSON(w, http.StatusBadRequest, "id is required")
		return
	}

	// service層に店舗IDを渡し、ホットペッパーAPIから詳細情報を取得する。
	response, err := h.service.GetGourmetDetail(shopID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// ホットペッパーAPIは存在しない店舗IDでも 200 OK を返し、
	// 結果配列（response.Results.Shop）が空になる仕様。
	// APIのレスポンスを正規化するため、ここで空チェックを行い 404 に変換する。
	//
	// 404 Not Found はリソースが存在しないことを示す。
	// 「店舗が見つかりません」はユーザーが理解できる日本語メッセージ。
	if len(response.Results.Shop) == 0 {
		writeErrorJSON(w, http.StatusNotFound, "店舗が見つかりません")
		return
	}

	// 検索結果配列の先頭要素が目的の店舗データ。
	// IDで検索しているため必ず1件か0件のみ返る。
	detail := toShopDetailSupplement(response.Results.Shop[0])
	body, err := json.Marshal(detail)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// toShopDetailSupplement は types.Shop（APIレスポンスの生データ）から
// 詳細画面に必要な補足情報（営業時間・予算・設備など）だけを抽出した
// types.ShopDetailSupplement 構造体に変換する関数。
//
// 全フィールドに strings.TrimSpace を適用する理由:
// ホットペッパーAPIのレスポンスには前後に空白や改行が入ることがある。
// フロントエンドでそのまま表示するとUIが崩れるため、ここで正規化しておく。
func toShopDetailSupplement(shop types.Shop) types.ShopDetailSupplement {
	return types.ShopDetailSupplement{
		Open:           strings.TrimSpace(shop.Open),
		Close:          strings.TrimSpace(shop.Close),
		BudgetMemo:     strings.TrimSpace(shop.BudgetMemo),
		Wifi:           strings.TrimSpace(shop.Wifi),
		PrivateRoom:    strings.TrimSpace(shop.PrivateRoom),
		NonSmoking:     strings.TrimSpace(shop.NonSmoking),
		Parking:        strings.TrimSpace(shop.Parking),
		Lunch:          strings.TrimSpace(shop.Lunch),
		Midnight:       strings.TrimSpace(shop.Midnight),
		ShopDetailMemo: strings.TrimSpace(shop.ShopDetailMemo),
	}
}
