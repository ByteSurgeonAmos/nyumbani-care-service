package api

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyumbanicare/internal/config"
	"github.com/nyumbanicare/internal/models"
	"github.com/nyumbanicare/internal/services"
	"gorm.io/gorm"
)

type PaystackPaymentRequest struct {
	Email   string  `json:"email" binding:"required,email"`
	Amount  float64 `json:"amount" binding:"required"`
	OrderID string  `json:"order_id" binding:"required"`
}

func ProcessPaystackPayment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req PaystackPaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate order exists and belongs to user
		var order models.TestKitOrder
		if err := db.Where("id = ? AND user_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		// Create payment service
		paymentSvc := services.NewPaymentService(&config.GetConfig().Payment)

		// Generate callback URL for this payment
		host := c.Request.Host
		protocol := "https"
		if host == "localhost" || host == "127.0.0.1" {
			protocol = "http"
		}
		callbackURL := fmt.Sprintf("%s://%s/api/payments/paystack/callback", protocol, host)

		// Create payment record
		payment, err := paymentSvc.InitiatePayment(&order, req.Email, callbackURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize payment"})
			return
		}

		// Save payment record
		if err := db.Create(payment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"payment_id": payment.ID,
			"reference":  payment.TransactionID,
			"status":     "pending",
			"message":    "Payment initialized successfully",
		})
	}
}

// Get payment status
func GetPaymentStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		paymentID := c.Param("id")
		var payment models.Payment
		query := db.Where("id = ?", paymentID)

		// Non-admin users can only see their own payments
		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.First(&payment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}

		c.JSON(http.StatusOK, payment)
	}
}

// List user payments
func ListUserPayments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var payments []models.Payment
		query := db.Preload("Order").Preload("Order.TestKit")

		// Non-admin users can only see their own payments
		role, _ := c.Get("role")
		if role != "admin" {
			query = query.Where("user_id = ?", userID)
		}

		if err := query.Find(&payments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
			return
		}

		c.JSON(http.StatusOK, payments)
	}
}

// Paystack webhook handler
func PaystackWebhook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) { // Read request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
			return
		}

		// Get configuration
		cfg := config.GetConfig()

		// Note: Paystack doesn't provide a webhook secret for verification
		// IP address validation can be implemented via middleware based
		// on Paystack's documented webhook IPs if needed

		// Create payment service
		paymentSvc := services.NewPaymentService(&cfg.Payment)

		// Process webhook
		reference, err := paymentSvc.HandlePaystackWebhook(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// If we have a transaction reference, update the payment and order
		if reference != "" {
			var payment models.Payment
			if err := db.Where("transaction_id = ?", reference).First(&payment).Error; err == nil {
				// Update payment status and add transaction info
				payment.Status = "completed"
				payment.UpdatedAt = time.Now()
				payment.Notes = fmt.Sprintf("Payment confirmed via webhook on %s", time.Now().Format(time.RFC3339))

				if err := db.Save(&payment).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
					return
				}

				// Update order status
				var order models.TestKitOrder
				if err := db.First(&order, "id = ?", payment.OrderID).Error; err == nil {
					order.PaymentStatus = "paid"
					order.Status = "confirmed"
					order.UpdatedAt = time.Now()

					if err := db.Save(&order).Error; err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
						return
					}
				}
			} else {
				// Log that a payment was not found
				c.JSON(http.StatusOK, gin.H{
					"status":  "skipped",
					"message": "Payment record not found for the provided reference",
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

// Paystack callback handler (for browser redirects after payment)
func PaystackCallback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		reference := c.Query("reference")
		if reference == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No reference provided"})
			return
		}

		// Create payment service
		paymentSvc := services.NewPaymentService(&config.GetConfig().Payment)

		// Verify payment
		verifyResp, err := paymentSvc.VerifyPayment(reference)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify payment"})
			return
		}
		if verifyResp.Status && verifyResp.Data.Status == "success" {
			// Payment successful, update database
			var payment models.Payment
			if err := db.Where("transaction_id = ?", reference).First(&payment).Error; err == nil {
				// Update payment status and add transaction info
				payment.Status = "completed"
				payment.UpdatedAt = time.Now()
				payment.Notes = fmt.Sprintf("Payment verified via callback on %s", time.Now().Format(time.RFC3339))

				if err := db.Save(&payment).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
					return
				}

				// Update order status
				var order models.TestKitOrder
				if err := db.First(&order, "id = ?", payment.OrderID).Error; err == nil {
					order.PaymentStatus = "paid"
					order.Status = "confirmed"
					order.UpdatedAt = time.Now()

					if err := db.Save(&order).Error; err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
						return
					}
				}

				// Redirect to success page
				c.Redirect(http.StatusFound, "/payments/success")
				return
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Payment record not found"})
				return
			}
		}

		// Payment failed or not found
		c.Redirect(http.StatusFound, "/payments/failed")
	}
}

// File upload handler
func UploadFile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get folder from query param, default to "general"
		folder := c.DefaultQuery("folder", "general")

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

		// Upload file to Cloudinary
		fileURL, err := storageSvc.UploadFile(header, folder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"file_url":  fileURL,
			"file_name": header.Filename,
			"message":   "File uploaded successfully to Cloudinary",
		})
	}
}
