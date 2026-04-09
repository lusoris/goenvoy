package trakt_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/metadata/tracking/trakt"
)

func newTestServer(t *testing.T, wantPath, wantKey string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if got := r.Header.Get("trakt-api-key"); got != wantKey {
			t.Errorf("trakt-api-key = %q, want %q", got, wantKey)
		}
		if got := r.Header.Get("trakt-api-version"); got != "2" {
			t.Errorf("trakt-api-version = %q, want %q", got, "2")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Pagination-Page", "1")
		w.Header().Set("X-Pagination-Limit", "10")
		w.Header().Set("X-Pagination-Page-Count", "5")
		w.Header().Set("X-Pagination-Item-Count", "50")
		_ = json.NewEncoder(w).Encode(response)
	}))
}

func TestGetMovie(t *testing.T) {
	ts := newTestServer(t, "/movies/the-dark-knight", "test-key", trakt.Movie{
		Title:   "The Dark Knight",
		Year:    2008,
		IDs:     trakt.IDs{Trakt: 120, Slug: "the-dark-knight", IMDb: "tt0468569", TMDb: 155},
		Genres:  []string{"action", "crime", "drama"},
		Runtime: 152,
	})
	defer ts.Close()

	c := trakt.New("test-key", trakt.WithBaseURL(ts.URL))
	movie, err := c.GetMovie(context.Background(), "the-dark-knight")
	if err != nil {
		t.Fatal(err)
	}
	if movie.Title != "The Dark Knight" {
		t.Errorf("Title = %q, want %q", movie.Title, "The Dark Knight")
	}
	if movie.Year != 2008 {
		t.Errorf("Year = %d, want %d", movie.Year, 2008)
	}
	if movie.IDs.IMDb != "tt0468569" {
		t.Errorf("IMDb = %q, want %q", movie.IDs.IMDb, "tt0468569")
	}
	if movie.Runtime != 152 {
		t.Errorf("Runtime = %d, want %d", movie.Runtime, 152)
	}
}

func TestGetShow(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad", "show-key", trakt.Show{
		Title:         "Breaking Bad",
		Year:          2008,
		IDs:           trakt.IDs{Trakt: 1388, Slug: "breaking-bad", IMDb: "tt0903747", TMDb: 1396, TVDb: 81189},
		Status:        "ended",
		AiredEpisodes: 62,
	})
	defer ts.Close()

	c := trakt.New("show-key", trakt.WithBaseURL(ts.URL))
	show, err := c.GetShow(context.Background(), "breaking-bad")
	if err != nil {
		t.Fatal(err)
	}
	if show.Title != "Breaking Bad" {
		t.Errorf("Title = %q, want %q", show.Title, "Breaking Bad")
	}
	if show.AiredEpisodes != 62 {
		t.Errorf("AiredEpisodes = %d, want %d", show.AiredEpisodes, 62)
	}
	if show.IDs.TVDb != 81189 {
		t.Errorf("TVDb = %d, want %d", show.IDs.TVDb, 81189)
	}
}

func TestGetEpisode(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad/seasons/1/episodes/1", "ep-key", trakt.Episode{
		Season:  1,
		Number:  1,
		Title:   "Pilot",
		IDs:     trakt.IDs{Trakt: 62085, IMDb: "tt0959621", TMDb: 62085, TVDb: 349232},
		Runtime: 58,
	})
	defer ts.Close()

	c := trakt.New("ep-key", trakt.WithBaseURL(ts.URL))
	ep, err := c.GetEpisode(context.Background(), "breaking-bad", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if ep.Title != "Pilot" {
		t.Errorf("Title = %q, want %q", ep.Title, "Pilot")
	}
	if ep.Season != 1 || ep.Number != 1 {
		t.Errorf("S%02dE%02d, want S01E01", ep.Season, ep.Number)
	}
}

func TestTrendingMovies(t *testing.T) {
	ts := newTestServer(t, "/movies/trending", "trend-key", []trakt.TrendingMovie{
		{Watchers: 85, Movie: trakt.Movie{Title: "Oppenheimer", Year: 2023, IDs: trakt.IDs{Trakt: 717468}}},
		{Watchers: 72, Movie: trakt.Movie{Title: "Barbie", Year: 2023, IDs: trakt.IDs{Trakt: 488552}}},
	})
	defer ts.Close()

	c := trakt.New("trend-key", trakt.WithBaseURL(ts.URL))
	movies, pg, err := c.TrendingMovies(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(movies) != 2 {
		t.Fatalf("len = %d, want 2", len(movies))
	}
	if movies[0].Watchers != 85 {
		t.Errorf("Watchers = %d, want 85", movies[0].Watchers)
	}
	if movies[0].Movie.Title != "Oppenheimer" {
		t.Errorf("Title = %q, want %q", movies[0].Movie.Title, "Oppenheimer")
	}
	if pg.PageCount != 5 {
		t.Errorf("PageCount = %d, want 5", pg.PageCount)
	}
	if pg.ItemCount != 50 {
		t.Errorf("ItemCount = %d, want 50", pg.ItemCount)
	}
}

func TestPopularShows(t *testing.T) {
	ts := newTestServer(t, "/shows/popular", "pop-key", []trakt.Show{
		{Title: "Game of Thrones", Year: 2011, IDs: trakt.IDs{Trakt: 1390, Slug: "game-of-thrones"}},
	})
	defer ts.Close()

	c := trakt.New("pop-key", trakt.WithBaseURL(ts.URL))
	shows, _, err := c.PopularShows(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 {
		t.Fatalf("len = %d, want 1", len(shows))
	}
	if shows[0].Title != "Game of Thrones" {
		t.Errorf("Title = %q, want %q", shows[0].Title, "Game of Thrones")
	}
}

func TestSearchText(t *testing.T) {
	ts := newTestServer(t, "/search/movie", "search-key", []trakt.SearchResult{
		{
			Type:  "movie",
			Score: 1000,
			Movie: &trakt.Movie{Title: "Inception", Year: 2010, IDs: trakt.IDs{Trakt: 16662}},
		},
	})
	defer ts.Close()

	c := trakt.New("search-key", trakt.WithBaseURL(ts.URL))
	results, _, err := c.SearchText(context.Background(), "inception", "movie", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("len = %d, want 1", len(results))
	}
	if results[0].Movie.Title != "Inception" {
		t.Errorf("Title = %q, want %q", results[0].Movie.Title, "Inception")
	}
	if results[0].Score != 1000 {
		t.Errorf("Score = %f, want 1000", results[0].Score)
	}
}

func TestSearchByID(t *testing.T) {
	ts := newTestServer(t, "/search/imdb/tt0468569", "id-key", []trakt.SearchResult{
		{
			Type:  "movie",
			Score: 1000,
			Movie: &trakt.Movie{Title: "The Dark Knight", Year: 2008, IDs: trakt.IDs{IMDb: "tt0468569"}},
		},
	})
	defer ts.Close()

	c := trakt.New("id-key", trakt.WithBaseURL(ts.URL))
	results, err := c.SearchByID(context.Background(), "imdb", "tt0468569", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("len = %d, want 1", len(results))
	}
	if results[0].Movie.IDs.IMDb != "tt0468569" {
		t.Errorf("IMDb = %q, want %q", results[0].Movie.IDs.IMDb, "tt0468569")
	}
}

func TestGetMovieRatings(t *testing.T) {
	ts := newTestServer(t, "/movies/the-dark-knight/ratings", "rate-key", trakt.Ratings{
		Rating: 9.0, Votes: 42000,
		Distribution: trakt.Distribution{Ten: 20000, Nine: 12000, Eight: 5000, Seven: 3000, Six: 1000, Five: 500, Four: 200, Three: 100, Two: 100, One: 100},
	})
	defer ts.Close()

	c := trakt.New("rate-key", trakt.WithBaseURL(ts.URL))
	r, err := c.GetMovieRatings(context.Background(), "the-dark-knight")
	if err != nil {
		t.Fatal(err)
	}
	if r.Rating != 9.0 {
		t.Errorf("Rating = %f, want 9.0", r.Rating)
	}
	if r.Votes != 42000 {
		t.Errorf("Votes = %d, want 42000", r.Votes)
	}
	if r.Distribution.Ten != 20000 {
		t.Errorf("Distribution.Ten = %d, want 20000", r.Distribution.Ten)
	}
}

func TestGetMoviePeople(t *testing.T) {
	ts := newTestServer(t, "/movies/the-dark-knight/people", "ppl-key", trakt.People{
		Cast: []trakt.CastMember{
			{Characters: []string{"Bruce Wayne"}, Person: trakt.Person{Name: "Christian Bale", IDs: trakt.IDs{Trakt: 1}}},
			{Characters: []string{"The Joker"}, Person: trakt.Person{Name: "Heath Ledger", IDs: trakt.IDs{Trakt: 2}}},
		},
	})
	defer ts.Close()

	c := trakt.New("ppl-key", trakt.WithBaseURL(ts.URL))
	p, err := c.GetMoviePeople(context.Background(), "the-dark-knight")
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Cast) != 2 {
		t.Fatalf("len(cast) = %d, want 2", len(p.Cast))
	}
	if p.Cast[0].Person.Name != "Christian Bale" {
		t.Errorf("Name = %q, want %q", p.Cast[0].Person.Name, "Christian Bale")
	}
	if p.Cast[1].Characters[0] != "The Joker" {
		t.Errorf("Character = %q, want %q", p.Cast[1].Characters[0], "The Joker")
	}
}

func TestGetShowSeasons(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad/seasons", "season-key", []trakt.Season{
		{Number: 0, Title: "Specials", IDs: trakt.IDs{Trakt: 3962}},
		{Number: 1, Title: "Season 1", IDs: trakt.IDs{Trakt: 3963}, EpisodeCount: 7},
		{Number: 2, Title: "Season 2", IDs: trakt.IDs{Trakt: 3964}, EpisodeCount: 13},
	})
	defer ts.Close()

	c := trakt.New("season-key", trakt.WithBaseURL(ts.URL))
	seasons, err := c.GetShowSeasons(context.Background(), "breaking-bad")
	if err != nil {
		t.Fatal(err)
	}
	if len(seasons) != 3 {
		t.Fatalf("len = %d, want 3", len(seasons))
	}
	if seasons[1].EpisodeCount != 7 {
		t.Errorf("EpisodeCount = %d, want 7", seasons[1].EpisodeCount)
	}
}

