package domain

import "time"

// SongGroup represents a collection of songs (e.g., album, playlist, or category)
type SongGroup struct {
	ID        int       `gorm:"primaryKey" json:"id,omitempty"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SongGroup) TableName() string {
	return "groups"
}
