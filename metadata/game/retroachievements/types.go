package retroachievements

// Game represents basic game information from RetroAchievements.
type Game struct {
	Title                 string `json:"Title"`
	GameTitle             string `json:"GameTitle"`
	ConsoleID             int    `json:"ConsoleID"`
	ConsoleName           string `json:"ConsoleName"`
	ForumTopicID          int    `json:"ForumTopicID"`
	Flags                 int    `json:"Flags"`
	GameIcon              string `json:"GameIcon"`
	ImageIcon             string `json:"ImageIcon"`
	ImageTitle            string `json:"ImageTitle"`
	ImageIngame           string `json:"ImageIngame"`
	ImageBoxArt           string `json:"ImageBoxArt"`
	Publisher             string `json:"Publisher"`
	Developer             string `json:"Developer"`
	Genre                 string `json:"Genre"`
	Released              string `json:"Released"`
	ReleasedAtGranularity string `json:"ReleasedAtGranularity"`
}

// GameExtended represents extended game information including achievements.
type GameExtended struct {
	ID                         int                    `json:"ID"`
	Title                      string                 `json:"Title"`
	ConsoleID                  int                    `json:"ConsoleID"`
	ConsoleName                string                 `json:"ConsoleName"`
	ForumTopicID               int                    `json:"ForumTopicID"`
	Flags                      int                    `json:"Flags"`
	GameIcon                   string                 `json:"GameIcon"`
	ImageIcon                  string                 `json:"ImageIcon"`
	ImageTitle                 string                 `json:"ImageTitle"`
	ImageIngame                string                 `json:"ImageIngame"`
	ImageBoxArt                string                 `json:"ImageBoxArt"`
	Publisher                  string                 `json:"Publisher"`
	Developer                  string                 `json:"Developer"`
	Genre                      string                 `json:"Genre"`
	Released                   string                 `json:"Released"`
	ReleasedAtGranularity      string                 `json:"ReleasedAtGranularity"`
	IsFinal                    bool                   `json:"IsFinal"`
	RichPresencePatch          string                 `json:"RichPresencePatch"`
	GuideURL                   string                 `json:"GuideURL"`
	Updated                    string                 `json:"Updated"`
	ParentGameID               int                    `json:"ParentGameID"`
	NumDistinctPlayers         int                    `json:"NumDistinctPlayers"`
	NumDistinctPlayersCasual   int                    `json:"NumDistinctPlayersCasual"`
	NumDistinctPlayersHardcore int                    `json:"NumDistinctPlayersHardcore"`
	NumAchievements            int                    `json:"NumAchievements"`
	Achievements               map[string]Achievement `json:"Achievements"`
	Claims                     []Claim                `json:"Claims"`
}

// Achievement represents a single achievement for a game.
type Achievement struct {
	ID                 int    `json:"ID"`
	NumAwarded         int    `json:"NumAwarded"`
	NumAwardedHardcore int    `json:"NumAwardedHardcore"`
	Title              string `json:"Title"`
	Description        string `json:"Description"`
	Points             int    `json:"Points"`
	TrueRatio          int    `json:"TrueRatio"`
	Author             string `json:"Author"`
	DateModified       string `json:"DateModified"`
	DateCreated        string `json:"DateCreated"`
	BadgeName          string `json:"BadgeName"`
	DisplayOrder       int    `json:"DisplayOrder"`
	MemAddr            string `json:"MemAddr"`
	Type               string `json:"type"`
}

// Claim represents a development claim on a game.
type Claim struct {
	User       string `json:"User"`
	SetType    int    `json:"SetType"`
	ClaimType  int    `json:"ClaimType"`
	Created    string `json:"Created"`
	Expiration string `json:"Expiration"`
}

// GameHash represents a ROM hash associated with a game.
type GameHash struct {
	MD5      string   `json:"MD5"`
	Name     string   `json:"Name"`
	Labels   []string `json:"Labels"`
	PatchURL string   `json:"PatchUrl"`
}

// HashResult wraps the response from GetGameHashes.
type HashResult struct {
	Results []GameHash `json:"Results"`
}

// Console represents a gaming platform/console.
type Console struct {
	ID           int    `json:"ID"`
	Name         string `json:"Name"`
	IconURL      string `json:"IconURL"`
	Active       bool   `json:"Active"`
	IsGameSystem bool   `json:"IsGameSystem"`
}
