// Package lastfm provides a client for the Last.fm API.
//
// The API uses a simple API key for read-only access. All requests go to
// ws.audioscrobbler.com/2.0/ with method and api_key query parameters.
//
// Usage:
//
//	c := lastfm.New("your-api-key")
//	artist, err := c.GetArtistInfo(context.Background(), "Radiohead")
package lastfm
