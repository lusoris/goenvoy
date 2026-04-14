package googlebooks_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/book/googlebooks/v2"
)

func setup(t *testing.T, handler http.HandlerFunc) *googlebooks.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return googlebooks.New("test-key", metadata.WithBaseURL(srv.URL))
}

func TestSearch(t *testing.T) {
	t.Parallel()

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
					ID:   "abc123",
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
	t.Parallel()

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
				{ID: "a", VolumeInfo: &googlebooks.VolumeInfo{Title: "Book A"}},
				{ID: "b", VolumeInfo: &googlebooks.VolumeInfo{Title: "Book B"}},
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
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/volumes/abc123" {
			t.Errorf("path = %q, want /volumes/abc123", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.Volume{
			ID:   "abc123",
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
	if vol.ID != "abc123" {
		t.Errorf("ID = %q, want abc123", vol.ID)
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
	t.Parallel()

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
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":{"code":403,"message":"forbidden"}}`))
	}))
	defer srv.Close()

	c := googlebooks.New("bad-key", metadata.WithBaseURL(srv.URL))
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
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	}))
	defer srv.Close()

	c := googlebooks.New("key", metadata.WithBaseURL(srv.URL))
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
	t.Parallel()

	c := googlebooks.New("key")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGetUserBookshelves(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/user123/bookshelves" {
			t.Errorf("path = %q, want /users/user123/bookshelves", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.BookshelvesResponse{
			Kind: "books#bookshelves",
			Items: []googlebooks.Bookshelf{
				{ID: 0, Title: "Favorites", Access: "PUBLIC", VolumeCount: 5},
				{ID: 3, Title: "Reading now", Access: "PUBLIC", VolumeCount: 2},
			},
		})
	})

	resp, err := c.GetUserBookshelves(context.Background(), "user123")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2", len(resp.Items))
	}
	if resp.Items[0].Title != "Favorites" {
		t.Errorf("Title = %q, want Favorites", resp.Items[0].Title)
	}
	if resp.Items[1].VolumeCount != 2 {
		t.Errorf("VolumeCount = %d, want 2", resp.Items[1].VolumeCount)
	}
}

func TestGetUserBookshelf(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/user123/bookshelves/0" {
			t.Errorf("path = %q, want /users/user123/bookshelves/0", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.Bookshelf{
			ID:          0,
			Title:       "Favorites",
			Access:      "PUBLIC",
			VolumeCount: 5,
		})
	})

	shelf, err := c.GetUserBookshelf(context.Background(), "user123", 0)
	if err != nil {
		t.Fatal(err)
	}
	if shelf.Title != "Favorites" {
		t.Errorf("Title = %q, want Favorites", shelf.Title)
	}
	if shelf.VolumeCount != 5 {
		t.Errorf("VolumeCount = %d, want 5", shelf.VolumeCount)
	}
}

func TestGetUserBookshelfVolumes(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/user123/bookshelves/0/volumes" {
			t.Errorf("path = %q, want /users/user123/bookshelves/0/volumes", r.URL.Path)
		}
		if got := r.URL.Query().Get("maxResults"); got != "10" {
			t.Errorf("maxResults = %q, want 10", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.VolumesResponse{
			TotalItems: 1,
			Items: []googlebooks.Volume{
				{ID: "vol1", VolumeInfo: &googlebooks.VolumeInfo{Title: "Test Book"}},
			},
		})
	})

	resp, err := c.GetUserBookshelfVolumes(context.Background(), "user123", 0, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if resp.TotalItems != 1 {
		t.Fatalf("TotalItems = %d, want 1", resp.TotalItems)
	}
	if resp.Items[0].VolumeInfo.Title != "Test Book" {
		t.Errorf("Title = %q, want Test Book", resp.Items[0].VolumeInfo.Title)
	}
}

