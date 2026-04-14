package igdb_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/game/igdb"
)

func Example() {
	c := igdb.New("your-client-id", "your-access-token")

	games, err := c.SearchGames(context.Background(), "zelda", 5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Games: %d\n", len(games))
}
