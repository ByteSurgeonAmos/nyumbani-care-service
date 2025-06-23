package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeTestResult  NotificationType = "test_result"
	NotificationTypeAppointment NotificationType = "appointment"
	NotificationTypeMedication  NotificationType = "medication"
	NotificationTypeSystem      NotificationType = "system"
)

type Notification struct {
	ID           uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	Type         NotificationType `gorm:"type:varchar(50);not null" json:"type"`
	Title        string           `gorm:"not null" json:"title"`
	Message      string           `gorm:"not null" json:"message"`
	ResourceID   *uuid.UUID       `gorm:"type:uuid" json:"resource_id,omitempty"` // Optional ID of related resource (test result, appointment, etc.)
	ResourceType string           `json:"resource_type,omitempty"`                // Type of the related resource
	IsRead       bool             `gorm:"default:false" json:"is_read"`
	IsSent       bool             `gorm:"default:false" json:"is_sent"`
	SendTo       string           `json:"send_to,omitempty"` // Email address for email notifications
	SendMethod   string           `json:"send_method"`       // 'app', 'email', 'sms'
	SentAt       *time.Time       `json:"sent_at,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
