package letterboxd_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/metadata"
	"github.com/lusoris/goenvoy/metadata/video/letterboxd"
)

func newTestServer(t *testing.T, wantPath string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-token")
		}
		if wantPath != "" && r.URL.RequestURI() != wantPath {
			t.Errorf("path = %q, want %q", r.URL.RequestURI(), wantPath)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
}

func newMutationTestServer(t *testing.T, wantPath, wantMethod string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != wantMethod {
			t.Errorf("method = %s, want %s", r.Method, wantMethod)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-token")
		}
		if wantPath != "" && r.URL.RequestURI() != wantPath {
			t.Errorf("path = %q, want %q", r.URL.RequestURI(), wantPath)
		}
		w.Header().Set("Content-Type", "application/json")
		if response != nil {
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}))
}

func newClient(t *testing.T, srv *httptest.Server) *letterboxd.Client {
	t.Helper()
	return letterboxd.New("test-token", metadata.WithBaseURL(srv.URL))
}

func TestNew(t *testing.T) {
	t.Parallel()
	c := letterboxd.New("token")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Film not found"})
	}))
	defer srv.Close()

	c := newClient(t, srv)
	_, err := c.GetFilm(context.Background(), "abc123")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *letterboxd.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusNotFound)
	}
}

func TestAPIErrorRawBody(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oops"))
	}))
	defer srv.Close()

	c := newClient(t, srv)
	_, err := c.GetFilm(context.Background(), "abc123")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *letterboxd.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.RawBody != "oops" {
		t.Errorf("RawBody = %q, want %q", apiErr.RawBody, "oops")
	}
}

// Films.

