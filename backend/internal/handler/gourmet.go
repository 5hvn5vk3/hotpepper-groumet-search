package handler

import (
	"encoding/json"
	"math"
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

		writeErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	params, parseErr := h.parseParams(r)

	if parseErr != nil {

		writeErrorJSON(w, http.StatusBadRequest, parseErr.Error())
		return
	}

	response, err := h.service.SearchGourmet(params)

	if err != nil {
		handleServiceError(w, err)
		return
	}

	body, err := json.Marshal(response)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "サービスが一時的に利用できません")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	w.Write(body)
}

func (h *GourmetHandler) parseParams(r *http.Request) (types.GourmetSearchParams, error) {

	query := r.URL.Query()

	address := strings.TrimSpace(query.Get("address"))
	keyword := strings.TrimSpace(query.Get("keyword"))

	lat, lng, parseErr := parseLatLngStrict(query.Get("lat"), query.Get("lng"))
	if parseErr != nil {
		return types.GourmetSearchParams{}, parseErr
	}

	rng, parseErr := parseOptionalIntStrict("range", query.Get("range"), 0)
	if parseErr != nil {
		return types.GourmetSearchParams{}, parseErr
	}

	start, parseErr := parseOptionalIntStrict("start", query.Get("start"), 1)
	if parseErr != nil {
		return types.GourmetSearchParams{}, parseErr
	}

	count, parseErr := parseOptionalIntStrict("count", query.Get("count"), 20)
	if parseErr != nil {
		return types.GourmetSearchParams{}, parseErr
	}

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

func parseOptionalIntStrict(name, value string, fallback int) (int, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, &ValidationError{Message: name + " must be an integer"}
	}

	return parsed, nil
}

func parseLatLngStrict(latStr, lngStr string) (float64, float64, error) {
	latTrimmed := strings.TrimSpace(latStr)
	lngTrimmed := strings.TrimSpace(lngStr)

	hasLat := latTrimmed != ""
	hasLng := lngTrimmed != ""

	if hasLat != hasLng {
		return 0, 0, &ValidationError{Message: "lat and lng must be provided together"}
	}

	if !hasLat {
		return 0, 0, nil
	}

	lat, err := strconv.ParseFloat(latTrimmed, 64)
	if err != nil {
		return 0, 0, &ValidationError{Message: "lat must be a number"}
	}
	if math.IsNaN(lat) || math.IsInf(lat, 0) {
		return 0, 0, &ValidationError{Message: "lat must be a finite number"}
	}

	lng, err := strconv.ParseFloat(lngTrimmed, 64)
	if err != nil {
		return 0, 0, &ValidationError{Message: "lng must be a number"}
	}
	if math.IsNaN(lng) || math.IsInf(lng, 0) {
		return 0, 0, &ValidationError{Message: "lng must be a finite number"}
	}

	return lat, lng, nil
}
