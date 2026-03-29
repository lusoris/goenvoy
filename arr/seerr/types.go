package seerr

// User represents a Seerr user account.
type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email,omitempty"`
	Username     string `json:"username,omitempty"`
	PlexToken    string `json:"plexToken,omitempty"`
	PlexUsername string `json:"plexUsername,omitempty"`
	UserType     int    `json:"userType,omitempty"`
	Permissions  int    `json:"permissions,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	UpdatedAt    string `json:"updatedAt,omitempty"`
	RequestCount int    `json:"requestCount,omitempty"`
}

// MediaRequest represents a media request in Seerr.
type MediaRequest struct {
	ID          int        `json:"id"`
	Status      int        `json:"status"`
	Media       *MediaInfo `json:"media,omitempty"`
	CreatedAt   string     `json:"createdAt,omitempty"`
	UpdatedAt   string     `json:"updatedAt,omitempty"`
	RequestedBy *User      `json:"requestedBy,omitempty"`
	ModifiedBy  *User      `json:"modifiedBy,omitempty"`
	Is4k        bool       `json:"is4k"`
	ServerID    int        `json:"serverId,omitempty"`
	ProfileID   int        `json:"profileId,omitempty"`
	RootFolder  string     `json:"rootFolder,omitempty"`
}

// MediaInfo represents the media metadata and availability status.
type MediaInfo struct {
	ID        int            `json:"id"`
	TmdbID    int            `json:"tmdbId,omitempty"`
	TvdbID    *int           `json:"tvdbId,omitempty"`
	Status    int            `json:"status"`
	Requests  []MediaRequest `json:"requests,omitempty"`
	CreatedAt string         `json:"createdAt,omitempty"`
	UpdatedAt string         `json:"updatedAt,omitempty"`
}

// MovieResult represents a movie search or discover result.
type MovieResult struct {
	ID               int        `json:"id"`
	MediaType        string     `json:"mediaType,omitempty"`
	Popularity       float64    `json:"popularity,omitempty"`
	PosterPath       string     `json:"posterPath,omitempty"`
	BackdropPath     string     `json:"backdropPath,omitempty"`
	VoteCount        int        `json:"voteCount,omitempty"`
	VoteAverage      float64    `json:"voteAverage,omitempty"`
	GenreIDs         []int      `json:"genreIds,omitempty"`
	Overview         string     `json:"overview,omitempty"`
	OriginalLanguage string     `json:"originalLanguage,omitempty"`
	Title            string     `json:"title,omitempty"`
	OriginalTitle    string     `json:"originalTitle,omitempty"`
	ReleaseDate      string     `json:"releaseDate,omitempty"`
	Adult            bool       `json:"adult,omitempty"`
	Video            bool       `json:"video,omitempty"`
	MediaInfo        *MediaInfo `json:"mediaInfo,omitempty"`
}

// TvResult represents a TV show search or discover result.
type TvResult struct {
	ID               int        `json:"id"`
	MediaType        string     `json:"mediaType,omitempty"`
	Popularity       float64    `json:"popularity,omitempty"`
	PosterPath       string     `json:"posterPath,omitempty"`
	BackdropPath     string     `json:"backdropPath,omitempty"`
	VoteCount        int        `json:"voteCount,omitempty"`
	VoteAverage      float64    `json:"voteAverage,omitempty"`
	GenreIDs         []int      `json:"genreIds,omitempty"`
	Overview         string     `json:"overview,omitempty"`
	OriginalLanguage string     `json:"originalLanguage,omitempty"`
	Name             string     `json:"name,omitempty"`
	OriginalName     string     `json:"originalName,omitempty"`
	OriginCountry    []string   `json:"originCountry,omitempty"`
	FirstAirDate     string     `json:"firstAirDate,omitempty"`
	MediaInfo        *MediaInfo `json:"mediaInfo,omitempty"`
}

// PersonResult represents a person search result.
type PersonResult struct {
	ID          int    `json:"id"`
	ProfilePath string `json:"profilePath,omitempty"`
	Adult       bool   `json:"adult,omitempty"`
	MediaType   string `json:"mediaType,omitempty"`
	Name        string `json:"name,omitempty"`
}

// SearchResults is a paginated response for search and discover endpoints.
// Results can contain MovieResult, TvResult, or PersonResult items.
type SearchResults struct {
	Page         int              `json:"page"`
	TotalPages   int              `json:"totalPages"`
	TotalResults int              `json:"totalResults"`
	Results      []map[string]any `json:"results,omitempty"`
}

// MovieDetails contains full details for a movie.
type MovieDetails struct {
	ID                  int                 `json:"id"`
	ImdbID              string              `json:"imdbId,omitempty"`
	Adult               bool                `json:"adult,omitempty"`
	BackdropPath        string              `json:"backdropPath,omitempty"`
	PosterPath          string              `json:"posterPath,omitempty"`
	Budget              int64               `json:"budget,omitempty"`
	Genres              []Genre             `json:"genres,omitempty"`
	Homepage            string              `json:"homepage,omitempty"`
	OriginalLanguage    string              `json:"originalLanguage,omitempty"`
	OriginalTitle       string              `json:"originalTitle,omitempty"`
	Overview            string              `json:"overview,omitempty"`
	Popularity          float64             `json:"popularity,omitempty"`
	ProductionCompanies []ProductionCompany `json:"productionCompanies,omitempty"`
	ReleaseDate         string              `json:"releaseDate,omitempty"`
	Revenue             *int64              `json:"revenue,omitempty"`
	Runtime             int                 `json:"runtime,omitempty"`
	Status              string              `json:"status,omitempty"`
	Tagline             string              `json:"tagline,omitempty"`
	Title               string              `json:"title,omitempty"`
	Video               bool                `json:"video,omitempty"`
	VoteAverage         float64             `json:"voteAverage,omitempty"`
	VoteCount           int                 `json:"voteCount,omitempty"`
	MediaInfo           *MediaInfo          `json:"mediaInfo,omitempty"`
}

// TvDetails contains full details for a TV show.
type TvDetails struct {
	ID               int                 `json:"id"`
	BackdropPath     string              `json:"backdropPath,omitempty"`
	PosterPath       string              `json:"posterPath,omitempty"`
	CreatedBy        []CreatedBy         `json:"createdBy,omitempty"`
	FirstAirDate     string              `json:"firstAirDate,omitempty"`
	Genres           []Genre             `json:"genres,omitempty"`
	Homepage         string              `json:"homepage,omitempty"`
	InProduction     bool                `json:"inProduction,omitempty"`
	LastAirDate      string              `json:"lastAirDate,omitempty"`
	Name             string              `json:"name,omitempty"`
	Networks         []ProductionCompany `json:"networks,omitempty"`
	NumberOfEpisodes int                 `json:"numberOfEpisodes,omitempty"`
	NumberOfSeasons  int                 `json:"numberOfSeason,omitempty"`
	OriginalLanguage string              `json:"originalLanguage,omitempty"`
	OriginalName     string              `json:"originalName,omitempty"`
	Overview         string              `json:"overview,omitempty"`
	Popularity       float64             `json:"popularity,omitempty"`
	Seasons          []Season            `json:"seasons,omitempty"`
	Status           string              `json:"status,omitempty"`
	Tagline          string              `json:"tagline,omitempty"`
	Type             string              `json:"type,omitempty"`
	VoteAverage      float64             `json:"voteAverage,omitempty"`
	VoteCount        int                 `json:"voteCount,omitempty"`
	MediaInfo        *MediaInfo          `json:"mediaInfo,omitempty"`
}

// Season represents a TV season.
type Season struct {
	ID           int       `json:"id"`
	AirDate      string    `json:"airDate,omitempty"`
	EpisodeCount int       `json:"episodeCount,omitempty"`
	Name         string    `json:"name,omitempty"`
	Overview     string    `json:"overview,omitempty"`
	PosterPath   string    `json:"posterPath,omitempty"`
	SeasonNumber int       `json:"seasonNumber"`
	Episodes     []Episode `json:"episodes,omitempty"`
}

// Episode represents a TV episode.
type Episode struct {
	ID            int     `json:"id"`
	Name          string  `json:"name,omitempty"`
	AirDate       string  `json:"airDate,omitempty"`
	EpisodeNumber int     `json:"episodeNumber"`
	Overview      string  `json:"overview,omitempty"`
	SeasonNumber  int     `json:"seasonNumber"`
	StillPath     string  `json:"stillPath,omitempty"`
	VoteAverage   float64 `json:"voteAverage,omitempty"`
	VoteCount     int     `json:"voteCount,omitempty"`
}

// Genre represents a media genre.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// ProductionCompany represents a production company or TV network.
type ProductionCompany struct {
	ID            int    `json:"id"`
	LogoPath      string `json:"logoPath,omitempty"`
	OriginCountry string `json:"originCountry,omitempty"`
	Name          string `json:"name,omitempty"`
}

// CreatedBy represents a TV show creator.
type CreatedBy struct {
	ID          int    `json:"id"`
	Name        string `json:"name,omitempty"`
	Gender      int    `json:"gender,omitempty"`
	ProfilePath string `json:"profilePath,omitempty"`
}

// Issue represents an issue report in Seerr.
type Issue struct {
	ID         int            `json:"id"`
	IssueType  int            `json:"issueType"`
	Media      *MediaInfo     `json:"media,omitempty"`
	CreatedBy  *User          `json:"createdBy,omitempty"`
	ModifiedBy *User          `json:"modifiedBy,omitempty"`
	Comments   []IssueComment `json:"comments,omitempty"`
}

// IssueComment represents a comment on an issue.
type IssueComment struct {
	ID      int    `json:"id"`
	User    *User  `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
}

