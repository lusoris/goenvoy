// Package radarr provides a client for the Radarr v3 API.
//
// Radarr is a movie collection manager for Usenet and BitTorrent users
// that monitors multiple RSS feeds for new movies and automatically
// grabs, sorts, and renames them. The v3 API applies to current
// versions of the Radarr application.
//
// The [Client] type wraps [arr.BaseClient] and exposes typed methods
// for every major Radarr resource: movies, movie files, collections,
// credits, calendar, queue, commands, history, and more.
package radarr
