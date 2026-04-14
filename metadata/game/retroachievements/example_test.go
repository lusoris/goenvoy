package retroachievements_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/game/retroachievements"
)

func Example() {
	c := retroachievements.New("your-api-key")

	game, err := c.GetGame(context.Background(), 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Game: %s (%s)\n", game.Title, game.ConsoleName)
}
