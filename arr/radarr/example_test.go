package radarr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/radarr"
)

func Example() {
	// Create a new Radarr client
	client, err := radarr.New("http://localhost:7878", "your-api-key")
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

	// Get all movies
	movies, err := client.GetAllMovies(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total movies: %d\n", len(movies))
}
