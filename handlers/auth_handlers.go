package handlers

import (
	"encoding/json"
	"log"
	"meeting-room-booking/config"
	"meeting-room-booking/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type RegistrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input (you can add more validations as needed)
	if input.Email == "" || input.Password == "" || input.Name == "" {
		http.Error(w, "Name, email, and password are required", http.StatusBadRequest)
		return
	}

	// Create user model
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password, // Storing plain password for simplicity, but this is NOT recommended for real-world applications.
	}

	// Save to database
	result := config.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		log.Printf("Error registering user: %v", result.Error)
		return
	}

	// Return response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

// Define a struct for the login request payload
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var LoginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&LoginReq)
	if err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the user credentials (this is a placeholder, you'll need to check against your database)
	user, err := models.ValidateUserCredentials(LoginReq.Email, LoginReq.Password)
	if err != nil {
		http.Error(w, "Invalid email or password : "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWTToken(user)
	if err != nil {
		http.Error(w, "Failed to generate JWT token : "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

const jwtSecret = "YOUR_SECRET_KEY" // Move this to an environment variable or config file in production

func generateJWTToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.UserID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires after 24 hours
	})

	return token.SignedString([]byte(jwtSecret))
}
