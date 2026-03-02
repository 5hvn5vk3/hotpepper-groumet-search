package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend/internal/service"
	"backend/internal/types"
)

type GourmetDetailHandler struct {
	service *service.HotpepperService
}

func NewGourmetDetailHandler(service *service.HotpepperService) *GourmetDetailHandler {
	return &GourmetDetailHandler{service: service}
}

func (h *GourmetDetailHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	shopID := strings.TrimSpace(r.URL.Query().Get("id"))
	if shopID == "" {
		writeErrorJSON(w, http.StatusBadRequest, "id is required")
		return
	}

	response, err := h.service.GetGourmetDetail(shopID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if len(response.Results.Shop) == 0 {
		writeErrorJSON(w, http.StatusNotFound, "店舗が見つかりません")
		return
	}

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
