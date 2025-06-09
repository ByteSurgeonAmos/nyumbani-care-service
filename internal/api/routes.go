package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nyumbanicare/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	router.Use(middleware.Logger())
	router.Use(middleware.SecurityLogger())
	router.Use(middleware.RateLimit())

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "nyumbani-care-api"})
	})

	router.GET("/api-docs", func(c *gin.Context) {
		c.File("docs/api.json")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	public := router.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", RegisterUser(db))
			auth.POST("/login", LoginUser(db))
			auth.POST("/refresh", RefreshToken(db))
		}
		{
		}
		testKits := public.Group("/test-kits")
		{
			testKits.GET("", ListTestKits(db))
			testKits.GET("/:id", GetTestKit(db))
		}

		labTests := public.Group("/lab-tests")
		{
			labTests.GET("", ListLabTests(db))
		}

		education := public.Group("/health-education")
		{
			education.GET("/articles", ListHealthArticles(db))
			education.GET("/articles/:id", GetHealthArticle(db))
		}
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		users := protected.Group("/users")
		{
			users.GET("/me", GetCurrentUser(db))
			users.PUT("/me", UpdateUser(db))
		}

		records := protected.Group("/medical-records")
		{
			records.GET("", ListMedicalRecords(db))
			records.GET("/:id", GetMedicalRecord(db))
			records.POST("", CreateMedicalRecord(db))
			records.PUT("/:id", UpdateMedicalRecord(db))
		}

		orders := protected.Group("/orders")
		{
			orders.POST("", CreateTestKitOrder(db))
			orders.GET("", ListUserOrders(db))
			orders.GET("/:id", GetOrder(db))
			orders.PUT("/:id/status", UpdateOrderStatus(db))
		}
		results := protected.Group("/test-results")
		{
			results.POST("", CreateTestResult(db))
			results.GET("", ListTestResults(db))
			results.GET("/:id", GetTestResult(db))
			results.PUT("/:id", UpdateTestResult(db))
		}

		testKitResults := protected.Group("/test-kits/results")
		{
			testKitResults.POST("/analyze", UploadAndAnalyzeTestKitResult(db))
			testKitResults.GET("", ListTestKitResults(db))
			testKitResults.GET("/:id", GetTestKitResult(db))
			testKitResults.PUT("/:id", UpdateTestKitResult(db))
		}

		consultations := protected.Group("/consultations")
		{
			consultations.POST("", CreateConsultation(db))
			consultations.GET("", ListConsultations(db))
			consultations.GET("/:id", GetConsultation(db))
			consultations.PUT("/:id", UpdateConsultation(db))
		}

		prescriptions := protected.Group("/prescriptions")
		{
			prescriptions.POST("", CreatePrescription(db))
			prescriptions.GET("", ListPrescriptions(db))
			prescriptions.PUT("/:id/status", UpdatePrescriptionStatus(db))
		}

		labBookings := protected.Group("/lab-bookings")
		{
			labBookings.POST("", CreateLabBooking(db))
			labBookings.GET("", ListLabBookings(db))
			labBookings.PUT("/:id/status", UpdateLabBookingStatus(db))
		}

		telehealth := protected.Group("/telehealth")
		{
			telehealth.POST("/sessions", CreateTelehealthSession(db))
			telehealth.GET("/sessions", ListTelehealthSessions(db))
		}

		symptoms := protected.Group("/symptoms")
		{
			symptoms.POST("/check", CreateSymptomCheck(db))
			symptoms.GET("/history", ListSymptomChecks(db))
		}
		caresense := protected.Group("/caresense")
		{
			caresense.POST("/analytics", GenerateCareSenseAnalytics(db))
			caresense.GET("/analytics", GetCareSenseAnalytics(db))
		}
		payments := protected.Group("/payments")
		{
			payments.POST("/paystack", ProcessPaystackPayment(db))
			payments.GET("/paystack/callback", PaystackCallback(db))

			payments.GET("", ListUserPayments(db))
			payments.GET("/:id", GetPaymentStatus(db))
		}

		webhooks := router.Group("/api/webhooks")
		{
			webhooks.POST("/paystack", PaystackWebhook(db))
		}

		uploads := protected.Group("/uploads")
		{
			uploads.POST("/file", UploadFile(db))
		}

		admin := protected.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/test-kits", CreateTestKit(db))
			admin.PUT("/test-kits/:id", UpdateTestKit(db))
			admin.DELETE("/test-kits/:id", DeleteTestKit(db))

			admin.GET("/users", ListUsers(db))
			admin.GET("/orders", ListAllOrders(db))

			admin.POST("/lab-tests", CreateLabTest(db))
			admin.PUT("/lab-tests/:id", UpdateLabTest(db))
			admin.DELETE("/lab-tests/:id", DeleteLabTest(db))

			admin.POST("/health-articles", CreateHealthArticle(db))
			admin.PUT("/health-articles/:id", UpdateHealthArticle(db))
			admin.DELETE("/health-articles/:id", DeleteHealthArticle(db))
		}
	}
}
