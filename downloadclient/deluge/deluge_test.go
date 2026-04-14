package deluge_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/downloadclient/deluge"
)

type rpcRequest struct {
	ID     int             `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

func newRPCServer(t *testing.T, wantMethod string, result any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("HTTP method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/json" {
			t.Errorf("path = %q, want /json", r.URL.Path)
		}
		var req rpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Method != wantMethod {
			t.Errorf("RPC method = %q, want %q", req.Method, wantMethod)
		}
		w.Header().Set("Content-Type", "application/json")
		resultJSON, _ := json.Marshal(result)
		resp := map[string]any{
			"id":     req.ID,
			"result": json.RawMessage(resultJSON),
			"error":  nil,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func loginServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req rpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{"id": req.ID, "result": true, "error": nil}
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func TestLogin(t *testing.T) {
	t.Parallel()

	ts := loginServer(t)
	defer ts.Close()

	c := deluge.New(ts.URL)
	if err := c.Login(context.Background(), "deluge"); err != nil {
		t.Fatal(err)
	}
}

func TestLoginFailed(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req rpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{"id": req.ID, "result": false, "error": nil}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := deluge.New(ts.URL)
	err := c.Login(context.Background(), "wrong")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *deluge.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestGetTorrentsStatus(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"abc123": map[string]any{
			"hash": "abc123", "name": "Ubuntu", "state": "Seeding",
			"progress": 100.0, "total_size": 4000000000,
		},
		"def456": map[string]any{
			"hash": "def456", "name": "Fedora", "state": "Downloading",
			"progress": 50.0, "total_size": 2000000000,
		},
	}
	ts := newRPCServer(t, "core.get_torrents_status", result)
	defer ts.Close()

	c := deluge.New(ts.URL)
	torrents, err := c.GetTorrentsStatus(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(torrents) != 2 {
		t.Fatalf("len = %d, want 2", len(torrents))
	}
	if torrents["abc123"].Name != "Ubuntu" {
		t.Errorf("Name = %q, want %q", torrents["abc123"].Name, "Ubuntu")
	}
}

func TestGetTorrentStatus(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"hash": "abc123", "name": "Ubuntu", "state": "Seeding", "progress": 100.0,
	}
	ts := newRPCServer(t, "core.get_torrent_status", result)
	defer ts.Close()

	c := deluge.New(ts.URL)
	torrent, err := c.GetTorrentStatus(context.Background(), "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if torrent.State != "Seeding" {
		t.Errorf("State = %q, want %q", torrent.State, "Seeding")
	}
}

func TestAddTorrentURL(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "core.add_torrent_url", "abc123def456")
	defer ts.Close()

	c := deluge.New(ts.URL)
	hash, err := c.AddTorrentURL(context.Background(), "magnet:?xt=urn:btih:abc123", nil)
	if err != nil {
		t.Fatal(err)
	}
	if hash != "abc123def456" {
		t.Errorf("hash = %q, want %q", hash, "abc123def456")
	}
}

func TestRemoveTorrent(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "core.remove_torrent", true)
	defer ts.Close()

	c := deluge.New(ts.URL)
	if err := c.RemoveTorrent(context.Background(), "abc123", true); err != nil {
		t.Fatal(err)
	}
}

func TestPauseTorrent(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "core.pause_torrent", nil)
	defer ts.Close()

	c := deluge.New(ts.URL)
	if err := c.PauseTorrent(context.Background(), "abc123"); err != nil {
		t.Fatal(err)
	}
}

func TestResumeTorrent(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "core.resume_torrent", nil)
	defer ts.Close()

	c := deluge.New(ts.URL)
	if err := c.ResumeTorrent(context.Background(), "def456"); err != nil {
		t.Fatal(err)
	}
}

func TestForceRecheck(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "core.force_recheck", nil)
	defer ts.Close()

	c := deluge.New(ts.URL)
	if err := c.ForceRecheck(context.Background(), []string{"abc123"}); err != nil {
		t.Fatal(err)
	}
}

func TestGetVersion(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "daemon.info", "2.1.1")
	defer ts.Close()

	c := deluge.New(ts.URL)
	v, err := c.GetVersion(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if v != "2.1.1" {
		t.Errorf("version = %q, want %q", v, "2.1.1")
	}
}

func TestGetSessionStatus(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"payload_download_rate": 50000, "payload_upload_rate": 10000,
		"dht_nodes": 200,
	}
	ts := newRPCServer(t, "core.get_session_status", result)
	defer ts.Close()

	c := deluge.New(ts.URL)
	stats, err := c.GetSessionStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if stats.DownloadRate != 50000 {
		t.Errorf("DownloadRate = %d, want 50000", stats.DownloadRate)
	}
}

func TestGetFreeSpace(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "core.get_free_space", 500000000000)
	defer ts.Close()

	c := deluge.New(ts.URL)
	space, err := c.GetFreeSpace(context.Background(), "/downloads")
	if err != nil {
		t.Fatal(err)
	}
	if space != 500000000000 {
		t.Errorf("space = %d, want 500000000000", space)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req rpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{
			"id":     req.ID,
			"result": nil,
			"error":  map[string]any{"message": "torrent not found", "code": 2},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := deluge.New(ts.URL)
	_, err := c.GetTorrentsStatus(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *deluge.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Message != "torrent not found" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "torrent not found")
	}
}

func TestConnected(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "web.connected", true)
	defer ts.Close()

	c := deluge.New(ts.URL)
	ok, err := c.Connected(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("connected = false, want true")
	}
}

func TestSetTorrentLabel(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "label.set_torrent", nil)
	defer ts.Close()

	c := deluge.New(ts.URL)
	if err := c.SetTorrentLabel(context.Background(), "abc123", "movies"); err != nil {
		t.Fatal(err)
	}
}

func TestMoveTorrent(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "core.move_storage", nil)
	defer ts.Close()

	c := deluge.New(ts.URL)
	if err := c.MoveTorrent(context.Background(), []string{"abc123"}, "/new/path"); err != nil {
		t.Fatal(err)
	}
}
