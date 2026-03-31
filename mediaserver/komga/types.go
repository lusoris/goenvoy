package komga

// Library represents a Komga library.
type Library struct {
	ID                        string `json:"id"`
	Name                      string `json:"name"`
	Root                      string `json:"root"`
	ImportComicInfoBook       bool   `json:"importComicInfoBook"`
	ImportComicInfoSeries     bool   `json:"importComicInfoSeries"`
	ImportComicInfoCollection bool   `json:"importComicInfoCollection"`
	ImportComicInfoReadList   bool   `json:"importComicInfoReadList"`
	ImportBarcodeISBN         bool   `json:"importBarcodeIsbn"`
	ScanForceModifiedTime     bool   `json:"scanForceModifiedTime"`
	ScanDeep                  bool   `json:"scanDeep"`
	RepairExtensions          bool   `json:"repairExtensions"`
	ConvertToCBZ              bool   `json:"convertToCbz"`
	EmptyTrashAfterScan       bool   `json:"emptyTrashAfterScan"`
	SeriesCover               string `json:"seriesCover"`
	HashFiles                 bool   `json:"hashFiles"`
	HashPages                 bool   `json:"hashPages"`
	AnalyzeDimensions         bool   `json:"analyzeDimensions"`
	Unavailable               bool   `json:"unavailable,omitempty"`
}

// Series represents a comic/manga series.
type Series struct {
	ID                   string         `json:"id"`
	LibraryID            string         `json:"libraryId"`
	Name                 string         `json:"name"`
	URL                  string         `json:"url,omitempty"`
	Created              string         `json:"created,omitempty"`
	LastModified         string         `json:"lastModified,omitempty"`
	FileLastModified     string         `json:"fileLastModified,omitempty"`
	BooksCount           int            `json:"booksCount"`
	BooksReadCount       int            `json:"booksReadCount"`
	BooksUnreadCount     int            `json:"booksUnreadCount"`
	BooksInProgressCount int            `json:"booksInProgressCount"`
	Metadata             SeriesMetadata `json:"metadata,omitempty"`
	Deleted              bool           `json:"deleted,omitempty"`
	OneShot              bool           `json:"oneshot,omitempty"`
}

// SeriesMetadata represents metadata for a series.
type SeriesMetadata struct {
	Status           string   `json:"status,omitempty"`
	StatusLock       bool     `json:"statusLock,omitempty"`
	Title            string   `json:"title,omitempty"`
	TitleLock        bool     `json:"titleLock,omitempty"`
	TitleSort        string   `json:"titleSort,omitempty"`
	TitleSortLock    bool     `json:"titleSortLock,omitempty"`
	Summary          string   `json:"summary,omitempty"`
	SummaryLock      bool     `json:"summaryLock,omitempty"`
	Publisher        string   `json:"publisher,omitempty"`
	ReadingDirection string   `json:"readingDirection,omitempty"`
	AgeRating        int      `json:"ageRating,omitempty"`
	Language         string   `json:"language,omitempty"`
	Genres           []string `json:"genres,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	TotalBookCount   int      `json:"totalBookCount,omitempty"`
}

// Book represents a single comic/manga book (issue/chapter).
type Book struct {
	ID               string        `json:"id"`
	SeriesID         string        `json:"seriesId"`
	LibraryID        string        `json:"libraryId"`
	Name             string        `json:"name"`
	URL              string        `json:"url,omitempty"`
	Number           int           `json:"number"`
	Created          string        `json:"created,omitempty"`
	LastModified     string        `json:"lastModified,omitempty"`
	FileLastModified string        `json:"fileLastModified,omitempty"`
	SizeBytes        int64         `json:"sizeBytes"`
	Size             string        `json:"size,omitempty"`
	Media            BookMedia     `json:"media,omitempty"`
	Metadata         BookMetadata  `json:"metadata,omitempty"`
	ReadProgress     *ReadProgress `json:"readProgress,omitempty"`
	Deleted          bool          `json:"deleted,omitempty"`
	FileHash         string        `json:"fileHash,omitempty"`
	OneShot          bool          `json:"oneshot,omitempty"`
}

// BookMedia represents media information for a book.
type BookMedia struct {
	Status     string `json:"status,omitempty"`
	MediaType  string `json:"mediaType,omitempty"`
	PagesCount int    `json:"pagesCount"`
	Comment    string `json:"comment,omitempty"`
}

// BookMetadata represents metadata for a book.
type BookMetadata struct {
	Title       string       `json:"title,omitempty"`
	TitleLock   bool         `json:"titleLock,omitempty"`
	Summary     string       `json:"summary,omitempty"`
	Number      string       `json:"number,omitempty"`
	NumberSort  float64      `json:"numberSort,omitempty"`
	ReleaseDate string       `json:"releaseDate,omitempty"`
	Authors     []PersonRole `json:"authors,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
	ISBN        string       `json:"isbn,omitempty"`
}

// PersonRole represents a person and their role.
type PersonRole struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// ReadProgress represents reading progress for a book.
type ReadProgress struct {
	Page         int    `json:"page"`
	Completed    bool   `json:"completed"`
	ReadDate     string `json:"readDate,omitempty"`
	Created      string `json:"created,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
}

// Collection represents a collection of series.
type Collection struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Ordered          bool     `json:"ordered"`
	SeriesIDs        []string `json:"seriesIds,omitempty"`
	CreatedDate      string   `json:"createdDate,omitempty"`
	LastModifiedDate string   `json:"lastModifiedDate,omitempty"`
	Filtered         bool     `json:"filtered,omitempty"`
}

// ReadList represents a read list of books.
type ReadList struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Summary          string   `json:"summary,omitempty"`
	Ordered          bool     `json:"ordered"`
	BookIDs          []string `json:"bookIds,omitempty"`
	CreatedDate      string   `json:"createdDate,omitempty"`
	LastModifiedDate string   `json:"lastModifiedDate,omitempty"`
	Filtered         bool     `json:"filtered,omitempty"`
}

// Page represents a paginated response.
type Page[T any] struct {
	Content          []T  `json:"content"`
	Pageable         any  `json:"pageable,omitempty"`
	TotalPages       int  `json:"totalPages"`
	TotalElements    int  `json:"totalElements"`
	Last             bool `json:"last"`
	Size             int  `json:"size"`
	Number           int  `json:"number"`
	NumberOfElements int  `json:"numberOfElements"`
	First            bool `json:"first"`
	Empty            bool `json:"empty"`
}

// User represents a Komga user.
type User struct {
	ID                 string   `json:"id"`
	Email              string   `json:"email"`
	Roles              []string `json:"roles,omitempty"`
	SharedAllLibraries bool     `json:"sharedAllLibraries"`
	SharedLibrariesIDs []string `json:"sharedLibrariesIds,omitempty"`
	LabelsAllow        []string `json:"labelsAllow,omitempty"`
	LabelsDeny         []string `json:"labelsDeny,omitempty"`
}
