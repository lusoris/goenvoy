package komga_test

import (
	"fmt"

	"github.com/lusoris/goenvoy/mediaserver/komga"
)

func Example() {
	c := komga.New("http://localhost:25600", "admin@example.com", "password")
	fmt.Println(c)
}
