package auth

import (
	"encoding/json"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"` // "faculty" | "hod"
}

type LoginResponse struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	role := strings.ToLower(strings.TrimSpace(req.Role))

	if !strings.HasSuffix(email, "@citchennai.net") {
		http.Error(w, "Only citchennai.net emails allowed", http.StatusForbidden)
		return
	}

	if role != "faculty" && role != "hod" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	resp := LoginResponse{
		Email: email,
		Role:  role,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
