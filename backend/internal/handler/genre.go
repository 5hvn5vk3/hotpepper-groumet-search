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

		writeErrorJSON(w, http.StatusMethodNotAllowed, "Method not allowed")

		return
	}

	response, err := h.service.GetGenreMaster()

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
