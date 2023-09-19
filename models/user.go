package models

import (
	"encoding/base64"
	"errors"
	"meeting-room-booking/config"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

type User struct {
	UserID    uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;unique;not null"`
	Password  string `gorm:"size:255;not null"`
	CreatedAt time.Time
}

// func ValidateUserCredentials(email, password string) (*User, error) {
// 	// Placeholder: Check the username and password against your database
// 	//code here:
// 	if email == "testuser" && password == "testpassword" {
// 		// For now, let's assume any username and password combination is valid
// 		return &User{Email: email}, nil
// 	}
// 	return nil, errors.New("invalid email or password")
// }

func ValidateUserCredentials(email, password string) (*User, error) {
	// var user User

	// //Query the database for a user with the given email
	// result := config.DB.Where("email = ?", email).First(&user)

	// //Check if a user with the given email was found and the password matches
	// if result.Error == nil && user.Password == password {
	// 	return &user, nil
	// }

	// return nil, errors.New("invalid email or password")

	var user User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	saltAndHash := strings.Split(user.Password, "$")
	if len(saltAndHash) != 2 {
		return nil, errors.New("stored password is not in the correct format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(saltAndHash[0])
	if err != nil {
		return nil, err
	}

	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	if saltAndHash[1] != base64.RawStdEncoding.EncodeToString(hashedPassword) {
		return nil, errors.New("invalid password")
	}
	return &user, nil
}
