package models

import "time"

// SectionStatistics mirrors the section_statistics table in Supabase.
type SectionStatistics struct {
	ID                   string    `gorm:"-"` // Missing
	Section              string    `gorm:"column:section;type:text;unique;not null"`
	TotalStudents        int       `gorm:"-"` // Missing
	PresentStudents      int       `gorm:"-"` // Missing
	AbsentStudents       int       `gorm:"-"` // Missing
	TotalCertificates    int       `gorm:"column:total_uploaded;type:int;default:0;not null"`
	LegitCertificates    int       `gorm:"column:legit_count;type:int;default:0;not null"`
	NotLegitCertificates int       `gorm:"column:not_legit_count;type:int;default:0;not null"`
	UpdatedAt            time.Time `gorm:"column:last_updated;type:timestamp with time zone;not null"`
}

func (SectionStatistics) TableName() string {
	return "section_statistics"
}
