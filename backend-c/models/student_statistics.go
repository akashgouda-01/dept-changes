package models

import "time"

// StudentStatistics mirrors the student_statistics table in Supabase.
type StudentStatistics struct {
	ID                     string    `gorm:"-"` // Missing
	StudentID              string    `gorm:"-"` // Missing
	RegisterNumber         string    `gorm:"column:reg_no;type:text;not null"`
	StudentName            string    `gorm:"-"` // Missing
	Section                string    `gorm:"-"` // Missing
	Semester               int       `gorm:"-"` // Missing
	TotalCertificates      int       `gorm:"column:total_uploaded;type:int;default:0;not null"`
	LegitCertificates      int       `gorm:"column:legit_count;type:int;default:0;not null"`
	NotLegitCertificates   int       `gorm:"column:not_legit_count;type:int;default:0;not null"`
	PendingCertificates    int       `gorm:"-"` // Missing
	MlVerifiedCertificates int       `gorm:"-"` // Missing
	UpdatedAt              time.Time `gorm:"column:last_updated;type:timestamp with time zone;not null"`
}

func (StudentStatistics) TableName() string {
	return "student_statistics"
}
