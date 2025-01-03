package service

import (
	"context"
	"testing"
	"time"

	"songs/internal/app/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSongRepo is a mock implementation of SongRepository
type MockSongRepo struct {
	mock.Mock
}

func (m *MockSongRepo) GetSong(ctx context.Context, id int) (*domain.Song, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Song), args.Error(1)
}

func (m *MockSongRepo) GetSongs(ctx context.Context, filter map[string]string, page, pageSize int) ([]*domain.Song, int64, error) {
	args := m.Called(ctx, filter, page, pageSize)
	return args.Get(0).([]*domain.Song), args.Get(1).(int64), args.Error(2)
}

func (m *MockSongRepo) CreateSong(ctx context.Context, song *domain.Song) (*domain.Song, error) {
	args := m.Called(ctx, song)
	return args.Get(0).(*domain.Song), args.Error(1)
}

func (m *MockSongRepo) UpdateSong(ctx context.Context, id int, song *domain.Song) (*domain.Song, error) {
	args := m.Called(ctx, id, song)
	return args.Get(0).(*domain.Song), args.Error(1)
}

func (m *MockSongRepo) PartialUpdateSong(ctx context.Context, id int, updates map[string]interface{}) (*domain.Song, error) {
	args := m.Called(ctx, id, updates)
	return args.Get(0).(*domain.Song), args.Error(1)
}

func (m *MockSongRepo) DeleteSong(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Add the missing GetSongVerses method
func (m *MockSongRepo) GetSongVerses(ctx context.Context, id int, page, size int) ([]string, int, error) {
	args := m.Called(ctx, id, page, size)
	return args.Get(0).([]string), args.Get(1).(int), args.Error(2)
}

func TestGetSong(t *testing.T) {
	mockRepo := new(MockSongRepo)
	service := NewSongService(mockRepo)

	ctx := context.Background()
	expectedSong := &domain.Song{
		ID:          1,
		Title:       "Test Song",
		ReleaseDate: time.Now(),
		Text:        "Test lyrics",
		Link:        "http://example.com",
	}

	mockRepo.On("GetSong", ctx, 1).Return(expectedSong, nil)

	song, err := service.GetSong(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedSong, song)
	mockRepo.AssertExpectations(t)
}

func TestGetSongs(t *testing.T) {
	mockRepo := new(MockSongRepo)
	service := NewSongService(mockRepo)

	ctx := context.Background()
	filter := map[string]string{"title": "Test"}
	page := 1
	pageSize := 10

	expectedSongs := []*domain.Song{
		{
			ID:    1,
			Title: "Test Song 1",
		},
		{
			ID:    2,
			Title: "Test Song 2",
		},
	}
	expectedTotal := int64(2)

	mockRepo.On("GetSongs", ctx, filter, page, pageSize).Return(expectedSongs, expectedTotal, nil)

	songs, total, err := service.GetSongs(ctx, filter, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, expectedSongs, songs)
	assert.Equal(t, expectedTotal, total)
	mockRepo.AssertExpectations(t)
}

func TestCreateSong(t *testing.T) {
	mockRepo := new(MockSongRepo)
	service := NewSongService(mockRepo)

	ctx := context.Background()
	newSong := &domain.Song{
		GroupID:     1,
		Title:       "New Test Song",
		ReleaseDate: time.Now(),
		Text:        "New test lyrics",
		Link:        "http://example.com/new",
	}

	expectedSong := &domain.Song{
		ID:          1,
		GroupID:     newSong.GroupID,
		Title:       newSong.Title,
		ReleaseDate: newSong.ReleaseDate,
		Text:        newSong.Text,
		Link:        newSong.Link,
	}

	mockRepo.On("CreateSong", ctx, newSong).Return(expectedSong, nil)

	createdSong, err := service.CreateSong(ctx, newSong)

	assert.NoError(t, err)
	assert.Equal(t, expectedSong, createdSong)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSong(t *testing.T) {
	mockRepo := new(MockSongRepo)
	service := NewSongService(mockRepo)

	ctx := context.Background()
	songID := 1
	updatedSong := &domain.Song{
		ID:          songID,
		GroupID:     1,
		Title:       "Updated Test Song",
		ReleaseDate: time.Now(),
		Text:        "Updated test lyrics",
		Link:        "http://example.com/updated",
	}

	mockRepo.On("UpdateSong", ctx, songID, updatedSong).Return(updatedSong, nil)

	result, err := service.UpdateSong(ctx, songID, updatedSong)

	assert.NoError(t, err)
	assert.Equal(t, updatedSong, result)
	mockRepo.AssertExpectations(t)
}

func TestPartialUpdateSong(t *testing.T) {
	mockRepo := new(MockSongRepo)
	service := NewSongService(mockRepo)

	ctx := context.Background()
	songID := 1
	updates := map[string]interface{}{
		"title": "Partially Updated Song",
		"text":  "New lyrics",
	}

	expectedSong := &domain.Song{
		ID:    songID,
		Title: "Partially Updated Song",
		Text:  "New lyrics",
	}

	mockRepo.On("PartialUpdateSong", ctx, songID, updates).Return(expectedSong, nil)

	result, err := service.PartialUpdateSong(ctx, songID, updates)

	assert.NoError(t, err)
	assert.Equal(t, expectedSong, result)
	mockRepo.AssertExpectations(t)
}

func TestDeleteSong(t *testing.T) {
	mockRepo := new(MockSongRepo)
	service := NewSongService(mockRepo)

	ctx := context.Background()
	songID := 1

	mockRepo.On("DeleteSong", ctx, songID).Return(nil)

	err := service.DeleteSong(ctx, songID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetSongVerses(t *testing.T) {
	mockRepo := new(MockSongRepo)
	service := NewSongService(mockRepo)

	ctx := context.Background()
	songID := 1
	page := 1
	size := 2

	expectedVerses := []string{
		"First verse",
		"Second verse",
	}
	expectedTotal := 4

	mockRepo.On("GetSongVerses", ctx, songID, page, size).Return(expectedVerses, expectedTotal, nil)

	verses, total, err := service.GetSongVerses(ctx, songID, page, size)

	assert.NoError(t, err)
	assert.Equal(t, expectedVerses, verses)
	assert.Equal(t, expectedTotal, total)
	mockRepo.AssertExpectations(t)
}

// Error cases
func TestGetSong_Error(t *testing.T) {
	mockRepo := new(MockSongRepo)
	service := NewSongService(mockRepo)

	ctx := context.Background()
	songID := 999

	mockRepo.On("GetSong", ctx, songID).Return(nil, assert.AnError)

	song, err := service.GetSong(ctx, songID)

	assert.Error(t, err)
	assert.Nil(t, song)
	mockRepo.AssertExpectations(t)
}

func TestDeleteSong_Error(t *testing.T) {
	mockRepo := new(MockSongRepo)
	service := NewSongService(mockRepo)

	ctx := context.Background()
	songID := 999

	mockRepo.On("DeleteSong", ctx, songID).Return(assert.AnError)

	err := service.DeleteSong(ctx, songID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
