package tvdb

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

	"github.com/lusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://api4.thetvdb.com/v4"

// Client is a TheTVDB API v4 client.
type Client struct {
	*metadata.BaseClient
	apiKey string
	pin    string

	mu    sync.Mutex
	token string
}

// New creates a TheTVDB [Client] using the given API key.
func New(apiKey string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "tvdb", opts...)
	return &Client{BaseClient: bc, apiKey: apiKey}
}

// NewWithPIN creates a TheTVDB [Client] with an API key and subscriber PIN.
func NewWithPIN(apiKey, pin string, opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "tvdb", opts...)
	return &Client{BaseClient: bc, apiKey: apiKey, pin: pin}
}

// Login authenticates with TheTVDB and caches the bearer token.
// It is called automatically on the first API request if no token is cached.
func (c *Client) Login(ctx context.Context) error {
	body := LoginRequest{APIKey: c.apiKey}
	if c.pin != "" {
		body.PIN = c.pin
	}

	payload, err := json.Marshal(body) // #nosec G117 -- API key is intentionally sent to the auth endpoint
	if err != nil {
		return fmt.Errorf("tvdb: marshal login body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL()+"/login", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("tvdb: create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent())

	resp, err := c.HTTPClient().Do(req)
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

	err := c.doGet(ctx, path, dst)
	if err == nil {
		return nil
	}

	// On 401 Unauthorized, re-login and retry once.
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized {
		c.mu.Lock()
		c.token = ""
		c.mu.Unlock()
		if loginErr := c.Login(ctx); loginErr != nil {
			return loginErr
		}
		return c.doGet(ctx, path, dst)
	}
	return err
}

func (c *Client) doGet(ctx context.Context, path string, dst any) error {
	c.mu.Lock()
	token := c.token
	c.mu.Unlock()

	c.SetAuth(func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+token)
	})

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

// Artwork.

