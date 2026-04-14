package igdb_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/game/igdb"
)

func setup(t *testing.T, handler http.HandlerFunc) *igdb.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return igdb.New("test-client-id", "test-token", metadata.WithBaseURL(srv.URL))
}

func TestSearchGames(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Client-Id"); got != "test-client-id" {
			t.Errorf("Client-ID = %q, want test-client-id", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q, want Bearer test-token", got)
		}
		body, _ := io.ReadAll(r.Body)
		if got := string(body); got != `search "zelda"; fields *; limit 5;` {
			t.Errorf("body = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1025, "name": "The Legend of Zelda", "slug": "the-legend-of-zelda"},
		})
	})

	games, err := c.SearchGames(context.Background(), "zelda", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fatalf("len = %d, want 1", len(games))
	}
	if games[0].Name != "The Legend of Zelda" {
		t.Errorf("Name = %q, want The Legend of Zelda", games[0].Name)
	}
	if games[0].ID != 1025 {
		t.Errorf("ID = %d, want 1025", games[0].ID)
	}
}

func TestGetGame(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games" {
			t.Errorf("path = %q, want /games", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if got := string(body); got != "fields *; where id = 1942;" {
			t.Errorf("body = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1942, "name": "The Witcher 3", "slug": "the-witcher-3", "rating": 92.5},
		})
	})

	game, err := c.GetGame(context.Background(), 1942)
	if err != nil {
		t.Fatal(err)
	}
	if game == nil {
		t.Fatal("expected non-nil game")
	}
	if game.Name != "The Witcher 3" {
		t.Errorf("Name = %q, want The Witcher 3", game.Name)
	}
	if game.Rating != 92.5 {
		t.Errorf("Rating = %f, want 92.5", game.Rating)
	}
}

func TestGetGameNotFound(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	})

	game, err := c.GetGame(context.Background(), 99999)
	if err != nil {
		t.Fatal(err)
	}
	if game != nil {
		t.Errorf("expected nil game, got %+v", game)
	}
}

func TestGetGames(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if got := string(body); got != "fields *; where id = (1,2,3);" {
			t.Errorf("body = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "Game One"},
			{"id": 2, "name": "Game Two"},
			{"id": 3, "name": "Game Three"},
		})
	})

	games, err := c.GetGames(context.Background(), []int{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 3 {
		t.Fatalf("len = %d, want 3", len(games))
	}
	if games[1].Name != "Game Two" {
		t.Errorf("Name = %q, want Game Two", games[1].Name)
	}
}

func TestGetPopularGames(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if got := string(body); got != "fields *; sort total_rating desc; where total_rating_count > 5; limit 10;" {
			t.Errorf("body = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1942, "name": "The Witcher 3", "total_rating": 95.0},
		})
	})

	games, err := c.GetPopularGames(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fatalf("len = %d, want 1", len(games))
	}
	if games[0].TotalRating != 95.0 {
		t.Errorf("TotalRating = %f, want 95.0", games[0].TotalRating)
	}
}

func TestGetPlatform(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platforms" {
			t.Errorf("path = %q, want /platforms", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 48, "name": "PlayStation 4", "slug": "ps4", "abbreviation": "PS4", "generation": 8},
		})
	})

	p, err := c.GetPlatform(context.Background(), 48)
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("expected non-nil platform")
	}
	if p.Name != "PlayStation 4" {
		t.Errorf("Name = %q, want PlayStation 4", p.Name)
	}
	if p.Generation != 8 {
		t.Errorf("Generation = %d, want 8", p.Generation)
	}
}

func TestGetPlatformNotFound(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	})

	p, err := c.GetPlatform(context.Background(), 99999)
	if err != nil {
		t.Fatal(err)
	}
	if p != nil {
		t.Errorf("expected nil platform, got %+v", p)
	}
}

func TestGetPlatforms(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if got := string(body); got != "fields *; limit 10; offset 0;" {
			t.Errorf("body = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 48, "name": "PlayStation 4"},
			{"id": 49, "name": "Xbox One"},
		})
	})

	platforms, err := c.GetPlatforms(context.Background(), 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(platforms) != 2 {
		t.Fatalf("len = %d, want 2", len(platforms))
	}
}

func TestGetGenre(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 12, "name": "Role-playing (RPG)", "slug": "role-playing-rpg"},
		})
	})

	g, err := c.GetGenre(context.Background(), 12)
	if err != nil {
		t.Fatal(err)
	}
	if g == nil {
		t.Fatal("expected non-nil genre")
	}
	if g.Name != "Role-playing (RPG)" {
		t.Errorf("Name = %q, want Role-playing (RPG)", g.Name)
	}
}

func TestGetGenreNotFound(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	})

	g, err := c.GetGenre(context.Background(), 99999)
	if err != nil {
		t.Fatal(err)
	}
	if g != nil {
		t.Errorf("expected nil genre, got %+v", g)
	}
}

func TestGetGenres(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 12, "name": "Role-playing (RPG)"},
			{"id": 31, "name": "Adventure"},
		})
	})

	genres, err := c.GetGenres(context.Background(), 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(genres) != 2 {
		t.Fatalf("len = %d, want 2", len(genres))
	}
}

