package sabnzbd_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/downloadclient/sabnzbd"
)

func Example() {
	// Create a new SABnzbd client
	client := sabnzbd.New("http://localhost:8080", "your-api-key")

	ctx := context.Background()

	// Get version
	version, err := client.GetVersion(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SABnzbd version: %s\n", version)

	// Get queue
	queue, err := client.GetQueue(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Queue slots: %d\n", len(queue.Slots))
}
