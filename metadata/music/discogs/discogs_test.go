package discogs

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return New("test-token", WithBaseURL(ts.URL))
}

func TestGetRelease(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Discogs token=test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(Release{ID: 249504, Title: "Nevermind", Year: 1991})
	})

	rel, err := c.GetRelease(context.Background(), 249504)
	if err != nil {
		t.Fatal(err)
	}
	if rel.Title != "Nevermind" || rel.Year != 1991 {
		t.Fatalf("unexpected release: %+v", rel)
	}
}

func TestGetArtist(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Artist{ID: 125246, Name: "Nirvana"})
	})

	a, err := c.GetArtist(context.Background(), 125246)
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "Nirvana" {
		t.Fatalf("unexpected artist: %+v", a)
	}
}

func TestGetArtistReleases(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(SearchResponse{
			Pagination: Pagination{Page: 1, Pages: 1, Items: 1},
			Results:    []SearchResult{{ID: 249504, Title: "Nevermind", Type: "release"}},
		})
	})

	sr, err := c.GetArtistReleases(context.Background(), 125246, 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(sr.Results) != 1 || sr.Results[0].Title != "Nevermind" {
		t.Fatalf("unexpected results: %+v", sr)
	}
}

func TestGetLabel(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Label{ID: 1, Name: "Planet E"})
	})

	l, err := c.GetLabel(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if l.Name != "Planet E" {
		t.Fatalf("unexpected label: %+v", l)
	}
}

func TestGetMasterRelease(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(MasterRelease{ID: 1000, Title: "Nevermind", Year: 1991})
	})

	m, err := c.GetMasterRelease(context.Background(), 1000)
	if err != nil {
		t.Fatal(err)
	}
	if m.Title != "Nevermind" {
		t.Fatalf("unexpected master: %+v", m)
	}
}

func TestSearch(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(SearchResponse{
			Pagination: Pagination{Page: 1, Pages: 1, Items: 1},
			Results:    []SearchResult{{ID: 249504, Title: "Nevermind", Type: "release"}},
		})
	})

	sr, err := c.Search(context.Background(), "nevermind", "release", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(sr.Results) != 1 || sr.Results[0].Title != "Nevermind" {
		t.Fatalf("unexpected results: %+v", sr)
	}
}

func TestAPIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	_, err := c.GetRelease(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected status: %d", apiErr.StatusCode)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c := New("token", WithHTTPClient(custom))
	if c.http != custom {
		t.Fatal("custom HTTP client not set")
	}
}

func TestWithUserAgent(t *testing.T) {
	c := New("token", WithUserAgent("myapp/2.0"))
	if c.userAgent != "myapp/2.0" {
		t.Fatal("user agent not set")
	}
}
