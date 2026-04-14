package jellyfin_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/mediaserver/jellyfin/v2"
)

func Example() {
	// Create a new Jellyfin client
	client := jellyfin.New("http://192.168.1.100:8096")

	ctx := context.Background()

	// Authenticate with username and password
	if err := client.AuthenticateByName(ctx, "username", "password"); err != nil {
		log.Fatal(err)
	}

	// Get current user
	user, err := client.GetCurrentUser(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User: %s\n", user.Name)

	// Get all items
	items, err := client.GetItems(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total items: %d\n", items.TotalRecordCount)
}
