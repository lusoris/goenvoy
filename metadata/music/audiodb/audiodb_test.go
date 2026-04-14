package audiodb_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/music/audiodb/v2"
)

func setup(t *testing.T, handler http.HandlerFunc) *audiodb.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return audiodb.New("2", metadata.WithBaseURL(srv.URL))
}

func TestSearchArtist(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("s"); got != "coldplay" {
			t.Errorf("s = %q, want coldplay", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"artists": []map[string]any{
				{"idArtist": "111239", "strArtist": "Coldplay", "strGenre": "Alternative Rock"},
			},
		})
	})

	artists, err := c.SearchArtist(context.Background(), "coldplay")
	if err != nil {
		t.Fatal(err)
	}
	if len(artists) != 1 {
		t.Fatalf("len = %d, want 1", len(artists))
	}
	if artists[0].StrArtist != "Coldplay" {
		t.Errorf("StrArtist = %q, want Coldplay", artists[0].StrArtist)
	}
	if artists[0].IDArtist != "111239" {
		t.Errorf("IDArtist = %q, want 111239", artists[0].IDArtist)
	}
}

func TestSearchAlbum(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"album": []map[string]any{
				{"idAlbum": "2115888", "strAlbum": "Parachutes", "strArtist": "Coldplay"},
			},
		})
	})

	albums, err := c.SearchAlbum(context.Background(), "coldplay")
	if err != nil {
		t.Fatal(err)
	}
	if len(albums) != 1 {
		t.Fatalf("len = %d, want 1", len(albums))
	}
	if albums[0].StrAlbum != "Parachutes" {
		t.Errorf("StrAlbum = %q, want Parachutes", albums[0].StrAlbum)
	}
}

func TestSearchTrack(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("s"); got != "coldplay" {
			t.Errorf("s = %q, want coldplay", got)
		}
		if got := r.URL.Query().Get("t"); got != "yellow" {
			t.Errorf("t = %q, want yellow", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"track": []map[string]any{
				{"idTrack": "32793800", "strTrack": "Yellow", "strArtist": "Coldplay"},
			},
		})
	})

	tracks, err := c.SearchTrack(context.Background(), "coldplay", "yellow")
	if err != nil {
		t.Fatal(err)
	}
	if len(tracks) != 1 {
		t.Fatalf("len = %d, want 1", len(tracks))
	}
	if tracks[0].StrTrack != "Yellow" {
		t.Errorf("StrTrack = %q, want Yellow", tracks[0].StrTrack)
	}
}

func TestGetArtist(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"artists": []map[string]any{
				{"idArtist": "111239", "strArtist": "Coldplay", "strCountry": "London, England"},
			},
		})
	})

	artist, err := c.GetArtist(context.Background(), "111239")
	if err != nil {
		t.Fatal(err)
	}
	if artist == nil {
		t.Fatal("expected non-nil artist")
	}
	if artist.StrCountry != "London, England" {
		t.Errorf("StrCountry = %q, want London, England", artist.StrCountry)
	}
}

func TestGetAlbum(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"album": []map[string]any{
				{"idAlbum": "2115888", "strAlbum": "Parachutes", "intYearReleased": "2000"},
			},
		})
	})

	album, err := c.GetAlbum(context.Background(), "2115888")
	if err != nil {
		t.Fatal(err)
	}
	if album == nil {
		t.Fatal("expected non-nil album")
	}
	if album.IntYearReleased != "2000" {
		t.Errorf("IntYearReleased = %q, want 2000", album.IntYearReleased)
	}
}

func TestGetAlbumsByArtist(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"album": []map[string]any{
				{"idAlbum": "2115888", "strAlbum": "Parachutes"},
				{"idAlbum": "2115889", "strAlbum": "A Rush of Blood to the Head"},
			},
		})
	})

	albums, err := c.GetAlbumsByArtist(context.Background(), "111239")
	if err != nil {
		t.Fatal(err)
	}
	if len(albums) != 2 {
		t.Fatalf("len = %d, want 2", len(albums))
	}
}

func TestGetTracksByAlbum(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"track": []map[string]any{
				{"idTrack": "32793800", "strTrack": "Yellow", "intTrackNumber": "4"},
			},
		})
	})

	tracks, err := c.GetTracksByAlbum(context.Background(), "2115888")
	if err != nil {
		t.Fatal(err)
	}
	if len(tracks) != 1 {
		t.Fatalf("len = %d, want 1", len(tracks))
	}
	if tracks[0].IntTrackNumber != "4" {
		t.Errorf("IntTrackNumber = %q, want 4", tracks[0].IntTrackNumber)
	}
}

func TestGetMusicVideos(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"mvids": []map[string]any{
				{"idTrack": "32793800", "strTrack": "Yellow", "strMusicVideo": "https://example.com/video"},
			},
		})
	})

	mvids, err := c.GetMusicVideos(context.Background(), "111239")
	if err != nil {
		t.Fatal(err)
	}
	if len(mvids) != 1 {
		t.Fatalf("len = %d, want 1", len(mvids))
	}
	if mvids[0].StrMusicVideo != "https://example.com/video" {
		t.Errorf("StrMusicVideo = %q", mvids[0].StrMusicVideo)
	}
}

func TestGetDiscography(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"album": []map[string]any{
				{"strAlbum": "Parachutes", "intYearReleased": "2000"},
				{"strAlbum": "A Rush of Blood to the Head", "intYearReleased": "2002"},
			},
		})
	})

	disco, err := c.GetDiscography(context.Background(), "coldplay")
	if err != nil {
		t.Fatal(err)
	}
	if len(disco) != 2 {
		t.Fatalf("len = %d, want 2", len(disco))
	}
	if disco[0].StrAlbum != "Parachutes" {
		t.Errorf("StrAlbum = %q, want Parachutes", disco[0].StrAlbum)
	}
}

func TestGetTopTracks(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"track": []map[string]any{
				{"idTrack": "1", "strTrack": "Yellow"},
			},
		})
	})

	tracks, err := c.GetTopTracks(context.Background(), "coldplay")
	if err != nil {
		t.Fatal(err)
	}
	if len(tracks) != 1 {
		t.Fatalf("len = %d, want 1", len(tracks))
	}
}

func TestGetTrending(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("country"); got != "us" {
			t.Errorf("country = %q, want us", got)
		}
		if got := r.URL.Query().Get("format"); got != "singles" {
			t.Errorf("format = %q, want singles", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"trending": []map[string]any{
				{"idArtist": "111239", "strArtist": "Coldplay", "intChartPlace": "1"},
			},
		})
	})

	trending, err := c.GetTrending(context.Background(), "us", "singles")
	if err != nil {
		t.Fatal(err)
	}
	if len(trending) != 1 {
		t.Fatalf("len = %d, want 1", len(trending))
	}
	if trending[0].IntChartPlace != "1" {
		t.Errorf("IntChartPlace = %q, want 1", trending[0].IntChartPlace)
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}))
	defer srv.Close()

	c := audiodb.New("2", metadata.WithBaseURL(srv.URL))
	_, err := c.SearchArtist(context.Background(), "coldplay")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *audiodb.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusInternalServerError)
	}
}

func TestAPIKeyInURL(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Path; got != "/2/search.php" {
			t.Errorf("path = %q, want /2/search.php", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"artists": nil})
	}))
	defer srv.Close()

	c := audiodb.New("2", metadata.WithBaseURL(srv.URL))
	_, err := c.SearchArtist(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
}
