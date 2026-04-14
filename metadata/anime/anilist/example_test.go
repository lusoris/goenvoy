package anilist_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/anime/anilist"
)

func Example() {
	// Create a new AniList client
	client := anilist.New()

	ctx := context.Background()

	// Search for anime
	results, err := client.SearchMedia(ctx, "Cowboy Bebop", anilist.MediaTypeAnime, 1, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d anime\n", len(results.Media))

	if len(results.Media) > 0 {
		fmt.Printf("Anime: %s\n", results.Media[0].Title.Romaji)
	}
}
