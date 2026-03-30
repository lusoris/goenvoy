package stash_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/stash"
)

func Example() {
	// Create a new Stash client
	client := stash.New("http://localhost:9999/graphql", "your-api-key")

	ctx := context.Background()

	// Get system status
	status, err := client.SystemStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Stash version: %s\n", status.Version)

	// Find scenes
	scenes, err := client.FindScenes(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total scenes: %d\n", scenes.Count)
}
