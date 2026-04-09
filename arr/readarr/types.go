package readarr

// Author represents an author in Readarr.
type Author struct {
	ID                  int               `json:"id"`
	AuthorMetadataID    int               `json:"authorMetadataId"`
	Status              string            `json:"status"`
	Ended               bool              `json:"ended,omitempty"`
	AuthorName          string            `json:"authorName"`
	AuthorNameLastFirst string            `json:"authorNameLastFirst,omitempty"`
	ForeignAuthorID     string            `json:"foreignAuthorId"`
	TitleSlug           string            `json:"titleSlug,omitempty"`
	Overview            string            `json:"overview,omitempty"`
	Disambiguation      string            `json:"disambiguation,omitempty"`
	Links               []Link            `json:"links,omitempty"`
	NextBook            *Book             `json:"nextBook,omitempty"`
	LastBook            *Book             `json:"lastBook,omitempty"`
	Images              []Image           `json:"images,omitempty"`
	RemotePoster        string            `json:"remotePoster,omitempty"`
	Path                string            `json:"path"`
	QualityProfileID    int               `json:"qualityProfileId"`
	MetadataProfileID   int               `json:"metadataProfileId"`
	Monitored           bool              `json:"monitored"`
	MonitorNewItems     string            `json:"monitorNewItems,omitempty"`
	RootFolderPath      string            `json:"rootFolderPath,omitempty"`
	Folder              string            `json:"folder,omitempty"`
	Genres              []string          `json:"genres,omitempty"`
	CleanName           string            `json:"cleanName,omitempty"`
	SortName            string            `json:"sortName,omitempty"`
	SortNameLastFirst   string            `json:"sortNameLastFirst,omitempty"`
	Tags                []int             `json:"tags"`
	Added               string            `json:"added"`
	AddOptions          *AddAuthorOptions `json:"addOptions,omitempty"`
	Ratings             Ratings           `json:"ratings,omitempty"`
	Statistics          *AuthorStatistics `json:"statistics,omitempty"`
}

// Book represents a book in Readarr.
type Book struct {
	ID               int             `json:"id"`
	Title            string          `json:"title"`
	AuthorTitle      string          `json:"authorTitle,omitempty"`
	SeriesTitle      string          `json:"seriesTitle,omitempty"`
	Disambiguation   string          `json:"disambiguation,omitempty"`
	Overview         string          `json:"overview,omitempty"`
	AuthorID         int             `json:"authorId"`
	ForeignBookID    string          `json:"foreignBookId"`
	ForeignEditionID string          `json:"foreignEditionId,omitempty"`
	TitleSlug        string          `json:"titleSlug,omitempty"`
	Monitored        bool            `json:"monitored"`
	AnyEditionOk     bool            `json:"anyEditionOk"`
	Ratings          Ratings         `json:"ratings,omitempty"`
	ReleaseDate      string          `json:"releaseDate,omitempty"`
	PageCount        int             `json:"pageCount,omitempty"`
	Genres           []string        `json:"genres,omitempty"`
	Author           *Author         `json:"author,omitempty"`
	Images           []Image         `json:"images,omitempty"`
	Links            []Link          `json:"links,omitempty"`
	Statistics       *BookStatistics `json:"statistics,omitempty"`
	Added            string          `json:"added,omitempty"`
	AddOptions       *AddBookOptions `json:"addOptions,omitempty"`
	RemoteCover      string          `json:"remoteCover,omitempty"`
	Editions         []Edition       `json:"editions,omitempty"`
}

// BookFile represents a downloaded book file on disk.
type BookFile struct {
	ID                  int          `json:"id"`
	AuthorID            int          `json:"authorId"`
	BookID              int          `json:"bookId"`
	Path                string       `json:"path"`
	Size                int64        `json:"size"`
	DateAdded           string       `json:"dateAdded"`
	Quality             QualityModel `json:"quality"`
	QualityCutoffNotMet bool         `json:"qualityCutoffNotMet"`
}

// Edition represents a specific edition of a book (hardcover, paperback, ebook, etc.).
type Edition struct {
	ID               int     `json:"id"`
	BookID           int     `json:"bookId"`
	ForeignEditionID string  `json:"foreignEditionId"`
	TitleSlug        string  `json:"titleSlug,omitempty"`
	ISBN13           string  `json:"isbn13,omitempty"`
	ASIN             string  `json:"asin,omitempty"`
	Title            string  `json:"title"`
	Language         string  `json:"language,omitempty"`
	Overview         string  `json:"overview,omitempty"`
	Format           string  `json:"format,omitempty"`
	IsEbook          bool    `json:"isEbook"`
	Disambiguation   string  `json:"disambiguation,omitempty"`
	Publisher        string  `json:"publisher,omitempty"`
	PageCount        int     `json:"pageCount,omitempty"`
	ReleaseDate      string  `json:"releaseDate,omitempty"`
	Images           []Image `json:"images,omitempty"`
	Links            []Link  `json:"links,omitempty"`
	Ratings          Ratings `json:"ratings,omitempty"`
	Monitored        bool    `json:"monitored"`
	ManualAdd        bool    `json:"manualAdd"`
	RemoteCover      string  `json:"remoteCover,omitempty"`
}