func TestCalendarMovies(t *testing.T) {
	ts := newTestServer(t, "/calendars/all/movies/2024-01-01/7", "cal-key", []trakt.CalendarMovie{
		{Released: "2024-01-03", Movie: trakt.Movie{Title: "Migration", Year: 2023, IDs: trakt.IDs{Trakt: 123}}},
	})
	defer ts.Close()

	c := trakt.New("cal-key", trakt.WithBaseURL(ts.URL))
	movies, err := c.CalendarMovies(context.Background(), "2024-01-01", 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(movies) != 1 {
		t.Fatalf("len = %d, want 1", len(movies))
	}
	if movies[0].Movie.Title != "Migration" {
		t.Errorf("Title = %q, want %q", movies[0].Movie.Title, "Migration")
	}
}

func TestGenres(t *testing.T) {
	ts := newTestServer(t, "/genres/movies", "genre-key", []trakt.Genre{
		{Name: "Action", Slug: "action"},
		{Name: "Adventure", Slug: "adventure"},
		{Name: "Comedy", Slug: "comedy"},
	})
	defer ts.Close()

	c := trakt.New("genre-key", trakt.WithBaseURL(ts.URL))
	genres, err := c.Genres(context.Background(), "movies")
	if err != nil {
		t.Fatal(err)
	}
	if len(genres) != 3 {
		t.Fatalf("len = %d, want 3", len(genres))
	}
	if genres[0].Slug != "action" {
		t.Errorf("Slug = %q, want %q", genres[0].Slug, "action")
	}
}

func TestGetPerson(t *testing.T) {
	ts := newTestServer(t, "/people/bryan-cranston", "person-key", trakt.Person{
		Name:     "Bryan Cranston",
		IDs:      trakt.IDs{Trakt: 297891, Slug: "bryan-cranston", IMDb: "nm0186505", TMDb: 17419},
		Birthday: "1956-03-07",
		Gender:   "male",
	})
	defer ts.Close()

	c := trakt.New("person-key", trakt.WithBaseURL(ts.URL))
	p, err := c.GetPerson(context.Background(), "bryan-cranston")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Bryan Cranston" {
		t.Errorf("Name = %q, want %q", p.Name, "Bryan Cranston")
	}
	if p.Birthday != "1956-03-07" {
		t.Errorf("Birthday = %q, want %q", p.Birthday, "1956-03-07")
	}
}

func TestGetMovieStats(t *testing.T) {
	ts := newTestServer(t, "/movies/the-dark-knight/stats", "stat-key", trakt.Stats{
		Watchers: 100000, Plays: 150000, Collectors: 80000, Votes: 42000,
	})
	defer ts.Close()

	c := trakt.New("stat-key", trakt.WithBaseURL(ts.URL))
	s, err := c.GetMovieStats(context.Background(), "the-dark-knight")
	if err != nil {
		t.Fatal(err)
	}
	if s.Watchers != 100000 {
		t.Errorf("Watchers = %d, want 100000", s.Watchers)
	}
	if s.Plays != 150000 {
		t.Errorf("Plays = %d, want 150000", s.Plays)
	}
}

func TestGetMovieStudios(t *testing.T) {
	ts := newTestServer(t, "/movies/the-dark-knight/studios", "studio-key", []trakt.Studio{
		{Name: "Warner Bros. Pictures", Country: "us", IDs: trakt.IDs{Trakt: 174}},
		{Name: "Legendary Pictures", Country: "us", IDs: trakt.IDs{Trakt: 923}},
	})
	defer ts.Close()

	c := trakt.New("studio-key", trakt.WithBaseURL(ts.URL))
	studios, err := c.GetMovieStudios(context.Background(), "the-dark-knight")
	if err != nil {
		t.Fatal(err)
	}
	if len(studios) != 2 {
		t.Fatalf("len = %d, want 2", len(studios))
	}
	if studios[0].Name != "Warner Bros. Pictures" {
		t.Errorf("Name = %q, want %q", studios[0].Name, "Warner Bros. Pictures")
	}
}

