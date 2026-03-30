package simkl_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/tracking/simkl"
)

func Example() {
	// Create a new Simkl client
	client := simkl.New("your-client-id")

	ctx := context.Background()

	// Search for shows
	results, err := client.Search(ctx, "Breaking Bad", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d results\n", len(results))
}
