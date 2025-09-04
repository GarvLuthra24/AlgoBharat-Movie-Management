package services

import (
	"algoBharat/backend/pkg/database"
	"algoBharat/backend/pkg/models"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// getJWTKey returns the JWT secret key from environment variables
func getJWTKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "my_secret_key" // Default fallback
	}
	return []byte(secret)
}

// Claims defines the custom claims for the JWT.
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UserServiceImpl struct{}

// Register handles the creation of a new user.
func (s *UserServiceImpl) Register(credentials Credentials) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	newUser := models.User{
		ID:           strconv.Itoa(rand.Intn(1000000)),
		Username:     credentials.Username,
		PasswordHash: string(hashedPassword),
		Role:         "user", // Default role
	}

	stmt, err := database.DB.Prepare("INSERT INTO users(id, username, password_hash, role) VALUES(?, ?, ?, ?)")
	if err != nil {
		return models.User{}, err
	}
	_, err = stmt.Exec(newUser.ID, newUser.Username, newUser.PasswordHash, newUser.Role)
	if err != nil {
		// This could be a unique constraint violation if the username is taken
		return models.User{}, fmt.Errorf("could not create user: %w", err)
	}

	return newUser, nil
}

// Login verifies a user's credentials and returns a JWT token if they are valid.
func (s *UserServiceImpl) Login(credentials Credentials) (string, error) {
	var user models.User
	row := database.DB.QueryRow("SELECT id, username, password_hash, role FROM users WHERE username = ?", credentials.Username)
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid username or password")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Create the JWT claims, which includes the username and expiry time
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    user.Role, // Using issuer to store the role
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(getJWTKey())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUsers retrieves all users from the database.
func (s *UserServiceImpl) GetUsers() ([]models.User, error) {
	rows, err := database.DB.Query("SELECT id, username, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateUserRole updates the role of a specific user.
func (s *UserServiceImpl) UpdateUserRole(userID string, newRole string) (models.User, error) {
	stmt, err := database.DB.Prepare("UPDATE users SET role = ? WHERE id = ?")
	if err != nil {
		return models.User{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(newRole, userID)
	if err != nil {
		return models.User{}, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return models.User{}, fmt.Errorf("user with ID %s not found or role not changed", userID)
	}

	// Fetch the updated user to return
	var updatedUser models.User
	row := database.DB.QueryRow("SELECT id, username, role FROM users WHERE id = ?", userID)
	if err := row.Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.Role); err != nil {
		return models.User{}, fmt.Errorf("failed to retrieve updated user: %w", err)
	}

	return updatedUser, nil
}
