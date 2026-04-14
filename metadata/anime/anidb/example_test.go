package anidb_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/anime/anidb"
)

func Example() {
	// Create a new AniDB client
	client := anidb.New("your-client-name", 1)

	ctx := context.Background()

	// Get anime info
	anime, err := client.GetAnime(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Anime ID: %d\n", anime.ID)
}
