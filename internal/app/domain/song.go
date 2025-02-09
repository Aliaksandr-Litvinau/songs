package domain

import "time"

type Song struct {
	ID          int
	GroupID     int
	Title       string
	ReleaseDate time.Time
	Text        string
	Link        string
}

type SongDetail struct {
	ReleaseDate string
	Text        string
	Link        string
}
