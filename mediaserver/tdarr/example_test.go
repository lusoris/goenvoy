package tdarr_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/mediaserver/tdarr/v2"
)

func Example() {
	client := tdarr.New("http://localhost:8265", tdarr.WithAPIKey("your-key"))

	ctx := context.Background()

	status, err := client.GetStatus(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Tdarr %s on %s\n", status.Version, status.Os)
}
