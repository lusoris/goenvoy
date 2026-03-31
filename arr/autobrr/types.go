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
