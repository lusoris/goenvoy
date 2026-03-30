package trakt_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/tracking/trakt"
)

func Example() {
	// Create a new Trakt client
	client := trakt.New("your-client-id")

	ctx := context.Background()

	// Search for shows
	results, err := client.Search(ctx, "game of thrones", "show")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d results\n", len(results))
}
