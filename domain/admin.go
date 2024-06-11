package domain

type Car struct {
	ID           uint `gorm:"primaryKey;autoIncrement"`
	Brand        string
	Model        string
	Year         string
	Color        string
	Variant      string
	Kms          int
	Ownership    int
	Bannerimage  string
	Transmission string
	Images       []Image `gorm:"foreignKey:CarID"`
	RegNo        string
	Status       string
	Price        int
}

type Image struct {
	ID    uint `gorm:"primaryKey:autoIncrement"`
	CarID uint `gorm:"index"`
	Path  string
}
