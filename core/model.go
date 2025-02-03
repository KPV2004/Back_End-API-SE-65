package core

import "gorm.io/gorm"

type User struct {
	gorm.Model
	// ID          uint
	Username    string
	Password    string
	Email       string
	Tel         string
	Firstname   string
	Lastname    string
	DateOfBirth string
	Gender      string
	isVerified  bool
}

type Verification struct {
	Otp   string
	Email string
}
