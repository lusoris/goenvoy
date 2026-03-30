package rtorrent_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/downloadclient/rtorrent"
)

func Example() {
	// Create a new rTorrent client
	client := rtorrent.New("http://localhost:8000/RPC2")

	ctx := context.Background()

	// Get rTorrent version
	version, err := client.GetVersion(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("rTorrent version: %s\n", version)

	// Get download list
	torrents, err := client.GetTorrents(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total torrents: %d\n", len(torrents))
}
