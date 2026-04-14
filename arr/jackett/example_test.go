package jackett_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/jackett"
)

func Example() {
	client := jackett.New("http://localhost:9117", "your-api-key")

	ctx := context.Background()

	// Search all indexers
	results, err := client.Search(ctx, "ubuntu", []int{4000})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Results: %d\n", len(results))
}
