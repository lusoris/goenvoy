package tmdb_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/video/tmdb"
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
	return tmdb.New("test-token", metadata.WithBaseURL(srv.URL))
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
	got, err := c.DiscoverMovies(context.Background(), "en-US", 1, url.Values{"sort_by": {"popularity.desc"}})
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

	c := tmdb.New("bad-token", metadata.WithBaseURL(srv.URL))
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

func TestSearchPeople(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.PersonResult]{
		Page: 1, TotalPages: 1, TotalResults: 1,
		Results: []tmdb.PersonResult{{ID: 2, Name: "Mark Hamill"}},
	}

	srv := newTestServer(t, "/search/person?query=Mark+Hamill&language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.SearchPeople(context.Background(), "Mark Hamill", "en-US", 1)
	if err != nil {
		t.Fatalf("SearchPeople: %v", err)
	}
	if got.Results[0].Name != "Mark Hamill" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetMovieRecommendations(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.MovieResult]{
		Page: 1, TotalResults: 1,
		Results: []tmdb.MovieResult{{ID: 12, Title: "Recommended"}},
	}

	srv := newTestServer(t, "/movie/11/recommendations?language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMovieRecommendations(context.Background(), 11, "en-US", 1)
	if err != nil {
		t.Fatalf("GetMovieRecommendations: %v", err)
	}
	if got.Results[0].Title != "Recommended" {
		t.Errorf("Title = %q", got.Results[0].Title)
	}
}

func TestGetMovieSimilar(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.MovieResult]{
		Page: 1, TotalResults: 1,
		Results: []tmdb.MovieResult{{ID: 13, Title: "Similar"}},
	}

	srv := newTestServer(t, "/movie/11/similar?language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMovieSimilar(context.Background(), 11, "en-US", 1)
	if err != nil {
		t.Fatalf("GetMovieSimilar: %v", err)
	}
	if got.Results[0].Title != "Similar" {
		t.Errorf("Title = %q", got.Results[0].Title)
	}
}

func TestGetTVCredits(t *testing.T) {
	t.Parallel()

	want := tmdb.Credits{
		ID:   1396,
		Cast: []tmdb.CastMember{{ID: 17419, Name: "Bryan Cranston", Character: "Walter White"}},
	}

	srv := newTestServer(t, "/tv/1396/credits?language=en-US", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetTVCredits(context.Background(), 1396, "en-US")
	if err != nil {
		t.Fatalf("GetTVCredits: %v", err)
	}
	if got.Cast[0].Name != "Bryan Cranston" {
		t.Errorf("Name = %q", got.Cast[0].Name)
	}
}

func TestGetTVImages(t *testing.T) {
	t.Parallel()

	want := tmdb.Images{
		ID:      1396,
		Posters: []tmdb.ImageItem{{FilePath: "/poster.jpg", Width: 500, Height: 750}},
	}

	srv := newTestServer(t, "/tv/1396/images", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetTVImages(context.Background(), 1396)
	if err != nil {
		t.Fatalf("GetTVImages: %v", err)
	}
	if len(got.Posters) != 1 {
		t.Errorf("Posters len = %d", len(got.Posters))
	}
}

func TestGetTVExternalIDs(t *testing.T) {
	t.Parallel()

	want := tmdb.ExternalIDs{ID: 1396, IMDbID: "tt0903747", TVDbID: 81189}

	srv := newTestServer(t, "/tv/1396/external_ids", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetTVExternalIDs(context.Background(), 1396)
	if err != nil {
		t.Fatalf("GetTVExternalIDs: %v", err)
	}
	if got.IMDbID != "tt0903747" {
		t.Errorf("IMDbID = %q", got.IMDbID)
	}
}

func TestGetTVRecommendations(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.TVResult]{
		Page: 1, TotalResults: 1,
		Results: []tmdb.TVResult{{ID: 1, Name: "Recommended Show"}},
	}

	srv := newTestServer(t, "/tv/1396/recommendations?language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetTVRecommendations(context.Background(), 1396, "en-US", 1)
	if err != nil {
		t.Fatalf("GetTVRecommendations: %v", err)
	}
	if got.Results[0].Name != "Recommended Show" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetTVSimilar(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.TVResult]{
		Page: 1, TotalResults: 1,
		Results: []tmdb.TVResult{{ID: 2, Name: "Similar Show"}},
	}

	srv := newTestServer(t, "/tv/1396/similar?language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetTVSimilar(context.Background(), 1396, "en-US", 1)
	if err != nil {
		t.Fatalf("GetTVSimilar: %v", err)
	}
	if got.Results[0].Name != "Similar Show" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetPersonExternalIDs(t *testing.T) {
	t.Parallel()

	want := tmdb.ExternalIDs{ID: 2, IMDbID: "nm0372176"}

	srv := newTestServer(t, "/person/2/external_ids", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetPersonExternalIDs(context.Background(), 2)
	if err != nil {
		t.Fatalf("GetPersonExternalIDs: %v", err)
	}
	if got.IMDbID != "nm0372176" {
		t.Errorf("IMDbID = %q", got.IMDbID)
	}
}

func TestDiscoverTV(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.TVResult]{
		Page: 1, TotalResults: 1,
		Results: []tmdb.TVResult{{ID: 100, Name: "Disco Show"}},
	}

	srv := newTestServer(t, "/discover/tv?language=en-US&page=1&sort_by=popularity.desc", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.DiscoverTV(context.Background(), "en-US", 1, url.Values{"sort_by": {"popularity.desc"}})
	if err != nil {
		t.Fatalf("DiscoverTV: %v", err)
	}
	if got.Results[0].Name != "Disco Show" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetGenresTV(t *testing.T) {
	t.Parallel()

	want := struct {
		Genres []tmdb.Genre `json:"genres"`
	}{Genres: []tmdb.Genre{{ID: 18, Name: "Drama"}, {ID: 80, Name: "Crime"}}}

	srv := newTestServer(t, "/genre/tv/list?language=en-US", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetGenresTV(context.Background(), "en-US")
	if err != nil {
		t.Fatalf("GetGenresTV: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}

func TestGetTopRatedMovies(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.MovieResult]{
		Page: 1, TotalResults: 1,
		Results: []tmdb.MovieResult{{ID: 278, Title: "The Shawshank Redemption"}},
	}

	srv := newTestServer(t, "/movie/top_rated?language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetTopRatedMovies(context.Background(), "en-US", 1)
	if err != nil {
		t.Fatalf("GetTopRatedMovies: %v", err)
	}
	if got.Results[0].Title != "The Shawshank Redemption" {
		t.Errorf("Title = %q", got.Results[0].Title)
	}
}

func TestGetNowPlayingMovies(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.MovieResult]{
		Page: 1, TotalResults: 1,
		Results: []tmdb.MovieResult{{ID: 500, Title: "In Theaters Now"}},
	}

	srv := newTestServer(t, "/movie/now_playing?language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetNowPlayingMovies(context.Background(), "en-US", 1)
	if err != nil {
		t.Fatalf("GetNowPlayingMovies: %v", err)
	}
	if got.Results[0].Title != "In Theaters Now" {
		t.Errorf("Title = %q", got.Results[0].Title)
	}
}

func TestGetUpcomingMovies(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.MovieResult]{
		Page: 1, TotalResults: 1,
		Results: []tmdb.MovieResult{{ID: 600, Title: "Coming Soon"}},
	}

	srv := newTestServer(t, "/movie/upcoming?language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetUpcomingMovies(context.Background(), "en-US", 1)
	if err != nil {
		t.Fatalf("GetUpcomingMovies: %v", err)
	}
	if got.Results[0].Title != "Coming Soon" {
		t.Errorf("Title = %q", got.Results[0].Title)
	}
}

func TestGetPopularTV(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.TVResult]{
		Page: 1, TotalResults: 1,
		Results: []tmdb.TVResult{{ID: 700, Name: "Popular Show"}},
	}

	srv := newTestServer(t, "/tv/popular?language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetPopularTV(context.Background(), "en-US", 1)
	if err != nil {
		t.Fatalf("GetPopularTV: %v", err)
	}
	if got.Results[0].Name != "Popular Show" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetTopRatedTV(t *testing.T) {
	t.Parallel()

	want := tmdb.PaginatedResult[tmdb.TVResult]{
		Page: 1, TotalResults: 1,
		Results: []tmdb.TVResult{{ID: 800, Name: "Top Show"}},
	}

	srv := newTestServer(t, "/tv/top_rated?language=en-US&page=1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetTopRatedTV(context.Background(), "en-US", 1)
	if err != nil {
		t.Fatalf("GetTopRatedTV: %v", err)
	}
	if got.Results[0].Name != "Top Show" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

// newMethodTestServer creates a test server that validates HTTP method, auth, and path.
func newMethodTestServer(t *testing.T, wantMethod, wantPath string, response any) *httptest.Server {
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
		}
	}))
}

// Movie extras tests.

func TestGetMovieVideos(t *testing.T) {
	t.Parallel()
	want := tmdb.VideosResponse{ID: 550, Results: []tmdb.Video{{ID: "v1", Key: "abc", Name: "Trailer", Site: "YouTube"}}}
	srv := newTestServer(t, "/movie/550/videos?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieVideos(context.Background(), 550, "en-US")
	if err != nil {
		t.Fatalf("GetMovieVideos: %v", err)
	}
	if got.Results[0].Name != "Trailer" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetMovieKeywords(t *testing.T) {
	t.Parallel()
	want := tmdb.KeywordsResponse{ID: 550, Keywords: []tmdb.Keyword{{ID: 1, Name: "fight"}}}
	srv := newTestServer(t, "/movie/550/keywords", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieKeywords(context.Background(), 550)
	if err != nil {
		t.Fatalf("GetMovieKeywords: %v", err)
	}
	if got.Keywords[0].Name != "fight" {
		t.Errorf("Name = %q", got.Keywords[0].Name)
	}
}

func TestGetMovieReviews(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.Review]{Page: 1, Results: []tmdb.Review{{ID: "r1", Author: "alice"}}}
	srv := newTestServer(t, "/movie/550/reviews?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieReviews(context.Background(), 550, "en-US", 1)
	if err != nil {
		t.Fatalf("GetMovieReviews: %v", err)
	}
	if got.Results[0].Author != "alice" {
		t.Errorf("Author = %q", got.Results[0].Author)
	}
}

func TestGetMovieReleaseDates(t *testing.T) {
	t.Parallel()
	want := tmdb.ReleaseDatesResponse{ID: 550, Results: []tmdb.ReleaseDateCountry{{ISO31661: "US"}}}
	srv := newTestServer(t, "/movie/550/release_dates", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieReleaseDates(context.Background(), 550)
	if err != nil {
		t.Fatalf("GetMovieReleaseDates: %v", err)
	}
	if got.Results[0].ISO31661 != "US" {
		t.Errorf("ISO = %q", got.Results[0].ISO31661)
	}
}

func TestGetMovieWatchProviders(t *testing.T) {
	t.Parallel()
	want := tmdb.WatchProvidersResponse{ID: 550, Results: map[string]tmdb.WatchProviderCountry{"US": {Link: "https://example.com"}}}
	srv := newTestServer(t, "/movie/550/watch/providers", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieWatchProviders(context.Background(), 550)
	if err != nil {
		t.Fatalf("GetMovieWatchProviders: %v", err)
	}
	if got.Results["US"].Link != "https://example.com" {
		t.Errorf("Link = %q", got.Results["US"].Link)
	}
}

func TestGetMovieAlternativeTitles(t *testing.T) {
	t.Parallel()
	want := tmdb.AlternativeTitlesResponse{ID: 550, Titles: []tmdb.AlternativeTitle{{Title: "El Club de la Pelea"}}}
	srv := newTestServer(t, "/movie/550/alternative_titles?country=MX", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieAlternativeTitles(context.Background(), 550, "MX")
	if err != nil {
		t.Fatalf("GetMovieAlternativeTitles: %v", err)
	}
	if got.Titles[0].Title != "El Club de la Pelea" {
		t.Errorf("Title = %q", got.Titles[0].Title)
	}
}

func TestGetMovieTranslations(t *testing.T) {
	t.Parallel()
	want := tmdb.TranslationsResponse{ID: 550, Translations: []tmdb.Translation{{EnglishName: "French"}}}
	srv := newTestServer(t, "/movie/550/translations", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieTranslations(context.Background(), 550)
	if err != nil {
		t.Fatalf("GetMovieTranslations: %v", err)
	}
	if got.Translations[0].EnglishName != "French" {
		t.Errorf("Name = %q", got.Translations[0].EnglishName)
	}
}

func TestGetMovieLists(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.ListSummary]{Page: 1, Results: []tmdb.ListSummary{{ID: 1, Name: "My List"}}}
	srv := newTestServer(t, "/movie/550/lists?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieLists(context.Background(), 550, "en-US", 1)
	if err != nil {
		t.Fatalf("GetMovieLists: %v", err)
	}
	if got.Results[0].Name != "My List" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetMovieAccountStates(t *testing.T) {
	t.Parallel()
	want := tmdb.AccountStates{ID: 550, Favorite: true}
	srv := newTestServer(t, "/movie/550/account_states", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieAccountStates(context.Background(), 550)
	if err != nil {
		t.Fatalf("GetMovieAccountStates: %v", err)
	}
	if !got.Favorite {
		t.Errorf("Favorite = false, want true")
	}
}

func TestRateMovie(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodPost, "/movie/550/rating", map[string]any{"status_code": 1})
	defer srv.Close()
	err := newClient(t, srv).RateMovie(context.Background(), 550, 8.5)
	if err != nil {
		t.Fatalf("RateMovie: %v", err)
	}
}

func TestDeleteMovieRating(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodDelete, "/movie/550/rating", nil)
	defer srv.Close()
	err := newClient(t, srv).DeleteMovieRating(context.Background(), 550)
	if err != nil {
		t.Fatalf("DeleteMovieRating: %v", err)
	}
}

// TV extras tests.

func TestGetTVVideos(t *testing.T) {
	t.Parallel()
	want := tmdb.VideosResponse{ID: 100, Results: []tmdb.Video{{ID: "v1", Name: "Teaser"}}}
	srv := newTestServer(t, "/tv/100/videos?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVVideos(context.Background(), 100, "en-US")
	if err != nil {
		t.Fatalf("GetTVVideos: %v", err)
	}
	if got.Results[0].Name != "Teaser" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetTVKeywords(t *testing.T) {
	t.Parallel()
	want := tmdb.KeywordsResponse{ID: 100, Results: []tmdb.Keyword{{ID: 2, Name: "drama"}}}
	srv := newTestServer(t, "/tv/100/keywords", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVKeywords(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetTVKeywords: %v", err)
	}
	if got.Results[0].Name != "drama" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetTVReviews(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.Review]{Page: 1, Results: []tmdb.Review{{ID: "r1", Author: "bob"}}}
	srv := newTestServer(t, "/tv/100/reviews?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVReviews(context.Background(), 100, "en-US", 1)
	if err != nil {
		t.Fatalf("GetTVReviews: %v", err)
	}
	if got.Results[0].Author != "bob" {
		t.Errorf("Author = %q", got.Results[0].Author)
	}
}

func TestGetTVWatchProviders(t *testing.T) {
	t.Parallel()
	want := tmdb.WatchProvidersResponse{ID: 100, Results: map[string]tmdb.WatchProviderCountry{"US": {Link: "https://tv.com"}}}
	srv := newTestServer(t, "/tv/100/watch/providers", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVWatchProviders(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetTVWatchProviders: %v", err)
	}
	if got.Results["US"].Link != "https://tv.com" {
		t.Errorf("Link = %q", got.Results["US"].Link)
	}
}

func TestGetTVAlternativeTitles(t *testing.T) {
	t.Parallel()
	want := tmdb.AlternativeTitlesResponse{ID: 100, Results: []tmdb.AlternativeTitle{{Title: "Alt"}}}
	srv := newTestServer(t, "/tv/100/alternative_titles", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVAlternativeTitles(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetTVAlternativeTitles: %v", err)
	}
	if got.Results[0].Title != "Alt" {
		t.Errorf("Title = %q", got.Results[0].Title)
	}
}

func TestGetTVTranslations(t *testing.T) {
	t.Parallel()
	want := tmdb.TranslationsResponse{ID: 100, Translations: []tmdb.Translation{{EnglishName: "Spanish"}}}
	srv := newTestServer(t, "/tv/100/translations", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVTranslations(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetTVTranslations: %v", err)
	}
	if got.Translations[0].EnglishName != "Spanish" {
		t.Errorf("Name = %q", got.Translations[0].EnglishName)
	}
}

func TestGetTVContentRatings(t *testing.T) {
	t.Parallel()
	want := tmdb.ContentRatingsResponse{ID: 100, Results: []tmdb.ContentRating{{ISO31661: "US", Rating: "TV-MA"}}}
	srv := newTestServer(t, "/tv/100/content_ratings", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVContentRatings(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetTVContentRatings: %v", err)
	}
	if got.Results[0].Rating != "TV-MA" {
		t.Errorf("Rating = %q", got.Results[0].Rating)
	}
}

func TestGetTVEpisodeGroups(t *testing.T) {
	t.Parallel()
	want := tmdb.EpisodeGroupsResponse{ID: 100, Results: []tmdb.EpisodeGroup{{ID: "g1", Name: "Seasons"}}}
	srv := newTestServer(t, "/tv/100/episode_groups", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVEpisodeGroups(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetTVEpisodeGroups: %v", err)
	}
	if got.Results[0].Name != "Seasons" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetTVAggregateCredits(t *testing.T) {
	t.Parallel()
	want := tmdb.AggregateCredits{ID: 100, Cast: []tmdb.AggregateCastMember{{ID: 1, Name: "Jane"}}}
	srv := newTestServer(t, "/tv/100/aggregate_credits?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVAggregateCredits(context.Background(), 100, "en-US")
	if err != nil {
		t.Fatalf("GetTVAggregateCredits: %v", err)
	}
	if got.Cast[0].Name != "Jane" {
		t.Errorf("Name = %q", got.Cast[0].Name)
	}
}

func TestGetTVAccountStates(t *testing.T) {
	t.Parallel()
	want := tmdb.AccountStates{ID: 100, Watchlist: true}
	srv := newTestServer(t, "/tv/100/account_states", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVAccountStates(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetTVAccountStates: %v", err)
	}
	if !got.Watchlist {
		t.Errorf("Watchlist = false")
	}
}

func TestRateTV(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodPost, "/tv/100/rating", map[string]any{"status_code": 1})
	defer srv.Close()
	err := newClient(t, srv).RateTV(context.Background(), 100, 9.0)
	if err != nil {
		t.Fatalf("RateTV: %v", err)
	}
}

func TestDeleteTVRating(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodDelete, "/tv/100/rating", nil)
	defer srv.Close()
	err := newClient(t, srv).DeleteTVRating(context.Background(), 100)
	if err != nil {
		t.Fatalf("DeleteTVRating: %v", err)
	}
}

// TV Season extras tests.

func TestGetTVSeasonCredits(t *testing.T) {
	t.Parallel()
	want := tmdb.Credits{ID: 1, Cast: []tmdb.CastMember{{ID: 1, Name: "Actor"}}}
	srv := newTestServer(t, "/tv/100/season/1/credits?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVSeasonCredits(context.Background(), 100, 1, "en-US")
	if err != nil {
		t.Fatalf("GetTVSeasonCredits: %v", err)
	}
	if got.Cast[0].Name != "Actor" {
		t.Errorf("Name = %q", got.Cast[0].Name)
	}
}

func TestGetTVSeasonImages(t *testing.T) {
	t.Parallel()
	want := tmdb.Images{ID: 1, Posters: []tmdb.ImageItem{{FilePath: "/poster.jpg"}}}
	srv := newTestServer(t, "/tv/100/season/1/images", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVSeasonImages(context.Background(), 100, 1)
	if err != nil {
		t.Fatalf("GetTVSeasonImages: %v", err)
	}
	if got.Posters[0].FilePath != "/poster.jpg" {
		t.Errorf("FilePath = %q", got.Posters[0].FilePath)
	}
}

func TestGetTVSeasonVideos(t *testing.T) {
	t.Parallel()
	want := tmdb.VideosResponse{ID: 1, Results: []tmdb.Video{{Name: "Recap"}}}
	srv := newTestServer(t, "/tv/100/season/1/videos?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVSeasonVideos(context.Background(), 100, 1, "en-US")
	if err != nil {
		t.Fatalf("GetTVSeasonVideos: %v", err)
	}
	if got.Results[0].Name != "Recap" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetTVSeasonExternalIDs(t *testing.T) {
	t.Parallel()
	want := tmdb.ExternalIDs{ID: 1, TVDbID: 123}
	srv := newTestServer(t, "/tv/100/season/1/external_ids", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVSeasonExternalIDs(context.Background(), 100, 1)
	if err != nil {
		t.Fatalf("GetTVSeasonExternalIDs: %v", err)
	}
	if got.TVDbID != 123 {
		t.Errorf("TVDbID = %d", got.TVDbID)
	}
}

func TestGetTVSeasonTranslations(t *testing.T) {
	t.Parallel()
	want := tmdb.TranslationsResponse{ID: 1, Translations: []tmdb.Translation{{EnglishName: "German"}}}
	srv := newTestServer(t, "/tv/100/season/1/translations", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVSeasonTranslations(context.Background(), 100, 1)
	if err != nil {
		t.Fatalf("GetTVSeasonTranslations: %v", err)
	}
	if got.Translations[0].EnglishName != "German" {
		t.Errorf("Name = %q", got.Translations[0].EnglishName)
	}
}

func TestGetTVSeasonWatchProviders(t *testing.T) {
	t.Parallel()
	want := tmdb.WatchProvidersResponse{ID: 1, Results: map[string]tmdb.WatchProviderCountry{"US": {Link: "https://s.com"}}}
	srv := newTestServer(t, "/tv/100/season/1/watch/providers", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVSeasonWatchProviders(context.Background(), 100, 1)
	if err != nil {
		t.Fatalf("GetTVSeasonWatchProviders: %v", err)
	}
	if got.Results["US"].Link != "https://s.com" {
		t.Errorf("Link = %q", got.Results["US"].Link)
	}
}

// TV Episode tests.

func TestGetTVEpisode(t *testing.T) {
	t.Parallel()
	want := tmdb.EpisodeDetails{ID: 1, Name: "Pilot", SeasonNumber: 1, EpisodeNumber: 1}
	srv := newTestServer(t, "/tv/100/season/1/episode/1?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVEpisode(context.Background(), 100, 1, 1, "en-US")
	if err != nil {
		t.Fatalf("GetTVEpisode: %v", err)
	}
	if got.Name != "Pilot" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestGetTVEpisodeCredits(t *testing.T) {
	t.Parallel()
	want := tmdb.Credits{ID: 1, Cast: []tmdb.CastMember{{ID: 5, Name: "Star"}}}
	srv := newTestServer(t, "/tv/100/season/1/episode/1/credits?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVEpisodeCredits(context.Background(), 100, 1, 1, "en-US")
	if err != nil {
		t.Fatalf("GetTVEpisodeCredits: %v", err)
	}
	if got.Cast[0].Name != "Star" {
		t.Errorf("Name = %q", got.Cast[0].Name)
	}
}

func TestGetTVEpisodeImages(t *testing.T) {
	t.Parallel()
	want := tmdb.Images{ID: 1, Backdrops: []tmdb.ImageItem{{FilePath: "/still.jpg"}}}
	srv := newTestServer(t, "/tv/100/season/1/episode/1/images", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVEpisodeImages(context.Background(), 100, 1, 1)
	if err != nil {
		t.Fatalf("GetTVEpisodeImages: %v", err)
	}
	if got.Backdrops[0].FilePath != "/still.jpg" {
		t.Errorf("FilePath = %q", got.Backdrops[0].FilePath)
	}
}

func TestGetTVEpisodeVideos(t *testing.T) {
	t.Parallel()
	want := tmdb.VideosResponse{ID: 1, Results: []tmdb.Video{{Name: "Preview"}}}
	srv := newTestServer(t, "/tv/100/season/1/episode/1/videos?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVEpisodeVideos(context.Background(), 100, 1, 1, "en-US")
	if err != nil {
		t.Fatalf("GetTVEpisodeVideos: %v", err)
	}
	if got.Results[0].Name != "Preview" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetTVEpisodeExternalIDs(t *testing.T) {
	t.Parallel()
	want := tmdb.ExternalIDs{ID: 1, IMDbID: "tt000001"}
	srv := newTestServer(t, "/tv/100/season/1/episode/1/external_ids", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVEpisodeExternalIDs(context.Background(), 100, 1, 1)
	if err != nil {
		t.Fatalf("GetTVEpisodeExternalIDs: %v", err)
	}
	if got.IMDbID != "tt000001" {
		t.Errorf("IMDbID = %q", got.IMDbID)
	}
}

func TestGetTVEpisodeTranslations(t *testing.T) {
	t.Parallel()
	want := tmdb.TranslationsResponse{ID: 1, Translations: []tmdb.Translation{{EnglishName: "Italian"}}}
	srv := newTestServer(t, "/tv/100/season/1/episode/1/translations", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVEpisodeTranslations(context.Background(), 100, 1, 1)
	if err != nil {
		t.Fatalf("GetTVEpisodeTranslations: %v", err)
	}
	if got.Translations[0].EnglishName != "Italian" {
		t.Errorf("Name = %q", got.Translations[0].EnglishName)
	}
}

func TestGetTVEpisodeAccountStates(t *testing.T) {
	t.Parallel()
	want := tmdb.AccountStates{ID: 1, Favorite: true}
	srv := newTestServer(t, "/tv/100/season/1/episode/1/account_states", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVEpisodeAccountStates(context.Background(), 100, 1, 1)
	if err != nil {
		t.Fatalf("GetTVEpisodeAccountStates: %v", err)
	}
	if !got.Favorite {
		t.Errorf("Favorite = false")
	}
}

func TestRateTVEpisode(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodPost, "/tv/100/season/1/episode/1/rating", map[string]any{"status_code": 1})
	defer srv.Close()
	err := newClient(t, srv).RateTVEpisode(context.Background(), 100, 1, 1, 7.5)
	if err != nil {
		t.Fatalf("RateTVEpisode: %v", err)
	}
}

func TestDeleteTVEpisodeRating(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodDelete, "/tv/100/season/1/episode/1/rating", nil)
	defer srv.Close()
	err := newClient(t, srv).DeleteTVEpisodeRating(context.Background(), 100, 1, 1)
	if err != nil {
		t.Fatalf("DeleteTVEpisodeRating: %v", err)
	}
}

// Person extras tests.

func TestGetPersonMovieCredits(t *testing.T) {
	t.Parallel()
	want := tmdb.PersonCredits{ID: 287, Cast: []tmdb.PersonCastCredit{{ID: 550, Title: "Fight Club"}}}
	srv := newTestServer(t, "/person/287/movie_credits?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetPersonMovieCredits(context.Background(), 287, "en-US")
	if err != nil {
		t.Fatalf("GetPersonMovieCredits: %v", err)
	}
	if got.Cast[0].Title != "Fight Club" {
		t.Errorf("Title = %q", got.Cast[0].Title)
	}
}

func TestGetPersonTVCredits(t *testing.T) {
	t.Parallel()
	want := tmdb.PersonCredits{ID: 287, Cast: []tmdb.PersonCastCredit{{ID: 1, Name: "Friends"}}}
	srv := newTestServer(t, "/person/287/tv_credits?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetPersonTVCredits(context.Background(), 287, "en-US")
	if err != nil {
		t.Fatalf("GetPersonTVCredits: %v", err)
	}
	if got.Cast[0].Name != "Friends" {
		t.Errorf("Name = %q", got.Cast[0].Name)
	}
}

func TestGetPersonCombinedCredits(t *testing.T) {
	t.Parallel()
	want := tmdb.PersonCredits{ID: 287, Cast: []tmdb.PersonCastCredit{{ID: 550, MediaType: "movie"}}}
	srv := newTestServer(t, "/person/287/combined_credits?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetPersonCombinedCredits(context.Background(), 287, "en-US")
	if err != nil {
		t.Fatalf("GetPersonCombinedCredits: %v", err)
	}
	if got.Cast[0].MediaType != "movie" {
		t.Errorf("MediaType = %q", got.Cast[0].MediaType)
	}
}

func TestGetPersonImages(t *testing.T) {
	t.Parallel()
	want := tmdb.PersonImages{ID: 287, Profiles: []tmdb.ImageItem{{FilePath: "/profile.jpg"}}}
	srv := newTestServer(t, "/person/287/images", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetPersonImages(context.Background(), 287)
	if err != nil {
		t.Fatalf("GetPersonImages: %v", err)
	}
	if got.Profiles[0].FilePath != "/profile.jpg" {
		t.Errorf("FilePath = %q", got.Profiles[0].FilePath)
	}
}

func TestGetPersonTaggedImages(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.TaggedImage]{Page: 1, Results: []tmdb.TaggedImage{{MediaType: "movie"}}}
	srv := newTestServer(t, "/person/287/tagged_images?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetPersonTaggedImages(context.Background(), 287, "en-US", 1)
	if err != nil {
		t.Fatalf("GetPersonTaggedImages: %v", err)
	}
	if got.Results[0].MediaType != "movie" {
		t.Errorf("MediaType = %q", got.Results[0].MediaType)
	}
}

// Search extras tests.

func TestSearchKeywords(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.Keyword]{Page: 1, Results: []tmdb.Keyword{{ID: 1, Name: "fight"}}}
	srv := newTestServer(t, "/search/keyword?query=fight&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).SearchKeywords(context.Background(), "fight", 1)
	if err != nil {
		t.Fatalf("SearchKeywords: %v", err)
	}
	if got.Results[0].Name != "fight" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestSearchCollections(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.CollectionResult]{Page: 1, Results: []tmdb.CollectionResult{{ID: 1, Name: "Star Wars"}}}
	srv := newTestServer(t, "/search/collection?query=Star+Wars&language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).SearchCollections(context.Background(), "Star Wars", "en-US", 1)
	if err != nil {
		t.Fatalf("SearchCollections: %v", err)
	}
	if got.Results[0].Name != "Star Wars" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestSearchCompanies(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.CompanyResult]{Page: 1, Results: []tmdb.CompanyResult{{ID: 1, Name: "Disney"}}}
	srv := newTestServer(t, "/search/company?query=Disney&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).SearchCompanies(context.Background(), "Disney", 1)
	if err != nil {
		t.Fatalf("SearchCompanies: %v", err)
	}
	if got.Results[0].Name != "Disney" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

// Collection tests.

func TestGetCollection(t *testing.T) {
	t.Parallel()
	want := tmdb.CollectionDetails{ID: 10, Name: "Star Wars Collection", Parts: []tmdb.CollectionPart{{ID: 11, Title: "Star Wars"}}}
	srv := newTestServer(t, "/collection/10?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetCollection(context.Background(), 10, "en-US")
	if err != nil {
		t.Fatalf("GetCollection: %v", err)
	}
	if got.Name != "Star Wars Collection" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestGetCollectionImages(t *testing.T) {
	t.Parallel()
	want := tmdb.Images{ID: 10, Posters: []tmdb.ImageItem{{FilePath: "/coll.jpg"}}}
	srv := newTestServer(t, "/collection/10/images", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetCollectionImages(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetCollectionImages: %v", err)
	}
	if got.Posters[0].FilePath != "/coll.jpg" {
		t.Errorf("FilePath = %q", got.Posters[0].FilePath)
	}
}

func TestGetCollectionTranslations(t *testing.T) {
	t.Parallel()
	want := tmdb.TranslationsResponse{ID: 10, Translations: []tmdb.Translation{{EnglishName: "English"}}}
	srv := newTestServer(t, "/collection/10/translations", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetCollectionTranslations(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetCollectionTranslations: %v", err)
	}
	if got.Translations[0].EnglishName != "English" {
		t.Errorf("Name = %q", got.Translations[0].EnglishName)
	}
}

// Account tests.

func TestGetAccountDetails(t *testing.T) {
	t.Parallel()
	want := tmdb.AccountDetails{ID: 1, Username: "testuser"}
	srv := newTestServer(t, "/account", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetAccountDetails(context.Background())
	if err != nil {
		t.Fatalf("GetAccountDetails: %v", err)
	}
	if got.Username != "testuser" {
		t.Errorf("Username = %q", got.Username)
	}
}

func TestGetFavoriteMovies(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.MovieResult]{Page: 1, Results: []tmdb.MovieResult{{ID: 1, Title: "Fav Movie"}}}
	srv := newTestServer(t, "/account/1/favorite/movies?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetFavoriteMovies(context.Background(), 1, "en-US", 1)
	if err != nil {
		t.Fatalf("GetFavoriteMovies: %v", err)
	}
	if got.Results[0].Title != "Fav Movie" {
		t.Errorf("Title = %q", got.Results[0].Title)
	}
}

func TestGetFavoriteTV(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.TVResult]{Page: 1, Results: []tmdb.TVResult{{ID: 1, Name: "Fav Show"}}}
	srv := newTestServer(t, "/account/1/favorite/tv?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetFavoriteTV(context.Background(), 1, "en-US", 1)
	if err != nil {
		t.Fatalf("GetFavoriteTV: %v", err)
	}
	if got.Results[0].Name != "Fav Show" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestAddFavorite(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodPost, "/account/1/favorite", map[string]any{"status_code": 1})
	defer srv.Close()
	err := newClient(t, srv).AddFavorite(context.Background(), 1, "movie", 550, true)
	if err != nil {
		t.Fatalf("AddFavorite: %v", err)
	}
}

func TestGetWatchlistMovies(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.MovieResult]{Page: 1, Results: []tmdb.MovieResult{{ID: 1, Title: "WL Movie"}}}
	srv := newTestServer(t, "/account/1/watchlist/movies?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetWatchlistMovies(context.Background(), 1, "en-US", 1)
	if err != nil {
		t.Fatalf("GetWatchlistMovies: %v", err)
	}
	if got.Results[0].Title != "WL Movie" {
		t.Errorf("Title = %q", got.Results[0].Title)
	}
}

func TestGetWatchlistTV(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.TVResult]{Page: 1, Results: []tmdb.TVResult{{ID: 1, Name: "WL Show"}}}
	srv := newTestServer(t, "/account/1/watchlist/tv?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetWatchlistTV(context.Background(), 1, "en-US", 1)
	if err != nil {
		t.Fatalf("GetWatchlistTV: %v", err)
	}
	if got.Results[0].Name != "WL Show" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestAddToWatchlist(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodPost, "/account/1/watchlist", map[string]any{"status_code": 1})
	defer srv.Close()
	err := newClient(t, srv).AddToWatchlist(context.Background(), 1, "movie", 550, true)
	if err != nil {
		t.Fatalf("AddToWatchlist: %v", err)
	}
}

func TestGetRatedMovies(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.RatedMovie]{Page: 1, Results: []tmdb.RatedMovie{{MovieResult: tmdb.MovieResult{ID: 1, Title: "Rated"}, Rating: 8}}}
	srv := newTestServer(t, "/account/1/rated/movies?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetRatedMovies(context.Background(), 1, "en-US", 1)
	if err != nil {
		t.Fatalf("GetRatedMovies: %v", err)
	}
	if got.Results[0].Rating != 8 {
		t.Errorf("Rating = %v", got.Results[0].Rating)
	}
}

func TestGetRatedTV(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.RatedTV]{Page: 1, Results: []tmdb.RatedTV{{TVResult: tmdb.TVResult{ID: 1, Name: "Rated Show"}, Rating: 9}}}
	srv := newTestServer(t, "/account/1/rated/tv?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetRatedTV(context.Background(), 1, "en-US", 1)
	if err != nil {
		t.Fatalf("GetRatedTV: %v", err)
	}
	if got.Results[0].Rating != 9 {
		t.Errorf("Rating = %v", got.Results[0].Rating)
	}
}

func TestGetRatedTVEpisodes(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.RatedEpisode]{Page: 1, Results: []tmdb.RatedEpisode{{ID: 1, Name: "Ep", Rating: 7}}}
	srv := newTestServer(t, "/account/1/rated/tv/episodes?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetRatedTVEpisodes(context.Background(), 1, "en-US", 1)
	if err != nil {
		t.Fatalf("GetRatedTVEpisodes: %v", err)
	}
	if got.Results[0].Rating != 7 {
		t.Errorf("Rating = %v", got.Results[0].Rating)
	}
}

// List tests.

func TestGetList(t *testing.T) {
	t.Parallel()
	want := tmdb.ListDetails{ID: "1", Name: "My List", ItemCount: 5}
	srv := newTestServer(t, "/list/1?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetList(context.Background(), 1, "en-US")
	if err != nil {
		t.Fatalf("GetList: %v", err)
	}
	if got.Name != "My List" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestCreateList(t *testing.T) {
	t.Parallel()
	want := tmdb.CreateListResponse{StatusCode: 1, Success: true, ListID: 42}
	srv := newMethodTestServer(t, http.MethodPost, "/list", want)
	defer srv.Close()
	got, err := newClient(t, srv).CreateList(context.Background(), "Test", "A list", "en")
	if err != nil {
		t.Fatalf("CreateList: %v", err)
	}
	if got.ListID != 42 {
		t.Errorf("ListID = %d", got.ListID)
	}
}

func TestAddToList(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodPost, "/list/1/add_item", map[string]any{"status_code": 12})
	defer srv.Close()
	err := newClient(t, srv).AddToList(context.Background(), 1, 550)
	if err != nil {
		t.Fatalf("AddToList: %v", err)
	}
}

func TestRemoveFromList(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodPost, "/list/1/remove_item", map[string]any{"status_code": 13})
	defer srv.Close()
	err := newClient(t, srv).RemoveFromList(context.Background(), 1, 550)
	if err != nil {
		t.Fatalf("RemoveFromList: %v", err)
	}
}

func TestClearList(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodPost, "/list/1/clear?confirm=true", map[string]any{"status_code": 12})
	defer srv.Close()
	err := newClient(t, srv).ClearList(context.Background(), 1, true)
	if err != nil {
		t.Fatalf("ClearList: %v", err)
	}
}

func TestDeleteList(t *testing.T) {
	t.Parallel()
	srv := newMethodTestServer(t, http.MethodDelete, "/list/1", nil)
	defer srv.Close()
	err := newClient(t, srv).DeleteList(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeleteList: %v", err)
	}
}

func TestCheckItemStatus(t *testing.T) {
	t.Parallel()
	want := tmdb.ItemStatus{ID: "1", ItemPresent: true}
	srv := newTestServer(t, "/list/1/item_status?movie_id=550", want)
	defer srv.Close()
	got, err := newClient(t, srv).CheckItemStatus(context.Background(), 1, 550)
	if err != nil {
		t.Fatalf("CheckItemStatus: %v", err)
	}
	if !got.ItemPresent {
		t.Errorf("ItemPresent = false")
	}
}

// Certification tests.

func TestGetMovieCertifications(t *testing.T) {
	t.Parallel()
	want := tmdb.CertificationsResponse{Certifications: map[string][]tmdb.Certification{"US": {{Certification: "PG-13"}}}}
	srv := newTestServer(t, "/certification/movie/list", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieCertifications(context.Background())
	if err != nil {
		t.Fatalf("GetMovieCertifications: %v", err)
	}
	if got.Certifications["US"][0].Certification != "PG-13" {
		t.Errorf("Certification = %q", got.Certifications["US"][0].Certification)
	}
}

func TestGetTVCertifications(t *testing.T) {
	t.Parallel()
	want := tmdb.CertificationsResponse{Certifications: map[string][]tmdb.Certification{"US": {{Certification: "TV-14"}}}}
	srv := newTestServer(t, "/certification/tv/list", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVCertifications(context.Background())
	if err != nil {
		t.Fatalf("GetTVCertifications: %v", err)
	}
	if got.Certifications["US"][0].Certification != "TV-14" {
		t.Errorf("Certification = %q", got.Certifications["US"][0].Certification)
	}
}

// Watch Provider tests.

func TestGetAvailableWatchProviderRegions(t *testing.T) {
	t.Parallel()
	want := tmdb.WatchProviderRegionsResponse{Results: []tmdb.WatchProviderRegion{{ISO31661: "US", EnglishName: "United States"}}}
	srv := newTestServer(t, "/watch/providers/regions?language=en-US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetAvailableWatchProviderRegions(context.Background(), "en-US")
	if err != nil {
		t.Fatalf("GetAvailableWatchProviderRegions: %v", err)
	}
	if got.Results[0].EnglishName != "United States" {
		t.Errorf("Name = %q", got.Results[0].EnglishName)
	}
}

func TestGetMovieWatchProviderList(t *testing.T) {
	t.Parallel()
	want := tmdb.WatchProviderListResponse{Results: []tmdb.WatchProviderListItem{{ProviderID: 8, ProviderName: "Netflix"}}}
	srv := newTestServer(t, "/watch/providers/movie?language=en-US&watch_region=US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieWatchProviderList(context.Background(), "en-US", "US")
	if err != nil {
		t.Fatalf("GetMovieWatchProviderList: %v", err)
	}
	if got.Results[0].ProviderName != "Netflix" {
		t.Errorf("Name = %q", got.Results[0].ProviderName)
	}
}

func TestGetTVWatchProviderList(t *testing.T) {
	t.Parallel()
	want := tmdb.WatchProviderListResponse{Results: []tmdb.WatchProviderListItem{{ProviderID: 9, ProviderName: "Hulu"}}}
	srv := newTestServer(t, "/watch/providers/tv?language=en-US&watch_region=US", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVWatchProviderList(context.Background(), "en-US", "US")
	if err != nil {
		t.Fatalf("GetTVWatchProviderList: %v", err)
	}
	if got.Results[0].ProviderName != "Hulu" {
		t.Errorf("Name = %q", got.Results[0].ProviderName)
	}
}

// Company & Keyword tests.

func TestGetCompany(t *testing.T) {
	t.Parallel()
	want := tmdb.CompanyDetails{ID: 1, Name: "Pixar", OriginCountry: "US"}
	srv := newTestServer(t, "/company/1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetCompany(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCompany: %v", err)
	}
	if got.Name != "Pixar" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestGetKeyword(t *testing.T) {
	t.Parallel()
	want := tmdb.Keyword{ID: 1, Name: "fight"}
	srv := newTestServer(t, "/keyword/1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetKeyword(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetKeyword: %v", err)
	}
	if got.Name != "fight" {
		t.Errorf("Name = %q", got.Name)
	}
}

// Changes tests.

func TestGetMovieChanges(t *testing.T) {
	t.Parallel()
	want := tmdb.ChangesResponse{Page: 1, Results: []tmdb.ChangeItem{{ID: 1}}}
	srv := newTestServer(t, "/movie/changes?end_date=2024-01-02&page=1&start_date=2024-01-01", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetMovieChanges(context.Background(), "2024-01-01", "2024-01-02", 1)
	if err != nil {
		t.Fatalf("GetMovieChanges: %v", err)
	}
	if got.Results[0].ID != 1 {
		t.Errorf("ID = %d", got.Results[0].ID)
	}
}

func TestGetTVChanges(t *testing.T) {
	t.Parallel()
	want := tmdb.ChangesResponse{Page: 1, Results: []tmdb.ChangeItem{{ID: 2}}}
	srv := newTestServer(t, "/tv/changes?end_date=2024-01-02&page=1&start_date=2024-01-01", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTVChanges(context.Background(), "2024-01-01", "2024-01-02", 1)
	if err != nil {
		t.Fatalf("GetTVChanges: %v", err)
	}
	if got.Results[0].ID != 2 {
		t.Errorf("ID = %d", got.Results[0].ID)
	}
}

func TestGetPersonChanges(t *testing.T) {
	t.Parallel()
	want := tmdb.ChangesResponse{Page: 1, Results: []tmdb.ChangeItem{{ID: 3}}}
	srv := newTestServer(t, "/person/changes?end_date=2024-01-02&page=1&start_date=2024-01-01", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetPersonChanges(context.Background(), "2024-01-01", "2024-01-02", 1)
	if err != nil {
		t.Fatalf("GetPersonChanges: %v", err)
	}
	if got.Results[0].ID != 3 {
		t.Errorf("ID = %d", got.Results[0].ID)
	}
}

// Misc tests.

func TestGetReview(t *testing.T) {
	t.Parallel()
	want := tmdb.Review{ID: "abc123", Author: "critic", Content: "Great movie!"}
	srv := newTestServer(t, "/review/abc123", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetReview(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("GetReview: %v", err)
	}
	if got.Author != "critic" {
		t.Errorf("Author = %q", got.Author)
	}
}

func TestGetNetwork(t *testing.T) {
	t.Parallel()
	want := tmdb.Network{ID: 213, Name: "Netflix"}
	srv := newTestServer(t, "/network/213", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetNetwork(context.Background(), 213)
	if err != nil {
		t.Fatalf("GetNetwork: %v", err)
	}
	if got.Name != "Netflix" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestGetOnTheAirTV(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.TVResult]{Page: 1, Results: []tmdb.TVResult{{ID: 1, Name: "On Air"}}}
	srv := newTestServer(t, "/tv/on_the_air?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetOnTheAirTV(context.Background(), "en-US", 1)
	if err != nil {
		t.Fatalf("GetOnTheAirTV: %v", err)
	}
	if got.Results[0].Name != "On Air" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetAiringTodayTV(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.TVResult]{Page: 1, Results: []tmdb.TVResult{{ID: 2, Name: "Today"}}}
	srv := newTestServer(t, "/tv/airing_today?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetAiringTodayTV(context.Background(), "en-US", 1)
	if err != nil {
		t.Fatalf("GetAiringTodayTV: %v", err)
	}
	if got.Results[0].Name != "Today" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetLatestMovie(t *testing.T) {
	t.Parallel()
	want := tmdb.MovieDetails{ID: 999, Title: "Latest Movie"}
	srv := newTestServer(t, "/movie/latest", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetLatestMovie(context.Background())
	if err != nil {
		t.Fatalf("GetLatestMovie: %v", err)
	}
	if got.Title != "Latest Movie" {
		t.Errorf("Title = %q", got.Title)
	}
}

func TestGetLatestTV(t *testing.T) {
	t.Parallel()
	want := tmdb.TVDetails{ID: 888, Name: "Latest TV"}
	srv := newTestServer(t, "/tv/latest", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetLatestTV(context.Background())
	if err != nil {
		t.Fatalf("GetLatestTV: %v", err)
	}
	if got.Name != "Latest TV" {
		t.Errorf("Name = %q", got.Name)
	}
}

func TestGetPopularPeople(t *testing.T) {
	t.Parallel()
	want := tmdb.PaginatedResult[tmdb.PersonResult]{Page: 1, Results: []tmdb.PersonResult{{ID: 1, Name: "Famous"}}}
	srv := newTestServer(t, "/person/popular?language=en-US&page=1", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetPopularPeople(context.Background(), "en-US", 1)
	if err != nil {
		t.Fatalf("GetPopularPeople: %v", err)
	}
	if got.Results[0].Name != "Famous" {
		t.Errorf("Name = %q", got.Results[0].Name)
	}
}

func TestGetLanguages(t *testing.T) {
	t.Parallel()
	want := []tmdb.Language{{ISO6391: "en", EnglishName: "English"}}
	srv := newTestServer(t, "/configuration/languages", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetLanguages(context.Background())
	if err != nil {
		t.Fatalf("GetLanguages: %v", err)
	}
	if got[0].EnglishName != "English" {
		t.Errorf("Name = %q", got[0].EnglishName)
	}
}

func TestGetCountries(t *testing.T) {
	t.Parallel()
	want := []tmdb.Country{{ISO31661: "US", EnglishName: "United States"}}
	srv := newTestServer(t, "/configuration/countries", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetCountries(context.Background())
	if err != nil {
		t.Fatalf("GetCountries: %v", err)
	}
	if got[0].EnglishName != "United States" {
		t.Errorf("Name = %q", got[0].EnglishName)
	}
}

func TestGetTimezones(t *testing.T) {
	t.Parallel()
	want := []tmdb.Timezone{{ISO31661: "US", Zones: []string{"America/New_York"}}}
	srv := newTestServer(t, "/configuration/timezones", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetTimezones(context.Background())
	if err != nil {
		t.Fatalf("GetTimezones: %v", err)
	}
	if got[0].Zones[0] != "America/New_York" {
		t.Errorf("Zone = %q", got[0].Zones[0])
	}
}

func TestGetJobs(t *testing.T) {
	t.Parallel()
	want := []tmdb.Department{{Department: "Directing", Jobs: []string{"Director"}}}
	srv := newTestServer(t, "/configuration/jobs", want)
	defer srv.Close()
	got, err := newClient(t, srv).GetJobs(context.Background())
	if err != nil {
		t.Fatalf("GetJobs: %v", err)
	}
	if got[0].Department != "Directing" {
		t.Errorf("Department = %q", got[0].Department)
	}
}
