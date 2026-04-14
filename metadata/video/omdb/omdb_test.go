package omdb_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/video/omdb"
)

func setup(t *testing.T, handler http.HandlerFunc) *omdb.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return omdb.New("test-key", metadata.WithBaseURL(srv.URL))
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
	t.Parallel()

	c := omdb.New("key")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestGetByIMDbID(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("i") != "tt0111161" {
			t.Fatalf("unexpected i param: %s", r.URL.Query().Get("i"))
		}
		if r.URL.Query().Get("apikey") != "test-key" {
			t.Fatal("missing apikey param")
		}
		respond(t, w, omdb.Title{
			Title:      "The Shawshank Redemption",
			Year:       "1994",
			IMDbID:     "tt0111161",
			Type:       "movie",
			IMDbRating: "9.3",
			Response:   "True",
			Ratings: []omdb.Rating{
				{Source: "Internet Movie Database", Value: "9.3/10"},
			},
		})
	})

	result, err := c.GetByIMDbID(context.Background(), "tt0111161", "")
	assertNoError(t, err)
	if result.Title != "The Shawshank Redemption" {
		t.Fatalf("got title %q", result.Title)
	}
	if result.IMDbRating != "9.3" {
		t.Fatalf("got rating %q", result.IMDbRating)
	}
	if len(result.Ratings) != 1 {
		t.Fatalf("got %d ratings", len(result.Ratings))
	}
}

func TestGetByIMDbIDWithPlot(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("plot") != "full" {
			t.Fatalf("unexpected plot param: %s", r.URL.Query().Get("plot"))
		}
		respond(t, w, omdb.Title{
			Title:    "Test",
			Plot:     "A very long plot description...",
			Response: "True",
		})
	})

	result, err := c.GetByIMDbID(context.Background(), "tt0111161", omdb.PlotFull)
	assertNoError(t, err)
	if result.Plot != "A very long plot description..." {
		t.Fatalf("got plot %q", result.Plot)
	}
}

func TestGetByTitle(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("t") != "Inception" {
			t.Fatalf("unexpected t param: %s", r.URL.Query().Get("t"))
		}
		respond(t, w, omdb.Title{
			Title:    "Inception",
			Year:     "2010",
			IMDbID:   "tt1375666",
			Type:     "movie",
			Response: "True",
		})
	})

	result, err := c.GetByTitle(context.Background(), "Inception", 0, "", "")
	assertNoError(t, err)
	if result.Title != "Inception" {
		t.Fatalf("got title %q", result.Title)
	}
}

func TestGetByTitleWithFilters(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("y") != "2010" {
			t.Fatalf("unexpected y param: %s", r.URL.Query().Get("y"))
		}
		if r.URL.Query().Get("type") != "movie" {
			t.Fatalf("unexpected type param: %s", r.URL.Query().Get("type"))
		}
		respond(t, w, omdb.Title{
			Title:    "Inception",
			Response: "True",
		})
	})

	result, err := c.GetByTitle(context.Background(), "Inception", 2010, omdb.MediaTypeMovie, "")
	assertNoError(t, err)
	if result.Title != "Inception" {
		t.Fatalf("got title %q", result.Title)
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("s") != "Batman" {
			t.Fatalf("unexpected s param: %s", r.URL.Query().Get("s"))
		}
		respond(t, w, omdb.SearchResponse{
			Search: []omdb.SearchResult{
				{Title: "Batman Begins", Year: "2005", IMDbID: "tt0372784", Type: "movie"},
				{Title: "The Dark Knight", Year: "2008", IMDbID: "tt0468569", Type: "movie"},
			},
			TotalResults: "530",
			Response:     "True",
		})
	})

	result, err := c.Search(context.Background(), "Batman", 0, "", 0)
	assertNoError(t, err)
	if len(result.Search) != 2 {
		t.Fatalf("got %d results", len(result.Search))
	}
	if result.TotalResults != "530" {
		t.Fatalf("got totalResults %q", result.TotalResults)
	}
}

