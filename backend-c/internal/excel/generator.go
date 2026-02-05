package excel

import (
	"bytes"
	"department-eduvault-backend/models"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

// BuildCertificatesWorkbook creates an Excel file from a slice of certificates.
// It creates sheets based on sections if multiple sections are present, or a single sheet if specified.
func BuildCertificatesWorkbook(certs []models.Certificate, defaultSheetName string) ([]byte, error) {
	f := excelize.NewFile()

	// Map certs by section
	certsBySection := make(map[string][]models.Certificate)
	for _, cert := range certs {
		section := cert.Section
		if section == "" {
			section = "Unknown"
		}
		certsBySection[section] = append(certsBySection[section], cert)
	}

	// Create a sheet for each section
	firstSheet := true
	for section, sectionCerts := range certsBySection {
		sheetName := fmt.Sprintf("Section %s", section)
		if len(sheetName) > 31 {
			sheetName = sheetName[:31] // Excel limit
		}

		index, err := f.NewSheet(sheetName)
		if err != nil {
			return nil, err
		}

		if firstSheet {
			// Remove default Sheet1 if we are adding our own
			f.DeleteSheet("Sheet1")
			f.SetActiveSheet(index)
			firstSheet = false
		}

		// Set headers
		headers := []string{
			"Register Number",
			"Student Name",
			"Section",
			"Drive Link",
			"Uploaded By",
			"Uploaded At",
			"ML Status",
			"ML Score",
			"Faculty Status",
			"Is Legit",
		}

		for i, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, header)
		}

		// Add data
		for i, cert := range sectionCerts {
			row := i + 2
			setCell := func(col int, val interface{}) {
				cell, _ := excelize.CoordinatesToCellName(col, row)
				_ = f.SetCellValue(sheetName, cell, val)
			}

			setCell(1, cert.RegisterNumber)
			setCell(2, cert.StudentName)
			setCell(3, cert.Section)
			setCell(4, cert.DriveLink)
			setCell(5, cert.UploadedBy)
			setCell(6, cert.UploadedAt.Format(time.RFC3339))
			setCell(7, cert.MLStatus)
			if cert.MLScore != nil {
				setCell(8, *cert.MLScore)
			}
			setCell(9, cert.FacultyStatus)
			if cert.IsLegit != nil {
				setCell(10, *cert.IsLegit)
			}
		}

		autoSizeColumns(f, sheetName, len(headers))
	}

	// If no certs, create a placeholder
	if len(certs) == 0 {
		f.NewSheet(defaultSheetName)
		f.DeleteSheet("Sheet1")
		f.SetCellValue(defaultSheetName, "A1", "No records found")
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// autoSizeColumns provides basic width adjustments for readability.
func autoSizeColumns(f *excelize.File, sheet string, columns int) {
	for col := 1; col <= columns; col++ {
		width := 18.0
		if col == 4 {
			width = 40.0 // drive link
		}
		colLetter, _ := excelize.ColumnNumberToName(col)
		_ = f.SetColWidth(sheet, colLetter, colLetter, width)
	}
	// Freeze header row for readability.
	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
		Selection: []excelize.Selection{
			{SQRef: "A2", ActiveCell: "A2", Pane: "bottomLeft"},
		},
	})
}
