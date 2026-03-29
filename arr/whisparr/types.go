package whisparr

// Gender represents the gender of a performer.
type Gender string

// Gender constants.
const (
	GenderFemale    Gender = "female"
	GenderMale      Gender = "male"
	GenderOther     Gender = "other"
	GenderTransMale Gender = "transMale"
)

// Image represents a cover image for a media item.
type Image struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
	RemoteURL string `json:"remoteUrl,omitempty"`
}

// Language identifies a language by ID and name.
type Language struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Ratings holds community rating information.
type Ratings struct {
	Votes int     `json:"votes"`
	Value float64 `json:"value"`
}

// Actor represents a performer credit on a v2 episode.
type Actor struct {
	TpdbID    int     `json:"tpdbId,omitempty"`
	Name      string  `json:"name"`
	Character string  `json:"character,omitempty"`
	Gender    Gender  `json:"gender,omitempty"`
	Images    []Image `json:"images,omitempty"`
}

// Series represents a site/series in Whisparr v2.
type Series struct {
	ID               int               `json:"id"`
	Title            string            `json:"title"`
	SortTitle        string            `json:"sortTitle,omitempty"`
	Status           string            `json:"status,omitempty"`
	Overview         string            `json:"overview,omitempty"`
	Network          string            `json:"network,omitempty"`
	Images           []Image           `json:"images,omitempty"`
	Seasons          []Season          `json:"seasons,omitempty"`
	Year             int               `json:"year,omitempty"`
	Path             string            `json:"path,omitempty"`
	QualityProfileID int               `json:"qualityProfileId"`
	Monitored        bool              `json:"monitored"`
	SeriesType       string            `json:"seriesType,omitempty"`
	Runtime          int               `json:"runtime,omitempty"`
	TvdbID           int               `json:"tvdbId,omitempty"`
	CleanTitle       string            `json:"cleanTitle,omitempty"`
	TitleSlug        string            `json:"titleSlug,omitempty"`
	RootFolderPath   string            `json:"rootFolderPath,omitempty"`
	Certification    string            `json:"certification,omitempty"`
	Genres           []string          `json:"genres,omitempty"`
	Tags             []int             `json:"tags,omitempty"`
	Added            string            `json:"added,omitempty"`
	AddOptions       *AddSeriesOptions `json:"addOptions,omitempty"`
	Ratings          Ratings           `json:"ratings,omitempty"`
	Statistics       *SeriesStatistics `json:"statistics,omitempty"`
}

// AddSeriesOptions controls behavior when adding a new series.
type AddSeriesOptions struct {
	Monitor                  string `json:"monitor,omitempty"`
	SearchForMissingEpisodes bool   `json:"searchForMissingEpisodes"`
}

// Season holds season-level metadata and monitoring status.
type Season struct {
	SeasonNumber int  `json:"seasonNumber"`
	Monitored    bool `json:"monitored"`
}

// SeriesStatistics contains episode counts for a series.
type SeriesStatistics struct {
	SeasonCount       int   `json:"seasonCount"`
	EpisodeFileCount  int   `json:"episodeFileCount"`
	EpisodeCount      int   `json:"episodeCount"`
	TotalEpisodeCount int   `json:"totalEpisodeCount"`
	SizeOnDisk        int64 `json:"sizeOnDisk"`
}

// Episode represents a scene/episode in Whisparr v2.
type Episode struct {
	ID                    int     `json:"id"`
	SeriesID              int     `json:"seriesId"`
	SeasonNumber          int     `json:"seasonNumber"`
	EpisodeNumber         int     `json:"episodeNumber"`
	Title                 string  `json:"title,omitempty"`
	AirDate               string  `json:"airDate,omitempty"`
	Overview              string  `json:"overview,omitempty"`
	EpisodeFileID         int     `json:"episodeFileId,omitempty"`
	HasFile               bool    `json:"hasFile"`
	Monitored             bool    `json:"monitored"`
	AbsoluteEpisodeNumber *int    `json:"absoluteEpisodeNumber,omitempty"`
	Actors                []Actor `json:"actors,omitempty"`
	Images                []Image `json:"images,omitempty"`
}

