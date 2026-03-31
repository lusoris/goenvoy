package tautulli

// Response is the top-level wrapper for all Tautulli API responses.
type Response struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

// Activity represents the current server activity.
type Activity struct {
	StreamCount             string    `json:"stream_count"`
	StreamCountDirectPlay   int       `json:"stream_count_direct_play"`
	StreamCountDirectStream int       `json:"stream_count_direct_stream"`
	StreamCountTranscode    int       `json:"stream_count_transcode"`
	TotalBandwidth          int       `json:"total_bandwidth"`
	LANBandwidth            int       `json:"lan_bandwidth"`
	WANBandwidth            int       `json:"wan_bandwidth"`
	Sessions                []Session `json:"sessions"`
}

// Session represents an active playback session.
type Session struct {
	SessionKey           string `json:"session_key"`
	SessionID            string `json:"session_id"`
	MediaType            string `json:"media_type"`
	RatingKey            string `json:"rating_key"`
	ParentRatingKey      string `json:"parent_rating_key"`
	GrandparentRatingKey string `json:"grandparent_rating_key"`
	Title                string `json:"title"`
	ParentTitle          string `json:"parent_title"`
	GrandparentTitle     string `json:"grandparent_title"`
	FullTitle            string `json:"full_title"`
	Year                 string `json:"year"`
	Thumb                string `json:"thumb"`
	Art                  string `json:"art"`
	Summary              string `json:"summary"`
	User                 string `json:"user"`
	UserID               int    `json:"user_id"`
	FriendlyName         string `json:"friendly_name"`
	Player               string `json:"player"`
	Product              string `json:"product"`
	Platform             string `json:"platform"`
	IPAddress            string `json:"ip_address"`
	State                string `json:"state"`
	ProgressPercent      string `json:"progress_percent"`
	TranscodeDecision    string `json:"transcode_decision"`
	Duration             string `json:"duration"`
	ViewOffset           string `json:"view_offset"`
	Bandwidth            string `json:"bandwidth"`
	QualityProfile       string `json:"quality_profile"`
	Container            string `json:"container"`
	VideoCodec           string `json:"video_codec"`
	VideoResolution      string `json:"video_resolution"`
	AudioCodec           string `json:"audio_codec"`
	AudioChannels        string `json:"audio_channels"`
}

// HistoryResponse wraps the paginated history response.
type HistoryResponse struct {
	RecordsTotal    int             `json:"recordsTotal"`
	RecordsFiltered int             `json:"recordsFiltered"`
	TotalDuration   string          `json:"total_duration"`
	FilterDuration  string          `json:"filter_duration"`
	Data            []HistoryRecord `json:"data"`
}

// HistoryRecord represents a single watch history entry.
type HistoryRecord struct {
	Date                  int64  `json:"date"`
	RowID                 int    `json:"row_id"`
	ReferenceID           int    `json:"reference_id"`
	UserID                int    `json:"user_id"`
	User                  string `json:"user"`
	FriendlyName          string `json:"friendly_name"`
	Platform              string `json:"platform"`
	Product               string `json:"product"`
	Player                string `json:"player"`
	IPAddress             string `json:"ip_address"`
	Live                  int    `json:"live"`
	MachineID             string `json:"machine_id"`
	Location              string `json:"location"`
	Secure                int    `json:"secure"`
	Relayed               int    `json:"relayed"`
	MediaType             string `json:"media_type"`
	RatingKey             int    `json:"rating_key"`
	ParentRatingKey       int    `json:"parent_rating_key"`
	GrandparentRatingKey  int    `json:"grandparent_rating_key"`
	FullTitle             string `json:"full_title"`
	Title                 string `json:"title"`
	ParentTitle           string `json:"parent_title"`
	GrandparentTitle      string `json:"grandparent_title"`
	OriginalTitle         string `json:"original_title"`
	Year                  int    `json:"year"`
	MediaIndex            int    `json:"media_index"`
	ParentMediaIndex      int    `json:"parent_media_index"`
	Thumb                 string `json:"thumb"`
	OriginallyAvailableAt string `json:"originally_available_at"`
	GUID                  string `json:"guid"`
	TranscodeDecision     string `json:"transcode_decision"`
	PercentComplete       int    `json:"percent_complete"`
	WatchedStatus         int    `json:"watched_status"`
	GroupCount            int    `json:"group_count"`
	GroupIDs              string `json:"group_ids"`
	State                 any    `json:"state"`
	SessionKey            any    `json:"session_key"`
	Started               int64  `json:"started"`
	Stopped               int64  `json:"stopped"`
	PausedCounter         int    `json:"paused_counter"`
	PlayDuration          int    `json:"play_duration"`
}

