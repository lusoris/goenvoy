// Package kavita provides a client for the Kavita REST API.
//
// Kavita is a self-hosted digital library for comics, manga, and ebooks.
// Authentication uses an API key exchanged for a JWT token.
//
// Usage:
//
//	c := kavita.New("http://localhost:5000", "your-api-key")
//	libraries, err := c.GetLibraries(context.Background())
package kavita
