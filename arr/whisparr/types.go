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

// Movie represents a movie or scene in Whisparr v3.
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

// AddMovieOptions controls behavior when adding a new movie in Whisparr v3.
type AddMovieOptions struct {
	Monitor        string `json:"monitor,omitempty"`
	SearchForMovie bool   `json:"searchForMovie"`
	AddMethod      string `json:"addMethod,omitempty"`
}

// MovieFile represents a downloaded movie/scene file in Whisparr v3.
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

// HistoryRecordV3 represents a history event in Whisparr v3.
type HistoryRecordV3 struct {
	ID          int    `json:"id"`
	MovieID     int    `json:"movieId"`
	SourceTitle string `json:"sourceTitle,omitempty"`
	Date        string `json:"date,omitempty"`
	EventType   string `json:"eventType,omitempty"`
}

// ParseResultV3 holds the result of a Whisparr v3 title parse.
type ParseResultV3 struct {
	Title           string                 `json:"title,omitempty"`
	ParsedMovieInfo map[string]interface{} `json:"parsedMovieInfo,omitempty"`
	Movie           *Movie                 `json:"movie,omitempty"`
}

// MovieEditorResource is the payload for bulk movie editing in Whisparr v3.
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

// Performer represents a performer in Whisparr v3.
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

// Studio represents a production studio in Whisparr v3.
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

// Credit represents a performer credit on a movie/scene in Whisparr v3.
type Credit struct {
	ID                 int     `json:"id"`
	PerformerName      string  `json:"personName,omitempty"`
	PerformerForeignID string  `json:"performerForeignId,omitempty"`
	MovieMetadataID    int     `json:"movieMetadataId,omitempty"`
	Images             []Image `json:"images,omitempty"`
	CreditType         string  `json:"type,omitempty"`
}

// ImportExclusion represents an item excluded from import in Whisparr v3.
type ImportExclusion struct {
	ID         int    `json:"id"`
	ForeignID  string `json:"foreignId,omitempty"`
	MovieTitle string `json:"movieTitle,omitempty"`
	MovieYear  int    `json:"movieYear,omitempty"`
	StashID    string `json:"stashId,omitempty"`
}

// ---------- V2 additional types ----------.

// SeriesEditorResource is the request body for batch series editing/deleting.
type SeriesEditorResource struct {
	SeriesIDs        []int  `json:"seriesIds"`
	Monitored        *bool  `json:"monitored,omitempty"`
	QualityProfileID *int   `json:"qualityProfileId,omitempty"`
	SeriesType       string `json:"seriesType,omitempty"`
	SeasonFolder     *bool  `json:"seasonFolder,omitempty"`
	RootFolderPath   string `json:"rootFolderPath,omitempty"`
	Tags             []int  `json:"tags,omitempty"`
	ApplyTags        string `json:"applyTags,omitempty"`
	MoveFiles        bool   `json:"moveFiles,omitempty"`
	DeleteFiles      bool   `json:"deleteFiles,omitempty"`
}

// EpisodesMonitoredResource is the payload for bulk episode monitoring updates.
type EpisodesMonitoredResource struct {
	EpisodeIDs []int `json:"episodeIds"`
	Monitored  bool  `json:"monitored"`
}

// EpisodeFileEditorResource is the payload for bulk episode file editing.
type EpisodeFileEditorResource struct {
	EpisodeFileIDs []int  `json:"episodeFileIds"`
	Languages      []int  `json:"languages,omitempty"`
	Quality        any    `json:"quality,omitempty"`
	SceneName      string `json:"sceneName,omitempty"`
}

// ImportListConfigResource holds import list global configuration.
type ImportListConfigResource struct {
	ID               int    `json:"id"`
	ListSyncLevel    string `json:"listSyncLevel,omitempty"`
	ListSyncTag      int    `json:"listSyncTag,omitempty"`
	ListSyncInterval int    `json:"listSyncInterval,omitempty"`
}

// LanguageProfileResource represents a language profile.
type LanguageProfileResource struct {
	ID             int                           `json:"id"`
	Name           string                        `json:"name"`
	UpgradeAllowed bool                          `json:"upgradeAllowed"`
	Cutoff         *Language                     `json:"cutoff,omitempty"`
	Languages      []LanguageProfileItemResource `json:"languages,omitempty"`
}

// LanguageProfileItemResource represents an item in a language profile.
type LanguageProfileItemResource struct {
	Language Language `json:"language"`
	Allowed  bool     `json:"allowed"`
}

// LocalizationLanguageResource is the current localization language.
type LocalizationLanguageResource struct {
	Identifier string `json:"identifier,omitempty"`
}

// QualityDefinitionLimitsResource holds quality definition limits.
type QualityDefinitionLimitsResource struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// ---------- V3 additional types ----------.

// AlternativeTitleResource represents an alternative title for a movie.
type AlternativeTitleResource struct {
	ID              int    `json:"id"`
	SourceType      string `json:"sourceType,omitempty"`
	MovieMetadataID int    `json:"movieMetadataId,omitempty"`
	Title           string `json:"title,omitempty"`
	CleanTitle      string `json:"cleanTitle,omitempty"`
}

// ExtraFileResource represents an extra file associated with a movie.
type ExtraFileResource struct {
	ID           int    `json:"id"`
	MovieID      int    `json:"movieId"`
	MovieFileID  int    `json:"movieFileId,omitempty"`
	RelativePath string `json:"relativePath,omitempty"`
	Extension    string `json:"extension,omitempty"`
	Type         string `json:"type,omitempty"`
}

// PerformerEditorResource is the payload for bulk performer editing.
type PerformerEditorResource struct {
	PerformerIDs     []int  `json:"performerIds"`
	Monitored        *bool  `json:"monitored,omitempty"`
	QualityProfileID *int   `json:"qualityProfileId,omitempty"`
	RootFolderPath   string `json:"rootFolderPath,omitempty"`
	Tags             []int  `json:"tags,omitempty"`
	ApplyTags        string `json:"applyTags,omitempty"`
	MoveFiles        bool   `json:"moveFiles,omitempty"`
}

// StudioEditorResource is the payload for bulk studio editing.
type StudioEditorResource struct {
	StudioIDs        []int  `json:"studioIds"`
	Monitored        *bool  `json:"monitored,omitempty"`
	QualityProfileID *int   `json:"qualityProfileId,omitempty"`
	RootFolderPath   string `json:"rootFolderPath,omitempty"`
	Tags             []int  `json:"tags,omitempty"`
	ApplyTags        string `json:"applyTags,omitempty"`
	MoveFiles        bool   `json:"moveFiles,omitempty"`
}

// MovieFileEditorResource is the payload for bulk movie file editing.
type MovieFileEditorResource struct {
	MovieFileIDs []int `json:"movieFileIds"`
	Languages    []int `json:"languages,omitempty"`
	Quality      any   `json:"quality,omitempty"`
}
