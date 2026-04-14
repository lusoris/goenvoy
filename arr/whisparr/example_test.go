package whisparr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/whisparr"
)

func Example() {
	// Create a new Whisparr v2 client (Sonarr-based).
	client, err := whisparr.New("http://localhost:6969", "your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Get system status.
	status, err := client.GetSystemStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s version %s\n", status.AppName, status.Version)

	// Get all series.
	series, err := client.GetAllSeries(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total series: %d\n", len(series))
}

func Example_v3() {
	// Create a new Whisparr v3 client (Radarr-based).
	client, err := whisparr.NewV3("http://localhost:6969", "your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Get all movies.
	movies, err := client.GetAllMovies(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total movies: %d\n", len(movies))
}
