// Package audiodb provides a client for TheAudioDB v1 API.
//
// TheAudioDB is a community-driven music metadata database providing artist,
// album, and track information including artwork.
//
// Usage:
//
//	c := audiodb.New("2") // free API key
//	artists, err := c.SearchArtist(context.Background(), "coldplay")
package audiodb
