package jackett

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	availableYes   = "yes"
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

// Client is a Jackett API client.
type Client struct {
	rawBaseURL string
	apiKey     string
	httpClient *http.Client
}

// New creates a Jackett [Client] for the instance at baseURL with the given API key.
func New(baseURL, apiKey string, opts ...Option) *Client {
	c := &Client{
		rawBaseURL: baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("jackett: HTTP %d: %s", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("jackett: HTTP %d", e.StatusCode)
}

func (c *Client) doRequest(ctx context.Context, u string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("jackett: create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jackett: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("jackett: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	return data, nil
}

func (c *Client) getJSON(ctx context.Context, path string, params url.Values, v any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("apikey", c.apiKey)
	u := c.rawBaseURL + path + "?" + params.Encode()

	data, err := c.doRequest(ctx, u)
	if err != nil {
		return err
	}

	if v != nil {
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("jackett: decode json: %w", err)
		}
	}
	return nil
}

func (c *Client) getXML(ctx context.Context, path string, params url.Values, v any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("apikey", c.apiKey)
	u := c.rawBaseURL + path + "?" + params.Encode()

	data, err := c.doRequest(ctx, u)
	if err != nil {
		return err
	}

	if v != nil {
		if err := xml.Unmarshal(data, v); err != nil {
			return fmt.Errorf("jackett: decode xml: %w", err)
		}
	}
	return nil
}

func (c *Client) torznabSearch(ctx context.Context, indexer string, params url.Values) ([]SearchResult, error) {
	path := "/api/v2.0/indexers/" + indexer + "/results/torznab"
	var rss rssResponse
	if err := c.getXML(ctx, path, params, &rss); err != nil {
		return nil, err
	}
	return convertItems(rss.Channel.Items), nil
}

func convertItems(items []rssItem) []SearchResult {
	results := make([]SearchResult, 0, len(items))
	for _, item := range items {
		r := SearchResult{
			Title:    item.Title,
			GUID:     item.GUID,
			Link:     item.Link,
			Comments: item.Comments,
			PubDate:  item.PubDate,
			Size:     item.Size,
		}
		for _, a := range item.Attrs {
			if a.XMLName.Local != "attr" {
				continue
			}
			switch a.Name {
			case "category":
				r.Category = a.Value
			case "description":
				r.CategoryDesc = a.Value
			case "seeders":
				r.Seeders, _ = strconv.Atoi(a.Value)
			case "peers":
				r.Peers, _ = strconv.Atoi(a.Value)
			case "infohash":
				r.InfoHash = a.Value
			case "magneturl":
				r.MagnetURL = a.Value
			case "minimumratio":
				r.MinimumRatio = a.Value
			case "minimumseedtime":
				r.MinimumSeedTime = a.Value
			case "downloadvolumefactor":
				r.DownloadVolumeFactor = a.Value
			case "uploadvolumefactor":
				r.UploadVolumeFactor = a.Value
			case "indexer":
				r.Indexer = a.Value
			}
		}
		results = append(results, r)
	}
	return results
}

func catParams(categories []int) string {
	if len(categories) == 0 {
		return ""
	}
	s := make([]string, len(categories))
	for i, c := range categories {
		s[i] = strconv.Itoa(c)
	}
	return strings.Join(s, ",")
}

// Search queries all indexers for the given term and optional category filter.
func (c *Client) Search(ctx context.Context, query string, categories []int) ([]SearchResult, error) {
	params := url.Values{"t": {"search"}, "q": {query}}
	if cats := catParams(categories); cats != "" {
		params.Set("cat", cats)
	}
	return c.torznabSearch(ctx, "all", params)
}

// SearchIndexer queries a specific indexer by ID.
func (c *Client) SearchIndexer(ctx context.Context, indexerID, query string, categories []int) ([]SearchResult, error) {
	params := url.Values{"t": {"search"}, "q": {query}}
	if cats := catParams(categories); cats != "" {
		params.Set("cat", cats)
	}
	return c.torznabSearch(ctx, indexerID, params)
}

// TVSearch performs a TV-specific search across all indexers.
func (c *Client) TVSearch(ctx context.Context, query string, season, episode int, imdbID string) ([]SearchResult, error) {
	params := url.Values{"t": {"tvsearch"}, "q": {query}}
	if season > 0 {
		params.Set("season", strconv.Itoa(season))
	}
	if episode > 0 {
		params.Set("ep", strconv.Itoa(episode))
	}
	if imdbID != "" {
		params.Set("imdbid", imdbID)
	}
	return c.torznabSearch(ctx, "all", params)
}

// MovieSearch performs a movie-specific search across all indexers.
func (c *Client) MovieSearch(ctx context.Context, query, imdbID string) ([]SearchResult, error) {
	params := url.Values{"t": {"movie"}, "q": {query}}
	if imdbID != "" {
		params.Set("imdbid", imdbID)
	}
	return c.torznabSearch(ctx, "all", params)
}

// MusicSearch performs a music-specific search across all indexers.
func (c *Client) MusicSearch(ctx context.Context, query string) ([]SearchResult, error) {
	params := url.Values{"t": {"music"}, "q": {query}}
	return c.torznabSearch(ctx, "all", params)
}

// BookSearch performs a book-specific search across all indexers.
func (c *Client) BookSearch(ctx context.Context, query string) ([]SearchResult, error) {
	params := url.Values{"t": {"book"}, "q": {query}}
	return c.torznabSearch(ctx, "all", params)
}

// GetCapabilities returns the server capabilities.
func (c *Client) GetCapabilities(ctx context.Context) (*Capabilities, error) {
	params := url.Values{"t": {"caps"}}
	var raw capsResponse
	if err := c.getXML(ctx, "/api/v2.0/indexers/all/results/torznab", params, &raw); err != nil {
		return nil, err
	}
	return convertCaps(&raw), nil
}

func convertCaps(raw *capsResponse) *Capabilities {
	caps := &Capabilities{
		Server: ServerInfo{Title: raw.Server.Title, Image: raw.Server.Image},
		Limits: Limits{Max: raw.Limits.Max, Default: raw.Limits.Default},
		Searching: Searching{
			SearchAvailable:      raw.Searching.Search.Available == availableYes,
			TVSearchAvailable:    raw.Searching.TVSearch.Available == availableYes,
			MovieSearchAvailable: raw.Searching.MovieSearch.Available == availableYes,
			MusicSearchAvailable: raw.Searching.MusicSearch.Available == availableYes,
			BookSearchAvailable:  raw.Searching.BookSearch.Available == availableYes,
		},
	}
	for _, c := range raw.Categories.Categories {
		cat := Category{ID: c.ID, Name: c.Name}
		for _, sc := range c.SubCategories {
			cat.SubCategories = append(cat.SubCategories, SubCategory(sc))
		}
		caps.Categories = append(caps.Categories, cat)
	}
	return caps
}

// GetIndexers returns the list of configured indexers.
func (c *Client) GetIndexers(ctx context.Context) ([]Indexer, error) {
	var out []Indexer
	params := url.Values{"configured": {"true"}}
	if err := c.getJSON(ctx, "/api/v2.0/indexers", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetServerConfig returns the Jackett server configuration.
func (c *Client) GetServerConfig(ctx context.Context) (*ServerConfig, error) {
	var out ServerConfig
	if err := c.getJSON(ctx, "/api/v2.0/server/config", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
