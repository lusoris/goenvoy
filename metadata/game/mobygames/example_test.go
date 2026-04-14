package mobygames_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/game/mobygames"
)

func Example() {
	c := mobygames.New("your-api-key")

	games, err := c.SearchGames(context.Background(), "zelda", 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Games: %d\n", len(games))
}
