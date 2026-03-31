package trakt

// IDs contains external identifiers for a media item.
type IDs struct {
	Trakt int    `json:"trakt"`
	Slug  string `json:"slug"`
	IMDb  string `json:"imdb,omitempty"`
	TMDb  int    `json:"tmdb,omitempty"`
	TVDb  int    `json:"tvdb,omitempty"`
}

// Movie represents a movie object returned by the API.
type Movie struct {
	Title                 string   `json:"title"`
	Year                  int      `json:"year"`
	IDs                   IDs      `json:"ids"`
	Tagline               string   `json:"tagline,omitempty"`
	Overview              string   `json:"overview,omitempty"`
	Released              string   `json:"released,omitempty"`
	Runtime               int      `json:"runtime,omitempty"`
	Country               string   `json:"country,omitempty"`
	Trailer               string   `json:"trailer,omitempty"`
	Homepage              string   `json:"homepage,omitempty"`
	Status                string   `json:"status,omitempty"`
	Rating                float64  `json:"rating,omitempty"`
	Votes                 int      `json:"votes,omitempty"`
	CommentCount          int      `json:"comment_count,omitempty"`
	Language              string   `json:"language,omitempty"`
	AvailableTranslations []string `json:"available_translations,omitempty"`
	Genres                []string `json:"genres,omitempty"`
	Certification         string   `json:"certification,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
}

// Show represents a TV show object returned by the API.
type Show struct {
	Title                 string   `json:"title"`
	Year                  int      `json:"year"`
	IDs                   IDs      `json:"ids"`
	Overview              string   `json:"overview,omitempty"`
	FirstAired            string   `json:"first_aired,omitempty"`
	Runtime               int      `json:"runtime,omitempty"`
	Certification         string   `json:"certification,omitempty"`
	Network               string   `json:"network,omitempty"`
	Country               string   `json:"country,omitempty"`
	Trailer               string   `json:"trailer,omitempty"`
	Homepage              string   `json:"homepage,omitempty"`
	Status                string   `json:"status,omitempty"`
	Rating                float64  `json:"rating,omitempty"`
	Votes                 int      `json:"votes,omitempty"`
	CommentCount          int      `json:"comment_count,omitempty"`
	Language              string   `json:"language,omitempty"`
	AvailableTranslations []string `json:"available_translations,omitempty"`
	Genres                []string `json:"genres,omitempty"`
	AiredEpisodes         int      `json:"aired_episodes,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
}

// Season represents a TV season.
type Season struct {
	Number        int     `json:"number"`
	IDs           IDs     `json:"ids"`
	Title         string  `json:"title,omitempty"`
	Overview      string  `json:"overview,omitempty"`
	Rating        float64 `json:"rating,omitempty"`
	Votes         int     `json:"votes,omitempty"`
	EpisodeCount  int     `json:"episode_count,omitempty"`
	AiredEpisodes int     `json:"aired_episodes,omitempty"`
	Network       string  `json:"network,omitempty"`
	FirstAired    string  `json:"first_aired,omitempty"`
}

// Episode represents a TV episode.
type Episode struct {
	Season                int      `json:"season"`
	Number                int      `json:"number"`
	Title                 string   `json:"title"`
	IDs                   IDs      `json:"ids"`
	Overview              string   `json:"overview,omitempty"`
	Rating                float64  `json:"rating,omitempty"`
	Votes                 int      `json:"votes,omitempty"`
	CommentCount          int      `json:"comment_count,omitempty"`
	FirstAired            string   `json:"first_aired,omitempty"`
	Runtime               int      `json:"runtime,omitempty"`
	AvailableTranslations []string `json:"available_translations,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
}

// Person represents a person (actor, director, etc.).
type Person struct {
	Name               string `json:"name"`
	IDs                IDs    `json:"ids"`
	Biography          string `json:"biography,omitempty"`
	Birthday           string `json:"birthday,omitempty"`
	Death              string `json:"death,omitempty"`
	Birthplace         string `json:"birthplace,omitempty"`
	Homepage           string `json:"homepage,omitempty"`
	Gender             string `json:"gender,omitempty"`
	KnownForDepartment string `json:"known_for_department,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
}

// TrendingMovie is a movie with its trending watcher count.
type TrendingMovie struct {
	Watchers int   `json:"watchers"`
	Movie    Movie `json:"movie"`
}

// TrendingShow is a show with its trending watcher count.
type TrendingShow struct {
	Watchers int  `json:"watchers"`
	Show     Show `json:"show"`
}

// PlayedMovie is a movie with its play/watch/collect count.
type PlayedMovie struct {
	WatcherCount   int   `json:"watcher_count"`
	PlayCount      int   `json:"play_count"`
	CollectedCount int   `json:"collected_count"`
	Movie          Movie `json:"movie"`
}

