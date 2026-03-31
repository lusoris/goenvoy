package nzbhydra

import (
	"encoding/xml"
	"time"
)

// SearchResult represents a single search result from a Newznab query.
type SearchResult struct {
	Title       string
	GUID        string
	Link        string
	Comments    string
	PubDate     string
	Size        int64
	Category    string
	Description string
	Indexer     string
}

// TVSearchOptions holds optional parameters for a TV search.
type TVSearchOptions struct {
	Season   int
	Episode  int
	TVDBID   string
	IMDBID   string
	TMDBID   string
	TVMazeID string
	RID      string
}

// Capabilities represents the capabilities of an NZBHydra2 server.
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
	BookSearchAvailable  bool
}

// Category represents a Newznab category.
type Category struct {
	ID            int
	Name          string
	SubCategories []SubCategory
}

// SubCategory represents a Newznab subcategory.
type SubCategory struct {
	ID   int
	Name string
}

// StatsRequest holds parameters for the statistics endpoint.
type StatsRequest struct {
	After                        time.Time `json:"after"`
	Before                       time.Time `json:"before"`
	IncludeAverageResponseTimes  bool      `json:"includeAverageResponseTimes"`
	IncludeSearchesPerDayOfWeek  bool      `json:"includeSearchesPerDayOfWeek"`
	IncludeDownloadsPerDayOfWeek bool      `json:"includeDownloadsPerDayOfWeek"`
	IncludeIndexerAPIAccessStats bool      `json:"includeIndexerApiAccessStats"`
}

// StatsResponse holds statistics returned by the stats endpoint.
type StatsResponse struct {
	AvgResponseTimes      []AvgResponseTime `json:"avgResponseTimes"`
	IndexerAPIAccessStats []IndexerAPIStat  `json:"indexerApiAccessStats"`
	SearchesPerDayOfWeek  map[string]int    `json:"searchesPerDayOfWeek"`
	DownloadsPerDayOfWeek map[string]int    `json:"downloadsPerDayOfWeek"`
}

// AvgResponseTime holds average response time for an indexer.
type AvgResponseTime struct {
	Indexer string `json:"indexer"`
	Avg     int    `json:"avgResponseTime"`
}

// IndexerAPIStat holds API access statistics for an indexer.
type IndexerAPIStat struct {
	Indexer      string `json:"indexer"`
	Successful   int    `json:"successful"`
	Unsuccessful int    `json:"unsuccessful"`
}

// HistoryRequest holds parameters for history query endpoints.
type HistoryRequest struct {
	Page        int           `json:"page"`
	Limit       int           `json:"limit"`
	FilterModel []FilterEntry `json:"filterModel,omitempty"`
	SortModel   []SortEntry   `json:"sortModel,omitempty"`
}

// FilterEntry defines a filter for history queries.
type FilterEntry struct {
	Field string `json:"field"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// SortEntry defines a sort order for history queries.
type SortEntry struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// PagedResponse is a generic paginated response.
type PagedResponse[T any] struct {
	Content          []T  `json:"content"`
	TotalElements    int  `json:"totalElements"`
	TotalPages       int  `json:"totalPages"`
	First            bool `json:"first"`
	Last             bool `json:"last"`
	NumberOfElements int  `json:"numberOfElements"`
	Number           int  `json:"number"`
	Size             int  `json:"size"`
}

// SearchHistoryEntry represents a single search history record.
type SearchHistoryEntry struct {
	ID         int    `json:"id"`
	Source     string `json:"source"`
	SearchType string `json:"searchType"`
	Time       string `json:"time"`
	Query      string `json:"query"`
	Season     string `json:"season"`
	Episode    string `json:"episode"`
	IP         string `json:"ip"`
	UserAgent  string `json:"userAgent"`
}

// DownloadHistoryEntry represents a single download history record.
type DownloadHistoryEntry struct {
	ID            int            `json:"id"`
	SearchResult  DownloadResult `json:"searchResult"`
	NzbAccessType string         `json:"nzbAccessType"`
	AccessSource  string         `json:"accessSource"`
	Time          string         `json:"time"`
	Status        string         `json:"status"`
	Error         string         `json:"error"`
	Age           int            `json:"age"`
}

// DownloadResult holds info about the downloaded NZB.
type DownloadResult struct {
	Title   string `json:"title"`
	Indexer string `json:"indexer"`
	Link    string `json:"link"`
}

// IndexerStatus represents the status of an NZBHydra2 indexer.
type IndexerStatus struct {
	Indexer       string `json:"indexer"`
	State         string `json:"state"`
	Level         string `json:"level"`
	DisabledUntil string `json:"disabledUntil"`
	LastError     string `json:"lastError"`
}

// rssResponse is the top-level XML structure for Newznab search results.
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
	Title       string    `xml:"title"`
	GUID        string    `xml:"guid"`
	Link        string    `xml:"link"`
	Comments    string    `xml:"comments"`
	PubDate     string    `xml:"pubDate"`
	Size        int64     `xml:"size"`
	Description string    `xml:"description"`
	Attrs       []rssAttr `xml:",any"`
}

// rssAttr captures newznab:attr elements from the XML feed.
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
