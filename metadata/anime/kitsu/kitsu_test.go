package kitsu

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/metadata"
)

// respondResource writes a JSON:API single-resource envelope.
func respondResource(w http.ResponseWriter, id, typ string, attrs any) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	env := map[string]any{
		"data": map[string]any{
			"id":         id,
			"type":       typ,
			"attributes": attrs,
		},
	}
	json.NewEncoder(w).Encode(env)
}

// respondCollection writes a JSON:API collection envelope.
func respondCollection(w http.ResponseWriter, typ string, items []map[string]any) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	data := make([]map[string]any, len(items))
	for i, item := range items {
		data[i] = map[string]any{
			"id":         item["id"],
			"type":       typ,
			"attributes": item,
		}
	}
	env := map[string]any{
		"data":  data,
		"links": map[string]string{},
	}
	json.NewEncoder(w).Encode(env)
}

func setup(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return New(metadata.WithBaseURL(ts.URL))
}

func TestGetAnime(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondResource(w, "1", "anime", map[string]any{
			"slug":           "cowboy-bebop",
			"canonicalTitle": "Cowboy Bebop",
			"subtype":        "TV",
			"status":         "finished",
			"episodeCount":   26,
			"averageRating":  "82.25",
			"startDate":      "1998-04-03",
			"nsfw":           false,
		})
	})

	anime, err := c.GetAnime(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetAnime() error: %v", err)
	}
	if anime.ID != "1" {
		t.Errorf("ID = %q, want %q", anime.ID, "1")
	}
	if anime.Slug != "cowboy-bebop" {
		t.Errorf("Slug = %q, want %q", anime.Slug, "cowboy-bebop")
	}
	if anime.CanonicalTitle != "Cowboy Bebop" {
		t.Errorf("CanonicalTitle = %q, want %q", anime.CanonicalTitle, "Cowboy Bebop")
	}
	if anime.Subtype != "TV" {
		t.Errorf("Subtype = %q, want %q", anime.Subtype, "TV")
	}
	if anime.Status != "finished" {
		t.Errorf("Status = %q, want %q", anime.Status, "finished")
	}
}

func TestSearchAnime(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("filter[text]") == "" {
			http.Error(w, "missing filter", http.StatusBadRequest)
			return
		}
		respondCollection(w, "anime", []map[string]any{
			{"id": "1", "slug": "cowboy-bebop", "canonicalTitle": "Cowboy Bebop"},
			{"id": "2", "slug": "cowboy-bebop-tengoku-no-tobira", "canonicalTitle": "Cowboy Bebop: Knockin' on Heaven's Door"},
		})
	})

	results, err := c.SearchAnime(context.Background(), "cowboy bebop", 10, 0)
	if err != nil {
		t.Fatalf("SearchAnime() error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	if results[0].ID != "1" {
		t.Errorf("results[0].ID = %q, want %q", results[0].ID, "1")
	}
}

func TestTrendingAnime(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondCollection(w, "anime", []map[string]any{
			{"id": "100", "slug": "jjk-3", "canonicalTitle": "Jujutsu Kaisen Season 3"},
		})
	})

	results, err := c.TrendingAnime(context.Background())
	if err != nil {
		t.Fatalf("TrendingAnime() error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].ID != "100" {
		t.Errorf("results[0].ID = %q, want %q", results[0].ID, "100")
	}
}

func TestGetAnimeEpisodes(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondCollection(w, "episodes", []map[string]any{
			{
				"id":             "229115",
				"canonicalTitle": "Asteroid Blues",
				"seasonNumber":   1,
				"number":         1,
				"airdate":        "1998-10-23",
				"length":         25,
			},
		})
	})

	epis, err := c.GetAnimeEpisodes(context.Background(), 1, 20, 0)
	if err != nil {
		t.Fatalf("GetAnimeEpisodes() error: %v", err)
	}
	if len(epis) != 1 {
		t.Fatalf("len(epis) = %d, want 1", len(epis))
	}
	if epis[0].CanonicalTitle != "Asteroid Blues" {
		t.Errorf("CanonicalTitle = %q, want %q", epis[0].CanonicalTitle, "Asteroid Blues")
	}
}