// Series represents a book series (e.g. "Harry Potter").
type Series struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// Image represents a cover image for a media item.
type Image struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
	RemoteURL string `json:"remoteUrl,omitempty"`
}

// Link represents a web link for an author or book.
type Link struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

// Ratings holds community rating data.
type Ratings struct {
	Votes int     `json:"votes"`
	Value float64 `json:"value"`
}

// AddAuthorOptions controls behavior when adding a new author.
type AddAuthorOptions struct {
	Monitor               string   `json:"monitor"`
	BooksToMonitor        []string `json:"booksToMonitor,omitempty"`
	Monitored             bool     `json:"monitored"`
	SearchForMissingBooks bool     `json:"searchForMissingBooks"`
}

// AddBookOptions controls behavior when adding a new book.
type AddBookOptions struct {
	AddType          string `json:"addType,omitempty"`
	SearchForNewBook bool   `json:"searchForNewBook"`
}

// AuthorStatistics contains aggregate statistics for an author.
type AuthorStatistics struct {
	BookFileCount      int     `json:"bookFileCount"`
	BookCount          int     `json:"bookCount"`
	AvailableBookCount int     `json:"availableBookCount"`
	TotalBookCount     int     `json:"totalBookCount"`
	SizeOnDisk         int64   `json:"sizeOnDisk"`
	PercentOfBooks     float64 `json:"percentOfBooks"`
}

// BookStatistics contains file counts and size information for a book.
type BookStatistics struct {
	BookFileCount  int     `json:"bookFileCount"`
	BookCount      int     `json:"bookCount"`
	TotalBookCount int     `json:"totalBookCount"`
	SizeOnDisk     int64   `json:"sizeOnDisk"`
	PercentOfBooks float64 `json:"percentOfBooks"`
}

// QualityModel pairs a quality definition with its revision.
type QualityModel struct {
	Quality  Quality  `json:"quality"`
	Revision Revision `json:"revision"`
}

// Quality identifies a quality tier.
type Quality struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Revision tracks repack and version information for a release.
type Revision struct {
	Version  int  `json:"version"`
	Real     int  `json:"real"`
	IsRepack bool `json:"isRepack"`
}

// CustomFormat describes a custom format definition.
type CustomFormat struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// HistoryRecord represents an event in the download history.
type HistoryRecord struct {
	ID                  int               `json:"id"`
	BookID              int               `json:"bookId"`
	AuthorID            int               `json:"authorId"`
	SourceTitle         string            `json:"sourceTitle"`
	Quality             QualityModel      `json:"quality"`
	CustomFormats       []CustomFormat    `json:"customFormats,omitempty"`
	CustomFormatScore   int               `json:"customFormatScore"`
	QualityCutoffNotMet bool              `json:"qualityCutoffNotMet"`
	Date                string            `json:"date"`
	DownloadID          string            `json:"downloadId,omitempty"`
	EventType           string            `json:"eventType"`
	Data                map[string]string `json:"data,omitempty"`
}

// ParseResult contains the result of parsing a release title.
type ParseResult struct {
	ID             int             `json:"id"`
	Title          string          `json:"title"`
	ParsedBookInfo *ParsedBookInfo `json:"parsedBookInfo,omitempty"`
	Author         *Author         `json:"author,omitempty"`
	Books          []Book          `json:"books,omitempty"`
}

// ParsedBookInfo holds the structured data extracted from a release title.
type ParsedBookInfo struct {
	BookTitle        string       `json:"bookTitle,omitempty"`
	AuthorName       string       `json:"authorName,omitempty"`
	Quality          QualityModel `json:"quality"`
	ReleaseDate      string       `json:"releaseDate,omitempty"`
	Discography      bool         `json:"discography"`
	DiscographyStart int          `json:"discographyStart,omitempty"`
	DiscographyEnd   int          `json:"discographyEnd,omitempty"`
	ReleaseGroup     string       `json:"releaseGroup,omitempty"`
	ReleaseHash      string       `json:"releaseHash,omitempty"`
	ReleaseVersion   string       `json:"releaseVersion,omitempty"`
}

// AuthorEditorResource is used for batch editing or deleting authors.
type AuthorEditorResource struct {
	AuthorIDs         []int  `json:"authorIds"`
	Monitored         *bool  `json:"monitored,omitempty"`
	MonitorNewItems   string `json:"monitorNewItems,omitempty"`
	QualityProfileID  *int   `json:"qualityProfileId,omitempty"`
	MetadataProfileID *int   `json:"metadataProfileId,omitempty"`
	RootFolderPath    string `json:"rootFolderPath,omitempty"`
	Tags              []int  `json:"tags,omitempty"`
	ApplyTags         string `json:"applyTags,omitempty"`
	MoveFiles         bool   `json:"moveFiles"`
	DeleteFiles       bool   `json:"deleteFiles"`
}

