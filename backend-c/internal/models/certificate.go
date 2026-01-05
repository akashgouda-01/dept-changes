package models

type Certificate struct {
	RegisterNumber string `json:"register_number"`
	Section        string `json:"section"`
	DriveLink      string `json:"drive_link"`
	Status         string `json:"status"` // pending | legit | not_legit
}
