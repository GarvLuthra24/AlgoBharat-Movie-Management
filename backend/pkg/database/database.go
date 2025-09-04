package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"algoBharat/backend/pkg/models"
)

var DB *sql.DB

func InitDB() {
	var err error

	// Get database configuration from environment variables
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "sqlite3" // Default to SQLite for backward compatibility
	}

	var dsn string

	if dbDriver == "mysql" {
		// MySQL configuration
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			dbHost = "localhost"
		}
		dbPort := os.Getenv("DB_PORT")
		if dbPort == "" {
			dbPort = "3306"
		}
		dbUser := os.Getenv("DB_USER")
		if dbUser == "" {
			dbUser = "root"
		}
		dbPassword := os.Getenv("DB_PASSWORD")
		if dbPassword == "" {
			dbPassword = "password"
		}
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "algoBharat"
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser, dbPassword, dbHost, dbPort, dbName)
	} else {
		// SQLite configuration (default)
		dsn = "./movies.db"
	}

	DB, err = sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to %s database", dbDriver)

	createTables()
	createDefaultAdmin() // Call the function to create default admin
}

func createTables() {
	createMoviesTable := `
	CREATE TABLE IF NOT EXISTS movies (
		id TEXT PRIMARY KEY,
		title TEXT,
		duration_minutes INTEGER
	);
	`
	createTheatresTable := `
	CREATE TABLE IF NOT EXISTS theatres (
		id TEXT PRIMARY KEY,
		name TEXT
	);
	`
	createHallsTable := `
	CREATE TABLE IF NOT EXISTS halls (
		id TEXT PRIMARY KEY,
		name TEXT,
		theatre_id TEXT,
		seat_map TEXT
	);
	`
	createShowsTable := `
	CREATE TABLE IF NOT EXISTS shows (
		id TEXT PRIMARY KEY,
		movie_id TEXT,
		hall_id TEXT,
		time TEXT,
		price REAL
	);
	`
	createSeatsTable := `
	CREATE TABLE IF NOT EXISTS seats (
		id TEXT PRIMARY KEY,
		row INTEGER,
		number INTEGER,
		hall_id TEXT,
		column INTEGER
	);
	`
	createBookingsTable := `
	CREATE TABLE IF NOT EXISTS bookings (
		id TEXT PRIMARY KEY,
		show_id TEXT,
		seat_ids TEXT
	);
	`
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL
	);
	`

	_, err := DB.Exec(createMoviesTable + createTheatresTable + createHallsTable + createShowsTable + createSeatsTable + createBookingsTable + createUsersTable)
	if err != nil {
		log.Fatal(err)
	}
}

func createDefaultAdmin() {
	// Check if admin user already exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", "admin").Scan(&count)
	if err != nil {
		log.Printf("Error checking for default admin: %v", err)
		return
	}

	if count == 0 {
		// Admin user does not exist, create it
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing default admin password: %v", err)
			return
		}

		adminUser := models.User{
			ID:           strconv.Itoa(1), // Assign a fixed ID for the default admin
			Username:     "admin",
			PasswordHash: string(hashedPassword),
			Role:         "admin",
		}

		stmt, err := DB.Prepare("INSERT INTO users(id, username, password_hash, role) VALUES(?, ?, ?, ?)")
		if err != nil {
			log.Printf("Error preparing default admin insert statement: %v", err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(adminUser.ID, adminUser.Username, adminUser.PasswordHash, adminUser.Role)
		if err != nil {
			log.Printf("Error inserting default admin: %v", err)
			return
		}
		log.Println("Default admin user created: admin/admin123")
	} else {
		log.Println("Default admin user already exists.")
	}
}
