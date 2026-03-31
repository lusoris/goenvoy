package googlebooks

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

const defaultBaseURL = "https://www.googleapis.com/books/v1"

// Client is a Google Books API client.
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// Option configures a [Client].
type Option func(*Client)

// WithHTTPClient sets a custom [http.Client].
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.http = c }
}

// WithBaseURL sets a custom base URL (useful for testing).
func WithBaseURL(u string) Option {
	return func(cl *Client) { cl.baseURL = strings.TrimRight(u, "/") }
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("googlebooks: %s: %s", e.Status, e.Body)
}

// New creates a Google Books [Client] with the given API key.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL: defaultBaseURL,
		apiKey:  apiKey,
		http:    http.DefaultClient,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) get(ctx context.Context, path string, params url.Values, v any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("key", c.apiKey)

	u := c.baseURL + path + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("googlebooks: create request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("googlebooks: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("googlebooks: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	return json.Unmarshal(body, v)
}

// Search searches for volumes by query string.
func (c *Client) Search(ctx context.Context, query string) (*VolumesResponse, error) {
	params := url.Values{}
	params.Set("q", query)
	var resp VolumesResponse
	if err := c.get(ctx, "/volumes", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchWithParams searches for volumes with detailed parameters.
func (c *Client) SearchWithParams(ctx context.Context, p *SearchParams) (*VolumesResponse, error) {
	params := url.Values{}
	if p.Query != "" {
		params.Set("q", p.Query)
	}
	if p.StartIndex > 0 {
		params.Set("startIndex", strconv.Itoa(p.StartIndex))
	}
	if p.MaxResults > 0 {
		params.Set("maxResults", strconv.Itoa(p.MaxResults))
	}
	if p.PrintType != "" {
		params.Set("printType", p.PrintType)
	}
	if p.OrderBy != "" {
		params.Set("orderBy", p.OrderBy)
	}
	if p.Filter != "" {
		params.Set("filter", p.Filter)
	}
	if p.LangRestrict != "" {
		params.Set("langRestrict", p.LangRestrict)
	}
	if p.Projection != "" {
		params.Set("projection", p.Projection)
	}
	var resp VolumesResponse
	if err := c.get(ctx, "/volumes", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVolume returns a single volume by ID.
func (c *Client) GetVolume(ctx context.Context, id string) (*Volume, error) {
	var resp Volume
	if err := c.get(ctx, "/volumes/"+id, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
