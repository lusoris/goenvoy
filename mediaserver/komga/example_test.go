package komga_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/mediaserver/komga"
)

func Example() {
	c := komga.New("http://localhost:25600", "admin@example.com", "password")

	ctx := context.Background()

	libs, err := c.GetLibraries(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Libraries: %d\n", len(libs))
}
