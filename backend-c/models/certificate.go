package models

import "time"

// MLStatus captures machine-learning verification states mapped to the Postgres enum ml_status.
type MLStatus string

const (
	MLStatusPending   MLStatus = "PENDING"
	MLStatusVerified  MLStatus = "VERIFIED"
	MLStatusDuplicate MLStatus = "DUPLICATE"
)

// FacultyStatus captures manual verification states mapped to the Postgres enum faculty_status.
type FacultyStatus string

const (
	FacultyStatusPending  FacultyStatus = "PENDING"
	FacultyStatusLegit    FacultyStatus = "LEGIT"
	FacultyStatusNotLegit FacultyStatus = "NOT_LEGIT"
)

// Certificate mirrors the existing certificates table in Supabase.
type Certificate struct {
	ID             string        `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	DriveLink      string        `gorm:"column:drive_link;type:text;not null"`
	RegisterNumber string        `gorm:"column:reg_no;type:text;not null"`
	Section        string        `gorm:"column:section;type:text;not null"`
	StudentName    string        `gorm:"column:student_name;type:text;not null"`
	UploadedBy     string        `gorm:"column:faculty_id;type:text;not null"`
	UploadedAt     time.Time     `gorm:"column:uploaded_at;type:timestamp with time zone;not null"`
	UpdatedAt      time.Time     `gorm:"column:updated_at;type:timestamp with time zone;default:now()"`
	MLStatus       MLStatus      `gorm:"column:ml_status;type:ml_status_enum;default:'PENDING';not null"`
	FacultyStatus  FacultyStatus `gorm:"column:faculty_status;type:faculty_status_enum;default:'PENDING';not null"`
	// Updated to map to DB columns
	IsLegit  *bool    `gorm:"column:is_legit;type:boolean"`
	MLScore  *float64 `gorm:"column:ml_score;type:double precision"`
	Archived bool     `gorm:"column:archived;type:boolean;default:false;not null"`
}

func (Certificate) TableName() string {
	return "faculty_certificates"
}
