package models

import "time"

// SectionStatistics mirrors the section_statistics table in Supabase.
type SectionStatistics struct {
	ID                string    `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Section           string    `gorm:"column:section;type:text;unique;not null"`
	TotalStudents     int       `gorm:"column:total_students;type:int;default:0;not null"`
	PresentStudents   int       `gorm:"column:present_students;type:int;default:0;not null"`
	AbsentStudents    int       `gorm:"column:absent_students;type:int;default:0;not null"`
	TotalCertificates int       `gorm:"column:total_certificates;type:int;default:0;not null"`
	LegitCertificates int       `gorm:"column:legit_certificates;type:int;default:0;not null"`
	UpdatedAt         time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null"`
}

func (SectionStatistics) TableName() string {
	return "section_statistics"
}
