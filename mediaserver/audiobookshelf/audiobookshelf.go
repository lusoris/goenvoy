package audiobookshelf

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

// Client is an Audiobookshelf API client.
type Client struct {
	rawBaseURL string
	token      string
	httpClient *http.Client
}

// New creates an Audiobookshelf [Client] for the instance at baseURL with the given token.
func New(baseURL, token string, opts ...Option) *Client {
	c := &Client{
		rawBaseURL: baseURL,
		token:      token,
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
		return fmt.Sprintf("audiobookshelf: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("audiobookshelf: HTTP %d", e.StatusCode)
}

func (c *Client) get(ctx context.Context, path string, params url.Values) ([]byte, error) {
	u, err := url.Parse(c.rawBaseURL + "/api" + path)
	if err != nil {
		return nil, fmt.Errorf("audiobookshelf: parse URL: %w", err)
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("audiobookshelf: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("audiobookshelf: GET %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("audiobookshelf: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, RawBody: string(body)}
	}
	return body, nil
}

// Libraries.

// GetLibraries returns all libraries.
func (c *Client) GetLibraries(ctx context.Context) ([]Library, error) {
	data, err := c.get(ctx, "/libraries", nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Libraries []Library `json:"libraries"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode libraries: %w", err)
	}
	return resp.Libraries, nil
}

// GetLibrary returns a library by ID.
func (c *Client) GetLibrary(ctx context.Context, libraryID string) (*Library, error) {
	data, err := c.get(ctx, "/libraries/"+url.PathEscape(libraryID), nil)
	if err != nil {
		return nil, err
	}
	var lib Library
	if err := json.Unmarshal(data, &lib); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode library: %w", err)
	}
	return &lib, nil
}

// GetLibraryItems returns items in a library with pagination.
func (c *Client) GetLibraryItems(ctx context.Context, libraryID string, page, limit int) (*LibraryItemsResponse, error) {
	p := url.Values{}
	if page > 0 {
		p.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}
	data, err := c.get(ctx, "/libraries/"+url.PathEscape(libraryID)+"/items", p)
	if err != nil {
		return nil, err
	}
	var resp LibraryItemsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode library items: %w", err)
	}
	return &resp, nil
}

// Items.

// GetItem returns a library item by ID.
func (c *Client) GetItem(ctx context.Context, itemID string) (*LibraryItem, error) {
	data, err := c.get(ctx, "/items/"+url.PathEscape(itemID), nil)
	if err != nil {
		return nil, err
	}
	var item LibraryItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode item: %w", err)
	}
	return &item, nil
}

// Users.

// GetUsers returns all users (admin only).
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	data, err := c.get(ctx, "/users", nil)
	if err != nil {
		return nil, err
	}
	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode users: %w", err)
	}
	return users, nil
}

// GetMe returns the authenticated user.
func (c *Client) GetMe(ctx context.Context) (*User, error) {
	data, err := c.get(ctx, "/me", nil)
	if err != nil {
		return nil, err
	}
	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode user: %w", err)
	}
	return &user, nil
}

// Collections.

// GetCollections returns all collections for a library.
func (c *Client) GetCollections(ctx context.Context, libraryID string) ([]Collection, error) {
	data, err := c.get(ctx, "/libraries/"+url.PathEscape(libraryID)+"/collections", nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Results []Collection `json:"results"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode collections: %w", err)
	}
	return resp.Results, nil
}

// Sessions.

// GetSessions returns active listening sessions.
func (c *Client) GetSessions(ctx context.Context) ([]PlaybackSession, error) {
	data, err := c.get(ctx, "/sessions", nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Sessions []PlaybackSession `json:"sessions"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode sessions: %w", err)
	}
	return resp.Sessions, nil
}

// Server.

// GetServerInfo returns server information.
func (c *Client) GetServerInfo(ctx context.Context) (*ServerInfo, error) {
	data, err := c.get(ctx, "/server", nil)
	if err != nil {
		return nil, err
	}
	var info ServerInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode server info: %w", err)
	}
	return &info, nil
}

// Search searches within a library.
func (c *Client) Search(ctx context.Context, libraryID, query string) ([]byte, error) {
	p := url.Values{"q": {query}}
	return c.get(ctx, "/libraries/"+url.PathEscape(libraryID)+"/search", p)
}

// GetMediaProgress returns listening progress for an item.
func (c *Client) GetMediaProgress(ctx context.Context, itemID string) (*MediaProgress, error) {
	data, err := c.get(ctx, "/me/progress/"+url.PathEscape(itemID), nil)
	if err != nil {
		return nil, err
	}
	var progress MediaProgress
	if err := json.Unmarshal(data, &progress); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode progress: %w", err)
	}
	return &progress, nil
}

// GetAuthors returns all authors in a library.
func (c *Client) GetAuthors(ctx context.Context, libraryID string) ([]Author, error) {
	data, err := c.get(ctx, "/libraries/"+url.PathEscape(libraryID)+"/authors", nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Authors []Author `json:"authors"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode authors: %w", err)
	}
	return resp.Authors, nil
}

// GetSeries returns all series in a library.
func (c *Client) GetSeries(ctx context.Context, libraryID string) ([]Series, error) {
	data, err := c.get(ctx, "/libraries/"+url.PathEscape(libraryID)+"/series", nil)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Results []Series `json:"results"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("audiobookshelf: decode series: %w", err)
	}
	return resp.Results, nil
}
