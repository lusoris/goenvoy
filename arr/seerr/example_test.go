package seerr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/seerr"
)

func Example() {
	// Create a new Seerr client
	client, err := seerr.New("http://localhost:5055", "your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Get status
	status, err := client.GetStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Seerr version %s\n", status.Version)

	// Get requests
	requests, _, err := client.GetRequests(ctx, 25, 0, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total requests: %d\n", len(requests))
}
