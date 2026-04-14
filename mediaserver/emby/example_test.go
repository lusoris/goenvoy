package emby_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/mediaserver/emby/v2"
)

func Example() {
	// Create a new Emby client
	client := emby.New("http://192.168.1.100:8096")

	ctx := context.Background()

	// Authenticate with username and password
	if err := client.AuthenticateByName(ctx, "username", "password"); err != nil {
		log.Fatal(err)
	}

	// Get system info
	info, err := client.GetSystemInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server: %s %s\n", info.ServerName, info.Version)

	// Get active sessions
	sessions, err := client.GetSessions(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Active sessions: %d\n", len(sessions))
}
