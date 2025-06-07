package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/nyumbanicare/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SimpleUser struct {
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"uniqueIndex"`
	Name  string
}

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

	// Test database connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)

	fmt.Printf("Connecting to database with DSN: %s\n", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successful!")

	// Test migration with a simple model
	fmt.Println("Testing migration with simple model...")
	if err := db.AutoMigrate(&SimpleUser{}); err != nil {
		log.Fatalf("Failed to migrate simple model: %v", err)
	}

	fmt.Println("Simple migration successful!")

	// Test database version
	var version string
	if err := db.Raw("SELECT version()").Scan(&version).Error; err != nil {
		log.Printf("Failed to get database version: %v", err)
	} else {
		fmt.Printf("Database version: %s\n", version)
	}

	fmt.Println("Database test completed successfully!")
}
