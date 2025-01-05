package transport

import (
	"context"
	"songs/internal/app/domain"
)

// SongService defines the interface for song-related operations
type SongService interface {
	// GetSong retrieves a song by ID
	GetSong(ctx context.Context, id int) (*domain.Song, error)

	// GetSongs retrieves a list of songs with optional filtering and pagination
	GetSongs(ctx context.Context, filter map[string]string, page, pageSize int) ([]*domain.Song, int64, error)

	// CreateSong creates a new song
	CreateSong(ctx context.Context, song *domain.Song) (*domain.Song, error)

	// UpdateSong updates an existing song
	UpdateSong(ctx context.Context, id int, song *domain.Song) (*domain.Song, error)

	// PartialUpdateSong updates specific fields of an existing song
	PartialUpdateSong(ctx context.Context, id int, updates map[string]interface{}) (*domain.Song, error)

	// DeleteSong deletes a song by ID
	DeleteSong(ctx context.Context, id int) error

	// GetSongVerses retrieves verses of a song with pagination
	GetSongVerses(ctx context.Context, id int, page, size int) ([]string, int, error)
}
