package metadata_test

import (
	"encoding/json"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
)

func TestRatingJSON(t *testing.T) {
	t.Parallel()

	r := metadata.Rating{Source: "imdb", Value: 8.7, Votes: 100000}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got metadata.Rating
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got != r {
		t.Errorf("got %+v, want %+v", got, r)
	}
}

func TestExternalIDJSON(t *testing.T) {
	t.Parallel()

	eid := metadata.ExternalID{Source: "tmdb", ID: "550"}

	data, err := json.Marshal(eid)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got metadata.ExternalID
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got != eid {
		t.Errorf("got %+v, want %+v", got, eid)
	}
}

func TestImageTypeConstants(t *testing.T) {
	t.Parallel()

	types := []metadata.ImageType{
		metadata.ImageTypePoster,
		metadata.ImageTypeBackdrop,
		metadata.ImageTypeLogo,
		metadata.ImageTypeStill,
	}

	seen := make(map[metadata.ImageType]bool)
	for _, it := range types {
		if it == "" {
			t.Error("empty image type constant")
		}
		if seen[it] {
			t.Errorf("duplicate image type: %s", it)
		}
		seen[it] = true
	}
}

func TestMediaTypeConstants(t *testing.T) {
	t.Parallel()

	types := []metadata.MediaType{
		metadata.MediaTypeMovie,
		metadata.MediaTypeSeries,
		metadata.MediaTypePerson,
	}

	seen := make(map[metadata.MediaType]bool)
	for _, mt := range types {
		if mt == "" {
			t.Error("empty media type constant")
		}
		if seen[mt] {
			t.Errorf("duplicate media type: %s", mt)
		}
		seen[mt] = true
	}
}

func TestSearchResultJSON(t *testing.T) {
	t.Parallel()

	sr := metadata.SearchResult{
		Title:     "Inception",
		Year:      2010,
		Type:      metadata.MediaTypeMovie,
		PosterURL: "https://example.com/poster.jpg",
		Overview:  "A thief who steals corporate secrets...",
		ExternalIDs: []metadata.ExternalID{
			{Source: "imdb", ID: "tt1375666"},
			{Source: "tmdb", ID: "27205"},
		},
	}

	data, err := json.Marshal(sr)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got metadata.SearchResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Title != sr.Title || got.Year != sr.Year || got.Type != sr.Type {
		t.Errorf("got %+v, want %+v", got, sr)
	}
	if len(got.ExternalIDs) != 2 {
		t.Errorf("got %d external IDs, want 2", len(got.ExternalIDs))
	}
}
