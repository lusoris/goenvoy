package sonarr_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golusoris/goenvoy/arr/sonarr"
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
		c, err := sonarr.New("http://localhost:8989", "test-key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c == nil {
			t.Fatal("expected non-nil client")
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		t.Parallel()
		_, err := sonarr.New("://bad", "test-key")
		if err == nil {
			t.Fatal("expected error for invalid URL")
		}
	})
}

func TestGetAllSeries(t *testing.T) {
	t.Parallel()

	want := []sonarr.Series{
		{ID: 1, Title: "Breaking Bad", TvdbID: 81189},
		{ID: 2, Title: "The Wire", TvdbID: 79126},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v3/series", want)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetAllSeries(context.Background())
	if err != nil {
		t.Fatalf("GetAllSeries: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Title != "Breaking Bad" {
		t.Errorf("Title = %q, want %q", got[0].Title, "Breaking Bad")
	}
}

func TestGetSeries(t *testing.T) {
	t.Parallel()

	want := sonarr.Series{ID: 1, Title: "Breaking Bad"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/series/1", want)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetSeries(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetSeries: %v", err)
	}
	if got.Title != "Breaking Bad" {
		t.Errorf("Title = %q, want %q", got.Title, "Breaking Bad")
	}
}

func TestAddSeries(t *testing.T) {
	t.Parallel()

	want := sonarr.Series{ID: 3, Title: "New Show", TvdbID: 99999}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var body sonarr.Series
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Title != "New Show" {
			t.Errorf("Title = %q, want %q", body.Title, "New Show")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.AddSeries(context.Background(), &sonarr.Series{
		Title:  "New Show",
		TvdbID: 99999,
	})
	if err != nil {
		t.Fatalf("AddSeries: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("ID = %d, want 3", got.ID)
	}
}

func TestDeleteSeries(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete,
		"/api/v3/series/1?deleteFiles=true&addImportListExclusion=false",
		nil)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteSeries(context.Background(), 1, true, false); err != nil {
		t.Fatalf("DeleteSeries: %v", err)
	}
}

func TestLookupSeries(t *testing.T) {
	t.Parallel()

	want := []sonarr.Series{{ID: 0, Title: "Breaking Bad", TvdbID: 81189}}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/series/lookup?term=breaking+bad",
		want)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.LookupSeries(context.Background(), "breaking bad")
	if err != nil {
		t.Fatalf("LookupSeries: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetEpisodes(t *testing.T) {
	t.Parallel()

	want := []sonarr.Episode{
		{ID: 10, SeriesID: 1, SeasonNumber: 1, EpisodeNumber: 1, Title: "Pilot"},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/episode?seriesId=1",
		want)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetEpisodes(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetEpisodes: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Title != "Pilot" {
		t.Errorf("Title = %q, want %q", got[0].Title, "Pilot")
	}
}

func TestGetEpisodeFiles(t *testing.T) {
	t.Parallel()

	want := []sonarr.EpisodeFile{
		{ID: 100, SeriesID: 1, RelativePath: "S01E01.mkv", Size: 1073741824},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/episodefile?seriesId=1",
		want)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetEpisodeFiles(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetEpisodeFiles: %v", err)
	}
	if got[0].Size != 1073741824 {
		t.Errorf("Size = %d, want 1073741824", got[0].Size)
	}
}

func TestSendCommand(t *testing.T) {
	t.Parallel()

	want := arr.CommandResponse{ID: 42, Name: "RefreshSeries"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var cmd arr.CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if cmd.Name != "RefreshSeries" {
			t.Errorf("Name = %q, want RefreshSeries", cmd.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.SendCommand(context.Background(), arr.CommandRequest{Name: "RefreshSeries"})
	if err != nil {
		t.Fatalf("SendCommand: %v", err)
	}
	if got.ID != 42 {
		t.Errorf("ID = %d, want 42", got.ID)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	want := sonarr.ParseResult{
		Title: "Breaking.Bad.S01E01.720p",
		ParsedEpisodeInfo: &sonarr.ParsedEpisodeInfo{
			SeriesTitle:    "Breaking Bad",
			SeasonNumber:   1,
			EpisodeNumbers: []int{1},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/parse?title=Breaking.Bad.S01E01.720p",
		want)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.Parse(context.Background(), "Breaking.Bad.S01E01.720p")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got.ParsedEpisodeInfo.SeriesTitle != "Breaking Bad" {
		t.Errorf("SeriesTitle = %q, want %q", got.ParsedEpisodeInfo.SeriesTitle, "Breaking Bad")
	}
}

func TestGetSystemStatus(t *testing.T) {
	t.Parallel()

	want := arr.StatusResponse{AppName: "Sonarr", Version: "4.0.0"}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/system/status",
		want)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetSystemStatus(context.Background())
	if err != nil {
		t.Fatalf("GetSystemStatus: %v", err)
	}
	if got.AppName != "Sonarr" {
		t.Errorf("AppName = %q, want %q", got.AppName, "Sonarr")
	}
}

func TestGetQueue(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[arr.QueueRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []arr.QueueRecord{
			{ID: 1, Title: "Breaking Bad - S01E01"},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/queue?page=1&pageSize=10",
		want)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetQueue(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetQueue: %v", err)
	}
	if got.TotalRecords != 1 {
		t.Errorf("TotalRecords = %d, want 1", got.TotalRecords)
	}
}

func TestGetTags(t *testing.T) {
	t.Parallel()

	want := []arr.Tag{{ID: 1, Label: "hd"}, {ID: 2, Label: "anime"}}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/tag",
		want)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
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

func TestCreateTag(t *testing.T) {
	t.Parallel()

	want := arr.Tag{ID: 3, Label: "new-tag"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var tag arr.Tag
		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if tag.Label != "new-tag" {
			t.Errorf("Label = %q, want %q", tag.Label, "new-tag")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.CreateTag(context.Background(), "new-tag")
	if err != nil {
		t.Fatalf("CreateTag: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("ID = %d, want 3", got.ID)
	}
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	}))
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "bad-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = c.GetAllSeries(context.Background())
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

func TestGetHistory(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[sonarr.HistoryRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []sonarr.HistoryRecord{
			{ID: 5, EpisodeID: 10, SeriesID: 1, EventType: "grabbed"},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/history?page=1&pageSize=10",
		want)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetHistory(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetHistory: %v", err)
	}
	if got.Records[0].EventType != "grabbed" {
		t.Errorf("EventType = %q, want %q", got.Records[0].EventType, "grabbed")
	}
}

func TestDeleteQueueItem(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete,
		"/api/v3/queue/5?removeFromClient=true&blocklist=false",
		nil)
	defer srv.Close()

	c, err := sonarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteQueueItem(context.Background(), 5, true, false); err != nil {
		t.Fatalf("DeleteQueueItem: %v", err)
	}
}

func TestUpdateSeries(t *testing.T) {
	t.Parallel()

	want := sonarr.Series{ID: 1, Title: "Updated"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.RequestURI() != "/api/v3/series/1?moveFiles=true" {
			t.Errorf("path = %q", r.URL.RequestURI())
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateSeries(context.Background(), &sonarr.Series{ID: 1, Title: "Updated"}, true)
	if err != nil {
		t.Fatalf("UpdateSeries: %v", err)
	}
	if got.Title != "Updated" {
		t.Errorf("Title = %q", got.Title)
	}
}

func TestGetEpisode(t *testing.T) {
	t.Parallel()

	want := sonarr.Episode{ID: 10, Title: "Pilot"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/episode/10", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetEpisode(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetEpisode: %v", err)
	}
	if got.Title != "Pilot" {
		t.Errorf("Title = %q", got.Title)
	}
}

func TestUpdateEpisode(t *testing.T) {
	t.Parallel()

	want := sonarr.Episode{ID: 10, Title: "Updated", Monitored: true}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateEpisode(context.Background(), &sonarr.Episode{ID: 10, Monitored: true})
	if err != nil {
		t.Fatalf("UpdateEpisode: %v", err)
	}
	if !got.Monitored {
		t.Error("expected Monitored=true")
	}
}

func TestMonitorEpisodes(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.MonitorEpisodes(context.Background(), []int{1, 2, 3}, true); err != nil {
		t.Fatalf("MonitorEpisodes: %v", err)
	}
}

func TestGetEpisodeFile(t *testing.T) {
	t.Parallel()

	want := sonarr.EpisodeFile{ID: 100, SeriesID: 1, Size: 1073741824}

	srv := newTestServer(t, http.MethodGet, "/api/v3/episodefile/100", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetEpisodeFile(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetEpisodeFile: %v", err)
	}
	if got.Size != 1073741824 {
		t.Errorf("Size = %d", got.Size)
	}
}

func TestDeleteEpisodeFile(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/episodefile/100", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.DeleteEpisodeFile(context.Background(), 100); err != nil {
		t.Fatalf("DeleteEpisodeFile: %v", err)
	}
}

func TestDeleteEpisodeFiles(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/episodefile/bulk", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.DeleteEpisodeFiles(context.Background(), []int{1, 2, 3}); err != nil {
		t.Fatalf("DeleteEpisodeFiles: %v", err)
	}
}

func TestGetCommands(t *testing.T) {
	t.Parallel()

	want := []arr.CommandResponse{{ID: 1, Name: "RefreshSeries"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/command", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetCommands(context.Background())
	if err != nil {
		t.Fatalf("GetCommands: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestGetCommand(t *testing.T) {
	t.Parallel()

	want := arr.CommandResponse{ID: 42, Name: "RefreshSeries"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/command/42", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetCommand(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetCommand: %v", err)
	}
	if got.Name != "RefreshSeries" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestGetCalendar(t *testing.T) {
	t.Parallel()

	want := []sonarr.Episode{{ID: 1, Title: "Upcoming"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/calendar?start=2026-01-01&end=2026-01-31&unmonitored=false", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetCalendar(context.Background(), "2026-01-01", "2026-01-31", false)
	if err != nil {
		t.Fatalf("GetCalendar: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestGetHealth(t *testing.T) {
	t.Parallel()

	want := []arr.HealthCheck{{Type: "warning", Message: "test"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/health", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestGetDiskSpace(t *testing.T) {
	t.Parallel()

	want := []arr.DiskSpace{{Path: "/data", FreeSpace: 1000}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/diskspace", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetDiskSpace(context.Background())
	if err != nil {
		t.Fatalf("GetDiskSpace: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestGetQualityProfiles(t *testing.T) {
	t.Parallel()

	want := []arr.QualityProfile{{ID: 1, Name: "Any"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/qualityprofile", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetQualityProfiles(context.Background())
	if err != nil {
		t.Fatalf("GetQualityProfiles: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestGetRootFolders(t *testing.T) {
	t.Parallel()

	want := []arr.RootFolder{{ID: 1, Path: "/tv"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/rootfolder", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetRootFolders(context.Background())
	if err != nil {
		t.Fatalf("GetRootFolders: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestUpdateSeasonPass(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.UpdateSeasonPass(context.Background(), sonarr.SeasonPassResource{}); err != nil {
		t.Fatalf("UpdateSeasonPass: %v", err)
	}
}

func TestDeleteCommand(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/command/1", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.DeleteCommand(context.Background(), 1); err != nil {
		t.Fatalf("DeleteCommand: %v", err)
	}
}

func TestUpdateEpisodeFile(t *testing.T) {
	t.Parallel()

	want := sonarr.EpisodeFile{ID: 1}

	srv := newTestServer(t, http.MethodPut, "/api/v3/episodefile/1", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateEpisodeFile(context.Background(), &sonarr.EpisodeFile{ID: 1})
	if err != nil {
		t.Fatalf("UpdateEpisodeFile: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestEditEpisodeFiles(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPut, "/api/v3/episodefile/editor", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.EditEpisodeFiles(context.Background(), &sonarr.EpisodeFileEditorResource{
		EpisodeFileIDs: []int{1, 2},
	}); err != nil {
		t.Fatalf("EditEpisodeFiles: %v", err)
	}
}

func TestUpdateCustomFormatsBulk(t *testing.T) {
	t.Parallel()

	want := []arr.CustomFormatResource{{ID: 1, Name: "test"}}

	srv := newTestServer(t, http.MethodPut, "/api/v3/customformat/bulk", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateCustomFormatsBulk(context.Background(), &arr.CustomFormatBulkResource{IDs: []int{1}})
	if err != nil {
		t.Fatalf("UpdateCustomFormatsBulk: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestDeleteCustomFormatsBulk(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/customformat/bulk", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.DeleteCustomFormatsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteCustomFormatsBulk: %v", err)
	}
}

func TestUpdateDownloadClientsBulk(t *testing.T) {
	t.Parallel()

	want := []arr.ProviderResource{{ID: 1}}

	srv := newTestServer(t, http.MethodPut, "/api/v3/downloadclient/bulk", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateDownloadClientsBulk(context.Background(), &arr.ProviderBulkResource{IDs: []int{1}})
	if err != nil {
		t.Fatalf("UpdateDownloadClientsBulk: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestDeleteDownloadClientsBulk(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/downloadclient/bulk", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.DeleteDownloadClientsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteDownloadClientsBulk: %v", err)
	}
}

func TestTestAllDownloadClients(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPost, "/api/v3/downloadclient/testall", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.TestAllDownloadClients(context.Background()); err != nil {
		t.Fatalf("TestAllDownloadClients: %v", err)
	}
}

func TestUpdateIndexersBulk(t *testing.T) {
	t.Parallel()

	want := []arr.ProviderResource{{ID: 1}}

	srv := newTestServer(t, http.MethodPut, "/api/v3/indexer/bulk", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateIndexersBulk(context.Background(), &arr.ProviderBulkResource{IDs: []int{1}})
	if err != nil {
		t.Fatalf("UpdateIndexersBulk: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestDeleteIndexersBulk(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/indexer/bulk", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.DeleteIndexersBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteIndexersBulk: %v", err)
	}
}

func TestTestAllIndexers(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPost, "/api/v3/indexer/testall", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.TestAllIndexers(context.Background()); err != nil {
		t.Fatalf("TestAllIndexers: %v", err)
	}
}

func TestUpdateImportListsBulk(t *testing.T) {
	t.Parallel()

	want := []arr.ProviderResource{{ID: 1}}

	srv := newTestServer(t, http.MethodPut, "/api/v3/importlist/bulk", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateImportListsBulk(context.Background(), &arr.ProviderBulkResource{IDs: []int{1}})
	if err != nil {
		t.Fatalf("UpdateImportListsBulk: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestDeleteImportListsBulk(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/importlist/bulk", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.DeleteImportListsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteImportListsBulk: %v", err)
	}
}

func TestTestAllImportLists(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPost, "/api/v3/importlist/testall", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.TestAllImportLists(context.Background()); err != nil {
		t.Fatalf("TestAllImportLists: %v", err)
	}
}

func TestGetImportListConfig(t *testing.T) {
	t.Parallel()

	want := sonarr.ImportListConfigResource{ID: 1, ListSyncLevel: "disabled"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/config/importlist", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetImportListConfig(context.Background())
	if err != nil {
		t.Fatalf("GetImportListConfig: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateImportListConfig(t *testing.T) {
	t.Parallel()

	want := sonarr.ImportListConfigResource{ID: 1, ListSyncLevel: "logOnly"}

	srv := newTestServer(t, http.MethodPut, "/api/v3/config/importlist/1", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateImportListConfig(context.Background(), &sonarr.ImportListConfigResource{ID: 1, ListSyncLevel: "logOnly"})
	if err != nil {
		t.Fatalf("UpdateImportListConfig: %v", err)
	}
	if got.ListSyncLevel != "logOnly" {
		t.Errorf("ListSyncLevel = %q, want logOnly", got.ListSyncLevel)
	}
}

func TestGetImportListExclusionsPaged(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[arr.ImportListExclusionResource]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records:      []arr.ImportListExclusionResource{{ID: 1, TvdbID: 123}},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v3/importlistexclusion/paged?page=1&pageSize=10", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetImportListExclusionsPaged(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetImportListExclusionsPaged: %v", err)
	}
	if got.TotalRecords != 1 {
		t.Errorf("TotalRecords = %d, want 1", got.TotalRecords)
	}
}

func TestDeleteImportListExclusionsBulk(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/importlistexclusion/bulk", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.DeleteImportListExclusionsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteImportListExclusionsBulk: %v", err)
	}
}

func TestTestAllNotifications(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPost, "/api/v3/notification/testall", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.TestAllNotifications(context.Background()); err != nil {
		t.Fatalf("TestAllNotifications: %v", err)
	}
}

func TestTestAllMetadataConsumers(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPost, "/api/v3/metadata/testall", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.TestAllMetadataConsumers(context.Background()); err != nil {
		t.Fatalf("TestAllMetadataConsumers: %v", err)
	}
}

func TestGetLanguage(t *testing.T) {
	t.Parallel()

	want := sonarr.Language{ID: 1, Name: "English"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/language/1", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetLanguage(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetLanguage: %v", err)
	}
	if got.Name != "English" {
		t.Errorf("Name = %q, want English", got.Name)
	}
}

func TestGetLocalization(t *testing.T) {
	t.Parallel()

	want := sonarr.LocalizationResource{ID: 1, Strings: map[string]string{"key": "val"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/localization", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetLocalization(context.Background())
	if err != nil {
		t.Fatalf("GetLocalization: %v", err)
	}
	if got.Strings["key"] != "val" {
		t.Errorf("Strings[key] = %q, want val", got.Strings["key"])
	}
}

func TestUpdateQualityDefinitions(t *testing.T) {
	t.Parallel()

	want := []arr.QualityDefinitionResource{{ID: 1, Title: "HDTV-720p"}}

	srv := newTestServer(t, http.MethodPut, "/api/v3/qualitydefinition/update", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateQualityDefinitions(context.Background(), []arr.QualityDefinitionResource{{ID: 1, Title: "HDTV-720p"}})
	if err != nil {
		t.Fatalf("UpdateQualityDefinitions: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestGetQualityProfileSchema(t *testing.T) {
	t.Parallel()

	want := arr.QualityProfile{ID: 1, Name: "schema"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/qualityprofile/schema", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetQualityProfileSchema(context.Background())
	if err != nil {
		t.Fatalf("GetQualityProfileSchema: %v", err)
	}
	if got.Name != "schema" {
		t.Errorf("Name = %q, want schema", got.Name)
	}
}

func TestUpdateRootFolder(t *testing.T) {
	t.Parallel()

	want := arr.RootFolder{ID: 1, Path: "/tv"}

	srv := newTestServer(t, http.MethodPut, "/api/v3/rootfolder/1", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateRootFolder(context.Background(), &arr.RootFolder{ID: 1, Path: "/tv"})
	if err != nil {
		t.Fatalf("UpdateRootFolder: %v", err)
	}
	if got.Path != "/tv" {
		t.Errorf("Path = %q, want /tv", got.Path)
	}
}

func TestBrowseFileSystem(t *testing.T) {
	t.Parallel()

	want := sonarr.FileSystemResource{
		Directories: []sonarr.FileSystemEntry{{Path: "/tv", Name: "tv"}},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v3/filesystem?path=%2Ftv&includeFiles=true", want)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.BrowseFileSystem(context.Background(), "/tv", true)
	if err != nil {
		t.Fatalf("BrowseFileSystem: %v", err)
	}
	if len(got.Directories) != 1 {
		t.Errorf("len(Directories) = %d", len(got.Directories))
	}
}

func TestPing(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodGet, "/ping", nil)
	defer srv.Close()

	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestGetCalendarByID(t *testing.T) {
	t.Parallel()
	want := sonarr.Episode{ID: 42, Title: "Pilot"}
	srv := newTestServer(t, http.MethodGet, "/api/v3/calendar/42", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetCalendarByID(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetCalendarByID: %v", err)
	}
	if got.ID != 42 {
		t.Errorf("ID = %d, want 42", got.ID)
	}
}

func TestGetWantedCutoffByID(t *testing.T) {
	t.Parallel()
	want := sonarr.Episode{ID: 7, Title: "Cutoff"}
	srv := newTestServer(t, http.MethodGet, "/api/v3/wanted/cutoff/7", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetWantedCutoffByID(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetWantedCutoffByID: %v", err)
	}
	if got.ID != 7 {
		t.Errorf("ID = %d, want 7", got.ID)
	}
}

func TestGetWantedMissingByID(t *testing.T) {
	t.Parallel()
	want := sonarr.Episode{ID: 3, Title: "Missing"}
	srv := newTestServer(t, http.MethodGet, "/api/v3/wanted/missing/3", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetWantedMissingByID(context.Background(), 3)
	if err != nil {
		t.Fatalf("GetWantedMissingByID: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("ID = %d, want 3", got.ID)
	}
}

func TestGetDownloadClientConfigByID(t *testing.T) {
	t.Parallel()
	want := arr.DownloadClientConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/downloadclient/1", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetDownloadClientConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetDownloadClientConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetHostConfigByID(t *testing.T) {
	t.Parallel()
	want := arr.HostConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/host/1", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetHostConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetHostConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetImportListConfigByID(t *testing.T) {
	t.Parallel()
	want := sonarr.ImportListConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/importlist/1", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetImportListConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetImportListConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetIndexerConfigByID(t *testing.T) {
	t.Parallel()
	want := arr.IndexerConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/indexer/1", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetIndexerConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetIndexerConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetMediaManagementConfigByID(t *testing.T) {
	t.Parallel()
	want := arr.MediaManagementConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/mediamanagement/1", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetMediaManagementConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMediaManagementConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetNamingConfigByID(t *testing.T) {
	t.Parallel()
	want := arr.NamingConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/naming/1", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetNamingConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetNamingConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetUIConfigByID(t *testing.T) {
	t.Parallel()
	want := arr.UIConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/ui/1", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetUIConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUIConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestDownloadClientAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v3/downloadclient/action/testAction", nil)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.DownloadClientAction(context.Background(), "testAction", nil); err != nil {
		t.Fatalf("DownloadClientAction: %v", err)
	}
}

func TestImportListAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v3/importlist/action/testAction", nil)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.ImportListAction(context.Background(), "testAction", nil); err != nil {
		t.Fatalf("ImportListAction: %v", err)
	}
}

func TestIndexerAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v3/indexer/action/testAction", nil)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.IndexerAction(context.Background(), "testAction", nil); err != nil {
		t.Fatalf("IndexerAction: %v", err)
	}
}

func TestMetadataAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v3/metadata/action/testAction", nil)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.MetadataAction(context.Background(), "testAction", nil); err != nil {
		t.Fatalf("MetadataAction: %v", err)
	}
}

func TestNotificationAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v3/notification/action/testAction", nil)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.NotificationAction(context.Background(), "testAction", nil); err != nil {
		t.Fatalf("NotificationAction: %v", err)
	}
}

func TestGetLanguageProfiles(t *testing.T) {
	t.Parallel()
	want := []sonarr.LanguageProfileResource{{ID: 1, Name: "English"}}
	srv := newTestServer(t, http.MethodGet, "/api/v3/languageprofile", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetLanguageProfiles(context.Background())
	if err != nil {
		t.Fatalf("GetLanguageProfiles: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetLanguageProfile(t *testing.T) {
	t.Parallel()
	want := sonarr.LanguageProfileResource{ID: 1, Name: "English"}
	srv := newTestServer(t, http.MethodGet, "/api/v3/languageprofile/1", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetLanguageProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetLanguageProfile: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestCreateLanguageProfile(t *testing.T) {
	t.Parallel()
	want := sonarr.LanguageProfileResource{ID: 1, Name: "English"}
	srv := newTestServer(t, http.MethodPost, "/api/v3/languageprofile", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.CreateLanguageProfile(context.Background(), &sonarr.LanguageProfileResource{Name: "English"})
	if err != nil {
		t.Fatalf("CreateLanguageProfile: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateLanguageProfile(t *testing.T) {
	t.Parallel()
	want := sonarr.LanguageProfileResource{ID: 1, Name: "Updated"}
	srv := newTestServer(t, http.MethodPut, "/api/v3/languageprofile/1", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateLanguageProfile(context.Background(), &sonarr.LanguageProfileResource{ID: 1, Name: "Updated"})
	if err != nil {
		t.Fatalf("UpdateLanguageProfile: %v", err)
	}
	if got.Name != "Updated" {
		t.Errorf("Name = %q, want Updated", got.Name)
	}
}

func TestDeleteLanguageProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v3/languageprofile/1", nil)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.DeleteLanguageProfile(context.Background(), 1); err != nil {
		t.Fatalf("DeleteLanguageProfile: %v", err)
	}
}

func TestGetLanguageProfileSchema(t *testing.T) {
	t.Parallel()
	want := sonarr.LanguageProfileResource{ID: 0, Name: "Schema"}
	srv := newTestServer(t, http.MethodGet, "/api/v3/languageprofile/schema", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetLanguageProfileSchema(context.Background())
	if err != nil {
		t.Fatalf("GetLanguageProfileSchema: %v", err)
	}
	if got.Name != "Schema" {
		t.Errorf("Name = %q, want Schema", got.Name)
	}
}

func TestGetLocalizationByID(t *testing.T) {
	t.Parallel()
	want := sonarr.LocalizationResource{ID: 1, Strings: map[string]string{"hello": "world"}}
	srv := newTestServer(t, http.MethodGet, "/api/v3/localization/1", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetLocalizationByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetLocalizationByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetLocalizationLanguages(t *testing.T) {
	t.Parallel()
	want := []sonarr.LocalizationLanguageResource{{Identifier: "en"}}
	srv := newTestServer(t, http.MethodGet, "/api/v3/localization/language", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetLocalizationLanguages(context.Background())
	if err != nil {
		t.Fatalf("GetLocalizationLanguages: %v", err)
	}
	if len(got) != 1 || got[0].Identifier != "en" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestGetNamingExamples(t *testing.T) {
	t.Parallel()
	want := arr.NamingConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/naming/examples", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetNamingExamples(context.Background())
	if err != nil {
		t.Fatalf("GetNamingExamples: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetQualityDefinitionLimits(t *testing.T) {
	t.Parallel()
	want := sonarr.QualityDefinitionLimitsResource{Min: 1, Max: 400}
	srv := newTestServer(t, http.MethodGet, "/api/v3/qualitydefinition/limits", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetQualityDefinitionLimits(context.Background())
	if err != nil {
		t.Fatalf("GetQualityDefinitionLimits: %v", err)
	}
	if got.Min != 1 || got.Max != 400 {
		t.Errorf("got Min=%d Max=%d, want 1/400", got.Min, got.Max)
	}
}

func TestUpdateEpisodeFilesBulk(t *testing.T) {
	t.Parallel()
	want := []sonarr.EpisodeFile{{ID: 1}, {ID: 2}}
	srv := newTestServer(t, http.MethodPut, "/api/v3/episodefile/bulk", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.UpdateEpisodeFilesBulk(context.Background(), &sonarr.EpisodeFileEditorResource{
		EpisodeFileIDs: []int{1, 2},
	})
	if err != nil {
		t.Fatalf("UpdateEpisodeFilesBulk: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("len = %d, want 2", len(got))
	}
}

func TestGetUpdateLogFileContent(t *testing.T) {
	t.Parallel()
	srv := newRawTestServer(t, http.MethodGet, "/api/v3/log/file/update/update.txt", "log content")
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetUpdateLogFileContent(context.Background(), "update.txt")
	if err != nil {
		t.Fatalf("GetUpdateLogFileContent: %v", err)
	}
	if got != "log content" {
		t.Errorf("content = %q, want %q", got, "log content")
	}
}

func TestGetSystemRoutesDuplicate(t *testing.T) {
	t.Parallel()
	want := []arr.SystemRouteResource{{Path: "/api/v3/test", Method: "GET"}}
	srv := newTestServer(t, http.MethodGet, "/api/v3/system/routes/duplicate", want)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetSystemRoutesDuplicate(context.Background())
	if err != nil {
		t.Fatalf("GetSystemRoutesDuplicate: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d, want 1", len(got))
	}
}

func TestGetLogFileContent(t *testing.T) {
	t.Parallel()
	srv := newRawTestServer(t, http.MethodGet, "/api/v3/log/file/sonarr.txt", "log line 1\nlog line 2")
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
	got, err := c.GetLogFileContent(context.Background(), "sonarr.txt")
	if err != nil {
		t.Fatalf("GetLogFileContent: %v", err)
	}
	if got != "log line 1\nlog line 2" {
		t.Errorf("content = %q, want %q", got, "log line 1\nlog line 2")
	}
}

func TestHeadPing(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodHead, "/ping", nil)
	defer srv.Close()
	c, _ := sonarr.New(srv.URL, "test-key")
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
	c, _ := sonarr.New(srv.URL, "test-key")
	if err := c.UploadBackup(context.Background(), "backup.zip", strings.NewReader("fake-backup-data")); err != nil {
		t.Fatalf("UploadBackup: %v", err)
	}
}
