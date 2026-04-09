package radarr

// Movie represents a movie in Radarr.
type Movie struct {
	ID                  int              `json:"id"`
	Title               string           `json:"title"`
	OriginalTitle       string           `json:"originalTitle,omitempty"`
	OriginalLanguage    Language         `json:"originalLanguage"`
	AlternateTitles     []AlternateTitle `json:"alternateTitles,omitempty"`
	SortTitle           string           `json:"sortTitle"`
	SizeOnDisk          int64            `json:"sizeOnDisk,omitempty"`
	Status              string           `json:"status"`
	Overview            string           `json:"overview"`
	InCinemas           string           `json:"inCinemas,omitempty"`
	PhysicalRelease     string           `json:"physicalRelease,omitempty"`
	DigitalRelease      string           `json:"digitalRelease,omitempty"`
	ReleaseDate         string           `json:"releaseDate,omitempty"`
	PhysicalReleaseNote string           `json:"physicalReleaseNote,omitempty"`
	Images              []Image          `json:"images"`
	Website             string           `json:"website,omitempty"`
	RemotePoster        string           `json:"remotePoster,omitempty"`
	Year                int              `json:"year"`
	YouTubeTrailerID    string           `json:"youTubeTrailerId,omitempty"`
	Studio              string           `json:"studio,omitempty"`
	Path                string           `json:"path"`
	QualityProfileID    int              `json:"qualityProfileId"`
	HasFile             bool             `json:"hasFile,omitempty"`
	MovieFileID         int              `json:"movieFileId,omitempty"`
	Monitored           bool             `json:"monitored"`
	MinimumAvailability string           `json:"minimumAvailability"`
	IsAvailable         bool             `json:"isAvailable,omitempty"`
	FolderName          string           `json:"folderName,omitempty"`
	Runtime             int              `json:"runtime"`
	CleanTitle          string           `json:"cleanTitle,omitempty"`
	ImdbID              string           `json:"imdbId,omitempty"`
	TmdbID              int              `json:"tmdbId"`
	TitleSlug           string           `json:"titleSlug"`
	RootFolderPath      string           `json:"rootFolderPath,omitempty"`
	Folder              string           `json:"folder,omitempty"`
	Certification       string           `json:"certification,omitempty"`
	Genres              []string         `json:"genres"`
	Tags                []int            `json:"tags"`
	Added               string           `json:"added"`
	AddOptions          *AddMovieOptions `json:"addOptions,omitempty"`
	Ratings             Ratings          `json:"ratings"`
	MovieFile           *MovieFile       `json:"movieFile,omitempty"`
	Collection          *MovieCollection `json:"collection,omitempty"`
	Popularity          float32          `json:"popularity,omitempty"`
	Statistics          *MovieStatistics `json:"statistics,omitempty"`
}

// AlternateTitle is an alternative name for a movie.
type AlternateTitle struct {
	SourceType      string `json:"sourceType,omitempty"`
	MovieMetadataID int    `json:"movieMetadataId,omitempty"`
	Title           string `json:"title"`
	CleanTitle      string `json:"cleanTitle,omitempty"`
}

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

