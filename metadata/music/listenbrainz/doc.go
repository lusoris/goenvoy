// Package listenbrainz provides a client for the ListenBrainz API.
//
// ListenBrainz is an open-source music listening data platform that tracks
// listening habits and provides statistics, recommendations, and social
// features.
//
// Usage:
//
//	c := listenbrainz.New() // no token needed for read-only access
//	listens, err := c.GetUserListens(context.Background(), "username", 25)
package listenbrainz
