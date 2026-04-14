package spotify_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/music/spotify"
)

func Example() {
	c := spotify.New("your-access-token")

	artist, err := c.GetArtist(context.Background(), "4gzpq5DPGxSnKTe4SA8HAU")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(artist.Name)
}
