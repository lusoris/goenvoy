// Package letterboxd provides a client for the Letterboxd v0 API.
//
// Letterboxd is a social film discovery and logging service.
// This client covers films, lists, log entries (diary/reviews),
// members, contributors, film collections, stories, search, and
// news endpoints.
//
// Authentication uses OAuth2 Bearer tokens. Obtain credentials from
// https://letterboxd.com/api-beta/. Public endpoints can be accessed
// with a client credentials token; member-authenticated endpoints
// require an authorization code flow token.
package letterboxd
