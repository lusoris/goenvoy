package audiodb_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/music/audiodb/v2"
)

func Example() {
	c := audiodb.New("2")

	artists, err := c.SearchArtist(context.Background(), "coldplay")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Artists: %d\n", len(artists))
}
