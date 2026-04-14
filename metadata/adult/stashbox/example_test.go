package stashbox_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/adult/stashbox"
)

func Example() {
	// Create a new StashBox client
	client := stashbox.New("https://stashdb.org/graphql", "your-api-key")

	ctx := context.Background()

	// Search for performers
	results, err := client.SearchPerformers(ctx, "Jane Doe", 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d performers\n", len(results))
}
