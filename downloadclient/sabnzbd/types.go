package sabnzbd

// Queue represents the SABnzbd download queue.
type Queue struct {
	Status     string `json:"status"`
	SpeedLimit string `json:"speedlimit"`
	Speed      string `json:"speed"`
	SizeLeft   string `json:"sizeleft"`
	Size       string `json:"size"`
	TimeLeft   string `json:"timeleft"`
	Slots      []Slot `json:"slots"`
	NoOfSlots  int    `json:"noofslots"`
	Paused     bool   `json:"paused"`
	PausedAll  bool   `json:"paused_all"`
	DiskSpace1 string `json:"diskspace1"`
	DiskSpace2 string `json:"diskspace2"`
}

// Slot represents a single download in the queue.
type Slot struct {
	ID         string `json:"nzo_id"`
	Filename   string `json:"filename"`
	Status     string `json:"status"`
	Size       string `json:"size"`
	SizeLeft   string `json:"sizeleft"`
	Percentage string `json:"percentage"`
	ETA        string `json:"eta"`
	TimeLeft   string `json:"timeleft"`
	Category   string `json:"cat"`
	Priority   string `json:"priority"`
	Index      int    `json:"index"`
	Script     string `json:"script"`
	MBLeft     string `json:"mbleft"`
	MB         string `json:"mb"`
	AvgAge     string `json:"avg_age"`
}

// HistoryResponse wraps the history listing.
type HistoryResponse struct {
	TotalSize string        `json:"total_size"`
	NoOfSlots int           `json:"noofslots"`
	Slots     []HistorySlot `json:"slots"`
}

// HistorySlot represents a completed download in history.
type HistorySlot struct {
	ID           string `json:"nzo_id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Size         string `json:"size"`
	Category     string `json:"category"`
	CompletedOn  int64  `json:"completed"`
	DownloadTime int64  `json:"download_time"`
	StoragePath  string `json:"storage"`
	Script       string `json:"script"`
	ScriptLine   string `json:"script_line"`
	FailMessage  string `json:"fail_message"`
}

// ServerStats holds SABnzbd server statistics.
type ServerStats struct {
	Total int64 `json:"total"`
	Day   int64 `json:"day"`
	Week  int64 `json:"week"`
	Month int64 `json:"month"`
}

// VersionInfo holds version information.
type VersionInfo struct {
	Version string `json:"version"`
}

// addResult is the response from adding an NZB.
type addResult struct {
	Status bool     `json:"status"`
	NZOIDs []string `json:"nzo_ids"`
}
