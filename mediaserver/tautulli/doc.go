// Package tautulli provides a client for the Tautulli API.
//
// Tautulli (https://tautulli.com) is a monitoring and tracking tool for
// Plex Media Server. The API provides access to activity streams, watch
// history, library statistics, user analytics, and server information.
//
// # Authentication
//
// All requests require an API key, which can be found in
// Tautulli Settings → Web Interface → API key.
//
// # Usage
//
//	client := tautulli.New("http://localhost:8181", "your-api-key")
//	activity, err := client.GetActivity(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
package tautulli
