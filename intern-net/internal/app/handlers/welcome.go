package handlers

import (
	"fmt"
	"intern-net/internal/app/services"
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	// Get JWT token from request context
	claims, ok := r.Context().Value("userClaims").(*services.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Return welcome message
	w.Write([]byte(fmt.Sprintf("Welcome, %s!", claims.Email)))
}
