package handlers

import (
	"algoBharat/backend/pkg/models"
	"algoBharat/backend/pkg/services"
	"algoBharat/backend/pkg/utils"
	"encoding/json"
	"net/http"
)

// ShowHandler handles HTTP requests for shows.
type ShowHandler struct {
	service services.ShowService
}

// NewShowHandler creates a new ShowHandler.
func NewShowHandler(service services.ShowService) *ShowHandler {
	return &ShowHandler{service: service}
}

// GetShows handles the GET /shows request.
func (h *ShowHandler) GetShows(w http.ResponseWriter, r *http.Request) {
	shows, err := h.service.GetShows()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, shows)
}

// CreateShow handles the POST /shows request.
func (h *ShowHandler) CreateShow(w http.ResponseWriter, r *http.Request) {
	var show models.Show
	if err := json.NewDecoder(r.Body).Decode(&show); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdShow, err := h.service.CreateShow(show)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, createdShow)
}
