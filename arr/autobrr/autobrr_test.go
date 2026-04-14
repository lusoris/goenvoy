package autobrr_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/arr/autobrr"
)

func newTestServer(t *testing.T, wantPath, wantMethod, wantKey string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if r.Method != wantMethod {
			t.Errorf("method = %q, want %q", r.Method, wantMethod)
		}
		if got := r.Header.Get("X-Api-Token"); got != wantKey {
			t.Errorf("X-API-Token = %q, want %q", got, wantKey)
		}
		w.Header().Set("Content-Type", "application/json")
		if response != nil {
			_ = json.NewEncoder(w).Encode(response)
		}
	}))
}

func TestGetFilters(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/filters", http.MethodGet, "test-key", []map[string]any{
		{"id": 1, "name": "Movies", "enabled": true},
		{"id": 2, "name": "TV", "enabled": false},
	})
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	filters, err := c.GetFilters(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(filters) != 2 {
		t.Fatalf("len(filters) = %d, want 2", len(filters))
	}
	if filters[0].Name != "Movies" {
		t.Errorf("filters[0].Name = %q, want Movies", filters[0].Name)
	}
}

func TestSetFilterEnabled(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/filters/1/enabled", http.MethodPut, "test-key", nil)
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	if err := c.SetFilterEnabled(context.Background(), 1, true); err != nil {
		t.Fatal(err)
	}
}

func TestGetIndexers(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/indexer", http.MethodGet, "test-key", []map[string]any{
		{"id": 1, "name": "Indexer1", "enabled": true},
	})
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	indexers, err := c.GetIndexers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(indexers) != 1 {
		t.Fatalf("len(indexers) = %d, want 1", len(indexers))
	}
}

func TestGetIRCNetworks(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/irc", http.MethodGet, "test-key", []map[string]any{
		{"id": 1, "name": "irc.example.com", "healthy": true},
	})
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	networks, err := c.GetIRCNetworks(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(networks) != 1 {
		t.Fatalf("len(networks) = %d, want 1", len(networks))
	}
	if !networks[0].Healthy {
		t.Error("networks[0].Healthy = false, want true")
	}
}

func TestGetFeeds(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/feeds", http.MethodGet, "test-key", []map[string]any{
		{"id": 1, "name": "Feed1", "enabled": true},
	})
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	feeds, err := c.GetFeeds(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(feeds) != 1 {
		t.Fatalf("len(feeds) = %d, want 1", len(feeds))
	}
}

func TestGetDownloadClients(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/download_clients", http.MethodGet, "test-key", []map[string]any{
		{"id": 1, "name": "qBittorrent", "type": "QBITTORRENT", "enabled": true},
	})
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	clients, err := c.GetDownloadClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(clients) != 1 {
		t.Fatalf("len(clients) = %d, want 1", len(clients))
	}
	if clients[0].Type != "QBITTORRENT" {
		t.Errorf("clients[0].Type = %q, want QBITTORRENT", clients[0].Type)
	}
}

func TestGetNotifications(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/notification", http.MethodGet, "test-key", []map[string]any{
		{"id": 1, "name": "Discord", "type": "DISCORD", "enabled": true},
	})
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	notifs, err := c.GetNotifications(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(notifs) != 1 {
		t.Fatalf("len(notifs) = %d, want 1", len(notifs))
	}
}

func TestGetConfig(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/config", http.MethodGet, "test-key", map[string]any{
		"log_level": "INFO",
		"version":   "1.30.0",
	})
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	cfg, err := c.GetConfig(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if cfg.LogLevel != "INFO" {
		t.Errorf("LogLevel = %q, want INFO", cfg.LogLevel)
	}
}

func TestLiveness(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/healthz/liveness", http.MethodGet, "test-key", nil)
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	if err := c.Liveness(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestReadiness(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/healthz/readiness", http.MethodGet, "test-key", nil)
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	if err := c.Readiness(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized"))
	}))
	defer ts.Close()

	c := autobrr.New(ts.URL, "bad-key")
	_, err := c.GetFilters(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *autobrr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	called := false
	custom := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			called = true
			return http.DefaultTransport.RoundTrip(r)
		}),
	}

	ts := newTestServer(t, "/api/filters", http.MethodGet, "k", []any{})
	defer ts.Close()

	c := autobrr.New(ts.URL, "k", autobrr.WithHTTPClient(custom))
	_, _ = c.GetFilters(context.Background())
	if !called {
		t.Error("custom HTTP client was not used")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestGetAPIKeys(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/keys", http.MethodGet, "test-key", []map[string]any{{"key": "abc"}})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	keys, err := c.GetAPIKeys(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 {
		t.Fatalf("len = %d, want 1", len(keys))
	}
}

