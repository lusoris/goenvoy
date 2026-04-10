package spotify_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/metadata/music/spotify"
	"github.com/lusoris/goenvoy/metadata"
)

func setup(t *testing.T, handler http.HandlerFunc) *spotify.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return spotify.New("test-token", metadata.WithBaseURL(srv.URL))
}

func TestAuthorizationHeader(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q, want Bearer test-token", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(spotify.Artist{ID: "1", Name: "Test"})
	})

	_, err := c.GetArtist(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearch(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q"); got != "coldplay" {
			t.Errorf("q = %q, want coldplay", got)
		}
		if got := r.URL.Query().Get("type"); got != "artist,track" {
			t.Errorf("type = %q, want artist,track", got)
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Errorf("limit = %q, want 10", got)
		}
		if got := r.URL.Query().Get("offset"); got != "0" {
			t.Errorf("offset = %q, want 0", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"artists": map[string]any{
				"items": []map[string]any{
					{"id": "4gzpq5DPGxSnKTe4SA8HAU", "name": "Coldplay"},
				},
				"total": 1, "limit": 10, "offset": 0,
			},
			"tracks": map[string]any{
				"items": []map[string]any{
					{"id": "t1", "name": "Yellow"},
				},
				"total": 1, "limit": 10, "offset": 0,
			},
		})
	})

	result, err := c.Search(context.Background(), "coldplay", []string{"artist", "track"}, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result.Artists == nil {
		t.Fatal("expected artists")
	}
	if len(result.Artists.Items) != 1 {
		t.Fatalf("artists len = %d, want 1", len(result.Artists.Items))
	}
	if result.Artists.Items[0].Name != "Coldplay" {
		t.Errorf("Name = %q, want Coldplay", result.Artists.Items[0].Name)
	}
	if result.Tracks == nil {
		t.Fatal("expected tracks")
	}
	if len(result.Tracks.Items) != 1 {
		t.Fatalf("tracks len = %d, want 1", len(result.Tracks.Items))
	}
}

func TestGetArtist(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/artists/4gzpq5DPGxSnKTe4SA8HAU" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":         "4gzpq5DPGxSnKTe4SA8HAU",
			"name":       "Coldplay",
			"genres":     []string{"alternative rock", "rock"},
			"popularity": 85,
			"followers":  map[string]any{"total": 50000000},
			"images":     []map[string]any{{"url": "https://img.example.com/1.jpg", "height": 640, "width": 640}},
		})
	})

	artist, err := c.GetArtist(context.Background(), "4gzpq5DPGxSnKTe4SA8HAU")
	if err != nil {
		t.Fatal(err)
	}
	if artist.Name != "Coldplay" {
		t.Errorf("Name = %q, want Coldplay", artist.Name)
	}
	if artist.Popularity != 85 {
		t.Errorf("Popularity = %d, want 85", artist.Popularity)
	}
	if artist.Followers.Total != 50000000 {
		t.Errorf("Followers = %d, want 50000000", artist.Followers.Total)
	}
	if len(artist.Genres) != 2 {
		t.Fatalf("Genres len = %d, want 2", len(artist.Genres))
	}
	if len(artist.Images) != 1 {
		t.Fatalf("Images len = %d, want 1", len(artist.Images))
	}
}

func TestGetArtistAlbums(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "5" {
			t.Errorf("limit = %q, want 5", r.URL.Query().Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{
				{"id": "a1", "name": "Parachutes", "album_type": "album", "total_tracks": 10},
				{"id": "a2", "name": "A Rush of Blood to the Head", "album_type": "album"},
			},
			"total": 2, "limit": 5, "offset": 0,
		})
	})

	result, err := c.GetArtistAlbums(context.Background(), "4gzpq5DPGxSnKTe4SA8HAU", 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Fatalf("len = %d, want 2", len(result.Items))
	}
	if result.Items[0].Name != "Parachutes" {
		t.Errorf("Name = %q, want Parachutes", result.Items[0].Name)
	}
	if result.Total != 2 {
		t.Errorf("Total = %d, want 2", result.Total)
	}
}

