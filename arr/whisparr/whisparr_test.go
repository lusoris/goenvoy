package whisparr_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golusoris/goenvoy/arr/v2"
	"github.com/golusoris/goenvoy/arr/whisparr"
)

func newV2TestServer(t *testing.T, wantPath, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path+"?"+r.URL.RawQuery != wantPath && r.URL.Path != wantPath {
			t.Errorf("path = %s?%s, want %s", r.URL.Path, r.URL.RawQuery, wantPath)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
}

func newV3TestServer(t *testing.T, wantPath, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path+"?"+r.URL.RawQuery != wantPath && r.URL.Path != wantPath {
			t.Errorf("path = %s?%s, want %s", r.URL.Path, r.URL.RawQuery, wantPath)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
}

func newV2RawTestServer(t *testing.T, method, wantPath, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("method = %s, want %s", r.Method, method)
		}
		if r.URL.Path != wantPath {
			t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, body)
	}))
}

func newV3RawTestServer(t *testing.T, method, wantPath, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("method = %s, want %s", r.Method, method)
		}
		if r.URL.Path != wantPath {
			t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, body)
	}))
}

func newV2MethodTestServer(t *testing.T, method, wantPath string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("method = %s, want %s", r.Method, method)
		}
		if r.URL.Path != wantPath {
			t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
		}
		w.WriteHeader(http.StatusOK)
	}))
}

func newV3MethodTestServer(t *testing.T, method, wantPath string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("method = %s, want %s", r.Method, method)
		}
		if r.URL.Path != wantPath {
			t.Errorf("path = %s, want %s", r.URL.Path, wantPath)
		}
		w.WriteHeader(http.StatusOK)
	}))
}

func TestNew(t *testing.T) {
	t.Parallel()

	_, err := whisparr.New("http://localhost:6969", "abc123")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
}

func TestNewInvalidURL(t *testing.T) {
	t.Parallel()

	_, err := whisparr.New("://bad", "key")
	if err == nil {
		t.Fatal("New() with bad URL should fail")
	}
}

func TestNewV3(t *testing.T) {
	t.Parallel()

	_, err := whisparr.NewV3("http://localhost:6969", "abc123")
	if err != nil {
		t.Fatalf("NewV3() error = %v", err)
	}
}

func TestNewV3InvalidURL(t *testing.T) {
	t.Parallel()

	_, err := whisparr.NewV3("://bad", "key")
	if err == nil {
		t.Fatal("NewV3() with bad URL should fail")
	}
}

func TestGetAllSeries(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/series", `[{"id":1,"title":"Test Site","monitored":true}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	series, err := c.GetAllSeries(context.Background())
	if err != nil {
		t.Fatalf("GetAllSeries() error = %v", err)
	}
	if len(series) != 1 || series[0].Title != "Test Site" {
		t.Errorf("got %+v", series)
	}
}

func TestGetSeries(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/series/1", `{"id":1,"title":"Site A"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	s, err := c.GetSeries(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetSeries() error = %v", err)
	}
	if s.Title != "Site A" {
		t.Errorf("title = %s", s.Title)
	}
}

func TestAddSeries(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var s whisparr.Series
		json.NewDecoder(r.Body).Decode(&s)
		s.ID = 42
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	}))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	out, err := c.AddSeries(context.Background(), &whisparr.Series{Title: "New Site"})
	if err != nil {
		t.Fatalf("AddSeries() error = %v", err)
	}
	if out.ID != 42 {
		t.Errorf("id = %d, want 42", out.ID)
	}
}

func TestDeleteSeries(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteSeries(context.Background(), 1, true, false); err != nil {
		t.Fatalf("DeleteSeries() error = %v", err)
	}
}

func TestLookupSeries(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/series/lookup?term=test", `[{"id":1,"title":"Found"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	res, err := c.LookupSeries(context.Background(), "test")
	if err != nil {
		t.Fatalf("LookupSeries() error = %v", err)
	}
	if len(res) != 1 {
		t.Errorf("got %d results", len(res))
	}
}

func TestGetEpisodes(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/episode?seriesId=1", `[{"id":1,"title":"Ep 1","actors":[{"name":"Jane","gender":"female"}]}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	eps, err := c.GetEpisodes(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetEpisodes() error = %v", err)
	}
	if len(eps) != 1 || len(eps[0].Actors) != 1 {
		t.Errorf("got %+v", eps)
	}
	if eps[0].Actors[0].Gender != whisparr.GenderFemale {
		t.Errorf("gender = %s, want female", eps[0].Actors[0].Gender)
	}
}

func TestGetEpisodeFiles(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/episodefile?seriesId=1", `[{"id":1,"size":1024}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	files, err := c.GetEpisodeFiles(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetEpisodeFiles() error = %v", err)
	}
	if len(files) != 1 {
		t.Errorf("got %d files", len(files))
	}
}

func TestV2SendCommand(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/command", `{"id":1,"name":"RefreshSeries","status":"queued"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.SendCommand(context.Background(), struct {
		Name string `json:"name"`
	}{Name: "RefreshSeries"})
	if err != nil {
		t.Fatalf("SendCommand() error = %v", err)
	}
}

func TestV2GetSystemStatus(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/system/status", `{"appName":"Whisparr","version":"2.2.0"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	status, err := c.GetSystemStatus(context.Background())
	if err != nil {
		t.Fatalf("GetSystemStatus() error = %v", err)
	}
	if status.AppName != "Whisparr" {
		t.Errorf("appName = %s", status.AppName)
	}
}

func TestV2GetHealth(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/health", `[{"type":"warning","message":"test"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	health, err := c.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth() error = %v", err)
	}
	if len(health) != 1 {
		t.Errorf("got %d health checks", len(health))
	}
}

func TestV2GetTags(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/tag", `[{"id":1,"label":"hd"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	tags, err := c.GetTags(context.Background())
	if err != nil {
		t.Fatalf("GetTags() error = %v", err)
	}
	if len(tags) != 1 || tags[0].Label != "hd" {
		t.Errorf("got %+v", tags)
	}
}

func TestV2GetQualityProfiles(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/qualityprofile", `[{"id":1,"name":"Any"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	profiles, err := c.GetQualityProfiles(context.Background())
	if err != nil {
		t.Fatalf("GetQualityProfiles() error = %v", err)
	}
	if len(profiles) != 1 {
		t.Errorf("got %d profiles", len(profiles))
	}
}

func TestV2GetRootFolders(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/rootfolder", `[{"id":1,"path":"/data"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	folders, err := c.GetRootFolders(context.Background())
	if err != nil {
		t.Fatalf("GetRootFolders() error = %v", err)
	}
	if len(folders) != 1 {
		t.Errorf("got %d folders", len(folders))
	}
}

func TestV3GetAllMovies(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/movie", `[{"id":1,"title":"Scene 1","stashId":"abc","itemType":"scene"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	movies, err := c.GetAllMovies(context.Background())
	if err != nil {
		t.Fatalf("GetAllMovies() error = %v", err)
	}
	if len(movies) != 1 || movies[0].StashID != "abc" {
		t.Errorf("got %+v", movies)
	}
}

func TestV3GetMovie(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/movie/1", `{"id":1,"title":"Movie A","code":"ABC-123"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	m, err := c.GetMovie(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMovie() error = %v", err)
	}
	if m.Code != "ABC-123" {
		t.Errorf("code = %s", m.Code)
	}
}

