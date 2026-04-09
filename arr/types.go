package arr

import "time"

// StatusResponse holds the system status returned by /api/v3/system/status.
type StatusResponse struct {
	AppName                string `json:"appName"`
	InstanceName           string `json:"instanceName"`
	Version                string `json:"version"`
	BuildTime              string `json:"buildTime"`
	IsDebug                bool   `json:"isDebug"`
	IsProduction           bool   `json:"isProduction"`
	IsAdmin                bool   `json:"isAdmin"`
	IsUserInteractive      bool   `json:"isUserInteractive"`
	StartupPath            string `json:"startupPath"`
	AppData                string `json:"appData"`
	OsName                 string `json:"osName"`
	OsVersion              string `json:"osVersion"`
	IsMonoRuntime          bool   `json:"isMonoRuntime"`
	IsMono                 bool   `json:"isMono"`
	IsLinux                bool   `json:"isLinux"`
	IsOsx                  bool   `json:"isOsx"`
	IsWindows              bool   `json:"isWindows"`
	IsDocker               bool   `json:"isDocker"`
	Branch                 string `json:"branch"`
	Authentication         string `json:"authentication"`
	SqliteVersion          string `json:"sqliteVersion"`
	MigrationVersion       int    `json:"migrationVersion"`
	URLBase                string `json:"urlBase"`
	RuntimeVersion         string `json:"runtimeVersion"`
	RuntimeName            string `json:"runtimeName"`
	StartTime              string `json:"startTime"`
	PackageVersion         string `json:"packageVersion"`
	PackageAuthor          string `json:"packageAuthor"`
	PackageUpdateMechanism string `json:"packageUpdateMechanism"`
}

// HealthCheck represents a single health-check entry from /api/v3/health.
type HealthCheck struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	WikiURL string `json:"wikiUrl"`
}

// StatusMessage is an embedded status note inside a queue record.
type StatusMessage struct {
	Title    string   `json:"title"`
	Messages []string `json:"messages"`
}

// QueueRecord represents an item in the download queue.
type QueueRecord struct {
	ID                      int             `json:"id"`
	Title                   string          `json:"title"`
	Size                    float64         `json:"size"`
	SizeLeft                float64         `json:"sizeleft"`
	Status                  string          `json:"status"`
	TrackedDownloadStatus   string          `json:"trackedDownloadStatus"`
	TrackedDownloadState    string          `json:"trackedDownloadState"`
	StatusMessages          []StatusMessage `json:"statusMessages"`
	DownloadID              string          `json:"downloadId"`
	Protocol                string          `json:"protocol"`
	DownloadClient          string          `json:"downloadClient"`
	Indexer                 string          `json:"indexer"`
	OutputPath              string          `json:"outputPath"`
	TimeleftEstimation      time.Duration   `json:"-"`
	EstimatedCompletionTime string          `json:"estimatedCompletionTime"`
	Added                   string          `json:"added"`
}

// PagingResource wraps a paginated response from an *arr API.
type PagingResource[T any] struct {
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
	SortKey       string `json:"sortKey"`
	SortDirection string `json:"sortDirection"`
	TotalRecords  int    `json:"totalRecords"`
	Records       []T    `json:"records"`
}

// CommandRequest is the payload sent to /api/v3/command.
type CommandRequest struct {
	Name string `json:"name"`
}

// CommandResponse is the reply from /api/v3/command.
type CommandResponse struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Message             string `json:"message"`
	Started             string `json:"started"`
	Ended               string `json:"ended"`
	Status              string `json:"status"`
	Priority            string `json:"priority"`
	Trigger             string `json:"trigger"`
	StateChangeTime     string `json:"stateChangeTime"`
	SendUpdatesToClient bool   `json:"sendUpdatesToClient"`
	UpdateScheduledTask bool   `json:"updateScheduledTask"`
}

// QualityProfile describes a quality profile configured in the *arr app.
type QualityProfile struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	UpgradeAllowed    bool   `json:"upgradeAllowed"`
	CutoffFormatScore int    `json:"cutoffFormatScore"`
	MinFormatScore    int    `json:"minFormatScore"`
}

