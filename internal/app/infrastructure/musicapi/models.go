package musicapi

import "songs/internal/app/domain"

// SongDetailResponse represents the Music API response structure
type SongDetailResponse struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// ToDomain converts API response to domain model
func (r *SongDetailResponse) ToDomain() *domain.SongDetail {
	return &domain.SongDetail{
		ReleaseDate: r.ReleaseDate,
		Text:        r.Text,
		Link:        r.Link,
	}
}
