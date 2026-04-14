package screenscraper_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/game/screenscraper"
)

func Example() {
	c := screenscraper.New("devid", "devpassword", "myapp",
		screenscraper.WithUser("user", "pass"))

	result, err := c.GetGameInfo(context.Background(), &screenscraper.GameInfoOptions{
		CRC: "ABCD1234",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Game: %s\n", result.Response.Game.ID)
}
