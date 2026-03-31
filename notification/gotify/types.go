package gotify

// Application represents a Gotify application.
type Application struct {
	Id          int    `json:"id"`
	Token       string `json:"token"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Internal    bool   `json:"internal"`
	Image       string `json:"image"`
}

// Message represents a Gotify push notification message.
type Message struct {
	Id       int            `json:"id"`
	AppId    int            `json:"appid"`
	Title    string         `json:"title"`
	Message  string         `json:"message"`
	Priority int            `json:"priority"`
	Date     string         `json:"date"`
	Extras   map[string]any `json:"extras,omitempty"`
}

// PagedMessages is a paginated list of messages.
type PagedMessages struct {
	Messages []Message `json:"messages"`
	Paging   *Paging   `json:"paging"`
}

// Paging holds pagination metadata.
type Paging struct {
	Size  int    `json:"size"`
	Limit int    `json:"limit"`
	Since int    `json:"since"`
	Next  string `json:"next"`
}

// ClientInfo represents a Gotify client.
type ClientInfo struct {
	Id       int    `json:"id"`
	Token    string `json:"token"`
	Name     string `json:"name"`
	LastUsed string `json:"lastUsed"`
}

// User represents a Gotify user.
type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

// Health represents the Gotify server health status.
type Health struct {
	Health   string `json:"health"`
	Database string `json:"database"`
}

// VersionInfo represents the Gotify server version.
type VersionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
}