func TestGetCompany(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 70, "name": "Nintendo", "slug": "nintendo", "country": 392},
		})
	})

	co, err := c.GetCompany(context.Background(), 70)
	if err != nil {
		t.Fatal(err)
	}
	if co == nil {
		t.Fatal("expected non-nil company")
	}
	if co.Name != "Nintendo" {
		t.Errorf("Name = %q, want Nintendo", co.Name)
	}
	if co.Country != 392 {
		t.Errorf("Country = %d, want 392", co.Country)
	}
}

func TestGetCompanyNotFound(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	})

	co, err := c.GetCompany(context.Background(), 99999)
	if err != nil {
		t.Fatal(err)
	}
	if co != nil {
		t.Errorf("expected nil company, got %+v", co)
	}
}

func TestSearchCompanies(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if got := string(body); got != `search "nintendo"; fields *; limit 5;` {
			t.Errorf("body = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 70, "name": "Nintendo"},
		})
	})

	companies, err := c.SearchCompanies(context.Background(), "nintendo", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(companies) != 1 {
		t.Fatalf("len = %d, want 1", len(companies))
	}
	if companies[0].Name != "Nintendo" {
		t.Errorf("Name = %q, want Nintendo", companies[0].Name)
	}
}

func TestGetGameCovers(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/covers" {
			t.Errorf("path = %q, want /covers", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if got := string(body); got != "fields *; where game = 1942;" {
			t.Errorf("body = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 100, "game": 1942, "image_id": "co1wyy", "width": 264, "height": 352},
		})
	})

	covers, err := c.GetGameCovers(context.Background(), 1942)
	if err != nil {
		t.Fatal(err)
	}
	if len(covers) != 1 {
		t.Fatalf("len = %d, want 1", len(covers))
	}
	if covers[0].ImageID != "co1wyy" {
		t.Errorf("ImageID = %q, want co1wyy", covers[0].ImageID)
	}
	if covers[0].Width != 264 {
		t.Errorf("Width = %d, want 264", covers[0].Width)
	}
}

func TestGetGameScreenshots(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/screenshots" {
			t.Errorf("path = %q, want /screenshots", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 200, "game": 1942, "image_id": "sc6m5p", "width": 1920, "height": 1080},
			{"id": 201, "game": 1942, "image_id": "sc6m5q", "width": 1920, "height": 1080},
		})
	})

	screenshots, err := c.GetGameScreenshots(context.Background(), 1942)
	if err != nil {
		t.Fatal(err)
	}
	if len(screenshots) != 2 {
		t.Fatalf("len = %d, want 2", len(screenshots))
	}
	if screenshots[0].ImageID != "sc6m5p" {
		t.Errorf("ImageID = %q, want sc6m5p", screenshots[0].ImageID)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("unauthorized"))
	}))
	defer srv.Close()

	c := igdb.New("bad-id", "bad-token", metadata.WithBaseURL(srv.URL))
	_, err := c.SearchGames(context.Background(), "zelda", 10)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *igdb.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
	if apiErr.Body != "unauthorized" {
		t.Errorf("Body = %q, want unauthorized", apiErr.Body)
	}
}

func TestAPIErrorString(t *testing.T) {
	t.Parallel()

	e := &igdb.APIError{StatusCode: 401, Status: "401 Unauthorized", Body: "bad token"}
	want := "igdb: 401 Unauthorized: bad token"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestRequestHeaders(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Client-Id"); got != "test-client-id" {
			t.Errorf("Client-ID = %q, want test-client-id", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q, want Bearer test-token", got)
		}
		if got := r.Header.Get("Content-Type"); got != "text/plain" {
			t.Errorf("Content-Type = %q, want text/plain", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	})

	_, err := c.SearchGames(context.Background(), "test", 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEndpointPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		wantPath string
		call     func(*igdb.Client) error
	}{
		{"SearchGames", "/games", func(c *igdb.Client) error {
			_, err := c.SearchGames(context.Background(), "test", 1)
			return err
		}},
		{"GetPlatforms", "/platforms", func(c *igdb.Client) error {
			_, err := c.GetPlatforms(context.Background(), 10, 0)
			return err
		}},
		{"GetGenres", "/genres", func(c *igdb.Client) error {
			_, err := c.GetGenres(context.Background(), 10, 0)
			return err
		}},
		{"GetCompany", "/companies", func(c *igdb.Client) error {
			_, err := c.GetCompany(context.Background(), 1)
			return err
		}},
		{"GetGameCovers", "/covers", func(c *igdb.Client) error {
			_, err := c.GetGameCovers(context.Background(), 1)
			return err
		}},
		{"GetGameScreenshots", "/screenshots", func(c *igdb.Client) error {
			_, err := c.GetGameScreenshots(context.Background(), 1)
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
				_, _ = w.Write([]byte("[]"))
			})
			if err := tt.call(c); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	custom := &http.Client{}
	c := igdb.New("id", "token", metadata.WithHTTPClient(custom))
	// Verify client was created without error.
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}
