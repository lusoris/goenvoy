// Package discogs provides a client for the Discogs API.
//
// The API uses a personal access token for authentication, passed via the
// Authorization header as "Discogs token=YOUR_TOKEN".
//
// Usage:
//
//	c := discogs.New("your-token")
//	release, err := c.GetRelease(context.Background(), 249504)
package discogs
