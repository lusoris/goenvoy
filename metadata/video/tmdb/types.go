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

// KnownForItem represents a media entry a person is known for.
type KnownForItem struct {
	MediaType        string  `json:"media_type"`
	ID               int     `json:"id"`
	Title            string  `json:"title,omitempty"`
	Name             string  `json:"name,omitempty"`
	OriginalTitle    string  `json:"original_title,omitempty"`
	OriginalName     string  `json:"original_name,omitempty"`
	Overview         string  `json:"overview,omitempty"`
	PosterPath       string  `json:"poster_path,omitempty"`
	BackdropPath     string  `json:"backdrop_path,omitempty"`
	ReleaseDate      string  `json:"release_date,omitempty"`
	FirstAirDate     string  `json:"first_air_date,omitempty"`
	VoteAverage      float64 `json:"vote_average,omitempty"`
	VoteCount        int     `json:"vote_count,omitempty"`
	GenreIDs         []int   `json:"genre_ids,omitempty"`
	OriginalLanguage string  `json:"original_language,omitempty"`
	Adult            bool    `json:"adult,omitempty"`
}

// PersonResult represents a person in search results.
type PersonResult struct {
	ID                 int            `json:"id"`
	Name               string         `json:"name,omitempty"`
	ProfilePath        string         `json:"profile_path,omitempty"`
	Adult              bool           `json:"adult,omitempty"`
	Popularity         float64        `json:"popularity,omitempty"`
	KnownForDepartment string         `json:"known_for_department,omitempty"`
	Gender             int            `json:"gender,omitempty"`
	KnownFor           []KnownForItem `json:"known_for,omitempty"`
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
	Stills    []ImageItem `json:"stills,omitempty"`
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
	TikTokID    string `json:"tiktok_id,omitempty"`
	YouTubeID   string `json:"youtube_id,omitempty"`
	FreebaseID  string `json:"freebase_id,omitempty"`
	FreebaseMID string `json:"freebase_mid,omitempty"`
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

// Video represents a video (trailer, teaser, featurette, etc.).
type Video struct {
	ID        string `json:"id,omitempty"`
	Key       string `json:"key,omitempty"`
	Name      string `json:"name,omitempty"`
	Site      string `json:"site,omitempty"`
	Size      int    `json:"size,omitempty"`
	Type      string `json:"type,omitempty"`
	Official  bool   `json:"official,omitempty"`
	ISO6391   string `json:"iso_639_1,omitempty"`
	ISO31661  string `json:"iso_3166_1,omitempty"`
	Published string `json:"published_at,omitempty"`
}

// VideosResponse contains the list of videos for a movie or TV show.
type VideosResponse struct {
	ID      int     `json:"id"`
	Results []Video `json:"results,omitempty"`
}

// Keyword represents a keyword tag.
type Keyword struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// KeywordsResponse contains the keywords for a movie or TV show.
type KeywordsResponse struct {
	ID       int       `json:"id"`
	Keywords []Keyword `json:"keywords,omitempty"`
	Results  []Keyword `json:"results,omitempty"`
}

// AuthorDetails contains the author details for a review.
type AuthorDetails struct {
	Name       string  `json:"name,omitempty"`
	Username   string  `json:"username,omitempty"`
	AvatarPath string  `json:"avatar_path,omitempty"`
	Rating     float64 `json:"rating,omitempty"`
}

// Review represents a user review.
type Review struct {
	ID            string        `json:"id,omitempty"`
	Author        string        `json:"author,omitempty"`
	AuthorDetails AuthorDetails `json:"author_details,omitempty"`
	Content       string        `json:"content,omitempty"`
	CreatedAt     string        `json:"created_at,omitempty"`
	UpdatedAt     string        `json:"updated_at,omitempty"`
	URL           string        `json:"url,omitempty"`
	MediaID       int           `json:"media_id,omitempty"`
	MediaTitle    string        `json:"media_title,omitempty"`
	MediaType     string        `json:"media_type,omitempty"`
}

// ReleaseDate represents a release date entry for a country.
type ReleaseDate struct {
	Certification string `json:"certification,omitempty"`
	ISO6391       string `json:"iso_639_1,omitempty"`
	Note          string `json:"note,omitempty"`
	ReleaseDate   string `json:"release_date,omitempty"`
	Type          int    `json:"type,omitempty"`
}

// ReleaseDateCountry groups release dates by country.
type ReleaseDateCountry struct {
	ISO31661     string        `json:"iso_3166_1,omitempty"`
	ReleaseDates []ReleaseDate `json:"release_dates,omitempty"`
}

// ReleaseDatesResponse contains release date information for a movie.
type ReleaseDatesResponse struct {
	ID      int                  `json:"id"`
	Results []ReleaseDateCountry `json:"results,omitempty"`
}

// WatchProvider represents a single watch provider.
type WatchProvider struct {
	DisplayPriority int    `json:"display_priority,omitempty"`
	LogoPath        string `json:"logo_path,omitempty"`
	ProviderID      int    `json:"provider_id"`
	ProviderName    string `json:"provider_name,omitempty"`
}

// WatchProviderCountry contains provider options for a specific country.
type WatchProviderCountry struct {
	Link     string          `json:"link,omitempty"`
	Flatrate []WatchProvider `json:"flatrate,omitempty"`
	Rent     []WatchProvider `json:"rent,omitempty"`
	Buy      []WatchProvider `json:"buy,omitempty"`
	Ads      []WatchProvider `json:"ads,omitempty"`
	Free     []WatchProvider `json:"free,omitempty"`
}

// WatchProvidersResponse contains watch provider information.
type WatchProvidersResponse struct {
	ID      int                             `json:"id"`
	Results map[string]WatchProviderCountry `json:"results,omitempty"`
}

// AlternativeTitle represents an alternative title.
type AlternativeTitle struct {
	ISO31661 string `json:"iso_3166_1,omitempty"`
	Title    string `json:"title,omitempty"`
	Type     string `json:"type,omitempty"`
}

// AlternativeTitlesResponse contains alternative titles.
type AlternativeTitlesResponse struct {
	ID      int                `json:"id"`
	Titles  []AlternativeTitle `json:"titles,omitempty"`
	Results []AlternativeTitle `json:"results,omitempty"`
}

// TranslationData contains translation text data.
type TranslationData struct {
	Homepage string `json:"homepage,omitempty"`
	Overview string `json:"overview,omitempty"`
	Runtime  int    `json:"runtime,omitempty"`
	Tagline  string `json:"tagline,omitempty"`
	Title    string `json:"title,omitempty"`
	Name     string `json:"name,omitempty"`
}

// Translation represents a translation entry.
type Translation struct {
	ISO31661    string          `json:"iso_3166_1,omitempty"`
	ISO6391     string          `json:"iso_639_1,omitempty"`
	Name        string          `json:"name,omitempty"`
	EnglishName string          `json:"english_name,omitempty"`
	Data        TranslationData `json:"data,omitempty"`
}

// TranslationsResponse contains translations.
type TranslationsResponse struct {
	ID           int           `json:"id"`
	Translations []Translation `json:"translations,omitempty"`
}

// ListSummary represents a list in context of movie lists.
type ListSummary struct {
	Description   string `json:"description,omitempty"`
	FavoriteCount int    `json:"favorite_count,omitempty"`
	ID            int    `json:"id"`
	ISO6391       string `json:"iso_639_1,omitempty"`
	ItemCount     int    `json:"item_count,omitempty"`
	ListType      string `json:"list_type,omitempty"`
	Name          string `json:"name,omitempty"`
	PosterPath    string `json:"poster_path,omitempty"`
}

// AccountStates contains rated/watchlist/favorite state for a media item.
type AccountStates struct {
	ID        int  `json:"id"`
	Favorite  bool `json:"favorite,omitempty"`
	Watchlist bool `json:"watchlist,omitempty"`
	Rated     any  `json:"rated,omitempty"`
}

// ContentRating represents a content rating entry.
type ContentRating struct {
	ISO31661    string   `json:"iso_3166_1,omitempty"`
	Rating      string   `json:"rating,omitempty"`
	Descriptors []string `json:"descriptors,omitempty"`
}

// ContentRatingsResponse contains content ratings for a TV show.
type ContentRatingsResponse struct {
	ID      int             `json:"id"`
	Results []ContentRating `json:"results,omitempty"`
}

// EpisodeGroup represents an episode group.
type EpisodeGroup struct {
	Description  string   `json:"description,omitempty"`
	EpisodeCount int      `json:"episode_count,omitempty"`
	GroupCount   int      `json:"group_count,omitempty"`
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name,omitempty"`
	Network      *Network `json:"network,omitempty"`
	Type         int      `json:"type,omitempty"`
}

// EpisodeGroupsResponse contains episode groups for a TV show.
type EpisodeGroupsResponse struct {
	ID      int            `json:"id"`
	Results []EpisodeGroup `json:"results,omitempty"`
}

// AggregateRole represents a role in aggregate credits.
type AggregateRole struct {
	CreditID     string `json:"credit_id,omitempty"`
	Character    string `json:"character,omitempty"`
	EpisodeCount int    `json:"episode_count,omitempty"`
}

// AggregateJob represents a job in aggregate credits.
type AggregateJob struct {
	CreditID     string `json:"credit_id,omitempty"`
	Job          string `json:"job,omitempty"`
	EpisodeCount int    `json:"episode_count,omitempty"`
}

// AggregateCastMember represents a cast member in aggregate credits.
type AggregateCastMember struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name,omitempty"`
	Adult              bool            `json:"adult,omitempty"`
	Gender             int             `json:"gender,omitempty"`
	KnownForDepartment string          `json:"known_for_department,omitempty"`
	OriginalName       string          `json:"original_name,omitempty"`
	Popularity         float64         `json:"popularity,omitempty"`
	ProfilePath        string          `json:"profile_path,omitempty"`
	Roles              []AggregateRole `json:"roles,omitempty"`
	TotalEpisodeCount  int             `json:"total_episode_count,omitempty"`
	Order              int             `json:"order"`
}

