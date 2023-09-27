package main

import (
	"log"
	"meeting-room-booking/config"
	"meeting-room-booking/handlers"
	"meeting-room-booking/models"
	"meeting-room-booking/seed"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Set timezone to UTC+7
	loc := time.FixedZone("UTC+7", 7*60*60)
	time.Local = loc

	config.InitDB()
	sqlDB, err := config.DB.DB()
	if err != nil {
		log.Fatalf("Error getting underlying SQL DB: %v", err)
	}
	defer sqlDB.Close()

	// AutoMigrate will create the tables and keep them updated with the model
	config.DB.AutoMigrate(&models.User{}, &models.Room{}, &models.Booking{})
	config.LoadEnv()

	// Insert mock rooms
	seed.InsertMockRooms()

	// Set up routes, middleware, and start the server here...

	r := mux.NewRouter()
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/rooms", handlers.ListRoomsHandler).Methods("GET")
	r.HandleFunc("/bookings", handlers.BookRoomHandler).Methods("POST")
	r.HandleFunc("/bookings", handlers.ListBookingsHandler).Methods("GET")

	// http.Handle("/", r)
	// http.ListenAndServe(":80", nil)

	log.Println("Starting server on port 80...")
	err = http.ListenAndServe(":80", r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// http.HandleFunc("/login", handlers.LoginHandler)
	// http.HandleFunc("/register", handlers.RegisterHandler)

	// log.Println("Starting server on port 80...")
	// err = http.ListenAndServe(":80", nil)
	// if err != nil {
	// 	log.Fatalf("Error starting server: %v", err)
	// }
}

func InsertMockRooms() {
	panic("unimplemented")
}
