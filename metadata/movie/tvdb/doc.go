// Package tvdb provides a client for TheTVDB API v4.
//
// TheTVDB is a community-maintained TV, movie, and people metadata
// database. This client covers series, movies, episodes, seasons,
// people, search, artwork, genres, languages, and update endpoints.
//
// Authentication uses a two-step flow: provide an API key (and
// optional subscriber PIN), and the client transparently obtains a
// JWT bearer token via the /login endpoint on the first request.
// Obtain an API key from https://thetvdb.com/dashboard/account/apikey.
package tvdb
