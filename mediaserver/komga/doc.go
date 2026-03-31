// Package komga provides a client for the Komga API.
//
// Komga (https://komga.org) is a free and open source comics/mangas media
// server. The API provides access to libraries, series, books, collections,
// read lists, and user management.
//
// # Authentication
//
// Requests use HTTP Basic Authentication with your Komga username and password.
//
// # Usage
//
//	client := komga.New("http://localhost:25600", "admin@example.com", "password")
//	libraries, err := client.GetLibraries(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
package komga
