package tpdb_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/adult/tpdb"
)

func newTestServer(t *testing.T, wantPath, wantToken string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		wantAuth := "Bearer " + wantToken
		if got := r.Header.Get("Authorization"); got != wantAuth {
			t.Errorf("Authorization = %q, want %q", got, wantAuth)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q, want %q", got, "application/json")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
}

type paginatedResponse struct {
	Data any             `json:"data"`
	Meta tpdb.Pagination `json:"meta"`
}

func newListServer(t *testing.T, wantPath, wantToken string, data any) *httptest.Server {
	t.Helper()
	resp := paginatedResponse{
		Data: data,
		Meta: tpdb.Pagination{Total: 1, PerPage: 25, CurrentPage: 1, LastPage: 1},
	}
	return newTestServer(t, wantPath, wantToken, resp)
}

type itemResponseWrapper struct {
	Data any `json:"data"`
}

func newItemServer(t *testing.T, wantPath, wantToken string, data any) *httptest.Server {
	t.Helper()
	return newTestServer(t, wantPath, wantToken, itemResponseWrapper{Data: data})
}

func TestSearchScenes(t *testing.T) {
	t.Parallel()

	scene := tpdb.Scene{
		ID:    42,
		PID:   "abc123",
		Title: "Test Scene",
		Date:  "2024-01-15",
		Tags:  []tpdb.TagRef{{ID: 1, Name: "tag1", Slug: "tag1"}},
	}
	ts := newListServer(t, "/scenes", "test-token", []tpdb.Scene{scene})
	defer ts.Close()

	c := tpdb.New("test-token", metadata.WithBaseURL(ts.URL))
	scenes, pg, err := c.SearchScenes(context.Background(), &tpdb.SceneSearchParams{Query: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if len(scenes) != 1 {
		t.Fatalf("len = %d, want 1", len(scenes))
	}
	if scenes[0].Title != "Test Scene" {
		t.Errorf("Title = %q, want %q", scenes[0].Title, "Test Scene")
	}
	if scenes[0].PID != "abc123" {
		t.Errorf("PID = %q, want %q", scenes[0].PID, "abc123")
	}
	if pg.Total != 1 {
		t.Errorf("Total = %d, want 1", pg.Total)
	}
}

func TestGetScene(t *testing.T) {
	t.Parallel()

	scene := tpdb.Scene{
		ID:          99,
		PID:         "def456",
		Title:       "Specific Scene",
		Description: "A description",
		Duration:    1800,
		Site:        &tpdb.SiteRef{ID: 5, Name: "TestSite"},
	}
	ts := newItemServer(t, "/scenes/def456", "scene-token", scene)
	defer ts.Close()

	c := tpdb.New("scene-token", metadata.WithBaseURL(ts.URL))
	s, err := c.GetScene(context.Background(), "def456")
	if err != nil {
		t.Fatal(err)
	}
	if s.Title != "Specific Scene" {
		t.Errorf("Title = %q, want %q", s.Title, "Specific Scene")
	}
	if s.Duration != 1800 {
		t.Errorf("Duration = %d, want %d", s.Duration, 1800)
	}
	if s.Site == nil || s.Site.Name != "TestSite" {
		t.Error("Site not parsed correctly")
	}
}

func TestGetSimilarScenes(t *testing.T) {
	t.Parallel()

	scenes := []tpdb.Scene{{ID: 10, Title: "Similar 1"}, {ID: 11, Title: "Similar 2"}}
	ts := newListServer(t, "/scenes/abc/similar", "sim-token", scenes)
	defer ts.Close()

	c := tpdb.New("sim-token", metadata.WithBaseURL(ts.URL))
	result, err := c.GetSimilarScenes(context.Background(), "abc")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
}

func TestSearchPerformers(t *testing.T) {
	t.Parallel()

	perf := tpdb.Performer{
		ID:   1,
		PID:  "perf1",
		Name: "Test Performer",
		Extras: tpdb.ExtraData{
			Gender:    "female",
			Ethnicity: "caucasian",
		},
	}
	ts := newListServer(t, "/performers", "perf-token", []tpdb.Performer{perf})
	defer ts.Close()

	c := tpdb.New("perf-token", metadata.WithBaseURL(ts.URL))
	perfs, pg, err := c.SearchPerformers(context.Background(), &tpdb.PerformerSearchParams{Query: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if len(perfs) != 1 {
		t.Fatalf("len = %d, want 1", len(perfs))
	}
	if perfs[0].Name != "Test Performer" {
		t.Errorf("Name = %q, want %q", perfs[0].Name, "Test Performer")
	}
	if perfs[0].Extras.Gender != "female" {
		t.Errorf("Gender = %q, want %q", perfs[0].Extras.Gender, "female")
	}
	if pg.Total != 1 {
		t.Errorf("Total = %d, want 1", pg.Total)
	}
}

func TestGetPerformer(t *testing.T) {
	t.Parallel()

	perf := tpdb.Performer{
		ID:       2,
		PID:      "perf2",
		Name:     "Specific Performer",
		FullName: "Specific J. Performer",
		Bio:      "A bio",
		Rating:   4.5,
		Aliases:  []string{"Alias1", "Alias2"},
	}
	ts := newItemServer(t, "/performers/perf2", "p-token", perf)
	defer ts.Close()

	c := tpdb.New("p-token", metadata.WithBaseURL(ts.URL))
	p, err := c.GetPerformer(context.Background(), "perf2")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Specific Performer" {
		t.Errorf("Name = %q, want %q", p.Name, "Specific Performer")
	}
	if p.Rating != 4.5 {
		t.Errorf("Rating = %f, want %f", p.Rating, 4.5)
	}
	if len(p.Aliases) != 2 {
		t.Errorf("len(Aliases) = %d, want 2", len(p.Aliases))
	}
}

func TestGetPerformerScenes(t *testing.T) {
	t.Parallel()

	scenes := []tpdb.Scene{{ID: 20, Title: "Scene With Performer"}}
	ts := newListServer(t, "/performers/p1/scenes", "ps-token", scenes)
	defer ts.Close()

	c := tpdb.New("ps-token", metadata.WithBaseURL(ts.URL))
	result, pg, err := c.GetPerformerScenes(context.Background(), "p1", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if pg.CurrentPage != 1 {
		t.Errorf("CurrentPage = %d, want 1", pg.CurrentPage)
	}
}

func TestSearchSites(t *testing.T) {
	t.Parallel()

	site := tpdb.Site{ID: 3, UUID: "uuid-123", Name: "Test Studio", URL: "https://example.com"}
	ts := newListServer(t, "/sites", "site-token", []tpdb.Site{site})
	defer ts.Close()

	c := tpdb.New("site-token", metadata.WithBaseURL(ts.URL))
	sites, pg, err := c.SearchSites(context.Background(), "test", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(sites) != 1 {
		t.Fatalf("len = %d, want 1", len(sites))
	}
	if sites[0].Name != "Test Studio" {
		t.Errorf("Name = %q, want %q", sites[0].Name, "Test Studio")
	}
	if pg.Total != 1 {
		t.Errorf("Total = %d, want 1", pg.Total)
	}
}

func TestGetSite(t *testing.T) {
	t.Parallel()

	site := tpdb.Site{ID: 4, UUID: "uuid-456", Name: "Specific Site", Rating: 3.8}
	ts := newItemServer(t, "/sites/uuid-456", "s-token", site)
	defer ts.Close()

	c := tpdb.New("s-token", metadata.WithBaseURL(ts.URL))
	s, err := c.GetSite(context.Background(), "uuid-456")
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != "Specific Site" {
		t.Errorf("Name = %q, want %q", s.Name, "Specific Site")
	}
	if s.Rating != 3.8 {
		t.Errorf("Rating = %f, want %f", s.Rating, 3.8)
	}
}

func TestListTags(t *testing.T) {
	t.Parallel()

	tags := []tpdb.Tag{{ID: 1, Name: "Brunette", Slug: "brunette"}}
	ts := newListServer(t, "/tags", "tag-token", tags)
	defer ts.Close()

	c := tpdb.New("tag-token", metadata.WithBaseURL(ts.URL))
	result, pg, err := c.ListTags(context.Background(), "brunette", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Slug != "brunette" {
		t.Errorf("Slug = %q, want %q", result[0].Slug, "brunette")
	}
	if pg.Total != 1 {
		t.Errorf("Total = %d, want 1", pg.Total)
	}
}

func TestGetTag(t *testing.T) {
	t.Parallel()

	tag := tpdb.Tag{ID: 7, Name: "Blonde", Slug: "blonde"}
	ts := newItemServer(t, "/tags/blonde", "gt-token", tag)
	defer ts.Close()

	c := tpdb.New("gt-token", metadata.WithBaseURL(ts.URL))
	result, err := c.GetTag(context.Background(), "blonde")
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Blonde" {
		t.Errorf("Name = %q, want %q", result.Name, "Blonde")
	}
}

func TestListDirectors(t *testing.T) {
	t.Parallel()

	dirs := []tpdb.Director{{ID: 1, Name: "John Doe", Slug: "john-doe"}}
	ts := newListServer(t, "/directors", "dir-token", dirs)
	defer ts.Close()

	c := tpdb.New("dir-token", metadata.WithBaseURL(ts.URL))
	result, pg, err := c.ListDirectors(context.Background(), "john", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Name != "John Doe" {
		t.Errorf("Name = %q, want %q", result[0].Name, "John Doe")
	}
	if pg.Total != 1 {
		t.Errorf("Total = %d, want 1", pg.Total)
	}
}

func TestSearchMovies(t *testing.T) {
	t.Parallel()

	movie := tpdb.Movie{ID: 5, Title: "Test Movie", SKU: "ABC-123"}
	ts := newListServer(t, "/movies", "mov-token", []tpdb.Movie{movie})
	defer ts.Close()

	c := tpdb.New("mov-token", metadata.WithBaseURL(ts.URL))
	movies, pg, err := c.SearchMovies(context.Background(), "test", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(movies) != 1 {
		t.Fatalf("len = %d, want 1", len(movies))
	}
	if movies[0].SKU != "ABC-123" {
		t.Errorf("SKU = %q, want %q", movies[0].SKU, "ABC-123")
	}
	if pg.Total != 1 {
		t.Errorf("Total = %d, want 1", pg.Total)
	}
}

func TestGetMovie(t *testing.T) {
	t.Parallel()

	movie := tpdb.Movie{ID: 6, Title: "Specific Movie", Duration: 7200}
	ts := newItemServer(t, "/movies/6", "gm-token", movie)
	defer ts.Close()

	c := tpdb.New("gm-token", metadata.WithBaseURL(ts.URL))
	m, err := c.GetMovie(context.Background(), "6")
	if err != nil {
		t.Fatal(err)
	}
	if m.Title != "Specific Movie" {
		t.Errorf("Title = %q, want %q", m.Title, "Specific Movie")
	}
	if m.Duration != 7200 {
		t.Errorf("Duration = %d, want %d", m.Duration, 7200)
	}
}

func TestSearchJav(t *testing.T) {
	t.Parallel()

	jav := tpdb.Jav{ID: 8, Title: "Test JAV", SKU: "JAV-001"}
	ts := newListServer(t, "/jav", "jav-token", []tpdb.Jav{jav})
	defer ts.Close()

	c := tpdb.New("jav-token", metadata.WithBaseURL(ts.URL))
	result, pg, err := c.SearchJav(context.Background(), "test", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].SKU != "JAV-001" {
		t.Errorf("SKU = %q, want %q", result[0].SKU, "JAV-001")
	}
	if pg.Total != 1 {
		t.Errorf("Total = %d, want 1", pg.Total)
	}
}

func TestGetJav(t *testing.T) {
	t.Parallel()

	jav := tpdb.Jav{ID: 9, Title: "Specific JAV", Duration: 5400}
	ts := newItemServer(t, "/jav/9", "gj-token", jav)
	defer ts.Close()

	c := tpdb.New("gj-token", metadata.WithBaseURL(ts.URL))
	j, err := c.GetJav(context.Background(), "9")
	if err != nil {
		t.Fatal(err)
	}
	if j.Title != "Specific JAV" {
		t.Errorf("Title = %q, want %q", j.Title, "Specific JAV")
	}
}

func TestGetChanges(t *testing.T) {
	t.Parallel()

	scenes := []tpdb.Scene{{ID: 100, Title: "Changed Scene"}}
	ts := newListServer(t, "/changelog", "ch-token", scenes)
	defer ts.Close()

	c := tpdb.New("ch-token", metadata.WithBaseURL(ts.URL))
	result, pg, err := c.GetChanges(context.Background(), "2024-01-01T00:00:00Z", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if pg.Total != 1 {
		t.Errorf("Total = %d, want 1", pg.Total)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Unauthenticated."})
	}))
	defer ts.Close()

	c := tpdb.New("bad-token", metadata.WithBaseURL(ts.URL))
	_, err := c.GetScene(context.Background(), "1")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *tpdb.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
	if apiErr.Message != "Unauthenticated." {
		t.Errorf("Message = %q, want %q", apiErr.Message, "Unauthenticated.")
	}
}

func TestAPIErrorRawBody(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}))
	defer ts.Close()

	c := tpdb.New("token", metadata.WithBaseURL(ts.URL))
	_, err := c.GetScene(context.Background(), "1")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *tpdb.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.RawBody != "internal error" {
		t.Errorf("RawBody = %q, want %q", apiErr.RawBody, "internal error")
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	ts := newItemServer(t, "/scenes/1", "custom-token", tpdb.Scene{ID: 1, Title: "Custom"})
	defer ts.Close()

	custom := &http.Client{}
	c := tpdb.New("custom-token", metadata.WithBaseURL(ts.URL), metadata.WithHTTPClient(custom))
	s, err := c.GetScene(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if s.Title != "Custom" {
		t.Errorf("Title = %q, want %q", s.Title, "Custom")
	}
}

func TestGetSimilarPerformers(t *testing.T) {
	t.Parallel()

	perfs := []tpdb.Performer{{ID: 50, Name: "Similar Performer"}}
	ts := newListServer(t, "/performers/p1/similar", "sp-token", perfs)
	defer ts.Close()

	c := tpdb.New("sp-token", metadata.WithBaseURL(ts.URL))
	result, err := c.GetSimilarPerformers(context.Background(), "p1")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Name != "Similar Performer" {
		t.Errorf("Name = %q, want %q", result[0].Name, "Similar Performer")
	}
}

func TestGetPerformerMovies(t *testing.T) {
	t.Parallel()

	movies := []tpdb.Movie{{ID: 60, Title: "Performer Movie"}}
	ts := newListServer(t, "/performers/p2/movies", "pm-token", movies)
	defer ts.Close()

	c := tpdb.New("pm-token", metadata.WithBaseURL(ts.URL))
	result, pg, err := c.GetPerformerMovies(context.Background(), "p2", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Title != "Performer Movie" {
		t.Errorf("Title = %q, want %q", result[0].Title, "Performer Movie")
	}
	if pg.PerPage != 25 {
		t.Errorf("PerPage = %d, want 25", pg.PerPage)
	}
}

func TestGetDirector(t *testing.T) {
	t.Parallel()

	dir := tpdb.Director{ID: 12, Name: "Test Director", Slug: "test-director"}
	ts := newItemServer(t, "/directors/test-director", "gd-token", dir)
	defer ts.Close()

	c := tpdb.New("gd-token", metadata.WithBaseURL(ts.URL))
	d, err := c.GetDirector(context.Background(), "test-director")
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != "Test Director" {
		t.Errorf("Name = %q, want %q", d.Name, "Test Director")
	}
	if d.Slug != "test-director" {
		t.Errorf("Slug = %q, want %q", d.Slug, "test-director")
	}
}

func TestFindSceneByHash(t *testing.T) {
	t.Parallel()

	scenes := []tpdb.Scene{{ID: 77, Title: "Hashed Scene"}}
	ts := newListServer(t, "/scenes", "hash-token", scenes)
	defer ts.Close()

	c := tpdb.New("hash-token", metadata.WithBaseURL(ts.URL))
	result, err := c.FindSceneByHash(context.Background(), "abc123def456")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Title != "Hashed Scene" {
		t.Errorf("Title = %q, want %q", result[0].Title, "Hashed Scene")
	}
}