func TestSearchWithFilters(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("y") != "2005" {
			t.Fatalf("unexpected y param: %s", r.URL.Query().Get("y"))
		}
		if r.URL.Query().Get("type") != "series" {
			t.Fatalf("unexpected type param: %s", r.URL.Query().Get("type"))
		}
		if r.URL.Query().Get("page") != "2" {
			t.Fatalf("unexpected page param: %s", r.URL.Query().Get("page"))
		}
		respond(t, w, omdb.SearchResponse{
			Search:       []omdb.SearchResult{},
			TotalResults: "0",
			Response:     "True",
		})
	})

	result, err := c.Search(context.Background(), "Batman", 2005, omdb.MediaTypeSeries, 2)
	assertNoError(t, err)
	if len(result.Search) != 0 {
		t.Fatalf("got %d results", len(result.Search))
	}
}

func TestGetSeason(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("i") != "tt0944947" {
			t.Fatalf("unexpected i param: %s", r.URL.Query().Get("i"))
		}
		if r.URL.Query().Get("Season") != "1" {
			t.Fatalf("unexpected Season param: %s", r.URL.Query().Get("Season"))
		}
		respond(t, w, omdb.SeasonResponse{
			Title:        "Game of Thrones",
			Season:       "1",
			TotalSeasons: "8",
			Episodes: []omdb.Episode{
				{Title: "Winter Is Coming", Released: "2011-04-17", Episode: "1", IMDbRating: "9.1", IMDbID: "tt1480055"},
				{Title: "The Kingsroad", Released: "2011-04-24", Episode: "2", IMDbRating: "8.8", IMDbID: "tt1668746"},
			},
			Response: "True",
		})
	})

	result, err := c.GetSeason(context.Background(), "tt0944947", 1)
	assertNoError(t, err)
	if result.Title != "Game of Thrones" {
		t.Fatalf("got title %q", result.Title)
	}
	if result.TotalSeasons != "8" {
		t.Fatalf("got totalSeasons %q", result.TotalSeasons)
	}
	if len(result.Episodes) != 2 {
		t.Fatalf("got %d episodes", len(result.Episodes))
	}
	if result.Episodes[0].Title != "Winter Is Coming" {
		t.Fatalf("got episode title %q", result.Episodes[0].Title)
	}
}

func TestGetEpisode(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("Season") != "1" {
			t.Fatalf("unexpected Season param: %s", r.URL.Query().Get("Season"))
		}
		if r.URL.Query().Get("Episode") != "1" {
			t.Fatalf("unexpected Episode param: %s", r.URL.Query().Get("Episode"))
		}
		respond(t, w, omdb.Title{
			Title:    "Winter Is Coming",
			Year:     "2011",
			IMDbID:   "tt1480055",
			Type:     "episode",
			Response: "True",
		})
	})

	result, err := c.GetEpisode(context.Background(), "tt0944947", 1, 1)
	assertNoError(t, err)
	if result.Title != "Winter Is Coming" {
		t.Fatalf("got title %q", result.Title)
	}
	if result.Type != "episode" {
		t.Fatalf("got type %q", result.Type)
	}
}

func TestAPIErrorResponse(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Response":"False","Error":"Movie not found!"}`))
	})

	_, err := c.GetByIMDbID(context.Background(), "tt9999999", "")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *omdb.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *omdb.APIError, got %T", err)
	}
	if apiErr.Message != "Movie not found!" {
		t.Fatalf("got message %q", apiErr.Message)
	}
}

func TestHTTPError(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := c.GetByIMDbID(context.Background(), "tt0111161", "")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *omdb.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *omdb.APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("got status %d", apiErr.StatusCode)
	}
}

func TestOptions(t *testing.T) {
	t.Parallel()

	c := omdb.New("key",
		metadata.WithUserAgent("custom-agent"),
		metadata.WithTimeout(60_000_000_000),
	)
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}
