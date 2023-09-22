package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"meeting-room-booking/config"
	"meeting-room-booking/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type BookingRequest struct {
	RoomID    uint      `json:"room_id"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

func BookRoomHandler(w http.ResponseWriter, r *http.Request) {
	var bookingRequest BookingRequest
	err := json.NewDecoder(r.Body).Decode(&bookingRequest)
	if err != nil {
		http.Error(w, "Invalid request payload :", http.StatusBadRequest)
		return
	}

	//validations here
	if bookingRequest.RoomID == 0 || bookingRequest.StartTime.IsZero() || bookingRequest.EndTime.IsZero() {
		http.Error(w, "Room ID, start time and end time are required", http.StatusBadRequest)
		return
	}

	//Check room availability
	if bookingRequest.StartTime.After(bookingRequest.EndTime) {
		http.Error(w, "Start time cannot be after end time", http.StatusBadRequest)
		return
	}

	//Extract token from Authorization header
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		http.Error(w, "Missing auth token", http.StatusForbidden)
		return
	}

	//The token usually comes in format `Bearer <token>`, hence we split by space
	splitToken := strings.Split(tokenHeader, " ")
	if len(splitToken) != 2 {
		http.Error(w, "Invalid token format", http.StatusForbidden)
		return
	}
	tokenString := splitToken[1]

	//Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Printf("Token validation error: %v", err)
		http.Error(w, "Invalid token :", http.StatusForbidden)
		return
	}

	//Extract user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}
	userID, ok := claims["userID"].(float64) // jwt-go library decodes numbers as float64
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusForbidden)
		return
	}

	//Room availability check
	var overlappingBookings []models.Booking
	config.DB.Where("room_id = ? AND ((start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?) OR (start_time <= ? AND end_time >= ?))", bookingRequest.RoomID, bookingRequest.StartTime, bookingRequest.EndTime, bookingRequest.StartTime, bookingRequest.EndTime, bookingRequest.StartTime, bookingRequest.EndTime).Find(&overlappingBookings)

	if len(overlappingBookings) > 0 {
		http.Error(w, "Room is already booked for the specified time", http.StatusBadRequest)
		return
	}

	//Create booking in the database
	booking := models.Booking{
		RoomID:    bookingRequest.RoomID,
		StartTime: bookingRequest.StartTime,
		EndTime:   bookingRequest.EndTime,
	}

	//Set the UserID in the booking object
	booking.UserID = uint(userID)

	result := config.DB.Create(&booking)
	if result.Error != nil {
		http.Error(w, "Error booking room", http.StatusInternalServerError)
		log.Printf("Error booking room: %v", result.Error)
		return
	}

	//Return response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Room booked successfully",
		"booking": booking,
	})
}
