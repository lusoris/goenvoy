package tpdb_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/adult/tpdb"
)

func Example() {
	// Create a new ThePornDB client
	client := tpdb.New("your-bearer-token")

	ctx := context.Background()

	// Search for performers
	results, _, err := client.SearchPerformers(ctx, &tpdb.PerformerSearchParams{Query: "Jane Doe"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d performers\n", len(results))
}
