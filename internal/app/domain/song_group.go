package domain

import "time"

type SongGroup struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SongGroup) TableName() string {
	return "groups"
}
