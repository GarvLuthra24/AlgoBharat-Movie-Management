package handlers

import (
	"algoBharat/backend/pkg/models"
	"algoBharat/backend/pkg/services"
	"algoBharat/backend/pkg/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// HallHandler handles HTTP requests for halls.
type HallHandler struct {
	service services.HallService
}

// NewHallHandler creates a new HallHandler.
func NewHallHandler(service services.HallService) *HallHandler {
	return &HallHandler{service: service}
}

// GetHalls handles the GET /halls request.
func (h *HallHandler) GetHalls(w http.ResponseWriter, r *http.Request) {
	theatreID := r.URL.Query().Get("theatreId") // Get theatreId from query parameter
	halls, err := h.service.GetHalls(theatreID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, halls)
}

// GetHall handles the GET /halls/{id} request.
func (h *HallHandler) GetHall(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hall, err := h.service.GetHall(params["id"])
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, hall)
}

// CreateHall handles the POST /halls request.
func (h *HallHandler) CreateHall(w http.ResponseWriter, r *http.Request) {
	var hall models.Hall
	if err := json.NewDecoder(r.Body).Decode(&hall); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdHall, err := h.service.CreateHall(hall)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, createdHall)
}

// UpdateHall handles the PUT /halls/{id} request.
func (h *HallHandler) UpdateHall(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var hall models.Hall
	if err := json.NewDecoder(r.Body).Decode(&hall); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hall.ID = params["id"] // Ensure the ID from URL is used
	updatedHall, err := h.service.UpdateHall(hall)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, updatedHall)
}

// DeleteHall handles the DELETE /halls/{id} request.
func (h *HallHandler) DeleteHall(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := h.service.DeleteHall(params["id"])
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Hall deleted successfully"})
}

// GetHallSeats handles the GET /halls/{id}/seats request.
func (h *HallHandler) GetHallSeats(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	seats, err := h.service.GetHallSeats(params["id"])
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, seats)
}