// AddMovieOptions controls behavior when adding a new movie.
type AddMovieOptions struct {
	IgnoreEpisodesWithFiles    bool   `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles bool   `json:"ignoreEpisodesWithoutFiles"`
	Monitor                    string `json:"monitor"`
	SearchForMovie             bool   `json:"searchForMovie"`
	AddMethod                  string `json:"addMethod,omitempty"`
}

// Ratings holds community rating data from different sources.
type Ratings struct {
	Imdb           *RatingChild `json:"imdb,omitempty"`
	Tmdb           *RatingChild `json:"tmdb,omitempty"`
	Metacritic     *RatingChild `json:"metacritic,omitempty"`
	RottenTomatoes *RatingChild `json:"rottenTomatoes,omitempty"`
}

// RatingChild contains votes and value for a single rating source.
type RatingChild struct {
	Votes int     `json:"votes"`
	Value float64 `json:"value"`
	Type  string  `json:"type,omitempty"`
}

// MovieCollection identifies the collection a movie belongs to.
type MovieCollection struct {
	Title  string `json:"title,omitempty"`
	TmdbID int    `json:"tmdbId"`
}

// MovieStatistics contains file counts and size information for a movie.
type MovieStatistics struct {
	MovieFileCount int      `json:"movieFileCount"`
	SizeOnDisk     int64    `json:"sizeOnDisk"`
	ReleaseGroups  []string `json:"releaseGroups,omitempty"`
}

// MovieFile represents a downloaded movie file on disk.
type MovieFile struct {
	ID                  int            `json:"id"`
	MovieID             int            `json:"movieId"`
	RelativePath        string         `json:"relativePath"`
	Path                string         `json:"path"`
	Size                int64          `json:"size"`
	DateAdded           string         `json:"dateAdded"`
	SceneName           string         `json:"sceneName,omitempty"`
	ReleaseGroup        string         `json:"releaseGroup,omitempty"`
	Edition             string         `json:"edition,omitempty"`
	Languages           []Language     `json:"languages"`
	Quality             QualityModel   `json:"quality"`
	CustomFormats       []CustomFormat `json:"customFormats,omitempty"`
	CustomFormatScore   *int           `json:"customFormatScore,omitempty"`
	IndexerFlags        *int           `json:"indexerFlags,omitempty"`
	MediaInfo           *MediaInfo     `json:"mediaInfo,omitempty"`
	OriginalFilePath    string         `json:"originalFilePath,omitempty"`
	QualityCutoffNotMet bool           `json:"qualityCutoffNotMet"`
}

// QualityModel pairs a quality definition with its revision.
type QualityModel struct {
	Quality  Quality  `json:"quality"`
	Revision Revision `json:"revision"`
}

// Quality identifies a quality tier (e.g. Bluray-1080p).
type Quality struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Source     string `json:"source"`
	Resolution int    `json:"resolution"`
}

// Revision tracks repack and version information for a release.
type Revision struct {
	Version  int  `json:"version"`
	Real     int  `json:"real"`
	IsRepack bool `json:"isRepack"`
}

// CustomFormat describes a custom format definition applied to a release.
type CustomFormat struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MediaInfo holds technical metadata about a video file.
type MediaInfo struct {
	AudioBitrate          int     `json:"audioBitrate"`
	AudioChannels         float64 `json:"audioChannels"`
	AudioCodec            string  `json:"audioCodec"`
	AudioLanguages        string  `json:"audioLanguages"`
	AudioStreamCount      int     `json:"audioStreamCount"`
	VideoBitDepth         int     `json:"videoBitDepth"`
	VideoBitrate          int     `json:"videoBitrate"`
	VideoCodec            string  `json:"videoCodec"`
	VideoFps              float64 `json:"videoFps"`
	VideoDynamicRange     string  `json:"videoDynamicRange"`
	VideoDynamicRangeType string  `json:"videoDynamicRangeType,omitempty"`
	Resolution            string  `json:"resolution"`
	RunTime               string  `json:"runTime"`
	ScanType              string  `json:"scanType"`
	Subtitles             string  `json:"subtitles"`
}

// Collection represents a movie collection in Radarr.
type Collection struct {
	ID                  int               `json:"id"`
	Title               string            `json:"title"`
	SortTitle           string            `json:"sortTitle,omitempty"`
	TmdbID              int               `json:"tmdbId"`
	Images              []Image           `json:"images,omitempty"`
	Overview            string            `json:"overview,omitempty"`
	Monitored           bool              `json:"monitored"`
	RootFolderPath      string            `json:"rootFolderPath,omitempty"`
	QualityProfileID    int               `json:"qualityProfileId,omitempty"`
	SearchOnAdd         bool              `json:"searchOnAdd"`
	MinimumAvailability string            `json:"minimumAvailability,omitempty"`
	Movies              []CollectionMovie `json:"movies,omitempty"`
	MissingMovies       int               `json:"missingMovies,omitempty"`
	Tags                []int             `json:"tags,omitempty"`
}

// CollectionMovie is a movie within a collection.
type CollectionMovie struct {
	TmdbID     int      `json:"tmdbId"`
	ImdbID     string   `json:"imdbId,omitempty"`
	Title      string   `json:"title"`
	CleanTitle string   `json:"cleanTitle,omitempty"`
	SortTitle  string   `json:"sortTitle,omitempty"`
	Status     string   `json:"status,omitempty"`
	Overview   string   `json:"overview,omitempty"`
	Runtime    int      `json:"runtime,omitempty"`
	Images     []Image  `json:"images,omitempty"`
	Year       int      `json:"year,omitempty"`
	Ratings    Ratings  `json:"ratings,omitempty"`
	Genres     []string `json:"genres,omitempty"`
	Folder     string   `json:"folder,omitempty"`
	IsExisting bool     `json:"isExisting"`
	IsExcluded bool     `json:"isExcluded"`
}

// ParseResult contains the result of parsing a release title.
type ParseResult struct {
	ID                int              `json:"id"`
	Title             string           `json:"title"`
	ParsedMovieInfo   *ParsedMovieInfo `json:"parsedMovieInfo,omitempty"`
	Movie             *Movie           `json:"movie,omitempty"`
	Languages         []Language       `json:"languages,omitempty"`
	CustomFormats     []CustomFormat   `json:"customFormats,omitempty"`
	CustomFormatScore int              `json:"customFormatScore"`
}

// ParsedMovieInfo holds the structured data extracted from a release title.
type ParsedMovieInfo struct {
	MovieTitles        []string     `json:"movieTitles,omitempty"`
	OriginalTitle      string       `json:"originalTitle,omitempty"`
	ReleaseTitle       string       `json:"releaseTitle,omitempty"`
	SimpleReleaseTitle string       `json:"simpleReleaseTitle,omitempty"`
	Quality            QualityModel `json:"quality"`
	Languages          []Language   `json:"languages,omitempty"`
	ReleaseGroup       string       `json:"releaseGroup,omitempty"`
	ReleaseHash        string       `json:"releaseHash,omitempty"`
	Edition            string       `json:"edition,omitempty"`
	Year               int          `json:"year,omitempty"`
	ImdbID             string       `json:"imdbId,omitempty"`
	TmdbID             int          `json:"tmdbId,omitempty"`
	HardcodedSubs      string       `json:"hardcodedSubs,omitempty"`
	MovieTitle         string       `json:"movieTitle,omitempty"`
}

// HistoryRecord represents an event in the download history.
type HistoryRecord struct {
	ID                  int               `json:"id"`
	MovieID             int               `json:"movieId"`
	SourceTitle         string            `json:"sourceTitle"`
	Languages           []Language        `json:"languages,omitempty"`
	Quality             QualityModel      `json:"quality"`
	CustomFormats       []CustomFormat    `json:"customFormats,omitempty"`
	CustomFormatScore   int               `json:"customFormatScore"`
	QualityCutoffNotMet bool              `json:"qualityCutoffNotMet"`
	Date                string            `json:"date"`
	DownloadID          string            `json:"downloadId,omitempty"`
	EventType           string            `json:"eventType"`
	Data                map[string]string `json:"data,omitempty"`
	Movie               *Movie            `json:"movie,omitempty"`
}

// Credit represents a cast or crew member for a movie.
type Credit struct {
	ID              int     `json:"id"`
	PersonName      string  `json:"personName"`
	CreditTmdbID    string  `json:"creditTmdbId,omitempty"`
	PersonTmdbID    int     `json:"personTmdbId"`
	MovieMetadataID int     `json:"movieMetadataId"`
	Images          []Image `json:"images,omitempty"`
	Department      string  `json:"department,omitempty"`
	Job             string  `json:"job,omitempty"`
	Character       string  `json:"character,omitempty"`
	Order           int     `json:"order"`
	Type            string  `json:"type"`
}

// MovieEditorResource is used for batch editing movies.
type MovieEditorResource struct {
	MovieIDs            []int  `json:"movieIds"`
	Monitored           *bool  `json:"monitored,omitempty"`
	QualityProfileID    *int   `json:"qualityProfileId,omitempty"`
	MinimumAvailability string `json:"minimumAvailability,omitempty"`
	RootFolderPath      string `json:"rootFolderPath,omitempty"`
	Tags                []int  `json:"tags,omitempty"`
	ApplyTags           string `json:"applyTags,omitempty"`
	MoveFiles           bool   `json:"moveFiles"`
	DeleteFiles         bool   `json:"deleteFiles"`
	AddImportExclusion  bool   `json:"addImportExclusion"`
}

// MovieFileListResource is the request body for bulk movie file operations.
type MovieFileListResource struct {
	MovieFileIDs []int `json:"movieFileIds"`
}

// ImportListExclusion represents a movie excluded from import lists.
type ImportListExclusion struct {
	ID         int    `json:"id"`
	TmdbID     int    `json:"tmdbId"`
	MovieTitle string `json:"movieTitle,omitempty"`
	MovieYear  int    `json:"movieYear,omitempty"`
}

// LocalizationResource represents localization strings.
type LocalizationResource struct {
	ID      int               `json:"id"`
	Strings map[string]string `json:"strings"`
}

// FileSystemResource represents the result of a filesystem browse.
type FileSystemResource struct {
	Directories []FileSystemEntry `json:"directories"`
	Files       []FileSystemEntry `json:"files"`
}

// FileSystemEntry is a single directory or file entry.
type FileSystemEntry struct {
	Path         string `json:"path"`
	RelativePath string `json:"relativePath,omitempty"`
	Name         string `json:"name"`
	Size         int64  `json:"size,omitempty"`
}

// ImportListConfigResource holds import list global configuration.
type ImportListConfigResource struct {
	ID            int    `json:"id"`
	ListSyncLevel string `json:"listSyncLevel"`
	ListSyncTag   int    `json:"listSyncTag"`
}

// MovieFileEditorResource is the request body for bulk updating movie file quality/language.
type MovieFileEditorResource struct {
	MovieFileIDs []int         `json:"movieFileIds"`
	Languages    []Language    `json:"languages,omitempty"`
	Quality      *QualityModel `json:"quality,omitempty"`
	ReleaseGroup string        `json:"releaseGroup,omitempty"`
}

// ExtraFileResource represents an extra file associated with a movie.
type ExtraFileResource struct {
	ID           int      `json:"id"`
	MovieID      int      `json:"movieId"`
	MovieFileID  int      `json:"movieFileId,omitempty"`
	RelativePath string   `json:"relativePath,omitempty"`
	Extension    string   `json:"extension,omitempty"`
	LanguageTags []string `json:"languageTags,omitempty"`
	Title        string   `json:"title,omitempty"`
	Type         string   `json:"type,omitempty"`
}

// MetadataConfigResource holds metadata global configuration.
type MetadataConfigResource struct {
	ID                   int    `json:"id"`
	CertificationCountry string `json:"certificationCountry,omitempty"`
}

// AlternativeTitleResource represents an alternative title returned by the API.
type AlternativeTitleResource struct {
	ID              int    `json:"id"`
	SourceType      string `json:"sourceType,omitempty"`
	MovieMetadataID int    `json:"movieMetadataId,omitempty"`
	Title           string `json:"title,omitempty"`
	CleanTitle      string `json:"cleanTitle,omitempty"`
}

// LocalizationLanguageResource represents a localization language identifier.
type LocalizationLanguageResource struct {
	Identifier string `json:"identifier"`
}

// QualityDefinitionLimitsResource holds min/max limits for quality definitions.
type QualityDefinitionLimitsResource struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
