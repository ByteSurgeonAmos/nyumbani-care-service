package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyumbanicare/internal/models"
	"gorm.io/gorm"
)

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

		bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking date format"})
			return
		}

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

		article.ViewCount++
		db.Save(&article)

		c.JSON(http.StatusOK, article)
	}
}

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
			Price:         50.0,
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

		role, _ := c.Get("role")
		if role == "admin" {
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

func GenerateCareSenseAnalytics(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req struct {
			AnalysisType string `json:"analysis_type" binding:"required"`
			TimeRange    string `json:"time_range,omitempty"` // 1m, 3m, 6m, 1y (defaults to 3m)
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		timeRange := req.TimeRange
		if timeRange == "" {
			timeRange = "3m"
		}

		var fromDate time.Time
		switch timeRange {
		case "1m":
			fromDate = time.Now().AddDate(0, -1, 0)
		case "3m":
			fromDate = time.Now().AddDate(0, -3, 0)
		case "6m":
			fromDate = time.Now().AddDate(0, -6, 0)
		case "1y":
			fromDate = time.Now().AddDate(-1, 0, 0)
		default:
			fromDate = time.Now().AddDate(0, -3, 0) // Default to 3 months
		}

		// Generate analytics based on analysis type
		var analytics models.CareSenseAnalytics
		var err error

		switch req.AnalysisType {
		case "health_trends":
			analytics, err = generateHealthTrendsAnalytics(db, userID.(uuid.UUID), fromDate, timeRange)
		case "risk_assessment":
			analytics, err = generateRiskAssessmentAnalytics(db, userID.(uuid.UUID), fromDate, timeRange)
		case "wellness_score":
			analytics, err = generateWellnessScoreAnalytics(db, userID.(uuid.UUID), fromDate, timeRange)
		case "comprehensive":
			analytics, err = generateComprehensiveAnalytics(db, userID.(uuid.UUID), fromDate, timeRange)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid analysis type. Supported types: health_trends, risk_assessment, wellness_score, comprehensive"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate analytics: " + err.Error()})
			return
		}

		if err := db.Create(&analytics).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save analytics"})
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

type AnalyticsData struct {
	HealthTrends        []HealthTrend       `json:"health_trends"`
	RiskFactors         []RiskFactor        `json:"risk_factors"`
	TestResults         []TestSummary       `json:"test_results"`
	SymptomPatterns     []SymptomPattern    `json:"symptom_patterns"`
	ConsultationStats   ConsultationStats   `json:"consultation_stats"`
	LabResults          []LabSummary        `json:"lab_results"`
	PrescriptionHistory PrescriptionHistory `json:"prescription_history"`
}

type HealthTrend struct {
	Metric    string    `json:"metric"`
	Values    []float64 `json:"values"`
	Dates     []string  `json:"dates"`
	Direction string    `json:"direction"` // improving, declining, stable
}

type RiskFactor struct {
	Category    string  `json:"category"`
	Level       string  `json:"level"` // low, medium, high
	Score       float64 `json:"score"`
	Description string  `json:"description"`
}

type TestSummary struct {
	TestType   string    `json:"test_type"`
	Result     string    `json:"result"`
	Confidence float64   `json:"confidence"`
	Date       time.Time `json:"date"`
	Status     string    `json:"status"`
}

type SymptomPattern struct {
	Symptoms     []string `json:"symptoms"`
	Frequency    int      `json:"frequency"`
	UrgencyLevel string   `json:"urgency_level"`
	Pattern      string   `json:"pattern"` // recurring, one-time, chronic
}

type ConsultationStats struct {
	TotalSessions     int       `json:"total_sessions"`
	CompletedSessions int       `json:"completed_sessions"`
	AverageRating     float64   `json:"average_rating"`
	LastConsultation  time.Time `json:"last_consultation"`
}

type LabSummary struct {
	TestName   string    `json:"test_name"`
	Status     string    `json:"status"`
	BookedDate time.Time `json:"booked_date"`
	Category   string    `json:"category"`
}

type PrescriptionHistory struct {
	TotalPrescriptions int      `json:"total_prescriptions"`
	ActiveMedications  int      `json:"active_medications"`
	CommonMedications  []string `json:"common_medications"`
	ComplianceRate     float64  `json:"compliance_rate"`
}

func generateHealthTrendsAnalytics(db *gorm.DB, userID uuid.UUID, fromDate time.Time, timeRange string) (models.CareSenseAnalytics, error) {
	var data AnalyticsData
	var insights []string
	var recommendations []string

	var testResults []models.TestKitResult
	if err := db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&testResults).Error; err != nil {
		return models.CareSenseAnalytics{}, fmt.Errorf("failed to fetch test results: %v", err)
	}

	var symptomChecks []models.SymptomCheck
	if err := db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&symptomChecks).Error; err != nil {
		return models.CareSenseAnalytics{}, fmt.Errorf("failed to fetch symptom checks: %v", err)
	}

	var telehealthSessions []models.TelehealthSession
	if err := db.Where("patient_id = ? AND created_at >= ?", userID, fromDate).Find(&telehealthSessions).Error; err != nil {
		return models.CareSenseAnalytics{}, fmt.Errorf("failed to fetch telehealth sessions: %v", err)
	}

	data.HealthTrends = generateHealthTrendsFromData(testResults, symptomChecks)
	data.TestResults = convertTestResults(testResults)
	data.SymptomPatterns = analyzeSymptomPatterns(symptomChecks)
	data.ConsultationStats = analyzeConsultationStats(telehealthSessions)

	insights = generateTrendInsights(data, timeRange)
	recommendations = generateTrendRecommendations(data)

	score := calculateTrendScore(data)

	analysisDataJSON, _ := json.Marshal(data)

	return models.CareSenseAnalytics{
		UserID:          userID,
		AnalysisType:    "health_trends",
		DataSource:      "test_results,symptoms,consultations",
		AnalysisData:    string(analysisDataJSON),
		Insights:        insights,
		Recommendations: recommendations,
		Score:           score,
		GeneratedAt:     time.Now(),
		ExpiresAt:       time.Now().AddDate(0, 1, 0),
	}, nil
}

