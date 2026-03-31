// Package steam provides a client for the Steam Store API and Steam Web API.
//
// Steam provides game metadata, player counts, news, and achievement data
// through two API surfaces: the Store API for game information and the Web
// API for community and stats data.
//
// Usage:
//
//	c := steam.New()
//	details, err := c.GetAppDetails(context.Background(), 730)
package steam
