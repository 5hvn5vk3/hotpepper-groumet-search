package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/service"
)

type GenreHandler struct {
	service *service.HotpepperService
}

func NewGenreHandler(service *service.HotpepperService) *GenreHandler {

	return &GenreHandler{service: service}
}

func (h *GenreHandler) Handle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	response, err := h.service.GetGenreMaster()

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
