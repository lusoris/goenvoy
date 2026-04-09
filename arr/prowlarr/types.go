package prowlarr

// Indexer represents a configured indexer in Prowlarr.
type Indexer struct {
	ID                 int                  `json:"id"`
	Name               string               `json:"name"`
	Fields             []Field              `json:"fields,omitempty"`
	ImplementationName string               `json:"implementationName,omitempty"`
	Implementation     string               `json:"implementation,omitempty"`
	ConfigContract     string               `json:"configContract,omitempty"`
	InfoLink           string               `json:"infoLink,omitempty"`
	Message            *ProviderMessage     `json:"message,omitempty"`
	Tags               []int                `json:"tags,omitempty"`
	IndexerURLs        []string             `json:"indexerUrls,omitempty"`
	LegacyURLs         []string             `json:"legacyUrls,omitempty"`
	DefinitionName     string               `json:"definitionName,omitempty"`
	Description        string               `json:"description,omitempty"`
	Language           string               `json:"language,omitempty"`
	Encoding           string               `json:"encoding,omitempty"`
	Enable             bool                 `json:"enable"`
	Redirect           bool                 `json:"redirect"`
	SupportsRss        bool                 `json:"supportsRss"`
	SupportsSearch     bool                 `json:"supportsSearch"`
	SupportsRedirect   bool                 `json:"supportsRedirect"`
	SupportsPagination bool                 `json:"supportsPagination"`
	AppProfileID       int                  `json:"appProfileId"`
	Protocol           string               `json:"protocol,omitempty"`
	Privacy            string               `json:"privacy,omitempty"`
	Capabilities       *IndexerCapability   `json:"capabilities,omitempty"`
	Priority           int                  `json:"priority"`
	DownloadClientID   int                  `json:"downloadClientId,omitempty"`
	Added              string               `json:"added,omitempty"`
	Status             *IndexerStatusDetail `json:"status,omitempty"`
	SortName           string               `json:"sortName,omitempty"`
}

// Application represents a connected PVR application (Sonarr, Radarr, etc.).
type Application struct {
	ID                 int              `json:"id"`
	Name               string           `json:"name"`
	Fields             []Field          `json:"fields,omitempty"`
	ImplementationName string           `json:"implementationName,omitempty"`
	Implementation     string           `json:"implementation,omitempty"`
	ConfigContract     string           `json:"configContract,omitempty"`
	InfoLink           string           `json:"infoLink,omitempty"`
	Message            *ProviderMessage `json:"message,omitempty"`
	Tags               []int            `json:"tags,omitempty"`
	SyncLevel          string           `json:"syncLevel,omitempty"`
}

// AppProfile defines search behavior settings for indexers.
type AppProfile struct {
	ID                      int    `json:"id"`
	Name                    string `json:"name"`
	EnableRss               bool   `json:"enableRss"`
	EnableAutomaticSearch   bool   `json:"enableAutomaticSearch"`
	EnableInteractiveSearch bool   `json:"enableInteractiveSearch"`
	MinimumSeeders          int    `json:"minimumSeeders"`
}

// Release represents a search result from an indexer.
type Release struct {
	ID               int               `json:"id"`
	GUID             string            `json:"guid,omitempty"`
	Age              int               `json:"age"`
	AgeHours         float64           `json:"ageHours"`
	AgeMinutes       float64           `json:"ageMinutes"`
	Size             int64             `json:"size"`
	Files            *int              `json:"files,omitempty"`
	Grabs            *int              `json:"grabs,omitempty"`
	IndexerID        int               `json:"indexerId"`
	Indexer          string            `json:"indexer,omitempty"`
	Title            string            `json:"title,omitempty"`
	SortTitle        string            `json:"sortTitle,omitempty"`
	ImdbID           int               `json:"imdbId,omitempty"`
	TmdbID           int               `json:"tmdbId,omitempty"`
	TvdbID           int               `json:"tvdbId,omitempty"`
	PublishDate      string            `json:"publishDate,omitempty"`
	CommentURL       string            `json:"commentUrl,omitempty"`
	DownloadURL      string            `json:"downloadUrl,omitempty"`
	InfoURL          string            `json:"infoUrl,omitempty"`
	PosterURL        string            `json:"posterUrl,omitempty"`
	IndexerFlags     []string          `json:"indexerFlags,omitempty"`
	Categories       []IndexerCategory `json:"categories,omitempty"`
	MagnetURL        string            `json:"magnetUrl,omitempty"`
	InfoHash         string            `json:"infoHash,omitempty"`
	Seeders          *int              `json:"seeders,omitempty"`
	Leechers         *int              `json:"leechers,omitempty"`
	Protocol         string            `json:"protocol,omitempty"`
	FileName         string            `json:"fileName,omitempty"`
	DownloadClientID *int              `json:"downloadClientId,omitempty"`
}

