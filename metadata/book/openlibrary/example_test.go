package openlibrary_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/book/openlibrary"
)

func Example() {
	c := openlibrary.New()

	results, err := c.Search(context.Background(), "the lord of the rings")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found: %d\n", results.NumFound)
}
