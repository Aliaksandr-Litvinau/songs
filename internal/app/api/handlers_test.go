package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"songs/internal/app/domain"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock service
type MockSongService struct {
	mock.Mock
}

func (m *MockSongService) CreateSong(ctx context.Context, song *domain.Song) (*domain.Song, error) {
	args := m.Called(ctx, song)
	return args.Get(0).(*domain.Song), args.Error(1)
}

func (m *MockSongService) GetSong(ctx context.Context, id int) (*domain.Song, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Song), args.Error(1)
}

func (m *MockSongService) ListSongs(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.Song, int64, error) {
	args := m.Called(ctx, page, pageSize, filters)
	return args.Get(0).([]domain.Song), args.Get(1).(int64), args.Error(2)
}

func (m *MockSongService) DeleteSong(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSongService) UpdateSong(ctx context.Context, id int, song *domain.Song) (*domain.Song, error) {
	args := m.Called(ctx, id, song)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Song), args.Error(1)
}

func (m *MockSongService) PartialUpdateSong(ctx context.Context, id int, updates map[string]interface{}) (*domain.Song, error) {
	args := m.Called(ctx, id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Song), args.Error(1)
}

func (m *MockSongService) GetSongs(ctx context.Context, filter map[string]string, page, pageSize int) ([]*domain.Song, int64, error) {
	args := m.Called(ctx, filter, page, pageSize)
	return args.Get(0).([]*domain.Song), args.Get(1).(int64), args.Error(2)
}

func (m *MockSongService) GetSongVerses(ctx context.Context, id int, page, size int) ([]string, int, error) {
	args := m.Called(ctx, id, page, size)
	return args.Get(0).([]string), args.Get(1).(int), args.Error(2)
}

func setupTestRouter(mockService *MockSongService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handler := NewHandler(mockService)

	// Register routes directly instead of using RegisterRoutes
	api := router.Group("/api")
	{
		api.GET("/songs", handler.GetSongs)
		api.GET("/songs/:id", handler.GetSong)
		api.POST("/songs", handler.CreateSong)
		api.PUT("/songs/:id", handler.UpdateSong)
		api.PATCH("/songs/:id", handler.PartialUpdateSong)
		api.DELETE("/songs/:id", handler.DeleteSong)
		api.GET("/songs/:id/verses", handler.GetSongVerses)
	}

	return router
}

func TestHandler_CreateSong(t *testing.T) {
	mockService := new(MockSongService)
	router := setupTestRouter(mockService)

	song := domain.Song{
		GroupID:     1,
		Title:       "Test Song",
		ReleaseDate: time.Now(),
		Text:        "Test lyrics",
		Link:        "https://example.com",
	}

	expectedSong := song
	expectedSong.ID = 1

	mockService.On("CreateSong", mock.Anything, mock.MatchedBy(func(s *domain.Song) bool {
		return s.Title == song.Title &&
			s.GroupID == song.GroupID &&
			s.Text == song.Text &&
			s.Link == song.Link
	})).Return(&expectedSong, nil)

	// Create request body
	body, _ := json.Marshal(song)
	req, _ := http.NewRequest(http.MethodPost, "/api/songs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response domain.Song
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedSong.ID, response.ID)

	mockService.AssertExpectations(t)
}

func TestHandler_GetSongs(t *testing.T) {
	mockService := new(MockSongService)
	router := setupTestRouter(mockService)

	songs := []*domain.Song{
		{ID: 1, Title: "Song 1"},
		{ID: 2, Title: "Song 2"},
	}

	mockService.On("GetSongs", mock.Anything, mock.Anything, 1, 10).
		Return(songs, int64(2), nil)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "/api/songs?page=1&page_size=10", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), response["total"])

	mockService.AssertExpectations(t)
}