func TestGetVolumeAssociatedList(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/volumes/vol1/associated" {
			t.Errorf("path = %q, want /volumes/vol1/associated", r.URL.Path)
		}
		if got := r.URL.Query().Get("association"); got != "end-of-sample" {
			t.Errorf("association = %q, want end-of-sample", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.VolumesResponse{
			TotalItems: 1,
			Items: []googlebooks.Volume{
				{ID: "vol2", VolumeInfo: &googlebooks.VolumeInfo{Title: "Related Book"}},
			},
		})
	})

	resp, err := c.GetVolumeAssociatedList(context.Background(), "vol1", "end-of-sample")
	if err != nil {
		t.Fatal(err)
	}
	if resp.TotalItems != 1 {
		t.Fatalf("TotalItems = %d, want 1", resp.TotalItems)
	}
	if resp.Items[0].VolumeInfo.Title != "Related Book" {
		t.Errorf("Title = %q, want Related Book", resp.Items[0].VolumeInfo.Title)
	}
}

func TestGetMyLibraryBookshelves(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mylibrary/bookshelves" {
			t.Errorf("path = %q, want /mylibrary/bookshelves", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.BookshelvesResponse{
			Kind: "books#bookshelves",
			Items: []googlebooks.Bookshelf{
				{ID: 0, Title: "Favorites", VolumeCount: 3},
			},
		})
	})

	resp, err := c.GetMyLibraryBookshelves(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1", len(resp.Items))
	}
	if resp.Items[0].Title != "Favorites" {
		t.Errorf("Title = %q, want Favorites", resp.Items[0].Title)
	}
}

func TestGetMyLibraryBookshelf(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mylibrary/bookshelves/2" {
			t.Errorf("path = %q, want /mylibrary/bookshelves/2", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.Bookshelf{
			ID:    2,
			Title: "To read",
		})
	})

	shelf, err := c.GetMyLibraryBookshelf(context.Background(), 2)
	if err != nil {
		t.Fatal(err)
	}
	if shelf.ID != 2 {
		t.Errorf("ID = %d, want 2", shelf.ID)
	}
	if shelf.Title != "To read" {
		t.Errorf("Title = %q, want To read", shelf.Title)
	}
}

func TestGetMyLibraryBookshelfVolumes(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mylibrary/bookshelves/2/volumes" {
			t.Errorf("path = %q, want /mylibrary/bookshelves/2/volumes", r.URL.Path)
		}
		if got := r.URL.Query().Get("startIndex"); got != "5" {
			t.Errorf("startIndex = %q, want 5", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.VolumesResponse{
			TotalItems: 10,
			Items: []googlebooks.Volume{
				{ID: "vol1", VolumeInfo: &googlebooks.VolumeInfo{Title: "My Book"}},
			},
		})
	})

	resp, err := c.GetMyLibraryBookshelfVolumes(context.Background(), 2, 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	if resp.TotalItems != 10 {
		t.Fatalf("TotalItems = %d, want 10", resp.TotalItems)
	}
}

