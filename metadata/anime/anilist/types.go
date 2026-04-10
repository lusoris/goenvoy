package anilist

// MediaType distinguishes anime from manga.
type MediaType string

const (
	// MediaTypeAnime represents anime entries.
	MediaTypeAnime MediaType = "ANIME"
	// MediaTypeManga represents manga entries.
	MediaTypeManga MediaType = "MANGA"
)

// FuzzyDate represents a date where any component may be unknown.
type FuzzyDate struct {
	Year  *int `json:"year"`
	Month *int `json:"month"`
	Day   *int `json:"day"`
}

// MediaTitle holds the official titles of a media in various languages.
type MediaTitle struct {
	Romaji        string `json:"romaji"`
	English       string `json:"english"`
	Native        string `json:"native"`
	UserPreferred string `json:"userPreferred"`
}

// MediaCoverImage holds cover image URLs and dominant color.
type MediaCoverImage struct {
	ExtraLarge string `json:"extraLarge"`
	Large      string `json:"large"`
	Medium     string `json:"medium"`
	Color      string `json:"color"`
}

// MediaTag describes an element or theme of a media entry.
type MediaTag struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Category         string `json:"category"`
	Rank             int    `json:"rank"`
	IsGeneralSpoiler bool   `json:"isGeneralSpoiler"`
	IsMediaSpoiler   bool   `json:"isMediaSpoiler"`
}

// AiringSchedule holds information about the next airing episode.
type AiringSchedule struct {
	AiringAt        int64 `json:"airingAt"`
	TimeUntilAiring int64 `json:"timeUntilAiring"`
	Episode         int   `json:"episode"`
}

// Media represents an anime or manga entry.
type Media struct {
	ID                int             `json:"id"`
	IDMal             *int            `json:"idMal"`
	Title             MediaTitle      `json:"title"`
	Type              string          `json:"type"`
	Format            string          `json:"format"`
	Status            string          `json:"status"`
	Description       string          `json:"description"`
	StartDate         FuzzyDate       `json:"startDate"`
	EndDate           FuzzyDate       `json:"endDate"`
	Season            string          `json:"season"`
	SeasonYear        *int            `json:"seasonYear"`
	Episodes          *int            `json:"episodes"`
	Duration          *int            `json:"duration"`
	Chapters          *int            `json:"chapters"`
	Volumes           *int            `json:"volumes"`
	CountryOfOrigin   string          `json:"countryOfOrigin"`
	IsLicensed        bool            `json:"isLicensed"`
	Source            string          `json:"source"`
	CoverImage        MediaCoverImage `json:"coverImage"`
	BannerImage       string          `json:"bannerImage"`
	Genres            []string        `json:"genres"`
	Synonyms          []string        `json:"synonyms"`
	AverageScore      *int            `json:"averageScore"`
	MeanScore         *int            `json:"meanScore"`
	Popularity        int             `json:"popularity"`
	Favorites         int             `json:"favourites"` //nolint:tagliatelle,misspell // API uses British spelling.
	IsAdult           bool            `json:"isAdult"`
	SiteURL           string          `json:"siteUrl"`
	Tags              []MediaTag      `json:"tags"`
	NextAiringEpisode *AiringSchedule `json:"nextAiringEpisode"`
}

// PageInfo holds pagination metadata.
type PageInfo struct {
	Total       int  `json:"total"`
	CurrentPage int  `json:"currentPage"`
	LastPage    int  `json:"lastPage"`
	HasNextPage bool `json:"hasNextPage"`
	PerPage     int  `json:"perPage"`
}

// MediaPage is a paginated list of media entries.
type MediaPage struct {
	PageInfo PageInfo `json:"pageInfo"`
	Media    []Media  `json:"media"`
}

// PersonName holds the name of a character or staff member.
type PersonName struct {
	First              string   `json:"first"`
	Middle             string   `json:"middle"`
	Last               string   `json:"last"`
	Full               string   `json:"full"`
	Native             string   `json:"native"`
	Alternative        []string `json:"alternative"`
	AlternativeSpoiler []string `json:"alternativeSpoiler"`
	UserPreferred      string   `json:"userPreferred"`
}

// PersonImage holds image URLs for a character or staff member.
type PersonImage struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
}

// Character represents a fictional character.
type Character struct {
	ID          int         `json:"id"`
	Name        PersonName  `json:"name"`
	Image       PersonImage `json:"image"`
	Description string      `json:"description"`
	Gender      string      `json:"gender"`
	DateOfBirth FuzzyDate   `json:"dateOfBirth"`
	Age         string      `json:"age"`
	SiteURL     string      `json:"siteUrl"`
	Favorites   int         `json:"favourites"` //nolint:tagliatelle,misspell // API uses British spelling.
}

