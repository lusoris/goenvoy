// Package googlebooks provides a client for the Google Books API v1.
//
// Google Books provides access to the world's largest digital book catalog.
// Authentication uses an API key passed as a query parameter.
//
// Usage:
//
//	c := googlebooks.New("your-api-key")
//	results, err := c.Search(context.Background(), "flowers for algernon")
package googlebooks
