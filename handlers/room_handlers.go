package handlers

import (
	"encoding/json"
	"meeting-room-booking/config"
	"meeting-room-booking/models"
	"net/http"
	"time"
)

type RoomWithStatus struct {
	models.Room
	Status string `json:"status"`
}

func GetRoomStatus(roomID uint) string {
	var bookings []models.Booking
	currentTime := time.Now()

	//Fetch bookings for the room that are active (i.e.  current time is between start and end time)
	config.DB.Where("room_id = ? AND start_time <= ? AND end_time >= ?", roomID, currentTime, currentTime).Find(&bookings)

	if len(bookings) > 0 {
		return "Booked"
	}
	return "Available"
}

func ListRoomsHandler(w http.ResponseWriter, r *http.Request) {
	var rooms []models.Room

	//Fetch rooms from database
	result := config.DB.Find(&rooms)
	if result.Error != nil {
		http.Error(w, "Error fetching rooms", http.StatusInternalServerError)
		return
	}

	//Create a slice to hold rooms with their statuses
	var roomWithStatus []RoomWithStatus

	//Determine the status of each room
	for _, room := range rooms {
		status := GetRoomStatus(room.RoomID)
		roomWithStatus = append(roomWithStatus, RoomWithStatus{room, status})
	}

	//Return the list of rooms with their statuses in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roomWithStatus)

	// //Return the list of rooms in the response
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(rooms)
}
