package openlibrary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://openlibrary.org"

// Client is an Open Library API client.
type Client struct {
	*metadata.BaseClient
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("openlibrary: %s: %s", e.Status, e.Body)
}

// New creates an Open Library [Client].
func New(opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "openlibrary", opts...)
	return &Client{BaseClient: bc}
}

func (c *Client) get(ctx context.Context, path string, params url.Values, v any) error {
	u := c.BaseURL() + path
	if params != nil {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("openlibrary: create request: %w", err)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("openlibrary: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("openlibrary: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(body)}
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("openlibrary: decode response: %w", err)
	}
	return nil
}

// Search searches for books by query string.
func (c *Client) Search(ctx context.Context, query string) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("q", query)
	var resp SearchResponse
	if err := c.get(ctx, "/search.json", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchWithParams searches for books with custom parameters (title, author, subject, isbn, etc.).
func (c *Client) SearchWithParams(ctx context.Context, params url.Values) (*SearchResponse, error) {
	var resp SearchResponse
	if err := c.get(ctx, "/search.json", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWork returns a work by its Open Library ID (e.g. "OL45804W").
func (c *Client) GetWork(ctx context.Context, id string) (*Work, error) {
	var resp Work
	if err := c.get(ctx, "/works/"+id+".json", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetEdition returns an edition by its Open Library ID (e.g. "OL7353617M").
func (c *Client) GetEdition(ctx context.Context, id string) (*Edition, error) {
	var resp Edition
	if err := c.get(ctx, "/books/"+id+".json", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAuthor returns an author by their Open Library ID (e.g. "OL34184A").
func (c *Client) GetAuthor(ctx context.Context, id string) (*Author, error) {
	var resp Author
	if err := c.get(ctx, "/authors/"+id+".json", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// authorWorksResp is the response wrapper for author works.
type authorWorksResp struct {
	Entries []SubjectWork `json:"entries"`
}

// GetAuthorWorks returns works by an author.
func (c *Client) GetAuthorWorks(ctx context.Context, id string) ([]SubjectWork, error) {
	var resp authorWorksResp
	if err := c.get(ctx, "/authors/"+id+"/works.json", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// GetSubject returns information about a subject (e.g. "fantasy").
func (c *Client) GetSubject(ctx context.Context, subject string) (*Subject, error) {
	var resp Subject
	if err := c.get(ctx, "/subjects/"+subject+".json", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetByISBN returns an edition by ISBN.
func (c *Client) GetByISBN(ctx context.Context, isbn string) (*Edition, error) {
	var resp Edition
	if err := c.get(ctx, "/isbn/"+isbn+".json", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
