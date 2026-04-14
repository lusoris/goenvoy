package autobrr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/autobrr"
)

func Example() {
	client := autobrr.New("http://localhost:7474", "your-api-key")

	ctx := context.Background()

	// Get all filters
	filters, err := client.GetFilters(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Filters: %d\n", len(filters))

	// Get IRC networks
	networks, err := client.GetIRCNetworks(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("IRC Networks: %d\n", len(networks))
}
