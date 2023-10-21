package handlers

import (
	"encoding/json"
	"intern-net/internal/app/models"
	"intern-net/internal/app/repositories"
	"intern-net/internal/app/services"
	"net/http"
)

type LoginHandler struct {
	UserRepository *repositories.UserRepository
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.Login(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func NewLoginHandler(userRepo *repositories.UserRepository) *LoginHandler {
	return &LoginHandler{
		UserRepository: userRepo,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials

	// Decode JSON data from the request body into the Credentials struct
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get user credentials from the database
	user, err := h.UserRepository.GetUserByEmail(r.Context(), creds.Email)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verify the password
	if !services.VerifyPassword(user.Password, creds.Password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token for the authenticated user
	token, err := services.GenerateToken(user.Email, user.Role)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond with the JWT token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}
