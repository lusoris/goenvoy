package tvmaze

import (
	"context"
	"net/url"
	"strconv"

	"github.com/golusoris/goenvoy/metadata"
)

const defaultBaseURL = "https://api.tvmaze.com"

// Client is a TVmaze API client. No authentication is required.
type Client struct {
	*metadata.BaseClient
}

// New creates a TVmaze [Client].
func New(opts ...metadata.Option) *Client {
	bc := metadata.NewBaseClient(defaultBaseURL, "tvmaze", opts...)
	return &Client{BaseClient: bc}
}

// SearchShows searches for shows by name.
func (c *Client) SearchShows(ctx context.Context, query string) ([]SearchShowResult, error) {
	var results []SearchShowResult
	err := c.Get(ctx, "/search/shows?q="+url.QueryEscape(query), &results)
	return results, err
}

// SearchShowSingle returns the single best-matching show for the query.
func (c *Client) SearchShowSingle(ctx context.Context, query string) (*Show, error) {
	var show Show
	if err := c.Get(ctx, "/singlesearch/shows?q="+url.QueryEscape(query), &show); err != nil {
		return nil, err
	}
	return &show, nil
}

// LookupShowByTheTVDB looks up a show by its TheTVDB ID.
func (c *Client) LookupShowByTheTVDB(ctx context.Context, thetvdbID int) (*Show, error) {
	var show Show
	if err := c.Get(ctx, "/lookup/shows?thetvdb="+strconv.Itoa(thetvdbID), &show); err != nil {
		return nil, err
	}
	return &show, nil
}

// LookupShowByIMDB looks up a show by its IMDb ID.
func (c *Client) LookupShowByIMDB(ctx context.Context, imdbID string) (*Show, error) {
	var show Show
	if err := c.Get(ctx, "/lookup/shows?imdb="+url.QueryEscape(imdbID), &show); err != nil {
		return nil, err
	}
	return &show, nil
}

// LookupShowByTVRage looks up a show by its TVRage ID.
func (c *Client) LookupShowByTVRage(ctx context.Context, tvrageID int) (*Show, error) {
	var show Show
	if err := c.Get(ctx, "/lookup/shows?tvrage="+strconv.Itoa(tvrageID), &show); err != nil {
		return nil, err
	}
	return &show, nil
}

// SearchPeople searches for people by name.
func (c *Client) SearchPeople(ctx context.Context, query string) ([]SearchPersonResult, error) {
	var results []SearchPersonResult
	err := c.Get(ctx, "/search/people?q="+url.QueryEscape(query), &results)
	return results, err
}

// GetShow returns a show by its TVmaze ID.
func (c *Client) GetShow(ctx context.Context, id int) (*Show, error) {
	var show Show
	if err := c.Get(ctx, "/shows/"+strconv.Itoa(id), &show); err != nil {
		return nil, err
	}
	return &show, nil
}

// GetShowEpisodes returns all episodes for a show.
func (c *Client) GetShowEpisodes(ctx context.Context, showID int) ([]Episode, error) {
	var episodes []Episode
	err := c.Get(ctx, "/shows/"+strconv.Itoa(showID)+"/episodes", &episodes)
	return episodes, err
}

// GetEpisodeByNumber returns a specific episode by season and episode number.
func (c *Client) GetEpisodeByNumber(ctx context.Context, showID, season, number int) (*Episode, error) {
	path := "/shows/" + strconv.Itoa(showID) +
		"/episodebynumber?season=" + strconv.Itoa(season) +
		"&number=" + strconv.Itoa(number)
	var ep Episode
	if err := c.Get(ctx, path, &ep); err != nil {
		return nil, err
	}
	return &ep, nil
}

// GetEpisodesByDate returns episodes that aired on the given date (ISO 8601).
func (c *Client) GetEpisodesByDate(ctx context.Context, showID int, date string) ([]Episode, error) {
	var episodes []Episode
	err := c.Get(ctx, "/shows/"+strconv.Itoa(showID)+"/episodesbydate?date="+url.QueryEscape(date), &episodes)
	return episodes, err
}

// GetShowSeasons returns all seasons for a show.
func (c *Client) GetShowSeasons(ctx context.Context, showID int) ([]Season, error) {
	var seasons []Season
	err := c.Get(ctx, "/shows/"+strconv.Itoa(showID)+"/seasons", &seasons)
	return seasons, err
}

// GetSeasonEpisodes returns all episodes for a specific season by season ID.
func (c *Client) GetSeasonEpisodes(ctx context.Context, seasonID int) ([]Episode, error) {
	var episodes []Episode
	err := c.Get(ctx, "/seasons/"+strconv.Itoa(seasonID)+"/episodes", &episodes)
	return episodes, err
}

// GetShowCast returns the cast for a show.
func (c *Client) GetShowCast(ctx context.Context, showID int) ([]CastMember, error) {
	var cast []CastMember
	err := c.Get(ctx, "/shows/"+strconv.Itoa(showID)+"/cast", &cast)
	return cast, err
}

// GetShowCrew returns the crew for a show.
func (c *Client) GetShowCrew(ctx context.Context, showID int) ([]CrewMember, error) {
	var crew []CrewMember
	err := c.Get(ctx, "/shows/"+strconv.Itoa(showID)+"/crew", &crew)
	return crew, err
}

// GetShowAKAs returns alternate names for a show.
func (c *Client) GetShowAKAs(ctx context.Context, showID int) ([]AKA, error) {
	var akas []AKA
	err := c.Get(ctx, "/shows/"+strconv.Itoa(showID)+"/akas", &akas)
	return akas, err
}

// GetShowImages returns all images for a show.
func (c *Client) GetShowImages(ctx context.Context, showID int) ([]ShowImage, error) {
	var images []ShowImage
	err := c.Get(ctx, "/shows/"+strconv.Itoa(showID)+"/images", &images)
	return images, err
}

// GetShowIndex returns a paginated list of all shows ordered by ID.
// Each page contains up to 250 shows. The page parameter is zero-indexed.
func (c *Client) GetShowIndex(ctx context.Context, page int) ([]Show, error) {
	var shows []Show
	err := c.Get(ctx, "/shows?page="+strconv.Itoa(page), &shows)
	return shows, err
}

// GetEpisode returns an episode by its TVmaze ID.
func (c *Client) GetEpisode(ctx context.Context, id int) (*Episode, error) {
	var ep Episode
	if err := c.Get(ctx, "/episodes/"+strconv.Itoa(id), &ep); err != nil {
		return nil, err
	}
	return &ep, nil
}

// GetPerson returns a person by their TVmaze ID.
func (c *Client) GetPerson(ctx context.Context, id int) (*Person, error) {
	var person Person
	if err := c.Get(ctx, "/people/"+strconv.Itoa(id), &person); err != nil {
		return nil, err
	}
	return &person, nil
}

// GetPersonCastCredits returns the cast credits for a person.
func (c *Client) GetPersonCastCredits(ctx context.Context, personID int) ([]CastCredit, error) {
	var credits []CastCredit
	err := c.Get(ctx, "/people/"+strconv.Itoa(personID)+"/castcredits", &credits)
	return credits, err
}

// GetPersonCrewCredits returns the crew credits for a person.
func (c *Client) GetPersonCrewCredits(ctx context.Context, personID int) ([]CrewCredit, error) {
	var credits []CrewCredit
	err := c.Get(ctx, "/people/"+strconv.Itoa(personID)+"/crewcredits", &credits)
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
	err := c.Get(ctx, path, &items)
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
	err := c.Get(ctx, path, &items)
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
	err := c.Get(ctx, path, &updates)
	return updates, err
}
