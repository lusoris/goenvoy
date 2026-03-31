package jellyfin

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

const (
	defaultTimeout    = 30 * time.Second
	defaultClient     = "goenvoy"
	defaultVersion    = "0.0.1"
	defaultDeviceID   = "goenvoy-client"
	defaultDeviceName = "GoEnvoy"
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

// WithAccessToken sets a pre-existing access token, skipping the need to call [Client.AuthenticateByName].
func WithAccessToken(token string) Option {
	return func(cl *Client) { cl.accessToken = token }
}

// WithDeviceID sets the device identifier sent in the Authorization header.
func WithDeviceID(id string) Option {
	return func(cl *Client) { cl.deviceID = id }
}

// Client is a Jellyfin Media Server API client.
type Client struct {
	rawBaseURL  string
	accessToken string
	httpClient  *http.Client
	clientName  string
	version     string
	deviceID    string
	deviceName  string
}

// New creates a Jellyfin [Client] for the server at baseURL.
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		rawBaseURL: baseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
		clientName: defaultClient,
		version:    defaultVersion,
		deviceID:   defaultDeviceID,
		deviceName: defaultDeviceName,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	RawBody    string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("jellyfin: HTTP %d: %s", e.StatusCode, e.Message)
	}
	if e.RawBody != "" {
		return fmt.Sprintf("jellyfin: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("jellyfin: HTTP %d", e.StatusCode)
}

func (c *Client) authHeader() string {
	h := fmt.Sprintf("MediaBrowser Client=%q, Device=%q, DeviceId=%q, Version=%q",
		c.clientName, c.deviceName, c.deviceID, c.version)
	if c.accessToken != "" {
		h += fmt.Sprintf(", Token=%q", c.accessToken)
	}
	return h
}

func (c *Client) doRequest(ctx context.Context, method, path string, body, dst any) error {
	u, err := url.Parse(c.rawBaseURL + path)
	if err != nil {
		return fmt.Errorf("jellyfin: parse URL: %w", err)
	}

	var reqBody io.Reader = http.NoBody
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("jellyfin: marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return fmt.Errorf("jellyfin: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.authHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("jellyfin: %s %s: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("jellyfin: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if jsonErr := json.Unmarshal(respBody, apiErr); jsonErr != nil {
			apiErr.RawBody = string(respBody)
		}
		return apiErr
	}

	if dst != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dst); err != nil {
			return fmt.Errorf("jellyfin: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) get(ctx context.Context, path string, params url.Values) ([]byte, error) {
	if params != nil {
		path += "?" + params.Encode()
	}
	u, err := url.Parse(c.rawBaseURL + path)
	if err != nil {
		return nil, fmt.Errorf("jellyfin: parse URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("jellyfin: create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.authHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jellyfin: GET %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("jellyfin: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if jsonErr := json.Unmarshal(body, apiErr); jsonErr != nil {
			apiErr.RawBody = string(body)
		}
		return nil, apiErr
	}
	return body, nil
}

// Authentication.

// AuthenticateByName authenticates with the server and stores the access token.
func (c *Client) AuthenticateByName(ctx context.Context, username, password string) error {
	var result AuthenticationResult
	err := c.doRequest(ctx, http.MethodPost, "/Users/AuthenticateByName", &authRequest{
		Username: username,
		Pw:       password,
	}, &result)
	if err != nil {
		return err
	}
	c.accessToken = result.AccessToken
	return nil
}

// System endpoints.

// GetPublicSystemInfo returns server info without authentication.
func (c *Client) GetPublicSystemInfo(ctx context.Context) (*PublicSystemInfo, error) {
	data, err := c.get(ctx, "/System/Info/Public", nil)
	if err != nil {
		return nil, err
	}
	var info PublicSystemInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("jellyfin: decode public info: %w", err)
	}
	return &info, nil
}

// GetSystemInfo returns server info (requires authentication).
func (c *Client) GetSystemInfo(ctx context.Context) (*SystemInfo, error) {
	data, err := c.get(ctx, "/System/Info", nil)
	if err != nil {
		return nil, err
	}
	var info SystemInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("jellyfin: decode system info: %w", err)
	}
	return &info, nil
}

// Ping tests server connectivity.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.get(ctx, "/System/Ping", nil)
	return err
}

// Users.

// GetUsers returns all users.
func (c *Client) GetUsers(ctx context.Context) ([]UserDto, error) {
	data, err := c.get(ctx, "/Users", nil)
	if err != nil {
		return nil, err
	}
	var users []UserDto
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("jellyfin: decode users: %w", err)
	}
	return users, nil
}

// GetCurrentUser returns the current authenticated user.
func (c *Client) GetCurrentUser(ctx context.Context) (*UserDto, error) {
	data, err := c.get(ctx, "/Users/Me", nil)
	if err != nil {
		return nil, err
	}
	var user UserDto
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("jellyfin: decode user: %w", err)
	}
	return &user, nil
}

// Sessions.

// GetSessions returns all active sessions.
func (c *Client) GetSessions(ctx context.Context) ([]SessionInfoDto, error) {
	data, err := c.get(ctx, "/Sessions", nil)
	if err != nil {
		return nil, err
	}
	var sessions []SessionInfoDto
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, fmt.Errorf("jellyfin: decode sessions: %w", err)
	}
	return sessions, nil
}

// Items/Libraries.

// GetItems returns items with optional query parameters.
func (c *Client) GetItems(ctx context.Context, params url.Values) (*ItemsResult, error) {
	if params == nil {
		params = url.Values{}
	}
	data, err := c.get(ctx, "/Items", params)
	if err != nil {
		return nil, err
	}
	var result ItemsResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("jellyfin: decode items: %w", err)
	}
	return &result, nil
}

// GetItemsByParent returns items in a specific folder/library.
func (c *Client) GetItemsByParent(ctx context.Context, parentID string, startIndex, limit int) (*ItemsResult, error) {
	p := url.Values{
		"ParentId": {parentID},
	}
	if startIndex >= 0 {
		p.Set("StartIndex", strconv.Itoa(startIndex))
	}
	if limit > 0 {
		p.Set("Limit", strconv.Itoa(limit))
	}
	return c.GetItems(ctx, p)
}

// GetItem returns a specific item by ID.
func (c *Client) GetItem(ctx context.Context, itemID string) (*BaseItemDto, error) {
	data, err := c.get(ctx, "/Items/"+url.PathEscape(itemID), nil)
	if err != nil {
		return nil, err
	}
	var item BaseItemDto
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("jellyfin: decode item: %w", err)
	}
	return &item, nil
}

// GetUserViews returns the user's library views.
func (c *Client) GetUserViews(ctx context.Context, userID string) (*ItemsResult, error) {
	data, err := c.get(ctx, "/UserViews", url.Values{"userId": {userID}})
	if err != nil {
		return nil, err
	}
	var result ItemsResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("jellyfin: decode views: %w", err)
	}
	return &result, nil
}
