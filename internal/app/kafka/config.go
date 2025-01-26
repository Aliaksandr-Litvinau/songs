package kafka

import "time"

type Config struct {
	Brokers []string
	Topic   string
}

const (
	TopicSongsUpdates = "songs.updates"
)

type Message struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}