func TestV3AddMovie(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var m whisparr.Movie
		json.NewDecoder(r.Body).Decode(&m)
		m.ID = 42
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	out, err := c.AddMovie(context.Background(), &whisparr.Movie{Title: "New Scene"})
	if err != nil {
		t.Fatalf("AddMovie() error = %v", err)
	}
	if out.ID != 42 {
		t.Errorf("id = %d, want 42", out.ID)
	}
}

func TestV3DeleteMovie(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteMovie(context.Background(), 1, true, false); err != nil {
		t.Fatalf("DeleteMovie() error = %v", err)
	}
}

func TestV3LookupScene(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/lookup/scene?term=test", `[{"id":1,"title":"Found"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	res, err := c.LookupScene(context.Background(), "test")
	if err != nil {
		t.Fatalf("LookupScene() error = %v", err)
	}
	if len(res) != 1 {
		t.Errorf("got %d results", len(res))
	}
}

func TestV3GetPerformers(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/performer", `[{"id":1,"name":"Jane Doe","gender":"female","stashId":"xyz"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	perfs, err := c.GetPerformers(context.Background())
	if err != nil {
		t.Fatalf("GetPerformers() error = %v", err)
	}
	if len(perfs) != 1 || perfs[0].Name != "Jane Doe" {
		t.Errorf("got %+v", perfs)
	}
	if perfs[0].Gender != whisparr.GenderFemale {
		t.Errorf("gender = %s, want female", perfs[0].Gender)
	}
}

func TestV3AddPerformer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var p whisparr.Performer
		json.NewDecoder(r.Body).Decode(&p)
		p.ID = 10
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	out, err := c.AddPerformer(context.Background(), &whisparr.Performer{Name: "New"})
	if err != nil {
		t.Fatalf("AddPerformer() error = %v", err)
	}
	if out.ID != 10 {
		t.Errorf("id = %d, want 10", out.ID)
	}
}

func TestV3DeletePerformer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeletePerformer(context.Background(), 1, false); err != nil {
		t.Fatalf("DeletePerformer() error = %v", err)
	}
}

func TestV3GetStudios(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/studio", `[{"id":1,"title":"Studio X","stashId":"s1"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	studios, err := c.GetStudios(context.Background())
	if err != nil {
		t.Fatalf("GetStudios() error = %v", err)
	}
	if len(studios) != 1 || studios[0].Title != "Studio X" {
		t.Errorf("got %+v", studios)
	}
}

func TestV3AddStudio(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var s whisparr.Studio
		json.NewDecoder(r.Body).Decode(&s)
		s.ID = 5
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	out, err := c.AddStudio(context.Background(), &whisparr.Studio{Title: "New Studio"})
	if err != nil {
		t.Fatalf("AddStudio() error = %v", err)
	}
	if out.ID != 5 {
		t.Errorf("id = %d, want 5", out.ID)
	}
}

func TestV3DeleteStudio(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteStudio(context.Background(), 1, false); err != nil {
		t.Fatalf("DeleteStudio() error = %v", err)
	}
}

func TestV3GetCredits(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/credit?movieId=1", `[{"id":1,"personName":"Jane","type":"cast"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	credits, err := c.GetCredits(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCredits() error = %v", err)
	}
	if len(credits) != 1 {
		t.Errorf("got %d credits", len(credits))
	}
}

func TestV3GetMoviesByPerformer(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/movie/listbyperformerforeignid?performerForeignId=abc", `[{"id":1,"title":"Scene 1"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	movies, err := c.GetMoviesByPerformer(context.Background(), "abc")
	if err != nil {
		t.Fatalf("GetMoviesByPerformer() error = %v", err)
	}
	if len(movies) != 1 {
		t.Errorf("got %d movies", len(movies))
	}
}

func TestV3GetSystemStatus(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/system/status", `{"appName":"Whisparr","version":"3.3.3"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	status, err := c.GetSystemStatus(context.Background())
	if err != nil {
		t.Fatalf("GetSystemStatus() error = %v", err)
	}
	if status.Version != "3.3.3" {
		t.Errorf("version = %s", status.Version)
	}
}

func TestV3GetImportExclusions(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/exclusions", `[{"id":1,"movieTitle":"Excluded"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	exclusions, err := c.GetImportExclusions(context.Background())
	if err != nil {
		t.Fatalf("GetImportExclusions() error = %v", err)
	}
	if len(exclusions) != 1 {
		t.Errorf("got %d exclusions", len(exclusions))
	}
}

func TestV2ErrorResponse(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	}))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "bad-key")
	_, err := c.GetAllSeries(context.Background())
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}

func TestV3ErrorResponse(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "bad-key")
	_, err := c.GetAllMovies(context.Background())
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}

// V2 untested methods.

func TestUpdateSeries(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		var s whisparr.Series
		json.NewDecoder(r.Body).Decode(&s)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	}))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	out, err := c.UpdateSeries(context.Background(), &whisparr.Series{ID: 1, Title: "Updated"}, true)
	if err != nil {
		t.Fatalf("UpdateSeries() error = %v", err)
	}
	if out.Title != "Updated" {
		t.Errorf("title = %s", out.Title)
	}
}

func TestGetEpisode(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/episode/5", `{"id":5,"title":"Scene 5"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	ep, err := c.GetEpisode(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetEpisode() error = %v", err)
	}
	if ep.ID != 5 {
		t.Errorf("id = %d, want 5", ep.ID)
	}
}

func TestDeleteEpisodeFile(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteEpisodeFile(context.Background(), 1); err != nil {
		t.Fatalf("DeleteEpisodeFile() error = %v", err)
	}
}

func TestV2GetCalendar(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/calendar?start=2026-01-01&end=2026-01-31&unmonitored=false", `[{"id":1,"title":"Upcoming"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	eps, err := c.GetCalendar(context.Background(), "2026-01-01", "2026-01-31", false)
	if err != nil {
		t.Fatalf("GetCalendar() error = %v", err)
	}
	if len(eps) != 1 {
		t.Errorf("got %d episodes", len(eps))
	}
}

func TestV2Parse(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/parse?title=test+scene", `{"title":"test scene"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	result, err := c.Parse(context.Background(), "test scene")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestV2GetDiskSpace(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/diskspace", `[{"path":"/data","freeSpace":1000}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	ds, err := c.GetDiskSpace(context.Background())
	if err != nil {
		t.Fatalf("GetDiskSpace() error = %v", err)
	}
	if len(ds) != 1 {
		t.Errorf("got %d disk spaces", len(ds))
	}
}

func TestV2GetQueue(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/queue?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	q, err := c.GetQueue(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetQueue() error = %v", err)
	}
	if q.Page != 1 {
		t.Errorf("page = %d", q.Page)
	}
}

func TestV2CreateTag(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":1,"label":"new-tag"}`))
	}))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	tag, err := c.CreateTag(context.Background(), "new-tag")
	if err != nil {
		t.Fatalf("CreateTag() error = %v", err)
	}
	if tag.Label != "new-tag" {
		t.Errorf("label = %s", tag.Label)
	}
}

func TestV2GetHistory(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/history?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	h, err := c.GetHistory(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetHistory() error = %v", err)
	}
	if h.Page != 1 {
		t.Errorf("page = %d", h.Page)
	}
}

func TestV2UpdateSeasonPass(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	err := c.UpdateSeasonPass(context.Background(), whisparr.SeasonPassResource{})
	if err != nil {
		t.Fatalf("UpdateSeasonPass() error = %v", err)
	}
}

// V3 untested methods.

