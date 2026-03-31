package lidarr_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/arr"
	"github.com/lusoris/goenvoy/arr/lidarr"
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