// AggregateCrewMember represents a crew member in aggregate credits.
type AggregateCrewMember struct {
	ID                 int            `json:"id"`
	Name               string         `json:"name,omitempty"`
	Adult              bool           `json:"adult,omitempty"`
	Gender             int            `json:"gender,omitempty"`
	KnownForDepartment string         `json:"known_for_department,omitempty"`
	OriginalName       string         `json:"original_name,omitempty"`
	Popularity         float64        `json:"popularity,omitempty"`
	ProfilePath        string         `json:"profile_path,omitempty"`
	Department         string         `json:"department,omitempty"`
	Jobs               []AggregateJob `json:"jobs,omitempty"`
	TotalEpisodeCount  int            `json:"total_episode_count,omitempty"`
}

// AggregateCredits contains aggregate cast and crew for a TV show.
type AggregateCredits struct {
	ID   int                   `json:"id"`
	Cast []AggregateCastMember `json:"cast,omitempty"`
	Crew []AggregateCrewMember `json:"crew,omitempty"`
}

// EpisodeDetails contains full details for a TV episode.
type EpisodeDetails struct {
	ID             int          `json:"id"`
	AirDate        string       `json:"air_date,omitempty"`
	Name           string       `json:"name,omitempty"`
	Overview       string       `json:"overview,omitempty"`
	ProductionCode string       `json:"production_code,omitempty"`
	Runtime        int          `json:"runtime,omitempty"`
	SeasonNumber   int          `json:"season_number"`
	EpisodeNumber  int          `json:"episode_number"`
	StillPath      string       `json:"still_path,omitempty"`
	VoteAverage    float64      `json:"vote_average,omitempty"`
	VoteCount      int          `json:"vote_count,omitempty"`
	ShowID         int          `json:"show_id,omitempty"`
	Crew           []CrewMember `json:"crew,omitempty"`
	GuestStars     []CastMember `json:"guest_stars,omitempty"`
}

