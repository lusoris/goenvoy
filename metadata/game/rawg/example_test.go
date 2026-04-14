package rawg_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/game/rawg"
)

func Example() {
	c := rawg.New("your-api-key")

	result, err := c.SearchGames(context.Background(), "zelda", 1, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Games: %d\n", result.Count)
}