// PageInfo contains pagination metadata.
type PageInfo struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	Results int `json:"results"`
}

// RequestCount contains aggregated request counts.
type RequestCount struct {
	Total      int `json:"total"`
	Movie      int `json:"movie"`
	TV         int `json:"tv"`
	Pending    int `json:"pending"`
	Approved   int `json:"approved"`
	Declined   int `json:"declined"`
	Processing int `json:"processing"`
	Available  int `json:"available"`
}

// IssueCount contains aggregated issue counts.
type IssueCount struct {
	Total     int `json:"total"`
	Video     int `json:"video"`
	Audio     int `json:"audio"`
	Subtitles int `json:"subtitles"`
	Others    int `json:"others"`
	Open      int `json:"open"`
	Closed    int `json:"closed"`
}

// StatusResponse contains the Seerr server status.
type StatusResponse struct {
	Version         string `json:"version,omitempty"`
	CommitTag       string `json:"commitTag,omitempty"`
	UpdateAvailable bool   `json:"updateAvailable"`
	CommitsBehind   int    `json:"commitsBehind"`
	RestartRequired bool   `json:"restartRequired"`
}

// CreateRequestBody is the payload when creating a new media request.
type CreateRequestBody struct {
	MediaType         string `json:"mediaType"`
	MediaID           int    `json:"mediaId"`
	TvdbID            int    `json:"tvdbId,omitempty"`
	Seasons           any    `json:"seasons,omitempty"`
	Is4k              bool   `json:"is4k,omitempty"`
	ServerID          int    `json:"serverId,omitempty"`
	ProfileID         int    `json:"profileId,omitempty"`
	RootFolder        string `json:"rootFolder,omitempty"`
	LanguageProfileID int    `json:"languageProfileId,omitempty"`
	UserID            *int   `json:"userId,omitempty"`
}

// CreateIssueBody is the payload when creating a new issue.
type CreateIssueBody struct {
	IssueType int    `json:"issueType"`
	Message   string `json:"message"`
	MediaID   int    `json:"mediaId"`
}

// UserQuota contains quota details for a user.
type UserQuota struct {
	Movie *QuotaDetail `json:"movie,omitempty"`
	TV    *QuotaDetail `json:"tv,omitempty"`
}

// QuotaDetail contains the details for a single quota type.
type QuotaDetail struct {
	Days       int  `json:"days"`
	Limit      int  `json:"limit"`
	Used       int  `json:"used"`
	Remaining  int  `json:"remaining"`
	Restricted bool `json:"restricted"`
}
