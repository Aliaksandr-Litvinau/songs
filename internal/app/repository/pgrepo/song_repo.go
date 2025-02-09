package pgrepo

import (
	"context"
	"errors"
	"songs/internal/app/domain"
	"songs/internal/app/repository/models"
	pg "songs/internal/pkg"
	"strings"

	"gorm.io/gorm"
)

// SongRepo implements repository pattern for songs
type SongRepo struct {
	db *pg.PostgresDB
}

// NewSongRepo creates a new song repository
func NewSongRepo(db *pg.PostgresDB) *SongRepo {
	return &SongRepo{
		db: db,
	}
}

// GetSong retrieves a song by ID
func (r *SongRepo) GetSong(ctx context.Context, id int) (*domain.Song, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidID
	}

	var dbSong models.Song
	result := r.db.WithContext(ctx).First(&dbSong, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrDatabase
	}

	song := dbSong.ToDomain()
	return &song, nil
}

// GetSongs retrieves songs with filtering and pagination
func (r *SongRepo) GetSongs(ctx context.Context, filter map[string]string, page, pageSize int) ([]*domain.Song, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, domain.ErrInvalidData
	}

	var total int64
	query := r.db.WithContext(ctx).Model(&models.Song{})

	// Apply filters
	if title, ok := filter["title"]; ok && title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if groupID, ok := filter["group_id"]; ok && groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, domain.ErrDatabase
	}

	// Get paginated results
	var dbSongs []models.Song
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&dbSongs).Error; err != nil {
		return nil, 0, domain.ErrDatabase
	}

	// Convert to domain models
	songs := make([]*domain.Song, len(dbSongs))
	for i, dbSong := range dbSongs {
		song := dbSong.ToDomain()
		songs[i] = &song
	}

	return songs, total, nil
}

// CreateSong creates a new song
func (r *SongRepo) CreateSong(ctx context.Context, song *domain.Song) (*domain.Song, error) {
	if song == nil {
		return nil, domain.ErrInvalidData
	}

	if err := validateSong(*song); err != nil {
		return nil, err
	}

	dbSong := models.ToDBModel(*song)

	if err := r.db.WithContext(ctx).Create(&dbSong).Error; err != nil {
		if isDuplicateError(err) {
			return nil, domain.ErrDuplicate
		}
		return nil, domain.ErrDatabase
	}

	result := dbSong.ToDomain()
	return &result, nil
}

// UpdateSong updates an existing song
func (r *SongRepo) UpdateSong(ctx context.Context, song *domain.Song) (*domain.Song, error) {
	if song == nil {
		return nil, domain.ErrInvalidData
	}

	if song.ID <= 0 {
		return nil, domain.ErrInvalidID
	}

	if err := validateSong(*song); err != nil {
		return nil, err
	}

	dbSong := models.ToDBModel(*song)

	result := r.db.WithContext(ctx).Save(&dbSong)
	if result.Error != nil {
		if isDuplicateError(result.Error) {
			return nil, domain.ErrDuplicate
		}
		return nil, domain.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrNotFound
	}

	updatedSong := dbSong.ToDomain()
	return &updatedSong, nil
}

// PartialUpdateSong updates specific fields of a song
func (r *SongRepo) PartialUpdateSong(ctx context.Context, id int, updates map[string]interface{}) (*domain.Song, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidID
	}

	if len(updates) == 0 {
		return nil, domain.ErrInvalidData
	}

	var updatedDBSong models.Song
	result := r.db.WithContext(ctx).Model(&models.Song{}).Where("id = ?", id).Updates(updates).First(&updatedDBSong)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		if isDuplicateError(result.Error) {
			return nil, domain.ErrDuplicate
		}
		return nil, domain.ErrDatabase
	}

	song := updatedDBSong.ToDomain()
	return &song, nil
}

// DeleteSong deletes a song by ID
func (r *SongRepo) DeleteSong(ctx context.Context, id int) error {
	if id <= 0 {
		return domain.ErrInvalidID
	}

	result := r.db.WithContext(ctx).Delete(&models.Song{}, id)
	if result.Error != nil {
		return domain.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetSongVerses retrieves verses of a song with pagination
func (r *SongRepo) GetSongVerses(ctx context.Context, id int, page, size int) ([]string, int, error) {
	if id <= 0 {
		return nil, 0, domain.ErrInvalidID
	}

	if page <= 0 || size <= 0 {
		return nil, 0, domain.ErrInvalidData
	}

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

// validateSong validates song fields
func validateSong(song domain.Song) error {
	if song.Title == "" {
		return domain.ErrRequired
	}
	if song.GroupID <= 0 {
		return domain.ErrInvalidID
	}
	if song.ReleaseDate.IsZero() {
		return domain.ErrRequired
	}
	return nil
}

// isDuplicateError checks if the error is a duplicate key error
func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "unique constraint") ||
		strings.Contains(err.Error(), "Duplicate entry")
}
