package tvdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	defaultBaseURL   = "https://api4.thetvdb.com/v4"
	defaultTimeout   = 30 * time.Second
	defaultUserAgent = "goenvoy/0.0.1"
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

// WithUserAgent sets the User-Agent header for all requests.
func WithUserAgent(ua string) Option {
	return func(cl *Client) { cl.userAgent = ua }
}

// WithBaseURL overrides the default TheTVDB API base URL.
func WithBaseURL(u string) Option {
	return func(cl *Client) { cl.rawBaseURL = u }
}

// Client is a TheTVDB API v4 client.
type Client struct {
	apiKey     string
	pin        string
	rawBaseURL string
	httpClient *http.Client
	userAgent  string

	mu    sync.Mutex
	token string
}

// New creates a TheTVDB [Client] using the given API key.
// An optional PIN can be provided for subscriber-level access.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		rawBaseURL: defaultBaseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
		userAgent:  defaultUserAgent,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// WithPIN returns an [Option] that sets a subscriber PIN.
func WithPIN(pin string) Option {
	return func(cl *Client) { cl.pin = pin }
}

// Login authenticates with TheTVDB and caches the bearer token.
// It is called automatically on the first API request if no token is cached.
func (c *Client) Login(ctx context.Context) error {
	body := LoginRequest{APIKey: c.apiKey}
	if c.pin != "" {
		body.PIN = c.pin
	}

	payload, err := json.Marshal(body) //nolint:gosec // API key is intentionally sent to the auth endpoint
	if err != nil {
		return fmt.Errorf("tvdb: marshal login body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.rawBaseURL+"/login", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("tvdb: create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("tvdb: login request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("tvdb: read login response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Message: "login failed"}
	}

	var result struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("tvdb: decode login response: %w", err)
	}

	c.mu.Lock()
	c.token = result.Data.Token
	c.mu.Unlock()
	return nil
}

func (c *Client) ensureToken(ctx context.Context) error {
	c.mu.Lock()
	hasToken := c.token != ""
	c.mu.Unlock()
	if hasToken {
		return nil
	}
	return c.Login(ctx)
}

func (c *Client) get(ctx context.Context, path string, dst any) error {
	if err := c.ensureToken(ctx); err != nil {
		return err
	}

	u, err := url.Parse(c.rawBaseURL + path)
	if err != nil {
		return fmt.Errorf("tvdb: parse URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("tvdb: create request: %w", err)
	}

	c.mu.Lock()
	token := c.token
	c.mu.Unlock()

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("tvdb: GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("tvdb: read response: %w", err)
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
			return fmt.Errorf("tvdb: decode response: %w", err)
		}
	}
	return nil
}

