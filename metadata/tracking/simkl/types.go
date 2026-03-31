package simkl

// IDs contains external identifiers for a media item.
type IDs struct {
	Simkl       int    `json:"simkl"`
	SimklID     int    `json:"simkl_id"`
	Slug        string `json:"slug,omitempty"`
	IMDb        string `json:"imdb,omitempty"`
	TMDb        string `json:"tmdb,omitempty"`
	TVDb        string `json:"tvdb,omitempty"`
	MAL         string `json:"mal,omitempty"`
	AniDB       string `json:"anidb,omitempty"`
	AniList     string `json:"anilist,omitempty"`
	Kitsu       string `json:"kitsu,omitempty"`
	LiveChart   string `json:"livechart,omitempty"`
	AniSearch   string `json:"anisearch,omitempty"`
	AnimePlanet string `json:"animeplanet,omitempty"`
	Netflix     string `json:"netflix,omitempty"`
	Hulu        string `json:"hulu,omitempty"`
	Crunchyroll string `json:"crunchyroll,omitempty"`
	TraktSlug   string `json:"traktslug,omitempty"`
	LetterSlug  string `json:"letterslug,omitempty"`
}

// RatingScore holds a rating value and vote count.
type RatingScore struct {
	Rating float64 `json:"rating"`
	Votes  int     `json:"votes"`
}

// Ratings holds ratings from various sources.
type Ratings struct {
	Simkl *RatingScore `json:"simkl,omitempty"`
	IMDb  *RatingScore `json:"imdb,omitempty"`
	MAL   *RatingScore `json:"mal,omitempty"`
}

// Trailer represents a video trailer.
type Trailer struct {
	Name    string `json:"name"`
	YouTube string `json:"youtube"`
	Size    int    `json:"size"`
}

// AirDate holds airing schedule information.
type AirDate struct {
	Day      string `json:"day"`
	Time     string `json:"time"`
	Timezone string `json:"timezone"`
}

// AlternativeTitle holds an alternative name.
type AlternativeTitle struct {
	Name string `json:"name"`
	Lang string `json:"lang"`
	Type string `json:"type"`
}

// ReleaseDate holds a release date entry.
type ReleaseDate struct {
	Type        int    `json:"type"`
	ReleaseDate string `json:"release_date"`
}

// LocalReleaseDate holds release dates for a country.
type LocalReleaseDate struct {
	ISO31661 string        `json:"iso_3166_1"`
	Results  []ReleaseDate `json:"results"`
}

// Movie represents a Simkl movie.
type Movie struct {
	Title           string             `json:"title"`
	Year            int                `json:"year"`
	Type            string             `json:"type,omitempty"`
	IDs             IDs                `json:"ids"`
	Rank            int                `json:"rank,omitempty"`
	Poster          string             `json:"poster,omitempty"`
	Fanart          string             `json:"fanart,omitempty"`
	Released        string             `json:"released,omitempty"`
	Runtime         int                `json:"runtime,omitempty"`
	Director        string             `json:"director,omitempty"`
	Certification   string             `json:"certification,omitempty"`
	Budget          int                `json:"budget,omitempty"`
	Revenue         int                `json:"revenue,omitempty"`
	Overview        string             `json:"overview,omitempty"`
	Genres          []string           `json:"genres,omitempty"`
	Countries       string             `json:"countries,omitempty"`
	Languages       string             `json:"languages,omitempty"`
	AltTitles       []AlternativeTitle `json:"alt_titles,omitempty"`
	Ratings         *Ratings           `json:"ratings,omitempty"`
	Trailers        []Trailer          `json:"trailers,omitempty"`
	ReleaseDates    []LocalReleaseDate `json:"release_dates,omitempty"`
	Recommendations []MovieShort       `json:"users_recommendations,omitempty"`
}

// MovieShort is a minimal movie reference.
type MovieShort struct {
	Title  string `json:"title"`
	Year   int    `json:"year"`
	Poster string `json:"poster,omitempty"`
	IDs    IDs    `json:"ids"`
}

// Show represents a Simkl TV show.
type Show struct {
	Title           string      `json:"title"`
	Year            int         `json:"year"`
	Type            string      `json:"type,omitempty"`
	IDs             IDs         `json:"ids"`
	YearStartEnd    string      `json:"year_start_end,omitempty"`
	Rank            int         `json:"rank,omitempty"`
	Poster          string      `json:"poster,omitempty"`
	Fanart          string      `json:"fanart,omitempty"`
	FirstAired      string      `json:"first_aired,omitempty"`
	LastAired       string      `json:"last_aired,omitempty"`
	Airs            *AirDate    `json:"airs,omitempty"`
	Runtime         int         `json:"runtime,omitempty"`
	Certification   string      `json:"certification,omitempty"`
	Overview        string      `json:"overview,omitempty"`
	Genres          []string    `json:"genres,omitempty"`
	Country         string      `json:"country,omitempty"`
	TotalEpisodes   int         `json:"total_episodes,omitempty"`
	Status          string      `json:"status,omitempty"`
	Network         string      `json:"network,omitempty"`
	Ratings         *Ratings    `json:"ratings,omitempty"`
	Trailers        []Trailer   `json:"trailers,omitempty"`
	Recommendations []ShowShort `json:"user_recommendations,omitempty"`
}

