package main

import (
	"algoBharat/backend/pkg/database"
	"algoBharat/backend/pkg/handlers"
	"algoBharat/backend/pkg/routes"
	"algoBharat/backend/pkg/services"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	database.InitDB()

	// Create services
	movieService := &services.MovieServiceImpl{}
	theatreService := &services.TheatreServiceImpl{}
	hallService := &services.HallServiceImpl{}
	showService := &services.ShowServiceImpl{}
	bookingService := &services.BookingServiceImpl{}
	analyticsService := &services.AnalyticsServiceImpl{}
	userService := &services.UserServiceImpl{}

	// Create handlers
	movieHandler := handlers.NewMovieHandler(movieService)
	theatreHandler := handlers.NewTheatreHandler(theatreService)
	hallHandler := handlers.NewHallHandler(hallService)
	showHandler := handlers.NewShowHandler(showService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	userHandler := handlers.NewUserHandler(userService)

	r := mux.NewRouter()

	// Register routes
	routes.RegisterRoutes(
		r,
		movieHandler,
		theatreHandler,
		hallHandler,
		showHandler,
		bookingHandler,
		analyticsHandler,
		userHandler,
	)

	// Configure CORS
	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:5173"
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{corsOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(r)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
