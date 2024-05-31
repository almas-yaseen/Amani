package models

type CarReceiver struct {
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Year         string `json:"year"`
	Color        string `json:"color"`
	Variant      string `json:"variant"`
	Kms          int    `json:"kms"`
	Ownership    int    `json:"ownership"`
	Transmission string `json:"transmission"`
	ImagePath    string `json:"image_path"`
	RegNo        string `json:"reg_no"`
	Status       string `json:"status"`
}