// GetSeries returns a series base record by its TheTVDB ID.
func (c *Client) GetSeries(ctx context.Context, id int) (*SeriesBase, error) {
	var out response[SeriesBase]
	if err := c.get(ctx, "/series/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetSeriesExtended returns a series extended record.
func (c *Client) GetSeriesExtended(ctx context.Context, id int) (*SeriesExtended, error) {
	var out response[SeriesExtended]
	if err := c.get(ctx, "/series/"+strconv.Itoa(id)+"/extended", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetSeriesEpisodes returns episodes for a series by season type.
func (c *Client) GetSeriesEpisodes(ctx context.Context, id int, seasonType string, page int) (*SeriesEpisodesResult, error) {
	var out response[SeriesEpisodesResult]
	path := "/series/" + strconv.Itoa(id) + "/episodes/" + url.PathEscape(seasonType) + "?page=" + strconv.Itoa(page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetSeriesArtworks returns artworks for a series.
func (c *Client) GetSeriesArtworks(ctx context.Context, id int) (*SeriesExtended, error) {
	var out response[SeriesExtended]
	if err := c.get(ctx, "/series/"+strconv.Itoa(id)+"/artworks", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetSeriesNextAired returns a series base record with the nextAired field.
func (c *Client) GetSeriesNextAired(ctx context.Context, id int) (*SeriesBase, error) {
	var out response[SeriesBase]
	if err := c.get(ctx, "/series/"+strconv.Itoa(id)+"/nextAired", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetSeriesTranslation returns a series translation for the given language.
func (c *Client) GetSeriesTranslation(ctx context.Context, id int, language string) (*Translation, error) {
	var out response[Translation]
	path := "/series/" + strconv.Itoa(id) + "/translations/" + url.PathEscape(language)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetMovie returns a movie base record by its TheTVDB ID.
func (c *Client) GetMovie(ctx context.Context, id int) (*MovieBase, error) {
	var out response[MovieBase]
	if err := c.get(ctx, "/movies/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetMovieExtended returns a movie extended record.
func (c *Client) GetMovieExtended(ctx context.Context, id int) (*MovieExtended, error) {
	var out response[MovieExtended]
	if err := c.get(ctx, "/movies/"+strconv.Itoa(id)+"/extended", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetMovieTranslation returns a movie translation for the given language.
func (c *Client) GetMovieTranslation(ctx context.Context, id int, language string) (*Translation, error) {
	var out response[Translation]
	path := "/movies/" + strconv.Itoa(id) + "/translations/" + url.PathEscape(language)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetEpisode returns an episode base record.
func (c *Client) GetEpisode(ctx context.Context, id int) (*EpisodeBase, error) {
	var out response[EpisodeBase]
	if err := c.get(ctx, "/episodes/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetEpisodeExtended returns an episode extended record.
func (c *Client) GetEpisodeExtended(ctx context.Context, id int) (*EpisodeExtended, error) {
	var out response[EpisodeExtended]
	if err := c.get(ctx, "/episodes/"+strconv.Itoa(id)+"/extended", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetEpisodeTranslation returns an episode translation for the given language.
func (c *Client) GetEpisodeTranslation(ctx context.Context, id int, language string) (*Translation, error) {
	var out response[Translation]
	path := "/episodes/" + strconv.Itoa(id) + "/translations/" + url.PathEscape(language)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetSeason returns a season base record.
func (c *Client) GetSeason(ctx context.Context, id int) (*SeasonBase, error) {
	var out response[SeasonBase]
	if err := c.get(ctx, "/seasons/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetSeasonExtended returns a season extended record.
func (c *Client) GetSeasonExtended(ctx context.Context, id int) (*SeasonExtended, error) {
	var out response[SeasonExtended]
	if err := c.get(ctx, "/seasons/"+strconv.Itoa(id)+"/extended", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetSeasonTranslation returns a season translation for the given language.
func (c *Client) GetSeasonTranslation(ctx context.Context, id int, language string) (*Translation, error) {
	var out response[Translation]
	path := "/seasons/" + strconv.Itoa(id) + "/translations/" + url.PathEscape(language)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetPerson returns a person base record.
func (c *Client) GetPerson(ctx context.Context, id int) (*PersonBase, error) {
	var out response[PersonBase]
	if err := c.get(ctx, "/people/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetPersonExtended returns a person extended record.
func (c *Client) GetPersonExtended(ctx context.Context, id int) (*PersonExtended, error) {
	var out response[PersonExtended]
	if err := c.get(ctx, "/people/"+strconv.Itoa(id)+"/extended", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetPersonTranslation returns a person translation for the given language.
func (c *Client) GetPersonTranslation(ctx context.Context, id int, language string) (*Translation, error) {
	var out response[Translation]
	path := "/people/" + strconv.Itoa(id) + "/translations/" + url.PathEscape(language)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetArtwork returns an artwork base record.
func (c *Client) GetArtwork(ctx context.Context, id int) (*ArtworkBase, error) {
	var out response[ArtworkBase]
	if err := c.get(ctx, "/artwork/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetArtworkExtended returns an artwork extended record.
func (c *Client) GetArtworkExtended(ctx context.Context, id int) (*ArtworkExtended, error) {
	var out response[ArtworkExtended]
	if err := c.get(ctx, "/artwork/"+strconv.Itoa(id)+"/extended", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetArtworkTypes returns all artwork type records.
func (c *Client) GetArtworkTypes(ctx context.Context) ([]ArtworkType, error) {
	var out response[[]ArtworkType]
	if err := c.get(ctx, "/artwork/types", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Search searches for series, movies, people, and companies.
func (c *Client) Search(ctx context.Context, query string, params *SearchParams) ([]SearchResult, error) {
	q := url.Values{}
	q.Set("query", query)
	if params != nil {
		if params.Type != "" {
			q.Set("type", params.Type)
		}
		if params.Year > 0 {
			q.Set("year", strconv.Itoa(params.Year))
		}
		if params.Company != "" {
			q.Set("company", params.Company)
		}
		if params.Country != "" {
			q.Set("country", params.Country)
		}
		if params.Director != "" {
			q.Set("director", params.Director)
		}
		if params.Language != "" {
			q.Set("language", params.Language)
		}
		if params.Network != "" {
			q.Set("network", params.Network)
		}
		if params.RemoteID != "" {
			q.Set("remote_id", params.RemoteID)
		}
		if params.Offset > 0 {
			q.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
	}

	var out response[[]SearchResult]
	if err := c.get(ctx, "/search?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// SearchByRemoteID searches by external remote ID (e.g. IMDB ID).
func (c *Client) SearchByRemoteID(ctx context.Context, remoteID string) ([]SearchByRemoteIDResult, error) {
	var out response[[]SearchByRemoteIDResult]
	if err := c.get(ctx, "/search/remoteid/"+url.PathEscape(remoteID), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetGenres returns all genre records.
func (c *Client) GetGenres(ctx context.Context) ([]Genre, error) {
	var out response[[]Genre]
	if err := c.get(ctx, "/genres", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetLanguages returns all language records.
func (c *Client) GetLanguages(ctx context.Context) ([]Language, error) {
	var out response[[]Language]
	if err := c.get(ctx, "/languages", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetContentRatings returns all content rating records.
func (c *Client) GetContentRatings(ctx context.Context) ([]ContentRating, error) {
	var out response[[]ContentRating]
	if err := c.get(ctx, "/content/ratings", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetUpdates returns entity updates since the given UNIX timestamp.
func (c *Client) GetUpdates(ctx context.Context, since int64, params *UpdatesParams) ([]EntityUpdate, error) {
	q := url.Values{}
	q.Set("since", strconv.FormatInt(since, 10))
	if params != nil {
		if params.Type != "" {
			q.Set("type", params.Type)
		}
		if params.Action != "" {
			q.Set("action", params.Action)
		}
		if params.Page > 0 {
			q.Set("page", strconv.Itoa(params.Page))
		}
	}

	var out response[[]EntityUpdate]
	if err := c.get(ctx, "/updates?"+q.Encode(), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetCharacter returns a character record by its ID.
func (c *Client) GetCharacter(ctx context.Context, id int) (*Character, error) {
	var out response[Character]
	if err := c.get(ctx, "/characters/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}
