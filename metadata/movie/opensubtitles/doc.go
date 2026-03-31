// Package opensubtitles provides a client for the OpenSubtitles REST API.
//
// The API requires an API key passed via the Api-Key header. Some endpoints
// (e.g. download, user info) also require a Bearer token obtained via login.
//
// Usage:
//
//	c := opensubtitles.New("your-api-key")
//	results, err := c.Search(context.Background(), &opensubtitles.SearchParams{Query: "inception"})
package opensubtitles