func TestCreateAPIKey(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/keys", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.CreateAPIKey(context.Background(), autobrr.APIKey{Name: "test"}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteAPIKey(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/keys/abc123", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteAPIKey(context.Background(), "abc123"); err != nil {
		t.Fatal(err)
	}
}

func TestGetFilter(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/filters/1", http.MethodGet, "test-key", map[string]any{"id": 1, "name": "Movies"})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	f, err := c.GetFilter(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if f.Name != "Movies" {
		t.Errorf("name = %q", f.Name)
	}
}

func TestCreateFilter(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/filters", http.MethodPost, "test-key", map[string]any{"id": 1, "name": "New"})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	f, err := c.CreateFilter(context.Background(), &autobrr.Filter{Name: "New"})
	if err != nil {
		t.Fatal(err)
	}
	if f.Name != "New" {
		t.Errorf("name = %q", f.Name)
	}
}

func TestUpdateFilter(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/filters/1", http.MethodPut, "test-key", map[string]any{"id": 1, "name": "Updated"})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	f, err := c.UpdateFilter(context.Background(), &autobrr.Filter{ID: 1, Name: "Updated"})
	if err != nil {
		t.Fatal(err)
	}
	if f.Name != "Updated" {
		t.Errorf("name = %q", f.Name)
	}
}

func TestDuplicateFilter(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/filters/1/duplicate", http.MethodGet, "test-key", map[string]any{"id": 2, "name": "Copy"})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	f, err := c.DuplicateFilter(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if f.ID != 2 {
		t.Errorf("id = %d", f.ID)
	}
}

func TestDeleteFilter(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/filters/1", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteFilter(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestGetFilterNotifications(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/filters/1/notifications", http.MethodGet, "test-key", []map[string]any{{"id": 1}})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	n, err := c.GetFilterNotifications(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(n) != 1 {
		t.Fatalf("len = %d", len(n))
	}
}

func TestUpdateFilterNotifications(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/filters/1/notifications", http.MethodPut, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.UpdateFilterNotifications(context.Background(), 1, []autobrr.FilterNotification{{ID: 1}}); err != nil {
		t.Fatal(err)
	}
}

func TestGetIndexerSchema(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/indexer/schema", http.MethodGet, "test-key", []map[string]any{{"name": "test"}})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	s, err := c.GetIndexerSchema(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 1 {
		t.Fatalf("len = %d", len(s))
	}
}

func TestGetIndexerOptions(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/indexer/options", http.MethodGet, "test-key", []map[string]any{{"id": 1}})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	o, err := c.GetIndexerOptions(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(o) != 1 {
		t.Fatalf("len = %d", len(o))
	}
}

func TestCreateIndexer(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/indexer", http.MethodPost, "test-key", map[string]any{"id": 1, "name": "New"})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	idx, err := c.CreateIndexer(context.Background(), autobrr.Indexer{Name: "New"})
	if err != nil {
		t.Fatal(err)
	}
	if idx.ID != 1 {
		t.Errorf("id = %d", idx.ID)
	}
}

