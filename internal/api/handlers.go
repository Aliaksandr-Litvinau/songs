package api

import (
	"gin/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var Db *gorm.DB

func GetSongs(c *gin.Context) {
	var songs []models.Song
	if err := Db.Find(&songs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve songs"})
		return
	}
	c.JSON(http.StatusOK, songs)
}

func AddSong(c *gin.Context) {
	var input models.Song
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := Db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save song"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Song created successfully", "song": input})
}
