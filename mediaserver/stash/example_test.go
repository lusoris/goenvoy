package stash_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/mediaserver/stash"
)

func Example() {
	// Create a new Stash client
	client := stash.New("http://localhost:9999/graphql", "your-api-key")

	ctx := context.Background()

	// Get system status
	status, err := client.GetSystemStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Stash status: %s\n", status.Status)

	// Find scenes
	scenes, count, err := client.FindScenes(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total scenes: %d\n", count)
	_ = scenes
}
