package readarr_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golusoris/goenvoy/arr/readarr"
	"github.com/golusoris/goenvoy/arr/v2"
)

func newTestServer(t *testing.T, method, wantPath string, body any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("method = %s, want %s", r.Method, method)
		}
		if r.Header.Get("X-Api-Key") == "" {
			t.Error("missing X-Api-Key header")
		}
		if wantPath != "" && r.URL.RequestURI() != wantPath {
			t.Errorf("path = %q, want %q", r.URL.RequestURI(), wantPath)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}))
}

func newRawTestServer(t *testing.T, method, wantPath, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("method = %s, want %s", r.Method, method)
		}
		if r.Header.Get("X-Api-Key") == "" {
			t.Error("missing X-Api-Key header")
		}
		if wantPath != "" && r.URL.RequestURI() != wantPath {
			t.Errorf("path = %q, want %q", r.URL.RequestURI(), wantPath)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, body)
	}))
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		c, err := readarr.New("http://localhost:8787", "test-key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c == nil {
			t.Fatal("expected non-nil client")
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		t.Parallel()
		_, err := readarr.New("://bad", "test-key")
		if err == nil {
			t.Fatal("expected error for invalid URL")
		}
	})
}

func TestGetAllAuthors(t *testing.T) {
	t.Parallel()

	want := []readarr.Author{
		{ID: 1, AuthorName: "Brandon Sanderson", ForeignAuthorID: "author-1"},
		{ID: 2, AuthorName: "Stephen King", ForeignAuthorID: "author-2"},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/author", want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetAllAuthors(context.Background())
	if err != nil {
		t.Fatalf("GetAllAuthors: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].AuthorName != "Brandon Sanderson" {
		t.Errorf("AuthorName = %q, want %q", got[0].AuthorName, "Brandon Sanderson")
	}
}

func TestGetAuthor(t *testing.T) {
	t.Parallel()

	want := readarr.Author{ID: 1, AuthorName: "Brandon Sanderson"}

	srv := newTestServer(t, http.MethodGet, "/api/v1/author/1", want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetAuthor(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetAuthor: %v", err)
	}
	if got.AuthorName != "Brandon Sanderson" {
		t.Errorf("AuthorName = %q, want %q", got.AuthorName, "Brandon Sanderson")
	}
}

func TestAddAuthor(t *testing.T) {
	t.Parallel()

	want := readarr.Author{ID: 3, AuthorName: "New Author", ForeignAuthorID: "abc-123"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var body readarr.Author
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.AuthorName != "New Author" {
			t.Errorf("AuthorName = %q, want %q", body.AuthorName, "New Author")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.AddAuthor(context.Background(), &readarr.Author{
		AuthorName:      "New Author",
		ForeignAuthorID: "abc-123",
	})
	if err != nil {
		t.Fatalf("AddAuthor: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("ID = %d, want 3", got.ID)
	}
}

func TestDeleteAuthor(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete,
		"/api/v1/author/1?deleteFiles=true&addImportListExclusion=false",
		nil)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteAuthor(context.Background(), 1, true, false); err != nil {
		t.Fatalf("DeleteAuthor: %v", err)
	}
}

func TestLookupAuthor(t *testing.T) {
	t.Parallel()

	want := []readarr.Author{{ID: 0, AuthorName: "Stephen King"}}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/author/lookup?term=stephen+king",
		want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.LookupAuthor(context.Background(), "stephen king")
	if err != nil {
		t.Fatalf("LookupAuthor: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetBooks(t *testing.T) {
	t.Parallel()

	want := []readarr.Book{
		{ID: 10, Title: "The Way of Kings", AuthorID: 1, ForeignBookID: "book-1"},
		{ID: 11, Title: "Words of Radiance", AuthorID: 1, ForeignBookID: "book-2"},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/book?authorId=1",
		want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetBooks(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetBooks: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Title != "The Way of Kings" {
		t.Errorf("Title = %q, want %q", got[0].Title, "The Way of Kings")
	}
}

func TestGetBook(t *testing.T) {
	t.Parallel()

	want := readarr.Book{ID: 10, Title: "The Way of Kings"}

	srv := newTestServer(t, http.MethodGet, "/api/v1/book/10", want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetBook(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetBook: %v", err)
	}
	if got.Title != "The Way of Kings" {
		t.Errorf("Title = %q, want %q", got.Title, "The Way of Kings")
	}
}

func TestDeleteBook(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete,
		"/api/v1/book/10?deleteFiles=false&addImportListExclusion=true",
		nil)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteBook(context.Background(), 10, false, true); err != nil {
		t.Fatalf("DeleteBook: %v", err)
	}
}

func TestLookupBook(t *testing.T) {
	t.Parallel()

	want := []readarr.Book{{ID: 0, Title: "The Way of Kings"}}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/book/lookup?term=way+of+kings",
		want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.LookupBook(context.Background(), "way of kings")
	if err != nil {
		t.Fatalf("LookupBook: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}

func TestGetBookFiles(t *testing.T) {
	t.Parallel()

	want := []readarr.BookFile{
		{ID: 200, AuthorID: 1, BookID: 10, Path: "/books/Sanderson/The Way of Kings.epub", Size: 5000000},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/bookfile?authorId=1",
		want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetBookFiles(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetBookFiles: %v", err)
	}
	if got[0].Size != 5000000 {
		t.Errorf("Size = %d, want 5000000", got[0].Size)
	}
}

func TestDeleteBookFile(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v1/bookfile/200", nil)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteBookFile(context.Background(), 200); err != nil {
		t.Fatalf("DeleteBookFile: %v", err)
	}
}

func TestGetEditions(t *testing.T) {
	t.Parallel()

	want := []readarr.Edition{
		{ID: 50, BookID: 10, Title: "The Way of Kings (Hardcover)", ForeignEditionID: "ed-1"},
		{ID: 51, BookID: 10, Title: "The Way of Kings (Kindle)", ForeignEditionID: "ed-2", IsEbook: true},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/edition?bookId=10",
		want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetEditions(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetEditions: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if !got[1].IsEbook {
		t.Error("expected second edition to be ebook")
	}
}

func TestSendCommand(t *testing.T) {
	t.Parallel()

	want := arr.CommandResponse{ID: 42, Name: "RefreshAuthor"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var cmd arr.CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if cmd.Name != "RefreshAuthor" {
			t.Errorf("Name = %q, want RefreshAuthor", cmd.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.SendCommand(context.Background(), arr.CommandRequest{Name: "RefreshAuthor"})
	if err != nil {
		t.Fatalf("SendCommand: %v", err)
	}
	if got.ID != 42 {
		t.Errorf("ID = %d, want 42", got.ID)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	want := readarr.ParseResult{
		Title: "Brandon Sanderson - The Way of Kings (2010) [EPUB]",
		ParsedBookInfo: &readarr.ParsedBookInfo{
			AuthorName: "Brandon Sanderson",
			BookTitle:  "The Way of Kings",
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/parse?title=Brandon+Sanderson+-+The+Way+of+Kings+%282010%29+%5BEPUB%5D",
		want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.Parse(context.Background(), "Brandon Sanderson - The Way of Kings (2010) [EPUB]")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got.ParsedBookInfo.AuthorName != "Brandon Sanderson" {
		t.Errorf("AuthorName = %q, want %q", got.ParsedBookInfo.AuthorName, "Brandon Sanderson")
	}
}

func TestGetSystemStatus(t *testing.T) {
	t.Parallel()

	want := arr.StatusResponse{AppName: "Readarr", Version: "0.3.0"}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/system/status",
		want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetSystemStatus(context.Background())
	if err != nil {
		t.Fatalf("GetSystemStatus: %v", err)
	}
	if got.AppName != "Readarr" {
		t.Errorf("AppName = %q, want %q", got.AppName, "Readarr")
	}
}

func TestGetQueue(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[arr.QueueRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []arr.QueueRecord{
			{ID: 1, Title: "Brandon Sanderson - The Way of Kings"},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/queue?page=1&pageSize=10",
		want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetQueue(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetQueue: %v", err)
	}
	if got.TotalRecords != 1 {
		t.Errorf("TotalRecords = %d, want 1", got.TotalRecords)
	}
}

func TestGetTags(t *testing.T) {
	t.Parallel()

	want := []arr.Tag{{ID: 1, Label: "fiction"}, {ID: 2, Label: "non-fiction"}}

	srv := newTestServer(t, http.MethodGet, "/api/v1/tag", want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetTags(context.Background())
	if err != nil {
		t.Fatalf("GetTags: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}

func TestCreateTag(t *testing.T) {
	t.Parallel()

	want := arr.Tag{ID: 3, Label: "new-tag"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var tag arr.Tag
		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if tag.Label != "new-tag" {
			t.Errorf("Label = %q, want %q", tag.Label, "new-tag")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.CreateTag(context.Background(), "new-tag")
	if err != nil {
		t.Fatalf("CreateTag: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("ID = %d, want 3", got.ID)
	}
}

func TestGetMetadataProfiles(t *testing.T) {
	t.Parallel()

	want := []readarr.MetadataProfile{
		{ID: 1, Name: "Standard"},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/metadataprofile", want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetMetadataProfiles(context.Background())
	if err != nil {
		t.Fatalf("GetMetadataProfiles: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Name != "Standard" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Standard")
	}
}

func TestGetHistory(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[readarr.HistoryRecord]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []readarr.HistoryRecord{
			{ID: 5, AuthorID: 1, BookID: 10, EventType: "grabbed"},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/history?page=1&pageSize=10",
		want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetHistory(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetHistory: %v", err)
	}
	if got.Records[0].EventType != "grabbed" {
		t.Errorf("EventType = %q, want %q", got.Records[0].EventType, "grabbed")
	}
}

func TestGetWantedMissing(t *testing.T) {
	t.Parallel()

	want := arr.PagingResource[readarr.Book]{
		Page:         1,
		PageSize:     10,
		TotalRecords: 1,
		Records: []readarr.Book{
			{ID: 10, Title: "The Way of Kings"},
		},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/wanted/missing?page=1&pageSize=10",
		want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetWantedMissing(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetWantedMissing: %v", err)
	}
	if got.Records[0].Title != "The Way of Kings" {
		t.Errorf("Title = %q, want %q", got.Records[0].Title, "The Way of Kings")
	}
}

func TestGetSeries(t *testing.T) {
	t.Parallel()

	want := []readarr.Series{
		{ID: 1, Title: "The Stormlight Archive"},
	}

	srv := newTestServer(t, http.MethodGet,
		"/api/v1/series?authorId=1",
		want)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := c.GetSeries(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetSeries: %v", err)
	}
	if got[0].Title != "The Stormlight Archive" {
		t.Errorf("Title = %q, want %q", got[0].Title, "The Stormlight Archive")
	}
}

func TestDeleteQueueItem(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete,
		"/api/v1/queue/5?removeFromClient=true&blocklist=false",
		nil)
	defer srv.Close()

	c, err := readarr.New(srv.URL, "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := c.DeleteQueueItem(context.Background(), 5, true, false); err != nil {
		t.Fatalf("DeleteQueueItem: %v", err)
	}
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	}))
	defer srv.Close()

	c, err := readarr.New(srv.URL, "bad-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = c.GetAllAuthors(context.Background())
	if err == nil {
		t.Fatal("expected error for 401 response")
	}

	var apiErr *arr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *arr.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
}

// ---------- Notification Tests ----------.

func TestGetNotifications(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/notification", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetNotifications(context.Background())
	if err != nil {
		t.Fatalf("GetNotifications: %v", err)
	}
	if len(out) != 1 || out[0].ID != 1 {
		t.Fatalf("unexpected result: %+v", out)
	}
}

func TestGetNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/notification/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetNotification(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetNotification: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateNotification(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("CreateNotification: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/notification/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateNotification(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateNotification: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/notification/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteNotification(context.Background(), 1); err != nil {
		t.Fatalf("DeleteNotification: %v", err)
	}
}

func TestGetNotificationSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/notification/schema", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetNotificationSchema(context.Background())
	if err != nil {
		t.Fatalf("GetNotificationSchema: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestTestNotification(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification/test", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.TestNotification(context.Background(), &arr.ProviderResource{}); err != nil {
		t.Fatalf("TestNotification: %v", err)
	}
}

func TestTestAllNotifications(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification/testall", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.TestAllNotifications(context.Background()); err != nil {
		t.Fatalf("TestAllNotifications: %v", err)
	}
}

func TestNotificationAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/notification/action/testAction", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.NotificationAction(context.Background(), "testAction", &arr.ProviderResource{}); err != nil {
		t.Fatalf("NotificationAction: %v", err)
	}
}

// ---------- Download Client Tests ----------.

func TestGetDownloadClients(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/downloadclient", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetDownloadClients(context.Background())
	if err != nil {
		t.Fatalf("GetDownloadClients: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/downloadclient/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetDownloadClient(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetDownloadClient: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateDownloadClient(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("CreateDownloadClient: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/downloadclient/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateDownloadClient(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateDownloadClient: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/downloadclient/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteDownloadClient(context.Background(), 1); err != nil {
		t.Fatalf("DeleteDownloadClient: %v", err)
	}
}

func TestGetDownloadClientSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/downloadclient/schema", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetDownloadClientSchema(context.Background())
	if err != nil {
		t.Fatalf("GetDownloadClientSchema: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestTestDownloadClient(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient/test", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.TestDownloadClient(context.Background(), &arr.ProviderResource{}); err != nil {
		t.Fatalf("TestDownloadClient: %v", err)
	}
}

func TestTestAllDownloadClients(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient/testall", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.TestAllDownloadClients(context.Background()); err != nil {
		t.Fatalf("TestAllDownloadClients: %v", err)
	}
}

func TestBulkUpdateDownloadClients(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/downloadclient/bulk", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.BulkUpdateDownloadClients(context.Background(), &arr.ProviderBulkResource{})
	if err != nil {
		t.Fatalf("BulkUpdateDownloadClients: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestBulkDeleteDownloadClients(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/downloadclient/bulk", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.BulkDeleteDownloadClients(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("BulkDeleteDownloadClients: %v", err)
	}
}

func TestDownloadClientAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/downloadclient/action/testAction", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DownloadClientAction(context.Background(), "testAction", &arr.ProviderResource{}); err != nil {
		t.Fatalf("DownloadClientAction: %v", err)
	}
}

// ---------- Indexer Tests ----------.

func TestGetIndexers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetIndexers(context.Background())
	if err != nil {
		t.Fatalf("GetIndexers: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetIndexer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetIndexer(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetIndexer: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateIndexer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexer", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateIndexer(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("CreateIndexer: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateIndexer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/indexer/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateIndexer(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateIndexer: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteIndexer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/indexer/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteIndexer(context.Background(), 1); err != nil {
		t.Fatalf("DeleteIndexer: %v", err)
	}
}

func TestGetIndexerSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexer/schema", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetIndexerSchema(context.Background())
	if err != nil {
		t.Fatalf("GetIndexerSchema: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestTestIndexer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexer/test", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.TestIndexer(context.Background(), &arr.ProviderResource{}); err != nil {
		t.Fatalf("TestIndexer: %v", err)
	}
}

func TestTestAllIndexers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexer/testall", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.TestAllIndexers(context.Background()); err != nil {
		t.Fatalf("TestAllIndexers: %v", err)
	}
}

func TestBulkUpdateIndexers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/indexer/bulk", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.BulkUpdateIndexers(context.Background(), &arr.ProviderBulkResource{})
	if err != nil {
		t.Fatalf("BulkUpdateIndexers: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestBulkDeleteIndexers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/indexer/bulk", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.BulkDeleteIndexers(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("BulkDeleteIndexers: %v", err)
	}
}

func TestIndexerAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/indexer/action/testAction", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.IndexerAction(context.Background(), "testAction", &arr.ProviderResource{}); err != nil {
		t.Fatalf("IndexerAction: %v", err)
	}
}

// ---------- Import List Tests ----------.

func TestGetImportLists(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/importlist", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetImportLists(context.Background())
	if err != nil {
		t.Fatalf("GetImportLists: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetImportList(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/importlist/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetImportList(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetImportList: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateImportList(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/importlist", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateImportList(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("CreateImportList: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateImportList(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/importlist/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateImportList(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateImportList: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteImportList(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/importlist/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteImportList(context.Background(), 1); err != nil {
		t.Fatalf("DeleteImportList: %v", err)
	}
}

func TestGetImportListSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/importlist/schema", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetImportListSchema(context.Background())
	if err != nil {
		t.Fatalf("GetImportListSchema: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestTestImportList(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/importlist/test", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.TestImportList(context.Background(), &arr.ProviderResource{}); err != nil {
		t.Fatalf("TestImportList: %v", err)
	}
}

func TestTestAllImportLists(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/importlist/testall", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.TestAllImportLists(context.Background()); err != nil {
		t.Fatalf("TestAllImportLists: %v", err)
	}
}

func TestBulkUpdateImportLists(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/importlist/bulk", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.BulkUpdateImportLists(context.Background(), &arr.ProviderBulkResource{})
	if err != nil {
		t.Fatalf("BulkUpdateImportLists: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestBulkDeleteImportLists(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/importlist/bulk", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.BulkDeleteImportLists(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("BulkDeleteImportLists: %v", err)
	}
}

func TestImportListAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/importlist/action/testAction", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.ImportListAction(context.Background(), "testAction", &arr.ProviderResource{}); err != nil {
		t.Fatalf("ImportListAction: %v", err)
	}
}

// ---------- Metadata Consumer Tests ----------.

func TestGetMetadataConsumers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/metadata", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetMetadataConsumers(context.Background())
	if err != nil {
		t.Fatalf("GetMetadataConsumers: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetMetadataConsumer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/metadata/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetMetadataConsumer(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMetadataConsumer: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateMetadataConsumer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/metadata", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateMetadataConsumer(context.Background(), &arr.ProviderResource{})
	if err != nil {
		t.Fatalf("CreateMetadataConsumer: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateMetadataConsumer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/metadata/1", arr.ProviderResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateMetadataConsumer(context.Background(), &arr.ProviderResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateMetadataConsumer: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteMetadataConsumer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/metadata/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteMetadataConsumer(context.Background(), 1); err != nil {
		t.Fatalf("DeleteMetadataConsumer: %v", err)
	}
}

func TestGetMetadataConsumerSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/metadata/schema", []arr.ProviderResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetMetadataConsumerSchema(context.Background())
	if err != nil {
		t.Fatalf("GetMetadataConsumerSchema: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestTestMetadataConsumer(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/metadata/test", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.TestMetadataConsumer(context.Background(), &arr.ProviderResource{}); err != nil {
		t.Fatalf("TestMetadataConsumer: %v", err)
	}
}

func TestTestAllMetadataConsumers(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/metadata/testall", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.TestAllMetadataConsumers(context.Background()); err != nil {
		t.Fatalf("TestAllMetadataConsumers: %v", err)
	}
}

func TestMetadataConsumerAction(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/metadata/action/testAction", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.MetadataConsumerAction(context.Background(), "testAction", &arr.ProviderResource{}); err != nil {
		t.Fatalf("MetadataConsumerAction: %v", err)
	}
}

// ---------- Config Tests ----------.

func TestGetDownloadClientConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/downloadclient", arr.DownloadClientConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetDownloadClientConfig(context.Background())
	if err != nil {
		t.Fatalf("GetDownloadClientConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateDownloadClientConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/downloadclient/1", arr.DownloadClientConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateDownloadClientConfig(context.Background(), &arr.DownloadClientConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateDownloadClientConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetDownloadClientConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/downloadclient/1", arr.DownloadClientConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetDownloadClientConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetDownloadClientConfigByID: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetIndexerConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/indexer", arr.IndexerConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetIndexerConfig(context.Background())
	if err != nil {
		t.Fatalf("GetIndexerConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateIndexerConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/indexer/1", arr.IndexerConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateIndexerConfig(context.Background(), &arr.IndexerConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateIndexerConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetIndexerConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/indexer/1", arr.IndexerConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetIndexerConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetIndexerConfigByID: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetNamingConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/naming", arr.NamingConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetNamingConfig(context.Background())
	if err != nil {
		t.Fatalf("GetNamingConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateNamingConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/naming/1", arr.NamingConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateNamingConfig(context.Background(), &arr.NamingConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateNamingConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetNamingConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/naming/1", arr.NamingConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetNamingConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetNamingConfigByID: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetNamingExamples(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/naming/examples", arr.NamingConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetNamingExamples(context.Background())
	if err != nil {
		t.Fatalf("GetNamingExamples: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetHostConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/host", arr.HostConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetHostConfig(context.Background())
	if err != nil {
		t.Fatalf("GetHostConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateHostConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/host/1", arr.HostConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateHostConfig(context.Background(), &arr.HostConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateHostConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetHostConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/host/1", arr.HostConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetHostConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetHostConfigByID: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetUIConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/ui", arr.UIConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetUIConfig(context.Background())
	if err != nil {
		t.Fatalf("GetUIConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateUIConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/ui/1", arr.UIConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateUIConfig(context.Background(), &arr.UIConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateUIConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetUIConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/ui/1", arr.UIConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetUIConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUIConfigByID: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetMediaManagementConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/mediamanagement", arr.MediaManagementConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetMediaManagementConfig(context.Background())
	if err != nil {
		t.Fatalf("GetMediaManagementConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateMediaManagementConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/mediamanagement/1", arr.MediaManagementConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateMediaManagementConfig(context.Background(), &arr.MediaManagementConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateMediaManagementConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetMediaManagementConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/mediamanagement/1", arr.MediaManagementConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetMediaManagementConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMediaManagementConfigByID: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

// ---------- Development Config Tests ----------.

func TestGetDevelopmentConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/development", readarr.DevelopmentConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetDevelopmentConfig(context.Background())
	if err != nil {
		t.Fatalf("GetDevelopmentConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetDevelopmentConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/development/1", readarr.DevelopmentConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetDevelopmentConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetDevelopmentConfigByID: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateDevelopmentConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/development/1", readarr.DevelopmentConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateDevelopmentConfig(context.Background(), &readarr.DevelopmentConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateDevelopmentConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

// ---------- Metadata Provider Config Tests ----------.

func TestGetMetadataProviderConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/metadataprovider", readarr.MetadataProviderConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetMetadataProviderConfig(context.Background())
	if err != nil {
		t.Fatalf("GetMetadataProviderConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetMetadataProviderConfigByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/config/metadataprovider/1", readarr.MetadataProviderConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetMetadataProviderConfigByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMetadataProviderConfigByID: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateMetadataProviderConfig(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/config/metadataprovider/1", readarr.MetadataProviderConfigResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateMetadataProviderConfig(context.Background(), &readarr.MetadataProviderConfigResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateMetadataProviderConfig: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

// ---------- Quality Profile Tests ----------.

func TestGetQualityProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/qualityprofile/1", arr.QualityProfile{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetQualityProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetQualityProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateQualityProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/qualityprofile", arr.QualityProfile{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateQualityProfile(context.Background(), &arr.QualityProfile{})
	if err != nil {
		t.Fatalf("CreateQualityProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateQualityProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/qualityprofile/1", arr.QualityProfile{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateQualityProfile(context.Background(), &arr.QualityProfile{ID: 1})
	if err != nil {
		t.Fatalf("UpdateQualityProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteQualityProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/qualityprofile/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteQualityProfile(context.Background(), 1); err != nil {
		t.Fatalf("DeleteQualityProfile: %v", err)
	}
}

func TestGetQualityProfileSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/qualityprofile/schema", arr.QualityProfile{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetQualityProfileSchema(context.Background())
	if err != nil {
		t.Fatalf("GetQualityProfileSchema: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

// ---------- Quality Definition Tests ----------.

func TestGetQualityDefinitions(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/qualitydefinition", []arr.QualityDefinitionResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetQualityDefinitions(context.Background())
	if err != nil {
		t.Fatalf("GetQualityDefinitions: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetQualityDefinition(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/qualitydefinition/1", arr.QualityDefinitionResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetQualityDefinition(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetQualityDefinition: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateQualityDefinition(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/qualitydefinition/1", arr.QualityDefinitionResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateQualityDefinition(context.Background(), &arr.QualityDefinitionResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateQualityDefinition: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestBulkUpdateQualityDefinitions(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/qualitydefinition/update", []arr.QualityDefinitionResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.BulkUpdateQualityDefinitions(context.Background(), []arr.QualityDefinitionResource{{ID: 1}})
	if err != nil {
		t.Fatalf("BulkUpdateQualityDefinitions: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

// ---------- Metadata Profile Tests ----------.

func TestGetMetadataProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/metadataprofile/1", readarr.MetadataProfile{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetMetadataProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMetadataProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateMetadataProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/metadataprofile", readarr.MetadataProfile{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateMetadataProfile(context.Background(), &readarr.MetadataProfile{})
	if err != nil {
		t.Fatalf("CreateMetadataProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateMetadataProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/metadataprofile/1", readarr.MetadataProfile{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateMetadataProfile(context.Background(), &readarr.MetadataProfile{ID: 1})
	if err != nil {
		t.Fatalf("UpdateMetadataProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteMetadataProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/metadataprofile/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteMetadataProfile(context.Background(), 1); err != nil {
		t.Fatalf("DeleteMetadataProfile: %v", err)
	}
}

func TestGetMetadataProfileSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/metadataprofile/schema", readarr.MetadataProfile{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetMetadataProfileSchema(context.Background())
	if err != nil {
		t.Fatalf("GetMetadataProfileSchema: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

// ---------- Tag Tests ----------.

func TestGetTag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/tag/1", arr.Tag{ID: 1, Label: "test"})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetTag(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTag: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateTag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/tag/1", arr.Tag{ID: 1, Label: "updated"})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateTag(context.Background(), &arr.Tag{ID: 1, Label: "updated"})
	if err != nil {
		t.Fatalf("UpdateTag: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteTag(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/tag/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteTag(context.Background(), 1); err != nil {
		t.Fatalf("DeleteTag: %v", err)
	}
}

func TestGetTagDetails(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/tag/detail", []arr.TagDetail{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetTagDetails(context.Background())
	if err != nil {
		t.Fatalf("GetTagDetails: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetTagDetail(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/tag/detail/1", arr.TagDetail{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetTagDetail(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTagDetail: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

// ---------- Root Folder Tests ----------.

func TestGetRootFolder(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/rootfolder/1", arr.RootFolder{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetRootFolder(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetRootFolder: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateRootFolder(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/rootfolder", arr.RootFolder{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateRootFolder(context.Background(), &arr.RootFolder{})
	if err != nil {
		t.Fatalf("CreateRootFolder: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateRootFolder(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/rootfolder/1", arr.RootFolder{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateRootFolder(context.Background(), &arr.RootFolder{ID: 1})
	if err != nil {
		t.Fatalf("UpdateRootFolder: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteRootFolder(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/rootfolder/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteRootFolder(context.Background(), 1); err != nil {
		t.Fatalf("DeleteRootFolder: %v", err)
	}
}

// ---------- Custom Filter Tests ----------.

func TestGetCustomFilters(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/customfilter", []arr.CustomFilterResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetCustomFilters(context.Background())
	if err != nil {
		t.Fatalf("GetCustomFilters: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetCustomFilter(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/customfilter/1", arr.CustomFilterResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetCustomFilter(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCustomFilter: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateCustomFilter(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/customfilter", arr.CustomFilterResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateCustomFilter(context.Background(), &arr.CustomFilterResource{})
	if err != nil {
		t.Fatalf("CreateCustomFilter: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateCustomFilter(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/customfilter/1", arr.CustomFilterResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateCustomFilter(context.Background(), &arr.CustomFilterResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateCustomFilter: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteCustomFilter(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/customfilter/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteCustomFilter(context.Background(), 1); err != nil {
		t.Fatalf("DeleteCustomFilter: %v", err)
	}
}

// ---------- Custom Format Tests ----------.

func TestGetCustomFormats(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/customformat", []arr.CustomFormatResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetCustomFormats(context.Background())
	if err != nil {
		t.Fatalf("GetCustomFormats: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetCustomFormat(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/customformat/1", arr.CustomFormatResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetCustomFormat(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCustomFormat: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateCustomFormat(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/customformat", arr.CustomFormatResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateCustomFormat(context.Background(), &arr.CustomFormatResource{})
	if err != nil {
		t.Fatalf("CreateCustomFormat: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateCustomFormat(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/customformat/1", arr.CustomFormatResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateCustomFormat(context.Background(), &arr.CustomFormatResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateCustomFormat: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteCustomFormat(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/customformat/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteCustomFormat(context.Background(), 1); err != nil {
		t.Fatalf("DeleteCustomFormat: %v", err)
	}
}

func TestGetCustomFormatSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/customformat/schema", []arr.CustomFormatResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetCustomFormatSchema(context.Background())
	if err != nil {
		t.Fatalf("GetCustomFormatSchema: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

// ---------- Delay Profile Tests ----------.

func TestGetDelayProfiles(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/delayprofile", []arr.DelayProfileResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetDelayProfiles(context.Background())
	if err != nil {
		t.Fatalf("GetDelayProfiles: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetDelayProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/delayprofile/1", arr.DelayProfileResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetDelayProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetDelayProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateDelayProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/delayprofile", arr.DelayProfileResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateDelayProfile(context.Background(), &arr.DelayProfileResource{})
	if err != nil {
		t.Fatalf("CreateDelayProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateDelayProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/delayprofile/1", arr.DelayProfileResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateDelayProfile(context.Background(), &arr.DelayProfileResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateDelayProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteDelayProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/delayprofile/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteDelayProfile(context.Background(), 1); err != nil {
		t.Fatalf("DeleteDelayProfile: %v", err)
	}
}

func TestReorderDelayProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/delayprofile/reorder/1?after=2", []arr.DelayProfileResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.ReorderDelayProfile(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("ReorderDelayProfile: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

// ---------- Release Profile Tests ----------.

func TestGetReleaseProfiles(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/releaseprofile", []arr.ReleaseProfileResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetReleaseProfiles(context.Background())
	if err != nil {
		t.Fatalf("GetReleaseProfiles: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetReleaseProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/releaseprofile/1", arr.ReleaseProfileResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetReleaseProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetReleaseProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateReleaseProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/releaseprofile", arr.ReleaseProfileResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateReleaseProfile(context.Background(), &arr.ReleaseProfileResource{})
	if err != nil {
		t.Fatalf("CreateReleaseProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateReleaseProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/releaseprofile/1", arr.ReleaseProfileResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateReleaseProfile(context.Background(), &arr.ReleaseProfileResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateReleaseProfile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteReleaseProfile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/releaseprofile/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteReleaseProfile(context.Background(), 1); err != nil {
		t.Fatalf("DeleteReleaseProfile: %v", err)
	}
}

// ---------- Remote Path Mapping Tests ----------.

func TestGetRemotePathMappings(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/remotepathmapping", []arr.RemotePathMappingResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetRemotePathMappings(context.Background())
	if err != nil {
		t.Fatalf("GetRemotePathMappings: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetRemotePathMapping(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/remotepathmapping/1", arr.RemotePathMappingResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetRemotePathMapping(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetRemotePathMapping: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateRemotePathMapping(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/remotepathmapping", arr.RemotePathMappingResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateRemotePathMapping(context.Background(), &arr.RemotePathMappingResource{})
	if err != nil {
		t.Fatalf("CreateRemotePathMapping: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateRemotePathMapping(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/remotepathmapping/1", arr.RemotePathMappingResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateRemotePathMapping(context.Background(), &arr.RemotePathMappingResource{ID: 1})
	if err != nil {
		t.Fatalf("UpdateRemotePathMapping: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteRemotePathMapping(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/remotepathmapping/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteRemotePathMapping(context.Background(), 1); err != nil {
		t.Fatalf("DeleteRemotePathMapping: %v", err)
	}
}

// ---------- Import List Exclusion Tests ----------.

func TestGetImportListExclusion(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/importlistexclusion/1", readarr.ImportListExclusion{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetImportListExclusion(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetImportListExclusion: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestCreateImportListExclusion(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/importlistexclusion", readarr.ImportListExclusion{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.CreateImportListExclusion(context.Background(), &readarr.ImportListExclusion{})
	if err != nil {
		t.Fatalf("CreateImportListExclusion: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestUpdateImportListExclusion(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/importlistexclusion/1", readarr.ImportListExclusion{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateImportListExclusion(context.Background(), &readarr.ImportListExclusion{ID: 1})
	if err != nil {
		t.Fatalf("UpdateImportListExclusion: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestDeleteImportListExclusion(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/importlistexclusion/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteImportListExclusion(context.Background(), 1); err != nil {
		t.Fatalf("DeleteImportListExclusion: %v", err)
	}
}

// ---------- Blocklist Tests ----------.

func TestGetBlocklist(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/blocklist?page=1&pageSize=10", arr.PagingResource[arr.BlocklistResource]{Page: 1, PageSize: 10})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetBlocklist(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetBlocklist: %v", err)
	}
	if out.Page != 1 {
		t.Fatalf("unexpected page: %d", out.Page)
	}
}

func TestDeleteBlocklistItem(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/blocklist/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteBlocklistItem(context.Background(), 1); err != nil {
		t.Fatalf("DeleteBlocklistItem: %v", err)
	}
}

func TestBulkDeleteBlocklist(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/blocklist/bulk", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.BulkDeleteBlocklist(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("BulkDeleteBlocklist: %v", err)
	}
}

// ---------- Queue Extended Tests ----------.

func TestBulkDeleteQueue(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/queue/bulk?removeFromClient=true&blocklist=false", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.BulkDeleteQueue(context.Background(), &arr.QueueBulkResource{IDs: []int{1}}, true, false); err != nil {
		t.Fatalf("BulkDeleteQueue: %v", err)
	}
}

func TestGrabQueueItem(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/queue/grab/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.GrabQueueItem(context.Background(), 1); err != nil {
		t.Fatalf("GrabQueueItem: %v", err)
	}
}

func TestGrabQueueItemsBulk(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/queue/grab/bulk", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.GrabQueueItemsBulk(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("GrabQueueItemsBulk: %v", err)
	}
}

func TestGetQueueDetails(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/queue/details", []arr.QueueRecord{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetQueueDetails(context.Background())
	if err != nil {
		t.Fatalf("GetQueueDetails: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetQueueStatus(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/queue/status", arr.QueueStatusResource{TotalCount: 5})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetQueueStatus(context.Background())
	if err != nil {
		t.Fatalf("GetQueueStatus: %v", err)
	}
	if out.TotalCount != 5 {
		t.Fatalf("unexpected count: %d", out.TotalCount)
	}
}

// ---------- History Extended Tests ----------.

func TestGetHistoryByAuthor(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/history/author?authorId=1", []readarr.HistoryRecord{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetHistoryByAuthor(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetHistoryByAuthor: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetHistorySince(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/history/since?date=2024-01-01", []readarr.HistoryRecord{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetHistorySince(context.Background(), "2024-01-01")
	if err != nil {
		t.Fatalf("GetHistorySince: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestMarkHistoryFailed(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/history/failed/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.MarkHistoryFailed(context.Background(), 1); err != nil {
		t.Fatalf("MarkHistoryFailed: %v", err)
	}
}

// ---------- Release Tests ----------.

func TestSearchReleases(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/release?bookId=1", []arr.ReleaseResource{{GUID: "abc"}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.SearchReleases(context.Background(), 1)
	if err != nil {
		t.Fatalf("SearchReleases: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGrabRelease(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/release", arr.ReleaseResource{GUID: "abc"})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GrabRelease(context.Background(), &arr.ReleaseResource{GUID: "abc"})
	if err != nil {
		t.Fatalf("GrabRelease: %v", err)
	}
	if out.GUID != "abc" {
		t.Fatalf("unexpected GUID: %s", out.GUID)
	}
}

func TestPushRelease(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/release/push", []arr.ReleaseResource{{GUID: "abc"}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.PushRelease(context.Background(), &arr.ReleasePushResource{})
	if err != nil {
		t.Fatalf("PushRelease: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

// ---------- Rename Tests ----------.

func TestGetRenamePreview(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/rename?authorId=1&bookId=2", []readarr.RenameBookResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetRenamePreview(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("GetRenamePreview: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

// ---------- Retag Tests ----------.

func TestGetRetagPreview(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/retag?authorId=1&bookId=2", []readarr.RetagBookResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetRetagPreview(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("GetRetagPreview: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

// ---------- Manual Import Tests ----------.

func TestGetManualImport(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/manualimport?folder=%2Fdata", []arr.ManualImportResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetManualImport(context.Background(), "/data")
	if err != nil {
		t.Fatalf("GetManualImport: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestReprocessManualImport(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/manualimport", []arr.ManualImportResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.ReprocessManualImport(context.Background(), []arr.ManualImportReprocessResource{{ID: 1}})
	if err != nil {
		t.Fatalf("ReprocessManualImport: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

// ---------- Backup Tests ----------.

func TestGetBackups(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/backup", []arr.Backup{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetBackups(context.Background())
	if err != nil {
		t.Fatalf("GetBackups: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestDeleteBackup(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/system/backup/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteBackup(context.Background(), 1); err != nil {
		t.Fatalf("DeleteBackup: %v", err)
	}
}

func TestRestoreBackup(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/system/backup/restore/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.RestoreBackup(context.Background(), 1); err != nil {
		t.Fatalf("RestoreBackup: %v", err)
	}
}

// ---------- Log Tests ----------.

func TestGetLogs(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/log?page=1&pageSize=10", arr.PagingResource[arr.LogRecord]{Page: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetLogs(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetLogs: %v", err)
	}
	if out.Page != 1 {
		t.Fatalf("unexpected page: %d", out.Page)
	}
}

func TestGetLogFiles(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/log/file", []arr.LogFileResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetLogFiles(context.Background())
	if err != nil {
		t.Fatalf("GetLogFiles: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetLogFileContent(t *testing.T) {
	t.Parallel()
	srv := newRawTestServer(t, http.MethodGet, "/api/v1/log/file/readarr.txt", "log content")
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetLogFileContent(context.Background(), "readarr.txt")
	if err != nil {
		t.Fatalf("GetLogFileContent: %v", err)
	}
	if out != "log content" {
		t.Errorf("content = %q, want %q", out, "log content")
	}
}

func TestGetUpdateLogFiles(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/log/file/update", []arr.LogFileResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetUpdateLogFiles(context.Background())
	if err != nil {
		t.Fatalf("GetUpdateLogFiles: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetUpdateLogFileContent(t *testing.T) {
	t.Parallel()
	srv := newRawTestServer(t, http.MethodGet, "/api/v1/log/file/update/readarr.txt", "update log content")
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetUpdateLogFileContent(context.Background(), "readarr.txt")
	if err != nil {
		t.Fatalf("GetUpdateLogFileContent: %v", err)
	}
	if out != "update log content" {
		t.Errorf("content = %q, want %q", out, "update log content")
	}
}

// ---------- System Tests ----------.

func TestGetTasks(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/task", []arr.TaskResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetTasks(context.Background())
	if err != nil {
		t.Fatalf("GetTasks: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetTask(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/task/1", arr.TaskResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetTask(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetUpdates(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/update", []arr.UpdateResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetUpdates(context.Background())
	if err != nil {
		t.Fatalf("GetUpdates: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetSystemRoutes(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/routes", []arr.SystemRouteResource{{Path: "/"}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetSystemRoutes(context.Background())
	if err != nil {
		t.Fatalf("GetSystemRoutes: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetSystemRoutesDuplicate(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/system/routes/duplicate", []arr.SystemRouteResource{})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetSystemRoutesDuplicate(context.Background())
	if err != nil {
		t.Fatalf("GetSystemRoutesDuplicate: %v", err)
	}
	if out == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestShutdown(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/system/shutdown", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestRestart(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/system/restart", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.Restart(context.Background()); err != nil {
		t.Fatalf("Restart: %v", err)
	}
}

func TestDeleteCommand(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodDelete, "/api/v1/command/1", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.DeleteCommand(context.Background(), 1); err != nil {
		t.Fatalf("DeleteCommand: %v", err)
	}
}

// ---------- Language Tests ----------.

func TestGetLanguages(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/language", []arr.LanguageResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetLanguages(context.Background())
	if err != nil {
		t.Fatalf("GetLanguages: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestGetLanguage(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/language/1", arr.LanguageResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetLanguage(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetLanguage: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

// ---------- Localization Tests ----------.

func TestGetLocalization(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/localization", readarr.LocalizationResource{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetLocalization(context.Background())
	if err != nil {
		t.Fatalf("GetLocalization: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

// ---------- Ping Tests ----------.

func TestPing(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/ping", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

// ---------- Indexer Flag Tests ----------.

func TestGetIndexerFlags(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/indexerflag", []arr.IndexerFlagResource{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetIndexerFlags(context.Background())
	if err != nil {
		t.Fatalf("GetIndexerFlags: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

// ---------- File System Tests ----------.

func TestBrowseFileSystem(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/filesystem?path=%2Fdata", readarr.FileSystemResource{})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	_, err := c.BrowseFileSystem(context.Background(), "/data")
	if err != nil {
		t.Fatalf("BrowseFileSystem: %v", err)
	}
}

func TestGetFileSystemType(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/filesystem/type?path=%2Fdata", "local")
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetFileSystemType(context.Background(), "/data")
	if err != nil {
		t.Fatalf("GetFileSystemType: %v", err)
	}
	if out == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestGetFileSystemMediaFiles(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/filesystem/mediafiles?path=%2Fdata", []readarr.FileSystemEntry{{Name: "test.epub"}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetFileSystemMediaFiles(context.Background(), "/data")
	if err != nil {
		t.Fatalf("GetFileSystemMediaFiles: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

// ---------- Book File Extended Tests ----------.

func TestUpdateBookFile(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/bookfile/1", readarr.BookFile{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.UpdateBookFile(context.Background(), &readarr.BookFile{ID: 1})
	if err != nil {
		t.Fatalf("UpdateBookFile: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestEditBookFilesBulk(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPut, "/api/v1/bookfile/editor", []readarr.BookFile{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.EditBookFilesBulk(context.Background(), &readarr.BookFileListResource{BookFileIDs: []int{1}})
	if err != nil {
		t.Fatalf("EditBookFilesBulk: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

// ---------- Bookshelf Tests ----------.

func TestBookshelf(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodPost, "/api/v1/bookshelf", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.Bookshelf(context.Background(), &readarr.BookshelfResource{}); err != nil {
		t.Fatalf("Bookshelf: %v", err)
	}
}

// ---------- Calendar By ID Tests ----------.

func TestGetCalendarByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/calendar/1", readarr.Book{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetCalendarByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCalendarByID: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

// ---------- Wanted By ID Tests ----------.

func TestGetWantedMissingByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/wanted/missing/1", readarr.Book{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetWantedMissingByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetWantedMissingByID: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

func TestGetWantedCutoffByID(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/wanted/cutoff/1", readarr.Book{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetWantedCutoffByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetWantedCutoffByID: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

// ---------- Book Overview Tests ----------.

func TestGetBookOverview(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/book/1/overview", readarr.Book{ID: 1})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.GetBookOverview(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetBookOverview: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("unexpected ID: %d", out.ID)
	}
}

// ---------- Search Tests ----------.

func TestSearch(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodGet, "/api/v1/search?term=tolkien", []readarr.Author{{ID: 1}})
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	out, err := c.Search(context.Background(), "tolkien")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("unexpected length: %d", len(out))
	}
}

func TestHeadPing(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t, http.MethodHead, "/ping", nil)
	defer srv.Close()
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.HeadPing(context.Background()); err != nil {
		t.Fatalf("HeadPing: %v", err)
	}
}

func TestUploadBackup(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.Header.Get("X-Api-Key") == "" {
			t.Error("missing X-Api-Key header")
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
	c, _ := readarr.New(srv.URL, "test-key")
	if err := c.UploadBackup(context.Background(), "backup.zip", strings.NewReader("fake-backup-data")); err != nil {
		t.Fatalf("UploadBackup: %v", err)
	}
}
