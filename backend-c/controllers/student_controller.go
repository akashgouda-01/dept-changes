package controllers

import (
	"department-eduvault-backend/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StudentController struct {
	service services.StudentService
}

func NewStudentController(service services.StudentService) *StudentController {
	return &StudentController{service: service}
}

// GetStudentsBySection returns students for a given section
func (ctrl *StudentController) GetStudentsBySection(c *gin.Context) {
	section := c.Query("section")
	if section == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Section is required"})
		return
	}

	// Ideally, we should also check if the logged-in faculty is assigned to this section.
	// For now, we trust the query or assume the frontend sends valid sections.

	fmt.Printf("Fetching students for section: %s\n", section)
	students, err := ctrl.service.GetStudentsBySection(c.Request.Context(), section)
	if err != nil {
		fmt.Printf("Error fetching students: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}
	fmt.Printf("Found %d students for section %s\n", len(students), section)
	if len(students) > 0 {
		fmt.Printf("Sample student: %+v\n", students[0])
	} else {
		fmt.Println("No students found. Checking DB table mapping...")
	}

	c.JSON(http.StatusOK, students)
}
