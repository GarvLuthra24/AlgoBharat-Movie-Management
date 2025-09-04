package services

import "algoBharat/backend/pkg/models"

// TheatreService defines the interface for theatre-related business logic.
type TheatreService interface {
	GetTheatres() ([]models.Theatre, error)
	GetTheatre(id string) (models.Theatre, error)
	CreateTheatre(theatre models.Theatre) (models.Theatre, error)
	UpdateTheatre(id string, theatre models.Theatre) (models.Theatre, error)
	DeleteTheatre(id string) error
}
