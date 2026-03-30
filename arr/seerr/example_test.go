package seerr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/arr/seerr"
)

func Example() {
	// Create a new Seerr client
	client := seerr.New("http://localhost:5055", "your-api-key")

	ctx := context.Background()

	// Get status
	status, err := client.GetStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Seerr version %s\n", status.Version)

	// Get requests
	requests, err := client.GetRequests(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total requests: %d\n", requests.PageInfo.Results)
}
