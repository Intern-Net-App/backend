package handlers

import (
	"encoding/json"
	"intern-net/internal/app/models"
	"intern-net/internal/app/repositories"
	"intern-net/internal/app/services"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SignupHandler handles user registration.
type SignupHandler struct {
	UserRepository *repositories.UserRepository
}

func (h *SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.Signup(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// NewSignupHandler creates a new SignupHandler instance with the provided user repository.
func NewSignupHandler(userRepo *repositories.UserRepository) *SignupHandler {
	return &SignupHandler{
		UserRepository: userRepo,
	}
}

func (h *SignupHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials

	// Decode JSON data from the request body into the Registration struct
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Check if a user with the same email exists
	existingUser, err := h.UserRepository.GetUserByEmail(r.Context(), creds.Email)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if existingUser != nil {
		http.Error(w, "Conflict", http.StatusConflict)
		return
	}

	// Hash the user's password before saving it to the database
	hashedPassword, err := services.HashPassword(creds.Password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a user based on registration data
	user := models.User{
		ID:       primitive.NewObjectID(),
		Email:    creds.Email,
		Password: hashedPassword,
		Role:     "User", // Defaulting to User role
	}

	// Save the user to the database
	err = h.UserRepository.CreateUser(r.Context(), &user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // User registration successful
}
