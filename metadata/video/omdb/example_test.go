package omdb_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/video/omdb"
)

func Example() {
	// Create a new OMDb client
	client := omdb.New("your-api-key")

	ctx := context.Background()

	// Search by title
	movie, err := client.GetByTitle(ctx, "The Matrix", 0, "", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Movie: %s (%s)\n", movie.Title, movie.Year)
}
