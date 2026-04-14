package nzbhydra_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golusoris/goenvoy/arr/nzbhydra"
)

const newznabXML = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:newznab="http://www.newznab.com/DTD/2010/feeds/attributes/">
  <channel>
    <item>
      <title>Test.NZB.Release-GRP</title>
      <guid>https://example.com/details/abc</guid>
      <link>https://example.com/getnzb/abc</link>
      <comments>https://example.com/details/abc#comments</comments>
      <pubDate>Mon, 01 Jan 2025 00:00:00 +0000</pubDate>
      <size>1073741824</size>
      <description>A test release</description>
      <newznab:attr name="category" value="2000"/>
      <newznab:attr name="indexer" value="NZBgeek"/>
    </item>
    <item>
      <title>Another.Release-GRP</title>
      <guid>https://example.com/details/def</guid>
      <link>https://example.com/getnzb/def</link>
      <pubDate>Tue, 02 Jan 2025 12:00:00 +0000</pubDate>
      <size>536870912</size>
      <description>Another release</description>
      <newznab:attr name="category" value="5000"/>
    </item>
  </channel>
</rss>`

const capsXML = `<?xml version="1.0" encoding="UTF-8"?>
<caps>
  <server title="NZBHydra2" image="https://example.com/logo.png"/>
  <limits max="100" default="50"/>
  <searching>
    <search available="yes"/>
    <tv-search available="yes"/>
    <movie-search available="yes"/>
    <book-search available="no"/>
  </searching>
  <categories>
    <category id="2000" name="Movies">
      <subcat id="2010" name="Movies/Foreign"/>
    </category>
    <category id="5000" name="TV"/>
  </categories>
