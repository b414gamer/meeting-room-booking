package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"meeting-room-booking/config"
	"meeting-room-booking/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
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

	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	//Room availability check
	var overlappingBookings []models.Booking
	config.DB.Where("room_id = ? AND ((start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?) OR (start_time <= ? AND end_time >= ?))", bookingRequest.RoomID, bookingRequest.StartTime, bookingRequest.EndTime, bookingRequest.StartTime, bookingRequest.EndTime, bookingRequest.StartTime, bookingRequest.EndTime).Find(&overlappingBookings)

	if len(overlappingBookings) > 0 {
		http.Error(w, "Room is already booked for the specified time", http.StatusBadRequest)
		return
	}

	// Adjust the time to UTC+7 before saving it to the database
	// Manually adjust the time to UTC+7 before saving it to the database
	adjustedStartTime := bookingRequest.StartTime.Add(-7 * time.Hour)
	adjustedEndTime := bookingRequest.EndTime.Add(-7 * time.Hour)

	//Create booking in the database
	booking := models.Booking{
		RoomID:    bookingRequest.RoomID,
		StartTime: adjustedStartTime,
		EndTime:   adjustedEndTime,
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

func ListBookingsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var bookings []models.Booking

	// Fetch bookings from the database for the specific user
	result := config.DB.Where("user_id = ?", uint(userID)).Find(&bookings)
	if result.Error != nil {
		http.Error(w, "Error fetching bookings", http.StatusInternalServerError)
		return
	}

	// Return the list of bookings in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func extractUserIDFromToken(r *http.Request) (uint, error) {
	// Extract token from Authorization header
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		return 0, fmt.Errorf("missing auth token")
	}

	// The token usually comes in format `Bearer <token>`, hence we split by space
	splitToken := strings.Split(tokenHeader, " ")
	if len(splitToken) != 2 {
		return 0, fmt.Errorf("invalid token format")
	}
	tokenString := splitToken[1]

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	// Extract user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	userID, ok := claims["userID"].(float64) // jwt-go library decodes numbers as float64
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	return uint(userID), nil
}

func UpdateBookingHandler(w http.ResponseWriter, r *http.Request) {
	// Extract booking ID from URL or another source
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	var updateRequest BookingRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate user and retrieve existing booking...
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var booking models.Booking
	result := config.DB.Where("booking_id = ? AND user_id = ?", bookingID, userID).First(&booking)
	if result.Error != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Validate new booking details...
	if updateRequest.StartTime.IsZero() || updateRequest.EndTime.IsZero() {
		http.Error(w, "Start time and end time are required", http.StatusBadRequest)
		return
	}

	if updateRequest.StartTime.After(updateRequest.EndTime) {
		http.Error(w, "Start time cannot be after end time", http.StatusBadRequest)
		return
	}

	// Check room availability for the new time slot...
	var overlappingBookings []models.Booking

	config.DB.Where("room_id = ? AND ((start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?) OR (start_time <= ? AND end_time >= ?))", booking.RoomID, updateRequest.StartTime, updateRequest.EndTime, updateRequest.StartTime, updateRequest.EndTime, updateRequest.StartTime, updateRequest.EndTime).Find(&overlappingBookings)

	if len(overlappingBookings) > 0 {
		http.Error(w, "Room is already booked for the specified time", http.StatusBadRequest)
		return
	}

	// Adjust the time to UTC+7 before saving it to the database
	adjustedStartTime := updateRequest.StartTime.Add(-7 * time.Hour)
	adjustedEndTime := updateRequest.EndTime.Add(-7 * time.Hour)

	log.Printf("Booking ID: %d, User ID: %d", bookingID, userID)

	// Update booking in the database...
	result = config.DB.Model(&booking).Updates(models.Booking{
		StartTime: adjustedStartTime,
		EndTime:   adjustedEndTime,
	})
	if result.Error != nil {
		http.Error(w, "Error updating booking", http.StatusInternalServerError)
		return
	}

	// Reload the updated booking object
	config.DB.First(&booking, bookingID)

	// Return response...
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Booking updated successfully",
		"booking": booking,
	})
}

func CancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	// Extract booking ID from URL
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	// Validate user and retrieve existing booking
	userID, err := extractUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Check if the booking exists and belongs to the user
	var booking models.Booking
	result := config.DB.Where("booking_id = ? AND user_id = ?", bookingID, userID).First(&booking)
	if result.Error != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Delete the booking
	result = config.DB.Delete(&booking)
	if result.Error != nil {
		http.Error(w, "Error deleting booking", http.StatusInternalServerError)
		return
	}

	// Return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Booking canceled successfully",
	})
}
