package transmission_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/downloadclient/transmission"
)

func Example() {
	// Create a new Transmission client
	client := transmission.New("http://localhost:9091/transmission/rpc",
		transmission.WithAuth("username", "password"))

	ctx := context.Background()

	// Get session stats
	stats, err := client.GetSessionStats(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Active torrents: %d\n", stats.ActiveTorrentCount)

	// List all torrents
	torrents, err := client.GetTorrents(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total torrents: %d\n", len(torrents))
}
