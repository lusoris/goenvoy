// Package openlibrary provides a client for the Open Library API.
//
// Open Library is an open, editable library catalog providing free access to
// book metadata. No authentication is required.
//
// Usage:
//
//	c := openlibrary.New()
//	results, err := c.Search(context.Background(), "the lord of the rings")
package openlibrary
