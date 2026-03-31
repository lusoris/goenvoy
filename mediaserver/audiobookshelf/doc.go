// Package audiobookshelf provides a client for the Audiobookshelf API.
//
// Audiobookshelf (https://www.audiobookshelf.org) is a self-hosted audiobook
// and podcast server. The API provides access to libraries, audiobooks,
// podcasts, users, playback sessions, and server information.
//
// # Authentication
//
// Requests require a Bearer token obtained via the /login endpoint or an
// API token from your user settings.
//
// # Usage
//
//	client := audiobookshelf.New("http://localhost:13378", "your-token")
//	libs, err := client.GetLibraries(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
package audiobookshelf
