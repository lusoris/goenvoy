package screenscraper

// GameInfoResponse is the top-level response from jeuInfos and jeuRecherche.
type GameInfoResponse struct {
	Header   Header   `json:"header"`
	Response Response `json:"response"`
}

// Header contains API response metadata.
type Header struct {
	APIVersion       string `json:"APIversion"`
	CommandRequested string `json:"commandRequested"`
	Success          string `json:"success"`
	Error            string `json:"error"`
}

// Response wraps the game data in the API response.
type Response struct {
	Game  GameInfo   `json:"jeu"`
	Games []GameInfo `json:"jeux"`
}

// GameInfo represents detailed game information from Screenscraper.
type GameInfo struct {
	ID              string              `json:"id"`
	RomID           string              `json:"romid"`
	Names           []RegionText        `json:"noms"`
	SystemID        string              `json:"systemeid"`
	System          SystemRef           `json:"systeme"`
	Publisher       Publisher           `json:"editeur"`
	Developer       Developer           `json:"developpeur"`
	Players         string              `json:"joueurs"`
	Note            Rating              `json:"note"`
	Genres          []GenreRef          `json:"genres"`
	Synopsis        []LangText          `json:"synopsis"`
	Dates           []RegionText        `json:"dates"`
	Medias          []Media             `json:"medias"`
	ROMs            []ROM               `json:"roms"`
	Familles        []FamilleRef        `json:"familles"`
	Classifications []ClassificationRef `json:"classifications"`
}

// RegionText holds text keyed by region.
type RegionText struct {
	Region string `json:"region"`
	Text   string `json:"text"`
}

// LangText holds text keyed by language.
type LangText struct {
	Langue string `json:"langue"`
	Text   string `json:"text"`
}

// SystemRef is a reference to a system.
type SystemRef struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// Publisher represents a game publisher.
type Publisher struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// Developer represents a game developer.
type Developer struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// Rating represents a game rating/score.
type Rating struct {
	Text string `json:"text"`
}

// GenreRef is a genre reference with localized names.
type GenreRef struct {
	ID    string     `json:"id"`
	Names []LangText `json:"noms"`
}

// FamilleRef is a game family/series reference.
type FamilleRef struct {
	ID    string     `json:"id"`
	Names []LangText `json:"noms"`
}

// ClassificationRef is a classification/rating reference.
type ClassificationRef struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Media represents a media asset (screenshot, box art, etc.).
type Media struct {
	Type   string `json:"type"`
	Parent string `json:"parent"`
	URL    string `json:"url"`
	Region string `json:"region"`
	CRC    string `json:"crc"`
	MD5    string `json:"md5"`
	SHA1   string `json:"sha1"`
	Size   string `json:"size"`
	Format string `json:"format"`
}

// ROM represents ROM file information.
type ROM struct {
	ID      string `json:"id"`
	RomSize string `json:"romsize"`
	RomName string `json:"romfilename"`
	RomCRC  string `json:"romcrc"`
	RomMD5  string `json:"rommd5"`
	RomSHA1 string `json:"romsha1"`
	Region  string `json:"romregion"`
}

// SystemsResponse is the top-level response from systemesListe.
type SystemsResponse struct {
	Header   Header        `json:"header"`
	Response SystemsResult `json:"response"`
}

// SystemsResult wraps the systems list.
type SystemsResult struct {
	Systems []System `json:"systemes"`
}

// System represents a gaming system/platform.
type System struct {
	ID         string       `json:"id"`
	Names      []RegionText `json:"noms"`
	Extensions string       `json:"extensions"`
	Company    string       `json:"compagnie"`
	Type       string       `json:"type"`
	Medias     []Media      `json:"medias"`
}

// GenresResponse is the top-level response from genresListe.
type GenresResponse struct {
	Header   Header       `json:"header"`
	Response GenresResult `json:"response"`
}

// GenresResult wraps the genres list.
type GenresResult struct {
	Genres []Genre `json:"genres"`
}

// Genre represents a genre with localized names.
type Genre struct {
	ID    string     `json:"id"`
	Names []LangText `json:"noms"`
}

// UserInfoResponse is the top-level response from ssuserInfos.
type UserInfoResponse struct {
	Header   Header   `json:"header"`
	Response UserInfo `json:"response"`
}

// UserInfo represents user account information.
type UserInfo struct {
	ID             string `json:"ssid"`
	NumRequests    string `json:"requeststoday"`
	MaxRequests    string `json:"maxrequestsperday"`
	VisitCount     string `json:"visites"`
	FavoriteRegion string `json:"favregion"`
}

// InfraInfoResponse is the response from ssinfraInfos.
type InfraInfoResponse struct {
	Header   Header    `json:"header"`
	Response InfraInfo `json:"response"`
}

// InfraInfo contains API infrastructure information.
type InfraInfo struct {
	MaxThreads string `json:"maxthreads"`
	CPU        string `json:"cpu"`
}
