package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyumbanicare/internal/config"
	"github.com/nyumbanicare/internal/models"
	"github.com/nyumbanicare/internal/services"
	"gorm.io/gorm"
)

// @Summary Create a notification
// @Description Create a new notification for a user
// @Tags Notifications
// @Accept json
// @Produce json
// @Param notification body models.Notification true "Notification object"
// @Success 201 {object} models.Notification "Created notification"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/notifications [post]
// @Security Bearer
func CreateNotification(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var notification models.Notification
		if err := c.BindJSON(&notification); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification data"})
			return
		}

		notification.UserID = userID.(uuid.UUID)
		notification.ID = uuid.New()
		notification.IsRead = false
		notification.IsSent = false

		if err := db.Create(&notification).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
			return
		}

		c.JSON(http.StatusCreated, notification)
	}
}

// @Summary Send an email notification
// @Description Send an email notification to a user or external recipient
// @Tags Notifications
// @Accept json
// @Produce json
// @Param emailNotification body map[string]interface{} true "Email Notification object"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/notifications/email [post]
// @Security Bearer
func SendEmailNotification(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var emailReq struct {
			To       string                 `json:"to" binding:"required,email"`
			Subject  string                 `json:"subject" binding:"required"`
			Message  string                 `json:"message" binding:"required"`
			Metadata map[string]interface{} `json:"metadata"`
		}

		if err := c.BindJSON(&emailReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification data"})
			return
		}

		// Create the notification record
		notification := models.Notification{
			ID:         uuid.New(),
			UserID:     userID.(uuid.UUID),
			Title:      emailReq.Subject,
			Message:    emailReq.Message,
			SendTo:     emailReq.To,
			SendMethod: "email",
			IsRead:     false,
			IsSent:     false,
		}

		// If metadata contains test result info, store it
		if emailReq.Metadata != nil {
			if resultType, ok := emailReq.Metadata["type"].(string); ok && resultType == "test_result" {
				notification.Type = models.NotificationTypeTestResult
				notification.ResourceType = "test_result"
				if resultID, ok := emailReq.Metadata["resultId"].(string); ok {
					resultUUID, err := uuid.Parse(resultID)
					if err == nil {
						notification.ResourceID = &resultUUID
					}
				}
			}
		}

		if err := db.Create(&notification).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification record"})
			return
		}
		emailSvc := services.NewEmailService(&config.GetConfig().Email)

		emailData := services.EmailData{
			To:      emailReq.To,
			Subject: emailReq.Subject,
			Content: emailReq.Message,
			Type:    "html",
		}

		if err := emailSvc.SendEmail(emailData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
			return
		}

		now := time.Now()
		notification.IsSent = true
		notification.SentAt = &now
		db.Save(&notification)

		c.JSON(http.StatusOK, gin.H{
			"message": "Email notification sent successfully",
			"id":      notification.ID,
		})
	}
}

func GetUserNotifications(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		unreadOnly := c.Query("unread") == "true"
		notificationType := c.Query("type")

		var notifications []models.Notification
		query := db.Where("user_id = ?", userID)

		if unreadOnly {
			query = query.Where("is_read = ?", false)
		}

		if notificationType != "" {
			query = query.Where("type = ?", notificationType)
		}

		if err := query.Order("created_at DESC").Find(&notifications).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
			return
		}

		c.JSON(http.StatusOK, notifications)
	}
}

func MarkNotificationAsRead(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		notificationID := c.Param("id")
		notifUUID, err := uuid.Parse(notificationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
			return
		}

		var notification models.Notification
		if err := db.Where("id = ? AND user_id = ?", notifUUID, userID).First(&notification).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}

		notification.IsRead = true
		if err := db.Save(&notification).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
	}
}

func MarkAllNotificationsAsRead(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		result := db.Model(&models.Notification{}).
			Where("user_id = ? AND is_read = ?", userID, false).
			Updates(map[string]interface{}{"is_read": true})

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "All notifications marked as read",
			"count":   result.RowsAffected,
		})
	}
}

func GetUnreadNotificationCount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var count int64
		if err := db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count notifications"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"count": count})
	}
}

func NotifyTestResultReady(db *gorm.DB, userID uuid.UUID, testResult *models.TestKitResult) error {
	notification := models.Notification{
		ID:           uuid.New(),
		UserID:       userID,
		Type:         models.NotificationTypeTestResult,
		Title:        "Your test result is ready",
		Message:      "Your test result has been analyzed and is now available.",
		ResourceID:   &testResult.ID,
		ResourceType: "test_result",
		IsRead:       false,
		IsSent:       false,
		SendMethod:   "app",
	}

	return db.Create(&notification).Error
}
