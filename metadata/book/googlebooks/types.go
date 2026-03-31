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
