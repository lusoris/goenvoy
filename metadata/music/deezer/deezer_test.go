package deezer_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/music/deezer"
)

func setup(t *testing.T, handler http.HandlerFunc) *deezer.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return deezer.New(metadata.WithBaseURL(srv.URL))
}

func TestSearch(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q"); got != "eminem" {
			t.Errorf("q = %q, want eminem", got)
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Errorf("limit = %q, want 10", got)
		}
		if got := r.URL.Query().Get("index"); got != "0" {
			t.Errorf("index = %q, want 0", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 916424, "title": "Lose Yourself", "duration": 326, "rank": 898543},
			},
			"total": 1,
		})
	})

	result, err := c.Search(context.Background(), "eminem", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("len = %d, want 1", len(result.Data))
	}
	if result.Data[0].Title != "Lose Yourself" {
		t.Errorf("Title = %q, want Lose Yourself", result.Data[0].Title)
	}
	if result.Data[0].Duration != 326 {
		t.Errorf("Duration = %d, want 326", result.Data[0].Duration)
	}
	if result.Total != 1 {
		t.Errorf("Total = %d, want 1", result.Total)
	}
}

func TestSearchAlbums(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Path; got != "/search/album" {
			t.Errorf("path = %q, want /search/album", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 103248, "title": "The Marshall Mathers LP"},
			},
			"total": 1,
		})
	})

	result, err := c.SearchAlbums(context.Background(), "eminem", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("len = %d, want 1", len(result.Data))
	}
	if result.Data[0].Title != "The Marshall Mathers LP" {
		t.Errorf("Title = %q", result.Data[0].Title)
	}
}

func TestSearchArtists(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Path; got != "/search/artist" {
			t.Errorf("path = %q, want /search/artist", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 13, "name": "Eminem", "nb_fan": 18000000},
			},
			"total": 1,
		})
	})

	result, err := c.SearchArtists(context.Background(), "eminem", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("len = %d, want 1", len(result.Data))
	}
	if result.Data[0].Name != "Eminem" {
		t.Errorf("Name = %q, want Eminem", result.Data[0].Name)
	}
	if result.Data[0].NbFan != 18000000 {
		t.Errorf("NbFan = %d, want 18000000", result.Data[0].NbFan)
	}
}

func TestGetArtist(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/artist/27" {
			t.Errorf("path = %q, want /artist/27", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":       27,
			"name":     "Daft Punk",
			"nb_album": 12,
			"nb_fan":   4600000,
			"link":     "https://www.deezer.com/artist/27",
			"picture":  "https://api.deezer.com/artist/27/image",
			"type":     "artist",
		})
	})

	artist, err := c.GetArtist(context.Background(), 27)
	if err != nil {
		t.Fatal(err)
	}
	if artist.Name != "Daft Punk" {
		t.Errorf("Name = %q, want Daft Punk", artist.Name)
	}
	if artist.NbAlbum != 12 {
		t.Errorf("NbAlbum = %d, want 12", artist.NbAlbum)
	}
	if artist.NbFan != 4600000 {
		t.Errorf("NbFan = %d, want 4600000", artist.NbFan)
	}
	if artist.Type != "artist" {
		t.Errorf("Type = %q, want artist", artist.Type)
	}
}

func TestGetArtistTopTracks(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Errorf("limit = %q, want 5", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 3135553, "title": "Get Lucky", "duration": 369, "rank": 950000},
				{"id": 3135554, "title": "Around the World", "duration": 420},
			},
		})
	})

	tracks, err := c.GetArtistTopTracks(context.Background(), 27, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(tracks) != 2 {
		t.Fatalf("len = %d, want 2", len(tracks))
	}
	if tracks[0].Title != "Get Lucky" {
		t.Errorf("Title = %q, want Get Lucky", tracks[0].Title)
	}
	if tracks[0].Duration != 369 {
		t.Errorf("Duration = %d, want 369", tracks[0].Duration)
	}
}

