package deluge

import "encoding/json"

// Torrent represents a torrent managed by Deluge.
type Torrent struct {
	Hash                string  `json:"hash"`
	Name                string  `json:"name"`
	State               string  `json:"state"`
	Progress            float64 `json:"progress"`
	SavePath            string  `json:"save_path"`
	TotalSize           int64   `json:"total_size"`
	TotalDone           int64   `json:"total_done"`
	TotalUploaded       int64   `json:"total_uploaded"`
	DownloadPayloadRate int64   `json:"download_payload_rate"`
	UploadPayloadRate   int64   `json:"upload_payload_rate"`
	NumPeers            int     `json:"num_peers"`
	NumSeeds            int     `json:"num_seeds"`
	ETA                 int64   `json:"eta"`
	Ratio               float64 `json:"ratio"`
	Label               string  `json:"label"`
	TimeAdded           float64 `json:"time_added"`
	ActiveTime          int64   `json:"active_time"`
	SeedingTime         int64   `json:"seeding_time"`
	IsFinished          bool    `json:"is_finished"`
	IsSeed              bool    `json:"is_seed"`
	TrackerHost         string  `json:"tracker_host"`
	TrackerStatus       string  `json:"tracker_status"`
	Comment             string  `json:"comment"`
}

// SessionStatus holds session statistics.
type SessionStatus struct {
	DownloadRate           int64 `json:"payload_download_rate"`
	UploadRate             int64 `json:"payload_upload_rate"`
	DHTNodes               int   `json:"dht_nodes"`
	HasIncomingConnections bool  `json:"has_incoming_connections"`
	TotalDownload          int64 `json:"total_payload_download"`
	TotalUpload            int64 `json:"total_payload_upload"`
}

// defaultTorrentFields are the keys requested from Deluge torrent status.
var defaultTorrentFields = []string{
	"hash", "name", "state", "progress", "save_path",
	"total_size", "total_done", "total_uploaded",
	"download_payload_rate", "upload_payload_rate",
	"num_peers", "num_seeds", "eta", "ratio", "label",
	"time_added", "active_time", "seeding_time",
	"is_finished", "is_seed", "tracker_host", "tracker_status", "comment",
}

// rpcRequest is the JSON-RPC request envelope for Deluge.
type rpcRequest struct {
	ID     int    `json:"id"`
	Method string `json:"method"`
	Params []any  `json:"params"`
}

// rpcResponse is the JSON-RPC response envelope.
type rpcResponse struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error"`
}

// rpcError represents a JSON-RPC error from Deluge.
type rpcError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
