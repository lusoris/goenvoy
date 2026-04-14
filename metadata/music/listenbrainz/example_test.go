package listenbrainz_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/music/listenbrainz"
)

func Example() {
	c := listenbrainz.New()

	resp, err := c.GetUserListens(context.Background(), "username", 5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listens: %d\n", resp.Payload.Count)
}

func Example_withToken() {
	c := listenbrainz.NewWithToken("my-token")

	err := c.SubmitListens(context.Background(), "single", []listenbrainz.Listen{
		{
			ListenedAt: 1700000000,
			TrackMetadata: listenbrainz.TrackMetadata{
				ArtistName:  "Radiohead",
				TrackName:   "Creep",
				ReleaseName: "Pablo Honey",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