func TestGetArtistAlbums(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/artist/27/albums" {
			t.Errorf("path = %q, want /artist/27/albums", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 302127, "title": "Random Access Memories"},
				{"id": 302128, "title": "Discovery"},
			},
		})
	})

	albums, err := c.GetArtistAlbums(context.Background(), 27)
	if err != nil {
		t.Fatal(err)
	}
	if len(albums) != 2 {
		t.Fatalf("len = %d, want 2", len(albums))
	}
	if albums[0].Title != "Random Access Memories" {
		t.Errorf("Title = %q", albums[0].Title)
	}
}

func TestGetRelatedArtists(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/artist/27/related" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 100, "name": "Justice"},
			},
		})
	})

	artists, err := c.GetRelatedArtists(context.Background(), 27)
	if err != nil {
		t.Fatal(err)
	}
	if len(artists) != 1 {
		t.Fatalf("len = %d, want 1", len(artists))
	}
	if artists[0].Name != "Justice" {
		t.Errorf("Name = %q, want Justice", artists[0].Name)
	}
}

func TestGetAlbum(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/album/302127" {
			t.Errorf("path = %q, want /album/302127", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":           302127,
			"title":        "Random Access Memories",
			"label":        "Columbia",
			"nb_tracks":    13,
			"fans":         250000,
			"duration":     4384,
			"release_date": "2013-05-17",
			"record_type":  "album",
			"artist":       map[string]any{"id": 27, "name": "Daft Punk"},
			"genres":       map[string]any{"data": []map[string]any{{"id": 113, "name": "Dance"}}},
			"tracks":       map[string]any{"data": []map[string]any{{"id": 3135553, "title": "Give Life Back to Music"}}},
		})
	})

	album, err := c.GetAlbum(context.Background(), 302127)
	if err != nil {
		t.Fatal(err)
	}
	if album.Title != "Random Access Memories" {
		t.Errorf("Title = %q", album.Title)
	}
	if album.Label != "Columbia" {
		t.Errorf("Label = %q, want Columbia", album.Label)
	}
	if album.NbTracks != 13 {
		t.Errorf("NbTracks = %d, want 13", album.NbTracks)
	}
	if album.Duration != 4384 {
		t.Errorf("Duration = %d, want 4384", album.Duration)
	}
	if album.Artist.Name != "Daft Punk" {
		t.Errorf("Artist.Name = %q", album.Artist.Name)
	}
	if len(album.Genres.Data) != 1 {
		t.Fatalf("Genres len = %d, want 1", len(album.Genres.Data))
	}
	if len(album.Tracks.Data) != 1 {
		t.Fatalf("Tracks len = %d, want 1", len(album.Tracks.Data))
	}
}

func TestGetAlbumTracks(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/album/302127/tracks" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 3135553, "title": "Give Life Back to Music", "duration": 274},
				{"id": 3135554, "title": "Get Lucky", "duration": 369},
			},
		})
	})

	tracks, err := c.GetAlbumTracks(context.Background(), 302127)
	if err != nil {
		t.Fatal(err)
	}
	if len(tracks) != 2 {
		t.Fatalf("len = %d, want 2", len(tracks))
	}
	if tracks[1].Title != "Get Lucky" {
		t.Errorf("Title = %q, want Get Lucky", tracks[1].Title)
	}
}

func TestGetTrack(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/track/3135553" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":              3135553,
			"title":           "Get Lucky",
			"title_short":     "Get Lucky",
			"duration":        369,
			"rank":            950000,
			"explicit_lyrics": false,
			"preview":         "https://cdns-preview.example.com/stream",
			"type":            "track",
			"artist":          map[string]any{"id": 27, "name": "Daft Punk"},
			"album":           map[string]any{"id": 302127, "title": "Random Access Memories"},
		})
	})

	track, err := c.GetTrack(context.Background(), 3135553)
	if err != nil {
		t.Fatal(err)
	}
	if track.Title != "Get Lucky" {
		t.Errorf("Title = %q, want Get Lucky", track.Title)
	}
	if track.Duration != 369 {
		t.Errorf("Duration = %d, want 369", track.Duration)
	}
	if track.Rank != 950000 {
		t.Errorf("Rank = %d, want 950000", track.Rank)
	}
	if track.Artist.Name != "Daft Punk" {
		t.Errorf("Artist.Name = %q", track.Artist.Name)
	}
	if track.Album.Title != "Random Access Memories" {
		t.Errorf("Album.Title = %q", track.Album.Title)
	}
}

