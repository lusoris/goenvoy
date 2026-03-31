package kavita_test

import (
	"context"
	"fmt"

	"github.com/lusoris/goenvoy/mediaserver/kavita"
)

func Example() {
	c := kavita.New("http://localhost:5000", "your-api-key")
	libs, err := c.GetLibraries(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(libs))
}
