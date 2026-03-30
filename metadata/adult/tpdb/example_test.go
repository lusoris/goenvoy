package tpdb_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/metadata/adult/tpdb"
)

func Example() {
	// Create a new ThePornDB client
	client := tpdb.New("your-bearer-token")

	ctx := context.Background()

	// Search for performers
	results, err := client.SearchPerformers(ctx, "Jane Doe")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d performers\n", len(results.Data))
}