// Tag is a simple label used for organizing items in *arr applications.
type Tag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// UnmappedFolder represents a folder not yet mapped to a root folder.
type UnmappedFolder struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// RootFolder represents a configured root folder in an *arr application.
type RootFolder struct {
	ID              int              `json:"id"`
	Path            string           `json:"path"`
	Accessible      bool             `json:"accessible"`
	FreeSpace       int64            `json:"freeSpace"`
	UnmappedFolders []UnmappedFolder `json:"unmappedFolders"`
}

// DiskSpace contains disk usage information.
type DiskSpace struct {
	Path       string `json:"path"`
	Label      string `json:"label"`
	FreeSpace  int64  `json:"freeSpace"`
	TotalSpace int64  `json:"totalSpace"`
}

// TagDetail extends a [Tag] with lists of IDs that use it.
type TagDetail struct {
	ID    int    `json:"id"`
	Label string `json:"label"`

	DelayProfileIDs   []int `json:"delayProfileIds"`
	ImportListIDs     []int `json:"importListIds"`
	NotificationIDs   []int `json:"notificationIds"`
	RestrictionIDs    []int `json:"restrictionIds"`
	IndexerIDs        []int `json:"indexerIds"`
	DownloadClientIDs []int `json:"downloadClientIds"`
	AutoTagIDs        []int `json:"autoTagIds"`
	SeriesIDs         []int `json:"seriesIds,omitempty"` // Sonarr
	MovieIDs          []int `json:"movieIds,omitempty"`  // Radarr
	ArtistIDs         []int `json:"artistIds,omitempty"` // Lidarr
	AuthorIDs         []int `json:"authorIds,omitempty"` // Readarr
}

// Backup represents a backup file available on the server.
type Backup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	Time string `json:"time"`
}

// BlocklistResource represents a blocklisted release.
type BlocklistResource struct {
	ID            int                `json:"id"`
	SeriesID      int                `json:"seriesId,omitempty"`
	MovieID       int                `json:"movieId,omitempty"`
	ArtistID      int                `json:"artistId,omitempty"`
	AuthorID      int                `json:"authorId,omitempty"`
	SourceTitle   string             `json:"sourceTitle"`
	Languages     []LanguageResource `json:"languages"`
	Quality       any                `json:"quality"`
	CustomFormats []any              `json:"customFormats"`
	Date          string             `json:"date"`
	Protocol      string             `json:"protocol"`
	Indexer       string             `json:"indexer"`
	Message       string             `json:"message"`
}

// LanguageResource is a simple id+name pair for languages in shared contexts.
type LanguageResource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// BlocklistBulkResource is the request body for bulk blocklist operations.
type BlocklistBulkResource struct {
	IDs []int `json:"ids"`
}

// CustomFilterResource represents a custom UI filter.
type CustomFilterResource struct {
	ID      int                     `json:"id"`
	Type    string                  `json:"type"`
	Label   string                  `json:"label"`
	Filters []CustomFilterSpecifier `json:"filters"`
}

// CustomFilterSpecifier is a single filter criterion inside a [CustomFilterResource].
type CustomFilterSpecifier struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
	Type  string   `json:"type,omitempty"`
}

// CustomFormatResource describes a custom format with its specifications.
type CustomFormatResource struct {
	ID                              int                         `json:"id"`
	Name                            string                      `json:"name"`
	IncludeCustomFormatWhenRenaming bool                        `json:"includeCustomFormatWhenRenaming"`
	Specifications                  []CustomFormatSpecification `json:"specifications"`
}

// CustomFormatSpecification is one condition inside a [CustomFormatResource].
type CustomFormatSpecification struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name"`
	Implementation     string          `json:"implementation"`
	ImplementationName string          `json:"implementationName"`
	Negate             bool            `json:"negate"`
	Required           bool            `json:"required"`
	Fields             []ProviderField `json:"fields"`
}

