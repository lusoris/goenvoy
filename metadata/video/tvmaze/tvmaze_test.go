package tvmaze_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/video/tvmaze"
)

func setup(t *testing.T, handler http.HandlerFunc) *tvmaze.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return tvmaze.New(metadata.WithBaseURL(srv.URL))
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

func intPtr(i int) *int { return &i }

func TestNew(t *testing.T) {
	t.Parallel()

	c := tvmaze.New()
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestSearchShows(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/shows" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "girls" {
			t.Fatalf("unexpected query: %s", r.URL.Query().Get("q"))
		}
		respond(t, w, []tvmaze.SearchShowResult{
			{Score: 0.9, Show: tvmaze.Show{ID: 139, Name: "Girls"}},
		})
	})

	results, err := c.SearchShows(context.Background(), "girls")
	assertNoError(t, err)
	if len(results) != 1 {
		t.Fatalf("got %d results", len(results))
	}
	if results[0].Show.Name != "Girls" {
		t.Fatalf("got name %q", results[0].Show.Name)
	}
}

func TestSearchShowSingle(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/singlesearch/shows" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, tvmaze.Show{ID: 139, Name: "Girls"})
	})

	show, err := c.SearchShowSingle(context.Background(), "girls")
	assertNoError(t, err)
	if show.Name != "Girls" {
		t.Fatalf("got name %q", show.Name)
	}
}

func TestLookupShowByTheTVDB(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/lookup/shows" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("thetvdb") != "264492" {
			t.Fatalf("unexpected thetvdb: %s", r.URL.Query().Get("thetvdb"))
		}
		respond(t, w, tvmaze.Show{ID: 1, Name: "Under the Dome"})
	})

	show, err := c.LookupShowByTheTVDB(context.Background(), 264492)
	assertNoError(t, err)
	if show.Name != "Under the Dome" {
		t.Fatalf("got name %q", show.Name)
	}
}

func TestLookupShowByIMDB(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("imdb") != "tt1553656" {
			t.Fatalf("unexpected imdb: %s", r.URL.Query().Get("imdb"))
		}
		respond(t, w, tvmaze.Show{ID: 1, Name: "Under the Dome"})
	})

	show, err := c.LookupShowByIMDB(context.Background(), "tt1553656")
	assertNoError(t, err)
	if show.ID != 1 {
		t.Fatalf("got id %d", show.ID)
	}
}

func TestLookupShowByTVRage(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tvrage") != "25988" {
			t.Fatalf("unexpected tvrage: %s", r.URL.Query().Get("tvrage"))
		}
		respond(t, w, tvmaze.Show{ID: 1, Name: "Under the Dome"})
	})

	show, err := c.LookupShowByTVRage(context.Background(), 25988)
	assertNoError(t, err)
	if show.ID != 1 {
		t.Fatalf("got id %d", show.ID)
	}
}

func TestSearchPeople(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/people" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []tvmaze.SearchPersonResult{
			{Score: 0.8, Person: tvmaze.Person{ID: 1, Name: "Mike Vogel"}},
		})
	})

	results, err := c.SearchPeople(context.Background(), "mike")
	assertNoError(t, err)
	if len(results) != 1 {
		t.Fatalf("got %d results", len(results))
	}
	if results[0].Person.Name != "Mike Vogel" {
		t.Fatalf("got name %q", results[0].Person.Name)
	}
}

func TestGetShow(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows/1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, tvmaze.Show{
			ID:     1,
			Name:   "Under the Dome",
			Status: "Ended",
			Genres: []string{"Drama", "Thriller"},
		})
	})

	show, err := c.GetShow(context.Background(), 1)
	assertNoError(t, err)
	if show.Status != "Ended" {
		t.Fatalf("got status %q", show.Status)
	}
	if len(show.Genres) != 2 {
		t.Fatalf("got %d genres", len(show.Genres))
	}
}

func TestGetShowEpisodes(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows/1/episodes" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []tvmaze.Episode{
			{ID: 1, Name: "Pilot", Season: 1, Number: intPtr(1)},
			{ID: 2, Name: "The Fire", Season: 1, Number: intPtr(2)},
		})
	})

	episodes, err := c.GetShowEpisodes(context.Background(), 1)
	assertNoError(t, err)
	if len(episodes) != 2 {
		t.Fatalf("got %d episodes", len(episodes))
	}
	if episodes[0].Name != "Pilot" {
		t.Fatalf("got name %q", episodes[0].Name)
	}
}

