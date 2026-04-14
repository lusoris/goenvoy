package fanart

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://webservice.fanart.tv/v3"

// Client is a Fanart.tv API v3 client.
type Client struct {
	*metadata.BaseClient
}

// New creates a Fanart.tv [Client] using the given project API key.
// An optional personal client key can be provided for higher rate limits.
func New(apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "fanart", opts...)
	bc.SetAuth(func(req *http.Request) {
		req.Header.Set("Api-Key", apiKey)
	})
	return &Client{BaseClient: bc}
}

// NewWithClientKey creates a Fanart.tv [Client] with both a project API key
// and a personal client key for higher rate limits.
func NewWithClientKey(apiKey, clientKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "fanart", opts...)
	bc.SetAuth(func(req *http.Request) {
		req.Header.Set("Api-Key", apiKey)
		req.Header.Set("Client-Key", clientKey)
	})
	return &Client{BaseClient: bc}
}

// APIError is returned when the API responds with a non-2xx status code.
type APIError struct {
	StatusCode   int    `json:"-"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error message"`
	// RawBody holds the raw response body when the error response could not be parsed as JSON.
	RawBody string `json:"-"`
}

func (e *APIError) Error() string {
	if e.ErrorMessage != "" {
		return "fanart: " + strconv.Itoa(e.StatusCode) + ": " + e.ErrorMessage
	}
	if e.RawBody != "" {
		return "fanart: " + strconv.Itoa(e.StatusCode) + ": " + e.RawBody
	}
	return "fanart: HTTP " + strconv.Itoa(e.StatusCode)
}

func (c *Client) get(ctx context.Context, path string, dst any) error {
	body, status, err := c.DoRaw(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return err
	}

	if status < 200 || status >= 300 {
		apiErr := &APIError{StatusCode: status}
		if err := json.Unmarshal(body, apiErr); err != nil {
			apiErr.RawBody = string(body)
		}
		return apiErr
	}

	if dst != nil && len(body) > 0 {
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("fanart: decode response: %w", err)
		}
	}

	return nil
}

// GetMovieImages returns all fan artwork for a movie by its TMDb or IMDb ID.
func (c *Client) GetMovieImages(ctx context.Context, id string) (*MovieImages, error) {
	var out MovieImages
	if err := c.get(ctx, "/movies/"+id, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetShowImages returns all fan artwork for a TV show by its TheTVDB ID.
func (c *Client) GetShowImages(ctx context.Context, id string) (*ShowImages, error) {
	var out ShowImages
	if err := c.get(ctx, "/tv/"+id, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetArtistImages returns all fan artwork for a music artist by MusicBrainz ID.
func (c *Client) GetArtistImages(ctx context.Context, mbid string) (*ArtistImages, error) {
	var out ArtistImages
	if err := c.get(ctx, "/music/"+mbid, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAlbumImages returns artwork for a music album by MusicBrainz release group ID.
func (c *Client) GetAlbumImages(ctx context.Context, mbid string) (*AlbumImagesResponse, error) {
	var out AlbumImagesResponse
	if err := c.get(ctx, "/music/albums/"+mbid, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLabelImages returns artwork for a music label by MusicBrainz label ID.
func (c *Client) GetLabelImages(ctx context.Context, mbid string) (*LabelImages, error) {
	var out LabelImages
	if err := c.get(ctx, "/music/labels/"+mbid, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLatestMovies returns movies that recently received new fan artwork.
// The optional date parameter filters results to entries after the given Unix timestamp.
func (c *Client) GetLatestMovies(ctx context.Context, date int64) ([]LatestMovie, error) {
	path := "/movies/latest"
	if date > 0 {
		path += "?date=" + strconv.FormatInt(date, 10)
	}
	var out []LatestMovie
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLatestShows returns TV shows that recently received new fan artwork.
// The optional date parameter filters results to entries after the given Unix timestamp.
func (c *Client) GetLatestShows(ctx context.Context, date int64) ([]LatestShow, error) {
	path := "/tv/latest"
	if date > 0 {
		path += "?date=" + strconv.FormatInt(date, 10)
	}
	var out []LatestShow
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLatestArtists returns music artists that recently received new fan artwork.
// The optional date parameter filters results to entries after the given Unix timestamp.
func (c *Client) GetLatestArtists(ctx context.Context, date int64) ([]LatestArtist, error) {
	path := "/music/latest"
	if date > 0 {
		path += "?date=" + strconv.FormatInt(date, 10)
	}
	var out []LatestArtist
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}
