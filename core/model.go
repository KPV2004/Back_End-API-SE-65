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

type TripLocation struct {
	PlaceID        string `json:"place_id" example:"place_id_1"`
	PlaceLabel     string `json:"place_label" example:"Place1"`
	CategorieLabel string `json:"categorie_label" example:"Shopping"`
	Introduction   string `json:"introduction" example:"intro to Place1"`
	ThumbnailURL   string `json:"thumbnail_url" example:"https://example.com/image.jpg"`
	Latitude       string `json:"latitude" example:"11.111"`
	Longtitude     string `json:"longtitude" example:"11.111"`
	TimeLocation   string `json:"time_location" example:"08:00AM"`
	Day            string `json:"day" example:"1"`
}

type Plan struct {
	PlanID        string         `gorm:"primaryKey" json:"plan_id" example:"4e63bbd1-1f39-41cd-a832-a18496ac4f11"`
	AuthorEmail   string         `json:"author_email" example:"user@example.com"`
	AuthorName    string         `json:"author_name" example:"author_name"`
	AuthorImg     string         `json:"author_img" example:"https://example.com/image.jpg"`
	TripName      string         `json:"trip_name" example:"BangkokTrip"`
	RegionLabel   string         `json:"region_label" example:"Central Thailand"`
	ProvinceLabel string         `json:"province_label" example:"Bangkok"`
	ProvinceID    string         `json:"province_id" example:"123"`
	StartDate     string         `json:"start_date" example:"2025-01-01"`
	StartTime     string         `json:"start_time" example:"16:53:44.581Z"`
	EndDate       string         `json:"end_date" example:"2025-01-01"`
	EndTime       string         `json:"end_time" example:"16:53:44.581Z"`
	Description   string         `json:"description" example:"1 Day Trip Description"`
	TripLocation  []TripLocation `json:"trip_location" gorm:"type:jsonb;serializer:json" example:"[{\"place_id\":\"place_id_1\", \"time_location\":\"08:00AM\", \"day\":\"Monday\"}, {\"place_id\":\"place_id_2\", \"time_location\":\"10:00AM\", \"day\":\"Monday\"}]"`
	Visibility    bool           `json:"visibility" example:"true"`
}

type Verification struct {
	Otp   string
	Email string
}

type Admin struct {
	Username string
	Password string
}
