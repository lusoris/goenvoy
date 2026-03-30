package radarr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/arr/radarr"
)

func Example() {
	// Create a new Radarr client
	client := radarr.New("http://localhost:7878", "your-api-key")

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
