package handlers

import (
	"encoding/json"
	"meeting-room-booking/config"
	"meeting-room-booking/models"
	"net/http"
)

func ListRoomsHandler(w http.ResponseWriter, r *http.Request) {
	var rooms []models.Room

	//Fetch rooms from database
	result := config.DB.Find(&rooms)
	if result.Error != nil {
		http.Error(w, "Error fetching rooms", http.StatusInternalServerError)
		return
	}

	//Return the list of rooms in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}