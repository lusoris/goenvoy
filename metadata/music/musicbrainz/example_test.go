package musicbrainz_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/music/musicbrainz"
	"github.com/lusoris/goenvoy/metadata"
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