func TestGetManga(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondResource(w, "42", "manga", map[string]any{
			"slug":           "guardian-dog",
			"canonicalTitle": "Guardian Dog",
			"subtype":        "manga",
			"status":         "finished",
			"chapterCount":   22,
			"volumeCount":    4,
		})
	})

	manga, err := c.GetManga(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetManga() error: %v", err)
	}
	if manga.ID != "42" {
		t.Errorf("ID = %q, want %q", manga.ID, "42")
	}
	if manga.Slug != "guardian-dog" {
		t.Errorf("Slug = %q, want %q", manga.Slug, "guardian-dog")
	}
	if manga.VolumeCount != 4 {
		t.Errorf("VolumeCount = %d, want 4", manga.VolumeCount)
	}
}

func TestSearchManga(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondCollection(w, "manga", []map[string]any{
			{"id": "10", "slug": "one-piece", "canonicalTitle": "One Piece"},
		})
	})

	results, err := c.SearchManga(context.Background(), "one piece", 10, 0)
	if err != nil {
		t.Fatalf("SearchManga() error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0].CanonicalTitle != "One Piece" {
		t.Errorf("CanonicalTitle = %q, want %q", results[0].CanonicalTitle, "One Piece")
	}
}

func TestTrendingManga(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondCollection(w, "manga", []map[string]any{
			{"id": "50", "slug": "solo-leveling", "canonicalTitle": "Solo Leveling"},
		})
	})

	results, err := c.TrendingManga(context.Background())
	if err != nil {
		t.Fatalf("TrendingManga() error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
}

func TestGetCharacter(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondResource(w, "7", "characters", map[string]any{
			"slug":          "jet-black",
			"name":          "Jet Black",
			"canonicalName": "Jet Black",
			"malId":         3,
			"description":   "Jet, known as the Black Dog.",
			"image":         map[string]any{"original": "https://example.com/jet.jpg"},
		})
	})

	ch, err := c.GetCharacter(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetCharacter() error: %v", err)
	}
	if ch.ID != "7" {
		t.Errorf("ID = %q, want %q", ch.ID, "7")
	}
	if ch.Name != "Jet Black" {
		t.Errorf("Name = %q, want %q", ch.Name, "Jet Black")
	}
	if ch.MalID != 3 {
		t.Errorf("MalID = %d, want 3", ch.MalID)
	}
}

func TestSearchCharacters(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondCollection(w, "characters", []map[string]any{
			{"id": "1", "name": "Jet Black", "slug": "jet-black"},
			{"id": "2", "name": "Spike Spiegel", "slug": "spike-spiegel"},
		})
	})

	results, err := c.SearchCharacters(context.Background(), "bebop", 10, 0)
	if err != nil {
		t.Fatalf("SearchCharacters() error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
}

func TestGetCategory(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondResource(w, "1", "categories", map[string]any{
			"title":           "Middle School",
			"slug":            "middle-school",
			"totalMediaCount": 111,
			"nsfw":            false,
			"childCount":      0,
		})
	})

	cat, err := c.GetCategory(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetCategory() error: %v", err)
	}
	if cat.ID != "1" {
		t.Errorf("ID = %q, want %q", cat.ID, "1")
	}
	if cat.Title != "Middle School" {
		t.Errorf("Title = %q, want %q", cat.Title, "Middle School")
	}
	if cat.TotalMediaCount != 111 {
		t.Errorf("TotalMediaCount = %d, want 111", cat.TotalMediaCount)
	}
}

func TestGetCategories(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondCollection(w, "categories", []map[string]any{
			{"id": "1", "title": "Middle School", "slug": "middle-school"},
			{"id": "2", "title": "High School", "slug": "high-school"},
		})
	})

	cats, err := c.GetCategories(context.Background(), 20, 0)
	if err != nil {
		t.Fatalf("GetCategories() error: %v", err)
	}
	if len(cats) != 2 {
		t.Fatalf("len(cats) = %d, want 2", len(cats))
	}
}

func TestGetUser(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondResource(w, "1", "users", map[string]any{
			"name":           "vikhyat",
			"slug":           "vikhyat",
			"followersCount": 1675,
			"followingCount": 1085,
			"about":          "",
			"status":         "registered",
		})
	})

	user, err := c.GetUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetUser() error: %v", err)
	}
	if user.ID != "1" {
		t.Errorf("ID = %q, want %q", user.ID, "1")
	}
	if user.Name != "vikhyat" {
		t.Errorf("Name = %q, want %q", user.Name, "vikhyat")
	}
	if user.FollowersCount != 1675 {
		t.Errorf("FollowersCount = %d, want 1675", user.FollowersCount)
	}
}

