package controllers

import (
	"department-eduvault-backend/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UserConfig struct {
	Role      string
	FacultyID string
	Password  string // Simple hardcoded password for now
}

// Whitelist configuration
var validUsers = map[string]UserConfig{
	"poornimalakshmir@citchennai.net":    {Role: "faculty", FacultyID: "CSE245", Password: "R.POORNIMALAKSHMI_CSE245"},
	"selvajothim@citchennai.net":         {Role: "faculty", FacultyID: "CSE086", Password: "SELVAJOTHIM_CSE086"},
	"nishanthinis@citchennai.net":        {Role: "faculty", FacultyID: "CSE262", Password: "NISHANTHINIS_CSE262"},
	"ayeeshanasreenm.cse@citchennai.net": {Role: "faculty", FacultyID: "CSE345", Password: "M.AYEESHANASREEN_CSE345"},
	"aadhishs.cse2024@citchennai.net":    {Role: "hod", FacultyID: "HOD01", Password: "AADHISHS_HOD01"},
}

func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[LOGIN DEBUG] Bind JSON Error: %v\n", err)
		c.Error(utils.NewValidationError("invalid request body", err))
		return
	}

	normalizedEmail := strings.ToLower(req.Email)
	fmt.Printf("[LOGIN DEBUG] Attempt: Email='%s', Role='%s'\n", normalizedEmail, req.Role)

	userConfig, exists := validUsers[normalizedEmail]

	if !exists {
		fmt.Printf("[LOGIN DEBUG] Email not found in whitelist\n")
		c.Error(utils.NewAuthenticationError("invalid email or password", nil))
		return
	}

	// Password Check
	if req.Password != userConfig.Password {
		fmt.Printf("[LOGIN DEBUG] Password Mismatch. Expected='%s', Got='%s'\n", userConfig.Password, req.Password)
		c.Error(utils.NewAuthenticationError("invalid email or password", nil))
		return
	}

	// Role Check
	if req.Role != userConfig.Role {
		fmt.Printf("[LOGIN DEBUG] Role Mismatch. Expected='%s', Got='%s'\n", userConfig.Role, req.Role)
		c.Error(utils.NewAuthenticationError("role mismatch for this user", nil))
		return
	}

	// Generate JWT
	token, err := utils.GenerateToken(normalizedEmail, userConfig.Role, userConfig.FacultyID)
	// ... rest of function
	if err != nil {
		c.Error(utils.NewInternalError("failed to generate token", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"email":      normalizedEmail,
			"role":       userConfig.Role,
			"faculty_id": userConfig.FacultyID,
		},
	})
}