// EpisodeFile represents a downloaded scene file in Whisparr v2.
type EpisodeFile struct {
	ID           int        `json:"id"`
	SeriesID     int        `json:"seriesId"`
	SeasonNumber int        `json:"seasonNumber"`
	RelativePath string     `json:"relativePath,omitempty"`
	Path         string     `json:"path,omitempty"`
	Size         int64      `json:"size"`
	DateAdded    string     `json:"dateAdded,omitempty"`
	ReleaseGroup string     `json:"releaseGroup,omitempty"`
	Languages    []Language `json:"languages,omitempty"`
}

// V2HistoryRecord represents a history event in Whisparr v2.
type V2HistoryRecord struct {
	ID          int    `json:"id"`
	EpisodeID   int    `json:"episodeId"`
	SeriesID    int    `json:"seriesId"`
	SourceTitle string `json:"sourceTitle,omitempty"`
	Date        string `json:"date,omitempty"`
	EventType   string `json:"eventType,omitempty"`
}

// V2ParseResult holds the result of a v2 title parse.
type V2ParseResult struct {
	Title             string                 `json:"title,omitempty"`
	ParsedEpisodeInfo map[string]interface{} `json:"parsedEpisodeInfo,omitempty"`
	Series            *Series                `json:"series,omitempty"`
	Episodes          []Episode              `json:"episodes,omitempty"`
}

// SeasonPassResource is the payload for updating monitoring on multiple series.
type SeasonPassResource struct {
	Series []SeasonPassSeries `json:"series"`
}

// SeasonPassSeries identifies a series and its monitored seasons.
type SeasonPassSeries struct {
	ID      int      `json:"id"`
	Seasons []Season `json:"seasons,omitempty"`
}

// Movie represents a movie or scene in Whisparr v3 (Eros).
type Movie struct {
	ID                  int              `json:"id"`
	Title               string           `json:"title"`
	OriginalTitle       string           `json:"originalTitle,omitempty"`
	SortTitle           string           `json:"sortTitle,omitempty"`
	SizeOnDisk          int64            `json:"sizeOnDisk,omitempty"`
	Status              string           `json:"status,omitempty"`
	Overview            string           `json:"overview,omitempty"`
	ReleaseDate         string           `json:"releaseDate,omitempty"`
	Images              []Image          `json:"images,omitempty"`
	Year                int              `json:"year,omitempty"`
	Path                string           `json:"path,omitempty"`
	QualityProfileID    int              `json:"qualityProfileId"`
	HasFile             bool             `json:"hasFile,omitempty"`
	MovieFileID         int              `json:"movieFileId,omitempty"`
	Monitored           bool             `json:"monitored"`
	MinimumAvailability string           `json:"minimumAvailability,omitempty"`
	Runtime             int              `json:"runtime,omitempty"`
	CleanTitle          string           `json:"cleanTitle,omitempty"`
	TitleSlug           string           `json:"titleSlug,omitempty"`
	RootFolderPath      string           `json:"rootFolderPath,omitempty"`
	Genres              []string         `json:"genres,omitempty"`
	Tags                []int            `json:"tags,omitempty"`
	Added               string           `json:"added,omitempty"`
	AddOptions          *AddMovieOptions `json:"addOptions,omitempty"`
	Ratings             Ratings          `json:"ratings,omitempty"`
	MovieFile           *MovieFile       `json:"movieFile,omitempty"`
	Code                string           `json:"code,omitempty"`
	StudioTitle         string           `json:"studioTitle,omitempty"`
	StudioForeignID     string           `json:"studioForeignId,omitempty"`
	ForeignID           string           `json:"foreignId,omitempty"`
	StashID             string           `json:"stashId,omitempty"`
	Credits             []Credit         `json:"credits,omitempty"`
	ItemType            string           `json:"itemType,omitempty"`
}

// AddMovieOptions controls behavior when adding a new movie in Eros.
type AddMovieOptions struct {
	Monitor        string `json:"monitor,omitempty"`
	SearchForMovie bool   `json:"searchForMovie"`
	AddMethod      string `json:"addMethod,omitempty"`
}

// MovieFile represents a downloaded movie/scene file in Eros.
type MovieFile struct {
	ID           int        `json:"id"`
	MovieID      int        `json:"movieId"`
	RelativePath string     `json:"relativePath,omitempty"`
	Path         string     `json:"path,omitempty"`
	Size         int64      `json:"size"`
	DateAdded    string     `json:"dateAdded,omitempty"`
	ReleaseGroup string     `json:"releaseGroup,omitempty"`
	Languages    []Language `json:"languages,omitempty"`
}