func TestGetFilm(t *testing.T) {
	t.Parallel()
	want := letterboxd.Film{ID: "abc1", Name: "Inception", ReleaseYear: 2010}
	srv := newTestServer(t, "/film/abc1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilm(context.Background(), "abc1")
	if err != nil {
		t.Fatalf("GetFilm: %v", err)
	}
	if got.Name != "Inception" {
		t.Errorf("Name = %q, want %q", got.Name, "Inception")
	}
}

func TestGetFilmByTMDbID(t *testing.T) {
	t.Parallel()
	want := letterboxd.Film{ID: "abc1", Name: "Inception"}
	srv := newTestServer(t, "/film/tmdb:27205", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilm(context.Background(), "tmdb:27205")
	if err != nil {
		t.Fatalf("GetFilm: %v", err)
	}
	if got.Name != "Inception" {
		t.Errorf("Name = %q, want %q", got.Name, "Inception")
	}
}

func TestGetFilmStatistics(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmStatistics{Rating: 4.3, Counts: &letterboxd.FilmStatCounts{Watches: 100000}}
	srv := newTestServer(t, "/film/abc1/statistics", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilmStatistics(context.Background(), "abc1")
	if err != nil {
		t.Fatalf("GetFilmStatistics: %v", err)
	}
	if got.Rating != 4.3 {
		t.Errorf("Rating = %f, want %f", got.Rating, 4.3)
	}
}

func TestGetFilmRelationship(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmRelationship{Watched: true, Liked: true, Rating: 4.5}
	srv := newTestServer(t, "/film/abc1/me", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilmRelationship(context.Background(), "abc1")
	if err != nil {
		t.Fatalf("GetFilmRelationship: %v", err)
	}
	if !got.Watched {
		t.Error("Watched = false, want true")
	}
}

func TestUpdateFilmRelationship(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmRelationshipUpdateResponse{
		Data: &letterboxd.FilmRelationship{Watched: true, Liked: true},
	}
	srv := newMutationTestServer(t, "/film/abc1/me", http.MethodPatch, want)
	defer srv.Close()

	c := newClient(t, srv)
	watched := true
	got, err := c.UpdateFilmRelationship(context.Background(), "abc1", letterboxd.FilmRelationshipUpdateRequest{Watched: &watched})
	if err != nil {
		t.Fatalf("UpdateFilmRelationship: %v", err)
	}
	if !got.Data.Watched {
		t.Error("Watched = false, want true")
	}
}

func TestGetFilmMemberRelationship(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmRelationship{Watched: true}
	srv := newTestServer(t, "/film/abc1/member/user1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilmMemberRelationship(context.Background(), "abc1", "user1")
	if err != nil {
		t.Fatalf("GetFilmMemberRelationship: %v", err)
	}
	if !got.Watched {
		t.Error("Watched = false, want true")
	}
}

func TestGetFilmMembers(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmMembersResponse{
		Items: []letterboxd.FilmMemberRelationshipItem{
			{Member: &letterboxd.MemberSummary{ID: "u1"}},
		},
	}
	srv := newTestServer(t, "/film/abc1/members", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilmMembers(context.Background(), "abc1", "", 0)
	if err != nil {
		t.Fatalf("GetFilmMembers: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetFilmFriends(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmFriendsResponse{
		Items: []letterboxd.FilmMemberRelationshipItem{
			{Member: &letterboxd.MemberSummary{ID: "friend1"}},
		},
	}
	srv := newTestServer(t, "/film/abc1/friends", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilmFriends(context.Background(), "abc1")
	if err != nil {
		t.Fatalf("GetFilmFriends: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestReportFilm(t *testing.T) {
	t.Parallel()
	srv := newMutationTestServer(t, "/film/abc1/report", http.MethodPost, nil)
	defer srv.Close()

	c := newClient(t, srv)
	err := c.ReportFilm(context.Background(), "abc1", letterboxd.ReportRequest{Reason: "spam"})
	if err != nil {
		t.Fatalf("ReportFilm: %v", err)
	}
}

func TestGetFilms(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmsResponse{
		Cursor: "next1",
		Items:  []letterboxd.FilmSummary{{ID: "f1", Name: "Test Film"}},
	}
	srv := newTestServer(t, "/films", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilms(context.Background(), "", 0)
	if err != nil {
		t.Fatalf("GetFilms: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetGenres(t *testing.T) {
	t.Parallel()
	want := letterboxd.GenresResponse{Items: []letterboxd.Genre{{ID: "action", Name: "Action"}}}
	srv := newTestServer(t, "/films/genres", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetGenres(context.Background())
	if err != nil {
		t.Fatalf("GetGenres: %v", err)
	}
	if len(got.Items) != 1 {
		t.Errorf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetFilmServices(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmServicesResponse{Items: []letterboxd.FilmService{{ID: "netflix", Name: "Netflix"}}}
	srv := newTestServer(t, "/films/film-services", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilmServices(context.Background())
	if err != nil {
		t.Fatalf("GetFilmServices: %v", err)
	}
	if len(got.Items) != 1 {
		t.Errorf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetCountries(t *testing.T) {
	t.Parallel()
	want := letterboxd.CountriesResponse{Items: []letterboxd.Country{{Code: "US", Name: "USA"}}}
	srv := newTestServer(t, "/films/countries", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetCountries(context.Background())
	if err != nil {
		t.Fatalf("GetCountries: %v", err)
	}
	if len(got.Items) != 1 {
		t.Errorf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetLanguages(t *testing.T) {
	t.Parallel()
	want := letterboxd.LanguagesResponse{Items: []letterboxd.Language{{Code: "en", Name: "English"}}}
	srv := newTestServer(t, "/films/languages", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetLanguages(context.Background())
	if err != nil {
		t.Fatalf("GetLanguages: %v", err)
	}
	if len(got.Items) != 1 {
		t.Errorf("Items = %d, want 1", len(got.Items))
	}
}

// Film collections.

func TestGetFilmCollection(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmCollection{ID: "col1", Name: "MCU"}
	srv := newTestServer(t, "/film-collection/col1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilmCollection(context.Background(), "col1")
	if err != nil {
		t.Fatalf("GetFilmCollection: %v", err)
	}
	if got.Name != "MCU" {
		t.Errorf("Name = %q, want %q", got.Name, "MCU")
	}
}

func TestGetFilmCollections(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmCollectionsResponse{
		Items: []letterboxd.FilmCollectionSummary{{ID: "col1", Name: "MCU"}},
	}
	srv := newTestServer(t, "/film-collections", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilmCollections(context.Background(), "", 0)
	if err != nil {
		t.Fatalf("GetFilmCollections: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

// Contributors.

func TestGetContributor(t *testing.T) {
	t.Parallel()
	want := letterboxd.Contributor{ID: "ct1", Name: "Christopher Nolan"}
	srv := newTestServer(t, "/contributor/ct1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetContributor(context.Background(), "ct1")
	if err != nil {
		t.Fatalf("GetContributor: %v", err)
	}
	if got.Name != "Christopher Nolan" {
		t.Errorf("Name = %q, want %q", got.Name, "Christopher Nolan")
	}
}

func TestGetContributions(t *testing.T) {
	t.Parallel()
	want := letterboxd.ContributionsResponse{
		Items: []letterboxd.ContributionItem{{Type: "Director", Film: &letterboxd.FilmSummary{ID: "f1"}}},
	}
	srv := newTestServer(t, "/contributor/ct1/contributions", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetContributions(context.Background(), "ct1", "", 0)
	if err != nil {
		t.Fatalf("GetContributions: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

// Lists.

func TestGetList(t *testing.T) {
	t.Parallel()
	want := letterboxd.List{ID: "lst1", Name: "Top 250", FilmCount: 250}
	srv := newTestServer(t, "/list/lst1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetList(context.Background(), "lst1")
	if err != nil {
		t.Fatalf("GetList: %v", err)
	}
	if got.Name != "Top 250" {
		t.Errorf("Name = %q, want %q", got.Name, "Top 250")
	}
}

func TestGetListEntries(t *testing.T) {
	t.Parallel()
	want := letterboxd.ListEntriesResponse{
		Items: []letterboxd.ListEntry{{Rank: 1, Film: &letterboxd.FilmSummary{ID: "f1"}}},
	}
	srv := newTestServer(t, "/list/lst1/entries", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetListEntries(context.Background(), "lst1", "", 0)
	if err != nil {
		t.Fatalf("GetListEntries: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetListComments(t *testing.T) {
	t.Parallel()
	want := letterboxd.CommentsResponse{
		Items: []letterboxd.Comment{{ID: "c1", Comment: "Great list!"}},
	}
	srv := newTestServer(t, "/list/lst1/comments", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetListComments(context.Background(), "lst1", "", 0)
	if err != nil {
		t.Fatalf("GetListComments: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetListStatistics(t *testing.T) {
	t.Parallel()
	want := letterboxd.ListStatistics{Counts: &letterboxd.ListStatCounts{Comments: 5, Likes: 100}}
	srv := newTestServer(t, "/list/lst1/statistics", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetListStatistics(context.Background(), "lst1")
	if err != nil {
		t.Fatalf("GetListStatistics: %v", err)
	}
	if got.Counts.Likes != 100 {
		t.Errorf("Likes = %d, want 100", got.Counts.Likes)
	}
}

func TestGetListRelationship(t *testing.T) {
	t.Parallel()
	want := letterboxd.ListRelationship{Liked: true}
	srv := newTestServer(t, "/list/lst1/me", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetListRelationship(context.Background(), "lst1")
	if err != nil {
		t.Fatalf("GetListRelationship: %v", err)
	}
	if !got.Liked {
		t.Error("Liked = false, want true")
	}
}

func TestUpdateListRelationship(t *testing.T) {
	t.Parallel()
	want := letterboxd.ListRelationshipUpdateResponse{
		Data: &letterboxd.ListRelationship{Liked: true},
	}
	srv := newMutationTestServer(t, "/list/lst1/me", http.MethodPatch, want)
	defer srv.Close()

	c := newClient(t, srv)
	liked := true
	got, err := c.UpdateListRelationship(context.Background(), "lst1", letterboxd.ListRelationshipUpdateRequest{Liked: &liked})
	if err != nil {
		t.Fatalf("UpdateListRelationship: %v", err)
	}
	if !got.Data.Liked {
		t.Error("Liked = false, want true")
	}
}

func TestGetLists(t *testing.T) {
	t.Parallel()
	want := letterboxd.ListsResponse{
		Items: []letterboxd.ListSummary{{ID: "lst1", Name: "Favorites"}},
	}
	srv := newTestServer(t, "/lists", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetLists(context.Background(), "", 0)
	if err != nil {
		t.Fatalf("GetLists: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetListTopics(t *testing.T) {
	t.Parallel()
	want := letterboxd.ListTopicsResponse{
		Items: []letterboxd.ListTopic{{Name: "Featured"}},
	}
	srv := newTestServer(t, "/lists/topics", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetListTopics(context.Background())
	if err != nil {
		t.Fatalf("GetListTopics: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestCreateList(t *testing.T) {
	t.Parallel()
	want := letterboxd.ListCreateResponse{
		Data: &letterboxd.List{ID: "new1", Name: "My List"},
	}
	srv := newMutationTestServer(t, "/lists", http.MethodPost, want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.CreateList(context.Background(), &letterboxd.ListCreationRequest{Name: "My List"})
	if err != nil {
		t.Fatalf("CreateList: %v", err)
	}
	if got.Data.Name != "My List" {
		t.Errorf("Name = %q, want %q", got.Data.Name, "My List")
	}
}

func TestUpdateList(t *testing.T) {
	t.Parallel()
	want := letterboxd.ListUpdateResponse{
		Data: &letterboxd.List{ID: "lst1", Name: "Updated"},
	}
	srv := newMutationTestServer(t, "/list/lst1", http.MethodPatch, want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.UpdateList(context.Background(), "lst1", &letterboxd.ListUpdateRequest{Name: "Updated"})
	if err != nil {
		t.Fatalf("UpdateList: %v", err)
	}
	if got.Data.Name != "Updated" {
		t.Errorf("Name = %q, want %q", got.Data.Name, "Updated")
	}
}

func TestDeleteList(t *testing.T) {
	t.Parallel()
	srv := newMutationTestServer(t, "/list/lst1", http.MethodDelete, nil)
	defer srv.Close()

	c := newClient(t, srv)
	err := c.DeleteList(context.Background(), "lst1")
	if err != nil {
		t.Fatalf("DeleteList: %v", err)
	}
}

func TestCreateListComment(t *testing.T) {
	t.Parallel()
	want := letterboxd.Comment{ID: "c1", Comment: "Nice list!"}
	srv := newMutationTestServer(t, "/list/lst1/comments", http.MethodPost, want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.CreateListComment(context.Background(), "lst1", letterboxd.CommentCreationRequest{Comment: "Nice list!"})
	if err != nil {
		t.Fatalf("CreateListComment: %v", err)
	}
	if got.Comment != "Nice list!" {
		t.Errorf("Comment = %q, want %q", got.Comment, "Nice list!")
	}
}

func TestReportList(t *testing.T) {
	t.Parallel()
	srv := newMutationTestServer(t, "/list/lst1/report", http.MethodPost, nil)
	defer srv.Close()

	c := newClient(t, srv)
	err := c.ReportList(context.Background(), "lst1", letterboxd.ReportRequest{Reason: "spam"})
	if err != nil {
		t.Fatalf("ReportList: %v", err)
	}
}

func TestForgetList(t *testing.T) {
	t.Parallel()
	srv := newMutationTestServer(t, "/list/lst1/forget", http.MethodPost, nil)
	defer srv.Close()

	c := newClient(t, srv)
	err := c.ForgetList(context.Background(), "lst1")
	if err != nil {
		t.Fatalf("ForgetList: %v", err)
	}
}

// Log entries.

func TestGetLogEntry(t *testing.T) {
	t.Parallel()
	want := letterboxd.LogEntry{ID: "le1", Rating: 4.5}
	srv := newTestServer(t, "/log-entry/le1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetLogEntry(context.Background(), "le1")
	if err != nil {
		t.Fatalf("GetLogEntry: %v", err)
	}
	if got.Rating != 4.5 {
		t.Errorf("Rating = %f, want 4.5", got.Rating)
	}
}

func TestGetLogEntries(t *testing.T) {
	t.Parallel()
	want := letterboxd.LogEntriesResponse{
		Items: []letterboxd.LogEntry{{ID: "le1"}},
	}
	srv := newTestServer(t, "/log-entries", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetLogEntries(context.Background(), "", 0)
	if err != nil {
		t.Fatalf("GetLogEntries: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetLogEntryComments(t *testing.T) {
	t.Parallel()
	want := letterboxd.CommentsResponse{
		Items: []letterboxd.Comment{{ID: "c1"}},
	}
	srv := newTestServer(t, "/log-entry/le1/comments", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetLogEntryComments(context.Background(), "le1", "", 0)
	if err != nil {
		t.Fatalf("GetLogEntryComments: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetLogEntryStatistics(t *testing.T) {
	t.Parallel()
	want := letterboxd.LogEntryStatistics{Counts: &letterboxd.LogEntryStatCounts{Likes: 50}}
	srv := newTestServer(t, "/log-entry/le1/statistics", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetLogEntryStatistics(context.Background(), "le1")
	if err != nil {
		t.Fatalf("GetLogEntryStatistics: %v", err)
	}
	if got.Counts.Likes != 50 {
		t.Errorf("Likes = %d, want 50", got.Counts.Likes)
	}
}

func TestGetLogEntryRelationship(t *testing.T) {
	t.Parallel()
	want := letterboxd.LogEntryRelationship{Liked: true}
	srv := newTestServer(t, "/log-entry/le1/me", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetLogEntryRelationship(context.Background(), "le1")
	if err != nil {
		t.Fatalf("GetLogEntryRelationship: %v", err)
	}
	if !got.Liked {
		t.Error("Liked = false, want true")
	}
}

func TestUpdateLogEntryRelationship(t *testing.T) {
	t.Parallel()
	want := letterboxd.LogEntryRelationshipUpdateResponse{
		Data: &letterboxd.LogEntryRelationship{Liked: true},
	}
	srv := newMutationTestServer(t, "/log-entry/le1/me", http.MethodPatch, want)
	defer srv.Close()

	c := newClient(t, srv)
	liked := true
	got, err := c.UpdateLogEntryRelationship(context.Background(), "le1", letterboxd.LogEntryRelationshipUpdateRequest{Liked: &liked})
	if err != nil {
		t.Fatalf("UpdateLogEntryRelationship: %v", err)
	}
	if !got.Data.Liked {
		t.Error("Liked = false, want true")
	}
}

func TestCreateLogEntry(t *testing.T) {
	t.Parallel()
	want := letterboxd.LogEntry{ID: "new1", Rating: 4.0}
	srv := newMutationTestServer(t, "/log-entries", http.MethodPost, want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.CreateLogEntry(context.Background(), &letterboxd.LogEntryCreationRequest{FilmID: "f1", Rating: 4.0})
	if err != nil {
		t.Fatalf("CreateLogEntry: %v", err)
	}
	if got.Rating != 4.0 {
		t.Errorf("Rating = %f, want 4.0", got.Rating)
	}
}

func TestUpdateLogEntry(t *testing.T) {
	t.Parallel()
	want := letterboxd.LogEntry{ID: "le1", Rating: 5.0}
	srv := newMutationTestServer(t, "/log-entry/le1", http.MethodPatch, want)
	defer srv.Close()

	c := newClient(t, srv)
	rating := 5.0
	got, err := c.UpdateLogEntry(context.Background(), "le1", &letterboxd.LogEntryUpdateRequest{Rating: &rating})
	if err != nil {
		t.Fatalf("UpdateLogEntry: %v", err)
	}
	if got.Rating != 5.0 {
		t.Errorf("Rating = %f, want 5.0", got.Rating)
	}
}

func TestDeleteLogEntry(t *testing.T) {
	t.Parallel()
	srv := newMutationTestServer(t, "/log-entry/le1", http.MethodDelete, nil)
	defer srv.Close()

	c := newClient(t, srv)
	err := c.DeleteLogEntry(context.Background(), "le1")
	if err != nil {
		t.Fatalf("DeleteLogEntry: %v", err)
	}
}

func TestCreateLogEntryComment(t *testing.T) {
	t.Parallel()
	want := letterboxd.Comment{ID: "c1", Comment: "Great review!"}
	srv := newMutationTestServer(t, "/log-entry/le1/comments", http.MethodPost, want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.CreateLogEntryComment(context.Background(), "le1", letterboxd.CommentCreationRequest{Comment: "Great review!"})
	if err != nil {
		t.Fatalf("CreateLogEntryComment: %v", err)
	}
	if got.Comment != "Great review!" {
		t.Errorf("Comment = %q, want %q", got.Comment, "Great review!")
	}
}

func TestReportLogEntry(t *testing.T) {
	t.Parallel()
	srv := newMutationTestServer(t, "/log-entry/le1/report", http.MethodPost, nil)
	defer srv.Close()

	c := newClient(t, srv)
	err := c.ReportLogEntry(context.Background(), "le1", letterboxd.ReportRequest{Reason: "spam"})
	if err != nil {
		t.Fatalf("ReportLogEntry: %v", err)
	}
}

// Members.

func TestGetMember(t *testing.T) {
	t.Parallel()
	want := letterboxd.Member{ID: "m1", Username: "cinephile", DisplayName: "Film Lover"}
	srv := newTestServer(t, "/member/m1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMember(context.Background(), "m1")
	if err != nil {
		t.Fatalf("GetMember: %v", err)
	}
	if got.Username != "cinephile" {
		t.Errorf("Username = %q, want %q", got.Username, "cinephile")
	}
}

func TestGetMemberStatistics(t *testing.T) {
	t.Parallel()
	want := letterboxd.MemberStatistics{Counts: &letterboxd.MemberStatCounts{Watches: 500}}
	srv := newTestServer(t, "/member/m1/statistics", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMemberStatistics(context.Background(), "m1")
	if err != nil {
		t.Fatalf("GetMemberStatistics: %v", err)
	}
	if got.Counts.Watches != 500 {
		t.Errorf("Watches = %d, want 500", got.Counts.Watches)
	}
}

func TestGetMemberRelationship(t *testing.T) {
	t.Parallel()
	want := letterboxd.MemberRelationship{Following: true}
	srv := newTestServer(t, "/member/m1/me", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMemberRelationship(context.Background(), "m1")
	if err != nil {
		t.Fatalf("GetMemberRelationship: %v", err)
	}
	if !got.Following {
		t.Error("Following = false, want true")
	}
}

func TestUpdateMemberRelationship(t *testing.T) {
	t.Parallel()
	want := letterboxd.MemberRelationshipUpdateResponse{
		Data: &letterboxd.MemberRelationship{Following: true},
	}
	srv := newMutationTestServer(t, "/member/m1/me", http.MethodPatch, want)
	defer srv.Close()

	c := newClient(t, srv)
	following := true
	got, err := c.UpdateMemberRelationship(context.Background(), "m1", letterboxd.MemberRelationshipUpdateRequest{Following: &following})
	if err != nil {
		t.Fatalf("UpdateMemberRelationship: %v", err)
	}
	if !got.Data.Following {
		t.Error("Following = false, want true")
	}
}

func TestGetMemberActivity(t *testing.T) {
	t.Parallel()
	want := letterboxd.ActivityResponse{
		Items: []letterboxd.ActivityItem{{Type: "DiaryEntryActivity"}},
	}
	srv := newTestServer(t, "/member/m1/activity", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMemberActivity(context.Background(), "m1", "", 0)
	if err != nil {
		t.Fatalf("GetMemberActivity: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetMemberWatchlist(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmsResponse{
		Items: []letterboxd.FilmSummary{{ID: "f1", Name: "Watchlist Film"}},
	}
	srv := newTestServer(t, "/member/m1/watchlist", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMemberWatchlist(context.Background(), "m1", "", 0)
	if err != nil {
		t.Fatalf("GetMemberWatchlist: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetMemberLogEntryTags(t *testing.T) {
	t.Parallel()
	want := letterboxd.TagsResponse{Items: []letterboxd.Tag{{Code: "horror", DisplayTag: "Horror"}}}
	srv := newTestServer(t, "/member/m1/log-entry-tags", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMemberLogEntryTags(context.Background(), "m1")
	if err != nil {
		t.Fatalf("GetMemberLogEntryTags: %v", err)
	}
	if len(got.Items) != 1 {
		t.Errorf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetMemberListTags(t *testing.T) {
	t.Parallel()
	want := letterboxd.TagsResponse{Items: []letterboxd.Tag{{Code: "favs", DisplayTag: "Favorites"}}}
	srv := newTestServer(t, "/member/m1/list-tags-2", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMemberListTags(context.Background(), "m1")
	if err != nil {
		t.Fatalf("GetMemberListTags: %v", err)
	}
	if len(got.Items) != 1 {
		t.Errorf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetMembers(t *testing.T) {
	t.Parallel()
	want := letterboxd.MembersResponse{
		Items: []letterboxd.MemberSummary{{ID: "m1", Username: "user1"}},
	}
	srv := newTestServer(t, "/members", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMembers(context.Background(), "", 0)
	if err != nil {
		t.Fatalf("GetMembers: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetPronouns(t *testing.T) {
	t.Parallel()
	want := letterboxd.PronounsResponse{Items: []letterboxd.Pronoun{{ID: "they", Label: "they/them"}}}
	srv := newTestServer(t, "/members/pronouns", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetPronouns(context.Background())
	if err != nil {
		t.Fatalf("GetPronouns: %v", err)
	}
	if len(got.Items) != 1 {
		t.Errorf("Items = %d, want 1", len(got.Items))
	}
}

func TestReportMember(t *testing.T) {
	t.Parallel()
	srv := newMutationTestServer(t, "/member/m1/report", http.MethodPost, nil)
	defer srv.Close()

	c := newClient(t, srv)
	err := c.ReportMember(context.Background(), "m1", letterboxd.ReportRequest{Reason: "spam"})
	if err != nil {
		t.Fatalf("ReportMember: %v", err)
	}
}

// Me.

func TestGetMe(t *testing.T) {
	t.Parallel()
	want := letterboxd.MemberAccount{
		Member: &letterboxd.Member{ID: "me1", Username: "myself"},
	}
	srv := newTestServer(t, "/me", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe: %v", err)
	}
	if got.Member.Username != "myself" {
		t.Errorf("Username = %q, want %q", got.Member.Username, "myself")
	}
}

func TestUpdateMe(t *testing.T) {
	t.Parallel()
	want := letterboxd.MemberSettingsUpdateResponse{
		Data: &letterboxd.MemberAccount{Member: &letterboxd.Member{Bio: "Updated bio"}},
	}
	srv := newMutationTestServer(t, "/me", http.MethodPatch, want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.UpdateMe(context.Background(), &letterboxd.MemberSettingsUpdateRequest{Bio: "Updated bio"})
	if err != nil {
		t.Fatalf("UpdateMe: %v", err)
	}
	if got.Data.Member.Bio != "Updated bio" {
		t.Errorf("Bio = %q, want %q", got.Data.Member.Bio, "Updated bio")
	}
}

// Comments.

func TestUpdateComment(t *testing.T) {
	t.Parallel()
	want := letterboxd.Comment{ID: "c1", Comment: "Edited"}
	srv := newMutationTestServer(t, "/comment/c1", http.MethodPatch, want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.UpdateComment(context.Background(), "c1", letterboxd.CommentUpdateRequest{Comment: "Edited"})
	if err != nil {
		t.Fatalf("UpdateComment: %v", err)
	}
	if got.Comment != "Edited" {
		t.Errorf("Comment = %q, want %q", got.Comment, "Edited")
	}
}

func TestDeleteComment(t *testing.T) {
	t.Parallel()
	srv := newMutationTestServer(t, "/comment/c1", http.MethodDelete, nil)
	defer srv.Close()

	c := newClient(t, srv)
	err := c.DeleteComment(context.Background(), "c1")
	if err != nil {
		t.Fatalf("DeleteComment: %v", err)
	}
}

func TestReportComment(t *testing.T) {
	t.Parallel()
	srv := newMutationTestServer(t, "/comment/c1/report", http.MethodPost, nil)
	defer srv.Close()

	c := newClient(t, srv)
	err := c.ReportComment(context.Background(), "c1", letterboxd.ReportRequest{Reason: "abuse"})
	if err != nil {
		t.Fatalf("ReportComment: %v", err)
	}
}

// Stories.

func TestGetStory(t *testing.T) {
	t.Parallel()
	want := letterboxd.Story{ID: "st1", Name: "Best of 2025"}
	srv := newTestServer(t, "/story/st1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetStory(context.Background(), "st1")
	if err != nil {
		t.Fatalf("GetStory: %v", err)
	}
	if got.Name != "Best of 2025" {
		t.Errorf("Name = %q, want %q", got.Name, "Best of 2025")
	}
}

func TestGetStories(t *testing.T) {
	t.Parallel()
	want := letterboxd.StoriesResponse{
		Items: []letterboxd.StorySummary{{ID: "st1", Name: "Story 1"}},
	}
	srv := newTestServer(t, "/stories", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetStories(context.Background(), "", 0)
	if err != nil {
		t.Fatalf("GetStories: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetStoryComments(t *testing.T) {
	t.Parallel()
	want := letterboxd.CommentsResponse{
		Items: []letterboxd.Comment{{ID: "c1"}},
	}
	srv := newTestServer(t, "/story/st1/comments", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetStoryComments(context.Background(), "st1", "", 0)
	if err != nil {
		t.Fatalf("GetStoryComments: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

func TestGetStoryStatistics(t *testing.T) {
	t.Parallel()
	want := letterboxd.StoryStatistics{Counts: &letterboxd.StoryStatCounts{Likes: 25}}
	srv := newTestServer(t, "/story/st1/statistics", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetStoryStatistics(context.Background(), "st1")
	if err != nil {
		t.Fatalf("GetStoryStatistics: %v", err)
	}
	if got.Counts.Likes != 25 {
		t.Errorf("Likes = %d, want 25", got.Counts.Likes)
	}
}

func TestGetStoryRelationship(t *testing.T) {
	t.Parallel()
	want := letterboxd.StoryRelationship{Liked: true}
	srv := newTestServer(t, "/story/st1/me", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetStoryRelationship(context.Background(), "st1")
	if err != nil {
		t.Fatalf("GetStoryRelationship: %v", err)
	}
	if !got.Liked {
		t.Error("Liked = false, want true")
	}
}

func TestUpdateStoryRelationship(t *testing.T) {
	t.Parallel()
	want := letterboxd.StoryRelationshipUpdateResponse{
		Data: &letterboxd.StoryRelationship{Liked: true},
	}
	srv := newMutationTestServer(t, "/story/st1/me", http.MethodPatch, want)
	defer srv.Close()

	c := newClient(t, srv)
	liked := true
	got, err := c.UpdateStoryRelationship(context.Background(), "st1", letterboxd.StoryRelationshipUpdateRequest{Liked: &liked})
	if err != nil {
		t.Fatalf("UpdateStoryRelationship: %v", err)
	}
	if !got.Data.Liked {
		t.Error("Liked = false, want true")
	}
}

func TestCreateStoryComment(t *testing.T) {
	t.Parallel()
	want := letterboxd.Comment{ID: "c1", Comment: "Nice story!"}
	srv := newMutationTestServer(t, "/story/st1/comments", http.MethodPost, want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.CreateStoryComment(context.Background(), "st1", letterboxd.CommentCreationRequest{Comment: "Nice story!"})
	if err != nil {
		t.Fatalf("CreateStoryComment: %v", err)
	}
	if got.Comment != "Nice story!" {
		t.Errorf("Comment = %q, want %q", got.Comment, "Nice story!")
	}
}

// Search.

func TestSearch(t *testing.T) {
	t.Parallel()
	want := letterboxd.SearchResponse{
		Items: []letterboxd.SearchItem{
			{Type: "FilmSearchItem", Score: 10.5, Film: &letterboxd.FilmSummary{ID: "f1", Name: "Inception"}},
		},
	}
	srv := newTestServer(t, "/search?input=Inception", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.Search(context.Background(), "Inception", "", 0)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
	if got.Items[0].Film.Name != "Inception" {
		t.Errorf("Film.Name = %q, want %q", got.Items[0].Film.Name, "Inception")
	}
}

func TestSearchWithPagination(t *testing.T) {
	t.Parallel()
	want := letterboxd.SearchResponse{Cursor: "page2"}
	srv := newTestServer(t, "/search?cursor=page1&input=test&perPage=10", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.Search(context.Background(), "test", "page1", 10)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if got.Cursor != "page2" {
		t.Errorf("Cursor = %q, want %q", got.Cursor, "page2")
	}
}

// News.

func TestGetNews(t *testing.T) {
	t.Parallel()
	want := letterboxd.NewsResponse{
		Items: []letterboxd.NewsItem{{ID: "n1", Title: "New Feature"}},
	}
	srv := newTestServer(t, "/news", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetNews(context.Background(), "", 0)
	if err != nil {
		t.Fatalf("GetNews: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

// Auth.

func TestCheckUsername(t *testing.T) {
	t.Parallel()
	want := letterboxd.UsernameCheckResponse{Result: "Available"}
	srv := newTestServer(t, "/auth/username-check?username=newuser", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.CheckUsername(context.Background(), "newuser")
	if err != nil {
		t.Fatalf("CheckUsername: %v", err)
	}
	if got.Result != "Available" {
		t.Errorf("Result = %q, want %q", got.Result, "Available")
	}
}

// Cursor pagination.

func TestCursorPagination(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmsResponse{
		Cursor: "next",
		Items:  []letterboxd.FilmSummary{{ID: "f1"}},
	}
	srv := newTestServer(t, "/films?cursor=abc&perPage=20", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetFilms(context.Background(), "abc", 20)
	if err != nil {
		t.Fatalf("GetFilms: %v", err)
	}
	if got.Cursor != "next" {
		t.Errorf("Cursor = %q, want %q", got.Cursor, "next")
	}
}

// GetLogEntryMembers.

func TestGetLogEntryMembers(t *testing.T) {
	t.Parallel()
	want := letterboxd.FilmMembersResponse{
		Items: []letterboxd.FilmMemberRelationshipItem{
			{Member: &letterboxd.MemberSummary{ID: "u1"}},
		},
	}
	srv := newTestServer(t, "/log-entry/le1/members", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetLogEntryMembers(context.Background(), "le1", "", 0)
	if err != nil {
		t.Fatalf("GetLogEntryMembers: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("Items = %d, want 1", len(got.Items))
	}
}

// AddFilmsToLists.

func TestAddFilmsToLists(t *testing.T) {
	t.Parallel()
	srv := newMutationTestServer(t, "/lists", http.MethodPatch, nil)
	defer srv.Close()

	c := newClient(t, srv)
	err := c.AddFilmsToLists(context.Background(), letterboxd.ListAddEntriesRequest{
		Lists: []string{"lst1"},
		Films: []letterboxd.ListEntryRequest{{Film: "f1"}},
	})
	if err != nil {
		t.Fatalf("AddFilmsToLists: %v", err)
	}
}
