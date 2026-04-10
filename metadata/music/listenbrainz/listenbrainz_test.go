package listenbrainz_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lusoris/goenvoy/metadata/music/listenbrainz"
	"github.com/lusoris/goenvoy/metadata"
)

func setup(t *testing.T, handler http.HandlerFunc) *listenbrainz.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return listenbrainz.NewWithToken("test-token", metadata.WithBaseURL(srv.URL))
}

func TestGetUserListens(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if got := r.URL.Query().Get("count"); got != "5" {
			t.Errorf("count = %q, want 5", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"count":            2,
				"user_name":        "testuser",
				"latest_listen_ts": 1700000000,
				"oldest_listen_ts": 1699000000,
				"listens": []map[string]any{
					{
						"listened_at": 1700000000,
						"inserted_at": 1700000001,
						"track_metadata": map[string]any{
							"artist_name":  "Radiohead",
							"track_name":   "Creep",
							"release_name": "Pablo Honey",
						},
					},
					{
						"listened_at": 1699500000,
						"inserted_at": 1699500001,
						"track_metadata": map[string]any{
							"artist_name":  "Nirvana",
							"track_name":   "Smells Like Teen Spirit",
							"release_name": "Nevermind",
						},
					},
				},
			},
		})
	})

	resp, err := c.GetUserListens(context.Background(), "testuser", 5)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Payload.Count != 2 {
		t.Errorf("Count = %d, want 2", resp.Payload.Count)
	}
	if resp.Payload.UserName != "testuser" {
		t.Errorf("UserName = %q, want testuser", resp.Payload.UserName)
	}
	if len(resp.Payload.Listens) != 2 {
		t.Fatalf("len(Listens) = %d, want 2", len(resp.Payload.Listens))
	}
	if resp.Payload.Listens[0].TrackMetadata.ArtistName != "Radiohead" {
		t.Errorf("ArtistName = %q, want Radiohead", resp.Payload.Listens[0].TrackMetadata.ArtistName)
	}
}

func TestGetListenCount(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"count": 42000,
			},
		})
	})

	count, err := c.GetListenCount(context.Background(), "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if count != 42000 {
		t.Errorf("count = %d, want 42000", count)
	}
}

func TestGetPlayingNow(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"count": 1,
				"listens": []map[string]any{
					{
						"track_metadata": map[string]any{
							"artist_name":  "Pink Floyd",
							"track_name":   "Comfortably Numb",
							"release_name": "The Wall",
						},
					},
				},
			},
		})
	})

	resp, err := c.GetPlayingNow(context.Background(), "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Payload.Count != 1 {
		t.Errorf("Count = %d, want 1", resp.Payload.Count)
	}
	if resp.Payload.Listens[0].TrackMetadata.TrackName != "Comfortably Numb" {
		t.Errorf("TrackName = %q, want Comfortably Numb", resp.Payload.Listens[0].TrackMetadata.TrackName)
	}
}

func TestGetUserTopArtists(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("range"); got != "week" {
			t.Errorf("range = %q, want week", got)
		}
		if got := r.URL.Query().Get("count"); got != "3" {
			t.Errorf("count = %q, want 3", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"artists": []map[string]any{
					{"artist_name": "Radiohead", "listen_count": 150, "artist_mbid": "abc-123"},
					{"artist_name": "Tool", "listen_count": 120},
					{"artist_name": "Deftones", "listen_count": 90},
				},
				"count":              3,
				"total_artist_count": 50,
				"range":              "week",
				"user_name":          "testuser",
			},
		})
	})

	resp, err := c.GetUserTopArtists(context.Background(), "testuser", "week", 3)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Payload.Count != 3 {
		t.Errorf("Count = %d, want 3", resp.Payload.Count)
	}
	if resp.Payload.TotalArtistCount != 50 {
		t.Errorf("TotalArtistCount = %d, want 50", resp.Payload.TotalArtistCount)
	}
	if resp.Payload.Artists[0].ArtistName != "Radiohead" {
		t.Errorf("ArtistName = %q, want Radiohead", resp.Payload.Artists[0].ArtistName)
	}
	if resp.Payload.Artists[0].ListenCount != 150 {
		t.Errorf("ListenCount = %d, want 150", resp.Payload.Artists[0].ListenCount)
	}
}

func TestGetUserTopReleases(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"releases": []map[string]any{
					{"artist_name": "Radiohead", "release_name": "OK Computer", "listen_count": 80},
				},
				"count":               1,
				"total_release_count": 25,
				"range":               "month",
				"user_name":           "testuser",
			},
		})
	})

	resp, err := c.GetUserTopReleases(context.Background(), "testuser", "month", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Payload.Releases) != 1 {
		t.Fatalf("len(Releases) = %d, want 1", len(resp.Payload.Releases))
	}
	if resp.Payload.Releases[0].ReleaseName != "OK Computer" {
		t.Errorf("ReleaseName = %q, want OK Computer", resp.Payload.Releases[0].ReleaseName)
	}
}

