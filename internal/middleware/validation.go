package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Validation middleware
func ValidateUUID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id != "" {
			if _, err := uuid.Parse(id); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// Business logic validation
func ValidateTestKitOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			TestKitID uuid.UUID `json:"test_kit_id"`
			Quantity  int       `json:"quantity"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}

		// Validate quantity
		if req.Quantity <= 0 || req.Quantity > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be between 1 and 10"})
			c.Abort()
			return
		}

		// Store validated data in context
		c.Set("validated_order", req)
		c.Next()
	}
}

func ValidateAppointment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ScheduledAt string `json:"scheduled_at"`
			Duration    int    `json:"duration"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}

		// Parse and validate scheduled time
		scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			c.Abort()
			return
		}

		// Check if appointment is in the future
		if scheduledAt.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Appointment must be scheduled for a future date"})
			c.Abort()
			return
		}

		// Check if appointment is within business hours (8 AM - 8 PM)
		hour := scheduledAt.Hour()
		if hour < 8 || hour > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Appointments must be scheduled between 8 AM and 8 PM"})
			c.Abort()
			return
		}

		// Validate duration
		if req.Duration < 15 || req.Duration > 120 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duration must be between 15 and 120 minutes"})
			c.Abort()
			return
		}

		c.Set("validated_appointment", req)
		c.Next()
	}
}

// Email validation
func ValidateEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}

		email := strings.ToLower(strings.TrimSpace(req.Email))
		if !isValidEmail(email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			c.Abort()
			return
		}

		c.Set("validated_email", email)
		c.Next()
	}
}

// Phone number validation for Kenya
func ValidateKenyanPhone() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			PhoneNumber string `json:"phone_number"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}

		phone := strings.ReplaceAll(req.PhoneNumber, " ", "")
		phone = strings.ReplaceAll(phone, "-", "")

		// Validate Kenyan phone number format
		if !isValidKenyanPhone(phone) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Kenyan phone number format"})
			c.Abort()
			return
		}

		c.Set("validated_phone", phone)
		c.Next()
	}
}

// Helper functions
func isValidEmail(email string) bool {
	// Basic email validation
	if len(email) < 5 || len(email) > 254 {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

func isValidKenyanPhone(phone string) bool {
	// Remove country code if present
	if strings.HasPrefix(phone, "+254") {
		phone = "0" + phone[4:]
	} else if strings.HasPrefix(phone, "254") {
		phone = "0" + phone[3:]
	}

	// Check format: should start with 07 or 01 and be 10 digits
	if len(phone) != 10 {
		return false
	}

	if !strings.HasPrefix(phone, "07") && !strings.HasPrefix(phone, "01") {
		return false
	}

	// Check if all characters after prefix are digits
	for i := 2; i < len(phone); i++ {
		if phone[i] < '0' || phone[i] > '9' {
			return false
		}
	}

	return true
}
