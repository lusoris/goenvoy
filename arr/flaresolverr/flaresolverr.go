package flaresolverr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 120 * time.Second

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

// Client is a FlareSolverr API client.
type Client struct {
	rawBaseURL string
	httpClient *http.Client
}

// New creates a FlareSolverr [Client] for the instance at baseURL.
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		rawBaseURL: baseURL,
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
		return fmt.Sprintf("flaresolverr: HTTP %d: %s", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("flaresolverr: HTTP %d", e.StatusCode)
}

func (c *Client) do(ctx context.Context, reqBody, v any) error {
	b, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("flaresolverr: encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.rawBaseURL+"/v1", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("flaresolverr: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("flaresolverr: POST /v1: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("flaresolverr: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	if v != nil {
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("flaresolverr: decode response: %w", err)
		}
	}

	return nil
}

func buildRequestBody(cmd string, opts *RequestOptions) map[string]any {
	body := map[string]any{"cmd": cmd}
	if opts == nil {
		return body
	}
	if opts.Session != "" {
		body["session"] = opts.Session
	}
	if opts.MaxTimeout > 0 {
		body["maxTimeout"] = opts.MaxTimeout
	}
	if len(opts.Cookies) > 0 {
		body["cookies"] = opts.Cookies
	}
	if opts.ReturnOnlyCookies {
		body["returnOnlyCookies"] = true
	}
	if opts.Proxy != nil {
		body["proxy"] = opts.Proxy
	}
	if opts.WaitInSeconds > 0 {
		body["waitInSeconds"] = opts.WaitInSeconds
	}
	if opts.DisableMedia {
		body["disableMedia"] = true
	}
	return body
}

// Get solves a Cloudflare challenge for the given URL via a GET request.
func (c *Client) Get(ctx context.Context, url string, opts *RequestOptions) (*Response, error) {
	body := buildRequestBody("request.get", opts)
	body["url"] = url
	var out Response
	if err := c.do(ctx, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Post solves a Cloudflare challenge for the given URL via a POST request.
func (c *Client) Post(ctx context.Context, url, postData string, opts *RequestOptions) (*Response, error) {
	body := buildRequestBody("request.post", opts)
	body["url"] = url
	body["postData"] = postData
	var out Response
	if err := c.do(ctx, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateSession creates a persistent browser session.
func (c *Client) CreateSession(ctx context.Context, session string, proxy *Proxy) (*Response, error) {
	body := map[string]any{"cmd": "sessions.create", "session": session}
	if proxy != nil {
		body["proxy"] = proxy
	}
	var out Response
	if err := c.do(ctx, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListSessions returns all active browser sessions.
func (c *Client) ListSessions(ctx context.Context) (*Response, error) {
	body := map[string]any{"cmd": "sessions.list"}
	var out Response
	if err := c.do(ctx, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DestroySession destroys a persistent browser session.
func (c *Client) DestroySession(ctx context.Context, session string) (*Response, error) {
	body := map[string]any{"cmd": "sessions.destroy", "session": session}
	var out Response
	if err := c.do(ctx, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
