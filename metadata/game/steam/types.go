package steam

import "encoding/json"

// AppDetails represents detailed information about a Steam application.
type AppDetails struct {
	Type                string           `json:"type"`
	Name                string           `json:"name"`
	SteamAppID          int              `json:"steam_appid"`
	RequiredAge         int              `json:"required_age"`
	IsFree              bool             `json:"is_free"`
	DetailedDescription string           `json:"detailed_description"`
	AboutTheGame        string           `json:"about_the_game"`
	ShortDescription    string           `json:"short_description"`
	SupportedLanguages  string           `json:"supported_languages"`
	HeaderImage         string           `json:"header_image"`
	Website             string           `json:"website"`
	Developers          []string         `json:"developers"`
	Publishers          []string         `json:"publishers"`
	PriceOverview       *PriceInfo       `json:"price_overview,omitempty"`
	Platforms           PlatformSupport  `json:"platforms"`
	Metacritic          *MetacriticInfo  `json:"metacritic,omitempty"`
	Categories          []CategoryInfo   `json:"categories"`
	Genres              []GenreInfo      `json:"genres"`
	Screenshots         []ScreenshotInfo `json:"screenshots"`
	Movies              []MovieInfo      `json:"movies"`
	ReleaseDate         ReleaseDateInfo  `json:"release_date"`
	Background          string           `json:"background"`
}

// PriceInfo represents pricing information for an app.
type PriceInfo struct {
	Currency         string `json:"currency"`
	Initial          int    `json:"initial"`
	Final            int    `json:"final"`
	DiscountPercent  int    `json:"discount_percent"`
	InitialFormatted string `json:"initial_formatted"`
	FinalFormatted   string `json:"final_formatted"`
}

// PlatformSupport indicates which platforms an app supports.
type PlatformSupport struct {
	Windows bool `json:"windows"`
	Mac     bool `json:"mac"`
	Linux   bool `json:"linux"`
}

// MetacriticInfo holds Metacritic score data.
type MetacriticInfo struct {
	Score int    `json:"score"`
	URL   string `json:"url"`
}

// CategoryInfo represents a Steam store category.
type CategoryInfo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

// GenreInfo represents a Steam store genre.
type GenreInfo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// ScreenshotInfo holds screenshot image URLs.
type ScreenshotInfo struct {
	ID            int    `json:"id"`
	PathThumbnail string `json:"path_thumbnail"`
	PathFull      string `json:"path_full"`
}

// MovieInfo holds video/trailer information.
type MovieInfo struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Thumbnail string            `json:"thumbnail"`
	Webm      map[string]string `json:"webm"`
	MP4       map[string]string `json:"mp4"`
}

// ReleaseDateInfo holds release date information.
type ReleaseDateInfo struct {
	ComingSoon bool   `json:"coming_soon"`
	Date       string `json:"date"`
}

// FeaturedResponse is the response from the featured games endpoint.
type FeaturedResponse struct {
	LargeCapsules []FeaturedItem `json:"large_capsules"`
	FeaturedWin   []FeaturedItem `json:"featured_win"`
	FeaturedMac   []FeaturedItem `json:"featured_mac"`
	FeaturedLinux []FeaturedItem `json:"featured_linux"`
}

// FeaturedItem represents a featured game.
type FeaturedItem struct {
	ID                int    `json:"id"`
	Type              int    `json:"type"`
	Name              string `json:"name"`
	Discounted        bool   `json:"discounted"`
	DiscountPercent   int    `json:"discount_percent"`
	OriginalPrice     int    `json:"original_price"`
	FinalPrice        int    `json:"final_price"`
	Currency          string `json:"currency"`
	LargeCapsuleImage string `json:"large_capsule_image"`
	SmallCapsuleImage string `json:"small_capsule_image"`
	WindowsAvailable  bool   `json:"windows_available"`
	MacAvailable      bool   `json:"mac_available"`
	LinuxAvailable    bool   `json:"linux_available"`
	HeaderImage       string `json:"header_image"`
}

// FeaturedCategories is the response from the featured categories endpoint.
type FeaturedCategories map[string]json.RawMessage

// AppListEntry represents a single app in the full Steam app list.
type AppListEntry struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}

// NewsItem represents a news article for a Steam app.
type NewsItem struct {
	GID       string   `json:"gid"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	Author    string   `json:"author"`
	Contents  string   `json:"contents"`
	FeedLabel string   `json:"feedlabel"`
	Date      int64    `json:"date"`
	FeedName  string   `json:"feedname"`
	Tags      []string `json:"tags"`
}

// Achievement represents a global achievement with its unlock percentage.
type Achievement struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
}

// appDetailsWrapper is used to parse the map[appid]{success, data} response.
type appDetailsWrapper struct {
	Success bool       `json:"success"`
	Data    AppDetails `json:"data"`
}

type appListResponse struct {
	AppList struct {
		Apps []AppListEntry `json:"apps"`
	} `json:"applist"`
}

type currentPlayersResponse struct {
	Response struct {
		PlayerCount int `json:"player_count"`
		Result      int `json:"result"`
	} `json:"response"`
}

type appNewsResponse struct {
	AppNews struct {
		AppID     int        `json:"appid"`
		NewsItems []NewsItem `json:"newsitems"`
	} `json:"appnews"`
}

type achievementsResponse struct {
	AchievementPercentages struct {
		Achievements []Achievement `json:"achievements"`
	} `json:"achievementpercentages"`
}
