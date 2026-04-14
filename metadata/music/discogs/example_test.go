package discogs_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/music/discogs"
)

func Example() {
	c := discogs.New("your-token")

	ctx := context.Background()

	artist, err := c.GetArtist(ctx, 108713)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Artist: %s\n", artist.Name)
}
