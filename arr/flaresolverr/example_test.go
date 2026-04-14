package flaresolverr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/flaresolverr"
)

func Example() {
	client := flaresolverr.New("http://localhost:8191")

	ctx := context.Background()

	// Solve a Cloudflare challenge
	resp, err := client.Get(ctx, "https://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %s\n", resp.Status)

	// List active sessions
	sessions, err := client.ListSessions(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sessions: %s\n", sessions.Message)
}
