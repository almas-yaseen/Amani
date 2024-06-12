package domain

type Car struct {
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Brand        string  `json:"brand"`
	Model        string  `json:"model"`
	Year         string  `json:"year"`
	Color        string  `json:"color"`
	Variant      string  `json:"variant"`
	Kms          int     `json:"kms"`
	Ownership    int     `json:"ownership"`
	Bannerimage  string  `json:"bannerimage"`
	Transmission string  `json:"transmission"`
	Images       []Image `gorm:"foreignKey:CarID" json:"images"`
	RegNo        string  `json:"regno"`
	Status       string  `json:"status"`
	Price        int     `json:"price"`
}
type Image struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	CarID uint   `json:"car_id"`
	Path  string `json:"path"`
}
