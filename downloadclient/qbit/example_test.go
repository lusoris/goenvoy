package qbit_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/downloadclient/qbit"
)

func Example() {
	// Create a new qBittorrent client
	client := qbit.New("http://localhost:8080")

	ctx := context.Background()

	// Login with username and password
	if err := client.Login(ctx, "admin", "adminpass"); err != nil {
		log.Fatal(err)
	}

	// Get application version
	version, err := client.Version(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("qBittorrent version: %s\n", version)

	// Get all torrents
	torrents, err := client.ListTorrents(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total torrents: %d\n", len(torrents))
}
