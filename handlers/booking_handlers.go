package handlers

import (
	"encoding/json"
	"log"
	"meeting-room-booking/config"
	"meeting-room-booking/models"
	"net/http"
	"time"
)

type BookingRequest struct {
	RoomID uint	  `json:"room_id"`
	StartTime time.Time `json:"startTime"`
	EndTime time.Time `json:"endTime"`
}

func BookRoomHandler(w http.ResponseWriter, r *http.Request) {
	var bookingRequest BookingRequest
	err := json.NewDecoder(r.Body).Decode(&bookingRequest)
	if err != nil {
		http.Error(w, "Invalid request payload :", http.StatusBadRequest)
		return
	}

	// TODO: Add validations here
	if bookingRequest.RoomID == 0 || bookingRequest.StartTime.IsZero() || bookingRequest.EndTime.IsZero() {
		http.Error(w, "Room ID, start time and end time are required", http.StatusBadRequest)
		return
	}

	// TODO: Check room availability
	if bookingRequest.StartTime.After(bookingRequest.EndTime) {
		http.Error(w, "Start time cannot be after end time", http.StatusBadRequest)
		return
	}

	// TODO: Create booking in the database
	booking := models.Booking{
		RoomID: bookingRequest.RoomID,
		StartTime: bookingRequest.StartTime,
		EndTime: bookingRequest.EndTime,
	}

	result := config.DB.Create(&booking)
	if result.Error != nil {
		http.Error(w, "Error booking room", http.StatusInternalServerError)
		log.Printf("Error booking room: %v", result.Error)
		return
	}

	// TODO: Return response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Room booked successfully",
		"booking": booking,
	})

}