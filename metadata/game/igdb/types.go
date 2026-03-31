package igdb

// Game represents an IGDB game entry.
type Game struct {
	ID                    int     `json:"id"`
	Name                  string  `json:"name"`
	Slug                  string  `json:"slug"`
	Summary               string  `json:"summary"`
	Storyline             string  `json:"storyline"`
	URL                   string  `json:"url"`
	Rating                float64 `json:"rating"`
	AggregatedRating      float64 `json:"aggregated_rating"`
	TotalRating           float64 `json:"total_rating"`
	RatingCount           int     `json:"rating_count"`
	AggregatedRatingCount int     `json:"aggregated_rating_count"`
	TotalRatingCount      int     `json:"total_rating_count"`
	FirstReleaseDate      int64   `json:"first_release_date"`
	Category              int     `json:"category"`
	Genres                []int   `json:"genres"`
	Platforms             []int   `json:"platforms"`
	Themes                []int   `json:"themes"`
	Cover                 int     `json:"cover"`
	Screenshots           []int   `json:"screenshots"`
	Videos                []int   `json:"videos"`
	SimilarGames          []int   `json:"similar_games"`
	Franchises            []int   `json:"franchises"`
	GameEngines           []int   `json:"game_engines"`
}

// Platform represents an IGDB gaming platform.
type Platform struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Abbreviation    string `json:"abbreviation"`
	AlternativeName string `json:"alternative_name"`
	Summary         string `json:"summary"`
	URL             string `json:"url"`
	Generation      int    `json:"generation"`
	Category        int    `json:"category"`
}

// Genre represents an IGDB game genre.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	URL  string `json:"url"`
}

// Company represents an IGDB game company.
type Company struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Country     int    `json:"country"`
	Developed   []int  `json:"developed"`
	Published   []int  `json:"published"`
}

// Cover represents an IGDB game cover image.
type Cover struct {
	ID      int    `json:"id"`
	Game    int    `json:"game"`
	ImageID string `json:"image_id"`
	URL     string `json:"url"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

// Screenshot represents an IGDB game screenshot.
type Screenshot struct {
	ID      int    `json:"id"`
	Game    int    `json:"game"`
	ImageID string `json:"image_id"`
	URL     string `json:"url"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

// GameVideo represents an IGDB game video.
type GameVideo struct {
	ID      int    `json:"id"`
	Game    int    `json:"game"`
	Name    string `json:"name"`
	VideoID string `json:"video_id"`
}
