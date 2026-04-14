package letterboxd_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/video/letterboxd"
)

func Example() {
	// Create a new Letterboxd client with an OAuth2 Bearer token.
	client := letterboxd.New("your-access-token")

	ctx := context.Background()

	// Search for a film.
	results, err := client.Search(ctx, "Inception", "", 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d results\n", len(results.Items))

	// Get film details by Letterboxd ID or external ID.
	film, err := client.GetFilm(ctx, "tmdb:27205")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Film: %s (%d)\n", film.Name, film.ReleaseYear)
}
