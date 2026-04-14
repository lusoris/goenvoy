package bazarr

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/golusoris/goenvoy/arr/v2"
)

// Client is a Bazarr API client.
type Client struct {
	base *arr.BaseClient
}

// New creates a Bazarr [Client] for the instance at baseURL.
func New(baseURL, apiKey string, opts ...arr.Option) (*Client, error) {
	base, err := arr.NewBaseClient(baseURL, apiKey, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{base: base}, nil
}

// dataWrapper is used to unwrap the top-level "data" key from Bazarr responses.
type dataWrapper[T any] struct {
	Data T `json:"data"`
}

// GetSeries returns series metadata. Pass seriesIDs to filter specific series,
// or empty to get a paginated list.
func (c *Client) GetSeries(ctx context.Context, start, length int, seriesIDs ...int) (*PagedResponse[Series], error) {
	var out PagedResponse[Series]
	params := url.Values{}
	params.Set("start", strconv.Itoa(start))
	params.Set("length", strconv.Itoa(length))
	for _, id := range seriesIDs {
		params.Add("seriesid[]", strconv.Itoa(id))
	}
	path := "/api/series?" + params.Encode()
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetSeriesProfile updates the language profile for one or more series.
func (c *Client) SetSeriesProfile(ctx context.Context, seriesIDs []int, profileIDs []string) error {
	body := map[string]any{
		"seriesid":  seriesIDs,
		"profileid": profileIDs,
	}
	return c.base.Post(ctx, "/api/series", body, nil)
}

// RunSeriesAction triggers an action on a specific series.
// Valid actions: "scan-disk", "search-missing", "search-wanted", "sync".
func (c *Client) RunSeriesAction(ctx context.Context, seriesID int, action string) error {
	body := map[string]any{
		"seriesid": seriesID,
		"action":   action,
	}
	return c.base.Patch(ctx, "/api/series", body, nil)
}

// GetEpisodes returns episodes for the given series or episode IDs.
func (c *Client) GetEpisodes(ctx context.Context, seriesIDs, episodeIDs []int) (*PagedResponse[Episode], error) {
	var out PagedResponse[Episode]
	params := url.Values{}
	for _, id := range seriesIDs {
		params.Add("seriesid[]", strconv.Itoa(id))
	}
	for _, id := range episodeIDs {
		params.Add("episodeid[]", strconv.Itoa(id))
	}
	path := "/api/episodes?" + params.Encode()
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetWantedEpisodes returns episodes missing subtitles.
func (c *Client) GetWantedEpisodes(ctx context.Context, start, length int) (*PagedResponse[WantedEpisode], error) {
	var out PagedResponse[WantedEpisode]
	path := fmt.Sprintf("/api/episodes/wanted?start=%d&length=%d", start, length)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetEpisodeHistory returns subtitle history events for episodes.
func (c *Client) GetEpisodeHistory(ctx context.Context, start, length int, episodeID *int) (*PagedResponse[EpisodeHistoryRecord], error) {
	var out PagedResponse[EpisodeHistoryRecord]
	params := url.Values{}
	params.Set("start", strconv.Itoa(start))
	params.Set("length", strconv.Itoa(length))
	if episodeID != nil {
		params.Set("episodeid", strconv.Itoa(*episodeID))
	}
	path := "/api/episodes/history?" + params.Encode()
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadEpisodeSubtitle triggers a subtitle download for a specific episode.
func (c *Client) DownloadEpisodeSubtitle(ctx context.Context, seriesID, episodeID int, language, forced, hi string) error {
	body := map[string]any{
		"seriesid":  seriesID,
		"episodeid": episodeID,
		"language":  language,
		"forced":    forced,
		"hi":        hi,
	}
	return c.base.Patch(ctx, "/api/episodes/subtitles", body, nil)
}

// DeleteEpisodeSubtitle removes a subtitle file for an episode.
func (c *Client) DeleteEpisodeSubtitle(ctx context.Context, seriesID, episodeID int, language, forced, hi, subtitlePath string) error {
	params := url.Values{}
	params.Set("seriesid", strconv.Itoa(seriesID))
	params.Set("episodeid", strconv.Itoa(episodeID))
	params.Set("language", language)
	params.Set("forced", forced)
	params.Set("hi", hi)
	params.Set("path", subtitlePath)
	path := "/api/episodes/subtitles?" + params.Encode()
	return c.base.Delete(ctx, path, nil, nil)
}

// GetMovies returns movie metadata. Pass radarrIDs to filter specific movies,
// or empty to get a paginated list.
func (c *Client) GetMovies(ctx context.Context, start, length int, radarrIDs ...int) (*PagedResponse[Movie], error) {
	var out PagedResponse[Movie]
	params := url.Values{}
	params.Set("start", strconv.Itoa(start))
	params.Set("length", strconv.Itoa(length))
	for _, id := range radarrIDs {
		params.Add("radarrid[]", strconv.Itoa(id))
	}
	path := "/api/movies?" + params.Encode()
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetMovieProfile updates the language profile for one or more movies.
func (c *Client) SetMovieProfile(ctx context.Context, radarrIDs []int, profileIDs []string) error {
	body := map[string]any{
		"radarrid":  radarrIDs,
		"profileid": profileIDs,
	}
	return c.base.Post(ctx, "/api/movies", body, nil)
}

// RunMovieAction triggers an action on a specific movie.
// Valid actions: "scan-disk", "search-missing", "search-wanted", "sync".
func (c *Client) RunMovieAction(ctx context.Context, radarrID int, action string) error {
	body := map[string]any{
		"radarrid": radarrID,
		"action":   action,
	}
	return c.base.Patch(ctx, "/api/movies", body, nil)
}

// GetWantedMovies returns movies missing subtitles.
func (c *Client) GetWantedMovies(ctx context.Context, start, length int) (*PagedResponse[WantedMovie], error) {
	var out PagedResponse[WantedMovie]
	path := fmt.Sprintf("/api/movies/wanted?start=%d&length=%d", start, length)
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMovieHistory returns subtitle history events for movies.
func (c *Client) GetMovieHistory(ctx context.Context, start, length int, radarrID *int) (*PagedResponse[MovieHistoryRecord], error) {
	var out PagedResponse[MovieHistoryRecord]
	params := url.Values{}
	params.Set("start", strconv.Itoa(start))
	params.Set("length", strconv.Itoa(length))
	if radarrID != nil {
		params.Set("radarrid", strconv.Itoa(*radarrID))
	}
	path := "/api/movies/history?" + params.Encode()
	if err := c.base.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadMovieSubtitle triggers a subtitle download for a specific movie.
func (c *Client) DownloadMovieSubtitle(ctx context.Context, radarrID int, language, forced, hi string) error {
	body := map[string]any{
		"radarrid": radarrID,
		"language": language,
		"forced":   forced,
		"hi":       hi,
	}
	return c.base.Patch(ctx, "/api/movies/subtitles", body, nil)
}

// DeleteMovieSubtitle removes a subtitle file for a movie.
func (c *Client) DeleteMovieSubtitle(ctx context.Context, radarrID int, language, forced, hi, subtitlePath string) error {
	params := url.Values{}
	params.Set("radarrid", strconv.Itoa(radarrID))
	params.Set("language", language)
	params.Set("forced", forced)
	params.Set("hi", hi)
	params.Set("path", subtitlePath)
	path := "/api/movies/subtitles?" + params.Encode()
	return c.base.Delete(ctx, path, nil, nil)
}

// GetProviders returns the status of all subtitle providers.
func (c *Client) GetProviders(ctx context.Context) ([]Provider, error) {
	var out dataWrapper[[]Provider]
	if err := c.base.Get(ctx, "/api/providers", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// ResetProviders resets all throttled subtitle providers.
func (c *Client) ResetProviders(ctx context.Context) error {
	body := map[string]string{"action": "reset"}
	return c.base.Post(ctx, "/api/providers", body, nil)
}

// GetBadges returns the badge counts for the Bazarr UI.
func (c *Client) GetBadges(ctx context.Context) (*BadgeCounts, error) {
	var out BadgeCounts
	if err := c.base.Get(ctx, "/api/badges", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSystemStatus returns Bazarr environment and version information.
func (c *Client) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	var out dataWrapper[SystemStatus]
	if err := c.base.Get(ctx, "/api/system/status", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetHealth returns current health issues.
func (c *Client) GetHealth(ctx context.Context) ([]string, error) {
	var out dataWrapper[[]string]
	if err := c.base.Get(ctx, "/api/system/health", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetLanguages returns all languages configured in Bazarr.
func (c *Client) GetLanguages(ctx context.Context) ([]Language, error) {
	var out dataWrapper[[]Language]
	if err := c.base.Get(ctx, "/api/system/languages", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetLanguageProfiles returns all language profiles.
func (c *Client) GetLanguageProfiles(ctx context.Context) ([]LanguageProfile, error) {
	var out dataWrapper[[]LanguageProfile]
	if err := c.base.Get(ctx, "/api/system/languages/profiles", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// SystemAction performs a system action. Valid actions: "shutdown", "restart".
func (c *Client) SystemAction(ctx context.Context, action string) error {
	body := map[string]string{"action": action}
	return c.base.Post(ctx, "/api/system", body, nil)
}

// Ping checks if the Bazarr instance is reachable.
func (c *Client) Ping(ctx context.Context) error {
	return c.base.Get(ctx, "/api/system/ping", nil)
}
