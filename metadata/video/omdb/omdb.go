package omdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://www.omdbapi.com"

// Client is an OMDb API client.
type Client struct {
	*metadata.BaseClient
	apiKey string
}

// New creates an OMDb [Client] using the given API key.
func New(apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "omdb", opts...)
	return &Client{BaseClient: bc, apiKey: apiKey}
}

// errorResponse is used to detect API errors in the JSON response.
type errorResponse struct {
	Response string `json:"Response"`
	Error    string `json:"Error"`
}

// APIError is returned when the OMDb API responds with Response=False or a non-2xx status.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return "omdb: " + e.Message
	}
	return "omdb: HTTP " + strconv.Itoa(e.StatusCode)
}

func (c *Client) get(ctx context.Context, params url.Values, dst any) error {
	params.Set("apikey", c.apiKey)
	params.Set("r", "json")

	u := c.BaseURL() + "/?" + params.Encode()

	body, status, err := c.DoRawURL(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return err
	}

	if status < 200 || status >= 300 {
		return &APIError{StatusCode: status, Message: string(body)}
	}

	var errResp errorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Response == "False" {
		return &APIError{StatusCode: status, Message: errResp.Error}
	}

	if dst != nil {
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("omdb: decode response: %w", err)
		}
	}

	return nil
}

// GetByIMDbID retrieves a title by its IMDb ID (e.g. "tt0111161").
// Pass PlotShort or PlotFull to control the plot length; pass "" for the default (short).
func (c *Client) GetByIMDbID(ctx context.Context, imdbID string, plot PlotLength) (*Title, error) {
	params := url.Values{"i": {imdbID}}
	if plot != "" {
		params.Set("plot", string(plot))
	}
	var out Title
	if err := c.get(ctx, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetByTitle retrieves the most popular match for the given title.
// Year and mediaType are optional filters; pass 0 and "" to omit.
func (c *Client) GetByTitle(ctx context.Context, title string, year int, mediaType MediaType, plot PlotLength) (*Title, error) {
	params := url.Values{"t": {title}}
	if year > 0 {
		params.Set("y", strconv.Itoa(year))
	}
	if mediaType != "" {
		params.Set("type", string(mediaType))
	}
	if plot != "" {
		params.Set("plot", string(plot))
	}
	var out Title
	if err := c.get(ctx, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Search searches for titles matching the given query string.
// Year, mediaType, and page are optional filters; pass 0, "", and 0 to omit.
func (c *Client) Search(ctx context.Context, query string, year int, mediaType MediaType, page int) (*SearchResponse, error) {
	params := url.Values{"s": {query}}
	if year > 0 {
		params.Set("y", strconv.Itoa(year))
	}
	if mediaType != "" {
		params.Set("type", string(mediaType))
	}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	var out SearchResponse
	if err := c.get(ctx, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSeason retrieves all episodes for a given season of a series.
// The series is identified by its IMDb ID.
func (c *Client) GetSeason(ctx context.Context, imdbID string, season int) (*SeasonResponse, error) {
	params := url.Values{
		"i":      {imdbID},
		"Season": {strconv.Itoa(season)},
	}
	var out SeasonResponse
	if err := c.get(ctx, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEpisode retrieves a specific episode of a series.
// The series is identified by its IMDb ID.
func (c *Client) GetEpisode(ctx context.Context, imdbID string, season, episode int) (*Title, error) {
	params := url.Values{
		"i":       {imdbID},
		"Season":  {strconv.Itoa(season)},
		"Episode": {strconv.Itoa(episode)},
	}
	var out Title
	if err := c.get(ctx, params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
