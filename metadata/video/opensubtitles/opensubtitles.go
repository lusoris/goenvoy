package opensubtitles

import (
	"bytes"
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

const defaultBaseURL = "https://api.opensubtitles.com/api/v1"

// Client is an OpenSubtitles API client.
type Client struct {
	*metadata.BaseClient
	apiKey string
	token  string // Bearer token for authenticated endpoints.
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("opensubtitles: %s: %s", e.Status, e.Body)
}

// New creates a new OpenSubtitles client.
//
// The apiKey is required for all API calls. A Bearer token can be
// obtained via [Client.Login] for endpoints that require authentication.
func New(apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "opensubtitles", opts...)
	c := &Client{BaseClient: bc, apiKey: apiKey}
	return c
}

func (c *Client) do(ctx context.Context, method, path string, body, v any) error {
	var bodyReader io.Reader = http.NoBody
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("opensubtitles: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	// Set auth headers via AuthFunc before each request.
	c.SetAuth(func(req *http.Request) {
		req.Header.Set("Api-Key", c.apiKey)
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
	})

	respBody, status, err := c.DoRaw(ctx, method, path, bodyReader)
	if err != nil {
		return err
	}

	if status < 200 || status > 299 {
		return &APIError{StatusCode: status, Status: strconv.Itoa(status), Body: string(respBody)}
	}

	if v != nil {
		if err := json.Unmarshal(respBody, v); err != nil {
			return fmt.Errorf("opensubtitles: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) get(ctx context.Context, path string, params url.Values, v any) error {
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	return c.do(ctx, http.MethodGet, path, nil, v)
}

// Search searches for subtitles matching the given parameters.
func (c *Client) Search(ctx context.Context, p *SearchParams) (*SearchResponse, error) {
	params := url.Values{}
	if p.Query != "" {
		params.Set("query", p.Query)
	}
	if p.IMDBID != 0 {
		params.Set("imdb_id", strconv.Itoa(p.IMDBID))
	}
	if p.TMDBID != 0 {
		params.Set("tmdb_id", strconv.Itoa(p.TMDBID))
	}
	if p.Languages != "" {
		params.Set("languages", p.Languages)
	}
	if p.MovieHash != "" {
		params.Set("moviehash", p.MovieHash)
	}
	if p.Type != "" {
		params.Set("type", p.Type)
	}
	if p.SeasonNumber != 0 {
		params.Set("season_number", strconv.Itoa(p.SeasonNumber))
	}
	if p.EpisodeNumber != 0 {
		params.Set("episode_number", strconv.Itoa(p.EpisodeNumber))
	}
	if p.Year != 0 {
		params.Set("year", strconv.Itoa(p.Year))
	}
	if p.Page != 0 {
		params.Set("page", strconv.Itoa(p.Page))
	}
	if p.OrderBy != "" {
		params.Set("order_by", p.OrderBy)
	}
	if p.OrderDirection != "" {
		params.Set("order_direction", p.OrderDirection)
	}

	var sr SearchResponse
	if err := c.get(ctx, "/subtitles", params, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// Download requests a download link for a subtitle file.
// Requires a Bearer token (set via WithToken or Login).
func (c *Client) Download(ctx context.Context, req DownloadRequest) (*DownloadResponse, error) {
	var dr DownloadResponse
	if err := c.do(ctx, http.MethodPost, "/download", req, &dr); err != nil {
		return nil, err
	}
	return &dr, nil
}

// Login authenticates a user and stores the Bearer token on the client.
func (c *Client) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	var lr LoginResponse
	if err := c.do(ctx, http.MethodPost, "/login", LoginRequest{Username: username, Password: password}, &lr); err != nil {
		return nil, err
	}
	c.token = lr.Token
	return &lr, nil
}

// SearchFeatures searches for movies/TV shows.
func (c *Client) SearchFeatures(ctx context.Context, query string) (*FeaturesResponse, error) {
	params := url.Values{}
	params.Set("query", query)
	var fr FeaturesResponse
	if err := c.get(ctx, "/features", params, &fr); err != nil {
		return nil, err
	}
	return &fr, nil
}

// GetLanguages returns all available subtitle languages.
func (c *Client) GetLanguages(ctx context.Context) ([]Language, error) {
	var resp struct {
		Data []Language `json:"data"`
	}
	if err := c.get(ctx, "/infos/languages", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetFormats returns all available subtitle formats.
func (c *Client) GetFormats(ctx context.Context) ([]SubtitleFormat, error) {
	var resp struct {
		Data []SubtitleFormat `json:"data"`
	}
	if err := c.get(ctx, "/infos/formats", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetUserInfo returns information about the authenticated user.
// Requires a Bearer token.
func (c *Client) GetUserInfo(ctx context.Context) (*UserInfo, error) {
	var resp struct {
		Data *UserInfo `json:"data"`
	}
	if err := c.get(ctx, "/infos/user", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Popular returns popular features.
func (c *Client) Popular(ctx context.Context, languages string) (*FeaturesResponse, error) {
	params := url.Values{}
	if languages != "" {
		params.Set("languages", strings.ToLower(languages))
	}
	var fr FeaturesResponse
	if err := c.get(ctx, "/discover/popular", params, &fr); err != nil {
		return nil, err
	}
	return &fr, nil
}

// Latest returns the latest uploaded subtitles.
func (c *Client) Latest(ctx context.Context, languages string) (*SearchResponse, error) {
	params := url.Values{}
	if languages != "" {
		params.Set("languages", strings.ToLower(languages))
	}
	var sr SearchResponse
	if err := c.get(ctx, "/discover/latest", params, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}