// ProviderField is a key/value field used in provider configuration (Notification,
// DownloadClient, Indexer, ImportList, Metadata).
type ProviderField struct {
	Order         int            `json:"order"`
	Name          string         `json:"name"`
	Label         string         `json:"label"`
	Unit          string         `json:"unit,omitempty"`
	HelpText      string         `json:"helpText,omitempty"`
	HelpLink      string         `json:"helpLink,omitempty"`
	Value         any            `json:"value"`
	Type          string         `json:"type"`
	Advanced      bool           `json:"advanced"`
	SelectOptions []SelectOption `json:"selectOptions,omitempty"`
}

// SelectOption is a choice inside a [ProviderField] of type "select".
type SelectOption struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
	Order int    `json:"order"`
	Hint  string `json:"hint,omitempty"`
}

// ProviderResource is the shared shape for Notification, DownloadClient, Indexer,
// ImportList, and Metadata provider configurations in *arr APIs.
type ProviderResource struct {
	ID                 int              `json:"id"`
	Name               string           `json:"name"`
	Implementation     string           `json:"implementation"`
	ImplementationName string           `json:"implementationName"`
	ConfigContract     string           `json:"configContract"`
	InfoLink           string           `json:"infoLink,omitempty"`
	Tags               []int            `json:"tags"`
	Fields             []ProviderField  `json:"fields"`
	Message            *ProviderMessage `json:"message,omitempty"`

	// Notification-specific
	OnGrab                                *bool `json:"onGrab,omitempty"`
	OnDownload                            *bool `json:"onDownload,omitempty"`
	OnUpgrade                             *bool `json:"onUpgrade,omitempty"`
	OnRename                              *bool `json:"onRename,omitempty"`
	OnSeriesAdd                           *bool `json:"onSeriesAdd,omitempty"`
	OnSeriesDelete                        *bool `json:"onSeriesDelete,omitempty"`
	OnEpisodeFileDelete                   *bool `json:"onEpisodeFileDelete,omitempty"`
	OnEpisodeFileDeleteForUpgrade         *bool `json:"onEpisodeFileDeleteForUpgrade,omitempty"`
	OnHealthIssue                         *bool `json:"onHealthIssue,omitempty"`
	OnHealthRestored                      *bool `json:"onHealthRestored,omitempty"`
	OnApplicationUpdate                   *bool `json:"onApplicationUpdate,omitempty"`
	OnManualInteractionRequired           *bool `json:"onManualInteractionRequired,omitempty"`
	IncludeHealthWarnings                 *bool `json:"includeHealthWarnings,omitempty"`
	SupportsOnGrab                        *bool `json:"supportsOnGrab,omitempty"`
	SupportsOnDownload                    *bool `json:"supportsOnDownload,omitempty"`
	SupportsOnUpgrade                     *bool `json:"supportsOnUpgrade,omitempty"`
	SupportsOnRename                      *bool `json:"supportsOnRename,omitempty"`
	SupportsOnSeriesAdd                   *bool `json:"supportsOnSeriesAdd,omitempty"`
	SupportsOnSeriesDelete                *bool `json:"supportsOnSeriesDelete,omitempty"`
	SupportsOnEpisodeFileDelete           *bool `json:"supportsOnEpisodeFileDelete,omitempty"`
	SupportsOnEpisodeFileDeleteForUpgrade *bool `json:"supportsOnEpisodeFileDeleteForUpgrade,omitempty"`
	SupportsOnHealthIssue                 *bool `json:"supportsOnHealthIssue,omitempty"`
	SupportsOnHealthRestored              *bool `json:"supportsOnHealthRestored,omitempty"`
	SupportsOnApplicationUpdate           *bool `json:"supportsOnApplicationUpdate,omitempty"`
	SupportsOnManualInteractionRequired   *bool `json:"supportsOnManualInteractionRequired,omitempty"`

	// Radarr movie notifications
	OnMovieAdded                *bool `json:"onMovieAdded,omitempty"`
	OnMovieDelete               *bool `json:"onMovieDelete,omitempty"`
	OnMovieFileDelete           *bool `json:"onMovieFileDelete,omitempty"`
	OnMovieFileDeleteForUpgrade *bool `json:"onMovieFileDeleteForUpgrade,omitempty"`

	// DownloadClient-specific
	Enable                   *bool  `json:"enable,omitempty"`
	Protocol                 string `json:"protocol,omitempty"`
	Priority                 *int   `json:"priority,omitempty"`
	RemoveCompletedDownloads *bool  `json:"removeCompletedDownloads,omitempty"`
	RemoveFailedDownloads    *bool  `json:"removeFailedDownloads,omitempty"`

	// Indexer-specific
	EnableRss               *bool `json:"enableRss,omitempty"`
	EnableAutomaticSearch   *bool `json:"enableAutomaticSearch,omitempty"`
	EnableInteractiveSearch *bool `json:"enableInteractiveSearch,omitempty"`
	SupportsRss             *bool `json:"supportsRss,omitempty"`
	SupportsSearch          *bool `json:"supportsSearch,omitempty"`

	// ImportList-specific
	EnableAuto         *bool  `json:"enableAuto,omitempty"`
	ShouldMonitor      string `json:"shouldMonitor,omitempty"`
	RootFolderPath     string `json:"rootFolderPath,omitempty"`
	QualityProfileID   *int   `json:"qualityProfileId,omitempty"`
	ListType           string `json:"listType,omitempty"`
	ListOrder          *int   `json:"listOrder,omitempty"`
	MinRefreshInterval string `json:"minRefreshInterval,omitempty"`

	// Metadata-specific
	MetadataEnable *bool `json:"metadataEnable,omitempty"`
}