func TestAnticipatedMovies(t *testing.T) {
	ts := newTestServer(t, "/movies/anticipated", "anti-key", []trakt.AnticipatedMovie{
		{ListCount: 5000, Movie: trakt.Movie{Title: "Dune: Part Two", Year: 2024, IDs: trakt.IDs{Trakt: 800100}}},
	})
	defer ts.Close()

	c := trakt.New("anti-key", trakt.WithBaseURL(ts.URL))
	movies, _, err := c.AnticipatedMovies(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(movies) != 1 {
		t.Fatalf("len = %d, want 1", len(movies))
	}
	if movies[0].ListCount != 5000 {
		t.Errorf("ListCount = %d, want 5000", movies[0].ListCount)
	}
}

func TestBoxOfficeMovies(t *testing.T) {
	ts := newTestServer(t, "/movies/boxoffice", "box-key", []trakt.BoxOfficeMovie{
		{Revenue: 100000000, Movie: trakt.Movie{Title: "Inside Out 2", Year: 2024, IDs: trakt.IDs{Trakt: 900123}}},
	})
	defer ts.Close()

	c := trakt.New("box-key", trakt.WithBaseURL(ts.URL))
	movies, err := c.BoxOfficeMovies(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(movies) != 1 {
		t.Fatalf("len = %d, want 1", len(movies))
	}
	if movies[0].Revenue != 100000000 {
		t.Errorf("Revenue = %d, want 100000000", movies[0].Revenue)
	}
}

func TestGetMovieTranslations(t *testing.T) {
	ts := newTestServer(t, "/movies/the-dark-knight/translations/de", "trans-key", []trakt.MovieTranslation{
		{Title: "The Dark Knight", Overview: "Batman erhebt sich...", Language: "de", Country: "de"},
	})
	defer ts.Close()

	c := trakt.New("trans-key", trakt.WithBaseURL(ts.URL))
	trans, err := c.GetMovieTranslations(context.Background(), "the-dark-knight", "de")
	if err != nil {
		t.Fatal(err)
	}
	if len(trans) != 1 {
		t.Fatalf("len = %d, want 1", len(trans))
	}
	if trans[0].Language != "de" {
		t.Errorf("Language = %q, want %q", trans[0].Language, "de")
	}
}

func TestNetworks(t *testing.T) {
	ts := newTestServer(t, "/networks", "net-key", []trakt.Network{
		{Name: "HBO", IDs: trakt.IDs{Trakt: 8}},
		{Name: "Netflix", IDs: trakt.IDs{Trakt: 213}},
	})
	defer ts.Close()

	c := trakt.New("net-key", trakt.WithBaseURL(ts.URL))
	nets, err := c.Networks(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(nets) != 2 {
		t.Fatalf("len = %d, want 2", len(nets))
	}
	if nets[0].Name != "HBO" {
		t.Errorf("Name = %q, want %q", nets[0].Name, "HBO")
	}
}

func TestAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":             "Unauthorized",
			"error_description": "invalid API key",
		})
	}))
	defer ts.Close()

	c := trakt.New("bad-key", trakt.WithBaseURL(ts.URL))
	_, err := c.GetMovie(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *trakt.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *trakt.APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
	if apiErr.Error_ != "Unauthorized" {
		t.Errorf("Error_ = %q, want %q", apiErr.Error_, "Unauthorized")
	}
	if apiErr.Description != "invalid API key" {
		t.Errorf("Description = %q, want %q", apiErr.Description, "invalid API key")
	}
}

func TestAPIErrorNonJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("<html>Bad Gateway</html>"))
	}))
	defer ts.Close()

	c := trakt.New("test", trakt.WithBaseURL(ts.URL))
	_, err := c.GetMovie(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *trakt.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *trakt.APIError, got %T", err)
	}
	if apiErr.StatusCode != 502 {
		t.Errorf("StatusCode = %d, want 502", apiErr.StatusCode)
	}
	if apiErr.RawBody != "<html>Bad Gateway</html>" {
		t.Errorf("RawBody = %q, want HTML body", apiErr.RawBody)
	}
}

func TestAPIErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		err  trakt.APIError
		want string
	}{
		{"full", trakt.APIError{StatusCode: 401, Error_: "Unauthorized", Description: "bad key"}, "trakt: HTTP 401: Unauthorized: bad key"},
		{"error only", trakt.APIError{StatusCode: 404, Error_: "Not Found"}, "trakt: HTTP 404: Not Found"},
		{"raw body", trakt.APIError{StatusCode: 502, RawBody: "gateway error"}, "trakt: HTTP 502: gateway error"},
		{"code only", trakt.APIError{StatusCode: 500}, "trakt: HTTP 500"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMostPlayedMovies(t *testing.T) {
	ts := newTestServer(t, "/movies/played/weekly", "played-key", []trakt.PlayedMovie{
		{WatcherCount: 5000, PlayCount: 8000, Movie: trakt.Movie{Title: "The Shawshank Redemption", Year: 1994, IDs: trakt.IDs{Trakt: 120}}},
	})
	defer ts.Close()

	c := trakt.New("played-key", trakt.WithBaseURL(ts.URL))
	movies, _, err := c.MostPlayedMovies(context.Background(), "weekly", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(movies) != 1 {
		t.Fatalf("len = %d, want 1", len(movies))
	}
	if movies[0].WatcherCount != 5000 {
		t.Errorf("WatcherCount = %d, want 5000", movies[0].WatcherCount)
	}
}

// OAuth2 tests.

func newOAuthServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func TestGetDeviceCode(t *testing.T) {
	t.Parallel()
	ts := newOAuthServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oauth/device/code" {
			t.Errorf("path = %q, want /oauth/device/code", r.URL.Path)
		}
		if got := r.Header.Get("trakt-api-key"); got != "cid" {
			t.Errorf("trakt-api-key = %q, want %q", got, "cid")
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["client_id"] != "cid" {
			t.Errorf("client_id = %q, want %q", body["client_id"], "cid")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trakt.DeviceCode{
			DeviceCode:      "dev-123",
			UserCode:        "A1B2C3",
			VerificationURL: "https://trakt.tv/activate",
			ExpiresIn:       600,
			Interval:        5,
		})
	})
	defer ts.Close()

	c := trakt.New("cid", trakt.WithBaseURL(ts.URL))
	dc, err := c.GetDeviceCode(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if dc.DeviceCode != "dev-123" {
		t.Errorf("DeviceCode = %q, want %q", dc.DeviceCode, "dev-123")
	}
	if dc.UserCode != "A1B2C3" {
		t.Errorf("UserCode = %q, want %q", dc.UserCode, "A1B2C3")
	}
	if dc.ExpiresIn != 600 {
		t.Errorf("ExpiresIn = %d, want 600", dc.ExpiresIn)
	}
}

func TestExchangeCode(t *testing.T) {
	t.Parallel()
	ts := newOAuthServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" {
			t.Errorf("path = %q, want /oauth/token", r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["grant_type"] != "authorization_code" {
			t.Errorf("grant_type = %q, want authorization_code", body["grant_type"])
		}
		if body["code"] != "auth-code-123" {
			t.Errorf("code = %q, want auth-code-123", body["code"])
		}
		if body["client_id"] != "cid" {
			t.Errorf("client_id = %q, want cid", body["client_id"])
		}
		if body["client_secret"] != "secret" {
			t.Errorf("client_secret = %q, want secret", body["client_secret"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trakt.Token{
			AccessToken:  "access-tok",
			TokenType:    "Bearer",
			ExpiresIn:    7776000,
			RefreshToken: "refresh-tok",
			Scope:        "public",
			CreatedAt:    1609459200,
		})
	})
	defer ts.Close()

	var callbackToken trakt.Token
	c := trakt.New("cid",
		trakt.WithBaseURL(ts.URL),
		trakt.WithClientSecret("secret"),
		trakt.WithTokenCallback(func(t trakt.Token) { callbackToken = t }),
	)
	tok, err := c.ExchangeCode(context.Background(), "auth-code-123", "urn:ietf:wg:oauth:2.0:oob")
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "access-tok" {
		t.Errorf("AccessToken = %q, want %q", tok.AccessToken, "access-tok")
	}
	if tok.RefreshToken != "refresh-tok" {
		t.Errorf("RefreshToken = %q, want %q", tok.RefreshToken, "refresh-tok")
	}
	if callbackToken.AccessToken != "access-tok" {
		t.Errorf("callback AccessToken = %q, want %q", callbackToken.AccessToken, "access-tok")
	}
}

func TestRefreshToken(t *testing.T) {
	t.Parallel()
	ts := newOAuthServer(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["grant_type"] != "refresh_token" {
			t.Errorf("grant_type = %q, want refresh_token", body["grant_type"])
		}
		if body["refresh_token"] != "old-refresh" {
			t.Errorf("refresh_token = %q, want old-refresh", body["refresh_token"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trakt.Token{
			AccessToken:  "new-access",
			RefreshToken: "new-refresh",
		})
	})
	defer ts.Close()

	c := trakt.New("cid",
		trakt.WithBaseURL(ts.URL),
		trakt.WithClientSecret("secret"),
		trakt.WithRefreshToken("old-refresh"),
	)
	tok, err := c.RefreshToken(context.Background(), "urn:ietf:wg:oauth:2.0:oob")
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "new-access" {
		t.Errorf("AccessToken = %q, want %q", tok.AccessToken, "new-access")
	}
}

func TestRefreshTokenMissing(t *testing.T) {
	t.Parallel()
	c := trakt.New("cid")
	_, err := c.RefreshToken(context.Background(), "")
	if err == nil {
		t.Fatal("expected error when no refresh token set")
	}
}

func TestRevokeToken(t *testing.T) {
	t.Parallel()
	ts := newOAuthServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/revoke" {
			t.Errorf("path = %q, want /oauth/revoke", r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["token"] != "my-token" {
			t.Errorf("token = %q, want my-token", body["token"])
		}
		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()

	c := trakt.New("cid",
		trakt.WithBaseURL(ts.URL),
		trakt.WithClientSecret("secret"),
		trakt.WithAccessToken("my-token"),
	)
	if err := c.RevokeToken(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestBearerTokenInHeader(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer user-tok" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer user-tok")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]trakt.Genre{{Name: "action", Slug: "action"}})
	}))
	defer ts.Close()

	c := trakt.New("cid", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("user-tok"))
	_, err := c.Genres(context.Background(), "movies")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNoBearerTokenWhenEmpty(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			t.Errorf("Authorization = %q, want empty", auth)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]trakt.Genre{{Name: "action", Slug: "action"}})
	}))
	defer ts.Close()

	c := trakt.New("cid", trakt.WithBaseURL(ts.URL))
	_, err := c.Genres(context.Background(), "movies")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetEpisodeRatings(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad/seasons/5/episodes/16/ratings", "erat-key", trakt.Ratings{
		Rating: 9.9, Votes: 30000,
	})
	defer ts.Close()

	c := trakt.New("erat-key", trakt.WithBaseURL(ts.URL))
	r, err := c.GetEpisodeRatings(context.Background(), "breaking-bad", 5, 16)
	if err != nil {
		t.Fatal(err)
	}
	if r.Rating != 9.9 {
		t.Errorf("Rating = %f, want 9.9", r.Rating)
	}
}

