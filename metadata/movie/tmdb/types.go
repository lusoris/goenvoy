package tmdb

import "fmt"

// PaginatedResult is a generic paginated response from TMDb.
type PaginatedResult[T any] struct {
	Page         int `json:"page"`
	Results      []T `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

// MovieResult represents a movie in search or discover results.
type MovieResult struct {
	ID               int     `json:"id"`
	Title            string  `json:"title,omitempty"`
	OriginalTitle    string  `json:"original_title,omitempty"`
	Overview         string  `json:"overview,omitempty"`
	PosterPath       string  `json:"poster_path,omitempty"`
	BackdropPath     string  `json:"backdrop_path,omitempty"`
	ReleaseDate      string  `json:"release_date,omitempty"`
	OriginalLanguage string  `json:"original_language,omitempty"`
	GenreIDs         []int   `json:"genre_ids,omitempty"`
	Popularity       float64 `json:"popularity,omitempty"`
	VoteAverage      float64 `json:"vote_average,omitempty"`
	VoteCount        int     `json:"vote_count,omitempty"`
	Adult            bool    `json:"adult,omitempty"`
	Video            bool    `json:"video,omitempty"`
}

// TVResult represents a TV show in search or discover results.
type TVResult struct {
	ID               int      `json:"id"`
	Name             string   `json:"name,omitempty"`
	OriginalName     string   `json:"original_name,omitempty"`
	Overview         string   `json:"overview,omitempty"`
	PosterPath       string   `json:"poster_path,omitempty"`
	BackdropPath     string   `json:"backdrop_path,omitempty"`
	FirstAirDate     string   `json:"first_air_date,omitempty"`
	OriginalLanguage string   `json:"original_language,omitempty"`
	OriginCountry    []string `json:"origin_country,omitempty"`
	GenreIDs         []int    `json:"genre_ids,omitempty"`
	Popularity       float64  `json:"popularity,omitempty"`
	VoteAverage      float64  `json:"vote_average,omitempty"`
	VoteCount        int      `json:"vote_count,omitempty"`
}

// PersonResult represents a person in search results.
type PersonResult struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name,omitempty"`
	ProfilePath        string  `json:"profile_path,omitempty"`
	Adult              bool    `json:"adult,omitempty"`
	Popularity         float64 `json:"popularity,omitempty"`
	KnownForDepartment string  `json:"known_for_department,omitempty"`
	Gender             int     `json:"gender,omitempty"`
}

// MultiResult represents a result from the multi-search endpoint.
// MediaType indicates whether it is a "movie", "tv", or "person".
type MultiResult struct {
	ID               int      `json:"id"`
	MediaType        string   `json:"media_type,omitempty"`
	Title            string   `json:"title,omitempty"`
	Name             string   `json:"name,omitempty"`
	OriginalTitle    string   `json:"original_title,omitempty"`
	OriginalName     string   `json:"original_name,omitempty"`
	Overview         string   `json:"overview,omitempty"`
	PosterPath       string   `json:"poster_path,omitempty"`
	BackdropPath     string   `json:"backdrop_path,omitempty"`
	ProfilePath      string   `json:"profile_path,omitempty"`
	ReleaseDate      string   `json:"release_date,omitempty"`
	FirstAirDate     string   `json:"first_air_date,omitempty"`
	OriginalLanguage string   `json:"original_language,omitempty"`
	OriginCountry    []string `json:"origin_country,omitempty"`
	GenreIDs         []int    `json:"genre_ids,omitempty"`
	Popularity       float64  `json:"popularity,omitempty"`
	VoteAverage      float64  `json:"vote_average,omitempty"`
	VoteCount        int      `json:"vote_count,omitempty"`
	Adult            bool     `json:"adult,omitempty"`
}

// Genre represents a movie or TV genre.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// ProductionCompany represents a production company.
type ProductionCompany struct {
	ID            int    `json:"id"`
	Name          string `json:"name,omitempty"`
	LogoPath      string `json:"logo_path,omitempty"`
	OriginCountry string `json:"origin_country,omitempty"`
}

