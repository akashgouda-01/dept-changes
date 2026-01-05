package models

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"` // faculty | hod
	StaffID  string `json:"staff_id"`
}
 