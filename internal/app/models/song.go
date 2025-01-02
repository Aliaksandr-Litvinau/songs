package models

import (
	"time"
)

// Group model
type Group struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique;not null"`
	Songs []Song `gorm:"foreignKey:GroupID"`
}

// Song model
type Song struct {
	ID          uint      `gorm:"primaryKey"`
	GroupID     uint      `gorm:"not null"`
	Title       string    `gorm:"not null"`
	ReleaseDate time.Time `gorm:"not null"`
	Text        string    `gorm:"type:text"`
	Link        string    `gorm:"type:varchar(255)"`
}
