package routes

import (
	"algoBharat/backend/pkg/handlers"
	"algoBharat/backend/pkg/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, movieHandler *handlers.MovieHandler, theatreHandler *handlers.TheatreHandler, hallHandler *handlers.HallHandler, showHandler *handlers.ShowHandler, bookingHandler *handlers.BookingHandler, analyticsHandler *handlers.AnalyticsHandler, userHandler *handlers.UserHandler) {

	// --- Public Routes --- (No authentication required)
	// Anyone can register or log in.
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Anyone can view movies, theatres, halls, and shows.
	r.HandleFunc("/movies", movieHandler.GetMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", movieHandler.GetMovie).Methods("GET")
	r.HandleFunc("/theatres", theatreHandler.GetTheatres).Methods("GET")
	r.HandleFunc("/theatres/{id}", theatreHandler.GetTheatre).Methods("GET")
	r.HandleFunc("/halls", hallHandler.GetHalls).Methods("GET")
	r.HandleFunc("/halls/{id}", hallHandler.GetHall).Methods("GET")
	r.HandleFunc("/halls/{id}/seats", hallHandler.GetHallSeats).Methods("GET")
	r.HandleFunc("/shows", showHandler.GetShows).Methods("GET")

	// --- Authenticated Routes --- (Requires a valid token, any role)
	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.AuthMiddleware)

	// Only logged-in users can create a booking.
	authRouter.HandleFunc("/bookings", bookingHandler.CreateBooking).Methods("POST")
	authRouter.HandleFunc("/bookings", bookingHandler.GetBookings).Methods("GET") // Added GET /bookings

	// --- Admin Routes --- (Requires a valid token with 'admin' role)
	adminRouter := r.PathPrefix("/").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware, middleware.AdminOnlyMiddleware)

	// Only admins can create, update, or delete movies, theatres, halls, and shows.
	adminRouter.HandleFunc("/movies", movieHandler.CreateMovie).Methods("POST")
	adminRouter.HandleFunc("/movies/{id}", movieHandler.UpdateMovie).Methods("PUT")
	adminRouter.HandleFunc("/movies/{id}", movieHandler.DeleteMovie).Methods("DELETE")

	adminRouter.HandleFunc("/theatres", theatreHandler.CreateTheatre).Methods("POST")
	adminRouter.HandleFunc("/theatres/{id}", theatreHandler.UpdateTheatre).Methods("PUT")
	adminRouter.HandleFunc("/theatres/{id}", theatreHandler.DeleteTheatre).Methods("DELETE")

	adminRouter.HandleFunc("/halls", hallHandler.CreateHall).Methods("POST")
	adminRouter.HandleFunc("/shows", showHandler.CreateShow).Methods("POST")

	// Only admins can view revenue analytics.
	adminRouter.HandleFunc("/analytics/movies/{id}/revenue", analyticsHandler.GetMovieRevenue).Methods("GET")

	// Admin User Management Routes
	adminRouter.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	adminRouter.HandleFunc("/users/{id}/role", userHandler.UpdateUserRole).Methods("PUT")
}
