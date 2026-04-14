package retroachievements_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/game/retroachievements"
)

func setup(t *testing.T, handler http.HandlerFunc) *retroachievements.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return retroachievements.New("test-key", metadata.WithBaseURL(srv.URL))
}

func TestGetGame(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if got := r.URL.Query().Get("y"); got != "test-key" {
			t.Errorf("y = %q, want test-key", got)
		}
		if got := r.URL.Query().Get("i"); got != "1" {
			t.Errorf("i = %q, want 1", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"Title":       "Sonic the Hedgehog",
			"ConsoleID":   1,
			"ConsoleName": "Mega Drive",
			"Developer":   "Sonic Team",
			"Publisher":   "Sega",
			"Genre":       "Platformer",
			"Released":    "1991-06-23",
			"ImageIcon":   "/Images/085573.png",
			"ImageBoxArt": "/Images/051872.png",
		})
	})

	game, err := c.GetGame(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if game.Title != "Sonic the Hedgehog" {
		t.Errorf("Title = %q", game.Title)
	}
	if game.ConsoleName != "Mega Drive" {
		t.Errorf("ConsoleName = %q", game.ConsoleName)
	}
	if game.Developer != "Sonic Team" {
		t.Errorf("Developer = %q", game.Developer)
	}
}

func TestGetGameExtended(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/API_GetGameExtended.php" {
			t.Errorf("path = %q, want /API_GetGameExtended.php", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"ID":                 1,
			"Title":              "Sonic the Hedgehog",
			"ConsoleID":          1,
			"ConsoleName":        "Mega Drive",
			"NumAchievements":    23,
			"NumDistinctPlayers": 5000,
			"Achievements": map[string]any{
				"1": map[string]any{
					"ID":          1,
					"Title":       "Green Hill Zone",
					"Description": "Complete Green Hill Zone",
					"Points":      10,
					"TrueRatio":   15,
					"Author":      "Admin",
				},
			},
		})
	})

	game, err := c.GetGameExtended(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if game.ID != 1 {
		t.Errorf("ID = %d, want 1", game.ID)
	}
	if game.NumAchievements != 23 {
		t.Errorf("NumAchievements = %d, want 23", game.NumAchievements)
	}
	ach, ok := game.Achievements["1"]
	if !ok {
		t.Fatal("expected achievement with key 1")
	}
	if ach.Title != "Green Hill Zone" {
		t.Errorf("Achievement.Title = %q", ach.Title)
	}
	if ach.Points != 10 {
		t.Errorf("Achievement.Points = %d, want 10", ach.Points)
	}
}

func TestGetGameHashes(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/API_GetGameHashes.php" {
			t.Errorf("path = %q, want /API_GetGameHashes.php", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"Results": []map[string]any{
				{"MD5": "abc123", "Name": "Sonic (USA).md", "Labels": []string{"nointro"}},
			},
		})
	})

	result, err := c.GetGameHashes(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Results) != 1 {
		t.Fatalf("len(Results) = %d, want 1", len(result.Results))
	}
	if result.Results[0].MD5 != "abc123" {
		t.Errorf("MD5 = %q", result.Results[0].MD5)
	}
	if result.Results[0].Name != "Sonic (USA).md" {
		t.Errorf("Name = %q", result.Results[0].Name)
	}
}

func TestGetConsoleIDs(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/API_GetConsoleIDs.php" {
			t.Errorf("path = %q, want /API_GetConsoleIDs.php", r.URL.Path)
		}
		if got := r.URL.Query().Get("a"); got != "1" {
			t.Errorf("a = %q, want 1", got)
		}
		if got := r.URL.Query().Get("g"); got != "1" {
			t.Errorf("g = %q, want 1", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"ID": 1, "Name": "Mega Drive", "IconURL": "/icon.png", "Active": true, "IsGameSystem": true},
			{"ID": 2, "Name": "SNES", "IconURL": "/snes.png", "Active": true, "IsGameSystem": true},
		})
	})

	consoles, err := c.GetConsoleIDs(context.Background(), true, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(consoles) != 2 {
		t.Fatalf("len(consoles) = %d, want 2", len(consoles))
	}
	if consoles[0].Name != "Mega Drive" {
		t.Errorf("Name = %q", consoles[0].Name)
	}
	if !consoles[0].IsGameSystem {
		t.Error("IsGameSystem = false, want true")
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Invalid API key"}`))
	})

	_, err := c.GetGame(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *retroachievements.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
}