func TestV3UpdateMovie(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		var m whisparr.Movie
		json.NewDecoder(r.Body).Decode(&m)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	out, err := c.UpdateMovie(context.Background(), &whisparr.Movie{ID: 1, Title: "Updated"}, true)
	if err != nil {
		t.Fatalf("UpdateMovie() error = %v", err)
	}
	if out.Title != "Updated" {
		t.Errorf("title = %s", out.Title)
	}
}

func TestV3LookupMovie(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/lookup/movie?term=test", `[{"id":1,"title":"Found"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	res, err := c.LookupMovie(context.Background(), "test")
	if err != nil {
		t.Fatalf("LookupMovie() error = %v", err)
	}
	if len(res) != 1 {
		t.Errorf("got %d results", len(res))
	}
}

func TestV3GetMoviesByStudio(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/movie/listbystudioforeignid?studioForeignId=s1", `[{"id":1,"title":"Scene 1"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	movies, err := c.GetMoviesByStudio(context.Background(), "s1")
	if err != nil {
		t.Fatalf("GetMoviesByStudio() error = %v", err)
	}
	if len(movies) != 1 {
		t.Errorf("got %d movies", len(movies))
	}
}

func TestV3GetMovieFile(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/moviefile/1", `{"id":1,"size":2048}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	f, err := c.GetMovieFile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMovieFile() error = %v", err)
	}
	if f.ID != 1 {
		t.Errorf("id = %d", f.ID)
	}
}

func TestV3DeleteMovieFile(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteMovieFile(context.Background(), 1); err != nil {
		t.Fatalf("DeleteMovieFile() error = %v", err)
	}
}

func TestV3EditMovies(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	err := c.EditMovies(context.Background(), &whisparr.MovieEditorResource{MovieIDs: []int{1, 2}})
	if err != nil {
		t.Fatalf("EditMovies() error = %v", err)
	}
}

func TestV3DeleteMovies(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	err := c.DeleteMovies(context.Background(), &whisparr.MovieEditorResource{MovieIDs: []int{1}})
	if err != nil {
		t.Fatalf("DeleteMovies() error = %v", err)
	}
}

func TestV3GetPerformer(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/performer/1", `{"id":1,"name":"Jane"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	p, err := c.GetPerformer(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetPerformer() error = %v", err)
	}
	if p.Name != "Jane" {
		t.Errorf("name = %s", p.Name)
	}
}

func TestV3UpdatePerformer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		var p whisparr.Performer
		json.NewDecoder(r.Body).Decode(&p)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	out, err := c.UpdatePerformer(context.Background(), &whisparr.Performer{ID: 1, Name: "Updated"})
	if err != nil {
		t.Fatalf("UpdatePerformer() error = %v", err)
	}
	if out.Name != "Updated" {
		t.Errorf("name = %s", out.Name)
	}
}

func TestV3GetStudio(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/studio/1", `{"id":1,"title":"Studio A"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	s, err := c.GetStudio(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetStudio() error = %v", err)
	}
	if s.Title != "Studio A" {
		t.Errorf("title = %s", s.Title)
	}
}

func TestV3UpdateStudio(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		var s whisparr.Studio
		json.NewDecoder(r.Body).Decode(&s)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	out, err := c.UpdateStudio(context.Background(), &whisparr.Studio{ID: 1, Title: "Updated"})
	if err != nil {
		t.Fatalf("UpdateStudio() error = %v", err)
	}
	if out.Title != "Updated" {
		t.Errorf("title = %s", out.Title)
	}
}

func TestV3GetCalendar(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/calendar?start=2026-01-01&end=2026-01-31&unmonitored=false", `[{"id":1,"title":"Upcoming"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	movies, err := c.GetCalendar(context.Background(), "2026-01-01", "2026-01-31", false)
	if err != nil {
		t.Fatalf("GetCalendar() error = %v", err)
	}
	if len(movies) != 1 {
		t.Errorf("got %d movies", len(movies))
	}
}

func TestV3SendCommand(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/command", `{"id":1,"name":"RefreshMovie","status":"queued"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.SendCommand(context.Background(), struct {
		Name string `json:"name"`
	}{Name: "RefreshMovie"})
	if err != nil {
		t.Fatalf("SendCommand() error = %v", err)
	}
}

func TestV3Parse(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/parse?title=test+movie", `{"title":"test movie"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	result, err := c.Parse(context.Background(), "test movie")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestV3GetHealth(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/health", `[{"type":"warning","message":"test"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	health, err := c.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth() error = %v", err)
	}
	if len(health) != 1 {
		t.Errorf("got %d health checks", len(health))
	}
}

func TestV3GetDiskSpace(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/diskspace", `[{"path":"/data","freeSpace":1000}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	ds, err := c.GetDiskSpace(context.Background())
	if err != nil {
		t.Fatalf("GetDiskSpace() error = %v", err)
	}
	if len(ds) != 1 {
		t.Errorf("got %d disk spaces", len(ds))
	}
}

func TestV3GetQueue(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/queue?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	q, err := c.GetQueue(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetQueue() error = %v", err)
	}
	if q.Page != 1 {
		t.Errorf("page = %d", q.Page)
	}
}

func TestV3GetQualityProfiles(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/qualityprofile", `[{"id":1,"name":"Any"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	profiles, err := c.GetQualityProfiles(context.Background())
	if err != nil {
		t.Fatalf("GetQualityProfiles() error = %v", err)
	}
	if len(profiles) != 1 {
		t.Errorf("got %d profiles", len(profiles))
	}
}

func TestV3GetTags(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/tag", `[{"id":1,"label":"hd"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	tags, err := c.GetTags(context.Background())
	if err != nil {
		t.Fatalf("GetTags() error = %v", err)
	}
	if len(tags) != 1 {
		t.Errorf("got %d tags", len(tags))
	}
}

func TestV3CreateTag(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":1,"label":"new-tag"}`))
	}))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	tag, err := c.CreateTag(context.Background(), "new-tag")
	if err != nil {
		t.Fatalf("CreateTag() error = %v", err)
	}
	if tag.Label != "new-tag" {
		t.Errorf("label = %s", tag.Label)
	}
}

func TestV3GetRootFolders(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/rootfolder", `[{"id":1,"path":"/data"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	folders, err := c.GetRootFolders(context.Background())
	if err != nil {
		t.Fatalf("GetRootFolders() error = %v", err)
	}
	if len(folders) != 1 {
		t.Errorf("got %d folders", len(folders))
	}
}

func TestV3GetHistory(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/history?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	h, err := c.GetHistory(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetHistory() error = %v", err)
	}
	if h.Page != 1 {
		t.Errorf("page = %d", h.Page)
	}
}

// V2 Extended Tests.

func TestV2GetAutoTags(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/autotagging", `[{"id":1,"name":"test"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	tags, err := c.GetAutoTags(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(tags) != 1 {
		t.Errorf("got %d", len(tags))
	}
}

func TestV2GetAutoTag(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/autotagging/1", `{"id":1,"name":"test"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	tag, err := c.GetAutoTag(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if tag.ID != 1 {
		t.Errorf("id = %d", tag.ID)
	}
}