// ShowShort is a minimal show reference.
type ShowShort struct {
	Title        string  `json:"title"`
	Year         int     `json:"year"`
	Poster       string  `json:"poster,omitempty"`
	UsersPercent float64 `json:"users_percent,omitempty"`
	UsersCount   int     `json:"users_count,omitempty"`
	IDs          IDs     `json:"ids"`
}

// Anime represents a Simkl anime.
type Anime struct {
	Title         string    `json:"title"`
	Year          int       `json:"year"`
	Type          string    `json:"type,omitempty"`
	AnimeType     string    `json:"anime_type,omitempty"`
	EnTitle       string    `json:"en_title,omitempty"`
	IDs           IDs       `json:"ids"`
	YearStartEnd  string    `json:"year_start_end,omitempty"`
	Rank          int       `json:"rank,omitempty"`
	Poster        string    `json:"poster,omitempty"`
	Fanart        string    `json:"fanart,omitempty"`
	FirstAired    string    `json:"first_aired,omitempty"`
	LastAired     string    `json:"last_aired,omitempty"`
	Airs          *AirDate  `json:"airs,omitempty"`
	Runtime       int       `json:"runtime,omitempty"`
	Certification string    `json:"certification,omitempty"`
	Overview      string    `json:"overview,omitempty"`
	Genres        []string  `json:"genres,omitempty"`
	Country       string    `json:"country,omitempty"`
	TotalEpisodes int       `json:"total_episodes,omitempty"`
	Status        string    `json:"status,omitempty"`
	Network       string    `json:"network,omitempty"`
	Ratings       *Ratings  `json:"ratings,omitempty"`
	Trailers      []Trailer `json:"trailers,omitempty"`
}

// Episode represents a TV show or anime episode.
type Episode struct {
	Title   string `json:"title"`
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
	Type    string `json:"type,omitempty"`
	IDs     IDs    `json:"ids,omitempty"`
	Img     string `json:"img,omitempty"`
	Date    string `json:"date,omitempty"`
	Desc    string `json:"description,omitempty"`
}

// EpisodeMinimal is a minimal episode in airing schedules.
type EpisodeMinimal struct {
	Title   string `json:"title"`
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
}

