package listenbrainz

// TrackMetadata holds metadata about a listened track.
type TrackMetadata struct {
	ArtistName     string         `json:"artist_name"`
	TrackName      string         `json:"track_name"`
	ReleaseName    string         `json:"release_name"`
	AdditionalInfo map[string]any `json:"additional_info,omitempty"`
}

// ListenInfo represents a single listen event returned by the API.
type ListenInfo struct {
	ListenedAt    int64         `json:"listened_at"`
	TrackMetadata TrackMetadata `json:"track_metadata"`
	InsertedAt    int64         `json:"inserted_at"`
}

// ListensResponse is the response from the listens endpoint.
type ListensResponse struct {
	Payload ListensPayload `json:"payload"`
}

// ListensPayload contains the listen data within a [ListensResponse].
type ListensPayload struct {
	Count          int          `json:"count"`
	Listens        []ListenInfo `json:"listens"`
	LatestListenTS int64        `json:"latest_listen_ts"`
	OldestListenTS int64        `json:"oldest_listen_ts"`
	UserName       string       `json:"user_name"`
}

// PlayingNow is the response from the playing-now endpoint.
type PlayingNow struct {
	Payload PlayingNowPayload `json:"payload"`
}

// PlayingNowPayload contains the currently playing track data.
type PlayingNowPayload struct {
	Listens []ListenInfo `json:"listens"`
	Count   int          `json:"count"`
}

// ArtistStatsResponse is the response from the user top artists endpoint.
type ArtistStatsResponse struct {
	Payload ArtistStatsPayload `json:"payload"`
}

// ArtistStatsPayload contains top artist statistics.
type ArtistStatsPayload struct {
	Artists          []ArtistStat `json:"artists"`
	Count            int          `json:"count"`
	TotalArtistCount int          `json:"total_artist_count"`
	Offset           int          `json:"offset"`
	Range            string       `json:"range"`
	UserName         string       `json:"user_name"`
	LastUpdated      int64        `json:"last_updated"`
}

// ArtistStat represents a single artist in a top-artists response.
type ArtistStat struct {
	ArtistName  string `json:"artist_name"`
	ArtistMBID  string `json:"artist_mbid,omitempty"`
	ListenCount int    `json:"listen_count"`
}

// ReleaseStatsResponse is the response from the user top releases endpoint.
type ReleaseStatsResponse struct {
	Payload ReleaseStatsPayload `json:"payload"`
}

// ReleaseStatsPayload contains top release statistics.
type ReleaseStatsPayload struct {
	Releases          []ReleaseStat `json:"releases"`
	Count             int           `json:"count"`
	TotalReleaseCount int           `json:"total_release_count"`
	Offset            int           `json:"offset"`
	Range             string        `json:"range"`
	UserName          string        `json:"user_name"`
	LastUpdated       int64         `json:"last_updated"`
}

// ReleaseStat represents a single release in a top-releases response.
type ReleaseStat struct {
	ArtistName  string `json:"artist_name"`
	ReleaseName string `json:"release_name"`
	ReleaseMBID string `json:"release_mbid,omitempty"`
	ListenCount int    `json:"listen_count"`
}

// RecordingStatsResponse is the response from the user top recordings endpoint.
type RecordingStatsResponse struct {
	Payload RecordingStatsPayload `json:"payload"`
}

// RecordingStatsPayload contains top recording statistics.
type RecordingStatsPayload struct {
	Recordings          []RecordingStat `json:"recordings"`
	Count               int             `json:"count"`
	TotalRecordingCount int             `json:"total_recording_count"`
	Offset              int             `json:"offset"`
	Range               string          `json:"range"`
	UserName            string          `json:"user_name"`
	LastUpdated         int64           `json:"last_updated"`
}

// RecordingStat represents a single recording in a top-recordings response.
type RecordingStat struct {
	ArtistName    string `json:"artist_name"`
	TrackName     string `json:"track_name"`
	RecordingMBID string `json:"recording_mbid,omitempty"`
	ReleaseName   string `json:"release_name"`
	ListenCount   int    `json:"listen_count"`
}

// ActivityResponse is the response from the listening-activity endpoint.
type ActivityResponse struct {
	Payload ActivityPayload `json:"payload"`
}

// ActivityPayload contains listening activity data.
type ActivityPayload struct {
	ListeningActivity []ActivityEntry `json:"listening_activity"`
	UserName          string          `json:"user_name"`
	Range             string          `json:"range"`
	LastUpdated       int64           `json:"last_updated"`
}

// ActivityEntry represents listening activity for a time period.
type ActivityEntry struct {
	ListenCount int    `json:"listen_count"`
	FromTS      int64  `json:"from_ts"`
	ToTS        int64  `json:"to_ts"`
	TimeRange   string `json:"time_range"`
}

// DailyActivityResponse is the response from the daily-activity endpoint.
type DailyActivityResponse struct {
	Payload DailyActivityPayload `json:"payload"`
}

// DailyActivityPayload contains daily listening activity data.
type DailyActivityPayload struct {
	DailyActivity map[string][]HourEntry `json:"daily_activity"`
	UserName      string                 `json:"user_name"`
	Range         string                 `json:"range"`
	LastUpdated   int64                  `json:"last_updated"`
}

// HourEntry represents listening activity for a single hour of the day.
type HourEntry struct {
	Hour        int `json:"hour"`
	ListenCount int `json:"listen_count"`
}

// SimilarUser represents a user similar to the queried user.
type SimilarUser struct {
	UserName   string  `json:"user_name"`
	Similarity float64 `json:"similarity"`
}

type similarUsersResponse struct {
	Payload []SimilarUser `json:"payload"`
}

// submitListensPayload is the request body for the submit-listens endpoint.
type submitListensPayload struct {
	ListenType string         `json:"listen_type"`
	Payload    []listenSubmit `json:"payload"`
}

// listenSubmit represents a single listen in a submit request.
type listenSubmit struct {
	ListenedAt    int64         `json:"listened_at,omitempty"`
	TrackMetadata TrackMetadata `json:"track_metadata"`
}

// Listen is a convenience type for submitting listens.
type Listen struct {
	ListenedAt    int64
	TrackMetadata TrackMetadata
}

type listenCountResponse struct {
	Payload struct {
		Count int64 `json:"count"`
	} `json:"payload"`
}

type latestImportResponse struct {
	LatestImport int64 `json:"latest_import"`
}