// PlayedShow is a show with its play/watch/collect count.
type PlayedShow struct {
	WatcherCount   int  `json:"watcher_count"`
	PlayCount      int  `json:"play_count"`
	CollectedCount int  `json:"collected_count"`
	Show           Show `json:"show"`
}

// AnticipatedMovie is a movie with its list count.
type AnticipatedMovie struct {
	ListCount int   `json:"list_count"`
	Movie     Movie `json:"movie"`
}

// AnticipatedShow is a show with its list count.
type AnticipatedShow struct {
	ListCount int  `json:"list_count"`
	Show      Show `json:"show"`
}

// BoxOfficeMovie is a movie with its revenue.
type BoxOfficeMovie struct {
	Revenue int   `json:"revenue"`
	Movie   Movie `json:"movie"`
}

// MovieTranslation represents a movie translation.
type MovieTranslation struct {
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Tagline  string `json:"tagline"`
	Language string `json:"language"`
	Country  string `json:"country"`
}

// ShowTranslation represents a show translation.
type ShowTranslation struct {
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Language string `json:"language"`
	Country  string `json:"country"`
}

// Ratings contains rating information for a media item.
type Ratings struct {
	Rating       float64      `json:"rating"`
	Votes        int          `json:"votes"`
	Distribution Distribution `json:"distribution"`
}

// Distribution maps rating values (1-10) to their counts.
type Distribution struct {
	One   int `json:"1"`
	Two   int `json:"2"`
	Three int `json:"3"`
	Four  int `json:"4"`
	Five  int `json:"5"`
	Six   int `json:"6"`
	Seven int `json:"7"`
	Eight int `json:"8"`
	Nine  int `json:"9"`
	Ten   int `json:"10"`
}

// Stats contains statistics for a media item.
type Stats struct {
	Watchers        int `json:"watchers"`
	Plays           int `json:"plays"`
	Collectors      int `json:"collectors"`
	Comments        int `json:"comments"`
	Lists           int `json:"lists"`
	Votes           int `json:"votes"`
	Favorited       int `json:"favorited"`
	Recommendations int `json:"recommendations"`
}

// CastMember represents a cast credit.
type CastMember struct {
	Characters []string `json:"characters"`
	Person     Person   `json:"person"`
}

// CrewMember represents a crew credit.
type CrewMember struct {
	Jobs   []string `json:"jobs"`
	Person Person   `json:"person"`
}

// People contains cast and crew for a media item.
type People struct {
	Cast []CastMember `json:"cast"`
	Crew *Crew        `json:"crew,omitempty"`
}

// Crew groups crew members by department.
type Crew struct {
	Production       []CrewMember `json:"production,omitempty"`
	Art              []CrewMember `json:"art,omitempty"`
	Crew             []CrewMember `json:"crew,omitempty"`
	CostumeAndMakeUp []CrewMember `json:"costume & make-up,omitempty"`
	Directing        []CrewMember `json:"directing,omitempty"`
	Writing          []CrewMember `json:"writing,omitempty"`
	Sound            []CrewMember `json:"sound,omitempty"`
	Camera           []CrewMember `json:"camera,omitempty"`
	VisualEffects    []CrewMember `json:"visual effects,omitempty"`
	Lighting         []CrewMember `json:"lighting,omitempty"`
	Editing          []CrewMember `json:"editing,omitempty"`
}

// Studio represents a production studio.
type Studio struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	IDs     IDs    `json:"ids"`
}

// SearchResult is a single result from the search endpoint.
type SearchResult struct {
	Type    string   `json:"type"`
	Score   float64  `json:"score"`
	Movie   *Movie   `json:"movie,omitempty"`
	Show    *Show    `json:"show,omitempty"`
	Episode *Episode `json:"episode,omitempty"`
	Person  *Person  `json:"person,omitempty"`
}

// CalendarMovie is a movie in a calendar list.
type CalendarMovie struct {
	Released string `json:"released"`
	Movie    Movie  `json:"movie"`
}

// CalendarShow is a show episode in a calendar list.
type CalendarShow struct {
	FirstAired string  `json:"first_aired"`
	Episode    Episode `json:"episode"`
	Show       Show    `json:"show"`
}

// Genre represents a content genre.
type Genre struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Certification represents a content certification.
type Certification struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// Alias represents a title alias.
type Alias struct {
	Title   string `json:"title"`
	Country string `json:"country"`
}

// MovieRelease represents a movie release date and certification.
type MovieRelease struct {
	Country       string `json:"country"`
	Certification string `json:"certification"`
	ReleaseDate   string `json:"release_date"`
	ReleaseType   string `json:"release_type"`
	Note          string `json:"note,omitempty"`
}

// Country represents a country.
type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Language represents a language.
type Language struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Network represents a TV network.
type Network struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	IDs     IDs    `json:"ids"`
}

// OAuth2 types.

// DeviceCode holds the response from the device code request.
type DeviceCode struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// Token holds OAuth2 access and refresh tokens.
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
}
