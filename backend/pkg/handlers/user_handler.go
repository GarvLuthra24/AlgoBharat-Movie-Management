package handlers

import (
	"algoBharat/backend/pkg/services"
	"algoBharat/backend/pkg/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	service services.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register handles the POST /register request.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var creds services.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.Register(creds)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, user)
}

// Login handles the POST /login request.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds services.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := h.service.Login(creds)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}

// GetUsers handles the GET /users request.
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, users)
}

// UpdateUserRole handles the PUT /users/{id}/role request.
func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	var requestBody struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedUser, err := h.service.UpdateUserRole(userID, requestBody.Role)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, updatedUser)
}
