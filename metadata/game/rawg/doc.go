// Package rawg provides a client for the RAWG Video Games Database API.
//
// RAWG provides information about video games, platforms, genres, publishers,
// developers, tags, and stores. Authentication is via an API key passed as a
// query parameter.
//
// Usage:
//
//	c := rawg.New("your-api-key")
//	result, err := c.SearchGames(context.Background(), "zelda", 1, 10)
package rawg
