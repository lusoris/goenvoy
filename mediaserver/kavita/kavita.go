package kavita

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultTimeout = 30 * time.Second

// Client is a Kavita API client.
type Client struct {
	baseURL string
	apiKey  string
	token   string
	mu      sync.RWMutex
	http    *http.Client
}

// Option configures the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.http = c }
}

// New creates a new Kavita client.
//
// The baseURL should include the protocol and host (e.g. "http://localhost:5000").
// Authentication uses an API key that is exchanged for a JWT token via the
// Plugin authenticate endpoint.
func New(baseURL, apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		http:    &http.Client{Timeout: defaultTimeout},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Authenticate exchanges the API key for a JWT token.
func (c *Client) Authenticate(ctx context.Context) error {
	body, err := json.Marshal(map[string]string{"apiKey": c.apiKey})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/Plugin/authenticate", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(rb)}
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rb, &result); err != nil {
		return err
	}
	c.mu.Lock()
	c.token = result.Token
	c.mu.Unlock()
	return nil
}

func (c *Client) do(ctx context.Context, method, path string, body, v any) error {
	c.mu.RLock()
	needAuth := c.token == ""
	c.mu.RUnlock()
	if needAuth {
		if err := c.Authenticate(ctx); err != nil {
			return err
		}
	}

	resp, err := c.send(ctx, method, path, body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		_ = resp.Body.Close()
		if err := c.Authenticate(ctx); err != nil {
			return err
		}
		resp, err = c.send(ctx, method, path, body)
		if err != nil {
			return err
		}
	}
	defer resp.Body.Close()

	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(rb)}
	}

	if v == nil {
		return nil
	}
	return json.Unmarshal(rb, v)
}

func (c *Client) getToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}

func (c *Client) send(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var reqBody io.Reader = http.NoBody
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.getToken())
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.http.Do(req)
}

// GetLibraries returns all libraries.
func (c *Client) GetLibraries(ctx context.Context) ([]Library, error) {
	var libs []Library
	if err := c.do(ctx, http.MethodGet, "/api/Library", nil, &libs); err != nil {
		return nil, err
	}
	return libs, nil
}

// GetLibrary returns a single library by ID.
func (c *Client) GetLibrary(ctx context.Context, id int) (*Library, error) {
	var lib Library
	if err := c.do(ctx, http.MethodGet, "/api/Library/"+strconv.Itoa(id), nil, &lib); err != nil {
		return nil, err
	}
	return &lib, nil
}

// ScanLibrary triggers a scan for the given library.
func (c *Client) ScanLibrary(ctx context.Context, id int) error {
	return c.do(ctx, http.MethodPost, "/api/Library/scan?libraryId="+strconv.Itoa(id), nil, nil)
}

// GetSeries returns all series for a library.
func (c *Client) GetSeries(ctx context.Context, libraryID int) ([]Series, error) {
	var series []Series
	if err := c.do(ctx, http.MethodPost, "/api/Series/v2", map[string]int{"libraryId": libraryID}, &series); err != nil {
		return nil, err
	}
	return series, nil
}

// GetOneSeries returns a single series by ID.
func (c *Client) GetOneSeries(ctx context.Context, id int) (*Series, error) {
	var s Series
	if err := c.do(ctx, http.MethodGet, "/api/Series/"+strconv.Itoa(id), nil, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetVolumes returns all volumes for a series.
func (c *Client) GetVolumes(ctx context.Context, seriesID int) ([]Volume, error) {
	var vols []Volume
	if err := c.do(ctx, http.MethodGet, "/api/Series/volumes?seriesId="+strconv.Itoa(seriesID), nil, &vols); err != nil {
		return nil, err
	}
	return vols, nil
}

// GetChapter returns a single chapter by ID.
func (c *Client) GetChapter(ctx context.Context, id int) (*Chapter, error) {
	var ch Chapter
	if err := c.do(ctx, http.MethodGet, "/api/Chapter?chapterId="+strconv.Itoa(id), nil, &ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

// GetCollections returns all collections.
func (c *Client) GetCollections(ctx context.Context) ([]Collection, error) {
	var cols []Collection
	if err := c.do(ctx, http.MethodGet, "/api/Collection", nil, &cols); err != nil {
		return nil, err
	}
	return cols, nil
}

// GetReadingLists returns all reading lists.
func (c *Client) GetReadingLists(ctx context.Context) ([]ReadingList, error) {
	var rls []ReadingList
	if err := c.do(ctx, http.MethodGet, "/api/ReadingList", nil, &rls); err != nil {
		return nil, err
	}
	return rls, nil
}

// GetUsers returns all users.
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	var users []User
	if err := c.do(ctx, http.MethodGet, "/api/Users", nil, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetServerInfo returns server information.
func (c *Client) GetServerInfo(ctx context.Context) (*ServerInfo, error) {
	var info ServerInfo
	if err := c.do(ctx, http.MethodGet, "/api/Server/server-info", nil, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// Search performs a search query and returns matching results.
func (c *Client) Search(ctx context.Context, query string) (*SearchResult, error) {
	var sr SearchResult
	if err := c.do(ctx, http.MethodGet, "/api/Search/search?queryString="+url.QueryEscape(query), nil, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}
