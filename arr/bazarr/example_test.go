package bazarr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/arr/bazarr"
)

func Example() {
	// Create a new Bazarr client
	client := bazarr.New("http://localhost:6767", "your-api-key")

	ctx := context.Background()

	// Get system status
	status, err := client.GetSystemStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bazarr version %s\n", status.Version)

	// Get series
	series, err := client.GetSeries(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total series: %d\n", len(series))
}