func TestGetUserTopRecordings(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"recordings": []map[string]any{
					{
						"artist_name":    "Tool",
						"track_name":     "Lateralus",
						"release_name":   "Lateralus",
						"listen_count":   45,
						"recording_mbid": "xyz-789",
					},
				},
				"count":                 1,
				"total_recording_count": 100,
				"range":                 "all_time",
				"user_name":             "testuser",
			},
		})
	})

	resp, err := c.GetUserTopRecordings(context.Background(), "testuser", "all_time", 10)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Payload.Recordings[0].TrackName != "Lateralus" {
		t.Errorf("TrackName = %q, want Lateralus", resp.Payload.Recordings[0].TrackName)
	}
	if resp.Payload.Recordings[0].ListenCount != 45 {
		t.Errorf("ListenCount = %d, want 45", resp.Payload.Recordings[0].ListenCount)
	}
}

func TestGetListeningActivity(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("range"); got != "week" {
			t.Errorf("range = %q, want week", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"listening_activity": []map[string]any{
					{"listen_count": 30, "from_ts": 1699000000, "to_ts": 1699100000, "time_range": "Monday"},
					{"listen_count": 45, "from_ts": 1699100000, "to_ts": 1699200000, "time_range": "Tuesday"},
				},
				"user_name": "testuser",
				"range":     "week",
			},
		})
	})

	resp, err := c.GetListeningActivity(context.Background(), "testuser", "week")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Payload.ListeningActivity) != 2 {
		t.Fatalf("len = %d, want 2", len(resp.Payload.ListeningActivity))
	}
	if resp.Payload.ListeningActivity[0].ListenCount != 30 {
		t.Errorf("ListenCount = %d, want 30", resp.Payload.ListeningActivity[0].ListenCount)
	}
}

func TestGetDailyActivity(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"daily_activity": map[string]any{
					"Monday": []map[string]any{
						{"hour": 10, "listen_count": 5},
						{"hour": 14, "listen_count": 12},
					},
				},
				"user_name": "testuser",
				"range":     "week",
			},
		})
	})

	resp, err := c.GetDailyActivity(context.Background(), "testuser", "week")
	if err != nil {
		t.Fatal(err)
	}
	monday, ok := resp.Payload.DailyActivity["Monday"]
	if !ok {
		t.Fatal("expected Monday key")
	}
	if len(monday) != 2 {
		t.Fatalf("len(Monday) = %d, want 2", len(monday))
	}
	if monday[0].Hour != 10 {
		t.Errorf("Hour = %d, want 10", monday[0].Hour)
	}
	if monday[1].ListenCount != 12 {
		t.Errorf("ListenCount = %d, want 12", monday[1].ListenCount)
	}
}

func TestGetSitewideArtists(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"payload": map[string]any{
				"artists": []map[string]any{
					{"artist_name": "Taylor Swift", "listen_count": 99999},
				},
				"count":              1,
				"total_artist_count": 1000,
				"range":              "week",
			},
		})
	})

	resp, err := c.GetSitewideArtists(context.Background(), "week", 10)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Payload.Artists[0].ArtistName != "Taylor Swift" {
		t.Errorf("ArtistName = %q, want Taylor Swift", resp.Payload.Artists[0].ArtistName)
	}
}

func TestGetSimilarUsers(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"payload": []map[string]any{
				{"user_name": "similar_user1", "similarity": 0.85},
				{"user_name": "similar_user2", "similarity": 0.72},
			},
		})
	})

	users, err := c.GetSimilarUsers(context.Background(), "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("len = %d, want 2", len(users))
	}
	if users[0].UserName != "similar_user1" {
		t.Errorf("UserName = %q, want similar_user1", users[0].UserName)
	}
	if users[0].Similarity != 0.85 {
		t.Errorf("Similarity = %f, want 0.85", users[0].Similarity)
	}
}

func TestGetLatestImport(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("user_name"); got != "testuser" {
			t.Errorf("user_name = %q, want testuser", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"latest_import": 1700000000,
		})
	})

	ts, err := c.GetLatestImport(context.Background(), "testuser")
	if err != nil {
		t.Fatal(err)
	}
	if ts != 1700000000 {
		t.Errorf("timestamp = %d, want 1700000000", ts)
	}
}

func TestSubmitListens(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Token test-token" {
			t.Errorf("Authorization = %q, want Token test-token", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if payload["listen_type"] != "single" {
			t.Errorf("listen_type = %v, want single", payload["listen_type"])
		}
		items, ok := payload["payload"].([]any)
		if !ok || len(items) != 1 {
			t.Fatalf("payload length = %d, want 1", len(items))
		}
		w.WriteHeader(http.StatusOK)
	})

	err := c.SubmitListens(context.Background(), "single", []listenbrainz.Listen{
		{
			ListenedAt: 1700000000,
			TrackMetadata: listenbrainz.TrackMetadata{
				ArtistName:  "Radiohead",
				TrackName:   "Karma Police",
				ReleaseName: "OK Computer",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAPIError(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	})

	_, err := c.GetListenCount(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *listenbrainz.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
}

func TestSubmitListensNoBody(t *testing.T) {
	c := setup(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err := c.SubmitListens(context.Background(), "import", []listenbrainz.Listen{})
	if err != nil {
		t.Fatal(err)
	}
}
