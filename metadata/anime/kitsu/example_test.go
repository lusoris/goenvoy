package kitsu_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/anime/kitsu"
)

func Example() {
	// Create a new Kitsu client
	client := kitsu.New()

	ctx := context.Background()

	// Get anime by ID
	anime, err := client.GetAnime(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Anime: %s\n", anime.CanonicalTitle)
}
