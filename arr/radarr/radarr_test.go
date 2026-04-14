package radarr_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golusoris/goenvoy/arr/radarr"
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
		c, err := radarr.New("http://localhost:7878", "test-key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c == nil {
			t.Fatal("expected non-nil client")
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		t.Parallel()
		_, err := radarr.New("://bad", "test-key")
		if err == nil {
			t.Fatal("expected error for invalid URL")
		}
	})
}

func TestGetAllMovies(t *testing.T) {
	t.Parallel()

	want := []radarr.Movie{
		{ID: 1, Title: "Inception", TmdbID: 27205},
		{ID: 2, Title: "The Matrix", TmdbID: 603},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v3/movie", want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetAllMovies(context.Background())
	if err != nil {
		t.Fatalf("GetAllMovies: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Title != "Inception" {
		t.Errorf("Title = %q, want %q", got[0].Title, "Inception")
	}
}

func TestGetMovie(t *testing.T) {
	t.Parallel()

	want := radarr.Movie{ID: 1, Title: "Inception"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/movie/1", want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetMovie(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMovie: %v", err)
	}
	if got.Title != "Inception" {
		t.Errorf("Title = %q, want %q", got.Title, "Inception")
	}
}

func TestAddMovie(t *testing.T) {
	t.Parallel()

	want := radarr.Movie{ID: 3, Title: "New Movie", TmdbID: 99999}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var body radarr.Movie
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Title != "New Movie" {
			t.Errorf("Title = %q, want %q", body.Title, "New Movie")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.AddMovie(context.Background(), &radarr.Movie{
		Title:  "New Movie",
		TmdbID: 99999,
	})
	if err != nil {
		t.Fatalf("AddMovie: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("ID = %d, want 3", got.ID)
	}
}

func TestDeleteMovie(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete,
		"/api/v3/movie/1?deleteFiles=true&addImportExclusion=false",
		nil)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteMovie(context.Background(), 1, true, false); err != nil {
		t.Fatalf("DeleteMovie: %v", err)
	}
}

func TestLookupMovie(t *testing.T) {
	t.Parallel()

	want := []radarr.Movie{{ID: 0, Title: "Inception", TmdbID: 27205}}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/movie/lookup?term=inception",
		want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.LookupMovie(context.Background(), "inception")
	if err != nil {
		t.Fatalf("LookupMovie: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestLookupMovieByTmdbID(t *testing.T) {
	t.Parallel()

	want := radarr.Movie{ID: 1, Title: "Inception", TmdbID: 27205}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/movie/lookup/tmdb?tmdbId=27205",
		want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.LookupMovieByTmdbID(context.Background(), 27205)
	if err != nil {
		t.Fatalf("LookupMovieByTmdbID: %v", err)
	}
	if got.TmdbID != 27205 {
		t.Errorf("TmdbID = %d, want 27205", got.TmdbID)
	}
}

func TestGetMovieFiles(t *testing.T) {
	t.Parallel()

	want := []radarr.MovieFile{
		{ID: 100, MovieID: 1, RelativePath: "Inception.2010.mkv", Size: 2147483648},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/moviefile?movieId=1",
		want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetMovieFiles(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMovieFiles: %v", err)
	}
	if got[0].Size != 2147483648 {
		t.Errorf("Size = %d, want 2147483648", got[0].Size)
	}
}

func TestGetCollections(t *testing.T) {
	t.Parallel()

	want := []radarr.Collection{
		{ID: 1, Title: "The Matrix Collection", TmdbID: 2344},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v3/collection", want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetCollections(context.Background())
	if err != nil {
		t.Fatalf("GetCollections: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Title != "The Matrix Collection" {
		t.Errorf("Title = %q, want %q", got[0].Title, "The Matrix Collection")
	}
}

func TestGetCredits(t *testing.T) {
	t.Parallel()

	want := []radarr.Credit{
		{ID: 1, PersonName: "Christopher Nolan", Type: "crew"},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/credit?movieId=1",
		want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetCredits(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCredits: %v", err)
	}
	if got[0].PersonName != "Christopher Nolan" {
		t.Errorf("PersonName = %q, want %q", got[0].PersonName, "Christopher Nolan")
	}
}

func TestSendCommand(t *testing.T) {
	t.Parallel()

	want := arr.CommandResponse{ID: 42, Name: "RefreshMovie"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var cmd arr.CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if cmd.Name != "RefreshMovie" {
			t.Errorf("Name = %q, want RefreshMovie", cmd.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.SendCommand(context.Background(), arr.CommandRequest{Name: "RefreshMovie"})
	if err != nil {
		t.Fatalf("SendCommand: %v", err)
	}
	if got.ID != 42 {
		t.Errorf("ID = %d, want 42", got.ID)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	want := radarr.ParseResult{
		Title: "Inception.2010.1080p",
		ParsedMovieInfo: &radarr.ParsedMovieInfo{
			MovieTitle: "Inception",
			Year:       2010,
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/parse?title=Inception.2010.1080p",
		want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.Parse(context.Background(), "Inception.2010.1080p")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got.ParsedMovieInfo.MovieTitle != "Inception" {
		t.Errorf("MovieTitle = %q, want %q", got.ParsedMovieInfo.MovieTitle, "Inception")
	}
}

func TestGetSystemStatus(t *testing.T) {
	t.Parallel()

	want := arr.StatusResponse{AppName: "Radarr", Version: "5.0.0"}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/system/status",
		want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetSystemStatus(context.Background())
	if err != nil {
		t.Fatalf("GetSystemStatus: %v", err)
	}
	if got.AppName != "Radarr" {
		t.Errorf("AppName = %q, want %q", got.AppName, "Radarr")
	}
}

func TestGetQueue(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[arr.QueueRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []arr.QueueRecord{
			{ID: 1, Title: "Inception (2010)"},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/queue?page=1&pageSize=10",
		want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
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

	want := []arr.Tag{{ID: 1, Label: "4k"}, {ID: 2, Label: "animation"}}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/tag",
		want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
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

	c, err := radarr.New(srv.URL, "test-key")
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

	c, err := radarr.New(srv.URL, "bad-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = c.GetAllMovies(context.Background())
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

	want := arr.PagingResource[radarr.HistoryRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []radarr.HistoryRecord{
			{ID: 5, MovieID: 1, EventType: "grabbed"},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v3/history?page=1&pageSize=10",
		want)
	defer srv.Close()

	c, err := radarr.New(srv.URL, "test-key")
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

	c, err := radarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteQueueItem(context.Background(), 5, true, false); err != nil {
		t.Fatalf("DeleteQueueItem: %v", err)
	}
}

func TestUpdateMovie(t *testing.T) {
	t.Parallel()

	want := radarr.Movie{ID: 1, Title: "Updated", TmdbID: 27205}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.RequestURI() != "/api/v3/movie/1?moveFiles=true" {
			t.Errorf("path = %q", r.URL.RequestURI())
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.UpdateMovie(context.Background(), &radarr.Movie{ID: 1, Title: "Updated"}, true)
	if err != nil {
		t.Fatalf("UpdateMovie: %v", err)
	}
	if got.Title != "Updated" {
		t.Errorf("Title = %q", got.Title)
	}
}

func TestLookupMovieByImdbID(t *testing.T) {
	t.Parallel()

	want := radarr.Movie{ID: 1, Title: "Inception", ImdbID: "tt1375666"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/movie/lookup/imdb?imdbId=tt1375666", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.LookupMovieByImdbID(context.Background(), "tt1375666")
	if err != nil {
		t.Fatalf("LookupMovieByImdbID: %v", err)
	}
	if got.ImdbID != "tt1375666" {
		t.Errorf("ImdbID = %q", got.ImdbID)
	}
}

func TestGetMovieFile(t *testing.T) {
	t.Parallel()

	want := radarr.MovieFile{ID: 100, MovieID: 1, Size: 2147483648}

	srv := newTestServer(t, http.MethodGet, "/api/v3/moviefile/100", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetMovieFile(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetMovieFile: %v", err)
	}
	if got.Size != 2147483648 {
		t.Errorf("Size = %d", got.Size)
	}
}

func TestDeleteMovieFile(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/moviefile/100", nil)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.DeleteMovieFile(context.Background(), 100); err != nil {
		t.Fatalf("DeleteMovieFile: %v", err)
	}
}

func TestDeleteMovieFiles(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/moviefile/bulk", nil)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.DeleteMovieFiles(context.Background(), []int{1, 2, 3}); err != nil {
		t.Fatalf("DeleteMovieFiles: %v", err)
	}
}

func TestGetCollection(t *testing.T) {
	t.Parallel()

	want := radarr.Collection{ID: 1, Title: "The Matrix Collection"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/collection/1", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetCollection(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCollection: %v", err)
	}
	if got.Title != "The Matrix Collection" {
		t.Errorf("Title = %q", got.Title)
	}
}

func TestUpdateCollection(t *testing.T) {
	t.Parallel()

	want := radarr.Collection{ID: 1, Title: "Updated Collection"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.UpdateCollection(context.Background(), &radarr.Collection{ID: 1, Title: "Updated Collection"})
	if err != nil {
		t.Fatalf("UpdateCollection: %v", err)
	}
	if got.Title != "Updated Collection" {
		t.Errorf("Title = %q", got.Title)
	}
}

func TestGetCalendar(t *testing.T) {
	t.Parallel()

	want := []radarr.Movie{{ID: 1, Title: "Upcoming"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/calendar?start=2026-01-01&end=2026-01-31&unmonitored=false", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetCalendar(context.Background(), "2026-01-01", "2026-01-31", false)
	if err != nil {
		t.Fatalf("GetCalendar: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d, want 1", len(got))
	}
}

func TestGetCommands(t *testing.T) {
	t.Parallel()

	want := []arr.CommandResponse{{ID: 1, Name: "RefreshMovie"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/command", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetCommands(context.Background())
	if err != nil {
		t.Fatalf("GetCommands: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d, want 1", len(got))
	}
}

func TestGetCommand(t *testing.T) {
	t.Parallel()

	want := arr.CommandResponse{ID: 42, Name: "RefreshMovie"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/command/42", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetCommand(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetCommand: %v", err)
	}
	if got.Name != "RefreshMovie" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestGetHealth(t *testing.T) {
	t.Parallel()

	want := []arr.HealthCheck{{Type: "warning", Message: "test"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/health", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d, want 1", len(got))
	}
}

func TestGetDiskSpace(t *testing.T) {
	t.Parallel()

	want := []arr.DiskSpace{{Path: "/data", FreeSpace: 1000}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/diskspace", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
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

	c, _ := radarr.New(srv.URL, "test-key")
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

	want := []arr.RootFolder{{ID: 1, Path: "/movies"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/rootfolder", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetRootFolders(context.Background())
	if err != nil {
		t.Fatalf("GetRootFolders: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestGetImportListExclusions(t *testing.T) {
	t.Parallel()

	want := []radarr.ImportListExclusion{{ID: 1, TmdbID: 550, MovieTitle: "Fight Club"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/exclusions", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetImportListExclusions(context.Background())
	if err != nil {
		t.Fatalf("GetImportListExclusions: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d", len(got))
	}
}

func TestEditMovies(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.EditMovies(context.Background(), &radarr.MovieEditorResource{MovieIDs: []int{1, 2}}); err != nil {
		t.Fatalf("EditMovies: %v", err)
	}
}

func TestDeleteMovies(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/movie/editor", nil)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.DeleteMovies(context.Background(), &radarr.MovieEditorResource{MovieIDs: []int{1}}); err != nil {
		t.Fatalf("DeleteMovies: %v", err)
	}
}

func TestDeleteCommand(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v3/command/1", nil)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.DeleteCommand(context.Background(), 1); err != nil {
		t.Fatalf("DeleteCommand: %v", err)
	}
}

func TestUpdateMovieFile(t *testing.T) {
	t.Parallel()

	want := radarr.MovieFile{ID: 1}

	srv := newTestServer(t, http.MethodPut, "/api/v3/moviefile/1", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.UpdateMovieFile(context.Background(), &radarr.MovieFile{ID: 1})
	if err != nil {
		t.Fatalf("UpdateMovieFile: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestEditMovieFiles(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPut, "/api/v3/moviefile/editor", nil)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.EditMovieFiles(context.Background(), &radarr.MovieFileEditorResource{
		MovieFileIDs: []int{1, 2},
	}); err != nil {
		t.Fatalf("EditMovieFiles: %v", err)
	}
}

func TestUpdateCustomFormatsBulk(t *testing.T) {
	t.Parallel()

	want := []arr.CustomFormatResource{{ID: 1, Name: "test"}}

	srv := newTestServer(t, http.MethodPut, "/api/v3/customformat/bulk", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
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

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.DeleteCustomFormatsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteCustomFormatsBulk: %v", err)
	}
}

func TestUpdateDownloadClientsBulk(t *testing.T) {
	t.Parallel()

	want := []arr.ProviderResource{{ID: 1}}

	srv := newTestServer(t, http.MethodPut, "/api/v3/downloadclient/bulk", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
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

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.DeleteDownloadClientsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteDownloadClientsBulk: %v", err)
	}
}

func TestTestAllDownloadClients(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPost, "/api/v3/downloadclient/testall", nil)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.TestAllDownloadClients(context.Background()); err != nil {
		t.Fatalf("TestAllDownloadClients: %v", err)
	}
}

func TestUpdateIndexersBulk(t *testing.T) {
	t.Parallel()

	want := []arr.ProviderResource{{ID: 1}}

	srv := newTestServer(t, http.MethodPut, "/api/v3/indexer/bulk", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
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

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.DeleteIndexersBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteIndexersBulk: %v", err)
	}
}

func TestTestAllIndexers(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPost, "/api/v3/indexer/testall", nil)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.TestAllIndexers(context.Background()); err != nil {
		t.Fatalf("TestAllIndexers: %v", err)
	}
}

func TestUpdateImportListsBulk(t *testing.T) {
	t.Parallel()

	want := []arr.ProviderResource{{ID: 1}}

	srv := newTestServer(t, http.MethodPut, "/api/v3/importlist/bulk", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
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

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.DeleteImportListsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteImportListsBulk: %v", err)
	}
}

func TestTestAllImportLists(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPost, "/api/v3/importlist/testall", nil)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.TestAllImportLists(context.Background()); err != nil {
		t.Fatalf("TestAllImportLists: %v", err)
	}
}

func TestGetImportListConfig(t *testing.T) {
	t.Parallel()

	want := radarr.ImportListConfigResource{ID: 1, ListSyncLevel: "disabled"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/config/importlist", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
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

	want := radarr.ImportListConfigResource{ID: 1, ListSyncLevel: "logOnly"}

	srv := newTestServer(t, http.MethodPut, "/api/v3/config/importlist/1", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.UpdateImportListConfig(context.Background(), &radarr.ImportListConfigResource{ID: 1, ListSyncLevel: "logOnly"})
	if err != nil {
		t.Fatalf("UpdateImportListConfig: %v", err)
	}
	if got.ListSyncLevel != "logOnly" {
		t.Errorf("ListSyncLevel = %q, want logOnly", got.ListSyncLevel)
	}
}

func TestGetImportListExclusionsPaged(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[radarr.ImportListExclusion]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records:      []radarr.ImportListExclusion{{ID: 1, TmdbID: 123}},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v3/exclusions/paged?page=1&pageSize=10", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
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

	srv := newTestServer(t, http.MethodDelete, "/api/v3/exclusions/bulk", nil)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.DeleteImportListExclusionsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteImportListExclusionsBulk: %v", err)
	}
}

func TestTestAllNotifications(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPost, "/api/v3/notification/testall", nil)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.TestAllNotifications(context.Background()); err != nil {
		t.Fatalf("TestAllNotifications: %v", err)
	}
}

func TestTestAllMetadataConsumers(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodPost, "/api/v3/metadata/testall", nil)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.TestAllMetadataConsumers(context.Background()); err != nil {
		t.Fatalf("TestAllMetadataConsumers: %v", err)
	}
}

func TestGetLanguage(t *testing.T) {
	t.Parallel()

	want := radarr.Language{ID: 1, Name: "English"}

	srv := newTestServer(t, http.MethodGet, "/api/v3/language/1", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
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

	want := radarr.LocalizationResource{ID: 1, Strings: map[string]string{"key": "val"}}

	srv := newTestServer(t, http.MethodGet, "/api/v3/localization", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
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

	want := []arr.QualityDefinitionResource{{ID: 1, Title: "Bluray-1080p"}}

	srv := newTestServer(t, http.MethodPut, "/api/v3/qualitydefinition/update", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.UpdateQualityDefinitions(context.Background(), []arr.QualityDefinitionResource{{ID: 1, Title: "Bluray-1080p"}})
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

	c, _ := radarr.New(srv.URL, "test-key")
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

	want := arr.RootFolder{ID: 1, Path: "/movies"}

	srv := newTestServer(t, http.MethodPut, "/api/v3/rootfolder/1", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.UpdateRootFolder(context.Background(), &arr.RootFolder{ID: 1, Path: "/movies"})
	if err != nil {
		t.Fatalf("UpdateRootFolder: %v", err)
	}
	if got.Path != "/movies" {
		t.Errorf("Path = %q, want /movies", got.Path)
	}
}

func TestBrowseFileSystem(t *testing.T) {
	t.Parallel()

	want := radarr.FileSystemResource{
		Directories: []radarr.FileSystemEntry{{Path: "/movies", Name: "movies"}},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v3/filesystem?path=%2Fmovies&includeFiles=true", want)
	defer srv.Close()

	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.BrowseFileSystem(context.Background(), "/movies", true)
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

	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestGetAlternativeTitles(t *testing.T) {
	t.Parallel()
	want := []radarr.AlternativeTitleResource{{ID: 1, Title: "Alt"}}
	srv := newTestServer(t, http.MethodGet, "/api/v3/alttitle", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetAlternativeTitles(context.Background())
	if err != nil {
		t.Fatalf("GetAlternativeTitles: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetAlternativeTitle(t *testing.T) {
	t.Parallel()
	want := radarr.AlternativeTitleResource{ID: 1, Title: "Alt"}
	srv := newTestServer(t, http.MethodGet, "/api/v3/alttitle/1", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetAlternativeTitle(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetAlternativeTitle: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetBlocklistMovies(t *testing.T) {
	t.Parallel()
	want := []arr.BlocklistResource{{ID: 1}}
	srv := newTestServer(t, http.MethodGet, "/api/v3/blocklist/movie", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetBlocklistMovies(context.Background())
	if err != nil {
		t.Fatalf("GetBlocklistMovies: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetCreditByID(t *testing.T) {
	t.Parallel()
	want := radarr.Credit{ID: 5, PersonName: "Actor"}
	srv := newTestServer(t, http.MethodGet, "/api/v3/credit/5", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetCreditByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetCreditByID: %v", err)
	}
	if got.ID != 5 {
		t.Errorf("ID = %d, want 5", got.ID)
	}
}

func TestGetDownloadClientConfigByID(t *testing.T) {
	t.Parallel()
	want := arr.DownloadClientConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/downloadclient/1", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
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
	c, _ := radarr.New(srv.URL, "test-key")
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
	want := radarr.ImportListConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/importlist/1", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
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
	c, _ := radarr.New(srv.URL, "test-key")
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
	c, _ := radarr.New(srv.URL, "test-key")
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
	c, _ := radarr.New(srv.URL, "test-key")
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
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetUIConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUIConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetExtraFiles(t *testing.T) {
	t.Parallel()
	want := []radarr.ExtraFileResource{{ID: 1, MovieID: 10}}
	srv := newTestServer(t, http.MethodGet, "/api/v3/extrafile?movieId=10", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetExtraFiles(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetExtraFiles: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetImportListExclusion(t *testing.T) {
	t.Parallel()
	want := radarr.ImportListExclusion{ID: 1, TmdbID: 123}
	srv := newTestServer(t, http.MethodGet, "/api/v3/exclusions/1", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetImportListExclusion(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetImportListExclusion: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestCreateImportListExclusionsBulk(t *testing.T) {
	t.Parallel()
	want := []radarr.ImportListExclusion{{ID: 1}, {ID: 2}}
	srv := newTestServer(t, http.MethodPost, "/api/v3/exclusions/bulk", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.CreateImportListExclusionsBulk(context.Background(), []radarr.ImportListExclusion{{TmdbID: 1}})
	if err != nil {
		t.Fatalf("CreateImportListExclusionsBulk: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("len = %d, want 2", len(got))
	}
}

func TestUpdateImportListExclusion(t *testing.T) {
	t.Parallel()
	want := radarr.ImportListExclusion{ID: 1, TmdbID: 456}
	srv := newTestServer(t, http.MethodPut, "/api/v3/exclusions/1", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.UpdateImportListExclusion(context.Background(), &radarr.ImportListExclusion{ID: 1, TmdbID: 456})
	if err != nil {
		t.Fatalf("UpdateImportListExclusion: %v", err)
	}
	if got.TmdbID != 456 {
		t.Errorf("TmdbID = %d, want 456", got.TmdbID)
	}
}

func TestGetImportListMovies(t *testing.T) {
	t.Parallel()
	want := []radarr.Movie{{ID: 1, Title: "Test"}}
	srv := newTestServer(t, http.MethodGet, "/api/v3/importlist/movie", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetImportListMovies(context.Background())
	if err != nil {
		t.Fatalf("GetImportListMovies: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestCreateImportListMovies(t *testing.T) {
	t.Parallel()
	want := []radarr.Movie{{ID: 1, Title: "Test"}}
	srv := newTestServer(t, http.MethodPost, "/api/v3/importlist/movie", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.CreateImportListMovies(context.Background(), []radarr.Movie{{Title: "Test"}})
	if err != nil {
		t.Fatalf("CreateImportListMovies: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestDownloadClientAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v3/downloadclient/action/testAction", nil)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.DownloadClientAction(context.Background(), "testAction", nil); err != nil {
		t.Fatalf("DownloadClientAction: %v", err)
	}
}

func TestImportListAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v3/importlist/action/testAction", nil)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.ImportListAction(context.Background(), "testAction", nil); err != nil {
		t.Fatalf("ImportListAction: %v", err)
	}
}

func TestIndexerAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v3/indexer/action/testAction", nil)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.IndexerAction(context.Background(), "testAction", nil); err != nil {
		t.Fatalf("IndexerAction: %v", err)
	}
}

func TestMetadataAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v3/metadata/action/testAction", nil)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.MetadataAction(context.Background(), "testAction", nil); err != nil {
		t.Fatalf("MetadataAction: %v", err)
	}
}

func TestNotificationAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v3/notification/action/testAction", nil)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.NotificationAction(context.Background(), "testAction", nil); err != nil {
		t.Fatalf("NotificationAction: %v", err)
	}
}

func TestGetLocalizationLanguages(t *testing.T) {
	t.Parallel()
	want := []radarr.LocalizationLanguageResource{{Identifier: "en"}}
	srv := newTestServer(t, http.MethodGet, "/api/v3/localization/language", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetLocalizationLanguages(context.Background())
	if err != nil {
		t.Fatalf("GetLocalizationLanguages: %v", err)
	}
	if len(got) != 1 || got[0].Identifier != "en" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestGetMetadataConfig(t *testing.T) {
	t.Parallel()
	want := radarr.MetadataConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/metadata", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetMetadataConfig(context.Background())
	if err != nil {
		t.Fatalf("GetMetadataConfig: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetMetadataConfigByID(t *testing.T) {
	t.Parallel()
	want := radarr.MetadataConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/metadata/1", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetMetadataConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMetadataConfigByID: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateMetadataConfig(t *testing.T) {
	t.Parallel()
	want := radarr.MetadataConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodPut, "/api/v3/config/metadata/1", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.UpdateMetadataConfig(context.Background(), &radarr.MetadataConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateMetadataConfig: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestUpdateMovieFilesBulk(t *testing.T) {
	t.Parallel()
	want := []radarr.MovieFile{{ID: 1}, {ID: 2}}
	srv := newTestServer(t, http.MethodPut, "/api/v3/moviefile/bulk", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.UpdateMovieFilesBulk(context.Background(), &radarr.MovieFileEditorResource{
		MovieFileIDs: []int{1, 2},
	})
	if err != nil {
		t.Fatalf("UpdateMovieFilesBulk: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("len = %d, want 2", len(got))
	}
}

func TestGetMovieFolder(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v3/movie/5/folder", nil)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.GetMovieFolder(context.Background(), 5); err != nil {
		t.Fatalf("GetMovieFolder: %v", err)
	}
}

func TestGetNamingExamples(t *testing.T) {
	t.Parallel()
	want := arr.NamingConfigResource{ID: 1}
	srv := newTestServer(t, http.MethodGet, "/api/v3/config/naming/examples", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
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
	want := radarr.QualityDefinitionLimitsResource{Min: 1, Max: 400}
	srv := newTestServer(t, http.MethodGet, "/api/v3/qualitydefinition/limits", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetQualityDefinitionLimits(context.Background())
	if err != nil {
		t.Fatalf("GetQualityDefinitionLimits: %v", err)
	}
	if got.Min != 1 || got.Max != 400 {
		t.Errorf("got Min=%d Max=%d, want 1/400", got.Min, got.Max)
	}
}

func TestGetSystemRoutesDuplicate(t *testing.T) {
	t.Parallel()
	want := []arr.SystemRouteResource{{Path: "/test", Method: "GET"}}
	srv := newTestServer(t, http.MethodGet, "/api/v3/system/routes/duplicate", want)
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetSystemRoutesDuplicate(context.Background())
	if err != nil {
		t.Fatalf("GetSystemRoutesDuplicate: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d, want 1", len(got))
	}
}

func TestGetUpdateLogFileContent(t *testing.T) {
	t.Parallel()
	srv := newRawTestServer(t, http.MethodGet, "/api/v3/log/file/update/update.txt", "log content")
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetUpdateLogFileContent(context.Background(), "update.txt")
	if err != nil {
		t.Fatalf("GetUpdateLogFileContent: %v", err)
	}
	if got != "log content" {
		t.Errorf("content = %q, want %q", got, "log content")
	}
}

func TestGetLogFileContent(t *testing.T) {
	t.Parallel()
	srv := newRawTestServer(t, http.MethodGet, "/api/v3/log/file/radarr.txt", "log line 1\nlog line 2")
	defer srv.Close()
	c, _ := radarr.New(srv.URL, "test-key")
	got, err := c.GetLogFileContent(context.Background(), "radarr.txt")
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
	c, _ := radarr.New(srv.URL, "test-key")
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
	c, _ := radarr.New(srv.URL, "test-key")
	if err := c.UploadBackup(context.Background(), "backup.zip", strings.NewReader("fake-backup-data")); err != nil {
		t.Fatalf("UploadBackup: %v", err)
	}
}
