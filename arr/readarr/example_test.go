package readarr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/readarr"
)

func Example() {
	// Create a new Readarr client
	client, err := readarr.New("http://localhost:8787", "your-api-key")
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

	// Get all authors
	authors, err := client.GetAllAuthors(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total authors: %d\n", len(authors))
}