// ProviderMessage is a warning/info message attached to a provider.
type ProviderMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// ProviderBulkResource is used for bulk update operations on providers.
type ProviderBulkResource struct {
	IDs       []int  `json:"ids"`
	Tags      []int  `json:"tags,omitempty"`
	ApplyTags string `json:"applyTags,omitempty"`
}

// CustomFormatBulkResource is used for bulk operations on custom formats.
type CustomFormatBulkResource struct {
	IDs                             []int `json:"ids"`
	IncludeCustomFormatWhenRenaming *bool `json:"includeCustomFormatWhenRenaming,omitempty"`
}

// ProviderTestResource wraps a provider body for test API calls.
type ProviderTestResource = ProviderResource

// DelayProfileResource represents a delay profile configuration.
type DelayProfileResource struct {
	ID                             int    `json:"id"`
	EnableUsenet                   bool   `json:"enableUsenet"`
	EnableTorrent                  bool   `json:"enableTorrent"`
	PreferredProtocol              string `json:"preferredProtocol"`
	UsenetDelay                    int    `json:"usenetDelay"`
	TorrentDelay                   int    `json:"torrentDelay"`
	BypassIfHighestQuality         bool   `json:"bypassIfHighestQuality"`
	BypassIfAboveCustomFormatScore bool   `json:"bypassIfAboveCustomFormatScore"`
	MinimumCustomFormatScore       int    `json:"minimumCustomFormatScore"`
	Order                          int    `json:"order"`
	Tags                           []int  `json:"tags"`
}

// QualityDefinitionResource represents a quality size limit entry.
type QualityDefinitionResource struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Weight        int     `json:"weight"`
	MinSize       float64 `json:"minSize"`
	MaxSize       float64 `json:"maxSize"`
	PreferredSize float64 `json:"preferredSize"`
}

