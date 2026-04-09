package autobrr

import "encoding/json"

// Filter represents an autobrr filter.
type Filter struct {
	ID                int              `json:"id"`
	Name              string           `json:"name"`
	Enabled           bool             `json:"enabled"`
	Priority          int              `json:"priority,omitempty"`
	UseRegex          bool             `json:"use_regex,omitempty"`
	Years             string           `json:"years,omitempty"`
	Resolutions       []string         `json:"resolutions,omitempty"`
	Sources           []string         `json:"sources,omitempty"`
	Codecs            []string         `json:"codecs,omitempty"`
	Containers        []string         `json:"containers,omitempty"`
	MatchHDR          []string         `json:"match_hdr,omitempty"`
	ExceptHDR         []string         `json:"except_hdr,omitempty"`
	MatchOther        []string         `json:"match_other,omitempty"`
	ExceptOther       []string         `json:"except_other,omitempty"`
	MatchReleases     string           `json:"match_releases,omitempty"`
	ExceptReleases    string           `json:"except_releases,omitempty"`
	Tags              string           `json:"tags,omitempty"`
	ExceptTags        string           `json:"except_tags,omitempty"`
	MatchLanguage     []string         `json:"match_language,omitempty"`
	ExceptLanguage    []string         `json:"except_language,omitempty"`
	Formats           []string         `json:"formats,omitempty"`
	Quality           []string         `json:"quality,omitempty"`
	Media             []string         `json:"media,omitempty"`
	MatchReleaseTypes []string         `json:"match_release_types,omitempty"`
	Origins           []string         `json:"origins,omitempty"`
	ExceptOrigins     []string         `json:"except_origins,omitempty"`
	SmartEpisode      bool             `json:"smart_episode,omitempty"`
	Indexers          []FilterIndexer  `json:"indexers,omitempty"`
	Actions           []FilterAction   `json:"actions,omitempty"`
	External          []FilterExternal `json:"external,omitempty"`
}

// FilterIndexer is an indexer entry within a filter.
type FilterIndexer struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}

// FilterAction is an action configured on a filter.
type FilterAction struct {
	Name                  string `json:"name"`
	Type                  string `json:"type"`
	Enabled               bool   `json:"enabled"`
	Category              string `json:"category,omitempty"`
	Tags                  string `json:"tags,omitempty"`
	ClientID              int    `json:"client_id,omitempty"`
	ReannounceInterval    int    `json:"reannounce_interval,omitempty"`
	ReannounceMaxAttempts int    `json:"reannounce_max_attempts,omitempty"`
}

// FilterExternal is an external action on a filter.
type FilterExternal struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Index               int    `json:"index"`
	Type                string `json:"type"`
	Enabled             bool   `json:"enabled"`
	WebhookHost         string `json:"webhook_host,omitempty"`
	WebhookMethod       string `json:"webhook_method,omitempty"`
	WebhookData         string `json:"webhook_data,omitempty"`
	WebhookExpectStatus int    `json:"webhook_expect_status,omitempty"`
}

// Indexer represents an autobrr indexer.
type Indexer struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier,omitempty"`
	Enabled    bool   `json:"enabled"`
	URL        string `json:"url,omitempty"`
}

// IRCNetwork represents an IRC network configuration.
type IRCNetwork struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Healthy bool   `json:"healthy"`
	Enabled bool   `json:"enabled,omitempty"`
	Server  string `json:"server,omitempty"`
	Port    int    `json:"port,omitempty"`
	TLS     bool   `json:"tls,omitempty"`
	Nick    string `json:"nick,omitempty"`
}

// Feed represents an RSS/torznab feed.
type Feed struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	URL     string `json:"url,omitempty"`
	Type    string `json:"type,omitempty"`
}

// DownloadClient represents a configured download client.
type DownloadClient struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	Enabled       bool            `json:"enabled"`
	Host          string          `json:"host,omitempty"`
	Port          int             `json:"port,omitempty"`
	TLS           bool            `json:"tls,omitempty"`
	TLSSkipVerify bool            `json:"tls_skip_verify,omitempty"`
	Username      string          `json:"username,omitempty"`
	Password      string          `json:"password,omitempty"`
	Settings      json.RawMessage `json:"settings,omitempty"`
}

// Notification represents a notification agent.
type Notification struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Enabled bool     `json:"enabled"`
	Events  []string `json:"events,omitempty"`
	Webhook string   `json:"webhook,omitempty"`
	APIKey  string   `json:"api_key,omitempty"`
	Token   string   `json:"token,omitempty"`
	Channel string   `json:"channel,omitempty"`
}

