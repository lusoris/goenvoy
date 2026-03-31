package discogs_test

import (
	"fmt"

	"github.com/lusoris/goenvoy/metadata/music/discogs"
)

func Example() {
	c := discogs.New("your-token")
	fmt.Println(c)
}
