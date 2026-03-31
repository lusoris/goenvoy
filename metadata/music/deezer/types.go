package deezer

// ArtistSimple is a simplified artist object.
type ArtistSimple struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Link    string `json:"link"`
	Picture string `json:"picture"`
	Type    string `json:"type"`
}

// Artist is a full Deezer artist object.
type Artist struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	PictureSmall  string `json:"picture_small"`
	PictureMedium string `json:"picture_medium"`
	PictureBig    string `json:"picture_big"`
	PictureXL     string `json:"picture_xl"`
	Tracklist     string `json:"tracklist"`
	Type          string `json:"type"`
	NbAlbum       int    `json:"nb_album"`
	NbFan         int    `json:"nb_fan"`
}

// AlbumSimple is a simplified album object.
type AlbumSimple struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	CoverSmall  string `json:"cover_small"`
	CoverMedium string `json:"cover_medium"`
	CoverBig    string `json:"cover_big"`
	CoverXL     string `json:"cover_xl"`
	Tracklist   string `json:"tracklist"`
	Type        string `json:"type"`
}

// Album is a full Deezer album object.
type Album struct {
	AlbumSimple
	Artist      ArtistSimple `json:"artist"`
	Genres      GenreList    `json:"genres"`
	Label       string       `json:"label"`
	Duration    int          `json:"duration"`
	NbTracks    int          `json:"nb_tracks"`
	Fans        int          `json:"fans"`
	ReleaseDate string       `json:"release_date"`
	RecordType  string       `json:"record_type"`
	Tracks      TrackList    `json:"tracks"`
}

// Track represents a Deezer track.
type Track struct {
	ID                    int          `json:"id"`
	Title                 string       `json:"title"`
	TitleShort            string       `json:"title_short"`
	TitleVersion          string       `json:"title_version"`
	Link                  string       `json:"link"`
	Duration              int          `json:"duration"`
	Rank                  int          `json:"rank"`
	ExplicitLyrics        bool         `json:"explicit_lyrics"`
	ExplicitContentLyrics int          `json:"explicit_content_lyrics"`
	ExplicitContentCover  int          `json:"explicit_content_cover"`
	Preview               string       `json:"preview"`
	MD5Image              string       `json:"md5_image"`
	Type                  string       `json:"type"`
	Artist                ArtistSimple `json:"artist"`
	Album                 AlbumSimple  `json:"album"`
}

// Genre represents a Deezer genre.
type Genre struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// GenreList wraps a list of genres.
type GenreList struct {
	Data []Genre `json:"data"`
}

// TrackList wraps a list of tracks.
type TrackList struct {
	Data []Track `json:"data"`
}

// AlbumList wraps a list of albums.
type AlbumList struct {
	Data []AlbumSimple `json:"data"`
}

// ArtistList wraps a list of artists.
type ArtistList struct {
	Data []Artist `json:"data"`
}

// PlaylistList wraps a list of playlists.
type PlaylistList struct {
	Data []Playlist `json:"data"`
}

// Playlist represents a Deezer playlist.
type Playlist struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Link        string      `json:"link"`
	Picture     string      `json:"picture"`
	Type        string      `json:"type"`
	NbTracks    int         `json:"nb_tracks"`
	Creator     CreatorInfo `json:"creator"`
}

// CreatorInfo represents a playlist creator.
type CreatorInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Chart contains chart data for tracks, albums, artists, and playlists.
type Chart struct {
	Tracks    TrackList    `json:"tracks"`
	Albums    AlbumList    `json:"albums"`
	Artists   ArtistList   `json:"artists"`
	Playlists PlaylistList `json:"playlists"`
}

// SearchResult contains track search results.
type SearchResult struct {
	Data  []Track `json:"data"`
	Total int     `json:"total"`
	Next  string  `json:"next"`
}

// AlbumSearchResult contains album search results.
type AlbumSearchResult struct {
	Data  []AlbumSimple `json:"data"`
	Total int           `json:"total"`
	Next  string        `json:"next"`
}

// ArtistSearchResult contains artist search results.
type ArtistSearchResult struct {
	Data  []Artist `json:"data"`
	Total int      `json:"total"`
	Next  string   `json:"next"`
}

// deezerError is the error structure in Deezer API responses.
type deezerError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// errorResp wraps a Deezer error response.
type errorResp struct {
	Error *deezerError `json:"error,omitempty"`
}

// artistTopResp wraps the artist top tracks response.
type artistTopResp struct {
	Data []Track `json:"data"`
}

// artistAlbumsResp wraps the artist albums response.
type artistAlbumsResp struct {
	Data []AlbumSimple `json:"data"`
}

// relatedArtistsResp wraps the related artists response.
type relatedArtistsResp struct {
	Data []Artist `json:"data"`
}

// albumTracksResp wraps the album tracks response.
type albumTracksResp struct {
	Data []Track `json:"data"`
}

// genresResp wraps the genres list response.
type genresResp struct {
	Data []Genre `json:"data"`
}
