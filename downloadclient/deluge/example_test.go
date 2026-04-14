package deluge_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/downloadclient/deluge"
)

func Example() {
	// Create a new Deluge client
	client := deluge.New("http://localhost:8112/json")

	ctx := context.Background()

	// Login with password
	if err := client.Login(ctx, "deluge"); err != nil {
		log.Fatal(err)
	}

	// Get Deluge version
	version, err := client.GetVersion(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deluge version: %s\n", version)

	// Get torrent status
	torrents, err := client.GetTorrentsStatus(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total torrents: %d\n", len(torrents))
}
