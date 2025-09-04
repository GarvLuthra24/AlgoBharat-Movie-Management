package handlers

import (
	"algoBharat/backend/pkg/services"
	"algoBharat/backend/pkg/utils"
	"encoding/json"
	"net/http"
	"strings"
)

// BookingHandler handles HTTP requests for bookings.
type BookingHandler struct {
	service services.BookingService
}

// NewBookingHandler creates a new BookingHandler.
func NewBookingHandler(service services.BookingService) *BookingHandler { // Corrected parameter type
	return &BookingHandler{service: service}
}

// CreateBooking handles the POST /bookings request.
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var request services.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdBooking, err := h.service.CreateBooking(request)
	if err != nil {
		// Check if it's a no contiguous seats error or show not found error
		if _, ok := err.(*services.ErrNoContiguousSeats); ok ||
			(err != nil && (strings.Contains(err.Error(), "no show found") || strings.Contains(err.Error(), "no contiguous seats"))) {
			alternatives, altErr := h.service.FindAlternativeShows(request.Time, request.NumSeats)
			if altErr != nil {
				utils.RespondError(w, http.StatusInternalServerError, "Seats are booked and failed to find alternatives")
				return
			}
			if len(alternatives) > 0 {
				utils.RespondJSON(w, http.StatusConflict, map[string]interface{}{
					"message":      "Could not book seats together for the requested show. Here are some alternatives for the same day:",
					"alternatives": alternatives,
				})
			} else {
				utils.RespondJSON(w, http.StatusConflict, map[string]string{"message": "Could not book seats together for the requested show, and no same-day alternatives are available."})
			}
		} else {
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondJSON(w, http.StatusCreated, createdBooking)
}

// GetBookings handles the GET /bookings request.
func (h *BookingHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	showID := r.URL.Query().Get("showId")
	if showID == "" {
		utils.RespondError(w, http.StatusBadRequest, "Missing showId query parameter")
		return
	}

	bookings, err := h.service.GetBookingsByShowID(showID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, bookings)
}
