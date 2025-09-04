package services

import "algoBharat/backend/pkg/models"

// HallService defines the interface for hall-related business logic.
type HallService interface {
	GetHalls(theatreID string) ([]models.Hall, error) // Added theatreID parameter
	GetHall(id string) (models.Hall, error)
	CreateHall(hall models.Hall) (models.Hall, error)
	GetHallSeats(hallID string) ([]models.Seat, error)
}
