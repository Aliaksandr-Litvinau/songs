package models

import "time"

type Song struct {
	ID          int       `gorm:"primaryKey" json:"id,omitempty"`
	GroupID     int       `gorm:"not null" json:"group_id"`
	Title       string    `gorm:"not null" json:"title"`
	ReleaseDate time.Time `gorm:"not null" json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

func (Song) TableName() string {
	return "songs"
}
