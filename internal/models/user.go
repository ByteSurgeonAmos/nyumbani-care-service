package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	Password      string         `gorm:"not null" json:"-"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	PhoneNumber   string         `gorm:"uniqueIndex" json:"phone_number"`
	DateOfBirth   time.Time      `json:"date_of_birth"`
	Gender        string         `json:"gender"`
	Address       string         `json:"address"`
	Role          string         `gorm:"default:'patient'" json:"role"` // patient, doctor, nurse, admin
	IsVerified    bool           `gorm:"default:false" json:"is_verified"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
} 