func TestGetEpisodeByNumber(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows/1/episodebynumber" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("season") != "1" || r.URL.Query().Get("number") != "1" {
			t.Fatal("missing season/number params")
		}
		respond(t, w, tvmaze.Episode{ID: 1, Name: "Pilot", Season: 1, Number: intPtr(1)})
	})

	ep, err := c.GetEpisodeByNumber(context.Background(), 1, 1, 1)
	assertNoError(t, err)
	if ep.Name != "Pilot" {
		t.Fatalf("got name %q", ep.Name)
	}
}

func TestGetEpisodesByDate(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows/1/episodesbydate" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("date") != "2013-06-24" {
			t.Fatalf("unexpected date: %s", r.URL.Query().Get("date"))
		}
		respond(t, w, []tvmaze.Episode{
			{ID: 1, Name: "Pilot"},
		})
	})

	episodes, err := c.GetEpisodesByDate(context.Background(), 1, "2013-06-24")
	assertNoError(t, err)
	if len(episodes) != 1 {
		t.Fatalf("got %d episodes", len(episodes))
	}
}

func TestGetShowSeasons(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows/1/seasons" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []tvmaze.Season{
			{ID: 1, Number: 1, EpisodeOrder: intPtr(13)},
			{ID: 2, Number: 2, EpisodeOrder: intPtr(13)},
		})
	})

	seasons, err := c.GetShowSeasons(context.Background(), 1)
	assertNoError(t, err)
	if len(seasons) != 2 {
		t.Fatalf("got %d seasons", len(seasons))
	}
}

func TestGetSeasonEpisodes(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/seasons/1/episodes" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []tvmaze.Episode{
			{ID: 1, Name: "Pilot"},
		})
	})

	episodes, err := c.GetSeasonEpisodes(context.Background(), 1)
	assertNoError(t, err)
	if len(episodes) != 1 {
		t.Fatalf("got %d episodes", len(episodes))
	}
}

func TestGetShowCast(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows/1/cast" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []tvmaze.CastMember{
			{
				Person:    tvmaze.Person{ID: 1, Name: "Mike Vogel"},
				Character: tvmaze.Character{ID: 1, Name: "Dale Barbara"},
			},
		})
	})

	cast, err := c.GetShowCast(context.Background(), 1)
	assertNoError(t, err)
	if len(cast) != 1 {
		t.Fatalf("got %d cast members", len(cast))
	}
	if cast[0].Person.Name != "Mike Vogel" {
		t.Fatalf("got person %q", cast[0].Person.Name)
	}
	if cast[0].Character.Name != "Dale Barbara" {
		t.Fatalf("got character %q", cast[0].Character.Name)
	}
}

func TestGetShowCrew(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows/1/crew" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []tvmaze.CrewMember{
			{Type: "Creator", Person: tvmaze.Person{ID: 15, Name: "Stephen King"}},
		})
	})

	crew, err := c.GetShowCrew(context.Background(), 1)
	assertNoError(t, err)
	if len(crew) != 1 {
		t.Fatalf("got %d crew", len(crew))
	}
	if crew[0].Type != "Creator" {
		t.Fatalf("got type %q", crew[0].Type)
	}
}

func TestGetShowAKAs(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows/1/akas" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []tvmaze.AKA{
			{Name: "Unter der Kuppel", Country: &tvmaze.Country{Name: "Germany", Code: "DE"}},
		})
	})

	akas, err := c.GetShowAKAs(context.Background(), 1)
	assertNoError(t, err)
	if len(akas) != 1 {
		t.Fatalf("got %d akas", len(akas))
	}
	if akas[0].Name != "Unter der Kuppel" {
		t.Fatalf("got name %q", akas[0].Name)
	}
}

func TestGetShowImages(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows/1/images" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []tvmaze.ShowImage{
			{
				ID:   1,
				Type: "poster",
				Main: true,
				Resolutions: map[string]tvmaze.ImageResolution{
					"original": {URL: "https://example.com/poster.jpg", Width: 680, Height: 1000},
				},
			},
		})
	})

	images, err := c.GetShowImages(context.Background(), 1)
	assertNoError(t, err)
	if len(images) != 1 {
		t.Fatalf("got %d images", len(images))
	}
	if images[0].Type != "poster" {
		t.Fatalf("got type %q", images[0].Type)
	}
	if !images[0].Main {
		t.Fatal("expected main=true")
	}
}

func TestGetShowIndex(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/shows" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "0" {
			t.Fatalf("unexpected page: %s", r.URL.Query().Get("page"))
		}
		respond(t, w, []tvmaze.Show{
			{ID: 1, Name: "Under the Dome"},
		})
	})

	shows, err := c.GetShowIndex(context.Background(), 0)
	assertNoError(t, err)
	if len(shows) != 1 {
		t.Fatalf("got %d shows", len(shows))
	}
}

