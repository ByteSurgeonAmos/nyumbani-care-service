package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyumbanicare/internal/models"
	"gorm.io/gorm"
)

func CreateMedicalRecordHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var requestData struct {
			Title        string  `json:"title" binding:"required"`
			Date         string  `json:"date" binding:"required"`
			RecordType   string  `json:"recordType" binding:"required"`
			Notes        string  `json:"notes"`
			TestResultID *string `json:"testResultId"`
			DoctorName   string  `json:"doctor"`
			Symptoms     string  `json:"symptoms"`
		}

		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record data"})
			return
		}

		recordDate, err := time.Parse(time.RFC3339, requestData.Date)
		if err != nil {
			recordDate, err = time.Parse("2006-01-02", requestData.Date)
			if err != nil {
				recordDate = time.Now()
			}
		}

		record := models.MedicalRecord{
			ID:     uuid.New(),
			UserID: userID.(uuid.UUID),
		}

		if requestData.TestResultID != nil && *requestData.TestResultID != "" {
			testResultID, err := uuid.Parse(*requestData.TestResultID)
			if err == nil {
				var testKitResult models.TestKitResult
				if err := db.Where("id = ?", testResultID).First(&testKitResult).Error; err == nil {
					testResult := models.TestResult{
						ID:              uuid.New(),
						MedicalRecordID: record.ID,
						TestType:        "Test Kit Analysis",
						TestDate:        recordDate,
						Result:          testKitResult.Result,
						Interpretation:  testKitResult.Notes,
						LabName:         "Nyumbani Care Analysis",
						DoctorNotes:     requestData.Notes,
					}

					if err := db.Create(&testResult).Error; err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test result record"})
						return
					}
				}
			}
		}

		if err := db.Create(&record).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medical record"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":  "Medical record created successfully",
			"recordId": record.ID,
		})
	}
}

func ListMedicalRecordsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var records []models.MedicalRecord
		if err := db.Where("user_id = ?", userID).Preload("TestResults").Find(&records).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch medical records"})
			return
		}

		c.JSON(http.StatusOK, records)
	}
}

func GetMedicalRecordHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		recordID := c.Param("id")
		var record models.MedicalRecord

		if err := db.Where("id = ? AND user_id = ?", recordID, userID).Preload("TestResults").Preload("Medications").Preload("Consultations").First(&record).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
			return
		}

		c.JSON(http.StatusOK, record)
	}
}

func UpdateMedicalRecordHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		recordID := c.Param("id")
		var record models.MedicalRecord

		if err := db.Where("id = ? AND user_id = ?", recordID, userID).First(&record).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
			return
		}

		var updateData models.MedicalRecord
		if err := c.BindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record data"})
			return
		}

		if updateData.BloodType != "" {
			record.BloodType = updateData.BloodType
		}
		if updateData.Allergies != nil {
			record.Allergies = updateData.Allergies
		}
		if updateData.ChronicConditions != nil {
			record.ChronicConditions = updateData.ChronicConditions
		}
		if updateData.FamilyHistory != "" {
			record.FamilyHistory = updateData.FamilyHistory
		}

		if err := db.Save(&record).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update medical record"})
			return
		}

		c.JSON(http.StatusOK, record)
	}
}
