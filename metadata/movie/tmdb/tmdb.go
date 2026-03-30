package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL   = "https://api.themoviedb.org/3"
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

// WithBaseURL overrides the default TMDb API base URL.
func WithBaseURL(u string) Option {
	return func(cl *Client) { cl.rawBaseURL = u }
}

// Client is a TMDb v3 API client.
type Client struct {
	accessToken string
	rawBaseURL  string
	httpClient  *http.Client
	userAgent   string
}

// New creates a TMDb [Client] using the given API Read Access Token (Bearer token).
func New(accessToken string, opts ...Option) *Client {
	c := &Client{
		accessToken: accessToken,
		rawBaseURL:  defaultBaseURL,
		httpClient:  &http.Client{Timeout: defaultTimeout},
		userAgent:   defaultUserAgent,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *Client) get(ctx context.Context, path string, dst any) error {
	u, err := url.Parse(c.rawBaseURL + path)
	if err != nil {
		return fmt.Errorf("tmdb: parse URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("tmdb: create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("tmdb: GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("tmdb: read response: %w", err)
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
			return fmt.Errorf("tmdb: decode response: %w", err)
		}
	}
	return nil
}

// GetConfiguration returns the API image configuration.
func (c *Client) GetConfiguration(ctx context.Context) (*Configuration, error) {
	var out Configuration
	if err := c.get(ctx, "/configuration", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SearchMovies searches for movies by title.
func (c *Client) SearchMovies(ctx context.Context, query, language string, page int) (*PaginatedResult[MovieResult], error) {
	var out PaginatedResult[MovieResult]
	path := fmt.Sprintf("/search/movie?query=%s&language=%s&page=%d", url.QueryEscape(query), url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SearchTV searches for TV shows by name.
func (c *Client) SearchTV(ctx context.Context, query, language string, page int) (*PaginatedResult[TVResult], error) {
	var out PaginatedResult[TVResult]
	path := fmt.Sprintf("/search/tv?query=%s&language=%s&page=%d", url.QueryEscape(query), url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SearchMulti searches for movies, TV shows, and people.
func (c *Client) SearchMulti(ctx context.Context, query, language string, page int) (*PaginatedResult[MultiResult], error) {
	var out PaginatedResult[MultiResult]
	path := fmt.Sprintf("/search/multi?query=%s&language=%s&page=%d", url.QueryEscape(query), url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SearchPeople searches for people by name.
func (c *Client) SearchPeople(ctx context.Context, query, language string, page int) (*PaginatedResult[PersonResult], error) {
	var out PaginatedResult[PersonResult]
	path := fmt.Sprintf("/search/person?query=%s&language=%s&page=%d", url.QueryEscape(query), url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovie returns full details for a movie by its TMDb ID.
func (c *Client) GetMovie(ctx context.Context, id int, language string) (*MovieDetails, error) {
	var out MovieDetails
	path := fmt.Sprintf("/movie/%d?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieCredits returns cast and crew for a movie.
func (c *Client) GetMovieCredits(ctx context.Context, id int, language string) (*Credits, error) {
	var out Credits
	path := fmt.Sprintf("/movie/%d/credits?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieImages returns images for a movie.
func (c *Client) GetMovieImages(ctx context.Context, id int) (*Images, error) {
	var out Images
	path := fmt.Sprintf("/movie/%d/images", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieExternalIDs returns external identifiers for a movie.
func (c *Client) GetMovieExternalIDs(ctx context.Context, id int) (*ExternalIDs, error) {
	var out ExternalIDs
	path := fmt.Sprintf("/movie/%d/external_ids", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieRecommendations returns recommended movies.
func (c *Client) GetMovieRecommendations(ctx context.Context, id int, language string, page int) (*PaginatedResult[MovieResult], error) {
	var out PaginatedResult[MovieResult]
	path := fmt.Sprintf("/movie/%d/recommendations?language=%s&page=%d", id, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieSimilar returns similar movies.
func (c *Client) GetMovieSimilar(ctx context.Context, id int, language string, page int) (*PaginatedResult[MovieResult], error) {
	var out PaginatedResult[MovieResult]
	path := fmt.Sprintf("/movie/%d/similar?language=%s&page=%d", id, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTV returns full details for a TV show by its TMDb ID.
func (c *Client) GetTV(ctx context.Context, id int, language string) (*TVDetails, error) {
	var out TVDetails
	path := fmt.Sprintf("/tv/%d?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVSeason returns details for a specific TV season.
func (c *Client) GetTVSeason(ctx context.Context, tvID, seasonNumber int, language string) (*SeasonDetails, error) {
	var out SeasonDetails
	path := fmt.Sprintf("/tv/%d/season/%d?language=%s", tvID, seasonNumber, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVCredits returns cast and crew for a TV show.
func (c *Client) GetTVCredits(ctx context.Context, id int, language string) (*Credits, error) {
	var out Credits
	path := fmt.Sprintf("/tv/%d/credits?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVImages returns images for a TV show.
func (c *Client) GetTVImages(ctx context.Context, id int) (*Images, error) {
	var out Images
	path := fmt.Sprintf("/tv/%d/images", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVExternalIDs returns external identifiers for a TV show.
func (c *Client) GetTVExternalIDs(ctx context.Context, id int) (*ExternalIDs, error) {
	var out ExternalIDs
	path := fmt.Sprintf("/tv/%d/external_ids", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVRecommendations returns recommended TV shows.
func (c *Client) GetTVRecommendations(ctx context.Context, id int, language string, page int) (*PaginatedResult[TVResult], error) {
	var out PaginatedResult[TVResult]
	path := fmt.Sprintf("/tv/%d/recommendations?language=%s&page=%d", id, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVSimilar returns similar TV shows.
func (c *Client) GetTVSimilar(ctx context.Context, id int, language string, page int) (*PaginatedResult[TVResult], error) {
	var out PaginatedResult[TVResult]
	path := fmt.Sprintf("/tv/%d/similar?language=%s&page=%d", id, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPerson returns full details for a person.
func (c *Client) GetPerson(ctx context.Context, id int, language string) (*PersonDetails, error) {
	var out PersonDetails
	path := fmt.Sprintf("/person/%d?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPersonExternalIDs returns external identifiers for a person.
func (c *Client) GetPersonExternalIDs(ctx context.Context, id int) (*ExternalIDs, error) {
	var out ExternalIDs
	path := fmt.Sprintf("/person/%d/external_ids", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DiscoverMovies returns a paginated list of movies matching discover filters.
// Filter parameters are appended as query string (e.g. "&sort_by=popularity.desc&year=2024").
func (c *Client) DiscoverMovies(ctx context.Context, language string, page int, extraParams string) (*PaginatedResult[MovieResult], error) {
	var out PaginatedResult[MovieResult]
	path := fmt.Sprintf("/discover/movie?language=%s&page=%d%s", url.QueryEscape(language), page, extraParams)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DiscoverTV returns a paginated list of TV shows matching discover filters.
func (c *Client) DiscoverTV(ctx context.Context, language string, page int, extraParams string) (*PaginatedResult[TVResult], error) {
	var out PaginatedResult[TVResult]
	path := fmt.Sprintf("/discover/tv?language=%s&page=%d%s", url.QueryEscape(language), page, extraParams)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTrending returns trending items for a media type and time window.
// mediaType can be "all", "movie", "tv", or "person".
// timeWindow can be "day" or "week".
func (c *Client) GetTrending(ctx context.Context, mediaType, timeWindow, language string, page int) (*PaginatedResult[MultiResult], error) {
	var out PaginatedResult[MultiResult]
	path := fmt.Sprintf("/trending/%s/%s?language=%s&page=%d", mediaType, timeWindow, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGenresMovie returns the list of movie genres.
func (c *Client) GetGenresMovie(ctx context.Context, language string) ([]Genre, error) {
	var out struct {
		Genres []Genre `json:"genres"`
	}
	path := "/genre/movie/list?language=" + url.QueryEscape(language)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out.Genres, nil
}

// GetGenresTV returns the list of TV genres.
func (c *Client) GetGenresTV(ctx context.Context, language string) ([]Genre, error) {
	var out struct {
		Genres []Genre `json:"genres"`
	}
	path := "/genre/tv/list?language=" + url.QueryEscape(language)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out.Genres, nil
}

// FindByExternalID finds movies, TV shows, and people by an external identifier.
// externalSource can be "imdb_id", "tvdb_id", etc.
func (c *Client) FindByExternalID(ctx context.Context, externalID, externalSource, language string) (*FindResult, error) {
	var out FindResult
	path := fmt.Sprintf("/find/%s?external_source=%s&language=%s", url.PathEscape(externalID), url.QueryEscape(externalSource), url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPopularMovies returns the current popular movies.
func (c *Client) GetPopularMovies(ctx context.Context, language string, page int) (*PaginatedResult[MovieResult], error) {
	var out PaginatedResult[MovieResult]
	path := fmt.Sprintf("/movie/popular?language=%s&page=%d", url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTopRatedMovies returns the top rated movies.
func (c *Client) GetTopRatedMovies(ctx context.Context, language string, page int) (*PaginatedResult[MovieResult], error) {
	var out PaginatedResult[MovieResult]
	path := fmt.Sprintf("/movie/top_rated?language=%s&page=%d", url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNowPlayingMovies returns movies currently in theaters.
func (c *Client) GetNowPlayingMovies(ctx context.Context, language string, page int) (*PaginatedResult[MovieResult], error) {
	var out PaginatedResult[MovieResult]
	path := fmt.Sprintf("/movie/now_playing?language=%s&page=%d", url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUpcomingMovies returns upcoming movies.
func (c *Client) GetUpcomingMovies(ctx context.Context, language string, page int) (*PaginatedResult[MovieResult], error) {
	var out PaginatedResult[MovieResult]
	path := fmt.Sprintf("/movie/upcoming?language=%s&page=%d", url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPopularTV returns the current popular TV shows.
func (c *Client) GetPopularTV(ctx context.Context, language string, page int) (*PaginatedResult[TVResult], error) {
	var out PaginatedResult[TVResult]
	path := fmt.Sprintf("/tv/popular?language=%s&page=%d", url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTopRatedTV returns the top rated TV shows.
func (c *Client) GetTopRatedTV(ctx context.Context, language string, page int) (*PaginatedResult[TVResult], error) {
	var out PaginatedResult[TVResult]
	path := fmt.Sprintf("/tv/top_rated?language=%s&page=%d", url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
