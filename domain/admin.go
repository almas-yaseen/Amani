package domain

import "time"

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
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Year         int       `json:"year"`
	Color        string    `json:"color"`
	CarType      string    `json:"car_type"`
	FuelType     string    `json:"fuel_type"`
	Variant      string    `json:"variant"`
	Kms          int       `json:"kms"`
	Ownership    int       `json:"ownership"`
	Bannerimage  string    `json:"bannerimage"`
	Transmission string    `json:"transmission"`
	Images       []Image   `gorm:"foreignKey:CarID" json:"images"`
	RegNo        string    `json:"regno"`
	Status       string    `json:"status"`
	Price        int       `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type Image struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	CarID uint   `json:"car_id"`
	Path  string `json:"path"`
}
