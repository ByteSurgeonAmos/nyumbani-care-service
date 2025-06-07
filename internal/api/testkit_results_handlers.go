package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyumbanicare/internal/config"
	"github.com/nyumbanicare/internal/models"
	"github.com/nyumbanicare/internal/services"
	"gorm.io/gorm"
)

// @Summary Upload and analyze test kit result
// @Description Upload test kit result image and get AI analysis
// @Tags TestKitResults
// @Accept multipart/form-data
// @Produce json
// @Param test_kit_id formData string true "Test Kit ID"
// @Param order_id formData string false "Order ID"
// @Param test_kit_type formData string true "Test Kit Type (e.g., 'covid', 'pregnancy', 'malaria')"
// @Param file formData file true "Test Kit Result Image"
// @Success 200 {object} models.TestKitResult "Test kit result analysis"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/test-kits/results/analyze [post]
// @Security Bearer
func UploadAndAnalyzeTestKitResult(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Get parameters
		testKitID := c.PostForm("test_kit_id")
		orderID := c.PostForm("order_id")
		testKitType := c.PostForm("test_kit_type")

		if testKitID == "" || testKitType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Test kit ID and type are required"})
			return
		}

		// Parse UUIDs
		testKitUUID, err := uuid.Parse(testKitID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test kit ID"})
			return
		}

		var orderUUID uuid.UUID
		if orderID != "" {
			orderUUID, err = uuid.Parse(orderID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
				return
			}
		}

		// Get uploaded file
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
		defer file.Close()

		// Validate file type and size
		if header.Size > 10*1024*1024 { // 10MB limit
			c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 10MB)"})
			return
		}

		// Create storage service for Cloudinary
		storageSvc, err := services.NewStorageService(&config.GetConfig().Storage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize storage service"})
			return
		}
		
		// Upload file to Cloudinary in test_results folder
		fileURL, err := storageSvc.UploadFile(header, "test_results")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}

		// Initialize AI service
		aiSvc := services.NewAIService(&config.GetConfig().External)
		
		// Analyze the test kit result
		analysisReq := &services.TestKitResultRequest{
			TestKitType: testKitType,
			ImageURL:    fileURL,
		}
		
		analysisResp, err := aiSvc.AnalyzeTestKitResult(analysisReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze test kit result"})
			return
		}

		// Create test kit result record
		result := models.TestKitResult{
			ID:              uuid.New(),
			UserID:          userID.(uuid.UUID),
			TestKitID:       testKitUUID,
			ImageURL:        fileURL,
			Result:          analysisResp.Result,
			AIConfidence:    analysisResp.Confidence,
			DetectedMarkers: analysisResp.DetectedMarkers,
			RecommendedSteps: analysisResp.RecommendedSteps,
			Notes:           analysisResp.Notes,
			Status:          "pending", // Pending healthcare professional review
		}

		// Add order ID if provided
		if orderID != "" {
			result.OrderID = orderUUID
		}

		// Save to database
		if err := db.Create(&result).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save test kit result"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Test kit result analyzed successfully",
			"result": result,
		})
	}
}

// @Summary Get test kit result by ID
// @Description Get a specific test kit result by ID
// @Tags TestKitResults
// @Produce json
// @Param id path string true "Test Kit Result ID"
// @Success 200 {object} models.TestKitResult "Test kit result"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/test-kits/results/{id} [get]
// @Security Bearer
func GetTestKitResult(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		resultID := c.Param("id")
		var result models.TestKitResult

		query := db.Preload("TestKit").Where("id = ?", resultID)
		
		// Non-admin users can only see their own results
		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.First(&result).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test kit result not found"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// @Summary List user's test kit results
// @Description Get a list of all test kit results for the current user
// @Tags TestKitResults
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {array} models.TestKitResult "List of test kit results"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/test-kits/results [get]
// @Security Bearer
func ListTestKitResults(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Pagination parameters
		page, _ := c.GetQuery("page")
		limit, _ := c.GetQuery("limit")

		var results []models.TestKitResult
		query := db.Preload("TestKit").Order("created_at DESC")
		
		// Non-admin users can only see their own results
		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		// Apply pagination if provided
		if page != "" && limit != "" {
			query = query.Scopes(Paginate(c))
		}

		if err := query.Find(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test kit results"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

// @Summary Update test kit result
// @Description Update a test kit result (admin or healthcare professional only)
// @Tags TestKitResults
// @Accept json
// @Produce json
// @Param id path string true "Test Kit Result ID"
// @Param result body models.TestKitResult true "Updated test kit result"
// @Success 200 {object} models.TestKitResult "Updated test kit result"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/test-kits/results/{id} [put]
// @Security Bearer
func UpdateTestKitResult(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		role, _ := c.Get("role")
		if role != "admin" && role != "healthcare_professional" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Insufficient permissions"})
			return
		}

		resultID := c.Param("id")
		var existingResult models.TestKitResult
		if err := db.First(&existingResult, "id = ?", resultID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test kit result not found"})
			return
		}

		var updateData struct {
			Result          string   `json:"result"`
			ReviewNotes     string   `json:"review_notes"`
			Status          string   `json:"status"`
			RecommendedSteps []string `json:"recommended_steps"`
		}

		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update only allowed fields
		updates := map[string]interface{}{
			"reviewed_by":      userID,
			"review_notes":     updateData.ReviewNotes,
			"status":           updateData.Status,
		}

		// Only update result if provided
		if updateData.Result != "" {
			updates["result"] = updateData.Result
		}

		// Only update recommended steps if provided
		if len(updateData.RecommendedSteps) > 0 {
			updates["recommended_steps"] = updateData.RecommendedSteps
		}

		if err := db.Model(&existingResult).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test kit result"})
			return
		}

		// Reload the updated record
		db.First(&existingResult, "id = ?", resultID)
		c.JSON(http.StatusOK, existingResult)
	}
}

// Helper function for pagination
func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("limit", "10")
		
		offset := 0
		limit := 10

		if p, err := parseInt(page); err == nil {
			offset = (p - 1) * limit
		}
		
		if l, err := parseInt(pageSize); err == nil {
			limit = l
		}

		return db.Offset(offset).Limit(limit)
	}
}

// Helper function to parse string to int
func parseInt(value string) (int, error) {
	var i int
	_, err := fmt.Sscanf(value, "%d", &i)
	return i, err
}