// APIKey represents an API key entry.
type APIKey struct {
	Name   string   `json:"name"`
	Key    string   `json:"key,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}

// Config represents the autobrr configuration.
type Config struct {
	Host            string `json:"host,omitempty"`
	Port            int    `json:"port,omitempty"`
	LogLevel        string `json:"log_level,omitempty"`
	LogPath         string `json:"log_path,omitempty"`
	CheckForUpdates bool   `json:"check_for_updates,omitempty"`
	Version         string `json:"version,omitempty"`
}

// ConfigUpdate represents fields that can be patched on the configuration.
type ConfigUpdate struct {
	Host            *string `json:"host,omitempty"`
	Port            *int    `json:"port,omitempty"`
	LogLevel        *string `json:"log_level,omitempty"`
	LogPath         *string `json:"log_path,omitempty"`
	CheckForUpdates *bool   `json:"check_for_updates,omitempty"`
}

// Action represents an autobrr action rule.
type Action struct {
	ID                       int      `json:"id"`
	Name                     string   `json:"name"`
	Type                     string   `json:"type"`
	Enabled                  bool     `json:"enabled"`
	FilterID                 int      `json:"filter_id,omitempty"`
	Category                 string   `json:"category,omitempty"`
	Tags                     string   `json:"tags,omitempty"`
	Label                    string   `json:"label,omitempty"`
	SavePath                 string   `json:"save_path,omitempty"`
	Paused                   bool     `json:"paused,omitempty"`
	IgnoreRules              bool     `json:"ignore_rules,omitempty"`
	SkipHashCheck            bool     `json:"skip_hash_check,omitempty"`
	ContentLayout            string   `json:"content_layout,omitempty"`
	LimitUploadSpeed         int64    `json:"limit_upload_speed,omitempty"`
	LimitDownloadSpeed       int64    `json:"limit_download_speed,omitempty"`
	LimitRatio               float64  `json:"limit_ratio,omitempty"`
	LimitSeedTime            int64    `json:"limit_seed_time,omitempty"`
	ClientID                 int      `json:"client_id,omitempty"`
	ReannounceInterval       int      `json:"reannounce_interval,omitempty"`
	ReannounceMaxAttempts    int      `json:"reannounce_max_attempts,omitempty"`
	WebhookHost              string   `json:"webhook_host,omitempty"`
	WebhookType              string   `json:"webhook_type,omitempty"`
	WebhookMethod            string   `json:"webhook_method,omitempty"`
	WebhookData              string   `json:"webhook_data,omitempty"`
	WebhookHeaders           []string `json:"webhook_headers,omitempty"`
	ExternalDownloadClientID int      `json:"external_download_client_id,omitempty"`
	ExternalDownloadClient   string   `json:"external_download_client,omitempty"`
}

// List represents an autobrr list (e.g. Trakt, MDBList, Steam, etc.).
type List struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Type               string `json:"type"`
	Enabled            bool   `json:"enabled"`
	URL                string `json:"url,omitempty"`
	APIKey             string `json:"api_key,omitempty"`
	MatchRelease       bool   `json:"match_release,omitempty"`
	TagsInclude        string `json:"tags_include,omitempty"`
	TagsExclude        string `json:"tags_exclude,omitempty"`
	IncludeUnmonitored bool   `json:"include_unmonitored,omitempty"`
	Filters            []int  `json:"filters,omitempty"`
	ClientID           int    `json:"client_id,omitempty"`
}

// Proxy represents a proxy configuration.
type Proxy struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
	Addr    string `json:"addr,omitempty"`
	User    string `json:"user,omitempty"`
	Pass    string `json:"pass,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}

// Release represents a release tracked by autobrr.
type Release struct {
	ID             int                   `json:"id"`
	FilterStatus   string                `json:"filter_status,omitempty"`
	Rejections     []string              `json:"rejections,omitempty"`
	Indexer        ReleaseIndexer        `json:"indexer,omitempty"`
	FilterName     string                `json:"filter_name,omitempty"`
	Protocol       string                `json:"protocol,omitempty"`
	Implementation string                `json:"implementation,omitempty"`
	Timestamp      string                `json:"timestamp,omitempty"`
	Name           string                `json:"name,omitempty"`
	Size           int64                 `json:"size,omitempty"`
	Title          string                `json:"title,omitempty"`
	Category       string                `json:"category,omitempty"`
	Season         int                   `json:"season,omitempty"`
	Episode        int                   `json:"episode,omitempty"`
	Year           int                   `json:"year,omitempty"`
	Resolution     string                `json:"resolution,omitempty"`
	Source         string                `json:"source,omitempty"`
	Codec          string                `json:"codec,omitempty"`
	Container      string                `json:"container,omitempty"`
	HDR            string                `json:"hdr,omitempty"`
	Group          string                `json:"group,omitempty"`
	ActionStatus   []ReleaseActionStatus `json:"action_status,omitempty"`
}

// ReleaseIndexer identifies the indexer for a release.
type ReleaseIndexer struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}

