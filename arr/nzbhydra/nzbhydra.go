package nzbhydra

import (
	"bytes"
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

// Client is an NZBHydra2 API client.
type Client struct {
	rawBaseURL string
	apiKey     string
	httpClient *http.Client
}

// New creates an NZBHydra2 [Client] for the instance at baseURL with the given API key.
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
		return fmt.Sprintf("nzbhydra: HTTP %d: %s", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("nzbhydra: HTTP %d", e.StatusCode)
}

func (c *Client) doGet(ctx context.Context, u string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("nzbhydra: create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nzbhydra: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("nzbhydra: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(data)}
	}

	return data, nil
}

func (c *Client) getXML(ctx context.Context, path string, params url.Values, v any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("apikey", c.apiKey)
	u := c.rawBaseURL + path + "?" + params.Encode()

	data, err := c.doGet(ctx, u)
	if err != nil {
		return err
	}

	if v != nil {
		if err := xml.Unmarshal(data, v); err != nil {
			return fmt.Errorf("nzbhydra: decode xml: %w", err)
		}
	}
	return nil
}

func (c *Client) getJSON(ctx context.Context, path string, params url.Values, v any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("apikey", c.apiKey)
	u := c.rawBaseURL + path + "?" + params.Encode()

	data, err := c.doGet(ctx, u)
	if err != nil {
		return err
	}

	if v != nil {
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("nzbhydra: decode json: %w", err)
		}
	}
	return nil
}

func (c *Client) postJSON(ctx context.Context, path string, body, v any) error {
	params := url.Values{"apikey": {c.apiKey}}
	u := c.rawBaseURL + path + "?" + params.Encode()

	var reqBody io.Reader = http.NoBody
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("nzbhydra: encode request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, reqBody)
	if err != nil {
		return fmt.Errorf("nzbhydra: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("nzbhydra: request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("nzbhydra: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status, Body: string(respData)}
	}

	if v != nil {
		if err := json.Unmarshal(respData, v); err != nil {
			return fmt.Errorf("nzbhydra: decode json: %w", err)
		}
	}
	return nil
}

func (c *Client) newznabSearch(ctx context.Context, params url.Values) ([]SearchResult, error) {
	var rss rssResponse
	if err := c.getXML(ctx, "/api", params, &rss); err != nil {
		return nil, err
	}
	return convertItems(rss.Channel.Items), nil
}

func convertItems(items []rssItem) []SearchResult {
	results := make([]SearchResult, 0, len(items))
	for i := range items {
		r := SearchResult{
			Title:       items[i].Title,
			GUID:        items[i].GUID,
			Link:        items[i].Link,
			Comments:    items[i].Comments,
			PubDate:     items[i].PubDate,
			Size:        items[i].Size,
			Description: items[i].Description,
		}
		for _, a := range items[i].Attrs {
			if a.XMLName.Local != "attr" {
				continue
			}
			switch a.Name {
			case "category":
				r.Category = a.Value
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

// Search queries NZBHydra2 for the given term and optional category filter.
func (c *Client) Search(ctx context.Context, query string, categories []int) ([]SearchResult, error) {
	params := url.Values{"t": {"search"}, "q": {query}}
	if cats := catParams(categories); cats != "" {
		params.Set("cat", cats)
	}
	return c.newznabSearch(ctx, params)
}

// TVSearch performs a TV-specific search with optional identifiers.
func (c *Client) TVSearch(ctx context.Context, query string, opts *TVSearchOptions) ([]SearchResult, error) {
	params := url.Values{"t": {"tvsearch"}, "q": {query}}
	if opts.Season > 0 {
		params.Set("season", strconv.Itoa(opts.Season))
	}
	if opts.Episode > 0 {
		params.Set("ep", strconv.Itoa(opts.Episode))
	}
	if opts.TVDBID != "" {
		params.Set("tvdbid", opts.TVDBID)
	}
	if opts.IMDBID != "" {
		params.Set("imdbid", opts.IMDBID)
	}
	if opts.TMDBID != "" {
		params.Set("tmdbid", opts.TMDBID)
	}
	if opts.TVMazeID != "" {
		params.Set("tvmazeid", opts.TVMazeID)
	}
	if opts.RID != "" {
		params.Set("rid", opts.RID)
	}
	return c.newznabSearch(ctx, params)
}

// MovieSearch performs a movie-specific search across all indexers.
func (c *Client) MovieSearch(ctx context.Context, query, imdbID, tmdbID string) ([]SearchResult, error) {
	params := url.Values{"t": {"movie"}, "q": {query}}
	if imdbID != "" {
		params.Set("imdbid", imdbID)
	}
	if tmdbID != "" {
		params.Set("tmdbid", tmdbID)
	}
	return c.newznabSearch(ctx, params)
}

// BookSearch performs a book-specific search across all indexers.
func (c *Client) BookSearch(ctx context.Context, query string) ([]SearchResult, error) {
	params := url.Values{"t": {"book"}, "q": {query}}
	return c.newznabSearch(ctx, params)
}

// GetCapabilities returns the server capabilities.
func (c *Client) GetCapabilities(ctx context.Context) (*Capabilities, error) {
	params := url.Values{"t": {"caps"}}
	var raw capsResponse
	if err := c.getXML(ctx, "/api", params, &raw); err != nil {
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

// GetStats retrieves statistics from NZBHydra2.
func (c *Client) GetStats(ctx context.Context, req StatsRequest) (*StatsResponse, error) {
	var out StatsResponse
	if err := c.postJSON(ctx, "/api/stats", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSearchHistory retrieves search history with pagination.
func (c *Client) GetSearchHistory(ctx context.Context, req HistoryRequest) (*PagedResponse[SearchHistoryEntry], error) {
	var out PagedResponse[SearchHistoryEntry]
	if err := c.postJSON(ctx, "/api/history/searches", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetDownloadHistory retrieves download history with pagination.
func (c *Client) GetDownloadHistory(ctx context.Context, req HistoryRequest) (*PagedResponse[DownloadHistoryEntry], error) {
	var out PagedResponse[DownloadHistoryEntry]
	if err := c.postJSON(ctx, "/api/history/downloads", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetIndexerStatuses returns the status of all configured indexers.
func (c *Client) GetIndexerStatuses(ctx context.Context) ([]IndexerStatus, error) {
	var out []IndexerStatus
	if err := c.getJSON(ctx, "/api/stats/indexers", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
