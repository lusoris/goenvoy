package qbit_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/downloadclient/qbit"
)

func newTestServer(t *testing.T, wantPath string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		w.Header().Set("Content-Type", "application/json")
		if response != nil {
			_ = json.NewEncoder(w).Encode(response)
		}
	}))
}

func newLoginServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/auth/login":
			if r.Method != http.MethodPost {
				t.Errorf("login method = %q, want POST", r.Method)
			}
			http.SetCookie(w, &http.Cookie{Name: "SID", Value: "test-session-id", Path: "/"})
			w.WriteHeader(http.StatusOK)
		case "/api/v2/auth/logout":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusForbidden)
		}
	}))
}

func TestLogin(t *testing.T) {
	t.Parallel()

	ts := newLoginServer(t)
	defer ts.Close()

	c := qbit.New(ts.URL)
	if err := c.Login(context.Background(), "admin", "adminadmin"); err != nil {
		t.Fatal(err)
	}
}

func TestLogout(t *testing.T) {
	t.Parallel()

	ts := newLoginServer(t)
	defer ts.Close()

	c := qbit.New(ts.URL)
	_ = c.Login(context.Background(), "admin", "adminadmin")
	if err := c.Logout(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestVersion(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("v4.6.7"))
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	v, err := c.Version(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if v != "v4.6.7" {
		t.Errorf("version = %q, want %q", v, "v4.6.7")
	}
}

func TestWebAPIVersion(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("2.10.5"))
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	v, err := c.WebAPIVersion(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if v != "2.10.5" {
		t.Errorf("version = %q, want %q", v, "2.10.5")
	}
}

func TestGetBuildInfo(t *testing.T) {
	t.Parallel()

	info := qbit.BuildInfo{Qt: "6.7.2", Libtorrent: "2.0.10.0", Boost: "1.86", OpenSSL: "3.3.1", Bitness: 64}
	ts := newTestServer(t, "/api/v2/app/buildInfo", info)
	defer ts.Close()

	c := qbit.New(ts.URL)
	b, err := c.GetBuildInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if b.Qt != "6.7.2" {
		t.Errorf("Qt = %q, want %q", b.Qt, "6.7.2")
	}
	if b.Bitness != 64 {
		t.Errorf("Bitness = %d, want %d", b.Bitness, 64)
	}
}

func TestGetPreferences(t *testing.T) {
	t.Parallel()

	prefs := qbit.Preferences{SavePath: "/downloads", DlLimit: 5000000, QueueingEnabled: true}
	ts := newTestServer(t, "/api/v2/app/preferences", prefs)
	defer ts.Close()

	c := qbit.New(ts.URL)
	p, err := c.GetPreferences(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if p.SavePath != "/downloads" {
		t.Errorf("SavePath = %q, want %q", p.SavePath, "/downloads")
	}
	if p.DlLimit != 5000000 {
		t.Errorf("DlLimit = %d, want %d", p.DlLimit, 5000000)
	}
}

func TestDefaultSavePath(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("/downloads/complete"))
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	p, err := c.DefaultSavePath(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if p != "/downloads/complete" {
		t.Errorf("path = %q, want %q", p, "/downloads/complete")
	}
}

func TestListTorrents(t *testing.T) {
	t.Parallel()

	torrents := []qbit.Torrent{
		{Hash: "abc123", Name: "Ubuntu 24.04", State: "downloading", Progress: 0.45, Size: 4000000000},
		{Hash: "def456", Name: "Fedora 40", State: "seeding", Progress: 1.0, Size: 2000000000},
	}
	ts := newTestServer(t, "/api/v2/torrents/info", torrents)
	defer ts.Close()

	c := qbit.New(ts.URL)
	result, err := c.ListTorrents(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result[0].Name != "Ubuntu 24.04" {
		t.Errorf("Name = %q, want %q", result[0].Name, "Ubuntu 24.04")
	}
	if result[1].State != "seeding" {
		t.Errorf("State = %q, want %q", result[1].State, "seeding")
	}
}

func TestListTorrentsWithOptions(t *testing.T) {
	t.Parallel()

	torrents := []qbit.Torrent{{Hash: "abc123", Name: "Test", Category: "movies"}}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("filter"); got != "downloading" {
			t.Errorf("filter = %q, want %q", got, "downloading")
		}
		if got := r.URL.Query().Get("category"); got != "movies" {
			t.Errorf("category = %q, want %q", got, "movies")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(torrents)
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	result, err := c.ListTorrents(context.Background(), &qbit.ListOptions{Filter: "downloading", Category: "movies"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestGetTorrentProperties(t *testing.T) {
	t.Parallel()

	props := qbit.TorrentProperties{SavePath: "/data/movies", TotalSize: 50000000, Seeds: 42, Peers: 10}
	ts := newTestServer(t, "/api/v2/torrents/properties", props)
	defer ts.Close()

	c := qbit.New(ts.URL)
	p, err := c.GetTorrentProperties(context.Background(), "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if p.SavePath != "/data/movies" {
		t.Errorf("SavePath = %q, want %q", p.SavePath, "/data/movies")
	}
	if p.Seeds != 42 {
		t.Errorf("Seeds = %d, want %d", p.Seeds, 42)
	}
}

func TestGetTorrentTrackers(t *testing.T) {
	t.Parallel()

	trackers := []qbit.Tracker{{URL: "udp://tracker.example.com:1337", Status: 2, NumPeers: 150}}
	ts := newTestServer(t, "/api/v2/torrents/trackers", trackers)
	defer ts.Close()

	c := qbit.New(ts.URL)
	result, err := c.GetTorrentTrackers(context.Background(), "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].NumPeers != 150 {
		t.Errorf("NumPeers = %d, want %d", result[0].NumPeers, 150)
	}
}

func TestGetTorrentFiles(t *testing.T) {
	t.Parallel()

	files := []qbit.TorrentFile{{Index: 0, Name: "movie.mkv", Size: 5000000000, Progress: 0.8, Priority: 1}}
	ts := newTestServer(t, "/api/v2/torrents/files", files)
	defer ts.Close()

	c := qbit.New(ts.URL)
	result, err := c.GetTorrentFiles(context.Background(), "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Name != "movie.mkv" {
		t.Errorf("Name = %q, want %q", result[0].Name, "movie.mkv")
	}
}

func TestAddTorrentURLs(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/v2/torrents/add" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/api/v2/torrents/add")
		}
		_ = r.ParseForm()
		if got := r.FormValue("urls"); got == "" {
			t.Error("urls is empty")
		}
		if got := r.FormValue("category"); got != "movies" {
			t.Errorf("category = %q, want %q", got, "movies")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	err := c.AddTorrentURLs(context.Background(), []string{"magnet:?xt=urn:btih:abc123"}, &qbit.AddTorrentOptions{Category: "movies"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTorrents(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		_ = r.ParseForm()
		if got := r.FormValue("hashes"); got != "abc123|def456" {
			t.Errorf("hashes = %q, want %q", got, "abc123|def456")
		}
		if got := r.FormValue("deleteFiles"); got != "true" {
			t.Errorf("deleteFiles = %q, want %q", got, "true")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	if err := c.DeleteTorrents(context.Background(), []string{"abc123", "def456"}, true); err != nil {
		t.Fatal(err)
	}
}

func TestPauseTorrents(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/torrents/pause" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	if err := c.PauseTorrents(context.Background(), []string{"abc123"}); err != nil {
		t.Fatal(err)
	}
}

func TestResumeTorrents(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/torrents/resume" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	if err := c.ResumeTorrents(context.Background(), []string{"abc123"}); err != nil {
		t.Fatal(err)
	}
}

func TestGetTransferInfo(t *testing.T) {
	t.Parallel()

	info := qbit.TransferInfo{
		DlInfoSpeed: 5000000, UpInfoSpeed: 1000000, DHTNodes: 450,
		ConnectionStatus: "connected", FreeSpaceOnDisk: 500000000000,
	}
	ts := newTestServer(t, "/api/v2/transfer/info", info)
	defer ts.Close()

	c := qbit.New(ts.URL)
	result, err := c.GetTransferInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.DlInfoSpeed != 5000000 {
		t.Errorf("DlInfoSpeed = %d, want %d", result.DlInfoSpeed, 5000000)
	}
	if result.ConnectionStatus != "connected" {
		t.Errorf("ConnectionStatus = %q, want %q", result.ConnectionStatus, "connected")
	}
}

func TestListCategories(t *testing.T) {
	t.Parallel()

	cats := map[string]*qbit.Category{
		"movies": {Name: "movies", SavePath: "/data/movies"},
		"tv":     {Name: "tv", SavePath: "/data/tv"},
	}
	ts := newTestServer(t, "/api/v2/torrents/categories", cats)
	defer ts.Close()

	c := qbit.New(ts.URL)
	result, err := c.ListCategories(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result["movies"].SavePath != "/data/movies" {
		t.Errorf("SavePath = %q, want %q", result["movies"].SavePath, "/data/movies")
	}
}

func TestListTags(t *testing.T) {
	t.Parallel()

	tags := []string{"4k", "remux", "web-dl"}
	ts := newTestServer(t, "/api/v2/torrents/tags", tags)
	defer ts.Close()

	c := qbit.New(ts.URL)
	result, err := c.ListTags(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Fatalf("len = %d, want 3", len(result))
	}
	if result[0] != "4k" {
		t.Errorf("tag = %q, want %q", result[0], "4k")
	}
}

func TestGetSyncMainData(t *testing.T) {
	t.Parallel()

	data := qbit.SyncMainData{
		RID:        1,
		FullUpdate: true,
		Torrents:   map[string]*qbit.Torrent{"abc": {Name: "Test Torrent", Hash: "abc"}},
		Tags:       []string{"hd"},
	}
	ts := newTestServer(t, "/api/v2/sync/maindata", data)
	defer ts.Close()

	c := qbit.New(ts.URL)
	result, err := c.GetSyncMainData(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if !result.FullUpdate {
		t.Error("FullUpdate = false, want true")
	}
	if result.Torrents["abc"].Name != "Test Torrent" {
		t.Errorf("Name = %q, want %q", result.Torrents["abc"].Name, "Test Torrent")
	}
}

func TestGetLog(t *testing.T) {
	t.Parallel()

	logs := []qbit.LogEntry{
		{ID: 1, Message: "qBittorrent started", Timestamp: 1700000000, Type: 1},
		{ID: 2, Message: "Torrent added", Timestamp: 1700000060, Type: 2},
	}
	ts := newTestServer(t, "/api/v2/log/main", logs)
	defer ts.Close()

	c := qbit.New(ts.URL)
	result, err := c.GetLog(context.Background(), -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result[0].Message != "qBittorrent started" {
		t.Errorf("Message = %q, want %q", result[0].Message, "qBittorrent started")
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Forbidden"))
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	_, err := c.ListTorrents(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *qbit.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusForbidden)
	}
}

func TestRecheckTorrents(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/torrents/recheck" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	if err := c.RecheckTorrents(context.Background(), []string{"abc123"}); err != nil {
		t.Fatal(err)
	}
}

func TestCreateCategory(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/torrents/createCategory" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_ = r.ParseForm()
		if got := r.FormValue("category"); got != "movies" {
			t.Errorf("category = %q, want %q", got, "movies")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	if err := c.CreateCategory(context.Background(), "movies", "/data/movies"); err != nil {
		t.Fatal(err)
	}
}

func TestGetPeerLog(t *testing.T) {
	t.Parallel()

	logs := []qbit.PeerLogEntry{
		{ID: 1, IP: "192.168.1.100", Timestamp: 1700000000, Blocked: false},
	}
	ts := newTestServer(t, "/api/v2/log/peers", logs)
	defer ts.Close()

	c := qbit.New(ts.URL)
	result, err := c.GetPeerLog(context.Background(), -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].IP != "192.168.1.100" {
		t.Errorf("IP = %q, want %q", result[0].IP, "192.168.1.100")
	}
}

func TestSetGlobalDownloadLimit(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/transfer/setDownloadLimit" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_ = r.ParseForm()
		if got := r.FormValue("limit"); got != "10000000" {
			t.Errorf("limit = %q, want %q", got, "10000000")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := qbit.New(ts.URL)
	if err := c.SetGlobalDownloadLimit(context.Background(), 10000000); err != nil {
		t.Fatal(err)
	}
}

func TestGetTorrentWebSeeds(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/webseeds", []qbit.WebSeed{{URL: "http://example.com"}})
	defer ts.Close()
	c := qbit.New(ts.URL)
	seeds, err := c.GetTorrentWebSeeds(context.Background(), "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if len(seeds) != 1 {
		t.Fatalf("len = %d", len(seeds))
	}
}

func TestReannounceTorrents(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/reannounce", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.ReannounceTorrents(context.Background(), []string{"abc"}); err != nil {
		t.Fatal(err)
	}
}

func TestSetTorrentLocation(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/setLocation", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.SetTorrentLocation(context.Background(), []string{"abc"}, "/new/path"); err != nil {
		t.Fatal(err)
	}
}

func TestRenameTorrent(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/rename", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.RenameTorrent(context.Background(), "abc", "newname"); err != nil {
		t.Fatal(err)
	}
}

func TestSetTorrentCategory(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/setCategory", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.SetTorrentCategory(context.Background(), []string{"abc"}, "movies"); err != nil {
		t.Fatal(err)
	}
}

func TestAddTorrentTags(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/addTags", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.AddTorrentTags(context.Background(), []string{"abc"}, []string{"tag1", "tag2"}); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveTorrentTags(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/removeTags", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.RemoveTorrentTags(context.Background(), []string{"abc"}, []string{"tag1"}); err != nil {
		t.Fatal(err)
	}
}

func TestGetGlobalDownloadLimit(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/transfer/downloadLimit", int64(5000000))
	defer ts.Close()
	c := qbit.New(ts.URL)
	limit, err := c.GetGlobalDownloadLimit(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if limit != 5000000 {
		t.Errorf("limit = %d", limit)
	}
}

func TestGetGlobalUploadLimit(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/transfer/uploadLimit", int64(3000000))
	defer ts.Close()
	c := qbit.New(ts.URL)
	limit, err := c.GetGlobalUploadLimit(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if limit != 3000000 {
		t.Errorf("limit = %d", limit)
	}
}

func TestSetGlobalUploadLimit(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/transfer/setUploadLimit", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.SetGlobalUploadLimit(context.Background(), 5000000); err != nil {
		t.Fatal(err)
	}
}

func TestSetPreferences(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/app/setPreferences", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.SetPreferences(context.Background(), map[string]any{"dl_limit": 1000}); err != nil {
		t.Fatal(err)
	}
}

func TestShutdown(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/app/shutdown", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestGetSpeedLimitsMode(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/transfer/speedLimitsMode" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte("1"))
	}))
	defer ts.Close()
	c := qbit.New(ts.URL)
	mode, err := c.GetSpeedLimitsMode(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !mode {
		t.Error("expected true")
	}
}

func TestToggleSpeedLimitsMode(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/transfer/toggleSpeedLimitsMode", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.ToggleSpeedLimitsMode(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestBanPeers(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/transfer/banPeers", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.BanPeers(context.Background(), []string{"1.2.3.4:6881"}); err != nil {
		t.Fatal(err)
	}
}

func TestGetSyncTorrentPeers(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/sync/torrentPeers", qbit.SyncTorrentPeers{RID: 1, FullData: true})
	defer ts.Close()
	c := qbit.New(ts.URL)
	result, err := c.GetSyncTorrentPeers(context.Background(), "abc", 0)
	if err != nil {
		t.Fatal(err)
	}
	if result.RID != 1 {
		t.Errorf("rid = %d", result.RID)
	}
}

func TestGetTorrentPieceStates(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/pieceStates", []int{0, 1, 2})
	defer ts.Close()
	c := qbit.New(ts.URL)
	states, err := c.GetTorrentPieceStates(context.Background(), "abc")
	if err != nil {
		t.Fatal(err)
	}
	if len(states) != 3 {
		t.Fatalf("len = %d", len(states))
	}
}

func TestGetTorrentPieceHashes(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/pieceHashes", []string{"hash1", "hash2"})
	defer ts.Close()
	c := qbit.New(ts.URL)
	hashes, err := c.GetTorrentPieceHashes(context.Background(), "abc")
	if err != nil {
		t.Fatal(err)
	}
	if len(hashes) != 2 {
		t.Fatalf("len = %d", len(hashes))
	}
}

func TestSetFilePriority(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/filePrio", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.SetFilePriority(context.Background(), "abc", []int{0, 1}, 7); err != nil {
		t.Fatal(err)
	}
}

func TestSetTorrentDownloadLimit(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/setDownloadLimit", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.SetTorrentDownloadLimit(context.Background(), []string{"abc"}, 100000); err != nil {
		t.Fatal(err)
	}
}

func TestSetTorrentUploadLimit(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/setUploadLimit", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.SetTorrentUploadLimit(context.Background(), []string{"abc"}, 50000); err != nil {
		t.Fatal(err)
	}
}

func TestSetShareLimits(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/setShareLimits", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.SetShareLimits(context.Background(), []string{"abc"}, 2.0, -1, -1); err != nil {
		t.Fatal(err)
	}
}

func TestIncreasePriority(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/increasePrio", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.IncreasePriority(context.Background(), []string{"abc"}); err != nil {
		t.Fatal(err)
	}
}

func TestDecreasePriority(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/decreasePrio", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.DecreasePriority(context.Background(), []string{"abc"}); err != nil {
		t.Fatal(err)
	}
}

