//go:build legacy_http
// +build legacy_http

package hod

import (
	"department-eduvault-backend/internal/db"
	"errors"
)

type StudentStats struct {
	RegNo    string
	Total    int
	Verified int
	Rejected int
	Pending  int
}

func GetStudentStatsByFaculty(facultyID string) ([]StudentStats, error) {
	dbc := db.GetDB()
	if dbc == nil {
		return nil, errors.New("database not initialized")
	}

	query := `
		SELECT
			reg_no,
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE faculty_status = 'LEGIT') AS verified,
			COUNT(*) FILTER (WHERE faculty_status = 'NOT_LEGIT') AS rejected,
			COUNT(*) FILTER (WHERE faculty_status = 'PENDING') AS pending
		FROM certificates
		WHERE faculty_id = 'FAC01'
		GROUP BY reg_no
		ORDER BY reg_no;
	`

	rows, err := dbc.Query(query, facultyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []StudentStats
	for rows.Next() {
		var s StudentStats
		if err := rows.Scan(
			&s.RegNo,
			&s.Total,
			&s.Verified,
			&s.Rejected,
			&s.Pending,
		); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}
