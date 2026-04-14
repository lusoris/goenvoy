package hasheous_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/game/hasheous"
)

func setup(t *testing.T, handler http.HandlerFunc) *hasheous.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return hasheous.New(metadata.WithBaseURL(srv.URL))
}

func TestLookupByHash(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/Lookup/ByHash" {
			t.Errorf("path = %q, want /Lookup/ByHash", r.URL.Path)
		}
		if got := r.URL.Query().Get("returnAllSources"); got != "true" {
			t.Errorf("returnAllSources = %q, want true", got)
		}
		if got := r.URL.Query().Get("returnFields"); got != "All" {
			t.Errorf("returnFields = %q, want All", got)
		}

		body, _ := io.ReadAll(r.Body)
		var req hasheous.HashLookupRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if req.SHA1 != "abc123" {
			t.Errorf("SHA1 = %q, want abc123", req.SHA1)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":   1234,
			"name": "Super Mario Bros.",
			"platform": map[string]any{
				"name": "Nintendo Entertainment System",
			},
			"signatures": map[string]any{
				"NoIntros": []map[string]any{
					{
						"game": map[string]any{"name": "Super Mario Bros.", "year": "1985"},
						"rom":  map[string]any{"name": "Super Mario Bros. (World).nes", "sha1": "abc123"},
					},
				},
			},
		})
	})

	result, err := c.LookupByHash(context.Background(), &hasheous.HashLookupRequest{SHA1: "abc123"}, true)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != 1234 {
		t.Errorf("ID = %d, want 1234", result.ID)
	}
	if result.Name != "Super Mario Bros." {
		t.Errorf("Name = %q", result.Name)
	}
	if result.Platform == nil || result.Platform.Name != "Nintendo Entertainment System" {
		t.Errorf("Platform = %v", result.Platform)
	}
	sigs, ok := result.Signatures["NoIntros"]
	if !ok || len(sigs) != 1 {
		t.Fatalf("NoIntros signatures = %v", result.Signatures)
	}
	if sigs[0].Game.Name != "Super Mario Bros." {
		t.Errorf("Game.Name = %q", sigs[0].Game.Name)
	}
}

func TestLookupBySHA1(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Lookup/ByHash/sha1/deadbeef" {
			t.Errorf("path = %q, want /Lookup/ByHash/sha1/deadbeef", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":   42,
			"name": "Tetris",
		})
	})

	result, err := c.LookupBySHA1(context.Background(), "deadbeef")
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Tetris" {
		t.Errorf("Name = %q, want Tetris", result.Name)
	}
}

func TestLookupByMD5(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Lookup/ByHash/md5/abc123md5" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": 1, "name": "Test"})
	})

	result, err := c.LookupByMD5(context.Background(), "abc123md5")
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Test" {
		t.Errorf("Name = %q", result.Name)
	}
}

func TestGetPlatforms(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Lookup/Platforms" {
			t.Errorf("path = %q, want /Lookup/Platforms", r.URL.Path)
		}
		if got := r.URL.Query().Get("PageNumber"); got != "1" {
			t.Errorf("PageNumber = %q, want 1", got)
		}
		if got := r.URL.Query().Get("PageSize"); got != "50" {
			t.Errorf("PageSize = %q, want 50", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{"NES", "SNES", "Game Boy"})
	})

	platforms, err := c.GetPlatforms(context.Background(), 1, 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(platforms) != 3 {
		t.Fatalf("len = %d, want 3", len(platforms))
	}
	if platforms[0] != "NES" {
		t.Errorf("platforms[0] = %q, want NES", platforms[0])
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"detail":"not found"}`))
	})

	_, err := c.LookupBySHA1(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *hasheous.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *hasheous.APIError", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
}
