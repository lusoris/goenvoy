package sabnzbd_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/downloadclient/sabnzbd"
)

func newServer(t *testing.T, wantMode string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got := r.URL.Query().Get("mode"); got != wantMode {
			t.Errorf("mode = %q, want %q", got, wantMode)
		}
		if got := r.URL.Query().Get("output"); got != "json" {
			t.Errorf("output = %q, want json", got)
		}
		if got := r.URL.Query().Get("apikey"); got != "test-key" {
			t.Errorf("apikey = %q, want test-key", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
}

func TestGetQueue(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"queue": map[string]any{
			"status": "Downloading", "speed": "5.0 M", "paused": false,
			"noofslots": 2,
			"slots": []map[string]any{
				{"nzo_id": "SABnzbd_nzo_abc", "filename": "Ubuntu.nzb", "status": "Downloading", "percentage": "75"},
				{"nzo_id": "SABnzbd_nzo_def", "filename": "Fedora.nzb", "status": "Queued", "percentage": "0"},
			},
		},
	}
	ts := newServer(t, "queue", result)
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "test-key")
	queue, err := c.GetQueue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if queue.Status != "Downloading" {
		t.Errorf("Status = %q, want %q", queue.Status, "Downloading")
	}
	if len(queue.Slots) != 2 {
		t.Fatalf("len(Slots) = %d, want 2", len(queue.Slots))
	}
	if queue.Slots[0].Filename != "Ubuntu.nzb" {
		t.Errorf("Filename = %q, want %q", queue.Slots[0].Filename, "Ubuntu.nzb")
	}
}

func TestAddURL(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"status":  true,
		"nzo_ids": []string{"SABnzbd_nzo_xyz"},
	}
	ts := newServer(t, "addurl", result)
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "test-key")
	ids, err := c.AddURL(context.Background(), "https://example.com/file.nzb", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != "SABnzbd_nzo_xyz" {
		t.Errorf("IDs = %v, want [SABnzbd_nzo_xyz]", ids)
	}
}

func TestPause(t *testing.T) {
	t.Parallel()

	ts := newServer(t, "pause", map[string]any{"status": "ok"})
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "test-key")
	if err := c.Pause(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestResume(t *testing.T) {
	t.Parallel()

	ts := newServer(t, "resume", map[string]any{"status": "ok"})
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "test-key")
	if err := c.Resume(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestGetHistory(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"history": map[string]any{
			"total_size": "10.5 GB", "noofslots": 1,
			"slots": []map[string]any{
				{
					"nzo_id": "SABnzbd_nzo_hist1", "name": "Completed Item",
					"status": "Completed", "size": "4.2 GB",
					"category": "movies", "storage": "/downloads/completed",
				},
			},
		},
	}
	ts := newServer(t, "history", result)
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "test-key")
	hist, err := c.GetHistory(context.Background(), 0, 50)
	if err != nil {
		t.Fatal(err)
	}
	if hist.TotalSize != "10.5 GB" {
		t.Errorf("TotalSize = %q, want %q", hist.TotalSize, "10.5 GB")
	}
	if len(hist.Slots) != 1 {
		t.Fatalf("len(Slots) = %d, want 1", len(hist.Slots))
	}
	if hist.Slots[0].Name != "Completed Item" {
		t.Errorf("Name = %q, want %q", hist.Slots[0].Name, "Completed Item")
	}
}

func TestGetVersion(t *testing.T) {
	t.Parallel()

	ts := newServer(t, "version", map[string]any{"version": "4.2.1"})
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "test-key")
	v, err := c.GetVersion(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if v != "4.2.1" {
		t.Errorf("version = %q, want %q", v, "4.2.1")
	}
}

func TestGetServerStats(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"total": 107374182400, "day": 5368709120, "week": 26843545600, "month": 85899345920,
	}
	ts := newServer(t, "server_stats", result)
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "test-key")
	stats, err := c.GetServerStats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if stats.Total != 107374182400 {
		t.Errorf("Total = %d, want %d", stats.Total, 107374182400)
	}
}

func TestGetWarnings(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"warnings": []string{"disk space low", "connection timeout"},
	}
	ts := newServer(t, "warnings", result)
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "test-key")
	warnings, err := c.GetWarnings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 2 {
		t.Fatalf("len = %d, want 2", len(warnings))
	}
	if warnings[0] != "disk space low" {
		t.Errorf("warning[0] = %q, want %q", warnings[0], "disk space low")
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"status": false,
		"error":  "API Key Required",
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	}))
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "bad-key")
	_, err := c.GetQueue(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *sabnzbd.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
	if apiErr.Message != "API Key Required" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "API Key Required")
	}
}

func TestDeleteItem(t *testing.T) {
	t.Parallel()

	ts := newServer(t, "queue", map[string]any{"status": true})
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "test-key")
	if err := c.DeleteItem(context.Background(), "SABnzbd_nzo_abc"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteHistory(t *testing.T) {
	t.Parallel()

	ts := newServer(t, "history", map[string]any{"status": true})
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "test-key")
	if err := c.DeleteHistory(context.Background(), "SABnzbd_nzo_abc"); err != nil {
		t.Fatal(err)
	}
}

func TestSetCategory(t *testing.T) {
	t.Parallel()

	ts := newServer(t, "change_cat", map[string]any{"status": true})
	defer ts.Close()

	c := sabnzbd.New(ts.URL, "test-key")
	if err := c.SetCategory(context.Background(), "SABnzbd_nzo_abc", "tv"); err != nil {
		t.Fatal(err)
	}
}
