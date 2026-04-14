package rawg_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/game/rawg"
)

func setup(t *testing.T, handler http.HandlerFunc) *rawg.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return rawg.New("test-key", metadata.WithBaseURL(srv.URL))
}

func TestSearchGames(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Errorf("key = %q, want test-key", got)
		}
		if got := r.URL.Query().Get("search"); got != "zelda" {
			t.Errorf("search = %q, want zelda", got)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("page = %q, want 1", got)
		}
		if got := r.URL.Query().Get("page_size"); got != "10" {
			t.Errorf("page_size = %q, want 10", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count": 50,
			"next":  "https://api.rawg.io/api/games?search=zelda&page=2",
			"results": []map[string]any{
				{"id": 22511, "name": "The Legend of Zelda: Breath of the Wild", "slug": "the-legend-of-zelda-breath-of-the-wild", "rating": 4.42},
			},
		})
	})

	result, err := c.SearchGames(context.Background(), "zelda", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Count != 50 {
		t.Errorf("Count = %d, want 50", result.Count)
	}
	if len(result.Results) != 1 {
		t.Fatalf("len(Results) = %d, want 1", len(result.Results))
	}
	if result.Results[0].Name != "The Legend of Zelda: Breath of the Wild" {
		t.Errorf("Name = %q", result.Results[0].Name)
	}
	if result.Results[0].Rating != 4.42 {
		t.Errorf("Rating = %f, want 4.42", result.Results[0].Rating)
	}
}

func TestGetGame(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games/3498" {
			t.Errorf("path = %q, want /games/3498", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":              3498,
			"name":            "Grand Theft Auto V",
			"slug":            "grand-theft-auto-v",
			"description_raw": "A huge open world game.",
			"metacritic":      97,
			"website":         "https://www.rockstargames.com/V/",
		})
	})

	game, err := c.GetGame(context.Background(), 3498)
	if err != nil {
		t.Fatal(err)
	}
	if game.Name != "Grand Theft Auto V" {
		t.Errorf("Name = %q, want Grand Theft Auto V", game.Name)
	}
	if game.Metacritic != 97 {
		t.Errorf("Metacritic = %d, want 97", game.Metacritic)
	}
	if game.DescriptionRaw != "A huge open world game." {
		t.Errorf("DescriptionRaw = %q", game.DescriptionRaw)
	}
}

func TestGetGameBySlug(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games/grand-theft-auto-v" {
			t.Errorf("path = %q, want /games/grand-theft-auto-v", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":   3498,
			"name": "Grand Theft Auto V",
			"slug": "grand-theft-auto-v",
		})
	})

	game, err := c.GetGameBySlug(context.Background(), "grand-theft-auto-v")
	if err != nil {
		t.Fatal(err)
	}
	if game.ID != 3498 {
		t.Errorf("ID = %d, want 3498", game.ID)
	}
}

func TestGetGameScreenshots(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games/3498/screenshots" {
			t.Errorf("path = %q, want /games/3498/screenshots", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count": 6,
			"results": []map[string]any{
				{"id": 1, "image": "https://example.com/ss1.jpg", "width": 1920, "height": 1080},
				{"id": 2, "image": "https://example.com/ss2.jpg", "width": 1920, "height": 1080},
			},
		})
	})

	result, err := c.GetGameScreenshots(context.Background(), 3498, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Count != 6 {
		t.Errorf("Count = %d, want 6", result.Count)
	}
	if len(result.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(result.Results))
	}
	if result.Results[0].Width != 1920 {
		t.Errorf("Width = %d, want 1920", result.Results[0].Width)
	}
}

func TestGetGameTrailers(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games/3498/movies" {
			t.Errorf("path = %q, want /games/3498/movies", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count": 1,
			"results": []map[string]any{
				{"id": 1, "name": "Trailer", "preview": map[string]any{"max": "https://example.com/preview.jpg"}, "data": map[string]any{"max": "https://example.com/video.mp4"}},
			},
		})
	})

	result, err := c.GetGameTrailers(context.Background(), 3498)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Results) != 1 {
		t.Fatalf("len(Results) = %d, want 1", len(result.Results))
	}
	if result.Results[0].Data.Max != "https://example.com/video.mp4" {
		t.Errorf("Data.Max = %q", result.Results[0].Data.Max)
	}
}

func TestGetGameAdditions(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games/3498/additions" {
			t.Errorf("path = %q, want /games/3498/additions", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count":   2,
			"results": []map[string]any{{"id": 100, "name": "DLC Pack"}},
		})
	})

	result, err := c.GetGameAdditions(context.Background(), 3498, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Count != 2 {
		t.Errorf("Count = %d, want 2", result.Count)
	}
	if result.Results[0].Name != "DLC Pack" {
		t.Errorf("Name = %q, want DLC Pack", result.Results[0].Name)
	}
}

func TestGetGameSeries(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games/3498/game-series" {
			t.Errorf("path = %q, want /games/3498/game-series", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count":   3,
			"results": []map[string]any{{"id": 801, "name": "Grand Theft Auto IV"}},
		})
	})

	result, err := c.GetGameSeries(context.Background(), 3498, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Count != 3 {
		t.Errorf("Count = %d, want 3", result.Count)
	}
}

func TestGetPlatforms(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platforms" {
			t.Errorf("path = %q, want /platforms", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count": 50,
			"results": []map[string]any{
				{"id": 4, "name": "PC", "slug": "pc", "games_count": 500000},
				{"id": 187, "name": "PlayStation 5", "slug": "playstation5", "games_count": 1000},
			},
		})
	})

	result, err := c.GetPlatforms(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Results) != 2 {
		t.Fatalf("len(Results) = %d, want 2", len(result.Results))
	}
	if result.Results[0].Name != "PC" {
		t.Errorf("Name = %q, want PC", result.Results[0].Name)
	}
}

