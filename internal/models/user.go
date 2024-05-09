package models

import "github.com/google/uuid"

type User struct {
	UserId      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
}

type OTPObject struct {
	Name        string `json:"name" db:"name"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	OTPPassword string `json:"otp_password" db:"otp_password"`
}
