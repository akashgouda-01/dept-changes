package models

type Student struct {
	RegisterNumber string `json:"register_number"`
	Name           string `json:"name"`
	Section        string `json:"section"`
}
