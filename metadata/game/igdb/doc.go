// Package igdb provides a client for the IGDB (Internet Game Database) v4 API.
//
// IGDB is powered by Twitch and uses OAuth2 client credentials for
// authentication. All requests use the APICalypse query language sent as
// POST request bodies.
//
// Users must obtain a Twitch access token via the client credentials flow
// before using this client. See https://api-docs.igdb.com for details.
//
// Usage:
//
//	c := igdb.New("your-client-id", "your-access-token")
//	games, err := c.SearchGames(context.Background(), "zelda", 10)
package igdb
