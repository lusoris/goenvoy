package omdb

// Title represents a full movie or series result from the OMDb API.
type Title struct {
	Title        string   `json:"Title"`
	Year         string   `json:"Year"`
	Rated        string   `json:"Rated"`
	Released     string   `json:"Released"`
	Runtime      string   `json:"Runtime"`
	Genre        string   `json:"Genre"`
	Director     string   `json:"Director"`
	Writer       string   `json:"Writer"`
	Actors       string   `json:"Actors"`
	Plot         string   `json:"Plot"`
	Language     string   `json:"Language"`
	Country      string   `json:"Country"`
	Awards       string   `json:"Awards"`
	Poster       string   `json:"Poster"`
	Ratings      []Rating `json:"Ratings"`
	Metascore    string   `json:"Metascore"`
	IMDbRating   string   `json:"imdbRating"`
	IMDbVotes    string   `json:"imdbVotes"`
	IMDbID       string   `json:"imdbID"`
	Type         string   `json:"Type"`
	DVD          string   `json:"DVD"`
	BoxOffice    string   `json:"BoxOffice"`
	Production   string   `json:"Production"`
	Website      string   `json:"Website"`
	TotalSeasons string   `json:"totalSeasons,omitempty"`
	Response     string   `json:"Response"`
}

// Rating is a single rating from a source (e.g. IMDb, Rotten Tomatoes).
type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

// SearchResult is a single item in a search response.
type SearchResult struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	IMDbID string `json:"imdbID"`
	Type   string `json:"Type"`
	Poster string `json:"Poster"`
}

// SearchResponse is the top-level response from a search query.
type SearchResponse struct {
	Search       []SearchResult `json:"Search"`
	TotalResults string         `json:"totalResults"`
	Response     string         `json:"Response"`
}

// SeasonResponse is the response for a season listing query.
type SeasonResponse struct {
	Title        string    `json:"Title"`
	Season       string    `json:"Season"`
	TotalSeasons string    `json:"totalSeasons"`
	Episodes     []Episode `json:"Episodes"`
	Response     string    `json:"Response"`
}

// Episode is a single episode entry in a season listing.
type Episode struct {
	Title      string `json:"Title"`
	Released   string `json:"Released"`
	Episode    string `json:"Episode"`
	IMDbRating string `json:"imdbRating"`
	IMDbID     string `json:"imdbID"`
}

// MediaType restricts search results to a specific type.
type MediaType string

const (
	// MediaTypeMovie filters results to movies.
	MediaTypeMovie MediaType = "movie"
	// MediaTypeSeries filters results to TV series.
	MediaTypeSeries MediaType = "series"
	// MediaTypeEpisode filters results to individual episodes.
	MediaTypeEpisode MediaType = "episode"
)

// PlotLength controls the verbosity of the plot synopsis.
type PlotLength string

const (
	// PlotShort returns a short plot synopsis (default).
	PlotShort PlotLength = "short"
	// PlotFull returns the full plot synopsis.
	PlotFull PlotLength = "full"
)
