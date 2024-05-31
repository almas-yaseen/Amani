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
	Transmission string
	ImagePath1   string
	ImagePath2   string
	ImagePath3   string
	ImagePath4   string
	ImagePath5   string
	ImagePath6   string
	ImagePath7   string

	RegNo  string
	Status string
}
