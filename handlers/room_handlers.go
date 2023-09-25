package handlers

import (
	"encoding/json"
	"log"
	"meeting-room-booking/config"
	"meeting-room-booking/models"
	"net/http"
	"time"
)

type RoomWithStatus struct {
	models.Room
	Status    string     `json:"status"`
	StartTime *time.Time `json:"startTime,omitempty"`
	EndTime   *time.Time `json:"endTime,omitempty"`
}

// func GetRoomStatus(roomID uint) string {
// 	var bookings []models.Booking
// 	currentTime := time.Now()

// 	//Fetch bookings for the room that are active (i.e.  current time is between start and end time)
// 	config.DB.Where("room_id = ? AND start_time <= ? AND end_time >= ?", roomID, currentTime, currentTime).Find(&bookings)

// 	if len(bookings) > 0 {
// 		return "Booked"
// 	}
// 	return "Available"
// }

func GetRoomStatus(roomID uint) (string, *time.Time, *time.Time) {
	var bookings []models.Booking
	currentTime := time.Now()

	// Fetch bookings for the room that are in the future or currently active
	config.DB.Where("room_id = ? AND end_time >= ?", roomID, currentTime).Order("start_time asc").Find(&bookings)

	// Debugging logs
	log.Printf("Current Time: %v", currentTime)
	log.Printf("Number of bookings found for RoomID %d: %d", roomID, len(bookings))
	for _, booking := range bookings {
		log.Printf("Booking for RoomID %d: StartTime: %v, EndTime: %v", booking.RoomID, booking.StartTime, booking.EndTime)
	}

	if len(bookings) > 0 {
		// If the start time of the first booking is in the future, we can show that the room is booked for a future time.
		if bookings[0].StartTime.After(currentTime) {
			return "Booked (Future)", &bookings[0].StartTime, &bookings[0].EndTime
		}
		// If the current time is between the start and end time of a booking, the room is currently booked.
		return "Booked", &bookings[0].StartTime, &bookings[0].EndTime
	}
	return "Available", nil, nil
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
		status, startTime, endTime := GetRoomStatus(room.RoomID)
		roomWithStatus = append(roomWithStatus, RoomWithStatus{room, status, startTime, endTime})
	}

	//Return the list of rooms with their statuses in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roomWithStatus)

	// //Return the list of rooms in the response
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(rooms)
}
