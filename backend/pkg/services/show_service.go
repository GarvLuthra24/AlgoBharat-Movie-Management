package services

import "algoBharat/backend/pkg/models"

// ShowService defines the interface for show-related business logic.
type ShowService interface {
	GetShows() ([]models.Show, error)
	CreateShow(show models.Show) (models.Show, error)
}
