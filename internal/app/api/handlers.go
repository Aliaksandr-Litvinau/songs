package api

import (
	"encoding/json"
	"fmt"
	"gin/internal/app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var Db *gorm.DB

// GetSong godoc
// @Summary Get a song by ID
// @Description Get details of a specific song
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} models.Song
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/songs/{id} [get]
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

// GetSongs godoc
// @Summary List songs
// @Description Get a list of songs with optional filtering and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Filter by group name"
// @Param song query string false "Filter by song name"
// @Param releaseDate query string false "Filter by release date"
// @Param text query string false "Filter by song text"
// @Param link query string false "Filter by song link"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/songs [get]
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

// AddSong godoc
// @Summary Add a new song
// @Description Add a new song to the database
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song object"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/songs [post]
func AddSong(c *gin.Context) {
	var input models.Song
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	apiURL := os.Getenv("MUSIC_INFO_API_URL") // TODO: added variable in the .env file
	if apiURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API URL not configured"})
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/info?group=%s&song=%s", apiURL, input.Group, input.Song))
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get song details from API"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid response from API"})
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode API response"})
		return
	}

	if err := Db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save song"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Song created successfully", "song": input})
}

// UpdateSong godoc
// @Summary Update a song
// @Description Update an existing song's details
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/songs/{id} [put]
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

// PartialUpdateSong godoc
// @Summary Partially update a song
// @Description Update specific fields of an existing song
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body map[string]interface{} true "Fields to update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/songs/{id} [patch]
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

// DeleteSong godoc
// @Summary Delete a song
// @Description Delete a song from the database
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	id := c.Param("id")

	if err := Db.Delete(&models.Song{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}

// Функция для разбиения текста на куплеты (мое предположение, что куплеты разделяются 2 знаками \n\n)
// TODO: Здесь я иду к аналитику и уточняю логику для получения куплетов, так как в ТЗ она не описана
func getVerses(text string) []string {
	return strings.Split(text, "\n\n")
}

// GetSongVerses godoc
// @Summary Get verses of a song
// @Description Get verses of a specific song with pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of verses per page" default(1)
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /api/songs/{id}/verses [get]
func GetSongVerses(c *gin.Context) {
	var song models.Song
	id := c.Param("id")

	if err := Db.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	verses := getVerses(song.Text)

	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 1
	}

	start := (page - 1) * size
	end := start + size
	if start >= len(verses) {
		c.JSON(http.StatusOK, gin.H{"verses": []string{}, "message": "No more verses"})
		return
	}

	if end > len(verses) {
		end = len(verses)
	}

	c.JSON(http.StatusOK, gin.H{
		"verses": verses[start:end],
		"page":   page,
		"size":   size,
		"total":  len(verses),
	})
}
