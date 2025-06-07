package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyumbanicare/internal/models"
	"gorm.io/gorm"
)

// Prescription handlers
func CreatePrescription(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req struct {
			ImageURL string `json:"image_url" binding:"required"`
			Notes    string `json:"notes"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		prescription := models.Prescription{
			UserID:   userID.(uuid.UUID),
			ImageURL: req.ImageURL,
			Notes:    req.Notes,
			Status:   "uploaded",
		}

		if err := db.Create(&prescription).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create prescription"})
			return
		}

		c.JSON(http.StatusCreated, prescription)
	}
}

func ListPrescriptions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var prescriptions []models.Prescription
		query := db.Preload("Medications")
		
		// Non-admin users can only see their own prescriptions
		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.Find(&prescriptions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prescriptions"})
			return
		}

		c.JSON(http.StatusOK, prescriptions)
	}
}

func UpdatePrescriptionStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Status        string `json:"status" binding:"required"`
			PharmacyNotes string `json:"pharmacy_notes"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var prescription models.Prescription
		if err := db.First(&prescription, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Prescription not found"})
			return
		}

		prescription.Status = req.Status
		if req.PharmacyNotes != "" {
			prescription.PharmacyNotes = req.PharmacyNotes
		}

		if err := db.Save(&prescription).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prescription"})
			return
		}

		c.JSON(http.StatusOK, prescription)
	}
}

// Lab Test handlers
func ListLabTests(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var labTests []models.LabTest
		if err := db.Where("available = ?", true).Find(&labTests).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lab tests"})
			return
		}

		c.JSON(http.StatusOK, labTests)
	}
}

func CreateLabBooking(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req struct {
			LabTestID              uuid.UUID `json:"lab_test_id" binding:"required"`
			BookingDate            string    `json:"booking_date" binding:"required"`
			PreferredTime          string    `json:"preferred_time"`
			SampleCollectionMethod string    `json:"sample_collection_method" binding:"required"`
			Address                string    `json:"address"`
			ContactNumber          string    `json:"contact_number" binding:"required"`
			SpecialInstructions    string    `json:"special_instructions"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse booking date
		bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking date format"})
			return
		}

		// Get lab test details
		var labTest models.LabTest
		if err := db.First(&labTest, "id = ? AND available = ?", req.LabTestID, true).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lab test not found"})
			return
		}

		booking := models.LabBooking{
			UserID:                 userID.(uuid.UUID),
			LabTestID:              req.LabTestID,
			BookingDate:            bookingDate,
			PreferredTime:          req.PreferredTime,
			Status:                 "booked",
			SampleCollectionMethod: req.SampleCollectionMethod,
			Address:                req.Address,
			ContactNumber:          req.ContactNumber,
			SpecialInstructions:    req.SpecialInstructions,
			PaymentStatus:          "pending",
			TotalPrice:             labTest.Price,
		}

		if err := db.Create(&booking).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lab booking"})
			return
		}

		c.JSON(http.StatusCreated, booking)
	}
}

func ListLabBookings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var bookings []models.LabBooking
		query := db.Preload("LabTest")
		
		// Non-admin users can only see their own bookings
		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.Find(&bookings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lab bookings"})
			return
		}

		c.JSON(http.StatusOK, bookings)
	}
}

func UpdateLabBookingStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Status         string `json:"status" binding:"required"`
			TrackingNumber string `json:"tracking_number"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var booking models.LabBooking
		if err := db.First(&booking, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lab booking not found"})
			return
		}

		booking.Status = req.Status
		if req.TrackingNumber != "" {
			booking.TrackingNumber = req.TrackingNumber
		}

		if err := db.Save(&booking).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status"})
			return
		}

		c.JSON(http.StatusOK, booking)
	}
}

// Health Education handlers
func ListHealthArticles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		category := c.Query("category")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var articles []models.HealthArticle
		var total int64

		query := db.Model(&models.HealthArticle{}).Where("published = ?", true)
		if category != "" {
			query = query.Where("category = ?", category)
		}

		query.Count(&total)
		if err := query.Preload("Author").Limit(limit).Offset(offset).Find(&articles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"articles":    articles,
			"total":       total,
			"page":        page,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		})
	}
}

func GetHealthArticle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var article models.HealthArticle
		if err := db.Preload("Author").First(&article, "id = ? AND published = ?", id, true).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}

		// Increment view count
		article.ViewCount++
		db.Save(&article)

		c.JSON(http.StatusOK, article)
	}
}

