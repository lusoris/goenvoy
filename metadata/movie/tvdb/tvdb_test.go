package tvdb_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/metadata/movie/tvdb"
)

// loginHandler handles the /login endpoint in the test server.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"data":   map[string]any{"token": "test-jwt-token"},
	})
}

func newTestServer(t *testing.T, wantPath string, response any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle login requests.
		if r.URL.Path == "/login" {
			loginHandler(w, r)
			return
		}

		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-jwt-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-jwt-token")
		}
		if wantPath != "" && r.URL.RequestURI() != wantPath {
			t.Errorf("path = %q, want %q", r.URL.RequestURI(), wantPath)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
}

func newClient(t *testing.T, srv *httptest.Server) *tvdb.Client {
	t.Helper()
	return tvdb.New("test-api-key", tvdb.WithBaseURL(srv.URL))
}

func TestNew(t *testing.T) {
	t.Parallel()
	c := tvdb.New("api-key")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/login" {
			t.Errorf("path = %s, want /login", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var body tvdb.LoginRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.APIKey != "test-key" {
			t.Errorf("apikey = %q, want %q", body.APIKey, "test-key")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status": "success",
			"data":   map[string]any{"token": "jwt-123"},
		})
	}))
	defer srv.Close()

	c := tvdb.New("test-key", tvdb.WithBaseURL(srv.URL))
	if err := c.Login(context.Background()); err != nil {
		t.Fatalf("Login: %v", err)
	}
}

func TestGetSeries(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data": map[string]any{
			"id":   float64(81189),
			"name": "Breaking Bad",
			"slug": "breaking-bad",
		},
	}
	srv := newTestServer(t, "/series/81189", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetSeries(context.Background(), 81189)
	if err != nil {
		t.Fatalf("GetSeries: %v", err)
	}
	if got.Name != "Breaking Bad" {
		t.Errorf("Name = %q, want %q", got.Name, "Breaking Bad")
	}
}

func TestGetSeriesExtended(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data": map[string]any{
			"id":       float64(81189),
			"name":     "Breaking Bad",
			"overview": "A chemistry teacher diagnosed with cancer.",
		},
	}
	srv := newTestServer(t, "/series/81189/extended", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetSeriesExtended(context.Background(), 81189)
	if err != nil {
		t.Fatalf("GetSeriesExtended: %v", err)
	}
	if got.Overview != "A chemistry teacher diagnosed with cancer." {
		t.Errorf("Overview = %q, want %q", got.Overview, "A chemistry teacher diagnosed with cancer.")
	}
}

func TestGetSeriesEpisodes(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data": map[string]any{
			"series":   map[string]any{"id": float64(81189), "name": "Breaking Bad"},
			"episodes": []any{map[string]any{"id": float64(1), "name": "Pilot"}},
		},
	}
	srv := newTestServer(t, "/series/81189/episodes/default?page=0", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetSeriesEpisodes(context.Background(), 81189, "default", 0)
	if err != nil {
		t.Fatalf("GetSeriesEpisodes: %v", err)
	}
	if len(got.Episodes) != 1 {
		t.Fatalf("episodes count = %d, want 1", len(got.Episodes))
	}
	if got.Episodes[0].Name != "Pilot" {
		t.Errorf("episode name = %q, want %q", got.Episodes[0].Name, "Pilot")
	}
}

func TestGetSeriesTranslation(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data": map[string]any{
			"name":     "Breaking Bad",
			"language": "eng",
		},
	}
	srv := newTestServer(t, "/series/81189/translations/eng", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetSeriesTranslation(context.Background(), 81189, "eng")
	if err != nil {
		t.Fatalf("GetSeriesTranslation: %v", err)
	}
	if got.Language != "eng" {
		t.Errorf("Language = %q, want %q", got.Language, "eng")
	}
}

func TestGetMovie(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   map[string]any{"id": float64(1), "name": "Inception"},
	}
	srv := newTestServer(t, "/movies/1", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMovie(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMovie: %v", err)
	}
	if got.Name != "Inception" {
		t.Errorf("Name = %q, want %q", got.Name, "Inception")
	}
}

func TestGetMovieExtended(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   map[string]any{"id": float64(1), "name": "Inception", "budget": "160000000"},
	}
	srv := newTestServer(t, "/movies/1/extended", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMovieExtended(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMovieExtended: %v", err)
	}
	if got.Budget != "160000000" {
		t.Errorf("Budget = %q, want %q", got.Budget, "160000000")
	}
}

