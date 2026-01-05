package models

import "time"

// StudentStatistics mirrors the student_statistics table in Supabase.
type StudentStatistics struct {
	ID                     string    `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	StudentID              string    `gorm:"column:student_id;type:uuid;not null"`
	RegisterNumber         string    `gorm:"column:register_number;type:text;not null"`
	StudentName            string    `gorm:"column:student_name;type:text;not null"`
	Section                string    `gorm:"column:section;type:text;not null"`
	Semester               int       `gorm:"column:semester;type:int;not null"`
	TotalCertificates      int       `gorm:"column:total_certificates;type:int;default:0;not null"`
	LegitCertificates      int       `gorm:"column:legit_certificates;type:int;default:0;not null"`
	NotLegitCertificates   int       `gorm:"column:not_legit_certificates;type:int;default:0;not null"`
	PendingCertificates    int       `gorm:"column:pending_certificates;type:int;default:0;not null"`
	MlVerifiedCertificates int       `gorm:"column:ml_verified_certificates;type:int;default:0;not null"`
	UpdatedAt              time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null"`
}

func (StudentStatistics) TableName() string {
	return "student_statistics"
}
