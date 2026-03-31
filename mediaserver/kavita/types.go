package kavita

import "fmt"

// Library represents a Kavita library.
type Library struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        int    `json:"type"`
	FolderPath  string `json:"folderPath"`
	LastScanned string `json:"lastScanned"`
}

// Series represents a comic, manga, or book series.
type Series struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	OriginalName   string `json:"originalName"`
	SortName       string `json:"sortName"`
	Pages          int    `json:"pages"`
	LibraryID      int    `json:"libraryId"`
	LibraryName    string `json:"libraryName"`
	Created        string `json:"created"`
	Format         int    `json:"format"`
	LatestReadDate string `json:"latestReadDate"`
	PagesRead      int    `json:"pagesRead"`
}

// Volume represents a volume within a series.
type Volume struct {
	ID        int    `json:"id"`
	Number    int    `json:"number"`
	Name      string `json:"name"`
	Pages     int    `json:"pages"`
	PagesRead int    `json:"pagesRead"`
	SeriesID  int    `json:"seriesId"`
	Created   string `json:"created"`
}

// Chapter represents a chapter within a volume.
type Chapter struct {
	ID        int    `json:"id"`
	Range     string `json:"range"`
	Number    string `json:"number"`
	Pages     int    `json:"pages"`
	IsSpecial bool   `json:"isSpecial"`
	Title     string `json:"title"`
	VolumeID  int    `json:"volumeId"`
	Created   string `json:"created"`
}

// Collection represents a collection of series.
type Collection struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Summary         string `json:"summary"`
	NormalizedTitle string `json:"normalizedTitle"`
	Promoted        bool   `json:"promoted"`
}

// ReadingList represents a user-created reading list.
type ReadingList struct {
	ID       int               `json:"id"`
	Title    string            `json:"title"`
	Summary  string            `json:"summary"`
	Promoted bool              `json:"promoted"`
	Items    []ReadingListItem `json:"items"`
}

// ReadingListItem represents an item within a reading list.
type ReadingListItem struct {
	ID        int `json:"id"`
	SeriesID  int `json:"seriesId"`
	ChapterID int `json:"chapterId"`
	Order     int `json:"order"`
}

// User represents a Kavita user.
type User struct {
	ID         int      `json:"id"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	Created    string   `json:"created"`
	LastActive string   `json:"lastActive"`
	Roles      []string `json:"roles"`
}

// ServerInfo represents Kavita server information.
type ServerInfo struct {
	Version   string `json:"version"`
	InstallID string `json:"installId"`
	Os        string `json:"os"`
}

// SearchResult represents the result of a search query.
type SearchResult struct {
	Series       []Series      `json:"series"`
	Collections  []Collection  `json:"collections"`
	ReadingLists []ReadingList `json:"readingLists"`
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("kavita: %s: %s", e.Status, e.Body)
}
