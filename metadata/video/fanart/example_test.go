package fanart_test

import (
	"context"
	"fmt"
	"log"

	"github.com/golusoris/goenvoy/metadata/video/fanart"
)

func Example() {
	// Create a new Fanart.tv client
	client := fanart.New("your-api-key")

	ctx := context.Background()

	// Get movie images
	images, err := client.GetMovieImages(ctx, "123456")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Movie posters: %d\n", len(images.MoviePoster))
}
