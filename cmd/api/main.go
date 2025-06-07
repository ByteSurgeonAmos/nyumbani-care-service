// Package main is the entry point for the Nyumbani Care API server
// @title Nyumbani Care Healthcare Platform API
// @version 1.0
// @description A comprehensive healthcare platform API for telemedicine, test kits, prescriptions, and more.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
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
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nyumbanicare/internal/api"
	"github.com/nyumbanicare/internal/config"
	"github.com/nyumbanicare/internal/database"
	_ "github.com/nyumbanicare/docs" // Import Swagger docs
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Initialize external services
	log.Println("Initializing external services...")
	
	// Initialize email service
	emailConfig := cfg.Email
	if emailConfig.Provider == "smtp" {
		log.Printf("Email configured with Dreamhost SMTP (%s)", emailConfig.SMTPHost)
	} else {
		log.Printf("Email using %s provider", emailConfig.Provider)
	}
	
	// Initialize storage service
	storageConfig := cfg.Storage
	if storageConfig.Provider == "cloudinary" && storageConfig.CloudName != "" {
		log.Printf("Storage configured with Cloudinary (%s)", storageConfig.CloudName)
	} else {
		log.Printf("Storage using %s provider", storageConfig.Provider)
	}
	
	// Initialize payment service
	paymentConfig := cfg.Payment
	if paymentConfig.PaystackSecretKey != "" {
		log.Println("Payment configured with Paystack")
	} else if paymentConfig.MPesaConsumerKey != "" {
		log.Println("Payment configured with M-Pesa")
	} else {
		log.Println("Warning: No payment provider configured")
	}
	
	// Initialize AI service
	if cfg.External.ChatGPTAPIKey != "" {
		log.Println("AI service configured with ChatGPT")
	} else {
		log.Println("Warning: ChatGPT API key not configured, using mock responses")
	}

	// Initialize router
	router := gin.Default()

	// Setup API routes
	api.SetupRoutes(router, db)
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		// Try alternative ports if 8080 is in use
		for i := 8081; i <= 8090; i++ {
			if _, err := net.Listen("tcp", fmt.Sprintf(":%d", i)); err == nil {
				port = fmt.Sprintf("%d", i)
				break
			}
		}
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 