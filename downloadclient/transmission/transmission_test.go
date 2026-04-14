package transmission_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/downloadclient/transmission"
)

type rpcRequest struct {
	Method    string          `json:"method"`
	Arguments json.RawMessage `json:"arguments"`
}

func newRPCServer(t *testing.T, wantMethod string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("HTTP method = %q, want POST", r.Method)
		}
		var req rpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Method != wantMethod {
			t.Errorf("RPC method = %q, want %q", req.Method, wantMethod)
		}
		w.Header().Set("Content-Type", "application/json")
		args, _ := json.Marshal(response)
		resp := map[string]any{
			"result":    "success",
			"arguments": json.RawMessage(args),
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func TestGetTorrents(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"torrents": []map[string]any{
			{"id": 1, "name": "Ubuntu 24.04", "status": 6, "percentDone": 1.0, "totalSize": 4000000000},
			{"id": 2, "name": "Fedora 40", "status": 4, "percentDone": 0.5, "totalSize": 2000000000},
		},
	}
	ts := newRPCServer(t, "torrent-get", result)
	defer ts.Close()

	c := transmission.New(ts.URL)
	torrents, err := c.GetTorrents(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(torrents) != 2 {
		t.Fatalf("len = %d, want 2", len(torrents))
	}
	if torrents[0].Name != "Ubuntu 24.04" {
		t.Errorf("Name = %q, want %q", torrents[0].Name, "Ubuntu 24.04")
	}
	if torrents[1].Status != 4 {
		t.Errorf("Status = %d, want 4", torrents[1].Status)
	}
}

func TestGetTorrentsWithIDs(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"torrents": []map[string]any{
			{"id": 5, "name": "Specific", "percentDone": 0.75},
		},
	}
	ts := newRPCServer(t, "torrent-get", result)
	defer ts.Close()

	c := transmission.New(ts.URL)
	torrents, err := c.GetTorrents(context.Background(), []int{5})
	if err != nil {
		t.Fatal(err)
	}
	if len(torrents) != 1 {
		t.Fatalf("len = %d, want 1", len(torrents))
	}
	if torrents[0].ID != 5 {
		t.Errorf("ID = %d, want 5", torrents[0].ID)
	}
}

func TestAddTorrentURL(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"torrent-added": map[string]any{
			"id": 10, "name": "New Torrent", "hashString": "abc123def456",
		},
	}
	ts := newRPCServer(t, "torrent-add", result)
	defer ts.Close()

	c := transmission.New(ts.URL)
	added, err := c.AddTorrentURL(context.Background(), "magnet:?xt=urn:btih:abc123", &transmission.AddTorrentOptions{
		DownloadDir: "/data/movies",
		Labels:      []string{"movie"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if added.Name != "New Torrent" {
		t.Errorf("Name = %q, want %q", added.Name, "New Torrent")
	}
	if added.ID != 10 {
		t.Errorf("ID = %d, want 10", added.ID)
	}
}

func TestStartTorrents(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "torrent-start", nil)
	defer ts.Close()

	c := transmission.New(ts.URL)
	if err := c.StartTorrents(context.Background(), []int{1, 2}); err != nil {
		t.Fatal(err)
	}
}

func TestStopTorrents(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "torrent-stop", nil)
	defer ts.Close()

	c := transmission.New(ts.URL)
	if err := c.StopTorrents(context.Background(), []int{3}); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveTorrents(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "torrent-remove", nil)
	defer ts.Close()

	c := transmission.New(ts.URL)
	if err := c.RemoveTorrents(context.Background(), []int{1}, true); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyTorrents(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "torrent-verify", nil)
	defer ts.Close()

	c := transmission.New(ts.URL)
	if err := c.VerifyTorrents(context.Background(), []int{1}); err != nil {
		t.Fatal(err)
	}
}

func TestReannounceTorrents(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "torrent-reannounce", nil)
	defer ts.Close()

	c := transmission.New(ts.URL)
	if err := c.ReannounceTorrents(context.Background(), []int{1}); err != nil {
		t.Fatal(err)
	}
}

