package steam_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/game/steam"
)

func Example() {
	c := steam.New()

	details, err := c.GetAppDetails(context.Background(), 730)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Game: %s\n", details.Name)
}

func Example_withAPIKey() {
	c := steam.NewWithAPIKey("your-api-key")

	count, err := c.GetCurrentPlayers(context.Background(), 730)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Players: %d\n", count)
}
