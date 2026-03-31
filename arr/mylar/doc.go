// Package mylar provides a client for the Mylar3 API.
//
// Mylar3 is an automated comic book downloader and manager.
// Authentication uses an API key passed as a query parameter.
//
// Usage:
//
//	c := mylar.New("http://localhost:8090", "your-api-key")
//	comics, err := c.GetIndex(context.Background())
package mylar
