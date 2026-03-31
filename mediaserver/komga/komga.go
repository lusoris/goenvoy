package komga

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Client is a Komga API client.
type Client struct {
	baseURL  string
	username string
	password string
	http     *http.Client
}

// Option configures the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.http = c }
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("komga: %s: %s", e.Status, e.Body)
}

// New creates a new Komga client.
//
// The baseURL should include the protocol and host (e.g. "http://localhost:25600").
// Authentication uses HTTP Basic Auth with the provided username and password.
func New(baseURL, username, password string, opts ...Option) *Client {
	c := &Client{
		baseURL:  strings.TrimRight(baseURL, "/"),
		username: username,
		password: password,
		http:     http.DefaultClient,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) get(ctx context.Context, path string, params url.Values, v any) error {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	return json.Unmarshal(body, v)
}

// GetLibraries returns all libraries.
func (c *Client) GetLibraries(ctx context.Context) ([]Library, error) {
	var libs []Library
	if err := c.get(ctx, "/api/v1/libraries", nil, &libs); err != nil {
		return nil, err
	}
	return libs, nil
}

// GetLibrary returns a single library by ID.
func (c *Client) GetLibrary(ctx context.Context, id string) (*Library, error) {
	var lib Library
	if err := c.get(ctx, "/api/v1/libraries/"+url.PathEscape(id), nil, &lib); err != nil {
		return nil, err
	}
	return &lib, nil
}

// GetSeries returns a paginated list of series, optionally filtered by library.
func (c *Client) GetSeries(ctx context.Context, libraryID string, page, size int) (*Page[Series], error) {
	params := url.Values{}
	if libraryID != "" {
		params.Set("library_id", libraryID)
	}
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(size))

	var p Page[Series]
	if err := c.get(ctx, "/api/v1/series", params, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetOneSeries returns a single series by ID.
func (c *Client) GetOneSeries(ctx context.Context, id string) (*Series, error) {
	var s Series
	if err := c.get(ctx, "/api/v1/series/"+url.PathEscape(id), nil, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetBooks returns a paginated list of books for a series.
func (c *Client) GetBooks(ctx context.Context, seriesID string, page, size int) (*Page[Book], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(size))

	var p Page[Book]
	if err := c.get(ctx, "/api/v1/series/"+url.PathEscape(seriesID)+"/books", params, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetBook returns a single book by ID.
func (c *Client) GetBook(ctx context.Context, id string) (*Book, error) {
	var b Book
	if err := c.get(ctx, "/api/v1/books/"+url.PathEscape(id), nil, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// GetCollections returns a paginated list of collections.
func (c *Client) GetCollections(ctx context.Context, page, size int) (*Page[Collection], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(size))

	var p Page[Collection]
	if err := c.get(ctx, "/api/v1/collections", params, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetCollection returns a single collection by ID.
func (c *Client) GetCollection(ctx context.Context, id string) (*Collection, error) {
	var col Collection
	if err := c.get(ctx, "/api/v1/collections/"+url.PathEscape(id), nil, &col); err != nil {
		return nil, err
	}
	return &col, nil
}

// GetReadLists returns a paginated list of read lists.
func (c *Client) GetReadLists(ctx context.Context, page, size int) (*Page[ReadList], error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(size))

	var p Page[ReadList]
	if err := c.get(ctx, "/api/v1/readlists", params, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetReadList returns a single read list by ID.
func (c *Client) GetReadList(ctx context.Context, id string) (*ReadList, error) {
	var rl ReadList
	if err := c.get(ctx, "/api/v1/readlists/"+url.PathEscape(id), nil, &rl); err != nil {
		return nil, err
	}
	return &rl, nil
}

// GetUsers returns all users (admin only).
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	var users []User
	if err := c.get(ctx, "/api/v1/users", nil, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// GetMe returns the authenticated user.
func (c *Client) GetMe(ctx context.Context) (*User, error) {
	var u User
	if err := c.get(ctx, "/api/v1/users/me", nil, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
