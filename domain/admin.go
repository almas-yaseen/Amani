package domain

import (
	"time"
)

const (
	CarStatusSold      = "Sold"
	CarStatusAvailable = "Available"
)

const (

	// cartype choice

	CarTypeSedan     = "sedan"
	CarTypeHatchback = "hatchback"
	CarTypeSuv       = "suv"
	CarTypeBike      = "bike"
	//  fuel type  choices
	FuelTypePetrol   = "petrol"
	FuelTypeDiesel   = "diesel"
	FuelTypeCNG      = "cng"
	FuelTypeElectric = "electric"
)

type Car struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Model          string    `json:"model"`
	BrandID        uint      `json:"brand_id"`
	Brand          Brand     `gorm:"foriegnKey:BrandID" json:"brand"`
	Year           int       `json:"year"`
	Color          string    `json:"color"`
	CarType        string    `json:"car_type"`
	FuelType       string    `json:"fuel_type"`
	Variant        string    `json:"variant"`
	Kms            int       `json:"kms"`
	Ownership      int       `json:"ownership"`
	Bannerimage    string    `json:"bannerimage"`
	Transmission   string    `json:"transmission"`
	Images         []Image   `gorm:"foreignKey:CarID" json:"images"`
	RegNo          string    `json:"regno"`
	Status         string    `json:"status"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Engine_size    string    `json:"engine_size"`      //new ones
	Insurance_date string    `json:"insurance_dating"` //new ones
	Location       string    `json:"location"`         //new ones
}

type Brand struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Image struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	CarID uint   `json:"car_id"`
	Path  string `json:"path"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type YoutubeLink struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        uint      `json:"id"`
	VideoLink string    `json:"video_link"`
}