func TestGetMovieTranslation(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   map[string]any{"name": "Inception", "language": "eng"},
	}
	srv := newTestServer(t, "/movies/1/translations/eng", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetMovieTranslation(context.Background(), 1, "eng")
	if err != nil {
		t.Fatalf("GetMovieTranslation: %v", err)
	}
	if got.Name != "Inception" {
		t.Errorf("Name = %q, want %q", got.Name, "Inception")
	}
}

func TestGetEpisode(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   map[string]any{"id": float64(127396), "name": "Pilot"},
	}
	srv := newTestServer(t, "/episodes/127396", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetEpisode(context.Background(), 127396)
	if err != nil {
		t.Fatalf("GetEpisode: %v", err)
	}
	if got.Name != "Pilot" {
		t.Errorf("Name = %q, want %q", got.Name, "Pilot")
	}
}

func TestGetEpisodeExtended(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   map[string]any{"id": float64(127396), "name": "Pilot", "productionCode": "1001"},
	}
	srv := newTestServer(t, "/episodes/127396/extended", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetEpisodeExtended(context.Background(), 127396)
	if err != nil {
		t.Fatalf("GetEpisodeExtended: %v", err)
	}
	if got.ProductionCode != "1001" {
		t.Errorf("ProductionCode = %q, want %q", got.ProductionCode, "1001")
	}
}

func TestGetSeason(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   map[string]any{"id": float64(100), "number": float64(1)},
	}
	srv := newTestServer(t, "/seasons/100", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetSeason(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetSeason: %v", err)
	}
	if got.Number != 1 {
		t.Errorf("Number = %d, want 1", got.Number)
	}
}

func TestGetSeasonExtended(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data": map[string]any{
			"id":       float64(100),
			"number":   float64(1),
			"episodes": []any{map[string]any{"id": float64(1), "name": "S01E01"}},
		},
	}
	srv := newTestServer(t, "/seasons/100/extended", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetSeasonExtended(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetSeasonExtended: %v", err)
	}
	if len(got.Episodes) != 1 {
		t.Fatalf("episodes count = %d, want 1", len(got.Episodes))
	}
}

func TestGetPerson(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   map[string]any{"id": float64(250930), "name": "Bryan Cranston"},
	}
	srv := newTestServer(t, "/people/250930", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetPerson(context.Background(), 250930)
	if err != nil {
		t.Fatalf("GetPerson: %v", err)
	}
	if got.Name != "Bryan Cranston" {
		t.Errorf("Name = %q, want %q", got.Name, "Bryan Cranston")
	}
}

func TestGetPersonExtended(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   map[string]any{"id": float64(250930), "name": "Bryan Cranston", "birthPlace": "Hollywood"},
	}
	srv := newTestServer(t, "/people/250930/extended", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetPersonExtended(context.Background(), 250930)
	if err != nil {
		t.Fatalf("GetPersonExtended: %v", err)
	}
	if got.BirthPlace != "Hollywood" {
		t.Errorf("BirthPlace = %q, want %q", got.BirthPlace, "Hollywood")
	}
}

func TestGetArtwork(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   map[string]any{"id": float64(42), "image": "https://artworks.thetvdb.com/poster.jpg"},
	}
	srv := newTestServer(t, "/artwork/42", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetArtwork(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetArtwork: %v", err)
	}
	if got.Image != "https://artworks.thetvdb.com/poster.jpg" {
		t.Errorf("Image = %q, want %q", got.Image, "https://artworks.thetvdb.com/poster.jpg")
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   []any{map[string]any{"name": "Breaking Bad", "type": "series", "tvdb_id": "81189"}},
	}
	srv := newTestServer(t, "", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.Search(context.Background(), "Breaking Bad", &tvdb.SearchParams{Type: "series"})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("results = %d, want 1", len(got))
	}
	if got[0].Name != "Breaking Bad" {
		t.Errorf("Name = %q, want %q", got[0].Name, "Breaking Bad")
	}
}

