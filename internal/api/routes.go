package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/nyumbanicare/internal/middleware"
	"gorm.io/gorm"
)

// @title Nyumbani Care Healthcare Platform API
// @version 1.0
// @description A comprehensive healthcare platform API for telemedicine, test kits, prescriptions, and more.
// @termsOfService http://swagger.io/terms/

/// @contact.name API Support
// @contact.url https://nyumbanicare.com/support
// @contact.email support@nyumbanicare.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Add logging middleware
	router.Use(middleware.Logger())
	router.Use(middleware.SecurityLogger())
	router.Use(middleware.RateLimit())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "nyumbani-care-api"})
	})	// Explicitly register docs page
	router.GET("/docs", func(c *gin.Context) {
		c.File("swagger_access.html")
	})
	
	// Serve API docs JSON directly
	router.GET("/api-docs", func(c *gin.Context) {
		c.File("docs/api.json")
	})
	
	// Standard Swagger setup
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve Swagger JSON manually from an embedded spec or file
	router.GET("/swagger.json", func(c *gin.Context) {
		// Try to read swagger.json from the docs directory
		swaggerFile := filepath.Join("docs", "swagger.json")
		if _, err := os.Stat(swaggerFile); !os.IsNotExist(err) {
			c.File(swaggerFile)
			return
		}

		// If file doesn't exist, serve a basic swagger spec
		c.JSON(http.StatusOK, map[string]interface{}{
			"swagger": "2.0",
			"info": map[string]interface{}{
				"title":       "Nyumbani Care Healthcare Platform API",
				"description": "A comprehensive healthcare platform API for telemedicine, test kits, prescriptions, and more.",
				"version":     "1.0",
			},
			"host":     c.Request.Host,
			"basePath": "/",
			"schemes":  []string{"http", "https"},
			"paths": map[string]interface{}{
				// Basic paths documentation
				"/api/v1/test-kits/results/analyze": map[string]interface{}{
					"post": map[string]interface{}{
						"summary":     "Upload and analyze test kit result",
						"description": "Upload test kit result image and get AI analysis",
						"tags":        []string{"TestKitResults"},
						"consumes":    []string{"multipart/form-data"},
						"produces":    []string{"application/json"},
						"parameters": []map[string]interface{}{
							{
								"name":     "test_kit_id",
								"in":       "formData",
								"required": true,
								"type":     "string",
							},
							{
								"name":     "test_kit_type",
								"in":       "formData",
								"required": true,
								"type":     "string",
							},
							{
								"name":     "file",
								"in":       "formData",
								"required": true,
								"type":     "file",
							},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{
								"description": "Test kit result analyzed successfully",
							},
						},
					},
				},
			},
		})
	})
	// Public routes
	public := router.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", RegisterUser(db))
			auth.POST("/login", LoginUser(db))
			auth.POST("/refresh", RefreshToken(db))
		}

		// Webhooks (public for external services)
		webhooks := public.Group("/webhooks")
		{
			webhooks.POST("/mpesa", MPesaWebhook(db))
		}		// Test kits
		testKits := public.Group("/test-kits")
		{
			testKits.GET("", ListTestKits(db))
			testKits.GET("/:id", GetTestKit(db))
		}

		// Lab tests (public listing)
		labTests := public.Group("/lab-tests")
		{
			labTests.GET("", ListLabTests(db))
		}

		// Health education (public access)
		education := public.Group("/health-education")
		{
			education.GET("/articles", ListHealthArticles(db))
			education.GET("/articles/:id", GetHealthArticle(db))
		}
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// User profile
		users := protected.Group("/users")
		{
			users.GET("/me", GetCurrentUser(db))
			users.PUT("/me", UpdateUser(db))
		}

		// Medical records
		records := protected.Group("/medical-records")
		{
			records.GET("", ListMedicalRecords(db))
			records.GET("/:id", GetMedicalRecord(db))
			records.POST("", CreateMedicalRecord(db))
			records.PUT("/:id", UpdateMedicalRecord(db))
		}

		// Test kit orders
		orders := protected.Group("/orders")
		{
			orders.POST("", CreateTestKitOrder(db))
			orders.GET("", ListUserOrders(db))
			orders.GET("/:id", GetOrder(db))
			orders.PUT("/:id/status", UpdateOrderStatus(db))
		}
		// Test results
		results := protected.Group("/test-results")
		{
			results.POST("", CreateTestResult(db))
			results.GET("", ListTestResults(db))
			results.GET("/:id", GetTestResult(db))
			results.PUT("/:id", UpdateTestResult(db))
		}
		
		// Test kit results analysis (AI-powered)
		testKitResults := protected.Group("/test-kits/results")
		{
			testKitResults.POST("/analyze", UploadAndAnalyzeTestKitResult(db))
			testKitResults.GET("", ListTestKitResults(db))
			testKitResults.GET("/:id", GetTestKitResult(db))
			testKitResults.PUT("/:id", UpdateTestKitResult(db))
		}

		// Consultations
		consultations := protected.Group("/consultations")
		{
			consultations.POST("", CreateConsultation(db))
			consultations.GET("", ListConsultations(db))
			consultations.GET("/:id", GetConsultation(db))
			consultations.PUT("/:id", UpdateConsultation(db))
		}

		// Prescriptions
		prescriptions := protected.Group("/prescriptions")
		{
			prescriptions.POST("", CreatePrescription(db))
			prescriptions.GET("", ListPrescriptions(db))
			prescriptions.PUT("/:id/status", UpdatePrescriptionStatus(db))
		}

		// Lab bookings
		labBookings := protected.Group("/lab-bookings")
		{
			labBookings.POST("", CreateLabBooking(db))
			labBookings.GET("", ListLabBookings(db))
			labBookings.PUT("/:id/status", UpdateLabBookingStatus(db))
		}

		// Telehealth sessions
		telehealth := protected.Group("/telehealth")
		{
			telehealth.POST("/sessions", CreateTelehealthSession(db))
			telehealth.GET("/sessions", ListTelehealthSessions(db))
		}

		// AI Symptom Checker
		symptoms := protected.Group("/symptoms")
		{
			symptoms.POST("/check", CreateSymptomCheck(db))
			symptoms.GET("/history", ListSymptomChecks(db))
		}
		// CareSense Analytics
		caresense := protected.Group("/caresense")
		{
			caresense.POST("/analytics", GenerateCareSenseAnalytics(db))
			caresense.GET("/analytics", GetCareSenseAnalytics(db))
		}
		// Payments
		payments := protected.Group("/payments")
		{
			// Primary payment method (Paystack)
			payments.POST("/paystack", ProcessPaystackPayment(db))
			payments.GET("/paystack/callback", PaystackCallback(db))
			
			// Legacy payment methods
			payments.POST("/mpesa", ProcessMPesaPayment(db))
			// Payment management
			payments.GET("", ListUserPayments(db))
			payments.GET("/:id", GetPaymentStatus(db))
		}
		
		// Webhooks (public endpoints)
		webhooks := router.Group("/api/webhooks")
		{
			webhooks.POST("/paystack", PaystackWebhook(db))
		}

		// File uploads
		uploads := protected.Group("/uploads")
		{
			uploads.POST("/file", UploadFile(db))
		}

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			// Test kit management
			admin.POST("/test-kits", CreateTestKit(db))
			admin.PUT("/test-kits/:id", UpdateTestKit(db))
			admin.DELETE("/test-kits/:id", DeleteTestKit(db))
			
			// User management
			admin.GET("/users", ListUsers(db))
			admin.GET("/orders", ListAllOrders(db))
			
			// Lab test management
			admin.POST("/lab-tests", CreateLabTest(db))
			admin.PUT("/lab-tests/:id", UpdateLabTest(db))
			admin.DELETE("/lab-tests/:id", DeleteLabTest(db))
			
			// Health education management
			admin.POST("/health-articles", CreateHealthArticle(db))
			admin.PUT("/health-articles/:id", UpdateHealthArticle(db))
			admin.DELETE("/health-articles/:id", DeleteHealthArticle(db))
		}	}
}