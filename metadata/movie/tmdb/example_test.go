package tmdb_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/movie/tmdb"
)

func Example() {
	// Create a new TMDb client
	client := tmdb.New("your-api-key")

	ctx := context.Background()

	// Search for a movie
	results, err := client.SearchMovie(ctx, "Inception", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d movies\n", len(results.Results))

	// Get movie details
	if len(results.Results) > 0 {
		movie, err := client.GetMovie(ctx, results.Results[0].ID, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Movie: %s\n", movie.Title)
	}
}