// TrendingMovie represents a trending movie entry.
type TrendingMovie struct {
	Title       string   `json:"title"`
	URL         string   `json:"url,omitempty"`
	Poster      string   `json:"poster,omitempty"`
	Fanart      string   `json:"fanart,omitempty"`
	IDs         IDs      `json:"ids"`
	ReleaseDate string   `json:"release_date,omitempty"`
	Rank        int      `json:"rank,omitempty"`
	DropRate    string   `json:"drop_rate,omitempty"`
	Watched     int      `json:"watched,omitempty"`
	PlanToWatch int      `json:"plan_to_watch,omitempty"`
	Ratings     *Ratings `json:"ratings,omitempty"`
	Country     string   `json:"country,omitempty"`
	Runtime     string   `json:"runtime,omitempty"`
	Status      string   `json:"status,omitempty"`
	DVDDate     string   `json:"dvd_date,omitempty"`
	Metadata    string   `json:"metadata,omitempty"`
	Overview    string   `json:"overview,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Theater     string   `json:"theater,omitempty"`
}

// TrendingShow represents a trending show entry.
type TrendingShow struct {
	Title         string   `json:"title,omitempty"`
	URL           string   `json:"url,omitempty"`
	Poster        string   `json:"poster,omitempty"`
	Fanart        string   `json:"fanart,omitempty"`
	IDs           IDs      `json:"ids"`
	ReleaseDate   string   `json:"release_date,omitempty"`
	Rank          int      `json:"rank,omitempty"`
	DropRate      string   `json:"drop_rate,omitempty"`
	Watched       int      `json:"watched,omitempty"`
	PlanToWatch   int      `json:"plan_to_watch,omitempty"`
	Ratings       *Ratings `json:"ratings,omitempty"`
	Country       string   `json:"country,omitempty"`
	Runtime       string   `json:"runtime,omitempty"`
	Status        string   `json:"status,omitempty"`
	TotalEpisodes int      `json:"total_episodes,omitempty"`
	Network       string   `json:"network,omitempty"`
	Metadata      string   `json:"metadata,omitempty"`
	Overview      string   `json:"overview,omitempty"`
	Genres        []string `json:"genres,omitempty"`
}

// TrendingAnime represents a trending anime entry.
type TrendingAnime struct {
	Title         string   `json:"title,omitempty"`
	URL           string   `json:"url,omitempty"`
	Poster        string   `json:"poster,omitempty"`
	Fanart        string   `json:"fanart,omitempty"`
	IDs           IDs      `json:"ids"`
	ReleaseDate   string   `json:"release_date,omitempty"`
	Rank          int      `json:"rank,omitempty"`
	DropRate      string   `json:"drop_rate,omitempty"`
	Watched       int      `json:"watched,omitempty"`
	PlanToWatch   int      `json:"plan_to_watch,omitempty"`
	Ratings       *Ratings `json:"ratings,omitempty"`
	Country       string   `json:"country,omitempty"`
	Runtime       string   `json:"runtime,omitempty"`
	Status        string   `json:"status,omitempty"`
	TotalEpisodes int      `json:"total_episodes,omitempty"`
	Network       string   `json:"network,omitempty"`
}

// GenreItem represents an item from genre filtered results.
type GenreItem struct {
	Title   string   `json:"title"`
	Year    int      `json:"year"`
	Date    string   `json:"date,omitempty"`
	URL     string   `json:"url,omitempty"`
	Poster  string   `json:"poster,omitempty"`
	Fanart  string   `json:"fanart,omitempty"`
	Rank    int      `json:"rank,omitempty"`
	IDs     IDs      `json:"ids"`
	Ratings *Ratings `json:"ratings,omitempty"`
}

// PremiereItem represents a premiering show.
type PremiereItem struct {
	Title   string   `json:"title"`
	Year    int      `json:"year"`
	Date    string   `json:"date,omitempty"`
	URL     string   `json:"url,omitempty"`
	Poster  string   `json:"poster,omitempty"`
	Rank    int      `json:"rank,omitempty"`
	IDs     IDs      `json:"ids"`
	Ratings *Ratings `json:"ratings,omitempty"`
}

// AiringItem represents an airing show or anime.
type AiringItem struct {
	Title   string          `json:"title"`
	Year    int             `json:"year"`
	Date    string          `json:"date,omitempty"`
	URL     string          `json:"url,omitempty"`
	Poster  string          `json:"poster,omitempty"`
	Rank    int             `json:"rank,omitempty"`
	IDs     IDs             `json:"ids"`
	Episode *EpisodeMinimal `json:"episode,omitempty"`
}

// BestItem represents an item from best-of lists.
type BestItem struct {
	Title   string   `json:"title"`
	Year    int      `json:"year"`
	Poster  string   `json:"poster,omitempty"`
	URL     string   `json:"url,omitempty"`
	IDs     IDs      `json:"ids"`
	Ratings *Ratings `json:"ratings,omitempty"`
}

// SearchResult represents a search result from the text search endpoint.
type SearchResult struct {
	Title         string   `json:"title"`
	Poster        string   `json:"poster,omitempty"`
	Year          int      `json:"year,omitempty"`
	Type          string   `json:"type,omitempty"`
	TitleEn       string   `json:"title_en,omitempty"`
	TitleRomaji   string   `json:"title_romaji,omitempty"`
	AllTitles     []string `json:"all_titles,omitempty"`
	URL           string   `json:"url,omitempty"`
	EpCount       int      `json:"ep_count,omitempty"`
	Rank          int      `json:"rank,omitempty"`
	Status        string   `json:"status,omitempty"`
	TotalEpisodes int      `json:"total_episodes,omitempty"`
	Ratings       *Ratings `json:"ratings,omitempty"`
	IDs           IDs      `json:"ids"`
}

// SearchIDResult represents a result from the ID lookup endpoint.
type SearchIDResult struct {
	Title         string `json:"title"`
	Poster        string `json:"poster,omitempty"`
	Year          int    `json:"year,omitempty"`
	Type          string `json:"type"`
	TotalEpisodes int    `json:"total_episodes,omitempty"`
	Status        string `json:"status,omitempty"`
	IDs           IDs    `json:"ids"`
}

// DeviceCode holds the response from the PIN/device code request.
type DeviceCode struct {
	Result          string `json:"result"`
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// CalendarShow represents a show in the TV calendar.
type CalendarShow struct {
	Title   string `json:"title"`
	Poster  string `json:"poster,omitempty"`
	Date    string `json:"date,omitempty"`
	Episode int    `json:"episode,omitempty"`
	Season  int    `json:"season,omitempty"`
	IDs     IDs    `json:"ids"`
}

// CalendarAnime represents an anime in the anime calendar.
type CalendarAnime struct {
	Title   string `json:"title"`
	Poster  string `json:"poster,omitempty"`
	Date    string `json:"date,omitempty"`
	Episode int    `json:"episode,omitempty"`
	IDs     IDs    `json:"ids"`
}

// CalendarMovie represents a movie release in the calendar.
type CalendarMovie struct {
	Title       string `json:"title"`
	Poster      string `json:"poster,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
	IDs         IDs    `json:"ids"`
}