// ErosHistoryRecord represents a history event in Eros.
type ErosHistoryRecord struct {
	ID          int    `json:"id"`
	MovieID     int    `json:"movieId"`
	SourceTitle string `json:"sourceTitle,omitempty"`
	Date        string `json:"date,omitempty"`
	EventType   string `json:"eventType,omitempty"`
}

// ErosParseResult holds the result of an Eros title parse.
type ErosParseResult struct {
	Title           string                 `json:"title,omitempty"`
	ParsedMovieInfo map[string]interface{} `json:"parsedMovieInfo,omitempty"`
	Movie           *Movie                 `json:"movie,omitempty"`
}

// MovieEditorResource is the payload for bulk movie editing in Eros.
type MovieEditorResource struct {
	MovieIDs           []int  `json:"movieIds"`
	Monitored          *bool  `json:"monitored,omitempty"`
	QualityProfileID   *int   `json:"qualityProfileId,omitempty"`
	RootFolderPath     string `json:"rootFolderPath,omitempty"`
	Tags               []int  `json:"tags,omitempty"`
	ApplyTags          string `json:"applyTags,omitempty"`
	MoveFiles          bool   `json:"moveFiles,omitempty"`
	DeleteFiles        bool   `json:"deleteFiles,omitempty"`
	AddImportExclusion bool   `json:"addImportExclusion,omitempty"`
}

// Performer represents a performer in Eros.
type Performer struct {
	ID               int     `json:"id"`
	ForeignID        string  `json:"foreignId,omitempty"`
	StashID          string  `json:"stashId,omitempty"`
	Name             string  `json:"name"`
	SortName         string  `json:"sortName,omitempty"`
	Gender           Gender  `json:"gender,omitempty"`
	Images           []Image `json:"images,omitempty"`
	Monitored        bool    `json:"monitored"`
	RootFolderPath   string  `json:"rootFolderPath,omitempty"`
	QualityProfileID int     `json:"qualityProfileId,omitempty"`
	SearchOnAdd      bool    `json:"searchOnAdd,omitempty"`
	Tags             []int   `json:"tags,omitempty"`
	Added            string  `json:"added,omitempty"`
	TotalSceneCount  int     `json:"totalSceneCount,omitempty"`
	SceneCount       int     `json:"sceneCount,omitempty"`
	SizeOnDisk       int64   `json:"sizeOnDisk,omitempty"`
}

// Studio represents a production studio in Eros.
type Studio struct {
	ID               int     `json:"id"`
	Title            string  `json:"title"`
	SortTitle        string  `json:"sortTitle,omitempty"`
	ForeignID        string  `json:"foreignId,omitempty"`
	StashID          string  `json:"stashId,omitempty"`
	Website          string  `json:"website,omitempty"`
	Network          string  `json:"network,omitempty"`
	Images           []Image `json:"images,omitempty"`
	Monitored        bool    `json:"monitored"`
	RootFolderPath   string  `json:"rootFolderPath,omitempty"`
	QualityProfileID int     `json:"qualityProfileId,omitempty"`
	SearchOnAdd      bool    `json:"searchOnAdd,omitempty"`
	Tags             []int   `json:"tags,omitempty"`
	TotalSceneCount  int     `json:"totalSceneCount,omitempty"`
	SceneCount       int     `json:"sceneCount,omitempty"`
	SizeOnDisk       int64   `json:"sizeOnDisk,omitempty"`
}

// Credit represents a performer credit on a movie/scene in Eros.
type Credit struct {
	ID                 int     `json:"id"`
	PerformerName      string  `json:"personName,omitempty"`
	PerformerForeignID string  `json:"performerForeignId,omitempty"`
	MovieMetadataID    int     `json:"movieMetadataId,omitempty"`
	Images             []Image `json:"images,omitempty"`
	CreditType         string  `json:"type,omitempty"`
}

// ImportExclusion represents an item excluded from import in Eros.
type ImportExclusion struct {
	ID         int    `json:"id"`
	ForeignID  string `json:"foreignId,omitempty"`
	MovieTitle string `json:"movieTitle,omitempty"`
	MovieYear  int    `json:"movieYear,omitempty"`
	StashID    string `json:"stashId,omitempty"`
}