func TestMoveTorrents(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "torrent-set-location", nil)
	defer ts.Close()

	c := transmission.New(ts.URL)
	if err := c.MoveTorrents(context.Background(), []int{1}, "/new/path", true); err != nil {
		t.Fatal(err)
	}
}

func TestGetSession(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"version": "4.0.5", "rpc-version": 18, "download-dir": "/downloads",
		"peer-port": 51413, "dht-enabled": true,
	}
	ts := newRPCServer(t, "session-get", result)
	defer ts.Close()

	c := transmission.New(ts.URL)
	s, err := c.GetSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if s.Version != "4.0.5" {
		t.Errorf("Version = %q, want %q", s.Version, "4.0.5")
	}
	if s.DownloadDir != "/downloads" {
		t.Errorf("DownloadDir = %q, want %q", s.DownloadDir, "/downloads")
	}
}

func TestGetSessionStats(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"activeTorrentCount": 5, "downloadSpeed": 50000, "uploadSpeed": 10000,
		"torrentCount": 100,
	}
	ts := newRPCServer(t, "session-stats", result)
	defer ts.Close()

	c := transmission.New(ts.URL)
	stats, err := c.GetSessionStats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if stats.ActiveTorrentCount != 5 {
		t.Errorf("ActiveTorrentCount = %d, want 5", stats.ActiveTorrentCount)
	}
	if stats.DownloadSpeed != 50000 {
		t.Errorf("DownloadSpeed = %d, want 50000", stats.DownloadSpeed)
	}
}

func TestGetFreeSpace(t *testing.T) {
	t.Parallel()

	result := map[string]any{"path": "/downloads", "size-bytes": 500000000000}
	ts := newRPCServer(t, "free-space", result)
	defer ts.Close()

	c := transmission.New(ts.URL)
	fs, err := c.GetFreeSpace(context.Background(), "/downloads")
	if err != nil {
		t.Fatal(err)
	}
	if fs.SizeBytes != 500000000000 {
		t.Errorf("SizeBytes = %d, want 500000000000", fs.SizeBytes)
	}
}

func TestTestPort(t *testing.T) {
	t.Parallel()

	result := map[string]any{"port-is-open": true}
	ts := newRPCServer(t, "port-test", result)
	defer ts.Close()

	c := transmission.New(ts.URL)
	open, err := c.TestPort(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !open {
		t.Error("port-is-open = false, want true")
	}
}

func TestSessionIDNegotiation(t *testing.T) {
	t.Parallel()

	attempt := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		if attempt == 1 {
			w.Header().Set("X-Transmission-Session-Id", "new-session-id")
			w.WriteHeader(http.StatusConflict)
			return
		}
		if got := r.Header.Get("X-Transmission-Session-Id"); got != "new-session-id" {
			t.Errorf("session ID = %q, want %q", got, "new-session-id")
		}
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{"result": "success", "arguments": map[string]any{"torrents": []any{}}}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := transmission.New(ts.URL)
	torrents, err := c.GetTorrents(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(torrents) != 0 {
		t.Errorf("len = %d, want 0", len(torrents))
	}
	if attempt != 2 {
		t.Errorf("attempts = %d, want 2", attempt)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{"result": "no such torrent", "arguments": map[string]any{}}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := transmission.New(ts.URL)
	_, err := c.GetTorrents(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *transmission.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Result != "no such torrent" {
		t.Errorf("Result = %q, want %q", apiErr.Result, "no such torrent")
	}
}

func TestHTTPError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("server error"))
	}))
	defer ts.Close()

	c := transmission.New(ts.URL)
	_, err := c.GetTorrents(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var httpErr *transmission.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected HTTPError, got %T", err)
	}
	if httpErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want %d", httpErr.StatusCode, http.StatusInternalServerError)
	}
}

func TestSetTorrentLabels(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "torrent-set", nil)
	defer ts.Close()

	c := transmission.New(ts.URL)
	if err := c.SetTorrentLabels(context.Background(), []int{1}, []string{"movies", "4k"}); err != nil {
		t.Fatal(err)
	}
}
