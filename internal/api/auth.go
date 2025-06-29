package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyumbanicare/internal/models"
	"github.com/nyumbanicare/internal/utils"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var existingUser models.User
		if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}
		dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date of birth format. Use YYYY-MM-DD"})
			return
		}
		user := models.User{
			ID:          uuid.New(),
			Email:       req.Email,
			Password:    hashedPassword,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			PhoneNumber: req.PhoneNumber,
			DateOfBirth: &dateOfBirth,
			Gender:      req.Gender,
			Address:     req.Address,
			Role:        "patient",
		}

		fmt.Printf("DEBUG: Creating user with ID: %s\n", user.ID.String())
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		fmt.Printf("DEBUG: User created with ID: %s\n", user.ID.String())

		token, err := utils.GenerateToken(user.ID, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		fmt.Printf("DEBUG: Generated token for user ID: %s\n", user.ID.String())

		c.JSON(http.StatusCreated, gin.H{
			"token": token,
			"user": gin.H{
				"id":         user.ID,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"role":       user.Role,
			},
		})
	}
}

func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if !utils.CheckPasswordHash(req.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := utils.GenerateToken(user.ID, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":         user.ID,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"role":       user.Role,
			},
		})
	}
}

func RefreshToken(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			return
		}

		token, err := utils.GenerateToken(userID.(uuid.UUID), role.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
