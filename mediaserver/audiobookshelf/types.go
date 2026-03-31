package audiobookshelf

// Library represents an Audiobookshelf library.
type Library struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Folders    []Folder        `json:"folders,omitempty"`
	Icon       string          `json:"icon,omitempty"`
	MediaType  string          `json:"mediaType"`
	Provider   string          `json:"provider,omitempty"`
	CreatedAt  int64           `json:"createdAt,omitempty"`
	LastUpdate int64           `json:"lastUpdate,omitempty"`
	Settings   LibrarySettings `json:"settings,omitempty"`
}

// LibrarySettings represents library-level settings.
type LibrarySettings struct {
	CoverAspectRatio          int    `json:"coverAspectRatio,omitempty"`
	DisableWatcher            bool   `json:"disableWatcher,omitempty"`
	SkipMatchingMediaWithASIN bool   `json:"skipMatchingMediaWithAsin,omitempty"`
	SkipMatchingMediaWithISBN bool   `json:"skipMatchingMediaWithIsbn,omitempty"`
	AutoScanCronExpression    string `json:"autoScanCronExpression,omitempty"`
}

// Folder represents a folder in a library.
type Folder struct {
	ID       string `json:"id"`
	FullPath string `json:"fullPath"`
}

// LibraryItem represents an item (audiobook or podcast) in a library.
type LibraryItem struct {
	ID          string `json:"id"`
	INO         string `json:"ino,omitempty"`
	LibraryID   string `json:"libraryId"`
	FolderID    string `json:"folderId,omitempty"`
	Path        string `json:"path,omitempty"`
	RelPath     string `json:"relPath,omitempty"`
	IsFile      bool   `json:"isFile,omitempty"`
	MTIMEMS     int64  `json:"mtimeMs,omitempty"`
	CTIMEMS     int64  `json:"ctimeMs,omitempty"`
	BirthtimeMS int64  `json:"birthtimeMs,omitempty"`
	AddedAt     int64  `json:"addedAt,omitempty"`
	UpdatedAt   int64  `json:"updatedAt,omitempty"`
	IsMissing   bool   `json:"isMissing,omitempty"`
	IsInvalid   bool   `json:"isInvalid,omitempty"`
	MediaType   string `json:"mediaType,omitempty"`
	Media       Media  `json:"media,omitempty"`
	NumFiles    int    `json:"numFiles,omitempty"`
	Size        int64  `json:"size,omitempty"`
}

// Media represents the media content of a library item.
type Media struct {
	Metadata    Metadata    `json:"metadata,omitempty"`
	CoverPath   string      `json:"coverPath,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	AudioFiles  []AudioFile `json:"audioFiles,omitempty"`
	Chapters    []Chapter   `json:"chapters,omitempty"`
	Duration    float64     `json:"duration,omitempty"`
	Size        int64       `json:"size,omitempty"`
	NumTracks   int         `json:"numTracks,omitempty"`
	NumChapters int         `json:"numChapters,omitempty"`
}

// Metadata represents metadata for a media item.
type Metadata struct {
	Title         string   `json:"title,omitempty"`
	Subtitle      string   `json:"subtitle,omitempty"`
	Authors       []Author `json:"authors,omitempty"`
	Narrators     []string `json:"narrators,omitempty"`
	Series        []Series `json:"series,omitempty"`
	Genres        []string `json:"genres,omitempty"`
	PublishedYear string   `json:"publishedYear,omitempty"`
	Publisher     string   `json:"publisher,omitempty"`
	Description   string   `json:"description,omitempty"`
	ISBN          string   `json:"isbn,omitempty"`
	ASIN          string   `json:"asin,omitempty"`
	Language      string   `json:"language,omitempty"`
	Explicit      bool     `json:"explicit,omitempty"`
}

// Author represents an audiobook author.
type Author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Series represents a book series.
type Series struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Sequence string `json:"sequence,omitempty"`
}

// AudioFile represents an audio file in a library item.
type AudioFile struct {
	Index    int     `json:"index"`
	INO      string  `json:"ino,omitempty"`
	Metadata any     `json:"metadata,omitempty"`
	AddedAt  int64   `json:"addedAt,omitempty"`
	Duration float64 `json:"duration,omitempty"`
	MimeType string  `json:"mimeType,omitempty"`
}

// Chapter represents a chapter in an audiobook.
type Chapter struct {
	ID    int     `json:"id"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Title string  `json:"title"`
}

