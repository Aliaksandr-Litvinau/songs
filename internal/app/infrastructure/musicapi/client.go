package musicapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"songs/internal/app/domain"
	"time"
)

// Client implements SongAPIClient interface
type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient creates a new instance of Music API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// GetSongDetail fetches detailed song information from the Music API
// Returns song details or an error if the request fails
func (c *Client) GetSongDetail(ctx context.Context, group, song string) (*domain.SongDetail, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", c.baseURL, group, url.QueryEscape(song))
	log.Printf("Making request to: %s", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("Failed to make request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Response status: %d", resp.StatusCode)

	var response SongDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Decoded response: %+v", response)
	return response.ToDomain(), nil
}