// PersonCastCredit represents a cast credit for a person.
type PersonCastCredit struct {
	ID            int     `json:"id"`
	Title         string  `json:"title,omitempty"`
	Name          string  `json:"name,omitempty"`
	OriginalTitle string  `json:"original_title,omitempty"`
	OriginalName  string  `json:"original_name,omitempty"`
	Character     string  `json:"character,omitempty"`
	Overview      string  `json:"overview,omitempty"`
	PosterPath    string  `json:"poster_path,omitempty"`
	BackdropPath  string  `json:"backdrop_path,omitempty"`
	MediaType     string  `json:"media_type,omitempty"`
	ReleaseDate   string  `json:"release_date,omitempty"`
	FirstAirDate  string  `json:"first_air_date,omitempty"`
	VoteAverage   float64 `json:"vote_average,omitempty"`
	VoteCount     int     `json:"vote_count,omitempty"`
	Popularity    float64 `json:"popularity,omitempty"`
	GenreIDs      []int   `json:"genre_ids,omitempty"`
	CreditID      string  `json:"credit_id,omitempty"`
	Adult         bool    `json:"adult,omitempty"`
	EpisodeCount  int     `json:"episode_count,omitempty"`
	Order         int     `json:"order,omitempty"`
}

// PersonCrewCredit represents a crew credit for a person.
type PersonCrewCredit struct {
	ID            int     `json:"id"`
	Title         string  `json:"title,omitempty"`
	Name          string  `json:"name,omitempty"`
	OriginalTitle string  `json:"original_title,omitempty"`
	OriginalName  string  `json:"original_name,omitempty"`
	Department    string  `json:"department,omitempty"`
	Job           string  `json:"job,omitempty"`
	Overview      string  `json:"overview,omitempty"`
	PosterPath    string  `json:"poster_path,omitempty"`
	BackdropPath  string  `json:"backdrop_path,omitempty"`
	MediaType     string  `json:"media_type,omitempty"`
	ReleaseDate   string  `json:"release_date,omitempty"`
	FirstAirDate  string  `json:"first_air_date,omitempty"`
	VoteAverage   float64 `json:"vote_average,omitempty"`
	VoteCount     int     `json:"vote_count,omitempty"`
	Popularity    float64 `json:"popularity,omitempty"`
	GenreIDs      []int   `json:"genre_ids,omitempty"`
	CreditID      string  `json:"credit_id,omitempty"`
	Adult         bool    `json:"adult,omitempty"`
	EpisodeCount  int     `json:"episode_count,omitempty"`
}

