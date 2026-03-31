package tautulli

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

// Client is a Tautulli API client.
type Client struct {
	rawBaseURL string
	apiKey     string
	httpClient *http.Client
}

// New creates a Tautulli [Client] for the instance at baseURL with the given API key.
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

// APIError is returned when the API responds with a non-2xx status or an error result.
type APIError struct {
	StatusCode int    `json:"-"`
	RawBody    string `json:"-"`
}

func (e *APIError) Error() string {
	if e.RawBody != "" {
		return fmt.Sprintf("tautulli: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("tautulli: HTTP %d", e.StatusCode)
}

// apiResponse wraps the standard Tautulli JSON envelope.
type apiResponse struct {
	Response struct {
		Result  string          `json:"result"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	} `json:"response"`
}

func (c *Client) get(ctx context.Context, cmd string, params url.Values) (json.RawMessage, error) {
	u, err := url.Parse(c.rawBaseURL + "/api/v2")
	if err != nil {
		return nil, fmt.Errorf("tautulli: parse URL: %w", err)
	}

	if params == nil {
		params = url.Values{}
	}
	params.Set("apikey", c.apiKey)
	params.Set("cmd", cmd)
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("tautulli: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tautulli: %s: %w", cmd, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("tautulli: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, RawBody: string(body)}
	}

	var r apiResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("tautulli: decode response: %w", err)
	}

	if r.Response.Result != "success" {
		return nil, &APIError{StatusCode: resp.StatusCode, RawBody: r.Response.Message}
	}

	return r.Response.Data, nil
}

// GetActivity returns current server activity and active sessions.
func (c *Client) GetActivity(ctx context.Context) (*Activity, error) {
	data, err := c.get(ctx, "get_activity", nil)
	if err != nil {
		return nil, err
	}
	var a Activity
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("tautulli: decode activity: %w", err)
	}
	return &a, nil
}

// GetHistory returns the watch history with optional pagination.
func (c *Client) GetHistory(ctx context.Context, start, length int) (*HistoryResponse, error) {
	p := url.Values{}
	if start >= 0 {
		p.Set("start", strconv.Itoa(start))
	}
	if length > 0 {
		p.Set("length", strconv.Itoa(length))
	}
	data, err := c.get(ctx, "get_history", p)
	if err != nil {
		return nil, err
	}
	var h HistoryResponse
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("tautulli: decode history: %w", err)
	}
	return &h, nil
}

// GetLibraries returns all library sections on the server.
func (c *Client) GetLibraries(ctx context.Context) ([]Library, error) {
	data, err := c.get(ctx, "get_libraries", nil)
	if err != nil {
		return nil, err
	}
	var libs []Library
	if err := json.Unmarshal(data, &libs); err != nil {
		return nil, fmt.Errorf("tautulli: decode libraries: %w", err)
	}
	return libs, nil
}

// GetLibrary returns detailed information for a library section.
func (c *Client) GetLibrary(ctx context.Context, sectionID string) (*LibraryDetail, error) {
	p := url.Values{"section_id": {sectionID}}
	data, err := c.get(ctx, "get_library", p)
	if err != nil {
		return nil, err
	}
	var lib LibraryDetail
	if err := json.Unmarshal(data, &lib); err != nil {
		return nil, fmt.Errorf("tautulli: decode library: %w", err)
	}
	return &lib, nil
}

// GetLibraryNames returns minimal library name entries.
func (c *Client) GetLibraryNames(ctx context.Context) ([]LibraryName, error) {
	data, err := c.get(ctx, "get_library_names", nil)
	if err != nil {
		return nil, err
	}
	var names []LibraryName
	if err := json.Unmarshal(data, &names); err != nil {
		return nil, fmt.Errorf("tautulli: decode library names: %w", err)
	}
	return names, nil
}

// GetUsers returns all users with access to the server.
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	data, err := c.get(ctx, "get_users", nil)
	if err != nil {
		return nil, err
	}
	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("tautulli: decode users: %w", err)
	}
	return users, nil
}

// GetUser returns detailed information for a user.
func (c *Client) GetUser(ctx context.Context, userID string) (*UserDetail, error) {
	p := url.Values{"user_id": {userID}}
	data, err := c.get(ctx, "get_user", p)
	if err != nil {
		return nil, err
	}
	var u UserDetail
	if err := json.Unmarshal(data, &u); err != nil {
		return nil, fmt.Errorf("tautulli: decode user: %w", err)
	}
	return &u, nil
}

// GetUserNames returns minimal user name entries.
func (c *Client) GetUserNames(ctx context.Context) ([]UserName, error) {
	data, err := c.get(ctx, "get_user_names", nil)
	if err != nil {
		return nil, err
	}
	var names []UserName
	if err := json.Unmarshal(data, &names); err != nil {
		return nil, fmt.Errorf("tautulli: decode user names: %w", err)
	}
	return names, nil
}

// GetHomeStats returns homepage watch statistics.
func (c *Client) GetHomeStats(ctx context.Context) ([]HomeStats, error) {
	data, err := c.get(ctx, "get_home_stats", nil)
	if err != nil {
		return nil, err
	}
	var stats []HomeStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("tautulli: decode home stats: %w", err)
	}
	return stats, nil
}

// GetServerInfo returns Plex server information.
func (c *Client) GetServerInfo(ctx context.Context) (*ServerInfo, error) {
	data, err := c.get(ctx, "get_server_info", nil)
	if err != nil {
		return nil, err
	}
	var info ServerInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("tautulli: decode server info: %w", err)
	}
	return &info, nil
}

// GetTautulliInfo returns information about the Tautulli installation.
func (c *Client) GetTautulliInfo(ctx context.Context) (*Info, error) {
	data, err := c.get(ctx, "get_tautulli_info", nil)
	if err != nil {
		return nil, err
	}
	var info Info
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("tautulli: decode tautulli info: %w", err)
	}
	return &info, nil
}

// GetMetadata returns metadata for a media item by rating key.
func (c *Client) GetMetadata(ctx context.Context, ratingKey string) (*ItemMetadata, error) {
	p := url.Values{"rating_key": {ratingKey}}
	data, err := c.get(ctx, "get_metadata", p)
	if err != nil {
		return nil, err
	}
	var m ItemMetadata
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("tautulli: decode metadata: %w", err)
	}
	return &m, nil
}

// GetRecentlyAdded returns recently added items.
func (c *Client) GetRecentlyAdded(ctx context.Context, count int) (*RecentlyAddedResponse, error) {
	p := url.Values{"count": {strconv.Itoa(count)}}
	data, err := c.get(ctx, "get_recently_added", p)
	if err != nil {
		return nil, err
	}
	var r RecentlyAddedResponse
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("tautulli: decode recently added: %w", err)
	}
	return &r, nil
}

// Search searches across all libraries.
func (c *Client) Search(ctx context.Context, query string) (*SearchResults, error) {
	p := url.Values{"query": {query}}
	data, err := c.get(ctx, "search", p)
	if err != nil {
		return nil, err
	}
	var s SearchResults
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("tautulli: decode search: %w", err)
	}
	return &s, nil
}

// ServerStatus returns the Plex server connection status.
func (c *Client) ServerStatus(ctx context.Context) (*ServerStatus, error) {
	data, err := c.get(ctx, "server_status", nil)
	if err != nil {
		return nil, err
	}
	var s ServerStatus
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("tautulli: decode server status: %w", err)
	}
	return &s, nil
}

// TerminateSession stops a streaming session.
func (c *Client) TerminateSession(ctx context.Context, sessionKey, message string) error {
	p := url.Values{"session_key": {sessionKey}}
	if message != "" {
		p.Set("message", message)
	}
	_, err := c.get(ctx, "terminate_session", p)
	return err
}

// GetGeoIPLookup returns geolocation info for an IP address.
func (c *Client) GetGeoIPLookup(ctx context.Context, ipAddress string) (*GeoIPInfo, error) {
	p := url.Values{"ip_address": {ipAddress}}
	data, err := c.get(ctx, "get_geoip_lookup", p)
	if err != nil {
		return nil, err
	}
	var g GeoIPInfo
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, fmt.Errorf("tautulli: decode geoip: %w", err)
	}
	return &g, nil
}

// GetUserWatchTimeStats returns watch time statistics for a user.
func (c *Client) GetUserWatchTimeStats(ctx context.Context, userID string) ([]WatchTimeStats, error) {
	p := url.Values{"user_id": {userID}}
	data, err := c.get(ctx, "get_user_watch_time_stats", p)
	if err != nil {
		return nil, err
	}
	var stats []WatchTimeStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("tautulli: decode watch time stats: %w", err)
	}
	return stats, nil
}

// GetLibraryWatchTimeStats returns watch time statistics for a library.
func (c *Client) GetLibraryWatchTimeStats(ctx context.Context, sectionID string) ([]WatchTimeStats, error) {
	p := url.Values{"section_id": {sectionID}}
	data, err := c.get(ctx, "get_library_watch_time_stats", p)
	if err != nil {
		return nil, err
	}
	var stats []WatchTimeStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("tautulli: decode library watch time stats: %w", err)
	}
	return stats, nil
}
