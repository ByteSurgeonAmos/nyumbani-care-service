package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyumbanicare/internal/models"
	"gorm.io/gorm"
)

func GetCurrentUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var userID uuid.UUID
		var err error

		switch v := userIDStr.(type) {
		case string:
			userID, err = uuid.Parse(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
				return
			}
		case uuid.UUID:
			userID = v
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
			return
		}
		var user models.User
		fmt.Printf("DEBUG: Looking for user with ID: %s\n", userID.String())
		if err := db.First(&user, "id = ?", userID).Error; err != nil {
			fmt.Printf("DEBUG: Database error: %v\n", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		fmt.Printf("DEBUG: Found user: %s (%s)\n", user.Email, user.ID.String())

		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var userID uuid.UUID
		var err error

		switch v := userIDStr.(type) {
		case string:
			userID, err = uuid.Parse(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
				return
			}
		case uuid.UUID:
			userID = v
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
			return
		}

		var req struct {
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			PhoneNumber string `json:"phone_number"`
			Address     string `json:"address"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.First(&user, "id = ?", userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		user.FirstName = req.FirstName
		user.LastName = req.LastName
		user.PhoneNumber = req.PhoneNumber
		user.Address = req.Address

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func ListTestKits(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var testKits []models.TestKit
		if err := db.Find(&testKits).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test kits"})
			return
		}

		c.JSON(http.StatusOK, testKits)
	}
}

func GetTestKit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var testKit models.TestKit
		if err := db.First(&testKit, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test kit not found"})
			return
		}

		c.JSON(http.StatusOK, testKit)
	}
}

func CreateTestKit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var testKit models.TestKit
		if err := c.ShouldBindJSON(&testKit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&testKit).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test kit"})
			return
		}

		c.JSON(http.StatusCreated, testKit)
	}
}

func UpdateTestKit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var testKit models.TestKit
		if err := db.First(&testKit, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test kit not found"})
			return
		}

		if err := c.ShouldBindJSON(&testKit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&testKit).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test kit"})
			return
		}

		c.JSON(http.StatusOK, testKit)
	}
}

func DeleteTestKit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.TestKit{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete test kit"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Test kit deleted successfully"})
	}
}

func CreateTestKitOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req struct {
			TestKitID       uuid.UUID `json:"test_kit_id" binding:"required"`
			Quantity        int       `json:"quantity" binding:"required,min=1"`
			PaymentMethod   string    `json:"payment_method" binding:"required"`
			ShippingAddress string    `json:"shipping_address" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get test kit details
		var testKit models.TestKit
		if err := db.First(&testKit, "id = ?", req.TestKitID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test kit not found"})
			return
		}

		// Check stock availability
		if testKit.Stock < req.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}
		order := models.TestKitOrder{
			UserID:          userID.(uuid.UUID),
			TestKitID:       req.TestKitID,
			Quantity:        req.Quantity,
			TotalPrice:      testKit.Price * float64(req.Quantity),
			Status:          "pending",
			PaymentStatus:   "pending",
			PaymentMethod:   req.PaymentMethod,
			ShippingAddress: req.ShippingAddress,
		}

		if err := db.Create(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}

		testKit.Stock -= req.Quantity
		db.Save(&testKit)

		c.JSON(http.StatusCreated, order)
	}
}

func ListUserOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		var orders []models.TestKitOrder
		if err := db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}

func GetOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		id := c.Param("id")
		var order models.TestKitOrder
		query := db.Preload("TestKit").Where("id = ?", id)

		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.First(&order).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func UpdateOrderStatus(db *gorm.DB) gin.HandlerFunc {
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

		var order models.TestKitOrder
		if err := db.First(&order, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		order.Status = req.Status
		if req.TrackingNumber != "" {
			order.TrackingNumber = req.TrackingNumber
		}

		if err := db.Save(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func ListMedicalRecords(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var records []models.MedicalRecord
		query := db.Preload("Medications").Preload("TestResults").Preload("Consultations")

		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.Find(&records).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch medical records"})
			return
		}

		c.JSON(http.StatusOK, records)
	}
}

func GetMedicalRecord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		id := c.Param("id")
		var record models.MedicalRecord
		query := db.Preload("Medications").Preload("TestResults").Preload("Consultations").Where("id = ?", id)

		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.First(&record).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
			return
		}

		c.JSON(http.StatusOK, record)
	}
}

func CreateMedicalRecord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var record models.MedicalRecord
		if err := c.ShouldBindJSON(&record); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		record.UserID = userID.(uuid.UUID)

		if err := db.Create(&record).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medical record"})
			return
		}

		c.JSON(http.StatusCreated, record)
	}
}

func UpdateMedicalRecord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		id := c.Param("id")
		var record models.MedicalRecord
		query := db.Where("id = ?", id)

		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.First(&record).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
			return
		}

		if err := c.ShouldBindJSON(&record); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&record).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update medical record"})
			return
		}

		c.JSON(http.StatusOK, record)
	}
}

func CreateTestResult(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var result models.TestKitResult
		if err := c.ShouldBindJSON(&result); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&result).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test result"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func ListTestResults(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		var results []models.TestKitResult
		query := db

		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.Find(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test results"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

func GetTestResult(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		id := c.Param("id")
		var result models.TestKitResult
		query := db.Where("id = ?", id)

		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.First(&result).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test result not found"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateTestResult(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var result models.TestKitResult
		if err := db.First(&result, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test result not found"})
			return
		}

		if err := c.ShouldBindJSON(&result); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&result).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test result"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func CreateConsultation(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var consultation models.Consultation
		if err := c.ShouldBindJSON(&consultation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&consultation).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create consultation"})
			return
		}

		c.JSON(http.StatusCreated, consultation)
	}
}

func ListConsultations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var consultations []models.Consultation
		query := db.Table("consultations").
			Joins("JOIN medical_records ON consultations.medical_record_id = medical_records.id")

		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("medical_records.user_id = ?", userID)
		}

		if err := query.Find(&consultations).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch consultations"})
			return
		}

		c.JSON(http.StatusOK, consultations)
	}
}

func GetConsultation(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		id := c.Param("id")
		var consultation models.Consultation
		query := db.Table("consultations").
			Joins("JOIN medical_records ON consultations.medical_record_id = medical_records.id").
			Where("consultations.id = ?", id)

		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("medical_records.user_id = ?", userID)
		}

		if err := query.First(&consultation).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Consultation not found"})
			return
		}

		c.JSON(http.StatusOK, consultation)
	}
}

func UpdateConsultation(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var consultation models.Consultation
		if err := db.First(&consultation, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Consultation not found"})
			return
		}

		if err := c.ShouldBindJSON(&consultation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&consultation).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update consultation"})
			return
		}

		c.JSON(http.StatusOK, consultation)
	}
}

func ListUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var users []models.User
		var total int64

		db.Model(&models.User{}).Count(&total)
		if err := db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"users":       users,
			"total":       total,
			"page":        page,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		})
	}
}

func ListAllOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var orders []models.TestKitOrder
		var total int64

		db.Model(&models.TestKitOrder{}).Count(&total)
		if err := db.Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"orders":      orders,
			"total":       total,
			"page":        page,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		})
	}
}
