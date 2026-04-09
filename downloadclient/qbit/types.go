package qbit

// Torrent represents a torrent in qBittorrent.
type Torrent struct {
	AddedOn           int64   `json:"added_on"`
	AmountLeft        int64   `json:"amount_left"`
	AutoTMM           bool    `json:"auto_tmm"`
	Availability      float64 `json:"availability"`
	Category          string  `json:"category"`
	Completed         int64   `json:"completed"`
	CompletionOn      int64   `json:"completion_on"`
	ContentPath       string  `json:"content_path"`
	DlLimit           int64   `json:"dl_limit"`
	DlSpeed           int64   `json:"dlspeed"`
	Downloaded        int64   `json:"downloaded"`
	DownloadedSession int64   `json:"downloaded_session"`
	ETA               int64   `json:"eta"`
	FLPiecePrio       bool    `json:"f_l_piece_prio"`
	ForceStart        bool    `json:"force_start"`
	Hash              string  `json:"hash"`
	IsPrivate         bool    `json:"isPrivate"`
	LastActivity      int64   `json:"last_activity"`
	MagnetURI         string  `json:"magnet_uri"`
	MaxRatio          float64 `json:"max_ratio"`
	MaxSeedingTime    int64   `json:"max_seeding_time"`
	Name              string  `json:"name"`
	NumComplete       int     `json:"num_complete"`
	NumIncomplete     int     `json:"num_incomplete"`
	NumLeechs         int     `json:"num_leechs"`
	NumSeeds          int     `json:"num_seeds"`
	Priority          int     `json:"priority"`
	Progress          float64 `json:"progress"`
	Ratio             float64 `json:"ratio"`
	SavePath          string  `json:"save_path"`
	SeedingTime       int64   `json:"seeding_time"`
	Size              int64   `json:"size"`
	State             string  `json:"state"`
	SuperSeeding      bool    `json:"super_seeding"`
	Tags              string  `json:"tags"`
	TimeActive        int64   `json:"time_active"`
	TotalSize         int64   `json:"total_size"`
	Tracker           string  `json:"tracker"`
	UpLimit           int64   `json:"up_limit"`
	Uploaded          int64   `json:"uploaded"`
	UploadedSession   int64   `json:"uploaded_session"`
	UpSpeed           int64   `json:"upspeed"`
}

// TorrentProperties holds detailed properties for a single torrent.
type TorrentProperties struct {
	SavePath               string  `json:"save_path"`
	CreationDate           int64   `json:"creation_date"`
	PieceSize              int64   `json:"piece_size"`
	Comment                string  `json:"comment"`
	TotalWasted            int64   `json:"total_wasted"`
	TotalUploaded          int64   `json:"total_uploaded"`
	TotalUploadedSession   int64   `json:"total_uploaded_session"`
	TotalDownloaded        int64   `json:"total_downloaded"`
	TotalDownloadedSession int64   `json:"total_downloaded_session"`
	UpLimit                int64   `json:"up_limit"`
	DlLimit                int64   `json:"dl_limit"`
	TimeElapsed            int64   `json:"time_elapsed"`
	SeedingTime            int64   `json:"seeding_time"`
	NbConnections          int     `json:"nb_connections"`
	NbConnectionsLimit     int     `json:"nb_connections_limit"`
	ShareRatio             float64 `json:"share_ratio"`
	AdditionDate           int64   `json:"addition_date"`
	CompletionDate         int64   `json:"completion_date"`
	CreatedBy              string  `json:"created_by"`
	DlSpeedAvg             int64   `json:"dl_speed_avg"`
	DlSpeed                int64   `json:"dl_speed"`
	ETA                    int64   `json:"eta"`
	LastSeen               int64   `json:"last_seen"`
	Peers                  int     `json:"peers"`
	PeersTotal             int     `json:"peers_total"`
	PiecesHave             int     `json:"pieces_have"`
	PiecesNum              int     `json:"pieces_num"`
	Reannounce             int64   `json:"reannounce"`
	Seeds                  int     `json:"seeds"`
	SeedsTotal             int     `json:"seeds_total"`
	TotalSize              int64   `json:"total_size"`
	UpSpeedAvg             int64   `json:"up_speed_avg"`
	UpSpeed                int64   `json:"up_speed"`
	IsPrivate              bool    `json:"isPrivate"`
}

// TorrentFile represents a file within a torrent.
type TorrentFile struct {
	Index        int     `json:"index"`
	Name         string  `json:"name"`
	Size         int64   `json:"size"`
	Progress     float64 `json:"progress"`
	Priority     int     `json:"priority"`
	IsSeed       bool    `json:"is_seed"`
	PieceRange   []int   `json:"piece_range"`
	Availability float64 `json:"availability"`
}

// Tracker represents a tracker associated with a torrent.
type Tracker struct {
	URL           string `json:"url"`
	Status        int    `json:"status"`
	Tier          int    `json:"tier"`
	NumPeers      int    `json:"num_peers"`
	NumSeeds      int    `json:"num_seeds"`
	NumLeeches    int    `json:"num_leeches"`
	NumDownloaded int    `json:"num_downloaded"`
	Msg           string `json:"msg"`
}

// WebSeed represents a web seed URL for a torrent.
type WebSeed struct {
	URL string `json:"url"`
}

// Category represents a torrent category.
type Category struct {
	Name     string `json:"name"`
	SavePath string `json:"savePath"`
}