// PersonCredits contains cast and crew credits for a person.
type PersonCredits struct {
	ID   int                `json:"id"`
	Cast []PersonCastCredit `json:"cast,omitempty"`
	Crew []PersonCrewCredit `json:"crew,omitempty"`
}

// PersonImages contains profile images for a person.
type PersonImages struct {
	ID       int         `json:"id"`
	Profiles []ImageItem `json:"profiles,omitempty"`
}

// TaggedImage represents a tagged image for a person.
type TaggedImage struct {
	ImageItem
	ID        string `json:"id,omitempty"`
	MediaType string `json:"media_type,omitempty"`
}

// MovieDetailsFull embeds MovieDetails with optional appended sub-resources.
type MovieDetailsFull struct {
	MovieDetails
	Credits      *Credits              `json:"credits,omitempty"`
	Images       *Images               `json:"images,omitempty"`
	ReleaseDates *ReleaseDatesResponse `json:"release_dates,omitempty"`
	Translations *TranslationsResponse `json:"translations,omitempty"`
	ExternalIDs  *ExternalIDs          `json:"external_ids,omitempty"`
	Videos       *VideosResponse       `json:"videos,omitempty"`
}

// TVDetailsFull embeds TVDetails with optional appended sub-resources.
type TVDetailsFull struct {
	TVDetails
	Credits        *Credits                `json:"credits,omitempty"`
	Images         *Images                 `json:"images,omitempty"`
	ContentRatings *ContentRatingsResponse `json:"content_ratings,omitempty"`
	Translations   *TranslationsResponse   `json:"translations,omitempty"`
	ExternalIDs    *ExternalIDs            `json:"external_ids,omitempty"`
	Videos         *VideosResponse         `json:"videos,omitempty"`
}

// SeasonDetailsFull embeds SeasonDetails with optional appended sub-resources.
type SeasonDetailsFull struct {
	SeasonDetails
	Credits      *Credits              `json:"credits,omitempty"`
	Images       *Images               `json:"images,omitempty"`
	Translations *TranslationsResponse `json:"translations,omitempty"`
	ExternalIDs  *ExternalIDs          `json:"external_ids,omitempty"`
	Videos       *VideosResponse       `json:"videos,omitempty"`
}

// EpisodeDetailsFull embeds EpisodeDetails with optional appended sub-resources.
type EpisodeDetailsFull struct {
	EpisodeDetails
	Credits      *Credits              `json:"credits,omitempty"`
	Images       *Images               `json:"images,omitempty"`
	Translations *TranslationsResponse `json:"translations,omitempty"`
	ExternalIDs  *ExternalIDs          `json:"external_ids,omitempty"`
	Videos       *VideosResponse       `json:"videos,omitempty"`
}

