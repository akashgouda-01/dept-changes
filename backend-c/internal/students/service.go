package students

import (
	"log"

	"department-eduvault-backend/internal/models"
	"department-eduvault-backend/internal/utils"
)

const studentFile = "data/students.json"

// facultyEmail -> []Student
func loadStudents() map[string][]models.Student {
	store := make(map[string][]models.Student)

	err := utils.ReadJSON(studentFile, &store)
	if err != nil {
		log.Println("ReadJSON error:", err)
	}

	return store
}

func saveStudents(store map[string][]models.Student) {
	err := utils.WriteJSON(studentFile, store)
	if err != nil {
		log.Println("WriteJSON error:", err)
	}
}

func GetStudentsByFaculty(email string) []models.Student {
	store := loadStudents()

	if students, exists := store[email]; exists {
		return students
	}

	return []models.Student{} // return empty array instead of nil
}

func AddStudent(email string, student models.Student) {
	store := loadStudents()
	store[email] = append(store[email], student)
	saveStudents(store)
}

func RemoveStudent(email, regNo string) {
	store := loadStudents()
	students := store[email]

	filtered := []models.Student{}
	for _, s := range students {
		if s.RegisterNumber != regNo {
			filtered = append(filtered, s)
		}
	}

	store[email] = filtered
	saveStudents(store)
}
