package qbit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTimeout   = 30 * time.Second
	defaultUserAgent = "goenvoy/0.0.1"
)

// Option configures a [Client].
type Option func(*Client)

// WithHTTPClient sets a custom [http.Client].
// The client must have a cookie jar configured for session management.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.httpClient = c }
}

// WithTimeout overrides the default HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(cl *Client) { cl.httpClient.Timeout = d }
}

// WithUserAgent sets the User-Agent header sent with every request.
func WithUserAgent(ua string) Option {
	return func(cl *Client) { cl.userAgent = ua }
}

// Client is a qBittorrent WebUI API client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

// New creates a qBittorrent [Client] for the given base URL (e.g. "http://localhost:8080").
func New(baseURL string, opts ...Option) *Client {
	jar, _ := cookiejar.New(nil)
	c := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: defaultTimeout, Jar: jar},
		userAgent:  defaultUserAgent,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int    `json:"-"`
	RawBody    string `json:"-"`
}

func (e *APIError) Error() string {
	if e.RawBody != "" {
		return fmt.Sprintf("qbit: HTTP %d: %s", e.StatusCode, e.RawBody)
	}
	return fmt.Sprintf("qbit: HTTP %d", e.StatusCode)
}

func (c *Client) doRequest(ctx context.Context, method, path string, form url.Values, dst any) error {
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("qbit: create request: %w", err)
	}

	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", c.userAgent)
	// Referer header is required by qBittorrent for CSRF protection.
	req.Header.Set("Referer", c.baseURL)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("qbit: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("qbit: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, RawBody: string(respBody)}
	}

	if dst != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, dst); err != nil {
			return fmt.Errorf("qbit: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) get(ctx context.Context, path string, params url.Values, dst any) error {
	if params != nil {
		path += "?" + params.Encode()
	}
	return c.doRequest(ctx, http.MethodGet, path, nil, dst)
}

func (c *Client) post(ctx context.Context, path string, form url.Values) error {
	return c.doRequest(ctx, http.MethodPost, path, form, nil)
}

func (c *Client) getRaw(ctx context.Context, path string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("qbit: create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", c.baseURL)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("qbit: GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("qbit: read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", &APIError{StatusCode: resp.StatusCode, RawBody: string(buf)}
	}
	return string(buf), nil
}

// Authentication.

// Login authenticates with the qBittorrent WebUI.
// The session cookie (SID) is automatically stored in the client's cookie jar.
func (c *Client) Login(ctx context.Context, username, password string) error {
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)
	return c.post(ctx, "/api/v2/auth/login", form)
}

// Logout ends the current session.
func (c *Client) Logout(ctx context.Context) error {
	return c.post(ctx, "/api/v2/auth/logout", nil)
}

// Application.

// Version returns the qBittorrent application version string.
func (c *Client) Version(ctx context.Context) (string, error) {
	return c.getRaw(ctx, "/api/v2/app/version")
}

// WebAPIVersion returns the WebUI API version string.
func (c *Client) WebAPIVersion(ctx context.Context) (string, error) {
	return c.getRaw(ctx, "/api/v2/app/webapiVersion")
}

