// Package flaresolverr provides a client for the FlareSolverr API.
//
// FlareSolverr is a proxy server to bypass Cloudflare and DDoS-GUARD protection.
// It uses a headless browser to solve challenges and return cookies/HTML.
//
// Usage:
//
//	c := flaresolverr.New("http://localhost:8191")
//	resp, err := c.Get(context.Background(), "https://example.com", nil)
package flaresolverr
