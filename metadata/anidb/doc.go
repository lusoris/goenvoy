// Package anidb provides a Go client for the AniDB HTTP XML API.
//
// AniDB is a community-maintained anime database. The HTTP API returns
// XML-encoded data and requires a registered client name and version.
//
// API reference: https://wiki.anidb.net/HTTP_API_Definition
//
// All API consumers must respect AniDB rate limits: no more than one
// request every two seconds, and identical requests must be cached locally.
package anidb
