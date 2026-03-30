package tvmaze

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultBaseURL   = "https://api.tvmaze.com"
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

// WithBaseURL overrides the default TVmaze API base URL.
func WithBaseURL(u string) Option {
	return func(cl *Client) { cl.rawBaseURL = u }
}

// Client is a TVmaze API client. No authentication is required.
type Client struct {
	rawBaseURL string
	httpClient *http.Client
	userAgent  string
}

// New creates a TVmaze [Client].
func New(opts ...Option) *Client {
	c := &Client{
		rawBaseURL: defaultBaseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
		userAgent:  defaultUserAgent,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// APIError is returned when the TVmaze API responds with a non-2xx status code.
type APIError struct {
	StatusCode int
	Status     string
}

func (e *APIError) Error() string {
	return "tvmaze: HTTP " + e.Status
}

func (c *Client) get(ctx context.Context, path string, dst any) error {
	u := c.rawBaseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return fmt.Errorf("tvmaze: create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("tvmaze: request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("tvmaze: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	if dst != nil {
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("tvmaze: decode response: %w", err)
		}
	}

	return nil
}

// SearchShows searches for shows by name.
func (c *Client) SearchShows(ctx context.Context, query string) ([]SearchShowResult, error) {
	var results []SearchShowResult
	err := c.get(ctx, "/search/shows?q="+url.QueryEscape(query), &results)
	return results, err
}

// SearchShowSingle returns the single best-matching show for the query.
func (c *Client) SearchShowSingle(ctx context.Context, query string) (*Show, error) {
	var show Show
	if err := c.get(ctx, "/singlesearch/shows?q="+url.QueryEscape(query), &show); err != nil {
		return nil, err
	}
	return &show, nil
}

// LookupShowByTheTVDB looks up a show by its TheTVDB ID.
func (c *Client) LookupShowByTheTVDB(ctx context.Context, thetvdbID int) (*Show, error) {
	var show Show
	if err := c.get(ctx, "/lookup/shows?thetvdb="+strconv.Itoa(thetvdbID), &show); err != nil {
		return nil, err
	}
	return &show, nil
}

// LookupShowByIMDB looks up a show by its IMDb ID.
func (c *Client) LookupShowByIMDB(ctx context.Context, imdbID string) (*Show, error) {
	var show Show
	if err := c.get(ctx, "/lookup/shows?imdb="+url.QueryEscape(imdbID), &show); err != nil {
		return nil, err
	}
	return &show, nil
}

// LookupShowByTVRage looks up a show by its TVRage ID.
func (c *Client) LookupShowByTVRage(ctx context.Context, tvrageID int) (*Show, error) {
	var show Show
	if err := c.get(ctx, "/lookup/shows?tvrage="+strconv.Itoa(tvrageID), &show); err != nil {
		return nil, err
	}
	return &show, nil
}

// SearchPeople searches for people by name.
func (c *Client) SearchPeople(ctx context.Context, query string) ([]SearchPersonResult, error) {
	var results []SearchPersonResult
	err := c.get(ctx, "/search/people?q="+url.QueryEscape(query), &results)
	return results, err
}

// GetShow returns a show by its TVmaze ID.
func (c *Client) GetShow(ctx context.Context, id int) (*Show, error) {
	var show Show
	if err := c.get(ctx, "/shows/"+strconv.Itoa(id), &show); err != nil {
		return nil, err
	}
	return &show, nil
}

// GetShowEpisodes returns all episodes for a show.
func (c *Client) GetShowEpisodes(ctx context.Context, showID int) ([]Episode, error) {
	var episodes []Episode
	err := c.get(ctx, "/shows/"+strconv.Itoa(showID)+"/episodes", &episodes)
	return episodes, err
}

// GetEpisodeByNumber returns a specific episode by season and episode number.
func (c *Client) GetEpisodeByNumber(ctx context.Context, showID, season, number int) (*Episode, error) {
	path := "/shows/" + strconv.Itoa(showID) +
		"/episodebynumber?season=" + strconv.Itoa(season) +
		"&number=" + strconv.Itoa(number)
	var ep Episode
	if err := c.get(ctx, path, &ep); err != nil {
		return nil, err
	}
	return &ep, nil
}

// GetEpisodesByDate returns episodes that aired on the given date (ISO 8601).
func (c *Client) GetEpisodesByDate(ctx context.Context, showID int, date string) ([]Episode, error) {
	var episodes []Episode
	err := c.get(ctx, "/shows/"+strconv.Itoa(showID)+"/episodesbydate?date="+url.QueryEscape(date), &episodes)
	return episodes, err
}

// GetShowSeasons returns all seasons for a show.
func (c *Client) GetShowSeasons(ctx context.Context, showID int) ([]Season, error) {
	var seasons []Season
	err := c.get(ctx, "/shows/"+strconv.Itoa(showID)+"/seasons", &seasons)
	return seasons, err
}

// GetSeasonEpisodes returns all episodes for a specific season by season ID.
func (c *Client) GetSeasonEpisodes(ctx context.Context, seasonID int) ([]Episode, error) {
	var episodes []Episode
	err := c.get(ctx, "/seasons/"+strconv.Itoa(seasonID)+"/episodes", &episodes)
	return episodes, err
}

// GetShowCast returns the cast for a show.
func (c *Client) GetShowCast(ctx context.Context, showID int) ([]CastMember, error) {
	var cast []CastMember
	err := c.get(ctx, "/shows/"+strconv.Itoa(showID)+"/cast", &cast)
	return cast, err
}

// GetShowCrew returns the crew for a show.
func (c *Client) GetShowCrew(ctx context.Context, showID int) ([]CrewMember, error) {
	var crew []CrewMember
	err := c.get(ctx, "/shows/"+strconv.Itoa(showID)+"/crew", &crew)
	return crew, err
}

// GetShowAKAs returns alternate names for a show.
func (c *Client) GetShowAKAs(ctx context.Context, showID int) ([]AKA, error) {
	var akas []AKA
	err := c.get(ctx, "/shows/"+strconv.Itoa(showID)+"/akas", &akas)
	return akas, err
}

// GetShowImages returns all images for a show.
func (c *Client) GetShowImages(ctx context.Context, showID int) ([]ShowImage, error) {
	var images []ShowImage
	err := c.get(ctx, "/shows/"+strconv.Itoa(showID)+"/images", &images)
	return images, err
}

// GetShowIndex returns a paginated list of all shows ordered by ID.
// Each page contains up to 250 shows. The page parameter is zero-indexed.
func (c *Client) GetShowIndex(ctx context.Context, page int) ([]Show, error) {
	var shows []Show
	err := c.get(ctx, "/shows?page="+strconv.Itoa(page), &shows)
	return shows, err
}

// GetEpisode returns an episode by its TVmaze ID.
func (c *Client) GetEpisode(ctx context.Context, id int) (*Episode, error) {
	var ep Episode
	if err := c.get(ctx, "/episodes/"+strconv.Itoa(id), &ep); err != nil {
		return nil, err
	}
	return &ep, nil
}

// GetPerson returns a person by their TVmaze ID.
func (c *Client) GetPerson(ctx context.Context, id int) (*Person, error) {
	var person Person
	if err := c.get(ctx, "/people/"+strconv.Itoa(id), &person); err != nil {
		return nil, err
	}
	return &person, nil
}

// GetPersonCastCredits returns the cast credits for a person.
func (c *Client) GetPersonCastCredits(ctx context.Context, personID int) ([]CastCredit, error) {
	var credits []CastCredit
	err := c.get(ctx, "/people/"+strconv.Itoa(personID)+"/castcredits", &credits)
	return credits, err
}

// GetPersonCrewCredits returns the crew credits for a person.
func (c *Client) GetPersonCrewCredits(ctx context.Context, personID int) ([]CrewCredit, error) {
	var credits []CrewCredit
	err := c.get(ctx, "/people/"+strconv.Itoa(personID)+"/crewcredits", &credits)
	return credits, err
}

// GetSchedule returns the daily TV schedule for a country on a given date.
// Both countryCode and date are optional; defaults to US and the current day.
func (c *Client) GetSchedule(ctx context.Context, countryCode, date string) ([]ScheduleItem, error) {
	params := url.Values{}
	if countryCode != "" {
		params.Set("country", countryCode)
	}
	if date != "" {
		params.Set("date", date)
	}
	path := "/schedule"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	var items []ScheduleItem
	err := c.get(ctx, path, &items)
	return items, err
}

// GetWebSchedule returns the web/streaming schedule for a country on a given date.
// Both countryCode and date are optional.
func (c *Client) GetWebSchedule(ctx context.Context, countryCode, date string) ([]ScheduleItem, error) {
	params := url.Values{}
	if countryCode != "" {
		params.Set("country", countryCode)
	}
	if date != "" {
		params.Set("date", date)
	}
	path := "/schedule/web"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	var items []ScheduleItem
	err := c.get(ctx, path, &items)
	return items, err
}

// GetShowUpdates returns a map of show IDs to their last-updated timestamps.
// The since parameter controls the time window (day, week, or month).
func (c *Client) GetShowUpdates(ctx context.Context, since UpdatePeriod) (map[string]int64, error) {
	path := "/updates/shows"
	if since != "" {
		path += "?since=" + string(since)
	}
	var updates map[string]int64
	err := c.get(ctx, path, &updates)
	return updates, err
}
