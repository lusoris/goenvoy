package openlibrary_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/book/openlibrary"
)

func setup(t *testing.T, handler http.HandlerFunc) *openlibrary.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return openlibrary.New(metadata.WithBaseURL(srv.URL))
}

func TestSearch(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q"); got != "the lord of the rings" {
			t.Errorf("q = %q, want the lord of the rings", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openlibrary.SearchResponse{
			NumFound: 1,
			Docs: []openlibrary.SearchDoc{
				{Key: "/works/OL27448W", Title: "The Lord of the Rings", EditionCount: 250},
			},
		})
	})

	resp, err := c.Search(context.Background(), "the lord of the rings")
	if err != nil {
		t.Fatal(err)
	}
	if resp.NumFound != 1 {
		t.Fatalf("NumFound = %d, want 1", resp.NumFound)
	}
	if resp.Docs[0].Title != "The Lord of the Rings" {
		t.Errorf("Title = %q", resp.Docs[0].Title)
	}
}

func TestSearchWithParams(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("author"); got != "tolkien" {
			t.Errorf("author = %q, want tolkien", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openlibrary.SearchResponse{
			NumFound: 2,
			Docs: []openlibrary.SearchDoc{
				{Key: "/works/OL27448W", Title: "The Lord of the Rings"},
				{Key: "/works/OL27479W", Title: "The Hobbit"},
			},
		})
	})

	params := url.Values{}
	params.Set("author", "tolkien")
	resp, err := c.SearchWithParams(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
	if resp.NumFound != 2 {
		t.Fatalf("NumFound = %d, want 2", resp.NumFound)
	}
}

func TestGetWork(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/works/OL27448W.json" {
			t.Errorf("path = %q, want /works/OL27448W.json", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openlibrary.Work{
			Key:      "/works/OL27448W",
			Title:    "The Lord of the Rings",
			Subjects: []string{"Fantasy", "Fiction"},
		})
	})

	work, err := c.GetWork(context.Background(), "OL27448W")
	if err != nil {
		t.Fatal(err)
	}
	if work.Title != "The Lord of the Rings" {
		t.Errorf("Title = %q", work.Title)
	}
	if len(work.Subjects) != 2 {
		t.Fatalf("len(Subjects) = %d, want 2", len(work.Subjects))
	}
}

func TestGetEdition(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/books/OL7353617M.json" {
			t.Errorf("path = %q, want /books/OL7353617M.json", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openlibrary.Edition{
			Key:           "/books/OL7353617M",
			Title:         "The Lord of the Rings",
			NumberOfPages: 1178,
			Publishers:    []string{"Houghton Mifflin"},
		})
	})

	edition, err := c.GetEdition(context.Background(), "OL7353617M")
	if err != nil {
		t.Fatal(err)
	}
	if edition.NumberOfPages != 1178 {
		t.Errorf("NumberOfPages = %d, want 1178", edition.NumberOfPages)
	}
	if len(edition.Publishers) != 1 || edition.Publishers[0] != "Houghton Mifflin" {
		t.Errorf("Publishers = %v", edition.Publishers)
	}
}

func TestGetAuthor(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/authors/OL34184A.json" {
			t.Errorf("path = %q, want /authors/OL34184A.json", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openlibrary.Author{
			Key:       "/authors/OL34184A",
			Name:      "J.R.R. Tolkien",
			BirthDate: "3 January 1892",
			DeathDate: "2 September 1973",
		})
	})

	author, err := c.GetAuthor(context.Background(), "OL34184A")
	if err != nil {
		t.Fatal(err)
	}
	if author.Name != "J.R.R. Tolkien" {
		t.Errorf("Name = %q", author.Name)
	}
	if author.BirthDate != "3 January 1892" {
		t.Errorf("BirthDate = %q", author.BirthDate)
	}
}

func TestGetAuthorWorks(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/authors/OL34184A/works.json" {
			t.Errorf("path = %q, want /authors/OL34184A/works.json", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"entries": []map[string]any{
				{"key": "/works/OL27448W", "title": "The Lord of the Rings"},
				{"key": "/works/OL27479W", "title": "The Hobbit"},
			},
		})
	})

	works, err := c.GetAuthorWorks(context.Background(), "OL34184A")
	if err != nil {
		t.Fatal(err)
	}
	if len(works) != 2 {
		t.Fatalf("len = %d, want 2", len(works))
	}
	if works[1].Title != "The Hobbit" {
		t.Errorf("Title = %q, want The Hobbit", works[1].Title)
	}
}

func TestGetSubject(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/subjects/fantasy.json" {
			t.Errorf("path = %q, want /subjects/fantasy.json", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openlibrary.Subject{
			Name:      "fantasy",
			WorkCount: 5000,
			Works: []openlibrary.SubjectWork{
				{Key: "/works/OL27448W", Title: "The Lord of the Rings", EditionCount: 250},
			},
		})
	})

	subject, err := c.GetSubject(context.Background(), "fantasy")
	if err != nil {
		t.Fatal(err)
	}
	if subject.Name != "fantasy" {
		t.Errorf("Name = %q", subject.Name)
	}
	if subject.WorkCount != 5000 {
		t.Errorf("WorkCount = %d, want 5000", subject.WorkCount)
	}
}

func TestGetByISBN(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/isbn/9780618640157.json" {
			t.Errorf("path = %q, want /isbn/9780618640157.json", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openlibrary.Edition{
			Key:    "/books/OL7353617M",
			Title:  "The Lord of the Rings",
			ISBN13: []string{"9780618640157"},
		})
	})

	edition, err := c.GetByISBN(context.Background(), "9780618640157")
	if err != nil {
		t.Fatal(err)
	}
	if edition.Title != "The Lord of the Rings" {
		t.Errorf("Title = %q", edition.Title)
	}
	if len(edition.ISBN13) != 1 || edition.ISBN13[0] != "9780618640157" {
		t.Errorf("ISBN13 = %v", edition.ISBN13)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	}))
	defer srv.Close()

	c := openlibrary.New(metadata.WithBaseURL(srv.URL))
	_, err := c.Search(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *openlibrary.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusNotFound)
	}
}

func TestAPIErrorOnWork(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("server error"))
	}))
	defer srv.Close()

	c := openlibrary.New(metadata.WithBaseURL(srv.URL))
	_, err := c.GetWork(context.Background(), "INVALID")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *openlibrary.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusInternalServerError)
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	c := openlibrary.New()
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}
