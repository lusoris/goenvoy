package fanart

// Image is a basic artwork image returned by Fanart.tv.
type Image struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Lang  string `json:"lang,omitempty"`
	Likes string `json:"likes"`
}

// SeasonImage is an artwork image associated with a specific season.
type SeasonImage struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Lang   string `json:"lang,omitempty"`
	Likes  string `json:"likes"`
	Season string `json:"season"`
}

// DiscImage is a disc/CD art image with disc number and type info.
type DiscImage struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Lang     string `json:"lang,omitempty"`
	Likes    string `json:"likes"`
	Disc     string `json:"disc"`
	DiscType string `json:"disc_type,omitempty"`
}

// CDArt is a CD disc art image for a music album.
type CDArt struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Likes string `json:"likes"`
	Disc  string `json:"disc"`
	Size  string `json:"size"`
}

// AlbumImages holds artwork for a single music album, keyed by MusicBrainz
// release group ID in the parent map.
type AlbumImages struct {
	AlbumCover []Image `json:"albumcover,omitempty"`
	CDArt      []CDArt `json:"cdart,omitempty"`
}

// LabelImage is a music label logo image.
type LabelImage struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Color string `json:"colour"` //nolint:misspell // API uses British spelling.
	Likes string `json:"likes"`
}

// MovieImages contains all fan artwork for a movie.
type MovieImages struct {
	Name            string      `json:"name"`
	TMDbID          string      `json:"tmdb_id"`
	IMDbID          string      `json:"imdb_id"`
	HDMovieLogo     []Image     `json:"hdmovielogo,omitempty"`
	MovieDisc       []DiscImage `json:"moviedisc,omitempty"`
	MovieLogo       []Image     `json:"movielogo,omitempty"`
	MoviePoster     []Image     `json:"movieposter,omitempty"`
	HDMovieClearArt []Image     `json:"hdmovieclearart,omitempty"`
	MovieArt        []Image     `json:"movieart,omitempty"`
	MovieBackground []Image     `json:"moviebackground,omitempty"`
	MovieBanner     []Image     `json:"moviebanner,omitempty"`
	MovieThumb      []Image     `json:"moviethumb,omitempty"`
}

// ShowImages contains all fan artwork for a TV show.
type ShowImages struct {
	Name           string        `json:"name"`
	TheTVDBID      string        `json:"thetvdb_id"`
	ClearLogo      []Image       `json:"clearlogo,omitempty"`
	HDTVLogo       []Image       `json:"hdtvlogo,omitempty"`
	ClearArt       []Image       `json:"clearart,omitempty"`
	ShowBackground []SeasonImage `json:"showbackground,omitempty"`
	TVThumb        []Image       `json:"tvthumb,omitempty"`
	SeasonPoster   []Image       `json:"seasonposter,omitempty"`
	SeasonThumb    []SeasonImage `json:"seasonthumb,omitempty"`
	HDClearArt     []Image       `json:"hdclearart,omitempty"`
	TVBanner       []Image       `json:"tvbanner,omitempty"`
	CharacterArt   []Image       `json:"characterart,omitempty"`
	TVPoster       []Image       `json:"tvposter,omitempty"`
	SeasonBanner   []SeasonImage `json:"seasonbanner,omitempty"`
}

// ArtistImages contains all fan artwork for a music artist.
type ArtistImages struct {
	Name             string                 `json:"name"`
	MBID             string                 `json:"mbid_id"`
	ArtistBackground []Image                `json:"artistbackground,omitempty"`
	ArtistThumb      []Image                `json:"artistthumb,omitempty"`
	MusicLogo        []Image                `json:"musiclogo,omitempty"`
	HDMusicLogo      []Image                `json:"hdmusiclogo,omitempty"`
	MusicBanner      []Image                `json:"musicbanner,omitempty"`
	Albums           map[string]AlbumImages `json:"albums,omitempty"`
}

// AlbumImagesResponse is the response for a specific album lookup.
type AlbumImagesResponse struct {
	Name   string                 `json:"name"`
	MBID   string                 `json:"mbid_id"`
	Albums map[string]AlbumImages `json:"albums,omitempty"`
}

// LabelImages contains artwork for a music label.
type LabelImages struct {
	Name       string       `json:"name"`
	ID         string       `json:"id"`
	MusicLabel []LabelImage `json:"musiclabel,omitempty"`
}

// LatestMovie represents a movie entry from the latest-movies endpoint.
type LatestMovie struct {
	TMDbID      string `json:"tmdb_id"`
	IMDbID      string `json:"imdb_id"`
	Name        string `json:"name"`
	NewImages   string `json:"new_images"`
	TotalImages string `json:"total_images"`
}

// LatestShow represents a TV show entry from the latest-shows endpoint.
type LatestShow struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	NewImages   string `json:"new_images"`
	TotalImages string `json:"total_images"`
}

// LatestArtist represents a music artist entry from the latest-artists endpoint.
type LatestArtist struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	NewImages   string `json:"new_images"`
	TotalImages string `json:"total_images"`
}
