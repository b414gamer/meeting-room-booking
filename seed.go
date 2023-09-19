// Description: Seed data for database
package main

import (
	"meeting-room-booking/config"
	"meeting-room-booking/models"
)

func InsertMockRooms() {
	var count int64
	config.DB.Model(&models.Room{}).Count(&count)

	if count == 0 {
		rooms := []models.Room{
			{Name: "Room A", Capacity: 10},
			{Name: "Room B", Capacity: 15},
			{Name: "Room C", Capacity: 20},
		}

		for _, room := range rooms {
			config.DB.Create(&room)
		}
	}
}