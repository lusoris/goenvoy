package autobrr_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/arr/autobrr"
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
		if got := r.Header.Get("X-API-Token"); got != wantKey {
			t.Errorf("X-API-Token = %q, want %q", got, wantKey)
		}
		w.Header().Set("Content-Type", "application/json")
		if response != nil {
			_ = json.NewEncoder(w).Encode(response)
		}
	}))
}

func TestGetFilters(t *testing.T) {
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
	ts := newTestServer(t, "/api/filters/1/enabled", http.MethodPut, "test-key", nil)
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	if err := c.SetFilterEnabled(context.Background(), 1, true); err != nil {
		t.Fatal(err)
	}
}

func TestGetIndexers(t *testing.T) {
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
	ts := newTestServer(t, "/api/healthz/liveness", http.MethodGet, "test-key", nil)
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	if err := c.Liveness(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestReadiness(t *testing.T) {
	ts := newTestServer(t, "/api/healthz/readiness", http.MethodGet, "test-key", nil)
	defer ts.Close()

	c := autobrr.New(ts.URL, "test-key")
	if err := c.Readiness(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestAPIError(t *testing.T) {
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