func TestV2CreateAutoTag(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/autotagging", `{"id":1,"name":"test"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateAutoTag(context.Background(), &arr.AutoTaggingResource{Name: "test"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2UpdateAutoTag(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/autotagging/1", `{"id":1,"name":"updated"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateAutoTag(context.Background(), &arr.AutoTaggingResource{ID: 1, Name: "updated"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteAutoTag(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteAutoTag(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetAutoTagSchema(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/autotagging/schema", `[{"name":"test"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	s, err := c.GetAutoTagSchema(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(s) != 1 {
		t.Errorf("got %d", len(s))
	}
}

func TestV2GetBackups(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/system/backup", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	b, err := c.GetBackups(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(b) != 1 {
		t.Errorf("got %d", len(b))
	}
}

func TestV2DeleteBackup(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteBackup(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2RestoreBackup(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/system/backup/restore/1", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.RestoreBackup(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetBlocklist(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/blocklist?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	b, err := c.GetBlocklist(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if b.Page != 1 {
		t.Errorf("page = %d", b.Page)
	}
}

func TestV2DeleteBlocklistItem(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteBlocklistItem(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2BulkDeleteBlocklist(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.BulkDeleteBlocklist(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetCalendarByID(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/calendar/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	ep, err := c.GetCalendarByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if ep.ID != 1 {
		t.Errorf("id = %d", ep.ID)
	}
}

func TestV2GetCommands(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/command", `[{"id":1,"name":"test"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	cmds, err := c.GetCommands(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(cmds) != 1 {
		t.Errorf("got %d", len(cmds))
	}
}

func TestV2GetCommand(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/command/1", `{"id":1,"name":"test"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	cmd, err := c.GetCommand(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cmd.ID != 1 {
		t.Errorf("id = %d", cmd.ID)
	}
}

func TestV2DeleteCommand(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteCommand(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetCustomFilters(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/customfilter", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	f, err := c.GetCustomFilters(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(f) != 1 {
		t.Errorf("got %d", len(f))
	}
}

func TestV2GetCustomFilter(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/customfilter/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	f, err := c.GetCustomFilter(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if f.ID != 1 {
		t.Errorf("id = %d", f.ID)
	}
}

func TestV2CreateCustomFilter(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/customfilter", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateCustomFilter(context.Background(), &arr.CustomFilterResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2UpdateCustomFilter(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/customfilter/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateCustomFilter(context.Background(), &arr.CustomFilterResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteCustomFilter(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteCustomFilter(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetCustomFormats(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/customformat", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	f, err := c.GetCustomFormats(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(f) != 1 {
		t.Errorf("got %d", len(f))
	}
}

func TestV2GetCustomFormat(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/customformat/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	f, err := c.GetCustomFormat(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if f.ID != 1 {
		t.Errorf("id = %d", f.ID)
	}
}

func TestV2CreateCustomFormat(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/customformat", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateCustomFormat(context.Background(), &arr.CustomFormatResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2UpdateCustomFormat(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/customformat/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateCustomFormat(context.Background(), &arr.CustomFormatResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteCustomFormat(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteCustomFormat(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetCustomFormatSchema(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/customformat/schema", `[{"name":"test"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	s, err := c.GetCustomFormatSchema(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(s) != 1 {
		t.Errorf("got %d", len(s))
	}
}

func TestV2GetWantedMissing(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/wanted/missing?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	m, err := c.GetWantedMissing(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if m.Page != 1 {
		t.Errorf("page = %d", m.Page)
	}
}

func TestV2GetWantedMissingByID(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/wanted/missing/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	ep, err := c.GetWantedMissingByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if ep.ID != 1 {
		t.Errorf("id = %d", ep.ID)
	}
}

func TestV2GetWantedCutoff(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/wanted/cutoff?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	m, err := c.GetWantedCutoff(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if m.Page != 1 {
		t.Errorf("page = %d", m.Page)
	}
}

func TestV2GetWantedCutoffByID(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/wanted/cutoff/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	ep, err := c.GetWantedCutoffByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if ep.ID != 1 {
		t.Errorf("id = %d", ep.ID)
	}
}

func TestV2GetDelayProfiles(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/delayprofile", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	d, err := c.GetDelayProfiles(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(d) != 1 {
		t.Errorf("got %d", len(d))
	}
}

func TestV2GetDelayProfile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/delayprofile/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	d, err := c.GetDelayProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if d.ID != 1 {
		t.Errorf("id = %d", d.ID)
	}
}

func TestV2CreateDelayProfile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/delayprofile", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateDelayProfile(context.Background(), &arr.DelayProfileResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2UpdateDelayProfile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/delayprofile/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateDelayProfile(context.Background(), &arr.DelayProfileResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteDelayProfile(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteDelayProfile(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2ReorderDelayProfile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/delayprofile/reorder/1?after=2", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.ReorderDelayProfile(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetDownloadClients(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/downloadclient", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	d, err := c.GetDownloadClients(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(d) != 1 {
		t.Errorf("got %d", len(d))
	}
}

func TestV2GetDownloadClient(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/downloadclient/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	d, err := c.GetDownloadClient(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if d.ID != 1 {
		t.Errorf("id = %d", d.ID)
	}
}

func TestV2CreateDownloadClient(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/downloadclient", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateDownloadClient(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2UpdateDownloadClient(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/downloadclient/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateDownloadClient(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteDownloadClient(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteDownloadClient(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetDownloadClientSchema(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/downloadclient/schema", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	s, err := c.GetDownloadClientSchema(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(s) != 1 {
		t.Errorf("got %d", len(s))
	}
}

func TestV2TestDownloadClient(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/downloadclient/test", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.TestDownloadClient(context.Background(), &arr.ProviderResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2TestAllDownloadClients(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/downloadclient/testall", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.TestAllDownloadClients(context.Background()); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2BulkUpdateDownloadClients(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/downloadclient/bulk", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.BulkUpdateDownloadClients(context.Background(), &arr.ProviderBulkResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2BulkDeleteDownloadClients(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.BulkDeleteDownloadClients(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DownloadClientAction(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/downloadclient/action/testAction", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DownloadClientAction(context.Background(), "testAction", &arr.ProviderResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetDownloadClientConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/downloadclient", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	cfg, err := c.GetDownloadClientConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV2UpdateDownloadClientConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/downloadclient/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateDownloadClientConfig(context.Background(), &arr.DownloadClientConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2UpdateEpisode(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/episode/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateEpisode(context.Background(), &whisparr.Episode{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2MonitorEpisodes(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/episode/monitor", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.MonitorEpisodes(context.Background(), &whisparr.EpisodesMonitoredResource{EpisodeIDs: []int{1}, Monitored: true})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetEpisodeFile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/episodefile/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	f, err := c.GetEpisodeFile(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if f.ID != 1 {
		t.Errorf("id = %d", f.ID)
	}
}

func TestV2UpdateEpisodeFile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/episodefile/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateEpisodeFile(context.Background(), &whisparr.EpisodeFile{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2EditEpisodeFiles(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/episodefile/editor", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.EditEpisodeFiles(context.Background(), &whisparr.EpisodeFileEditorResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2BulkDeleteEpisodeFiles(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.BulkDeleteEpisodeFiles(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2BulkUpdateEpisodeFiles(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/episodefile/bulk", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.BulkUpdateEpisodeFiles(context.Background(), &whisparr.EpisodeFileEditorResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2BrowseFileSystem(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/filesystem?path=%2Fdata", `{"directories":[]}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.BrowseFileSystem(context.Background(), "/data")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetHostConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/host", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	cfg, err := c.GetHostConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV2UpdateHostConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/host/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateHostConfig(context.Background(), &arr.HostConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetHistorySince(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/history/since?date=2024-01-01", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	h, err := c.GetHistorySince(context.Background(), "2024-01-01")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(h) != 1 {
		t.Errorf("got %d", len(h))
	}
}

func TestV2GetHistorySeries(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/history/series?seriesId=1", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	h, err := c.GetHistorySeries(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(h) != 1 {
		t.Errorf("got %d", len(h))
	}
}

func TestV2MarkHistoryFailed(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/history/failed/1", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.MarkHistoryFailed(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetImportLists(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/importlist", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	l, err := c.GetImportLists(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(l) != 1 {
		t.Errorf("got %d", len(l))
	}
}

func TestV2GetImportList(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/importlist/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	l, err := c.GetImportList(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if l.ID != 1 {
		t.Errorf("id = %d", l.ID)
	}
}

func TestV2CreateImportList(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/importlist", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateImportList(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2UpdateImportList(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/importlist/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateImportList(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteImportList(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteImportList(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetImportListSchema(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/importlist/schema", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	s, err := c.GetImportListSchema(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(s) != 1 {
		t.Errorf("got %d", len(s))
	}
}

func TestV2TestImportList(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/importlist/test", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.TestImportList(context.Background(), &arr.ProviderResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2TestAllImportLists(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/importlist/testall", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.TestAllImportLists(context.Background()); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetImportListConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/importlist", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	cfg, err := c.GetImportListConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV2UpdateImportListConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/importlist/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateImportListConfig(context.Background(), &whisparr.ImportListConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetImportListExclusions(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/importlistexclusion", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	e, err := c.GetImportListExclusions(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(e) != 1 {
		t.Errorf("got %d", len(e))
	}
}

func TestV2GetImportListExclusion(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/importlistexclusion/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	e, err := c.GetImportListExclusion(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if e.ID != 1 {
		t.Errorf("id = %d", e.ID)
	}
}

func TestV2CreateImportListExclusion(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/importlistexclusion", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateImportListExclusion(context.Background(), &arr.ImportListExclusionResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2UpdateImportListExclusion(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/importlistexclusion/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateImportListExclusion(context.Background(), &arr.ImportListExclusionResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteImportListExclusion(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteImportListExclusion(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetIndexers(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/indexer", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	idx, err := c.GetIndexers(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(idx) != 1 {
		t.Errorf("got %d", len(idx))
	}
}

func TestV2GetIndexer(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/indexer/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	idx, err := c.GetIndexer(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if idx.ID != 1 {
		t.Errorf("id = %d", idx.ID)
	}
}

func TestV2CreateIndexer(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/indexer", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateIndexer(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2UpdateIndexer(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/indexer/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateIndexer(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteIndexer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteIndexer(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetIndexerSchema(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/indexer/schema", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	s, err := c.GetIndexerSchema(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(s) != 1 {
		t.Errorf("got %d", len(s))
	}
}

func TestV2TestIndexer(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/indexer/test", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.TestIndexer(context.Background(), &arr.ProviderResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2TestAllIndexers(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/indexer/testall", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.TestAllIndexers(context.Background()); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetIndexerConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/indexer", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	cfg, err := c.GetIndexerConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV2UpdateIndexerConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/indexer/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateIndexerConfig(context.Background(), &arr.IndexerConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetLanguages(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/language", `[{"id":1,"name":"English"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	l, err := c.GetLanguages(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(l) != 1 {
		t.Errorf("got %d", len(l))
	}
}

func TestV2GetLanguage(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/language/1", `{"id":1,"name":"English"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	l, err := c.GetLanguage(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if l.ID != 1 {
		t.Errorf("id = %d", l.ID)
	}
}

func TestV2GetLanguageProfiles(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/languageprofile", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	lp, err := c.GetLanguageProfiles(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(lp) != 1 {
		t.Errorf("got %d", len(lp))
	}
}

func TestV2GetLanguageProfile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/languageprofile/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	lp, err := c.GetLanguageProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if lp.ID != 1 {
		t.Errorf("id = %d", lp.ID)
	}
}

func TestV2CreateLanguageProfile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/languageprofile", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateLanguageProfile(context.Background(), &whisparr.LanguageProfileResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2UpdateLanguageProfile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/languageprofile/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateLanguageProfile(context.Background(), &whisparr.LanguageProfileResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteLanguageProfile(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteLanguageProfile(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetLocalization(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/localization", `{"key":"value"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	l, err := c.GetLocalization(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if l["key"] != "value" {
		t.Errorf("got %v", l)
	}
}

func TestV2GetLogs(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/log?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	l, err := c.GetLogs(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if l.Page != 1 {
		t.Errorf("page = %d", l.Page)
	}
}

func TestV2GetLogFiles(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/log/file", `[{"filename":"test.log"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	f, err := c.GetLogFiles(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(f) != 1 {
		t.Errorf("got %d", len(f))
	}
}

func TestV2GetManualImport(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/manualimport?folder=%2Fdata", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	m, err := c.GetManualImport(context.Background(), "/data")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(m) != 1 {
		t.Errorf("got %d", len(m))
	}
}

func TestV2GetMediaManagementConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/mediamanagement", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	cfg, err := c.GetMediaManagementConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV2UpdateMediaManagementConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/mediamanagement/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateMediaManagementConfig(context.Background(), &arr.MediaManagementConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetMetadata(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/metadata", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	m, err := c.GetMetadata(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(m) != 1 {
		t.Errorf("got %d", len(m))
	}
}

func TestV2GetMetadataByID(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/metadata/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	m, err := c.GetMetadataByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if m.ID != 1 {
		t.Errorf("id = %d", m.ID)
	}
}

func TestV2CreateMetadata(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/metadata", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateMetadata(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteMetadata(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteMetadata(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetNamingConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/naming", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	cfg, err := c.GetNamingConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV2UpdateNamingConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/naming/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateNamingConfig(context.Background(), &arr.NamingConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetNotifications(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/notification", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	n, err := c.GetNotifications(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(n) != 1 {
		t.Errorf("got %d", len(n))
	}
}

func TestV2GetNotification(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/notification/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	n, err := c.GetNotification(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if n.ID != 1 {
		t.Errorf("id = %d", n.ID)
	}
}

func TestV2CreateNotification(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/notification", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateNotification(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteNotification(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteNotification(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetQualityDefinitions(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/qualitydefinition", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	d, err := c.GetQualityDefinitions(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(d) != 1 {
		t.Errorf("got %d", len(d))
	}
}

func TestV2GetQualityDefinition(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/qualitydefinition/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	d, err := c.GetQualityDefinition(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if d.ID != 1 {
		t.Errorf("id = %d", d.ID)
	}
}

func TestV2UpdateQualityDefinition(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/qualitydefinition/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateQualityDefinition(context.Background(), &arr.QualityDefinitionResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetQualityProfile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/qualityprofile/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	p, err := c.GetQualityProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if p.ID != 1 {
		t.Errorf("id = %d", p.ID)
	}
}

func TestV2CreateQualityProfile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/qualityprofile", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateQualityProfile(context.Background(), &arr.QualityProfile{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteQualityProfile(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteQualityProfile(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteQueueItem(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteQueueItem(context.Background(), 1, true, false); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GrabQueueItem(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/queue/grab/1", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.GrabQueueItem(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetQueueDetails(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/queue/details", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	d, err := c.GetQueueDetails(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(d) != 1 {
		t.Errorf("got %d", len(d))
	}
}

func TestV2GetQueueStatus(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/queue/status", `{"totalCount":5}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	s, err := c.GetQueueStatus(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if s.TotalCount != 5 {
		t.Errorf("totalCount = %d", s.TotalCount)
	}
}

func TestV2SearchReleases(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/release?episodeId=1", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	r, err := c.SearchReleases(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(r) != 1 {
		t.Errorf("got %d", len(r))
	}
}

func TestV2GrabRelease(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/release", `{"guid":"abc"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.GrabRelease(context.Background(), &arr.ReleaseResource{GUID: "abc"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetReleaseProfiles(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/releaseprofile", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	rp, err := c.GetReleaseProfiles(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(rp) != 1 {
		t.Errorf("got %d", len(rp))
	}
}

func TestV2GetReleaseProfile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/releaseprofile/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	rp, err := c.GetReleaseProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if rp.ID != 1 {
		t.Errorf("id = %d", rp.ID)
	}
}

func TestV2CreateReleaseProfile(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/releaseprofile", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateReleaseProfile(context.Background(), &arr.ReleaseProfileResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteReleaseProfile(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteReleaseProfile(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetRemotePathMappings(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/remotepathmapping", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	r, err := c.GetRemotePathMappings(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(r) != 1 {
		t.Errorf("got %d", len(r))
	}
}

func TestV2GetRemotePathMapping(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/remotepathmapping/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	r, err := c.GetRemotePathMapping(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if r.ID != 1 {
		t.Errorf("id = %d", r.ID)
	}
}

func TestV2CreateRemotePathMapping(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/remotepathmapping", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateRemotePathMapping(context.Background(), &arr.RemotePathMappingResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteRemotePathMapping(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteRemotePathMapping(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetRenamePreview(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/rename?seriesId=1", `[{"episodeFileId":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	r, err := c.GetRenamePreview(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(r) != 1 {
		t.Errorf("got %d", len(r))
	}
}

func TestV2GetRootFolder(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/rootfolder/1", `{"id":1,"path":"/data"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	r, err := c.GetRootFolder(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if r.ID != 1 {
		t.Errorf("id = %d", r.ID)
	}
}

func TestV2CreateRootFolder(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/rootfolder", `{"id":1,"path":"/data"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.CreateRootFolder(context.Background(), &arr.RootFolder{Path: "/data"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteRootFolder(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteRootFolder(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2EditSeries(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/series/editor", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.EditSeries(context.Background(), &whisparr.SeriesEditorResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteSeriesBulk(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteSeriesBulk(context.Background(), &whisparr.SeriesEditorResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2ImportSeries(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/series/import", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.ImportSeries(context.Background(), []whisparr.Series{{Title: "Test"}}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetSystemRoutes(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/system/routes", `[{"path":"/api/v3/test"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	r, err := c.GetSystemRoutes(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(r) != 1 {
		t.Errorf("got %d", len(r))
	}
}

func TestV2Shutdown(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/system/shutdown", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2Restart(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/system/restart", `{}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.Restart(context.Background()); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetTag(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/tag/1", `{"id":1,"label":"test"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	tag, err := c.GetTag(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if tag.ID != 1 {
		t.Errorf("id = %d", tag.ID)
	}
}

func TestV2UpdateTag(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/tag/1", `{"id":1,"label":"updated"}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateTag(context.Background(), &arr.Tag{ID: 1, Label: "updated"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2DeleteTag(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.DeleteTag(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetTagDetails(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/tag/detail", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	d, err := c.GetTagDetails(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(d) != 1 {
		t.Errorf("got %d", len(d))
	}
}

func TestV2GetTagDetail(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/tag/detail/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	d, err := c.GetTagDetail(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if d.ID != 1 {
		t.Errorf("id = %d", d.ID)
	}
}

func TestV2GetTasks(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/system/task", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	tasks, err := c.GetTasks(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("got %d", len(tasks))
	}
}

func TestV2GetTask(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/system/task/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	task, err := c.GetTask(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if task.ID != 1 {
		t.Errorf("id = %d", task.ID)
	}
}

func TestV2GetUIConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/ui", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	cfg, err := c.GetUIConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV2UpdateUIConfig(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/config/ui/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	_, err := c.UpdateUIConfig(context.Background(), &arr.UIConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV2GetUpdates(t *testing.T) {
	t.Parallel()

	ts := newV2TestServer(t, "/api/v3/update", `[{"version":"2.0.0"}]`)
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	u, err := c.GetUpdates(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(u) != 1 {
		t.Errorf("got %d", len(u))
	}
}

// V3 Extended Tests.

func TestV3GetAlternativeTitles(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/alttitle?movieId=1", `[{"id":1,"title":"Alt"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	a, err := c.GetAlternativeTitles(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(a) != 1 {
		t.Errorf("got %d", len(a))
	}
}

func TestV3GetAlternativeTitle(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/alttitle/1", `{"id":1,"title":"Alt"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	a, err := c.GetAlternativeTitle(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if a.ID != 1 {
		t.Errorf("id = %d", a.ID)
	}
}

func TestV3GetAutoTags(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/autotagging", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	tags, err := c.GetAutoTags(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(tags) != 1 {
		t.Errorf("got %d", len(tags))
	}
}

func TestV3GetAutoTag(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/autotagging/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	tag, err := c.GetAutoTag(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if tag.ID != 1 {
		t.Errorf("id = %d", tag.ID)
	}
}

func TestV3CreateAutoTag(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/autotagging", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateAutoTag(context.Background(), &arr.AutoTaggingResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteAutoTag(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteAutoTag(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetBackups(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/system/backup", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	b, err := c.GetBackups(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(b) != 1 {
		t.Errorf("got %d", len(b))
	}
}

func TestV3DeleteBackup(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteBackup(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetBlocklist(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/blocklist?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	b, err := c.GetBlocklist(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if b.Page != 1 {
		t.Errorf("page = %d", b.Page)
	}
}

func TestV3DeleteBlocklistItem(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteBlocklistItem(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetBlocklistMovie(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/blocklist/movie?movieId=1", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	b, err := c.GetBlocklistMovie(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(b) != 1 {
		t.Errorf("got %d", len(b))
	}
}

func TestV3GetCalendarByID(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/calendar/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	m, err := c.GetCalendarByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if m.ID != 1 {
		t.Errorf("id = %d", m.ID)
	}
}

func TestV3GetCommands(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/command", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	cmds, err := c.GetCommands(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(cmds) != 1 {
		t.Errorf("got %d", len(cmds))
	}
}

func TestV3GetCommand(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/command/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	cmd, err := c.GetCommand(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cmd.ID != 1 {
		t.Errorf("id = %d", cmd.ID)
	}
}

func TestV3DeleteCommand(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteCommand(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetCredit(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/credit/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	cr, err := c.GetCredit(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cr.ID != 1 {
		t.Errorf("id = %d", cr.ID)
	}
}

func TestV3GetCustomFilters(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/customfilter", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	f, err := c.GetCustomFilters(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(f) != 1 {
		t.Errorf("got %d", len(f))
	}
}

func TestV3GetCustomFilter(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/customfilter/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	f, err := c.GetCustomFilter(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if f.ID != 1 {
		t.Errorf("id = %d", f.ID)
	}
}

func TestV3CreateCustomFilter(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/customfilter", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateCustomFilter(context.Background(), &arr.CustomFilterResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteCustomFilter(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteCustomFilter(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetCustomFormats(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/customformat", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	f, err := c.GetCustomFormats(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(f) != 1 {
		t.Errorf("got %d", len(f))
	}
}

func TestV3CreateCustomFormat(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/customformat", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateCustomFormat(context.Background(), &arr.CustomFormatResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteCustomFormat(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteCustomFormat(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetDelayProfiles(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/delayprofile", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	d, err := c.GetDelayProfiles(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(d) != 1 {
		t.Errorf("got %d", len(d))
	}
}

func TestV3CreateDelayProfile(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/delayprofile", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateDelayProfile(context.Background(), &arr.DelayProfileResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteDelayProfile(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteDelayProfile(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetDownloadClients(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/downloadclient", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	d, err := c.GetDownloadClients(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(d) != 1 {
		t.Errorf("got %d", len(d))
	}
}

func TestV3CreateDownloadClient(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/downloadclient", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateDownloadClient(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteDownloadClient(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteDownloadClient(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetDownloadClientConfig(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/config/downloadclient", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	cfg, err := c.GetDownloadClientConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV3GetExtraFiles(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/extrafile?movieId=1", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	ef, err := c.GetExtraFiles(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(ef) != 1 {
		t.Errorf("got %d", len(ef))
	}
}

func TestV3BrowseFileSystem(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/filesystem?path=%2Fdata", `{"directories":[]}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.BrowseFileSystem(context.Background(), "/data")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetHealthByID(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/health/1", `{"source":"test","type":"warning","message":"msg"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	h, err := c.GetHealthByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if h.Source != "test" {
		t.Errorf("source = %s", h.Source)
	}
}

func TestV3GetHistorySince(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/history/since?date=2024-01-01", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	h, err := c.GetHistorySince(context.Background(), "2024-01-01")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(h) != 1 {
		t.Errorf("got %d", len(h))
	}
}

func TestV3GetHistoryMovie(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/history/movie?movieId=1", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	h, err := c.GetHistoryMovie(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(h) != 1 {
		t.Errorf("got %d", len(h))
	}
}

func TestV3MarkHistoryFailed(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/history/failed/1", `{}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.MarkHistoryFailed(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetHostConfig(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/config/host", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	cfg, err := c.GetHostConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV3GetImportLists(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/importlist", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	l, err := c.GetImportLists(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(l) != 1 {
		t.Errorf("got %d", len(l))
	}
}

func TestV3CreateImportList(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/importlist", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateImportList(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteImportList(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteImportList(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetImportListMovies(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/importlist/movie", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	m, err := c.GetImportListMovies(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(m) != 1 {
		t.Errorf("got %d", len(m))
	}
}

func TestV3GetImportExclusion(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/exclusions/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	e, err := c.GetImportExclusion(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if e.ID != 1 {
		t.Errorf("id = %d", e.ID)
	}
}

func TestV3CreateImportExclusion(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/exclusions", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateImportExclusion(context.Background(), &whisparr.ImportExclusion{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteImportExclusion(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteImportExclusion(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetIndexers(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/indexer", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	idx, err := c.GetIndexers(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(idx) != 1 {
		t.Errorf("got %d", len(idx))
	}
}

func TestV3CreateIndexer(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/indexer", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateIndexer(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteIndexer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteIndexer(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetIndexerConfig(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/config/indexer", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	cfg, err := c.GetIndexerConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV3GetIndexerFlags(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/indexerflag", `[{"id":1,"name":"freeleech"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	f, err := c.GetIndexerFlags(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(f) != 1 {
		t.Errorf("got %d", len(f))
	}
}

func TestV3GetLocalization(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/localization", `{"key":"value"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	l, err := c.GetLocalization(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if l["key"] != "value" {
		t.Errorf("got %v", l)
	}
}

func TestV3GetLogs(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/log?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	l, err := c.GetLogs(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if l.Page != 1 {
		t.Errorf("page = %d", l.Page)
	}
}

func TestV3GetLogFiles(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/log/file", `[{"filename":"test.log"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	f, err := c.GetLogFiles(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(f) != 1 {
		t.Errorf("got %d", len(f))
	}
}

func TestV3GetManualImport(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/manualimport?folder=%2Fdata", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	m, err := c.GetManualImport(context.Background(), "/data")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(m) != 1 {
		t.Errorf("got %d", len(m))
	}
}

func TestV3GetMediaManagementConfig(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/config/mediamanagement", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	cfg, err := c.GetMediaManagementConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV3GetMetadata(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/metadata", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	m, err := c.GetMetadata(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(m) != 1 {
		t.Errorf("got %d", len(m))
	}
}

func TestV3CreateMetadata(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/metadata", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateMetadata(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteMetadata(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteMetadata(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetMovieFiles(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/moviefile?movieId=1", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	f, err := c.GetMovieFiles(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(f) != 1 {
		t.Errorf("got %d", len(f))
	}
}

func TestV3UpdateMovieFile(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/moviefile/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.UpdateMovieFile(context.Background(), &whisparr.MovieFile{ID: 1})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3EditMovieFiles(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/moviefile/editor", `{}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.EditMovieFiles(context.Background(), &whisparr.MovieFileEditorResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3BulkDeleteMovieFiles(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.BulkDeleteMovieFiles(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3LookupMovieByTMDB(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/lookup/movie/tmdb?tmdbId=123", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	m, err := c.LookupMovieByTMDB(context.Background(), 123)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if m.ID != 1 {
		t.Errorf("id = %d", m.ID)
	}
}

func TestV3LookupMovieByIMDB(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/lookup/movie/imdb?imdbId=tt123", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	m, err := c.LookupMovieByIMDB(context.Background(), "tt123")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if m.ID != 1 {
		t.Errorf("id = %d", m.ID)
	}
}

func TestV3GetNamingConfig(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/config/naming", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	cfg, err := c.GetNamingConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV3GetNotifications(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/notification", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	n, err := c.GetNotifications(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(n) != 1 {
		t.Errorf("got %d", len(n))
	}
}

func TestV3CreateNotification(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/notification", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateNotification(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteNotification(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteNotification(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3EditPerformers(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/performer/editor", `{}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.EditPerformers(context.Background(), &whisparr.PerformerEditorResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeletePerformersBulk(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeletePerformersBulk(context.Background(), &whisparr.PerformerEditorResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3EditStudios(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/studio/editor", `{}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.EditStudios(context.Background(), &whisparr.StudioEditorResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteStudiosBulk(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteStudiosBulk(context.Background(), &whisparr.StudioEditorResource{}); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetQualityDefinitions(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/qualitydefinition", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	d, err := c.GetQualityDefinitions(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(d) != 1 {
		t.Errorf("got %d", len(d))
	}
}

func TestV3GetQualityProfile(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/qualityprofile/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	p, err := c.GetQualityProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if p.ID != 1 {
		t.Errorf("id = %d", p.ID)
	}
}

func TestV3CreateQualityProfile(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/qualityprofile", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateQualityProfile(context.Background(), &arr.QualityProfile{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteQualityProfile(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteQualityProfile(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteQueueItem(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteQueueItem(context.Background(), 1, true, false); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GrabQueueItem(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/queue/grab/1", `{}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.GrabQueueItem(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetQueueDetails(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/queue/details", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	d, err := c.GetQueueDetails(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(d) != 1 {
		t.Errorf("got %d", len(d))
	}
}

func TestV3GetQueueStatus(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/queue/status", `{"totalCount":5}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	s, err := c.GetQueueStatus(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if s.TotalCount != 5 {
		t.Errorf("totalCount = %d", s.TotalCount)
	}
}

func TestV3SearchReleases(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/release?movieId=1", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	r, err := c.SearchReleases(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(r) != 1 {
		t.Errorf("got %d", len(r))
	}
}

func TestV3GrabRelease(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/release", `{"guid":"abc"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.GrabRelease(context.Background(), &arr.ReleaseResource{GUID: "abc"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetRemotePathMappings(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/remotepathmapping", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	r, err := c.GetRemotePathMappings(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(r) != 1 {
		t.Errorf("got %d", len(r))
	}
}

func TestV3CreateRemotePathMapping(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/remotepathmapping", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateRemotePathMapping(context.Background(), &arr.RemotePathMappingResource{})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteRemotePathMapping(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteRemotePathMapping(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetRenamePreview(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/rename?movieId=1", `[{"movieFileId":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	r, err := c.GetRenamePreview(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(r) != 1 {
		t.Errorf("got %d", len(r))
	}
}

func TestV3GetRootFolder(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/rootfolder/1", `{"id":1,"path":"/data"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	r, err := c.GetRootFolder(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if r.ID != 1 {
		t.Errorf("id = %d", r.ID)
	}
}

func TestV3CreateRootFolder(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/rootfolder", `{"id":1,"path":"/data"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.CreateRootFolder(context.Background(), &arr.RootFolder{Path: "/data"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteRootFolder(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteRootFolder(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetSystemRoutes(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/system/routes", `[{"path":"/api/v3/test"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	r, err := c.GetSystemRoutes(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(r) != 1 {
		t.Errorf("got %d", len(r))
	}
}

func TestV3Shutdown(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/system/shutdown", `{}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3Restart(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/system/restart", `{}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.Restart(context.Background()); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetTag(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/tag/1", `{"id":1,"label":"test"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	tag, err := c.GetTag(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if tag.ID != 1 {
		t.Errorf("id = %d", tag.ID)
	}
}

func TestV3UpdateTag(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/tag/1", `{"id":1,"label":"updated"}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	_, err := c.UpdateTag(context.Background(), &arr.Tag{ID: 1, Label: "updated"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3DeleteTag(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.DeleteTag(context.Background(), 1); err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestV3GetTagDetails(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/tag/detail", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	d, err := c.GetTagDetails(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(d) != 1 {
		t.Errorf("got %d", len(d))
	}
}

func TestV3GetTasks(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/system/task", `[{"id":1}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	tasks, err := c.GetTasks(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("got %d", len(tasks))
	}
}

func TestV3GetUIConfig(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/config/ui", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	cfg, err := c.GetUIConfig(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if cfg.ID != 1 {
		t.Errorf("id = %d", cfg.ID)
	}
}

func TestV3GetUpdates(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/update", `[{"version":"3.0.0"}]`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	u, err := c.GetUpdates(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(u) != 1 {
		t.Errorf("got %d", len(u))
	}
}

func TestV3GetWantedMissing(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/wanted/missing?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	m, err := c.GetWantedMissing(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if m.Page != 1 {
		t.Errorf("page = %d", m.Page)
	}
}

func TestV3GetWantedMissingByID(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/wanted/missing/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	m, err := c.GetWantedMissingByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if m.ID != 1 {
		t.Errorf("id = %d", m.ID)
	}
}

func TestV3GetWantedCutoff(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/wanted/cutoff?page=1&pageSize=10", `{"page":1,"pageSize":10,"totalRecords":0,"records":[]}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	m, err := c.GetWantedCutoff(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if m.Page != 1 {
		t.Errorf("page = %d", m.Page)
	}
}

func TestV3GetWantedCutoffByID(t *testing.T) {
	t.Parallel()

	ts := newV3TestServer(t, "/api/v3/wanted/cutoff/1", `{"id":1}`)
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	m, err := c.GetWantedCutoffByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if m.ID != 1 {
		t.Errorf("id = %d", m.ID)
	}
}

// ---------- V2 Ping / HeadPing / UploadBackup ----------.

func TestV2Ping(t *testing.T) {
	t.Parallel()

	ts := newV2MethodTestServer(t, http.MethodGet, "/ping")
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestV2HeadPing(t *testing.T) {
	t.Parallel()

	ts := newV2MethodTestServer(t, http.MethodHead, "/ping")
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	if err := c.HeadPing(context.Background()); err != nil {
		t.Fatalf("HeadPing: %v", err)
	}
}

func TestV2UploadBackup(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
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
	c, _ := whisparr.New(srv.URL, "test-key")
	if err := c.UploadBackup(context.Background(), "backup.zip", strings.NewReader("fake")); err != nil {
		t.Fatalf("UploadBackup: %v", err)
	}
}

func TestV2GetLogFileContent(t *testing.T) {
	t.Parallel()

	ts := newV2RawTestServer(t, http.MethodGet, "/api/v3/log/file/whisparr.txt", "log content")
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	got, err := c.GetLogFileContent(context.Background(), "whisparr.txt")
	if err != nil {
		t.Fatalf("GetLogFileContent: %v", err)
	}
	if got != "log content" {
		t.Errorf("content = %q, want %q", got, "log content")
	}
}

func TestV2GetUpdateLogFileContent(t *testing.T) {
	t.Parallel()

	ts := newV2RawTestServer(t, http.MethodGet, "/api/v3/log/file/update/update.txt", "update log")
	defer ts.Close()
	c, _ := whisparr.New(ts.URL, "test-key")
	got, err := c.GetUpdateLogFileContent(context.Background(), "update.txt")
	if err != nil {
		t.Fatalf("GetUpdateLogFileContent: %v", err)
	}
	if got != "update log" {
		t.Errorf("content = %q, want %q", got, "update log")
	}
}

// ---------- V3 Ping / HeadPing / UploadBackup ----------.

func TestV3Ping(t *testing.T) {
	t.Parallel()

	ts := newV3MethodTestServer(t, http.MethodGet, "/ping")
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestV3HeadPing(t *testing.T) {
	t.Parallel()

	ts := newV3MethodTestServer(t, http.MethodHead, "/ping")
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	if err := c.HeadPing(context.Background()); err != nil {
		t.Fatalf("HeadPing: %v", err)
	}
}

func TestV3UploadBackup(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
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
	c, _ := whisparr.NewV3(srv.URL, "test-key")
	if err := c.UploadBackup(context.Background(), "backup.zip", strings.NewReader("fake")); err != nil {
		t.Fatalf("UploadBackup: %v", err)
	}
}

func TestV3GetLogFileContent(t *testing.T) {
	t.Parallel()

	ts := newV3RawTestServer(t, http.MethodGet, "/api/v3/log/file/whisparr.txt", "log content")
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	got, err := c.GetLogFileContent(context.Background(), "whisparr.txt")
	if err != nil {
		t.Fatalf("GetLogFileContent: %v", err)
	}
	if got != "log content" {
		t.Errorf("content = %q, want %q", got, "log content")
	}
}

func TestV3GetUpdateLogFileContent(t *testing.T) {
	t.Parallel()

	ts := newV3RawTestServer(t, http.MethodGet, "/api/v3/log/file/update/update.txt", "update log")
	defer ts.Close()
	c, _ := whisparr.NewV3(ts.URL, "test-key")
	got, err := c.GetUpdateLogFileContent(context.Background(), "update.txt")
	if err != nil {
		t.Fatalf("GetUpdateLogFileContent: %v", err)
	}
	if got != "update log" {
		t.Errorf("content = %q, want %q", got, "update log")
	}
}