func TestCalendarShows(t *testing.T) {
	ts := newTestServer(t, "/calendars/all/shows/2024-03-01/7", "calshow-key", []trakt.CalendarShow{
		{
			FirstAired: "2024-03-04",
			Episode:    trakt.Episode{Season: 2, Number: 1, Title: "Premiere"},
			Show:       trakt.Show{Title: "Shogun", Year: 2024, IDs: trakt.IDs{Trakt: 999}},
		},
	})
	defer ts.Close()

	c := trakt.New("calshow-key", trakt.WithBaseURL(ts.URL))
	shows, err := c.CalendarShows(context.Background(), "2024-03-01", 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 {
		t.Fatalf("len = %d, want 1", len(shows))
	}
	if shows[0].Show.Title != "Shogun" {
		t.Errorf("Title = %q, want %q", shows[0].Show.Title, "Shogun")
	}
}

func TestContextCancellation(t *testing.T) {
	ts := newTestServer(t, "/movies/test", "cancel-key", trakt.Movie{Title: "Test"})
	defer ts.Close()

	c := trakt.New("cancel-key", trakt.WithBaseURL(ts.URL))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.GetMovie(ctx, "test")
	if err == nil {
		t.Fatal("expected error from canceled context")
	}
}

// newAuthServer creates a test server that validates the Authorization Bearer header.
func newAuthServer(t *testing.T, wantMethod, wantPath, wantKey, wantToken string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != wantMethod {
			t.Errorf("method = %s, want %s", r.Method, wantMethod)
		}
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if got := r.Header.Get("trakt-api-key"); got != wantKey {
			t.Errorf("trakt-api-key = %q, want %q", got, wantKey)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+wantToken {
			t.Errorf("Authorization = %q, want %q", got, "Bearer "+wantToken)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Pagination-Page", "1")
		w.Header().Set("X-Pagination-Limit", "10")
		w.Header().Set("X-Pagination-Page-Count", "5")
		w.Header().Set("X-Pagination-Item-Count", "50")
		_ = json.NewEncoder(w).Encode(response)
	}))
}

// Tests for existing untested methods.

func TestGetMovieAliases(t *testing.T) {
	ts := newTestServer(t, "/movies/the-dark-knight/aliases", "alias-key", []trakt.Alias{
		{Title: "Il Cavaliere Oscuro", Country: "it"},
	})
	defer ts.Close()

	c := trakt.New("alias-key", trakt.WithBaseURL(ts.URL))
	aliases, err := c.GetMovieAliases(context.Background(), "the-dark-knight")
	if err != nil {
		t.Fatal(err)
	}
	if len(aliases) != 1 {
		t.Fatalf("len = %d, want 1", len(aliases))
	}
	if aliases[0].Title != "Il Cavaliere Oscuro" {
		t.Errorf("Title = %q, want %q", aliases[0].Title, "Il Cavaliere Oscuro")
	}
}

func TestGetMovieReleases(t *testing.T) {
	ts := newTestServer(t, "/movies/the-dark-knight/releases/us", "rel-key", []trakt.MovieRelease{
		{Country: "us", Certification: "PG-13", ReleaseDate: "2008-07-18", ReleaseType: "theatrical"},
	})
	defer ts.Close()

	c := trakt.New("rel-key", trakt.WithBaseURL(ts.URL))
	releases, err := c.GetMovieReleases(context.Background(), "the-dark-knight", "us")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("len = %d, want 1", len(releases))
	}
	if releases[0].Certification != "PG-13" {
		t.Errorf("Certification = %q, want %q", releases[0].Certification, "PG-13")
	}
}