// GetArtworkStatuses returns all artwork status records.
func (c *Client) GetArtworkStatuses(ctx context.Context) ([]ArtworkStatus, error) {
	var out response[[]ArtworkStatus]
	if err := c.get(ctx, "/artwork/statuses", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Awards.

// GetAwards returns all award base records.
func (c *Client) GetAwards(ctx context.Context) ([]AwardBase, error) {
	var out response[[]AwardBase]
	if err := c.get(ctx, "/awards", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetAward returns an award base record by ID.
func (c *Client) GetAward(ctx context.Context, id int) (*AwardBase, error) {
	var out response[AwardBase]
	if err := c.get(ctx, "/awards/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetAwardExtended returns an award extended record.
func (c *Client) GetAwardExtended(ctx context.Context, id int) (*AwardExtended, error) {
	var out response[AwardExtended]
	if err := c.get(ctx, "/awards/"+strconv.Itoa(id)+"/extended", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetAwardCategory returns an award category base record.
func (c *Client) GetAwardCategory(ctx context.Context, id int) (*AwardCategory, error) {
	var out response[AwardCategory]
	if err := c.get(ctx, "/awards/categories/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetAwardCategoryExtended returns an award category extended record.
func (c *Client) GetAwardCategoryExtended(ctx context.Context, id int) (*AwardCategoryExtended, error) {
	var out response[AwardCategoryExtended]
	if err := c.get(ctx, "/awards/categories/"+strconv.Itoa(id)+"/extended", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// Companies.

// GetCompanies returns a paginated list of company records.
func (c *Client) GetCompanies(ctx context.Context, page int) ([]Company, error) {
	var out response[[]Company]
	if err := c.get(ctx, "/companies?page="+strconv.Itoa(page), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetCompanyTypes returns all company type records.
func (c *Client) GetCompanyTypes(ctx context.Context) ([]CompanyType, error) {
	var out response[[]CompanyType]
	if err := c.get(ctx, "/companies/types", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetCompany returns a company record by ID.
func (c *Client) GetCompany(ctx context.Context, id int) (*Company, error) {
	var out response[Company]
	if err := c.get(ctx, "/companies/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// Countries.

// GetCountries returns all country records.
func (c *Client) GetCountries(ctx context.Context) ([]Country, error) {
	var out response[[]Country]
	if err := c.get(ctx, "/countries", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Entity Types.

// GetEntityTypes returns all entity type records.
func (c *Client) GetEntityTypes(ctx context.Context) ([]EntityType, error) {
	var out response[[]EntityType]
	if err := c.get(ctx, "/entities", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Episodes (paginated).

// GetEpisodes returns a paginated list of episodes.
func (c *Client) GetEpisodes(ctx context.Context, page int) ([]EpisodeBase, error) {
	var out response[[]EpisodeBase]
	if err := c.get(ctx, "/episodes?page="+strconv.Itoa(page), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Genders.

// GetGenders returns all gender records.
func (c *Client) GetGenders(ctx context.Context) ([]Gender, error) {
	var out response[[]Gender]
	if err := c.get(ctx, "/genders", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Genres.

// GetGenre returns a single genre record by ID.
func (c *Client) GetGenre(ctx context.Context, id int) (*Genre, error) {
	var out response[Genre]
	if err := c.get(ctx, "/genres/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// Inspiration Types.

// GetInspirationTypes returns all inspiration type records.
func (c *Client) GetInspirationTypes(ctx context.Context) ([]InspirationType, error) {
	var out response[[]InspirationType]
	if err := c.get(ctx, "/inspiration/types", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Lists.

// GetLists returns a paginated list of list records.
func (c *Client) GetLists(ctx context.Context, page int) ([]ListBase, error) {
	var out response[[]ListBase]
	if err := c.get(ctx, "/lists?page="+strconv.Itoa(page), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetList returns a list record by ID.
func (c *Client) GetList(ctx context.Context, id int) (*ListBase, error) {
	var out response[ListBase]
	if err := c.get(ctx, "/lists/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetListBySlug returns a list record by slug.
func (c *Client) GetListBySlug(ctx context.Context, slug string) (*ListBase, error) {
	var out response[ListBase]
	if err := c.get(ctx, "/lists/slug/"+url.PathEscape(slug), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetListExtended returns a list extended record.
func (c *Client) GetListExtended(ctx context.Context, id int) (*ListExtended, error) {
	var out response[ListExtended]
	if err := c.get(ctx, "/lists/"+strconv.Itoa(id)+"/extended", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetListTranslation returns a list translation for the given language.
func (c *Client) GetListTranslation(ctx context.Context, id int, language string) (*Translation, error) {
	var out response[Translation]
	path := "/lists/" + strconv.Itoa(id) + "/translations/" + url.PathEscape(language)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// Movies (paginated, filter, slug, statuses).

// GetMovies returns a paginated list of movies.
func (c *Client) GetMovies(ctx context.Context, page int) ([]MovieBase, error) {
	var out response[[]MovieBase]
	if err := c.get(ctx, "/movies?page="+strconv.Itoa(page), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// FilterMovies returns movies matching the given filter parameters.
func (c *Client) FilterMovies(ctx context.Context, params *FilterParams) ([]MovieBase, error) {
	var out response[[]MovieBase]
	if err := c.get(ctx, "/movies/filter?"+params.encode(), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetMovieBySlug returns a movie record by slug.
func (c *Client) GetMovieBySlug(ctx context.Context, slug string) (*MovieBase, error) {
	var out response[MovieBase]
	if err := c.get(ctx, "/movies/slug/"+url.PathEscape(slug), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetMovieStatuses returns all movie status records.
func (c *Client) GetMovieStatuses(ctx context.Context) ([]Status, error) {
	var out response[[]Status]
	if err := c.get(ctx, "/movies/statuses", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// People (paginated, types).

// GetPeople returns a paginated list of people.
func (c *Client) GetPeople(ctx context.Context, page int) ([]PersonBase, error) {
	var out response[[]PersonBase]
	if err := c.get(ctx, "/people?page="+strconv.Itoa(page), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetPeopleTypes returns all people type records.
func (c *Client) GetPeopleTypes(ctx context.Context) ([]PeopleType, error) {
	var out response[[]PeopleType]
	if err := c.get(ctx, "/people/types", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Seasons (paginated, types).

// GetSeasons returns a paginated list of seasons.
func (c *Client) GetSeasons(ctx context.Context, page int) ([]SeasonBase, error) {
	var out response[[]SeasonBase]
	if err := c.get(ctx, "/seasons?page="+strconv.Itoa(page), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetSeasonTypes returns all season type records.
func (c *Client) GetSeasonTypes(ctx context.Context) ([]SeasonType, error) {
	var out response[[]SeasonType]
	if err := c.get(ctx, "/seasons/types", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Series (paginated, filter, slug, statuses, episodes+lang).

// GetAllSeries returns a paginated list of series.
func (c *Client) GetAllSeries(ctx context.Context, page int) ([]SeriesBase, error) {
	var out response[[]SeriesBase]
	if err := c.get(ctx, "/series?page="+strconv.Itoa(page), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetSeriesEpisodesWithLang returns translated episodes for a series by season type and language.
func (c *Client) GetSeriesEpisodesWithLang(ctx context.Context, id int, seasonType, language string, page int) (*SeriesEpisodesResult, error) {
	var out response[SeriesEpisodesResult]
	path := fmt.Sprintf("/series/%d/episodes/%s/%s?page=%d", id, url.PathEscape(seasonType), url.PathEscape(language), page)
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// FilterSeries returns series matching the given filter parameters.
func (c *Client) FilterSeries(ctx context.Context, params *FilterParams) ([]SeriesBase, error) {
	var out response[[]SeriesBase]
	if err := c.get(ctx, "/series/filter?"+params.encode(), &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// GetSeriesBySlug returns a series record by slug.
func (c *Client) GetSeriesBySlug(ctx context.Context, slug string) (*SeriesBase, error) {
	var out response[SeriesBase]
	if err := c.get(ctx, "/series/slug/"+url.PathEscape(slug), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetSeriesStatuses returns all series status records.
func (c *Client) GetSeriesStatuses(ctx context.Context) ([]Status, error) {
	var out response[[]Status]
	if err := c.get(ctx, "/series/statuses", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// Source Types.

// GetSourceTypes returns all source type records.
func (c *Client) GetSourceTypes(ctx context.Context) ([]SourceType, error) {
	var out response[[]SourceType]
	if err := c.get(ctx, "/sources/types", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// User.

// GetUserInfo returns the current user's info.
func (c *Client) GetUserInfo(ctx context.Context) (*UserInfo, error) {
	var out response[UserInfo]
	if err := c.get(ctx, "/user", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetUserByID returns a user's info by their ID.
func (c *Client) GetUserByID(ctx context.Context, id int) (*UserInfo, error) {
	var out response[UserInfo]
	if err := c.get(ctx, "/user/"+strconv.Itoa(id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// GetUserFavorites returns the current user's favorites.
func (c *Client) GetUserFavorites(ctx context.Context) (*Favorites, error) {
	var out response[Favorites]
	if err := c.get(ctx, "/user/favorites", &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// AddUserFavorites adds items to the current user's favorites.
func (c *Client) AddUserFavorites(ctx context.Context, record *FavoriteRecord) error {
	return c.post(ctx, "/user/favorites", record)
}

// post is a helper for POST requests.
func (c *Client) post(ctx context.Context, path string, body any) error {
	if err := c.ensureToken(ctx); err != nil {
		return err
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("tvdb: marshal body: %w", err)
	}

	c.mu.Lock()
	token := c.token
	c.mu.Unlock()

	c.SetAuth(func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
	})

	respBody, status, err := c.DoRaw(ctx, http.MethodPost, path, bytes.NewReader(data))
	if err != nil {
		return err
	}

	if status < 200 || status >= 300 {
		apiErr := &APIError{StatusCode: status}
		if err := json.Unmarshal(respBody, apiErr); err != nil {
			apiErr.RawBody = string(respBody)
		}
		return apiErr
	}
	return nil
}
