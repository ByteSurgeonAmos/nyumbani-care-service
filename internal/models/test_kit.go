package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestKit struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Price        float64        `json:"price"`
	Category     string         `json:"category"`
	Stock        int            `json:"stock"`
	ImageURL     string         `json:"image_url"`
	Instructions string         `json:"instructions"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type TestKitOrder struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	TestKitID       uuid.UUID      `gorm:"type:uuid;not null" json:"test_kit_id"`
	Quantity        int            `json:"quantity"`
	TotalPrice      float64        `json:"total_price"`
	Status          string         `json:"status"`         // pending, paid, shipped, delivered, cancelled
	PaymentStatus   string         `json:"payment_status"` // pending, paid, refunded
	PaymentMethod   string         `json:"payment_method"`
	ShippingAddress string         `json:"shipping_address"`
	TrackingNumber  string         `json:"tracking_number"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TestKitResult moved to extended_models.go with enhanced fields for AI analysis

func (t *TestKit) BeforeCreate(tx *gorm.DB) error {
	// Generate a new UUID if one is not provided
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	// Set default timestamps if not provided
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = time.Now()
	}

	return nil
}

func (o *TestKitOrder) BeforeCreate(tx *gorm.DB) error {
	// Generate a new UUID if one is not provided
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	// Set default timestamps if not provided
	if o.CreatedAt.IsZero() {
		o.CreatedAt = time.Now()
	}
	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = time.Now()
	}

	return nil
}