func TestTopPriority(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/topPrio", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.TopPriority(context.Background(), []string{"abc"}); err != nil {
		t.Fatal(err)
	}
}

func TestBottomPriority(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/bottomPrio", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.BottomPriority(context.Background(), []string{"abc"}); err != nil {
		t.Fatal(err)
	}
}

func TestSetForceStart(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/setForceStart", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.SetForceStart(context.Background(), []string{"abc"}, true); err != nil {
		t.Fatal(err)
	}
}

func TestSetSuperSeeding(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/setSuperSeeding", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.SetSuperSeeding(context.Background(), []string{"abc"}, true); err != nil {
		t.Fatal(err)
	}
}

func TestSetAutoManagement(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/setAutoManagement", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.SetAutoManagement(context.Background(), []string{"abc"}, true); err != nil {
		t.Fatal(err)
	}
}

func TestToggleSequentialDownload(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/toggleSequentialDownload", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.ToggleSequentialDownload(context.Background(), []string{"abc"}); err != nil {
		t.Fatal(err)
	}
}

func TestToggleFirstLastPiecePrio(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/toggleFirstLastPiecePrio", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.ToggleFirstLastPiecePrio(context.Background(), []string{"abc"}); err != nil {
		t.Fatal(err)
	}
}

func TestAddTrackers(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/addTrackers", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.AddTrackers(context.Background(), "abc", []string{"http://tracker1.com/announce"}); err != nil {
		t.Fatal(err)
	}
}

func TestEditTracker(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/editTracker", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.EditTracker(context.Background(), "abc", "http://old.com/announce", "http://new.com/announce"); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveTrackers(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/removeTrackers", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.RemoveTrackers(context.Background(), "abc", []string{"http://tracker1.com/announce"}); err != nil {
		t.Fatal(err)
	}
}

func TestEditCategory(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/editCategory", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.EditCategory(context.Background(), "movies", "/new/path"); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveCategories(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/removeCategories", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.RemoveCategories(context.Background(), []string{"movies", "tv"}); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTags(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/createTags", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.CreateTags(context.Background(), []string{"tag1", "tag2"}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTags(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/deleteTags", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.DeleteTags(context.Background(), []string{"tag1"}); err != nil {
		t.Fatal(err)
	}
}

func TestRenameFile(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/renameFile", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.RenameFile(context.Background(), "abc", "old.mkv", "new.mkv"); err != nil {
		t.Fatal(err)
	}
}

func TestRenameFolder(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v2/torrents/renameFolder", nil)
	defer ts.Close()
	c := qbit.New(ts.URL)
	if err := c.RenameFolder(context.Background(), "abc", "OldFolder", "NewFolder"); err != nil {
		t.Fatal(err)
	}
}
