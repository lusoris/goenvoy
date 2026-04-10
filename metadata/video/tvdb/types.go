package tvdb

import (
	"net/url"
	"strconv"
)

// response is the standard TheTVDB v4 API envelope.
type response[T any] struct {
	Data   T      `json:"data"`
	Status string `json:"status"`
}

// Links contains pagination links returned by list endpoints.
type Links struct {
	Prev       *string `json:"prev"`
	Self       string  `json:"self"`
	Next       *string `json:"next"`
	TotalItems int     `json:"total_items"`
	PageSize   int     `json:"page_size"`
}

// LoginRequest is the body for POST /login.
type LoginRequest struct {
	APIKey string `json:"apikey"`
	PIN    string `json:"pin,omitempty"`
}

// Alias represents an alternative name for an entity.
type Alias struct {
	Language string `json:"language"`
	Name     string `json:"name"`
}

// Status represents a series or movie status.
type Status struct {
	ID         *int64 `json:"id"`
	KeepUpdate bool   `json:"keepUpdated"`
	Name       string `json:"name"`
	RecordType string `json:"recordType"`
}

// Genre is a genre record.
type Genre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Language is a language record.
type Language struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	NativeName string `json:"nativeName"`
	ShortCode  string `json:"shortCode"`
}

// RemoteID is a cross-reference to an external service (e.g. IMDB).
type RemoteID struct {
	ID         string `json:"id"`
	Type       int64  `json:"type"`
	SourceName string `json:"sourceName"`
}

// TagOption is a tag option record.
type TagOption struct {
	HelpText string `json:"helpText"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Tag      int64  `json:"tag"`
	TagName  string `json:"tagName"`
}

// Trailer is a trailer record.
type Trailer struct {
	ID       int64  `json:"id"`
	Language string `json:"language"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Runtime  int    `json:"runtime"`
}

// ContentRating is a content rating record.
type ContentRating struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Country     string `json:"country"`
	ContentType string `json:"contentType"`
	Order       int    `json:"order"`
	FullName    string `json:"fullName"`
}

// Translation is a translation record.
type Translation struct {
	Aliases   []string `json:"aliases"`
	IsAlias   bool     `json:"isAlias"`
	IsPrimary bool     `json:"isPrimary"`
	Language  string   `json:"language"`
	Name      string   `json:"name"`
	Overview  string   `json:"overview"`
	Tagline   string   `json:"tagline,omitempty"`
}

// Company is a company record (network, studio, etc.).
type Company struct {
	ID                 int64       `json:"id"`
	Name               string      `json:"name"`
	Slug               string      `json:"slug"`
	Country            string      `json:"country"`
	ActiveDate         string      `json:"activeDate"`
	InactiveDate       string      `json:"inactiveDate"`
	PrimaryCompanyType *int64      `json:"primaryCompanyType"`
	Aliases            []Alias     `json:"aliases,omitempty"`
	TagOptions         []TagOption `json:"tagOptions,omitempty"`
}

// Character is a character record.
type Character struct {
	ID                   int64       `json:"id"`
	Name                 string      `json:"name"`
	PeopleID             int         `json:"peopleId"`
	SeriesID             *int        `json:"seriesId"`
	MovieID              *int        `json:"movieId"`
	EpisodeID            *int        `json:"episodeId"`
	Type                 int64       `json:"type"`
	Image                string      `json:"image"`
	Sort                 int64       `json:"sort"`
	IsFeatured           bool        `json:"isFeatured"`
	URL                  string      `json:"url"`
	PersonName           string      `json:"personName"`
	PeopleType           string      `json:"peopleType"`
	PersonImgURL         string      `json:"personImgURL"`
	NameTranslations     []string    `json:"nameTranslations"`
	OverviewTranslations []string    `json:"overviewTranslations"`
	Aliases              []Alias     `json:"aliases,omitempty"`
	TagOptions           []TagOption `json:"tagOptions,omitempty"`
}

