// Package gotify provides a client for the Gotify REST API.
//
// Gotify is a simple server for sending and receiving push notifications.
// Authentication uses an application token (for sending messages) or client
// token (for management) passed in the X-Gotify-Key header.
//
// Usage:
//
//	c := gotify.New("http://localhost:80", "your-app-token")
//	err := c.CreateMessage(context.Background(), "Title", "Hello!", 5)
package gotify
