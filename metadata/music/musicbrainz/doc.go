// Package musicbrainz provides a client for the MusicBrainz API.
//
// MusicBrainz (https://musicbrainz.org) is an open music encyclopedia.
// The API provides access to artists, releases, release groups, recordings,
// labels, works, areas, events, genres, instruments, places, series, and URLs.
//
// # Authentication
//
// No API key is required. However, MusicBrainz requires a meaningful User-Agent
// string for identification. The default is "goenvoy/0.0.1 (https://github.com/golusoris/goenvoy)".
//
// # Rate Limiting
//
// MusicBrainz enforces a rate limit of 1 request per second. This client does
// not enforce rate limiting internally; callers should throttle requests.
//
// # Usage
//
//	client := musicbrainz.New()
//	artist, err := client.LookupArtist(ctx, "5b11f4ce-a62d-471e-81fc-a69a8278c7da", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
package musicbrainz
