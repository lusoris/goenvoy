package audiodb_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/music/audiodb"
)

func Example() {
	c := audiodb.New("2")

	artists, err := c.SearchArtist(context.Background(), "coldplay")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Artists: %d\n", len(artists))
}
