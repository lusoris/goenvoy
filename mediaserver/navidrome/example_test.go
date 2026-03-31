package navidrome_test

import (
	"fmt"

	"github.com/lusoris/goenvoy/mediaserver/navidrome"
)

func Example() {
	c := navidrome.New("http://localhost:4533", "admin", "password")
	fmt.Println(c)
}
