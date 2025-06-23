package database

import (
	"fmt"

	"github.com/google/uuid"
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
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		fmt.Printf("Warning: Could not create uuid-ossp extension: %v\n", err)
	}

	if err := db.Exec("SELECT uuid_generate_v4()").Error; err != nil {
		fmt.Printf("Warning: UUID functions may not be available: %v\n", err)
	}

	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	fmt.Println("Database connected successfully!")

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	fmt.Println("Starting conservative database migration...")

	coreModels := []interface{}{
		&models.User{},
		&models.TestKit{},
		&models.Notification{},
	}

	for _, model := range coreModels {
		fmt.Printf("Migrating %T...\n", model)
		if err := db.AutoMigrate(model); err != nil {
			fmt.Printf("Warning: Failed to migrate %T: %v\n", model, err)
		} else {
			fmt.Printf("Successfully migrated %T\n", model)
		}
	}
	relatedModels := []interface{}{
		&models.TestKitOrder{},
		&models.TestKitResult{},
		&models.Payment{},
		&models.MedicalRecord{},
	}

	for _, model := range relatedModels {
		fmt.Printf("Migrating %T...\n", model)
		if err := db.AutoMigrate(model); err != nil {
			fmt.Printf("Warning: Failed to migrate %T: %v\n", model, err)
		} else {
			fmt.Printf("Successfully migrated %T\n", model)
		}
	}

	fmt.Println("Creating health_articles table manually...")
	createHealthArticlesTable(db)
	var count int64
	db.Model(&models.TestKit{}).Count(&count)
	if count == 0 {
		fmt.Println("Creating default test kit data...")
		defaultKit := models.TestKit{
			Name:         "COVID-19 Home Test Kit",
			Description:  "A rapid antigen test kit for COVID-19",
			Price:        19.99,
			Stock:        100,
			Category:     "COVID-19",
			Instructions: "Follow enclosed instructions for use",
		}

		if err := db.Create(&defaultKit).Error; err != nil {
			fmt.Printf("Warning: Failed to create default test kit: %v\n", err)
		}
	}

	db.Model(&models.HealthArticle{}).Count(&count)
	if count == 0 {
		fmt.Println("Creating default health article...")
		var userCount int64
		db.Model(&models.User{}).Count(&userCount)

		var authorID uuid.UUID
		if userCount == 0 {
			admin := models.User{
				Email:       "admin@nyumbanicare.com",
				Password:    "$2a$10$IrCaCe.tUQFwQI9sG/vjz.JzKjNwCptatAyXHwW7lyXWgQCFrZMD.", // hashed "admin123"
				FirstName:   "Admin",
				LastName:    "User",
				Role:        "admin",
				IsVerified:  true,
				PhoneNumber: "254700000000",
			}
			if err := db.Create(&admin).Error; err == nil {
				authorID = admin.ID
			} else {
				fmt.Printf("Warning: Failed to create admin user: %v\n", err)
				// Generate a random ID as fallback
				authorID = uuid.New()
			}
		} else {
			var user models.User
			if err := db.First(&user).Error; err == nil {
				authorID = user.ID
			} else {
				authorID = uuid.New()
			}
		}

		article := models.HealthArticle{
			Title:     "Understanding COVID-19 Testing",
			Content:   "COVID-19 testing is essential for controlling the spread of the virus. This article explains the different types of tests available and when to use them.",
			Summary:   "Learn about the different COVID-19 testing options",
			Category:  "COVID-19",
			Tags:      []string{"covid", "testing", "health"},
			AuthorID:  authorID,
			Published: true,
			ReadTime:  5,
		}

		if err := db.Create(&article).Error; err != nil {
			fmt.Printf("Warning: Failed to create default health article: %v\n", err)
		}
	}

	return nil
}

func createHealthArticlesTable(db *gorm.DB) {
	var count int64
	db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'health_articles'").Count(&count)

	if count == 0 {
		fmt.Println("Creating health_articles table...")
		sqlSchema := `
		CREATE TABLE health_articles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(255) NOT NULL,
			content TEXT,
			summary TEXT,
			category VARCHAR(100),
			author_id UUID,
			image_url TEXT,
			video_url TEXT,
			read_time INTEGER DEFAULT 0,
			published BOOLEAN DEFAULT false,
			view_count INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		);
		`

		if err := db.Exec(sqlSchema).Error; err != nil {
			fmt.Printf("Warning: Failed to create health_articles table: %v\n", err)
		} else {
			fmt.Println("Successfully created health_articles table.")

			sampleArticleSQL := `
			INSERT INTO health_articles (
				title, content, summary, category, published, read_time, created_at, updated_at
			) VALUES (
				'Understanding COVID-19 Testing',
				'COVID-19 testing is essential for controlling the spread of the virus. This article explains the different types of tests available and when to use them.',
				'Learn about different COVID-19 testing options',
				'COVID-19',
				true,
				5,
				NOW(),
				NOW()
			);
			`

			if err := db.Exec(sampleArticleSQL).Error; err != nil {
				fmt.Printf("Warning: Failed to create sample health article: %v\n", err)
			} else {
				fmt.Println("Successfully created sample health article.")
			}
		}
	} else {
		fmt.Println("Health articles table already exists.")
	}
}