func TestGetArtistTopTracks(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("market"); got != "US" {
			t.Errorf("market = %q, want US", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"tracks": []map[string]any{
				{"id": "t1", "name": "Yellow", "popularity": 80},
				{"id": "t2", "name": "The Scientist", "popularity": 78},
			},
		})
	})

	tracks, err := c.GetArtistTopTracks(context.Background(), "4gzpq5DPGxSnKTe4SA8HAU", "US")
	if err != nil {
		t.Fatal(err)
	}
	if len(tracks) != 2 {
		t.Fatalf("len = %d, want 2", len(tracks))
	}
	if tracks[0].Name != "Yellow" {
		t.Errorf("Name = %q, want Yellow", tracks[0].Name)
	}
}

func TestGetRelatedArtists(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"artists": []map[string]any{
				{"id": "a1", "name": "Radiohead"},
				{"id": "a2", "name": "Muse"},
			},
		})
	})

	artists, err := c.GetRelatedArtists(context.Background(), "4gzpq5DPGxSnKTe4SA8HAU")
	if err != nil {
		t.Fatal(err)
	}
	if len(artists) != 2 {
		t.Fatalf("len = %d, want 2", len(artists))
	}
	if artists[0].Name != "Radiohead" {
		t.Errorf("Name = %q, want Radiohead", artists[0].Name)
	}
}

func TestGetAlbum(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/albums/2ix8vWvvSp2Yo7rKMiWpkg" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":           "2ix8vWvvSp2Yo7rKMiWpkg",
			"name":         "Parachutes",
			"album_type":   "album",
			"total_tracks": 10,
			"label":        "Parlophone",
			"popularity":   75,
			"genres":       []string{"alternative rock"},
			"copyrights":   []map[string]any{{"text": "2000 Parlophone", "type": "C"}},
			"tracks": map[string]any{
				"items": []map[string]any{
					{"id": "t1", "name": "Don't Panic", "track_number": 1},
				},
				"total": 10, "limit": 50, "offset": 0,
			},
		})
	})

	album, err := c.GetAlbum(context.Background(), "2ix8vWvvSp2Yo7rKMiWpkg")
	if err != nil {
		t.Fatal(err)
	}
	if album.Name != "Parachutes" {
		t.Errorf("Name = %q, want Parachutes", album.Name)
	}
	if album.Label != "Parlophone" {
		t.Errorf("Label = %q, want Parlophone", album.Label)
	}
	if album.Popularity != 75 {
		t.Errorf("Popularity = %d, want 75", album.Popularity)
	}
	if len(album.Tracks.Items) != 1 {
		t.Fatalf("Tracks len = %d, want 1", len(album.Tracks.Items))
	}
	if len(album.Copyrights) != 1 {
		t.Fatalf("Copyrights len = %d, want 1", len(album.Copyrights))
	}
}

func TestGetAlbumTracks(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "20" {
			t.Errorf("limit = %q, want 20", r.URL.Query().Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{
				{"id": "t1", "name": "Don't Panic", "track_number": 1, "disc_number": 1, "duration_ms": 137000},
			},
			"total": 10, "limit": 20, "offset": 0,
		})
	})

	result, err := c.GetAlbumTracks(context.Background(), "2ix8vWvvSp2Yo7rKMiWpkg", 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("len = %d, want 1", len(result.Items))
	}
	if result.Items[0].DurationMS != 137000 {
		t.Errorf("DurationMS = %d, want 137000", result.Items[0].DurationMS)
	}
}

func TestGetTrack(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tracks/3AJwUDP919kvQ9QcozQPxg" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":           "3AJwUDP919kvQ9QcozQPxg",
			"name":         "Yellow",
			"duration_ms":  266773,
			"explicit":     false,
			"popularity":   80,
			"album":        map[string]any{"id": "a1", "name": "Parachutes"},
			"external_ids": map[string]any{"isrc": "GBAYE0000651"},
		})
	})

	track, err := c.GetTrack(context.Background(), "3AJwUDP919kvQ9QcozQPxg")
	if err != nil {
		t.Fatal(err)
	}
	if track.Name != "Yellow" {
		t.Errorf("Name = %q, want Yellow", track.Name)
	}
	if track.DurationMS != 266773 {
		t.Errorf("DurationMS = %d, want 266773", track.DurationMS)
	}
	if track.Popularity != 80 {
		t.Errorf("Popularity = %d, want 80", track.Popularity)
	}
	if track.Album.Name != "Parachutes" {
		t.Errorf("Album.Name = %q, want Parachutes", track.Album.Name)
	}
}

