package lastfm

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return New("test-key", metadata.WithBaseURL(ts.URL))
}

func TestGetArtistInfo(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("api_key") != "test-key" {
			http.Error(w, "bad key", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"artist": Artist{Name: "Radiohead", URL: "https://www.last.fm/music/Radiohead"},
		})
	})

	a, err := c.GetArtistInfo(context.Background(), "Radiohead")
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "Radiohead" {
		t.Fatalf("unexpected artist: %+v", a)
	}
}

func TestGetAlbumInfo(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"album": Album{Name: "OK Computer", Artist: "Radiohead"},
		})
	})

	a, err := c.GetAlbumInfo(context.Background(), "Radiohead", "OK Computer")
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "OK Computer" {
		t.Fatalf("unexpected album: %+v", a)
	}
}

func TestGetTrackInfo(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"track": Track{Name: "Creep", Duration: "238"},
		})
	})

	tr, err := c.GetTrackInfo(context.Background(), "Radiohead", "Creep")
	if err != nil {
		t.Fatal(err)
	}
	if tr.Name != "Creep" {
		t.Fatalf("unexpected track: %+v", tr)
	}
}

func TestGetSimilarArtists(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"similarartists": map[string]any{
				"artist": []Artist{{Name: "Muse"}},
			},
		})
	})

	artists, err := c.GetSimilarArtists(context.Background(), "Radiohead", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(artists) != 1 || artists[0].Name != "Muse" {
		t.Fatalf("unexpected artists: %+v", artists)
	}
}

func TestGetTopAlbums(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"topalbums": map[string]any{
				"album": []Album{{Name: "OK Computer"}},
			},
		})
	})

	albums, err := c.GetTopAlbums(context.Background(), "Radiohead", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(albums) != 1 || albums[0].Name != "OK Computer" {
		t.Fatalf("unexpected albums: %+v", albums)
	}
}

func TestGetTopTracks(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"toptracks": map[string]any{
				"track": []Track{{Name: "Creep"}},
			},
		})
	})

	tracks, err := c.GetTopTracks(context.Background(), "Radiohead", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(tracks) != 1 || tracks[0].Name != "Creep" {
		t.Fatalf("unexpected tracks: %+v", tracks)
	}
}

func TestGetChartTopArtists(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"artists": map[string]any{
				"artist": []Artist{{Name: "The Weeknd"}},
			},
		})
	})

	artists, err := c.GetChartTopArtists(context.Background(), 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(artists) != 1 || artists[0].Name != "The Weeknd" {
		t.Fatalf("unexpected artists: %+v", artists)
	}
}

func TestGetTopTags(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"toptags": map[string]any{
				"tag": []Tag{{Name: "rock"}},
			},
		})
	})

	tags, err := c.GetTopTags(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 1 || tags[0].Name != "rock" {
		t.Fatalf("unexpected tags: %+v", tags)
	}
}

func TestSearchArtist(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"results": map[string]any{
				"artistmatches": map[string]any{
					"artist": []Artist{{Name: "Radiohead"}},
				},
			},
		})
	})

	artists, err := c.SearchArtist(context.Background(), "Radiohead", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(artists) != 1 || artists[0].Name != "Radiohead" {
		t.Fatalf("unexpected artists: %+v", artists)
	}
}

func TestSearchAlbum(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"results": map[string]any{
				"albummatches": map[string]any{
					"album": []Album{{Name: "OK Computer"}},
				},
			},
		})
	})

	albums, err := c.SearchAlbum(context.Background(), "OK Computer", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(albums) != 1 || albums[0].Name != "OK Computer" {
		t.Fatalf("unexpected albums: %+v", albums)
	}
}

func TestSearchTrack(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"results": map[string]any{
				"trackmatches": map[string]any{
					"track": []Track{{Name: "Creep"}},
				},
			},
		})
	})

	tracks, err := c.SearchTrack(context.Background(), "Creep", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(tracks) != 1 || tracks[0].Name != "Creep" {
		t.Fatalf("unexpected tracks: %+v", tracks)
	}
}

func TestLastFMError(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"error":   6,
			"message": "Artist not found",
		})
	})

	_, err := c.GetArtistInfo(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	})

	_, err := c.GetArtistInfo(context.Background(), "test")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusForbidden {
		t.Fatalf("unexpected status: %d", apiErr.StatusCode)
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	custom := &http.Client{}
	c := New("key", metadata.WithHTTPClient(custom))
	if c.HTTPClient() != custom {
		t.Fatal("custom HTTP client not set")
	}
}
