package opensubtitles_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/video/opensubtitles"
)

func Example() {
	c := opensubtitles.New("your-api-key")

	ctx := context.Background()

	results, err := c.Search(ctx, &opensubtitles.SearchParams{
		Query: "Inception",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Subtitles found: %d\n", results.TotalCount)
}
