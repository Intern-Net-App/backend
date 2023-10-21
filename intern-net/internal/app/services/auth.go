package services

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// Define structure of JWT claims.
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// Generates a JWT token for the user.
func GenerateToken(email, role string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &Claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // Token expires in 24 hours.
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Authentication Middleware ensures that a request is authenticated first
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromRequest(r)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := validateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userClaims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Retrieves the JWT Token from the request.
func extractTokenFromRequest(r *http.Request) string {
	// Get value of Authorization header
	authHeader := r.Header.Get("Authorization")

	// Header value should be in format 'Bearer {token}'
	if authHeader == "" {
		return ""
	}

	// Split header value to extract token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	// Return Token
	return parts[1]
}

// Validates and Parses JWT Token
func validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