func generateRiskAssessmentAnalytics(db *gorm.DB, userID uuid.UUID, fromDate time.Time, timeRange string) (models.CareSenseAnalytics, error) {
	var data AnalyticsData
	var insights []string
	var recommendations []string

	var testResults []models.TestKitResult
	var symptomChecks []models.SymptomCheck
	var prescriptions []models.Prescription
	var labBookings []models.LabBooking

	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&testResults)
	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&symptomChecks)
	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&prescriptions)
	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&labBookings)

	data.RiskFactors = assessRiskFactors(testResults, symptomChecks, prescriptions)
	data.TestResults = convertTestResults(testResults)
	data.SymptomPatterns = analyzeSymptomPatterns(symptomChecks)
	data.LabResults = convertLabBookings(labBookings)

	insights = generateRiskInsights(data, timeRange)
	recommendations = generateRiskRecommendations(data)

	score := calculateRiskScore(data)

	analysisDataJSON, _ := json.Marshal(data)

	return models.CareSenseAnalytics{
		UserID:          userID,
		AnalysisType:    "risk_assessment",
		DataSource:      "test_results,symptoms,prescriptions,lab_tests",
		AnalysisData:    string(analysisDataJSON),
		Insights:        insights,
		Recommendations: recommendations,
		Score:           score,
		GeneratedAt:     time.Now(),
		ExpiresAt:       time.Now().AddDate(0, 1, 0),
	}, nil
}

