package service

import (
	"context"
	"songs/internal/app/domain"
)

// SongServiceImpl implements the SongService interface
type SongServiceImpl struct {
	repo SongRepository
}

// SongRepository defines the interface for song repository operations
type SongRepository interface {
	GetSong(ctx context.Context, id int) (*domain.Song, error)
	GetSongs(ctx context.Context, filter map[string]string, page, pageSize int) ([]*domain.Song, int64, error)
	CreateSong(ctx context.Context, song *domain.Song) (*domain.Song, error)
	UpdateSong(ctx context.Context, id int, song *domain.Song) (*domain.Song, error)
	PartialUpdateSong(ctx context.Context, id int, updates map[string]interface{}) (*domain.Song, error)
	DeleteSong(ctx context.Context, id int) error
	GetSongVerses(ctx context.Context, id int, page, size int) ([]string, int, error)
}

// NewSongService creates a new instance of SongService
func NewSongService(repo SongRepository) SongService {
	return &SongServiceImpl{
		repo: repo,
	}
}

// GetSong retrieves a song by ID
func (s *SongServiceImpl) GetSong(ctx context.Context, id int) (*domain.Song, error) {
	return s.repo.GetSong(ctx, id)
}

// GetSongs retrieves a list of songs with filtering and pagination
func (s *SongServiceImpl) GetSongs(ctx context.Context, filter map[string]string, page, pageSize int) ([]*domain.Song, int64, error) {
	return s.repo.GetSongs(ctx, filter, page, pageSize)
}

// CreateSong creates a new song
func (s *SongServiceImpl) CreateSong(ctx context.Context, song *domain.Song) (*domain.Song, error) {
	return s.repo.CreateSong(ctx, song)
}

// UpdateSong updates an existing song
func (s *SongServiceImpl) UpdateSong(ctx context.Context, id int, song *domain.Song) (*domain.Song, error) {
	return s.repo.UpdateSong(ctx, id, song)
}

// PartialUpdateSong updates specific fields of a song
func (s *SongServiceImpl) PartialUpdateSong(ctx context.Context, id int, updates map[string]interface{}) (*domain.Song, error) {
	return s.repo.PartialUpdateSong(ctx, id, updates)
}

// DeleteSong deletes a song by ID
func (s *SongServiceImpl) DeleteSong(ctx context.Context, id int) error {
	return s.repo.DeleteSong(ctx, id)
}

// GetSongVerses retrieves verses of a song with pagination
func (s *SongServiceImpl) GetSongVerses(ctx context.Context, id int, page, size int) ([]string, int, error) {
	return s.repo.GetSongVerses(ctx, id, page, size)
}