package models

type Room struct {
    RoomID   uint   `gorm:"primaryKey"`
    Name     string `gorm:"size:255;not null"`
    Capacity int    `gorm:"not null"`
}
