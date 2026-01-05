package config

import "errors"

var HodEmails = map[string]bool{
	"hod.cse@citchennai.net": true,
}

var FacultyEmails = map[string]bool{
	"faculty1.cse@citchennai.net": true,
	"faculty2.cse@citchennai.net": true,
}

// GetRoleByEmail returns role if allowed, else error
func GetRoleByEmail(email string) (string, error) {
	if HodEmails[email] {
		return "hod", nil
	}

	if FacultyEmails[email] {
		return "faculty", nil
	}

	return "", errors.New("unauthorized email")
}
