package api

import (
	"gin/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var Db *gorm.DB

func GetSong(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		var song models.Song
		if err := Db.First(&song, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve song"})
			}
			return
		}
		c.JSON(http.StatusOK, song)
		return
	}
}

func GetSongs(c *gin.Context) {
	var songs []models.Song
	query := Db.Model(&models.Song{})

	// Фильтрация с экранированием имен столбцов
	if group := c.Query("group"); group != "" {
		query = query.Where("\"group\" LIKE ?", "%"+group+"%")
	}
	if song := c.Query("song"); song != "" {
		query = query.Where("\"song\" LIKE ?", "%"+song+"%")
	}
	if releaseDate := c.Query("releaseDate"); releaseDate != "" {
		query = query.Where("\"release_date\" = ?", releaseDate)
	}
	if text := c.Query("text"); text != "" {
		query = query.Where("\"text\" LIKE ?", "%"+text+"%")
	}
	if link := c.Query("link"); link != "" {
		query = query.Where("\"link\" LIKE ?", "%"+link+"%")
	}

	// Пагинация
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&songs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve songs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songs": songs,
		"total": total,
		"page":  page,
		"pages": (int(total) + pageSize - 1) / pageSize,
	})
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

func UpdateSong(c *gin.Context) {
	var input models.Song
	id := c.Param("id")

	var song models.Song
	if err := Db.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	song.Group = input.Group
	song.Song = input.Song
	song.ReleaseDate = input.ReleaseDate
	song.Text = input.Text
	song.Link = input.Link

	if err := Db.Save(&song).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully", "song": song})
}

func PartialUpdateSong(c *gin.Context) {
	id := c.Param("id")

	var song models.Song
	if err := Db.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if group, ok := input["group"]; ok {
		song.Group = group.(string)
	}
	if songName, ok := input["song"]; ok {
		song.Song = songName.(string)
	}
	if releaseDate, ok := input["releaseDate"]; ok {
		song.ReleaseDate = releaseDate.(string)
	}
	if text, ok := input["text"]; ok {
		song.Text = text.(string)
	}
	if link, ok := input["link"]; ok {
		song.Link = link.(string)
	}

	if err := Db.Save(&song).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song partially updated successfully", "song": song})
}

func DeleteSong(c *gin.Context) {
	id := c.Param("id")

	if err := Db.Delete(&models.Song{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