</caps>`

func newXMLServer(t *testing.T, xmlBody string) *nzbhydra.Client {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("apikey"); got != "test-key" {
			t.Errorf("apikey = %q, want test-key", got)
		}
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(xmlBody))
	}))
	t.Cleanup(ts.Close)
	return nzbhydra.New(ts.URL, "test-key")
}

func newJSONPostServer(t *testing.T, wantPath string, response any) *nzbhydra.Client {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if wantPath != "" && r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if got := r.URL.Query().Get("apikey"); got != "test-key" {
			t.Errorf("apikey = %q, want test-key", got)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "application/json") {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		w.Header().Set("Content-Type", "application/json")
		if response != nil {
			_ = json.NewEncoder(w).Encode(response)
		}
	}))
	t.Cleanup(ts.Close)
	return nzbhydra.New(ts.URL, "test-key")
}

func newJSONGetServer(t *testing.T, wantPath string, response any) *nzbhydra.Client {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if wantPath != "" && r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if got := r.URL.Query().Get("apikey"); got != "test-key" {
			t.Errorf("apikey = %q, want test-key", got)
		}
		w.Header().Set("Content-Type", "application/json")
		if response != nil {
			_ = json.NewEncoder(w).Encode(response)
		}
	}))
	t.Cleanup(ts.Close)
	return nzbhydra.New(ts.URL, "test-key")
}

func TestSearch(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("t"); got != "search" {
			t.Errorf("t = %q, want search", got)
		}
		if got := q.Get("q"); got != "ubuntu" {
			t.Errorf("q = %q, want ubuntu", got)
		}
		if got := q.Get("cat"); got != "2000,5000" {
			t.Errorf("cat = %q, want 2000,5000", got)
		}
		if r.URL.Path != "/api" {
			t.Errorf("path = %q, want /api", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(newznabXML))
	}))
	defer ts.Close()

	c := nzbhydra.New(ts.URL, "test-key")
	results, err := c.Search(context.Background(), "ubuntu", []int{2000, 5000})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	r := results[0]
	if r.Title != "Test.NZB.Release-GRP" {
		t.Errorf("Title = %q, want Test.NZB.Release-GRP", r.Title)
	}
	if r.Size != 1073741824 {
		t.Errorf("Size = %d, want 1073741824", r.Size)
	}
	if r.Category != "2000" {
		t.Errorf("Category = %q, want 2000", r.Category)
	}
	if r.Indexer != "NZBgeek" {
		t.Errorf("Indexer = %q, want NZBgeek", r.Indexer)
	}
	if r.Description != "A test release" {
		t.Errorf("Description = %q, want A test release", r.Description)
	}
}

func TestSearchNilCategories(t *testing.T) {
	t.Parallel()

	c := newXMLServer(t, newznabXML)
	results, err := c.Search(context.Background(), "test", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
}

func TestTVSearch(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("t"); got != "tvsearch" {
			t.Errorf("t = %q, want tvsearch", got)
		}
		if got := q.Get("season"); got != "3" {
			t.Errorf("season = %q, want 3", got)
		}
		if got := q.Get("ep"); got != "5" {
			t.Errorf("ep = %q, want 5", got)
		}
		if got := q.Get("tvdbid"); got != "12345" {
			t.Errorf("tvdbid = %q, want 12345", got)
		}
		if got := q.Get("imdbid"); got != "tt0903747" {
			t.Errorf("imdbid = %q, want tt0903747", got)
		}
		if got := q.Get("tmdbid"); got != "1396" {
			t.Errorf("tmdbid = %q, want 1396", got)
		}
		if got := q.Get("tvmazeid"); got != "169" {
			t.Errorf("tvmazeid = %q, want 169", got)
		}
		if got := q.Get("rid"); got != "999" {
			t.Errorf("rid = %q, want 999", got)
		}
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(newznabXML))
	}))
	defer ts.Close()

	c := nzbhydra.New(ts.URL, "test-key")
	results, err := c.TVSearch(context.Background(), "breaking bad", &nzbhydra.TVSearchOptions{
		Season:   3,
		Episode:  5,
		TVDBID:   "12345",
		IMDBID:   "tt0903747",
		TMDBID:   "1396",
		TVMazeID: "169",
		RID:      "999",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
}

func TestTVSearchMinimalParams(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Has("season") {
			t.Error("season should not be set")
		}
		if q.Has("ep") {
			t.Error("ep should not be set")
		}
		if q.Has("tvdbid") {
			t.Error("tvdbid should not be set")
		}
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(newznabXML))
	}))
	defer ts.Close()

	c := nzbhydra.New(ts.URL, "test-key")
	_, err := c.TVSearch(context.Background(), "test", &nzbhydra.TVSearchOptions{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestMovieSearch(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("t"); got != "movie" {
			t.Errorf("t = %q, want movie", got)
		}
		if got := q.Get("imdbid"); got != "tt1375666" {
			t.Errorf("imdbid = %q, want tt1375666", got)
		}
		if got := q.Get("tmdbid"); got != "27205" {
			t.Errorf("tmdbid = %q, want 27205", got)
		}
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(newznabXML))
	}))
	defer ts.Close()

	c := nzbhydra.New(ts.URL, "test-key")
	results, err := c.MovieSearch(context.Background(), "inception", "tt1375666", "27205")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
}

func TestBookSearch(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("t"); got != "book" {
			t.Errorf("t = %q, want book", got)
		}
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(newznabXML))
	}))
	defer ts.Close()

	c := nzbhydra.New(ts.URL, "test-key")
	_, err := c.BookSearch(context.Background(), "dune")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCapabilities(t *testing.T) {
	t.Parallel()

	c := newXMLServer(t, capsXML)
	caps, err := c.GetCapabilities(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if caps.Server.Title != "NZBHydra2" {
		t.Errorf("Server.Title = %q, want NZBHydra2", caps.Server.Title)
	}
	if caps.Limits.Max != 100 {
		t.Errorf("Limits.Max = %d, want 100", caps.Limits.Max)
	}
	if !caps.Searching.SearchAvailable {
		t.Error("SearchAvailable = false, want true")
	}
	if !caps.Searching.TVSearchAvailable {
		t.Error("TVSearchAvailable = false, want true")
	}
	if !caps.Searching.MovieSearchAvailable {
		t.Error("MovieSearchAvailable = false, want true")
	}
	if caps.Searching.BookSearchAvailable {
		t.Error("BookSearchAvailable = true, want false")
	}
	if len(caps.Categories) != 2 {
		t.Fatalf("len(Categories) = %d, want 2", len(caps.Categories))
	}
	if caps.Categories[0].ID != 2000 {
		t.Errorf("Categories[0].ID = %d, want 2000", caps.Categories[0].ID)
	}
	if len(caps.Categories[0].SubCategories) != 1 {
		t.Fatalf("len(SubCategories) = %d, want 1", len(caps.Categories[0].SubCategories))
	}
	if caps.Categories[0].SubCategories[0].Name != "Movies/Foreign" {
		t.Errorf("SubCategories[0].Name = %q, want Movies/Foreign", caps.Categories[0].SubCategories[0].Name)
	}
}

func TestGetStats(t *testing.T) {
	t.Parallel()

	stats := map[string]any{
		"avgResponseTimes":      []map[string]any{{"indexer": "NZBgeek", "avgResponseTime": 200}},
		"indexerApiAccessStats": []map[string]any{{"indexer": "NZBgeek", "successful": 100, "unsuccessful": 5}},
		"searchesPerDayOfWeek":  map[string]any{"Monday": 10, "Tuesday": 15},
		"downloadsPerDayOfWeek": map[string]any{"Monday": 5},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/api/stats" {
			t.Errorf("path = %q, want /api/stats", r.URL.Path)
		}

		// Verify request body is valid JSON.
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var req map[string]any
		if err := json.Unmarshal(body, &req); err != nil {
			t.Errorf("invalid request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(stats)
	}))
	defer ts.Close()

	c := nzbhydra.New(ts.URL, "test-key")
	result, err := c.GetStats(context.Background(), nzbhydra.StatsRequest{
		IncludeAverageResponseTimes: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.AvgResponseTimes) != 1 {
		t.Fatalf("len(AvgResponseTimes) = %d, want 1", len(result.AvgResponseTimes))
	}
	if result.AvgResponseTimes[0].Indexer != "NZBgeek" {
		t.Errorf("Indexer = %q, want NZBgeek", result.AvgResponseTimes[0].Indexer)
	}
	if result.AvgResponseTimes[0].Avg != 200 {
		t.Errorf("Avg = %d, want 200", result.AvgResponseTimes[0].Avg)
	}
	if len(result.IndexerAPIAccessStats) != 1 {
		t.Fatalf("len(IndexerAPIAccessStats) = %d, want 1", len(result.IndexerAPIAccessStats))
	}
	if result.IndexerAPIAccessStats[0].Successful != 100 {
		t.Errorf("Successful = %d, want 100", result.IndexerAPIAccessStats[0].Successful)
	}
}

func TestGetSearchHistory(t *testing.T) {
	t.Parallel()

	resp := map[string]any{
		"content": []map[string]any{
			{"id": 1, "source": "INTERNAL", "searchType": "SEARCH", "time": "2025-01-01T00:00:00Z", "query": "test", "ip": "127.0.0.1"},
			{"id": 2, "source": "API", "searchType": "TVSEARCH", "time": "2025-01-02T00:00:00Z", "query": "breaking bad", "season": "5", "episode": "3"},
		},
		"totalElements":    2,
		"totalPages":       1,
		"first":            true,
		"last":             true,
		"numberOfElements": 2,
		"number":           0,
		"size":             25,
	}
	c := newJSONPostServer(t, "/api/history/searches", resp)

	result, err := c.GetSearchHistory(context.Background(), nzbhydra.HistoryRequest{
		Page: 0, Limit: 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Content) != 2 {
		t.Fatalf("len(Content) = %d, want 2", len(result.Content))
	}
	if result.Content[0].Source != "INTERNAL" {
		t.Errorf("Source = %q, want INTERNAL", result.Content[0].Source)
	}
	if result.Content[1].Query != "breaking bad" {
		t.Errorf("Query = %q, want breaking bad", result.Content[1].Query)
	}
	if result.TotalElements != 2 {
		t.Errorf("TotalElements = %d, want 2", result.TotalElements)
	}
	if !result.First {
		t.Error("First = false, want true")
	}
	if !result.Last {
		t.Error("Last = false, want true")
	}
}

func TestGetDownloadHistory(t *testing.T) {
	t.Parallel()

	resp := map[string]any{
		"content": []map[string]any{
			{
				"id":            1,
				"searchResult":  map[string]any{"title": "Test.Release", "indexer": "NZBgeek", "link": "https://example.com/nzb"},
				"nzbAccessType": "PROXY",
				"accessSource":  "INTERNAL",
				"time":          "2025-01-01T00:00:00Z",
				"status":        "CONTENT_DOWNLOAD_SUCCESSFUL",
				"age":           5,
			},
		},
		"totalElements":    1,
		"totalPages":       1,
		"first":            true,
		"last":             true,
		"numberOfElements": 1,
		"number":           0,
		"size":             25,
	}
	c := newJSONPostServer(t, "/api/history/downloads", resp)

	result, err := c.GetDownloadHistory(context.Background(), nzbhydra.HistoryRequest{
		Page: 0, Limit: 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("len(Content) = %d, want 1", len(result.Content))
	}
	if result.Content[0].SearchResult.Title != "Test.Release" {
		t.Errorf("Title = %q, want Test.Release", result.Content[0].SearchResult.Title)
	}
	if result.Content[0].NzbAccessType != "PROXY" {
		t.Errorf("NzbAccessType = %q, want PROXY", result.Content[0].NzbAccessType)
	}
	if result.Content[0].Status != "CONTENT_DOWNLOAD_SUCCESSFUL" {
		t.Errorf("Status = %q, want CONTENT_DOWNLOAD_SUCCESSFUL", result.Content[0].Status)
	}
	if result.Content[0].Age != 5 {
		t.Errorf("Age = %d, want 5", result.Content[0].Age)
	}
}

func TestGetIndexerStatuses(t *testing.T) {
	t.Parallel()

	statuses := []map[string]any{
		{"indexer": "NZBgeek", "state": "ENABLED", "level": "NORMAL", "disabledUntil": "", "lastError": ""},
		{"indexer": "DrunkenSlug", "state": "DISABLED", "level": "ERROR", "disabledUntil": "2025-01-02T00:00:00Z", "lastError": "Connection timeout"},
	}
	c := newJSONGetServer(t, "/api/stats/indexers", statuses)

	result, err := c.GetIndexerStatuses(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("len(result) = %d, want 2", len(result))
	}
	if result[0].Indexer != "NZBgeek" {
		t.Errorf("Indexer = %q, want NZBgeek", result[0].Indexer)
	}
	if result[0].State != "ENABLED" {
		t.Errorf("State = %q, want ENABLED", result[0].State)
	}
	if result[1].State != "DISABLED" {
		t.Errorf("State = %q, want DISABLED", result[1].State)
	}
	if result[1].LastError != "Connection timeout" {
		t.Errorf("LastError = %q, want Connection timeout", result[1].LastError)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized"))
	}))
	defer ts.Close()

	c := nzbhydra.New(ts.URL, "bad-key")
	_, err := c.Search(context.Background(), "test", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *nzbhydra.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
	if apiErr.Body != "Unauthorized" {
		t.Errorf("Body = %q, want Unauthorized", apiErr.Body)
	}
}

func TestAPIErrorNoBody(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := nzbhydra.New(ts.URL, "key")
	_, err := c.GetIndexerStatuses(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *nzbhydra.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if !strings.Contains(apiErr.Error(), "HTTP 500") {
		t.Errorf("error = %q, want HTTP 500 substring", apiErr.Error())
	}
}

func TestPostAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Forbidden"))
	}))
	defer ts.Close()

	c := nzbhydra.New(ts.URL, "bad-key")
	_, err := c.GetStats(context.Background(), nzbhydra.StatsRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *nzbhydra.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusForbidden)
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

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(newznabXML))
	}))
	defer ts.Close()

	c := nzbhydra.New(ts.URL, "k", nzbhydra.WithHTTPClient(custom))
	_, _ = c.Search(context.Background(), "test", nil)
	if !called {
		t.Error("custom HTTP client was not used")
	}
}

func TestWithTimeout(t *testing.T) {
	t.Parallel()

	c := nzbhydra.New("http://localhost:5076", "key", nzbhydra.WithTimeout(5*1e9))
	_ = c
}

func TestEmptySearchResults(t *testing.T) {
	t.Parallel()

	emptyXML := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"><channel></channel></rss>`
	c := newXMLServer(t, emptyXML)
	results, err := c.Search(context.Background(), "nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("len(results) = %d, want 0", len(results))
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
