package steam_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golusoris/goenvoy/metadata/game/steam"
)

func setup(t *testing.T, handler http.HandlerFunc) *steam.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := steam.NewWithAPIKey("test-key")
	c.SetStoreURL(srv.URL)
	c.SetWebAPIURL(srv.URL)
	return c
}

func TestGetAppDetails(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("appids"); got != "730" {
			t.Errorf("appids = %q, want 730", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"730": map[string]any{
				"success": true,
				"data": map[string]any{
					"type":         "game",
					"name":         "Counter-Strike 2",
					"steam_appid":  730,
					"required_age": 0,
					"is_free":      true,
					"header_image": "https://cdn.akamai.steamstatic.com/steam/apps/730/header.jpg",
					"developers":   []string{"Valve"},
					"publishers":   []string{"Valve"},
					"platforms":    map[string]any{"windows": true, "mac": false, "linux": true},
					"release_date": map[string]any{"coming_soon": false, "date": "21 Aug, 2012"},
				},
			},
		})
	})

	details, err := c.GetAppDetails(context.Background(), 730)
	if err != nil {
		t.Fatal(err)
	}
	if details.Name != "Counter-Strike 2" {
		t.Errorf("Name = %q, want Counter-Strike 2", details.Name)
	}
	if details.SteamAppID != 730 {
		t.Errorf("SteamAppID = %d, want 730", details.SteamAppID)
	}
	if !details.IsFree {
		t.Error("IsFree = false, want true")
	}
	if !details.Platforms.Windows {
		t.Error("Platforms.Windows = false, want true")
	}
	if !details.Platforms.Linux {
		t.Error("Platforms.Linux = false, want true")
	}
	if len(details.Developers) != 1 || details.Developers[0] != "Valve" {
		t.Errorf("Developers = %v, want [Valve]", details.Developers)
	}
}

func TestGetAppDetailsNotFound(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"99999": map[string]any{
				"success": false,
			},
		})
	})

	_, err := c.GetAppDetails(context.Background(), 99999)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q, want to contain 'not found'", err.Error())
	}
}

func TestGetMultipleAppDetails(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("appids")
		if !strings.Contains(ids, "730") || !strings.Contains(ids, "440") {
			t.Errorf("appids = %q, want to contain 730 and 440", ids)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"730": map[string]any{
				"success": true,
				"data":    map[string]any{"name": "Counter-Strike 2", "steam_appid": 730},
			},
			"440": map[string]any{
				"success": true,
				"data":    map[string]any{"name": "Team Fortress 2", "steam_appid": 440},
			},
		})
	})

	results, err := c.GetMultipleAppDetails(context.Background(), []int{730, 440})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("len = %d, want 2", len(results))
	}
	if results[730].Name != "Counter-Strike 2" {
		t.Errorf("730 Name = %q, want Counter-Strike 2", results[730].Name)
	}
	if results[440].Name != "Team Fortress 2" {
		t.Errorf("440 Name = %q, want Team Fortress 2", results[440].Name)
	}
}

func TestGetFeatured(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/featured") {
			t.Errorf("path = %q, want /featured", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"large_capsules": []map[string]any{
				{"id": 730, "name": "Counter-Strike 2", "discounted": false, "final_price": 0},
			},
			"featured_win": []map[string]any{
				{"id": 440, "name": "Team Fortress 2", "final_price": 0},
			},
			"featured_mac":   []any{},
			"featured_linux": []any{},
		})
	})

	resp, err := c.GetFeatured(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.LargeCapsules) != 1 {
		t.Fatalf("len(LargeCapsules) = %d, want 1", len(resp.LargeCapsules))
	}
	if resp.LargeCapsules[0].Name != "Counter-Strike 2" {
		t.Errorf("Name = %q, want Counter-Strike 2", resp.LargeCapsules[0].Name)
	}
	if len(resp.FeaturedWin) != 1 {
		t.Fatalf("len(FeaturedWin) = %d, want 1", len(resp.FeaturedWin))
	}
}

