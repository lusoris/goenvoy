package tmdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://api.themoviedb.org/3"

// Client is a TMDb v3 API client.
type Client struct {
	*metadata.BaseClient
}

// RequestOption configures a single API request.
type RequestOption func(*requestConfig)

type requestConfig struct {
	appendToResponse string
}

// WithAppendToResponse bundles additional sub-resources into the response.
// fields is a comma-separated list, e.g. "credits,images,videos".
func WithAppendToResponse(fields string) RequestOption {
	return func(cfg *requestConfig) { cfg.appendToResponse = fields }
}

// New creates a TMDb [Client] using the given API Read Access Token (Bearer token).
func New(accessToken string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "tmdb", opts...)
	bc.SetAuth(func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	})
	return &Client{BaseClient: bc}
}

func (c *Client) get(ctx context.Context, path string, dst any) error {
	body, status, err := c.DoRaw(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return err
	}

	if status < 200 || status >= 300 {
		apiErr := &APIError{StatusCode: status}
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

func (c *Client) doJSON(ctx context.Context, method, path string, payload, dst any) error {
	var bodyReader io.Reader = http.NoBody
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("tmdb: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	body, status, err := c.DoRaw(ctx, method, path, bodyReader)
	if err != nil {
		return err
	}

	if status < 200 || status >= 300 {
		apiErr := &APIError{StatusCode: status}
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

func (c *Client) post(ctx context.Context, path string, payload, dst any) error {
	return c.doJSON(ctx, http.MethodPost, path, payload, dst)
}

func (c *Client) del(ctx context.Context, path string) error {
	return c.doJSON(ctx, http.MethodDelete, path, nil, nil)
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

func applyRequestOpts(path string, opts []RequestOption) string {
	var cfg requestConfig
	for _, o := range opts {
		o(&cfg)
	}
	if cfg.appendToResponse != "" {
		sep := "&"
		if !strings.Contains(path, "?") {
			sep = "?"
		}
		path += sep + "append_to_response=" + url.QueryEscape(cfg.appendToResponse)
	}
	return path
}

// GetMovieFull returns movie details with optional appended sub-resources.
func (c *Client) GetMovieFull(ctx context.Context, id int, language string, opts ...RequestOption) (*MovieDetailsFull, error) {
	var out MovieDetailsFull
	path := applyRequestOpts(fmt.Sprintf("/movie/%d?language=%s", id, url.QueryEscape(language)), opts)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVFull returns TV show details with optional appended sub-resources.
func (c *Client) GetTVFull(ctx context.Context, id int, language string, opts ...RequestOption) (*TVDetailsFull, error) {
	var out TVDetailsFull
	path := applyRequestOpts(fmt.Sprintf("/tv/%d?language=%s", id, url.QueryEscape(language)), opts)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVSeasonFull returns season details with optional appended sub-resources.
func (c *Client) GetTVSeasonFull(ctx context.Context, tvID, seasonNumber int, language string, opts ...RequestOption) (*SeasonDetailsFull, error) {
	var out SeasonDetailsFull
	path := applyRequestOpts(fmt.Sprintf("/tv/%d/season/%d?language=%s", tvID, seasonNumber, url.QueryEscape(language)), opts)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVEpisodeFull returns episode details with optional appended sub-resources.
func (c *Client) GetTVEpisodeFull(ctx context.Context, tvID, seasonNumber, episodeNumber int, language string, opts ...RequestOption) (*EpisodeDetailsFull, error) {
	var out EpisodeDetailsFull
	path := applyRequestOpts(fmt.Sprintf("/tv/%d/season/%d/episode/%d?language=%s", tvID, seasonNumber, episodeNumber, url.QueryEscape(language)), opts)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPersonFull returns person details with optional appended sub-resources.
func (c *Client) GetPersonFull(ctx context.Context, id int, language string, opts ...RequestOption) (*PersonDetailsFull, error) {
	var out PersonDetailsFull
	path := applyRequestOpts(fmt.Sprintf("/person/%d?language=%s", id, url.QueryEscape(language)), opts)
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
// Extra filter parameters are safely encoded and appended to the query string.
func (c *Client) DiscoverMovies(ctx context.Context, language string, page int, extraParams url.Values) (*PaginatedResult[MovieResult], error) {
	var out PaginatedResult[MovieResult]
	params := url.Values{}
	params.Set("language", language)
	params.Set("page", strconv.Itoa(page))
	for k, vs := range extraParams {
		for _, v := range vs {
			params.Add(k, v)
		}
	}
	path := "/discover/movie?" + params.Encode()
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DiscoverTV returns a paginated list of TV shows matching discover filters.
// Extra filter parameters are safely encoded and appended to the query string.
func (c *Client) DiscoverTV(ctx context.Context, language string, page int, extraParams url.Values) (*PaginatedResult[TVResult], error) {
	var out PaginatedResult[TVResult]
	params := url.Values{}
	params.Set("language", language)
	params.Set("page", strconv.Itoa(page))
	for k, vs := range extraParams {
		for _, v := range vs {
			params.Add(k, v)
		}
	}
	path := "/discover/tv?" + params.Encode()
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
	path := fmt.Sprintf("/trending/%s/%s?language=%s&page=%d", url.PathEscape(mediaType), url.PathEscape(timeWindow), url.QueryEscape(language), page)
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

// Movie extras.

// GetMovieVideos returns videos (trailers, teasers, etc.) for a movie.
func (c *Client) GetMovieVideos(ctx context.Context, id int, language string) (*VideosResponse, error) {
	var out VideosResponse
	path := fmt.Sprintf("/movie/%d/videos?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieKeywords returns keywords for a movie.
func (c *Client) GetMovieKeywords(ctx context.Context, id int) (*KeywordsResponse, error) {
	var out KeywordsResponse
	path := fmt.Sprintf("/movie/%d/keywords", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieReviews returns reviews for a movie.
func (c *Client) GetMovieReviews(ctx context.Context, id int, language string, page int) (*PaginatedResult[Review], error) {
	var out PaginatedResult[Review]
	path := fmt.Sprintf("/movie/%d/reviews?language=%s&page=%d", id, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieReleaseDates returns release date info for a movie.
func (c *Client) GetMovieReleaseDates(ctx context.Context, id int) (*ReleaseDatesResponse, error) {
	var out ReleaseDatesResponse
	path := fmt.Sprintf("/movie/%d/release_dates", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieWatchProviders returns watch providers for a movie.
func (c *Client) GetMovieWatchProviders(ctx context.Context, id int) (*WatchProvidersResponse, error) {
	var out WatchProvidersResponse
	path := fmt.Sprintf("/movie/%d/watch/providers", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieAlternativeTitles returns alternative titles for a movie.
func (c *Client) GetMovieAlternativeTitles(ctx context.Context, id int, country string) (*AlternativeTitlesResponse, error) {
	var out AlternativeTitlesResponse
	path := fmt.Sprintf("/movie/%d/alternative_titles", id)
	if country != "" {
		path += "?country=" + url.QueryEscape(country)
	}
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieTranslations returns translations for a movie.
func (c *Client) GetMovieTranslations(ctx context.Context, id int) (*TranslationsResponse, error) {
	var out TranslationsResponse
	path := fmt.Sprintf("/movie/%d/translations", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieLists returns lists that contain a movie.
func (c *Client) GetMovieLists(ctx context.Context, id int, language string, page int) (*PaginatedResult[ListSummary], error) {
	var out PaginatedResult[ListSummary]
	path := fmt.Sprintf("/movie/%d/lists?language=%s&page=%d", id, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieAccountStates returns the account states (rated, watchlist, favorite) for a movie.
func (c *Client) GetMovieAccountStates(ctx context.Context, id int) (*AccountStates, error) {
	var out AccountStates
	path := fmt.Sprintf("/movie/%d/account_states", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RateMovie rates a movie. Rating is a value from 0.5 to 10.0 in 0.5 increments.
func (c *Client) RateMovie(ctx context.Context, id int, rating float64) error {
	path := fmt.Sprintf("/movie/%d/rating", id)
	return c.post(ctx, path, map[string]float64{"value": rating}, nil)
}

// DeleteMovieRating removes the rating for a movie.
func (c *Client) DeleteMovieRating(ctx context.Context, id int) error {
	return c.del(ctx, fmt.Sprintf("/movie/%d/rating", id))
}

// TV Show extras.

// GetTVVideos returns videos for a TV show.
func (c *Client) GetTVVideos(ctx context.Context, id int, language string) (*VideosResponse, error) {
	var out VideosResponse
	path := fmt.Sprintf("/tv/%d/videos?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVKeywords returns keywords for a TV show.
func (c *Client) GetTVKeywords(ctx context.Context, id int) (*KeywordsResponse, error) {
	var out KeywordsResponse
	path := fmt.Sprintf("/tv/%d/keywords", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVReviews returns reviews for a TV show.
func (c *Client) GetTVReviews(ctx context.Context, id int, language string, page int) (*PaginatedResult[Review], error) {
	var out PaginatedResult[Review]
	path := fmt.Sprintf("/tv/%d/reviews?language=%s&page=%d", id, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVWatchProviders returns watch providers for a TV show.
func (c *Client) GetTVWatchProviders(ctx context.Context, id int) (*WatchProvidersResponse, error) {
	var out WatchProvidersResponse
	path := fmt.Sprintf("/tv/%d/watch/providers", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVAlternativeTitles returns alternative titles for a TV show.
func (c *Client) GetTVAlternativeTitles(ctx context.Context, id int) (*AlternativeTitlesResponse, error) {
	var out AlternativeTitlesResponse
	path := fmt.Sprintf("/tv/%d/alternative_titles", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVTranslations returns translations for a TV show.
func (c *Client) GetTVTranslations(ctx context.Context, id int) (*TranslationsResponse, error) {
	var out TranslationsResponse
	path := fmt.Sprintf("/tv/%d/translations", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVContentRatings returns content ratings for a TV show.
func (c *Client) GetTVContentRatings(ctx context.Context, id int) (*ContentRatingsResponse, error) {
	var out ContentRatingsResponse
	path := fmt.Sprintf("/tv/%d/content_ratings", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVEpisodeGroups returns episode groups for a TV show.
func (c *Client) GetTVEpisodeGroups(ctx context.Context, id int) (*EpisodeGroupsResponse, error) {
	var out EpisodeGroupsResponse
	path := fmt.Sprintf("/tv/%d/episode_groups", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVAggregateCredits returns aggregate credits for a TV show.
func (c *Client) GetTVAggregateCredits(ctx context.Context, id int, language string) (*AggregateCredits, error) {
	var out AggregateCredits
	path := fmt.Sprintf("/tv/%d/aggregate_credits?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVAccountStates returns account states for a TV show.
func (c *Client) GetTVAccountStates(ctx context.Context, id int) (*AccountStates, error) {
	var out AccountStates
	path := fmt.Sprintf("/tv/%d/account_states", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RateTV rates a TV show.
func (c *Client) RateTV(ctx context.Context, id int, rating float64) error {
	path := fmt.Sprintf("/tv/%d/rating", id)
	return c.post(ctx, path, map[string]float64{"value": rating}, nil)
}

// DeleteTVRating removes the rating for a TV show.
func (c *Client) DeleteTVRating(ctx context.Context, id int) error {
	return c.del(ctx, fmt.Sprintf("/tv/%d/rating", id))
}

// TV Season extras.

// GetTVSeasonCredits returns credits for a TV season.
func (c *Client) GetTVSeasonCredits(ctx context.Context, tvID, seasonNumber int, language string) (*Credits, error) {
	var out Credits
	path := fmt.Sprintf("/tv/%d/season/%d/credits?language=%s", tvID, seasonNumber, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVSeasonImages returns images for a TV season.
func (c *Client) GetTVSeasonImages(ctx context.Context, tvID, seasonNumber int) (*Images, error) {
	var out Images
	path := fmt.Sprintf("/tv/%d/season/%d/images", tvID, seasonNumber)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVSeasonVideos returns videos for a TV season.
func (c *Client) GetTVSeasonVideos(ctx context.Context, tvID, seasonNumber int, language string) (*VideosResponse, error) {
	var out VideosResponse
	path := fmt.Sprintf("/tv/%d/season/%d/videos?language=%s", tvID, seasonNumber, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVSeasonExternalIDs returns external IDs for a TV season.
func (c *Client) GetTVSeasonExternalIDs(ctx context.Context, tvID, seasonNumber int) (*ExternalIDs, error) {
	var out ExternalIDs
	path := fmt.Sprintf("/tv/%d/season/%d/external_ids", tvID, seasonNumber)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVSeasonTranslations returns translations for a TV season.
func (c *Client) GetTVSeasonTranslations(ctx context.Context, tvID, seasonNumber int) (*TranslationsResponse, error) {
	var out TranslationsResponse
	path := fmt.Sprintf("/tv/%d/season/%d/translations", tvID, seasonNumber)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVSeasonWatchProviders returns watch providers for a TV season.
func (c *Client) GetTVSeasonWatchProviders(ctx context.Context, tvID, seasonNumber int) (*WatchProvidersResponse, error) {
	var out WatchProvidersResponse
	path := fmt.Sprintf("/tv/%d/season/%d/watch/providers", tvID, seasonNumber)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TV Episode.

// GetTVEpisode returns details for a specific TV episode.
func (c *Client) GetTVEpisode(ctx context.Context, tvID, seasonNumber, episodeNumber int, language string) (*EpisodeDetails, error) {
	var out EpisodeDetails
	path := fmt.Sprintf("/tv/%d/season/%d/episode/%d?language=%s", tvID, seasonNumber, episodeNumber, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVEpisodeCredits returns credits for a TV episode.
func (c *Client) GetTVEpisodeCredits(ctx context.Context, tvID, seasonNumber, episodeNumber int, language string) (*Credits, error) {
	var out Credits
	path := fmt.Sprintf("/tv/%d/season/%d/episode/%d/credits?language=%s", tvID, seasonNumber, episodeNumber, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVEpisodeImages returns images for a TV episode.
func (c *Client) GetTVEpisodeImages(ctx context.Context, tvID, seasonNumber, episodeNumber int) (*Images, error) {
	var out Images
	path := fmt.Sprintf("/tv/%d/season/%d/episode/%d/images", tvID, seasonNumber, episodeNumber)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVEpisodeVideos returns videos for a TV episode.
func (c *Client) GetTVEpisodeVideos(ctx context.Context, tvID, seasonNumber, episodeNumber int, language string) (*VideosResponse, error) {
	var out VideosResponse
	path := fmt.Sprintf("/tv/%d/season/%d/episode/%d/videos?language=%s", tvID, seasonNumber, episodeNumber, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVEpisodeExternalIDs returns external IDs for a TV episode.
func (c *Client) GetTVEpisodeExternalIDs(ctx context.Context, tvID, seasonNumber, episodeNumber int) (*ExternalIDs, error) {
	var out ExternalIDs
	path := fmt.Sprintf("/tv/%d/season/%d/episode/%d/external_ids", tvID, seasonNumber, episodeNumber)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVEpisodeTranslations returns translations for a TV episode.
func (c *Client) GetTVEpisodeTranslations(ctx context.Context, tvID, seasonNumber, episodeNumber int) (*TranslationsResponse, error) {
	var out TranslationsResponse
	path := fmt.Sprintf("/tv/%d/season/%d/episode/%d/translations", tvID, seasonNumber, episodeNumber)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVEpisodeAccountStates returns account states for a TV episode.
func (c *Client) GetTVEpisodeAccountStates(ctx context.Context, tvID, seasonNumber, episodeNumber int) (*AccountStates, error) {
	var out AccountStates
	path := fmt.Sprintf("/tv/%d/season/%d/episode/%d/account_states", tvID, seasonNumber, episodeNumber)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RateTVEpisode rates a TV episode.
func (c *Client) RateTVEpisode(ctx context.Context, tvID, seasonNumber, episodeNumber int, rating float64) error {
	path := fmt.Sprintf("/tv/%d/season/%d/episode/%d/rating", tvID, seasonNumber, episodeNumber)
	return c.post(ctx, path, map[string]float64{"value": rating}, nil)
}

// DeleteTVEpisodeRating removes the rating for a TV episode.
func (c *Client) DeleteTVEpisodeRating(ctx context.Context, tvID, seasonNumber, episodeNumber int) error {
	return c.del(ctx, fmt.Sprintf("/tv/%d/season/%d/episode/%d/rating", tvID, seasonNumber, episodeNumber))
}

// Person extras.

// GetPersonMovieCredits returns movie credits for a person.
func (c *Client) GetPersonMovieCredits(ctx context.Context, id int, language string) (*PersonCredits, error) {
	var out PersonCredits
	path := fmt.Sprintf("/person/%d/movie_credits?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPersonTVCredits returns TV credits for a person.
func (c *Client) GetPersonTVCredits(ctx context.Context, id int, language string) (*PersonCredits, error) {
	var out PersonCredits
	path := fmt.Sprintf("/person/%d/tv_credits?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPersonCombinedCredits returns combined movie and TV credits for a person.
func (c *Client) GetPersonCombinedCredits(ctx context.Context, id int, language string) (*PersonCredits, error) {
	var out PersonCredits
	path := fmt.Sprintf("/person/%d/combined_credits?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPersonImages returns images (profile photos) for a person.
func (c *Client) GetPersonImages(ctx context.Context, id int) (*PersonImages, error) {
	var out PersonImages
	path := fmt.Sprintf("/person/%d/images", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPersonTaggedImages returns tagged images for a person.
func (c *Client) GetPersonTaggedImages(ctx context.Context, id int, language string, page int) (*PaginatedResult[TaggedImage], error) {
	var out PaginatedResult[TaggedImage]
	path := fmt.Sprintf("/person/%d/tagged_images?language=%s&page=%d", id, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Search extras.

// SearchKeywords searches for keywords by query.
func (c *Client) SearchKeywords(ctx context.Context, query string, page int) (*PaginatedResult[Keyword], error) {
	var out PaginatedResult[Keyword]
	path := fmt.Sprintf("/search/keyword?query=%s&page=%d", url.QueryEscape(query), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SearchCollections searches for collections by query.
func (c *Client) SearchCollections(ctx context.Context, query, language string, page int) (*PaginatedResult[CollectionResult], error) {
	var out PaginatedResult[CollectionResult]
	path := fmt.Sprintf("/search/collection?query=%s&language=%s&page=%d", url.QueryEscape(query), url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SearchCompanies searches for companies by query.
func (c *Client) SearchCompanies(ctx context.Context, query string, page int) (*PaginatedResult[CompanyResult], error) {
	var out PaginatedResult[CompanyResult]
	path := fmt.Sprintf("/search/company?query=%s&page=%d", url.QueryEscape(query), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Collections.

// GetCollection returns details for a movie collection.
func (c *Client) GetCollection(ctx context.Context, id int, language string) (*CollectionDetails, error) {
	var out CollectionDetails
	path := fmt.Sprintf("/collection/%d?language=%s", id, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCollectionImages returns images for a collection.
func (c *Client) GetCollectionImages(ctx context.Context, id int) (*Images, error) {
	var out Images
	path := fmt.Sprintf("/collection/%d/images", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCollectionTranslations returns translations for a collection.
func (c *Client) GetCollectionTranslations(ctx context.Context, id int) (*TranslationsResponse, error) {
	var out TranslationsResponse
	path := fmt.Sprintf("/collection/%d/translations", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Account.

// GetAccountDetails returns the account details for the authenticated user.
func (c *Client) GetAccountDetails(ctx context.Context) (*AccountDetails, error) {
	var out AccountDetails
	if err := c.get(ctx, "/account", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFavoriteMovies returns the user's favorite movies.
func (c *Client) GetFavoriteMovies(ctx context.Context, accountID int, language string, page int) (*PaginatedResult[MovieResult], error) {
	var out PaginatedResult[MovieResult]
	path := fmt.Sprintf("/account/%d/favorite/movies?language=%s&page=%d", accountID, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetFavoriteTV returns the user's favorite TV shows.
func (c *Client) GetFavoriteTV(ctx context.Context, accountID int, language string, page int) (*PaginatedResult[TVResult], error) {
	var out PaginatedResult[TVResult]
	path := fmt.Sprintf("/account/%d/favorite/tv?language=%s&page=%d", accountID, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddFavorite adds or removes a movie/TV show from favorites.
func (c *Client) AddFavorite(ctx context.Context, accountID int, mediaType string, mediaID int, favorite bool) error {
	path := fmt.Sprintf("/account/%d/favorite", accountID)
	return c.post(ctx, path, map[string]any{
		"media_type": mediaType,
		"media_id":   mediaID,
		"favorite":   favorite,
	}, nil)
}

// GetWatchlistMovies returns the user's movie watchlist.
func (c *Client) GetWatchlistMovies(ctx context.Context, accountID int, language string, page int) (*PaginatedResult[MovieResult], error) {
	var out PaginatedResult[MovieResult]
	path := fmt.Sprintf("/account/%d/watchlist/movies?language=%s&page=%d", accountID, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWatchlistTV returns the user's TV watchlist.
func (c *Client) GetWatchlistTV(ctx context.Context, accountID int, language string, page int) (*PaginatedResult[TVResult], error) {
	var out PaginatedResult[TVResult]
	path := fmt.Sprintf("/account/%d/watchlist/tv?language=%s&page=%d", accountID, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddToWatchlist adds or removes a movie/TV show from the watchlist.
func (c *Client) AddToWatchlist(ctx context.Context, accountID int, mediaType string, mediaID int, watchlist bool) error {
	path := fmt.Sprintf("/account/%d/watchlist", accountID)
	return c.post(ctx, path, map[string]any{
		"media_type": mediaType,
		"media_id":   mediaID,
		"watchlist":  watchlist,
	}, nil)
}

// GetRatedMovies returns the user's rated movies.
func (c *Client) GetRatedMovies(ctx context.Context, accountID int, language string, page int) (*PaginatedResult[RatedMovie], error) {
	var out PaginatedResult[RatedMovie]
	path := fmt.Sprintf("/account/%d/rated/movies?language=%s&page=%d", accountID, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRatedTV returns the user's rated TV shows.
func (c *Client) GetRatedTV(ctx context.Context, accountID int, language string, page int) (*PaginatedResult[RatedTV], error) {
	var out PaginatedResult[RatedTV]
	path := fmt.Sprintf("/account/%d/rated/tv?language=%s&page=%d", accountID, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRatedTVEpisodes returns the user's rated TV episodes.
func (c *Client) GetRatedTVEpisodes(ctx context.Context, accountID int, language string, page int) (*PaginatedResult[RatedEpisode], error) {
	var out PaginatedResult[RatedEpisode]
	path := fmt.Sprintf("/account/%d/rated/tv/episodes?language=%s&page=%d", accountID, url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Lists.

// GetList returns details for a list.
func (c *Client) GetList(ctx context.Context, listID int, language string) (*ListDetails, error) {
	var out ListDetails
	path := fmt.Sprintf("/list/%d?language=%s", listID, url.QueryEscape(language))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateList creates a new list.
func (c *Client) CreateList(ctx context.Context, name, description, language string) (*CreateListResponse, error) {
	var out CreateListResponse
	if err := c.post(ctx, "/list", map[string]string{
		"name":        name,
		"description": description,
		"language":    language,
	}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddToList adds a movie to a list.
func (c *Client) AddToList(ctx context.Context, listID, mediaID int) error {
	path := fmt.Sprintf("/list/%d/add_item", listID)
	return c.post(ctx, path, map[string]int{"media_id": mediaID}, nil)
}

// RemoveFromList removes a movie from a list.
func (c *Client) RemoveFromList(ctx context.Context, listID, mediaID int) error {
	path := fmt.Sprintf("/list/%d/remove_item", listID)
	return c.post(ctx, path, map[string]int{"media_id": mediaID}, nil)
}

// ClearList removes all items from a list.
func (c *Client) ClearList(ctx context.Context, listID int, confirm bool) error {
	path := fmt.Sprintf("/list/%d/clear?confirm=%t", listID, confirm)
	return c.post(ctx, path, nil, nil)
}

// DeleteList deletes a list.
func (c *Client) DeleteList(ctx context.Context, listID int) error {
	return c.del(ctx, fmt.Sprintf("/list/%d", listID))
}

// CheckItemStatus checks if a movie is on a list.
func (c *Client) CheckItemStatus(ctx context.Context, listID, movieID int) (*ItemStatus, error) {
	var out ItemStatus
	path := fmt.Sprintf("/list/%d/item_status?movie_id=%d", listID, movieID)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Certifications.

// GetMovieCertifications returns the list of movie certifications by country.
func (c *Client) GetMovieCertifications(ctx context.Context) (*CertificationsResponse, error) {
	var out CertificationsResponse
	if err := c.get(ctx, "/certification/movie/list", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVCertifications returns the list of TV certifications by country.
func (c *Client) GetTVCertifications(ctx context.Context) (*CertificationsResponse, error) {
	var out CertificationsResponse
	if err := c.get(ctx, "/certification/tv/list", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Watch Providers.

// GetAvailableWatchProviderRegions returns available watch provider regions.
func (c *Client) GetAvailableWatchProviderRegions(ctx context.Context, language string) (*WatchProviderRegionsResponse, error) {
	var out WatchProviderRegionsResponse
	path := "/watch/providers/regions?language=" + url.QueryEscape(language)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieWatchProviderList returns the list of available movie watch providers.
func (c *Client) GetMovieWatchProviderList(ctx context.Context, language, watchRegion string) (*WatchProviderListResponse, error) {
	var out WatchProviderListResponse
	path := fmt.Sprintf("/watch/providers/movie?language=%s&watch_region=%s", url.QueryEscape(language), url.QueryEscape(watchRegion))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVWatchProviderList returns the list of available TV watch providers.
func (c *Client) GetTVWatchProviderList(ctx context.Context, language, watchRegion string) (*WatchProviderListResponse, error) {
	var out WatchProviderListResponse
	path := fmt.Sprintf("/watch/providers/tv?language=%s&watch_region=%s", url.QueryEscape(language), url.QueryEscape(watchRegion))
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Companies & Keywords.

// GetCompany returns details for a production company.
func (c *Client) GetCompany(ctx context.Context, id int) (*CompanyDetails, error) {
	var out CompanyDetails
	path := fmt.Sprintf("/company/%d", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetKeyword returns a keyword by ID.
func (c *Client) GetKeyword(ctx context.Context, id int) (*Keyword, error) {
	var out Keyword
	path := fmt.Sprintf("/keyword/%d", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Changes.

// GetMovieChanges returns recently changed movie IDs.
func (c *Client) GetMovieChanges(ctx context.Context, startDate, endDate string, page int) (*ChangesResponse, error) {
	var out ChangesResponse
	params := url.Values{}
	if startDate != "" {
		params.Set("start_date", startDate)
	}
	if endDate != "" {
		params.Set("end_date", endDate)
	}
	params.Set("page", strconv.Itoa(page))
	path := "/movie/changes?" + params.Encode()
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTVChanges returns recently changed TV show IDs.
func (c *Client) GetTVChanges(ctx context.Context, startDate, endDate string, page int) (*ChangesResponse, error) {
	var out ChangesResponse
	params := url.Values{}
	if startDate != "" {
		params.Set("start_date", startDate)
	}
	if endDate != "" {
		params.Set("end_date", endDate)
	}
	params.Set("page", strconv.Itoa(page))
	path := "/tv/changes?" + params.Encode()
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPersonChanges returns recently changed person IDs.
func (c *Client) GetPersonChanges(ctx context.Context, startDate, endDate string, page int) (*ChangesResponse, error) {
	var out ChangesResponse
	params := url.Values{}
	if startDate != "" {
		params.Set("start_date", startDate)
	}
	if endDate != "" {
		params.Set("end_date", endDate)
	}
	params.Set("page", strconv.Itoa(page))
	path := "/person/changes?" + params.Encode()
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Reviews.

// GetReview returns a single review by ID.
func (c *Client) GetReview(ctx context.Context, reviewID string) (*Review, error) {
	var out Review
	path := "/review/" + url.PathEscape(reviewID)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Networks.

// GetNetwork returns details for a TV network.
func (c *Client) GetNetwork(ctx context.Context, id int) (*Network, error) {
	var out Network
	path := fmt.Sprintf("/network/%d", id)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetOnTheAirTV returns TV shows currently on the air.
func (c *Client) GetOnTheAirTV(ctx context.Context, language string, page int) (*PaginatedResult[TVResult], error) {
	var out PaginatedResult[TVResult]
	path := fmt.Sprintf("/tv/on_the_air?language=%s&page=%d", url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAiringTodayTV returns TV shows airing today.
func (c *Client) GetAiringTodayTV(ctx context.Context, language string, page int) (*PaginatedResult[TVResult], error) {
	var out PaginatedResult[TVResult]
	path := fmt.Sprintf("/tv/airing_today?language=%s&page=%d", url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLatestMovie returns the latest movie added to TMDb.
func (c *Client) GetLatestMovie(ctx context.Context) (*MovieDetails, error) {
	var out MovieDetails
	if err := c.get(ctx, "/movie/latest", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLatestTV returns the latest TV show added to TMDb.
func (c *Client) GetLatestTV(ctx context.Context) (*TVDetails, error) {
	var out TVDetails
	if err := c.get(ctx, "/tv/latest", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPopularPeople returns popular people.
func (c *Client) GetPopularPeople(ctx context.Context, language string, page int) (*PaginatedResult[PersonResult], error) {
	var out PaginatedResult[PersonResult]
	path := fmt.Sprintf("/person/popular?language=%s&page=%d", url.QueryEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLanguages returns the list of languages supported by TMDb.
func (c *Client) GetLanguages(ctx context.Context) ([]Language, error) {
	var out []Language
	if err := c.get(ctx, "/configuration/languages", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetCountries returns the list of countries supported by TMDb.
func (c *Client) GetCountries(ctx context.Context) ([]Country, error) {
	var out []Country
	if err := c.get(ctx, "/configuration/countries", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTimezones returns the list of timezones supported by TMDb.
func (c *Client) GetTimezones(ctx context.Context) ([]Timezone, error) {
	var out []Timezone
	if err := c.get(ctx, "/configuration/timezones", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetJobs returns the list of jobs and departments.
func (c *Client) GetJobs(ctx context.Context) ([]Department, error) {
	var out []Department
	if err := c.get(ctx, "/configuration/jobs", &out); err != nil {
		return nil, err
	}
	return out, nil
}
