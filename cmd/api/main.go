package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/nyumbanicare/docs"
	"github.com/nyumbanicare/internal/api"
	"github.com/nyumbanicare/internal/config"
	"github.com/nyumbanicare/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Initializing external services...")

	emailConfig := cfg.Email
	if emailConfig.Provider == "smtp" {
		log.Printf("Email configured with Dreamhost SMTP (%s)", emailConfig.SMTPHost)
	} else {
		log.Printf("Email using %s provider", emailConfig.Provider)
	}

	storageConfig := cfg.Storage
	if storageConfig.Provider == "cloudinary" && storageConfig.CloudName != "" {
		log.Printf("Storage configured with Cloudinary (%s)", storageConfig.CloudName)
	} else {
		log.Printf("Storage using %s provider", storageConfig.Provider)
	}
	paymentConfig := cfg.Payment
	if paymentConfig.PaystackSecretKey != "" {
		log.Println("Payment configured with Paystack")
	} else {
		log.Println("Warning: No payment provider configured")
	}

	if cfg.External.ChatGPTAPIKey != "" {
		log.Println("AI service configured with ChatGPT")
	} else {
		log.Println("Warning: ChatGPT API key not configured, using mock responses")
	}

	router := gin.Default()

	api.SetupRoutes(router, db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		for i := 8081; i <= 8090; i++ {
			if _, err := net.Listen("tcp", fmt.Sprintf(":%d", i)); err == nil {
				port = fmt.Sprintf("%d", i)
				break
			}
		}
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
