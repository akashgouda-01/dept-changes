package routes

import (
	"net/http"

	"department-eduvault-backend/internal/auth"
	"department-eduvault-backend/internal/certificates"
	"department-eduvault-backend/internal/faculty"
	"department-eduvault-backend/internal/hod"
	"department-eduvault-backend/internal/students"
)

// withCORS adds minimal CORS headers so the React frontend can talk to this API.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Login
	mux.HandleFunc("/login", auth.LoginHandler)

	// Dashboards
	mux.HandleFunc("/faculty/dashboard", faculty.FacultyDashboard)
	mux.HandleFunc("/hod/dashboard", hod.HodDashboard)

	// Student management
	mux.HandleFunc("/faculty/students", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			students.GetStudents(w, r)
		case http.MethodPost:
			students.AddStudentHandler(w, r)
		case http.MethodDelete:
			students.RemoveStudentHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Certificate management
	mux.HandleFunc("/faculty/certificates", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			certificates.UploadCertificates(w, r)
		case http.MethodGet:
			certificates.GetCertificates(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/faculty/certificates/status", certificates.UpdateStatus)

	return withCORS(mux)
}
