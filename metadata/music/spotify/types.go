package spotify

// Image represents a Spotify image resource.
type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// FollowerInfo contains follower count information.
type FollowerInfo struct {
	Total int `json:"total"`
}

// Copyright represents a copyright statement.
type Copyright struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

// ArtistSimple is a simplified artist object.
type ArtistSimple struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	ExternalURLs map[string]string `json:"external_urls"`
	URI          string            `json:"uri"`
	Href         string            `json:"href"`
	Type         string            `json:"type"`
}

// Artist is a full Spotify artist object.
type Artist struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Genres       []string          `json:"genres"`
	Popularity   int               `json:"popularity"`
	Followers    FollowerInfo      `json:"followers"`
	Images       []Image           `json:"images"`
	ExternalURLs map[string]string `json:"external_urls"`
	URI          string            `json:"uri"`
	Href         string            `json:"href"`
	Type         string            `json:"type"`
}

// AlbumSimple is a simplified album object.
type AlbumSimple struct {
	ID                   string            `json:"id"`
	Name                 string            `json:"name"`
	Artists              []ArtistSimple    `json:"artists"`
	Images               []Image           `json:"images"`
	ReleaseDate          string            `json:"release_date"`
	ReleaseDatePrecision string            `json:"release_date_precision"`
	AlbumType            string            `json:"album_type"`
	TotalTracks          int               `json:"total_tracks"`
	ExternalURLs         map[string]string `json:"external_urls"`
	URI                  string            `json:"uri"`
	Href                 string            `json:"href"`
	Type                 string            `json:"type"`
}

// Album is a full Spotify album object.
type Album struct {
	AlbumSimple
	Tracks     Paged[TrackSimple] `json:"tracks"`
	Genres     []string           `json:"genres"`
	Label      string             `json:"label"`
	Popularity int                `json:"popularity"`
	Copyrights []Copyright        `json:"copyrights"`
}

// TrackSimple is a simplified track object.
type TrackSimple struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Artists      []ArtistSimple    `json:"artists"`
	DiscNumber   int               `json:"disc_number"`
	TrackNumber  int               `json:"track_number"`
	DurationMS   int               `json:"duration_ms"`
	Explicit     bool              `json:"explicit"`
	ExternalURLs map[string]string `json:"external_urls"`
	URI          string            `json:"uri"`
	Href         string            `json:"href"`
	Type         string            `json:"type"`
	PreviewURL   string            `json:"preview_url"`
}

// Track is a full Spotify track object.
type Track struct {
	TrackSimple
	Album       AlbumSimple       `json:"album"`
	Popularity  int               `json:"popularity"`
	ExternalIDs map[string]string `json:"external_ids"`
}

// AudioFeatures contains audio analysis features for a track.
type AudioFeatures struct {
	ID               string  `json:"id"`
	Danceability     float64 `json:"danceability"`
	Energy           float64 `json:"energy"`
	Key              float64 `json:"key"`
	Loudness         float64 `json:"loudness"`
	Mode             float64 `json:"mode"`
	Speechiness      float64 `json:"speechiness"`
	Acousticness     float64 `json:"acousticness"`
	Instrumentalness float64 `json:"instrumentalness"`
	Liveness         float64 `json:"liveness"`
	Valence          float64 `json:"valence"`
	Tempo            float64 `json:"tempo"`
	DurationMS       int     `json:"duration_ms"`
	TimeSignature    int     `json:"time_signature"`
	URI              string  `json:"uri"`
	TrackHref        string  `json:"track_href"`
	AnalysisURL      string  `json:"analysis_url"`
	Type             string  `json:"type"`
}

// Paged is a generic paginated response.
type Paged[T any] struct {
	Items    []T    `json:"items"`
	Total    int    `json:"total"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Href     string `json:"href"`
}

// SearchResult contains results from a Spotify search.
type SearchResult struct {
	Artists *Paged[Artist]      `json:"artists,omitempty"`
	Albums  *Paged[AlbumSimple] `json:"albums,omitempty"`
	Tracks  *Paged[Track]       `json:"tracks,omitempty"`
}

// Category represents a Spotify browse category.
type Category struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Icons []Image `json:"icons"`
	Href  string  `json:"href"`
}

// RecommendationSeeds contains seed parameters for recommendations.
type RecommendationSeeds struct {
	SeedArtists []string
	SeedGenres  []string
	SeedTracks  []string
}

// Seed represents a recommendation seed object.
type Seed struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	InitialPoolSize    int    `json:"initialPoolSize"`
	AfterFilteringSize int    `json:"afterFilteringSize"`
	AfterRelinkingSize int    `json:"afterRelinkingSize"`
}

// Recommendations contains recommendation results.
type Recommendations struct {
	Seeds  []Seed  `json:"seeds"`
	Tracks []Track `json:"tracks"`
}

// topTracksResp wraps the top tracks response.
type topTracksResp struct {
	Tracks []Track `json:"tracks"`
}

// relatedArtistsResp wraps the related artists response.
type relatedArtistsResp struct {
	Artists []Artist `json:"artists"`
}

// categoriesResp wraps the browse categories response.
type categoriesResp struct {
	Categories Paged[Category] `json:"categories"`
}

// newReleasesResp wraps the new releases response.
type newReleasesResp struct {
	Albums Paged[AlbumSimple] `json:"albums"`
}
