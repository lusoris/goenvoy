package openlibrary

// SearchResponse represents the response from a search query.
type SearchResponse struct {
	NumFound int         `json:"numFound"`
	Start    int         `json:"start"`
	Docs     []SearchDoc `json:"docs"`
}

// SearchDoc represents a single document in search results.
type SearchDoc struct {
	Key              string   `json:"key"`
	Title            string   `json:"title"`
	AuthorName       []string `json:"author_name"`
	FirstPublishYear int      `json:"first_publish_year"`
	ISBN             []string `json:"isbn"`
	Publisher        []string `json:"publisher"`
	Language         []string `json:"language"`
	Subject          []string `json:"subject"`
	CoverI           int      `json:"cover_i"`
	EditionCount     int      `json:"edition_count"`
}

// Work represents an Open Library work.
type Work struct {
	Key         string       `json:"key"`
	Title       string       `json:"title"`
	Description any          `json:"description"`
	Subjects    []string     `json:"subjects"`
	Authors     []AuthorRole `json:"authors"`
	Covers      []int        `json:"covers"`
	Created     *ChangeEntry `json:"created"`
}

// AuthorRole represents an author entry on a work.
type AuthorRole struct {
	Author *AuthorRef `json:"author"`
	Type   *TypeRef   `json:"type"`
}

// AuthorRef is a reference to an author by key.
type AuthorRef struct {
	Key string `json:"key"`
}

// TypeRef is a reference to a type by key.
type TypeRef struct {
	Key string `json:"key"`
}

// ChangeEntry records a change timestamp.
type ChangeEntry struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Edition represents an Open Library edition.
type Edition struct {
	Key           string    `json:"key"`
	Title         string    `json:"title"`
	Publishers    []string  `json:"publishers"`
	PublishDate   string    `json:"publish_date"`
	ISBN10        []string  `json:"isbn_10"`
	ISBN13        []string  `json:"isbn_13"`
	NumberOfPages int       `json:"number_of_pages"`
	Covers        []int     `json:"covers"`
	Works         []WorkRef `json:"works"`
}

// WorkRef is a reference to a work by key.
type WorkRef struct {
	Key string `json:"key"`
}

// Author represents an Open Library author.
type Author struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Bio       any    `json:"bio"`
	BirthDate string `json:"birth_date"`
	DeathDate string `json:"death_date"`
	Photos    []int  `json:"photos"`
	Links     []Link `json:"links"`
}

// Link represents a link on an author.
type Link struct {
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Type  *TypeRef `json:"type"`
}

// Subject represents an Open Library subject.
type Subject struct {
	Name      string        `json:"name"`
	WorkCount int           `json:"work_count"`
	Works     []SubjectWork `json:"works"`
}

// SubjectWork represents a work within a subject listing.
type SubjectWork struct {
	Key          string          `json:"key"`
	Title        string          `json:"title"`
	Authors      []SubjectAuthor `json:"authors"`
	EditionCount int             `json:"edition_count"`
	CoverID      int             `json:"cover_id"`
}

// SubjectAuthor represents an author in a subject work.
type SubjectAuthor struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}
