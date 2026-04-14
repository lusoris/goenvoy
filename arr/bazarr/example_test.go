package bazarr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/bazarr"
)

func Example() {
	// Create a new Bazarr client
	client, err := bazarr.New("http://localhost:6767", "your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Get system status
	status, err := client.GetSystemStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bazarr version %s\n", status.BazarrVersion)

	// Get series
	series, err := client.GetSeries(ctx, 0, 25)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total series: %d\n", series.Total)
}
