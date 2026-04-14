package spotify

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

const defaultBaseURL = "https://api.spotify.com/v1"

// Client is a Spotify Web API client.
type Client struct {
	*metadata.BaseClient
	accessToken string
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("spotify: %s: %s", e.Status, e.Body)
}

// New creates a Spotify [Client] with the given OAuth2 access token.
func New(accessToken string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "spotify", opts...)
	c := &Client{BaseClient: bc, accessToken: accessToken}
	bc.SetAuth(func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	})
	return c
}

func (c *Client) get(ctx context.Context, path string, v any) error {
	u := c.BaseURL() + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("spotify: create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("spotify: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("spotify: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("spotify: decode response: %w", err)
	}
	return nil
}

// Search searches for artists, albums, and/or tracks.
func (c *Client) Search(ctx context.Context, query string, types []string, limit, offset int) (*SearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("type", strings.Join(types, ","))
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))

	var result SearchResult
	if err := c.get(ctx, "/search?"+params.Encode(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetArtist returns an artist by ID.
func (c *Client) GetArtist(ctx context.Context, id string) (*Artist, error) {
	var artist Artist
	if err := c.get(ctx, "/artists/"+url.PathEscape(id), &artist); err != nil {
		return nil, err
	}
	return &artist, nil
}

// GetArtistAlbums returns albums for an artist.
func (c *Client) GetArtistAlbums(ctx context.Context, artistID string, limit, offset int) (*Paged[AlbumSimple], error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))

	var result Paged[AlbumSimple]
	if err := c.get(ctx, "/artists/"+url.PathEscape(artistID)+"/albums?"+params.Encode(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetArtistTopTracks returns top tracks for an artist in a given market.
func (c *Client) GetArtistTopTracks(ctx context.Context, artistID, market string) ([]Track, error) {
	params := url.Values{}
	params.Set("market", market)

	var resp topTracksResp
	if err := c.get(ctx, "/artists/"+url.PathEscape(artistID)+"/top-tracks?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return resp.Tracks, nil
}

// GetRelatedArtists returns artists related to the given artist.
func (c *Client) GetRelatedArtists(ctx context.Context, artistID string) ([]Artist, error) {
	var resp relatedArtistsResp
	if err := c.get(ctx, "/artists/"+url.PathEscape(artistID)+"/related-artists", &resp); err != nil {
		return nil, err
	}
	return resp.Artists, nil
}

// GetAlbum returns an album by ID.
func (c *Client) GetAlbum(ctx context.Context, id string) (*Album, error) {
	var album Album
	if err := c.get(ctx, "/albums/"+url.PathEscape(id), &album); err != nil {
		return nil, err
	}
	return &album, nil
}

// GetAlbumTracks returns tracks on an album.
func (c *Client) GetAlbumTracks(ctx context.Context, albumID string, limit, offset int) (*Paged[TrackSimple], error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))

	var result Paged[TrackSimple]
	if err := c.get(ctx, "/albums/"+url.PathEscape(albumID)+"/tracks?"+params.Encode(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTrack returns a track by ID.
func (c *Client) GetTrack(ctx context.Context, id string) (*Track, error) {
	var track Track
	if err := c.get(ctx, "/tracks/"+url.PathEscape(id), &track); err != nil {
		return nil, err
	}
	return &track, nil
}

// GetAudioFeatures returns audio features for a track.
func (c *Client) GetAudioFeatures(ctx context.Context, trackID string) (*AudioFeatures, error) {
	var features AudioFeatures
	if err := c.get(ctx, "/audio-features/"+url.PathEscape(trackID), &features); err != nil {
		return nil, err
	}
	return &features, nil
}

// GetNewReleases returns new album releases.
func (c *Client) GetNewReleases(ctx context.Context, limit, offset int) (*Paged[AlbumSimple], error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))

	var resp newReleasesResp
	if err := c.get(ctx, "/browse/new-releases?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return &resp.Albums, nil
}

// GetCategories returns browse categories.
func (c *Client) GetCategories(ctx context.Context, limit, offset int) (*Paged[Category], error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))

	var resp categoriesResp
	if err := c.get(ctx, "/browse/categories?"+params.Encode(), &resp); err != nil {
		return nil, err
	}
	return &resp.Categories, nil
}

// GetRecommendations returns track recommendations based on seed artists, genres, and tracks.
func (c *Client) GetRecommendations(ctx context.Context, seeds RecommendationSeeds) (*Recommendations, error) {
	params := url.Values{}
	if len(seeds.SeedArtists) > 0 {
		params.Set("seed_artists", strings.Join(seeds.SeedArtists, ","))
	}
	if len(seeds.SeedGenres) > 0 {
		params.Set("seed_genres", strings.Join(seeds.SeedGenres, ","))
	}
	if len(seeds.SeedTracks) > 0 {
		params.Set("seed_tracks", strings.Join(seeds.SeedTracks, ","))
	}

	var result Recommendations
	if err := c.get(ctx, "/recommendations?"+params.Encode(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
