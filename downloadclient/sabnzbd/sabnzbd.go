package sabnzbd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultTimeout   = 30 * time.Second
	defaultUserAgent = "goenvoy/0.0.1"
)

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

// WithUserAgent sets the User-Agent header sent with every request.
func WithUserAgent(ua string) Option {
	return func(cl *Client) { cl.userAgent = ua }
}

// Client is a SABnzbd API client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	userAgent  string
}

// New creates a SABnzbd [Client] for the given base URL and API key.
func New(baseURL, apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: defaultTimeout},
		userAgent:  defaultUserAgent,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// APIError is returned when the SABnzbd API returns an error.
type APIError struct {
	Message string
}

func (e *APIError) Error() string {
	return "sabnzbd: " + e.Message
}

func (c *Client) get(ctx context.Context, mode string, params url.Values, dst any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("apikey", c.apiKey)
	params.Set("mode", mode)
	params.Set("output", "json")

	reqURL := c.baseURL + "/api?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("sabnzbd: create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sabnzbd: GET %s: %w", mode, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sabnzbd: HTTP %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("sabnzbd: read response: %w", err)
	}

	// Check for SABnzbd error responses like {"status": false, "error": "..."}
	var errCheck struct {
		Status bool   `json:"status"`
		Error  string `json:"error"`
	}
	if json.Unmarshal(body, &errCheck) == nil && errCheck.Error != "" {
		return &APIError{Message: errCheck.Error}
	}

	if dst != nil {
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("sabnzbd: decode response: %w", err)
		}
	}
	return nil
}

// Queue methods.

// GetQueue returns the current download queue.
func (c *Client) GetQueue(ctx context.Context) (*Queue, error) {
	var wrapper struct {
		Queue Queue `json:"queue"`
	}
	if err := c.get(ctx, "queue", nil, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Queue, nil
}

// AddURL adds an NZB by URL.
func (c *Client) AddURL(ctx context.Context, nzbURL string, params url.Values) ([]string, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("name", nzbURL)
	var result addResult
	if err := c.get(ctx, "addurl", params, &result); err != nil {
		return nil, err
	}
	return result.NZOIDs, nil
}

// Pause pauses the entire download queue.
func (c *Client) Pause(ctx context.Context) error {
	return c.get(ctx, "pause", nil, nil)
}

// Resume resumes the download queue.
func (c *Client) Resume(ctx context.Context) error {
	return c.get(ctx, "resume", nil, nil)
}

// PauseItem pauses a single download by its NZO ID.
func (c *Client) PauseItem(ctx context.Context, nzoID string) error {
	params := url.Values{
		"name":  {"pause"},
		"value": {nzoID},
	}
	return c.get(ctx, "queue", params, nil)
}

// ResumeItem resumes a single download by its NZO ID.
func (c *Client) ResumeItem(ctx context.Context, nzoID string) error {
	params := url.Values{
		"name":  {"resume"},
		"value": {nzoID},
	}
	return c.get(ctx, "queue", params, nil)
}

// DeleteItem deletes a download from the queue by NZO ID.
func (c *Client) DeleteItem(ctx context.Context, nzoID string) error {
	params := url.Values{
		"name":  {"delete"},
		"value": {nzoID},
	}
	return c.get(ctx, "queue", params, nil)
}

// SetSpeedLimit sets the download speed limit.
// Pass 0 to remove the limit. Value is in KiB/s or a percentage (e.g. "50%").
func (c *Client) SetSpeedLimit(ctx context.Context, limit string) error {
	params := url.Values{
		"name":  {"speedlimit"},
		"value": {limit},
	}
	return c.get(ctx, "config", params, nil)
}

// History methods.

// GetHistory returns the download history.
func (c *Client) GetHistory(ctx context.Context, start, limit int) (*HistoryResponse, error) {
	params := url.Values{
		"start": {strconv.Itoa(start)},
		"limit": {strconv.Itoa(limit)},
	}
	var wrapper struct {
		History HistoryResponse `json:"history"`
	}
	if err := c.get(ctx, "history", params, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.History, nil
}

// DeleteHistory deletes an item from history by NZO ID.
func (c *Client) DeleteHistory(ctx context.Context, nzoID string) error {
	params := url.Values{
		"name":  {"delete"},
		"value": {nzoID},
	}
	return c.get(ctx, "history", params, nil)
}

// Server methods.

// GetVersion returns the SABnzbd version.
func (c *Client) GetVersion(ctx context.Context) (string, error) {
	var result struct {
		Version string `json:"version"`
	}
	if err := c.get(ctx, "version", nil, &result); err != nil {
		return "", err
	}
	return result.Version, nil
}

// GetServerStats returns download statistics.
func (c *Client) GetServerStats(ctx context.Context) (*ServerStats, error) {
	var result ServerStats
	if err := c.get(ctx, "server_stats", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetWarnings returns the list of warning messages.
func (c *Client) GetWarnings(ctx context.Context) ([]string, error) {
	var result struct {
		Warnings []string `json:"warnings"`
	}
	if err := c.get(ctx, "warnings", nil, &result); err != nil {
		return nil, err
	}
	return result.Warnings, nil
}

// SetCategory sets the category for a download.
func (c *Client) SetCategory(ctx context.Context, nzoID, category string) error {
	params := url.Values{
		"name":   {"cat"},
		"value":  {nzoID},
		"value2": {category},
	}
	return c.get(ctx, "change_cat", params, nil)
}
