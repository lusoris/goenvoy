package discogs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const defaultBaseURL = "https://api.discogs.com"

// Client is a Discogs API client.
type Client struct {
	baseURL   string
	token     string
	userAgent string
	http      *http.Client
}

// Option configures the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.http = c }
}

// WithBaseURL sets a custom base URL (useful for testing).
func WithBaseURL(u string) Option {
	return func(cl *Client) { cl.baseURL = strings.TrimRight(u, "/") }
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(ua string) Option {
	return func(cl *Client) { cl.userAgent = ua }
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("discogs: %s: %s", e.Status, e.Body)
}

// New creates a new Discogs client.
//
// The token is a Discogs personal access token used for authentication.
func New(token string, opts ...Option) *Client {
	c := &Client{
		baseURL:   defaultBaseURL,
		token:     token,
		userAgent: "goenvoy/1.0",
		http:      http.DefaultClient,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) get(ctx context.Context, path string, params url.Values, v any) error {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Discogs token="+c.token)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
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

	return json.Unmarshal(body, v)
}

// GetRelease returns a release by ID.
func (c *Client) GetRelease(ctx context.Context, id int) (*Release, error) {
	var r Release
	if err := c.get(ctx, fmt.Sprintf("/releases/%d", id), nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetArtist returns an artist by ID.
func (c *Client) GetArtist(ctx context.Context, id int) (*Artist, error) {
	var a Artist
	if err := c.get(ctx, fmt.Sprintf("/artists/%d", id), nil, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// GetArtistReleases returns releases for an artist.
func (c *Client) GetArtistReleases(ctx context.Context, id, page, perPage int) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))

	var sr SearchResponse
	if err := c.get(ctx, fmt.Sprintf("/artists/%d/releases", id), params, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// GetLabel returns a label by ID.
func (c *Client) GetLabel(ctx context.Context, id int) (*Label, error) {
	var l Label
	if err := c.get(ctx, fmt.Sprintf("/labels/%d", id), nil, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// GetLabelReleases returns releases for a label.
func (c *Client) GetLabelReleases(ctx context.Context, id, page, perPage int) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))

	var sr SearchResponse
	if err := c.get(ctx, fmt.Sprintf("/labels/%d/releases", id), params, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// GetMasterRelease returns a master release by ID.
func (c *Client) GetMasterRelease(ctx context.Context, id int) (*MasterRelease, error) {
	var m MasterRelease
	if err := c.get(ctx, fmt.Sprintf("/masters/%d", id), nil, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMasterVersions returns all versions of a master release.
func (c *Client) GetMasterVersions(ctx context.Context, id, page, perPage int) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))

	var sr SearchResponse
	if err := c.get(ctx, fmt.Sprintf("/masters/%d/versions", id), params, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// Search performs a database search.
func (c *Client) Search(ctx context.Context, query, searchType string, page, perPage int) (*SearchResponse, error) {
	params := url.Values{}
	if query != "" {
		params.Set("q", query)
	}
	if searchType != "" {
		params.Set("type", searchType)
	}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))

	var sr SearchResponse
	if err := c.get(ctx, "/database/search", params, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}
