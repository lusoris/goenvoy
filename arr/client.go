package arr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout   = 30 * time.Second
	defaultUserAgent = "goenvoy/0.0.1"
)

// Option configures a [BaseClient].
type Option func(*BaseClient)

// WithHTTPClient sets a custom [http.Client] for the [BaseClient].
func WithHTTPClient(c *http.Client) Option {
	return func(b *BaseClient) { b.httpClient = c }
}

// WithTimeout overrides the default HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(b *BaseClient) { b.httpClient.Timeout = d }
}

// WithUserAgent sets the User-Agent header sent with every request.
func WithUserAgent(ua string) Option {
	return func(b *BaseClient) { b.userAgent = ua }
}

// BaseClient is a low-level HTTP client shared by all *arr service clients.
// It handles authentication, JSON marshaling, and error wrapping.
type BaseClient struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client
	userAgent  string
}

// NewBaseClient creates a [BaseClient] targeting the given base URL.
// The apiKey is sent in the X-Api-Key header on every request.
func NewBaseClient(baseURL, apiKey string, opts ...Option) (*BaseClient, error) {
	baseURL = strings.TrimRight(baseURL, "/")

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("arr: invalid base URL %q: %w", baseURL, err)
	}

	c := &BaseClient{
		baseURL:    u,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: defaultTimeout, Transport: http.DefaultTransport.(*http.Transport).Clone()},
		userAgent:  defaultUserAgent,
	}

	for _, o := range opts {
		o(c)
	}

	return c, nil
}

// Get performs an authenticated GET request and decodes the JSON response into dst.
func (c *BaseClient) Get(ctx context.Context, path string, dst any) error {
	return c.do(ctx, http.MethodGet, path, nil, dst)
}

// Post performs an authenticated POST request with a JSON body and decodes the response into dst.
func (c *BaseClient) Post(ctx context.Context, path string, body, dst any) error {
	return c.do(ctx, http.MethodPost, path, body, dst)
}

// Put performs an authenticated PUT request with a JSON body and decodes the response into dst.
func (c *BaseClient) Put(ctx context.Context, path string, body, dst any) error {
	return c.do(ctx, http.MethodPut, path, body, dst)
}

// Delete performs an authenticated DELETE request with an optional JSON body and decodes the response into dst.
func (c *BaseClient) Delete(ctx context.Context, path string, body, dst any) error {
	return c.do(ctx, http.MethodDelete, path, body, dst)
}

// Patch performs an authenticated PATCH request with a JSON body and decodes the response into dst.
func (c *BaseClient) Patch(ctx context.Context, path string, body, dst any) error {
	return c.do(ctx, http.MethodPatch, path, body, dst)
}

// Head performs an authenticated HEAD request and returns nil on a 2xx status.
func (c *BaseClient) Head(ctx context.Context, path string) error {
	ref, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("arr: parse path: %w", err)
	}
	u := c.baseURL.ResolveReference(ref)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("arr: create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("arr: HEAD %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Method:     http.MethodHead,
			Path:       path,
		}
	}
	return nil
}

// GetRaw performs an authenticated GET request and returns the raw response body as bytes.
// Use this for endpoints that return non-JSON content (e.g. plain text log files).
func (c *BaseClient) GetRaw(ctx context.Context, path string) ([]byte, error) {
	ref, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("arr: parse path: %w", err)
	}
	u := c.baseURL.ResolveReference(ref)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("arr: create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("arr: GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("arr: read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Method:     http.MethodGet,
			Path:       path,
			Body:       body,
		}
	}
	return body, nil
}

// Upload performs an authenticated multipart file upload via POST.
// The file content is sent as a form field with the given fieldName and fileName.
func (c *BaseClient) Upload(ctx context.Context, path, fieldName, fileName string, fileData io.Reader) error {
	ref, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("arr: parse path: %w", err)
	}
	u := c.baseURL.ResolveReference(ref)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile(fieldName, fileName)
	if err != nil {
		return fmt.Errorf("arr: create form file: %w", err)
	}
	if _, err := io.Copy(part, fileData); err != nil {
		return fmt.Errorf("arr: copy file data: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("arr: close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &buf)
	if err != nil {
		return fmt.Errorf("arr: create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("arr: POST %s: %w", path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("arr: read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Method:     http.MethodPost,
			Path:       path,
			Body:       respBody,
		}
	}
	return nil
}

// do is the internal method that executes every HTTP request.
func (c *BaseClient) do(ctx context.Context, method, path string, body, dst any) error {
	ref, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("arr: parse path: %w", err)
	}
	u := c.baseURL.ResolveReference(ref)

	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("arr: marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return fmt.Errorf("arr: create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("arr: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("arr: read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Method:     method,
			Path:       path,
			Body:       respBody,
		}
	}

	if dst != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dst); err != nil {
			return fmt.Errorf("arr: decode response: %w", err)
		}
	}

	return nil
}
