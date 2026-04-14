package mal_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/anime/mal"
)

func Example() {
	// Create a new MyAnimeList client
	client := mal.New("your-client-id")

	ctx := context.Background()

	// Search anime
	results, _, err := client.SearchAnime(ctx, "Steins Gate", nil, 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d anime\n", len(results))
}