// ProductionCountry represents a country of production.
type ProductionCountry struct {
	ISO31661 string `json:"iso_3166_1,omitempty"`
	Name     string `json:"name,omitempty"`
}

// SpokenLanguage represents a spoken language.
type SpokenLanguage struct {
	ISO6391     string `json:"iso_639_1,omitempty"`
	EnglishName string `json:"english_name,omitempty"`
	Name        string `json:"name,omitempty"`
}

// Collection is a brief collection reference.
type Collection struct {
	ID           int    `json:"id"`
	Name         string `json:"name,omitempty"`
	PosterPath   string `json:"poster_path,omitempty"`
	BackdropPath string `json:"backdrop_path,omitempty"`
}

// MovieDetails contains full details for a movie.
type MovieDetails struct {
	ID                  int                 `json:"id"`
	IMDbID              string              `json:"imdb_id,omitempty"`
	Title               string              `json:"title,omitempty"`
	OriginalTitle       string              `json:"original_title,omitempty"`
	Overview            string              `json:"overview,omitempty"`
	Tagline             string              `json:"tagline,omitempty"`
	Status              string              `json:"status,omitempty"`
	Homepage            string              `json:"homepage,omitempty"`
	PosterPath          string              `json:"poster_path,omitempty"`
	BackdropPath        string              `json:"backdrop_path,omitempty"`
	ReleaseDate         string              `json:"release_date,omitempty"`
	OriginalLanguage    string              `json:"original_language,omitempty"`
	OriginCountry       []string            `json:"origin_country,omitempty"`
	Adult               bool                `json:"adult,omitempty"`
	Video               bool                `json:"video,omitempty"`
	Budget              int64               `json:"budget,omitempty"`
	Revenue             int64               `json:"revenue,omitempty"`
	Runtime             int                 `json:"runtime,omitempty"`
	Popularity          float64             `json:"popularity,omitempty"`
	VoteAverage         float64             `json:"vote_average,omitempty"`
	VoteCount           int                 `json:"vote_count,omitempty"`
	Genres              []Genre             `json:"genres,omitempty"`
	BelongsToCollection *Collection         `json:"belongs_to_collection,omitempty"`
	ProductionCompanies []ProductionCompany `json:"production_companies,omitempty"`
	ProductionCountries []ProductionCountry `json:"production_countries,omitempty"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages,omitempty"`
}

// CreatedBy represents a TV series creator.
type CreatedBy struct {
	ID          int    `json:"id"`
	Name        string `json:"name,omitempty"`
	Gender      int    `json:"gender,omitempty"`
	ProfilePath string `json:"profile_path,omitempty"`
	CreditID    string `json:"credit_id,omitempty"`
}

// Network represents a TV network.
type Network struct {
	ID            int    `json:"id"`
	Name          string `json:"name,omitempty"`
	LogoPath      string `json:"logo_path,omitempty"`
	OriginCountry string `json:"origin_country,omitempty"`
}

// SeasonSummary is the brief season info returned in TVDetails.
type SeasonSummary struct {
	ID           int     `json:"id"`
	AirDate      string  `json:"air_date,omitempty"`
	EpisodeCount int     `json:"episode_count,omitempty"`
	Name         string  `json:"name,omitempty"`
	Overview     string  `json:"overview,omitempty"`
	PosterPath   string  `json:"poster_path,omitempty"`
	SeasonNumber int     `json:"season_number"`
	VoteAverage  float64 `json:"vote_average,omitempty"`
}

