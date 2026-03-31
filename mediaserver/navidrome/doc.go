// Package navidrome provides a client for the Navidrome music server using the
// Subsonic/OpenSubsonic API.
//
// Navidrome implements the Subsonic API v1.16.1 with OpenSubsonic extensions.
// Authentication uses token-based auth where token = md5(password + salt).
//
// Usage:
//
//	c := navidrome.New("http://localhost:4533", "admin", "password")
//	artists, err := c.GetArtists(context.Background())
package navidrome