func TestUpdateIndexer(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/indexer/1", http.MethodPut, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.UpdateIndexer(context.Background(), autobrr.Indexer{ID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteIndexer(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/indexer/1", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteIndexer(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestSetIndexerEnabled(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/indexer/1/enabled", http.MethodPatch, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.SetIndexerEnabled(context.Background(), 1, true); err != nil {
		t.Fatal(err)
	}
}

func TestTestIndexerAPI(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/indexer/1/api/test", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.TestIndexerAPI(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestRestartIRCNetwork(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/irc/network/1/restart", http.MethodGet, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.RestartIRCNetwork(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestCreateIRCNetwork(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/irc", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.CreateIRCNetwork(context.Background(), &autobrr.IRCNetwork{Name: "test"}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateIRCNetwork(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/irc/network/1", http.MethodPut, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.UpdateIRCNetwork(context.Background(), &autobrr.IRCNetwork{ID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteIRCNetwork(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/irc/network/1", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteIRCNetwork(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestSendIRCCommand(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/irc/network/1/cmd", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.SendIRCCommand(context.Background(), autobrr.SendIRCCmdRequest{NetworkID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestReprocessAnnounce(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/irc/network/1/channel/#test/announce/process", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.ReprocessAnnounce(context.Background(), 1, "#test", "hello"); err != nil {
		t.Fatal(err)
	}
}

func TestSetFeedEnabled(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/feeds/1/enabled", http.MethodPatch, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.SetFeedEnabled(context.Background(), 1, true); err != nil {
		t.Fatal(err)
	}
}

func TestCreateFeed(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/feeds", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.CreateFeed(context.Background(), autobrr.Feed{Name: "test"}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateFeed(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/feeds/1", http.MethodPut, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.UpdateFeed(context.Background(), autobrr.Feed{ID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteFeed(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/feeds/1", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteFeed(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteFeedCache(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/feeds/1/cache", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteFeedCache(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestForceRunFeed(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/feeds/1/forcerun", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.ForceRunFeed(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestTestFeed(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/feeds/test", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.TestFeed(context.Background(), autobrr.Feed{}); err != nil {
		t.Fatal(err)
	}
}

func TestGetFeedCaps(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/feeds/1/caps", http.MethodGet, "test-key", map[string]any{"caps": true})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	caps, err := c.GetFeedCaps(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(caps) == 0 {
		t.Error("empty caps")
	}
}

func TestCreateDownloadClient(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/download_clients", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.CreateDownloadClient(context.Background(), &autobrr.DownloadClient{Name: "test"}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateDownloadClient(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/download_clients", http.MethodPut, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.UpdateDownloadClient(context.Background(), &autobrr.DownloadClient{ID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDownloadClient(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/download_clients/1", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteDownloadClient(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestTestDownloadClient(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/download_clients/test", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.TestDownloadClient(context.Background(), &autobrr.DownloadClient{}); err != nil {
		t.Fatal(err)
	}
}

func TestGetDownloadClientArrTags(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/download_clients/1/arr/tags", http.MethodGet, "test-key", []map[string]any{{"id": 1, "label": "tag1"}})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	tags, err := c.GetDownloadClientArrTags(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 1 {
		t.Fatalf("len = %d", len(tags))
	}
}

func TestCreateNotification(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/notification", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.CreateNotification(context.Background(), &autobrr.Notification{Name: "test"}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateNotification(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/notification/1", http.MethodPut, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.UpdateNotification(context.Background(), &autobrr.Notification{ID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteNotification(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/notification/1", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteNotification(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestTestNotification(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/notification/test", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.TestNotification(context.Background(), &autobrr.Notification{}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateConfig(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/config", http.MethodPatch, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.UpdateConfig(context.Background(), autobrr.ConfigUpdate{}); err != nil {
		t.Fatal(err)
	}
}

func TestCreateAction(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/actions", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.CreateAction(context.Background(), &autobrr.Action{Name: "test"}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateAction(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/actions/1", http.MethodPut, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.UpdateAction(context.Background(), &autobrr.Action{ID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteAction(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/actions/1", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteAction(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestToggleActionEnabled(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/actions/1/toggleEnabled", http.MethodPatch, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.ToggleActionEnabled(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestGetLists(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/lists", http.MethodGet, "test-key", []map[string]any{{"id": 1, "name": "test"}})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	lists, err := c.GetLists(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(lists) != 1 {
		t.Fatalf("len = %d", len(lists))
	}
}

func TestGetList(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/lists/1", http.MethodGet, "test-key", map[string]any{"id": 1, "name": "test"})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	l, err := c.GetList(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if l.Name != "test" {
		t.Errorf("name = %q", l.Name)
	}
}

func TestCreateList(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/lists", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.CreateList(context.Background(), &autobrr.List{Name: "test"}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateList(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/lists/1", http.MethodPut, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.UpdateList(context.Background(), &autobrr.List{ID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteList(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/lists/1", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteList(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestRefreshList(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/lists/1/refresh", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.RefreshList(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestRefreshAllLists(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/lists/refresh", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.RefreshAllLists(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestTestList(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/lists/test", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.TestList(context.Background(), &autobrr.List{}); err != nil {
		t.Fatal(err)
	}
}

func TestGetProxies(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/proxy", http.MethodGet, "test-key", []map[string]any{{"id": 1}})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	p, err := c.GetProxies(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 1 {
		t.Fatalf("len = %d", len(p))
	}
}

func TestGetProxy(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/proxy/1", http.MethodGet, "test-key", map[string]any{"id": 1, "name": "test"})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	p, err := c.GetProxy(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != 1 {
		t.Errorf("id = %d", p.ID)
	}
}

func TestCreateProxy(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/proxy", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.CreateProxy(context.Background(), &autobrr.Proxy{Name: "test"}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateProxy(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/proxy/1", http.MethodPut, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.UpdateProxy(context.Background(), &autobrr.Proxy{ID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteProxy(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/proxy/1", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteProxy(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestTestProxy(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/proxy/test", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.TestProxy(context.Background(), &autobrr.Proxy{}); err != nil {
		t.Fatal(err)
	}
}

func TestGetReleases(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release", http.MethodGet, "test-key", map[string]any{"data": []any{}, "count": 0})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	r, err := c.GetReleases(context.Background(), 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("nil response")
	}
}

func TestGetRecentReleases(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/recent", http.MethodGet, "test-key", map[string]any{"data": []any{}, "count": 0})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	r, err := c.GetRecentReleases(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("nil response")
	}
}

func TestGetReleaseStats(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/stats", http.MethodGet, "test-key", map[string]any{"totalCount": 5})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	s, err := c.GetReleaseStats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("nil stats")
	}
}

func TestGetReleaseIndexers(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/indexers", http.MethodGet, "test-key", []string{"idx1", "idx2"})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	idx, err := c.GetReleaseIndexers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(idx) != 2 {
		t.Fatalf("len = %d", len(idx))
	}
}

func TestDeleteReleases(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteReleases(context.Background(), autobrr.ReleaseDeleteParams{OlderThan: 30}); err != nil {
		t.Fatal(err)
	}
}

func TestReplayReleaseAction(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/1/actions/2/retry", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.ReplayReleaseAction(context.Background(), 1, 2); err != nil {
		t.Fatal(err)
	}
}

func TestGetReleaseCleanupJobs(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/cleanup-jobs", http.MethodGet, "test-key", []map[string]any{{"id": 1}})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	jobs, err := c.GetReleaseCleanupJobs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 {
		t.Fatalf("len = %d", len(jobs))
	}
}

func TestGetReleaseCleanupJob(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/cleanup-jobs/1", http.MethodGet, "test-key", map[string]any{"id": 1})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	j, err := c.GetReleaseCleanupJob(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if j.ID != 1 {
		t.Errorf("id = %d", j.ID)
	}
}

func TestCreateReleaseCleanupJob(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/cleanup-jobs", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.CreateReleaseCleanupJob(context.Background(), &autobrr.ReleaseCleanupJob{}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateReleaseCleanupJob(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/cleanup-jobs/1", http.MethodPut, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.UpdateReleaseCleanupJob(context.Background(), &autobrr.ReleaseCleanupJob{ID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteReleaseCleanupJob(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/cleanup-jobs/1", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteReleaseCleanupJob(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestToggleReleaseCleanupJobEnabled(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/cleanup-jobs/1/enabled", http.MethodPatch, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.ToggleReleaseCleanupJobEnabled(context.Background(), 1, true); err != nil {
		t.Fatal(err)
	}
}

func TestForceRunReleaseCleanupJob(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/cleanup-jobs/1/run", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.ForceRunReleaseCleanupJob(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestGetReleaseDuplicateProfiles(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/profiles/duplicate", http.MethodGet, "test-key", []map[string]any{{"id": 1}})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	p, err := c.GetReleaseDuplicateProfiles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(p) != 1 {
		t.Fatalf("len = %d", len(p))
	}
}

func TestCreateReleaseDuplicateProfile(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/profiles/duplicate", http.MethodPost, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.CreateReleaseDuplicateProfile(context.Background(), autobrr.ReleaseProfileDuplicate{}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteReleaseDuplicateProfile(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/release/profiles/duplicate/1", http.MethodDelete, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.DeleteReleaseDuplicateProfile(context.Background(), 1); err != nil {
		t.Fatal(err)
	}
}

func TestGetLogFiles(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/logs/files", http.MethodGet, "test-key", map[string]any{"files": []map[string]any{{"name": "log1.txt"}}, "count": 1})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	lf, err := c.GetLogFiles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if lf == nil {
		t.Fatal("nil")
	}
}

func TestGetLogFile(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("log content"))
	}))
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	data, err := c.GetLogFile(context.Background(), "test.log")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("empty data")
	}
}

func TestCheckForUpdates(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/updates/check", http.MethodGet, "test-key", nil)
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	if err := c.CheckForUpdates(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestGetLatestRelease(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/updates/latest", http.MethodGet, "test-key", map[string]any{"tag_name": "v1.0.0"})
	defer ts.Close()
	c := autobrr.New(ts.URL, "test-key")
	r, err := c.GetLatestRelease(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("nil")
	}
}