// TransferInfo holds global transfer statistics.
type TransferInfo struct {
	DlInfoSpeed       int64  `json:"dl_info_speed"`
	DlInfoData        int64  `json:"dl_info_data"`
	UpInfoSpeed       int64  `json:"up_info_speed"`
	UpInfoData        int64  `json:"up_info_data"`
	DlRateLimit       int64  `json:"dl_rate_limit"`
	UpRateLimit       int64  `json:"up_rate_limit"`
	DHTNodes          int    `json:"dht_nodes"`
	ConnectionStatus  string `json:"connection_status"`
	Queueing          bool   `json:"queueing"`
	UseAltSpeedLimits bool   `json:"use_alt_speed_limits"`
	RefreshInterval   int    `json:"refresh_interval"`
	FreeSpaceOnDisk   int64  `json:"free_space_on_disk"`
}

// BuildInfo holds qBittorrent build information.
type BuildInfo struct {
	Qt         string `json:"qt"`
	Libtorrent string `json:"libtorrent"`
	Boost      string `json:"boost"`
	OpenSSL    string `json:"openssl"`
	Bitness    int    `json:"bitness"`
}

// Preferences holds qBittorrent application preferences.
type Preferences struct {
	SavePath           string  `json:"save_path"`
	TempPathEnabled    bool    `json:"temp_path_enabled"`
	TempPath           string  `json:"temp_path"`
	DlLimit            int64   `json:"dl_limit"`
	UpLimit            int64   `json:"up_limit"`
	MaxActiveDownloads int     `json:"max_active_downloads"`
	MaxActiveTorrents  int     `json:"max_active_torrents"`
	MaxActiveUploads   int     `json:"max_active_uploads"`
	ListenPort         int     `json:"listen_port"`
	DHT                bool    `json:"dht"`
	PEX                bool    `json:"pex"`
	LSD                bool    `json:"lsd"`
	QueueingEnabled    bool    `json:"queueing_enabled"`
	MaxRatioEnabled    bool    `json:"max_ratio_enabled"`
	MaxRatio           float64 `json:"max_ratio"`
	MaxRatioAct        int     `json:"max_ratio_act"`
	AltDlLimit         int64   `json:"alt_dl_limit"`
	AltUpLimit         int64   `json:"alt_up_limit"`
	SchedulerEnabled   bool    `json:"scheduler_enabled"`
	AutoTMMEnabled     bool    `json:"auto_tmm_enabled"`
	WebUIPort          int     `json:"web_ui_port"`
	WebUIAddress       string  `json:"web_ui_address"`
}

// SyncMainData holds the response from the sync/maindata endpoint.
type SyncMainData struct {
	RID               int                  `json:"rid"`
	FullUpdate        bool                 `json:"full_update"`
	Torrents          map[string]*Torrent  `json:"torrents"`
	TorrentsRemoved   []string             `json:"torrents_removed"`
	Categories        map[string]*Category `json:"categories"`
	CategoriesRemoved []string             `json:"categories_removed"`
	Tags              []string             `json:"tags"`
	TagsRemoved       []string             `json:"tags_removed"`
	ServerState       *TransferInfo        `json:"server_state"`
}

// LogEntry represents a log message from qBittorrent.
type LogEntry struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Type      int    `json:"type"`
}

// PeerLogEntry represents a peer log entry.
type PeerLogEntry struct {
	ID        int    `json:"id"`
	IP        string `json:"ip"`
	Timestamp int64  `json:"timestamp"`
	Blocked   bool   `json:"blocked"`
	Reason    string `json:"reason"`
}

// AddTorrentOptions holds optional parameters for adding a torrent.
type AddTorrentOptions struct {
	SavePath           string  `json:"savepath,omitempty"`
	Category           string  `json:"category,omitempty"`
	Tags               string  `json:"tags,omitempty"`
	SkipChecking       bool    `json:"skip_checking,omitempty"`
	Paused             bool    `json:"paused,omitempty"`
	RootFolder         bool    `json:"root_folder,omitempty"`
	Rename             string  `json:"rename,omitempty"`
	UpLimit            int64   `json:"upLimit,omitempty"`
	DlLimit            int64   `json:"dlLimit,omitempty"`
	RatioLimit         float64 `json:"ratioLimit,omitempty"`
	SeedingTimeLimit   int64   `json:"seedingTimeLimit,omitempty"`
	AutoTMM            bool    `json:"autoTMM,omitempty"`
	SequentialDownload bool    `json:"sequentialDownload,omitempty"`
	FirstLastPiecePrio bool    `json:"firstLastPiecePrio,omitempty"`
}

// ListOptions holds optional parameters for listing torrents.
type ListOptions struct {
	Filter   string `json:"filter,omitempty"`
	Category string `json:"category,omitempty"`
	Tag      string `json:"tag,omitempty"`
	Sort     string `json:"sort,omitempty"`
	Reverse  bool   `json:"reverse,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	Hashes   string `json:"hashes,omitempty"`
}

// SyncTorrentPeers holds the response from the sync/torrentPeers endpoint.
type SyncTorrentPeers struct {
	RID      int                  `json:"rid"`
	FullData bool                 `json:"full_update"`
	Peers    map[string]*PeerInfo `json:"peers"`
	Removed  []string             `json:"peers_removed"`
}

// PeerInfo holds information about a connected peer.
type PeerInfo struct {
	IP          string  `json:"ip"`
	Port        int     `json:"port"`
	Client      string  `json:"client"`
	Progress    float64 `json:"progress"`
	DlSpeed     int64   `json:"dl_speed"`
	UpSpeed     int64   `json:"up_speed"`
	Downloaded  int64   `json:"downloaded"`
	Uploaded    int64   `json:"uploaded"`
	Connection  string  `json:"connection"`
	Flags       string  `json:"flags"`
	FlagsDesc   string  `json:"flags_desc"`
	Relevance   float64 `json:"relevance"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
}