// Telehealth handlers
func CreateTelehealthSession(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req struct {
			ProviderID   uuid.UUID `json:"provider_id" binding:"required"`
			ProviderType string    `json:"provider_type" binding:"required"`
			SessionType  string    `json:"session_type" binding:"required"`
			ScheduledAt  string    `json:"scheduled_at" binding:"required"`
			Duration     int       `json:"duration" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse scheduled time
		scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled time format"})
			return
		}

		session := models.TelehealthSession{
			PatientID:     userID.(uuid.UUID),
			ProviderID:    req.ProviderID,
			ProviderType:  req.ProviderType,
			SessionType:   req.SessionType,
			ScheduledAt:   scheduledAt,
			Duration:      req.Duration,
			Status:        "scheduled",
			PaymentStatus: "pending",
			Price:         50.0, // Default price, should be configurable
		}

		if err := db.Create(&session).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create telehealth session"})
			return
		}

		c.JSON(http.StatusCreated, session)
	}
}

func ListTelehealthSessions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var sessions []models.TelehealthSession
		query := db.Preload("Provider")
		
		// Users can see sessions where they are patient or provider
		role, _ := c.Get("role")
		if role == "admin" {
			// Admin can see all sessions
		} else if role == "doctor" || role == "nurse" {
			query = query.Where("provider_id = ? OR patient_id = ?", userID, userID)
		} else {
			query = query.Where("patient_id = ?", userID)
		}

		if err := query.Find(&sessions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch telehealth sessions"})
			return
		}

		c.JSON(http.StatusOK, sessions)
	}
}

// AI Symptom Checker handlers
func CreateSymptomCheck(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req struct {
			Symptoms []string `json:"symptoms" binding:"required"`
			Severity string   `json:"severity" binding:"required"`
			Duration string   `json:"duration" binding:"required"`
			Age      int      `json:"age" binding:"required"`
			Gender   string   `json:"gender" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// TODO: Integrate with AI service for symptom analysis
		// For now, we'll create a basic response
		symptomCheck := models.SymptomCheck{
			UserID:           userID.(uuid.UUID),
			Symptoms:         req.Symptoms,
			Severity:         req.Severity,
			Duration:         req.Duration,
			Age:              req.Age,
			Gender:           req.Gender,
			Results:          `{"analysis": "Basic symptom analysis", "conditions": []}`,
			Recommendations:  "Please consult with a healthcare provider for proper diagnosis",
			UrgencyLevel:     "medium",
			FollowUpRequired: true,
		}

		if err := db.Create(&symptomCheck).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create symptom check"})
			return
		}

		c.JSON(http.StatusCreated, symptomCheck)
	}
}

func ListSymptomChecks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var symptomChecks []models.SymptomCheck
		query := db.Model(&models.SymptomCheck{})
		
		// Non-admin users can only see their own symptom checks
		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.Find(&symptomChecks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch symptom checks"})
			return
		}

		c.JSON(http.StatusOK, symptomChecks)
	}
}

// CareSense Analytics handlers
func GenerateCareSenseAnalytics(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req struct {
			AnalysisType string `json:"analysis_type" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// TODO: Implement actual analytics generation
		// For now, we'll create a sample analytics record
		analytics := models.CareSenseAnalytics{
			UserID:          userID.(uuid.UUID),
			AnalysisType:    req.AnalysisType,
			DataSource:      "multiple_sources",
			AnalysisData:    `{"trends": [], "patterns": []}`,
			Insights:        []string{"Your health trends show steady improvement", "Consider regular checkups"},
			Recommendations: []string{"Maintain current lifestyle", "Schedule yearly health screening"},
			Score:           75.5,
			GeneratedAt:     time.Now(),
			ExpiresAt:       time.Now().AddDate(0, 1, 0), // Expires in 1 month
		}

		if err := db.Create(&analytics).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate analytics"})
			return
		}

		c.JSON(http.StatusCreated, analytics)
	}
}

func GetCareSenseAnalytics(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		analysisType := c.Query("type")
		var analytics []models.CareSenseAnalytics
		query := db.Where("user_id = ? AND expires_at > ?", userID, time.Now())
		
		if analysisType != "" {
			query = query.Where("analysis_type = ?", analysisType)
		}

		if err := query.Order("generated_at DESC").Find(&analytics).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch analytics"})
			return
		}

		c.JSON(http.StatusOK, analytics)
	}
}
