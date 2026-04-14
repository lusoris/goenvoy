package anilist_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/anime/anilist"
)

func setup(t *testing.T, handler http.HandlerFunc) *anilist.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return anilist.New(metadata.WithBaseURL(srv.URL))
}

func respondGraphQL(t *testing.T, w http.ResponseWriter, data any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]any{"data": data}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		t.Fatal(err)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func decodeRequest(t *testing.T, r *http.Request) map[string]any {
	t.Helper()
	var req map[string]any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.Fatal(err)
	}
	return req
}

func TestNew(t *testing.T) {
	t.Parallel()

	c := anilist.New()
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestQuery(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("unexpected content-type: %s", r.Header.Get("Content-Type"))
		}
		respondGraphQL(t, w, map[string]any{
			"Media": map[string]any{"id": float64(1), "title": map[string]any{"romaji": "Test"}},
		})
	})

	var resp struct {
		Media struct {
			ID    int `json:"id"`
			Title struct {
				Romaji string `json:"romaji"`
			} `json:"title"`
		} `json:"Media"`
	}
	err := c.Query(context.Background(), `query ($id: Int) { Media(id: $id) { id title { romaji } } }`, map[string]any{"id": 1}, &resp)
	assertNoError(t, err)
	if resp.Media.Title.Romaji != "Test" {
		t.Fatalf("got title %q", resp.Media.Title.Romaji)
	}
}

func TestGetMedia(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		req := decodeRequest(t, r)
		vars, _ := req["variables"].(map[string]any)
		if vars["id"] != float64(1) {
			t.Fatalf("unexpected id: %v", vars["id"])
		}
		respondGraphQL(t, w, map[string]any{
			"Media": map[string]any{
				"id":     float64(1),
				"idMal":  float64(1),
				"title":  map[string]any{"romaji": "Cowboy Bebop", "english": "Cowboy Bebop"},
				"type":   "ANIME",
				"format": "TV",
				"status": "FINISHED",
				"genres": []string{"Action", "Sci-Fi"},
			},
		})
	})

	media, err := c.GetMedia(context.Background(), 1)
	assertNoError(t, err)
	if media.Title.Romaji != "Cowboy Bebop" {
		t.Fatalf("got title %q", media.Title.Romaji)
	}
	if media.Type != "ANIME" {
		t.Fatalf("got type %q", media.Type)
	}
	if len(media.Genres) != 2 {
		t.Fatalf("got %d genres", len(media.Genres))
	}
}

func TestGetMediaByMalID(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		req := decodeRequest(t, r)
		vars, _ := req["variables"].(map[string]any)
		if vars["idMal"] != float64(1) {
			t.Fatalf("unexpected idMal: %v", vars["idMal"])
		}
		if vars["type"] != "ANIME" {
			t.Fatalf("unexpected type: %v", vars["type"])
		}
		respondGraphQL(t, w, map[string]any{
			"Media": map[string]any{
				"id":    float64(1),
				"idMal": float64(1),
				"title": map[string]any{"romaji": "Cowboy Bebop"},
			},
		})
	})

	media, err := c.GetMediaByMalID(context.Background(), 1, anilist.MediaTypeAnime)
	assertNoError(t, err)
	if media.Title.Romaji != "Cowboy Bebop" {
		t.Fatalf("got title %q", media.Title.Romaji)
	}
}

func TestSearchMedia(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		req := decodeRequest(t, r)
		vars, _ := req["variables"].(map[string]any)
		if vars["search"] != "cowboy" {
			t.Fatalf("unexpected search: %v", vars["search"])
		}
		if vars["type"] != "ANIME" {
			t.Fatalf("unexpected type: %v", vars["type"])
		}
		respondGraphQL(t, w, map[string]any{
			"Page": map[string]any{
				"pageInfo": map[string]any{
					"currentPage": float64(1),
					"hasNextPage": true,
					"perPage":     float64(10),
				},
				"media": []any{
					map[string]any{
						"id":    float64(1),
						"title": map[string]any{"romaji": "Cowboy Bebop"},
						"type":  "ANIME",
					},
				},
			},
		})
	})

	page, err := c.SearchMedia(context.Background(), "cowboy", anilist.MediaTypeAnime, 1, 10)
	assertNoError(t, err)
	if !page.PageInfo.HasNextPage {
		t.Fatal("expected hasNextPage=true")
	}
	if len(page.Media) != 1 {
		t.Fatalf("got %d media", len(page.Media))
	}
	if page.Media[0].Title.Romaji != "Cowboy Bebop" {
		t.Fatalf("got title %q", page.Media[0].Title.Romaji)
	}
}

func TestGetCharacter(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		req := decodeRequest(t, r)
		vars, _ := req["variables"].(map[string]any)
		if vars["id"] != float64(1) {
			t.Fatalf("unexpected id: %v", vars["id"])
		}
		respondGraphQL(t, w, map[string]any{
			"Character": map[string]any{
				"id":          float64(1),
				"name":        map[string]any{"full": "Spike Spiegel"},
				"gender":      "Male",
				"age":         "27",
				"description": "A bounty hunter.",
			},
		})
	})

	ch, err := c.GetCharacter(context.Background(), 1)
	assertNoError(t, err)
	if ch.Name.Full != "Spike Spiegel" {
		t.Fatalf("got name %q", ch.Name.Full)
	}
	if ch.Gender != "Male" {
		t.Fatalf("got gender %q", ch.Gender)
	}
}

