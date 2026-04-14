package steamgriddb_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/game/steamgriddb"
)

func Example() {
	c := steamgriddb.New("your-api-key")

	game, err := c.GetGameByID(context.Background(), 12345)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Game: %s\n", game.Name)
}
