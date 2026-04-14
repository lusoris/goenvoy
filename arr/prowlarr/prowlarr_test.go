package prowlarr_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golusoris/goenvoy/arr/prowlarr/v2"
	"github.com/golusoris/goenvoy/arr/v2"
)

func newTestServer(t *testing.T, method, wantPath string, body any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("method = %s, want %s", r.Method, method)
		}
		if r.Header.Get("X-Api-Key") == "" {
			t.Error("missing X-Api-Key header")
		}
		if wantPath != "" && r.URL.RequestURI() != wantPath {
			t.Errorf("path = %q, want %q", r.URL.RequestURI(), wantPath)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}))
}

func newRawTestServer(t *testing.T, method, wantPath, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("method = %s, want %s", r.Method, method)
		}
		if r.Header.Get("X-Api-Key") == "" {
			t.Error("missing X-Api-Key header")
		}
		if wantPath != "" && r.URL.RequestURI() != wantPath {
			t.Errorf("path = %q, want %q", r.URL.RequestURI(), wantPath)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, body)
	}))
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		c, err := prowlarr.New("http://localhost:9696", "test-key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c == nil {
			t.Fatal("expected non-nil client")
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		t.Parallel()
		_, err := prowlarr.New("://bad", "test-key")
		if err == nil {
			t.Fatal("expected error for invalid URL")
		}
	})
}

func TestGetIndexers(t *testing.T) {
	t.Parallel()

	want := []prowlarr.Indexer{
		{ID: 1, Name: "NZBgeek", Enable: true, Protocol: "usenet"},
		{ID: 2, Name: "TorrentLeech", Enable: true, Protocol: "torrent"},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer", want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetIndexers(context.Background())
	if err != nil {
		t.Fatalf("GetIndexers: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Name != "NZBgeek" {
		t.Errorf("Name = %q, want %q", got[0].Name, "NZBgeek")
	}
}

func TestGetIndexer(t *testing.T) {
	t.Parallel()

	want := prowlarr.Indexer{ID: 1, Name: "NZBgeek", Enable: true}

	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer/1", want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetIndexer(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetIndexer: %v", err)
	}
	if got.Name != "NZBgeek" {
		t.Errorf("Name = %q, want %q", got.Name, "NZBgeek")
	}
}