// ReleaseResource represents a release found during a search.
type ReleaseResource struct {
	GUID                     string  `json:"guid"`
	Quality                  any     `json:"quality"`
	CustomFormats            []any   `json:"customFormats"`
	CustomFormatScore        int     `json:"customFormatScore"`
	QualityWeight            int     `json:"qualityWeight"`
	Age                      int     `json:"age"`
	AgeHours                 float64 `json:"ageHours"`
	AgeMinutes               float64 `json:"ageMinutes"`
	Size                     int64   `json:"size"`
	IndexerID                int     `json:"indexerId"`
	Indexer                  string  `json:"indexer"`
	ReleaseGroup             string  `json:"releaseGroup"`
	SubGroup                 string  `json:"subGroup"`
	ReleaseHash              string  `json:"releaseHash"`
	Title                    string  `json:"title"`
	FullSeason               bool    `json:"fullSeason"`
	SceneSource              bool    `json:"sceneSource"`
	SeasonNumber             int     `json:"seasonNumber"`
	Languages                []any   `json:"languages"`
	MappedSeriesID           int     `json:"mappedSeriesId,omitempty"`
	MappedMovieID            int     `json:"mappedMovieId,omitempty"`
	Approved                 bool    `json:"approved"`
	TemporarilyRejected      bool    `json:"temporarilyRejected"`
	Rejected                 bool    `json:"rejected"`
	Rejections               []any   `json:"rejections"`
	PublishDate              string  `json:"publishDate"`
	CommentURL               string  `json:"commentUrl"`
	DownloadURL              string  `json:"downloadUrl"`
	InfoURL                  string  `json:"infoUrl"`
	DownloadAllowed          bool    `json:"downloadAllowed"`
	ReleaseWeight            int     `json:"releaseWeight"`
	Seeders                  int     `json:"seeders"`
	Leechers                 int     `json:"leechers"`
	Protocol                 string  `json:"protocol"`
	IsDaily                  bool    `json:"isDaily"`
	IsAbsoluteNumbering      bool    `json:"isAbsoluteNumbering"`
	IsPossibleSpecialEpisode bool    `json:"isPossibleSpecialEpisode"`
	Special                  bool    `json:"special"`
}

// ReleasePushResource is the body for manually pushing a release.
type ReleasePushResource struct {
	Title       string `json:"title"`
	DownloadURL string `json:"downloadUrl"`
	Protocol    string `json:"protocol"`
	PublishDate string `json:"publishDate"`
}

// ReleaseProfileResource configures preferred/required/ignored terms.
type ReleaseProfileResource struct {
	ID        int             `json:"id"`
	Name      string          `json:"name,omitempty"`
	Enabled   bool            `json:"enabled"`
	Required  []string        `json:"required"`
	Ignored   []string        `json:"ignored"`
	Preferred []PreferredTerm `json:"preferred"`
	IndexerID int             `json:"indexerId"`
	Tags      []int           `json:"tags"`
}

// PreferredTerm is a scored term inside a [ReleaseProfileResource].
type PreferredTerm struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

// RemotePathMappingResource maps a remote host's path to a local path.
type RemotePathMappingResource struct {
	ID         int    `json:"id"`
	Host       string `json:"host"`
	RemotePath string `json:"remotePath"`
	LocalPath  string `json:"localPath"`
}

// RenameEpisodeResource represents a proposed file rename.
type RenameEpisodeResource struct {
	EpisodeID      int    `json:"episodeId"`
	SeasonNumber   int    `json:"seasonNumber"`
	EpisodeNumbers []int  `json:"episodeNumbers"`
	ExistingPath   string `json:"existingPath"`
	NewPath        string `json:"newPath"`
}

// RenameMovieResource represents a proposed movie file rename.
type RenameMovieResource struct {
	MovieID      int    `json:"movieId"`
	MovieFileID  int    `json:"movieFileId"`
	ExistingPath string `json:"existingPath"`
	NewPath      string `json:"newPath"`
}

// ManualImportResource represents a file available for manual import.
type ManualImportResource struct {
	ID           int    `json:"id"`
	Path         string `json:"path"`
	RelativePath string `json:"relativePath"`
	FolderName   string `json:"folderName"`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	Quality      any    `json:"quality"`
	Languages    []any  `json:"languages"`
	ReleaseGroup string `json:"releaseGroup"`
	Rejections   []any  `json:"rejections"`
}

// ManualImportReprocessResource is the body for manual import confirmation.
type ManualImportReprocessResource struct {
	ID           int    `json:"id"`
	Path         string `json:"path"`
	SeriesID     int    `json:"seriesId,omitempty"`
	MovieID      int    `json:"movieId,omitempty"`
	SeasonNumber int    `json:"seasonNumber,omitempty"`
	EpisodeIDs   []int  `json:"episodeIds,omitempty"`
	Quality      any    `json:"quality"`
	Languages    []any  `json:"languages"`
	ReleaseGroup string `json:"releaseGroup"`
	DownloadID   string `json:"downloadId"`
}

