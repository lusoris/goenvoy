package sonarr

// Series represents a TV series in Sonarr.
type Series struct {
	ID                int               `json:"id"`
	Title             string            `json:"title"`
	AlternateTitles   []AlternateTitle  `json:"alternateTitles"`
	SortTitle         string            `json:"sortTitle"`
	Status            string            `json:"status"`
	Ended             bool              `json:"ended"`
	Overview          string            `json:"overview"`
	NextAiring        string            `json:"nextAiring,omitempty"`
	PreviousAiring    string            `json:"previousAiring,omitempty"`
	Network           string            `json:"network"`
	AirTime           string            `json:"airTime"`
	Images            []Image           `json:"images"`
	OriginalLanguage  Language          `json:"originalLanguage"`
	RemotePoster      string            `json:"remotePoster,omitempty"`
	Seasons           []Season          `json:"seasons"`
	Year              int               `json:"year"`
	Path              string            `json:"path"`
	QualityProfileID  int               `json:"qualityProfileId"`
	SeasonFolder      bool              `json:"seasonFolder"`
	Monitored         bool              `json:"monitored"`
	MonitorNewItems   string            `json:"monitorNewItems"`
	UseSceneNumbering bool              `json:"useSceneNumbering"`
	Runtime           int               `json:"runtime"`
	TvdbID            int               `json:"tvdbId"`
	TvRageID          int               `json:"tvRageId"`
	TvMazeID          int               `json:"tvMazeId"`
	TmdbID            int               `json:"tmdbId"`
	FirstAired        string            `json:"firstAired,omitempty"`
	LastAired         string            `json:"lastAired,omitempty"`
	SeriesType        string            `json:"seriesType"`
	CleanTitle        string            `json:"cleanTitle"`
	ImdbID            string            `json:"imdbId"`
	TitleSlug         string            `json:"titleSlug"`
	RootFolderPath    string            `json:"rootFolderPath"`
	Folder            string            `json:"folder"`
	Certification     string            `json:"certification"`
	Genres            []string          `json:"genres"`
	Tags              []int             `json:"tags"`
	Added             string            `json:"added"`
	AddOptions        *AddSeriesOptions `json:"addOptions,omitempty"`
	Ratings           Ratings           `json:"ratings"`
	Statistics        *SeriesStatistics `json:"statistics,omitempty"`
}

