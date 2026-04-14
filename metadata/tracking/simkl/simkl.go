package simkl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/golusoris/goenvoy/metadata"
)

const (
	defaultBaseURL = "https://api.simkl.com"
	defaultCalURL  = "https://data.simkl.in"
	defaultFilter  = "all"
)

// Client is a Simkl API client.
type Client struct {
	*metadata.BaseClient
	clientID     string
	clientSecret string
	calBaseURL   string
	onToken      TokenCallback

	mu          sync.RWMutex
	accessToken string
}

// SetClientSecret sets the client secret needed for OAuth2 flows.
func (c *Client) SetClientSecret(secret string) { c.clientSecret = secret }

// SetCalendarURL overrides the default Simkl calendar (CDN) base URL.
func (c *Client) SetCalendarURL(u string) { c.calBaseURL = u }

// SetAccessToken sets a pre-existing OAuth2 access token.
func (c *Client) SetAccessToken(token string) {
	c.mu.Lock()
	c.accessToken = token
	c.mu.Unlock()
}

// TokenCallback is called when a new access token is obtained.
type TokenCallback func(accessToken string)

// SetTokenCallback sets a callback invoked when a new token is obtained.
func (c *Client) SetTokenCallback(cb TokenCallback) { c.onToken = cb }

// New creates a Simkl [Client] using the given client ID (API key).
func New(clientID string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "simkl", opts...)
	return &Client{
		BaseClient: bc,
		clientID:   clientID,
		calBaseURL: defaultCalURL,
	}
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int    `json:"-"`
	Error_     string `json:"error"` //nolint:revive // field name conflicts with Error() method below
	Code       int    `json:"code"`
	Message    string `json:"message"`
	// RawBody holds the raw response body when the error could not be parsed as JSON.
	RawBody string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("simkl: HTTP %d: %s: %s", e.StatusCode, e.Error_, e.Message)
	}
	if e.Error_ != "" {
		return fmt.Sprintf("simkl: HTTP %d: %s", e.StatusCode, e.Error_)
	}
	if e.RawBody != "" {
		return fmt.Sprintf("simkl: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("simkl: HTTP %d", e.StatusCode)
}

func (c *Client) doGet(ctx context.Context, baseURL, path string, params url.Values, dst any) error {
	u, err := url.Parse(baseURL + path)
	if err != nil {
		return fmt.Errorf("simkl: parse URL: %w", err)
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("simkl: create request: %w", err)
	}

	req.Header.Set("Simkl-Api-Key", c.clientID)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("simkl: GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("simkl: read response: %w", err)
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
			return fmt.Errorf("simkl: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) get(ctx context.Context, path string, params url.Values, dst any) error {
	return c.doGet(ctx, c.BaseURL(), path, params, dst)
}

func (c *Client) getCal(ctx context.Context, path string, dst any) error {
	return c.doGet(ctx, c.calBaseURL, path, nil, dst)
}

func pageParams(page, limit int) url.Values {
	p := url.Values{}
	if page > 0 {
		p.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		p.Set("limit", strconv.Itoa(limit))
	}
	return p
}

// Movies.

// GetMovie returns detail information about a movie.
// The id can be a Simkl ID or an IMDB ID.
func (c *Client) GetMovie(ctx context.Context, id string) (*Movie, error) {
	var out Movie
	p := url.Values{}
	p.Set("extended", "full")
	if err := c.get(ctx, "/movies/"+url.PathEscape(id), p, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TrendingMovies returns trending movies, optionally filtered by time interval.
// Valid intervals: "today", "week", "month" (empty string for default).
func (c *Client) TrendingMovies(ctx context.Context, interval string) ([]TrendingMovie, error) {
	path := "/movies/trending"
	if interval != "" {
		path += "/" + url.PathEscape(interval)
	}
	var out []TrendingMovie
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TV Shows.

// GetShow returns detail information about a TV show.
func (c *Client) GetShow(ctx context.Context, id string) (*Show, error) {
	var out Show
	p := url.Values{}
	p.Set("extended", "full")
	if err := c.get(ctx, "/tv/"+url.PathEscape(id), p, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetShowEpisodes returns episodes for a TV show.
func (c *Client) GetShowEpisodes(ctx context.Context, id string) ([]Episode, error) {
	var out []Episode
	p := url.Values{}
	p.Set("extended", "full")
	if err := c.get(ctx, "/tv/episodes/"+url.PathEscape(id), p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TrendingShows returns trending TV shows.
func (c *Client) TrendingShows(ctx context.Context, interval string) ([]TrendingShow, error) {
	path := "/tv/trending"
	if interval != "" {
		path += "/" + url.PathEscape(interval)
	}
	var out []TrendingShow
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ShowGenres returns shows filtered by genre and other criteria.
func (c *Client) ShowGenres(ctx context.Context, genre string, page, limit int) ([]GenreItem, error) {
	if genre == "" {
		genre = defaultFilter
	}
	var out []GenreItem
	if err := c.get(ctx, "/tv/genres/"+url.PathEscape(genre), pageParams(page, limit), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ShowPremieres returns latest TV show premieres.
// param should be "new" or "soon".
func (c *Client) ShowPremieres(ctx context.Context, param string, page, limit int) ([]PremiereItem, error) {
	if param == "" {
		param = "new"
	}
	var out []PremiereItem
	if err := c.get(ctx, "/tv/premieres/"+url.PathEscape(param), pageParams(page, limit), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AiringShows returns currently airing TV shows.
func (c *Client) AiringShows(ctx context.Context, date, sort string) ([]AiringItem, error) {
	p := url.Values{}
	if date != "" {
		p.Set("date", date)
	}
	if sort != "" {
		p.Set("sort", sort)
	}
	var out []AiringItem
	if err := c.get(ctx, "/tv/airing", p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BestShows returns the best of all TV shows.
// filter: "year", "month", "all", "voted", "watched".
func (c *Client) BestShows(ctx context.Context, filter string) ([]BestItem, error) {
	if filter == "" {
		filter = defaultFilter
	}
	var out []BestItem
	if err := c.get(ctx, "/tv/best/"+url.PathEscape(filter), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Anime.

// GetAnime returns detail information about an anime.
func (c *Client) GetAnime(ctx context.Context, id string) (*Anime, error) {
	var out Anime
	p := url.Values{}
	p.Set("extended", "full")
	if err := c.get(ctx, "/anime/"+url.PathEscape(id), p, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAnimeEpisodes returns episodes for an anime.
func (c *Client) GetAnimeEpisodes(ctx context.Context, id string) ([]Episode, error) {
	var out []Episode
	p := url.Values{}
	p.Set("extended", "full")
	if err := c.get(ctx, "/anime/episodes/"+url.PathEscape(id), p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TrendingAnime returns trending anime.
func (c *Client) TrendingAnime(ctx context.Context, interval string) ([]TrendingAnime, error) {
	path := "/anime/trending"
	if interval != "" {
		path += "/" + url.PathEscape(interval)
	}
	var out []TrendingAnime
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AnimeGenres returns anime filtered by genre.
func (c *Client) AnimeGenres(ctx context.Context, genre string, page, limit int) ([]GenreItem, error) {
	if genre == "" {
		genre = defaultFilter
	}
	var out []GenreItem
	if err := c.get(ctx, "/anime/genres/"+url.PathEscape(genre), pageParams(page, limit), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AnimePremieres returns latest anime premieres.
func (c *Client) AnimePremieres(ctx context.Context, param string, page, limit int) ([]PremiereItem, error) {
	if param == "" {
		param = "new"
	}
	var out []PremiereItem
	if err := c.get(ctx, "/anime/premieres/"+url.PathEscape(param), pageParams(page, limit), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AiringAnime returns currently airing anime.
func (c *Client) AiringAnime(ctx context.Context, date, sort string) ([]AiringItem, error) {
	p := url.Values{}
	if date != "" {
		p.Set("date", date)
	}
	if sort != "" {
		p.Set("sort", sort)
	}
	var out []AiringItem
	if err := c.get(ctx, "/anime/airing", p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BestAnime returns the best of all anime.
func (c *Client) BestAnime(ctx context.Context, filter string) ([]BestItem, error) {
	if filter == "" {
		filter = defaultFilter
	}
	var out []BestItem
	if err := c.get(ctx, "/anime/best/"+url.PathEscape(filter), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Search.

// SearchByID resolves an external ID to Simkl entries.
// Supported ID types: imdb, tvdb, tmdb, anidb, mal, anilist, hulu, netflix,
// crunchyroll, kitsu, livechart, anisearch, animeplanet, traktslug, letterboxd.
func (c *Client) SearchByID(ctx context.Context, idType, idValue string) ([]SearchIDResult, error) {
	p := url.Values{}
	p.Set(idType, idValue)
	var out []SearchIDResult
	if err := c.get(ctx, "/search/id", p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SearchText searches for movies, shows, or anime by text query.
// searchType should be "tv", "anime", or "movie".
func (c *Client) SearchText(ctx context.Context, searchType, query string, page, limit int) ([]SearchResult, error) {
	p := pageParams(page, limit)
	p.Set("q", query)
	p.Set("extended", "full")
	var out []SearchResult
	if err := c.get(ctx, "/search/"+url.PathEscape(searchType), p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Calendar.
// Calendar endpoints use a separate CDN URL and do not require API authentication.

// CalendarShows returns the next 33 days of TV show releases.
func (c *Client) CalendarShows(ctx context.Context) ([]CalendarShow, error) {
	var out []CalendarShow
	if err := c.getCal(ctx, "/calendar/tv.json", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarShowsMonth returns TV show releases for a specific month.
func (c *Client) CalendarShowsMonth(ctx context.Context, year, month int) ([]CalendarShow, error) {
	path := fmt.Sprintf("/calendar/%d/%d/tv.json", year, month)
	var out []CalendarShow
	if err := c.getCal(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarAnime returns the next 33 days of anime releases.
func (c *Client) CalendarAnime(ctx context.Context) ([]CalendarAnime, error) {
	var out []CalendarAnime
	if err := c.getCal(ctx, "/calendar/anime.json", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarAnimeMonth returns anime releases for a specific month.
func (c *Client) CalendarAnimeMonth(ctx context.Context, year, month int) ([]CalendarAnime, error) {
	path := fmt.Sprintf("/calendar/%d/%d/anime.json", year, month)
	var out []CalendarAnime
	if err := c.getCal(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarMovies returns the next 33 days of movie releases.
func (c *Client) CalendarMovies(ctx context.Context) ([]CalendarMovie, error) {
	var out []CalendarMovie
	if err := c.getCal(ctx, "/calendar/movie_release.json", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarMoviesMonth returns movie releases for a specific month.
func (c *Client) CalendarMoviesMonth(ctx context.Context, year, month int) ([]CalendarMovie, error) {
	path := fmt.Sprintf("/calendar/%d/%d/movie_release.json", year, month)
	var out []CalendarMovie
	if err := c.getCal(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// OAuth2.

func (c *Client) post(ctx context.Context, path string, body, dst any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("simkl: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL()+path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("simkl: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Simkl-Api-Key", c.clientID)
	req.Header.Set("User-Agent", c.UserAgent())

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("simkl: POST %s: %w", path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("simkl: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if jsonErr := json.Unmarshal(respBody, apiErr); jsonErr != nil {
			apiErr.RawBody = string(respBody)
		}
		return apiErr
	}

	if dst != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dst); err != nil {
			return fmt.Errorf("simkl: decode response: %w", err)
		}
	}
	return nil
}

// GetDeviceCode starts the Simkl PIN-based device code flow.
// Display the returned UserCode and VerificationURL to the user.
func (c *Client) GetDeviceCode(ctx context.Context) (*DeviceCode, error) {
	var out DeviceCode
	err := c.post(ctx, "/oauth/pin", map[string]string{
		"client_id": c.clientID,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// PollDeviceToken polls for the PIN/device token after the user has authorized the app.
// It blocks until the token is obtained, the code expires, or the context is canceled.
func (c *Client) PollDeviceToken(ctx context.Context, code *DeviceCode) (string, error) {
	interval := code.Interval
	if interval < 5 {
		interval = 5
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	deadline := time.Now().Add(time.Duration(code.ExpiresIn) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("simkl: device poll: %w", ctx.Err())
		case <-ticker.C:
			if time.Now().After(deadline) {
				return "", errors.New("simkl: device code expired")
			}
			var result struct {
				Result      string `json:"result"`
				AccessToken string `json:"access_token"`
				Message     string `json:"message"`
			}
			err := c.post(ctx, "/oauth/pin/"+url.PathEscape(code.UserCode), map[string]string{
				"client_id": c.clientID,
			}, &result)
			if err != nil {
				return "", err
			}
			switch result.Result {
			case "OK":
				c.mu.Lock()
				c.accessToken = result.AccessToken
				c.mu.Unlock()
				if c.onToken != nil {
					c.onToken(result.AccessToken)
				}
				return result.AccessToken, nil
			case "KO":
				// Not yet authorized, keep polling.
				continue
			default:
				continue
			}
		}
	}
}

// ExchangeCode exchanges an authorization code for an access token.
func (c *Client) ExchangeCode(ctx context.Context, code, redirectURI string) (string, error) {
	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	err := c.post(ctx, "/oauth/token", map[string]string{
		"code":          code,
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"redirect_uri":  redirectURI,
		"grant_type":    "authorization_code",
	}, &result)
	if err != nil {
		return "", err
	}
	c.mu.Lock()
	c.accessToken = result.AccessToken
	c.mu.Unlock()
	if c.onToken != nil {
		c.onToken(result.AccessToken)
	}
	return result.AccessToken, nil
}

func (c *Client) del(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL()+path, http.NoBody)
	if err != nil {
		return fmt.Errorf("simkl: create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Simkl-Api-Key", c.clientID)
	req.Header.Set("User-Agent", c.UserAgent())

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("simkl: DELETE %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("simkl: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if jsonErr := json.Unmarshal(body, apiErr); jsonErr != nil {
			apiErr.RawBody = string(body)
		}
		return apiErr
	}
	return nil
}

// Ratings.

// GetRatingByID returns ratings for a movie, TV show, or anime by Simkl ID.
// fields can include: "rank", "droprate", "simkl", "ext", "has_trailer", "reactions", "year".
func (c *Client) GetRatingByID(ctx context.Context, simklID int, fields string) (*RatingInfo, error) {
	p := url.Values{}
	p.Set("simkl", strconv.Itoa(simklID))
	if fields != "" {
		p.Set("fields", fields)
	}
	var out RatingInfo
	if err := c.get(ctx, "/ratings", p, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWatchlistRatings returns ratings for items in a user's watchlist.
// mediaType can be "tv", "anime", or "movies".
// status can be "all", "watching", "plantowatch", "completed", "dropped", "hold".
func (c *Client) GetWatchlistRatings(ctx context.Context, mediaType, status, fields string) ([]RatingInfo, error) {
	p := url.Values{}
	if status != "" {
		p.Set("user_watchlist", status)
	}
	if fields != "" {
		p.Set("fields", fields)
	}
	var out []RatingInfo
	path := "/ratings"
	if mediaType != "" {
		path += "/" + url.PathEscape(mediaType) + "/"
	}
	if err := c.get(ctx, path, p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Search (additional).

// SearchRandom finds a random item based on filters.
func (c *Client) SearchRandom(ctx context.Context, params *RandomSearchParams) (*RandomResult, error) {
	var out RandomResult
	if err := c.post(ctx, "/search/random", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Movie Genres.

// MovieGenres returns movies filtered by genre and other criteria.
func (c *Client) MovieGenres(ctx context.Context, genre string, page, limit int) ([]GenreItem, error) {
	if genre == "" {
		genre = defaultFilter
	}
	var out []GenreItem
	if err := c.get(ctx, "/movies/genres/"+url.PathEscape(genre), pageParams(page, limit), &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Movie Best.

// BestMovies returns the best of all movies.
func (c *Client) BestMovies(ctx context.Context, filter string) ([]BestItem, error) {
	if filter == "" {
		filter = defaultFilter
	}
	var out []BestItem
	if err := c.get(ctx, "/movies/best/"+url.PathEscape(filter), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Scrobble.

// ScrobbleStart starts a scrobble session for the given item.
func (c *Client) ScrobbleStart(ctx context.Context, item *ScrobbleRequest) (*ScrobbleResponse, error) {
	var out ScrobbleResponse
	if err := c.post(ctx, "/scrobble/start", item, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ScrobblePause pauses the current scrobble session.
func (c *Client) ScrobblePause(ctx context.Context, item *ScrobbleRequest) (*ScrobbleResponse, error) {
	var out ScrobbleResponse
	if err := c.post(ctx, "/scrobble/pause", item, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ScrobbleStop stops the current scrobble session.
// If progress >= 80%, the item is marked as watched.
func (c *Client) ScrobbleStop(ctx context.Context, item *ScrobbleRequest) (*ScrobbleResponse, error) {
	var out ScrobbleResponse
	if err := c.post(ctx, "/scrobble/stop", item, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ScrobbleCheckin checks in to an item without real-time tracking.
func (c *Client) ScrobbleCheckin(ctx context.Context, item *ScrobbleRequest) (*ScrobbleResponse, error) {
	var out ScrobbleResponse
	if err := c.post(ctx, "/scrobble/checkin", item, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Sync.

// GetLastActivity returns the user's last activity timestamps for efficient syncing.
func (c *Client) GetLastActivity(ctx context.Context) (*LastActivity, error) {
	var out LastActivity
	if err := c.get(ctx, "/sync/activities", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAllItems returns all items in the user's watchlist.
// mediaType can be empty, "shows", "movies", or "anime".
// status can be empty, "watching", "plantowatch", "completed", "hold", or "dropped".
// dateFrom is an optional ISO 8601 timestamp for incremental sync.
func (c *Client) GetAllItems(ctx context.Context, mediaType, status, extended, dateFrom string) (*WatchlistResponse, error) {
	path := "/sync/all-items/"
	if mediaType != "" {
		path += url.PathEscape(mediaType) + "/"
	}
	if status != "" {
		path += url.PathEscape(status)
	}
	p := url.Values{}
	if extended != "" {
		p.Set("extended", extended)
	}
	if dateFrom != "" {
		p.Set("date_from", dateFrom)
	}
	var out WatchlistResponse
	if err := c.get(ctx, path, p, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddToHistory adds items to the user's watched history.
func (c *Client) AddToHistory(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/history", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveFromHistory removes items from the user's watched history.
func (c *Client) RemoveFromHistory(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/history/remove", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSyncRatings returns all user ratings, optionally filtered by type and rating.
// mediaType: "", "shows", "movies", or "anime". rating: "" or "1" through "10".
// dateFrom: optional ISO 8601 timestamp.
func (c *Client) GetSyncRatings(ctx context.Context, mediaType, rating, dateFrom string) (*WatchlistResponse, error) {
	path := "/sync/ratings/"
	if mediaType != "" {
		path += url.PathEscape(mediaType) + "/"
	}
	if rating != "" {
		path += url.PathEscape(rating)
	}
	p := url.Values{}
	if dateFrom != "" {
		p.Set("date_from", dateFrom)
	}
	var out WatchlistResponse
	if err := c.get(ctx, path, p, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddRatings adds ratings for movies and/or shows.
func (c *Client) AddRatings(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/ratings", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveRatings removes ratings for movies and/or shows.
func (c *Client) RemoveRatings(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/ratings/remove", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddToList adds items to a specific Simkl watchlist (e.g., plantowatch, watching).
func (c *Client) AddToList(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/add-to-list", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPlaybacks retrieves paused playback sessions.
// mediaType: "", "movies", or "episodes".
func (c *Client) GetPlaybacks(ctx context.Context, mediaType string) ([]PlaybackSession, error) {
	path := "/sync/playback/"
	if mediaType != "" {
		path += url.PathEscape(mediaType)
	}
	var out []PlaybackSession
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeletePlayback deletes a specific paused playback session by its ID.
func (c *Client) DeletePlayback(ctx context.Context, id int64) error {
	return c.del(ctx, "/sync/playback/"+strconv.FormatInt(id, 10))
}

// CheckIfWatched checks whether items have been watched by the user.
func (c *Client) CheckIfWatched(ctx context.Context, items []WatchedCheckItem, extended string) ([]WatchedCheckResult, error) {
	p := url.Values{}
	if extended != "" {
		p.Set("extended", extended)
	}
	u, err := url.Parse(c.BaseURL() + "/sync/watched/")
	if err != nil {
		return nil, fmt.Errorf("simkl: parse URL: %w", err)
	}
	if len(p) > 0 {
		u.RawQuery = p.Encode()
	}

	payload, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("simkl: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("simkl: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Simkl-Api-Key", c.clientID)
	req.Header.Set("User-Agent", c.UserAgent())

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("simkl: POST /sync/watched/: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("simkl: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if jsonErr := json.Unmarshal(body, apiErr); jsonErr != nil {
			apiErr.RawBody = string(body)
		}
		return nil, apiErr
	}

	var out []WatchedCheckResult
	if len(body) > 0 {
		if err := json.Unmarshal(body, &out); err != nil {
			return nil, fmt.Errorf("simkl: decode response: %w", err)
		}
	}
	return out, nil
}

// Users.

// GetUserStats returns watched statistics for a user.
func (c *Client) GetUserStats(ctx context.Context, userID int) (*UserStats, error) {
	var out UserStats
	if err := c.get(ctx, fmt.Sprintf("/users/%d/stats", userID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUserSettings returns the authenticated user's settings.
func (c *Client) GetUserSettings(ctx context.Context) (*UserSettings, error) {
	var out UserSettings
	if err := c.get(ctx, "/users/settings", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLastWatchedArts returns images from a user's last watched item.
func (c *Client) GetLastWatchedArts(ctx context.Context, userID int) (*LastWatchedArt, error) {
	var out LastWatchedArt
	if err := c.get(ctx, fmt.Sprintf("/users/recently-watched-background/%d", userID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