// LogRecord is a single log entry returned by the /log endpoint.
type LogRecord struct {
	ID            int    `json:"id"`
	Time          string `json:"time"`
	Exception     string `json:"exception,omitempty"`
	ExceptionType string `json:"exceptionType,omitempty"`
	Level         string `json:"level"`
	Logger        string `json:"logger"`
	Message       string `json:"message"`
	Method        string `json:"method,omitempty"`
}

// LogFileResource describes an available log file.
type LogFileResource struct {
	ID            int    `json:"id"`
	Filename      string `json:"filename"`
	LastWriteTime string `json:"lastWriteTime"`
	ContentsURL   string `json:"contentsUrl"`
	DownloadURL   string `json:"downloadUrl"`
}

// NamingConfigResource holds the file naming configuration.
type NamingConfigResource struct {
	ID                       int    `json:"id"`
	RenameEpisodes           bool   `json:"renameEpisodes,omitempty"`
	ReplaceIllegalCharacters bool   `json:"replaceIllegalCharacters"`
	ColonReplacementFormat   string `json:"colonReplacementFormat"`
	StandardEpisodeFormat    string `json:"standardEpisodeFormat,omitempty"`
	DailyEpisodeFormat       string `json:"dailyEpisodeFormat,omitempty"`
	AnimeEpisodeFormat       string `json:"animeEpisodeFormat,omitempty"`
	SeriesFolderFormat       string `json:"seriesFolderFormat,omitempty"`
	SeasonFolderFormat       string `json:"seasonFolderFormat,omitempty"`
	SpecialsFolderFormat     string `json:"specialsFolderFormat,omitempty"`
	MultiEpisodeStyle        int    `json:"multiEpisodeStyle,omitempty"`
	// Radarr naming fields
	RenameMovies        bool   `json:"renameMovies,omitempty"`
	StandardMovieFormat string `json:"standardMovieFormat,omitempty"`
	MovieFolderFormat   string `json:"movieFolderFormat,omitempty"`
	// Lidarr naming fields
	RenameTracks         bool   `json:"renameTracks,omitempty"`
	StandardTrackFormat  string `json:"standardTrackFormat,omitempty"`
	MultiDiscTrackFormat string `json:"multiDiscTrackFormat,omitempty"`
	ArtistFolderFormat   string `json:"artistFolderFormat,omitempty"`
	// Readarr naming fields
	RenameBooks        bool   `json:"renameBooks,omitempty"`
	StandardBookFormat string `json:"standardBookFormat,omitempty"`
	AuthorFolderFormat string `json:"authorFolderFormat,omitempty"`
}

// HostConfigResource holds host/general configuration.
type HostConfigResource struct {
	ID                        int    `json:"id"`
	BindAddress               string `json:"bindAddress"`
	Port                      int    `json:"port"`
	SslPort                   int    `json:"sslPort"`
	EnableSsl                 bool   `json:"enableSsl"`
	LaunchBrowser             bool   `json:"launchBrowser"`
	AuthenticationMethod      string `json:"authenticationMethod"`
	AuthenticationRequired    string `json:"authenticationRequired"`
	AnalyticsEnabled          bool   `json:"analyticsEnabled"`
	Username                  string `json:"username,omitempty"`
	Password                  string `json:"password,omitempty"`
	PasswordConfirmation      string `json:"passwordConfirmation,omitempty"`
	LogLevel                  string `json:"logLevel"`
	ConsoleLogLevel           string `json:"consoleLogLevel"`
	Branch                    string `json:"branch"`
	ApiKey                    string `json:"apiKey"`
	SslCertPath               string `json:"sslCertPath"`
	SslCertPassword           string `json:"sslCertPassword"`
	URLBase                   string `json:"urlBase"`
	InstanceName              string `json:"instanceName"`
	UpdateAutomatically       bool   `json:"updateAutomatically"`
	UpdateMechanism           string `json:"updateMechanism"`
	UpdateScriptPath          string `json:"updateScriptPath"`
	ProxyEnabled              bool   `json:"proxyEnabled"`
	ProxyType                 string `json:"proxyType"`
	ProxyHostname             string `json:"proxyHostname"`
	ProxyPort                 int    `json:"proxyPort"`
	ProxyUsername             string `json:"proxyUsername,omitempty"`
	ProxyPassword             string `json:"proxyPassword,omitempty"`
	ProxyBypassFilter         string `json:"proxyBypassFilter"`
	ProxyBypassLocalAddresses bool   `json:"proxyBypassLocalAddresses"`
	CertificateValidation     string `json:"certificateValidation"`
	BackupFolder              string `json:"backupFolder"`
	BackupInterval            int    `json:"backupInterval"`
	BackupRetention           int    `json:"backupRetention"`
}

