// Package deezer provides a client for the Deezer API.
//
// Deezer is a music streaming service providing access to artist, album, track,
// genre, chart, and radio metadata. No authentication is required for public data.
//
// Usage:
//
//	c := deezer.New()
//	artist, err := c.GetArtist(context.Background(), 27)
package deezer