// PersonDetailsFull embeds PersonDetails with optional appended sub-resources.
type PersonDetailsFull struct {
	PersonDetails
	ExternalIDs *ExternalIDs  `json:"external_ids,omitempty"`
	Images      *PersonImages `json:"images,omitempty"`
}

// CollectionResult represents a collection in search results.
type CollectionResult struct {
	ID               int    `json:"id"`
	Name             string `json:"name,omitempty"`
	PosterPath       string `json:"poster_path,omitempty"`
	BackdropPath     string `json:"backdrop_path,omitempty"`
	Adult            bool   `json:"adult,omitempty"`
	OriginalLanguage string `json:"original_language,omitempty"`
	OriginalName     string `json:"original_name,omitempty"`
	Overview         string `json:"overview,omitempty"`
}

// CompanyResult represents a company in search results.
type CompanyResult struct {
	ID            int    `json:"id"`
	Name          string `json:"name,omitempty"`
	LogoPath      string `json:"logo_path,omitempty"`
	OriginCountry string `json:"origin_country,omitempty"`
}

// CollectionPart represents a movie part of a collection.
type CollectionPart struct {
	ID               int     `json:"id"`
	Title            string  `json:"title,omitempty"`
	OriginalTitle    string  `json:"original_title,omitempty"`
	Overview         string  `json:"overview,omitempty"`
	PosterPath       string  `json:"poster_path,omitempty"`
	BackdropPath     string  `json:"backdrop_path,omitempty"`
	ReleaseDate      string  `json:"release_date,omitempty"`
	OriginalLanguage string  `json:"original_language,omitempty"`
	Adult            bool    `json:"adult,omitempty"`
	Video            bool    `json:"video,omitempty"`
	GenreIDs         []int   `json:"genre_ids,omitempty"`
	Popularity       float64 `json:"popularity,omitempty"`
	VoteAverage      float64 `json:"vote_average,omitempty"`
	VoteCount        int     `json:"vote_count,omitempty"`
	MediaType        string  `json:"media_type,omitempty"`
}

// CollectionDetails contains full details for a movie collection.
type CollectionDetails struct {
	ID           int              `json:"id"`
	Name         string           `json:"name,omitempty"`
	Overview     string           `json:"overview,omitempty"`
	PosterPath   string           `json:"poster_path,omitempty"`
	BackdropPath string           `json:"backdrop_path,omitempty"`
	Parts        []CollectionPart `json:"parts,omitempty"`
}

// AccountDetails contains the account information for the authenticated user.
type AccountDetails struct {
	ID           int     `json:"id"`
	Name         string  `json:"name,omitempty"`
	Username     string  `json:"username,omitempty"`
	ISO6391      string  `json:"iso_639_1,omitempty"`
	ISO31661     string  `json:"iso_3166_1,omitempty"`
	IncludeAdult bool    `json:"include_adult,omitempty"`
	Avatar       *Avatar `json:"avatar,omitempty"`
}

// Avatar contains avatar information.
type Avatar struct {
	Gravatar *Gravatar   `json:"gravatar,omitempty"`
	TMDb     *UserAvatar `json:"tmdb,omitempty"`
}

// Gravatar contains a gravatar hash.
type Gravatar struct {
	Hash string `json:"hash,omitempty"`
}

// UserAvatar contains a TMDb avatar path.
type UserAvatar struct {
	AvatarPath string `json:"avatar_path,omitempty"`
}

// RatedMovie represents a rated movie in the user's rated list.
type RatedMovie struct {
	MovieResult
	Rating float64 `json:"rating,omitempty"`
}

// RatedTV represents a rated TV show in the user's rated list.
type RatedTV struct {
	TVResult
	Rating float64 `json:"rating,omitempty"`
}

// RatedEpisode represents a rated TV episode.
type RatedEpisode struct {
	ID            int     `json:"id"`
	Name          string  `json:"name,omitempty"`
	Overview      string  `json:"overview,omitempty"`
	AirDate       string  `json:"air_date,omitempty"`
	EpisodeNumber int     `json:"episode_number"`
	SeasonNumber  int     `json:"season_number"`
	StillPath     string  `json:"still_path,omitempty"`
	VoteAverage   float64 `json:"vote_average,omitempty"`
	VoteCount     int     `json:"vote_count,omitempty"`
	ShowID        int     `json:"show_id,omitempty"`
	Rating        float64 `json:"rating,omitempty"`
}

