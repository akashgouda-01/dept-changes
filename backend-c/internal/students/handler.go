package students

import (
	"encoding/json"
	"net/http"

	"department-eduvault-backend/internal/models"
)

// GET /faculty/students?email=faculty1.cse@citchennai.net
func GetStudents(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Faculty email required", http.StatusBadRequest)
		return
	}

	students := GetStudentsByFaculty(email)
	json.NewEncoder(w).Encode(students)
}

// POST /faculty/students
func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		FacultyEmail string         `json:"faculty_email"`
		Student      models.Student `json:"student"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	AddStudent(payload.FacultyEmail, payload.Student)
	w.Write([]byte("Student added successfully"))
}

// DELETE /faculty/students
func RemoveStudentHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		FacultyEmail  string `json:"faculty_email"`
		RegisterNumber string `json:"register_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	RemoveStudent(payload.FacultyEmail, payload.RegisterNumber)
	w.Write([]byte("Student removed successfully"))
}
