package lastfm_test

import (
	"fmt"

	"github.com/lusoris/goenvoy/metadata/music/lastfm"
)

func Example() {
	c := lastfm.New("your-api-key")
	fmt.Println(c)
}
