package services

import "algoBharat/backend/pkg/models"

// BookingRequest represents the user's request to book seats.
type BookingRequest struct {
	MovieID  string `json:"movieId"`
	HallID   string `json:"hallId"`
	Time     string `json:"time"`
	NumSeats int    `json:"numSeats"`
}

// BookingService defines the interface for booking-related business logic.
type BookingService interface {
	// CreateBooking attempts to find and book a contiguous block of seats.
	CreateBooking(request BookingRequest) (models.Booking, error)
	// FindAlternativeShows finds other shows on the same day with enough consecutive seats.
	FindAlternativeShows(originalTime string, numSeats int) ([]models.Show, error)
	// GetBookingsByShowID retrieves all bookings for a specific show.
	GetBookingsByShowID(showID string) ([]models.Booking, error)
}