// BookEditorResource is used for batch editing or deleting books.
type BookEditorResource struct {
	BookIDs                []int `json:"bookIds"`
	Monitored              *bool `json:"monitored,omitempty"`
	DeleteFiles            bool  `json:"deleteFiles"`
	AddImportListExclusion bool  `json:"addImportListExclusion"`
}

// BooksMonitoredResource is used to set the monitored status of books.
type BooksMonitoredResource struct {
	BookIDs   []int `json:"bookIds"`
	Monitored bool  `json:"monitored"`
}

// BookFileListResource is the request body for bulk book file operations.
type BookFileListResource struct {
	BookFileIDs []int `json:"bookFileIds"`
}

// ImportListExclusion represents an author excluded from import lists.
type ImportListExclusion struct {
	ID         int    `json:"id"`
	ForeignID  string `json:"foreignId,omitempty"`
	AuthorName string `json:"authorName,omitempty"`
}

// MetadataProfile describes a metadata profile for filtering book types.
type MetadataProfile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DevelopmentConfigResource holds development configuration.
type DevelopmentConfigResource struct {
	ID                 int    `json:"id"`
	MetadataSource     string `json:"metadataSource,omitempty"`
	ConsoleLogLevel    string `json:"consoleLogLevel,omitempty"`
	LogSQL             bool   `json:"logSql"`
	LogRotate          int    `json:"logRotate"`
	FilterSentryEvents bool   `json:"filterSentryEvents"`
}

// MetadataProviderConfigResource holds metadata provider configuration.
type MetadataProviderConfigResource struct {
	ID             int    `json:"id"`
	WriteAudioTags string `json:"writeAudioTags,omitempty"`
	ScrubAudioTags bool   `json:"scrubAudioTags"`
	WriteBookTags  string `json:"writeBookTags,omitempty"`
	UpdateCovers   bool   `json:"updateCovers"`
	EmbedMetadata  bool   `json:"embedMetadata"`
}

// BookshelfResource is the request body for the bookshelf endpoint.
type BookshelfResource struct {
	Authors           []BookshelfAuthorResource `json:"authors,omitempty"`
	MonitoringOptions MonitoringOptions         `json:"monitoringOptions"`
	MonitorNewItems   string                    `json:"monitorNewItems,omitempty"`
}

// BookshelfAuthorResource is an author entry in a bookshelf request.
type BookshelfAuthorResource struct {
	ID        int    `json:"id"`
	Monitored *bool  `json:"monitored,omitempty"`
	Books     []Book `json:"books,omitempty"`
}

// MonitoringOptions controls what to monitor when adding or updating.
type MonitoringOptions struct {
	Monitor        string   `json:"monitor,omitempty"`
	BooksToMonitor []string `json:"booksToMonitor,omitempty"`
	Monitored      bool     `json:"monitored"`
}

// RenameBookResource represents a rename preview for a book file.
type RenameBookResource struct {
	ID           int    `json:"id"`
	AuthorID     int    `json:"authorId"`
	BookID       int    `json:"bookId"`
	BookFileID   int    `json:"bookFileId"`
	ExistingPath string `json:"existingPath,omitempty"`
	NewPath      string `json:"newPath,omitempty"`
}

// RetagBookResource represents a retag preview for a book file.
type RetagBookResource struct {
	ID           int             `json:"id"`
	AuthorID     int             `json:"authorId"`
	BookID       int             `json:"bookId"`
	TrackNumbers []int           `json:"trackNumbers,omitempty"`
	BookFileID   int             `json:"bookFileId"`
	Path         string          `json:"path,omitempty"`
	Changes      []TagDifference `json:"changes,omitempty"`
}

// TagDifference represents a single tag change.
type TagDifference struct {
	Field    string `json:"field,omitempty"`
	OldValue string `json:"oldValue,omitempty"`
	NewValue string `json:"newValue,omitempty"`
}

// FileSystemResource represents a filesystem browse response.
type FileSystemResource struct {
	Directories []FileSystemEntry `json:"directories,omitempty"`
	Files       []FileSystemEntry `json:"files,omitempty"`
	Parent      string            `json:"parent,omitempty"`
}

// FileSystemEntry represents a single directory or file entry.
type FileSystemEntry struct {
	Type         string `json:"type,omitempty"`
	Name         string `json:"name,omitempty"`
	Path         string `json:"path,omitempty"`
	Size         int64  `json:"size,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
}

// LocalizationResource represents localization string data.
type LocalizationResource struct {
	ID      int               `json:"id"`
	Strings map[string]string `json:"strings,omitempty"`
}
