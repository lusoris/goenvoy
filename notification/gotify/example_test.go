package gotify_test

import (
	"context"
	"fmt"
	"log"

	"github.com/lusoris/goenvoy/notification/gotify"
)

func Example() {
	client := gotify.New("http://localhost:80", "your-app-token")

	ctx := context.Background()

	msg, err := client.CreateMessage(ctx, "Deploy", "v1.2.0 deployed successfully", 5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent message %d\n", msg.Id)
}
