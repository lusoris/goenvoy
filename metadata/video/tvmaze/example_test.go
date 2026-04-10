package tvmaze_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/video/tvmaze"
)

func Example() {
	// Create a new TVmaze client
	client := tvmaze.New()

	ctx := context.Background()

	// Search for a show
	results, err := client.SearchShows(ctx, "Game of Thrones")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d shows\n", len(results))

	if len(results) > 0 {
		fmt.Printf("Show: %s\n", results[0].Show.Name)
	}
}
