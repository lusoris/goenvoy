package musicbrainz_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/music/musicbrainz"
)

func Example() {
	// Create a new MusicBrainz client
	client := musicbrainz.New(metadata.WithUserAgent("your-app-name/1.0"))

	ctx := context.Background()

	// Search for artists
	results, err := client.SearchArtists(ctx, "The Beatles", 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d artists\n", len(results.Entities))
}