// User represents an Audiobookshelf user.
type User struct {
	ID                  string   `json:"id"`
	Username            string   `json:"username"`
	Type                string   `json:"type"`
	Token               string   `json:"token,omitempty"`
	IsActive            bool     `json:"isActive"`
	IsLocked            bool     `json:"isLocked"`
	LastSeen            int64    `json:"lastSeen,omitempty"`
	CreatedAt           int64    `json:"createdAt,omitempty"`
	LibrariesAccessible []string `json:"librariesAccessible,omitempty"`
	ItemTagsSelected    []string `json:"itemTagsSelected,omitempty"`
}

// PlaybackSession represents an active listening session.
type PlaybackSession struct {
	ID            string  `json:"id"`
	UserID        string  `json:"userId"`
	LibraryItemID string  `json:"libraryItemId"`
	EpisodeID     string  `json:"episodeId,omitempty"`
	MediaType     string  `json:"mediaType"`
	DisplayTitle  string  `json:"displayTitle,omitempty"`
	DisplayAuthor string  `json:"displayAuthor,omitempty"`
	Duration      float64 `json:"duration"`
	PlayMethod    int     `json:"playMethod"`
	StartTime     float64 `json:"startTime"`
	CurrentTime   float64 `json:"currentTime"`
	StartedAt     int64   `json:"startedAt,omitempty"`
	UpdatedAt     int64   `json:"updatedAt,omitempty"`
}

// ServerInfo represents Audiobookshelf server information.
type ServerInfo struct {
	IsInit        bool   `json:"isInit"`
	Language      string `json:"language,omitempty"`
	ConfigPath    string `json:"ConfigPath,omitempty"`
	MetadataPath  string `json:"MetadataPath,omitempty"`
	Source        string `json:"Source,omitempty"`
	Version       string `json:"version"`
	LatestVersion string `json:"latestVersion,omitempty"`
	HasUpdate     bool   `json:"hasUpdate"`
}

// LibraryItemsResponse wraps a paginated library items response.
type LibraryItemsResponse struct {
	Results []LibraryItem `json:"results"`
	Total   int           `json:"total"`
	Limit   int           `json:"limit"`
	Page    int           `json:"page"`
}

// Collection represents a collection of library items.
type Collection struct {
	ID          string        `json:"id"`
	LibraryID   string        `json:"libraryId"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Books       []LibraryItem `json:"books,omitempty"`
	CreatedAt   int64         `json:"createdAt,omitempty"`
	UpdatedAt   int64         `json:"updatedAt,omitempty"`
}

// MediaProgress represents user progress for a media item.
type MediaProgress struct {
	ID               string  `json:"id"`
	LibraryItemID    string  `json:"libraryItemId"`
	EpisodeID        string  `json:"episodeId,omitempty"`
	Duration         float64 `json:"duration"`
	Progress         float64 `json:"progress"`
	CurrentTime      float64 `json:"currentTime"`
	IsFinished       bool    `json:"isFinished"`
	HideFromContinue bool    `json:"hideFromContinueListening"`
	LastUpdate       int64   `json:"lastUpdate,omitempty"`
	StartedAt        int64   `json:"startedAt,omitempty"`
	FinishedAt       int64   `json:"finishedAt,omitempty"`
}

// LoginResponse represents the response from the login endpoint.
type LoginResponse struct {
	User User `json:"user"`
}
