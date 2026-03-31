package navidrome

import (
	"context"
	"crypto/md5" //nolint:gosec // required by Subsonic API protocol
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	apiVersion = "1.16.1"
	clientName = "goenvoy"
	saltLength = 12
)

// Client is a Navidrome/Subsonic API client.
type Client struct {
	baseURL  string
	username string
	password string
	http     *http.Client
}

// Option configures the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.http = c }
}

// APIError is returned when the API responds with a non-2xx HTTP status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("navidrome: %s: %s", e.Status, e.Body)
}

// New creates a new Navidrome client.
//
// The baseURL should include the protocol and host (e.g. "http://localhost:4533").
// Authentication uses Subsonic token-based auth (md5(password + salt)).
func New(baseURL, username, password string, opts ...Option) *Client {
	c := &Client{
		baseURL:  strings.TrimRight(baseURL, "/"),
		username: username,
		password: password,
		http:     http.DefaultClient,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// authParams returns the common authentication query parameters.
func (c *Client) authParams() (url.Values, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	s := hex.EncodeToString(salt)
	h := md5.Sum([]byte(c.password + s)) //nolint:gosec // required by Subsonic API protocol
	t := hex.EncodeToString(h[:])

	params := url.Values{}
	params.Set("u", c.username)
	params.Set("t", t)
	params.Set("s", s)
	params.Set("v", apiVersion)
	params.Set("c", clientName)
	params.Set("f", "json")
	return params, nil
}

func (c *Client) get(ctx context.Context, endpoint string, extra url.Values) (*responseBody, error) {
	params, err := c.authParams()
	if err != nil {
		return nil, err
	}
	for k, vs := range extra {
		for _, v := range vs {
			params.Add(k, v)
		}
	}

	u := c.baseURL + "/rest/" + endpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	var sr subsonicResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		return nil, err
	}

	if sr.Response.Status != "ok" {
		if sr.Response.Error != nil {
			return nil, sr.Response.Error
		}
		return nil, fmt.Errorf("navidrome: unexpected status %q", sr.Response.Status)
	}

	return &sr.Response, nil
}

// Ping tests connectivity with the server.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.get(ctx, "ping", nil)
	return err
}

// GetArtists returns all artists.
func (c *Client) GetArtists(ctx context.Context) (*ArtistsID3, error) {
	r, err := c.get(ctx, "getArtists", nil)
	if err != nil {
		return nil, err
	}
	return r.Artists, nil
}

// GetArtist returns a single artist by ID including their albums.
func (c *Client) GetArtist(ctx context.Context, id string) (*ArtistID3, error) {
	p := url.Values{}
	p.Set("id", id)
	r, err := c.get(ctx, "getArtist", p)
	if err != nil {
		return nil, err
	}
	return r.Artist, nil
}

// GetAlbum returns a single album by ID including its songs.
func (c *Client) GetAlbum(ctx context.Context, id string) (*AlbumID3, error) {
	p := url.Values{}
	p.Set("id", id)
	r, err := c.get(ctx, "getAlbum", p)
	if err != nil {
		return nil, err
	}
	return r.Album, nil
}

// GetAlbumList2 returns a list of albums by type (e.g. "newest", "random", "recent",
// "frequent", "starred", "alphabeticalByName", "alphabeticalByArtist", "byYear", "byGenre").
func (c *Client) GetAlbumList2(ctx context.Context, listType string, size, offset int) ([]AlbumID3, error) {
	p := url.Values{}
	p.Set("type", listType)
	p.Set("size", strconv.Itoa(size))
	p.Set("offset", strconv.Itoa(offset))
	r, err := c.get(ctx, "getAlbumList2", p)
	if err != nil {
		return nil, err
	}
	if r.AlbumList2 == nil {
		return nil, nil
	}
	return r.AlbumList2.Album, nil
}

// GetSong returns a single song by ID.
func (c *Client) GetSong(ctx context.Context, id string) (*Song, error) {
	p := url.Values{}
	p.Set("id", id)
	r, err := c.get(ctx, "getSong", p)
	if err != nil {
		return nil, err
	}
	return r.Song, nil
}

// GetRandomSongs returns random songs.
func (c *Client) GetRandomSongs(ctx context.Context, size int) ([]Song, error) {
	p := url.Values{}
	p.Set("size", strconv.Itoa(size))
	r, err := c.get(ctx, "getRandomSongs", p)
	if err != nil {
		return nil, err
	}
	if r.RandomSongs == nil {
		return nil, nil
	}
	return r.RandomSongs.Song, nil
}

