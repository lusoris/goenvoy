package fanart_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/metadata"
	"github.com/lusoris/goenvoy/metadata/video/fanart"
)

func setup(t *testing.T, handler http.HandlerFunc) *fanart.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return fanart.NewWithClientKey("test-api-key", "test-client-key",
		metadata.WithBaseURL(srv.URL),
	)
}

func respond(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Fatal(err)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew(t *testing.T) {
	c := fanart.New("key")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGetMovieImages(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/movies/120" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Api-Key") != "test-api-key" {
			t.Fatal("missing api-key header")
		}
		if r.Header.Get("Client-Key") != "test-client-key" {
			t.Fatal("missing client-key header")
		}
		respond(t, w, fanart.MovieImages{
			Name:   "The Fellowship of the Ring",
			TMDbID: "120",
			IMDbID: "tt0120737",
			HDMovieLogo: []fanart.Image{
				{ID: "50927", URL: "https://example.com/logo.png", Lang: "en", Likes: "7"},
			},
			MoviePoster: []fanart.Image{
				{ID: "57317", URL: "https://example.com/poster.jpg", Lang: "en", Likes: "4"},
			},
		})
	})

	result, err := c.GetMovieImages(context.Background(), "120")
	assertNoError(t, err)
	if result.Name != "The Fellowship of the Ring" {
		t.Fatalf("got name %q", result.Name)
	}
	if result.TMDbID != "120" {
		t.Fatalf("got tmdb_id %q", result.TMDbID)
	}
	if len(result.HDMovieLogo) != 1 {
		t.Fatalf("got %d hdmovielogos", len(result.HDMovieLogo))
	}
	if result.HDMovieLogo[0].Likes != "7" {
		t.Fatalf("got likes %q", result.HDMovieLogo[0].Likes)
	}
}

func TestGetMovieImagesWithDisc(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respond(t, w, fanart.MovieImages{
			Name: "Test Movie",
			MovieDisc: []fanart.DiscImage{
				{ID: "29003", URL: "https://example.com/disc.png", Lang: "en", Likes: "5", Disc: "1", DiscType: "bluray"},
			},
		})
	})

	result, err := c.GetMovieImages(context.Background(), "120")
	assertNoError(t, err)
	if len(result.MovieDisc) != 1 {
		t.Fatalf("got %d discs", len(result.MovieDisc))
	}
	if result.MovieDisc[0].DiscType != "bluray" {
		t.Fatalf("got disc_type %q", result.MovieDisc[0].DiscType)
	}
}

func TestGetShowImages(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tv/75682" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, fanart.ShowImages{
			Name:      "Bones",
			TheTVDBID: "75682",
			ClearLogo: []fanart.Image{
				{ID: "2112", URL: "https://example.com/logo.png", Lang: "en", Likes: "7"},
			},
			ShowBackground: []fanart.SeasonImage{
				{ID: "19374", URL: "https://example.com/bg.jpg", Lang: "en", Likes: "2", Season: "7"},
			},
			SeasonBanner: []fanart.SeasonImage{
				{ID: "37718", URL: "https://example.com/banner.jpg", Lang: "en", Likes: "0", Season: "1"},
			},
		})
	})

	result, err := c.GetShowImages(context.Background(), "75682")
	assertNoError(t, err)
	if result.Name != "Bones" {
		t.Fatalf("got name %q", result.Name)
	}
	if result.TheTVDBID != "75682" {
		t.Fatalf("got thetvdb_id %q", result.TheTVDBID)
	}
	if len(result.ClearLogo) != 1 {
		t.Fatalf("got %d clearlogos", len(result.ClearLogo))
	}
	if len(result.ShowBackground) != 1 {
		t.Fatalf("got %d showbackgrounds", len(result.ShowBackground))
	}
	if result.ShowBackground[0].Season != "7" {
		t.Fatalf("got season %q", result.ShowBackground[0].Season)
	}
}

func TestGetArtistImages(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/music/f4a31f0a-51dd-4fa7-986d-3095c40c5ed9" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, fanart.ArtistImages{
			Name: "Evanescence",
			MBID: "f4a31f0a-51dd-4fa7-986d-3095c40c5ed9",
			ArtistBackground: []fanart.Image{
				{ID: "6", URL: "https://example.com/bg.jpg", Likes: "4"},
			},
			HDMusicLogo: []fanart.Image{
				{ID: "50850", URL: "https://example.com/logo.png", Likes: "2"},
			},
			Albums: map[string]fanart.AlbumImages{
				"2187d248-1a3b-35d0-a4ec-bead586ff547": {
					AlbumCover: []fanart.Image{
						{ID: "43", URL: "https://example.com/cover.jpg", Likes: "1"},
					},
					CDArt: []fanart.CDArt{
						{ID: "17739", URL: "https://example.com/cd.png", Likes: "0", Disc: "1", Size: "1000"},
					},
				},
			},
		})
	})

	result, err := c.GetArtistImages(context.Background(), "f4a31f0a-51dd-4fa7-986d-3095c40c5ed9")
	assertNoError(t, err)
	if result.Name != "Evanescence" {
		t.Fatalf("got name %q", result.Name)
	}
	if len(result.ArtistBackground) != 1 {
		t.Fatalf("got %d backgrounds", len(result.ArtistBackground))
	}
	if len(result.Albums) != 1 {
		t.Fatalf("got %d albums", len(result.Albums))
	}
	album := result.Albums["2187d248-1a3b-35d0-a4ec-bead586ff547"]
	if len(album.CDArt) != 1 {
		t.Fatalf("got %d cdarts", len(album.CDArt))
	}
	if album.CDArt[0].Size != "1000" {
		t.Fatalf("got size %q", album.CDArt[0].Size)
	}
}