// HistoryRecord represents an indexer event in Prowlarr history.
type HistoryRecord struct {
	ID         int               `json:"id"`
	IndexerID  int               `json:"indexerId"`
	Date       string            `json:"date"`
	DownloadID string            `json:"downloadId,omitempty"`
	Successful bool              `json:"successful"`
	EventType  string            `json:"eventType"`
	Data       map[string]string `json:"data,omitempty"`
}

// IndexerStats contains aggregated statistics for all indexers.
type IndexerStats struct {
	ID         int                  `json:"id"`
	Indexers   []IndexerStatistic   `json:"indexers,omitempty"`
	UserAgents []UserAgentStatistic `json:"userAgents,omitempty"`
	Hosts      []HostStatistic      `json:"hosts,omitempty"`
}

// IndexerStatistic contains usage statistics for a single indexer.
type IndexerStatistic struct {
	IndexerID                 int    `json:"indexerId"`
	IndexerName               string `json:"indexerName,omitempty"`
	AverageResponseTime       int    `json:"averageResponseTime"`
	AverageGrabResponseTime   int    `json:"averageGrabResponseTime"`
	NumberOfQueries           int    `json:"numberOfQueries"`
	NumberOfGrabs             int    `json:"numberOfGrabs"`
	NumberOfRssQueries        int    `json:"numberOfRssQueries"`
	NumberOfAuthQueries       int    `json:"numberOfAuthQueries"`
	NumberOfFailedQueries     int    `json:"numberOfFailedQueries"`
	NumberOfFailedGrabs       int    `json:"numberOfFailedGrabs"`
	NumberOfFailedRssQueries  int    `json:"numberOfFailedRssQueries"`
	NumberOfFailedAuthQueries int    `json:"numberOfFailedAuthQueries"`
}

// UserAgentStatistic contains query/grab stats for a user agent.
type UserAgentStatistic struct {
	UserAgent       string `json:"userAgent,omitempty"`
	NumberOfQueries int    `json:"numberOfQueries"`
	NumberOfGrabs   int    `json:"numberOfGrabs"`
}

// HostStatistic contains query/grab stats for a host.
type HostStatistic struct {
	Host            string `json:"host,omitempty"`
	NumberOfQueries int    `json:"numberOfQueries"`
	NumberOfGrabs   int    `json:"numberOfGrabs"`
}

// IndexerStatus represents the current status of an indexer.
type IndexerStatus struct {
	ID                int    `json:"id"`
	IndexerID         int    `json:"indexerId"`
	DisabledTill      string `json:"disabledTill,omitempty"`
	MostRecentFailure string `json:"mostRecentFailure,omitempty"`
	InitialFailure    string `json:"initialFailure,omitempty"`
}

// IndexerStatusDetail is the inline status embedded in an IndexerResource.
type IndexerStatusDetail struct {
	ID                int    `json:"id"`
	IndexerID         int    `json:"indexerId"`
	DisabledTill      string `json:"disabledTill,omitempty"`
	MostRecentFailure string `json:"mostRecentFailure,omitempty"`
	InitialFailure    string `json:"initialFailure,omitempty"`
}

// IndexerCapability describes the search capabilities of an indexer.
type IndexerCapability struct {
	ID                int               `json:"id"`
	LimitsMax         *int              `json:"limitsMax,omitempty"`
	LimitsDefault     *int              `json:"limitsDefault,omitempty"`
	Categories        []IndexerCategory `json:"categories,omitempty"`
	SupportsRawSearch bool              `json:"supportsRawSearch"`
}

