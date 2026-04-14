package lastfm_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/music/lastfm"
)

func Example() {
	c := lastfm.New("your-api-key")

	ctx := context.Background()

	artist, err := c.GetArtistInfo(ctx, "Radiohead")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Artist: %s\n", artist.Name)
}