func TestAddVolumeToBookshelf(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/mylibrary/bookshelves/0/addVolume" {
			t.Errorf("path = %q, want /mylibrary/bookshelves/0/addVolume", r.URL.Path)
		}
		if got := r.URL.Query().Get("volumeId"); got != "vol123" {
			t.Errorf("volumeId = %q, want vol123", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.AddVolumeToBookshelf(context.Background(), 0, "vol123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveVolumeFromBookshelf(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/mylibrary/bookshelves/0/removeVolume" {
			t.Errorf("path = %q, want /mylibrary/bookshelves/0/removeVolume", r.URL.Path)
		}
		if got := r.URL.Query().Get("volumeId"); got != "vol123" {
			t.Errorf("volumeId = %q, want vol123", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.RemoveVolumeFromBookshelf(context.Background(), 0, "vol123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestClearBookshelf(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/mylibrary/bookshelves/0/clearVolumes" {
			t.Errorf("path = %q, want /mylibrary/bookshelves/0/clearVolumes", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.ClearBookshelf(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMyLibraryAnnotations(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mylibrary/annotations" {
			t.Errorf("path = %q, want /mylibrary/annotations", r.URL.Path)
		}
		if got := r.URL.Query().Get("volumeId"); got != "vol1" {
			t.Errorf("volumeId = %q, want vol1", got)
		}
		if got := r.URL.Query().Get("maxResults"); got != "20" {
			t.Errorf("maxResults = %q, want 20", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.AnnotationsResponse{
			Kind:       "books#annotations",
			TotalItems: 1,
			Items: []googlebooks.Annotation{
				{
					ID:           "ann1",
					VolumeID:     "vol1",
					SelectedText: "highlighted text",
				},
			},
		})
	})

	resp, err := c.GetMyLibraryAnnotations(context.Background(), "vol1", 0, 20)
	if err != nil {
		t.Fatal(err)
	}
	if resp.TotalItems != 1 {
		t.Fatalf("TotalItems = %d, want 1", resp.TotalItems)
	}
	if resp.Items[0].SelectedText != "highlighted text" {
		t.Errorf("SelectedText = %q", resp.Items[0].SelectedText)
	}
}

func TestGetMyLibraryReadingPositions(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mylibrary/readingpositions/vol1" {
			t.Errorf("path = %q, want /mylibrary/readingpositions/vol1", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.ReadingPosition{
			Kind:     "books#readingPosition",
			VolumeID: "vol1",
			Position: "cfi(/6/4[chap01ref]!/4[body01])",
		})
	})

	pos, err := c.GetMyLibraryReadingPositions(context.Background(), "vol1")
	if err != nil {
		t.Fatal(err)
	}
	if pos.VolumeID != "vol1" {
		t.Errorf("VolumeID = %q, want vol1", pos.VolumeID)
	}
	if pos.Position != "cfi(/6/4[chap01ref]!/4[body01])" {
		t.Errorf("Position = %q", pos.Position)
	}
}

func TestGetSeries(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/series/get" {
			t.Errorf("path = %q, want /series/get", r.URL.Path)
		}
		if got := r.URL.Query().Get("series_id"); got != "series123" {
			t.Errorf("series_id = %q, want series123", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.SeriesResponse{
			Kind: "books#series",
			Series: []googlebooks.SeriesInfo{
				{
					SeriesID:   "series123",
					Title:      "Harry Potter",
					SeriesType: "BOOK_SERIES",
				},
			},
		})
	})

	resp, err := c.GetSeries(context.Background(), "series123")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Series) != 1 {
		t.Fatalf("len(Series) = %d, want 1", len(resp.Series))
	}
	if resp.Series[0].Title != "Harry Potter" {
		t.Errorf("Title = %q, want Harry Potter", resp.Series[0].Title)
	}
}

func TestGetSeriesMembers(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/series/membership/get" {
			t.Errorf("path = %q, want /series/membership/get", r.URL.Path)
		}
		if got := r.URL.Query().Get("series_id"); got != "series123" {
			t.Errorf("series_id = %q, want series123", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(googlebooks.VolumesResponse{
			TotalItems: 7,
			Items: []googlebooks.Volume{
				{ID: "hp1", VolumeInfo: &googlebooks.VolumeInfo{Title: "Philosopher's Stone"}},
				{ID: "hp2", VolumeInfo: &googlebooks.VolumeInfo{Title: "Chamber of Secrets"}},
			},
		})
	})

	resp, err := c.GetSeriesMembers(context.Background(), "series123")
	if err != nil {
		t.Fatal(err)
	}
	if resp.TotalItems != 7 {
		t.Fatalf("TotalItems = %d, want 7", resp.TotalItems)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2", len(resp.Items))
	}
}

func TestPostError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"code":401,"message":"unauthorized"}}`))
	}))
	defer srv.Close()

	c := googlebooks.New("bad-key", metadata.WithBaseURL(srv.URL))
	err := c.AddVolumeToBookshelf(context.Background(), 0, "vol1")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *googlebooks.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
}