func TestSearchUsers(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondCollection(w, "users", []map[string]any{
			{"id": "1", "name": "vikhyat", "slug": "vikhyat"},
		})
	})

	results, err := c.SearchUsers(context.Background(), "vikhyat", 10, 0)
	if err != nil {
		t.Fatalf("SearchUsers() error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
}

func TestAPIError(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"title":"Record not found","status":"404"}]}`))
	})

	_, err := c.GetAnime(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusNotFound)
	}
}

func TestJSONAPIHeaders(t *testing.T) {
	var gotAccept string
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		gotAccept = r.Header.Get("Accept")
		respondResource(w, "1", "anime", map[string]any{"slug": "test"})
	})

	_, err := c.GetAnime(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAccept != jsonAPIMediaType {
		t.Errorf("Accept header = %q, want %q", gotAccept, jsonAPIMediaType)
	}
}

func TestWithUserAgent(t *testing.T) {
	var gotUA string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		respondResource(w, "1", "anime", map[string]any{"slug": "test"})
	}))
	t.Cleanup(ts.Close)

	c := New(metadata.WithBaseURL(ts.URL), metadata.WithUserAgent("test-agent/1.0"))
	_, err := c.GetAnime(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotUA != "test-agent/1.0" {
		t.Errorf("User-Agent = %q, want %q", gotUA, "test-agent/1.0")
	}
}

// OAuth2 tests.

func TestAuthenticate(t *testing.T) {
	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if ct != "application/x-www-form-urlencoded" {
			t.Errorf("Content-Type = %q, want application/x-www-form-urlencoded", ct)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.PostForm.Get("grant_type") != "password" {
			t.Errorf("grant_type = %q, want password", r.PostForm.Get("grant_type"))
		}
		if r.PostForm.Get("username") != "user@example.com" {
			t.Errorf("username = %q", r.PostForm.Get("username"))
		}
		if r.PostForm.Get("password") != "s3cret" {
			t.Errorf("password = %q", r.PostForm.Get("password"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Token{
			AccessToken:  "kitsu-access",
			TokenType:    "Bearer",
			ExpiresIn:    2592000,
			RefreshToken: "kitsu-refresh",
			Scope:        "public",
			CreatedAt:    1609459200,
		})
	}))
	t.Cleanup(authSrv.Close)

	var callbackToken Token
	c := New()
	c.SetTokenCallback(func(tok Token) { callbackToken = tok })
	c.authURL = authSrv.URL

	tok, err := c.Authenticate(context.Background(), "user@example.com", "s3cret")
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "kitsu-access" {
		t.Errorf("AccessToken = %q, want %q", tok.AccessToken, "kitsu-access")
	}
	if tok.RefreshToken != "kitsu-refresh" {
		t.Errorf("RefreshToken = %q, want %q", tok.RefreshToken, "kitsu-refresh")
	}
	if callbackToken.AccessToken != "kitsu-access" {
		t.Errorf("callback AccessToken = %q, want %q", callbackToken.AccessToken, "kitsu-access")
	}
}

func TestRefreshToken(t *testing.T) {
	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.PostForm.Get("grant_type") != "refresh_token" {
			t.Errorf("grant_type = %q, want refresh_token", r.PostForm.Get("grant_type"))
		}
		if r.PostForm.Get("refresh_token") != "old-rt" {
			t.Errorf("refresh_token = %q, want old-rt", r.PostForm.Get("refresh_token"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Token{
			AccessToken:  "new-access",
			RefreshToken: "new-refresh",
		})
	}))
	t.Cleanup(authSrv.Close)

	c := New()
	c.SetRefreshToken("old-rt")
	c.authURL = authSrv.URL

	tok, err := c.RefreshToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "new-access" {
		t.Errorf("AccessToken = %q, want %q", tok.AccessToken, "new-access")
	}
}

func TestRefreshTokenMissing(t *testing.T) {
	c := New()
	_, err := c.RefreshToken(context.Background())
	if err == nil {
		t.Fatal("expected error when no refresh token set")
	}
}

func TestBearerTokenInHeader(t *testing.T) {
	var gotAuth string
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		respondResource(w, "1", "anime", map[string]any{"slug": "test"})
	})
	c.mu.Lock()
	c.accessToken = "kitsu-tok"
	c.mu.Unlock()

	_, err := c.GetAnime(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer kitsu-tok" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer kitsu-tok")
	}
}

func TestAuthenticateError(t *testing.T) {
	authSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid_grant","error_description":"Invalid credentials"}`))
	}))
	t.Cleanup(authSrv.Close)

	c := New()
	c.authURL = authSrv.URL

	_, err := c.Authenticate(context.Background(), "bad", "wrong")
	if err == nil {
		t.Fatal("expected error for bad credentials")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
}