func TestGetAlbumImages(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/music/albums/9ba659df-5814-32f6-b95f-02b738698e7c" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, fanart.AlbumImagesResponse{
			Name: "Evanescence",
			MBID: "f4a31f0a-51dd-4fa7-986d-3095c40c5ed9",
			Albums: map[string]fanart.AlbumImages{
				"9ba659df-5814-32f6-b95f-02b738698e7c": {
					CDArt: []fanart.CDArt{
						{ID: "12420", URL: "https://example.com/cd.png", Likes: "0", Disc: "1", Size: "1000"},
					},
					AlbumCover: []fanart.Image{
						{ID: "116236", URL: "https://example.com/cover.jpg", Likes: "0"},
					},
				},
			},
		})
	})

	result, err := c.GetAlbumImages(context.Background(), "9ba659df-5814-32f6-b95f-02b738698e7c")
	assertNoError(t, err)
	if result.Name != "Evanescence" {
		t.Fatalf("got name %q", result.Name)
	}
	album := result.Albums["9ba659df-5814-32f6-b95f-02b738698e7c"]
	if len(album.AlbumCover) != 1 {
		t.Fatalf("got %d covers", len(album.AlbumCover))
	}
}

func TestGetLabelImages(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/music/labels/e832b688-546b-45e3-83e5-9f8db5dcde1d" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, fanart.LabelImages{
			Name: "Profile Records",
			ID:   "e832b688-546b-45e3-83e5-9f8db5dcde1d",
			MusicLabel: []fanart.LabelImage{
				{ID: "120425", URL: "https://example.com/label.png", Color: "colour", Likes: "0"}, //nolint:misspell // API uses British spelling.
				{ID: "120426", URL: "https://example.com/label2.png", Color: "white", Likes: "0"},
			},
		})
	})

	result, err := c.GetLabelImages(context.Background(), "e832b688-546b-45e3-83e5-9f8db5dcde1d")
	assertNoError(t, err)
	if result.Name != "Profile Records" {
		t.Fatalf("got name %q", result.Name)
	}
	if len(result.MusicLabel) != 2 {
		t.Fatalf("got %d labels", len(result.MusicLabel))
	}
	if result.MusicLabel[0].Color != "colour" { //nolint:misspell // API uses British spelling.
		t.Fatalf("got color %q", result.MusicLabel[0].Color)
	}
}

func TestGetLatestMovies(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/movies/latest" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []fanart.LatestMovie{
			{TMDbID: "17610", IMDbID: "tt1045778", Name: "Year One", NewImages: "1", TotalImages: "7"},
			{TMDbID: "5651", IMDbID: "tt0075376", Name: "Up!", NewImages: "1", TotalImages: "1"},
		})
	})

	result, err := c.GetLatestMovies(context.Background(), 0)
	assertNoError(t, err)
	if len(result) != 2 {
		t.Fatalf("got %d movies", len(result))
	}
	if result[0].Name != "Year One" {
		t.Fatalf("got name %q", result[0].Name)
	}
}

func TestGetLatestMoviesWithDate(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("date") != "1609459200" {
			t.Fatalf("unexpected date param: %s", r.URL.Query().Get("date"))
		}
		respond(t, w, []fanart.LatestMovie{})
	})

	result, err := c.GetLatestMovies(context.Background(), 1609459200)
	assertNoError(t, err)
	if len(result) != 0 {
		t.Fatalf("got %d movies", len(result))
	}
}

func TestGetLatestShows(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tv/latest" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []fanart.LatestShow{
			{ID: "79473", Name: "Witchblade", NewImages: "2", TotalImages: "7"},
		})
	})

	result, err := c.GetLatestShows(context.Background(), 0)
	assertNoError(t, err)
	if len(result) != 1 {
		t.Fatalf("got %d shows", len(result))
	}
	if result[0].ID != "79473" {
		t.Fatalf("got id %q", result[0].ID)
	}
}

func TestGetLatestArtists(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/music/latest" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []fanart.LatestArtist{
			{ID: "932cb13f-677c-4e63-beb1-fca45574d280", Name: "Zazie", NewImages: "1", TotalImages: "23"},
		})
	})

	result, err := c.GetLatestArtists(context.Background(), 0)
	assertNoError(t, err)
	if len(result) != 1 {
		t.Fatalf("got %d artists", len(result))
	}
	if result[0].Name != "Zazie" {
		t.Fatalf("got name %q", result[0].Name)
	}
}

func TestAPIError(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"status":"error","error message":"Not found"}`))
	})

	_, err := c.GetMovieImages(context.Background(), "999999")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *fanart.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *fanart.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("got status %d", apiErr.StatusCode)
	}
	if apiErr.ErrorMessage != "Not found" {
		t.Fatalf("got message %q", apiErr.ErrorMessage)
	}
}

func TestOptions(t *testing.T) {
	c := fanart.NewWithClientKey("key", "ck",
		metadata.WithUserAgent("custom-agent"),
		metadata.WithTimeout(60_000_000_000),
	)
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNoClientKeyHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Client-Key") != "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		respond(&testing.T{}, w, fanart.MovieImages{Name: "Test"})
	}))
	defer srv.Close()

	c := fanart.New("key", metadata.WithBaseURL(srv.URL))
	result, err := c.GetMovieImages(context.Background(), "1")
	assertNoError(t, err)
	if result.Name != "Test" {
		t.Fatalf("got name %q", result.Name)
	}
}
