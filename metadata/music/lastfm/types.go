package lastfm

// Artist represents a Last.fm artist.
type Artist struct {
	Name       string          `json:"name"`
	MBID       string          `json:"mbid,omitempty"`
	URL        string          `json:"url,omitempty"`
	Image      []Image         `json:"image,omitempty"`
	Streamable string          `json:"streamable,omitempty"`
	Stats      *ArtistStats    `json:"stats,omitempty"`
	Similar    *SimilarArtists `json:"similar,omitempty"`
	Tags       *Tags           `json:"tags,omitempty"`
	Bio        *Bio            `json:"bio,omitempty"`
}

// ArtistStats contains playcount and listener counts.
type ArtistStats struct {
	Listeners string `json:"listeners"`
	Playcount string `json:"playcount"`
}

// SimilarArtists wraps a list of similar artists.
type SimilarArtists struct {
	Artist []Artist `json:"artist,omitempty"`
}

// Album represents a Last.fm album.
type Album struct {
	Name      string  `json:"name"`
	Artist    string  `json:"artist,omitempty"`
	MBID      string  `json:"mbid,omitempty"`
	URL       string  `json:"url,omitempty"`
	Image     []Image `json:"image,omitempty"`
	Listeners string  `json:"listeners,omitempty"`
	Playcount string  `json:"playcount,omitempty"`
	Tracks    *Tracks `json:"tracks,omitempty"`
	Tags      *Tags   `json:"tags,omitempty"`
	Wiki      *Wiki   `json:"wiki,omitempty"`
}

// Track represents a Last.fm track.
type Track struct {
	Name       string  `json:"name"`
	Artist     any     `json:"artist,omitempty"` // Can be string or object.
	MBID       string  `json:"mbid,omitempty"`
	URL        string  `json:"url,omitempty"`
	Duration   string  `json:"duration,omitempty"`
	Listeners  string  `json:"listeners,omitempty"`
	Playcount  string  `json:"playcount,omitempty"`
	Streamable any     `json:"streamable,omitempty"`
	Image      []Image `json:"image,omitempty"`
	Album      *Album  `json:"album,omitempty"`
	TopTags    *Tags   `json:"toptags,omitempty"`
	Wiki       *Wiki   `json:"wiki,omitempty"`
}

// Tracks wraps a list of tracks.
type Tracks struct {
	Track []Track `json:"track,omitempty"`
}

// Tags wraps a list of tags.
type Tags struct {
	Tag []Tag `json:"tag,omitempty"`
}

// Tag represents a Last.fm tag/genre.
type Tag struct {
	Name  string `json:"name"`
	URL   string `json:"url,omitempty"`
	Count int    `json:"count,omitempty"`
}

// Image represents an image with a size designation.
type Image struct {
	Text string `json:"#text"`
	Size string `json:"size"`
}

// Bio represents an artist biography.
type Bio struct {
	Published string `json:"published,omitempty"`
	Summary   string `json:"summary,omitempty"`
	Content   string `json:"content,omitempty"`
}

// Wiki represents album/track wiki content.
type Wiki struct {
	Published string `json:"published,omitempty"`
	Summary   string `json:"summary,omitempty"`
	Content   string `json:"content,omitempty"`
}

// ChartArtists wraps a chart of top artists.
type ChartArtists struct {
	Artist []Artist `json:"artist,omitempty"`
}

// ChartTracks wraps a chart of top tracks.
type ChartTracks struct {
	Track []Track `json:"track,omitempty"`
}

// TopAlbums wraps an artist's top albums.
type TopAlbums struct {
	Album []Album `json:"album,omitempty"`
}

// TopTracks wraps an artist's top tracks.
type TopTracks struct {
	Track []Track `json:"track,omitempty"`
}

// TopTags wraps a list of top tags.
type TopTags struct {
	Tag []Tag `json:"tag,omitempty"`
}
