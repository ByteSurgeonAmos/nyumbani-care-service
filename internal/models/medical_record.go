package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MedicalRecord struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID            uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	BloodType         string         `json:"blood_type"`
	Allergies         []string       `gorm:"type:text[]" json:"allergies"`
	ChronicConditions []string       `gorm:"type:text[]" json:"chronic_conditions"`
	FamilyHistory     string         `json:"family_history"`
	Medications       []Medication   `gorm:"foreignKey:MedicalRecordID" json:"medications"`
	TestResults       []TestResult   `gorm:"foreignKey:MedicalRecordID" json:"test_results"`
	Consultations     []Consultation `gorm:"foreignKey:MedicalRecordID" json:"consultations"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

type Medication struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	MedicalRecordID   uuid.UUID      `gorm:"type:uuid;not null" json:"medical_record_id"`
	Name              string         `json:"name"`
	Dosage            string         `json:"dosage"`
	Frequency         string         `json:"frequency"`
	StartDate         time.Time      `json:"start_date"`
	EndDate           *time.Time     `json:"end_date,omitempty"`
	PrescribedBy      string         `json:"prescribed_by"`
	Notes             string         `json:"notes"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

type TestResult struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	MedicalRecordID uuid.UUID      `gorm:"type:uuid;not null" json:"medical_record_id"`
	TestType        string         `json:"test_type"`
	TestDate        time.Time      `json:"test_date"`
	Result          string         `json:"result"`
	Interpretation  string         `json:"interpretation"`
	LabName         string         `json:"lab_name"`
	DoctorNotes     string         `json:"doctor_notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type Consultation struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	MedicalRecordID uuid.UUID      `gorm:"type:uuid;not null" json:"medical_record_id"`
	DoctorID        uuid.UUID      `gorm:"type:uuid;not null" json:"doctor_id"`
	Date            time.Time      `json:"date"`
	Type            string         `json:"type"` // in-person, telehealth
	Diagnosis       string         `json:"diagnosis"`
	Treatment       string         `json:"treatment"`
	Notes           string         `json:"notes"`
	FollowUpDate    *time.Time     `json:"follow_up_date,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *MedicalRecord) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *Medication) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (t *TestResult) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (c *Consultation) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
} 