func generateWellnessScoreAnalytics(db *gorm.DB, userID uuid.UUID, fromDate time.Time, timeRange string) (models.CareSenseAnalytics, error) {
	var data AnalyticsData
	var insights []string
	var recommendations []string

	var testResults []models.TestKitResult
	var symptomChecks []models.SymptomCheck
	var telehealthSessions []models.TelehealthSession
	var prescriptions []models.Prescription
	var labBookings []models.LabBooking

	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&testResults)
	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&symptomChecks)
	db.Where("patient_id = ? AND created_at >= ?", userID, fromDate).Find(&telehealthSessions)
	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&prescriptions)
	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&labBookings)

	data.TestResults = convertTestResults(testResults)
	data.SymptomPatterns = analyzeSymptomPatterns(symptomChecks)
	data.ConsultationStats = analyzeConsultationStats(telehealthSessions)
	data.PrescriptionHistory = analyzePrescriptionHistory(prescriptions)
	data.LabResults = convertLabBookings(labBookings)

	insights = generateWellnessInsights(data, timeRange)
	recommendations = generateWellnessRecommendations(data)

	score := calculateWellnessScore(data)

	analysisDataJSON, _ := json.Marshal(data)

	return models.CareSenseAnalytics{
		UserID:          userID,
		AnalysisType:    "wellness_score",
		DataSource:      "comprehensive_health_data",
		AnalysisData:    string(analysisDataJSON),
		Insights:        insights,
		Recommendations: recommendations,
		Score:           score,
		GeneratedAt:     time.Now(),
		ExpiresAt:       time.Now().AddDate(0, 1, 0),
	}, nil
}

func generateComprehensiveAnalytics(db *gorm.DB, userID uuid.UUID, fromDate time.Time, timeRange string) (models.CareSenseAnalytics, error) {
	var data AnalyticsData
	var insights []string
	var recommendations []string

	var testResults []models.TestKitResult
	var symptomChecks []models.SymptomCheck
	var telehealthSessions []models.TelehealthSession
	var prescriptions []models.Prescription
	var labBookings []models.LabBooking

	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&testResults)
	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&symptomChecks)
	db.Where("patient_id = ? AND created_at >= ?", userID, fromDate).Find(&telehealthSessions)
	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&prescriptions)
	db.Where("user_id = ? AND created_at >= ?", userID, fromDate).Find(&labBookings)

	data.HealthTrends = generateHealthTrendsFromData(testResults, symptomChecks)
	data.RiskFactors = assessRiskFactors(testResults, symptomChecks, prescriptions)
	data.TestResults = convertTestResults(testResults)
	data.SymptomPatterns = analyzeSymptomPatterns(symptomChecks)
	data.ConsultationStats = analyzeConsultationStats(telehealthSessions)
	data.PrescriptionHistory = analyzePrescriptionHistory(prescriptions)
	data.LabResults = convertLabBookings(labBookings)

	insights = generateComprehensiveInsights(data, timeRange)
	recommendations = generateComprehensiveRecommendations(data)

	score := calculateComprehensiveScore(data)

	analysisDataJSON, _ := json.Marshal(data)

	return models.CareSenseAnalytics{
		UserID:          userID,
		AnalysisType:    "comprehensive",
		DataSource:      "all_health_data_sources",
		AnalysisData:    string(analysisDataJSON),
		Insights:        insights,
		Recommendations: recommendations,
		Score:           score,
		GeneratedAt:     time.Now(),
		ExpiresAt:       time.Now().AddDate(0, 1, 0),
	}, nil
}

func generateHealthTrendsFromData(testResults []models.TestKitResult, symptomChecks []models.SymptomCheck) []HealthTrend {
	var trends []HealthTrend

	if len(testResults) > 0 {
		confidenceTrend := HealthTrend{
			Metric: "test_confidence",
			Values: make([]float64, len(testResults)),
			Dates:  make([]string, len(testResults)),
		}

		for i, result := range testResults {
			confidenceTrend.Values[i] = result.AIConfidence
			confidenceTrend.Dates[i] = result.CreatedAt.Format("2006-01-02")
		}

		confidenceTrend.Direction = calculateTrendDirection(confidenceTrend.Values)
		trends = append(trends, confidenceTrend)
	}

	if len(symptomChecks) > 0 {
		urgencyTrend := HealthTrend{
			Metric: "symptom_urgency",
			Values: make([]float64, len(symptomChecks)),
			Dates:  make([]string, len(symptomChecks)),
		}

		for i, check := range symptomChecks {
			urgencyTrend.Values[i] = urgencyToScore(check.UrgencyLevel)
			urgencyTrend.Dates[i] = check.CreatedAt.Format("2006-01-02")
		}

		urgencyTrend.Direction = calculateTrendDirection(urgencyTrend.Values)
		trends = append(trends, urgencyTrend)
	}

	return trends
}

