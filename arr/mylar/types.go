package mylar

// Comic represents a comic series in Mylar3.
type Comic struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Year        string `json:"year"`
	Publisher   string `json:"publisher"`
	LatestIssue string `json:"latestIssue"`
	TotalIssues int    `json:"totalIssues"`
	ComicImage  string `json:"comicImage"`
}

// Issue represents a single comic issue.
type Issue struct {
	Id          string `json:"id"`
	IssueName   string `json:"issueName"`
	IssueNumber string `json:"issueNumber"`
	Status      string `json:"status"`
	ComicId     string `json:"comicId"`
	ReleaseDate string `json:"releaseDate"`
}

// Upcoming represents an upcoming comic issue.
type Upcoming struct {
	Id          string `json:"id"`
	IssueName   string `json:"issueName"`
	IssueNumber string `json:"issueNumber"`
	Status      string `json:"status"`
	ComicId     string `json:"comicId"`
	ReleaseDate string `json:"releaseDate"`
	ComicName   string `json:"comicName"`
}

// WantedIssue represents a wanted comic issue.
type WantedIssue struct {
	Id          string `json:"id"`
	IssueName   string `json:"issueName"`
	IssueNumber string `json:"issueNumber"`
	Status      string `json:"status"`
	ComicId     string `json:"comicId"`
	ReleaseDate string `json:"releaseDate"`
}

// HistoryEntry represents an entry in the download history.
type HistoryEntry struct {
	Id          string `json:"id"`
	ComicName   string `json:"comicName"`
	IssueNumber string `json:"issueNumber"`
	Date        string `json:"date"`
	Status      string `json:"status"`
	Provider    string `json:"provider"`
}

// LogEntry represents a single log line.
type LogEntry struct {
	Message   string `json:"message"`
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
}

// SearchResult represents a comic found via search.
type SearchResult struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Year      string `json:"year"`
	Publisher string `json:"publisher"`
	Issues    int    `json:"issues"`
	Image     string `json:"image"`
}

// VersionInfo contains Mylar3 version information.
type VersionInfo struct {
	Version       string `json:"version"`
	LatestVersion string `json:"latestVersion"`
	Commits       string `json:"commits"`
}

// StoryArc represents a story arc containing issues.
type StoryArc struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Publisher string  `json:"publisher"`
	Issues    []Issue `json:"issues"`
}

// ReadList represents a reading list.
type ReadList struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Issues []Issue `json:"issues"`
}

// Provider represents a search provider.
type Provider struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}