func TestGetMovieReleasesAllCountries(t *testing.T) {
	ts := newTestServer(t, "/movies/the-dark-knight/releases", "relall-key", []trakt.MovieRelease{
		{Country: "us"},
		{Country: "gb"},
	})
	defer ts.Close()

	c := trakt.New("relall-key", trakt.WithBaseURL(ts.URL))
	releases, err := c.GetMovieReleases(context.Background(), "the-dark-knight", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 2 {
		t.Fatalf("len = %d, want 2", len(releases))
	}
}

func TestPopularMovies(t *testing.T) {
	ts := newTestServer(t, "/movies/popular", "pop-key", []trakt.Movie{
		{Title: "Inception", Year: 2010, IDs: trakt.IDs{Trakt: 16662}},
	})
	defer ts.Close()

	c := trakt.New("pop-key", trakt.WithBaseURL(ts.URL))
	movies, pg, err := c.PopularMovies(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(movies) != 1 {
		t.Fatalf("len = %d, want 1", len(movies))
	}
	if movies[0].Title != "Inception" {
		t.Errorf("Title = %q, want %q", movies[0].Title, "Inception")
	}
	if pg.Page != 1 {
		t.Errorf("Page = %d, want 1", pg.Page)
	}
}

func TestMostWatchedMovies(t *testing.T) {
	ts := newTestServer(t, "/movies/watched/weekly", "mw-key", []trakt.PlayedMovie{
		{WatcherCount: 500, Movie: trakt.Movie{Title: "Fight Club"}},
	})
	defer ts.Close()

	c := trakt.New("mw-key", trakt.WithBaseURL(ts.URL))
	movies, _, err := c.MostWatchedMovies(context.Background(), "weekly", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(movies) != 1 {
		t.Fatalf("len = %d, want 1", len(movies))
	}
}

func TestGetShowAliases(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad/aliases", "sa-key", []trakt.Alias{
		{Title: "Totál Szívás", Country: "hu"},
	})
	defer ts.Close()

	c := trakt.New("sa-key", trakt.WithBaseURL(ts.URL))
	aliases, err := c.GetShowAliases(context.Background(), "breaking-bad")
	if err != nil {
		t.Fatal(err)
	}
	if len(aliases) != 1 {
		t.Fatalf("len = %d, want 1", len(aliases))
	}
}

func TestGetShowTranslations(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad/translations/de", "st-key", []trakt.ShowTranslation{
		{Title: "Breaking Bad", Overview: "Ein Chemielehrer...", Language: "de"},
	})
	defer ts.Close()

	c := trakt.New("st-key", trakt.WithBaseURL(ts.URL))
	translations, err := c.GetShowTranslations(context.Background(), "breaking-bad", "de")
	if err != nil {
		t.Fatal(err)
	}
	if len(translations) != 1 {
		t.Fatalf("len = %d, want 1", len(translations))
	}
	if translations[0].Language != "de" {
		t.Errorf("Language = %q, want %q", translations[0].Language, "de")
	}
}

func TestGetShowPeople(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad/people", "sp-key", trakt.People{
		Cast: []trakt.CastMember{{Characters: []string{"Walter White"}, Person: trakt.Person{Name: "Bryan Cranston"}}},
	})
	defer ts.Close()

	c := trakt.New("sp-key", trakt.WithBaseURL(ts.URL))
	people, err := c.GetShowPeople(context.Background(), "breaking-bad")
	if err != nil {
		t.Fatal(err)
	}
	if len(people.Cast) != 1 {
		t.Fatalf("len(cast) = %d, want 1", len(people.Cast))
	}
	if people.Cast[0].Characters[0] != "Walter White" {
		t.Errorf("Character = %q, want %q", people.Cast[0].Characters[0], "Walter White")
	}
}

func TestGetShowRatings(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad/ratings", "sr-key", trakt.Ratings{Rating: 9.4, Votes: 80000})
	defer ts.Close()

	c := trakt.New("sr-key", trakt.WithBaseURL(ts.URL))
	r, err := c.GetShowRatings(context.Background(), "breaking-bad")
	if err != nil {
		t.Fatal(err)
	}
	if r.Rating != 9.4 {
		t.Errorf("Rating = %f, want 9.4", r.Rating)
	}
}

func TestGetShowStats(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad/stats", "ss-key", trakt.Stats{
		Watchers: 50000, Plays: 200000, Collectors: 30000, Comments: 500, Lists: 10000, Votes: 80000,
	})
	defer ts.Close()

	c := trakt.New("ss-key", trakt.WithBaseURL(ts.URL))
	s, err := c.GetShowStats(context.Background(), "breaking-bad")
	if err != nil {
		t.Fatal(err)
	}
	if s.Watchers != 50000 {
		t.Errorf("Watchers = %d, want 50000", s.Watchers)
	}
}

func TestGetShowStudios(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad/studios", "stu-key", []trakt.Studio{
		{Name: "Sony Pictures Television", IDs: trakt.IDs{Trakt: 1}},
	})
	defer ts.Close()

	c := trakt.New("stu-key", trakt.WithBaseURL(ts.URL))
	studios, err := c.GetShowStudios(context.Background(), "breaking-bad")
	if err != nil {
		t.Fatal(err)
	}
	if len(studios) != 1 {
		t.Fatalf("len = %d, want 1", len(studios))
	}
}

