package excel

import (
	"bytes"
	"fmt"
	"time"

	"department-eduvault-backend/models"
	"github.com/xuri/excelize/v2"
)

// BuildCertificatesWorkbook renders certificates into a single-sheet XLSX file.
func BuildCertificatesWorkbook(certs []models.Certificate, sheetName string) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Certificates"
	if sheetName != "" {
		sheet = sheetName
	}
	f.SetSheetName(f.GetSheetName(0), sheet)

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

	// Header row
	for idx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(idx+1, 1)
		_ = f.SetCellValue(sheet, cell, header)
	}

	// Data rows
	for i, cert := range certs {
		row := i + 2 // data starts at row 2
		setCell := func(col int, val interface{}) {
			cell, _ := excelize.CoordinatesToCellName(col, row)
			_ = f.SetCellValue(sheet, cell, val)
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
		} else {
			setCell(8, "")
		}
		setCell(9, cert.FacultyStatus)
		if cert.IsLegit != nil {
			setCell(10, *cert.IsLegit)
		} else {
			setCell(10, "")
		}
	}

	autoSizeColumns(f, sheet, len(headers))

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("write workbook: %w", err)
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