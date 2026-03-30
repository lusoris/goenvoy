package tvmaze

// Country represents a geographic location.
type Country struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Timezone string `json:"timezone"`
}

// ImageURL holds medium and original resolution URLs for an image.
type ImageURL struct {
	Medium   string `json:"medium"`
	Original string `json:"original"`
}

// Rating contains an average rating value. Average is nil when unrated.
type Rating struct {
	Average *float64 `json:"average"`
}

// Network represents a broadcast TV network.
type Network struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Country      *Country `json:"country"`
	OfficialSite string   `json:"officialSite"`
}

// WebChannel represents a web or streaming service.
type WebChannel struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Country      *Country `json:"country"`
	OfficialSite string   `json:"officialSite"`
}

// ShowSchedule holds the regular airing schedule for a show.
type ShowSchedule struct {
	Time string   `json:"time"`
	Days []string `json:"days"`
}

// Externals holds external service IDs (TVRage, TheTVDB, IMDb).
type Externals struct {
	TVRage  *int   `json:"tvrage"`
	TheTVDB *int   `json:"thetvdb"`
	IMDB    string `json:"imdb"`
}

// Show represents a TV show.
type Show struct {
	ID             int          `json:"id"`
	URL            string       `json:"url"`
	Name           string       `json:"name"`
	Type           string       `json:"type"`
	Language       string       `json:"language"`
	Genres         []string     `json:"genres"`
	Status         string       `json:"status"`
	Runtime        *int         `json:"runtime"`
	AverageRuntime *int         `json:"averageRuntime"`
	Premiered      string       `json:"premiered"`
	Ended          string       `json:"ended"`
	OfficialSite   string       `json:"officialSite"`
	Schedule       ShowSchedule `json:"schedule"`
	Rating         Rating       `json:"rating"`
	Weight         int          `json:"weight"`
	Network        *Network     `json:"network"`
	WebChannel     *WebChannel  `json:"webChannel"`
	Externals      Externals    `json:"externals"`
	Image          *ImageURL    `json:"image"`
	Summary        string       `json:"summary"`
	Updated        int64        `json:"updated"`
}

// Episode represents a TV episode.
type Episode struct {
	ID       int       `json:"id"`
	URL      string    `json:"url"`
	Name     string    `json:"name"`
	Season   int       `json:"season"`
	Number   *int      `json:"number"`
	Type     string    `json:"type"`
	Airdate  string    `json:"airdate"`
	Airtime  string    `json:"airtime"`
	Airstamp string    `json:"airstamp"`
	Runtime  *int      `json:"runtime"`
	Rating   Rating    `json:"rating"`
	Image    *ImageURL `json:"image"`
	Summary  string    `json:"summary"`
}

// Season represents a TV season.
type Season struct {
	ID           int         `json:"id"`
	URL          string      `json:"url"`
	Number       int         `json:"number"`
	Name         string      `json:"name"`
	EpisodeOrder *int        `json:"episodeOrder"`
	PremiereDate string      `json:"premiereDate"`
	EndDate      string      `json:"endDate"`
	Network      *Network    `json:"network"`
	WebChannel   *WebChannel `json:"webChannel"`
	Image        *ImageURL   `json:"image"`
	Summary      string      `json:"summary"`
}

// Person represents a person (actor, crew member, etc.).
type Person struct {
	ID       int       `json:"id"`
	URL      string    `json:"url"`
	Name     string    `json:"name"`
	Country  *Country  `json:"country"`
	Birthday string    `json:"birthday"`
	Deathday string    `json:"deathday"`
	Gender   string    `json:"gender"`
	Image    *ImageURL `json:"image"`
	Updated  int64     `json:"updated"`
}

// Character represents a fictional character.
type Character struct {
	ID    int       `json:"id"`
	URL   string    `json:"url"`
	Name  string    `json:"name"`
	Image *ImageURL `json:"image"`
}

// CastMember pairs a person with the character they play.
type CastMember struct {
	Person    Person    `json:"person"`
	Character Character `json:"character"`
	Self      bool      `json:"self"`
	Voice     bool      `json:"voice"`
}

// CrewMember pairs a crew role type with a person.
type CrewMember struct {
	Type   string `json:"type"`
	Person Person `json:"person"`
}

// AKA is an alternate name for a show in a specific country.
type AKA struct {
	Name    string   `json:"name"`
	Country *Country `json:"country"`
}

// ImageResolution describes a specific resolution of an image.
type ImageResolution struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// ShowImage is an image associated with a show in multiple resolutions.
type ShowImage struct {
	ID          int                        `json:"id"`
	Type        string                     `json:"type"`
	Main        bool                       `json:"main"`
	Resolutions map[string]ImageResolution `json:"resolutions"`
}

// SearchShowResult is a single result from a show search.
type SearchShowResult struct {
	Score float64 `json:"score"`
	Show  Show    `json:"show"`
}

// SearchPersonResult is a single result from a person search.
type SearchPersonResult struct {
	Score  float64 `json:"score"`
	Person Person  `json:"person"`
}

// ScheduleItem is an episode with an embedded show from a schedule query.
type ScheduleItem struct {
	Episode
	Show Show `json:"show"`
}

// Link is a HAL-style hyperlink reference.
type Link struct {
	Href string `json:"href"`
	Name string `json:"name,omitempty"`
}

// CastCreditLinks holds navigation links for a cast credit.
type CastCreditLinks struct {
	Show      Link `json:"show"`
	Character Link `json:"character"`
}

// CastCredit represents a cast role a person plays in a show.
type CastCredit struct {
	Self  bool            `json:"self"`
	Voice bool            `json:"voice"`
	Links CastCreditLinks `json:"_links"`
}

// CrewCreditLinks holds navigation links for a crew credit.
type CrewCreditLinks struct {
	Show Link `json:"show"`
}

// CrewCredit represents a crew role a person holds on a show.
type CrewCredit struct {
	Type  string          `json:"type"`
	Links CrewCreditLinks `json:"_links"`
}

// UpdatePeriod controls the time window for update queries.
type UpdatePeriod string

const (
	// UpdateDay includes updates from the last 24 hours.
	UpdateDay UpdatePeriod = "day"
	// UpdateWeek includes updates from the last 7 days.
	UpdateWeek UpdatePeriod = "week"
	// UpdateMonth includes updates from the last 30 days.
	UpdateMonth UpdatePeriod = "month"
)
