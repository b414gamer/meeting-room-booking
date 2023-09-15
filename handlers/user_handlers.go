package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"meeting-room-booking/config"
	"meeting-room-booking/models"
	"net/http"

	"golang.org/x/crypto/argon2"
)

type RegistrationInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HashPassword(password string) string {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		log.Fatalf("Failed to generate random salt: %v", err)
	}

	// Use Argon2 to hash the password
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Concatenate the salt and the hash
	combined := append(salt, hash...)

	// Return the base64 encoded hash
	return base64.RawStdEncoding.EncodeToString(combined)
}

func VerifyPassword(stored string, password string) bool {
	// Decode the stored value
	decoded, err := base64.RawStdEncoding.DecodeString(stored)
	if err != nil {
		log.Fatalf("Failed to decode store value: %v", err)
	}
	// Extract the salt and hash
	salt := decoded[:16]
	storedHash := decoded[16:]

	// Compute the hash for the provided password
	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Compare the computed hash with the stored hash
	return string(computedHash) == string(storedHash)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	//Parse request body
	var input RegistrationInput
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

	//Hash the password (you can use a library like bcrypt but I use Argon2)
	hashedPassword := HashPassword(input.Password)

	//Create user model
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	//Save to database
	result := config.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		log.Printf("Error registering user: %v", result.Error)
		return
	}

	//Return response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}
