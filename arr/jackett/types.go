package jackett

import "encoding/xml"

// SearchResult represents a single search result from a torznab query.
type SearchResult struct {
	Title                string
	GUID                 string
	Link                 string
	Comments             string
	PubDate              string
	Size                 int64
	Category             string
	CategoryDesc         string
	Seeders              int
	Peers                int
	InfoHash             string
	MagnetURL            string
	MinimumRatio         string
	MinimumSeedTime      string
	DownloadVolumeFactor string
	UploadVolumeFactor   string
	Indexer              string
}

// Capabilities represents the capabilities of a Jackett server or indexer.
type Capabilities struct {
	Server     ServerInfo
	Limits     Limits
	Searching  Searching
	Categories []Category
}

// ServerInfo contains server identification.
type ServerInfo struct {
	Title string
	Image string
}

// Limits describes search result limits.
type Limits struct {
	Max     int
	Default int
}

// Searching describes available search types.
type Searching struct {
	SearchAvailable      bool
	TVSearchAvailable    bool
	MovieSearchAvailable bool
	MusicSearchAvailable bool
	BookSearchAvailable  bool
}

// Category represents a torznab category.
type Category struct {
	ID            int
	Name          string
	SubCategories []SubCategory
}

// SubCategory represents a torznab subcategory.
type SubCategory struct {
	ID   int
	Name string
}

// Indexer represents a configured Jackett indexer.
type Indexer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Configured bool   `json:"configured"`
	Language   string `json:"language"`
	SiteLink   string `json:"site_link"`
}

// ServerConfig represents Jackett server configuration.
type ServerConfig struct {
	APIKey       string `json:"api_key"`
	BlackholeDir string `json:"blackholedir"`
	Port         int    `json:"port"`
	InstanceID   string `json:"instance_id"`
}

// rssResponse is the top-level XML structure for torznab search results.
type rssResponse struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

// rssChannel is the channel element in an RSS response.
type rssChannel struct {
	Items []rssItem `xml:"item"`
}

// rssItem is a single result item in the RSS feed.
type rssItem struct {
	Title    string    `xml:"title"`
	GUID     string    `xml:"guid"`
	Link     string    `xml:"link"`
	Comments string    `xml:"comments"`
	PubDate  string    `xml:"pubDate"`
	Size     int64     `xml:"size"`
	Attrs    []rssAttr `xml:",any"`
}

// rssAttr captures torznab:attr elements from the XML feed.
type rssAttr struct {
	XMLName xml.Name
	Name    string `xml:"name,attr"`
	Value   string `xml:"value,attr"`
}

// capsResponse is the XML structure for the capabilities endpoint.
type capsResponse struct {
	XMLName    xml.Name       `xml:"caps"`
	Server     capsServer     `xml:"server"`
	Limits     capsLimits     `xml:"limits"`
	Searching  capsSearching  `xml:"searching"`
	Categories capsCategories `xml:"categories"`
}

type capsServer struct {
	Title string `xml:"title,attr"`
	Image string `xml:"image,attr"`
}

type capsLimits struct {
	Max     int `xml:"max,attr"`
	Default int `xml:"default,attr"`
}

type capsSearching struct {
	Search      capsSearchType `xml:"search"`
	TVSearch    capsSearchType `xml:"tv-search"`
	MovieSearch capsSearchType `xml:"movie-search"`
	MusicSearch capsSearchType `xml:"music-search"`
	BookSearch  capsSearchType `xml:"book-search"`
}

type capsSearchType struct {
	Available string `xml:"available,attr"`
}

type capsCategories struct {
	Categories []capsCategory `xml:"category"`
}

type capsCategory struct {
	ID            int               `xml:"id,attr"`
	Name          string            `xml:"name,attr"`
	SubCategories []capsSubCategory `xml:"subcat"`
}

type capsSubCategory struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}
