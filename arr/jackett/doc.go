// Package jackett provides a client for the Jackett API.
//
// Jackett is a Torznab/Newznab proxy that translates queries from
// *arr applications to tracker site queries.
// Authentication uses an API key passed as a query parameter.
//
// Usage:
//
//	c := jackett.New("http://localhost:9117", "your-api-key")
//	results, err := c.Search(context.Background(), "ubuntu", nil)
package jackett
