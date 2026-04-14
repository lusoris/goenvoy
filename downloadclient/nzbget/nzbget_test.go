package nzbget_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/downloadclient/nzbget"
)

type rpcRequest struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      int             `json:"id"`
}

func newRPCServer(t *testing.T, wantMethod string, result any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("HTTP method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/jsonrpc" {
			t.Errorf("path = %q, want /jsonrpc", r.URL.Path)
		}
		u, p, ok := r.BasicAuth()
		if !ok || u != "nzbget" || p != "pass123" {
			t.Errorf("auth = %q/%q/%v, want nzbget/pass123/true", u, p, ok)
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
			"jsonrpc": "2.0",
			"result":  json.RawMessage(resultJSON),
			"id":      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func TestListGroups(t *testing.T) {
	t.Parallel()

	result := []map[string]any{
		{"NZBID": 1, "NZBName": "Ubuntu.nzb", "Status": "DOWNLOADING", "Category": "linux"},
		{"NZBID": 2, "NZBName": "Fedora.nzb", "Status": "QUEUED", "Category": "linux"},
	}
	ts := newRPCServer(t, "listgroups", result)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	groups, err := c.ListGroups(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatalf("len = %d, want 2", len(groups))
	}
	if groups[0].Name != "Ubuntu.nzb" {
		t.Errorf("Name = %q, want %q", groups[0].Name, "Ubuntu.nzb")
	}
	if groups[1].Status != "QUEUED" {
		t.Errorf("Status = %q, want %q", groups[1].Status, "QUEUED")
	}
}

func TestAppend(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "append", 42)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	id, err := c.Append(context.Background(), "file.nzb", "https://example.com/file.nzb", "movies", 0)
	if err != nil {
		t.Fatal(err)
	}
	if id != 42 {
		t.Errorf("id = %d, want 42", id)
	}
}

func TestEditQueue(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "editqueue", true)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	ok, err := c.EditQueue(context.Background(), "GroupPause", "", []int{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("ok = false, want true")
	}
}

func TestPauseDownload(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "pausedownload", true)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	if err := c.PauseDownload(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestResumeDownload(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "resumedownload", true)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	if err := c.ResumeDownload(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestGetStatus(t *testing.T) {
	t.Parallel()

	result := map[string]any{
		"ServerPaused": false, "DownloadRate": 50000000,
		"FreeDiskSpaceLo": 500000000, "ThreadCount": 4,
	}
	ts := newRPCServer(t, "status", result)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	status, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if status.ServerPaused {
		t.Error("ServerPaused = true, want false")
	}
	if status.DownloadRate != 50000000 {
		t.Errorf("DownloadRate = %d, want 50000000", status.DownloadRate)
	}
}

func TestGetVersion(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "version", "21.1")
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	v, err := c.GetVersion(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if v != "21.1" {
		t.Errorf("version = %q, want %q", v, "21.1")
	}
}

func TestGetHistory(t *testing.T) {
	t.Parallel()

	result := []map[string]any{
		{"NZBID": 10, "Name": "Completed", "Status": "SUCCESS/ALL", "Category": "tv"},
	}
	ts := newRPCServer(t, "history", result)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	items, err := c.GetHistory(context.Background(), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("len = %d, want 1", len(items))
	}
	if items[0].Name != "Completed" {
		t.Errorf("Name = %q, want %q", items[0].Name, "Completed")
	}
}

func TestGetConfig(t *testing.T) {
	t.Parallel()

	result := []map[string]any{
		{"Name": "MainDir", "Value": "/downloads"},
		{"Name": "TempDir", "Value": "/tmp/nzbget"},
	}
	ts := newRPCServer(t, "config", result)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	entries, err := c.GetConfig(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("len = %d, want 2", len(entries))
	}
	if entries[0].Name != "MainDir" {
		t.Errorf("Name = %q, want %q", entries[0].Name, "MainDir")
	}
}

func TestGetLog(t *testing.T) {
	t.Parallel()

	result := []map[string]any{
		{"ID": 1, "Kind": "INFO", "Text": "Starting download"},
		{"ID": 2, "Kind": "WARNING", "Text": "Speed limited"},
	}
	ts := newRPCServer(t, "log", result)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	entries, err := c.GetLog(context.Background(), 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("len = %d, want 2", len(entries))
	}
	if entries[1].Kind != "WARNING" {
		t.Errorf("Kind = %q, want %q", entries[1].Kind, "WARNING")
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
			"jsonrpc": "2.0",
			"error":   map[string]any{"code": -32601, "message": "Method not found"},
			"id":      req.ID,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	_, err := c.ListGroups(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *nzbget.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Code != -32601 {
		t.Errorf("Code = %d, want -32601", apiErr.Code)
	}
}

func TestSetDownloadRate(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "rate", true)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	if err := c.SetDownloadRate(context.Background(), 5000); err != nil {
		t.Fatal(err)
	}
}

func TestScanNZBDir(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "scan", true)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	if err := c.ScanNZBDir(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestSetCategory(t *testing.T) {
	t.Parallel()

	ts := newRPCServer(t, "editqueue", true)
	defer ts.Close()

	c := nzbget.New(ts.URL, "nzbget", "pass123")
	ok, err := c.SetCategory(context.Background(), 1, "movies")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("ok = false, want true")
	}
}
