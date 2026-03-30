package mal_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/anime/mal"
)

func Example() {
	// Create a new MyAnimeList client
	client := mal.New("your-client-id")

	ctx := context.Background()

	// Search anime
	results, err := client.SearchAnime(ctx, "Steins Gate", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d anime\n", len(results.Data))
}
