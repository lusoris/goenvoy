package komga

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return New(ts.URL, "admin@example.com", "password")
}

func TestGetLibraries(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != "admin@example.com" || p != "password" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode([]Library{{ID: "1", Name: "Comics"}})
	})

	libs, err := c.GetLibraries(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(libs) != 1 || libs[0].Name != "Comics" {
		t.Fatalf("unexpected libraries: %+v", libs)
	}
}

func TestGetLibrary(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Library{ID: "1", Name: "Manga"})
	})

	lib, err := c.GetLibrary(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if lib.Name != "Manga" {
		t.Fatalf("unexpected library: %+v", lib)
	}
}

func TestGetSeries(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Page[Series]{
			Content:       []Series{{ID: "s1", Name: "One Piece"}},
			TotalElements: 1,
			TotalPages:    1,
		})
	})

	p, err := c.GetSeries(context.Background(), "", 0, 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Content) != 1 || p.Content[0].Name != "One Piece" {
		t.Fatalf("unexpected series: %+v", p)
	}
}

func TestGetOneSeries(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Series{ID: "s1", Name: "Naruto"})
	})

	s, err := c.GetOneSeries(context.Background(), "s1")
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != "Naruto" {
		t.Fatalf("unexpected series: %+v", s)
	}
}

func TestGetBooks(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Page[Book]{
			Content:       []Book{{ID: "b1", Name: "Chapter 1"}},
			TotalElements: 1,
			TotalPages:    1,
		})
	})

	p, err := c.GetBooks(context.Background(), "s1", 0, 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Content) != 1 || p.Content[0].Name != "Chapter 1" {
		t.Fatalf("unexpected books: %+v", p)
	}
}

func TestGetBook(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Book{ID: "b1", Name: "Vol 1"})
	})

	b, err := c.GetBook(context.Background(), "b1")
	if err != nil {
		t.Fatal(err)
	}
	if b.Name != "Vol 1" {
		t.Fatalf("unexpected book: %+v", b)
	}
}

func TestGetCollections(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Page[Collection]{
			Content:       []Collection{{ID: "c1", Name: "Favorites"}},
			TotalElements: 1,
			TotalPages:    1,
		})
	})

	p, err := c.GetCollections(context.Background(), 0, 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Content) != 1 || p.Content[0].Name != "Favorites" {
		t.Fatalf("unexpected collections: %+v", p)
	}
}

func TestGetCollection(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Collection{ID: "c1", Name: "DC Comics"})
	})

	col, err := c.GetCollection(context.Background(), "c1")
	if err != nil {
		t.Fatal(err)
	}
	if col.Name != "DC Comics" {
		t.Fatalf("unexpected collection: %+v", col)
	}
}

func TestGetReadLists(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(Page[ReadList]{
			Content:       []ReadList{{ID: "rl1", Name: "To Read"}},
			TotalElements: 1,
			TotalPages:    1,
		})
	})

	p, err := c.GetReadLists(context.Background(), 0, 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Content) != 1 || p.Content[0].Name != "To Read" {
		t.Fatalf("unexpected read lists: %+v", p)
	}
}

func TestGetReadList(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(ReadList{ID: "rl1", Name: "Weekend Reads"})
	})

	rl, err := c.GetReadList(context.Background(), "rl1")
	if err != nil {
		t.Fatal(err)
	}
	if rl.Name != "Weekend Reads" {
		t.Fatalf("unexpected read list: %+v", rl)
	}
}

func TestGetUsers(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode([]User{{ID: "u1", Email: "admin@example.com"}})
	})

	users, err := c.GetUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 || users[0].Email != "admin@example.com" {
		t.Fatalf("unexpected users: %+v", users)
	}
}

func TestGetMe(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(User{ID: "u1", Email: "me@example.com"})
	})

	u, err := c.GetMe(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if u.Email != "me@example.com" {
		t.Fatalf("unexpected user: %+v", u)
	}
}

func TestAPIError(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	})

	_, err := c.GetLibraries(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Fatalf("unexpected status code: %d", apiErr.StatusCode)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c := New("http://localhost", "user", "pass", WithHTTPClient(custom))
	if c.http != custom {
		t.Fatal("custom HTTP client not set")
	}
}
