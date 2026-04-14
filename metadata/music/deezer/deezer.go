package deezer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://api.deezer.com"

// Client is a Deezer API client.
type Client struct {
	*metadata.BaseClient
	accessToken string
}

// APIError is returned when the Deezer API responds with an error.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
	Type       string
	Message    string
	Code       int
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("deezer: %s: %s (code %d)", e.Type, e.Message, e.Code)
	}
	return fmt.Sprintf("deezer: %s: %s", e.Status, e.Body)
}

// New creates a Deezer [Client].
func New(opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "deezer", opts...)
	return &Client{BaseClient: bc}
}

// NewWithToken creates a Deezer [Client] with an access token for user-specific data.
func NewWithToken(accessToken string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "deezer", opts...)
	return &Client{BaseClient: bc, accessToken: accessToken}
}

func (c *Client) get(ctx context.Context, path string, v any) error {
	u := c.BaseURL() + path

	// Append access token if set.
	if c.accessToken != "" {
		sep := "?"
		if strings.Contains(path, "?") {
			sep = "&"
		}
		u += sep + "access_token=" + url.QueryEscape(c.accessToken)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("deezer: create request: %w", err)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("deezer: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("deezer: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	// Deezer returns 200 with error field for API-level errors.
	var errResp errorResp
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Type:       errResp.Error.Type,
			Message:    errResp.Error.Message,
			Code:       errResp.Error.Code,
		}
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("deezer: decode response: %w", err)
	}
	return nil
}

// Search searches for tracks.
func (c *Client) Search(ctx context.Context, query string, limit, index int) (*SearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("limit", strconv.Itoa(limit))
	params.Set("index", strconv.Itoa(index))

	var result SearchResult
	if err := c.get(ctx, "/search?"+params.Encode(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchAlbums searches for albums.
func (c *Client) SearchAlbums(ctx context.Context, query string, limit, index int) (*AlbumSearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("limit", strconv.Itoa(limit))
	params.Set("index", strconv.Itoa(index))

	var result AlbumSearchResult
	if err := c.get(ctx, "/search/album?"+params.Encode(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchArtists searches for artists.
func (c *Client) SearchArtists(ctx context.Context, query string, limit, index int) (*ArtistSearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("limit", strconv.Itoa(limit))
	params.Set("index", strconv.Itoa(index))

	var result ArtistSearchResult
	if err := c.get(ctx, "/search/artist?"+params.Encode(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetArtist returns an artist by ID.
func (c *Client) GetArtist(ctx context.Context, id int) (*Artist, error) {
	var artist Artist
	if err := c.get(ctx, "/artist/"+strconv.Itoa(id), &artist); err != nil {
		return nil, err
	}
	return &artist, nil
}

// GetArtistTopTracks returns top tracks for an artist.
func (c *Client) GetArtistTopTracks(ctx context.Context, artistID, limit int) ([]Track, error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))

	var resp artistTopResp
	if err := c.get(ctx, "/artist/"+strconv.Itoa(artistID)+"/top?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetArtistAlbums returns albums for an artist.
func (c *Client) GetArtistAlbums(ctx context.Context, artistID int) ([]AlbumSimple, error) {
	var resp artistAlbumsResp
	if err := c.get(ctx, "/artist/"+strconv.Itoa(artistID)+"/albums", &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetRelatedArtists returns artists related to the given artist.
func (c *Client) GetRelatedArtists(ctx context.Context, artistID int) ([]Artist, error) {
	var resp relatedArtistsResp
	if err := c.get(ctx, "/artist/"+strconv.Itoa(artistID)+"/related", &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetAlbum returns an album by ID.
func (c *Client) GetAlbum(ctx context.Context, id int) (*Album, error) {
	var album Album
	if err := c.get(ctx, "/album/"+strconv.Itoa(id), &album); err != nil {
		return nil, err
	}
	return &album, nil
}

// GetAlbumTracks returns tracks on an album.
func (c *Client) GetAlbumTracks(ctx context.Context, albumID int) ([]Track, error) {
	var resp albumTracksResp
	if err := c.get(ctx, "/album/"+strconv.Itoa(albumID)+"/tracks", &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetTrack returns a track by ID.
func (c *Client) GetTrack(ctx context.Context, id int) (*Track, error) {
	var track Track
	if err := c.get(ctx, "/track/"+strconv.Itoa(id), &track); err != nil {
		return nil, err
	}
	return &track, nil
}

// GetGenres returns all genres.
func (c *Client) GetGenres(ctx context.Context) ([]Genre, error) {
	var resp genresResp
	if err := c.get(ctx, "/genre", &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetChart returns chart data (top tracks, albums, artists, playlists).
func (c *Client) GetChart(ctx context.Context) (*Chart, error) {
	var chart Chart
	if err := c.get(ctx, "/chart", &chart); err != nil {
		return nil, err
	}
	return &chart, nil
}
