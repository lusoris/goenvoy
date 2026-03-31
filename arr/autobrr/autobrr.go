package autobrr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

// Client is an autobrr API client.
type Client struct {
	rawBaseURL string
	apiKey     string
	httpClient *http.Client
}

// New creates an autobrr [Client] for the instance at baseURL with the given API key.
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
	RawBody    string `json:"-"`
}

func (e *APIError) Error() string {
	if e.RawBody != "" {
		return fmt.Sprintf("autobrr: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("autobrr: HTTP %d", e.StatusCode)
}

func (c *Client) do(ctx context.Context, method, path string, reqBody any) ([]byte, error) {
	u, err := url.Parse(c.rawBaseURL + "/api" + path)
	if err != nil {
		return nil, fmt.Errorf("autobrr: parse URL: %w", err)
	}

	var body io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("autobrr: encode request: %w", err)
		}
		body = bytes.NewReader(b)
	} else {
		body = http.NoBody
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("autobrr: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Token", c.apiKey)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("autobrr: %s %s: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("autobrr: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, RawBody: string(data)}
	}

	return data, nil
}

func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	return c.do(ctx, http.MethodGet, path, nil)
}

// Health checks.

// Liveness checks if the application is running.
func (c *Client) Liveness(ctx context.Context) error {
	_, err := c.get(ctx, "/healthz/liveness")
	return err
}

// Readiness checks if the application and dependencies are ready.
func (c *Client) Readiness(ctx context.Context) error {
	_, err := c.get(ctx, "/healthz/readiness")
	return err
}

// Filters.

// GetFilters returns all filters.
func (c *Client) GetFilters(ctx context.Context) ([]Filter, error) {
	data, err := c.get(ctx, "/filters")
	if err != nil {
		return nil, err
	}
	var filters []Filter
	if err := json.Unmarshal(data, &filters); err != nil {
		return nil, fmt.Errorf("autobrr: decode filters: %w", err)
	}
	return filters, nil
}

// SetFilterEnabled enables or disables a filter.
func (c *Client) SetFilterEnabled(ctx context.Context, filterID int, enabled bool) error {
	_, err := c.do(ctx, http.MethodPut, "/filters/"+strconv.Itoa(filterID)+"/enabled", map[string]bool{"enabled": enabled})
	return err
}

// Indexers.

// GetIndexers returns all indexers.
func (c *Client) GetIndexers(ctx context.Context) ([]Indexer, error) {
	data, err := c.get(ctx, "/indexer")
	if err != nil {
		return nil, err
	}
	var indexers []Indexer
	if err := json.Unmarshal(data, &indexers); err != nil {
		return nil, fmt.Errorf("autobrr: decode indexers: %w", err)
	}
	return indexers, nil
}

// SetIndexerEnabled enables or disables an indexer.
func (c *Client) SetIndexerEnabled(ctx context.Context, indexerID int, enabled bool) error {
	_, err := c.do(ctx, http.MethodPatch, "/indexer/"+strconv.Itoa(indexerID)+"/enabled", map[string]bool{"enabled": enabled})
	return err
}

// IRC networks.

// GetIRCNetworks returns all IRC networks.
func (c *Client) GetIRCNetworks(ctx context.Context) ([]IRCNetwork, error) {
	data, err := c.get(ctx, "/irc")
	if err != nil {
		return nil, err
	}
	var networks []IRCNetwork
	if err := json.Unmarshal(data, &networks); err != nil {
		return nil, fmt.Errorf("autobrr: decode IRC networks: %w", err)
	}
	return networks, nil
}

// RestartIRCNetwork restarts a specific IRC network.
func (c *Client) RestartIRCNetwork(ctx context.Context, networkID int) error {
	_, err := c.get(ctx, "/irc/network/"+strconv.Itoa(networkID)+"/restart")
	return err
}

// Feeds.

// GetFeeds returns all feeds.
func (c *Client) GetFeeds(ctx context.Context) ([]Feed, error) {
	data, err := c.get(ctx, "/feeds")
	if err != nil {
		return nil, err
	}
	var feeds []Feed
	if err := json.Unmarshal(data, &feeds); err != nil {
		return nil, fmt.Errorf("autobrr: decode feeds: %w", err)
	}
	return feeds, nil
}

// SetFeedEnabled enables or disables a feed.
func (c *Client) SetFeedEnabled(ctx context.Context, feedID int, enabled bool) error {
	_, err := c.do(ctx, http.MethodPatch, "/feeds/"+strconv.Itoa(feedID)+"/enabled", map[string]bool{"enabled": enabled})
	return err
}

// Download clients.

// GetDownloadClients returns all download clients.
func (c *Client) GetDownloadClients(ctx context.Context) ([]DownloadClient, error) {
	data, err := c.get(ctx, "/download_clients")
	if err != nil {
		return nil, err
	}
	var clients []DownloadClient
	if err := json.Unmarshal(data, &clients); err != nil {
		return nil, fmt.Errorf("autobrr: decode download clients: %w", err)
	}
	return clients, nil
}

// Notifications.

// GetNotifications returns all notification agents.
func (c *Client) GetNotifications(ctx context.Context) ([]Notification, error) {
	data, err := c.get(ctx, "/notification")
	if err != nil {
		return nil, err
	}
	var notifs []Notification
	if err := json.Unmarshal(data, &notifs); err != nil {
		return nil, fmt.Errorf("autobrr: decode notifications: %w", err)
	}
	return notifs, nil
}

// Config.

// GetConfig returns the autobrr configuration.
func (c *Client) GetConfig(ctx context.Context) (*Config, error) {
	data, err := c.get(ctx, "/config")
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("autobrr: decode config: %w", err)
	}
	return &cfg, nil
}

// API keys.

// GetAPIKeys returns all API keys.
func (c *Client) GetAPIKeys(ctx context.Context) ([]APIKey, error) {
	data, err := c.get(ctx, "/keys")
	if err != nil {
		return nil, err
	}
	var keys []APIKey
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("autobrr: decode API keys: %w", err)
	}
	return keys, nil
}
