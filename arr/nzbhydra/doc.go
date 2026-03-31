// Package nzbhydra provides a client for the NZBHydra2 API.
//
// NZBHydra2 is a meta search application for NZB indexers that provides
// a Newznab-compatible search API plus statistics and history endpoints.
// Authentication uses an API key passed as a query parameter.
//
// Usage:
//
//	c := nzbhydra.New("http://localhost:5076", "your-api-key")
//	results, err := c.Search(context.Background(), "ubuntu", nil)
package nzbhydra