func TestSearchByRemoteID(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data": []any{map[string]any{
			"series": map[string]any{"id": float64(81189), "name": "Breaking Bad"},
		}},
	}
	srv := newTestServer(t, "/search/remoteid/tt0903747", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.SearchByRemoteID(context.Background(), "tt0903747")
	if err != nil {
		t.Fatalf("SearchByRemoteID: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("results = %d, want 1", len(got))
	}
	if got[0].Series == nil {
		t.Fatal("expected series result")
	}
	if got[0].Series.Name != "Breaking Bad" {
		t.Errorf("Name = %q, want %q", got[0].Series.Name, "Breaking Bad")
	}
}

func TestGetGenres(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   []any{map[string]any{"id": float64(1), "name": "Drama", "slug": "drama"}},
	}
	srv := newTestServer(t, "/genres", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetGenres(context.Background())
	if err != nil {
		t.Fatalf("GetGenres: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Drama" {
		t.Errorf("genres = %+v, want [{Drama}]", got)
	}
}

func TestGetLanguages(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   []any{map[string]any{"id": "eng", "name": "English", "shortCode": "en"}},
	}
	srv := newTestServer(t, "/languages", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetLanguages(context.Background())
	if err != nil {
		t.Fatalf("GetLanguages: %v", err)
	}
	if len(got) != 1 || got[0].Name != "English" {
		t.Errorf("languages = %+v, want [{English}]", got)
	}
}

func TestGetContentRatings(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   []any{map[string]any{"id": float64(1), "name": "TV-MA", "country": "usa"}},
	}
	srv := newTestServer(t, "/content/ratings", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetContentRatings(context.Background())
	if err != nil {
		t.Fatalf("GetContentRatings: %v", err)
	}
	if len(got) != 1 || got[0].Name != "TV-MA" {
		t.Errorf("ratings = %+v, want [{TV-MA}]", got)
	}
}

func TestGetUpdates(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data": []any{map[string]any{
			"recordId": float64(81189), "method": "update", "methodInt": float64(2), "timeStamp": float64(1700000000),
		}},
	}
	srv := newTestServer(t, "", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetUpdates(context.Background(), 1700000000, nil)
	if err != nil {
		t.Fatalf("GetUpdates: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("updates = %d, want 1", len(got))
	}
	if got[0].RecordID != 81189 {
		t.Errorf("RecordID = %d, want 81189", got[0].RecordID)
	}
}

func TestGetCharacter(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"status": "success",
		"data":   map[string]any{"id": float64(100), "name": "Walter White"},
	}
	srv := newTestServer(t, "/characters/100", want)
	defer srv.Close()

	c := newClient(t, srv)
	got, err := c.GetCharacter(context.Background(), 100)
	if err != nil {
		t.Fatalf("GetCharacter: %v", err)
	}
	if got.Name != "Walter White" {
		t.Errorf("Name = %q, want %q", got.Name, "Walter White")
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
			loginHandler(w, r)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"message": "not found"})
	}))
	defer srv.Close()

	c := newClient(t, srv)
	_, err := c.GetSeries(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *tvdb.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
}

func TestLoginError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"message": "invalid credentials"})
	}))
	defer srv.Close()

	c := tvdb.New("bad-key", tvdb.WithBaseURL(srv.URL))
	err := c.Login(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *tvdb.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T: %v", err, err)
	}
}

func TestTokenRefreshOn401(t *testing.T) {
	t.Parallel()
	loginCount := 0
	requestCount := 0

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
			loginCount++
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"status": "success",
				"data":   map[string]any{"token": "jwt-token-" + string(rune('0'+loginCount))},
			})
			return
		}
		requestCount++
		// First data request returns 401 (expired token), second succeeds.
		if requestCount == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]any{"message": "token expired"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status": "success",
			"data":   map[string]any{"id": float64(1), "name": "Test Series"},
		})
	}))
	defer srv.Close()

	c := tvdb.New("test-key", tvdb.WithBaseURL(srv.URL))
	series, err := c.GetSeries(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetSeries: %v", err)
	}
	if series.Name != "Test Series" {
		t.Errorf("Name = %q, want %q", series.Name, "Test Series")
	}
	// Should have logged in twice: once for initial token, once after 401.
	if loginCount != 2 {
		t.Errorf("loginCount = %d, want 2", loginCount)
	}
	// Should have made 2 data requests: first failed (401), second succeeded.
	if requestCount != 2 {
		t.Errorf("requestCount = %d, want 2", requestCount)
	}
}

func TestTokenRefreshOn401LoginFails(t *testing.T) {
	t.Parallel()
	loginCount := 0

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
			loginCount++
			if loginCount == 1 {
				// First login succeeds.
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{
					"status": "success",
					"data":   map[string]any{"token": "jwt-1"},
				})
				return
			}
			// Second login fails.
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]any{"message": "invalid"})
			return
		}
		// Always return 401 for data requests.
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"message": "token expired"})
	}))
	defer srv.Close()

	c := tvdb.New("test-key", tvdb.WithBaseURL(srv.URL))
	_, err := c.GetSeries(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error when re-login fails")
	}
}
