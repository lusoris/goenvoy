package tvdb_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/video/tvdb"
)

func Example() {
	// Create a new TheTVDB client
	client := tvdb.New("your-api-key")

	ctx := context.Background()

	// Login to get token
	if err := client.Login(ctx); err != nil {
		log.Fatal(err)
	}

	// Search for a series
	results, err := client.Search(ctx, "Breaking Bad", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d series\n", len(results))
}