func TestGetAudioFeatures(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/audio-features/3AJwUDP919kvQ9QcozQPxg" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":               "3AJwUDP919kvQ9QcozQPxg",
			"danceability":     0.429,
			"energy":           0.661,
			"key":              6.0,
			"loudness":         -7.299,
			"tempo":            173.372,
			"duration_ms":      266773,
			"time_signature":   4,
			"acousticness":     0.00474,
			"instrumentalness": 0.000107,
			"liveness":         0.0894,
			"valence":          0.267,
		})
	})

	features, err := c.GetAudioFeatures(context.Background(), "3AJwUDP919kvQ9QcozQPxg")
	if err != nil {
		t.Fatal(err)
	}
	if features.Danceability != 0.429 {
		t.Errorf("Danceability = %f, want 0.429", features.Danceability)
	}
	if features.Tempo != 173.372 {
		t.Errorf("Tempo = %f, want 173.372", features.Tempo)
	}
	if features.DurationMS != 266773 {
		t.Errorf("DurationMS = %d, want 266773", features.DurationMS)
	}
	if features.TimeSignature != 4 {
		t.Errorf("TimeSignature = %d, want 4", features.TimeSignature)
	}
}

func TestGetNewReleases(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "5" {
			t.Errorf("limit = %q, want 5", r.URL.Query().Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"albums": map[string]any{
				"items": []map[string]any{
					{"id": "a1", "name": "New Album"},
				},
				"total": 1, "limit": 5, "offset": 0,
			},
		})
	})

	result, err := c.GetNewReleases(context.Background(), 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("len = %d, want 1", len(result.Items))
	}
	if result.Items[0].Name != "New Album" {
		t.Errorf("Name = %q, want New Album", result.Items[0].Name)
	}
}

func TestGetCategories(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"categories": map[string]any{
				"items": []map[string]any{
					{"id": "pop", "name": "Pop", "icons": []map[string]any{{"url": "https://example.com/pop.jpg"}}},
					{"id": "rock", "name": "Rock"},
				},
				"total": 2, "limit": 20, "offset": 0,
			},
		})
	})

	result, err := c.GetCategories(context.Background(), 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Fatalf("len = %d, want 2", len(result.Items))
	}
	if result.Items[0].ID != "pop" {
		t.Errorf("ID = %q, want pop", result.Items[0].ID)
	}
}

func TestGetRecommendations(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("seed_artists"); got != "artist1" {
			t.Errorf("seed_artists = %q, want artist1", got)
		}
		if got := r.URL.Query().Get("seed_genres"); got != "rock,pop" {
			t.Errorf("seed_genres = %q, want rock,pop", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"seeds": []map[string]any{
				{"id": "artist1", "type": "ARTIST", "initialPoolSize": 500},
			},
			"tracks": []map[string]any{
				{"id": "t1", "name": "Recommended Track"},
			},
		})
	})

	result, err := c.GetRecommendations(context.Background(), spotify.RecommendationSeeds{
		SeedArtists: []string{"artist1"},
		SeedGenres:  []string{"rock", "pop"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Seeds) != 1 {
		t.Fatalf("seeds len = %d, want 1", len(result.Seeds))
	}
	if len(result.Tracks) != 1 {
		t.Fatalf("tracks len = %d, want 1", len(result.Tracks))
	}
	if result.Tracks[0].Name != "Recommended Track" {
		t.Errorf("Name = %q, want Recommended Track", result.Tracks[0].Name)
	}
}

func TestAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"status":401,"message":"Invalid access token"}}`))
	}))
	defer srv.Close()

	c := spotify.New("bad-token", metadata.WithBaseURL(srv.URL))
	_, err := c.GetArtist(context.Background(), "1")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *spotify.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c := spotify.New("token", metadata.WithHTTPClient(custom))
	// Verify the client was created without error.
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestAPIErrorFormat(t *testing.T) {
	e := &spotify.APIError{StatusCode: 404, Status: "404 Not Found", Body: "not found"}
	if got := e.Error(); got != "spotify: 404 Not Found: not found" {
		t.Errorf("Error() = %q", got)
	}
}

func TestSearchPath(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if !hasPrefix(r.URL.Path, "/search") {
			t.Errorf("path = %q, want /search prefix", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(spotify.SearchResult{})
	})

	_, err := c.Search(context.Background(), "test", []string{"track"}, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
