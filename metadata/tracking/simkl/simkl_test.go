package simkl_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/metadata/tracking/simkl"
)

func newTestServer(t *testing.T, wantPath, wantKey string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != wantPath {
			t.Errorf("path = %q, want %q", r.URL.Path, wantPath)
		}
		if wantKey != "" {
			if got := r.Header.Get("simkl-api-key"); got != wantKey {
				t.Errorf("simkl-api-key = %q, want %q", got, wantKey)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
}

func TestGetMovie(t *testing.T) {
	movie := simkl.Movie{
		Title:    "The Dark Knight",
		Year:     2008,
		IDs:      simkl.IDs{Simkl: 120, Slug: "the-dark-knight", IMDb: "tt0468569"},
		Runtime:  152,
		Genres:   []string{"action", "crime", "drama"},
		Overview: "When the menace known as the Joker wreaks havoc...",
	}
	ts := newTestServer(t, "/movies/the-dark-knight", "test-key", movie)
	defer ts.Close()

	c := simkl.New("test-key", simkl.WithBaseURL(ts.URL))
	m, err := c.GetMovie(context.Background(), "the-dark-knight")
	if err != nil {
		t.Fatal(err)
	}
	if m.Title != "The Dark Knight" {
		t.Errorf("Title = %q, want %q", m.Title, "The Dark Knight")
	}
	if m.Runtime != 152 {
		t.Errorf("Runtime = %d, want %d", m.Runtime, 152)
	}
	if m.IDs.IMDb != "tt0468569" {
		t.Errorf("IMDb = %q, want %q", m.IDs.IMDb, "tt0468569")
	}
}

func TestTrendingMovies(t *testing.T) {
	movies := []simkl.TrendingMovie{
		{Title: "Movie 1", Rank: 1, IDs: simkl.IDs{Simkl: 1}},
		{Title: "Movie 2", Rank: 2, IDs: simkl.IDs{Simkl: 2}},
	}
	ts := newTestServer(t, "/movies/trending/week", "trend-key", movies)
	defer ts.Close()

	c := simkl.New("trend-key", simkl.WithBaseURL(ts.URL))
	result, err := c.TrendingMovies(context.Background(), "week")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result[0].Title != "Movie 1" {
		t.Errorf("Title = %q, want %q", result[0].Title, "Movie 1")
	}
}

func TestGetShow(t *testing.T) {
	show := simkl.Show{
		Title:         "Breaking Bad",
		Year:          2008,
		IDs:           simkl.IDs{Simkl: 1388, Slug: "breaking-bad"},
		Status:        "ended",
		TotalEpisodes: 62,
	}
	ts := newTestServer(t, "/tv/breaking-bad", "show-key", show)
	defer ts.Close()

	c := simkl.New("show-key", simkl.WithBaseURL(ts.URL))
	s, err := c.GetShow(context.Background(), "breaking-bad")
	if err != nil {
		t.Fatal(err)
	}
	if s.Title != "Breaking Bad" {
		t.Errorf("Title = %q, want %q", s.Title, "Breaking Bad")
	}
	if s.TotalEpisodes != 62 {
		t.Errorf("TotalEpisodes = %d, want %d", s.TotalEpisodes, 62)
	}
}

func TestGetShowEpisodes(t *testing.T) {
	eps := []simkl.Episode{
		{Title: "Pilot", Season: 1, Episode: 1},
		{Title: "Cat's in the Bag", Season: 1, Episode: 2},
	}
	ts := newTestServer(t, "/tv/episodes/1388", "ep-key", eps)
	defer ts.Close()

	c := simkl.New("ep-key", simkl.WithBaseURL(ts.URL))
	result, err := c.GetShowEpisodes(context.Background(), "1388")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result[0].Title != "Pilot" {
		t.Errorf("Title = %q, want %q", result[0].Title, "Pilot")
	}
}

func TestTrendingShows(t *testing.T) {
	shows := []simkl.TrendingShow{{Title: "Show 1", Rank: 1, IDs: simkl.IDs{SimklID: 1}}}
	ts := newTestServer(t, "/tv/trending", "ts-key", shows)
	defer ts.Close()

	c := simkl.New("ts-key", simkl.WithBaseURL(ts.URL))
	result, err := c.TrendingShows(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestShowGenres(t *testing.T) {
	items := []simkl.GenreItem{{Title: "Drama Show", Year: 2024, Rank: 1}}
	ts := newTestServer(t, "/tv/genres/drama", "sg-key", items)
	defer ts.Close()

	c := simkl.New("sg-key", simkl.WithBaseURL(ts.URL))
	result, err := c.ShowGenres(context.Background(), "drama", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Title != "Drama Show" {
		t.Errorf("Title = %q, want %q", result[0].Title, "Drama Show")
	}
}

func TestShowPremieres(t *testing.T) {
	items := []simkl.PremiereItem{{Title: "New Show", Year: 2024}}
	ts := newTestServer(t, "/tv/premieres/new", "sp-key", items)
	defer ts.Close()

	c := simkl.New("sp-key", simkl.WithBaseURL(ts.URL))
	result, err := c.ShowPremieres(context.Background(), "new", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestAiringShows(t *testing.T) {
	items := []simkl.AiringItem{{Title: "Airing Show", Year: 2024}}
	ts := newTestServer(t, "/tv/airing", "as-key", items)
	defer ts.Close()

	c := simkl.New("as-key", simkl.WithBaseURL(ts.URL))
	result, err := c.AiringShows(context.Background(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestBestShows(t *testing.T) {
	items := []simkl.BestItem{{Title: "Best Show", Year: 2024, Ratings: &simkl.Ratings{Simkl: &simkl.RatingScore{Rating: 9.5, Votes: 1000}}}}
	ts := newTestServer(t, "/tv/best/all", "bs-key", items)
	defer ts.Close()

	c := simkl.New("bs-key", simkl.WithBaseURL(ts.URL))
	result, err := c.BestShows(context.Background(), "all")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Ratings == nil || result[0].Ratings.Simkl.Rating != 9.5 {
		t.Error("Ratings not parsed correctly")
	}
}

func TestGetAnime(t *testing.T) {
	anime := simkl.Anime{
		Title:     "Attack on Titan",
		Year:      2013,
		AnimeType: "tv",
		IDs:       simkl.IDs{Simkl: 38, MAL: "16498"},
		Status:    "ended",
	}
	ts := newTestServer(t, "/anime/38", "anime-key", anime)
	defer ts.Close()

	c := simkl.New("anime-key", simkl.WithBaseURL(ts.URL))
	a, err := c.GetAnime(context.Background(), "38")
	if err != nil {
		t.Fatal(err)
	}
	if a.Title != "Attack on Titan" {
		t.Errorf("Title = %q, want %q", a.Title, "Attack on Titan")
	}
	if a.AnimeType != "tv" {
		t.Errorf("AnimeType = %q, want %q", a.AnimeType, "tv")
	}
}

func TestGetAnimeEpisodes(t *testing.T) {
	eps := []simkl.Episode{{Title: "To You, 2000 Years Later", Season: 1, Episode: 1}}
	ts := newTestServer(t, "/anime/episodes/38", "ae-key", eps)
	defer ts.Close()

	c := simkl.New("ae-key", simkl.WithBaseURL(ts.URL))
	result, err := c.GetAnimeEpisodes(context.Background(), "38")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestTrendingAnime(t *testing.T) {
	items := []simkl.TrendingAnime{{Title: "Trending Anime", Rank: 1}}
	ts := newTestServer(t, "/anime/trending/today", "ta-key", items)
	defer ts.Close()

	c := simkl.New("ta-key", simkl.WithBaseURL(ts.URL))
	result, err := c.TrendingAnime(context.Background(), "today")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestAnimePremieres(t *testing.T) {
	items := []simkl.PremiereItem{{Title: "New Anime", Year: 2024}}
	ts := newTestServer(t, "/anime/premieres/new", "ap-key", items)
	defer ts.Close()

	c := simkl.New("ap-key", simkl.WithBaseURL(ts.URL))
	result, err := c.AnimePremieres(context.Background(), "new", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestAiringAnime(t *testing.T) {
	items := []simkl.AiringItem{{Title: "Airing Anime", Episode: &simkl.EpisodeMinimal{Title: "ep1", Episode: 5}}}
	ts := newTestServer(t, "/anime/airing", "aa-key", items)
	defer ts.Close()

	c := simkl.New("aa-key", simkl.WithBaseURL(ts.URL))
	result, err := c.AiringAnime(context.Background(), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Episode == nil || result[0].Episode.Episode != 5 {
		t.Error("Episode not parsed correctly")
	}
}

func TestBestAnime(t *testing.T) {
	items := []simkl.BestItem{{Title: "Best Anime", Year: 2020}}
	ts := newTestServer(t, "/anime/best/voted", "ba-key", items)
	defer ts.Close()

	c := simkl.New("ba-key", simkl.WithBaseURL(ts.URL))
	result, err := c.BestAnime(context.Background(), "voted")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestSearchByID(t *testing.T) {
	results := []simkl.SearchIDResult{
		{Title: "Found Show", Type: "tv", IDs: simkl.IDs{Simkl: 100, Slug: "found-show"}},
	}
	ts := newTestServer(t, "/search/id", "sid-key", results)
	defer ts.Close()

	c := simkl.New("sid-key", simkl.WithBaseURL(ts.URL))
	res, err := c.SearchByID(context.Background(), "imdb", "tt0903747")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("len = %d, want 1", len(res))
	}
	if res[0].Title != "Found Show" {
		t.Errorf("Title = %q, want %q", res[0].Title, "Found Show")
	}
}

func TestSearchText(t *testing.T) {
	results := []simkl.SearchResult{
		{Title: "Breaking Bad", Type: "tv", Year: 2008, IDs: simkl.IDs{Simkl: 1388}},
	}
	ts := newTestServer(t, "/search/tv", "st-key", results)
	defer ts.Close()

	c := simkl.New("st-key", simkl.WithBaseURL(ts.URL))
	res, err := c.SearchText(context.Background(), "tv", "breaking bad", 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("len = %d, want 1", len(res))
	}
	if res[0].Year != 2008 {
		t.Errorf("Year = %d, want %d", res[0].Year, 2008)
	}
}

func TestCalendarShows(t *testing.T) {
	items := []simkl.CalendarShow{{Title: "Cal Show", Date: "2024-06-01"}}
	ts := newTestServer(t, "/calendar/tv.json", "", items)
	defer ts.Close()

	c := simkl.New("cal-key", simkl.WithCalendarURL(ts.URL))
	result, err := c.CalendarShows(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
	if result[0].Title != "Cal Show" {
		t.Errorf("Title = %q, want %q", result[0].Title, "Cal Show")
	}
}

func TestCalendarAnime(t *testing.T) {
	items := []simkl.CalendarAnime{{Title: "Cal Anime", Date: "2024-06-02"}}
	ts := newTestServer(t, "/calendar/anime.json", "", items)
	defer ts.Close()

	c := simkl.New("cal-key", simkl.WithCalendarURL(ts.URL))
	result, err := c.CalendarAnime(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestCalendarMovies(t *testing.T) {
	items := []simkl.CalendarMovie{{Title: "Cal Movie", ReleaseDate: "2024-07-01"}}
	ts := newTestServer(t, "/calendar/movie_release.json", "", items)
	defer ts.Close()

	c := simkl.New("cal-key", simkl.WithCalendarURL(ts.URL))
	result, err := c.CalendarMovies(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestCalendarShowsMonth(t *testing.T) {
	items := []simkl.CalendarShow{{Title: "June Show", Date: "2024-06-15"}}
	ts := newTestServer(t, "/calendar/2024/6/tv.json", "", items)
	defer ts.Close()

	c := simkl.New("cal-key", simkl.WithCalendarURL(ts.URL))
	result, err := c.CalendarShowsMonth(context.Background(), 2024, 6)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusPreconditionFailed)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error":   "client_id_failed",
			"code":    412,
			"message": "Invalid client_id",
		})
	}))
	defer ts.Close()

	c := simkl.New("bad-key", simkl.WithBaseURL(ts.URL))
	_, err := c.GetMovie(context.Background(), "1")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *simkl.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusPreconditionFailed)
	}
	if apiErr.Error_ != "client_id_failed" {
		t.Errorf("Error_ = %q, want %q", apiErr.Error_, "client_id_failed")
	}
	if apiErr.Message != "Invalid client_id" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "Invalid client_id")
	}
}

// OAuth2 tests.

func TestGetDeviceCode(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oauth/pin" {
			t.Errorf("path = %q, want /oauth/pin", r.URL.Path)
		}
		if got := r.Header.Get("simkl-api-key"); got != "cid" {
			t.Errorf("simkl-api-key = %q, want %q", got, "cid")
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["client_id"] != "cid" {
			t.Errorf("client_id = %q, want %q", body["client_id"], "cid")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(simkl.DeviceCode{
			Result:          "OK",
			DeviceCode:      "dev-code",
			UserCode:        "ABC123",
			VerificationURL: "https://simkl.com/pin/ABC123",
			ExpiresIn:       900,
			Interval:        5,
		})
	}))
	defer ts.Close()

	c := simkl.New("cid", simkl.WithBaseURL(ts.URL))
	dc, err := c.GetDeviceCode(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if dc.UserCode != "ABC123" {
		t.Errorf("UserCode = %q, want %q", dc.UserCode, "ABC123")
	}
	if dc.DeviceCode != "dev-code" {
		t.Errorf("DeviceCode = %q, want %q", dc.DeviceCode, "dev-code")
	}
}

func TestExchangeCode(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" {
			t.Errorf("path = %q, want /oauth/token", r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["grant_type"] != "authorization_code" {
			t.Errorf("grant_type = %q, want authorization_code", body["grant_type"])
		}
		if body["code"] != "auth-code" {
			t.Errorf("code = %q, want auth-code", body["code"])
		}
		if body["client_id"] != "cid" {
			t.Errorf("client_id = %q, want cid", body["client_id"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": "simkl-access-tok",
			"token_type":   "Bearer",
		})
	}))
	defer ts.Close()

	var callbackToken string
	c := simkl.New("cid",
		simkl.WithBaseURL(ts.URL),
		simkl.WithClientSecret("secret"),
		simkl.WithTokenCallback(func(tok string) { callbackToken = tok }),
	)
	tok, err := c.ExchangeCode(context.Background(), "auth-code", "urn:ietf:wg:oauth:2.0:oob")
	if err != nil {
		t.Fatal(err)
	}
	if tok != "simkl-access-tok" {
		t.Errorf("token = %q, want %q", tok, "simkl-access-tok")
	}
	if callbackToken != "simkl-access-tok" {
		t.Errorf("callback token = %q, want %q", callbackToken, "simkl-access-tok")
	}
}

func TestBearerTokenInHeader(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer my-tok" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer my-tok")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(simkl.Movie{Title: "Test"})
	}))
	defer ts.Close()

	c := simkl.New("cid", simkl.WithBaseURL(ts.URL), simkl.WithAccessToken("my-tok"))
	_, err := c.GetMovie(context.Background(), "1")
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
		json.NewEncoder(w).Encode(simkl.Movie{Title: "Test"})
	}))
	defer ts.Close()

	c := simkl.New("cid", simkl.WithBaseURL(ts.URL))
	_, err := c.GetMovie(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAPIErrorRawBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}))
	defer ts.Close()

	c := simkl.New("key", simkl.WithBaseURL(ts.URL))
	_, err := c.GetMovie(context.Background(), "1")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *simkl.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.RawBody != "internal error" {
		t.Errorf("RawBody = %q, want %q", apiErr.RawBody, "internal error")
	}
}

func TestAnimeGenres(t *testing.T) {
	items := []simkl.GenreItem{{Title: "Action Anime", Year: 2024}}
	ts := newTestServer(t, "/anime/genres/action", "ag-key", items)
	defer ts.Close()

	c := simkl.New("ag-key", simkl.WithBaseURL(ts.URL))
	result, err := c.AnimeGenres(context.Background(), "action", 1, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("len = %d, want 1", len(result))
	}
}

func TestWithCustomHTTPClient(t *testing.T) {
	movie := simkl.Movie{Title: "Custom", Year: 2024, IDs: simkl.IDs{Simkl: 1}}
	ts := newTestServer(t, "/movies/1", "custom-key", movie)
	defer ts.Close()

	custom := &http.Client{}
	c := simkl.New("custom-key", simkl.WithBaseURL(ts.URL), simkl.WithHTTPClient(custom))
	m, err := c.GetMovie(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if m.Title != "Custom" {
		t.Errorf("Title = %q, want %q", m.Title, "Custom")
	}
}
