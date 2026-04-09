package googlebooks

// VolumesResponse represents the response from a volumes search.
type VolumesResponse struct {
	Kind       string   `json:"kind"`
	TotalItems int      `json:"totalItems"`
	Items      []Volume `json:"items"`
}

// Volume represents a Google Books volume.
type Volume struct {
	Kind       string      `json:"kind"`
	Id         string      `json:"id"`
	SelfLink   string      `json:"selfLink"`
	VolumeInfo *VolumeInfo `json:"volumeInfo"`
	SaleInfo   *SaleInfo   `json:"saleInfo,omitempty"`
	AccessInfo *AccessInfo `json:"accessInfo,omitempty"`
}

// VolumeInfo contains detailed information about a volume.
type VolumeInfo struct {
	Title               string               `json:"title"`
	Subtitle            string               `json:"subtitle"`
	Authors             []string             `json:"authors"`
	Publisher           string               `json:"publisher"`
	PublishedDate       string               `json:"publishedDate"`
	Description         string               `json:"description"`
	IndustryIdentifiers []IndustryIdentifier `json:"industryIdentifiers"`
	PageCount           int                  `json:"pageCount"`
	PrintType           string               `json:"printType"`
	Categories          []string             `json:"categories"`
	AverageRating       float64              `json:"averageRating"`
	RatingsCount        int                  `json:"ratingsCount"`
	MaturityRating      string               `json:"maturityRating"`
	ImageLinks          *ImageLinks          `json:"imageLinks"`
	Language            string               `json:"language"`
	PreviewLink         string               `json:"previewLink"`
	InfoLink            string               `json:"infoLink"`
}

// IndustryIdentifier represents an ISBN or other identifier.
type IndustryIdentifier struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

// ImageLinks contains thumbnail URLs.
type ImageLinks struct {
	SmallThumbnail string `json:"smallThumbnail"`
	Thumbnail      string `json:"thumbnail"`
}

// SaleInfo contains sale-related information.
type SaleInfo struct {
	Country     string `json:"country"`
	Saleability string `json:"saleability"`
	IsEbook     bool   `json:"isEbook"`
	ListPrice   *Price `json:"listPrice,omitempty"`
	RetailPrice *Price `json:"retailPrice,omitempty"`
}

// Price represents a monetary value.
type Price struct {
	Amount       float64 `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
}

// AccessInfo contains access-related information.
type AccessInfo struct {
	Country                string      `json:"country"`
	Viewability            string      `json:"viewability"`
	Embeddable             bool        `json:"embeddable"`
	PublicDomain           bool        `json:"publicDomain"`
	TextToSpeechPermission string      `json:"textToSpeechPermission"`
	Epub                   *FormatInfo `json:"epub"`
	Pdf                    *FormatInfo `json:"pdf"`
}

// FormatInfo represents availability for a specific format.
type FormatInfo struct {
	IsAvailable  bool   `json:"isAvailable"`
	AcsTokenLink string `json:"acsTokenLink,omitempty"`
}

// SearchParams holds parameters for an advanced volume search.
type SearchParams struct {
	Query        string
	StartIndex   int
	MaxResults   int
	PrintType    string
	OrderBy      string
	Filter       string
	LangRestrict string
	Projection   string
}

// BookshelvesResponse represents the response from a bookshelves request.
type BookshelvesResponse struct {
	Kind  string      `json:"kind"`
	Items []Bookshelf `json:"items"`
}

// Bookshelf represents a user's bookshelf.
type Bookshelf struct {
	Kind           string `json:"kind"`
	ID             int    `json:"id"`
	SelfLink       string `json:"selfLink"`
	Title          string `json:"title"`
	Description    string `json:"description,omitempty"`
	Access         string `json:"access"`
	Updated        string `json:"updated"`
	Created        string `json:"created"`
	VolumeCount    int    `json:"volumeCount"`
	VolumesLastUpdated string `json:"volumesLastUpdated,omitempty"`
}

// AnnotationsResponse represents the response from an annotations request.
type AnnotationsResponse struct {
	Kind       string       `json:"kind"`
	TotalItems int          `json:"totalItems"`
	Items      []Annotation `json:"items"`
}

// Annotation represents a highlight or note in a volume.
type Annotation struct {
	Kind             string `json:"kind"`
	ID               string `json:"id"`
	SelfLink         string `json:"selfLink"`
	VolumeID         string `json:"volumeId"`
	LayerID          string `json:"layerId"`
	SelectedText     string `json:"selectedText,omitempty"`
	HighlightStyle   string `json:"highlightStyle,omitempty"`
	Data             string `json:"data,omitempty"`
	Created          string `json:"created"`
	Updated          string `json:"updated"`
	PageIDs          []string `json:"pageIds,omitempty"`
	BeforeSelectedText string `json:"beforeSelectedText,omitempty"`
	AfterSelectedText  string `json:"afterSelectedText,omitempty"`
}

// ReadingPosition represents the reading position in a volume.
type ReadingPosition struct {
	Kind               string `json:"kind"`
	VolumeID           string `json:"volumeId"`
	Position           string `json:"position"`
	Updated            string `json:"updated"`
	PdfPosition        string `json:"pdfPosition,omitempty"`
	EpubCfiPosition    string `json:"epubCfiPosition,omitempty"`
	GbImagePosition    string `json:"gbImagePosition,omitempty"`
	GbTextPosition     string `json:"gbTextPosition,omitempty"`
}

// SeriesResponse represents the response from a series request.
type SeriesResponse struct {
	Kind   string       `json:"kind"`
	Series []SeriesInfo `json:"series"`
}

// SeriesInfo represents a book series.
type SeriesInfo struct {
	Kind                  string `json:"kind"`
	SeriesID              string `json:"seriesId"`
	Title                 string `json:"title"`
	BannerImageURL        string `json:"bannerImageUrl,omitempty"`
	ImageURL              string `json:"imageUrl,omitempty"`
	SeriesType            string `json:"seriesType,omitempty"`
	EligibleForSubscription bool  `json:"eligibleForSubscription,omitempty"`
}
