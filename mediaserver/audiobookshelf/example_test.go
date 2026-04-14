package audiobookshelf_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/mediaserver/audiobookshelf"
)

func Example() {
	client := audiobookshelf.New("http://localhost:13378", "your-token")

	ctx := context.Background()

	// Get all libraries
	libraries, err := client.GetLibraries(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Libraries: %d\n", len(libraries))

	// Get server info
	info, err := client.GetServerInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Version: %s\n", info.Version)
}
