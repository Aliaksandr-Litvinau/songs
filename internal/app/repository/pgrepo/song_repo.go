package pgrepo

import (
	"context"
	"errors"
	"fmt"
	"songs/internal/app/domain"
	"strings"

	"gorm.io/gorm"
)

// SongRepo implements repository pattern for songs
type SongRepo struct {
	db *gorm.DB
}

// NewSongRepo creates a new song repository
func NewSongRepo(db *gorm.DB) *SongRepo {
	return &SongRepo{
		db: db,
	}
}

// GetSong retrieves a song by ID
func (r SongRepo) GetSong(ctx context.Context, id int) (*domain.Song, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid song ID")
	}

	var song domain.Song
	result := r.db.WithContext(ctx).First(&song, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("song not found")
		}
		return nil, fmt.Errorf("failed to get song: %w", result.Error)
	}

	return &song, nil
}

// GetSongs retrieves a list of songs with filtering and pagination
func (r SongRepo) GetSongs(ctx context.Context, filter map[string]string, page, pageSize int) ([]*domain.Song, int64, error) {
	var songs []*domain.Song
	var total int64
	query := r.db.WithContext(ctx).Model(&domain.Song{})

	// Apply filters
	for key, value := range filter {
		switch key {
		case "title":
			query = query.Where("title LIKE ?", "%"+value+"%")
		case "group_id":
			query = query.Where("group_id = ?", value)
		case "link":
			query = query.Where("link LIKE ?", "%"+value+"%")
		}
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count songs: %w", err)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	result := query.
		Offset(offset).
		Limit(pageSize).
		Order("id asc").
		Find(&songs)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to get songs: %w", result.Error)
	}

	return songs, total, nil
}

// CreateSong creates a new song
func (r SongRepo) CreateSong(ctx context.Context, song *domain.Song) (*domain.Song, error) {
	result := r.db.WithContext(ctx).Create(song)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create song: %w", result.Error)
	}

	return song, nil
}

// UpdateSong updates an existing song
func (r *SongRepo) UpdateSong(ctx context.Context, id int, song *domain.Song) (*domain.Song, error) {
	var existingSong domain.Song
	if err := r.db.WithContext(ctx).First(&existingSong, id).Error; err != nil {
		return nil, err
	}

	song.ID = id

	if err := r.db.WithContext(ctx).Save(song).Error; err != nil {
		return nil, fmt.Errorf("failed to update song: %w", err)
	}

	var updatedSong domain.Song
	if err := r.db.WithContext(ctx).First(&updatedSong, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated song: %w", err)
	}

	return &updatedSong, nil
}

// PartialUpdateSong updates specific fields of a song
func (r *SongRepo) PartialUpdateSong(ctx context.Context, id int, updates map[string]interface{}) (*domain.Song, error) {
	var song domain.Song
	if err := r.db.WithContext(ctx).First(&song, id).Error; err != nil {
		return nil, err
	}

	result := r.db.WithContext(ctx).Model(&song).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update song: %w", result.Error)
	}

	if err := r.db.WithContext(ctx).First(&song, id).Error; err != nil {
		return nil, err
	}

	return &song, nil
}

// DeleteSong deletes a song by ID
func (r SongRepo) DeleteSong(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&domain.Song{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete song: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("song not found")
	}

	return nil
}

// GetSongVerses retrieves verses of a song with pagination
func (r SongRepo) GetSongVerses(ctx context.Context, id int, page, size int) ([]string, int, error) {
	song, err := r.GetSong(ctx, id)
	if err != nil {
		return nil, 0, err
	}

	// Split text into verses (assuming verses are separated by blank lines)
	verses := strings.Split(song.Text, "\n\n")
	totalVerses := len(verses)

	// Calculate pagination
	start := (page - 1) * size
	end := start + size
	if end > totalVerses {
		end = totalVerses
	}
	if start >= totalVerses {
		return []string{}, totalVerses, nil
	}

	return verses[start:end], totalVerses, nil
}
