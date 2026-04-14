package shoko_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/anime/shoko"
)

func newTestServer(t *testing.T, wantPath, wantAPIKey string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if wantAPIKey != "" {
			if got := r.Header.Get("apikey"); got != wantAPIKey {
				t.Errorf("apikey = %q, want %q", got, wantAPIKey)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
}

func TestLogin(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/auth" {
			t.Errorf("path = %q, want /api/auth", r.URL.Path)
		}
		var body struct {
			User   string `json:"user"`
			Pass   string `json:"pass"`
			Device string `json:"device"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.User != "admin" {
			t.Errorf("user = %q, want admin", body.User)
		}
		if body.Pass != "secret" {
			t.Errorf("pass = %q, want secret", body.Pass)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"apikey": "test-api-key-123"})
	}))
	defer ts.Close()

	c := shoko.New(ts.URL)
	if err := c.Login(context.Background(), "admin", "secret"); err != nil {
		t.Fatal(err)
	}
}

func TestLoginError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"title": "Invalid credentials"})
	}))
	defer ts.Close()

	c := shoko.New(ts.URL)
	err := c.Login(context.Background(), "bad", "creds")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *shoko.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *shoko.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
}

func TestGetSeries(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Series/1", "my-key", shoko.Series{
		IDs:  shoko.SeriesIDs{ID: 1, AniDB: 16498},
		Name: "Frieren: Beyond Journey\u0027s End",
		Size: 28,
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("my-key"))
	s, err := c.GetSeries(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != "Frieren: Beyond Journey\u0027s End" {
		t.Errorf("Name = %q, want Frieren", s.Name)
	}
	if s.IDs.AniDB != 16498 {
		t.Errorf("AniDB = %d, want 16498", s.IDs.AniDB)
	}
	if s.Size != 28 {
		t.Errorf("Size = %d, want 28", s.Size)
	}
}

func TestListSeries(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Series", "list-key", []shoko.Series{
		{IDs: shoko.SeriesIDs{ID: 1}, Name: "Series One", Size: 12},
		{IDs: shoko.SeriesIDs{ID: 2}, Name: "Series Two", Size: 24},
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("list-key"))
	result, err := c.ListSeries(context.Background(), 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result[0].Name != "Series One" {
		t.Errorf("Name = %q, want Series One", result[0].Name)
	}
}

func TestSearchSeries(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("search"); got != "frieren" {
			t.Errorf("search = %q, want frieren", got)
		}
		if got := r.URL.Query().Get("fuzzy"); got != "true" {
			t.Errorf("fuzzy = %q, want true", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]shoko.Series{
			{IDs: shoko.SeriesIDs{ID: 1}, Name: "Frieren", Size: 28},
		})
	}))
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("search-key"))
	result, err := c.SearchSeries(context.Background(), "frieren", true, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Name != "Frieren" {
		t.Errorf("Name = %q, want Frieren", result[0].Name)
	}
}

func TestGetSeriesAniDB(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Series/1/AniDB", "anidb-key", shoko.AniDBAnime{
		ID:    16498,
		Type:  "TV",
		Title: "Sousou no Frieren",
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("anidb-key"))
	a, err := c.GetSeriesAniDB(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if a.Title != "Sousou no Frieren" {
		t.Errorf("Title = %q, want Sousou no Frieren", a.Title)
	}
	if a.Type != "TV" {
		t.Errorf("Type = %q, want TV", a.Type)
	}
}

func TestGetSeriesTags(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Series/1/Tags", "tag-key", []shoko.Tag{
		{Name: "fantasy", Weight: 600, Source: "AniDB"},
		{Name: "adventure", Weight: 400, Source: "AniDB"},
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("tag-key"))
	tags, err := c.GetSeriesTags(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 {
		t.Fatalf("len = %d, want 2", len(tags))
	}
	if tags[0].Name != "fantasy" {
		t.Errorf("Name = %q, want fantasy", tags[0].Name)
	}
	if tags[0].Weight != 600 {
		t.Errorf("Weight = %d, want 600", tags[0].Weight)
	}
}

func TestGetSeriesEpisodes(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Series/1/Episode", "ep-key", []shoko.Episode{
		{IDs: shoko.EpisodeIDs{ID: 10, AniDB: 280000}, Name: "The Journey Begins"},
		{IDs: shoko.EpisodeIDs{ID: 11, AniDB: 280001}, Name: "It Didn\u0027t Have to Be Magic"},
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("ep-key"))
	eps, err := c.GetSeriesEpisodes(context.Background(), 1, 1, 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 2 {
		t.Fatalf("len = %d, want 2", len(eps))
	}
	if eps[0].IDs.AniDB != 280000 {
		t.Errorf("AniDB = %d, want 280000", eps[0].IDs.AniDB)
	}
}

func TestGetAniDBAnime(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Series/AniDB/16498", "anime-key", shoko.AniDBAnime{
		ID:      16498,
		ShokoID: 1,
		Type:    "TV",
		Title:   "Sousou no Frieren",
		AirDate: "2023-09-29",
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("anime-key"))
	a, err := c.GetAniDBAnime(context.Background(), 16498)
	if err != nil {
		t.Fatal(err)
	}
	if a.ShokoID != 1 {
		t.Errorf("ShokoID = %d, want 1", a.ShokoID)
	}
	if a.AirDate != "2023-09-29" {
		t.Errorf("AirDate = %q, want 2023-09-29", a.AirDate)
	}
}

func TestGetSeriesByAniDBID(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Series/AniDB/16498/Series", "lookup-key", shoko.Series{
		IDs:  shoko.SeriesIDs{ID: 1, AniDB: 16498},
		Name: "Frieren",
		Size: 28,
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("lookup-key"))
	s, err := c.GetSeriesByAniDBID(context.Background(), 16498)
	if err != nil {
		t.Fatal(err)
	}
	if s.IDs.ID != 1 {
		t.Errorf("ID = %d, want 1", s.IDs.ID)
	}
}

func TestGetAniDBRelations(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Series/AniDB/16498/Relations", "rel-key", []shoko.AniDBRelation{
		{RelatedID: 16499, Type: "Sequel"},
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("rel-key"))
	rels, err := c.GetAniDBRelations(context.Background(), 16498)
	if err != nil {
		t.Fatal(err)
	}
	if len(rels) != 1 {
		t.Fatalf("len = %d, want 1", len(rels))
	}
	if rels[0].Type != "Sequel" {
		t.Errorf("Type = %q, want Sequel", rels[0].Type)
	}
}

func TestGetEpisode(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Episode/10", "ep2-key", shoko.Episode{
		IDs:  shoko.EpisodeIDs{ID: 10, ParentSeries: 1, AniDB: 280000},
		Name: "The Journey Begins",
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("ep2-key"))
	ep, err := c.GetEpisode(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if ep.IDs.ParentSeries != 1 {
		t.Errorf("ParentSeries = %d, want 1", ep.IDs.ParentSeries)
	}
	if ep.Name != "The Journey Begins" {
		t.Errorf("Name = %q, want The Journey Begins", ep.Name)
	}
}

func TestGetAniDBEpisode(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Episode/AniDB/280000", "aep-key", shoko.AniDBEpisode{
		ID:            280000,
		AnimeID:       16498,
		Type:          "Normal",
		EpisodeNumber: 1,
		Title:         "The Journey Begins",
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("aep-key"))
	ep, err := c.GetAniDBEpisode(context.Background(), 280000)
	if err != nil {
		t.Fatal(err)
	}
	if ep.EpisodeNumber != 1 {
		t.Errorf("EpisodeNumber = %d, want 1", ep.EpisodeNumber)
	}
	if ep.Type != "Normal" {
		t.Errorf("Type = %q, want Normal", ep.Type)
	}
}

func TestGetFile(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/File/100", "file-key", shoko.File{
		ID:   100,
		Size: 1400000000,
		Hashes: &shoko.FileHashes{
			ED2K:  "abc123def456",
			CRC32: "AABB1122",
		},
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("file-key"))
	f, err := c.GetFile(context.Background(), 100)
	if err != nil {
		t.Fatal(err)
	}
	if f.Size != 1400000000 {
		t.Errorf("Size = %d, want 1400000000", f.Size)
	}
	if f.Hashes.ED2K != "abc123def456" {
		t.Errorf("ED2K = %q, want abc123def456", f.Hashes.ED2K)
	}
}

func TestListManagedFolders(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/ManagedFolder", "mf-key", []shoko.ManagedFolder{
		{ID: 1, Path: "/anime", FileCount: 500},
		{ID: 2, Path: "/movies", FileCount: 200},
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("mf-key"))
	folders, err := c.ListManagedFolders(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(folders) != 2 {
		t.Fatalf("len = %d, want 2", len(folders))
	}
	if folders[0].Path != "/anime" {
		t.Errorf("Path = %q, want /anime", folders[0].Path)
	}
	if folders[0].FileCount != 500 {
		t.Errorf("FileCount = %d, want 500", folders[0].FileCount)
	}
}

func TestGetDashboardStats(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Dashboard/Stats", "dash-key", shoko.DashboardStats{
		SeriesCount:     150,
		FileCount:       3000,
		WatchedEpisodes: 500,
	})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("dash-key"))
	stats, err := c.GetDashboardStats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if stats.SeriesCount != 150 {
		t.Errorf("SeriesCount = %d, want 150", stats.SeriesCount)
	}
	if stats.FileCount != 3000 {
		t.Errorf("FileCount = %d, want 3000", stats.FileCount)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"title": "Forbidden"})
	}))
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("bad"))
	_, err := c.GetSeries(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *shoko.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *shoko.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode = %d, want 403", apiErr.StatusCode)
	}
	if apiErr.Title != "Forbidden" {
		t.Errorf("Title = %q, want Forbidden", apiErr.Title)
	}
}

func TestAPIErrorNonJSON(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("<html>Bad Gateway</html>"))
	}))
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("test"))
	_, err := c.GetSeries(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *shoko.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *shoko.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadGateway {
		t.Errorf("StatusCode = %d, want 502", apiErr.StatusCode)
	}
	if apiErr.RawBody != "<html>Bad Gateway</html>" {
		t.Errorf("RawBody = %q, want HTML body", apiErr.RawBody)
	}
}

func TestAPIErrorMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  shoko.APIError
		want string
	}{
		{"with title", shoko.APIError{StatusCode: 404, Title: "Not Found"}, "shoko: HTTP 404: Not Found"},
		{"raw body", shoko.APIError{StatusCode: 502, RawBody: "gateway error"}, "shoko: HTTP 502: gateway error"},
		{"code only", shoko.APIError{StatusCode: 500}, "shoko: HTTP 500"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRunImport(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Action/RunImport", "import-key", nil)
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("import-key"))
	if err := c.RunImport(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestContextCancellation(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Series/1", "cancel-key", shoko.Series{})
	defer ts.Close()

	c := shoko.New(ts.URL, shoko.WithAPIKey("cancel-key"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.GetSeries(ctx, 1)
	if err == nil {
		t.Fatal("expected error from canceled context")
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "/api/v3/Series/1", "custom-key", shoko.Series{
		IDs:  shoko.SeriesIDs{ID: 1},
		Name: "Test",
	})
	defer ts.Close()

	customClient := &http.Client{}
	c := shoko.New(ts.URL, shoko.WithHTTPClient(customClient), shoko.WithAPIKey("custom-key"))
	s, err := c.GetSeries(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != "Test" {
		t.Errorf("Name = %q, want Test", s.Name)
	}
}
