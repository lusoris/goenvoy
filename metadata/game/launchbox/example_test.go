package launchbox_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/game/launchbox"
)

func Example() {
	c := launchbox.New()

	if err := c.Download(context.Background()); err != nil {
		log.Fatal(err)
	}

	games := c.SearchGames("Mario", "")
	for i := range games {
		fmt.Printf("%s (%s)\n", games[i].Name, games[i].Platform)
	}
}
