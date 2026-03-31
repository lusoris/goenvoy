package tdarr

// Status represents the Tdarr server status.
type Status struct {
	Status  string `json:"status"`
	Os      string `json:"os"`
	Version string `json:"version"`
}

// Node represents a Tdarr processing node.
type Node struct {
	Id       string            `json:"_id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Workers  map[string]Worker `json:"workers"`
	NodeType string            `json:"nodeType"`
}

// Worker represents a transcoding worker on a node.
type Worker struct {
	Id         string  `json:"id"`
	File       string  `json:"file"`
	Percentage float64 `json:"percentage"`
	ETA        string  `json:"ETA"`
	Status     string  `json:"status"`
	WorkerType string  `json:"workerType"`
}

// DBFile represents a media file entry in the Tdarr database.
type DBFile struct {
	Id                     string  `json:"_id"`
	File                   string  `json:"file"`
	Codec                  string  `json:"codec"`
	Container              string  `json:"container"`
	Resolution             string  `json:"resolution"`
	FileSize               int64   `json:"fileSize"`
	TranscodeDecisionMaker string  `json:"transcodeDecisionMaker"`
	HealthCheck            string  `json:"healthCheck"`
	Bitrate                int64   `json:"bitrate"`
	Duration               float64 `json:"duration"`
	LibraryId              string  `json:"libraryId"`
}

// SearchDBRequest is the request body for the search-db endpoint.
type SearchDBRequest struct {
	Data *SearchDBData `json:"data"`
}

// SearchDBData holds the search parameters.
type SearchDBData struct {
	Collection string         `json:"collection"`
	Limit      int            `json:"limit"`
	Skip       int            `json:"skip"`
	Filters    []SearchFilter `json:"filters"`
}

// SearchFilter defines a single filter criterion for database searches.
type SearchFilter struct {
	Field     string `json:"field"`
	Value     string `json:"value"`
	Condition string `json:"condition"`
}

// CrudDBRequest is the request body for the cruddb endpoint.
type CrudDBRequest struct {
	Data *CrudDBData `json:"data"`
}

// CrudDBData holds the CRUD operation parameters.
type CrudDBData struct {
	Collection string           `json:"collection"`
	Mode       string           `json:"mode"`
	DocID      string           `json:"docID,omitempty"`
	Docs       []map[string]any `json:"docs,omitempty"`
}

// ResStats holds resolution statistics.
type ResStats struct {
	Pie map[string]int `json:"pie"`
}

// DBStatuses holds database table counts.
type DBStatuses struct {
	Table1Count int `json:"table1Count"`
	Table2Count int `json:"table2Count"`
	Table3Count int `json:"table3Count"`
	Table4Count int `json:"table4Count"`
	Table5Count int `json:"table5Count"`
	Table6Count int `json:"table6Count"`
}

// ScanFilesRequest is the request body for the scan-files endpoint.
type ScanFilesRequest struct {
	Data *ScanFilesData `json:"data"`
}

// ScanFilesData holds the scan parameters.
type ScanFilesData struct {
	LibraryId  string `json:"libraryId"`
	FolderPath string `json:"folderPath,omitempty"`
}
