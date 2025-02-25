package core

import (
	"github.com/lib/pq"
)

type User struct {
	UserID      string         `gorm:"primaryKey" json:"userid" example:"5e63bbd1-1f39-41cd-a832-a18496ac4f11"`
	Image       string         `json:"image" example:"https://example.com/image.jpg"`
	Username    string         `json:"username" example:"test_user"`
	Email       string         `json:"email" example:"user@example.com"`
	Tel         string         `json:"tel" example:"06xxxxxxxx"`
	Firstname   string         `json:"firstname" example:"John"`
	Lastname    string         `json:"lastname" example:"Doe"`
	DateOfBirth string         `json:"date_of_birth" example:"2000-01-01"`
	Gender      string         `json:"gender" example:"none"`
	UserPlanID  pq.StringArray `gorm:"type:text[]" json:"userplan_id" example:"[\"plan_id\", \"plan_id\"]"`
}
type Location struct {
	PlaceID string `json:"place_id" example:"123"`
}
type Plan struct {
	PlanID       string     `json:"plan_id" example:"none"`
	AuthorEmail  string     `json:"author_email" example:"user@example.com"`
	TripName     string     `json:"TripName" example:"BangkokTrip"`
	RegionID     string     `json:"region_id" example:"Central Thailand"`
	ProvinceID   string     `json:"province_id" example:"Bangkok"`
	StartDate    string     `json:"start_date" example:"2025-01-01"`
	StartTime    string     `json:"start_time" example:"16:53:44.581Z"`
	EndDate      string     `json:"end_date" example:"2025-01-01"`
	EndTime      string     `json:"end_time" example:"16:53:44.581Z"`
	TripLocation []Location `json:"trip_location" example:"[\"place_id\", \"place_id\"]"`
	Visibility   bool       `json:"visibility" example:"true"`
}

type Verification struct {
	Otp   string
	Email string
}

type Admin struct {
	Username string
	Password string
}
