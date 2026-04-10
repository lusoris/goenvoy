package tpdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://api.theporndb.net"

// Client is a ThePornDB (TPDB) API client.
type Client struct {
	*metadata.BaseClient
	apiToken string
}

// New creates a TPDB [Client] using the given API token.
func New(apiToken string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "tpdb", opts...)
	c := &Client{BaseClient: bc, apiToken: apiToken}
	bc.SetAuth(func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+apiToken)
	})
	return c
}


// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	// RawBody holds the raw response body when the error could not be parsed as JSON.
	RawBody string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("tpdb: HTTP %d: %s", e.StatusCode, e.Message)
	}
	if e.RawBody != "" {
		return fmt.Sprintf("tpdb: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("tpdb: HTTP %d", e.StatusCode)
}

// listResponse wraps a paginated list response from the API.
type listResponse[T any] struct {
	Data []T        `json:"data"`
	Meta Pagination `json:"meta"`
}

// itemResponse wraps a single-item response.
type itemResponse[T any] struct {
	Data T `json:"data"`
}

func (c *Client) doRequest(ctx context.Context, method, path string, params url.Values, dst any) error {
	u, err := url.Parse(c.BaseURL() + path)
	if err != nil {
		return fmt.Errorf("tpdb: parse URL: %w", err)
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("tpdb: create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("tpdb: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("tpdb: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if err := json.Unmarshal(body, apiErr); err != nil {
			apiErr.RawBody = string(body)
		}
		return apiErr
	}

	if dst != nil && len(body) > 0 {
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("tpdb: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) get(ctx context.Context, path string, params url.Values, dst any) error {
	return c.doRequest(ctx, http.MethodGet, path, params, dst)
}

func pageParams(page, perPage int) url.Values {
	p := url.Values{}
	if page > 0 {
		p.Set("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		p.Set("per_page", strconv.Itoa(perPage))
	}
	return p
}

// Scenes.

// SceneSearchParams configures scene search filters.
type SceneSearchParams struct {
	Query     string
	Page      int
	PerPage   int
	Performer string
	Site      string
	Tag       string
	Director  string
	DateFrom  string
	DateTo    string
	Year      int
	Duration  string
	Hash      string
	Category  string
	OrderBy   string
}

func (p *SceneSearchParams) values() url.Values {
	v := pageParams(p.Page, p.PerPage)
	if p.Query != "" {
		v.Set("q", p.Query)
	}
	if p.Performer != "" {
		v.Set("performer", p.Performer)
	}
	if p.Site != "" {
		v.Set("site", p.Site)
	}
	if p.Tag != "" {
		v.Set("tag", p.Tag)
	}
	if p.Director != "" {
		v.Set("director", p.Director)
	}
	if p.DateFrom != "" {
		v.Set("date_from", p.DateFrom)
	}
	if p.DateTo != "" {
		v.Set("date_to", p.DateTo)
	}
	if p.Year > 0 {
		v.Set("year", strconv.Itoa(p.Year))
	}
	if p.Duration != "" {
		v.Set("duration", p.Duration)
	}
	if p.Hash != "" {
		v.Set("hash", p.Hash)
	}
	if p.Category != "" {
		v.Set("category", p.Category)
	}
	if p.OrderBy != "" {
		v.Set("order_by", p.OrderBy)
	}
	return v
}

// SearchScenes searches for scenes matching the given criteria.
func (c *Client) SearchScenes(ctx context.Context, params *SceneSearchParams) ([]Scene, *Pagination, error) {
	var resp listResponse[Scene]
	if err := c.get(ctx, "/scenes", params.values(), &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// GetScene returns a single scene by its ID or slug.
func (c *Client) GetScene(ctx context.Context, id string) (*Scene, error) {
	var resp itemResponse[Scene]
	if err := c.get(ctx, "/scenes/"+url.PathEscape(id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// FindSceneByHash finds a scene by its content hash.
func (c *Client) FindSceneByHash(ctx context.Context, hash string) ([]Scene, error) {
	var resp listResponse[Scene]
	p := url.Values{}
	p.Set("hash", hash)
	if err := c.get(ctx, "/scenes", p, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetSimilarScenes returns scenes similar to the given scene.
func (c *Client) GetSimilarScenes(ctx context.Context, id string) ([]Scene, error) {
	var resp listResponse[Scene]
	if err := c.get(ctx, "/scenes/"+url.PathEscape(id)+"/similar", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Performers.

// PerformerSearchParams configures performer search filters.
type PerformerSearchParams struct {
	Query     string
	Page      int
	PerPage   int
	Gender    string
	Ethnicity string
	Age       string
	Height    string
	Weight    string
	HairColor string
	EyeColor  string
	CupSize   string
	Tattoos   string
	Piercings string
	OrderBy   string
}

func (p *PerformerSearchParams) values() url.Values {
	v := pageParams(p.Page, p.PerPage)
	if p.Query != "" {
		v.Set("q", p.Query)
	}
	if p.Gender != "" {
		v.Set("gender", p.Gender)
	}
	if p.Ethnicity != "" {
		v.Set("ethnicity", p.Ethnicity)
	}
	if p.Age != "" {
		v.Set("age", p.Age)
	}
	if p.Height != "" {
		v.Set("height", p.Height)
	}
	if p.Weight != "" {
		v.Set("weight", p.Weight)
	}
	if p.HairColor != "" {
		v.Set("hair_color", p.HairColor)
	}
	if p.EyeColor != "" {
		v.Set("eye_color", p.EyeColor)
	}
	if p.CupSize != "" {
		v.Set("cup_size", p.CupSize)
	}
	if p.Tattoos != "" {
		v.Set("tattoos", p.Tattoos)
	}
	if p.Piercings != "" {
		v.Set("piercings", p.Piercings)
	}
	if p.OrderBy != "" {
		v.Set("order_by", p.OrderBy)
	}
	return v
}

// SearchPerformers searches for performers matching the given criteria.
func (c *Client) SearchPerformers(ctx context.Context, params *PerformerSearchParams) ([]Performer, *Pagination, error) {
	var resp listResponse[Performer]
	if err := c.get(ctx, "/performers", params.values(), &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// GetPerformer returns a single performer by ID or slug.
func (c *Client) GetPerformer(ctx context.Context, id string) (*Performer, error) {
	var resp itemResponse[Performer]
	if err := c.get(ctx, "/performers/"+url.PathEscape(id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetSimilarPerformers returns performers similar to the given performer.
func (c *Client) GetSimilarPerformers(ctx context.Context, id string) ([]Performer, error) {
	var resp listResponse[Performer]
	if err := c.get(ctx, "/performers/"+url.PathEscape(id)+"/similar", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetPerformerScenes returns scenes featuring a performer.
func (c *Client) GetPerformerScenes(ctx context.Context, id string, page, perPage int) ([]Scene, *Pagination, error) {
	var resp listResponse[Scene]
	if err := c.get(ctx, "/performers/"+url.PathEscape(id)+"/scenes", pageParams(page, perPage), &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// GetPerformerMovies returns movies featuring a performer.
func (c *Client) GetPerformerMovies(ctx context.Context, id string, page, perPage int) ([]Movie, *Pagination, error) {
	var resp listResponse[Movie]
	if err := c.get(ctx, "/performers/"+url.PathEscape(id)+"/movies", pageParams(page, perPage), &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// Sites.

// SearchSites searches for sites/studios.
func (c *Client) SearchSites(ctx context.Context, query string, page, perPage int) ([]Site, *Pagination, error) {
	var resp listResponse[Site]
	p := pageParams(page, perPage)
	if query != "" {
		p.Set("q", query)
	}
	if err := c.get(ctx, "/sites", p, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// GetSite returns a single site by its ID or UUID.
func (c *Client) GetSite(ctx context.Context, id string) (*Site, error) {
	var resp itemResponse[Site]
	if err := c.get(ctx, "/sites/"+url.PathEscape(id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Tags.

// ListTags returns a paginated list of tags.
func (c *Client) ListTags(ctx context.Context, query string, page, perPage int) ([]Tag, *Pagination, error) {
	var resp listResponse[Tag]
	p := pageParams(page, perPage)
	if query != "" {
		p.Set("q", query)
	}
	if err := c.get(ctx, "/tags", p, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// GetTag returns a single tag by its ID or slug.
func (c *Client) GetTag(ctx context.Context, id string) (*Tag, error) {
	var resp itemResponse[Tag]
	if err := c.get(ctx, "/tags/"+url.PathEscape(id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Directors.

// ListDirectors returns a paginated list of directors.
func (c *Client) ListDirectors(ctx context.Context, query string, page, perPage int) ([]Director, *Pagination, error) {
	var resp listResponse[Director]
	p := pageParams(page, perPage)
	if query != "" {
		p.Set("q", query)
	}
	if err := c.get(ctx, "/directors", p, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// GetDirector returns a single director by ID or slug.
func (c *Client) GetDirector(ctx context.Context, id string) (*Director, error) {
	var resp itemResponse[Director]
	if err := c.get(ctx, "/directors/"+url.PathEscape(id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Movies/DVDs.

// SearchMovies searches for movies/DVDs.
func (c *Client) SearchMovies(ctx context.Context, query string, page, perPage int) ([]Movie, *Pagination, error) {
	var resp listResponse[Movie]
	p := pageParams(page, perPage)
	if query != "" {
		p.Set("q", query)
	}
	if err := c.get(ctx, "/movies", p, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// GetMovie returns a single movie/DVD by ID or slug.
func (c *Client) GetMovie(ctx context.Context, id string) (*Movie, error) {
	var resp itemResponse[Movie]
	if err := c.get(ctx, "/movies/"+url.PathEscape(id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// JAV (Japanese Adult Video).

// SearchJav searches for JAV content.
func (c *Client) SearchJav(ctx context.Context, query string, page, perPage int) ([]Jav, *Pagination, error) {
	var resp listResponse[Jav]
	p := pageParams(page, perPage)
	if query != "" {
		p.Set("q", query)
	}
	if err := c.get(ctx, "/jav", p, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// GetJav returns a single JAV entry by ID or slug.
func (c *Client) GetJav(ctx context.Context, id string) (*Jav, error) {
	var resp itemResponse[Jav]
	if err := c.get(ctx, "/jav/"+url.PathEscape(id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Hashes.

// LookupByHash resolves content by hash, searching across scenes, movies, and JAV.
func (c *Client) LookupByHash(ctx context.Context, hash string) ([]Scene, error) {
	var resp listResponse[Scene]
	p := url.Values{}
	p.Set("hash", hash)
	if err := c.get(ctx, "/scenes", p, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Changes.

// GetChanges returns recently changed scene IDs since the given timestamp.
// The timestamp should be in RFC3339 format.
func (c *Client) GetChanges(ctx context.Context, since string, page, perPage int) ([]Scene, *Pagination, error) {
	var resp listResponse[Scene]
	p := pageParams(page, perPage)
	if since != "" {
		p.Set("timestamp", since)
	}
	if err := c.get(ctx, "/changelog", p, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}