func assessRiskFactors(testResults []models.TestKitResult, symptomChecks []models.SymptomCheck, prescriptions []models.Prescription) []RiskFactor {
	var riskFactors []RiskFactor

	highUrgencyCount := 0
	for _, check := range symptomChecks {
		if check.UrgencyLevel == "high" || check.UrgencyLevel == "emergency" {
			highUrgencyCount++
		}
	}

	if highUrgencyCount > 0 {
		level := "medium"
		score := float64(highUrgencyCount) * 20.0
		if score > 60 {
			level = "high"
		}

		riskFactors = append(riskFactors, RiskFactor{
			Category:    "symptom_urgency",
			Level:       level,
			Score:       math.Min(score, 100),
			Description: fmt.Sprintf("Reported %d high urgency symptoms recently", highUrgencyCount),
		})
	}

	positiveResults := 0
	for _, result := range testResults {
		if result.Result == "positive" {
			positiveResults++
		}
	}

	if positiveResults > 0 {
		score := float64(positiveResults) * 30.0
		level := "medium"
		if score > 50 {
			level = "high"
		}

		riskFactors = append(riskFactors, RiskFactor{
			Category:    "test_results",
			Level:       level,
			Score:       math.Min(score, 100),
			Description: fmt.Sprintf("Has %d positive test results requiring attention", positiveResults),
		})
	}

	if len(prescriptions) > 3 {
		riskFactors = append(riskFactors, RiskFactor{
			Category:    "medication_complexity",
			Level:       "medium",
			Score:       40.0,
			Description: fmt.Sprintf("Managing %d prescriptions simultaneously", len(prescriptions)),
		})
	}

	return riskFactors
}

func convertTestResults(testResults []models.TestKitResult) []TestSummary {
	var summaries []TestSummary
	for _, result := range testResults {
		summaries = append(summaries, TestSummary{
			TestType:   "test_kit",
			Result:     result.Result,
			Confidence: result.AIConfidence,
			Date:       result.CreatedAt,
			Status:     result.Status,
		})
	}
	return summaries
}

func analyzeSymptomPatterns(symptomChecks []models.SymptomCheck) []SymptomPattern {
	symptomFreq := make(map[string]int)
	urgencyMap := make(map[string]string)

	for _, check := range symptomChecks {
		for _, symptom := range check.Symptoms {
			symptomFreq[symptom]++
			if urgencyMap[symptom] == "" || check.UrgencyLevel == "high" || check.UrgencyLevel == "emergency" {
				urgencyMap[symptom] = check.UrgencyLevel
			}
		}
	}

	var patterns []SymptomPattern
	for symptom, freq := range symptomFreq {
		pattern := "one-time"
		if freq > 2 {
			pattern = "recurring"
		}
		if freq > 5 {
			pattern = "chronic"
		}

		patterns = append(patterns, SymptomPattern{
			Symptoms:     []string{symptom},
			Frequency:    freq,
			UrgencyLevel: urgencyMap[symptom],
			Pattern:      pattern,
		})
	}

	return patterns
}

func analyzeConsultationStats(sessions []models.TelehealthSession) ConsultationStats {
	stats := ConsultationStats{
		TotalSessions: len(sessions),
	}

	completed := 0
	var lastSession time.Time

	for _, session := range sessions {
		if session.Status == "completed" {
			completed++
		}
		if session.CreatedAt.After(lastSession) {
			lastSession = session.CreatedAt
		}
	}

	stats.CompletedSessions = completed
	stats.LastConsultation = lastSession
	if len(sessions) > 0 {
		stats.AverageRating = 4.2
	}

	return stats
}

func analyzePrescriptionHistory(prescriptions []models.Prescription) PrescriptionHistory {
	history := PrescriptionHistory{
		TotalPrescriptions: len(prescriptions),
	}

	active := 0
	for _, prescription := range prescriptions {
		if prescription.Status == "approved" || prescription.Status == "dispensed" {
			active++
		}
	}

	history.ActiveMedications = active
	history.ComplianceRate = 85.0
	history.CommonMedications = []string{"Pain Relief", "Antibiotics"}

	return history
}

func convertLabBookings(bookings []models.LabBooking) []LabSummary {
	var summaries []LabSummary
	for _, booking := range bookings {
		summaries = append(summaries, LabSummary{
			TestName:   "Lab Test",
			Status:     booking.Status,
			BookedDate: booking.BookingDate,
			Category:   "laboratory",
		})
	}
	return summaries
}

