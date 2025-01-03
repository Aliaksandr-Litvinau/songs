package transport

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"songs/internal/app/service"
	"strconv"
)

type Handler struct {
	songService service.SongService
}

func NewHandler(songService service.SongService) *Handler {
	return &Handler{
		songService: songService,
	}
}

// GetSong godoc
// @Summary Get a song by ID
// @Description Get details of a specific song
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} SongResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/songs/{id} [get]
func (h *Handler) GetSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	song, err := h.songService.GetSong(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// GetSongs godoc
// @Summary List songs
// @Description Get a list of songs with optional filtering and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Filter by group name"
// @Param title query string false "Filter by song title"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/songs [get]
func (h *Handler) GetSongs(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// Create filter map
	filter := make(map[string]string)
	if title := c.Query("title"); title != "" {
		filter["title"] = title
	}
	if groupID := c.Query("group_id"); groupID != "" {
		filter["group_id"] = groupID
	}

	songs, total, err := h.songService.GetSongs(c.Request.Context(), filter, page, pageSize)
	if err != nil {
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

// CreateSong godoc
// @Summary Add a new song
// @Description Add a new song to the database
// @Tags songs
// @Accept json
// @Produce json
// @Param song body SongRequest true "Song object"
// @Success 201 {object} SongResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/songs [post]
func (h *Handler) CreateSong(c *gin.Context) {
	var req SongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	song, err := ToSongDomain(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdSong, err := h.songService.CreateSong(c.Request.Context(), song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create song"})
		return
	}

	c.JSON(http.StatusCreated, ToSongResponse(createdSong))
}

// UpdateSong godoc
// @Summary Update a song
// @Description Update an existing song's details
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body SongRequest true "Updated song object"
// @Success 200 {object} SongResponse
// @Failure 400,404 {object} map[string]string
// @Router /api/v1/songs/{id} [put]
func (h *Handler) UpdateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var req SongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	song, err := ToSongDomain(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSong, err := h.songService.UpdateSong(c.Request.Context(), id, song)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, ToSongResponse(updatedSong))
}

// PartialUpdateSong godoc
// @Summary Partially update a song
// @Description Update specific fields of an existing song
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} SongResponse
// @Failure 400,404 {object} map[string]string
// @Router /api/v1/songs/{id} [patch]
func (h *Handler) PartialUpdateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedSong, err := h.songService.PartialUpdateSong(c.Request.Context(), id, updates)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, updatedSong)
}

// DeleteSong godoc
// @Summary Delete a song
// @Description Delete a song from the database
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]string
// @Failure 404,500 {object} map[string]string
// @Router /api/v1/songs/{id} [delete]
func (h *Handler) DeleteSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	if err := h.songService.DeleteSong(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
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
// @Router /api/v1/songs/{id}/verses [get]
func (h *Handler) GetSongVerses(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "1"))

	verses, total, err := h.songService.GetSongVerses(c.Request.Context(), id, page, size)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"verses": verses,
		"total":  total,
		"page":   page,
		"size":   size,
	})
}
