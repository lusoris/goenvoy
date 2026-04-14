package navidrome_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/mediaserver/navidrome"
)

func Example() {
	c := navidrome.New("http://localhost:4533", "admin", "password")

	ctx := context.Background()

	if err := c.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is reachable")

	artists, err := c.GetArtists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Artists: %d\n", len(artists.Index))
}
