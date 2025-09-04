package services

// AnalyticsService defines the interface for analytics-related business logic.
type AnalyticsService interface {
	GetMovieRevenue(movieID string) (float64, error)
}
