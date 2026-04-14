package kavita_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/mediaserver/kavita"
)

func Example() {
	c := kavita.New("http://localhost:5000", "your-api-key")
	libs, err := c.GetLibraries(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(libs))
}
