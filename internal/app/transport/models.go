package transport

import (
	"fmt"
	"time"
)

type SongRequest struct {
	GroupID     int    `json:"group_id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (r *SongRequest) Validate() error {
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	if r.Text == "" {
		return fmt.Errorf("text is required")
	}
	if _, err := time.Parse(time.RFC3339, r.ReleaseDate); err != nil {
		return fmt.Errorf("invalid release_date format, expected RFC3339")
	}
	return nil
}

type SongResponse struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"group_id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