func TestAddIndexer(t *testing.T) {
	t.Parallel()

	want := prowlarr.Indexer{ID: 3, Name: "Jackett", Enable: true}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var body prowlarr.Indexer
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Name != "Jackett" {
			t.Errorf("Name = %q, want %q", body.Name, "Jackett")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.AddIndexer(context.Background(), &prowlarr.Indexer{Name: "Jackett"})
	if err != nil {
		t.Fatalf("AddIndexer: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("ID = %d, want 3", got.ID)
	}
}

func TestDeleteIndexer(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v1/indexer/1", nil)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteIndexer(context.Background(), 1); err != nil {
		t.Fatalf("DeleteIndexer: %v", err)
	}
}

func TestGetIndexerCategories(t *testing.T) {
	t.Parallel()

	want := []prowlarr.IndexerCategory{
		{ID: 2000, Name: "Movies"},
		{ID: 5000, Name: "TV", SubCategories: []prowlarr.IndexerCategory{
			{ID: 5010, Name: "TV/WEB-DL"},
		}},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer/categories", want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetIndexerCategories(context.Background())
	if err != nil {
		t.Fatalf("GetIndexerCategories: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if len(got[1].SubCategories) != 1 {
		t.Errorf("SubCategories len = %d, want 1", len(got[1].SubCategories))
	}
}

func TestGetApplications(t *testing.T) {
	t.Parallel()

	want := []prowlarr.Application{
		{ID: 1, Name: "Sonarr", SyncLevel: "fullSync"},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/applications", want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetApplications(context.Background())
	if err != nil {
		t.Fatalf("GetApplications: %v", err)
	}
	if got[0].Name != "Sonarr" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Sonarr")
	}
}

func TestDeleteApplication(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v1/applications/1", nil)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteApplication(context.Background(), 1); err != nil {
		t.Fatalf("DeleteApplication: %v", err)
	}
}

func TestGetAppProfiles(t *testing.T) {
	t.Parallel()

	want := []prowlarr.AppProfile{
		{ID: 1, Name: "Standard", EnableRss: true, EnableAutomaticSearch: true, EnableInteractiveSearch: true},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/appprofile", want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetAppProfiles(context.Background())
	if err != nil {
		t.Fatalf("GetAppProfiles: %v", err)
	}
	if got[0].Name != "Standard" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Standard")
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()

	want := []prowlarr.Release{
		{ID: 1, Title: "Ubuntu 24.04 LTS", Size: 4500000000, IndexerID: 1, Indexer: "NZBgeek"},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/search?categories=2000&categories=2010&indexerIds=1&limit=25&query=ubuntu&type=search",
		want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.Search(context.Background(), &prowlarr.SearchOptions{
		Query:      "ubuntu",
		Type:       "search",
		IndexerIDs: []int{1},
		Categories: []int{2000, 2010},
		Limit:      25,
	})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Title != "Ubuntu 24.04 LTS" {
		t.Errorf("Title = %q, want %q", got[0].Title, "Ubuntu 24.04 LTS")
	}
}

func TestGrabRelease(t *testing.T) {
	t.Parallel()

	want := prowlarr.Release{ID: 5, Title: "Grabbed", IndexerID: 1}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GrabRelease(context.Background(), &prowlarr.Release{
		GUID:      "abc-123",
		IndexerID: 1,
	})
	if err != nil {
		t.Fatalf("GrabRelease: %v", err)
	}
	if got.ID != 5 {
		t.Errorf("ID = %d, want 5", got.ID)
	}
}

func TestGetDownloadClients(t *testing.T) {
	t.Parallel()

	want := []prowlarr.DownloadClientResource{
		{ID: 1, Name: "qBittorrent", Enable: true, Protocol: "torrent"},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/downloadclient", want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetDownloadClients(context.Background())
	if err != nil {
		t.Fatalf("GetDownloadClients: %v", err)
	}
	if got[0].Name != "qBittorrent" {
		t.Errorf("Name = %q, want %q", got[0].Name, "qBittorrent")
	}
}

func TestSendCommand(t *testing.T) {
	t.Parallel()

	want := arr.CommandResponse{ID: 1, Name: "ApplicationIndexerSync"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var cmd arr.CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if cmd.Name != "ApplicationIndexerSync" {
			t.Errorf("Name = %q, want ApplicationIndexerSync", cmd.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.SendCommand(context.Background(), arr.CommandRequest{Name: "ApplicationIndexerSync"})
	if err != nil {
		t.Fatalf("SendCommand: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetSystemStatus(t *testing.T) {
	t.Parallel()

	want := arr.StatusResponse{AppName: "Prowlarr", Version: "1.25.0"}

	srv := newTestServer(t, http.MethodGet, "/api/v1/system/status", want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetSystemStatus(context.Background())
	if err != nil {
		t.Fatalf("GetSystemStatus: %v", err)
	}
	if got.AppName != "Prowlarr" {
		t.Errorf("AppName = %q, want %q", got.AppName, "Prowlarr")
	}
}

func TestGetTags(t *testing.T) {
	t.Parallel()

	want := []arr.Tag{{ID: 1, Label: "usenet"}, {ID: 2, Label: "torrent"}}

	srv := newTestServer(t, http.MethodGet, "/api/v1/tag", want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetTags(context.Background())
	if err != nil {
		t.Fatalf("GetTags: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}

func TestGetHistory(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[prowlarr.HistoryRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []prowlarr.HistoryRecord{
			{ID: 1, IndexerID: 1, EventType: "releaseGrabbed", Successful: true},
		},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/history?page=1&pageSize=10", want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetHistory(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetHistory: %v", err)
	}
	if got.Records[0].EventType != "releaseGrabbed" {
		t.Errorf("EventType = %q, want %q", got.Records[0].EventType, "releaseGrabbed")
	}
}

func TestGetHistoryByIndexer(t *testing.T) {
	t.Parallel()

	want := []prowlarr.HistoryRecord{
		{ID: 5, IndexerID: 1, EventType: "indexerQuery"},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/history/indexer?indexerId=1&limit=50",
		want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetHistoryByIndexer(context.Background(), 1, 50)
	if err != nil {
		t.Fatalf("GetHistoryByIndexer: %v", err)
	}
	if got[0].EventType != "indexerQuery" {
		t.Errorf("EventType = %q, want %q", got[0].EventType, "indexerQuery")
	}
}

func TestGetIndexerStats(t *testing.T) {
	t.Parallel()

	want := prowlarr.IndexerStats{
		ID: 0,
		Indexers: []prowlarr.IndexerStatistic{
			{IndexerID: 1, IndexerName: "NZBgeek", NumberOfQueries: 100, NumberOfGrabs: 10},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/indexerstats?startDate=2025-01-01&endDate=2025-01-31",
		want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetIndexerStats(context.Background(), "2025-01-01", "2025-01-31")
	if err != nil {
		t.Fatalf("GetIndexerStats: %v", err)
	}
	if got.Indexers[0].NumberOfQueries != 100 {
		t.Errorf("NumberOfQueries = %d, want 100", got.Indexers[0].NumberOfQueries)
	}
}

func TestGetIndexerStatuses(t *testing.T) {
	t.Parallel()

	want := []prowlarr.IndexerStatus{
		{ID: 1, IndexerID: 5, DisabledTill: "2025-01-01T12:00:00Z"},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/indexerstatus", want)
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetIndexerStatuses(context.Background())
	if err != nil {
		t.Fatalf("GetIndexerStatuses: %v", err)
	}
	if got[0].IndexerID != 5 {
		t.Errorf("IndexerID = %d, want 5", got[0].IndexerID)
	}
}

// ---------- Indexer Extended ----------.

func TestGetIndexerSchema(t *testing.T) {
	t.Parallel()
	want := []prowlarr.Indexer{{ID: 0, Name: "Newznab", Protocol: "usenet"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer/schema", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetIndexerSchema(context.Background())
	if err != nil {
		t.Fatalf("GetIndexerSchema: %v", err)
	}
	if got[0].Name != "Newznab" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Newznab")
	}
}

func TestTestIndexer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexer/test", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.TestIndexer(context.Background(), &prowlarr.Indexer{ID: 1}); err != nil {
		t.Fatalf("TestIndexer: %v", err)
	}
}

func TestTestAllIndexers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexer/testall", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.TestAllIndexers(context.Background()); err != nil {
		t.Fatalf("TestAllIndexers: %v", err)
	}
}

func TestBulkUpdateIndexers(t *testing.T) {
	t.Parallel()
	want := []prowlarr.Indexer{{ID: 1, Name: "Updated"}}
	srv := newTestServer(t, http.MethodPut, "/api/v1/indexer/bulk", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.BulkUpdateIndexers(context.Background(), &prowlarr.IndexerBulkResource{IDs: []int{1}})
	if err != nil {
		t.Fatalf("BulkUpdateIndexers: %v", err)
	}
	if got[0].Name != "Updated" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Updated")
	}
}

func TestBulkDeleteIndexers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/indexer/bulk", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.BulkDeleteIndexers(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("BulkDeleteIndexers: %v", err)
	}
}

func TestIndexerAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexer/action/testAction", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.IndexerAction(context.Background(), "testAction", &prowlarr.Indexer{ID: 1}); err != nil {
		t.Fatalf("IndexerAction: %v", err)
	}
}

// ---------- Indexer Proxies ----------.

func TestGetIndexerProxies(t *testing.T) {
	t.Parallel()
	want := []prowlarr.IndexerProxyResource{{ID: 1, Name: "FlareSolverr"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexerproxy", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetIndexerProxies(context.Background())
	if err != nil {
		t.Fatalf("GetIndexerProxies: %v", err)
	}
	if got[0].Name != "FlareSolverr" {
		t.Errorf("Name = %q, want %q", got[0].Name, "FlareSolverr")
	}
}

func TestGetIndexerProxy(t *testing.T) {
	t.Parallel()
	want := prowlarr.IndexerProxyResource{ID: 1, Name: "FlareSolverr"}
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexerproxy/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetIndexerProxy(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetIndexerProxy: %v", err)
	}
	if got.Name != "FlareSolverr" {
		t.Errorf("Name = %q, want %q", got.Name, "FlareSolverr")
	}
}

func TestCreateIndexerProxy(t *testing.T) {
	t.Parallel()
	want := prowlarr.IndexerProxyResource{ID: 1, Name: "FlareSolverr"}
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexerproxy", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.CreateIndexerProxy(context.Background(), &prowlarr.IndexerProxyResource{Name: "FlareSolverr"})
	if err != nil {
		t.Fatalf("CreateIndexerProxy: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateIndexerProxy(t *testing.T) {
	t.Parallel()
	want := prowlarr.IndexerProxyResource{ID: 1, Name: "Updated"}
	srv := newTestServer(t, http.MethodPut, "/api/v1/indexerproxy/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.UpdateIndexerProxy(context.Background(), &prowlarr.IndexerProxyResource{ID: 1, Name: "Updated"})
	if err != nil {
		t.Fatalf("UpdateIndexerProxy: %v", err)
	}
	if got.Name != "Updated" {
		t.Errorf("Name = %q, want %q", got.Name, "Updated")
	}
}

func TestDeleteIndexerProxy(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/indexerproxy/1", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.DeleteIndexerProxy(context.Background(), 1); err != nil {
		t.Fatalf("DeleteIndexerProxy: %v", err)
	}
}

func TestGetIndexerProxySchema(t *testing.T) {
	t.Parallel()
	want := []prowlarr.IndexerProxyResource{{ID: 0, Name: "FlareSolverr"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexerproxy/schema", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetIndexerProxySchema(context.Background())
	if err != nil {
		t.Fatalf("GetIndexerProxySchema: %v", err)
	}
	if got[0].Name != "FlareSolverr" {
		t.Errorf("Name = %q, want %q", got[0].Name, "FlareSolverr")
	}
}

func TestTestIndexerProxy(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexerproxy/test", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.TestIndexerProxy(context.Background(), &prowlarr.IndexerProxyResource{ID: 1}); err != nil {
		t.Fatalf("TestIndexerProxy: %v", err)
	}
}

func TestTestAllIndexerProxies(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexerproxy/testall", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.TestAllIndexerProxies(context.Background()); err != nil {
		t.Fatalf("TestAllIndexerProxies: %v", err)
	}
}

func TestIndexerProxyAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexerproxy/action/testAction", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.IndexerProxyAction(context.Background(), "testAction", &prowlarr.IndexerProxyResource{ID: 1}); err != nil {
		t.Fatalf("IndexerProxyAction: %v", err)
	}
}

// ---------- Applications Extended ----------.

func TestGetApplicationSchema(t *testing.T) {
	t.Parallel()
	want := []prowlarr.Application{{ID: 0, Name: "Sonarr"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/applications/schema", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetApplicationSchema(context.Background())
	if err != nil {
		t.Fatalf("GetApplicationSchema: %v", err)
	}
	if got[0].Name != "Sonarr" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Sonarr")
	}
}

func TestTestApplication(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/applications/test", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.TestApplication(context.Background(), &prowlarr.Application{ID: 1}); err != nil {
		t.Fatalf("TestApplication: %v", err)
	}
}

func TestTestAllApplications(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/applications/testall", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.TestAllApplications(context.Background()); err != nil {
		t.Fatalf("TestAllApplications: %v", err)
	}
}

func TestBulkUpdateApplications(t *testing.T) {
	t.Parallel()
	want := []prowlarr.Application{{ID: 1, Name: "Updated"}}
	srv := newTestServer(t, http.MethodPut, "/api/v1/applications/bulk", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.BulkUpdateApplications(context.Background(), &prowlarr.ApplicationBulkResource{IDs: []int{1}})
	if err != nil {
		t.Fatalf("BulkUpdateApplications: %v", err)
	}
	if got[0].Name != "Updated" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Updated")
	}
}

func TestBulkDeleteApplications(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/applications/bulk", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.BulkDeleteApplications(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("BulkDeleteApplications: %v", err)
	}
}

func TestApplicationAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/applications/action/testAction", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.ApplicationAction(context.Background(), "testAction", &prowlarr.Application{ID: 1}); err != nil {
		t.Fatalf("ApplicationAction: %v", err)
	}
}

// ---------- App Profile Extended ----------.

func TestGetAppProfileSchema(t *testing.T) {
	t.Parallel()
	want := prowlarr.AppProfile{ID: 0, Name: "Standard"}
	srv := newTestServer(t, http.MethodGet, "/api/v1/appprofile/schema", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetAppProfileSchema(context.Background())
	if err != nil {
		t.Fatalf("GetAppProfileSchema: %v", err)
	}
	if got.Name != "Standard" {
		t.Errorf("Name = %q, want %q", got.Name, "Standard")
	}
}

// ---------- Download Clients Extended ----------.

func TestGetDownloadClient(t *testing.T) {
	t.Parallel()
	want := prowlarr.DownloadClientResource{ID: 1, Name: "SABnzbd"}
	srv := newTestServer(t, http.MethodGet, "/api/v1/downloadclient/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetDownloadClient(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetDownloadClient: %v", err)
	}
	if got.Name != "SABnzbd" {
		t.Errorf("Name = %q, want %q", got.Name, "SABnzbd")
	}
}

func TestCreateDownloadClient(t *testing.T) {
	t.Parallel()
	want := prowlarr.DownloadClientResource{ID: 1, Name: "SABnzbd"}
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.CreateDownloadClient(context.Background(), &prowlarr.DownloadClientResource{Name: "SABnzbd"})
	if err != nil {
		t.Fatalf("CreateDownloadClient: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateDownloadClient(t *testing.T) {
	t.Parallel()
	want := prowlarr.DownloadClientResource{ID: 1, Name: "Updated"}
	srv := newTestServer(t, http.MethodPut, "/api/v1/downloadclient/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.UpdateDownloadClient(context.Background(), &prowlarr.DownloadClientResource{ID: 1, Name: "Updated"})
	if err != nil {
		t.Fatalf("UpdateDownloadClient: %v", err)
	}
	if got.Name != "Updated" {
		t.Errorf("Name = %q, want %q", got.Name, "Updated")
	}
}

func TestDeleteDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/downloadclient/1", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.DeleteDownloadClient(context.Background(), 1); err != nil {
		t.Fatalf("DeleteDownloadClient: %v", err)
	}
}

func TestGetDownloadClientSchema(t *testing.T) {
	t.Parallel()
	want := []prowlarr.DownloadClientResource{{ID: 0, Name: "SABnzbd"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/downloadclient/schema", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetDownloadClientSchema(context.Background())
	if err != nil {
		t.Fatalf("GetDownloadClientSchema: %v", err)
	}
	if got[0].Name != "SABnzbd" {
		t.Errorf("Name = %q, want %q", got[0].Name, "SABnzbd")
	}
}

func TestTestDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient/test", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.TestDownloadClient(context.Background(), &prowlarr.DownloadClientResource{ID: 1}); err != nil {
		t.Fatalf("TestDownloadClient: %v", err)
	}
}

func TestTestAllDownloadClients(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient/testall", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.TestAllDownloadClients(context.Background()); err != nil {
		t.Fatalf("TestAllDownloadClients: %v", err)
	}
}

func TestBulkUpdateDownloadClients(t *testing.T) {
	t.Parallel()
	want := []prowlarr.DownloadClientResource{{ID: 1, Name: "Updated"}}
	srv := newTestServer(t, http.MethodPut, "/api/v1/downloadclient/bulk", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.BulkUpdateDownloadClients(context.Background(), &arr.ProviderBulkResource{IDs: []int{1}})
	if err != nil {
		t.Fatalf("BulkUpdateDownloadClients: %v", err)
	}
	if got[0].Name != "Updated" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Updated")
	}
}

func TestBulkDeleteDownloadClients(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/downloadclient/bulk", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.BulkDeleteDownloadClients(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("BulkDeleteDownloadClients: %v", err)
	}
}

func TestDownloadClientAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient/action/testAction", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.DownloadClientAction(context.Background(), "testAction", &prowlarr.DownloadClientResource{ID: 1}); err != nil {
		t.Fatalf("DownloadClientAction: %v", err)
	}
}

// ---------- Notifications ----------.

func TestGetNotifications(t *testing.T) {
	t.Parallel()
	want := []arr.ProviderResource{{ID: 1, Name: "Email"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/notification", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetNotifications(context.Background())
	if err != nil {
		t.Fatalf("GetNotifications: %v", err)
	}
	if got[0].Name != "Email" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Email")
	}
}

func TestGetNotification(t *testing.T) {
	t.Parallel()
	want := arr.ProviderResource{ID: 1, Name: "Email"}
	srv := newTestServer(t, http.MethodGet, "/api/v1/notification/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetNotification(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetNotification: %v", err)
	}
	if got.Name != "Email" {
		t.Errorf("Name = %q, want %q", got.Name, "Email")
	}
}

func TestCreateNotification(t *testing.T) {
	t.Parallel()
	want := arr.ProviderResource{ID: 1, Name: "Email"}
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.CreateNotification(context.Background(), &arr.ProviderResource{Name: "Email"})
	if err != nil {
		t.Fatalf("CreateNotification: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateNotification(t *testing.T) {
	t.Parallel()
	want := arr.ProviderResource{ID: 1, Name: "Updated"}
	srv := newTestServer(t, http.MethodPut, "/api/v1/notification/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.UpdateNotification(context.Background(), &arr.ProviderResource{ID: 1, Name: "Updated"})
	if err != nil {
		t.Fatalf("UpdateNotification: %v", err)
	}
	if got.Name != "Updated" {
		t.Errorf("Name = %q, want %q", got.Name, "Updated")
	}
}

func TestDeleteNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/notification/1", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.DeleteNotification(context.Background(), 1); err != nil {
		t.Fatalf("DeleteNotification: %v", err)
	}
}

func TestGetNotificationSchema(t *testing.T) {
	t.Parallel()
	want := []arr.ProviderResource{{ID: 0, Name: "Email"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/notification/schema", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetNotificationSchema(context.Background())
	if err != nil {
		t.Fatalf("GetNotificationSchema: %v", err)
	}
	if got[0].Name != "Email" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Email")
	}
}

func TestTestNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification/test", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.TestNotification(context.Background(), &arr.ProviderResource{ID: 1}); err != nil {
		t.Fatalf("TestNotification: %v", err)
	}
}

func TestTestAllNotifications(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification/testall", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.TestAllNotifications(context.Background()); err != nil {
		t.Fatalf("TestAllNotifications: %v", err)
	}
}

func TestNotificationAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification/action/testAction", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.NotificationAction(context.Background(), "testAction", &arr.ProviderResource{ID: 1}); err != nil {
		t.Fatalf("NotificationAction: %v", err)
	}
}

// ---------- Config ----------.

func TestGetDownloadClientConfig(t *testing.T) {
	t.Parallel()
	want := arr.DownloadClientConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/downloadclient", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetDownloadClientConfig(context.Background())
	if err != nil {
		t.Fatalf("GetDownloadClientConfig: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetDownloadClientConfigByID(t *testing.T) {
	t.Parallel()
	want := arr.DownloadClientConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/downloadclient/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetDownloadClientConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetDownloadClientConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateDownloadClientConfig(t *testing.T) {
	t.Parallel()
	want := arr.DownloadClientConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/downloadclient/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.UpdateDownloadClientConfig(context.Background(), &arr.DownloadClientConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateDownloadClientConfig: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetHostConfig(t *testing.T) {
	t.Parallel()
	want := arr.HostConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/host", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetHostConfig(context.Background())
	if err != nil {
		t.Fatalf("GetHostConfig: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetHostConfigByID(t *testing.T) {
	t.Parallel()
	want := arr.HostConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/host/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetHostConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetHostConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateHostConfig(t *testing.T) {
	t.Parallel()
	want := arr.HostConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/host/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.UpdateHostConfig(context.Background(), &arr.HostConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateHostConfig: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetUIConfig(t *testing.T) {
	t.Parallel()
	want := arr.UIConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/ui", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetUIConfig(context.Background())
	if err != nil {
		t.Fatalf("GetUIConfig: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetUIConfigByID(t *testing.T) {
	t.Parallel()
	want := arr.UIConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/ui/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetUIConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUIConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateUIConfig(t *testing.T) {
	t.Parallel()
	want := arr.UIConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/ui/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.UpdateUIConfig(context.Background(), &arr.UIConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateUIConfig: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetDevelopmentConfig(t *testing.T) {
	t.Parallel()
	want := prowlarr.DevelopmentConfigResource{ID: 1, LogIndexerResponse: true}
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/development", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetDevelopmentConfig(context.Background())
	if err != nil {
		t.Fatalf("GetDevelopmentConfig: %v", err)
	}
	if !got.LogIndexerResponse {
		t.Error("expected LogIndexerResponse=true")
	}
}

func TestGetDevelopmentConfigByID(t *testing.T) {
	t.Parallel()
	want := prowlarr.DevelopmentConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/development/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetDevelopmentConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetDevelopmentConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateDevelopmentConfig(t *testing.T) {
	t.Parallel()
	want := prowlarr.DevelopmentConfigResource{ID: 1, LogSQL: true}
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/development/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.UpdateDevelopmentConfig(context.Background(), &prowlarr.DevelopmentConfigResource{ID: 1, LogSQL: true})
	if err != nil {
		t.Fatalf("UpdateDevelopmentConfig: %v", err)
	}
	if !got.LogSQL {
		t.Error("expected LogSQL=true")
	}
}

// ---------- Custom Filters ----------.

func TestGetCustomFilters(t *testing.T) {
	t.Parallel()
	want := []arr.CustomFilterResource{{ID: 1, Label: "My Filter"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/customfilter", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetCustomFilters(context.Background())
	if err != nil {
		t.Fatalf("GetCustomFilters: %v", err)
	}
	if got[0].Label != "My Filter" {
		t.Errorf("Label = %q, want %q", got[0].Label, "My Filter")
	}
}

func TestGetCustomFilter(t *testing.T) {
	t.Parallel()
	want := arr.CustomFilterResource{ID: 1, Label: "My Filter"}
	srv := newTestServer(t, http.MethodGet, "/api/v1/customfilter/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetCustomFilter(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCustomFilter: %v", err)
	}
	if got.Label != "My Filter" {
		t.Errorf("Label = %q, want %q", got.Label, "My Filter")
	}
}

func TestCreateCustomFilter(t *testing.T) {
	t.Parallel()
	want := arr.CustomFilterResource{ID: 1, Label: "New Filter"}
	srv := newTestServer(t, http.MethodPost, "/api/v1/customfilter", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.CreateCustomFilter(context.Background(), &arr.CustomFilterResource{Label: "New Filter"})
	if err != nil {
		t.Fatalf("CreateCustomFilter: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateCustomFilter(t *testing.T) {
	t.Parallel()
	want := arr.CustomFilterResource{ID: 1, Label: "Updated"}
	srv := newTestServer(t, http.MethodPut, "/api/v1/customfilter/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.UpdateCustomFilter(context.Background(), &arr.CustomFilterResource{ID: 1, Label: "Updated"})
	if err != nil {
		t.Fatalf("UpdateCustomFilter: %v", err)
	}
	if got.Label != "Updated" {
		t.Errorf("Label = %q, want %q", got.Label, "Updated")
	}
}

func TestDeleteCustomFilter(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/customfilter/1", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.DeleteCustomFilter(context.Background(), 1); err != nil {
		t.Fatalf("DeleteCustomFilter: %v", err)
	}
}

// ---------- Tags Extended ----------.

func TestGetTag(t *testing.T) {
	t.Parallel()
	want := arr.Tag{ID: 1, Label: "test"}
	srv := newTestServer(t, http.MethodGet, "/api/v1/tag/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetTag(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTag: %v", err)
	}
	if got.Label != "test" {
		t.Errorf("Label = %q, want %q", got.Label, "test")
	}
}

func TestUpdateTag(t *testing.T) {
	t.Parallel()
	want := arr.Tag{ID: 1, Label: "updated"}
	srv := newTestServer(t, http.MethodPut, "/api/v1/tag/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.UpdateTag(context.Background(), &arr.Tag{ID: 1, Label: "updated"})
	if err != nil {
		t.Fatalf("UpdateTag: %v", err)
	}
	if got.Label != "updated" {
		t.Errorf("Label = %q, want %q", got.Label, "updated")
	}
}

func TestDeleteTag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/tag/1", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.DeleteTag(context.Background(), 1); err != nil {
		t.Fatalf("DeleteTag: %v", err)
	}
}

func TestGetTagDetails(t *testing.T) {
	t.Parallel()
	want := []arr.TagDetail{{ID: 1, Label: "test"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/tag/detail", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetTagDetails(context.Background())
	if err != nil {
		t.Fatalf("GetTagDetails: %v", err)
	}
	if got[0].Label != "test" {
		t.Errorf("Label = %q, want %q", got[0].Label, "test")
	}
}

func TestGetTagDetail(t *testing.T) {
	t.Parallel()
	want := arr.TagDetail{ID: 1, Label: "test"}
	srv := newTestServer(t, http.MethodGet, "/api/v1/tag/detail/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetTagDetail(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTagDetail: %v", err)
	}
	if got.Label != "test" {
		t.Errorf("Label = %q, want %q", got.Label, "test")
	}
}

// ---------- Backups ----------.

func TestGetBackups(t *testing.T) {
	t.Parallel()
	want := []arr.Backup{{ID: 1, Name: "backup.zip"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/backup", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetBackups(context.Background())
	if err != nil {
		t.Fatalf("GetBackups: %v", err)
	}
	if got[0].Name != "backup.zip" {
		t.Errorf("Name = %q, want %q", got[0].Name, "backup.zip")
	}
}

func TestDeleteBackup(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/system/backup/1", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.DeleteBackup(context.Background(), 1); err != nil {
		t.Fatalf("DeleteBackup: %v", err)
	}
}

func TestRestoreBackup(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/system/backup/restore/1", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.RestoreBackup(context.Background(), 1); err != nil {
		t.Fatalf("RestoreBackup: %v", err)
	}
}

// ---------- Logs ----------.

func TestGetLogs(t *testing.T) {
	t.Parallel()
	want := arr.PagingResource[arr.LogRecord]{Page: 1, PageSize: 10, TotalRecords: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v1/log?page=1&pageSize=10", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetLogs(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetLogs: %v", err)
	}
	if got.Page != 1 {
		t.Errorf("Page = %d, want 1", got.Page)
	}
}

func TestGetLogFiles(t *testing.T) {
	t.Parallel()
	want := []arr.LogFileResource{{ID: 1, Filename: "prowlarr.txt"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/log/file", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetLogFiles(context.Background())
	if err != nil {
		t.Fatalf("GetLogFiles: %v", err)
	}
	if got[0].Filename != "prowlarr.txt" {
		t.Errorf("Filename = %q, want %q", got[0].Filename, "prowlarr.txt")
	}
}

func TestGetLogFileContent(t *testing.T) {
	t.Parallel()
	srv := newRawTestServer(t, http.MethodGet, "/api/v1/log/file/prowlarr.txt", "log content")
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetLogFileContent(context.Background(), "prowlarr.txt")
	if err != nil {
		t.Fatalf("GetLogFileContent: %v", err)
	}
	if got != "log content" {
		t.Errorf("content = %q, want %q", got, "log content")
	}
}

func TestGetUpdateLogFiles(t *testing.T) {
	t.Parallel()
	want := []arr.LogFileResource{{ID: 1, Filename: "update.txt"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/log/file/update", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetUpdateLogFiles(context.Background())
	if err != nil {
		t.Fatalf("GetUpdateLogFiles: %v", err)
	}
	if got[0].Filename != "update.txt" {
		t.Errorf("Filename = %q, want %q", got[0].Filename, "update.txt")
	}
}

func TestGetUpdateLogFileContent(t *testing.T) {
	t.Parallel()
	srv := newRawTestServer(t, http.MethodGet, "/api/v1/log/file/update/update.txt", "update log")
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetUpdateLogFileContent(context.Background(), "update.txt")
	if err != nil {
		t.Fatalf("GetUpdateLogFileContent: %v", err)
	}
	if got != "update log" {
		t.Errorf("content = %q, want %q", got, "update log")
	}
}

// ---------- System ----------.

func TestGetTasks(t *testing.T) {
	t.Parallel()
	want := []arr.TaskResource{{ID: 1, Name: "RssSync"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/task", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetTasks(context.Background())
	if err != nil {
		t.Fatalf("GetTasks: %v", err)
	}
	if got[0].Name != "RssSync" {
		t.Errorf("Name = %q, want %q", got[0].Name, "RssSync")
	}
}

func TestGetTask(t *testing.T) {
	t.Parallel()
	want := arr.TaskResource{ID: 1, Name: "RssSync"}
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/task/1", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetTask(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}
	if got.Name != "RssSync" {
		t.Errorf("Name = %q, want %q", got.Name, "RssSync")
	}
}

func TestGetUpdates(t *testing.T) {
	t.Parallel()
	want := []arr.UpdateResource{{Version: "1.0.0", Installed: true}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/update", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetUpdates(context.Background())
	if err != nil {
		t.Fatalf("GetUpdates: %v", err)
	}
	if got[0].Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", got[0].Version, "1.0.0")
	}
}

func TestGetSystemRoutes(t *testing.T) {
	t.Parallel()
	want := []arr.SystemRouteResource{{Method: "GET", Path: "/api/v1/indexer"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/routes", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetSystemRoutes(context.Background())
	if err != nil {
		t.Fatalf("GetSystemRoutes: %v", err)
	}
	if got[0].Method != http.MethodGet {
		t.Errorf("Method = %q, want %q", got[0].Method, "GET")
	}
}

func TestGetSystemRoutesDuplicate(t *testing.T) {
	t.Parallel()
	want := []arr.SystemRouteResource{{Method: "GET", Path: "/duplicate"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/routes/duplicate", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetSystemRoutesDuplicate(context.Background())
	if err != nil {
		t.Fatalf("GetSystemRoutesDuplicate: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestShutdown(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/system/shutdown", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestRestart(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/system/restart", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.Restart(context.Background()); err != nil {
		t.Fatalf("Restart: %v", err)
	}
}

func TestDeleteCommand(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/command/1", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.DeleteCommand(context.Background(), 1); err != nil {
		t.Fatalf("DeleteCommand: %v", err)
	}
}

// ---------- Localization ----------.

func TestGetLocalization(t *testing.T) {
	t.Parallel()
	want := map[string]string{"hello": "world"}
	srv := newTestServer(t, http.MethodGet, "/api/v1/localization", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetLocalization(context.Background())
	if err != nil {
		t.Fatalf("GetLocalization: %v", err)
	}
	if got["hello"] != "world" {
		t.Errorf("got %q, want %q", got["hello"], "world")
	}
}

func TestGetLocalizationOptions(t *testing.T) {
	t.Parallel()
	want := []prowlarr.LocalizationOption{{Name: "English", Value: "en"}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/localization/options", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetLocalizationOptions(context.Background())
	if err != nil {
		t.Fatalf("GetLocalizationOptions: %v", err)
	}
	if got[0].Value != "en" {
		t.Errorf("Value = %q, want %q", got[0].Value, "en")
	}
}

// ---------- Ping ----------.

func TestPing(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/ping", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

// ---------- File System ----------.

func TestBrowseFileSystem(t *testing.T) {
	t.Parallel()
	want := map[string]any{"parent": "/", "directories": []any{}}
	srv := newTestServer(t, http.MethodGet, "/api/v1/filesystem?path=%2Ftmp", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.BrowseFileSystem(context.Background(), "/tmp")
	if err != nil {
		t.Fatalf("BrowseFileSystem: %v", err)
	}
	if got["parent"] != "/" {
		t.Errorf("parent = %v, want /", got["parent"])
	}
}

func TestGetFileSystemType(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/filesystem/type?path=%2Ftmp", "ext4")
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetFileSystemType(context.Background(), "/tmp")
	if err != nil {
		t.Fatalf("GetFileSystemType: %v", err)
	}
	if got != "ext4" {
		t.Errorf("type = %q, want %q", got, "ext4")
	}
}

// ---------- Search Extended ----------.

func TestGrabReleasesBulk(t *testing.T) {
	t.Parallel()
	want := prowlarr.Release{GUID: "abc123", Title: "Test Release"}
	srv := newTestServer(t, http.MethodPost, "/api/v1/search/bulk", want)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GrabReleasesBulk(context.Background(), []prowlarr.Release{{GUID: "abc123"}})
	if err != nil {
		t.Fatalf("GrabReleasesBulk: %v", err)
	}
	if got.Title != "Test Release" {
		t.Errorf("Title = %q, want %q", got.Title, "Test Release")
	}
}

// ---------- Newznab ----------.

func TestGetIndexerNewznab(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer/1/newznab", "<xml>test</xml>")
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.GetIndexerNewznab(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetIndexerNewznab: %v", err)
	}
	if got != "<xml>test</xml>" {
		t.Errorf("content = %q, want %q", got, "<xml>test</xml>")
	}
}

func TestDownloadIndexerRelease(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer/1/download", "binary-data")
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	got, err := c.DownloadIndexerRelease(context.Background(), 1)
	if err != nil {
		t.Fatalf("DownloadIndexerRelease: %v", err)
	}
	if got != "binary-data" {
		t.Errorf("content = %q, want %q", got, "binary-data")
	}
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	}))
	defer srv.Close()

	c, err := prowlarr.New(srv.URL, "bad-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = c.GetIndexers(context.Background())
	if err == nil {
		t.Fatal("expected error for 401 response")
	}

	var apiErr *arr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *arr.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
}

func TestHeadPing(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodHead, "/ping", nil)
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.HeadPing(context.Background()); err != nil {
		t.Fatalf("HeadPing: %v", err)
	}
}

func TestUploadBackup(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.Header.Get("X-Api-Key") == "" {
			t.Error("missing X-Api-Key header")
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("Content-Type = %q, want multipart/form-data", ct)
		}
		f, fh, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile: %v", err)
		}
		defer f.Close()
		if fh.Filename != "backup.zip" {
			t.Errorf("filename = %q, want %q", fh.Filename, "backup.zip")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	c, _ := prowlarr.New(srv.URL, "test-key")
	if err := c.UploadBackup(context.Background(), "backup.zip", strings.NewReader("fake-backup-data")); err != nil {
		t.Fatalf("UploadBackup: %v", err)
	}
}
