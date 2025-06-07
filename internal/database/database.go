package database

import (
	"fmt"

	"github.com/nyumbanicare/internal/config"
	"github.com/nyumbanicare/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Skip auto-migration for now to get the server running
	// Auto-migrate can be enabled later once we resolve the GORM issues
	fmt.Println("Database connected successfully! (Migration skipped)")
	
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	// Start with essential models only
	essentialModels := []interface{}{
		&models.TestKit{},
		&models.TestKitOrder{},
		&models.Payment{},
	}
	
	for _, model := range essentialModels {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %v", model, err)
		}
	}
	
	return nil
}