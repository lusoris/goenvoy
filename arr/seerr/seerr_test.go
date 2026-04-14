package seerr_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/arr/seerr"
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

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		c, err := seerr.New("http://localhost:5055", "test-key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c == nil {
			t.Fatal("expected non-nil client")
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		t.Parallel()
		_, err := seerr.New("://bad", "test-key")
		if err == nil {
			t.Fatal("expected error for invalid URL")
		}
	})
}

func TestGetStatus(t *testing.T) {
	t.Parallel()

	want := seerr.StatusResponse{Version: "1.33.2", UpdateAvailable: false}

	srv := newTestServer(t, http.MethodGet, "/api/v1/status", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus: %v", err)
	}
	if got.Version != "1.33.2" {
		t.Errorf("Version = %q, want %q", got.Version, "1.33.2")
	}
}

func TestGetMe(t *testing.T) {
	t.Parallel()

	want := seerr.User{ID: 1, Email: "admin@example.com", Permissions: 2}

	srv := newTestServer(t, http.MethodGet, "/api/v1/auth/me", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe: %v", err)
	}
	if got.Email != "admin@example.com" {
		t.Errorf("Email = %q, want %q", got.Email, "admin@example.com")
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()

	want := seerr.SearchResults{Page: 1, TotalPages: 5, TotalResults: 100}

	srv := newTestServer(t, http.MethodGet, "/api/v1/search?query=Mulan&page=1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.Search(context.Background(), "Mulan", 1)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if got.TotalResults != 100 {
		t.Errorf("TotalResults = %d, want 100", got.TotalResults)
	}
}

func TestDiscoverMovies(t *testing.T) {
	t.Parallel()

	want := seerr.SearchResults{Page: 1, TotalPages: 10, TotalResults: 200}

	srv := newTestServer(t, http.MethodGet, "/api/v1/discover/movies?page=1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.DiscoverMovies(context.Background(), 1)
	if err != nil {
		t.Fatalf("DiscoverMovies: %v", err)
	}
	if got.TotalResults != 200 {
		t.Errorf("TotalResults = %d, want 200", got.TotalResults)
	}
}

func TestDiscoverTrending(t *testing.T) {
	t.Parallel()

	want := seerr.SearchResults{Page: 1, TotalPages: 1, TotalResults: 20}

	srv := newTestServer(t, http.MethodGet, "/api/v1/discover/trending?page=1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.DiscoverTrending(context.Background(), 1)
	if err != nil {
		t.Fatalf("DiscoverTrending: %v", err)
	}
	if got.TotalResults != 20 {
		t.Errorf("TotalResults = %d, want 20", got.TotalResults)
	}
}

func TestGetRequests(t *testing.T) {
	t.Parallel()

	want := struct {
		PageInfo seerr.PageInfo       `json:"pageInfo"`
		Results  []seerr.MediaRequest `json:"results"`
	}{
		PageInfo: seerr.PageInfo{Page: 1, Pages: 1, Results: 1},
		Results: []seerr.MediaRequest{
			{ID: 1, Status: 2, Is4k: false},
		},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/request?take=20&skip=0&filter=pending", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	results, pageInfo, err := c.GetRequests(context.Background(), 20, 0, "pending")
	if err != nil {
		t.Fatalf("GetRequests: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len = %d, want 1", len(results))
	}
	if pageInfo.Results != 1 {
		t.Errorf("PageInfo.Results = %d, want 1", pageInfo.Results)
	}
}

func TestCreateRequest(t *testing.T) {
	t.Parallel()

	want := seerr.MediaRequest{ID: 5, Status: 1}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var body seerr.CreateRequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.MediaType != "movie" {
			t.Errorf("MediaType = %q, want %q", body.MediaType, "movie")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.CreateRequest(context.Background(), &seerr.CreateRequestBody{
		MediaType: "movie",
		MediaID:   337401,
	})
	if err != nil {
		t.Fatalf("CreateRequest: %v", err)
	}
	if got.ID != 5 {
		t.Errorf("ID = %d, want 5", got.ID)
	}
}

func TestDeleteRequest(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v1/request/1", nil)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	if err := c.DeleteRequest(context.Background(), 1); err != nil {
		t.Fatalf("DeleteRequest: %v", err)
	}
}

func TestApproveRequest(t *testing.T) {
	t.Parallel()

	want := seerr.MediaRequest{ID: 1, Status: 2}

	srv := newTestServer(t, http.MethodPost, "/api/v1/request/1/approve", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.ApproveRequest(context.Background(), 1)
	if err != nil {
		t.Fatalf("ApproveRequest: %v", err)
	}
	if got.Status != 2 {
		t.Errorf("Status = %d, want 2", got.Status)
	}
}

func TestDeclineRequest(t *testing.T) {
	t.Parallel()

	want := seerr.MediaRequest{ID: 1, Status: 3}

	srv := newTestServer(t, http.MethodPost, "/api/v1/request/1/decline", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.DeclineRequest(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeclineRequest: %v", err)
	}
	if got.Status != 3 {
		t.Errorf("Status = %d, want 3", got.Status)
	}
}

func TestGetRequestCount(t *testing.T) {
	t.Parallel()

	want := seerr.RequestCount{Total: 50, Movie: 30, TV: 20, Pending: 5}

	srv := newTestServer(t, http.MethodGet, "/api/v1/request/count", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetRequestCount(context.Background())
	if err != nil {
		t.Fatalf("GetRequestCount: %v", err)
	}
	if got.Total != 50 {
		t.Errorf("Total = %d, want 50", got.Total)
	}
}

func TestGetMovie(t *testing.T) {
	t.Parallel()

	want := seerr.MovieDetails{
		ID: 337401, Title: "Mulan", VoteAverage: 7.0,
		Genres: []seerr.Genre{{ID: 28, Name: "Action"}},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/movie/337401", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetMovie(context.Background(), 337401)
	if err != nil {
		t.Fatalf("GetMovie: %v", err)
	}
	if got.Title != "Mulan" {
		t.Errorf("Title = %q, want %q", got.Title, "Mulan")
	}
}

func TestGetTV(t *testing.T) {
	t.Parallel()

	want := seerr.TvDetails{
		ID: 76479, Name: "The Boys", Status: "Returning Series",
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/tv/76479", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetTV(context.Background(), 76479)
	if err != nil {
		t.Fatalf("GetTV: %v", err)
	}
	if got.Name != "The Boys" {
		t.Errorf("Name = %q, want %q", got.Name, "The Boys")
	}
}

func TestGetTVSeason(t *testing.T) {
	t.Parallel()

	want := seerr.Season{
		ID: 1, Name: "Season 1", SeasonNumber: 1,
		Episodes: []seerr.Episode{{ID: 100, Name: "Pilot", EpisodeNumber: 1, SeasonNumber: 1}},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/tv/76479/season/1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetTVSeason(context.Background(), 76479, 1)
	if err != nil {
		t.Fatalf("GetTVSeason: %v", err)
	}
	if len(got.Episodes) != 1 {
		t.Fatalf("Episodes len = %d, want 1", len(got.Episodes))
	}
	if got.Episodes[0].Name != "Pilot" {
		t.Errorf("Episode Name = %q, want %q", got.Episodes[0].Name, "Pilot")
	}
}

func TestGetMedia(t *testing.T) {
	t.Parallel()

	want := struct {
		PageInfo seerr.PageInfo    `json:"pageInfo"`
		Results  []seerr.MediaInfo `json:"results"`
	}{
		PageInfo: seerr.PageInfo{Page: 1, Pages: 1, Results: 2},
		Results: []seerr.MediaInfo{
			{ID: 1, TmdbID: 337401, Status: 5},
			{ID: 2, TmdbID: 76479, Status: 3},
		},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/media?take=20&skip=0", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	results, pageInfo, err := c.GetMedia(context.Background(), 20, 0)
	if err != nil {
		t.Fatalf("GetMedia: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("len = %d, want 2", len(results))
	}
	if pageInfo.Results != 2 {
		t.Errorf("PageInfo.Results = %d, want 2", pageInfo.Results)
	}
}

func TestDeleteMedia(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v1/media/1", nil)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	if err := c.DeleteMedia(context.Background(), 1); err != nil {
		t.Fatalf("DeleteMedia: %v", err)
	}
}

func TestGetUsers(t *testing.T) {
	t.Parallel()

	want := struct {
		PageInfo seerr.PageInfo `json:"pageInfo"`
		Results  []seerr.User   `json:"results"`
	}{
		PageInfo: seerr.PageInfo{Page: 1, Pages: 1, Results: 1},
		Results:  []seerr.User{{ID: 1, Email: "admin@example.com"}},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/user?take=20&skip=0", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	results, _, err := c.GetUsers(context.Background(), 20, 0)
	if err != nil {
		t.Fatalf("GetUsers: %v", err)
	}
	if results[0].Email != "admin@example.com" {
		t.Errorf("Email = %q, want %q", results[0].Email, "admin@example.com")
	}
}

func TestGetUserQuota(t *testing.T) {
	t.Parallel()

	want := seerr.UserQuota{
		Movie: &seerr.QuotaDetail{Days: 7, Limit: 10, Used: 3, Remaining: 7},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/user/1/quota", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetUserQuota(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUserQuota: %v", err)
	}
	if got.Movie.Remaining != 7 {
		t.Errorf("Movie.Remaining = %d, want 7", got.Movie.Remaining)
	}
}

func TestGetIssues(t *testing.T) {
	t.Parallel()

	want := struct {
		PageInfo seerr.PageInfo `json:"pageInfo"`
		Results  []seerr.Issue  `json:"results"`
	}{
		PageInfo: seerr.PageInfo{Page: 1, Pages: 1, Results: 1},
		Results:  []seerr.Issue{{ID: 1, IssueType: 1}},
	}

	srv := newTestServer(t, http.MethodGet, "/api/v1/issue?take=20&skip=0&filter=open", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	results, _, err := c.GetIssues(context.Background(), 20, 0, "open")
	if err != nil {
		t.Fatalf("GetIssues: %v", err)
	}
	if results[0].IssueType != 1 {
		t.Errorf("IssueType = %d, want 1", results[0].IssueType)
	}
}

func TestCreateIssue(t *testing.T) {
	t.Parallel()

	want := seerr.Issue{ID: 10, IssueType: 2}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var body seerr.CreateIssueBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.IssueType != 2 {
			t.Errorf("IssueType = %d, want 2", body.IssueType)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.CreateIssue(context.Background(), &seerr.CreateIssueBody{
		IssueType: 2,
		Message:   "Audio sync issue",
		MediaID:   337401,
	})
	if err != nil {
		t.Fatalf("CreateIssue: %v", err)
	}
	if got.ID != 10 {
		t.Errorf("ID = %d, want 10", got.ID)
	}
}

func TestAddIssueComment(t *testing.T) {
	t.Parallel()

	want := seerr.Issue{ID: 1, IssueType: 1}

	srv := newTestServer(t, http.MethodPost, "/api/v1/issue/1/comment", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.AddIssueComment(context.Background(), 1, "Still broken")
	if err != nil {
		t.Fatalf("AddIssueComment: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
}

func TestGetIssueCount(t *testing.T) {
	t.Parallel()

	want := seerr.IssueCount{Total: 15, Open: 10, Closed: 5, Video: 3, Audio: 2, Subtitles: 5, Others: 5}

	srv := newTestServer(t, http.MethodGet, "/api/v1/issue/count", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetIssueCount(context.Background())
	if err != nil {
		t.Fatalf("GetIssueCount: %v", err)
	}
	if got.Open != 10 {
		t.Errorf("Open = %d, want 10", got.Open)
	}
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	}))
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "bad-key")
	_, err := c.GetStatus(context.Background())
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

func TestDiscoverTV(t *testing.T) {
	t.Parallel()

	want := seerr.SearchResults{Page: 1, TotalPages: 10, TotalResults: 200}

	srv := newTestServer(t, http.MethodGet, "/api/v1/discover/tv?page=1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.DiscoverTV(context.Background(), 1)
	if err != nil {
		t.Fatalf("DiscoverTV: %v", err)
	}
	if got.TotalResults != 200 {
		t.Errorf("TotalResults = %d, want 200", got.TotalResults)
	}
}

func TestDiscoverUpcomingMovies(t *testing.T) {
	t.Parallel()

	want := seerr.SearchResults{Page: 1, TotalResults: 50}

	srv := newTestServer(t, http.MethodGet, "/api/v1/discover/movies/upcoming?page=1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.DiscoverUpcomingMovies(context.Background(), 1)
	if err != nil {
		t.Fatalf("DiscoverUpcomingMovies: %v", err)
	}
	if got.TotalResults != 50 {
		t.Errorf("TotalResults = %d", got.TotalResults)
	}
}

func TestDiscoverUpcomingTV(t *testing.T) {
	t.Parallel()

	want := seerr.SearchResults{Page: 1, TotalResults: 30}

	srv := newTestServer(t, http.MethodGet, "/api/v1/discover/tv/upcoming?page=1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.DiscoverUpcomingTV(context.Background(), 1)
	if err != nil {
		t.Fatalf("DiscoverUpcomingTV: %v", err)
	}
	if got.TotalResults != 30 {
		t.Errorf("TotalResults = %d", got.TotalResults)
	}
}

func TestGetRequest(t *testing.T) {
	t.Parallel()

	want := seerr.MediaRequest{ID: 1, Status: 2}

	srv := newTestServer(t, http.MethodGet, "/api/v1/request/1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetRequest(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetRequest: %v", err)
	}
	if got.Status != 2 {
		t.Errorf("Status = %d", got.Status)
	}
}

func TestRetryRequest(t *testing.T) {
	t.Parallel()

	want := seerr.MediaRequest{ID: 1, Status: 1}

	srv := newTestServer(t, http.MethodPost, "/api/v1/request/1/retry", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.RetryRequest(context.Background(), 1)
	if err != nil {
		t.Fatalf("RetryRequest: %v", err)
	}
	if got.Status != 1 {
		t.Errorf("Status = %d", got.Status)
	}
}

func TestGetMovieRecommendations(t *testing.T) {
	t.Parallel()

	want := seerr.SearchResults{Page: 1, TotalResults: 10}

	srv := newTestServer(t, http.MethodGet, "/api/v1/movie/337401/recommendations?page=1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetMovieRecommendations(context.Background(), 337401, 1)
	if err != nil {
		t.Fatalf("GetMovieRecommendations: %v", err)
	}
	if got.TotalResults != 10 {
		t.Errorf("TotalResults = %d", got.TotalResults)
	}
}

func TestGetMovieSimilar(t *testing.T) {
	t.Parallel()

	want := seerr.SearchResults{Page: 1, TotalResults: 15}

	srv := newTestServer(t, http.MethodGet, "/api/v1/movie/337401/similar?page=1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetMovieSimilar(context.Background(), 337401, 1)
	if err != nil {
		t.Fatalf("GetMovieSimilar: %v", err)
	}
	if got.TotalResults != 15 {
		t.Errorf("TotalResults = %d", got.TotalResults)
	}
}

func TestGetTVRecommendations(t *testing.T) {
	t.Parallel()

	want := seerr.SearchResults{Page: 1, TotalResults: 8}

	srv := newTestServer(t, http.MethodGet, "/api/v1/tv/76479/recommendations?page=1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetTVRecommendations(context.Background(), 76479, 1)
	if err != nil {
		t.Fatalf("GetTVRecommendations: %v", err)
	}
	if got.TotalResults != 8 {
		t.Errorf("TotalResults = %d", got.TotalResults)
	}
}

func TestGetTVSimilar(t *testing.T) {
	t.Parallel()

	want := seerr.SearchResults{Page: 1, TotalResults: 12}

	srv := newTestServer(t, http.MethodGet, "/api/v1/tv/76479/similar?page=1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetTVSimilar(context.Background(), 76479, 1)
	if err != nil {
		t.Fatalf("GetTVSimilar: %v", err)
	}
	if got.TotalResults != 12 {
		t.Errorf("TotalResults = %d", got.TotalResults)
	}
}

func TestUpdateMediaStatus(t *testing.T) {
	t.Parallel()

	want := seerr.MediaInfo{ID: 1, Status: 5}

	srv := newTestServer(t, http.MethodPost, "/api/v1/media/1/available", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.UpdateMediaStatus(context.Background(), 1, "available", false)
	if err != nil {
		t.Fatalf("UpdateMediaStatus: %v", err)
	}
	if got.Status != 5 {
		t.Errorf("Status = %d", got.Status)
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	want := seerr.User{ID: 1, Email: "admin@example.com"}

	srv := newTestServer(t, http.MethodGet, "/api/v1/user/1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.Email != "admin@example.com" {
		t.Errorf("Email = %q", got.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v1/user/1", nil)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	if err := c.DeleteUser(context.Background(), 1); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
}

func TestGetIssue(t *testing.T) {
	t.Parallel()

	want := seerr.Issue{ID: 1, IssueType: 1}

	srv := newTestServer(t, http.MethodGet, "/api/v1/issue/1", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.GetIssue(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetIssue: %v", err)
	}
	if got.IssueType != 1 {
		t.Errorf("IssueType = %d", got.IssueType)
	}
}

func TestDeleteIssue(t *testing.T) {
	t.Parallel()

	srv := newTestServer(t, http.MethodDelete, "/api/v1/issue/1", nil)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	if err := c.DeleteIssue(context.Background(), 1); err != nil {
		t.Fatalf("DeleteIssue: %v", err)
	}
}

func TestResolveIssue(t *testing.T) {
	t.Parallel()

	want := seerr.Issue{ID: 1, IssueType: 1}

	srv := newTestServer(t, http.MethodPost, "/api/v1/issue/1/resolved", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.ResolveIssue(context.Background(), 1)
	if err != nil {
		t.Fatalf("ResolveIssue: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d", got.ID)
	}
}

func TestReopenIssue(t *testing.T) {
	t.Parallel()

	want := seerr.Issue{ID: 1, IssueType: 1}

	srv := newTestServer(t, http.MethodPost, "/api/v1/issue/1/open", want)
	defer srv.Close()

	c, _ := seerr.New(srv.URL, "test-key")
	got, err := c.ReopenIssue(context.Background(), 1)
	if err != nil {
		t.Fatalf("ReopenIssue: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d", got.ID)
	}
}
