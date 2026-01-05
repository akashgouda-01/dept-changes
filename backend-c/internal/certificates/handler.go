package certificates

import (
	"encoding/json"
	"net/http"

	"department-eduvault-backend/internal/models"
)

// POST /faculty/certificates
func UploadCertificates(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		FacultyEmail string              `json:"faculty_email"`
		Certificates []models.Certificate `json:"certificates"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if len(payload.Certificates) > 10 {
		http.Error(w, "Maximum 10 certificates allowed", http.StatusBadRequest)
		return
	}

	AddCertificates(payload.FacultyEmail, payload.Certificates)
	w.Write([]byte("Certificates uploaded successfully"))
}

// GET /faculty/certificates?email=faculty1.cse@citchennai.net
func GetCertificates(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Faculty email required", http.StatusBadRequest)
		return
	}

	certs := GetCertificatesByFaculty(email)
	json.NewEncoder(w).Encode(certs)
}

// PUT /faculty/certificates/status
func UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		FacultyEmail string `json:"faculty_email"`
		RegisterNumber string `json:"register_number"`
		DriveLink string `json:"drive_link"`
		Status string `json:"status"` // legit | not_legit
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	UpdateCertificateStatus(
		payload.FacultyEmail,
		payload.RegisterNumber,
		payload.DriveLink,
		payload.Status,
	)

	w.Write([]byte("Certificate status updated"))
}
