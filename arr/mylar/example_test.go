package mylar_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/arr/mylar/v2"
)

func Example() {
	client := mylar.New("http://localhost:8090", "your-api-key")

	ctx := context.Background()

	// Get all comic series
	comics, err := client.GetIndex(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Comics: %d\n", len(comics))

	// Get version info
	ver, err := client.GetVersion(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Version: %s\n", ver.Version)
}