// ReleaseActionStatus represents the push status of an action on a release.
type ReleaseActionStatus struct {
	ID         int      `json:"id"`
	Status     string   `json:"status"`
	Action     string   `json:"action"`
	ActionID   int      `json:"action_id"`
	Type       string   `json:"type"`
	Client     string   `json:"client,omitempty"`
	Filter     string   `json:"filter,omitempty"`
	Rejections []string `json:"rejections,omitempty"`
	Timestamp  string   `json:"timestamp,omitempty"`
}

// ReleaseFindResponse is the paginated response from the release endpoint.
type ReleaseFindResponse struct {
	Data       []Release `json:"data"`
	NextCursor int       `json:"next_cursor"`
	Count      int       `json:"count"`
}

// ReleaseStats contains aggregate release statistics.
type ReleaseStats struct {
	TotalCount          int `json:"total_count"`
	FilteredCount       int `json:"filtered_count"`
	FilterRejectedCount int `json:"filter_rejected_count"`
	PushApprovedCount   int `json:"push_approved_count"`
	PushRejectedCount   int `json:"push_rejected_count"`
	PushErrorCount      int `json:"push_error_count"`
}

// ReleaseDeleteParams holds query parameters for bulk-deleting releases.
type ReleaseDeleteParams struct {
	OlderThan       int      `json:"olderThan,omitempty"`
	Indexers        []string `json:"indexer,omitempty"`
	ReleaseStatuses []string `json:"releaseStatus,omitempty"`
}

// FilterNotification links a notification agent to a filter.
type FilterNotification struct {
	ID             int  `json:"id"`
	FilterID       int  `json:"filter_id"`
	NotificationID int  `json:"notification_id"`
	Enabled        bool `json:"enabled"`
}

// SendIRCCmdRequest is the payload for sending an IRC command.
type SendIRCCmdRequest struct {
	NetworkID int    `json:"network_id"`
	Server    string `json:"server,omitempty"`
	Channel   string `json:"channel"`
	Message   string `json:"msg"`
}

// LogFile represents a log file entry.
type LogFile struct {
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// LogFileResponse is the response from the logs endpoint.
type LogFileResponse struct {
	Files []LogFile `json:"files"`
	Count int       `json:"count"`
}

// IndexerDefinition is a full indexer definition (from /indexer/schema).
type IndexerDefinition struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	Identifier     string           `json:"identifier,omitempty"`
	Implementation string           `json:"implementation,omitempty"`
	BaseURL        string           `json:"base_url,omitempty"`
	Enabled        bool             `json:"enabled"`
	Description    string           `json:"description,omitempty"`
	Language       string           `json:"language,omitempty"`
	Privacy        string           `json:"privacy,omitempty"`
	Protocol       string           `json:"protocol,omitempty"`
	URLs           []string         `json:"urls,omitempty"`
	Settings       []IndexerSetting `json:"settings,omitempty"`
}

// IndexerSetting is a configurable field on an indexer definition.
type IndexerSetting struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value,omitempty"`
	Label string `json:"label,omitempty"`
	Help  string `json:"help,omitempty"`
}

// ReleaseCleanupJob represents a scheduled release cleanup job.
type ReleaseCleanupJob struct {
	ID              int      `json:"id"`
	Enabled         bool     `json:"enabled"`
	Name            string   `json:"name,omitempty"`
	OlderThan       int      `json:"older_than,omitempty"`
	Indexers        []string `json:"indexers,omitempty"`
	ReleaseStatuses []string `json:"release_statuses,omitempty"`
}

// ReleaseProfileDuplicate represents a duplicate release profile.
type ReleaseProfileDuplicate struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// GithubRelease represents a GitHub release (from /updates/latest).
type GithubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
}

// ArrTag represents a tag from an *arr download client.
type ArrTag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}
