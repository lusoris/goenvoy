package navidrome

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// helper wraps a responseBody in the subsonic-response envelope.
func subsonicJSON(t *testing.T, w http.ResponseWriter, rb *responseBody) {
	t.Helper()
	rb.Status = "ok"
	rb.Version = "1.16.1"
	json.NewEncoder(w).Encode(subsonicResponse{Response: *rb})
}

func newTestServer(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return New(ts.URL, "admin", "password")
}

func TestPing(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("u") != "admin" {
			http.Error(w, "bad auth", http.StatusUnauthorized)
			return
		}
		subsonicJSON(t, w, &responseBody{})
	})

	if err := c.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestGetArtists(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			Artists: &ArtistsID3{
				Index: []IndexID3{{Name: "A", Artist: []ArtistID3{{ID: "1", Name: "ABBA"}}}},
			},
		})
	})

	artists, err := c.GetArtists(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(artists.Index) != 1 || artists.Index[0].Artist[0].Name != "ABBA" {
		t.Fatalf("unexpected artists: %+v", artists)
	}
}

func TestGetArtist(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			Artist: &ArtistID3{ID: "1", Name: "Radiohead", AlbumCount: 9},
		})
	})

	a, err := c.GetArtist(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "Radiohead" {
		t.Fatalf("unexpected artist: %+v", a)
	}
}

func TestGetAlbum(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			Album: &AlbumID3{ID: "a1", Name: "OK Computer", SongCount: 12},
		})
	})

	album, err := c.GetAlbum(context.Background(), "a1")
	if err != nil {
		t.Fatal(err)
	}
	if album.Name != "OK Computer" {
		t.Fatalf("unexpected album: %+v", album)
	}
}

func TestGetAlbumList2(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			AlbumList2: &AlbumList2{Album: []AlbumID3{{ID: "a1", Name: "Abbey Road"}}},
		})
	})

	albums, err := c.GetAlbumList2(context.Background(), "newest", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(albums) != 1 || albums[0].Name != "Abbey Road" {
		t.Fatalf("unexpected albums: %+v", albums)
	}
}

func TestGetSong(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			Song: &Song{ID: "s1", Title: "Paranoid Android"},
		})
	})

	song, err := c.GetSong(context.Background(), "s1")
	if err != nil {
		t.Fatal(err)
	}
	if song.Title != "Paranoid Android" {
		t.Fatalf("unexpected song: %+v", song)
	}
}

func TestGetRandomSongs(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			RandomSongs: &Songs{Song: []Song{{ID: "s1", Title: "Creep"}}},
		})
	})

	songs, err := c.GetRandomSongs(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(songs) != 1 || songs[0].Title != "Creep" {
		t.Fatalf("unexpected songs: %+v", songs)
	}
}

func TestSearch3(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			SearchResult: &SearchResult3{
				Artist: []ArtistID3{{ID: "1", Name: "Beatles"}},
				Album:  []AlbumID3{{ID: "a1", Name: "Help!"}},
				Song:   []Song{{ID: "s1", Title: "Yesterday"}},
			},
		})
	})

	res, err := c.Search3(context.Background(), "beatles", 5, 5, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Artist) != 1 || res.Artist[0].Name != "Beatles" {
		t.Fatalf("unexpected search results: %+v", res)
	}
}

func TestGetPlaylists(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			Playlists: &Playlists{Playlist: []Playlist{{ID: "p1", Name: "Favorites"}}},
		})
	})

	playlists, err := c.GetPlaylists(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(playlists) != 1 || playlists[0].Name != "Favorites" {
		t.Fatalf("unexpected playlists: %+v", playlists)
	}
}

func TestGetPlaylist(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			Playlist: &Playlist{ID: "p1", Name: "Road Trip", SongCount: 20},
		})
	})

	pl, err := c.GetPlaylist(context.Background(), "p1")
	if err != nil {
		t.Fatal(err)
	}
	if pl.Name != "Road Trip" {
		t.Fatalf("unexpected playlist: %+v", pl)
	}
}

func TestGetGenres(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			Genres: &Genres{Genre: []Genre{{Value: "Rock", SongCount: 100, AlbumCount: 20}}},
		})
	})

	genres, err := c.GetGenres(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(genres) != 1 || genres[0].Value != "Rock" {
		t.Fatalf("unexpected genres: %+v", genres)
	}
}

func TestGetScanStatus(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			ScanStatus: &ScanStatus{Scanning: false, Count: 5000},
		})
	})

	ss, err := c.GetScanStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if ss.Scanning || ss.Count != 5000 {
		t.Fatalf("unexpected scan status: %+v", ss)
	}
}

func TestGetStarred2(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		subsonicJSON(t, w, &responseBody{
			Starred2: &Starred2{
				Song: []Song{{ID: "s1", Title: "Bohemian Rhapsody"}},
			},
		})
	})

	starred, err := c.GetStarred2(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(starred.Song) != 1 || starred.Song[0].Title != "Bohemian Rhapsody" {
		t.Fatalf("unexpected starred: %+v", starred)
	}
}

func TestSubsonicError(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(subsonicResponse{
			Response: responseBody{
				Status:  "failed",
				Version: "1.16.1",
				Error:   &SubsonicError{Code: 40, Message: "Wrong username or password"},
			},
		})
	})

	err := c.Ping(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var se *SubsonicError
	if !errors.As(err, &se) {
		t.Fatalf("expected *SubsonicError, got %T", err)
	}
	if se.Code != 40 {
		t.Fatalf("unexpected error code: %d", se.Code)
	}
}

func TestAPIError(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
	})

	err := c.Ping(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("unexpected status code: %d", apiErr.StatusCode)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c := New("http://localhost", "u", "p", WithHTTPClient(custom))
	if c.http != custom {
		t.Fatal("custom HTTP client not set")
	}
}
