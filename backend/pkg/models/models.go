package models

// Movie represents a movie
type Movie struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	DurationMinutes int    `json:"duration_minutes"` // Added DurationMinutes field
}

// Theatre represents a movie theatre
type Theatre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Hall represents a hall in a theatre
type Hall struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	TheatreID string         `json:"theatre_id"`
	SeatMap   map[string][]int `json:"seat_map"` // Updated to support column-based layout
}

// Show represents a movie show
type Show struct {
	ID      string  `json:"id"`
	MovieID string  `json:"movie_id"`
	HallID  string  `json:"hall_id"`
	Time    string  `json:"time"`
	Price   float64 `json:"price"` // Added Price field
}

// Seat represents a seat in a hall
type Seat struct {
	ID      string `json:"id"`
	Row     int    `json:"row"`
	Number  int    `json:"number"`
	HallID  string `json:"hall_id"`
	Column  int    `json:"column"` // Added to store the column number
}

// Booking represents a ticket booking
type Booking struct {
	ID      string   `json:"id"`
	ShowID  string   `json:"show_id"`
	SeatIDs []string `json:"seat_ids"`
}

// User represents an application user.
type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // Do not expose password hash in JSON responses
	Role         string `json:"role"` // e.g., "user", "admin"
}