// GetBuildInfo returns qBittorrent build information.
func (c *Client) GetBuildInfo(ctx context.Context) (*BuildInfo, error) {
	var out BuildInfo
	if err := c.get(ctx, "/api/v2/app/buildInfo", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPreferences returns the current qBittorrent preferences.
func (c *Client) GetPreferences(ctx context.Context) (*Preferences, error) {
	var out Preferences
	if err := c.get(ctx, "/api/v2/app/preferences", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DefaultSavePath returns the default save path for downloads.
func (c *Client) DefaultSavePath(ctx context.Context) (string, error) {
	return c.getRaw(ctx, "/api/v2/app/defaultSavePath")
}

// Torrents.

// ListTorrents returns a list of torrents matching the given options.
// Pass nil for opts to list all torrents.
func (c *Client) ListTorrents(ctx context.Context, opts *ListOptions) ([]Torrent, error) {
	p := url.Values{}
	if opts != nil {
		if opts.Filter != "" {
			p.Set("filter", opts.Filter)
		}
		if opts.Category != "" {
			p.Set("category", opts.Category)
		}
		if opts.Tag != "" {
			p.Set("tag", opts.Tag)
		}
		if opts.Sort != "" {
			p.Set("sort", opts.Sort)
		}
		if opts.Reverse {
			p.Set("reverse", "true")
		}
		if opts.Limit > 0 {
			p.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			p.Set("offset", strconv.Itoa(opts.Offset))
		}
		if opts.Hashes != "" {
			p.Set("hashes", opts.Hashes)
		}
	}
	var out []Torrent
	if err := c.get(ctx, "/api/v2/torrents/info", p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTorrentProperties returns detailed properties for a specific torrent.
func (c *Client) GetTorrentProperties(ctx context.Context, hash string) (*TorrentProperties, error) {
	var out TorrentProperties
	p := url.Values{}
	p.Set("hash", hash)
	if err := c.get(ctx, "/api/v2/torrents/properties", p, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTorrentTrackers returns trackers for a specific torrent.
func (c *Client) GetTorrentTrackers(ctx context.Context, hash string) ([]Tracker, error) {
	var out []Tracker
	p := url.Values{}
	p.Set("hash", hash)
	if err := c.get(ctx, "/api/v2/torrents/trackers", p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTorrentWebSeeds returns web seeds for a specific torrent.
func (c *Client) GetTorrentWebSeeds(ctx context.Context, hash string) ([]WebSeed, error) {
	var out []WebSeed
	p := url.Values{}
	p.Set("hash", hash)
	if err := c.get(ctx, "/api/v2/torrents/webseeds", p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTorrentFiles returns files within a specific torrent.
func (c *Client) GetTorrentFiles(ctx context.Context, hash string) ([]TorrentFile, error) {
	var out []TorrentFile
	p := url.Values{}
	p.Set("hash", hash)
	if err := c.get(ctx, "/api/v2/torrents/files", p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddTorrentURLs adds torrents by URL (magnet links or HTTP URLs to .torrent files).
func (c *Client) AddTorrentURLs(ctx context.Context, urls []string, opts *AddTorrentOptions) error {
	form := url.Values{}
	form.Set("urls", strings.Join(urls, "\n"))
	if opts != nil {
		setAddOptions(form, opts)
	}
	return c.post(ctx, "/api/v2/torrents/add", form)
}

func setAddOptions(form url.Values, opts *AddTorrentOptions) {
	if opts.SavePath != "" {
		form.Set("savepath", opts.SavePath)
	}
	if opts.Category != "" {
		form.Set("category", opts.Category)
	}
	if opts.Tags != "" {
		form.Set("tags", opts.Tags)
	}
	if opts.SkipChecking {
		form.Set("skip_checking", "true")
	}
	if opts.Paused {
		form.Set("paused", "true")
	}
	if opts.RootFolder {
		form.Set("root_folder", "true")
	}
	if opts.Rename != "" {
		form.Set("rename", opts.Rename)
	}
	if opts.UpLimit > 0 {
		form.Set("upLimit", strconv.FormatInt(opts.UpLimit, 10))
	}
	if opts.DlLimit > 0 {
		form.Set("dlLimit", strconv.FormatInt(opts.DlLimit, 10))
	}
	if opts.AutoTMM {
		form.Set("autoTMM", "true")
	}
	if opts.SequentialDownload {
		form.Set("sequentialDownload", "true")
	}
	if opts.FirstLastPiecePrio {
		form.Set("firstLastPiecePrio", "true")
	}
}

// DeleteTorrents removes torrents by their hashes.
// If deleteFiles is true, downloaded data is also deleted from disk.
func (c *Client) DeleteTorrents(ctx context.Context, hashes []string, deleteFiles bool) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	form.Set("deleteFiles", strconv.FormatBool(deleteFiles))
	return c.post(ctx, "/api/v2/torrents/delete", form)
}

// PauseTorrents pauses torrents by their hashes.
// Use "all" as single element to pause all torrents.
func (c *Client) PauseTorrents(ctx context.Context, hashes []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	return c.post(ctx, "/api/v2/torrents/pause", form)
}

// ResumeTorrents resumes torrents by their hashes.
// Use "all" as single element to resume all torrents.
func (c *Client) ResumeTorrents(ctx context.Context, hashes []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	return c.post(ctx, "/api/v2/torrents/resume", form)
}

// RecheckTorrents rechecks torrents by their hashes.
func (c *Client) RecheckTorrents(ctx context.Context, hashes []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	return c.post(ctx, "/api/v2/torrents/recheck", form)
}

// ReannounceTorrents reannounces torrents to their trackers.
func (c *Client) ReannounceTorrents(ctx context.Context, hashes []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	return c.post(ctx, "/api/v2/torrents/reannounce", form)
}

// SetTorrentLocation moves torrents to the specified path.
func (c *Client) SetTorrentLocation(ctx context.Context, hashes []string, location string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	form.Set("location", location)
	return c.post(ctx, "/api/v2/torrents/setLocation", form)
}

// RenameTorrent renames a torrent.
func (c *Client) RenameTorrent(ctx context.Context, hash, name string) error {
	form := url.Values{}
	form.Set("hash", hash)
	form.Set("name", name)
	return c.post(ctx, "/api/v2/torrents/rename", form)
}

// Categories and Tags.

// ListCategories returns all categories.
func (c *Client) ListCategories(ctx context.Context) (map[string]*Category, error) {
	out := make(map[string]*Category)
	if err := c.get(ctx, "/api/v2/torrents/categories", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateCategory creates a new category with an optional save path.
func (c *Client) CreateCategory(ctx context.Context, name, savePath string) error {
	form := url.Values{}
	form.Set("category", name)
	form.Set("savePath", savePath)
	return c.post(ctx, "/api/v2/torrents/createCategory", form)
}

// SetTorrentCategory assigns a category to torrents.
func (c *Client) SetTorrentCategory(ctx context.Context, hashes []string, category string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	form.Set("category", category)
	return c.post(ctx, "/api/v2/torrents/setCategory", form)
}

// ListTags returns all tags.
func (c *Client) ListTags(ctx context.Context) ([]string, error) {
	var out []string
	if err := c.get(ctx, "/api/v2/torrents/tags", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddTorrentTags adds tags to torrents.
func (c *Client) AddTorrentTags(ctx context.Context, hashes, tags []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	form.Set("tags", strings.Join(tags, ","))
	return c.post(ctx, "/api/v2/torrents/addTags", form)
}

// RemoveTorrentTags removes tags from torrents.
func (c *Client) RemoveTorrentTags(ctx context.Context, hashes, tags []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	form.Set("tags", strings.Join(tags, ","))
	return c.post(ctx, "/api/v2/torrents/removeTags", form)
}

// Transfer.

// GetTransferInfo returns global transfer statistics.
func (c *Client) GetTransferInfo(ctx context.Context) (*TransferInfo, error) {
	var out TransferInfo
	if err := c.get(ctx, "/api/v2/transfer/info", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGlobalDownloadLimit returns the global download speed limit in bytes/second.
// A value of 0 means no limit.
func (c *Client) GetGlobalDownloadLimit(ctx context.Context) (int64, error) {
	var out int64
	if err := c.get(ctx, "/api/v2/transfer/downloadLimit", nil, &out); err != nil {
		return 0, err
	}
	return out, nil
}

// SetGlobalDownloadLimit sets the global download speed limit in bytes/second.
// Set to 0 to remove the limit.
func (c *Client) SetGlobalDownloadLimit(ctx context.Context, limit int64) error {
	form := url.Values{}
	form.Set("limit", strconv.FormatInt(limit, 10))
	return c.post(ctx, "/api/v2/transfer/setDownloadLimit", form)
}

// GetGlobalUploadLimit returns the global upload speed limit in bytes/second.
func (c *Client) GetGlobalUploadLimit(ctx context.Context) (int64, error) {
	var out int64
	if err := c.get(ctx, "/api/v2/transfer/uploadLimit", nil, &out); err != nil {
		return 0, err
	}
	return out, nil
}

// SetGlobalUploadLimit sets the global upload speed limit in bytes/second.
func (c *Client) SetGlobalUploadLimit(ctx context.Context, limit int64) error {
	form := url.Values{}
	form.Set("limit", strconv.FormatInt(limit, 10))
	return c.post(ctx, "/api/v2/transfer/setUploadLimit", form)
}

// Sync.

// GetSyncMainData returns the sync main data.
// Pass rid=0 for a full update, or the previous rid for incremental updates.
func (c *Client) GetSyncMainData(ctx context.Context, rid int) (*SyncMainData, error) {
	var out SyncMainData
	p := url.Values{}
	p.Set("rid", strconv.Itoa(rid))
	if err := c.get(ctx, "/api/v2/sync/maindata", p, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Log.

// GetLog returns application log entries.
func (c *Client) GetLog(ctx context.Context, lastKnownID int) ([]LogEntry, error) {
	var out []LogEntry
	p := url.Values{}
	p.Set("last_known_id", strconv.Itoa(lastKnownID))
	if err := c.get(ctx, "/api/v2/log/main", p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetPeerLog returns peer log entries.
func (c *Client) GetPeerLog(ctx context.Context, lastKnownID int) ([]PeerLogEntry, error) {
	var out []PeerLogEntry
	p := url.Values{}
	p.Set("last_known_id", strconv.Itoa(lastKnownID))
	if err := c.get(ctx, "/api/v2/log/peers", p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Application (continued).

// SetPreferences updates qBittorrent preferences.
// Pass a JSON-serialisable map or struct with the preferences to change.
func (c *Client) SetPreferences(ctx context.Context, prefs any) error {
	data, err := json.Marshal(prefs)
	if err != nil {
		return fmt.Errorf("qbit: encode preferences: %w", err)
	}
	form := url.Values{}
	form.Set("json", string(data))
	return c.post(ctx, "/api/v2/app/setPreferences", form)
}

// Shutdown shuts down the qBittorrent application.
func (c *Client) Shutdown(ctx context.Context) error {
	return c.post(ctx, "/api/v2/app/shutdown", nil)
}

// Transfer (continued).

// GetSpeedLimitsMode returns whether alternative speed limits are enabled.
// Returns true if alternative speed limits are active.
func (c *Client) GetSpeedLimitsMode(ctx context.Context) (bool, error) {
	raw, err := c.getRaw(ctx, "/api/v2/transfer/speedLimitsMode")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(raw) == "1", nil
}

// ToggleSpeedLimitsMode toggles alternative speed limits on or off.
func (c *Client) ToggleSpeedLimitsMode(ctx context.Context) error {
	return c.post(ctx, "/api/v2/transfer/toggleSpeedLimitsMode", nil)
}

// BanPeers bans one or more peers. Each peer should be in "host:port" format.
func (c *Client) BanPeers(ctx context.Context, peers []string) error {
	form := url.Values{}
	form.Set("peers", strings.Join(peers, "|"))
	return c.post(ctx, "/api/v2/transfer/banPeers", form)
}

// Sync (continued).

// GetSyncTorrentPeers returns peer data for a specific torrent.
// Pass rid=0 for a full update, or the previous rid for incremental updates.
func (c *Client) GetSyncTorrentPeers(ctx context.Context, hash string, rid int) (*SyncTorrentPeers, error) {
	var out SyncTorrentPeers
	p := url.Values{}
	p.Set("hash", hash)
	p.Set("rid", strconv.Itoa(rid))
	if err := c.get(ctx, "/api/v2/sync/torrentPeers", p, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Torrent management (continued).

// GetTorrentPieceStates returns piece states for a torrent.
// Each value is: 0 = not downloaded, 1 = downloading, 2 = downloaded.
func (c *Client) GetTorrentPieceStates(ctx context.Context, hash string) ([]int, error) {
	var out []int
	p := url.Values{}
	p.Set("hash", hash)
	if err := c.get(ctx, "/api/v2/torrents/pieceStates", p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetTorrentPieceHashes returns piece hashes for a torrent.
func (c *Client) GetTorrentPieceHashes(ctx context.Context, hash string) ([]string, error) {
	var out []string
	p := url.Values{}
	p.Set("hash", hash)
	if err := c.get(ctx, "/api/v2/torrents/pieceHashes", p, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetFilePriority sets the priority of specific files within a torrent.
// File IDs correspond to the index field in [TorrentFile].
// Priority values: 0 = do not download, 1 = normal, 6 = high, 7 = maximum.
func (c *Client) SetFilePriority(ctx context.Context, hash string, fileIDs []int, priority int) error {
	ids := make([]string, len(fileIDs))
	for i, id := range fileIDs {
		ids[i] = strconv.Itoa(id)
	}
	form := url.Values{}
	form.Set("hash", hash)
	form.Set("id", strings.Join(ids, "|"))
	form.Set("priority", strconv.Itoa(priority))
	return c.post(ctx, "/api/v2/torrents/filePrio", form)
}

// SetTorrentDownloadLimit sets download speed limits for specific torrents (bytes/second).
func (c *Client) SetTorrentDownloadLimit(ctx context.Context, hashes []string, limit int64) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	form.Set("limit", strconv.FormatInt(limit, 10))
	return c.post(ctx, "/api/v2/torrents/setDownloadLimit", form)
}

// SetTorrentUploadLimit sets upload speed limits for specific torrents (bytes/second).
func (c *Client) SetTorrentUploadLimit(ctx context.Context, hashes []string, limit int64) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	form.Set("limit", strconv.FormatInt(limit, 10))
	return c.post(ctx, "/api/v2/torrents/setUploadLimit", form)
}

// SetShareLimits sets the share ratio and seeding time limits for specific torrents.
// Use -1 for ratioLimit/seedingTimeLimit to use global values, -2 for no limit.
func (c *Client) SetShareLimits(ctx context.Context, hashes []string, ratioLimit float64, seedingTimeLimit, inactiveSeedingTimeLimit int64) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	form.Set("ratioLimit", strconv.FormatFloat(ratioLimit, 'f', -1, 64))
	form.Set("seedingTimeLimit", strconv.FormatInt(seedingTimeLimit, 10))
	form.Set("inactiveSeedingTimeLimit", strconv.FormatInt(inactiveSeedingTimeLimit, 10))
	return c.post(ctx, "/api/v2/torrents/setShareLimits", form)
}

// IncreasePriority increases the priority of torrents in the queue.
func (c *Client) IncreasePriority(ctx context.Context, hashes []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	return c.post(ctx, "/api/v2/torrents/increasePrio", form)
}

// DecreasePriority decreases the priority of torrents in the queue.
func (c *Client) DecreasePriority(ctx context.Context, hashes []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	return c.post(ctx, "/api/v2/torrents/decreasePrio", form)
}

// TopPriority moves torrents to the top of the queue.
func (c *Client) TopPriority(ctx context.Context, hashes []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	return c.post(ctx, "/api/v2/torrents/topPrio", form)
}

// BottomPriority moves torrents to the bottom of the queue.
func (c *Client) BottomPriority(ctx context.Context, hashes []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	return c.post(ctx, "/api/v2/torrents/bottomPrio", form)
}

// SetForceStart enables or disables force-start for torrents.
func (c *Client) SetForceStart(ctx context.Context, hashes []string, value bool) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	form.Set("value", strconv.FormatBool(value))
	return c.post(ctx, "/api/v2/torrents/setForceStart", form)
}

// SetSuperSeeding enables or disables super-seeding for torrents.
func (c *Client) SetSuperSeeding(ctx context.Context, hashes []string, value bool) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	form.Set("value", strconv.FormatBool(value))
	return c.post(ctx, "/api/v2/torrents/setSuperSeeding", form)
}

// SetAutoManagement enables or disables automatic torrent management for torrents.
func (c *Client) SetAutoManagement(ctx context.Context, hashes []string, enable bool) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	form.Set("enable", strconv.FormatBool(enable))
	return c.post(ctx, "/api/v2/torrents/setAutoManagement", form)
}

// ToggleSequentialDownload toggles sequential download for torrents.
func (c *Client) ToggleSequentialDownload(ctx context.Context, hashes []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	return c.post(ctx, "/api/v2/torrents/toggleSequentialDownload", form)
}

// ToggleFirstLastPiecePrio toggles first/last piece priority for torrents.
func (c *Client) ToggleFirstLastPiecePrio(ctx context.Context, hashes []string) error {
	form := url.Values{}
	form.Set("hashes", strings.Join(hashes, "|"))
	return c.post(ctx, "/api/v2/torrents/toggleFirstLastPiecePrio", form)
}

// Tracker management.

// AddTrackers adds trackers to a torrent.
func (c *Client) AddTrackers(ctx context.Context, hash string, urls []string) error {
	form := url.Values{}
	form.Set("hash", hash)
	form.Set("urls", strings.Join(urls, "\n"))
	return c.post(ctx, "/api/v2/torrents/addTrackers", form)
}

// EditTracker replaces a tracker URL for a torrent.
func (c *Client) EditTracker(ctx context.Context, hash, origURL, newURL string) error {
	form := url.Values{}
	form.Set("hash", hash)
	form.Set("origUrl", origURL)
	form.Set("newUrl", newURL)
	return c.post(ctx, "/api/v2/torrents/editTracker", form)
}

// RemoveTrackers removes trackers from a torrent.
func (c *Client) RemoveTrackers(ctx context.Context, hash string, urls []string) error {
	form := url.Values{}
	form.Set("hash", hash)
	form.Set("urls", strings.Join(urls, "|"))
	return c.post(ctx, "/api/v2/torrents/removeTrackers", form)
}

// Category management (continued).

// EditCategory edits an existing category.
func (c *Client) EditCategory(ctx context.Context, name, savePath string) error {
	form := url.Values{}
	form.Set("category", name)
	form.Set("savePath", savePath)
	return c.post(ctx, "/api/v2/torrents/editCategory", form)
}

// RemoveCategories removes one or more categories.
func (c *Client) RemoveCategories(ctx context.Context, categories []string) error {
	form := url.Values{}
	form.Set("categories", strings.Join(categories, "\n"))
	return c.post(ctx, "/api/v2/torrents/removeCategories", form)
}

// Tag management (continued).

// CreateTags creates new tags.
func (c *Client) CreateTags(ctx context.Context, tags []string) error {
	form := url.Values{}
	form.Set("tags", strings.Join(tags, ","))
	return c.post(ctx, "/api/v2/torrents/createTags", form)
}

// DeleteTags deletes tags.
func (c *Client) DeleteTags(ctx context.Context, tags []string) error {
	form := url.Values{}
	form.Set("tags", strings.Join(tags, ","))
	return c.post(ctx, "/api/v2/torrents/deleteTags", form)
}

// Content renaming.

// RenameFile renames a file within a torrent.
func (c *Client) RenameFile(ctx context.Context, hash, oldPath, newPath string) error {
	form := url.Values{}
	form.Set("hash", hash)
	form.Set("oldPath", oldPath)
	form.Set("newPath", newPath)
	return c.post(ctx, "/api/v2/torrents/renameFile", form)
}

// RenameFolder renames a folder within a torrent.
func (c *Client) RenameFolder(ctx context.Context, hash, oldPath, newPath string) error {
	form := url.Values{}
	form.Set("hash", hash)
	form.Set("oldPath", oldPath)
	form.Set("newPath", newPath)
	return c.post(ctx, "/api/v2/torrents/renameFolder", form)
}