func TestGetPlatform(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platforms/4" {
			t.Errorf("path = %q, want /platforms/4", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id": 4, "name": "PC", "slug": "pc", "games_count": 500000,
		})
	})

	p, err := c.GetPlatform(context.Background(), 4)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "PC" {
		t.Errorf("Name = %q, want PC", p.Name)
	}
	if p.GamesCount != 500000 {
		t.Errorf("GamesCount = %d, want 500000", p.GamesCount)
	}
}

func TestGetGenres(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/genres" {
			t.Errorf("path = %q, want /genres", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count": 19,
			"results": []map[string]any{
				{"id": 4, "name": "Action", "slug": "action", "games_count": 170000},
			},
		})
	})

	result, err := c.GetGenres(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.Count != 19 {
		t.Errorf("Count = %d, want 19", result.Count)
	}
	if result.Results[0].Name != "Action" {
		t.Errorf("Name = %q, want Action", result.Results[0].Name)
	}
}

func TestGetPublishers(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/publishers" {
			t.Errorf("path = %q, want /publishers", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count":   100,
			"results": []map[string]any{{"id": 354, "name": "Nintendo", "slug": "nintendo"}},
		})
	})

	result, err := c.GetPublishers(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Results[0].Name != "Nintendo" {
		t.Errorf("Name = %q, want Nintendo", result.Results[0].Name)
	}
}

func TestGetDevelopers(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/developers" {
			t.Errorf("path = %q, want /developers", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count":   200,
			"results": []map[string]any{{"id": 405, "name": "CD PROJEKT RED", "slug": "cd-projekt-red"}},
		})
	})

	result, err := c.GetDevelopers(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Results[0].Name != "CD PROJEKT RED" {
		t.Errorf("Name = %q, want CD PROJEKT RED", result.Results[0].Name)
	}
}

func TestGetTags(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tags" {
			t.Errorf("path = %q, want /tags", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count":   400,
			"results": []map[string]any{{"id": 31, "name": "Singleplayer", "slug": "singleplayer", "language": "eng"}},
		})
	})

	result, err := c.GetTags(context.Background(), 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if result.Results[0].Name != "Singleplayer" {
		t.Errorf("Name = %q, want Singleplayer", result.Results[0].Name)
	}
}

func TestGetStores(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stores" {
			t.Errorf("path = %q, want /stores", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count": 10,
			"results": []map[string]any{
				{"id": 1, "name": "Steam", "slug": "steam", "domain": "store.steampowered.com", "games_count": 80000},
			},
		})
	})

	result, err := c.GetStores(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.Count != 10 {
		t.Errorf("Count = %d, want 10", result.Count)
	}
	if result.Results[0].Domain != "store.steampowered.com" {
		t.Errorf("Domain = %q, want store.steampowered.com", result.Results[0].Domain)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("invalid key"))
	}))
	defer srv.Close()

	c := rawg.New("bad-key", metadata.WithBaseURL(srv.URL))
	_, err := c.SearchGames(context.Background(), "zelda", 1, 10)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *rawg.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusForbidden)
	}
	if apiErr.Body != "invalid key" {
		t.Errorf("Body = %q, want invalid key", apiErr.Body)
	}
}

func TestAPIErrorString(t *testing.T) {
	t.Parallel()

	e := &rawg.APIError{StatusCode: 403, Status: "403 Forbidden", Body: "invalid key"}
	want := "rawg: 403 Forbidden: invalid key"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAPIKeyInQueryParams(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Errorf("key = %q, want test-key", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"count": 0, "results": []any{}})
	})

	_, err := c.GetGenres(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	custom := &http.Client{}
	c := rawg.New("key", metadata.WithHTTPClient(custom))
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestEndpointPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		wantPath string
		call     func(*rawg.Client) error
	}{
		{"SearchGames", "/games", func(c *rawg.Client) error {
			_, err := c.SearchGames(context.Background(), "test", 1, 10)
			return err
		}},
		{"GetGame", "/games/1", func(c *rawg.Client) error {
			_, err := c.GetGame(context.Background(), 1)
			return err
		}},
		{"GetPlatforms", "/platforms", func(c *rawg.Client) error {
			_, err := c.GetPlatforms(context.Background(), 1, 10)
			return err
		}},
		{"GetGenres", "/genres", func(c *rawg.Client) error {
			_, err := c.GetGenres(context.Background())
			return err
		}},
		{"GetPublishers", "/publishers", func(c *rawg.Client) error {
			_, err := c.GetPublishers(context.Background(), 1, 10)
			return err
		}},
		{"GetDevelopers", "/developers", func(c *rawg.Client) error {
			_, err := c.GetDevelopers(context.Background(), 1, 10)
			return err
		}},
		{"GetTags", "/tags", func(c *rawg.Client) error {
			_, err := c.GetTags(context.Background(), 1, 10)
			return err
		}},
		{"GetStores", "/stores", func(c *rawg.Client) error {
			_, err := c.GetStores(context.Background())
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := setup(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.wantPath {
					t.Errorf("path = %q, want %q", r.URL.Path, tt.wantPath)
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"count": 0, "results": []any{}})
			})
			if err := tt.call(c); err != nil {
				t.Fatal(err)
			}
		})
	}
}
