package models

import (
	"songs/internal/app/domain"
	"time"
)

type Song struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID     int       `gorm:"not null" json:"group_id"`
	Title       string    `gorm:"not null" json:"title"`
	ReleaseDate time.Time `gorm:"not null" json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

func (s *Song) TableName() string {
	return "songs"
}

func (s *Song) ToDomain() domain.Song {
	return domain.Song{
		ID:          s.ID,
		GroupID:     s.GroupID,
		Title:       s.Title,
		ReleaseDate: s.ReleaseDate,
		Text:        s.Text,
		Link:        s.Link,
	}
}

func ToDBModel(s domain.Song) Song {
	return Song{
		ID:          s.ID,
		GroupID:     s.GroupID,
		Title:       s.Title,
		ReleaseDate: s.ReleaseDate,
		Text:        s.Text,
		Link:        s.Link,
	}
}