// ArtworkBase is a base artwork record.
type ArtworkBase struct {
	ID           int64   `json:"id"`
	Image        string  `json:"image"`
	Thumbnail    string  `json:"thumbnail"`
	Language     string  `json:"language"`
	Type         int64   `json:"type"`
	Score        float64 `json:"score"`
	Width        int64   `json:"width"`
	Height       int64   `json:"height"`
	IncludesText bool    `json:"includesText"`
}

// ArtworkExtended is an extended artwork record.
type ArtworkExtended struct {
	ArtworkBase
	ThumbnailWidth  int64          `json:"thumbnailWidth"`
	ThumbnailHeight int64          `json:"thumbnailHeight"`
	UpdatedAt       int64          `json:"updatedAt"`
	EpisodeID       int            `json:"episodeId"`
	MovieID         int            `json:"movieId"`
	SeriesID        int            `json:"seriesId"`
	SeasonID        int            `json:"seasonId"`
	NetworkID       int            `json:"networkId"`
	PeopleID        int            `json:"peopleId"`
	Status          *ArtworkStatus `json:"status"`
	TagOptions      []TagOption    `json:"tagOptions,omitempty"`
}

// ArtworkStatus is an artwork status record.
type ArtworkStatus struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ArtworkType is an artwork type record.
type ArtworkType struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	RecordType  string `json:"recordType"`
	Slug        string `json:"slug"`
	ImageFormat string `json:"imageFormat"`
	Width       int64  `json:"width"`
	Height      int64  `json:"height"`
	ThumbWidth  int64  `json:"thumbWidth"`
	ThumbHeight int64  `json:"thumbHeight"`
}

// SeasonType is a season type record.
type SeasonType struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	AlternateName string `json:"alternateName"`
}

// Companies groups companies by their role.
type Companies struct {
	Studio         []Company `json:"studio,omitempty"`
	Network        []Company `json:"network,omitempty"`
	Production     []Company `json:"production,omitempty"`
	Distributor    []Company `json:"distributor,omitempty"`
	SpecialEffects []Company `json:"special_effects,omitempty"`
}

// SeriesAirsDays indicates which days a series airs.
type SeriesAirsDays struct {
	Monday    bool `json:"monday"`
	Tuesday   bool `json:"tuesday"`
	Wednesday bool `json:"wednesday"`
	Thursday  bool `json:"thursday"`
	Friday    bool `json:"friday"`
	Saturday  bool `json:"saturday"`
	Sunday    bool `json:"sunday"`
}

// SeriesBase is the base record for a series.
type SeriesBase struct {
	ID                   int      `json:"id"`
	Name                 string   `json:"name"`
	Slug                 string   `json:"slug"`
	Image                string   `json:"image"`
	FirstAired           string   `json:"firstAired"`
	LastAired            string   `json:"lastAired"`
	NextAired            string   `json:"nextAired"`
	Score                float64  `json:"score"`
	Status               *Status  `json:"status"`
	OriginalCountry      string   `json:"originalCountry"`
	OriginalLanguage     string   `json:"originalLanguage"`
	Country              string   `json:"country"`
	DefaultSeasonType    int64    `json:"defaultSeasonType"`
	IsOrderRandomized    bool     `json:"isOrderRandomized"`
	LastUpdated          string   `json:"lastUpdated"`
	AverageRuntime       *int     `json:"averageRuntime"`
	NameTranslations     []string `json:"nameTranslations"`
	OverviewTranslations []string `json:"overviewTranslations"`
	Aliases              []Alias  `json:"aliases,omitempty"`
	Year                 string   `json:"year"`
}

// SeriesExtended is the extended record for a series.
type SeriesExtended struct {
	SeriesBase
	Abbreviation    string            `json:"abbreviation"`
	AirsDays        *SeriesAirsDays   `json:"airsDays"`
	AirsTime        string            `json:"airsTime"`
	Artworks        []ArtworkExtended `json:"artworks,omitempty"`
	Characters      []Character       `json:"characters,omitempty"`
	Companies       []Company         `json:"companies,omitempty"`
	ContentRatings  []ContentRating   `json:"contentRatings,omitempty"`
	Episodes        []EpisodeBase     `json:"episodes,omitempty"`
	Genres          []Genre           `json:"genres,omitempty"`
	Lists           []ListBase        `json:"lists,omitempty"`
	Overview        string            `json:"overview"`
	OriginalNetwork *Company          `json:"originalNetwork"`
	LatestNetwork   *Company          `json:"latestNetwork"`
	RemoteIDs       []RemoteID        `json:"remoteIds,omitempty"`
	Seasons         []SeasonBase      `json:"seasons,omitempty"`
	SeasonTypes     []SeasonType      `json:"seasonTypes,omitempty"`
	Tags            []TagOption       `json:"tags,omitempty"`
	Trailers        []Trailer         `json:"trailers,omitempty"`
}