// IndexerCategory represents a Newznab/Torznab category.
type IndexerCategory struct {
	ID            int               `json:"id"`
	Name          string            `json:"name,omitempty"`
	Description   string            `json:"description,omitempty"`
	SubCategories []IndexerCategory `json:"subCategories,omitempty"`
}

// Field represents a configurable field for a provider (indexer, application, etc.).
type Field struct {
	Order    int    `json:"order"`
	Name     string `json:"name,omitempty"`
	Label    string `json:"label,omitempty"`
	HelpText string `json:"helpText,omitempty"`
	Value    any    `json:"value,omitempty"`
	Type     string `json:"type,omitempty"`
	Advanced bool   `json:"advanced"`
}

// ProviderMessage contains informational or warning messages from a provider.
type ProviderMessage struct {
	Message string `json:"message,omitempty"`
	Type    string `json:"type,omitempty"`
}

// DownloadClientResource represents a download client configured in Prowlarr.
type DownloadClientResource struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	Fields             []Field `json:"fields,omitempty"`
	ImplementationName string  `json:"implementationName,omitempty"`
	Implementation     string  `json:"implementation,omitempty"`
	ConfigContract     string  `json:"configContract,omitempty"`
	Tags               []int   `json:"tags,omitempty"`
	Enable             bool    `json:"enable"`
	Protocol           string  `json:"protocol,omitempty"`
	Priority           int     `json:"priority"`
}

// DevelopmentConfigResource represents Prowlarr development configuration.
type DevelopmentConfigResource struct {
	ID                 int    `json:"id"`
	ConsoleLogLevel    string `json:"consoleLogLevel,omitempty"`
	LogSql             bool   `json:"logSql"`
	LogIndexerResponse bool   `json:"logIndexerResponse"`
	LogRotate          int    `json:"logRotate"`
	FilterSentryEvents bool   `json:"filterSentryEvents"`
}

// IndexerProxyResource represents a configured indexer proxy in Prowlarr.
type IndexerProxyResource struct {
	ID                    int                    `json:"id"`
	Name                  string                 `json:"name,omitempty"`
	Fields                []Field                `json:"fields,omitempty"`
	ImplementationName    string                 `json:"implementationName,omitempty"`
	Implementation        string                 `json:"implementation,omitempty"`
	ConfigContract        string                 `json:"configContract,omitempty"`
	InfoLink              string                 `json:"infoLink,omitempty"`
	Message               *ProviderMessage       `json:"message,omitempty"`
	Tags                  []int                  `json:"tags,omitempty"`
	Presets               []IndexerProxyResource `json:"presets,omitempty"`
	Link                  string                 `json:"link,omitempty"`
	OnHealthIssue         bool                   `json:"onHealthIssue"`
	SupportsOnHealthIssue bool                   `json:"supportsOnHealthIssue"`
	IncludeHealthWarnings bool                   `json:"includeHealthWarnings"`
	TestCommand           string                 `json:"testCommand,omitempty"`
}

// LocalizationOption represents a localization language option.
type LocalizationOption struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// ApplicationBulkResource represents a bulk update for applications.
type ApplicationBulkResource struct {
	IDs       []int  `json:"ids,omitempty"`
	Tags      []int  `json:"tags,omitempty"`
	ApplyTags string `json:"applyTags,omitempty"`
	SyncLevel string `json:"syncLevel,omitempty"`
}

// IndexerBulkResource represents a bulk update for indexers.
type IndexerBulkResource struct {
	IDs             []int    `json:"ids,omitempty"`
	Tags            []int    `json:"tags,omitempty"`
	ApplyTags       string   `json:"applyTags,omitempty"`
	Enable          *bool    `json:"enable,omitempty"`
	AppProfileID    *int     `json:"appProfileId,omitempty"`
	Priority        *int     `json:"priority,omitempty"`
	MinimumSeeders  *int     `json:"minimumSeeders,omitempty"`
	SeedRatio       *float64 `json:"seedRatio,omitempty"`
	SeedTime        *int     `json:"seedTime,omitempty"`
	PackSeedTime    *int     `json:"packSeedTime,omitempty"`
	PreferMagnetURL *bool    `json:"preferMagnetUrl,omitempty"`
}
