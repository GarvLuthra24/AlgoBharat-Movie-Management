package services

import "algoBharat/backend/pkg/models"

// Credentials represents the user's login credentials.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserService defines the interface for user-related business logic.
type UserService interface {
	Register(credentials Credentials) (models.User, error)
	Login(credentials Credentials) (string, error) // Returns a JWT token string
	GetUsers() ([]models.User, error)
	UpdateUserRole(userID string, newRole string) (models.User, error)
}
