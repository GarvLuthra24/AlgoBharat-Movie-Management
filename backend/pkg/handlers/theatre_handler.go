package handlers

import (
	"algoBharat/backend/pkg/models"
	"algoBharat/backend/pkg/services"
	"algoBharat/backend/pkg/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// TheatreHandler handles HTTP requests for theatres.
type TheatreHandler struct {
	service services.TheatreService
}

// NewTheatreHandler creates a new TheatreHandler.
func NewTheatreHandler(service services.TheatreService) *TheatreHandler {
	return &TheatreHandler{service: service}
}

// GetTheatres handles the GET /theatres request.
func (h *TheatreHandler) GetTheatres(w http.ResponseWriter, r *http.Request) {
	theatres, err := h.service.GetTheatres()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, theatres)
}

// GetTheatre handles the GET /theatres/{id} request.
func (h *TheatreHandler) GetTheatre(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	theatre, err := h.service.GetTheatre(params["id"])
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, theatre)
}

// CreateTheatre handles the POST /theatres request.
func (h *TheatreHandler) CreateTheatre(w http.ResponseWriter, r *http.Request) {
	var theatre models.Theatre
	if err := json.NewDecoder(r.Body).Decode(&theatre); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdTheatre, err := h.service.CreateTheatre(theatre)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, createdTheatre)
}

// UpdateTheatre handles the PUT /theatres/{id} request.
func (h *TheatreHandler) UpdateTheatre(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var theatre models.Theatre
	if err := json.NewDecoder(r.Body).Decode(&theatre); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedTheatre, err := h.service.UpdateTheatre(params["id"], theatre)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, updatedTheatre)
}

// DeleteTheatre handles the DELETE /theatres/{id} request.
func (h *TheatreHandler) DeleteTheatre(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if err := h.service.DeleteTheatre(params["id"]); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Theatre deleted successfully"})
}
