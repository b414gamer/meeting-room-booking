package handlers

import (
	"os"

	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"meeting-room-booking/config"
	"meeting-room-booking/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/argon2"
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

	if !IsValidPassword(input.Password) {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Create user model
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: HashPassword(input.Password),
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

func HashPassword(password string) string {
	//Generate a Salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatalf("Failed to generate salt: %v", err)
	}

	//Hash the password using Argon2
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	//Return the base64 encoded hash
	return base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash)
}

func IsValidPassword(password string) bool {
	// Example: Ensure password is at least 8 characters long
	return len(password) >= 8
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
		http.Error(w, "Invalid credentials : "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Validate the user credentials (this is a placeholder, you'll need to check against your database)
	user, err := models.ValidateUserCredentials(LoginReq.Email, LoginReq.Password)
	if err != nil {
		http.Error(w, "Invalid credentials : "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWTToken(user)
	if err != nil {
		http.Error(w, "Invalid credentials : "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

var jwtSecret = os.Getenv("JWT_SECRET") // Move this to an environment variable or config file in production

func generateJWTToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.UserID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires after 24 hours
	})

	return token.SignedString([]byte(jwtSecret))
}
