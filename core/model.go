package core

type User struct {
	Image       string `json:"image" example:"https://example.com/image.jpg"`
	Username    string `json:"username" example:"test_user"`
	Email       string `json:"email" example:"user@example.com"`
	Tel         string `json:"tel" example:"06xxxxxxxx"`
	Firstname   string `json:"firstname" example:"John"`
	Lastname    string `json:"lastname" example:"Doe"`
	DateOfBirth string `json:"date_of_birth" example:"2000-01-01"`
	Gender      string `json:"gender" example:"none"`
}

type Verification struct {
	Otp   string
	Email string
}

type Admin struct {
	Username string
	Password string
}
