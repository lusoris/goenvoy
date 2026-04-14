package mobygames_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/game/mobygames"
)

func setup(t *testing.T, handler http.HandlerFunc) *mobygames.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return mobygames.New("test-key", metadata.WithBaseURL(srv.URL))
}

func TestSearchGames(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if got := r.URL.Query().Get("api_key"); got != "test-key" {
			t.Errorf("api_key = %q, want test-key", got)
		}
		if got := r.URL.Query().Get("title"); got != "zelda" {
			t.Errorf("title = %q, want zelda", got)
		}
		if got := r.URL.Query().Get("format"); got != "normal" {
			t.Errorf("format = %q, want normal", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"games": []map[string]any{
				{"game_id": 123, "title": "The Legend of Zelda", "moby_score": 4.5},
			},
		})
	})

	games, err := c.SearchGames(context.Background(), "zelda", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fatalf("len(games) = %d, want 1", len(games))
	}
	if games[0].Title != "The Legend of Zelda" {
		t.Errorf("Title = %q", games[0].Title)
	}
	if games[0].MobyScore != 4.5 {
		t.Errorf("MobyScore = %f, want 4.5", games[0].MobyScore)
	}
}

func TestGetGame(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games/123" {
			t.Errorf("path = %q, want /games/123", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"game_id":     123,
			"title":       "The Legend of Zelda",
			"description": "An action-adventure game.",
			"moby_score":  4.5,
			"num_votes":   100,
		})
	})

	game, err := c.GetGame(context.Background(), 123)
	if err != nil {
		t.Fatal(err)
	}
	if game.Title != "The Legend of Zelda" {
		t.Errorf("Title = %q", game.Title)
	}
	if game.Description != "An action-adventure game." {
		t.Errorf("Description = %q", game.Description)
	}
}

func TestGetGenres(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/genres" {
			t.Errorf("path = %q, want /genres", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"genre_id": 1, "genre_name": "Action", "genre_category": "Basic Genres"},
		})
	})

	genres, err := c.GetGenres(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(genres) != 1 {
		t.Fatalf("len(genres) = %d, want 1", len(genres))
	}
	if genres[0].GenreName != "Action" {
		t.Errorf("GenreName = %q", genres[0].GenreName)
	}
}

func TestGetPlatforms(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platforms" {
			t.Errorf("path = %q, want /platforms", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"platform_id": 1, "platform_name": "PC"},
		})
	})

	platforms, err := c.GetPlatforms(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(platforms) != 1 {
		t.Fatalf("len(platforms) = %d, want 1", len(platforms))
	}
	if platforms[0].PlatformName != "PC" {
		t.Errorf("PlatformName = %q", platforms[0].PlatformName)
	}
}

func TestGetGameScreenshots(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games/123/platforms/1/screenshots" {
			t.Errorf("path = %q, want /games/123/platforms/1/screenshots", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"screenshots": []map[string]any{
				{"image": "https://example.com/ss.jpg", "width": 640, "height": 480, "caption": "Title screen"},
			},
		})
	})

	ss, err := c.GetGameScreenshots(context.Background(), 123, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) != 1 {
		t.Fatalf("len(screenshots) = %d, want 1", len(ss))
	}
	if ss[0].Caption != "Title screen" {
		t.Errorf("Caption = %q", ss[0].Caption)
	}
}

func TestGetGameCovers(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games/123/platforms/1/covers" {
			t.Errorf("path = %q, want /games/123/platforms/1/covers", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"cover_groups": []map[string]any{
				{"covers": []map[string]any{
					{"image": "https://example.com/cover.jpg", "width": 300, "height": 400, "scan_of": "Front Cover"},
				}},
			},
		})
	})

	groups, err := c.GetGameCovers(context.Background(), 123, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 {
		t.Fatalf("len(groups) = %d, want 1", len(groups))
	}
	if groups[0].Covers[0].ScanOf != "Front Cover" {
		t.Errorf("ScanOf = %q", groups[0].Covers[0].ScanOf)
	}
}

func TestGetRecentGames(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games/recent" {
			t.Errorf("path = %q, want /games/recent", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"games": []map[string]any{
				{"game_id": 456, "title": "New Game"},
			},
		})
	})

	games, err := c.GetRecentGames(context.Background(), 30, 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fatalf("len(games) = %d, want 1", len(games))
	}
	if games[0].Title != "New Game" {
		t.Errorf("Title = %q", games[0].Title)
	}
}

func TestGetRandomGames(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/games/random" {
			t.Errorf("path = %q, want /games/random", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"games": []map[string]any{
				{"game_id": 789, "title": "Random Game"},
			},
		})
	})

	games, err := c.GetRandomGames(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fatalf("len(games) = %d, want 1", len(games))
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Invalid API key"}`))
	})

	_, err := c.GetGenres(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *mobygames.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
}
