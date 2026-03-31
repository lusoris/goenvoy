package gotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const defaultTimeout = 30 * time.Second

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int    `json:"-"`
	Status     string `json:"-"`
	Body       string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("gotify: HTTP %d: %s", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("gotify: HTTP %d", e.StatusCode)
}

// Client is a Gotify REST API client.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// Option configures a [Client].
type Option func(*Client)

// WithHTTPClient sets a custom [http.Client].
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.httpClient = c }
}

// New creates a Gotify [Client] for the instance at baseURL with the given token.
func New(baseURL, token string, opts ...Option) *Client {
	c := &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) do(ctx context.Context, method, path string, body, v any) error {
	var reqBody io.Reader = http.NoBody
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("gotify: marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("gotify: create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Gotify-Key", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gotify: %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("gotify: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(raw)}
	}

	if v != nil {
		if err := json.Unmarshal(raw, v); err != nil {
			return fmt.Errorf("gotify: decode response: %w", err)
		}
	}
	return nil
}

// CreateMessage sends a push notification message.
func (c *Client) CreateMessage(ctx context.Context, title, message string, priority int) (*Message, error) {
	body := map[string]any{
		"title":    title,
		"message":  message,
		"priority": priority,
	}
	var out Message
	if err := c.do(ctx, http.MethodPost, "/message", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMessages returns a paginated list of messages.
func (c *Client) GetMessages(ctx context.Context) (*PagedMessages, error) {
	var out PagedMessages
	if err := c.do(ctx, http.MethodGet, "/message", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteMessages deletes all messages.
func (c *Client) DeleteMessages(ctx context.Context) error {
	return c.do(ctx, http.MethodDelete, "/message", nil, nil)
}

// DeleteMessage deletes a single message by ID.
func (c *Client) DeleteMessage(ctx context.Context, id int) error {
	return c.do(ctx, http.MethodDelete, "/message/"+strconv.Itoa(id), nil, nil)
}

// GetApplications returns all applications.
func (c *Client) GetApplications(ctx context.Context) ([]Application, error) {
	var out []Application
	if err := c.do(ctx, http.MethodGet, "/application", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateApplication creates a new application.
func (c *Client) CreateApplication(ctx context.Context, name, description string) (*Application, error) {
	body := map[string]string{"name": name, "description": description}
	var out Application
	if err := c.do(ctx, http.MethodPost, "/application", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteApplication deletes an application by ID.
func (c *Client) DeleteApplication(ctx context.Context, id int) error {
	return c.do(ctx, http.MethodDelete, "/application/"+strconv.Itoa(id), nil, nil)
}

// GetClients returns all clients.
func (c *Client) GetClients(ctx context.Context) ([]ClientInfo, error) {
	var out []ClientInfo
	if err := c.do(ctx, http.MethodGet, "/client", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateClient creates a new client.
func (c *Client) CreateClient(ctx context.Context, name string) (*ClientInfo, error) {
	body := map[string]string{"name": name}
	var out ClientInfo
	if err := c.do(ctx, http.MethodPost, "/client", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteClient deletes a client by ID.
func (c *Client) DeleteClient(ctx context.Context, id int) error {
	return c.do(ctx, http.MethodDelete, "/client/"+strconv.Itoa(id), nil, nil)
}

// GetCurrentUser returns the current authenticated user.
func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	var out User
	if err := c.do(ctx, http.MethodGet, "/current/user", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetHealth returns the server health status.
func (c *Client) GetHealth(ctx context.Context) (*Health, error) {
	var out Health
	if err := c.do(ctx, http.MethodGet, "/health", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetVersion returns the server version information.
func (c *Client) GetVersion(ctx context.Context) (*VersionInfo, error) {
	var out VersionInfo
	if err := c.do(ctx, http.MethodGet, "/version", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
