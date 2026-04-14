package bazarr_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/arr/bazarr"
	"github.com/golusoris/goenvoy/arr/v2"
)

func newTestServer(t *testing.T, wantPath string, body any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
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
		c, err := bazarr.New("http://localhost:6767", "test-key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c == nil {
			t.Fatal("expected non-nil client")
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		t.Parallel()
		_, err := bazarr.New("://bad", "test-key")
		if err == nil {
			t.Fatal("expected error for invalid URL")
		}
	})
}

func TestGetSeries(t *testing.T) {
	t.Parallel()

	want := bazarr.PagedResponse[bazarr.Series]{
		Data: []bazarr.Series{
			{SonarrSeriesID: 1, Title: "Breaking Bad", Monitored: true},
		},
		Total: 1,
	}

	srv := newTestServer(t, "/api/series?length=-1&start=0", want)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetSeries(context.Background(), 0, -1)
	if err != nil {
		t.Fatalf("GetSeries: %v", err)
	}
	if len(got.Data) != 1 {
		t.Fatalf("len = %d, want 1", len(got.Data))
	}
	if got.Data[0].Title != "Breaking Bad" {
		t.Errorf("Title = %q, want %q", got.Data[0].Title, "Breaking Bad")
	}
}

func TestRunSeriesAction(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body["action"] != "scan-disk" {
			t.Errorf("action = %v, want scan-disk", body["action"])
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.RunSeriesAction(context.Background(), 1, "scan-disk"); err != nil {
		t.Fatalf("RunSeriesAction: %v", err)
	}
}

func TestGetEpisodes(t *testing.T) {
	t.Parallel()

	want := bazarr.PagedResponse[bazarr.Episode]{
		Data: []bazarr.Episode{
			{SonarrEpisodeID: 100, Title: "Pilot", Season: 1, EpisodeNumber: 1},
		},
	}

	srv := newTestServer(t, "/api/episodes?seriesid%5B%5D=1", want)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetEpisodes(context.Background(), []int{1}, nil)
	if err != nil {
		t.Fatalf("GetEpisodes: %v", err)
	}
	if got.Data[0].Title != "Pilot" {
		t.Errorf("Title = %q, want %q", got.Data[0].Title, "Pilot")
	}
}

func TestGetWantedEpisodes(t *testing.T) {
	t.Parallel()

	want := bazarr.PagedResponse[bazarr.WantedEpisode]{
		Data: []bazarr.WantedEpisode{
			{SeriesTitle: "Dexter", EpisodeNumber: "1x01", SonarrEpisodeID: 10},
		},
		Total: 1,
	}

	srv := newTestServer(t, "/api/episodes/wanted?start=0&length=25", want)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetWantedEpisodes(context.Background(), 0, 25)
	if err != nil {
		t.Fatalf("GetWantedEpisodes: %v", err)
	}
	if got.Data[0].SeriesTitle != "Dexter" {
		t.Errorf("SeriesTitle = %q, want %q", got.Data[0].SeriesTitle, "Dexter")
	}
}

func TestGetEpisodeHistory(t *testing.T) {
	t.Parallel()

	want := bazarr.PagedResponse[bazarr.EpisodeHistoryRecord]{
		Data: []bazarr.EpisodeHistoryRecord{
			{SeriesTitle: "Lost", Provider: "opensubtitles", Action: 1},
		},
		Total: 1,
	}

	srv := newTestServer(t, "/api/episodes/history?length=10&start=0", want)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetEpisodeHistory(context.Background(), 0, 10, nil)
	if err != nil {
		t.Fatalf("GetEpisodeHistory: %v", err)
	}
	if got.Data[0].Provider != "opensubtitles" {
		t.Errorf("Provider = %q, want %q", got.Data[0].Provider, "opensubtitles")
	}
}

func TestDownloadEpisodeSubtitle(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DownloadEpisodeSubtitle(context.Background(), 1, 100, "en", "false", "false"); err != nil {
		t.Fatalf("DownloadEpisodeSubtitle: %v", err)
	}
}

func TestGetMovies(t *testing.T) {
	t.Parallel()

	want := bazarr.PagedResponse[bazarr.Movie]{
		Data: []bazarr.Movie{
			{RadarrID: 1, Title: "Inception", Monitored: true},
		},
		Total: 1,
	}

	srv := newTestServer(t, "/api/movies?length=-1&start=0", want)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetMovies(context.Background(), 0, -1)
	if err != nil {
		t.Fatalf("GetMovies: %v", err)
	}
	if got.Data[0].Title != "Inception" {
		t.Errorf("Title = %q, want %q", got.Data[0].Title, "Inception")
	}
}

func TestRunMovieAction(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.RunMovieAction(context.Background(), 1, "search-missing"); err != nil {
		t.Fatalf("RunMovieAction: %v", err)
	}
}

func TestGetWantedMovies(t *testing.T) {
	t.Parallel()

	want := bazarr.PagedResponse[bazarr.WantedMovie]{
		Data: []bazarr.WantedMovie{
			{Title: "The Matrix", RadarrID: 5},
		},
		Total: 1,
	}

	srv := newTestServer(t, "/api/movies/wanted?start=0&length=25", want)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetWantedMovies(context.Background(), 0, 25)
	if err != nil {
		t.Fatalf("GetWantedMovies: %v", err)
	}
	if got.Data[0].Title != "The Matrix" {
		t.Errorf("Title = %q, want %q", got.Data[0].Title, "The Matrix")
	}
}

func TestGetProviders(t *testing.T) {
	t.Parallel()

	want := map[string]any{
		"data": []bazarr.Provider{
			{Name: "opensubtitlescom", Status: "Good", Retry: "-"},
		},
	}

	srv := newTestServer(t, "/api/providers", want)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetProviders(context.Background())
	if err != nil {
		t.Fatalf("GetProviders: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Name != "opensubtitlescom" {
		t.Errorf("Name = %q, want %q", got[0].Name, "opensubtitlescom")
	}
}

func TestResetProviders(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.ResetProviders(context.Background()); err != nil {
		t.Fatalf("ResetProviders: %v", err)
	}
}

func TestGetBadges(t *testing.T) {
	t.Parallel()

	want := bazarr.BadgeCounts{
		Episodes:  5,
		Movies:    3,
		Providers: 0,
		Status:    0,
	}

	srv := newTestServer(t, "/api/badges", want)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetBadges(context.Background())
	if err != nil {
		t.Fatalf("GetBadges: %v", err)
	}
	if got.Episodes != 5 {
		t.Errorf("Episodes = %d, want 5", got.Episodes)
	}
}

func TestGetSystemStatus(t *testing.T) {
	t.Parallel()

	want := map[string]any{
		"data": bazarr.SystemStatus{
			BazarrVersion:   "1.4.0",
			OperatingSystem: "Linux",
			PythonVersion:   "3.11.0",
		},
	}

	srv := newTestServer(t, "/api/system/status", want)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetSystemStatus(context.Background())
	if err != nil {
		t.Fatalf("GetSystemStatus: %v", err)
	}
	if got.BazarrVersion != "1.4.0" {
		t.Errorf("BazarrVersion = %q, want %q", got.BazarrVersion, "1.4.0")
	}
}

func TestGetHealth(t *testing.T) {
	t.Parallel()

	want := map[string]any{
		"data": []string{"Sonarr connection failed"},
	}

	srv := newTestServer(t, "/api/system/health", want)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetLanguages(t *testing.T) {
	t.Parallel()

	want := map[string]any{
		"data": []bazarr.Language{
			{Code2: "en", Code3: "eng", Name: "English", Enabled: true},
		},
	}

	srv := newTestServer(t, "/api/system/languages", want)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetLanguages(context.Background())
	if err != nil {
		t.Fatalf("GetLanguages: %v", err)
	}
	if got[0].Code2 != "en" {
		t.Errorf("Code2 = %q, want %q", got[0].Code2, "en")
	}
}

func TestPing(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, "/api/system/ping", nil)
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	}))
	defer srv.Close()

	c, err := bazarr.New(srv.URL, "bad-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = c.GetSeries(context.Background(), 0, -1)
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
