package services

import "algoBharat/backend/pkg/models"

// MovieService defines the interface for movie-related business logic.
type MovieService interface {
	GetMovies() ([]models.Movie, error)
	GetMovie(id string) (models.Movie, error)
	CreateMovie(movie models.Movie) (models.Movie, error)
	UpdateMovie(id string, movie models.Movie) (models.Movie, error)
	DeleteMovie(id string) error
}
