package models

import "time"

// StudentStatistics mirrors the student_statistics table in Supabase.
type StudentStatistics struct {
	// Primary key is RegisterNumber
	RegisterNumber string `gorm:"column:reg_no;type:text;primaryKey"`
	StudentName    string `gorm:"column:student_name;type:text"`
	Section        string `gorm:"column:section;type:text"`
	// Semester               int       `gorm:"-"` // Not tracked in stats for now
	TotalCertificates      int       `gorm:"column:total_uploaded;type:int;default:0;not null"`
	LegitCertificates      int       `gorm:"column:legit_count;type:int;default:0;not null"`
	NotLegitCertificates   int       `gorm:"column:not_legit_count;type:int;default:0;not null"`
	PendingCertificates    int       `gorm:"column:pending_certificates;type:int;default:0;not null"`
	MlVerifiedCertificates int       `gorm:"column:ml_verified_certificates;type:int;default:0;not null"`
	UpdatedAt              time.Time `gorm:"column:last_updated;type:timestamp with time zone;not null"`
}

func (StudentStatistics) TableName() string {
	return "student_statistics"
}
