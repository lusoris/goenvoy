package plex_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/mediaserver/plex"
)

func Example() {
	// Create a new Plex client
	client := plex.New("http://192.168.1.100:32400", plex.WithToken("your-plex-token"))

	ctx := context.Background()

	// Get server identity (no authentication required)
	identity, err := client.GetIdentity(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server: %s\n", identity.FriendlyName)

	// Get libraries
	libraries, err := client.GetLibraries(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Libraries: %d\n", len(libraries))
}