// MovieBase is the base record for a movie.
type MovieBase struct {
	ID                   int      `json:"id"`
	Name                 string   `json:"name"`
	Slug                 string   `json:"slug"`
	Image                string   `json:"image"`
	Score                float64  `json:"score"`
	Status               *Status  `json:"status"`
	Runtime              *int     `json:"runtime"`
	LastUpdated          string   `json:"lastUpdated"`
	Year                 string   `json:"year"`
	NameTranslations     []string `json:"nameTranslations"`
	OverviewTranslations []string `json:"overviewTranslations"`
	Aliases              []Alias  `json:"aliases,omitempty"`
}

// Release is a movie release record.
type Release struct {
	Country string `json:"country"`
	Date    string `json:"date"`
	Detail  string `json:"detail"`
}

// Inspiration is a movie inspiration record.
type Inspiration struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	TypeName string `json:"type_name"`
	URL      string `json:"url"`
}

// MovieExtended is the extended record for a movie.
type MovieExtended struct {
	MovieBase
	Artworks            []ArtworkBase       `json:"artworks,omitempty"`
	AudioLanguages      []string            `json:"audioLanguages,omitempty"`
	Awards              []AwardBase         `json:"awards,omitempty"`
	BoxOffice           string              `json:"boxOffice"`
	BoxOfficeUS         string              `json:"boxOfficeUS"`
	Budget              string              `json:"budget"`
	Characters          []Character         `json:"characters,omitempty"`
	Companies           *Companies          `json:"companies"`
	ContentRatings      []ContentRating     `json:"contentRatings,omitempty"`
	FirstRelease        *Release            `json:"first_release"`
	Genres              []Genre             `json:"genres,omitempty"`
	Inspirations        []Inspiration       `json:"inspirations,omitempty"`
	Lists               []ListBase          `json:"lists,omitempty"`
	OriginalCountry     string              `json:"originalCountry"`
	OriginalLanguage    string              `json:"originalLanguage"`
	ProductionCountries []ProductionCountry `json:"production_countries,omitempty"`
	Releases            []Release           `json:"releases,omitempty"`
	RemoteIDs           []RemoteID          `json:"remoteIds,omitempty"`
	SpokenLanguages     []string            `json:"spoken_languages,omitempty"`
	Studios             []StudioBase        `json:"studios,omitempty"`
	SubtitleLanguages   []string            `json:"subtitleLanguages,omitempty"`
	TagOptions          []TagOption         `json:"tagOptions,omitempty"`
	Trailers            []Trailer           `json:"trailers,omitempty"`
}

// StudioBase is a studio record.
type StudioBase struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ParentStudio int    `json:"parentStudio"`
}

// ProductionCountry is a production country record.
type ProductionCountry struct {
	ID      int64  `json:"id"`
	Country string `json:"country"`
	Name    string `json:"name"`
}

// AwardBase is a base award record.
type AwardBase struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ListBase is a list record.
type ListBase struct {
	ID                   int64       `json:"id"`
	Name                 string      `json:"name"`
	Overview             string      `json:"overview"`
	URL                  string      `json:"url"`
	IsOfficial           bool        `json:"isOfficial"`
	NameTranslations     []string    `json:"nameTranslations"`
	OverviewTranslations []string    `json:"overviewTranslations"`
	Aliases              []Alias     `json:"aliases,omitempty"`
	Score                int         `json:"score"`
	Image                string      `json:"image"`
	ImageIsFallback      bool        `json:"imageIsFallback"`
	RemoteIDs            []RemoteID  `json:"remoteIds,omitempty"`
	Tags                 []TagOption `json:"tags,omitempty"`
}

