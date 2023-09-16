package models

import (
	"errors"
	"time"
)

type User struct {
	UserID    uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;unique;not null"`
	Password  string `gorm:"size:255;not null"`
	CreatedAt time.Time
}

func ValidateUserCredentials(username, password string) (*User, error) {
	// Placeholder: Check the username and password against your database
	//code here:
	if username == "testuser" && password == "testpassword" {
		// For now, let's assume any username and password combination is valid
		return &User{Name: username}, nil
	}
	return nil, errors.New("invalid username or password")
}
