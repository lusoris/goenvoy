package mylar

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultTimeout = 30 * time.Second

// Option configures a [Client].
type Option func(*Client)

// WithHTTPClient sets a custom [http.Client].
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.httpClient = c }
}

// WithTimeout overrides the default HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(cl *Client) { cl.httpClient.Timeout = d }
}

// Client is a Mylar3 API client.
type Client struct {
	rawBaseURL string
	apiKey     string
	httpClient *http.Client
}

// New creates a Mylar3 [Client] for the instance at baseURL with the given API key.
func New(baseURL, apiKey string, opts ...Option) *Client {
	c := &Client{
		rawBaseURL: baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int    `json:"-"`
	Status     string `json:"-"`
	Body       string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("mylar: HTTP %d: %s", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("mylar: HTTP %d", e.StatusCode)
}

func (c *Client) get(ctx context.Context, cmd string, params url.Values, v any) error {
	u := c.rawBaseURL + "/api?apikey=" + c.apiKey + "&cmd=" + cmd
	if params != nil {
		u += "&" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("mylar: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("mylar: %s: %w", cmd, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("mylar: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	if v != nil {
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("mylar: decode %s: %w", cmd, err)
		}
	}

	return nil
}

// GetIndex returns all comic series.
func (c *Client) GetIndex(ctx context.Context) ([]Comic, error) {
	var out []Comic
	if err := c.get(ctx, "getIndex", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetComic returns a single comic series by ID.
func (c *Client) GetComic(ctx context.Context, id string) (*Comic, error) {
	var out Comic
	if err := c.get(ctx, "getComic", url.Values{"id": {id}}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUpcoming returns upcoming comic issues.
func (c *Client) GetUpcoming(ctx context.Context) ([]Upcoming, error) {
	var out []Upcoming
	if err := c.get(ctx, "getUpcoming", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetWanted returns all wanted issues.
func (c *Client) GetWanted(ctx context.Context) ([]WantedIssue, error) {
	var out []WantedIssue
	if err := c.get(ctx, "getWanted", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetHistory returns the download history.
func (c *Client) GetHistory(ctx context.Context) ([]HistoryEntry, error) {
	var out []HistoryEntry
	if err := c.get(ctx, "getHistory", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetLogs returns recent log entries.
func (c *Client) GetLogs(ctx context.Context) ([]LogEntry, error) {
	var out []LogEntry
	if err := c.get(ctx, "getLogs", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FindComic searches for a comic by name.
func (c *Client) FindComic(ctx context.Context, name string) ([]SearchResult, error) {
	var out []SearchResult
	if err := c.get(ctx, "findComic", url.Values{"name": {name}}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddComic adds a comic series by ID.
func (c *Client) AddComic(ctx context.Context, id string) error {
	return c.get(ctx, "addComic", url.Values{"id": {id}}, nil)
}

// DeleteComic removes a comic series by ID.
func (c *Client) DeleteComic(ctx context.Context, id string) error {
	return c.get(ctx, "delComic", url.Values{"id": {id}}, nil)
}

// PauseComic pauses monitoring of a comic series.
func (c *Client) PauseComic(ctx context.Context, id string) error {
	return c.get(ctx, "pauseComic", url.Values{"id": {id}}, nil)
}

// ResumeComic resumes monitoring of a comic series.
func (c *Client) ResumeComic(ctx context.Context, id string) error {
	return c.get(ctx, "resumeComic", url.Values{"id": {id}}, nil)
}

// RefreshComic refreshes metadata for a comic series.
func (c *Client) RefreshComic(ctx context.Context, id string) error {
	return c.get(ctx, "refreshComic", url.Values{"id": {id}}, nil)
}

// ForceSearch forces an issue search for a comic series.
func (c *Client) ForceSearch(ctx context.Context, id string) error {
	return c.get(ctx, "forceSearch", url.Values{"id": {id}}, nil)
}

// GetVersion returns Mylar3 version information.
func (c *Client) GetVersion(ctx context.Context) (*VersionInfo, error) {
	var out VersionInfo
	if err := c.get(ctx, "getVersion", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetReadList returns all read lists.
func (c *Client) GetReadList(ctx context.Context) ([]ReadList, error) {
	var out []ReadList
	if err := c.get(ctx, "getReadList", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetStoryArc returns a story arc by ID.
func (c *Client) GetStoryArc(ctx context.Context, id string) (*StoryArc, error) {
	var out StoryArc
	if err := c.get(ctx, "getStoryArc", url.Values{"id": {id}}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
