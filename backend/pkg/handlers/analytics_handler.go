package handlers

import (
	"algoBharat/backend/pkg/services"
	"algoBharat/backend/pkg/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// AnalyticsHandler handles HTTP requests for analytics.
type AnalyticsHandler struct {
	service services.AnalyticsService
}

// NewAnalyticsHandler creates a new AnalyticsHandler.
func NewAnalyticsHandler(service services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

// GetMovieRevenue handles the GET /analytics/movies/{id}/revenue request.
func (h *AnalyticsHandler) GetMovieRevenue(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieID := params["id"]

	revenue, err := h.service.GetMovieRevenue(movieID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"movie_id":      movieID,
		"total_revenue": revenue,
	})
}
