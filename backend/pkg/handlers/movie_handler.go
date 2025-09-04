package handlers

import (
	"algoBharat/backend/pkg/models"
	"algoBharat/backend/pkg/services"
	"algoBharat/backend/pkg/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// MovieHandler handles HTTP requests for movies.
type MovieHandler struct {
	service services.MovieService
}

// NewMovieHandler creates a new MovieHandler.
func NewMovieHandler(service services.MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

// GetMovies handles the GET /movies request.
func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.service.GetMovies()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, movies)
}

// GetMovie handles the GET /movies/{id} request.
func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movie, err := h.service.GetMovie(params["id"])
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, movie)
}

// CreateMovie handles the POST /movies request.
func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdMovie, err := h.service.CreateMovie(movie)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, createdMovie)
}

// UpdateMovie handles the PUT /movies/{id} request.
func (h *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedMovie, err := h.service.UpdateMovie(params["id"], movie)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, updatedMovie)
}

// DeleteMovie handles the DELETE /movies/{id} request.
func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if err := h.service.DeleteMovie(params["id"]); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Movie deleted successfully"})
}