func TestSearchCharacters(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondGraphQL(t, w, map[string]any{
			"Page": map[string]any{
				"pageInfo": map[string]any{"hasNextPage": false, "currentPage": float64(1)},
				"characters": []any{
					map[string]any{"id": float64(1), "name": map[string]any{"full": "Spike Spiegel"}},
				},
			},
		})
	})

	page, err := c.SearchCharacters(context.Background(), "spike", 1, 10)
	assertNoError(t, err)
	if len(page.Characters) != 1 {
		t.Fatalf("got %d characters", len(page.Characters))
	}
}

func TestGetStaff(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondGraphQL(t, w, map[string]any{
			"Staff": map[string]any{
				"id":       float64(95508),
				"name":     map[string]any{"full": "Youko Kanno"},
				"homeTown": "Sendai, Miyagi, Japan",
			},
		})
	})

	staff, err := c.GetStaff(context.Background(), 95508)
	assertNoError(t, err)
	if staff.Name.Full != "Youko Kanno" {
		t.Fatalf("got name %q", staff.Name.Full)
	}
	if staff.HomeTown != "Sendai, Miyagi, Japan" {
		t.Fatalf("got hometown %q", staff.HomeTown)
	}
}

func TestSearchStaff(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondGraphQL(t, w, map[string]any{
			"Page": map[string]any{
				"pageInfo": map[string]any{"hasNextPage": false},
				"staff": []any{
					map[string]any{"id": float64(95508), "name": map[string]any{"full": "Youko Kanno"}},
				},
			},
		})
	})

	page, err := c.SearchStaff(context.Background(), "kanno", 1, 10)
	assertNoError(t, err)
	if len(page.Staff) != 1 {
		t.Fatalf("got %d staff", len(page.Staff))
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		req := decodeRequest(t, r)
		vars, _ := req["variables"].(map[string]any)
		if vars["id"] != float64(1) {
			t.Fatalf("unexpected id: %v", vars["id"])
		}
		respondGraphQL(t, w, map[string]any{
			"User": map[string]any{
				"id":   float64(1),
				"name": "TestUser",
			},
		})
	})

	user, err := c.GetUser(context.Background(), 1)
	assertNoError(t, err)
	if user.Name != "TestUser" {
		t.Fatalf("got name %q", user.Name)
	}
}

func TestGetUserByName(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		req := decodeRequest(t, r)
		vars, _ := req["variables"].(map[string]any)
		if vars["name"] != "TestUser" {
			t.Fatalf("unexpected name: %v", vars["name"])
		}
		respondGraphQL(t, w, map[string]any{
			"User": map[string]any{
				"id":   float64(1),
				"name": "TestUser",
			},
		})
	})

	user, err := c.GetUserByName(context.Background(), "TestUser")
	assertNoError(t, err)
	if user.ID != 1 {
		t.Fatalf("got id %d", user.ID)
	}
}

func TestGetGenres(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondGraphQL(t, w, map[string]any{
			"GenreCollection": []string{"Action", "Comedy", "Drama"},
		})
	})

	genres, err := c.GetGenres(context.Background())
	assertNoError(t, err)
	if len(genres) != 3 {
		t.Fatalf("got %d genres", len(genres))
	}
	if genres[0] != "Action" {
		t.Fatalf("got genre %q", genres[0])
	}
}

func TestGetTags(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		respondGraphQL(t, w, map[string]any{
			"MediaTagCollection": []any{
				map[string]any{"id": float64(1), "name": "Action", "description": "Combat scenes."},
			},
		})
	})

	tags, err := c.GetTags(context.Background())
	assertNoError(t, err)
	if len(tags) != 1 {
		t.Fatalf("got %d tags", len(tags))
	}
	if tags[0].Name != "Action" {
		t.Fatalf("got tag %q", tags[0].Name)
	}
}

func TestGraphQLError(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"data":   nil,
			"errors": []map[string]any{{"message": "Not Found.", "status": float64(404)}},
		}); err != nil {
			t.Fatal(err)
		}
	})

	_, err := c.GetMedia(context.Background(), 99999)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *anilist.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Errors[0].Status != 404 {
		t.Fatalf("got status %d", apiErr.Errors[0].Status)
	}
}

func TestHTTPError(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	})

	_, err := c.GetMedia(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
	var httpErr *anilist.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected HTTPError, got %T: %v", err, err)
	}
	if httpErr.StatusCode != http.StatusInternalServerError {
		t.Fatalf("got status %d", httpErr.StatusCode)
	}
}

func TestWithAccessToken(t *testing.T) {
	t.Parallel()

	var gotAuth string
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		respondGraphQL(t, w, map[string]any{
			"Media": map[string]any{"id": float64(1)},
		})
	})
	// The setup function doesn't pass WithAccessToken, so we need to create a new client.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		respondGraphQL(t, w, map[string]any{
			"Media": map[string]any{"id": float64(1)},
		})
	}))
	t.Cleanup(srv.Close)

	_ = c // suppress unused

	cl := anilist.NewWithToken("my-tok", metadata.WithBaseURL(srv.URL))
	_, err := cl.GetMedia(context.Background(), 1)
	assertNoError(t, err)
	if gotAuth != "Bearer my-tok" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer my-tok")
	}
}

func TestNoAccessTokenHeader(t *testing.T) {
	t.Parallel()

	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		respondGraphQL(t, w, map[string]any{
			"Media": map[string]any{"id": float64(1)},
		})
	}))
	t.Cleanup(srv.Close)

	cl := anilist.New(metadata.WithBaseURL(srv.URL))
	_, err := cl.GetMedia(context.Background(), 1)
	assertNoError(t, err)
	if gotAuth != "" {
		t.Errorf("Authorization = %q, want empty", gotAuth)
	}
}
