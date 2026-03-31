package tdarr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 30 * time.Second

// Option configures a [Client].
type Option func(*Client)

// WithAPIKey sets the optional API key for authentication.
func WithAPIKey(key string) Option {
	return func(c *Client) { c.apiKey = key }
}

// WithHTTPClient sets a custom [http.Client].
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.httpClient = c }
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int    `json:"-"`
	Status     string `json:"-"`
	Body       string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("tdarr: HTTP %d: %s", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("tdarr: HTTP %d", e.StatusCode)
}

// Client is a Tdarr v2 API client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// New creates a Tdarr [Client] for the instance at baseURL.
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		baseURL:    baseURL,
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
			return fmt.Errorf("tdarr: marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("tdarr: create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("tdarr: %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("tdarr: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(raw)}
	}

	if v != nil {
		if err := json.Unmarshal(raw, v); err != nil {
			return fmt.Errorf("tdarr: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) doText(ctx context.Context, method, path string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("tdarr: create request: %w", err)
	}
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("tdarr: %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("tdarr: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(raw)}
	}
	return string(raw), nil
}

// GetStatus returns the current Tdarr server status.
func (c *Client) GetStatus(ctx context.Context) (*Status, error) {
	var out Status
	if err := c.do(ctx, http.MethodGet, "/api/v2/status", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNodes returns all registered Tdarr nodes keyed by node ID.
func (c *Client) GetNodes(ctx context.Context) (map[string]Node, error) {
	var out map[string]Node
	if err := c.do(ctx, http.MethodGet, "/api/v2/get-nodes", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SearchDB searches the Tdarr database for files matching the given filters.
func (c *Client) SearchDB(ctx context.Context, collection string, limit, skip int, filters []SearchFilter) ([]DBFile, error) {
	req := SearchDBRequest{
		Data: &SearchDBData{
			Collection: collection,
			Limit:      limit,
			Skip:       skip,
			Filters:    filters,
		},
	}
	var out []DBFile
	if err := c.do(ctx, http.MethodPost, "/api/v2/search-db", req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CrudDB performs a CRUD operation on the Tdarr database.
func (c *Client) CrudDB(ctx context.Context, collection, mode, docID string) ([]map[string]any, error) {
	req := CrudDBRequest{
		Data: &CrudDBData{
			Collection: collection,
			Mode:       mode,
			DocID:      docID,
		},
	}
	var out []map[string]any
	if err := c.do(ctx, http.MethodPost, "/api/v2/cruddb", req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetResStats returns resolution statistics for a library.
func (c *Client) GetResStats(ctx context.Context, libraryId string) (*ResStats, error) {
	body := map[string]any{"data": map[string]string{"libraryId": libraryId}}
	var out ResStats
	if err := c.do(ctx, http.MethodPost, "/api/v2/get-res-stats", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDBStatuses returns database table counts for a library.
func (c *Client) GetDBStatuses(ctx context.Context, libraryId string) (*DBStatuses, error) {
	body := map[string]any{"data": map[string]string{"libraryId": libraryId}}
	var out DBStatuses
	if err := c.do(ctx, http.MethodPost, "/api/v2/get-db-statuses", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ScanFiles triggers a file scan for a library.
func (c *Client) ScanFiles(ctx context.Context, libraryId, folderPath string) error {
	req := ScanFilesRequest{
		Data: &ScanFilesData{
			LibraryId:  libraryId,
			FolderPath: folderPath,
		},
	}
	return c.do(ctx, http.MethodPost, "/api/v2/scan-files", req, nil)
}

// CancelWorkerItem cancels the current item being processed by a worker.
func (c *Client) CancelWorkerItem(ctx context.Context, nodeId, workerId string) error {
	body := map[string]any{"data": map[string]string{"nodeId": nodeId, "workerId": workerId}}
	return c.do(ctx, http.MethodPost, "/api/v2/cancel-worker-item", body, nil)
}

// KillWorker terminates a worker on a node.
func (c *Client) KillWorker(ctx context.Context, nodeId, workerId string) error {
	body := map[string]any{"data": map[string]string{"nodeId": nodeId, "workerId": workerId}}
	return c.do(ctx, http.MethodPost, "/api/v2/kill-worker", body, nil)
}

// GetServerLog returns the Tdarr server log as plain text.
func (c *Client) GetServerLog(ctx context.Context) (string, error) {
	return c.doText(ctx, http.MethodGet, "/api/v2/get-server-log")
}