// Library represents a Plex library section as seen by Tautulli.
type Library struct {
	SectionID   string `json:"section_id"`
	SectionName string `json:"section_name"`
	SectionType string `json:"section_type"`
	Count       string `json:"count"`
	ParentCount string `json:"parent_count"`
	ChildCount  string `json:"child_count"`
	IsActive    int    `json:"is_active"`
	Thumb       string `json:"thumb"`
	Art         string `json:"art"`
}

// LibraryDetail represents detailed library information.
type LibraryDetail struct {
	SectionID       string `json:"section_id"`
	SectionName     string `json:"section_name"`
	SectionType     string `json:"section_type"`
	Count           int    `json:"count"`
	ParentCount     any    `json:"parent_count"`
	ChildCount      any    `json:"child_count"`
	IsActive        int    `json:"is_active"`
	DoNotify        int    `json:"do_notify"`
	DoNotifyCreated int    `json:"do_notify_created"`
	KeepHistory     int    `json:"keep_history"`
	DeletedSection  int    `json:"deleted_section"`
	LastAccessed    int64  `json:"last_accessed"`
	RowID           int    `json:"row_id"`
	ServerID        string `json:"server_id"`
	LibraryThumb    string `json:"library_thumb"`
	LibraryArt      string `json:"library_art"`
}

// User represents a Plex user as seen by Tautulli.
type User struct {
	UserID          string   `json:"user_id"`
	Username        string   `json:"username"`
	FriendlyName    string   `json:"friendly_name"`
	Email           string   `json:"email"`
	Thumb           string   `json:"thumb"`
	IsActive        int      `json:"is_active"`
	IsAdmin         int      `json:"is_admin"`
	IsHomeUser      int      `json:"is_home_user"`
	IsAllowSync     int      `json:"is_allow_sync"`
	IsRestricted    int      `json:"is_restricted"`
	DoNotify        int      `json:"do_notify"`
	KeepHistory     int      `json:"keep_history"`
	AllowGuest      int      `json:"allow_guest"`
	SharedLibraries []string `json:"shared_libraries"`
	RowID           int      `json:"row_id"`
	DeletedUser     int      `json:"deleted_user"`
}

// UserDetail represents detailed user information.
type UserDetail struct {
	UserID          int      `json:"user_id"`
	Username        string   `json:"username"`
	FriendlyName    string   `json:"friendly_name"`
	Email           string   `json:"email"`
	Thumb           string   `json:"user_thumb"`
	IsActive        int      `json:"is_active"`
	IsAdmin         int      `json:"is_admin"`
	IsHomeUser      int      `json:"is_home_user"`
	IsAllowSync     int      `json:"is_allow_sync"`
	IsRestricted    int      `json:"is_restricted"`
	DoNotify        int      `json:"do_notify"`
	KeepHistory     int      `json:"keep_history"`
	AllowGuest      int      `json:"allow_guest"`
	SharedLibraries []string `json:"shared_libraries"`
	RowID           int      `json:"row_id"`
	DeletedUser     int      `json:"deleted_user"`
	LastSeen        int64    `json:"last_seen"`
}

// HomeStats represents a single homepage statistics category.
type HomeStats struct {
	StatID   string `json:"stat_id"`
	StatType string `json:"stat_type,omitempty"`
	Rows     []any  `json:"rows"`
}

// ServerInfo represents Plex server information from Tautulli.
type ServerInfo struct {
	PMSIdentifier string `json:"pms_identifier"`
	PMSIP         string `json:"pms_ip"`
	PMSIsRemote   int    `json:"pms_is_remote"`
	PMSName       string `json:"pms_name"`
	PMSPlatform   string `json:"pms_platform"`
	PMSPlexPass   int    `json:"pms_plexpass"`
	PMSPort       int    `json:"pms_port"`
	PMSSSL        int    `json:"pms_ssl"`
	PMSURL        string `json:"pms_url"`
	PMSURLManual  int    `json:"pms_url_manual"`
	PMSVersion    string `json:"pms_version"`
}

