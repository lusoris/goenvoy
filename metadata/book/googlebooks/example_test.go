package googlebooks_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/book/googlebooks/v2"
)

func Example() {
	c := googlebooks.New("your-api-key")

	results, err := c.Search(context.Background(), "flowers for algernon")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found: %d\n", results.TotalItems)
}
