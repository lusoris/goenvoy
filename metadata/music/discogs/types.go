package discogs

// Release represents a Discogs release (album/single/EP).
type Release struct {
	ID                int         `json:"id"`
	Title             string      `json:"title"`
	Year              int         `json:"year,omitempty"`
	ResourceURL       string      `json:"resource_url,omitempty"`
	URI               string      `json:"uri,omitempty"`
	Artists           []ArtistRef `json:"artists,omitempty"`
	ArtistsSort       string      `json:"artists_sort,omitempty"`
	Labels            []LabelRef  `json:"labels,omitempty"`
	Formats           []Format    `json:"formats,omitempty"`
	DataQuality       string      `json:"data_quality,omitempty"`
	Community         *Community  `json:"community,omitempty"`
	DateAdded         string      `json:"date_added,omitempty"`
	DateChanged       string      `json:"date_changed,omitempty"`
	NumForSale        int         `json:"num_for_sale,omitempty"`
	LowestPrice       float64     `json:"lowest_price,omitempty"`
	MasterID          int         `json:"master_id,omitempty"`
	MasterURL         string      `json:"master_url,omitempty"`
	Country           string      `json:"country,omitempty"`
	Released          string      `json:"released,omitempty"`
	Notes             string      `json:"notes,omitempty"`
	ReleasedFormatted string      `json:"released_formatted,omitempty"`
	Genres            []string    `json:"genres,omitempty"`
	Styles            []string    `json:"styles,omitempty"`
	Tracklist         []Track     `json:"tracklist,omitempty"`
	ExtraArtists      []ArtistRef `json:"extraartists,omitempty"`
	Images            []Image     `json:"images,omitempty"`
	Thumb             string      `json:"thumb,omitempty"`
	Videos            []Video     `json:"videos,omitempty"`
}

// Artist represents a Discogs artist.
type Artist struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	RealName       string   `json:"realname,omitempty"`
	ResourceURL    string   `json:"resource_url,omitempty"`
	URI            string   `json:"uri,omitempty"`
	Profile        string   `json:"profile,omitempty"`
	URLs           []string `json:"urls,omitempty"`
	NameVariations []string `json:"namevariations,omitempty"`
	Members        []Member `json:"members,omitempty"`
	Groups         []Member `json:"groups,omitempty"`
	Images         []Image  `json:"images,omitempty"`
	DataQuality    string   `json:"data_quality,omitempty"`
}

// ArtistRef is a reference to an artist within a release.
type ArtistRef struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ANV         string `json:"anv,omitempty"`
	Join        string `json:"join,omitempty"`
	Role        string `json:"role,omitempty"`
	ResourceURL string `json:"resource_url,omitempty"`
	Tracks      string `json:"tracks,omitempty"`
}

// Label represents a Discogs label.
type Label struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	ResourceURL string     `json:"resource_url,omitempty"`
	URI         string     `json:"uri,omitempty"`
	Profile     string     `json:"profile,omitempty"`
	ContactInfo string     `json:"contact_info,omitempty"`
	URLs        []string   `json:"urls,omitempty"`
	Images      []Image    `json:"images,omitempty"`
	DataQuality string     `json:"data_quality,omitempty"`
	Sublabels   []LabelRef `json:"sublabels,omitempty"`
	ParentLabel *LabelRef  `json:"parent_label,omitempty"`
}

// LabelRef is a reference to a label.
type LabelRef struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CatNo       string `json:"catno,omitempty"`
	EntityType  string `json:"entity_type,omitempty"`
	ResourceURL string `json:"resource_url,omitempty"`
}

// MasterRelease represents a master release (canonical version).
type MasterRelease struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	MainRelease int         `json:"main_release,omitempty"`
	Year        int         `json:"year,omitempty"`
	ResourceURL string      `json:"resource_url,omitempty"`
	URI         string      `json:"uri,omitempty"`
	Artists     []ArtistRef `json:"artists,omitempty"`
	Genres      []string    `json:"genres,omitempty"`
	Styles      []string    `json:"styles,omitempty"`
	Tracklist   []Track     `json:"tracklist,omitempty"`
	Images      []Image     `json:"images,omitempty"`
	Videos      []Video     `json:"videos,omitempty"`
	DataQuality string      `json:"data_quality,omitempty"`
}

// Track represents a track on a release.
type Track struct {
	Position     string      `json:"position,omitempty"`
	Type         string      `json:"type_,omitempty"`
	Title        string      `json:"title"`
	Duration     string      `json:"duration,omitempty"`
	ExtraArtists []ArtistRef `json:"extraartists,omitempty"`
}

// Format describes the physical format of a release.
type Format struct {
	Name         string   `json:"name"`
	Qty          string   `json:"qty,omitempty"`
	Text         string   `json:"text,omitempty"`
	Descriptions []string `json:"descriptions,omitempty"`
}

// Image represents an image resource.
type Image struct {
	Type        string `json:"type"`
	URI         string `json:"uri,omitempty"`
	ResourceURL string `json:"resource_url,omitempty"`
	URI150      string `json:"uri150,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
}

// Video represents a video resource.
type Video struct {
	URI         string `json:"uri"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	Embed       bool   `json:"embed,omitempty"`
}

// Community contains community data for a release.
type Community struct {
	Have        int     `json:"have"`
	Want        int     `json:"want"`
	Rating      *Rating `json:"rating,omitempty"`
	Status      string  `json:"status,omitempty"`
	DataQuality string  `json:"data_quality,omitempty"`
}

// Rating represents community rating.
type Rating struct {
	Count   int     `json:"count"`
	Average float64 `json:"average"`
}

// Member represents a group member or group association.
type Member struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Active      bool   `json:"active,omitempty"`
	ResourceURL string `json:"resource_url,omitempty"`
}

// SearchResponse is the response from a database search.
type SearchResponse struct {
	Pagination Pagination     `json:"pagination"`
	Results    []SearchResult `json:"results"`
}

// Pagination contains pagination info.
type Pagination struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	PerPage int `json:"per_page"`
	Items   int `json:"items"`
}

// SearchResult represents a single search result.
type SearchResult struct {
	ID          int      `json:"id"`
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Thumb       string   `json:"thumb,omitempty"`
	CoverImage  string   `json:"cover_image,omitempty"`
	ResourceURL string   `json:"resource_url,omitempty"`
	URI         string   `json:"uri,omitempty"`
	Country     string   `json:"country,omitempty"`
	Year        string   `json:"year,omitempty"`
	Genre       []string `json:"genre,omitempty"`
	Style       []string `json:"style,omitempty"`
	Format      []string `json:"format,omitempty"`
	Label       []string `json:"label,omitempty"`
	Barcode     []string `json:"barcode,omitempty"`
	CatNo       string   `json:"catno,omitempty"`
}
