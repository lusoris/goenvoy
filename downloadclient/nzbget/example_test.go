package nzbget_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/downloadclient/nzbget"
)

func Example() {
	// Create a new NZBGet client
	client := nzbget.New("http://localhost:6789",
		nzbget.WithAuth("username", "password"))

	ctx := context.Background()

	// Get version
	version, err := client.Version(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("NZBGet version: %s\n", version)

	// List download queue
	queue, err := client.ListGroups(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Queue items: %d\n", len(queue))
}