func TestGetFeaturedCategories(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"0":{"id":"cat_newreleases","name":"New Releases"},"status":1}`))
	})

	resp, err := c.GetFeaturedCategories(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if _, ok := (*resp)["0"]; !ok {
		t.Error("expected key '0' in response")
	}
}

func TestGetAppList(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "GetAppList") {
			t.Errorf("path = %q, want to contain GetAppList", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"applist": map[string]any{
				"apps": []map[string]any{
					{"appid": 10, "name": "Counter-Strike"},
					{"appid": 20, "name": "Team Fortress Classic"},
					{"appid": 730, "name": "Counter-Strike 2"},
				},
			},
		})
	})

	apps, err := c.GetAppList(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 3 {
		t.Fatalf("len = %d, want 3", len(apps))
	}
	if apps[2].AppID != 730 {
		t.Errorf("AppID = %d, want 730", apps[2].AppID)
	}
	if apps[2].Name != "Counter-Strike 2" {
		t.Errorf("Name = %q, want Counter-Strike 2", apps[2].Name)
	}
}

func TestGetCurrentPlayers(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("appid"); got != "730" {
			t.Errorf("appid = %q, want 730", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"response": map[string]any{
				"player_count": 850000,
				"result":       1,
			},
		})
	})

	count, err := c.GetCurrentPlayers(context.Background(), 730)
	if err != nil {
		t.Fatal(err)
	}
	if count != 850000 {
		t.Errorf("count = %d, want 850000", count)
	}
}

func TestGetAppNews(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("appid"); got != "730" {
			t.Errorf("appid = %q, want 730", got)
		}
		if got := r.URL.Query().Get("count"); got != "3" {
			t.Errorf("count = %q, want 3", got)
		}
		if got := r.URL.Query().Get("maxlength"); got != "300" {
			t.Errorf("maxlength = %q, want 300", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"appnews": map[string]any{
				"appid": 730,
				"newsitems": []map[string]any{
					{
						"gid":       "123456",
						"title":     "Update Released",
						"url":       "https://example.com/news/1",
						"author":    "Valve",
						"contents":  "New update is live...",
						"feedlabel": "Community Announcements",
						"date":      1700000000,
						"feedname":  "steam_community_announcements",
						"tags":      []string{"patchnotes"},
					},
				},
			},
		})
	})

	news, err := c.GetAppNews(context.Background(), 730, 3, 300)
	if err != nil {
		t.Fatal(err)
	}
	if len(news) != 1 {
		t.Fatalf("len = %d, want 1", len(news))
	}
	if news[0].Title != "Update Released" {
		t.Errorf("Title = %q, want Update Released", news[0].Title)
	}
	if news[0].Author != "Valve" {
		t.Errorf("Author = %q, want Valve", news[0].Author)
	}
	if news[0].GID != "123456" {
		t.Errorf("GID = %q, want 123456", news[0].GID)
	}
}

func TestGetGlobalAchievements(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("gameid"); got != "730" {
			t.Errorf("gameid = %q, want 730", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"achievementpercentages": map[string]any{
				"achievements": []map[string]any{
					{"name": "Win a Round", "percent": 85.5},
					{"name": "Win a Match", "percent": 72.3},
				},
			},
		})
	})

	achievements, err := c.GetGlobalAchievements(context.Background(), 730)
	if err != nil {
		t.Fatal(err)
	}
	if len(achievements) != 2 {
		t.Fatalf("len = %d, want 2", len(achievements))
	}
	if achievements[0].Name != "Win a Round" {
		t.Errorf("Name = %q, want Win a Round", achievements[0].Name)
	}
	if achievements[0].Percent != 85.5 {
		t.Errorf("Percent = %f, want 85.5", achievements[0].Percent)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	})

	_, err := c.GetCurrentPlayers(context.Background(), 730)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *steam.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want 500", apiErr.StatusCode)
	}
}

func TestWithAPIKey(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Errorf("key = %q, want test-key", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"response": map[string]any{
				"player_count": 100,
				"result":       1,
			},
		})
	})

	_, err := c.GetCurrentPlayers(context.Background(), 730)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMultipleAppDetailsPartialSuccess(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"730": map[string]any{
				"success": true,
				"data":    map[string]any{"name": "Counter-Strike 2", "steam_appid": 730},
			},
			"99999": map[string]any{
				"success": false,
			},
		})
	})

	results, err := c.GetMultipleAppDetails(context.Background(), []int{730, 99999})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("len = %d, want 1 (only successful)", len(results))
	}
	if results[730] == nil {
		t.Fatal("expected 730 to be present")
	}
	if results[99999] != nil {
		t.Error("expected 99999 to be nil")
	}
}
