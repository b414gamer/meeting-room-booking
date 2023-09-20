package models

import "time"

type Booking struct {
	BookingID uint      `gorm:"primaryKey"`
	RoomID    uint      `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	CreatedAt time.Time
}