// ListDetails contains the full details for a list.
type ListDetails struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	PosterPath    string `json:"poster_path,omitempty"`
	ISO6391       string `json:"iso_639_1,omitempty"`
	CreatedBy     string `json:"created_by,omitempty"`
	FavoriteCount int    `json:"favorite_count,omitempty"`
	ItemCount     int    `json:"item_count,omitempty"`
	Items         []any  `json:"items,omitempty"`
}

// CreateListResponse is the response to creating a list.
type CreateListResponse struct {
	StatusMessage string `json:"status_message,omitempty"`
	Success       bool   `json:"success,omitempty"`
	StatusCode    int    `json:"status_code,omitempty"`
	ListID        int    `json:"list_id,omitempty"`
}

// ItemStatus represents the status of a movie on a list.
type ItemStatus struct {
	ID          string `json:"id,omitempty"`
	ItemPresent bool   `json:"item_present,omitempty"`
}

// Certification represents a certification entry.
type Certification struct {
	Certification string `json:"certification,omitempty"`
	Meaning       string `json:"meaning,omitempty"`
	Order         int    `json:"order,omitempty"`
}

// CertificationsResponse contains certifications by country.
type CertificationsResponse struct {
	Certifications map[string][]Certification `json:"certifications,omitempty"`
}

// WatchProviderRegion represents a region where watch providers are available.
type WatchProviderRegion struct {
	ISO31661    string `json:"iso_3166_1,omitempty"`
	EnglishName string `json:"english_name,omitempty"`
	NativeName  string `json:"native_name,omitempty"`
}

// WatchProviderRegionsResponse contains available watch provider regions.
type WatchProviderRegionsResponse struct {
	Results []WatchProviderRegion `json:"results,omitempty"`
}

// WatchProviderListItem represents a watch provider in the provider list.
type WatchProviderListItem struct {
	DisplayPriorities map[string]int `json:"display_priorities,omitempty"`
	DisplayPriority   int            `json:"display_priority,omitempty"`
	LogoPath          string         `json:"logo_path,omitempty"`
	ProviderID        int            `json:"provider_id"`
	ProviderName      string         `json:"provider_name,omitempty"`
}

// WatchProviderListResponse contains a list of available watch providers.
type WatchProviderListResponse struct {
	Results []WatchProviderListItem `json:"results,omitempty"`
}

// CompanyDetails contains full details for a production company.
type CompanyDetails struct {
	ID            int                `json:"id"`
	Name          string             `json:"name,omitempty"`
	Description   string             `json:"description,omitempty"`
	Headquarters  string             `json:"headquarters,omitempty"`
	Homepage      string             `json:"homepage,omitempty"`
	LogoPath      string             `json:"logo_path,omitempty"`
	OriginCountry string             `json:"origin_country,omitempty"`
	ParentCompany *ProductionCompany `json:"parent_company,omitempty"`
}

// ChangeItem represents a changed item in the changes endpoint.
type ChangeItem struct {
	ID    int  `json:"id"`
	Adult bool `json:"adult,omitempty"`
}

// ChangesResponse contains a list of changed items.
type ChangesResponse struct {
	Results      []ChangeItem `json:"results,omitempty"`
	Page         int          `json:"page"`
	TotalPages   int          `json:"total_pages"`
	TotalResults int          `json:"total_results"`
}

// Language represents a language from the configuration endpoint.
type Language struct {
	ISO6391     string `json:"iso_639_1,omitempty"`
	EnglishName string `json:"english_name,omitempty"`
	Name        string `json:"name,omitempty"`
}

// Country represents a country from the configuration endpoint.
type Country struct {
	ISO31661    string `json:"iso_3166_1,omitempty"`
	EnglishName string `json:"english_name,omitempty"`
	NativeName  string `json:"native_name,omitempty"`
}

// Timezone represents a timezone from the configuration endpoint.
type Timezone struct {
	ISO31661 string   `json:"iso_3166_1,omitempty"`
	Zones    []string `json:"zones,omitempty"`
}

// Department represents a department and its jobs.
type Department struct {
	Department string   `json:"department,omitempty"`
	Jobs       []string `json:"jobs,omitempty"`
}