// Info represents information about the Tautulli server itself.
type Info struct {
	InstallType         string `json:"tautulli_install_type"`
	Version             string `json:"tautulli_version"`
	Branch              string `json:"tautulli_branch"`
	Commit              string `json:"tautulli_commit"`
	Platform            string `json:"tautulli_platform"`
	PlatformRelease     string `json:"tautulli_platform_release"`
	PlatformVersion     string `json:"tautulli_platform_version"`
	PlatformLinuxDistro string `json:"tautulli_platform_linux_distro"`
	PlatformDeviceName  string `json:"tautulli_platform_device_name"`
	PythonVersion       string `json:"tautulli_python_version"`
}

// ItemMetadata represents a media item's metadata from Tautulli.
type ItemMetadata struct {
	MediaType             string   `json:"media_type"`
	RatingKey             string   `json:"rating_key"`
	ParentRatingKey       string   `json:"parent_rating_key"`
	GrandparentRatingKey  string   `json:"grandparent_rating_key"`
	Title                 string   `json:"title"`
	ParentTitle           string   `json:"parent_title"`
	GrandparentTitle      string   `json:"grandparent_title"`
	OriginalTitle         string   `json:"original_title"`
	SortTitle             string   `json:"sort_title"`
	EditionTitle          string   `json:"edition_title"`
	MediaIndex            string   `json:"media_index"`
	ParentMediaIndex      string   `json:"parent_media_index"`
	Studio                string   `json:"studio"`
	ContentRating         string   `json:"content_rating"`
	Summary               string   `json:"summary"`
	Tagline               string   `json:"tagline"`
	Rating                string   `json:"rating"`
	AudienceRating        string   `json:"audience_rating"`
	UserRating            string   `json:"user_rating"`
	Duration              string   `json:"duration"`
	Year                  string   `json:"year"`
	Thumb                 string   `json:"thumb"`
	Art                   string   `json:"art"`
	Banner                string   `json:"banner"`
	GUID                  string   `json:"guid"`
	SectionID             string   `json:"section_id"`
	LibraryName           string   `json:"library_name"`
	OriginallyAvailableAt string   `json:"originally_available_at"`
	AddedAt               string   `json:"added_at"`
	UpdatedAt             string   `json:"updated_at"`
	LastViewedAt          string   `json:"last_viewed_at"`
	Actors                []string `json:"actors"`
	Directors             []string `json:"directors"`
	Writers               []string `json:"writers"`
	Genres                []string `json:"genres"`
	Labels                []string `json:"labels"`
	Collections           []string `json:"collections"`
}

// RecentlyAddedResponse wraps the recently added items.
type RecentlyAddedResponse struct {
	RecentlyAdded []ItemMetadata `json:"recently_added"`
}

// SearchResults represents search results from the PMS.
type SearchResults struct {
	ResultsCount int `json:"results_count"`
	ResultsList  any `json:"results_list"`
}

// WatchTimeStats represents watch time statistics for a user or library.
type WatchTimeStats struct {
	QueryDays  int `json:"query_days"`
	TotalPlays int `json:"total_plays"`
	TotalTime  int `json:"total_time"`
}

// GeoIPInfo represents geolocation information for an IP address.
type GeoIPInfo struct {
	City       string  `json:"city"`
	Code       string  `json:"code"`
	Continent  string  `json:"continent"`
	Country    string  `json:"country"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	PostalCode string  `json:"postal_code"`
	Region     string  `json:"region"`
	Timezone   string  `json:"timezone"`
}

// LibraryName represents a minimal library entry.
type LibraryName struct {
	SectionID   int    `json:"section_id"`
	SectionName string `json:"section_name"`
	SectionType string `json:"section_type"`
}

// UserName represents a minimal user entry.
type UserName struct {
	UserID       int    `json:"user_id"`
	FriendlyName string `json:"friendly_name"`
}

// ServerStatus represents the connection status.
type ServerStatus struct {
	Connected bool `json:"connected"`
}
