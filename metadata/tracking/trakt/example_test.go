package trakt_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/tracking/trakt"
)

func Example() {
	// Create a new Trakt client
	client := trakt.New("your-client-id")

	ctx := context.Background()

	// Search for shows
	results, _, err := client.SearchText(ctx, "game of thrones", "show", 1, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d results\n", len(results))
}
