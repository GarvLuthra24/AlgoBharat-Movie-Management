package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// getJWTKey returns the JWT secret key from environment variables
func getJWTKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "my_secret_key" // Default fallback
	}
	return []byte(secret)
}

// AuthMiddleware verifies the JWT token from the Authorization header.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the token from the header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 2. The header should be in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims := &jwt.RegisteredClaims{}

		// 3. Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return getJWTKey(), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 4. Token is valid. Store user info in the request context for downstream handlers.
		ctx := context.WithValue(r.Context(), "userID", claims.Subject)
		ctx = context.WithValue(ctx, "userRole", claims.Issuer) // We stored the role in the Issuer field
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnlyMiddleware checks if the user has the 'admin' role.
// This middleware MUST run AFTER AuthMiddleware.
func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the user role from the context (set by AuthMiddleware)
		role := r.Context().Value("userRole") // Corrected from GetValue to Value
		if role == nil {
			http.Error(w, "User role not found in context", http.StatusInternalServerError)
			return
		}

		// 2. Check if the role is 'admin'
		if role.(string) != "admin" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}

		// 3. User is an admin. Proceed to the next handler.
		next.ServeHTTP(w, r)
	})
}
