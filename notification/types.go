package notification

// Message represents a generic notification message.
type Message struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Priority int    `json:"priority"`
}
