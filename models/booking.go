package models

import "time"

type Booking struct {
	BookingID uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	RoomID    uint      `gorm:"not null"`
	Date      time.Time `gorm:"not null"`
	CreatedAt time.Time
}
