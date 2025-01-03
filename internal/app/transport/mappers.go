package transport

import (
	"songs/internal/app/domain"
	"time"
)

func ToSongDomain(req SongRequest) (*domain.Song, error) {
	releaseDate, err := time.Parse(time.RFC3339, req.ReleaseDate)
	if err != nil {
		return nil, err
	}

	return &domain.Song{
		GroupID:     req.GroupID,
		Title:       req.Title,
		ReleaseDate: releaseDate,
		Text:        req.Text,
		Link:        req.Link,
	}, nil
}

func ToSongResponse(song *domain.Song) SongResponse {
	return SongResponse{
		ID:          song.ID,
		GroupID:     song.GroupID,
		Title:       song.Title,
		ReleaseDate: song.ReleaseDate.Format(time.RFC3339),
		Text:        song.Text,
		Link:        song.Link,
	}
}
