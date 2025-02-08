package service

import (
	"context"
	"fmt"
	"log"
	"songs/internal/app/domain"
	"strconv"
)

type MusicAPIClient interface {
	GetSongDetail(ctx context.Context, group string, song string) (*domain.SongDetail, error)
}

type SongEnricherService struct {
	client MusicAPIClient
}

func NewSongEnricherService(client MusicAPIClient) *SongEnricherService {
	return &SongEnricherService{
		client: client,
	}
}

func (s *SongEnricherService) EnrichSong(ctx context.Context, song *domain.Song) error {
	groupID := strconv.Itoa(song.GroupID)
	title := song.Title

	log.Printf("Requesting details for group: %s, song: %s", groupID, title)
	details, err := s.client.GetSongDetail(ctx, groupID, title)
	if err != nil {
		return fmt.Errorf("get song details: %w", err)
	}

	song.Text = details.Text
	song.Link = details.Link

	return nil
}
