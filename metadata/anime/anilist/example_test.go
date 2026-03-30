package anilist_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/anime/anilist"
)

func Example() {
	// Create a new AniList client
	client := anilist.New()

	ctx := context.Background()

	// Search for anime
	results, err := client.SearchAnime(ctx, "Cowboy Bebop", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d anime\n", len(results.Media))

	if len(results.Media) > 0 {
		fmt.Printf("Anime: %s\n", results.Media[0].Title.Romaji)
	}
}
