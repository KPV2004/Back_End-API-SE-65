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
}

type Verification struct {
	Otp   string
	Email string
}

type Admin struct {
	gorm.Model
	Username string
	Password string
}