// GetTopSongs returns top songs for a given artist name.
func (c *Client) GetTopSongs(ctx context.Context, artist string, count int) ([]Song, error) {
	p := url.Values{}
	p.Set("artist", artist)
	p.Set("count", strconv.Itoa(count))
	r, err := c.get(ctx, "getTopSongs", p)
	if err != nil {
		return nil, err
	}
	if r.TopSongs == nil {
		return nil, nil
	}
	return r.TopSongs.Song, nil
}

// Search3 searches for artists, albums, and songs matching a query.
func (c *Client) Search3(ctx context.Context, query string, artistCount, albumCount, songCount int) (*SearchResult3, error) {
	p := url.Values{}
	p.Set("query", query)
	p.Set("artistCount", strconv.Itoa(artistCount))
	p.Set("albumCount", strconv.Itoa(albumCount))
	p.Set("songCount", strconv.Itoa(songCount))
	r, err := c.get(ctx, "search3", p)
	if err != nil {
		return nil, err
	}
	return r.SearchResult, nil
}

// GetPlaylists returns all playlists.
func (c *Client) GetPlaylists(ctx context.Context) ([]Playlist, error) {
	r, err := c.get(ctx, "getPlaylists", nil)
	if err != nil {
		return nil, err
	}
	if r.Playlists == nil {
		return nil, nil
	}
	return r.Playlists.Playlist, nil
}

// GetPlaylist returns a single playlist by ID including its entries.
func (c *Client) GetPlaylist(ctx context.Context, id string) (*Playlist, error) {
	p := url.Values{}
	p.Set("id", id)
	r, err := c.get(ctx, "getPlaylist", p)
	if err != nil {
		return nil, err
	}
	return r.Playlist, nil
}

// GetNowPlaying returns what is currently being played.
func (c *Client) GetNowPlaying(ctx context.Context) ([]NowPlayingEntry, error) {
	r, err := c.get(ctx, "getNowPlaying", nil)
	if err != nil {
		return nil, err
	}
	if r.NowPlaying == nil {
		return nil, nil
	}
	return r.NowPlaying.Entry, nil
}

// GetGenres returns all genres.
func (c *Client) GetGenres(ctx context.Context) ([]Genre, error) {
	r, err := c.get(ctx, "getGenres", nil)
	if err != nil {
		return nil, err
	}
	if r.Genres == nil {
		return nil, nil
	}
	return r.Genres.Genre, nil
}

// GetScanStatus returns the current media library scan status.
func (c *Client) GetScanStatus(ctx context.Context) (*ScanStatus, error) {
	r, err := c.get(ctx, "getScanStatus", nil)
	if err != nil {
		return nil, err
	}
	return r.ScanStatus, nil
}

// StartScan initiates a media library scan.
func (c *Client) StartScan(ctx context.Context) (*ScanStatus, error) {
	r, err := c.get(ctx, "startScan", nil)
	if err != nil {
		return nil, err
	}
	return r.ScanStatus, nil
}

// Scrobble registers the playback of a song.
func (c *Client) Scrobble(ctx context.Context, id string, submission bool) error {
	p := url.Values{}
	p.Set("id", id)
	if submission {
		p.Set("submission", "true")
	} else {
		p.Set("submission", "false")
	}
	_, err := c.get(ctx, "scrobble", p)
	return err
}

// Star stars a song, album, or artist.
func (c *Client) Star(ctx context.Context, id, albumID, artistID string) error {
	p := url.Values{}
	if id != "" {
		p.Set("id", id)
	}
	if albumID != "" {
		p.Set("albumId", albumID)
	}
	if artistID != "" {
		p.Set("artistId", artistID)
	}
	_, err := c.get(ctx, "star", p)
	return err
}

// Unstar removes a star from a song, album, or artist.
func (c *Client) Unstar(ctx context.Context, id, albumID, artistID string) error {
	p := url.Values{}
	if id != "" {
		p.Set("id", id)
	}
	if albumID != "" {
		p.Set("albumId", albumID)
	}
	if artistID != "" {
		p.Set("artistId", artistID)
	}
	_, err := c.get(ctx, "unstar", p)
	return err
}

// GetStarred2 returns all starred artists, albums, and songs.
func (c *Client) GetStarred2(ctx context.Context) (*Starred2, error) {
	r, err := c.get(ctx, "getStarred2", nil)
	if err != nil {
		return nil, err
	}
	return r.Starred2, nil
}
