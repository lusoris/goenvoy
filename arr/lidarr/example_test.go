package lidarr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/lidarr"
)

func Example() {
	// Create a new Lidarr client
	client, err := lidarr.New("http://localhost:8686", "your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Get system status
	status, err := client.GetSystemStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s version %s\n", status.AppName, status.Version)

	// Get all artists
	artists, err := client.GetAllArtists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total artists: %d\n", len(artists))
}
