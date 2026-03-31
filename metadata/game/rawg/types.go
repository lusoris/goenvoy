package rawg

// PagedResult is a generic paginated API response.
type PagedResult[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

// GameListItem represents a game in list/search results.
type GameListItem struct {
	ID               int              `json:"id"`
	Slug             string           `json:"slug"`
	Name             string           `json:"name"`
	Released         string           `json:"released"`
	BackgroundImage  string           `json:"background_image"`
	Rating           float64          `json:"rating"`
	RatingTop        int              `json:"rating_top"`
	Metacritic       int              `json:"metacritic"`
	Playtime         int              `json:"playtime"`
	Added            int              `json:"added"`
	AddedByStatus    AddedByStatus    `json:"added_by_status"`
	Platforms        []PlatformInfo   `json:"platforms"`
	Genres           []GenreInfo      `json:"genres"`
	Tags             []TagInfo        `json:"tags"`
	Stores           []StoreInfo      `json:"stores"`
	ShortScreenshots []ScreenshotInfo `json:"short_screenshots"`
}

// AddedByStatus tracks how many users have a game in each status.
type AddedByStatus struct {
	Yet     int `json:"yet"`
	Owned   int `json:"owned"`
	Beaten  int `json:"beaten"`
	Toplay  int `json:"toplay"`
	Dropped int `json:"dropped"`
	Playing int `json:"playing"`
}

// Game represents a full game detail response.
type Game struct {
	ID                  int                  `json:"id"`
	Slug                string               `json:"slug"`
	Name                string               `json:"name"`
	Released            string               `json:"released"`
	BackgroundImage     string               `json:"background_image"`
	Rating              float64              `json:"rating"`
	RatingTop           int                  `json:"rating_top"`
	Metacritic          int                  `json:"metacritic"`
	Playtime            int                  `json:"playtime"`
	Added               int                  `json:"added"`
	AddedByStatus       AddedByStatus        `json:"added_by_status"`
	Platforms           []PlatformInfo       `json:"platforms"`
	Genres              []GenreInfo          `json:"genres"`
	Tags                []TagInfo            `json:"tags"`
	Stores              []StoreInfo          `json:"stores"`
	Description         string               `json:"description"`
	DescriptionRaw      string               `json:"description_raw"`
	MetacriticURL       string               `json:"metacritic_url"`
	Website             string               `json:"website"`
	RedditURL           string               `json:"reddit_url"`
	RedditDescription   string               `json:"reddit_description"`
	RedditName          string               `json:"reddit_name"`
	MetacriticPlatforms []MetacriticPlatform `json:"metacritic_platforms"`
	Developers          []Developer          `json:"developers"`
	Publishers          []Publisher          `json:"publishers"`
}

// PlatformInfo wraps platform detail with release info.
type PlatformInfo struct {
	Platform       PlatformDetail `json:"platform"`
	ReleasedAt     string         `json:"released_at"`
	RequirementsEn any            `json:"requirements_en"`
	RequirementsRu any            `json:"requirements_ru"`
}

// PlatformDetail is the core platform information.
type PlatformDetail struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// GenreInfo represents a genre reference.
type GenreInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// TagInfo represents a tag reference.
type TagInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Language   string `json:"language"`
	GamesCount int    `json:"games_count"`
}

// StoreInfo wraps store detail.
type StoreInfo struct {
	ID    int         `json:"id"`
	Store StoreDetail `json:"store"`
}

// StoreDetail is the core store information.
type StoreDetail struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Domain string `json:"domain"`
}

// ScreenshotInfo represents a short screenshot in list results.
type ScreenshotInfo struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
}

// Screenshot represents a full screenshot.
type Screenshot struct {
	ID     int    `json:"id"`
	Image  string `json:"image"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Trailer represents a game trailer.
type Trailer struct {
	ID      int          `json:"id"`
	Name    string       `json:"name"`
	Preview TrailerMedia `json:"preview"`
	Data    TrailerMedia `json:"data"`
}

// TrailerMedia contains media URLs for a trailer.
type TrailerMedia struct {
	Max string `json:"max"`
}

// Platform represents a full platform resource.
type Platform struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	GamesCount      int    `json:"games_count"`
	Image           string `json:"image"`
	ImageBackground string `json:"image_background"`
	Description     string `json:"description"`
}

// Genre represents a full genre resource.
type Genre struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	GamesCount      int    `json:"games_count"`
	ImageBackground string `json:"image_background"`
}

// Publisher represents a game publisher.
type Publisher struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	GamesCount      int    `json:"games_count"`
	ImageBackground string `json:"image_background"`
}

// Developer represents a game developer.
type Developer struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	GamesCount      int    `json:"games_count"`
	ImageBackground string `json:"image_background"`
}

// Tag represents a game tag.
type Tag struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Language        string `json:"language"`
	GamesCount      int    `json:"games_count"`
	ImageBackground string `json:"image_background"`
}

// Store represents a full store resource.
type Store struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Domain          string `json:"domain"`
	GamesCount      int    `json:"games_count"`
	ImageBackground string `json:"image_background"`
}

// MetacriticPlatform holds a metacritic score for a specific platform.
type MetacriticPlatform struct {
	MetaScore int            `json:"metascore"`
	URL       string         `json:"url"`
	Platform  PlatformDetail `json:"platform"`
}
