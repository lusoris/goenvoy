package tmdb_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/metadata/movie/tmdb"
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

func newClient(t *testing.T, srv *httptest.Server) *tmdb.Client {
	t.Helper()
	return tmdb.New("test-token", tmdb.WithBaseURL(srv.URL))
}

func TestNew(t *testing.T) {
	t.Parallel()
	c := tmdb.New("token")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGetConfiguration(t *testing.T) {
	t.Parallel()

	want := tmdb.Configuration{
		Images: tmdb.ImageConfiguration{
			SecureBaseURL: "https://image.tmdb.org/t/p/",
			PosterSizes:   []string{"w92", "w154", "w185", "w342", "w500", "w780", "original"},
		},
	}

	srv := newTestServer(t, "/configuration", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetConfiguration(context.Background())
	if err != nil {
		t.Fatalf("GetConfiguration: %v", err)
	}
	if got.Images.SecureBaseURL != "https://image.tmdb.org/t/p/" {
		t.Errorf("SecureBaseURL = %q, want %q", got.Images.SecureBaseURL, "https://image.tmdb.org/t/p/")
	}
}

func TestSearchMovies(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.MovieResult]{
		Page: 1, TotalPages: 1, TotalResults: 1,
		Results: []tmdb.MovieResult{{ID: 11, Title: "Star Wars"}},
	}

	srv := newTestServer(t, "/search/movie?query=Star+Wars&language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.SearchMovies(context.Background(), "Star Wars", "en-US", 1)
	if err != nil {
		t.Fatalf("SearchMovies: %v", err)
	}
	if got.Results[0].Title != "Star Wars" {
		t.Errorf("Title = %q, want %q", got.Results[0].Title, "Star Wars")
	}
}

func TestSearchTV(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.TVResult]{
		Page: 1, TotalPages: 1, TotalResults: 1,
		Results: []tmdb.TVResult{{ID: 1396, Name: "Breaking Bad"}},
	}

	srv := newTestServer(t, "/search/tv?query=Breaking+Bad&language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.SearchTV(context.Background(), "Breaking Bad", "en-US", 1)
	if err != nil {
		t.Fatalf("SearchTV: %v", err)
	}
	if got.Results[0].Name != "Breaking Bad" {
		t.Errorf("Name = %q, want %q", got.Results[0].Name, "Breaking Bad")
	}
}

func TestSearchMulti(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.MultiResult]{
		Page: 1, TotalPages: 1, TotalResults: 2,
		Results: []tmdb.MultiResult{
			{ID: 11, MediaType: "movie", Title: "Star Wars"},
			{ID: 1, MediaType: "person", Name: "Mark Hamill"},
		},
	}

	srv := newTestServer(t, "/search/multi?query=Star+Wars&language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.SearchMulti(context.Background(), "Star Wars", "en-US", 1)
	if err != nil {
		t.Fatalf("SearchMulti: %v", err)
	}
	if len(got.Results) != 2 {
		t.Fatalf("len = %d, want 2", len(got.Results))
	}
}

func TestGetMovie(t *testing.T) {
	t.Parallel()

	want := tmdb.MovieDetails{
		ID: 11, Title: "Star Wars", IMDbID: "tt0076759",
		Genres: []tmdb.Genre{{ID: 12, Name: "Adventure"}, {ID: 28, Name: "Action"}},
	}

	srv := newTestServer(t, "/movie/11?language=en-US", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMovie(context.Background(), 11, "en-US")
	if err != nil {
		t.Fatalf("GetMovie: %v", err)
	}
	if got.Title != "Star Wars" {
		t.Errorf("Title = %q, want %q", got.Title, "Star Wars")
	}
	if len(got.Genres) != 2 {
		t.Errorf("Genres len = %d, want 2", len(got.Genres))
	}
}

func TestGetMovieCredits(t *testing.T) {
	t.Parallel()

	want := tmdb.Credits{
		ID:   11,
		Cast: []tmdb.CastMember{{ID: 2, Name: "Mark Hamill", Character: "Luke Skywalker"}},
		Crew: []tmdb.CrewMember{{ID: 1, Name: "George Lucas", Job: "Director"}},
	}

	srv := newTestServer(t, "/movie/11/credits?language=en-US", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMovieCredits(context.Background(), 11, "en-US")
	if err != nil {
		t.Fatalf("GetMovieCredits: %v", err)
	}
	if got.Cast[0].Character != "Luke Skywalker" {
		t.Errorf("Character = %q, want %q", got.Cast[0].Character, "Luke Skywalker")
	}
}

func TestGetMovieImages(t *testing.T) {
	t.Parallel()

	want := tmdb.Images{
		ID:      11,
		Posters: []tmdb.ImageItem{{FilePath: "/poster.jpg", Width: 500, Height: 750}},
	}

	srv := newTestServer(t, "/movie/11/images", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMovieImages(context.Background(), 11)
	if err != nil {
		t.Fatalf("GetMovieImages: %v", err)
	}
	if len(got.Posters) != 1 {
		t.Fatalf("Posters len = %d, want 1", len(got.Posters))
	}
}

func TestGetMovieExternalIDs(t *testing.T) {
	t.Parallel()

	want := tmdb.ExternalIDs{ID: 11, IMDbID: "tt0076759"}

	srv := newTestServer(t, "/movie/11/external_ids", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMovieExternalIDs(context.Background(), 11)
	if err != nil {
		t.Fatalf("GetMovieExternalIDs: %v", err)
	}
	if got.IMDbID != "tt0076759" {
		t.Errorf("IMDbID = %q, want %q", got.IMDbID, "tt0076759")
	}
}

func TestGetTV(t *testing.T) {
	t.Parallel()

	want := tmdb.TVDetails{
		ID: 1396, Name: "Breaking Bad", Status: "Ended",
		NumberOfSeasons: 5, NumberOfEpisodes: 62,
	}

	srv := newTestServer(t, "/tv/1396?language=en-US", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetTV(context.Background(), 1396, "en-US")
	if err != nil {
		t.Fatalf("GetTV: %v", err)
	}
	if got.Name != "Breaking Bad" {
		t.Errorf("Name = %q, want %q", got.Name, "Breaking Bad")
	}
	if got.NumberOfSeasons != 5 {
		t.Errorf("NumberOfSeasons = %d, want 5", got.NumberOfSeasons)
	}
}

func TestGetTVSeason(t *testing.T) {
	t.Parallel()

	want := tmdb.SeasonDetails{
		ID: 3572, Name: "Season 1", SeasonNumber: 1,
		Episodes: []tmdb.Episode{
			{ID: 62085, Name: "Pilot", EpisodeNumber: 1, SeasonNumber: 1},
		},
	}

	srv := newTestServer(t, "/tv/1396/season/1?language=en-US", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetTVSeason(context.Background(), 1396, 1, "en-US")
	if err != nil {
		t.Fatalf("GetTVSeason: %v", err)
	}
	if len(got.Episodes) != 1 {
		t.Fatalf("Episodes len = %d, want 1", len(got.Episodes))
	}
}

func TestGetPerson(t *testing.T) {
	t.Parallel()

	want := tmdb.PersonDetails{ID: 2, Name: "Mark Hamill", KnownForDepartment: "Acting"}

	srv := newTestServer(t, "/person/2?language=en-US", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetPerson(context.Background(), 2, "en-US")
	if err != nil {
		t.Fatalf("GetPerson: %v", err)
	}
	if got.Name != "Mark Hamill" {
		t.Errorf("Name = %q, want %q", got.Name, "Mark Hamill")
	}
}

func TestGetTrending(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.MultiResult]{
		Page: 1, TotalPages: 1, TotalResults: 1,
		Results: []tmdb.MultiResult{{ID: 11, MediaType: "movie", Title: "Star Wars"}},
	}

	srv := newTestServer(t, "/trending/all/day?language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetTrending(context.Background(), "all", "day", "en-US", 1)
	if err != nil {
		t.Fatalf("GetTrending: %v", err)
	}
	if got.Results[0].Title != "Star Wars" {
		t.Errorf("Title = %q, want %q", got.Results[0].Title, "Star Wars")
	}
}

func TestGetGenresMovie(t *testing.T) {
	t.Parallel()

	want := struct {
		Genres []tmdb.Genre `json:"genres"`
	}{Genres: []tmdb.Genre{{ID: 28, Name: "Action"}, {ID: 12, Name: "Adventure"}}}

	srv := newTestServer(t, "/genre/movie/list?language=en-US", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetGenresMovie(context.Background(), "en-US")
	if err != nil {
		t.Fatalf("GetGenresMovie: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}

func TestFindByExternalID(t *testing.T) {
	t.Parallel()

	want := tmdb.FindResult{
		MovieResults: []tmdb.MovieResult{{ID: 11, Title: "Star Wars"}},
	}

	srv := newTestServer(t, "/find/tt0076759?external_source=imdb_id&language=en-US", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.FindByExternalID(context.Background(), "tt0076759", "imdb_id", "en-US")
	if err != nil {
		t.Fatalf("FindByExternalID: %v", err)
	}
	if len(got.MovieResults) != 1 {
		t.Fatalf("MovieResults len = %d, want 1", len(got.MovieResults))
	}
}

func TestGetPopularMovies(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.MovieResult]{
		Page: 1, TotalPages: 500, TotalResults: 10000,
		Results: []tmdb.MovieResult{{ID: 1, Title: "Popular Movie"}},
	}

	srv := newTestServer(t, "/movie/popular?language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetPopularMovies(context.Background(), "en-US", 1)
	if err != nil {
		t.Fatalf("GetPopularMovies: %v", err)
	}
	if got.TotalResults != 10000 {
		t.Errorf("TotalResults = %d, want 10000", got.TotalResults)
	}
}

func TestDiscoverMovies(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.MovieResult]{
		Page: 1, TotalPages: 1, TotalResults: 1,
		Results: []tmdb.MovieResult{{ID: 100, Title: "Disco Movie"}},
	}

	srv := newTestServer(t, "/discover/movie?language=en-US&page=1&sort_by=popularity.desc", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.DiscoverMovies(context.Background(), "en-US", 1, "&sort_by=popularity.desc")
	if err != nil {
		t.Fatalf("DiscoverMovies: %v", err)
	}
	if got.Results[0].Title != "Disco Movie" {
		t.Errorf("Title = %q, want %q", got.Results[0].Title, "Disco Movie")
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(tmdb.APIError{
			StatusMessage: "Invalid API key",
			ErrorCode:     7,
		})
	}))
	defer srv.Close()

	c := tmdb.New("bad-token", tmdb.WithBaseURL(srv.URL))
	_, err := c.GetMovie(context.Background(), 11, "en-US")
	if err == nil {
		t.Fatal("expected error for 401 response")
	}

	var apiErr *tmdb.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *tmdb.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
	if apiErr.StatusMessage != "Invalid API key" {
		t.Errorf("StatusMessage = %q, want %q", apiErr.StatusMessage, "Invalid API key")
	}
}
