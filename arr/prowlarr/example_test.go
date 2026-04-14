package prowlarr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/prowlarr/v2"
)

func Example() {
	// Create a new Prowlarr client
	client, err := prowlarr.New("http://localhost:9696", "your-api-key")
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

	// Get all indexers
	indexers, err := client.GetIndexers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total indexers: %d\n", len(indexers))
}
