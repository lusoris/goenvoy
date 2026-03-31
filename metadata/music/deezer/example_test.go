package deezer_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/music/deezer"
)

func Example() {
	c := deezer.New()

	artist, err := c.GetArtist(context.Background(), 27)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(artist.Name)
}
