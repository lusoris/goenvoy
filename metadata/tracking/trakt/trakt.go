package trakt

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
)

const (
	defaultBaseURL   = "https://api.trakt.tv"
	defaultTimeout   = 30 * time.Second
	defaultUserAgent = "goenvoy/0.0.1"
	apiVersion       = "2"
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

// WithUserAgent sets the User-Agent header sent with every request.
func WithUserAgent(ua string) Option {
	return func(cl *Client) { cl.userAgent = ua }
}

// WithBaseURL overrides the default Trakt API base URL.
// The URL must not have a trailing slash.
func WithBaseURL(u string) Option {
	return func(cl *Client) { cl.rawBaseURL = u }
}

// WithClientSecret sets the client secret needed for OAuth2 flows.
func WithClientSecret(secret string) Option {
	return func(cl *Client) { cl.clientSecret = secret }
}

// WithAccessToken sets a pre-existing OAuth2 access token for user-authenticated requests.
func WithAccessToken(token string) Option {
	return func(cl *Client) {
		cl.mu.Lock()
		cl.accessToken = token
		cl.mu.Unlock()
	}
}

// WithRefreshToken sets a pre-existing OAuth2 refresh token.
func WithRefreshToken(token string) Option {
	return func(cl *Client) {
		cl.mu.Lock()
		cl.refreshToken = token
		cl.mu.Unlock()
	}
}

// TokenCallback is called whenever a new token pair is obtained (via refresh or exchange).
// Store the tokens persistently so they survive restarts.
type TokenCallback func(token Token)

// WithTokenCallback sets a callback invoked whenever tokens are refreshed or exchanged.
func WithTokenCallback(cb TokenCallback) Option {
	return func(cl *Client) { cl.onToken = cb }
}

// Client is a Trakt API v2 client.
type Client struct {
	clientID     string
	clientSecret string
	rawBaseURL   string
	httpClient   *http.Client
	userAgent    string
	onToken      TokenCallback

	mu           sync.RWMutex
	accessToken  string
	refreshToken string
}

// New creates a Trakt [Client] using the given client ID (API key).
func New(clientID string, opts ...Option) *Client {
	c := &Client{
		clientID:   clientID,
		rawBaseURL: defaultBaseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
		userAgent:  defaultUserAgent,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode  int    `json:"-"`
	Error_      string `json:"error"`
	Description string `json:"error_description"`
	// RawBody holds the raw response body when the error could not be parsed as JSON.
	RawBody string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("trakt: HTTP %d: %s: %s", e.StatusCode, e.Error_, e.Description)
	}
	if e.Error_ != "" {
		return fmt.Sprintf("trakt: HTTP %d: %s", e.StatusCode, e.Error_)
	}
	if e.RawBody != "" {
		return fmt.Sprintf("trakt: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("trakt: HTTP %d", e.StatusCode)
}

// PaginationHeaders contains pagination information from response headers.
type PaginationHeaders struct {
	Page      int
	Limit     int
	PageCount int
	ItemCount int
}

func parsePaginationHeaders(h http.Header) PaginationHeaders {
	atoi := func(s string) int {
		v, _ := strconv.Atoi(s)
		return v
	}
	return PaginationHeaders{
		Page:      atoi(h.Get("X-Pagination-Page")),
		Limit:     atoi(h.Get("X-Pagination-Limit")),
		PageCount: atoi(h.Get("X-Pagination-Page-Count")),
		ItemCount: atoi(h.Get("X-Pagination-Item-Count")),
	}
}

func (c *Client) get(ctx context.Context, path string, params url.Values, dst any) (*PaginationHeaders, error) {
	u, err := url.Parse(c.rawBaseURL + path)
	if err != nil {
		return nil, fmt.Errorf("trakt: parse URL: %w", err)
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("trakt: create request: %w", err)
	}

	req.Header.Set("trakt-api-key", c.clientID)
	req.Header.Set("trakt-api-version", apiVersion)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("trakt: GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("trakt: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if err := json.Unmarshal(body, apiErr); err != nil {
			apiErr.RawBody = string(body)
		}
		return nil, apiErr
	}

	pg := parsePaginationHeaders(resp.Header)

	if dst != nil && len(body) > 0 {
		if err := json.Unmarshal(body, dst); err != nil {
			return nil, fmt.Errorf("trakt: decode response: %w", err)
		}
	}
	return &pg, nil
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

func extendedParams(page, limit int) url.Values {
	p := pageParams(page, limit)
	p.Set("extended", "full")
	return p
}

// Movies.

// GetMovie returns a single movie by its Trakt slug or ID.
func (c *Client) GetMovie(ctx context.Context, idOrSlug string) (*Movie, error) {
	var out Movie
	_, err := c.get(ctx, "/movies/"+url.PathEscape(idOrSlug), url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieAliases returns all title aliases for a movie.
func (c *Client) GetMovieAliases(ctx context.Context, idOrSlug string) ([]Alias, error) {
	var out []Alias
	_, err := c.get(ctx, "/movies/"+url.PathEscape(idOrSlug)+"/aliases", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetMovieReleases returns release information for a movie in a given country.
// Pass an empty country for all countries.
func (c *Client) GetMovieReleases(ctx context.Context, idOrSlug, country string) ([]MovieRelease, error) {
	path := "/movies/" + url.PathEscape(idOrSlug) + "/releases"
	if country != "" {
		path += "/" + url.PathEscape(country)
	}
	var out []MovieRelease
	_, err := c.get(ctx, path, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetMovieTranslations returns translations for a movie.
// Pass an empty language for all languages.
func (c *Client) GetMovieTranslations(ctx context.Context, idOrSlug, language string) ([]MovieTranslation, error) {
	path := "/movies/" + url.PathEscape(idOrSlug) + "/translations"
	if language != "" {
		path += "/" + url.PathEscape(language)
	}
	var out []MovieTranslation
	_, err := c.get(ctx, path, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetMoviePeople returns the cast and crew for a movie.
func (c *Client) GetMoviePeople(ctx context.Context, idOrSlug string) (*People, error) {
	var out People
	_, err := c.get(ctx, "/movies/"+url.PathEscape(idOrSlug)+"/people", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieRatings returns the rating and vote distribution for a movie.
func (c *Client) GetMovieRatings(ctx context.Context, idOrSlug string) (*Ratings, error) {
	var out Ratings
	_, err := c.get(ctx, "/movies/"+url.PathEscape(idOrSlug)+"/ratings", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieStats returns stats for a movie.
func (c *Client) GetMovieStats(ctx context.Context, idOrSlug string) (*Stats, error) {
	var out Stats
	_, err := c.get(ctx, "/movies/"+url.PathEscape(idOrSlug)+"/stats", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieStudios returns production studios for a movie.
func (c *Client) GetMovieStudios(ctx context.Context, idOrSlug string) ([]Studio, error) {
	var out []Studio
	_, err := c.get(ctx, "/movies/"+url.PathEscape(idOrSlug)+"/studios", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrendingMovies returns trending movies with pagination.
func (c *Client) TrendingMovies(ctx context.Context, page, limit int) ([]TrendingMovie, *PaginationHeaders, error) {
	var out []TrendingMovie
	pg, err := c.get(ctx, "/movies/trending", extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// PopularMovies returns popular movies with pagination.
func (c *Client) PopularMovies(ctx context.Context, page, limit int) ([]Movie, *PaginationHeaders, error) {
	var out []Movie
	pg, err := c.get(ctx, "/movies/popular", extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// MostPlayedMovies returns the most played movies for a given period (weekly, monthly, yearly, all).
func (c *Client) MostPlayedMovies(ctx context.Context, period string, page, limit int) ([]PlayedMovie, *PaginationHeaders, error) {
	var out []PlayedMovie
	pg, err := c.get(ctx, "/movies/played/"+url.PathEscape(period), extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// MostWatchedMovies returns the most watched movies for a given period.
func (c *Client) MostWatchedMovies(ctx context.Context, period string, page, limit int) ([]PlayedMovie, *PaginationHeaders, error) {
	var out []PlayedMovie
	pg, err := c.get(ctx, "/movies/watched/"+url.PathEscape(period), extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// AnticipatedMovies returns the most anticipated movies.
func (c *Client) AnticipatedMovies(ctx context.Context, page, limit int) ([]AnticipatedMovie, *PaginationHeaders, error) {
	var out []AnticipatedMovie
	pg, err := c.get(ctx, "/movies/anticipated", extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// BoxOfficeMovies returns the current weekend box office.
func (c *Client) BoxOfficeMovies(ctx context.Context) ([]BoxOfficeMovie, error) {
	var out []BoxOfficeMovie
	_, err := c.get(ctx, "/movies/boxoffice", url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Shows.

// GetShow returns a single show by its Trakt slug or ID.
func (c *Client) GetShow(ctx context.Context, idOrSlug string) (*Show, error) {
	var out Show
	_, err := c.get(ctx, "/shows/"+url.PathEscape(idOrSlug), url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetShowAliases returns all title aliases for a show.
func (c *Client) GetShowAliases(ctx context.Context, idOrSlug string) ([]Alias, error) {
	var out []Alias
	_, err := c.get(ctx, "/shows/"+url.PathEscape(idOrSlug)+"/aliases", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetShowTranslations returns translations for a show.
func (c *Client) GetShowTranslations(ctx context.Context, idOrSlug, language string) ([]ShowTranslation, error) {
	path := "/shows/" + url.PathEscape(idOrSlug) + "/translations"
	if language != "" {
		path += "/" + url.PathEscape(language)
	}
	var out []ShowTranslation
	_, err := c.get(ctx, path, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetShowPeople returns the cast and crew for a show.
func (c *Client) GetShowPeople(ctx context.Context, idOrSlug string) (*People, error) {
	var out People
	_, err := c.get(ctx, "/shows/"+url.PathEscape(idOrSlug)+"/people", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetShowRatings returns the rating for a show.
func (c *Client) GetShowRatings(ctx context.Context, idOrSlug string) (*Ratings, error) {
	var out Ratings
	_, err := c.get(ctx, "/shows/"+url.PathEscape(idOrSlug)+"/ratings", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetShowStats returns stats for a show.
func (c *Client) GetShowStats(ctx context.Context, idOrSlug string) (*Stats, error) {
	var out Stats
	_, err := c.get(ctx, "/shows/"+url.PathEscape(idOrSlug)+"/stats", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetShowStudios returns production studios for a show.
func (c *Client) GetShowStudios(ctx context.Context, idOrSlug string) ([]Studio, error) {
	var out []Studio
	_, err := c.get(ctx, "/shows/"+url.PathEscape(idOrSlug)+"/studios", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrendingShows returns trending shows with pagination.
func (c *Client) TrendingShows(ctx context.Context, page, limit int) ([]TrendingShow, *PaginationHeaders, error) {
	var out []TrendingShow
	pg, err := c.get(ctx, "/shows/trending", extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// PopularShows returns popular shows with pagination.
func (c *Client) PopularShows(ctx context.Context, page, limit int) ([]Show, *PaginationHeaders, error) {
	var out []Show
	pg, err := c.get(ctx, "/shows/popular", extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// MostPlayedShows returns the most played shows for a given period.
func (c *Client) MostPlayedShows(ctx context.Context, period string, page, limit int) ([]PlayedShow, *PaginationHeaders, error) {
	var out []PlayedShow
	pg, err := c.get(ctx, "/shows/played/"+url.PathEscape(period), extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// MostWatchedShows returns the most watched shows for a given period.
func (c *Client) MostWatchedShows(ctx context.Context, period string, page, limit int) ([]PlayedShow, *PaginationHeaders, error) {
	var out []PlayedShow
	pg, err := c.get(ctx, "/shows/watched/"+url.PathEscape(period), extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// AnticipatedShows returns the most anticipated shows.
func (c *Client) AnticipatedShows(ctx context.Context, page, limit int) ([]AnticipatedShow, *PaginationHeaders, error) {
	var out []AnticipatedShow
	pg, err := c.get(ctx, "/shows/anticipated", extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// Seasons and episodes.

// GetShowSeasons returns all seasons for a show.
func (c *Client) GetShowSeasons(ctx context.Context, idOrSlug string) ([]Season, error) {
	var out []Season
	_, err := c.get(ctx, "/shows/"+url.PathEscape(idOrSlug)+"/seasons", url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetSeasonEpisodes returns all episodes for a specific season of a show.
func (c *Client) GetSeasonEpisodes(ctx context.Context, idOrSlug string, season int) ([]Episode, error) {
	var out []Episode
	path := "/shows/" + url.PathEscape(idOrSlug) + "/seasons/" + strconv.Itoa(season)
	_, err := c.get(ctx, path, url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetEpisode returns a single episode by show, season number, and episode number.
func (c *Client) GetEpisode(ctx context.Context, idOrSlug string, season, episode int) (*Episode, error) {
	var out Episode
	path := "/shows/" + url.PathEscape(idOrSlug) + "/seasons/" + strconv.Itoa(season) + "/episodes/" + strconv.Itoa(episode)
	_, err := c.get(ctx, path, url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEpisodeRatings returns the ratings for a specific episode.
func (c *Client) GetEpisodeRatings(ctx context.Context, idOrSlug string, season, episode int) (*Ratings, error) {
	var out Ratings
	path := "/shows/" + url.PathEscape(idOrSlug) + "/seasons/" + strconv.Itoa(season) + "/episodes/" + strconv.Itoa(episode) + "/ratings"
	_, err := c.get(ctx, path, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEpisodeStats returns the stats for a specific episode.
func (c *Client) GetEpisodeStats(ctx context.Context, idOrSlug string, season, episode int) (*Stats, error) {
	var out Stats
	path := "/shows/" + url.PathEscape(idOrSlug) + "/seasons/" + strconv.Itoa(season) + "/episodes/" + strconv.Itoa(episode) + "/stats"
	_, err := c.get(ctx, path, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// People.

// GetPerson returns a single person by their Trakt slug or ID.
func (c *Client) GetPerson(ctx context.Context, idOrSlug string) (*Person, error) {
	var out Person
	_, err := c.get(ctx, "/people/"+url.PathEscape(idOrSlug), url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Search.

// SearchText searches Trakt by text query.
// searchType can be "movie", "show", "episode", "person", or a comma-separated combination.
func (c *Client) SearchText(ctx context.Context, query, searchType string, page, limit int) ([]SearchResult, *PaginationHeaders, error) {
	params := pageParams(page, limit)
	params.Set("query", query)
	var out []SearchResult
	pg, err := c.get(ctx, "/search/"+url.PathEscape(searchType), params, &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// SearchByID searches by an external ID (imdb, tmdb, tvdb, trakt).
// idType should be one of: "imdb", "tmdb", "tvdb", "trakt".
// searchType filters result types (e.g. "movie", "show" or "" for all).
func (c *Client) SearchByID(ctx context.Context, idType, id, searchType string) ([]SearchResult, error) {
	params := url.Values{}
	if searchType != "" {
		params.Set("type", searchType)
	}
	var out []SearchResult
	_, err := c.get(ctx, "/search/"+url.PathEscape(idType)+"/"+url.PathEscape(id), params, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Calendars.

// CalendarMovies returns movies with releases in the given date range.
// startDate format: YYYY-MM-DD, days is the number of days (1-33).
func (c *Client) CalendarMovies(ctx context.Context, startDate string, days int) ([]CalendarMovie, error) {
	path := "/calendars/all/movies/" + url.PathEscape(startDate) + "/" + strconv.Itoa(days)
	var out []CalendarMovie
	_, err := c.get(ctx, path, url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarShows returns show episodes airing in the given date range.
func (c *Client) CalendarShows(ctx context.Context, startDate string, days int) ([]CalendarShow, error) {
	path := "/calendars/all/shows/" + url.PathEscape(startDate) + "/" + strconv.Itoa(days)
	var out []CalendarShow
	_, err := c.get(ctx, path, url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarNewShows returns new show premieres in the given date range.
func (c *Client) CalendarNewShows(ctx context.Context, startDate string, days int) ([]CalendarShow, error) {
	path := "/calendars/all/shows/new/" + url.PathEscape(startDate) + "/" + strconv.Itoa(days)
	var out []CalendarShow
	_, err := c.get(ctx, path, url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarSeasonPremieres returns season premieres in the given date range.
func (c *Client) CalendarSeasonPremieres(ctx context.Context, startDate string, days int) ([]CalendarShow, error) {
	path := "/calendars/all/shows/premieres/" + url.PathEscape(startDate) + "/" + strconv.Itoa(days)
	var out []CalendarShow
	_, err := c.get(ctx, path, url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Reference data.

// Genres returns all genres for the given type ("movies" or "shows").
func (c *Client) Genres(ctx context.Context, mediaType string) ([]Genre, error) {
	var out []Genre
	_, err := c.get(ctx, "/genres/"+url.PathEscape(mediaType), nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Certifications returns all certifications for the given type ("movies" or "shows").
func (c *Client) Certifications(ctx context.Context, mediaType string) ([]Certification, error) {
	var out []Certification
	_, err := c.get(ctx, "/certifications/"+url.PathEscape(mediaType), nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Countries returns all countries.
func (c *Client) Countries(ctx context.Context, mediaType string) ([]Country, error) {
	var out []Country
	_, err := c.get(ctx, "/countries/"+url.PathEscape(mediaType), nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Languages returns all languages.
func (c *Client) Languages(ctx context.Context, mediaType string) ([]Language, error) {
	var out []Language
	_, err := c.get(ctx, "/languages/"+url.PathEscape(mediaType), nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Networks returns all TV networks.
func (c *Client) Networks(ctx context.Context) ([]Network, error) {
	var out []Network
	_, err := c.get(ctx, "/networks", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Recently updated.

// GetUpdatedMovies returns movies updated since the given date (YYYY-MM-DD).
func (c *Client) GetUpdatedMovies(ctx context.Context, startDate string, page, limit int) ([]UpdatedMovie, *PaginationHeaders, error) {
	var out []UpdatedMovie
	pg, err := c.get(ctx, "/movies/updates/"+url.PathEscape(startDate), extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// GetUpdatedShows returns shows updated since the given date (YYYY-MM-DD).
func (c *Client) GetUpdatedShows(ctx context.Context, startDate string, page, limit int) ([]UpdatedShow, *PaginationHeaders, error) {
	var out []UpdatedShow
	pg, err := c.get(ctx, "/shows/updates/"+url.PathEscape(startDate), extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// User profile (authenticated).

// GetProfile returns the authenticated user's profile.
func (c *Client) GetProfile(ctx context.Context) (*UserProfile, error) {
	var out UserProfile
	_, err := c.get(ctx, "/users/me", url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUserStats returns the authenticated user's stats.
func (c *Client) GetUserStats(ctx context.Context) (*UserStats, error) {
	var out UserStats
	_, err := c.get(ctx, "/users/me/stats", nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Watchlist (authenticated).

// GetWatchlist returns the user's watchlist items filtered by type.
// Pass an empty mediaType for all items.
func (c *Client) GetWatchlist(ctx context.Context, mediaType string, page, limit int) ([]WatchlistItem, *PaginationHeaders, error) {
	path := "/sync/watchlist"
	if mediaType != "" {
		path += "/" + url.PathEscape(mediaType)
	}
	var out []WatchlistItem
	pg, err := c.get(ctx, path, extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// AddToWatchlist adds items to the user's watchlist.
func (c *Client) AddToWatchlist(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/watchlist", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveFromWatchlist removes items from the user's watchlist.
func (c *Client) RemoveFromWatchlist(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/watchlist/remove", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Collection (authenticated).

// GetCollection returns the user's collection filtered by type ("movies" or "shows").
func (c *Client) GetCollection(ctx context.Context, mediaType string, page, limit int) ([]CollectionItem, *PaginationHeaders, error) {
	var out []CollectionItem
	pg, err := c.get(ctx, "/sync/collection/"+url.PathEscape(mediaType), extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// AddToCollection adds items to the user's collection.
func (c *Client) AddToCollection(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/collection", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveFromCollection removes items from the user's collection.
func (c *Client) RemoveFromCollection(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/collection/remove", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// History (authenticated).

// GetHistory returns the user's watch history filtered by type.
func (c *Client) GetHistory(ctx context.Context, mediaType string, page, limit int) ([]HistoryItem, *PaginationHeaders, error) {
	path := "/sync/history"
	if mediaType != "" {
		path += "/" + url.PathEscape(mediaType)
	}
	var out []HistoryItem
	pg, err := c.get(ctx, path, extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// AddToHistory adds items to the user's watch history.
func (c *Client) AddToHistory(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/history", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveFromHistory removes items from the user's watch history.
func (c *Client) RemoveFromHistory(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/history/remove", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Ratings (authenticated).

// GetRatings returns the user's ratings filtered by type.
func (c *Client) GetRatings(ctx context.Context, mediaType string) ([]RatedItem, error) {
	path := "/sync/ratings"
	if mediaType != "" {
		path += "/" + url.PathEscape(mediaType)
	}
	var out []RatedItem
	_, err := c.get(ctx, path, url.Values{"extended": {"full"}}, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AddRatings adds ratings for items.
func (c *Client) AddRatings(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/ratings", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveRatings removes ratings for items.
func (c *Client) RemoveRatings(ctx context.Context, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/sync/ratings/remove", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// User lists (authenticated).

// GetUserLists returns all custom lists for the authenticated user.
func (c *Client) GetUserLists(ctx context.Context) ([]UserList, error) {
	var out []UserList
	_, err := c.get(ctx, "/users/me/lists", nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CreateList creates a new custom list.
func (c *Client) CreateList(ctx context.Context, list *UserList) (*UserList, error) {
	var out UserList
	if err := c.post(ctx, "/users/me/lists", list, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateList updates an existing custom list.
func (c *Client) UpdateList(ctx context.Context, idOrSlug string, list *UserList) error {
	return c.put(ctx, "/users/me/lists/"+url.PathEscape(idOrSlug), list)
}

// DeleteList deletes a custom list.
func (c *Client) DeleteList(ctx context.Context, idOrSlug string) error {
	return c.del(ctx, "/users/me/lists/"+url.PathEscape(idOrSlug))
}

// GetListItems returns all items in a custom list.
func (c *Client) GetListItems(ctx context.Context, idOrSlug string, page, limit int) ([]ListItem, *PaginationHeaders, error) {
	var out []ListItem
	pg, err := c.get(ctx, "/users/me/lists/"+url.PathEscape(idOrSlug)+"/items", extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// AddListItems adds items to a custom list.
func (c *Client) AddListItems(ctx context.Context, idOrSlug string, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/users/me/lists/"+url.PathEscape(idOrSlug)+"/items", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveListItems removes items from a custom list.
func (c *Client) RemoveListItems(ctx context.Context, idOrSlug string, items *SyncItems) (*SyncResponse, error) {
	var out SyncResponse
	if err := c.post(ctx, "/users/me/lists/"+url.PathEscape(idOrSlug)+"/items/remove", items, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Scrobble (authenticated).

// ScrobbleStart starts watching an item.
func (c *Client) ScrobbleStart(ctx context.Context, req *ScrobbleRequest) (*ScrobbleResponse, error) {
	var out ScrobbleResponse
	if err := c.post(ctx, "/scrobble/start", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ScrobblePause pauses watching an item.
func (c *Client) ScrobblePause(ctx context.Context, req *ScrobbleRequest) (*ScrobbleResponse, error) {
	var out ScrobbleResponse
	if err := c.post(ctx, "/scrobble/pause", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ScrobbleStop stops watching an item.
func (c *Client) ScrobbleStop(ctx context.Context, req *ScrobbleRequest) (*ScrobbleResponse, error) {
	var out ScrobbleResponse
	if err := c.post(ctx, "/scrobble/stop", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Checkin (authenticated).

// Checkin checks in to a movie or episode.
func (c *Client) Checkin(ctx context.Context, req *CheckinRequest) (*CheckinResponse, error) {
	var out CheckinResponse
	if err := c.post(ctx, "/checkin", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelCheckin cancels the active checkin.
func (c *Client) CancelCheckin(ctx context.Context) error {
	return c.del(ctx, "/checkin")
}

// Recommendations (authenticated).

// GetMovieRecommendations returns personalized movie recommendations for the user.
func (c *Client) GetMovieRecommendations(ctx context.Context, page, limit int) ([]Movie, *PaginationHeaders, error) {
	var out []Movie
	pg, err := c.get(ctx, "/recommendations/movies", extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// GetShowRecommendations returns personalized show recommendations for the user.
func (c *Client) GetShowRecommendations(ctx context.Context, page, limit int) ([]Show, *PaginationHeaders, error) {
	var out []Show
	pg, err := c.get(ctx, "/recommendations/shows", extendedParams(page, limit), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, pg, nil
}

// OAuth2.

func (c *Client) post(ctx context.Context, path string, body, dst any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("trakt: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.rawBaseURL+path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("trakt: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("trakt-api-key", c.clientID)
	req.Header.Set("trakt-api-version", apiVersion)
	req.Header.Set("User-Agent", c.userAgent)

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("trakt: POST %s: %w", path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("trakt: read response: %w", err)
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
			return fmt.Errorf("trakt: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) del(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.rawBaseURL+path, http.NoBody)
	if err != nil {
		return fmt.Errorf("trakt: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("trakt-api-key", c.clientID)
	req.Header.Set("trakt-api-version", apiVersion)
	req.Header.Set("User-Agent", c.userAgent)

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("trakt: DELETE %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("trakt: read response: %w", err)
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

func (c *Client) put(ctx context.Context, path string, body any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("trakt: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.rawBaseURL+path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("trakt: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("trakt-api-key", c.clientID)
	req.Header.Set("trakt-api-version", apiVersion)
	req.Header.Set("User-Agent", c.userAgent)

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("trakt: PUT %s: %w", path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("trakt: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if jsonErr := json.Unmarshal(respBody, apiErr); jsonErr != nil {
			apiErr.RawBody = string(respBody)
		}
		return apiErr
	}
	return nil
}

// GetDeviceCode starts the OAuth2 device code flow.
// Display the returned UserCode and VerificationURL to the user.
func (c *Client) GetDeviceCode(ctx context.Context) (*DeviceCode, error) {
	var out DeviceCode
	err := c.post(ctx, "/oauth/device/code", map[string]string{
		"client_id": c.clientID,
	}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// PollDeviceToken polls for the device token after the user has authorized the app.
// It blocks until the token is obtained, the code expires, or the context is canceled.
// The interval between polls is taken from the DeviceCode response.
func (c *Client) PollDeviceToken(ctx context.Context, code *DeviceCode) (*Token, error) {
	ticker := time.NewTicker(time.Duration(code.Interval) * time.Second)
	defer ticker.Stop()

	deadline := time.Now().Add(time.Duration(code.ExpiresIn) * time.Second)
	body := map[string]string{
		"code":          code.DeviceCode,
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			if time.Now().After(deadline) {
				return nil, errors.New("trakt: device code expired")
			}
			var token Token
			err := c.post(ctx, "/oauth/device/token", body, &token)
			if err != nil {
				var apiErr *APIError
				if ok := errorAs(err, &apiErr); ok {
					switch apiErr.StatusCode {
					case http.StatusBadRequest: // 400 = pending
						continue
					case http.StatusNotFound: // 404 = invalid code
						return nil, errors.New("trakt: invalid device code")
					case http.StatusConflict: // 409 = already approved
						continue
					case http.StatusGone: // 410 = expired
						return nil, errors.New("trakt: device code expired")
					case http.StatusTeapot: // 418 = denied
						return nil, errors.New("trakt: user denied authorization")
					case http.StatusTooManyRequests: // 429 = slow down
						time.Sleep(time.Duration(code.Interval) * time.Second)
						continue
					}
				}
				return nil, err
			}
			c.storeToken(&token)
			return &token, nil
		}
	}
}

// ExchangeCode exchanges an authorization code for access and refresh tokens.
func (c *Client) ExchangeCode(ctx context.Context, code, redirectURI string) (*Token, error) {
	var token Token
	err := c.post(ctx, "/oauth/token", map[string]string{
		"code":          code,
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"redirect_uri":  redirectURI,
		"grant_type":    "authorization_code",
	}, &token)
	if err != nil {
		return nil, err
	}
	c.storeToken(&token)
	return &token, nil
}

// RefreshToken uses the refresh token to obtain a new access token.
func (c *Client) RefreshToken(ctx context.Context, redirectURI string) (*Token, error) {
	c.mu.RLock()
	rt := c.refreshToken
	c.mu.RUnlock()
	if rt == "" {
		return nil, errors.New("trakt: no refresh token available")
	}

	var token Token
	err := c.post(ctx, "/oauth/token", map[string]string{
		"refresh_token": rt,
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"redirect_uri":  redirectURI,
		"grant_type":    "refresh_token",
	}, &token)
	if err != nil {
		return nil, err
	}
	c.storeToken(&token)
	return &token, nil
}

// RevokeToken revokes the current access token.
func (c *Client) RevokeToken(ctx context.Context) error {
	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()

	return c.post(ctx, "/oauth/revoke", map[string]string{
		"token":         token,
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
	}, nil)
}

func (c *Client) storeToken(t *Token) {
	c.mu.Lock()
	c.accessToken = t.AccessToken
	c.refreshToken = t.RefreshToken
	c.mu.Unlock()
	if c.onToken != nil {
		c.onToken(*t)
	}
}

func errorAs(err error, target **APIError) bool {
	return errors.As(err, target)
}