// UIConfigResource holds UI settings.
type UIConfigResource struct {
	ID                       int    `json:"id"`
	FirstDayOfWeek           int    `json:"firstDayOfWeek"`
	CalendarWeekColumnHeader string `json:"calendarWeekColumnHeader"`
	ShortDateFormat          string `json:"shortDateFormat"`
	LongDateFormat           string `json:"longDateFormat"`
	TimeFormat               string `json:"timeFormat"`
	ShowRelativeDates        bool   `json:"showRelativeDates"`
	EnableColorImpairedMode  bool   `json:"enableColorImpairedMode"`
	Theme                    string `json:"theme"`
	UILanguage               int    `json:"uiLanguage"`
	MovieInfoLanguage        int    `json:"movieInfoLanguage,omitempty"`
	MovieRuntimeFormat       string `json:"movieRuntimeFormat,omitempty"`
}

// MediaManagementConfigResource holds media management settings.
type MediaManagementConfigResource struct {
	ID                                        int    `json:"id"`
	AutoUnmonitorPreviouslyDownloadedEpisodes bool   `json:"autoUnmonitorPreviouslyDownloadedEpisodes,omitempty"`
	AutoUnmonitorPreviouslyDownloadedMovies   bool   `json:"autoUnmonitorPreviouslyDownloadedMovies,omitempty"`
	RecycleBin                                string `json:"recycleBin"`
	RecycleBinCleanupDays                     int    `json:"recycleBinCleanupDays"`
	DownloadPropersAndRepacks                 string `json:"downloadPropersAndRepacks"`
	CreateEmptySeriesFolders                  bool   `json:"createEmptySeriesFolders,omitempty"`
	CreateEmptyMovieFolders                   bool   `json:"createEmptyMovieFolders,omitempty"`
	CreateEmptyArtistFolders                  bool   `json:"createEmptyArtistFolders,omitempty"`
	CreateEmptyAuthorFolders                  bool   `json:"createEmptyAuthorFolders,omitempty"`
	DeleteEmptyFolders                        bool   `json:"deleteEmptyFolders"`
	FileDate                                  string `json:"fileDate"`
	RescanAfterRefresh                        string `json:"rescanAfterRefresh"`
	AutoRenameFolders                         bool   `json:"autoRenameFolders"`
	PathsDefaultStatic                        bool   `json:"pathsDefaultStatic"`
	SetPermissionsLinux                       bool   `json:"setPermissionsLinux"`
	ChmodFolder                               string `json:"chmodFolder"`
	ChownGroup                                string `json:"chownGroup"`
	SkipFreeSpaceCheckWhenImporting           bool   `json:"skipFreeSpaceCheckWhenImporting"`
	MinimumFreeSpaceWhenImporting             int    `json:"minimumFreeSpaceWhenImporting"`
	CopyUsingHardlinks                        bool   `json:"copyUsingHardlinks"`
	UseScriptImport                           bool   `json:"useScriptImport"`
	ScriptImportPath                          string `json:"scriptImportPath,omitempty"`
	ImportExtraFiles                          bool   `json:"importExtraFiles"`
	ExtraFileExtensions                       string `json:"extraFileExtensions"`
	EnableMediaInfo                           bool   `json:"enableMediaInfo"`
}