// CharacterPage is a paginated list of characters.
type CharacterPage struct {
	PageInfo   PageInfo    `json:"pageInfo"`
	Characters []Character `json:"characters"`
}

// Staff represents a person who worked on a media entry.
type Staff struct {
	ID          int         `json:"id"`
	Name        PersonName  `json:"name"`
	Image       PersonImage `json:"image"`
	Description string      `json:"description"`
	Gender      string      `json:"gender"`
	DateOfBirth FuzzyDate   `json:"dateOfBirth"`
	DateOfDeath FuzzyDate   `json:"dateOfDeath"`
	Age         string      `json:"age"`
	HomeTown    string      `json:"homeTown"`
	SiteURL     string      `json:"siteUrl"`
	Favorites   int         `json:"favourites"` //nolint:tagliatelle,misspell // API uses British spelling.
}

// StaffPage is a paginated list of staff members.
type StaffPage struct {
	PageInfo PageInfo `json:"pageInfo"`
	Staff    []Staff  `json:"staff"`
}

// UserAvatar holds avatar image URLs.
type UserAvatar struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
}

// User represents an AniList user.
type User struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	About       string     `json:"about"`
	Avatar      UserAvatar `json:"avatar"`
	BannerImage string     `json:"bannerImage"`
	SiteURL     string     `json:"siteUrl"`
	CreatedAt   int64      `json:"createdAt"`
}

// GraphQLError is a single error returned by the AniList GraphQL API.
type GraphQLError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// APIError is returned when the AniList API responds with one or more GraphQL errors.
type APIError struct {
	Errors []GraphQLError `json:"errors"`
}

func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return "anilist: " + e.Errors[0].Message
	}
	return "anilist: unknown error"
}

// HTTPError is returned when the server responds with a non-JSON error body.
type HTTPError struct {
	StatusCode int
	Status     string
}

func (e *HTTPError) Error() string {
	return "anilist: HTTP " + e.Status
}

// Studio represents an animation studio.
type Studio struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// StudioEdge connects a studio to a media entry.
type StudioEdge struct {
	Node   Studio `json:"node"`
	IsMain bool   `json:"isMain"`
}

// StudioConnection contains studio edges.
type StudioConnection struct {
	Edges []StudioEdge `json:"edges"`
}

// CharacterEdge connects a character to a media entry with role and voice actors.
type CharacterEdge struct {
	Node        Character `json:"node"`
	Role        string    `json:"role"`
	VoiceActors []Staff   `json:"voiceActors"`
}

// CharacterConnectionDetailed contains character edges with roles and voice actors.
type CharacterConnectionDetailed struct {
	Edges []CharacterEdge `json:"edges"`
}

// StaffEdge connects a staff member to a media entry.
type StaffEdge struct {
	Node Staff  `json:"node"`
	Role string `json:"role"`
}

// StaffConnectionDetailed contains staff edges with roles.
type StaffConnectionDetailed struct {
	Edges []StaffEdge `json:"edges"`
}

// MediaEdge connects related media entries.
type MediaEdge struct {
	Node         Media  `json:"node"`
	RelationType string `json:"relationType"`
}

// MediaConnectionDetailed contains related media edges.
type MediaConnectionDetailed struct {
	Edges []MediaEdge `json:"edges"`
}

// ExternalLink is a link to an external site for a media entry.
type ExternalLink struct {
	ID       int     `json:"id"`
	URL      *string `json:"url"`
	Site     string  `json:"site"`
	SiteID   *int    `json:"siteId"`
	Type     string  `json:"type"`
	Language *string `json:"language"`
}

// StreamingEpisode is a streaming source for an episode.
type StreamingEpisode struct {
	Title     *string `json:"title"`
	Thumbnail *string `json:"thumbnail"`
	URL       *string `json:"url"`
	Site      *string `json:"site"`
}

// Trailer holds a media trailer reference.
type Trailer struct {
	ID        *string `json:"id"`
	Site      *string `json:"site"`
	Thumbnail *string `json:"thumbnail"`
}

// MediaDetailed extends Media with connection fields (studios, characters, staff, relations, etc.).
type MediaDetailed struct {
	Media
	Studios           StudioConnection             `json:"studios"`
	Characters        CharacterConnectionDetailed  `json:"characters"`
	Staff             StaffConnectionDetailed      `json:"staff"`
	Relations         MediaConnectionDetailed      `json:"relations"`
	ExternalLinks     []ExternalLink               `json:"externalLinks"`
	StreamingEpisodes []StreamingEpisode           `json:"streamingEpisodes"`
	Trailer           *Trailer                     `json:"trailer"`
}
