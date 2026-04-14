package simkl_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/tracking/simkl"
)

func Example() {
	// Create a new Simkl client
	client := simkl.New("your-client-id")

	ctx := context.Background()

	// Search for shows
	results, err := client.SearchText(ctx, "tv", "Breaking Bad", 1, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d results\n", len(results))
}
