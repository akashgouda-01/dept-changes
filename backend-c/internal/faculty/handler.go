package faculty

import "net/http"

func FacultyDashboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Faculty Dashboard API"))
}
