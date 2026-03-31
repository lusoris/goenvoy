package opensubtitles_test

import (
	"fmt"

	"github.com/lusoris/goenvoy/metadata/movie/opensubtitles"
)

func Example() {
	c := opensubtitles.New("your-api-key")
	fmt.Println(c)
}