func generateTrendInsights(data AnalyticsData, timeRange string) []string {
	var insights []string

	if len(data.HealthTrends) > 0 {
		for _, trend := range data.HealthTrends {
			if trend.Direction == "improving" {
				insights = append(insights, fmt.Sprintf("Your %s shows improvement over the past %s", trend.Metric, timeRange))
			} else if trend.Direction == "declining" {
				insights = append(insights, fmt.Sprintf("Your %s shows a declining trend over the past %s - consider consulting a healthcare provider", trend.Metric, timeRange))
			}
		}
	}

	if len(data.TestResults) > 0 {
		insights = append(insights, fmt.Sprintf("You have completed %d health assessments in the past %s", len(data.TestResults), timeRange))
	}

	if data.ConsultationStats.TotalSessions > 0 {
		insights = append(insights, fmt.Sprintf("You've had %d telehealth consultations with a %d%% completion rate",
			data.ConsultationStats.TotalSessions,
			int((float64(data.ConsultationStats.CompletedSessions)/float64(data.ConsultationStats.TotalSessions))*100)))
	}

	return insights
}

func generateRiskInsights(data AnalyticsData, timeRange string) []string {
	var insights []string

	highRiskCount := 0
	for _, risk := range data.RiskFactors {
		if risk.Level == "high" {
			highRiskCount++
			insights = append(insights, fmt.Sprintf("High risk factor identified: %s", risk.Description))
		}
	}

	if highRiskCount == 0 {
		insights = append(insights, "Your current risk assessment shows low to moderate risk levels")
	}

	if len(data.SymptomPatterns) > 0 {
		chronicSymptoms := 0
		for _, pattern := range data.SymptomPatterns {
			if pattern.Pattern == "chronic" {
				chronicSymptoms++
			}
		}
		if chronicSymptoms > 0 {
			insights = append(insights, fmt.Sprintf("You have %d chronic symptom patterns that may need ongoing management", chronicSymptoms))
		}
	}

	return insights
}

func generateWellnessInsights(data AnalyticsData, timeRange string) []string {
	var insights []string

	totalActivities := len(data.TestResults) + data.ConsultationStats.TotalSessions + len(data.LabResults)

	if totalActivities > 5 {
		insights = append(insights, "You're actively engaged in managing your health with regular monitoring and consultations")
	} else if totalActivities > 2 {
		insights = append(insights, "You show good health awareness with moderate engagement in health activities")
	} else {
		insights = append(insights, "Consider increasing your health monitoring activities for better wellness tracking")
	}

	if data.PrescriptionHistory.ComplianceRate > 80 {
		insights = append(insights, "Excellent medication compliance rate supports your overall wellness")
	}

	return insights
}

func generateComprehensiveInsights(data AnalyticsData, timeRange string) []string {
	insights := []string{}

	trends := generateTrendInsights(data, timeRange)
	risks := generateRiskInsights(data, timeRange)
	wellness := generateWellnessInsights(data, timeRange)

	insights = append(insights, trends...)
	insights = append(insights, risks...)
	insights = append(insights, wellness...)

	insights = append(insights, fmt.Sprintf("Comprehensive analysis over %s shows overall health engagement across multiple areas", timeRange))

	return insights
}

func generateTrendRecommendations(data AnalyticsData) []string {
	var recommendations []string

	for _, trend := range data.HealthTrends {
		if trend.Direction == "declining" {
			recommendations = append(recommendations, "Schedule a consultation to discuss declining health trends")
		}
	}

	if len(data.TestResults) == 0 {
		recommendations = append(recommendations, "Consider regular health screenings to establish baseline trends")
	}

	recommendations = append(recommendations, "Continue monitoring your health trends regularly")
	return recommendations
}

func generateRiskRecommendations(data AnalyticsData) []string {
	var recommendations []string

	for _, risk := range data.RiskFactors {
		if risk.Level == "high" {
			recommendations = append(recommendations, fmt.Sprintf("Address high-risk factor: %s", risk.Category))
		}
	}

	if len(data.RiskFactors) > 3 {
		recommendations = append(recommendations, "Consider comprehensive health evaluation due to multiple risk factors")
	}

	recommendations = append(recommendations, "Maintain regular health monitoring to track risk factors")
	return recommendations
}

