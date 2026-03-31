package googlebooks_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/metadata/book/googlebooks"
)

func setup(t *testing.T, handler http.HandlerFunc) *googlebooks.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return googlebooks.New("test-key", googlebooks.WithBaseURL(srv.URL))
}

func TestSearch(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q"); got != "flowers for algernon" {
			t.Errorf("q = %q, want flowers for algernon", got)
		}
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Errorf("key = %q, want test-key", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.VolumesResponse{
			Kind:       "books#volumes",
			TotalItems: 1,
			Items: []googlebooks.Volume{
				{
					Id:   "abc123",
					Kind: "books#volume",
					VolumeInfo: &googlebooks.VolumeInfo{
						Title:   "Flowers for Algernon",
						Authors: []string{"Daniel Keyes"},
					},
				},
			},
		})
	})

	resp, err := c.Search(context.Background(), "flowers for algernon")
	if err != nil {
		t.Fatal(err)
	}
	if resp.TotalItems != 1 {
		t.Fatalf("TotalItems = %d, want 1", resp.TotalItems)
	}
	if resp.Items[0].VolumeInfo.Title != "Flowers for Algernon" {
		t.Errorf("Title = %q", resp.Items[0].VolumeInfo.Title)
	}
}

func TestSearchWithParams(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q"); got != "science fiction" {
			t.Errorf("q = %q, want science fiction", got)
		}
		if got := r.URL.Query().Get("maxResults"); got != "5" {
			t.Errorf("maxResults = %q, want 5", got)
		}
		if got := r.URL.Query().Get("orderBy"); got != "relevance" {
			t.Errorf("orderBy = %q, want relevance", got)
		}
		if got := r.URL.Query().Get("langRestrict"); got != "en" {
			t.Errorf("langRestrict = %q, want en", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.VolumesResponse{
			TotalItems: 2,
			Items: []googlebooks.Volume{
				{Id: "a", VolumeInfo: &googlebooks.VolumeInfo{Title: "Book A"}},
				{Id: "b", VolumeInfo: &googlebooks.VolumeInfo{Title: "Book B"}},
			},
		})
	})

	resp, err := c.SearchWithParams(context.Background(), &googlebooks.SearchParams{
		Query:        "science fiction",
		MaxResults:   5,
		OrderBy:      "relevance",
		LangRestrict: "en",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.TotalItems != 2 {
		t.Fatalf("TotalItems = %d, want 2", resp.TotalItems)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2", len(resp.Items))
	}
}

func TestGetVolume(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/volumes/abc123" {
			t.Errorf("path = %q, want /volumes/abc123", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.Volume{
			Id:   "abc123",
			Kind: "books#volume",
			VolumeInfo: &googlebooks.VolumeInfo{
				Title:     "Flowers for Algernon",
				Authors:   []string{"Daniel Keyes"},
				PageCount: 311,
				ImageLinks: &googlebooks.ImageLinks{
					Thumbnail: "https://example.com/thumb.jpg",
				},
			},
			SaleInfo: &googlebooks.SaleInfo{
				Country:     "US",
				Saleability: "FOR_SALE",
				IsEbook:     true,
			},
		})
	})

	vol, err := c.GetVolume(context.Background(), "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if vol.Id != "abc123" {
		t.Errorf("Id = %q, want abc123", vol.Id)
	}
	if vol.VolumeInfo.PageCount != 311 {
		t.Errorf("PageCount = %d, want 311", vol.VolumeInfo.PageCount)
	}
	if vol.VolumeInfo.ImageLinks == nil || vol.VolumeInfo.ImageLinks.Thumbnail != "https://example.com/thumb.jpg" {
		t.Error("ImageLinks not parsed correctly")
	}
	if vol.SaleInfo == nil || !vol.SaleInfo.IsEbook {
		t.Error("SaleInfo not parsed correctly")
	}
}

func TestAPIKeyIsSent(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Errorf("key = %q, want test-key", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.VolumesResponse{TotalItems: 0})
	})

	_, err := c.Search(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":{"code":403,"message":"forbidden"}}`))
	}))
	defer srv.Close()

	c := googlebooks.New("bad-key", googlebooks.WithBaseURL(srv.URL))
	_, err := c.Search(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *googlebooks.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusForbidden)
	}
}

func TestAPIErrorOnGetVolume(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	}))
	defer srv.Close()

	c := googlebooks.New("key", googlebooks.WithBaseURL(srv.URL))
	_, err := c.GetVolume(context.Background(), "invalid")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *googlebooks.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusNotFound)
	}
}

func TestNew(t *testing.T) {
	c := googlebooks.New("key")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}
