package audiodb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/lusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://www.theaudiodb.com/api/v1/json"

// Client is a TheAudioDB API client.
type Client struct {
	*metadata.BaseClient
	apiKey string
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("audiodb: %s: %s", e.Status, e.Body)
}

// New creates a TheAudioDB [Client] with the given API key.
func New(apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "audiodb", opts...)
	return &Client{BaseClient: bc, apiKey: apiKey}
}


// Response wrappers matching the API JSON structure.
type artistsResp struct {
	Artists []Artist `json:"artists"`
}

type albumsResp struct {
	Album []Album `json:"album"`
}

type tracksResp struct {
	Track []Track `json:"track"`
}

type mvidsResp struct {
	Mvids []MusicVideo `json:"mvids"`
}

type discographyResp struct {
	Album []Discography `json:"album"`
}

type trendingResp struct {
	Trending []Trending `json:"trending"`
}

func (c *Client) get(ctx context.Context, endpoint string, v any) error {
	u := c.BaseURL() + "/" + c.apiKey + "/" + endpoint

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("audiodb: create request: %w", err)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("audiodb: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("audiodb: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	return json.Unmarshal(body, v)
}

// SearchArtist searches for artists by name.
func (c *Client) SearchArtist(ctx context.Context, name string) ([]Artist, error) {
	var resp artistsResp
	if err := c.get(ctx, "search.php?s="+url.QueryEscape(name), &resp); err != nil {
		return nil, err
	}
	return resp.Artists, nil
}

// SearchAlbum searches for albums by artist name.
func (c *Client) SearchAlbum(ctx context.Context, artist string) ([]Album, error) {
	var resp albumsResp
	if err := c.get(ctx, "searchalbum.php?s="+url.QueryEscape(artist), &resp); err != nil {
		return nil, err
	}
	return resp.Album, nil
}

// SearchTrack searches for a track by artist and track name.
func (c *Client) SearchTrack(ctx context.Context, artist, track string) ([]Track, error) {
	var resp tracksResp
	if err := c.get(ctx, "searchtrack.php?s="+url.QueryEscape(artist)+"&t="+url.QueryEscape(track), &resp); err != nil {
		return nil, err
	}
	return resp.Track, nil
}

// GetArtist returns an artist by ID.
func (c *Client) GetArtist(ctx context.Context, id string) (*Artist, error) {
	var resp artistsResp
	if err := c.get(ctx, "artist.php?i="+url.QueryEscape(id), &resp); err != nil {
		return nil, err
	}
	if len(resp.Artists) == 0 {
		return nil, nil
	}
	return &resp.Artists[0], nil
}

// GetAlbum returns an album by ID.
func (c *Client) GetAlbum(ctx context.Context, id string) (*Album, error) {
	var resp albumsResp
	if err := c.get(ctx, "album.php?m="+url.QueryEscape(id), &resp); err != nil {
		return nil, err
	}
	if len(resp.Album) == 0 {
		return nil, nil
	}
	return &resp.Album[0], nil
}

// GetAlbumsByArtist returns all albums by an artist ID.
func (c *Client) GetAlbumsByArtist(ctx context.Context, artistID string) ([]Album, error) {
	var resp albumsResp
	if err := c.get(ctx, "album.php?i="+url.QueryEscape(artistID), &resp); err != nil {
		return nil, err
	}
	return resp.Album, nil
}

// GetTracksByAlbum returns all tracks on an album.
func (c *Client) GetTracksByAlbum(ctx context.Context, albumID string) ([]Track, error) {
	var resp tracksResp
	if err := c.get(ctx, "track.php?m="+url.QueryEscape(albumID), &resp); err != nil {
		return nil, err
	}
	return resp.Track, nil
}

// GetMusicVideos returns music videos for an artist.
func (c *Client) GetMusicVideos(ctx context.Context, artistID string) ([]MusicVideo, error) {
	var resp mvidsResp
	if err := c.get(ctx, "mvid.php?i="+url.QueryEscape(artistID), &resp); err != nil {
		return nil, err
	}
	return resp.Mvids, nil
}

// GetDiscography returns the discography for an artist.
func (c *Client) GetDiscography(ctx context.Context, artist string) ([]Discography, error) {
	var resp discographyResp
	if err := c.get(ctx, "discography.php?s="+url.QueryEscape(artist), &resp); err != nil {
		return nil, err
	}
	return resp.Album, nil
}

// GetTopTracks returns the top 10 tracks for an artist.
func (c *Client) GetTopTracks(ctx context.Context, artist string) ([]Track, error) {
	var resp tracksResp
	if err := c.get(ctx, "track-top10.php?s="+url.QueryEscape(artist), &resp); err != nil {
		return nil, err
	}
	return resp.Track, nil
}

// GetTrending returns trending music for a country and format.
func (c *Client) GetTrending(ctx context.Context, country, format string) ([]Trending, error) {
	var resp trendingResp
	if err := c.get(ctx, "trending.php?country="+url.QueryEscape(country)+"&type=itunes&format="+url.QueryEscape(format), &resp); err != nil {
		return nil, err
	}
	return resp.Trending, nil
}
