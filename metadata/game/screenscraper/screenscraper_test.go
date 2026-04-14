package screenscraper_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/game/screenscraper"
)

func setup(t *testing.T, handler http.HandlerFunc) *screenscraper.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return screenscraper.New("testdev", "testpass", "testapp",
		metadata.WithBaseURL(srv.URL),
		screenscraper.WithUser("user", "userpass"),
	)
}

func TestGetGameInfo(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if got := r.URL.Query().Get("devid"); got != "testdev" {
			t.Errorf("devid = %q, want testdev", got)
		}
		if got := r.URL.Query().Get("devpassword"); got != "testpass" {
			t.Errorf("devpassword = %q, want testpass", got)
		}
		if got := r.URL.Query().Get("softname"); got != "testapp" {
			t.Errorf("softname = %q, want testapp", got)
		}
		if got := r.URL.Query().Get("ssid"); got != "user" {
			t.Errorf("ssid = %q, want user", got)
		}
		if got := r.URL.Query().Get("output"); got != "json" {
			t.Errorf("output = %q, want json", got)
		}
		if got := r.URL.Query().Get("crc"); got != "ABCD1234" {
			t.Errorf("crc = %q, want ABCD1234", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"header": map[string]any{
				"APIversion": "2",
				"success":    "true",
			},
			"response": map[string]any{
				"jeu": map[string]any{
					"id": "3",
					"noms": []map[string]any{
						{"region": "us", "text": "Sonic the Hedgehog"},
						{"region": "jp", "text": "ソニック・ザ・ヘッジホッグ"},
					},
					"systemeid": "1",
					"editeur":   map[string]any{"id": "10", "text": "Sega"},
					"joueurs":   "1",
				},
			},
		})
	})

	result, err := c.GetGameInfo(context.Background(), &screenscraper.GameInfoOptions{
		CRC: "ABCD1234",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Response.Game.ID != "3" {
		t.Errorf("ID = %q, want 3", result.Response.Game.ID)
	}
	if len(result.Response.Game.Names) != 2 {
		t.Fatalf("len(Names) = %d, want 2", len(result.Response.Game.Names))
	}
	if result.Response.Game.Names[0].Text != "Sonic the Hedgehog" {
		t.Errorf("Names[0].Text = %q", result.Response.Game.Names[0].Text)
	}
	if result.Response.Game.Publisher.Text != "Sega" {
		t.Errorf("Publisher.Text = %q", result.Response.Game.Publisher.Text)
	}
}

func TestSearchGames(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("recherche"); got != "sonic" {
			t.Errorf("recherche = %q, want sonic", got)
		}
		if got := r.URL.Query().Get("systemeid"); got != "1" {
			t.Errorf("systemeid = %q, want 1", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"header": map[string]any{"success": "true"},
			"response": map[string]any{
				"jeux": []map[string]any{
					{"id": "3", "noms": []map[string]any{{"region": "us", "text": "Sonic the Hedgehog"}}},
				},
			},
		})
	})

	result, err := c.SearchGames(context.Background(), "sonic", "1")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Response.Games) != 1 {
		t.Fatalf("len(Games) = %d, want 1", len(result.Response.Games))
	}
	if result.Response.Games[0].ID != "3" {
		t.Errorf("ID = %q", result.Response.Games[0].ID)
	}
}

func TestGetSystems(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/systemesListe.php" {
			t.Errorf("path = %q, want /systemesListe.php", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"header": map[string]any{"success": "true"},
			"response": map[string]any{
				"systemes": []map[string]any{
					{"id": "1", "noms": []map[string]any{{"region": "us", "text": "Mega Drive"}}},
				},
			},
		})
	})

	result, err := c.GetSystems(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Response.Systems) != 1 {
		t.Fatalf("len(Systems) = %d, want 1", len(result.Response.Systems))
	}
	if result.Response.Systems[0].ID != "1" {
		t.Errorf("ID = %q", result.Response.Systems[0].ID)
	}
}

func TestGetGenres(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/genresListe.php" {
			t.Errorf("path = %q, want /genresListe.php", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"header": map[string]any{"success": "true"},
			"response": map[string]any{
				"genres": []map[string]any{
					{"id": "1", "noms": []map[string]any{{"langue": "en", "text": "Action"}}},
				},
			},
		})
	})

	result, err := c.GetGenres(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Response.Genres) != 1 {
		t.Fatalf("len(Genres) = %d, want 1", len(result.Response.Genres))
	}
	if result.Response.Genres[0].ID != "1" {
		t.Errorf("ID = %q", result.Response.Genres[0].ID)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`API closed`))
	})

	_, err := c.GetSystems(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *screenscraper.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Errorf("StatusCode = %d, want 403", apiErr.StatusCode)
	}
}
