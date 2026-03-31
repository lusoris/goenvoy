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
