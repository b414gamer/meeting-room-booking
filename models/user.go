package models

import "time"

type User struct {
	UserID    uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;unique;not null"`
	Password  string `gorm:"size:255;not null"`
	CreatedAt time.Time
}
