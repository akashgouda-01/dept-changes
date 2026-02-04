package models

// Student mirrors the students table in the database.
type Student struct {
	ID             uint   `gorm:"primaryKey"`
	RegisterNumber string `gorm:"column:register_number;unique;not null"`
	Name           string `gorm:"column:name;not null"`
	Email          string `gorm:"column:email"`
	Section        string `gorm:"column:section;not null"`
	Semester       int    `gorm:"column:semester"`
	IsPresent      bool   `gorm:"column:is_present;default:true;not null"`
	FacultyEmail   string `gorm:"column:faculty_email;not null"`
}

func (Student) TableName() string {
	return "students"
}