// EpisodeBase is the base record for an episode.
type EpisodeBase struct {
	ID                   int      `json:"id"`
	SeriesID             int64    `json:"seriesId"`
	Name                 string   `json:"name"`
	Aired                string   `json:"aired"`
	Runtime              *int     `json:"runtime"`
	Overview             string   `json:"overview"`
	Image                string   `json:"image"`
	ImageType            *int     `json:"imageType"`
	IsMovie              int64    `json:"isMovie"`
	Number               int      `json:"number"`
	SeasonNumber         int      `json:"seasonNumber"`
	AbsoluteNumber       int      `json:"absoluteNumber"`
	AirsAfterSeason      int      `json:"airsAfterSeason"`
	AirsBeforeEpisode    int      `json:"airsBeforeEpisode"`
	AirsBeforeSeason     int      `json:"airsBeforeSeason"`
	FinaleType           string   `json:"finaleType"`
	LastUpdated          string   `json:"lastUpdated"`
	Year                 string   `json:"year"`
	NameTranslations     []string `json:"nameTranslations"`
	OverviewTranslations []string `json:"overviewTranslations"`
	SeasonName           string   `json:"seasonName"`
	LinkedMovie          int      `json:"linkedMovie"`
}

// EpisodeExtended is the extended record for an episode.
type EpisodeExtended struct {
	EpisodeBase
	Awards         []AwardBase     `json:"awards,omitempty"`
	Characters     []Character     `json:"characters,omitempty"`
	Companies      []Company       `json:"companies,omitempty"`
	ContentRatings []ContentRating `json:"contentRatings,omitempty"`
	Networks       []Company       `json:"networks,omitempty"`
	Studios        []Company       `json:"studios,omitempty"`
	Nominations    []interface{}   `json:"nominations,omitempty"`
	ProductionCode string          `json:"productionCode"`
	RemoteIDs      []RemoteID      `json:"remoteIds,omitempty"`
	TagOptions     []TagOption     `json:"tagOptions,omitempty"`
	Trailers       []Trailer       `json:"trailers,omitempty"`
}

// SeasonBase is the base record for a season.
type SeasonBase struct {
	ID                   int         `json:"id"`
	SeriesID             int64       `json:"seriesId"`
	Name                 string      `json:"name"`
	Number               int64       `json:"number"`
	Image                string      `json:"image"`
	ImageType            int         `json:"imageType"`
	LastUpdated          string      `json:"lastUpdated"`
	Year                 string      `json:"year"`
	Type                 *SeasonType `json:"type"`
	Companies            *Companies  `json:"companies"`
	NameTranslations     []string    `json:"nameTranslations"`
	OverviewTranslations []string    `json:"overviewTranslations"`
}

// SeasonExtended is the extended record for a season.
type SeasonExtended struct {
	SeasonBase
	Artwork      []ArtworkBase `json:"artwork,omitempty"`
	Episodes     []EpisodeBase `json:"episodes,omitempty"`
	Trailers     []Trailer     `json:"trailers,omitempty"`
	TagOptions   []TagOption   `json:"tagOptions,omitempty"`
	Translations []Translation `json:"translations,omitempty"`
}

// PersonBase is the base record for a person.
type PersonBase struct {
	ID                   int      `json:"id"`
	Name                 string   `json:"name"`
	Image                string   `json:"image"`
	Score                int64    `json:"score"`
	LastUpdated          string   `json:"lastUpdated"`
	NameTranslations     []string `json:"nameTranslations"`
	OverviewTranslations []string `json:"overviewTranslations"`
	Aliases              []Alias  `json:"aliases,omitempty"`
}

// Biography is a person biography record.
type Biography struct {
	Biography string `json:"biography"`
	Language  string `json:"language"`
}

// PersonExtended is the extended record for a person.
type PersonExtended struct {
	PersonBase
	Birth       string      `json:"birth"`
	BirthPlace  string      `json:"birthPlace"`
	Death       string      `json:"death"`
	Gender      int         `json:"gender"`
	Slug        string      `json:"slug"`
	Awards      []AwardBase `json:"awards,omitempty"`
	Biographies []Biography `json:"biographies,omitempty"`
	Characters  []Character `json:"characters,omitempty"`
	RemoteIDs   []RemoteID  `json:"remoteIds,omitempty"`
	TagOptions  []TagOption `json:"tagOptions,omitempty"`
}

