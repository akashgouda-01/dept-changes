package certificates

import (
	"log"

	"department-eduvault-backend/internal/models"
	"department-eduvault-backend/internal/utils"
)

const certFile = "data/certificates.json"

// facultyEmail -> []Certificate
func loadCertificates() map[string][]models.Certificate {
	store := make(map[string][]models.Certificate)

	err := utils.ReadJSON(certFile, &store)
	if err != nil {
		log.Println("ReadJSON error:", err)
	}

	return store
}

func saveCertificates(store map[string][]models.Certificate) {
	err := utils.WriteJSON(certFile, store)
	if err != nil {
		log.Println("WriteJSON error:", err)
	}
}

func AddCertificates(email string, certs []models.Certificate) {
	store := loadCertificates()

	for i := range certs {
		certs[i].Status = "pending"
	}

	store[email] = append(store[email], certs...)
	saveCertificates(store)
}

func GetCertificatesByFaculty(email string) []models.Certificate {
	store := loadCertificates()

	if certs, ok := store[email]; ok {
		return certs
	}

	return []models.Certificate{}
}

func UpdateCertificateStatus(email, regNo, driveLink, status string) {
	store := loadCertificates()
	certs := store[email]

	for i, c := range certs {
		if c.RegisterNumber == regNo && c.DriveLink == driveLink {
			certs[i].Status = status
			break
		}
	}

	store[email] = certs
	saveCertificates(store)
}
