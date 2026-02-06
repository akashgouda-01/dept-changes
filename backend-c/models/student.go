package models

// Student mirrors the students table in the database.
type Student struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	RegisterNumber string `gorm:"column:register_number;unique" json:"registerNumber"`
	Name           string `gorm:"column:student_name" json:"name"`
	Email          string `gorm:"column:official_mail_id" json:"email"`
	Section        string `gorm:"column:section" json:"section"`
	Semester       int    `gorm:"column:semester" json:"semester"`
	IsPresent      bool   `gorm:"column:is_present;default:true" json:"isPresent"`
	FacultyEmail   string `gorm:"column:faculty_email" json:"facultyEmail"`
}

func (Student) TableName() string {
	return "STUDENTS_DATA"
}