// SearchResult is the result from the search endpoint.
type SearchResult struct {
	Aliases         []string   `json:"aliases,omitempty"`
	Companies       []string   `json:"companies,omitempty"`
	CompanyType     string     `json:"companyType"`
	Country         string     `json:"country"`
	Director        string     `json:"director"`
	FirstAirTime    string     `json:"first_air_time"`
	Genres          []string   `json:"genres,omitempty"`
	ID              string     `json:"id"`
	ImageURL        string     `json:"image_url"`
	Name            string     `json:"name"`
	IsOfficial      bool       `json:"is_official"`
	Network         string     `json:"network"`
	ObjectID        string     `json:"objectID"`
	OfficialList    string     `json:"officialList"`
	Overview        string     `json:"overview"`
	Poster          string     `json:"poster"`
	Posters         []string   `json:"posters,omitempty"`
	PrimaryLanguage string     `json:"primary_language"`
	RemoteIDs       []RemoteID `json:"remote_ids,omitempty"`
	Status          string     `json:"status"`
	Slug            string     `json:"slug"`
	Studios         []string   `json:"studios,omitempty"`
	Title           string     `json:"title"`
	Thumbnail       string     `json:"thumbnail"`
	TVDBID          string     `json:"tvdb_id"`
	Type            string     `json:"type"`
	Year            string     `json:"year"`
}

// SearchByRemoteIDResult contains entity records matched by a remote ID.
type SearchByRemoteIDResult struct {
	Series  *SeriesBase  `json:"series,omitempty"`
	People  *PersonBase  `json:"people,omitempty"`
	Movie   *MovieBase   `json:"movie,omitempty"`
	Episode *EpisodeBase `json:"episode,omitempty"`
}

// EntityUpdate is a record update entry from the /updates endpoint.
type EntityUpdate struct {
	EntityType        string `json:"entityType"`
	Method            string `json:"method"`
	MethodInt         int    `json:"methodInt"`
	ExtraInfo         string `json:"extraInfo"`
	UserID            int    `json:"userId"`
	RecordType        string `json:"recordType"`
	RecordID          int64  `json:"recordId"`
	TimeStamp         int64  `json:"timeStamp"`
	SeriesID          int64  `json:"seriesId"`
	MergeToID         int64  `json:"mergeToId"`
	MergeToEntityType string `json:"mergeToEntityType"`
}

// SeriesEpisodesResult is the combined response from the series episodes endpoint.
type SeriesEpisodesResult struct {
	Series   *SeriesBase   `json:"series"`
	Episodes []EpisodeBase `json:"episodes"`
}

// SearchParams holds optional query parameters for the search endpoint.
type SearchParams struct {
	Type     string
	Year     int
	Company  string
	Country  string
	Director string
	Language string
	Network  string
	RemoteID string
	Offset   int
	Limit    int
}

// UpdatesParams holds optional query parameters for the updates endpoint.
type UpdatesParams struct {
	Type   string
	Action string
	Page   int
}

