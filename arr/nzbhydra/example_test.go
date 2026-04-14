package nzbhydra_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/nzbhydra"
)

func Example() {
	client := nzbhydra.New("http://localhost:5076", "your-api-key")

	ctx := context.Background()

	// Search all indexers
	results, err := client.Search(ctx, "ubuntu", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Results: %d\n", len(results))
}
