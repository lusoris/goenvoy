// Package tdarr provides a client for the Tdarr v2 API.
//
// Tdarr is a distributed media transcoding and health checking system.
// Authentication uses an optional API key passed in the x-api-key header.
//
// Usage:
//
//	c := tdarr.New("http://localhost:8265", tdarr.WithAPIKey("your-key"))
//	status, err := c.GetStatus(context.Background())
package tdarr
