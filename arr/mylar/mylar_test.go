package mylar_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/arr/mylar/v2"
)

func newTestServer(t *testing.T, _ string, response any) *mylar.Client {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Query validation is done in individual tests via the full-handler variant.
		w.Header().Set("Content-Type", "application/json")
		if response != nil {
			_ = json.NewEncoder(w).Encode(response)
		}
	}))
	t.Cleanup(ts.Close)
	return mylar.New(ts.URL, "test-key")
}

func newTestServerFull(t *testing.T, wantCmd, wantKey string, response any) *mylar.Client {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("apikey"); got != wantKey {
			t.Errorf("apikey = %q, want %q", got, wantKey)
		}
		if got := q.Get("cmd"); got != wantCmd {
			t.Errorf("cmd = %q, want %q", got, wantCmd)
		}
		w.Header().Set("Content-Type", "application/json")
		if response != nil {
			_ = json.NewEncoder(w).Encode(response)
		}
	}))
	t.Cleanup(ts.Close)
	return mylar.New(ts.URL, wantKey)
}

func TestGetIndex(t *testing.T) {
	t.Parallel()

	c := newTestServerFull(t, "getIndex", "test-key", []map[string]any{
		{"id": "1", "name": "Spider-Man", "status": "Active", "totalIssues": 50},
		{"id": "2", "name": "Batman", "status": "Ended", "totalIssues": 100},
	})

	comics, err := c.GetIndex(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(comics) != 2 {
		t.Fatalf("len(comics) = %d, want 2", len(comics))
	}
	if comics[0].Name != "Spider-Man" {
		t.Errorf("comics[0].Name = %q, want Spider-Man", comics[0].Name)
	}
	if comics[1].TotalIssues != 100 {
		t.Errorf("comics[1].TotalIssues = %d, want 100", comics[1].TotalIssues)
	}
}

func TestGetComic(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "getComic", map[string]any{
		"id": "42", "name": "X-Men", "status": "Active", "year": "1991",
	})

	comic, err := c.GetComic(context.Background(), "42")
	if err != nil {
		t.Fatal(err)
	}
	if comic.Name != "X-Men" {
		t.Errorf("Name = %q, want X-Men", comic.Name)
	}
	if comic.Year != "1991" {
		t.Errorf("Year = %q, want 1991", comic.Year)
	}
}

func TestGetUpcoming(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "getUpcoming", []map[string]any{
		{"id": "1", "comicName": "Spider-Man", "issueNumber": "51"},
	})

	upcoming, err := c.GetUpcoming(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(upcoming) != 1 {
		t.Fatalf("len(upcoming) = %d, want 1", len(upcoming))
	}
	if upcoming[0].ComicName != "Spider-Man" {
		t.Errorf("ComicName = %q, want Spider-Man", upcoming[0].ComicName)
	}
}

func TestGetWanted(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "getWanted", []map[string]any{
		{"id": "10", "issueName": "Issue 5", "status": "Wanted"},
	})

	wanted, err := c.GetWanted(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(wanted) != 1 {
		t.Fatalf("len(wanted) = %d, want 1", len(wanted))
	}
	if wanted[0].Status != "Wanted" {
		t.Errorf("Status = %q, want Wanted", wanted[0].Status)
	}
}

func TestGetHistory(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "getHistory", []map[string]any{
		{"id": "1", "comicName": "Batman", "issueNumber": "10", "status": "Downloaded"},
	})

	history, err := c.GetHistory(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 {
		t.Fatalf("len(history) = %d, want 1", len(history))
	}
	if history[0].ComicName != "Batman" {
		t.Errorf("ComicName = %q, want Batman", history[0].ComicName)
	}
}

func TestGetLogs(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "getLogs", []map[string]any{
		{"message": "Started scan", "level": "INFO", "timestamp": "2025-01-01T00:00:00Z"},
	})

	logs, err := c.GetLogs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 1 {
		t.Fatalf("len(logs) = %d, want 1", len(logs))
	}
	if logs[0].Level != "INFO" {
		t.Errorf("Level = %q, want INFO", logs[0].Level)
	}
}

func TestFindComic(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "findComic", []map[string]any{
		{"id": "100", "name": "Saga", "year": "2012", "publisher": "Image", "issues": 66},
	})

	results, err := c.FindComic(context.Background(), "Saga")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].Publisher != "Image" {
		t.Errorf("Publisher = %q, want Image", results[0].Publisher)
	}
}

func TestAddComic(t *testing.T) {
	t.Parallel()

	c := newTestServerFull(t, "addComic", "test-key", nil)
	if err := c.AddComic(context.Background(), "99"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteComic(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "delComic", nil)
	if err := c.DeleteComic(context.Background(), "99"); err != nil {
		t.Fatal(err)
	}
}

func TestPauseComic(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "pauseComic", nil)
	if err := c.PauseComic(context.Background(), "42"); err != nil {
		t.Fatal(err)
	}
}

func TestResumeComic(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "resumeComic", nil)
	if err := c.ResumeComic(context.Background(), "42"); err != nil {
		t.Fatal(err)
	}
}

func TestGetVersion(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "getVersion", map[string]any{
		"version": "0.7.0", "latestVersion": "0.7.1", "commits": "abc123",
	})

	ver, err := c.GetVersion(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if ver.Version != "0.7.0" {
		t.Errorf("Version = %q, want 0.7.0", ver.Version)
	}
}

func TestGetReadList(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "getReadList", []map[string]any{
		{"id": "1", "name": "Summer Reading", "issues": []map[string]any{
			{"id": "10", "issueNumber": "1"},
		}},
	})

	lists, err := c.GetReadList(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(lists) != 1 {
		t.Fatalf("len(lists) = %d, want 1", len(lists))
	}
	if lists[0].Name != "Summer Reading" {
		t.Errorf("Name = %q, want Summer Reading", lists[0].Name)
	}
}

func TestGetStoryArc(t *testing.T) {
	t.Parallel()

	c := newTestServer(t, "getStoryArc", map[string]any{
		"id": "5", "name": "Civil War", "publisher": "Marvel",
		"issues": []map[string]any{{"id": "20", "issueNumber": "1"}},
	})

	arc, err := c.GetStoryArc(context.Background(), "5")
	if err != nil {
		t.Fatal(err)
	}
	if arc.Name != "Civil War" {
		t.Errorf("Name = %q, want Civil War", arc.Name)
	}
	if len(arc.Issues) != 1 {
		t.Fatalf("len(Issues) = %d, want 1", len(arc.Issues))
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized"))
	}))
	defer ts.Close()

	c := mylar.New(ts.URL, "bad-key")
	_, err := c.GetIndex(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *mylar.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
}

func TestQueryParams(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("apikey"); got != "my-key" {
			t.Errorf("apikey = %q, want my-key", got)
		}
		if got := q.Get("cmd"); got != "getComic" {
			t.Errorf("cmd = %q, want getComic", got)
		}
		if got := q.Get("id"); got != "77" {
			t.Errorf("id = %q, want 77", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "77", "name": "Test"})
	}))
	defer ts.Close()

	c := mylar.New(ts.URL, "my-key")
	comic, err := c.GetComic(context.Background(), "77")
	if err != nil {
		t.Fatal(err)
	}
	if comic.ID != "77" {
		t.Errorf("ID = %q, want 77", comic.ID)
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
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]any{})
	}))
	defer ts.Close()

	c := mylar.New(ts.URL, "k", mylar.WithHTTPClient(custom))
	_, _ = c.GetIndex(context.Background())
	if !called {
		t.Error("custom HTTP client was not used")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
