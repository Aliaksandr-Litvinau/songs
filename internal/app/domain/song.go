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
