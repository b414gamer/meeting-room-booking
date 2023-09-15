package main

import (
	"log"
	"meeting-room-booking/config"
	"meeting-room-booking/models"
)

func main() {
	config.InitDB()
	sqlDB, err := config.DB.DB()
	if err != nil {
		log.Fatalf("Error getting underlying SQL DB: %v", err)
	}
	defer sqlDB.Close()

	// AutoMigrate will create the tables and keep them updated with the model
	config.DB.AutoMigrate(&models.User{}, &models.Room{}, &models.Booking{})

	// Set up routes, middleware, and start the server here...
}
