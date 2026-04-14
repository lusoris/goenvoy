package tdarr_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/mediaserver/tdarr/v2"
)

func newGetServer(t *testing.T, wantPath string, resp any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func newPostServer(t *testing.T, wantPath string, resp any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func TestAPIKeyHeader(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Api-Key"); got != "secret-key" {
			t.Errorf("x-api-key = %q, want secret-key", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tdarr.Status{Status: "ok"})
	}))
	defer ts.Close()

	c := tdarr.New(ts.URL, tdarr.WithAPIKey("secret-key"))
	_, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetStatus(t *testing.T) {
	t.Parallel()

	ts := newGetServer(t, "/api/v2/status", tdarr.Status{
		Status: "good", Os: "linux", Version: "2.17.01",
	})
	defer ts.Close()

	c := tdarr.New(ts.URL)
	s, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if s.Version != "2.17.01" {
		t.Errorf("Version = %q, want 2.17.01", s.Version)
	}
}

func TestGetNodes(t *testing.T) {
	t.Parallel()

	nodes := map[string]tdarr.Node{
		"node1": {ID: "node1", Name: "Main", Port: 8266},
	}
	ts := newGetServer(t, "/api/v2/get-nodes", nodes)
	defer ts.Close()

	c := tdarr.New(ts.URL)
	out, err := c.GetNodes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if n, ok := out["node1"]; !ok {
		t.Fatal("missing node1")
	} else if n.Name != "Main" {
		t.Errorf("Name = %q, want Main", n.Name)
	}
}

func TestSearchDB(t *testing.T) {
	t.Parallel()

	files := []tdarr.DBFile{
		{ID: "f1", File: "/media/movie.mkv", Codec: "h264"},
	}
	ts := newPostServer(t, "/api/v2/search-db", files)
	defer ts.Close()

	c := tdarr.New(ts.URL)
	out, err := c.SearchDB(context.Background(), "staging", 10, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("len = %d, want 1", len(out))
	}
	if out[0].Codec != "h264" {
		t.Errorf("Codec = %q, want h264", out[0].Codec)
	}
}

func TestCrudDB(t *testing.T) {
	t.Parallel()

	docs := []map[string]any{{"_id": "doc1", "name": "test"}}
	ts := newPostServer(t, "/api/v2/cruddb", docs)
	defer ts.Close()

	c := tdarr.New(ts.URL)
	out, err := c.CrudDB(context.Background(), "staging", "read", "doc1")
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("len = %d, want 1", len(out))
	}
}

func TestGetResStats(t *testing.T) {
	t.Parallel()

	ts := newPostServer(t, "/api/v2/get-res-stats", tdarr.ResStats{
		Pie: map[string]int{"1080p": 50, "720p": 30},
	})
	defer ts.Close()

	c := tdarr.New(ts.URL)
	out, err := c.GetResStats(context.Background(), "lib1")
	if err != nil {
		t.Fatal(err)
	}
	if out.Pie["1080p"] != 50 {
		t.Errorf("Pie[1080p] = %d, want 50", out.Pie["1080p"])
	}
}

func TestGetDBStatuses(t *testing.T) {
	t.Parallel()

	ts := newPostServer(t, "/api/v2/get-db-statuses", tdarr.DBStatuses{
		Table1Count: 100, Table2Count: 50,
	})
	defer ts.Close()

	c := tdarr.New(ts.URL)
	out, err := c.GetDBStatuses(context.Background(), "lib1")
	if err != nil {
		t.Fatal(err)
	}
	if out.Table1Count != 100 {
		t.Errorf("Table1Count = %d, want 100", out.Table1Count)
	}
}

func TestScanFiles(t *testing.T) {
	t.Parallel()

	ts := newPostServer(t, "/api/v2/scan-files", nil)
	defer ts.Close()

	c := tdarr.New(ts.URL)
	if err := c.ScanFiles(context.Background(), "lib1", "/media"); err != nil {
		t.Fatal(err)
	}
}

func TestCancelWorkerItem(t *testing.T) {
	t.Parallel()

	ts := newPostServer(t, "/api/v2/cancel-worker-item", nil)
	defer ts.Close()

	c := tdarr.New(ts.URL)
	if err := c.CancelWorkerItem(context.Background(), "node1", "worker1"); err != nil {
		t.Fatal(err)
	}
}

func TestKillWorker(t *testing.T) {
	t.Parallel()

	ts := newPostServer(t, "/api/v2/kill-worker", nil)
	defer ts.Close()

	c := tdarr.New(ts.URL)
	if err := c.KillWorker(context.Background(), "node1", "worker1"); err != nil {
		t.Fatal(err)
	}
}

func TestGetServerLog(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("log line 1\nlog line 2"))
	}))
	defer ts.Close()

	c := tdarr.New(ts.URL)
	log, err := c.GetServerLog(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if log != "log line 1\nlog line 2" {
		t.Errorf("log = %q, want log line 1\\nlog line 2", log)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized"))
	}))
	defer ts.Close()

	c := tdarr.New(ts.URL, tdarr.WithAPIKey("bad-key"))
	_, err := c.GetStatus(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *tdarr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	ts := newGetServer(t, "/api/v2/status", tdarr.Status{Status: "ok"})
	defer ts.Close()

	custom := &http.Client{}
	c := tdarr.New(ts.URL, tdarr.WithHTTPClient(custom))
	_, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