func TestGetGenres(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/genre" {
			t.Errorf("path = %q, want /genre", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 132, "name": "Pop", "picture": "https://example.com/pop.jpg"},
				{"id": 116, "name": "Rap/Hip Hop"},
				{"id": 152, "name": "Rock"},
			},
		})
	})

	genres, err := c.GetGenres(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(genres) != 3 {
		t.Fatalf("len = %d, want 3", len(genres))
	}
	if genres[0].Name != "Pop" {
		t.Errorf("Name = %q, want Pop", genres[0].Name)
	}
	if genres[0].ID != 132 {
		t.Errorf("ID = %d, want 132", genres[0].ID)
	}
}

func TestGetChart(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chart" {
			t.Errorf("path = %q, want /chart", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"tracks": map[string]any{
				"data": []map[string]any{
					{"id": 1, "title": "Top Track"},
				},
			},
			"albums": map[string]any{
				"data": []map[string]any{
					{"id": 100, "title": "Top Album"},
				},
			},
			"artists": map[string]any{
				"data": []map[string]any{
					{"id": 10, "name": "Top Artist"},
				},
			},
			"playlists": map[string]any{
				"data": []map[string]any{
					{"id": 50, "title": "Top Playlist"},
				},
			},
		})
	})

	chart, err := c.GetChart(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(chart.Tracks.Data) != 1 {
		t.Fatalf("tracks len = %d, want 1", len(chart.Tracks.Data))
	}
	if chart.Tracks.Data[0].Title != "Top Track" {
		t.Errorf("Track Title = %q", chart.Tracks.Data[0].Title)
	}
	if len(chart.Albums.Data) != 1 {
		t.Fatalf("albums len = %d, want 1", len(chart.Albums.Data))
	}
	if len(chart.Artists.Data) != 1 {
		t.Fatalf("artists len = %d, want 1", len(chart.Artists.Data))
	}
	if len(chart.Playlists.Data) != 1 {
		t.Fatalf("playlists len = %d, want 1", len(chart.Playlists.Data))
	}
}

func TestDeezerAPIError(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Deezer returns 200 with error in body.
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"type":    "DataException",
				"message": "no data",
				"code":    800,
			},
		})
	})

	_, err := c.GetArtist(context.Background(), 999999999)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *deezer.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Type != "DataException" {
		t.Errorf("Type = %q, want DataException", apiErr.Type)
	}
	if apiErr.Message != "no data" {
		t.Errorf("Message = %q, want no data", apiErr.Message)
	}
	if apiErr.Code != 800 {
		t.Errorf("Code = %d, want 800", apiErr.Code)
	}
}

func TestHTTPError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}))
	defer srv.Close()

	c := deezer.New(metadata.WithBaseURL(srv.URL))
	_, err := c.GetArtist(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *deezer.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusInternalServerError)
	}
}

func TestWithAccessToken(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("access_token"); got != "my-token" {
			t.Errorf("access_token = %q, want my-token", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{},
		})
	})
	// Recreate with token — need to get the base URL from the existing setup.
	_ = c // unused, we need to pass token manually
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("access_token"); got != "my-token" {
			t.Errorf("access_token = %q, want my-token", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{},
		})
	}))
	defer srv.Close()

	ct := deezer.NewWithToken("my-token", metadata.WithBaseURL(srv.URL))
	_, err := ct.GetGenres(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	custom := &http.Client{}
	c := deezer.New(metadata.WithHTTPClient(custom))
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestAPIErrorFormat(t *testing.T) {
	t.Parallel()

	e := &deezer.APIError{Type: "DataException", Message: "no data", Code: 800}
	if got := e.Error(); got != "deezer: DataException: no data (code 800)" {
		t.Errorf("Error() = %q", got)
	}

	e2 := &deezer.APIError{StatusCode: 500, Status: "500 Internal Server Error", Body: "fail"}
	if got := e2.Error(); got != "deezer: 500 Internal Server Error: fail" {
		t.Errorf("Error() = %q", got)
	}
}
