package anime_test

import (
	"encoding/json"
	"testing"

	"github.com/golusoris/goenvoy/anime"
)

func TestSeriesJSON(t *testing.T) {
	t.Parallel()

	s := anime.Series{
		ID:   1,
		Name: "My Anime",
		Size: 12,
		AniDB: &anime.AniDBInfo{
			AniDBID:      100,
			Type:         "TV",
			Title:        "My Anime",
			EpisodeCount: 12,
			Rating:       8.5,
		},
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got anime.Series
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ID != s.ID || got.Name != s.Name || got.Size != s.Size {
		t.Errorf("got %+v, want %+v", got, s)
	}
	if got.AniDB == nil || got.AniDB.AniDBID != 100 {
		t.Error("AniDB info lost in round-trip")
	}
}

func TestEpisodeJSON(t *testing.T) {
	t.Parallel()

	e := anime.Episode{
		ID:       1,
		SeriesID: 10,
		Name:     "Episode 1",
		Number:   1,
		Type:     "Normal",
		AirDate:  "2024-01-01",
	}

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got anime.Episode
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got != e {
		t.Errorf("got %+v, want %+v", got, e)
	}
}
