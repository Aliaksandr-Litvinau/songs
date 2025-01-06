package models

import "time"

type SongGroup struct {
	ID        int       `gorm:"primaryKey" json:"id,omitempty"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SongGroup) TableName() string {
	return "groups"
}
