package tautulli_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/mediaserver/tautulli"
)

func Example() {
	client := tautulli.New("http://localhost:8181", "your-api-key")

	ctx := context.Background()

	// Get current activity
	activity, err := client.GetActivity(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Active streams: %s\n", activity.StreamCount)

	// Get libraries
	libraries, err := client.GetLibraries(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Libraries: %d\n", len(libraries))
}