func TestTrendingShows(t *testing.T) {
	ts := newTestServer(t, "/shows/trending", "ts-key", []trakt.TrendingShow{
		{Watchers: 100, Show: trakt.Show{Title: "Shogun", Year: 2024}},
	})
	defer ts.Close()

	c := trakt.New("ts-key", trakt.WithBaseURL(ts.URL))
	shows, _, err := c.TrendingShows(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 {
		t.Fatalf("len = %d, want 1", len(shows))
	}
	if shows[0].Show.Title != "Shogun" {
		t.Errorf("Title = %q, want %q", shows[0].Show.Title, "Shogun")
	}
}

func TestMostPlayedShows(t *testing.T) {
	ts := newTestServer(t, "/shows/played/weekly", "mps-key", []trakt.PlayedShow{
		{WatcherCount: 200, Show: trakt.Show{Title: "House of the Dragon"}},
	})
	defer ts.Close()

	c := trakt.New("mps-key", trakt.WithBaseURL(ts.URL))
	shows, _, err := c.MostPlayedShows(context.Background(), "weekly", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 {
		t.Fatalf("len = %d, want 1", len(shows))
	}
}

func TestMostWatchedShows(t *testing.T) {
	ts := newTestServer(t, "/shows/watched/monthly", "mws-key", []trakt.PlayedShow{
		{WatcherCount: 300, Show: trakt.Show{Title: "The Bear"}},
	})
	defer ts.Close()

	c := trakt.New("mws-key", trakt.WithBaseURL(ts.URL))
	shows, _, err := c.MostWatchedShows(context.Background(), "monthly", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 {
		t.Fatalf("len = %d, want 1", len(shows))
	}
}

func TestAnticipatedShows(t *testing.T) {
	ts := newTestServer(t, "/shows/anticipated", "as-key", []trakt.AnticipatedShow{
		{ListCount: 5000, Show: trakt.Show{Title: "The Last of Us"}},
	})
	defer ts.Close()

	c := trakt.New("as-key", trakt.WithBaseURL(ts.URL))
	shows, _, err := c.AnticipatedShows(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 {
		t.Fatalf("len = %d, want 1", len(shows))
	}
}

func TestGetSeasonEpisodes(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad/seasons/5", "se-key", []trakt.Episode{
		{Season: 5, Number: 1, Title: "Live Free or Die"},
		{Season: 5, Number: 16, Title: "Felina"},
	})
	defer ts.Close()

	c := trakt.New("se-key", trakt.WithBaseURL(ts.URL))
	eps, err := c.GetSeasonEpisodes(context.Background(), "breaking-bad", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(eps) != 2 {
		t.Fatalf("len = %d, want 2", len(eps))
	}
	if eps[1].Title != "Felina" {
		t.Errorf("Title = %q, want %q", eps[1].Title, "Felina")
	}
}

func TestGetEpisodeStats(t *testing.T) {
	ts := newTestServer(t, "/shows/breaking-bad/seasons/5/episodes/16/stats", "es-key", trakt.Stats{
		Watchers: 10000, Plays: 50000,
	})
	defer ts.Close()

	c := trakt.New("es-key", trakt.WithBaseURL(ts.URL))
	s, err := c.GetEpisodeStats(context.Background(), "breaking-bad", 5, 16)
	if err != nil {
		t.Fatal(err)
	}
	if s.Watchers != 10000 {
		t.Errorf("Watchers = %d, want 10000", s.Watchers)
	}
}

func TestCalendarNewShows(t *testing.T) {
	ts := newTestServer(t, "/calendars/all/shows/new/2024-01-01/30", "cns-key", []trakt.CalendarShow{
		{FirstAired: "2024-01-15", Show: trakt.Show{Title: "New Show"}},
	})
	defer ts.Close()

	c := trakt.New("cns-key", trakt.WithBaseURL(ts.URL))
	shows, err := c.CalendarNewShows(context.Background(), "2024-01-01", 30)
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 {
		t.Fatalf("len = %d, want 1", len(shows))
	}
}

func TestCalendarSeasonPremieres(t *testing.T) {
	ts := newTestServer(t, "/calendars/all/shows/premieres/2024-03-01/14", "csp-key", []trakt.CalendarShow{
		{FirstAired: "2024-03-10", Show: trakt.Show{Title: "Premiere Show"}},
	})
	defer ts.Close()

	c := trakt.New("csp-key", trakt.WithBaseURL(ts.URL))
	shows, err := c.CalendarSeasonPremieres(context.Background(), "2024-03-01", 14)
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 {
		t.Fatalf("len = %d, want 1", len(shows))
	}
}

func TestCertifications(t *testing.T) {
	ts := newTestServer(t, "/certifications/movies", "cert-key", []trakt.Certification{
		{Name: "PG-13", Slug: "pg-13", Description: "Parents Strongly Cautioned"},
	})
	defer ts.Close()

	c := trakt.New("cert-key", trakt.WithBaseURL(ts.URL))
	certs, err := c.Certifications(context.Background(), "movies")
	if err != nil {
		t.Fatal(err)
	}
	if len(certs) != 1 {
		t.Fatalf("len = %d, want 1", len(certs))
	}
	if certs[0].Name != "PG-13" {
		t.Errorf("Name = %q, want %q", certs[0].Name, "PG-13")
	}
}

func TestCountries(t *testing.T) {
	ts := newTestServer(t, "/countries/movies", "co-key", []trakt.Country{
		{Name: "United States", Code: "us"},
	})
	defer ts.Close()

	c := trakt.New("co-key", trakt.WithBaseURL(ts.URL))
	countries, err := c.Countries(context.Background(), "movies")
	if err != nil {
		t.Fatal(err)
	}
	if len(countries) != 1 {
		t.Fatalf("len = %d, want 1", len(countries))
	}
	if countries[0].Code != "us" {
		t.Errorf("Code = %q, want %q", countries[0].Code, "us")
	}
}

func TestLanguages(t *testing.T) {
	ts := newTestServer(t, "/languages/movies", "lang-key", []trakt.Language{
		{Name: "English", Code: "en"},
	})
	defer ts.Close()

	c := trakt.New("lang-key", trakt.WithBaseURL(ts.URL))
	langs, err := c.Languages(context.Background(), "movies")
	if err != nil {
		t.Fatal(err)
	}
	if len(langs) != 1 {
		t.Fatalf("len = %d, want 1", len(langs))
	}
	if langs[0].Code != "en" {
		t.Errorf("Code = %q, want %q", langs[0].Code, "en")
	}
}

// Tests for new user-authenticated methods.

func TestGetUpdatedMovies(t *testing.T) {
	ts := newTestServer(t, "/movies/updates/2024-01-01", "upd-key", []trakt.UpdatedMovie{
		{UpdatedAt: "2024-01-02T10:00:00.000Z", Movie: trakt.Movie{Title: "Updated Film", IDs: trakt.IDs{Trakt: 1}}},
	})
	defer ts.Close()

	c := trakt.New("upd-key", trakt.WithBaseURL(ts.URL))
	movies, pg, err := c.GetUpdatedMovies(context.Background(), "2024-01-01", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(movies) != 1 {
		t.Fatalf("len = %d, want 1", len(movies))
	}
	if movies[0].Movie.Title != "Updated Film" {
		t.Errorf("Title = %q, want %q", movies[0].Movie.Title, "Updated Film")
	}
	if pg.Page != 1 {
		t.Errorf("Page = %d, want 1", pg.Page)
	}
}

func TestGetUpdatedShows(t *testing.T) {
	ts := newTestServer(t, "/shows/updates/2024-06-01", "upds-key", []trakt.UpdatedShow{
		{UpdatedAt: "2024-06-02T08:00:00.000Z", Show: trakt.Show{Title: "Updated Show"}},
	})
	defer ts.Close()

	c := trakt.New("upds-key", trakt.WithBaseURL(ts.URL))
	shows, _, err := c.GetUpdatedShows(context.Background(), "2024-06-01", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 {
		t.Fatalf("len = %d, want 1", len(shows))
	}
	if shows[0].Show.Title != "Updated Show" {
		t.Errorf("Title = %q, want %q", shows[0].Show.Title, "Updated Show")
	}
}

func TestGetProfile(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/users/me", "prof-key", "user-tok", trakt.UserProfile{
		Username: "sean", Name: "Sean Rudford", VIP: true, JoinedAt: "2010-09-25T17:49:25.000Z",
	})
	defer ts.Close()

	c := trakt.New("prof-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("user-tok"))
	profile, err := c.GetProfile(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if profile.Username != "sean" {
		t.Errorf("Username = %q, want %q", profile.Username, "sean")
	}
	if !profile.VIP {
		t.Error("VIP = false, want true")
	}
}

func TestGetUserStats(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/users/me/stats", "stats-key", "user-tok", trakt.UserStats{
		Movies:   trakt.UserMovieStats{Plays: 500, Watched: 480, Minutes: 60000},
		Episodes: trakt.UserEpisodeStats{Plays: 5000, Watched: 4500, Minutes: 200000},
	})
	defer ts.Close()

	c := trakt.New("stats-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("user-tok"))
	stats, err := c.GetUserStats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if stats.Movies.Plays != 500 {
		t.Errorf("Movies.Plays = %d, want 500", stats.Movies.Plays)
	}
	if stats.Episodes.Watched != 4500 {
		t.Errorf("Episodes.Watched = %d, want 4500", stats.Episodes.Watched)
	}
}

func TestGetWatchlist(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/sync/watchlist/movies", "wl-key", "wl-tok", []trakt.WatchlistItem{
		{Rank: 1, ListedAt: "2024-01-01T00:00:00.000Z", Type: "movie", Movie: &trakt.Movie{Title: "Dune: Part Three", IDs: trakt.IDs{Trakt: 1}}},
	})
	defer ts.Close()

	c := trakt.New("wl-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("wl-tok"))
	items, pg, err := c.GetWatchlist(context.Background(), "movies", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("len = %d, want 1", len(items))
	}
	if items[0].Movie.Title != "Dune: Part Three" {
		t.Errorf("Title = %q, want %q", items[0].Movie.Title, "Dune: Part Three")
	}
	if pg.ItemCount != 50 {
		t.Errorf("ItemCount = %d, want 50", pg.ItemCount)
	}
}

func TestGetWatchlistAll(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/sync/watchlist", "wla-key", "wla-tok", []trakt.WatchlistItem{
		{Rank: 1, Type: "movie"},
		{Rank: 2, Type: "show"},
	})
	defer ts.Close()

	c := trakt.New("wla-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("wla-tok"))
	items, _, err := c.GetWatchlist(context.Background(), "", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("len = %d, want 2", len(items))
	}
}

func TestAddToWatchlist(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/sync/watchlist", "aw-key", "aw-tok", trakt.SyncResponse{
		Added: &trakt.SyncCount{Movies: 1},
	})
	defer ts.Close()

	c := trakt.New("aw-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("aw-tok"))
	resp, err := c.AddToWatchlist(context.Background(), &trakt.SyncItems{
		Movies: []trakt.SyncMovie{{IDs: trakt.IDs{Trakt: 120}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Added.Movies != 1 {
		t.Errorf("Added.Movies = %d, want 1", resp.Added.Movies)
	}
}

func TestRemoveFromWatchlist(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/sync/watchlist/remove", "rw-key", "rw-tok", trakt.SyncResponse{
		Deleted: &trakt.SyncCount{Movies: 1},
	})
	defer ts.Close()

	c := trakt.New("rw-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("rw-tok"))
	resp, err := c.RemoveFromWatchlist(context.Background(), &trakt.SyncItems{
		Movies: []trakt.SyncMovie{{IDs: trakt.IDs{Trakt: 120}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Deleted.Movies != 1 {
		t.Errorf("Deleted.Movies = %d, want 1", resp.Deleted.Movies)
	}
}

func TestGetCollection(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/sync/collection/movies", "gc-key", "gc-tok", []trakt.CollectionItem{
		{CollectedAt: "2024-01-01T00:00:00.000Z", Movie: &trakt.Movie{Title: "Inception"}},
	})
	defer ts.Close()

	c := trakt.New("gc-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("gc-tok"))
	items, _, err := c.GetCollection(context.Background(), "movies", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("len = %d, want 1", len(items))
	}
	if items[0].Movie.Title != "Inception" {
		t.Errorf("Title = %q, want %q", items[0].Movie.Title, "Inception")
	}
}

func TestAddToCollection(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/sync/collection", "ac-key", "ac-tok", trakt.SyncResponse{
		Added: &trakt.SyncCount{Movies: 2, Shows: 1},
	})
	defer ts.Close()

	c := trakt.New("ac-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("ac-tok"))
	resp, err := c.AddToCollection(context.Background(), &trakt.SyncItems{
		Movies: []trakt.SyncMovie{{IDs: trakt.IDs{Trakt: 120}}, {IDs: trakt.IDs{Trakt: 121}}},
		Shows:  []trakt.SyncShow{{IDs: trakt.IDs{Trakt: 1388}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Added.Movies != 2 {
		t.Errorf("Added.Movies = %d, want 2", resp.Added.Movies)
	}
}

func TestRemoveFromCollection(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/sync/collection/remove", "rc-key", "rc-tok", trakt.SyncResponse{
		Deleted: &trakt.SyncCount{Movies: 1},
	})
	defer ts.Close()

	c := trakt.New("rc-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("rc-tok"))
	resp, err := c.RemoveFromCollection(context.Background(), &trakt.SyncItems{
		Movies: []trakt.SyncMovie{{IDs: trakt.IDs{Trakt: 120}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Deleted.Movies != 1 {
		t.Errorf("Deleted.Movies = %d, want 1", resp.Deleted.Movies)
	}
}

func TestGetHistory(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/sync/history/movies", "gh-key", "gh-tok", []trakt.HistoryItem{
		{ID: 123, WatchedAt: "2024-01-15T20:00:00.000Z", Action: "watch", Type: "movie", Movie: &trakt.Movie{Title: "Oppenheimer"}},
	})
	defer ts.Close()

	c := trakt.New("gh-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("gh-tok"))
	items, _, err := c.GetHistory(context.Background(), "movies", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("len = %d, want 1", len(items))
	}
	if items[0].Movie.Title != "Oppenheimer" {
		t.Errorf("Title = %q, want %q", items[0].Movie.Title, "Oppenheimer")
	}
}

func TestGetHistoryAll(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/sync/history", "gha-key", "gha-tok", []trakt.HistoryItem{
		{ID: 1, Type: "movie"},
		{ID: 2, Type: "episode"},
	})
	defer ts.Close()

	c := trakt.New("gha-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("gha-tok"))
	items, _, err := c.GetHistory(context.Background(), "", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("len = %d, want 2", len(items))
	}
}

func TestAddToHistory(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/sync/history", "ah-key", "ah-tok", trakt.SyncResponse{
		Added: &trakt.SyncCount{Movies: 1, Episodes: 3},
	})
	defer ts.Close()

	c := trakt.New("ah-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("ah-tok"))
	resp, err := c.AddToHistory(context.Background(), &trakt.SyncItems{
		Movies: []trakt.SyncMovie{{IDs: trakt.IDs{Trakt: 120}, WatchedAt: "2024-01-01T00:00:00.000Z"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Added.Movies != 1 {
		t.Errorf("Added.Movies = %d, want 1", resp.Added.Movies)
	}
}

func TestRemoveFromHistory(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/sync/history/remove", "rh-key", "rh-tok", trakt.SyncResponse{
		Deleted: &trakt.SyncCount{Episodes: 2},
	})
	defer ts.Close()

	c := trakt.New("rh-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("rh-tok"))
	resp, err := c.RemoveFromHistory(context.Background(), &trakt.SyncItems{
		Episodes: []trakt.SyncEpisode{{IDs: trakt.IDs{Trakt: 1}}, {IDs: trakt.IDs{Trakt: 2}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Deleted.Episodes != 2 {
		t.Errorf("Deleted.Episodes = %d, want 2", resp.Deleted.Episodes)
	}
}

func TestGetRatings(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/sync/ratings/movies", "gr-key", "gr-tok", []trakt.RatedItem{
		{RatedAt: "2024-01-01T00:00:00.000Z", Rating: 10, Type: "movie", Movie: &trakt.Movie{Title: "The Shawshank Redemption"}},
	})
	defer ts.Close()

	c := trakt.New("gr-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("gr-tok"))
	items, err := c.GetRatings(context.Background(), "movies")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("len = %d, want 1", len(items))
	}
	if items[0].Rating != 10 {
		t.Errorf("Rating = %d, want 10", items[0].Rating)
	}
}

func TestGetRatingsAll(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/sync/ratings", "gra-key", "gra-tok", []trakt.RatedItem{
		{Rating: 8, Type: "movie"},
		{Rating: 9, Type: "show"},
	})
	defer ts.Close()

	c := trakt.New("gra-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("gra-tok"))
	items, err := c.GetRatings(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("len = %d, want 2", len(items))
	}
}

func TestAddRatings(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/sync/ratings", "ar-key", "ar-tok", trakt.SyncResponse{
		Added: &trakt.SyncCount{Movies: 1},
	})
	defer ts.Close()

	c := trakt.New("ar-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("ar-tok"))
	resp, err := c.AddRatings(context.Background(), &trakt.SyncItems{
		Movies: []trakt.SyncMovie{{IDs: trakt.IDs{Trakt: 120}, Rating: 10}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Added.Movies != 1 {
		t.Errorf("Added.Movies = %d, want 1", resp.Added.Movies)
	}
}

func TestRemoveRatings(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/sync/ratings/remove", "rr-key", "rr-tok", trakt.SyncResponse{
		Deleted: &trakt.SyncCount{Movies: 1},
	})
	defer ts.Close()

	c := trakt.New("rr-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("rr-tok"))
	resp, err := c.RemoveRatings(context.Background(), &trakt.SyncItems{
		Movies: []trakt.SyncMovie{{IDs: trakt.IDs{Trakt: 120}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Deleted.Movies != 1 {
		t.Errorf("Deleted.Movies = %d, want 1", resp.Deleted.Movies)
	}
}

func TestGetUserLists(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/users/me/lists", "ul-key", "ul-tok", []trakt.UserList{
		{Name: "Marvel", Description: "MCU films", Privacy: "public", ItemCount: 30, IDs: trakt.IDs{Trakt: 55}},
	})
	defer ts.Close()

	c := trakt.New("ul-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("ul-tok"))
	lists, err := c.GetUserLists(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(lists) != 1 {
		t.Fatalf("len = %d, want 1", len(lists))
	}
	if lists[0].Name != "Marvel" {
		t.Errorf("Name = %q, want %q", lists[0].Name, "Marvel")
	}
}

func TestCreateList(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/users/me/lists", "cl-key", "cl-tok", trakt.UserList{
		Name: "Horror", Privacy: "private", IDs: trakt.IDs{Trakt: 100},
	})
	defer ts.Close()

	c := trakt.New("cl-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("cl-tok"))
	list, err := c.CreateList(context.Background(), &trakt.UserList{Name: "Horror", Privacy: "private"})
	if err != nil {
		t.Fatal(err)
	}
	if list.Name != "Horror" {
		t.Errorf("Name = %q, want %q", list.Name, "Horror")
	}
}

func TestUpdateList(t *testing.T) {
	ts := newAuthServer(t, http.MethodPut, "/users/me/lists/horror", "upl-key", "upl-tok", nil)
	defer ts.Close()

	c := trakt.New("upl-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("upl-tok"))
	err := c.UpdateList(context.Background(), "horror", &trakt.UserList{Name: "Horror Films", Privacy: "public"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteList(t *testing.T) {
	ts := newAuthServer(t, http.MethodDelete, "/users/me/lists/horror", "dl-key", "dl-tok", nil)
	defer ts.Close()

	c := trakt.New("dl-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("dl-tok"))
	err := c.DeleteList(context.Background(), "horror")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetListItems(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/users/me/lists/marvel/items", "gli-key", "gli-tok", []trakt.ListItem{
		{Rank: 1, Type: "movie", Movie: &trakt.Movie{Title: "Iron Man"}},
		{Rank: 2, Type: "movie", Movie: &trakt.Movie{Title: "The Avengers"}},
	})
	defer ts.Close()

	c := trakt.New("gli-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("gli-tok"))
	items, _, err := c.GetListItems(context.Background(), "marvel", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("len = %d, want 2", len(items))
	}
	if items[0].Movie.Title != "Iron Man" {
		t.Errorf("Title = %q, want %q", items[0].Movie.Title, "Iron Man")
	}
}

func TestAddListItems(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/users/me/lists/marvel/items", "ali-key", "ali-tok", trakt.SyncResponse{
		Added: &trakt.SyncCount{Movies: 1},
	})
	defer ts.Close()

	c := trakt.New("ali-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("ali-tok"))
	resp, err := c.AddListItems(context.Background(), "marvel", &trakt.SyncItems{
		Movies: []trakt.SyncMovie{{IDs: trakt.IDs{Trakt: 120}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Added.Movies != 1 {
		t.Errorf("Added.Movies = %d, want 1", resp.Added.Movies)
	}
}

func TestRemoveListItems(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/users/me/lists/marvel/items/remove", "rli-key", "rli-tok", trakt.SyncResponse{
		Deleted: &trakt.SyncCount{Movies: 1},
	})
	defer ts.Close()

	c := trakt.New("rli-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("rli-tok"))
	resp, err := c.RemoveListItems(context.Background(), "marvel", &trakt.SyncItems{
		Movies: []trakt.SyncMovie{{IDs: trakt.IDs{Trakt: 120}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Deleted.Movies != 1 {
		t.Errorf("Deleted.Movies = %d, want 1", resp.Deleted.Movies)
	}
}

func TestScrobbleStart(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/scrobble/start", "ss-key", "ss-tok", trakt.ScrobbleResponse{
		ID: 1, Action: "start", Movie: &trakt.Movie{Title: "Inception"},
	})
	defer ts.Close()

	c := trakt.New("ss-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("ss-tok"))
	resp, err := c.ScrobbleStart(context.Background(), &trakt.ScrobbleRequest{
		Movie:    &trakt.SyncMovie{IDs: trakt.IDs{Trakt: 16662}},
		Progress: 2.5,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Action != "start" {
		t.Errorf("Action = %q, want %q", resp.Action, "start")
	}
	if resp.Movie.Title != "Inception" {
		t.Errorf("Title = %q, want %q", resp.Movie.Title, "Inception")
	}
}

func TestScrobblePause(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/scrobble/pause", "sp-key", "sp-tok", trakt.ScrobbleResponse{
		ID: 2, Action: "pause",
	})
	defer ts.Close()

	c := trakt.New("sp-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("sp-tok"))
	resp, err := c.ScrobblePause(context.Background(), &trakt.ScrobbleRequest{
		Movie:    &trakt.SyncMovie{IDs: trakt.IDs{Trakt: 16662}},
		Progress: 50.0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Action != "pause" {
		t.Errorf("Action = %q, want %q", resp.Action, "pause")
	}
}

func TestScrobbleStop(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/scrobble/stop", "st-key", "st-tok", trakt.ScrobbleResponse{
		ID: 3, Action: "scrobble",
	})
	defer ts.Close()

	c := trakt.New("st-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("st-tok"))
	resp, err := c.ScrobbleStop(context.Background(), &trakt.ScrobbleRequest{
		Movie:    &trakt.SyncMovie{IDs: trakt.IDs{Trakt: 16662}},
		Progress: 95.0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Action != "scrobble" {
		t.Errorf("Action = %q, want %q", resp.Action, "scrobble")
	}
}

func TestCheckin(t *testing.T) {
	ts := newAuthServer(t, http.MethodPost, "/checkin", "ci-key", "ci-tok", trakt.CheckinResponse{
		ID: 10, WatchedAt: "2024-06-15T20:00:00.000Z", Movie: &trakt.Movie{Title: "Interstellar"},
	})
	defer ts.Close()

	c := trakt.New("ci-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("ci-tok"))
	resp, err := c.Checkin(context.Background(), &trakt.CheckinRequest{
		Movie:   &trakt.SyncMovie{IDs: trakt.IDs{Trakt: 157336}},
		Message: "Watching on the big screen!",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Movie.Title != "Interstellar" {
		t.Errorf("Title = %q, want %q", resp.Movie.Title, "Interstellar")
	}
	if resp.ID != 10 {
		t.Errorf("ID = %d, want 10", resp.ID)
	}
}

func TestCancelCheckin(t *testing.T) {
	ts := newAuthServer(t, http.MethodDelete, "/checkin", "cc-key", "cc-tok", nil)
	defer ts.Close()

	c := trakt.New("cc-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("cc-tok"))
	err := c.CancelCheckin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMovieRecommendations(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/recommendations/movies", "mr-key", "mr-tok", []trakt.Movie{
		{Title: "Arrival", Year: 2016, IDs: trakt.IDs{Trakt: 212691}},
		{Title: "Ex Machina", Year: 2014, IDs: trakt.IDs{Trakt: 184309}},
	})
	defer ts.Close()

	c := trakt.New("mr-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("mr-tok"))
	movies, pg, err := c.GetMovieRecommendations(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(movies) != 2 {
		t.Fatalf("len = %d, want 2", len(movies))
	}
	if movies[0].Title != "Arrival" {
		t.Errorf("Title = %q, want %q", movies[0].Title, "Arrival")
	}
	if pg.PageCount != 5 {
		t.Errorf("PageCount = %d, want 5", pg.PageCount)
	}
}

func TestGetShowRecommendations(t *testing.T) {
	ts := newAuthServer(t, http.MethodGet, "/recommendations/shows", "sr-key", "sr-tok", []trakt.Show{
		{Title: "Severance", Year: 2022, IDs: trakt.IDs{Trakt: 168110}},
	})
	defer ts.Close()

	c := trakt.New("sr-key", trakt.WithBaseURL(ts.URL), trakt.WithAccessToken("sr-tok"))
	shows, _, err := c.GetShowRecommendations(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 {
		t.Fatalf("len = %d, want 1", len(shows))
	}
	if shows[0].Title != "Severance" {
		t.Errorf("Title = %q, want %q", shows[0].Title, "Severance")
	}
}