// TVDetails contains full details for a TV show.
type TVDetails struct {
	ID                  int                 `json:"id"`
	Name                string              `json:"name,omitempty"`
	OriginalName        string              `json:"original_name,omitempty"`
	Overview            string              `json:"overview,omitempty"`
	Tagline             string              `json:"tagline,omitempty"`
	Status              string              `json:"status,omitempty"`
	Type                string              `json:"type,omitempty"`
	Homepage            string              `json:"homepage,omitempty"`
	PosterPath          string              `json:"poster_path,omitempty"`
	BackdropPath        string              `json:"backdrop_path,omitempty"`
	FirstAirDate        string              `json:"first_air_date,omitempty"`
	LastAirDate         string              `json:"last_air_date,omitempty"`
	OriginalLanguage    string              `json:"original_language,omitempty"`
	OriginCountry       []string            `json:"origin_country,omitempty"`
	InProduction        bool                `json:"in_production,omitempty"`
	NumberOfEpisodes    int                 `json:"number_of_episodes,omitempty"`
	NumberOfSeasons     int                 `json:"number_of_seasons,omitempty"`
	Popularity          float64             `json:"popularity,omitempty"`
	VoteAverage         float64             `json:"vote_average,omitempty"`
	VoteCount           int                 `json:"vote_count,omitempty"`
	Genres              []Genre             `json:"genres,omitempty"`
	CreatedBy           []CreatedBy         `json:"created_by,omitempty"`
	Networks            []Network           `json:"networks,omitempty"`
	Seasons             []SeasonSummary     `json:"seasons,omitempty"`
	ProductionCompanies []ProductionCompany `json:"production_companies,omitempty"`
	ProductionCountries []ProductionCountry `json:"production_countries,omitempty"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages,omitempty"`
	EpisodeRunTime      []int               `json:"episode_run_time,omitempty"`
	Languages           []string            `json:"languages,omitempty"`
}

// Episode represents a TV episode.
type Episode struct {
	ID             int     `json:"id"`
	Name           string  `json:"name,omitempty"`
	Overview       string  `json:"overview,omitempty"`
	AirDate        string  `json:"air_date,omitempty"`
	EpisodeNumber  int     `json:"episode_number"`
	SeasonNumber   int     `json:"season_number"`
	StillPath      string  `json:"still_path,omitempty"`
	VoteAverage    float64 `json:"vote_average,omitempty"`
	VoteCount      int     `json:"vote_count,omitempty"`
	ProductionCode string  `json:"production_code,omitempty"`
	Runtime        int     `json:"runtime,omitempty"`
	ShowID         int     `json:"show_id,omitempty"`
}

// SeasonDetails contains full details for a TV season.
type SeasonDetails struct {
	ID           int       `json:"id"`
	AirDate      string    `json:"air_date,omitempty"`
	Name         string    `json:"name,omitempty"`
	Overview     string    `json:"overview,omitempty"`
	PosterPath   string    `json:"poster_path,omitempty"`
	SeasonNumber int       `json:"season_number"`
	VoteAverage  float64   `json:"vote_average,omitempty"`
	Episodes     []Episode `json:"episodes,omitempty"`
}

// PersonDetails contains full details for a person.
type PersonDetails struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name,omitempty"`
	Biography          string   `json:"biography,omitempty"`
	Birthday           string   `json:"birthday,omitempty"`
	Deathday           string   `json:"deathday,omitempty"`
	PlaceOfBirth       string   `json:"place_of_birth,omitempty"`
	ProfilePath        string   `json:"profile_path,omitempty"`
	IMDbID             string   `json:"imdb_id,omitempty"`
	Homepage           string   `json:"homepage,omitempty"`
	KnownForDepartment string   `json:"known_for_department,omitempty"`
	AlsoKnownAs        []string `json:"also_known_as,omitempty"`
	Gender             int      `json:"gender,omitempty"`
	Adult              bool     `json:"adult,omitempty"`
	Popularity         float64  `json:"popularity,omitempty"`
}

// CastMember represents an actor in a credits response.
type CastMember struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name,omitempty"`
	Character          string  `json:"character,omitempty"`
	ProfilePath        string  `json:"profile_path,omitempty"`
	Order              int     `json:"order"`
	Gender             int     `json:"gender,omitempty"`
	Popularity         float64 `json:"popularity,omitempty"`
	CreditID           string  `json:"credit_id,omitempty"`
	KnownForDepartment string  `json:"known_for_department,omitempty"`
	OriginalName       string  `json:"original_name,omitempty"`
	Adult              bool    `json:"adult,omitempty"`
}