// AlternateTitle is an alternative name for a series.
type AlternateTitle struct {
	Title             string `json:"title"`
	SeasonNumber      *int   `json:"seasonNumber,omitempty"`
	SceneSeasonNumber *int   `json:"sceneSeasonNumber,omitempty"`
	SceneOrigin       string `json:"sceneOrigin,omitempty"`
	Comment           string `json:"comment,omitempty"`
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

// Season holds season-level metadata and monitoring status.
type Season struct {
	SeasonNumber int               `json:"seasonNumber"`
	Monitored    bool              `json:"monitored"`
	Statistics   *SeasonStatistics `json:"statistics,omitempty"`
}

// SeasonStatistics contains file and episode counts for one season.
type SeasonStatistics struct {
	EpisodeFileCount  int      `json:"episodeFileCount"`
	EpisodeCount      int      `json:"episodeCount"`
	TotalEpisodeCount int      `json:"totalEpisodeCount"`
	SizeOnDisk        int64    `json:"sizeOnDisk"`
	ReleaseGroups     []string `json:"releaseGroups"`
	PercentOfEpisodes float64  `json:"percentOfEpisodes"`
	PreviousAiring    string   `json:"previousAiring,omitempty"`
	NextAiring        string   `json:"nextAiring,omitempty"`
}

// Ratings holds community rating data.
type Ratings struct {
	Votes int     `json:"votes"`
	Value float64 `json:"value"`
}

// SeriesStatistics contains aggregate file and episode counts for a series.
type SeriesStatistics struct {
	SeasonCount       int      `json:"seasonCount"`
	EpisodeFileCount  int      `json:"episodeFileCount"`
	EpisodeCount      int      `json:"episodeCount"`
	TotalEpisodeCount int      `json:"totalEpisodeCount"`
	SizeOnDisk        int64    `json:"sizeOnDisk"`
	ReleaseGroups     []string `json:"releaseGroups"`
	PercentOfEpisodes float64  `json:"percentOfEpisodes"`
}

// AddSeriesOptions controls behavior when adding a new series.
type AddSeriesOptions struct {
	IgnoreEpisodesWithFiles      bool   `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles   bool   `json:"ignoreEpisodesWithoutFiles"`
	Monitor                      string `json:"monitor"`
	SearchForMissingEpisodes     bool   `json:"searchForMissingEpisodes"`
	SearchForCutoffUnmetEpisodes bool   `json:"searchForCutoffUnmetEpisodes"`
}

// Episode represents a single TV episode in Sonarr.
type Episode struct {
	ID                         int     `json:"id"`
	SeriesID                   int     `json:"seriesId"`
	TvdbID                     int     `json:"tvdbId"`
	EpisodeFileID              int     `json:"episodeFileId"`
	SeasonNumber               int     `json:"seasonNumber"`
	EpisodeNumber              int     `json:"episodeNumber"`
	Title                      string  `json:"title"`
	AirDate                    string  `json:"airDate"`
	AirDateUtc                 string  `json:"airDateUtc,omitempty"`
	LastSearchTime             string  `json:"lastSearchTime,omitempty"`
	Runtime                    int     `json:"runtime"`
	FinaleType                 string  `json:"finaleType,omitempty"`
	Overview                   string  `json:"overview"`
	HasFile                    bool    `json:"hasFile"`
	Monitored                  bool    `json:"monitored"`
	AbsoluteEpisodeNumber      int     `json:"absoluteEpisodeNumber,omitempty"`
	SceneAbsoluteEpisodeNumber int     `json:"sceneAbsoluteEpisodeNumber,omitempty"`
	SceneEpisodeNumber         int     `json:"sceneEpisodeNumber,omitempty"`
	SceneSeasonNumber          int     `json:"sceneSeasonNumber,omitempty"`
	UnverifiedSceneNumbering   bool    `json:"unverifiedSceneNumbering"`
	GrabDate                   string  `json:"grabDate,omitempty"`
	Series                     *Series `json:"series,omitempty"`
	Images                     []Image `json:"images"`
}

// EpisodeFile represents a downloaded episode file on disk.
type EpisodeFile struct {
	ID                  int            `json:"id"`
	SeriesID            int            `json:"seriesId"`
	SeasonNumber        int            `json:"seasonNumber"`
	RelativePath        string         `json:"relativePath"`
	Path                string         `json:"path"`
	Size                int64          `json:"size"`
	DateAdded           string         `json:"dateAdded"`
	SceneName           string         `json:"sceneName,omitempty"`
	ReleaseGroup        string         `json:"releaseGroup,omitempty"`
	Languages           []Language     `json:"languages"`
	Quality             QualityModel   `json:"quality"`
	CustomFormats       []CustomFormat `json:"customFormats"`
	CustomFormatScore   int            `json:"customFormatScore"`
	MediaInfo           *MediaInfo     `json:"mediaInfo,omitempty"`
	QualityCutoffNotMet bool           `json:"qualityCutoffNotMet"`
}

// QualityModel pairs a quality definition with its revision.
type QualityModel struct {
	Quality  Quality  `json:"quality"`
	Revision Revision `json:"revision"`
}

// Quality identifies a quality tier (e.g. HDTV-720p).
type Quality struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Source     string `json:"source"`
	Resolution int    `json:"resolution"`
}

// Revision tracks repack and proper versioning for a release.
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
	AudioBitrate      int     `json:"audioBitrate"`
	AudioChannels     float64 `json:"audioChannels"`
	AudioCodec        string  `json:"audioCodec"`
	AudioLanguages    string  `json:"audioLanguages"`
	AudioStreamCount  int     `json:"audioStreamCount"`
	VideoBitDepth     int     `json:"videoBitDepth"`
	VideoBitrate      int     `json:"videoBitrate"`
	VideoCodec        string  `json:"videoCodec"`
	VideoFps          float64 `json:"videoFps"`
	VideoDynamicRange string  `json:"videoDynamicRange"`
	Resolution        string  `json:"resolution"`
	RunTime           string  `json:"runTime"`
	ScanType          string  `json:"scanType"`
	Subtitles         string  `json:"subtitles"`
}

// ParseResult contains the result of parsing a release title.
type ParseResult struct {
	ID                int                `json:"id"`
	Title             string             `json:"title"`
	ParsedEpisodeInfo *ParsedEpisodeInfo `json:"parsedEpisodeInfo"`
	Series            *Series            `json:"series,omitempty"`
	Episodes          []Episode          `json:"episodes"`
}

// ParsedEpisodeInfo holds the structured data extracted from a release title.
type ParsedEpisodeInfo struct {
	ReleaseTitle             string       `json:"releaseTitle"`
	SeriesTitle              string       `json:"seriesTitle"`
	Quality                  QualityModel `json:"quality"`
	SeasonNumber             int          `json:"seasonNumber"`
	EpisodeNumbers           []int        `json:"episodeNumbers"`
	AbsoluteEpisodeNumbers   []int        `json:"absoluteEpisodeNumbers"`
	Languages                []Language   `json:"languages"`
	FullSeason               bool         `json:"fullSeason"`
	Special                  bool         `json:"special"`
	ReleaseGroup             string       `json:"releaseGroup"`
	ReleaseHash              string       `json:"releaseHash"`
	IsDaily                  bool         `json:"isDaily"`
	IsAbsoluteNumbering      bool         `json:"isAbsoluteNumbering"`
	IsPossibleSpecialEpisode bool         `json:"isPossibleSpecialEpisode"`
}

// HistoryRecord represents an event in the download history.
type HistoryRecord struct {
	ID                  int               `json:"id"`
	EpisodeID           int               `json:"episodeId"`
	SeriesID            int               `json:"seriesId"`
	SourceTitle         string            `json:"sourceTitle"`
	Languages           []Language        `json:"languages"`
	Quality             QualityModel      `json:"quality"`
	CustomFormats       []CustomFormat    `json:"customFormats"`
	CustomFormatScore   int               `json:"customFormatScore"`
	QualityCutoffNotMet bool              `json:"qualityCutoffNotMet"`
	Date                string            `json:"date"`
	DownloadID          string            `json:"downloadId"`
	EventType           string            `json:"eventType"`
	Data                map[string]string `json:"data"`
	Episode             *Episode          `json:"episode,omitempty"`
	Series              *Series           `json:"series,omitempty"`
}

// EpisodesMonitoredResource is the request body for batch monitoring episodes.
type EpisodesMonitoredResource struct {
	EpisodeIDs []int `json:"episodeIds"`
	Monitored  bool  `json:"monitored"`
}

// EpisodeFileListResource is the request body for bulk episode file operations.
type EpisodeFileListResource struct {
	EpisodeFileIDs []int `json:"episodeFileIds"`
}

// SeasonPassResource is the request body for batch updating season monitoring.
type SeasonPassResource struct {
	Series            []SeasonPassSeries `json:"series"`
	MonitoringOptions *MonitoringOptions `json:"monitoringOptions,omitempty"`
}

// SeasonPassSeries identifies a series and seasons to update.
type SeasonPassSeries struct {
	ID        int      `json:"id"`
	Seasons   []Season `json:"seasons"`
	Monitored bool     `json:"monitored"`
}

// MonitoringOptions configures which episodes to monitor.
type MonitoringOptions struct {
	Monitor string `json:"monitor"`
}

// SeriesEditorResource is the request body for batch series editing/deleting.
type SeriesEditorResource struct {
	SeriesIDs              []int  `json:"seriesIds"`
	Monitored              *bool  `json:"monitored,omitempty"`
	MonitorNewItems        string `json:"monitorNewItems,omitempty"`
	QualityProfileID       *int   `json:"qualityProfileId,omitempty"`
	SeriesType             string `json:"seriesType,omitempty"`
	SeasonFolder           *bool  `json:"seasonFolder,omitempty"`
	RootFolderPath         string `json:"rootFolderPath,omitempty"`
	Tags                   []int  `json:"tags,omitempty"`
	ApplyTags              string `json:"applyTags,omitempty"`
	MoveFiles              bool   `json:"moveFiles,omitempty"`
	DeleteFiles            bool   `json:"deleteFiles,omitempty"`
	AddImportListExclusion bool   `json:"addImportListExclusion,omitempty"`
}

// EpisodeFileEditorResource is the request body for bulk updating episode file quality/language.
type EpisodeFileEditorResource struct {
	EpisodeFileIDs []int         `json:"episodeFileIds"`
	Languages      []Language    `json:"languages,omitempty"`
	Quality        *QualityModel `json:"quality,omitempty"`
	SceneName      string        `json:"sceneName,omitempty"`
	ReleaseGroup   string        `json:"releaseGroup,omitempty"`
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
	ID       int       `json:"id"`
	Language *Language `json:"language,omitempty"`
	Allowed  bool      `json:"allowed"`
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
