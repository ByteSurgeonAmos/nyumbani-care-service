package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyumbanicare/internal/models"
	"gorm.io/gorm"
)

// Admin Lab Test handlers
func CreateLabTest(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var labTest models.LabTest
		if err := c.ShouldBindJSON(&labTest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&labTest).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lab test"})
			return
		}

		c.JSON(http.StatusCreated, labTest)
	}
}

func UpdateLabTest(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var labTest models.LabTest
		if err := db.First(&labTest, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lab test not found"})
			return
		}

		if err := c.ShouldBindJSON(&labTest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&labTest).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lab test"})
			return
		}

		c.JSON(http.StatusOK, labTest)
	}
}

func DeleteLabTest(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.LabTest{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lab test"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Lab test deleted successfully"})
	}
}

// Admin Health Article handlers
func CreateHealthArticle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var article models.HealthArticle
		if err := c.ShouldBindJSON(&article); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		article.AuthorID = userID.(uuid.UUID)

		if err := db.Create(&article).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create health article"})
			return
		}

		c.JSON(http.StatusCreated, article)
	}
}

func UpdateHealthArticle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var article models.HealthArticle
		if err := db.First(&article, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Health article not found"})
			return
		}

		if err := c.ShouldBindJSON(&article); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&article).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update health article"})
			return
		}

		c.JSON(http.StatusOK, article)
	}
}

func DeleteHealthArticle(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.HealthArticle{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete health article"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Health article deleted successfully"})
	}
}
