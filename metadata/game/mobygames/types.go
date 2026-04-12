package mobygames

// Game represents a game from the MobyGames API (normal format).
type Game struct {
	GameID            int              `json:"game_id"`
	Title             string           `json:"title"`
	MobyURL           string           `json:"moby_url"`
	AlternateTitles   []AlternateTitle `json:"alternate_titles"`
	Description       string           `json:"description"`
	Genres            []Genre          `json:"genres"`
	MobyScore         float64          `json:"moby_score"`
	NumVotes          int              `json:"num_votes"`
	OfficialURL       string           `json:"official_url"`
	Platforms         []PlatformRef    `json:"platforms"`
	SampleCover       *SampleCover     `json:"sample_cover"`
	SampleScreenshots []Screenshot     `json:"sample_screenshots"`
}

// GameBrief represents a game in brief format.
type GameBrief struct {
	GameID int    `json:"game_id"`
	Title  string `json:"title"`
}

// AlternateTitle is an alternate game title.
type AlternateTitle struct {
	Description string `json:"description"`
	Title       string `json:"title"`
}

// Genre represents a game genre.
type Genre struct {
	GenreID       int    `json:"genre_id"`
	GenreName     string `json:"genre_name"`
	GenreCategory string `json:"genre_category"`
}

// PlatformRef is a platform reference in a game response.
type PlatformRef struct {
	PlatformID   int    `json:"platform_id"`
	PlatformName string `json:"platform_name"`
	FirstRelease string `json:"first_release_date"`
}

// Platform represents a platform from the platforms endpoint.
type Platform struct {
	PlatformID   int    `json:"platform_id"`
	PlatformName string `json:"platform_name"`
}

// SampleCover is a cover image included in game responses.
type SampleCover struct {
	Height       int      `json:"height"`
	Width        int      `json:"width"`
	Image        string   `json:"image"`
	ThumbnailURL string   `json:"thumbnail_image"`
	Platforms    []string `json:"platforms"`
}

// Screenshot represents a game screenshot.
type Screenshot struct {
	Height  int    `json:"height"`
	Width   int    `json:"width"`
	Image   string `json:"image"`
	Caption string `json:"caption"`
}

// Cover represents a cover image from the covers endpoint.
type Cover struct {
	Height       int      `json:"height"`
	Width        int      `json:"width"`
	Image        string   `json:"image"`
	ThumbnailURL string   `json:"thumbnail_image"`
	ScanOf       string   `json:"scan_of"`
	Platforms    []string `json:"platforms"`
	Countries    []string `json:"countries"`
}

// CoverGroup wraps covers for a specific platform.
type CoverGroup struct {
	Covers []Cover `json:"covers"`
}

// PlatformDetail represents detailed platform info for a game.
type PlatformDetail struct {
	PlatformID   int    `json:"platform_id"`
	PlatformName string `json:"platform_name"`
	FirstRelease string `json:"first_release_date"`
}

// Group represents a game group from MobyGames.
type Group struct {
	GroupID          int    `json:"group_id"`
	GroupName        string `json:"group_name"`
	GroupDescription string `json:"group_description"`
}

// GamesResult is the response from the games list endpoint.
type GamesResult struct {
	Games []Game `json:"games"`
}
