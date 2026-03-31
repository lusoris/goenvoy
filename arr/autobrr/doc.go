// Package autobrr provides a client for the autobrr API.
//
// autobrr (https://autobrr.com) is a torrent/usenet automation tool for
// IRC announcements and RSS feeds. The API provides access to filters,
// indexers, IRC networks, feeds, download clients, and release history.
//
// # Authentication
//
// All requests require an API key passed via the X-API-Token header.
// Generate an API key from Settings → API keys in the autobrr dashboard.
//
// # Usage
//
//	client := autobrr.New("http://localhost:7474", "your-api-key")
//	filters, err := client.GetFilters(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
package autobrr
