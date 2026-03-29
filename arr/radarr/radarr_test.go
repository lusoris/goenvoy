package radarr_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/arr"
	"github.com/lusoris/goenvoy/arr/radarr"
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