func generateWellnessRecommendations(data AnalyticsData) []string {
	recommendations := []string{}

	if data.ConsultationStats.TotalSessions == 0 {
		recommendations = append(recommendations, "Consider scheduling a telehealth consultation for personalized health guidance")
	}

	if len(data.LabResults) == 0 {
		recommendations = append(recommendations, "Annual lab work can provide valuable health insights")
	}

	if data.PrescriptionHistory.ComplianceRate < 80 {
		recommendations = append(recommendations, "Focus on improving medication compliance for better health outcomes")
	}

	recommendations = append(recommendations, "Maintain a balanced approach to preventive care and health monitoring")
	return recommendations
}

func generateComprehensiveRecommendations(data AnalyticsData) []string {
	recommendations := []string{}

	trends := generateTrendRecommendations(data)
	risks := generateRiskRecommendations(data)
	wellness := generateWellnessRecommendations(data)

	recommendations = append(recommendations, trends...)
	recommendations = append(recommendations, risks...)
	recommendations = append(recommendations, wellness...)

	recommendations = append(recommendations, "Create a comprehensive health plan incorporating all identified areas")
	recommendations = append(recommendations, "Schedule quarterly reviews to track progress across all health metrics")

	return recommendations
}

func calculateTrendScore(data AnalyticsData) float64 {
	score := 50.0

	for _, trend := range data.HealthTrends {
		if trend.Direction == "improving" {
			score += 15.0
		} else if trend.Direction == "declining" {
			score -= 10.0
		}
	}

	if len(data.TestResults) > 0 {
		score += float64(len(data.TestResults)) * 2.0
	}

	if data.ConsultationStats.CompletedSessions > 0 {
		score += 10.0
	}

	return math.Min(math.Max(score, 0), 100)
}

func calculateRiskScore(data AnalyticsData) float64 {
	score := 100.0
	for _, risk := range data.RiskFactors {
		switch risk.Level {
		case "high":
			score -= 25.0
		case "medium":
			score -= 15.0
		case "low":
			score -= 5.0
		}
	}

	return math.Min(math.Max(score, 0), 100)
}

func calculateWellnessScore(data AnalyticsData) float64 {
	score := 0.0
	maxScore := 0.0

	maxScore += 25.0
	if len(data.TestResults) > 0 {
		score += 25.0
	}

	maxScore += 20.0
	if data.ConsultationStats.TotalSessions > 0 {
		completionRate := float64(data.ConsultationStats.CompletedSessions) / float64(data.ConsultationStats.TotalSessions)
		score += 20.0 * completionRate
	}

	maxScore += 20.0
	score += 20.0 * (data.PrescriptionHistory.ComplianceRate / 100.0)

	maxScore += 15.0
	if len(data.LabResults) > 0 {
		score += 15.0
	}

	maxScore += 20.0
	riskScore := calculateRiskScore(data)
	score += 20.0 * (riskScore / 100.0)

	return (score / maxScore) * 100
}

func calculateComprehensiveScore(data AnalyticsData) float64 {
	trendScore := calculateTrendScore(data)
	riskScore := calculateRiskScore(data)
	wellnessScore := calculateWellnessScore(data)

	return (trendScore*0.3 + riskScore*0.3 + wellnessScore*0.4)
}

func calculateTrendDirection(values []float64) string {
	if len(values) < 2 {
		return "stable"
	}

	n := len(values)
	var sumX, sumY, sumXY, sumX2 float64

	for i, y := range values {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope := (float64(n)*sumXY - sumX*sumY) / (float64(n)*sumX2 - sumX*sumX)

	if slope > 0.1 {
		return "improving"
	} else if slope < -0.1 {
		return "declining"
	}
	return "stable"
}

func urgencyToScore(urgency string) float64 {
	switch urgency {
	case "low":
		return 1.0
	case "medium":
		return 2.0
	case "high":
		return 3.0
	case "emergency":
		return 4.0
	default:
		return 1.0
	}
}
