package tautulli_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/mediaserver/tautulli"
)

type apiResp struct {
	Response struct {
		Result  string `json:"result"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	} `json:"response"`
}

func newTestServer(t *testing.T, wantCmd, wantKey string, data any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2" {
			t.Errorf("path = %q, want /api/v2", r.URL.Path)
		}
		if got := r.URL.Query().Get("cmd"); got != wantCmd {
			t.Errorf("cmd = %q, want %q", got, wantCmd)
		}
		if got := r.URL.Query().Get("apikey"); got != wantKey {
			t.Errorf("apikey = %q, want %q", got, wantKey)
		}
		resp := apiResp{}
		resp.Response.Result = "success"
		resp.Response.Data = data
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func TestGetActivity(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_activity", "test-key", map[string]any{
		"stream_count": "2",
		"sessions":     []any{},
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	a, err := c.GetActivity(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if a.StreamCount != "2" {
		t.Errorf("StreamCount = %q, want 2", a.StreamCount)
	}
}

func TestGetHistory(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_history", "test-key", map[string]any{
		"recordsTotal":    100,
		"recordsFiltered": 10,
		"total_duration":  "5 hrs",
		"data":            []any{},
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	h, err := c.GetHistory(context.Background(), 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if h.RecordsTotal != 100 {
		t.Errorf("RecordsTotal = %d, want 100", h.RecordsTotal)
	}
}

func TestGetLibraries(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_libraries", "test-key", []map[string]any{
		{"section_id": "1", "section_name": "Movies", "section_type": "movie"},
		{"section_id": "2", "section_name": "TV Shows", "section_type": "show"},
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	libs, err := c.GetLibraries(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(libs) != 2 {
		t.Fatalf("len(libraries) = %d, want 2", len(libs))
	}
	if libs[0].SectionName != "Movies" {
		t.Errorf("libraries[0].SectionName = %q, want Movies", libs[0].SectionName)
	}
}

func TestGetLibrary(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_library", "test-key", map[string]any{
		"section_id":   "1",
		"section_name": "Movies",
		"section_type": "movie",
		"count":        887,
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	lib, err := c.GetLibrary(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if lib.SectionName != "Movies" {
		t.Errorf("SectionName = %q, want Movies", lib.SectionName)
	}
}

func TestGetUsers(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_users", "test-key", []map[string]any{
		{"user_id": "133788", "username": "JonSnow"},
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	users, err := c.GetUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatalf("len(users) = %d, want 1", len(users))
	}
	if users[0].Username != "JonSnow" {
		t.Errorf("Username = %q, want JonSnow", users[0].Username)
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_user", "test-key", map[string]any{
		"user_id":       133788,
		"username":      "JonSnow",
		"friendly_name": "Jon Snow",
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	u, err := c.GetUser(context.Background(), "133788")
	if err != nil {
		t.Fatal(err)
	}
	if u.FriendlyName != "Jon Snow" {
		t.Errorf("FriendlyName = %q, want Jon Snow", u.FriendlyName)
	}
}

func TestGetServerInfo(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_server_info", "test-key", map[string]any{
		"pms_name":    "Winterfell-Server",
		"pms_version": "1.20.0.3133",
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	info, err := c.GetServerInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.PMSName != "Winterfell-Server" {
		t.Errorf("PMSName = %q, want Winterfell-Server", info.PMSName)
	}
}

func TestGetTautulliInfo(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_tautulli_info", "test-key", map[string]any{
		"tautulli_version": "v2.8.1",
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	info, err := c.GetTautulliInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.Version != "v2.8.1" {
		t.Errorf("Version = %q, want v2.8.1", info.Version)
	}
}

func TestGetMetadata(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_metadata", "test-key", map[string]any{
		"title":      "The Red Woman",
		"media_type": "episode",
		"year":       "2016",
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	m, err := c.GetMetadata(context.Background(), "153037")
	if err != nil {
		t.Fatal(err)
	}
	if m.Title != "The Red Woman" {
		t.Errorf("Title = %q, want The Red Woman", m.Title)
	}
}

func TestGetRecentlyAdded(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_recently_added", "test-key", map[string]any{
		"recently_added": []map[string]any{
			{"title": "Deadpool", "media_type": "movie"},
		},
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	r, err := c.GetRecentlyAdded(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.RecentlyAdded) != 1 {
		t.Fatalf("len(recently_added) = %d, want 1", len(r.RecentlyAdded))
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "search", "test-key", map[string]any{
		"results_count": 5,
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	s, err := c.Search(context.Background(), "Thrones")
	if err != nil {
		t.Fatal(err)
	}
	if s.ResultsCount != 5 {
		t.Errorf("ResultsCount = %d, want 5", s.ResultsCount)
	}
}

func TestServerStatus(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "server_status", "test-key", map[string]any{
		"connected": true,
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	s, err := c.ServerStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !s.Connected {
		t.Error("Connected = false, want true")
	}
}

func TestTerminateSession(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "terminate_session", "test-key", nil)
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	if err := c.TerminateSession(context.Background(), "27", "test"); err != nil {
		t.Fatal(err)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Invalid API key"))
	}))
	defer ts.Close()

	c := tautulli.New(ts.URL, "bad-key")
	_, err := c.GetActivity(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *tautulli.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
}

func TestAPIErrorResult(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := apiResp{}
		resp.Response.Result = "error"
		resp.Response.Message = "something went wrong"
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	_, err := c.GetActivity(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *tautulli.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
}

func TestGetGeoIPLookup(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_geoip_lookup", "test-key", map[string]any{
		"city":    "Mountain View",
		"country": "United States",
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	g, err := c.GetGeoIPLookup(context.Background(), "8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	if g.City != "Mountain View" {
		t.Errorf("City = %q, want Mountain View", g.City)
	}
}

func TestGetUserWatchTimeStats(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, "get_user_watch_time_stats", "test-key", []map[string]any{
		{"query_days": 7, "total_plays": 3, "total_time": 15694},
	})
	defer ts.Close()

	c := tautulli.New(ts.URL, "test-key")
	stats, err := c.GetUserWatchTimeStats(context.Background(), "133788")
	if err != nil {
		t.Fatal(err)
	}
	if len(stats) != 1 {
		t.Fatalf("len(stats) = %d, want 1", len(stats))
	}
	if stats[0].TotalPlays != 3 {
		t.Errorf("TotalPlays = %d, want 3", stats[0].TotalPlays)
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	called := false
	custom := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			called = true
			return http.DefaultTransport.RoundTrip(r)
		}),
	}

	ts := newTestServer(t, "get_activity", "k", map[string]any{"stream_count": "0"})
	defer ts.Close()

	c := tautulli.New(ts.URL, "k", tautulli.WithHTTPClient(custom))
	_, _ = c.GetActivity(context.Background())
	if !called {
		t.Error("custom HTTP client was not used")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
