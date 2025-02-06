package transport

import (
	"errors"
	"net/http"
	"songs/internal/app/common"
	"songs/internal/app/common/server"
	"songs/internal/app/domain"
	"strconv"
)

type Handler struct {
	songService SongService
}

func NewHandler(songService SongService) *Handler {
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
func (h *Handler) GetSong(r common.RequestReader, w http.ResponseWriter) error {
	songIDStr, err := r.PathParam("id")
	if err != nil {
		server.BadRequest("invalid-song-id", err, w)
		return nil
	}

	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		server.BadRequest("invalid-song-id", err, w)
		return nil
	}

	song, err := h.songService.GetSong(r.Context(), songID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.NotFound("song-not-found", err, w)
			return nil
		}
		server.RespondWithError(err, w)
		return nil
	}

	response := ToSongResponse(song)
	server.RespondOK(response, w)
	return nil
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
func (h *Handler) GetSongs(r common.RequestReader, w http.ResponseWriter) error {
	pageStr := r.DefaultQueryParam("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSizeStr := r.DefaultQueryParam("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	filter := make(map[string]string)
	if title := r.QueryParam("title"); title != "" {
		filter["title"] = title
	}
	if groupID := r.QueryParam("group_id"); groupID != "" {
		filter["group_id"] = groupID
	}

	songs, total, err := h.songService.GetSongs(r.Context(), filter, page, pageSize)
	if err != nil {
		server.RespondWithError(err, w)
		return nil
	}

	server.RespondOK(map[string]interface{}{
		"songs": songs,
		"total": total,
		"page":  page,
		"pages": (int(total) + pageSize - 1) / pageSize,
	}, w)
	return nil
}

// CreateSong godoc
// @Summary Add a new song
// @Description Add a new song to the database
// @Tags songs
// @Accept json
// @Produce json
// @Param song body SongRequest true "Song object"
// @Success 200 {object} SongResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/songs [post]
func (h *Handler) CreateSong(r common.RequestReader, w http.ResponseWriter) error {
	var req SongRequest
	if err := r.DecodeBody(&req); err != nil {
		server.BadRequest("invalid-request-body", err, w)
		return nil
	}

	if err := req.Validate(); err != nil {
		server.BadRequest("validation-failed", err, w)
		return nil
	}

	song, err := ToSongDomain(req)
	if err != nil {
		server.BadRequest("invalid-song-data", domain.ErrInvalidData, w)
		return nil
	}

	createdSong, err := h.songService.CreateSong(r.Context(), song)
	if err != nil {
		server.RespondWithError(err, w)
		return nil
	}

	response := ToSongResponse(createdSong)
	server.RespondOK(response, w)
	return nil
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
func (h *Handler) UpdateSong(r common.RequestReader, w http.ResponseWriter) error {
	var req SongRequest
	if err := r.DecodeBody(&req); err != nil {
		server.BadRequest("invalid-request-body", err, w)
		return nil
	}

	song, err := ToSongDomain(req)
	if err != nil {
		server.BadRequest("invalid-song-data", domain.ErrInvalidData, w)
		return nil
	}

	updatedSong, err := h.songService.UpdateSong(r.Context(), song)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.NotFound("song-not-found", err, w)
			return nil
		}
		server.RespondWithError(err, w)
		return nil
	}

	response := ToSongResponse(updatedSong)
	server.RespondOK(response, w)
	return nil
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
func (h *Handler) PartialUpdateSong(r common.RequestReader, w http.ResponseWriter) error {
	idStr, err := r.PathParam("id")
	if err != nil {
		server.BadRequest("invalid-song-id", domain.ErrInvalidID, w)
		return nil
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.BadRequest("invalid-song-id", domain.ErrInvalidID, w)
		return nil
	}

	var updates map[string]interface{}
	if err := r.DecodeBody(&updates); err != nil {
		server.BadRequest("invalid-request-body", err, w)
		return nil
	}

	updatedSong, err := h.songService.PartialUpdateSong(r.Context(), id, updates)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.NotFound("song-not-found", err, w)
			return nil
		}
		server.RespondWithError(err, w)
		return nil
	}

	response := ToSongResponse(updatedSong)
	server.RespondOK(response, w)
	return nil
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
func (h *Handler) DeleteSong(r common.RequestReader, w http.ResponseWriter) error {
	idStr, err := r.PathParam("id")
	if err != nil {
		server.BadRequest("invalid-song-id", domain.ErrInvalidID, w)
		return nil
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		server.BadRequest("invalid-song-id", domain.ErrInvalidID, w)
		return nil
	}

	if err := h.songService.DeleteSong(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.NotFound("song-not-found", err, w)
			return nil
		}
		server.RespondWithError(err, w)
		return nil
	}

	server.RespondOK("Deleted song", w)
	return nil
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
func (h *Handler) GetSongVerses(r common.RequestReader, w http.ResponseWriter) error {
	songIDStr, err := r.PathParam("id")
	if err != nil {
		server.BadRequest("invalid-song-id", err, w)
		return nil
	}

	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		server.BadRequest("invalid-song-id", err, w)
		return nil
	}

	pageStr := r.DefaultQueryParam("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	sizeStr := r.DefaultQueryParam("size", "1")
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 1
	}

	verses, total, err := h.songService.GetSongVerses(r.Context(), songID, page, size)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.NotFound("song-not-found", err, w)
			return nil
		}
		server.RespondWithError(err, w)
		return nil
	}

	server.RespondOK(map[string]interface{}{
		"verses": verses,
		"total":  total,
		"page":   page,
		"size":   size,
	}, w)
	return nil
}
