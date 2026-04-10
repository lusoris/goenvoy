package lastfm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://ws.audioscrobbler.com/2.0/"

// Client is a Last.fm API client.
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
	return fmt.Sprintf("lastfm: %s: %s", e.Status, e.Body)
}

// Error represents an error returned by the Last.fm API.
type Error struct {
	Code    int    `json:"error"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("lastfm: error %d: %s", e.Code, e.Message)
}

// New creates a new Last.fm client.
func New(apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "lastfm", opts...)
	return &Client{BaseClient: bc, apiKey: apiKey}
}


func (c *Client) get(ctx context.Context, method string, extra url.Values, v any) error {
	params := url.Values{}
	params.Set("method", method)
	params.Set("api_key", c.apiKey)
	params.Set("format", "json")
	for k, vs := range extra {
		for _, val := range vs {
			params.Add(k, val)
		}
	}

	u := c.BaseURL() + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	// Check for Last.fm API error in JSON body.
	var lfmErr Error
	if json.Unmarshal(body, &lfmErr) == nil && lfmErr.Code != 0 {
		return &lfmErr
	}

	return json.Unmarshal(body, v)
}

// GetArtistInfo returns metadata for an artist.
func (c *Client) GetArtistInfo(ctx context.Context, artist string) (*Artist, error) {
	p := url.Values{}
	p.Set("artist", artist)
	p.Set("autocorrect", "1")

	var resp struct {
		Artist *Artist `json:"artist"`
	}
	if err := c.get(ctx, "artist.getinfo", p, &resp); err != nil {
		return nil, err
	}
	return resp.Artist, nil
}

// GetAlbumInfo returns metadata for an album.
func (c *Client) GetAlbumInfo(ctx context.Context, artist, album string) (*Album, error) {
	p := url.Values{}
	p.Set("artist", artist)
	p.Set("album", album)
	p.Set("autocorrect", "1")

	var resp struct {
		Album *Album `json:"album"`
	}
	if err := c.get(ctx, "album.getinfo", p, &resp); err != nil {
		return nil, err
	}
	return resp.Album, nil
}

// GetTrackInfo returns metadata for a track.
func (c *Client) GetTrackInfo(ctx context.Context, artist, track string) (*Track, error) {
	p := url.Values{}
	p.Set("artist", artist)
	p.Set("track", track)
	p.Set("autocorrect", "1")

	var resp struct {
		Track *Track `json:"track"`
	}
	if err := c.get(ctx, "track.getinfo", p, &resp); err != nil {
		return nil, err
	}
	return resp.Track, nil
}

// GetSimilarArtists returns artists similar to the given artist.
func (c *Client) GetSimilarArtists(ctx context.Context, artist string, limit int) ([]Artist, error) {
	p := url.Values{}
	p.Set("artist", artist)
	p.Set("autocorrect", "1")
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}

	var resp struct {
		SimilarArtists struct {
			Artist []Artist `json:"artist"`
		} `json:"similarartists"`
	}
	if err := c.get(ctx, "artist.getsimilar", p, &resp); err != nil {
		return nil, err
	}
	return resp.SimilarArtists.Artist, nil
}

// GetTopAlbums returns an artist's top albums.
func (c *Client) GetTopAlbums(ctx context.Context, artist string, limit int) ([]Album, error) {
	p := url.Values{}
	p.Set("artist", artist)
	p.Set("autocorrect", "1")
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}

	var resp struct {
		TopAlbums TopAlbums `json:"topalbums"`
	}
	if err := c.get(ctx, "artist.gettopalbums", p, &resp); err != nil {
		return nil, err
	}
	return resp.TopAlbums.Album, nil
}

// GetTopTracks returns an artist's top tracks.
func (c *Client) GetTopTracks(ctx context.Context, artist string, limit int) ([]Track, error) {
	p := url.Values{}
	p.Set("artist", artist)
	p.Set("autocorrect", "1")
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}

	var resp struct {
		TopTracks TopTracks `json:"toptracks"`
	}
	if err := c.get(ctx, "artist.gettoptracks", p, &resp); err != nil {
		return nil, err
	}
	return resp.TopTracks.Track, nil
}

// GetChartTopArtists returns the top artists chart.
func (c *Client) GetChartTopArtists(ctx context.Context, limit int) ([]Artist, error) {
	p := url.Values{}
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}

	var resp struct {
		Artists ChartArtists `json:"artists"`
	}
	if err := c.get(ctx, "chart.gettopartists", p, &resp); err != nil {
		return nil, err
	}
	return resp.Artists.Artist, nil
}

// GetChartTopTracks returns the top tracks chart.
func (c *Client) GetChartTopTracks(ctx context.Context, limit int) ([]Track, error) {
	p := url.Values{}
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}

	var resp struct {
		Tracks ChartTracks `json:"tracks"`
	}
	if err := c.get(ctx, "chart.gettoptracks", p, &resp); err != nil {
		return nil, err
	}
	return resp.Tracks.Track, nil
}

// GetTopTags returns the top global tags.
func (c *Client) GetTopTags(ctx context.Context) ([]Tag, error) {
	var resp struct {
		TopTags TopTags `json:"toptags"`
	}
	if err := c.get(ctx, "chart.gettoptags", nil, &resp); err != nil {
		return nil, err
	}
	return resp.TopTags.Tag, nil
}

// SearchArtist searches for artists by name.
func (c *Client) SearchArtist(ctx context.Context, artist string, limit int) ([]Artist, error) {
	p := url.Values{}
	p.Set("artist", artist)
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}

	var resp struct {
		Results struct {
			ArtistMatches struct {
				Artist []Artist `json:"artist"`
			} `json:"artistmatches"`
		} `json:"results"`
	}
	if err := c.get(ctx, "artist.search", p, &resp); err != nil {
		return nil, err
	}
	return resp.Results.ArtistMatches.Artist, nil
}

// SearchAlbum searches for albums by name.
func (c *Client) SearchAlbum(ctx context.Context, album string, limit int) ([]Album, error) {
	p := url.Values{}
	p.Set("album", album)
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}

	var resp struct {
		Results struct {
			AlbumMatches struct {
				Album []Album `json:"album"`
			} `json:"albummatches"`
		} `json:"results"`
	}
	if err := c.get(ctx, "album.search", p, &resp); err != nil {
		return nil, err
	}
	return resp.Results.AlbumMatches.Album, nil
}

// SearchTrack searches for tracks by name.
func (c *Client) SearchTrack(ctx context.Context, track string, limit int) ([]Track, error) {
	p := url.Values{}
	p.Set("track", track)
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}

	var resp struct {
		Results struct {
			TrackMatches struct {
				Track []Track `json:"track"`
			} `json:"trackmatches"`
		} `json:"results"`
	}
	if err := c.get(ctx, "track.search", p, &resp); err != nil {
		return nil, err
	}
	return resp.Results.TrackMatches.Track, nil
}