func TestGetEpisode(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/episodes/1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, tvmaze.Episode{ID: 1, Name: "Pilot", Season: 1, Number: intPtr(1)})
	})

	ep, err := c.GetEpisode(context.Background(), 1)
	assertNoError(t, err)
	if ep.Name != "Pilot" {
		t.Fatalf("got name %q", ep.Name)
	}
}

func TestGetPerson(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/people/1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, tvmaze.Person{
			ID:       1,
			Name:     "Mike Vogel",
			Birthday: "1979-07-17",
			Gender:   "Male",
		})
	})

	person, err := c.GetPerson(context.Background(), 1)
	assertNoError(t, err)
	if person.Name != "Mike Vogel" {
		t.Fatalf("got name %q", person.Name)
	}
	if person.Gender != "Male" {
		t.Fatalf("got gender %q", person.Gender)
	}
}

func TestGetPersonCastCredits(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/people/1/castcredits" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []tvmaze.CastCredit{
			{
				Self:  false,
				Voice: false,
				Links: tvmaze.CastCreditLinks{
					Show:      tvmaze.Link{Href: "/shows/1", Name: "Under the Dome"},
					Character: tvmaze.Link{Href: "/characters/1", Name: "Dale Barbara"},
				},
			},
		})
	})

	credits, err := c.GetPersonCastCredits(context.Background(), 1)
	assertNoError(t, err)
	if len(credits) != 1 {
		t.Fatalf("got %d credits", len(credits))
	}
	if credits[0].Links.Show.Name != "Under the Dome" {
		t.Fatalf("got show %q", credits[0].Links.Show.Name)
	}
}

func TestGetPersonCrewCredits(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/people/15/crewcredits" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []tvmaze.CrewCredit{
			{
				Type: "Creator",
				Links: tvmaze.CrewCreditLinks{
					Show: tvmaze.Link{Href: "/shows/1", Name: "Under the Dome"},
				},
			},
		})
	})

	credits, err := c.GetPersonCrewCredits(context.Background(), 15)
	assertNoError(t, err)
	if len(credits) != 1 {
		t.Fatalf("got %d credits", len(credits))
	}
	if credits[0].Type != "Creator" {
		t.Fatalf("got type %q", credits[0].Type)
	}
}

func TestGetSchedule(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/schedule" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("country") != "US" {
			t.Fatalf("unexpected country: %s", r.URL.Query().Get("country"))
		}
		if r.URL.Query().Get("date") != "2024-01-01" {
			t.Fatalf("unexpected date: %s", r.URL.Query().Get("date"))
		}
		respond(t, w, []tvmaze.ScheduleItem{
			{
				Episode: tvmaze.Episode{ID: 1, Name: "Pilot"},
				Show:    tvmaze.Show{ID: 1, Name: "Under the Dome"},
			},
		})
	})

	items, err := c.GetSchedule(context.Background(), "US", "2024-01-01")
	assertNoError(t, err)
	if len(items) != 1 {
		t.Fatalf("got %d items", len(items))
	}
	if items[0].Show.Name != "Under the Dome" {
		t.Fatalf("got show %q", items[0].Show.Name)
	}
}

func TestGetScheduleNoParams(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/schedule" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("country") != "" || r.URL.Query().Get("date") != "" {
			t.Fatal("unexpected query params")
		}
		respond(t, w, []tvmaze.ScheduleItem{})
	})

	items, err := c.GetSchedule(context.Background(), "", "")
	assertNoError(t, err)
	if len(items) != 0 {
		t.Fatalf("got %d items", len(items))
	}
}

func TestGetWebSchedule(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/schedule/web" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		respond(t, w, []tvmaze.ScheduleItem{
			{
				Episode: tvmaze.Episode{ID: 100, Name: "Episode 1"},
				Show:    tvmaze.Show{ID: 50, Name: "Test Show"},
			},
		})
	})

	items, err := c.GetWebSchedule(context.Background(), "", "")
	assertNoError(t, err)
	if len(items) != 1 {
		t.Fatalf("got %d items", len(items))
	}
}

func TestGetShowUpdates(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/updates/shows" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("since") != "day" {
			t.Fatalf("unexpected since: %s", r.URL.Query().Get("since"))
		}
		respond(t, w, map[string]int64{"1": 1704067200, "2": 1704153600})
	})

	updates, err := c.GetShowUpdates(context.Background(), tvmaze.UpdateDay)
	assertNoError(t, err)
	if len(updates) != 2 {
		t.Fatalf("got %d updates", len(updates))
	}
	if updates["1"] != 1704067200 {
		t.Fatalf("got timestamp %d", updates["1"])
	}
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := c.GetShow(context.Background(), 99999)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *metadata.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("got status %d", apiErr.StatusCode)
	}
}
