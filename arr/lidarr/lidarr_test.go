package lidarr_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golusoris/goenvoy/arr/lidarr"
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
		c, err := lidarr.New("http://localhost:8686", "test-key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c == nil {
			t.Fatal("expected non-nil client")
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		t.Parallel()
		_, err := lidarr.New("://bad", "test-key")
		if err == nil {
			t.Fatal("expected error for invalid URL")
		}
	})
}

func TestGetAllArtists(t *testing.T) {
	t.Parallel()

	want := []lidarr.Artist{
		{ID: 1, ArtistName: "Radiohead", ForeignArtistID: "a74b1b7f-71a5-4011-9441-d0b5e4122711"},
		{ID: 2, ArtistName: "Pink Floyd", ForeignArtistID: "83d91898-7763-47d7-b03b-b92132571559"},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/artist", want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetAllArtists(context.Background())
	if err != nil {
		t.Fatalf("GetAllArtists: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].ArtistName != "Radiohead" {
		t.Errorf("ArtistName = %q, want %q", got[0].ArtistName, "Radiohead")
	}
}

func TestGetArtist(t *testing.T) {
	t.Parallel()

	want := lidarr.Artist{ID: 1, ArtistName: "Radiohead"}

	srv := newTestServer(t, http.MethodGet, "/api/v1/artist/1", want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetArtist(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetArtist: %v", err)
	}
	if got.ArtistName != "Radiohead" {
		t.Errorf("ArtistName = %q, want %q", got.ArtistName, "Radiohead")
	}
}

func TestAddArtist(t *testing.T) {
	t.Parallel()

	want := lidarr.Artist{ID: 3, ArtistName: "New Artist", ForeignArtistID: "abc-123"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var body lidarr.Artist
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.ArtistName != "New Artist" {
			t.Errorf("ArtistName = %q, want %q", body.ArtistName, "New Artist")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.AddArtist(context.Background(), &lidarr.Artist{
		ArtistName:      "New Artist",
		ForeignArtistID: "abc-123",
	})
	if err != nil {
		t.Fatalf("AddArtist: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("ID = %d, want 3", got.ID)
	}
}

func TestDeleteArtist(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete,
		"/api/v1/artist/1?deleteFiles=true&addImportListExclusion=false",
		nil)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteArtist(context.Background(), 1, true, false); err != nil {
		t.Fatalf("DeleteArtist: %v", err)
	}
}

func TestLookupArtist(t *testing.T) {
	t.Parallel()

	want := []lidarr.Artist{{ID: 0, ArtistName: "Radiohead"}}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/artist/lookup?term=radiohead",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.LookupArtist(context.Background(), "radiohead")
	if err != nil {
		t.Fatalf("LookupArtist: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetAlbums(t *testing.T) {
	t.Parallel()

	want := []lidarr.Album{
		{ID: 10, Title: "OK Computer", ArtistID: 1, ForeignAlbumID: "album-1"},
		{ID: 11, Title: "Kid A", ArtistID: 1, ForeignAlbumID: "album-2"},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/album?artistId=1",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetAlbums(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetAlbums: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Title != "OK Computer" {
		t.Errorf("Title = %q, want %q", got[0].Title, "OK Computer")
	}
}

func TestGetAlbum(t *testing.T) {
	t.Parallel()

	want := lidarr.Album{ID: 10, Title: "OK Computer"}

	srv := newTestServer(t, http.MethodGet, "/api/v1/album/10", want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetAlbum(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetAlbum: %v", err)
	}
	if got.Title != "OK Computer" {
		t.Errorf("Title = %q, want %q", got.Title, "OK Computer")
	}
}

func TestDeleteAlbum(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete,
		"/api/v1/album/10?deleteFiles=false&addImportListExclusion=true",
		nil)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteAlbum(context.Background(), 10, false, true); err != nil {
		t.Fatalf("DeleteAlbum: %v", err)
	}
}

func TestLookupAlbum(t *testing.T) {
	t.Parallel()

	want := []lidarr.Album{{ID: 0, Title: "OK Computer"}}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/album/lookup?term=ok+computer",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.LookupAlbum(context.Background(), "ok computer")
	if err != nil {
		t.Fatalf("LookupAlbum: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetTracks(t *testing.T) {
	t.Parallel()

	want := []lidarr.Track{
		{ID: 100, ArtistID: 1, AlbumID: 10, Title: "Airbag", TrackNumber: "1", Duration: 284},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/track?artistId=1&albumId=10",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetTracks(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetTracks: %v", err)
	}
	if got[0].Title != "Airbag" {
		t.Errorf("Title = %q, want %q", got[0].Title, "Airbag")
	}
}

func TestGetTrackFiles(t *testing.T) {
	t.Parallel()

	want := []lidarr.TrackFile{
		{ID: 200, ArtistID: 1, AlbumID: 10, Path: "/music/Radiohead/OK Computer/01 - Airbag.flac", Size: 40000000},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/trackfile?artistId=1",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetTrackFiles(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTrackFiles: %v", err)
	}
	if got[0].Size != 40000000 {
		t.Errorf("Size = %d, want 40000000", got[0].Size)
	}
}

func TestDeleteTrackFile(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v1/trackfile/200", nil)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteTrackFile(context.Background(), 200); err != nil {
		t.Fatalf("DeleteTrackFile: %v", err)
	}
}

func TestSendCommand(t *testing.T) {
	t.Parallel()

	want := arr.CommandResponse{ID: 42, Name: "RefreshArtist"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var cmd arr.CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if cmd.Name != "RefreshArtist" {
			t.Errorf("Name = %q, want RefreshArtist", cmd.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.SendCommand(context.Background(), arr.CommandRequest{Name: "RefreshArtist"})
	if err != nil {
		t.Fatalf("SendCommand: %v", err)
	}
	if got.ID != 42 {
		t.Errorf("ID = %d, want 42", got.ID)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	want := lidarr.ParseResult{
		Title: "Radiohead - OK Computer (1997) [FLAC]",
		ParsedAlbumInfo: &lidarr.ParsedAlbumInfo{
			ArtistName: "Radiohead",
			AlbumTitle: "OK Computer",
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/parse?title=Radiohead+-+OK+Computer+%281997%29+%5BFLAC%5D",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.Parse(context.Background(), "Radiohead - OK Computer (1997) [FLAC]")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got.ParsedAlbumInfo.ArtistName != "Radiohead" {
		t.Errorf("ArtistName = %q, want %q", got.ParsedAlbumInfo.ArtistName, "Radiohead")
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()

	want := []lidarr.SearchResult{
		{ID: 1, ForeignID: "abc-123", Artist: &lidarr.Artist{ArtistName: "Radiohead"}},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/search?term=radiohead",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.Search(context.Background(), "radiohead")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Artist.ArtistName != "Radiohead" {
		t.Errorf("ArtistName = %q, want %q", got[0].Artist.ArtistName, "Radiohead")
	}
}

func TestGetSystemStatus(t *testing.T) {
	t.Parallel()

	want := arr.StatusResponse{AppName: "Lidarr", Version: "2.0.0"}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/system/status",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetSystemStatus(context.Background())
	if err != nil {
		t.Fatalf("GetSystemStatus: %v", err)
	}
	if got.AppName != "Lidarr" {
		t.Errorf("AppName = %q, want %q", got.AppName, "Lidarr")
	}
}

func TestGetQueue(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[arr.QueueRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []arr.QueueRecord{
			{ID: 1, Title: "Radiohead - OK Computer"},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/queue?page=1&pageSize=10",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
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

	want := []arr.Tag{{ID: 1, Label: "flac"}, {ID: 2, Label: "vinyl"}}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/tag",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
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

	c, err := lidarr.New(srv.URL, "test-key")
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

func TestGetMetadataProfiles(t *testing.T) {
	t.Parallel()

	want := []lidarr.MetadataProfile{
		{ID: 1, Name: "Standard"},
		{ID: 2, Name: "None"},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/metadataprofile", want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetMetadataProfiles(context.Background())
	if err != nil {
		t.Fatalf("GetMetadataProfiles: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Name != "Standard" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Standard")
	}
}

func TestGetHistory(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[lidarr.HistoryRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []lidarr.HistoryRecord{
			{ID: 5, ArtistID: 1, AlbumID: 10, EventType: "grabbed"},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/history?page=1&pageSize=10",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
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

func TestGetWantedMissing(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[lidarr.Album]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []lidarr.Album{
			{ID: 10, Title: "OK Computer"},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/wanted/missing?page=1&pageSize=10",
		want)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetWantedMissing(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetWantedMissing: %v", err)
	}
	if got.Records[0].Title != "OK Computer" {
		t.Errorf("Title = %q, want %q", got.Records[0].Title, "OK Computer")
	}
}

func TestDeleteQueueItem(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete,
		"/api/v1/queue/5?removeFromClient=true&blocklist=false",
		nil)
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteQueueItem(context.Background(), 5, true, false); err != nil {
		t.Fatalf("DeleteQueueItem: %v", err)
	}
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	}))
	defer srv.Close()

	c, err := lidarr.New(srv.URL, "bad-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = c.GetAllArtists(context.Background())
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

func TestUpdateArtist(t *testing.T) {
	t.Parallel()

	want := lidarr.Artist{ID: 1, ArtistName: "Updated"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.RequestURI() != "/api/v1/artist/1?moveFiles=true" {
			t.Errorf("path = %q", r.URL.RequestURI())
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateArtist(context.Background(), &lidarr.Artist{ID: 1, ArtistName: "Updated"}, true)
	if err != nil {
		t.Fatalf("UpdateArtist: %v", err)
	}
	if got.ArtistName != "Updated" {
		t.Errorf("ArtistName = %q", got.ArtistName)
	}
}

func TestAddAlbum(t *testing.T) {
	t.Parallel()

	want := lidarr.Album{ID: 20, Title: "New Album"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.AddAlbum(context.Background(), &lidarr.Album{Title: "New Album"})
	if err != nil {
		t.Fatalf("AddAlbum: %v", err)
	}
	if got.ID != 20 {
		t.Errorf("ID = %d", got.ID)
	}
}

func TestUpdateAlbum(t *testing.T) {
	t.Parallel()

	want := lidarr.Album{ID: 10, Title: "Updated Album"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateAlbum(context.Background(), &lidarr.Album{ID: 10, Title: "Updated Album"})
	if err != nil {
		t.Fatalf("UpdateAlbum: %v", err)
	}
	if got.Title != "Updated Album" {
		t.Errorf("Title = %q", got.Title)
	}
}

func TestMonitorAlbums(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.MonitorAlbums(context.Background(), &lidarr.AlbumsMonitoredResource{
		AlbumIDs:  []int{10, 11},
		Monitored: true,
	}); err != nil {
		t.Fatalf("MonitorAlbums: %v", err)
	}
}

func TestGetTrack(t *testing.T) {
	t.Parallel()

	want := lidarr.Track{ID: 100, Title: "Airbag"}

	srv := newTestServer(t, http.MethodGet, "/api/v1/track/100", want)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetTrack(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetTrack: %v", err)
	}
	if got.Title != "Airbag" {
		t.Errorf("Title = %q", got.Title)
	}
}

func TestGetTrackFile(t *testing.T) {
	t.Parallel()

	want := lidarr.TrackFile{ID: 200, Size: 40000000}

	srv := newTestServer(t, http.MethodGet, "/api/v1/trackfile/200", want)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetTrackFile(context.Background(), 200)
	if err != nil {
		t.Fatalf("GetTrackFile: %v", err)
	}
	if got.Size != 40000000 {
		t.Errorf("Size = %d", got.Size)
	}
}

func TestDeleteTrackFiles(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v1/trackfile/bulk", nil)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteTrackFiles(context.Background(), []int{1, 2, 3}); err != nil {
		t.Fatalf("DeleteTrackFiles: %v", err)
	}
}

func TestGetCalendar(t *testing.T) {
	t.Parallel()

	want := []lidarr.Album{{ID: 10, Title: "Upcoming"}}

	srv := newTestServer(t, http.MethodGet, "/api/v1/calendar?start=2026-01-01&end=2026-01-31&unmonitored=false", want)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetCalendar(context.Background(), "2026-01-01", "2026-01-31", false)
	if err != nil {
		t.Fatalf("GetCalendar: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestGetCommands(t *testing.T) {
	t.Parallel()

	want := []arr.CommandResponse{{ID: 1, Name: "RefreshArtist"}}

	srv := newTestServer(t, http.MethodGet, "/api/v1/command", want)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
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

	want := arr.CommandResponse{ID: 42, Name: "RefreshArtist"}

	srv := newTestServer(t, http.MethodGet, "/api/v1/command/42", want)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetCommand(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetCommand: %v", err)
	}
	if got.Name != "RefreshArtist" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestGetHealth(t *testing.T) {
	t.Parallel()

	want := []arr.HealthCheck{{Type: "warning", Message: "test"}}

	srv := newTestServer(t, http.MethodGet, "/api/v1/health", want)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
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

	srv := newTestServer(t, http.MethodGet, "/api/v1/diskspace", want)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
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

	srv := newTestServer(t, http.MethodGet, "/api/v1/qualityprofile", want)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
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

	want := []arr.RootFolder{{ID: 1, Path: "/music"}}

	srv := newTestServer(t, http.MethodGet, "/api/v1/rootfolder", want)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetRootFolders(context.Background())
	if err != nil {
		t.Fatalf("GetRootFolders: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestGetWantedCutoff(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[lidarr.Album]{
		Page: 1, PageSize: 10, TotalRecords: 1,
		Records: []lidarr.Album{{ID: 10, Title: "OK Computer"}},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/wanted/cutoff?page=1&pageSize=10", want)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetWantedCutoff(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetWantedCutoff: %v", err)
	}
	if got.TotalRecords != 1 {
		t.Errorf("TotalRecords = %d", got.TotalRecords)
	}
}

func TestGetImportListExclusions(t *testing.T) {
	t.Parallel()

	want := []lidarr.ImportListExclusion{{ID: 1, ForeignID: "abc-123", ArtistName: "Test"}}

	srv := newTestServer(t, http.MethodGet, "/api/v1/importlistexclusion", want)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetImportListExclusions(context.Background())
	if err != nil {
		t.Fatalf("GetImportListExclusions: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestEditArtists(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.EditArtists(context.Background(), &lidarr.ArtistEditorResource{ArtistIDs: []int{1, 2}}); err != nil {
		t.Fatalf("EditArtists: %v", err)
	}
}

func TestDeleteArtists(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v1/artist/editor", nil)
	defer srv.Close()

	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteArtists(context.Background(), &lidarr.ArtistEditorResource{ArtistIDs: []int{1}}); err != nil {
		t.Fatalf("DeleteArtists: %v", err)
	}
}

// ========== Notifications ==========.

func TestGetNotifications(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/notification", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetNotifications(context.Background())
	if err != nil {
		t.Fatalf("GetNotifications: %v", err)
	}
	if len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestGetNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/notification/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetNotification(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetNotification: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification", arr.ProviderResource{ID: 1, Name: "test"})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateNotification(context.Background(), &arr.ProviderResource{Name: "test"})
	if err != nil {
		t.Fatalf("CreateNotification: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/notification/1", arr.ProviderResource{ID: 1, Name: "updated"})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateNotification(context.Background(), &arr.ProviderResource{ID: 1, Name: "updated"})
	if err != nil {
		t.Fatalf("UpdateNotification: %v", err)
	}
	if got.Name != "updated" {
		t.Fatalf("Name = %q, want %q", got.Name, "updated")
	}
}

func TestDeleteNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/notification/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteNotification(context.Background(), 1); err != nil {
		t.Fatalf("DeleteNotification: %v", err)
	}
}

func TestGetNotificationSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/notification/schema", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetNotificationSchema(context.Background())
	if err != nil {
		t.Fatalf("GetNotificationSchema: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestTestNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification/test", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.TestNotification(context.Background(), &arr.ProviderResource{ID: 1}); err != nil {
		t.Fatalf("TestNotification: %v", err)
	}
}

func TestTestAllNotifications(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification/testall", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.TestAllNotifications(context.Background()); err != nil {
		t.Fatalf("TestAllNotifications: %v", err)
	}
}

// ========== Download Clients ==========.

func TestGetDownloadClients(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/downloadclient", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetDownloadClients(context.Background())
	if err != nil {
		t.Fatalf("GetDownloadClients: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/downloadclient/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetDownloadClient(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetDownloadClient: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateDownloadClient(context.Background(), &arr.ProviderResource{Name: "test"})
	if err != nil {
		t.Fatalf("CreateDownloadClient: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/downloadclient/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateDownloadClient(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateDownloadClient: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/downloadclient/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteDownloadClient(context.Background(), 1); err != nil {
		t.Fatalf("DeleteDownloadClient: %v", err)
	}
}

func TestGetDownloadClientSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/downloadclient/schema", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetDownloadClientSchema(context.Background())
	if err != nil {
		t.Fatalf("GetDownloadClientSchema: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestTestDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient/test", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.TestDownloadClient(context.Background(), &arr.ProviderResource{ID: 1}); err != nil {
		t.Fatalf("TestDownloadClient: %v", err)
	}
}

func TestTestAllDownloadClients(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient/testall", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.TestAllDownloadClients(context.Background()); err != nil {
		t.Fatalf("TestAllDownloadClients: %v", err)
	}
}

func TestUpdateDownloadClientsBulk(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/downloadclient/bulk", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateDownloadClientsBulk(context.Background(), &arr.ProviderBulkResource{IDs: []int{1}})
	if err != nil {
		t.Fatalf("UpdateDownloadClientsBulk: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestDeleteDownloadClientsBulk(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/downloadclient/bulk", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteDownloadClientsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteDownloadClientsBulk: %v", err)
	}
}

// ========== Indexers ==========.

func TestGetIndexers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetIndexers(context.Background())
	if err != nil {
		t.Fatalf("GetIndexers: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetIndexer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetIndexer(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetIndexer: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateIndexer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexer", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateIndexer(context.Background(), &arr.ProviderResource{Name: "test"})
	if err != nil {
		t.Fatalf("CreateIndexer: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateIndexer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/indexer/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateIndexer(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateIndexer: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteIndexer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/indexer/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteIndexer(context.Background(), 1); err != nil {
		t.Fatalf("DeleteIndexer: %v", err)
	}
}

func TestGetIndexerSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer/schema", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetIndexerSchema(context.Background())
	if err != nil {
		t.Fatalf("GetIndexerSchema: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestTestIndexer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexer/test", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.TestIndexer(context.Background(), &arr.ProviderResource{ID: 1}); err != nil {
		t.Fatalf("TestIndexer: %v", err)
	}
}

func TestTestAllIndexers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexer/testall", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.TestAllIndexers(context.Background()); err != nil {
		t.Fatalf("TestAllIndexers: %v", err)
	}
}

func TestUpdateIndexersBulk(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/indexer/bulk", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateIndexersBulk(context.Background(), &arr.ProviderBulkResource{IDs: []int{1}})
	if err != nil {
		t.Fatalf("UpdateIndexersBulk: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestDeleteIndexersBulk(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/indexer/bulk", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteIndexersBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteIndexersBulk: %v", err)
	}
}

// ========== Import Lists ==========.

func TestGetImportLists(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/importlist", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetImportLists(context.Background())
	if err != nil {
		t.Fatalf("GetImportLists: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetImportList(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/importlist/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetImportList(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetImportList: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateImportList(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/importlist", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateImportList(context.Background(), &arr.ProviderResource{Name: "test"})
	if err != nil {
		t.Fatalf("CreateImportList: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateImportList(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/importlist/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateImportList(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateImportList: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteImportList(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/importlist/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteImportList(context.Background(), 1); err != nil {
		t.Fatalf("DeleteImportList: %v", err)
	}
}

func TestGetImportListSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/importlist/schema", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetImportListSchema(context.Background())
	if err != nil {
		t.Fatalf("GetImportListSchema: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestTestImportList(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/importlist/test", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.TestImportList(context.Background(), &arr.ProviderResource{ID: 1}); err != nil {
		t.Fatalf("TestImportList: %v", err)
	}
}

func TestTestAllImportLists(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/importlist/testall", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.TestAllImportLists(context.Background()); err != nil {
		t.Fatalf("TestAllImportLists: %v", err)
	}
}

func TestUpdateImportListsBulk(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/importlist/bulk", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateImportListsBulk(context.Background(), &arr.ProviderBulkResource{IDs: []int{1}})
	if err != nil {
		t.Fatalf("UpdateImportListsBulk: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestDeleteImportListsBulk(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/importlist/bulk", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteImportListsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteImportListsBulk: %v", err)
	}
}

// ========== Metadata Consumers ==========.

func TestGetMetadataConsumers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/metadata", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetMetadataConsumers(context.Background())
	if err != nil {
		t.Fatalf("GetMetadataConsumers: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetMetadataConsumer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/metadata/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetMetadataConsumer(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMetadataConsumer: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateMetadataConsumer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/metadata", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateMetadataConsumer(context.Background(), &arr.ProviderResource{Name: "test"})
	if err != nil {
		t.Fatalf("CreateMetadataConsumer: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateMetadataConsumer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/metadata/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateMetadataConsumer(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateMetadataConsumer: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteMetadataConsumer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/metadata/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteMetadataConsumer(context.Background(), 1); err != nil {
		t.Fatalf("DeleteMetadataConsumer: %v", err)
	}
}

func TestGetMetadataSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/metadata/schema", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetMetadataSchema(context.Background())
	if err != nil {
		t.Fatalf("GetMetadataSchema: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestTestMetadataConsumer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/metadata/test", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.TestMetadataConsumer(context.Background(), &arr.ProviderResource{ID: 1}); err != nil {
		t.Fatalf("TestMetadataConsumer: %v", err)
	}
}

func TestTestAllMetadataConsumers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/metadata/testall", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.TestAllMetadataConsumers(context.Background()); err != nil {
		t.Fatalf("TestAllMetadataConsumers: %v", err)
	}
}

// ========== Config Endpoints ==========.

func TestGetDownloadClientConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/downloadclient", arr.DownloadClientConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetDownloadClientConfig(context.Background())
	if err != nil {
		t.Fatalf("GetDownloadClientConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateDownloadClientConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/downloadclient/1", arr.DownloadClientConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateDownloadClientConfig(context.Background(), &arr.DownloadClientConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateDownloadClientConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetIndexerConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/indexer", arr.IndexerConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetIndexerConfig(context.Background())
	if err != nil {
		t.Fatalf("GetIndexerConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateIndexerConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/indexer/1", arr.IndexerConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateIndexerConfig(context.Background(), &arr.IndexerConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateIndexerConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetNamingConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/naming", arr.NamingConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetNamingConfig(context.Background())
	if err != nil {
		t.Fatalf("GetNamingConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateNamingConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/naming/1", arr.NamingConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateNamingConfig(context.Background(), &arr.NamingConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateNamingConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetHostConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/host", arr.HostConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetHostConfig(context.Background())
	if err != nil {
		t.Fatalf("GetHostConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateHostConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/host/1", arr.HostConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateHostConfig(context.Background(), &arr.HostConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateHostConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetUIConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/ui", arr.UIConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetUIConfig(context.Background())
	if err != nil {
		t.Fatalf("GetUIConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateUIConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/ui/1", arr.UIConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateUIConfig(context.Background(), &arr.UIConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateUIConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetMediaManagementConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/mediamanagement", arr.MediaManagementConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetMediaManagementConfig(context.Background())
	if err != nil {
		t.Fatalf("GetMediaManagementConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateMediaManagementConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/mediamanagement/1", arr.MediaManagementConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateMediaManagementConfig(context.Background(), &arr.MediaManagementConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateMediaManagementConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetImportListConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/importlist", lidarr.ImportListConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetImportListConfig(context.Background())
	if err != nil {
		t.Fatalf("GetImportListConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateImportListConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/importlist/1", lidarr.ImportListConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateImportListConfig(context.Background(), &lidarr.ImportListConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateImportListConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

// ========== Quality Profiles ==========.

func TestGetQualityProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/qualityprofile/1", arr.QualityProfile{ID: 1, Name: "Flac"})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetQualityProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetQualityProfile: %v", err)
	}
	if got.Name != "Flac" {
		t.Fatalf("Name = %q, want %q", got.Name, "Flac")
	}
}

func TestCreateQualityProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/qualityprofile", arr.QualityProfile{ID: 1, Name: "Flac"})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateQualityProfile(context.Background(), &arr.QualityProfile{Name: "Flac"})
	if err != nil {
		t.Fatalf("CreateQualityProfile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateQualityProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/qualityprofile/1", arr.QualityProfile{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateQualityProfile(context.Background(), &arr.QualityProfile{ID: 1})
	if err != nil {
		t.Fatalf("UpdateQualityProfile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteQualityProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/qualityprofile/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteQualityProfile(context.Background(), 1); err != nil {
		t.Fatalf("DeleteQualityProfile: %v", err)
	}
}

func TestGetQualityProfileSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/qualityprofile/schema", arr.QualityProfile{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetQualityProfileSchema(context.Background())
	if err != nil {
		t.Fatalf("GetQualityProfileSchema: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

// ========== Quality Definitions ==========.

func TestGetQualityDefinitions(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/qualitydefinition", []arr.QualityDefinitionResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetQualityDefinitions(context.Background())
	if err != nil {
		t.Fatalf("GetQualityDefinitions: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetQualityDefinition(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/qualitydefinition/1", arr.QualityDefinitionResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetQualityDefinition(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetQualityDefinition: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateQualityDefinition(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/qualitydefinition/1", arr.QualityDefinitionResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateQualityDefinition(context.Background(), &arr.QualityDefinitionResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateQualityDefinition: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateQualityDefinitions(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/qualitydefinition/update", []arr.QualityDefinitionResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateQualityDefinitions(context.Background(), []arr.QualityDefinitionResource{{ID: 1}})
	if err != nil {
		t.Fatalf("UpdateQualityDefinitions: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

// ========== Metadata Profiles ==========.

func TestGetMetadataProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/metadataprofile/1", lidarr.MetadataProfile{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetMetadataProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMetadataProfile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateMetadataProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/metadataprofile", lidarr.MetadataProfile{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateMetadataProfile(context.Background(), &lidarr.MetadataProfile{Name: "test"})
	if err != nil {
		t.Fatalf("CreateMetadataProfile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateMetadataProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/metadataprofile/1", lidarr.MetadataProfile{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateMetadataProfile(context.Background(), &lidarr.MetadataProfile{ID: 1})
	if err != nil {
		t.Fatalf("UpdateMetadataProfile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteMetadataProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/metadataprofile/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteMetadataProfile(context.Background(), 1); err != nil {
		t.Fatalf("DeleteMetadataProfile: %v", err)
	}
}

// ========== Tags Extended ==========.

func TestGetTag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/tag/1", arr.Tag{ID: 1, Label: "test"})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetTag(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTag: %v", err)
	}
	if got.Label != "test" {
		t.Fatalf("Label = %q, want %q", got.Label, "test")
	}
}

func TestUpdateTag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/tag/1", arr.Tag{ID: 1, Label: "updated"})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateTag(context.Background(), &arr.Tag{ID: 1, Label: "updated"})
	if err != nil {
		t.Fatalf("UpdateTag: %v", err)
	}
	if got.Label != "updated" {
		t.Fatalf("Label = %q, want %q", got.Label, "updated")
	}
}

func TestDeleteTag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/tag/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteTag(context.Background(), 1); err != nil {
		t.Fatalf("DeleteTag: %v", err)
	}
}

func TestGetTagDetails(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/tag/detail", []arr.TagDetail{{ID: 1, Label: "test"}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetTagDetails(context.Background())
	if err != nil {
		t.Fatalf("GetTagDetails: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetTagDetail(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/tag/detail/1", arr.TagDetail{ID: 1, Label: "test"})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetTagDetail(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTagDetail: %v", err)
	}
	if got.Label != "test" {
		t.Fatalf("Label = %q, want %q", got.Label, "test")
	}
}

// ========== Root Folders Extended ==========.

func TestGetRootFolder(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/rootfolder/1", arr.RootFolder{ID: 1, Path: "/music"})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetRootFolder(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetRootFolder: %v", err)
	}
	if got.Path != "/music" {
		t.Fatalf("Path = %q, want %q", got.Path, "/music")
	}
}

func TestCreateRootFolder(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/rootfolder", arr.RootFolder{ID: 1, Path: "/music"})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateRootFolder(context.Background(), &arr.RootFolder{Path: "/music"})
	if err != nil {
		t.Fatalf("CreateRootFolder: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateRootFolder(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/rootfolder/1", arr.RootFolder{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateRootFolder(context.Background(), &arr.RootFolder{ID: 1})
	if err != nil {
		t.Fatalf("UpdateRootFolder: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteRootFolder(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/rootfolder/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteRootFolder(context.Background(), 1); err != nil {
		t.Fatalf("DeleteRootFolder: %v", err)
	}
}

// ========== Custom Filters ==========.

func TestGetCustomFilters(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/customfilter", []arr.CustomFilterResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetCustomFilters(context.Background())
	if err != nil {
		t.Fatalf("GetCustomFilters: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetCustomFilter(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/customfilter/1", arr.CustomFilterResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetCustomFilter(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCustomFilter: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateCustomFilter(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/customfilter", arr.CustomFilterResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateCustomFilter(context.Background(), &arr.CustomFilterResource{})
	if err != nil {
		t.Fatalf("CreateCustomFilter: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateCustomFilter(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/customfilter/1", arr.CustomFilterResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateCustomFilter(context.Background(), &arr.CustomFilterResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateCustomFilter: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteCustomFilter(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/customfilter/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteCustomFilter(context.Background(), 1); err != nil {
		t.Fatalf("DeleteCustomFilter: %v", err)
	}
}

// ========== Custom Formats ==========.

func TestGetCustomFormats(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/customformat", []arr.CustomFormatResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetCustomFormats(context.Background())
	if err != nil {
		t.Fatalf("GetCustomFormats: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetCustomFormat(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/customformat/1", arr.CustomFormatResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetCustomFormat(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCustomFormat: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateCustomFormat(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/customformat", arr.CustomFormatResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateCustomFormat(context.Background(), &arr.CustomFormatResource{})
	if err != nil {
		t.Fatalf("CreateCustomFormat: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateCustomFormat(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/customformat/1", arr.CustomFormatResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateCustomFormat(context.Background(), &arr.CustomFormatResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateCustomFormat: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteCustomFormat(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/customformat/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteCustomFormat(context.Background(), 1); err != nil {
		t.Fatalf("DeleteCustomFormat: %v", err)
	}
}

func TestGetCustomFormatSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/customformat/schema", []arr.CustomFormatResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetCustomFormatSchema(context.Background())
	if err != nil {
		t.Fatalf("GetCustomFormatSchema: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestDeleteCustomFormatsBulk(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/customformat/bulk", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteCustomFormatsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteCustomFormatsBulk: %v", err)
	}
}

// ========== Delay Profiles ==========.

func TestGetDelayProfiles(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/delayprofile", []arr.DelayProfileResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetDelayProfiles(context.Background())
	if err != nil {
		t.Fatalf("GetDelayProfiles: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetDelayProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/delayprofile/1", arr.DelayProfileResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetDelayProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetDelayProfile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateDelayProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/delayprofile", arr.DelayProfileResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateDelayProfile(context.Background(), &arr.DelayProfileResource{})
	if err != nil {
		t.Fatalf("CreateDelayProfile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateDelayProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/delayprofile/1", arr.DelayProfileResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateDelayProfile(context.Background(), &arr.DelayProfileResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateDelayProfile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteDelayProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/delayprofile/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteDelayProfile(context.Background(), 1); err != nil {
		t.Fatalf("DeleteDelayProfile: %v", err)
	}
}

func TestReorderDelayProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/delayprofile/reorder/1?after=2", []arr.DelayProfileResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.ReorderDelayProfile(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("ReorderDelayProfile: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

// ========== Release Profiles ==========.

func TestGetReleaseProfiles(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/releaseprofile", []arr.ReleaseProfileResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetReleaseProfiles(context.Background())
	if err != nil {
		t.Fatalf("GetReleaseProfiles: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetReleaseProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/releaseprofile/1", arr.ReleaseProfileResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetReleaseProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetReleaseProfile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateReleaseProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/releaseprofile", arr.ReleaseProfileResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateReleaseProfile(context.Background(), &arr.ReleaseProfileResource{})
	if err != nil {
		t.Fatalf("CreateReleaseProfile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateReleaseProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/releaseprofile/1", arr.ReleaseProfileResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateReleaseProfile(context.Background(), &arr.ReleaseProfileResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateReleaseProfile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteReleaseProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/releaseprofile/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteReleaseProfile(context.Background(), 1); err != nil {
		t.Fatalf("DeleteReleaseProfile: %v", err)
	}
}

// ========== Remote Path Mappings ==========.

func TestGetRemotePathMappings(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/remotepathmapping", []arr.RemotePathMappingResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetRemotePathMappings(context.Background())
	if err != nil {
		t.Fatalf("GetRemotePathMappings: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetRemotePathMapping(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/remotepathmapping/1", arr.RemotePathMappingResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetRemotePathMapping(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetRemotePathMapping: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateRemotePathMapping(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/remotepathmapping", arr.RemotePathMappingResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateRemotePathMapping(context.Background(), &arr.RemotePathMappingResource{})
	if err != nil {
		t.Fatalf("CreateRemotePathMapping: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateRemotePathMapping(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/remotepathmapping/1", arr.RemotePathMappingResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateRemotePathMapping(context.Background(), &arr.RemotePathMappingResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateRemotePathMapping: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteRemotePathMapping(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/remotepathmapping/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteRemotePathMapping(context.Background(), 1); err != nil {
		t.Fatalf("DeleteRemotePathMapping: %v", err)
	}
}

// ========== Auto Tagging ==========.

func TestGetAutoTagging(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/autotagging", []arr.AutoTaggingResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetAutoTagging(context.Background())
	if err != nil {
		t.Fatalf("GetAutoTagging: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetAutoTag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/autotagging/1", arr.AutoTaggingResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetAutoTag(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetAutoTag: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateAutoTag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/autotagging", arr.AutoTaggingResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateAutoTag(context.Background(), &arr.AutoTaggingResource{})
	if err != nil {
		t.Fatalf("CreateAutoTag: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateAutoTag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/autotagging/1", arr.AutoTaggingResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateAutoTag(context.Background(), &arr.AutoTaggingResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateAutoTag: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteAutoTag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/autotagging/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteAutoTag(context.Background(), 1); err != nil {
		t.Fatalf("DeleteAutoTag: %v", err)
	}
}

func TestGetAutoTagSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/autotagging/schema", []arr.AutoTaggingResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetAutoTagSchema(context.Background())
	if err != nil {
		t.Fatalf("GetAutoTagSchema: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

// ========== Import List Exclusions Extended ==========.

func TestGetImportListExclusion(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/importlistexclusion/1", lidarr.ImportListExclusion{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetImportListExclusion(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetImportListExclusion: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestCreateImportListExclusion(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/importlistexclusion", lidarr.ImportListExclusion{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.CreateImportListExclusion(context.Background(), &lidarr.ImportListExclusion{})
	if err != nil {
		t.Fatalf("CreateImportListExclusion: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateImportListExclusion(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/importlistexclusion/1", lidarr.ImportListExclusion{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateImportListExclusion(context.Background(), &lidarr.ImportListExclusion{ID: 1})
	if err != nil {
		t.Fatalf("UpdateImportListExclusion: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestDeleteImportListExclusion(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/importlistexclusion/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteImportListExclusion(context.Background(), 1); err != nil {
		t.Fatalf("DeleteImportListExclusion: %v", err)
	}
}

func TestGetImportListExclusionsPaged(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/importlistexclusion/paged?page=1&pageSize=10",
		arr.PagingResource[lidarr.ImportListExclusion]{Page: 1, PageSize: 10, TotalRecords: 1, Records: []lidarr.ImportListExclusion{{ID: 1}}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetImportListExclusionsPaged(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetImportListExclusionsPaged: %v", err)
	}
	if got.TotalRecords != 1 {
		t.Fatalf("TotalRecords = %d, want 1", got.TotalRecords)
	}
}

func TestDeleteImportListExclusionsBulk(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/importlistexclusion/bulk", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteImportListExclusionsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteImportListExclusionsBulk: %v", err)
	}
}

// ========== Blocklist ==========.

func TestGetBlocklist(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/blocklist?page=1&pageSize=10",
		arr.PagingResource[arr.BlocklistResource]{Page: 1, PageSize: 10, TotalRecords: 1, Records: []arr.BlocklistResource{{ID: 1}}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetBlocklist(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetBlocklist: %v", err)
	}
	if got.TotalRecords != 1 {
		t.Fatalf("TotalRecords = %d, want 1", got.TotalRecords)
	}
}

func TestDeleteBlocklistItem(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/blocklist/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteBlocklistItem(context.Background(), 1); err != nil {
		t.Fatalf("DeleteBlocklistItem: %v", err)
	}
}

func TestDeleteBlocklistBulk(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/blocklist/bulk", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteBlocklistBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteBlocklistBulk: %v", err)
	}
}

// ========== Queue Extended ==========.

func TestDeleteQueueItems(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/queue/bulk?removeFromClient=true&blocklist=false", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteQueueItems(context.Background(), []int{1, 2}, true, false); err != nil {
		t.Fatalf("DeleteQueueItems: %v", err)
	}
}

func TestGrabQueueItem(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/queue/grab/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.GrabQueueItem(context.Background(), 1); err != nil {
		t.Fatalf("GrabQueueItem: %v", err)
	}
}

func TestGrabQueueItems(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/queue/grab/bulk", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.GrabQueueItems(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("GrabQueueItems: %v", err)
	}
}

func TestGetQueueDetails(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/queue/details", []arr.QueueRecord{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetQueueDetails(context.Background())
	if err != nil {
		t.Fatalf("GetQueueDetails: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetQueueStatus(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/queue/status", arr.QueueStatusResource{})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	_, err := c.GetQueueStatus(context.Background())
	if err != nil {
		t.Fatalf("GetQueueStatus: %v", err)
	}
}

// ========== History Extended ==========.

func TestGetHistoryArtist(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/history/artist?artistId=1", []lidarr.HistoryRecord{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetHistoryArtist(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetHistoryArtist: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetHistorySince(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/history/since?date=2024-01-01", []lidarr.HistoryRecord{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetHistorySince(context.Background(), "2024-01-01")
	if err != nil {
		t.Fatalf("GetHistorySince: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestMarkHistoryFailed(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/history/failed/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.MarkHistoryFailed(context.Background(), 1); err != nil {
		t.Fatalf("MarkHistoryFailed: %v", err)
	}
}

// ========== Releases ==========.

func TestSearchReleases(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/release?albumId=1", []arr.ReleaseResource{{GUID: "abc"}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.SearchReleases(context.Background(), 1)
	if err != nil {
		t.Fatalf("SearchReleases: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGrabRelease(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/release", arr.ReleaseResource{GUID: "abc"})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GrabRelease(context.Background(), &arr.ReleaseResource{GUID: "abc"})
	if err != nil {
		t.Fatalf("GrabRelease: %v", err)
	}
	if got.GUID != "abc" {
		t.Fatalf("GUID = %q, want %q", got.GUID, "abc")
	}
}

func TestPushRelease(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/release/push", []arr.ReleaseResource{{GUID: "abc"}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.PushRelease(context.Background(), &arr.ReleasePushResource{Title: "test"})
	if err != nil {
		t.Fatalf("PushRelease: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

// ========== Rename ==========.

func TestGetRenameList(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/rename?artistId=1", []lidarr.RenameTrackResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetRenameList(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetRenameList: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

// ========== Retag ==========.

func TestGetRetag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/retag?artistId=1", []lidarr.RetagResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetRetag(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetRetag: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

// ========== Manual Import ==========.

func TestGetManualImport(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/manualimport?folder=%2Fmusic", []arr.ManualImportResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetManualImport(context.Background(), "/music")
	if err != nil {
		t.Fatalf("GetManualImport: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestProcessManualImport(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/manualimport", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.ProcessManualImport(context.Background(), []arr.ManualImportReprocessResource{{ID: 1}}); err != nil {
		t.Fatalf("ProcessManualImport: %v", err)
	}
}

// ========== Backups ==========.

func TestGetBackups(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/backup", []arr.Backup{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetBackups(context.Background())
	if err != nil {
		t.Fatalf("GetBackups: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestDeleteBackup(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/system/backup/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteBackup(context.Background(), 1); err != nil {
		t.Fatalf("DeleteBackup: %v", err)
	}
}

func TestRestoreBackup(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/system/backup/restore/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.RestoreBackup(context.Background(), 1); err != nil {
		t.Fatalf("RestoreBackup: %v", err)
	}
}

// ========== Logs ==========.

func TestGetLogs(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/log?page=1&pageSize=10",
		arr.PagingResource[arr.LogRecord]{Page: 1, PageSize: 10, TotalRecords: 1, Records: []arr.LogRecord{{ID: 1}}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetLogs(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetLogs: %v", err)
	}
	if got.TotalRecords != 1 {
		t.Fatalf("TotalRecords = %d, want 1", got.TotalRecords)
	}
}

func TestGetLogFiles(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/log/file", []arr.LogFileResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetLogFiles(context.Background())
	if err != nil {
		t.Fatalf("GetLogFiles: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetUpdateLogFiles(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/log/file/update", []arr.LogFileResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetUpdateLogFiles(context.Background())
	if err != nil {
		t.Fatalf("GetUpdateLogFiles: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

// ========== System ==========.

func TestGetTasks(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/task", []arr.TaskResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetTasks(context.Background())
	if err != nil {
		t.Fatalf("GetTasks: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetTask(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/task/1", arr.TaskResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetTask(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetUpdates(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/update", []arr.UpdateResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetUpdates(context.Background())
	if err != nil {
		t.Fatalf("GetUpdates: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetSystemRoutes(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/routes", []arr.SystemRouteResource{{Path: "/api/v1/test"}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetSystemRoutes(context.Background())
	if err != nil {
		t.Fatalf("GetSystemRoutes: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestShutdown(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/system/shutdown", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestRestart(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/system/restart", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.Restart(context.Background()); err != nil {
		t.Fatalf("Restart: %v", err)
	}
}

func TestDeleteCommand(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/command/1", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.DeleteCommand(context.Background(), 1); err != nil {
		t.Fatalf("DeleteCommand: %v", err)
	}
}

func TestGetLanguages(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/language", []arr.LanguageResource{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetLanguages(context.Background())
	if err != nil {
		t.Fatalf("GetLanguages: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetLanguage(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/language/1", arr.LanguageResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetLanguage(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetLanguage: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetLocalization(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/localization", lidarr.LocalizationResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetLocalization(context.Background())
	if err != nil {
		t.Fatalf("GetLocalization: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestPing(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/ping", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

// ========== File System ==========.

func TestBrowseFileSystem(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/filesystem?path=%2Fmusic&includeFiles=true",
		lidarr.FileSystemResource{})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	_, err := c.BrowseFileSystem(context.Background(), "/music", true)
	if err != nil {
		t.Fatalf("BrowseFileSystem: %v", err)
	}
}

func TestGetFileSystemType(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/filesystem/type?path=%2Fmusic", "ext4")
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetFileSystemType(context.Background(), "/music")
	if err != nil {
		t.Fatalf("GetFileSystemType: %v", err)
	}
	if got != "ext4" {
		t.Fatalf("type = %q, want %q", got, "ext4")
	}
}

// ========== Track File Update ==========.

func TestUpdateTrackFile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/trackfile/1", lidarr.TrackFile{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateTrackFile(context.Background(), &lidarr.TrackFile{ID: 1})
	if err != nil {
		t.Fatalf("UpdateTrackFile: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestEditTrackFiles(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/trackfile/editor", []lidarr.TrackFile{{ID: 1}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.EditTrackFiles(context.Background(), &lidarr.TrackFileEditorResource{TrackFileIDs: []int{1}})
	if err != nil {
		t.Fatalf("EditTrackFiles: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

// ========== Album Studio ==========.

func TestAlbumStudio(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/albumstudio", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.AlbumStudio(context.Background(), &lidarr.AlbumStudioResource{}); err != nil {
		t.Fatalf("AlbumStudio: %v", err)
	}
}

// ========== Calendar By ID ==========.

func TestGetCalendarByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/calendar/1", lidarr.Album{ID: 1, Title: "OK Computer"})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetCalendarByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCalendarByID: %v", err)
	}
	if got.Title != "OK Computer" {
		t.Fatalf("Title = %q, want %q", got.Title, "OK Computer")
	}
}

// ========== Wanted By ID ==========.

func TestGetWantedMissingByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/wanted/missing/1", lidarr.Album{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetWantedMissingByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetWantedMissingByID: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetWantedCutoffByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/wanted/cutoff/1", lidarr.Album{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetWantedCutoffByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetWantedCutoffByID: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

// ========== Config By ID ==========.

func TestGetDownloadClientConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/downloadclient/1", arr.DownloadClientConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetDownloadClientConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetDownloadClientConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetHostConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/host/1", arr.HostConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetHostConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetHostConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetIndexerConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/indexer/1", arr.IndexerConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetIndexerConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetIndexerConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetMediaManagementConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/mediamanagement/1", arr.MediaManagementConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetMediaManagementConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMediaManagementConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetNamingConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/naming/1", arr.NamingConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetNamingConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetNamingConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetUIConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/ui/1", arr.UIConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetUIConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUIConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

// ========== Metadata Provider Config ==========.

func TestGetMetadataProviderConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/metadataprovider", lidarr.MetadataProviderConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetMetadataProviderConfig(context.Background())
	if err != nil {
		t.Fatalf("GetMetadataProviderConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestGetMetadataProviderConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/metadataprovider/1", lidarr.MetadataProviderConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetMetadataProviderConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMetadataProviderConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateMetadataProviderConfig(t *testing.T) {
	t.Parallel()
	in := &lidarr.MetadataProviderConfigResource{ID: 1, MetadataSource: "test"}
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/metadataprovider/1", *in)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.UpdateMetadataProviderConfig(context.Background(), in)
	if err != nil {
		t.Fatalf("UpdateMetadataProviderConfig: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

// ========== Indexer Flags ==========.

func TestGetIndexerFlags(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexerflag", []arr.IndexerFlagResource{{ID: 1, Name: "flag1"}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetIndexerFlags(context.Background())
	if err != nil {
		t.Fatalf("GetIndexerFlags: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

// ========== Metadata Profile Schema ==========.

func TestGetMetadataProfileSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/metadataprofile/schema", lidarr.MetadataProfile{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetMetadataProfileSchema(context.Background())
	if err != nil {
		t.Fatalf("GetMetadataProfileSchema: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

// ========== Naming Examples ==========.

func TestGetNamingExamples(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/naming/examples", arr.NamingConfigResource{ID: 1})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetNamingExamples(context.Background())
	if err != nil {
		t.Fatalf("GetNamingExamples: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("ID = %d, want 1", got.ID)
	}
}

// ========== System Routes Duplicate ==========.

func TestGetSystemRoutesDuplicate(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/routes/duplicate", []arr.SystemRouteResource{{Path: "/dup"}})
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetSystemRoutesDuplicate(context.Background())
	if err != nil {
		t.Fatalf("GetSystemRoutesDuplicate: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

// ========== Update Log File Content ==========.

func TestGetUpdateLogFileContent(t *testing.T) {
	t.Parallel()
	srv := newRawTestServer(t, http.MethodGet, "/api/v1/log/file/update/test.log", "update log content")
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetUpdateLogFileContent(context.Background(), "test.log")
	if err != nil {
		t.Fatalf("GetUpdateLogFileContent: %v", err)
	}
	if got != "update log content" {
		t.Fatalf("got %q, want %q", got, "update log content")
	}
}

// ========== Provider Actions ==========.

func TestDownloadClientAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient/action/testAction", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	err := c.DownloadClientAction(context.Background(), "testAction", &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("DownloadClientAction: %v", err)
	}
}

func TestImportListAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/importlist/action/testAction", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	err := c.ImportListAction(context.Background(), "testAction", &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("ImportListAction: %v", err)
	}
}

func TestIndexerAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexer/action/testAction", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	err := c.IndexerAction(context.Background(), "testAction", &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("IndexerAction: %v", err)
	}
}

func TestMetadataConsumerAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/metadata/action/testAction", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	err := c.MetadataConsumerAction(context.Background(), "testAction", &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("MetadataConsumerAction: %v", err)
	}
}

func TestNotificationAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification/action/testAction", nil)
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	err := c.NotificationAction(context.Background(), "testAction", &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("NotificationAction: %v", err)
	}
}

// ========== Error Handling ==========.

func TestNotificationError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	_, err := c.GetNotifications(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *arr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
}

func TestGetLogFileContent(t *testing.T) {
	t.Parallel()
	srv := newRawTestServer(t, http.MethodGet, "/api/v1/log/file/lidarr.txt", "log line 1\nlog line 2")
	defer srv.Close()
	c, _ := lidarr.New(srv.URL, "test-key")
	got, err := c.GetLogFileContent(context.Background(), "lidarr.txt")
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
	c, _ := lidarr.New(srv.URL, "test-key")
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
	c, _ := lidarr.New(srv.URL, "test-key")
	if err := c.UploadBackup(context.Background(), "backup.zip", strings.NewReader("fake-backup-data")); err != nil {
		t.Fatalf("UploadBackup: %v", err)
	}
}
