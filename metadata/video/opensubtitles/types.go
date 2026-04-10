package opensubtitles

// SearchParams contains parameters for subtitle search.
type SearchParams struct {
	Query           string `json:"query,omitempty"`
	IMDBID          int    `json:"imdb_id,omitempty"`
	TMDBID          int    `json:"tmdb_id,omitempty"`
	Languages       string `json:"languages,omitempty"`
	MovieHash       string `json:"moviehash,omitempty"`
	Type            string `json:"type,omitempty"` // movie, episode, all
	SeasonNumber    int    `json:"season_number,omitempty"`
	EpisodeNumber   int    `json:"episode_number,omitempty"`
	ParentFeatureID int    `json:"parent_feature_id,omitempty"`
	ParentIMDBID    int    `json:"parent_imdb_id,omitempty"`
	ParentTMDBID    int    `json:"parent_tmdb_id,omitempty"`
	Year            int    `json:"year,omitempty"`
	Page            int    `json:"page,omitempty"`
	OrderBy         string `json:"order_by,omitempty"`
	OrderDirection  string `json:"order_direction,omitempty"`
}

// SearchResponse is the response from a subtitle search.
type SearchResponse struct {
	TotalPages int        `json:"total_pages"`
	TotalCount int        `json:"total_count"`
	PerPage    int        `json:"per_page"`
	Page       int        `json:"page"`
	Data       []Subtitle `json:"data"`
}

// Subtitle represents a single subtitle result.
type Subtitle struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes SubtitleAttributes `json:"attributes"`
}

// SubtitleAttributes contains the details of a subtitle.
type SubtitleAttributes struct {
	SubtitleID        string          `json:"subtitle_id"`
	Language          string          `json:"language"`
	DownloadCount     int             `json:"download_count"`
	NewDownloadCount  int             `json:"new_download_count"`
	HearingImpaired   bool            `json:"hearing_impaired"`
	HD                bool            `json:"hd"`
	FPS               float64         `json:"fps"`
	Votes             int             `json:"votes"`
	Points            int             `json:"points,omitempty"`
	Ratings           float64         `json:"ratings"`
	FromTrusted       bool            `json:"from_trusted"`
	ForeignPartsOnly  bool            `json:"foreign_parts_only"`
	AITranslated      bool            `json:"ai_translated"`
	MachineTranslated bool            `json:"machine_translated"`
	UploadDate        string          `json:"upload_date"`
	Release           string          `json:"release"`
	Comments          string          `json:"comments"`
	LegacySubtitleID  int             `json:"legacy_subtitle_id,omitempty"`
	Uploader          *Uploader       `json:"uploader,omitempty"`
	FeatureDetails    *FeatureDetails `json:"feature_details,omitempty"`
	URL               string          `json:"url,omitempty"`
	Files             []File          `json:"files,omitempty"`
	MovieHashMatch    bool            `json:"moviehash_match,omitempty"`
}

// Uploader represents the subtitle uploader.
type Uploader struct {
	UploaderID int    `json:"uploader_id"`
	Name       string `json:"name"`
	Rank       string `json:"rank"`
}

// FeatureDetails contains movie/episode details associated with a subtitle.
type FeatureDetails struct {
	FeatureID     int    `json:"feature_id"`
	FeatureType   string `json:"feature_type"`
	Year          int    `json:"year"`
	Title         string `json:"title"`
	MovieName     string `json:"movie_name"`
	IMDBID        int    `json:"imdb_id"`
	TMDBID        int    `json:"tmdb_id"`
	SeasonNumber  int    `json:"season_number,omitempty"`
	EpisodeNumber int    `json:"episode_number,omitempty"`
	ParentIMDBID  int    `json:"parent_imdb_id,omitempty"`
	ParentTitle   string `json:"parent_title,omitempty"`
}

// File represents a subtitle file.
type File struct {
	FileID   int    `json:"file_id"`
	CDNumber int    `json:"cd_number"`
	FileName string `json:"file_name"`
}

// DownloadRequest is the request body for downloading a subtitle.
type DownloadRequest struct {
	FileID        int    `json:"file_id"`
	SubFormat     string `json:"sub_format,omitempty"`
	FileName      string `json:"file_name,omitempty"`
	InFPS         int    `json:"in_fps,omitempty"`
	OutFPS        int    `json:"out_fps,omitempty"`
	Timeshift     int    `json:"timeshift,omitempty"`
	ForceDownload bool   `json:"force_download,omitempty"`
}

// DownloadResponse is returned after requesting a download link.
type DownloadResponse struct {
	Link         string `json:"link"`
	FileName     string `json:"file_name"`
	Requests     int    `json:"requests"`
	Remaining    int    `json:"remaining"`
	Message      string `json:"message"`
	ResetTime    string `json:"reset_time"`
	ResetTimeUTC string `json:"reset_time_utc"`
}

// LoginRequest is the request body for user login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse is returned after a successful login.
type LoginResponse struct {
	User                *UserInfo `json:"user,omitempty"`
	BaseURL             string    `json:"base_url,omitempty"`
	Token               string    `json:"token"`
	Status              int       `json:"status"`
	AllowedTranslations int       `json:"allowed_translations,omitempty"`
	AllowedDownloads    int       `json:"allowed_downloads,omitempty"`
}

// UserInfo represents authenticated user information.
type UserInfo struct {
	AllowedDownloads    int    `json:"allowed_downloads"`
	AllowedTranslations int    `json:"allowed_translations"`
	Level               string `json:"level"`
	UserID              int    `json:"user_id"`
	ExtInstalled        bool   `json:"ext_installed"`
	VIP                 bool   `json:"vip"`
}

// FeaturesResponse is returned from a features search.
type FeaturesResponse struct {
	TotalPages int       `json:"total_pages"`
	TotalCount int       `json:"total_count"`
	PerPage    int       `json:"per_page"`
	Page       int       `json:"page"`
	Data       []Feature `json:"data"`
}

// Feature represents a movie or TV show.
type Feature struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes FeatureAttributes `json:"attributes"`
}

// FeatureAttributes contains the details of a feature.
type FeatureAttributes struct {
	Title           string         `json:"title"`
	OriginalTitle   string         `json:"original_title,omitempty"`
	Year            string         `json:"year,omitempty"`
	SubtitlesCount  int            `json:"subtitles_count,omitempty"`
	SubtitlesCounts map[string]int `json:"subtitles_counts,omitempty"`
	FeatureID       string         `json:"feature_id"`
	IMDBID          int            `json:"imdb_id,omitempty"`
	TMDBID          int            `json:"tmdb_id,omitempty"`
	FeatureType     string         `json:"feature_type,omitempty"`
	URL             string         `json:"url,omitempty"`
	ImgURL          string         `json:"img_url,omitempty"`
}

// Language represents an available subtitle language.
type Language struct {
	LanguageCode string `json:"language_code"`
	LanguageName string `json:"language_name"`
}

// SubtitleFormat represents an available subtitle format.
type SubtitleFormat struct {
	FormatName string `json:"format_name"`
	Extension  string `json:"extension"`
}
