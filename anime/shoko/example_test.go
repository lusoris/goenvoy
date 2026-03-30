package shoko_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/anime/shoko"
)

func Example() {
	// Create a new Shoko client
	client := shoko.New("http://localhost:8111")

	ctx := context.Background()

	// Login with username and password
	if err := client.Login(ctx, "username", "password"); err != nil {
		log.Fatal(err)
	}

	// Get dashboard stats
	stats, err := client.GetDashboardStats(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total series: %d\n", stats.SeriesCount)

	// List all series
	series, err := client.ListSeries(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Series count: %d\n", series.Total)
}
