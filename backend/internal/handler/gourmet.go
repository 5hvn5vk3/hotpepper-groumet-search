package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"backend/internal/service"
	"backend/internal/types"
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {

	return e.Message
}

type GourmetHandler struct {
	service *service.HotpepperService
}

func NewGourmetHandler(service *service.HotpepperService) *GourmetHandler {

	return &GourmetHandler{service: service}
}

func (h *GourmetHandler) Handle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var (
		response *types.GourmetSearchResponse
		err      error
	)

	shopID := strings.TrimSpace(r.URL.Query().Get("id"))
	if shopID != "" {
		response, err = h.service.GetGourmetDetail(shopID)
	} else {

		params, parseErr := h.parseParams(r)

		if parseErr != nil {

			http.Error(w, parseErr.Error(), http.StatusBadRequest)
			return
		}

		response, err = h.service.SearchGourmet(params)
	}

	if err != nil {

		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	w.Write(body)
}

func (h *GourmetHandler) parseParams(r *http.Request) (types.GourmetSearchParams, error) {

	query := r.URL.Query()

	latStr := query.Get("lat")
	lngStr := query.Get("lng")
	rangeStr := query.Get("range")
	address := strings.TrimSpace(query.Get("address"))
	keyword := strings.TrimSpace(query.Get("keyword"))

	var lat, lng float64
	var rng int

	if latStr != "" && lngStr != "" {
		if v, err := strconv.ParseFloat(latStr, 64); err == nil {
			lat = v
		}
		if v, err := strconv.ParseFloat(lngStr, 64); err == nil {
			lng = v
		}
	}
	if rangeStr != "" {
		rng = parseIntOrDefault(rangeStr, 0)
	}

	start := parseIntOrDefault(query.Get("start"), 1)
	count := parseIntOrDefault(query.Get("count"), 20)

	hasLocation := lat != 0 && lng != 0
	hasAddress := address != ""
	hasKeyword := keyword != ""
	if !hasLocation && !hasAddress && !hasKeyword {
		return types.GourmetSearchParams{}, &ValidationError{Message: "either lat/lng, address, or keyword must be provided"}
	}

	return types.GourmetSearchParams{
		Address: address,
		Genre:   query.Get("genre"),
		Keyword: keyword,
		Start:   start,
		Count:   count,
		Lat:     lat,
		Lng:     lng,
		Range:   rng,
	}, nil
}

func parseIntOrDefault(value string, fallback int) int {

	if value == "" {
		return fallback
	}

	if v, err := strconv.Atoi(value); err == nil {

		return v
	}

	return fallback
}