// DownloadClientConfigResource holds download client global settings.
type DownloadClientConfigResource struct {
	ID                                        int    `json:"id"`
	DownloadClientWorkingFolders              string `json:"downloadClientWorkingFolders"`
	EnableCompletedDownloadHandling           bool   `json:"enableCompletedDownloadHandling"`
	AutoRedownloadFailed                      bool   `json:"autoRedownloadFailed"`
	AutoRedownloadFailedFromInteractiveSearch bool   `json:"autoRedownloadFailedFromInteractiveSearch"`
}

// IndexerConfigResource holds indexer global settings.
type IndexerConfigResource struct {
	ID              int `json:"id"`
	MinimumAge      int `json:"minimumAge"`
	Retention       int `json:"retention"`
	MaximumSize     int `json:"maximumSize"`
	RssSyncInterval int `json:"rssSyncInterval"`
}

// IndexerFlagResource describes an indexer flag.
type IndexerFlagResource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TaskResource represents a scheduled task.
type TaskResource struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	TaskName      string `json:"taskName"`
	Interval      int    `json:"interval"`
	LastExecution string `json:"lastExecution"`
	LastStartTime string `json:"lastStartTime"`
	NextExecution string `json:"nextExecution"`
	LastDuration  string `json:"lastDuration"`
}

// UpdateResource represents an available application update.
type UpdateResource struct {
	ID          int    `json:"id"`
	Version     string `json:"version"`
	Branch      string `json:"branch"`
	ReleaseDate string `json:"releaseDate"`
	FileName    string `json:"fileName"`
	URL         string `json:"url"`
	Installed   bool   `json:"installed"`
	InstalledOn string `json:"installedOn,omitempty"`
	Installable bool   `json:"installable"`
	Latest      bool   `json:"latest"`
	Changes     any    `json:"changes"`
	Hash        string `json:"hash"`
}

// QueueStatusResource reports the overall queue status.
type QueueStatusResource struct {
	TotalCount      int  `json:"totalCount"`
	Count           int  `json:"count"`
	UnknownCount    int  `json:"unknownCount"`
	Errors          bool `json:"errors"`
	Warnings        bool `json:"warnings"`
	UnknownErrors   bool `json:"unknownErrors"`
	UnknownWarnings bool `json:"unknownWarnings"`
}

// QueueBulkResource is the request body for bulk queue operations.
type QueueBulkResource struct {
	IDs []int `json:"ids"`
}

// AutoTaggingResource represents an auto-tagging rule.
type AutoTaggingResource struct {
	ID                      int                        `json:"id"`
	Name                    string                     `json:"name"`
	RemoveTagsAutomatically bool                       `json:"removeTagsAutomatically"`
	Tags                    []int                      `json:"tags"`
	Specifications          []AutoTaggingSpecification `json:"specifications"`
}

// AutoTaggingSpecification is a condition in an auto-tagging rule.
type AutoTaggingSpecification struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name"`
	Implementation     string          `json:"implementation"`
	ImplementationName string          `json:"implementationName"`
	Negate             bool            `json:"negate"`
	Required           bool            `json:"required"`
	Fields             []ProviderField `json:"fields"`
}

// ImportListExclusionResource represents an import list exclusion entry.
type ImportListExclusionResource struct {
	ID     int    `json:"id"`
	TvdbID int    `json:"tvdbId,omitempty"` // Sonarr
	TmdbID int    `json:"tmdbId,omitempty"` // Radarr
	Title  string `json:"title"`
	Year   int    `json:"year,omitempty"`
}

// SystemRouteResource describes an API route.
type SystemRouteResource struct {
	Path     string `json:"path"`
	Method   string `json:"httpMethod"`
	IsPublic bool   `json:"isPublic"`
	IsDebug  bool   `json:"isDebug"`
}

// CalendarFeedParams holds optional parameters for the iCal calendar feed.
type CalendarFeedParams struct {
	PastDays    int
	FutureDays  int
	Unmonitored bool
	Tags        []int
}
