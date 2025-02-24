package core

type User struct {
	UserID      string   `json:"userid" example:"5e63bbd1-1f39-41cd-a832-a18496ac4f11"`
	Image       string   `json:"image" example:"https://example.com/image.jpg"`
	Username    string   `json:"username" example:"test_user"`
	Email       string   `json:"email" example:"user@example.com"`
	Tel         string   `json:"tel" example:"06xxxxxxxx"`
	Firstname   string   `json:"firstname" example:"John"`
	Lastname    string   `json:"lastname" example:"Doe"`
	DateOfBirth string   `json:"date_of_birth" example:"2000-01-01"`
	Gender      string   `json:"gender" example:"none"`
	UserPlanID  []string `json:"userplan_id" example:"[\"plan_id\", \"plan_id\"]"`
}

type Plan struct {
	PlanID     string `json:"plan_id" example:"none"`
	TripName   string `json:"TripName" example:"BangkokTrip"`
	RegionID   string `json:"region_id" example:"Central Thailand"`
	ProvinceID string `json:"province_id" example:"Bangkok"`
	StartDate  string `json:"start_date" example:"2025-01-01"`
	StartTime  string `json:"start_time" example:"2025-01-01"`
	EndDate    string `json:"end_date" example:"2025-01-01"`
	EndTime    string `json:"end_time" example:"2025-01-01"`
	Visibility bool   `json:"visibility" example:"true"`
}

type Verification struct {
	Otp   string
	Email string
}

type Admin struct {
	Username string
	Password string
}
