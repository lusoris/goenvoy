package navidrome

// subsonicResponse is the top-level response envelope for all Subsonic API calls.
type subsonicResponse struct {
	Response responseBody `json:"subsonic-response"`
}

// responseBody contains the status and data fields of a Subsonic response.
type responseBody struct {
	Status        string         `json:"status"`
	Version       string         `json:"version"`
	Type          string         `json:"type,omitempty"`
	ServerVersion string         `json:"serverVersion,omitempty"`
	OpenSubsonic  bool           `json:"openSubsonic,omitempty"`
	Error         *SubsonicError `json:"error,omitempty"`

	// Data fields — each endpoint populates one of these.
	Artists      *ArtistsID3    `json:"artists,omitempty"`
	Artist       *ArtistID3     `json:"artist,omitempty"`
	Album        *AlbumID3      `json:"album,omitempty"`
	AlbumList2   *AlbumList2    `json:"albumList2,omitempty"`
	Song         *Song          `json:"song,omitempty"`
	RandomSongs  *Songs         `json:"randomSongs,omitempty"`
	TopSongs     *TopSongs      `json:"topSongs,omitempty"`
	SearchResult *SearchResult3 `json:"searchResult3,omitempty"`
	Playlists    *Playlists     `json:"playlists,omitempty"`
	Playlist     *Playlist      `json:"playlist,omitempty"`
	NowPlaying   *NowPlaying    `json:"nowPlaying,omitempty"`
	Genres       *Genres        `json:"genres,omitempty"`
	ScanStatus   *ScanStatus    `json:"scanStatus,omitempty"`
	Starred2     *Starred2      `json:"starred2,omitempty"`
}

// SubsonicError represents an error returned by the Subsonic API.
type SubsonicError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *SubsonicError) Error() string {
	return e.Message
}

// ArtistsID3 represents the top-level artists index.
type ArtistsID3 struct {
	Index           []IndexID3 `json:"index,omitempty"`
	IgnoredArticles string     `json:"ignoredArticles,omitempty"`
}

// IndexID3 represents an alphabetical index of artists.
type IndexID3 struct {
	Name   string      `json:"name"`
	Artist []ArtistID3 `json:"artist,omitempty"`
}

// ArtistID3 represents an artist.
type ArtistID3 struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	CoverArt       string     `json:"coverArt,omitempty"`
	ArtistImageURL string     `json:"artistImageUrl,omitempty"`
	AlbumCount     int        `json:"albumCount,omitempty"`
	Starred        string     `json:"starred,omitempty"`
	MusicBrainzID  string     `json:"musicBrainzId,omitempty"`
	Album          []AlbumID3 `json:"album,omitempty"`
}

// AlbumID3 represents an album.
type AlbumID3 struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Artist        string `json:"artist,omitempty"`
	ArtistID      string `json:"artistId,omitempty"`
	CoverArt      string `json:"coverArt,omitempty"`
	SongCount     int    `json:"songCount"`
	Duration      int    `json:"duration"`
	Created       string `json:"created,omitempty"`
	Year          int    `json:"year,omitempty"`
	Genre         string `json:"genre,omitempty"`
	Starred       string `json:"starred,omitempty"`
	MusicBrainzID string `json:"musicBrainzId,omitempty"`
	Song          []Song `json:"song,omitempty"`
}

// AlbumList2 wraps a list of albums.
type AlbumList2 struct {
	Album []AlbumID3 `json:"album,omitempty"`
}

// Song represents a song/track.
type Song struct {
	ID            string `json:"id"`
	Parent        string `json:"parent,omitempty"`
	IsDir         bool   `json:"isDir,omitempty"`
	Title         string `json:"title"`
	Album         string `json:"album,omitempty"`
	Artist        string `json:"artist,omitempty"`
	Track         int    `json:"track,omitempty"`
	Year          int    `json:"year,omitempty"`
	Genre         string `json:"genre,omitempty"`
	CoverArt      string `json:"coverArt,omitempty"`
	Size          int64  `json:"size,omitempty"`
	ContentType   string `json:"contentType,omitempty"`
	Suffix        string `json:"suffix,omitempty"`
	Duration      int    `json:"duration,omitempty"`
	BitRate       int    `json:"bitRate,omitempty"`
	Path          string `json:"path,omitempty"`
	DiscNumber    int    `json:"discNumber,omitempty"`
	Created       string `json:"created,omitempty"`
	AlbumID       string `json:"albumId,omitempty"`
	ArtistID      string `json:"artistId,omitempty"`
	Type          string `json:"type,omitempty"`
	Starred       string `json:"starred,omitempty"`
	MusicBrainzID string `json:"musicBrainzId,omitempty"`
}

// Songs wraps a list of songs.
type Songs struct {
	Song []Song `json:"song,omitempty"`
}

// TopSongs wraps a list of top songs.
type TopSongs struct {
	Song []Song `json:"song,omitempty"`
}

// SearchResult3 represents search results from search3.
type SearchResult3 struct {
	Artist []ArtistID3 `json:"artist,omitempty"`
	Album  []AlbumID3  `json:"album,omitempty"`
	Song   []Song      `json:"song,omitempty"`
}

// Playlists wraps a list of playlists.
type Playlists struct {
	Playlist []Playlist `json:"playlist,omitempty"`
}

// Playlist represents a playlist.
type Playlist struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Comment   string `json:"comment,omitempty"`
	Owner     string `json:"owner,omitempty"`
	Public    bool   `json:"public,omitempty"`
	SongCount int    `json:"songCount"`
	Duration  int    `json:"duration"`
	Created   string `json:"created,omitempty"`
	Changed   string `json:"changed,omitempty"`
	CoverArt  string `json:"coverArt,omitempty"`
	Entry     []Song `json:"entry,omitempty"`
}

// NowPlaying wraps a list of now-playing entries.
type NowPlaying struct {
	Entry []NowPlayingEntry `json:"entry,omitempty"`
}

// NowPlayingEntry represents a currently playing song.
type NowPlayingEntry struct {
	Song
	Username   string `json:"username,omitempty"`
	MinutesAgo int    `json:"minutesAgo,omitempty"`
	PlayerID   int    `json:"playerId,omitempty"`
	PlayerName string `json:"playerName,omitempty"`
}

// Genres wraps a list of genres.
type Genres struct {
	Genre []Genre `json:"genre,omitempty"`
}

// Genre represents a music genre.
type Genre struct {
	SongCount  int    `json:"songCount"`
	AlbumCount int    `json:"albumCount"`
	Value      string `json:"value"`
}

// ScanStatus represents the media library scan status.
type ScanStatus struct {
	Scanning    bool   `json:"scanning"`
	Count       int64  `json:"count,omitempty"`
	FolderCount int    `json:"folderCount,omitempty"`
	LastScan    string `json:"lastScan,omitempty"`
}

// Starred2 represents starred items.
type Starred2 struct {
	Artist []ArtistID3 `json:"artist,omitempty"`
	Album  []AlbumID3  `json:"album,omitempty"`
	Song   []Song      `json:"song,omitempty"`
}
