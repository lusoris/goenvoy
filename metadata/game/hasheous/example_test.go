package hasheous_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/game/hasheous"
)

func Example() {
	c := hasheous.New()

	result, err := c.LookupBySHA1(context.Background(), "da39a3ee5e6b4b0d3255bfef95601890afd80709")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Game: %s\n", result.Name)
}
