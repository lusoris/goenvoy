package whisparr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/arr/whisparr"
)

func Example() {
	// Create a new Whisparr client
	client := whisparr.New("http://localhost:6969", "your-api-key")

	ctx := context.Background()

	// Get system status
	status, err := client.GetSystemStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s version %s\n", status.AppName, status.Version)

	// Get all sites
	sites, err := client.GetAllSites(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total sites: %d\n", len(sites))
}
