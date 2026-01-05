package hod

import "net/http"

func HodDashboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOD Dashboard API"))
}