// CrewMember represents a crew member in a credits response.
type CrewMember struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name,omitempty"`
	Department         string  `json:"department,omitempty"`
	Job                string  `json:"job,omitempty"`
	ProfilePath        string  `json:"profile_path,omitempty"`
	Gender             int     `json:"gender,omitempty"`
	Popularity         float64 `json:"popularity,omitempty"`
	CreditID           string  `json:"credit_id,omitempty"`
	KnownForDepartment string  `json:"known_for_department,omitempty"`
	OriginalName       string  `json:"original_name,omitempty"`
	Adult              bool    `json:"adult,omitempty"`
}

// Credits contains the cast and crew for a movie or TV show.
type Credits struct {
	ID   int          `json:"id"`
	Cast []CastMember `json:"cast,omitempty"`
	Crew []CrewMember `json:"crew,omitempty"`
}

// ImageItem represents a single image.
type ImageItem struct {
	FilePath    string  `json:"file_path,omitempty"`
	Width       int     `json:"width,omitempty"`
	Height      int     `json:"height,omitempty"`
	AspectRatio float64 `json:"aspect_ratio,omitempty"`
	VoteAverage float64 `json:"vote_average,omitempty"`
	VoteCount   int     `json:"vote_count,omitempty"`
	ISO6391     string  `json:"iso_639_1,omitempty"`
}

// Images contains categorized images for a movie or TV show.
type Images struct {
	ID        int         `json:"id"`
	Backdrops []ImageItem `json:"backdrops,omitempty"`
	Posters   []ImageItem `json:"posters,omitempty"`
	Logos     []ImageItem `json:"logos,omitempty"`
}

// ExternalIDs contains cross-platform identifiers.
type ExternalIDs struct {
	ID          int    `json:"id"`
	IMDbID      string `json:"imdb_id,omitempty"`
	FacebookID  string `json:"facebook_id,omitempty"`
	InstagramID string `json:"instagram_id,omitempty"`
	TwitterID   string `json:"twitter_id,omitempty"`
	WikidataID  string `json:"wikidata_id,omitempty"`
	TVDbID      int    `json:"tvdb_id,omitempty"`
	TVRageID    int    `json:"tvrage_id,omitempty"`
}

// Configuration contains the API image base URL and sizes.
type Configuration struct {
	Images     ImageConfiguration `json:"images"`
	ChangeKeys []string           `json:"change_keys,omitempty"`
}

// ImageConfiguration contains image base URL and available sizes.
type ImageConfiguration struct {
	BaseURL       string   `json:"base_url,omitempty"`
	SecureBaseURL string   `json:"secure_base_url,omitempty"`
	BackdropSizes []string `json:"backdrop_sizes,omitempty"`
	LogoSizes     []string `json:"logo_sizes,omitempty"`
	PosterSizes   []string `json:"poster_sizes,omitempty"`
	ProfileSizes  []string `json:"profile_sizes,omitempty"`
	StillSizes    []string `json:"still_sizes,omitempty"`
}

// FindResult contains the results of a find-by-external-ID query.
type FindResult struct {
	MovieResults  []MovieResult  `json:"movie_results,omitempty"`
	TVResults     []TVResult     `json:"tv_results,omitempty"`
	PersonResults []PersonResult `json:"person_results,omitempty"`
}

// APIError is returned when the TMDb API responds with a non-2xx status.
type APIError struct {
	StatusCode    int    `json:"-"`
	StatusMessage string `json:"status_message,omitempty"`
	ErrorCode     int    `json:"status_code,omitempty"`
	// RawBody holds the raw response body when the error response could not be parsed as JSON.
	RawBody string `json:"-"`
}

func (e *APIError) Error() string {
	if e.StatusMessage != "" {
		return fmt.Sprintf("tmdb: HTTP %d: %s (code %d)", e.StatusCode, e.StatusMessage, e.ErrorCode)
	}
	if e.RawBody != "" {
		return fmt.Sprintf("tmdb: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("tmdb: HTTP %d", e.StatusCode)
}