// APIError represents an error response from TheTVDB API.
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	// RawBody holds the raw response body when the error response could not be parsed as JSON.
	RawBody string `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return "tvdb: " + e.Message
	}
	if e.RawBody != "" {
		return "tvdb: HTTP " + httpStatusText(e.StatusCode) + ": " + e.RawBody
	}
	return "tvdb: HTTP " + httpStatusText(e.StatusCode)
}

func httpStatusText(code int) string {
	switch code {
	case 400:
		return "400 Bad Request"
	case 401:
		return "401 Unauthorized"
	case 404:
		return "404 Not Found"
	default:
		return strconv.Itoa(code)
	}
}

// AwardNominee is an award nominee record.
type AwardNominee struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Character *Character   `json:"character,omitempty"`
	Details   string       `json:"details"`
	Episode   *EpisodeBase `json:"episode,omitempty"`
	Movie     *MovieBase   `json:"movie,omitempty"`
	Series    *SeriesBase  `json:"series,omitempty"`
	Year      string       `json:"year"`
	Category  string       `json:"category"`
	IsWinner  bool         `json:"isWinner"`
}

// AwardCategory is a base award category record.
type AwardCategory struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	AllowCoNominees bool      `json:"allowCoNominees"`
	ForSeries       bool      `json:"forSeries"`
	ForMovies       bool      `json:"forMovies"`
	Award           AwardBase `json:"award"`
}

// AwardCategoryExtended is the extended award category record.
type AwardCategoryExtended struct {
	AwardCategory
	Nominees []AwardNominee `json:"nominees,omitempty"`
}

// AwardExtended is the extended award record.
type AwardExtended struct {
	AwardBase
	Score      int             `json:"score"`
	Categories []AwardCategory `json:"categories,omitempty"`
}

// CompanyType is a company type record.
type CompanyType struct {
	ID   int64  `json:"companyTypeId"`
	Name string `json:"companyTypeName"`
}

// Country is a country record.
type Country struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortCode string `json:"shortCode"`
}

// EntityType is an entity type record.
type EntityType struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	HasSpecials bool   `json:"hasSpecials"`
}

// Gender is a gender record.
type Gender struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// InspirationType is an inspiration type record.
type InspirationType struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ReferenceURL string `json:"reference_url"`
	URL          string `json:"url"`
}

// ListExtended is the extended record for a list.
type ListExtended struct {
	ListBase
	Entities []ListEntity `json:"entities,omitempty"`
}

// ListEntity is an entity within a list.
type ListEntity struct {
	Order    int64 `json:"order"`
	SeriesID int64 `json:"seriesId,omitempty"`
	MovieID  int64 `json:"movieId,omitempty"`
}

// PeopleType is a people type record.
type PeopleType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// SourceType is a source type record.
type SourceType struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Prefix  string `json:"prefix"`
	PostFix string `json:"postfix"`
	Sort    int    `json:"sort"`
}

// UserInfo contains authenticated user information.
type UserInfo struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Language             string `json:"language"`
	FavoritesDisplayMode string `json:"favoritesDisplaymode"`
}

// Favorites is the user's favorites list.
type Favorites struct {
	Series   []int `json:"series,omitempty"`
	Movies   []int `json:"movies,omitempty"`
	Episodes []int `json:"episodes,omitempty"`
	Artwork  []int `json:"artwork,omitempty"`
	People   []int `json:"people,omitempty"`
	Lists    []int `json:"lists,omitempty"`
}

// FavoriteRecord is used to add items to user favorites via POST.
type FavoriteRecord struct {
	Series   int `json:"series,omitempty"`
	Movies   int `json:"movies,omitempty"`
	Episodes int `json:"episodes,omitempty"`
	Artwork  int `json:"artwork,omitempty"`
	People   int `json:"people,omitempty"`
	Lists    int `json:"lists,omitempty"`
}

// FilterParams holds optional query parameters for filter endpoints.
type FilterParams struct {
	Country       string
	Language      string
	Company       int
	ContentRating int
	Genre         int
	Year          int
	Sort          string
	SortType      string
	Status        int
}

func (p *FilterParams) encode() string {
	if p == nil {
		return ""
	}
	q := url.Values{}
	if p.Country != "" {
		q.Set("country", p.Country)
	}
	if p.Language != "" {
		q.Set("lang", p.Language)
	}
	if p.Company > 0 {
		q.Set("company", strconv.Itoa(p.Company))
	}
	if p.ContentRating > 0 {
		q.Set("contentRating", strconv.Itoa(p.ContentRating))
	}
	if p.Genre > 0 {
		q.Set("genre", strconv.Itoa(p.Genre))
	}
	if p.Year > 0 {
		q.Set("year", strconv.Itoa(p.Year))
	}
	if p.Sort != "" {
		q.Set("sort", p.Sort)
	}
	if p.SortType != "" {
		q.Set("sortType", p.SortType)
	}
	if p.Status > 0 {
		q.Set("status", strconv.Itoa(p.Status))
	}
	return q.Encode()
}
