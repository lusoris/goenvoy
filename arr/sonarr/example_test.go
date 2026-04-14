package sonarr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/sonarr"
)

func Example() {
	// Create a new Sonarr client
	client, err := sonarr.New("http://localhost:8989", "your-api-key")
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

	// Get all series
	series, err := client.GetAllSeries(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total series: %d\n", len(series))
}
