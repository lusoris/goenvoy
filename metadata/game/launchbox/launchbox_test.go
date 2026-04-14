package launchbox_test

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golusoris/goenvoy/metadata"
	"github.com/golusoris/goenvoy/metadata/game/launchbox"
)

func makeTestZip(t *testing.T) []byte {
	t.Helper()

	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	metadataXML := `<?xml version="1.0" encoding="utf-8"?>
<LaunchBox>
  <Game>
    <DatabaseID>1</DatabaseID>
    <Name>Super Mario Bros.</Name>
    <Platform>Nintendo Entertainment System</Platform>
    <Developer>Nintendo</Developer>
    <Publisher>Nintendo</Publisher>
    <Genres>Platform</Genres>
    <MaxPlayers>2</MaxPlayers>
    <Overview>A classic platformer.</Overview>
  </Game>
  <Game>
    <DatabaseID>2</DatabaseID>
    <Name>The Legend of Zelda</Name>
    <Platform>Nintendo Entertainment System</Platform>
    <Developer>Nintendo</Developer>
    <Publisher>Nintendo</Publisher>
  </Game>
  <Game>
    <DatabaseID>3</DatabaseID>
    <Name>Sonic the Hedgehog</Name>
    <Platform>Sega Genesis</Platform>
    <Developer>Sonic Team</Developer>
    <Publisher>Sega</Publisher>
  </Game>
  <GameAlternateName>
    <DatabaseID>1</DatabaseID>
    <AlternateNameID>100</AlternateNameID>
    <AlternateName>Super Mario Bros</AlternateName>
    <Region>Japan</Region>
  </GameAlternateName>
  <GameImage>
    <DatabaseID>1</DatabaseID>
    <FileName>Images/Nintendo Entertainment System/Super Mario Bros./Box - Front.jpg</FileName>
    <Type>Box - Front</Type>
    <Region>North America</Region>
  </GameImage>
</LaunchBox>`

	platformsXML := `<?xml version="1.0" encoding="utf-8"?>
<LaunchBox>
  <Platform>
    <Name>Nintendo Entertainment System</Name>
    <Emulated>true</Emulated>
    <Developer>Nintendo</Developer>
    <Manufacturer>Nintendo</Manufacturer>
  </Platform>
  <Platform>
    <Name>Sega Genesis</Name>
    <Emulated>true</Emulated>
    <Manufacturer>Sega</Manufacturer>
  </Platform>
</LaunchBox>`

	f1, err := w.Create("Metadata.xml")
	if err != nil {
		t.Fatal(err)
	}
	f1.Write([]byte(metadataXML))

	f2, err := w.Create("Platforms.xml")
	if err != nil {
		t.Fatal(err)
	}
	f2.Write([]byte(platformsXML))

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func setup(t *testing.T) *launchbox.Client {
	t.Helper()

	zipData := makeTestZip(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Write(zipData)
	}))
	t.Cleanup(srv.Close)

	c := launchbox.New(metadata.WithBaseURL(srv.URL))
	if err := c.Download(context.Background()); err != nil {
		t.Fatal(err)
	}
	return c
}

func TestDownloadAndGameCount(t *testing.T) {
	t.Parallel()

	c := setup(t)
	if c.GameCount() != 3 {
		t.Errorf("GameCount = %d, want 3", c.GameCount())
	}
}

func TestGetGameByID(t *testing.T) {
	t.Parallel()

	c := setup(t)

	game := c.GetGameByID(1)
	if game == nil {
		t.Fatal("game not found")
	}
	if game.Name != "Super Mario Bros." {
		t.Errorf("Name = %q", game.Name)
	}
	if game.Platform != "Nintendo Entertainment System" {
		t.Errorf("Platform = %q", game.Platform)
	}
	if game.Developer != "Nintendo" {
		t.Errorf("Developer = %q", game.Developer)
	}

	if c.GetGameByID(999) != nil {
		t.Error("expected nil for unknown ID")
	}
}

func TestSearchGames(t *testing.T) {
	t.Parallel()

	c := setup(t)

	// Search by name.
	results := c.SearchGames("mario", "")
	if len(results) != 1 {
		t.Fatalf("len = %d, want 1", len(results))
	}
	if results[0].DatabaseID != 1 {
		t.Errorf("DatabaseID = %d, want 1", results[0].DatabaseID)
	}

	// Search with platform filter.
	results = c.SearchGames("", "Nintendo Entertainment System")
	if len(results) != 2 {
		t.Errorf("len = %d, want 2", len(results))
	}

	// Search no match.
	results = c.SearchGames("nonexistent", "")
	if len(results) != 0 {
		t.Errorf("len = %d, want 0", len(results))
	}
}

func TestGetAlternateNames(t *testing.T) {
	t.Parallel()

	c := setup(t)

	names := c.GetAlternateNames(1)
	if len(names) != 1 {
		t.Fatalf("len = %d, want 1", len(names))
	}
	if names[0].Name != "Super Mario Bros" {
		t.Errorf("Name = %q", names[0].Name)
	}
	if names[0].Region != "Japan" {
		t.Errorf("Region = %q", names[0].Region)
	}
}

func TestGetImages(t *testing.T) {
	t.Parallel()

	c := setup(t)

	images := c.GetImages(1)
	if len(images) != 1 {
		t.Fatalf("len = %d, want 1", len(images))
	}
	if images[0].Type != "Box - Front" {
		t.Errorf("Type = %q", images[0].Type)
	}
}

func TestGetPlatforms(t *testing.T) {
	t.Parallel()

	c := setup(t)

	platforms := c.GetPlatforms()
	if len(platforms) != 2 {
		t.Fatalf("len = %d, want 2", len(platforms))
	}
	if platforms[0].Name != "Nintendo Entertainment System" {
		t.Errorf("Name = %q", platforms[0].Name)
	}
}

func TestDownloadError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	t.Cleanup(srv.Close)

	c := launchbox.New(metadata.WithBaseURL(srv.URL))
	err := c.Download(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *launchbox.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *launchbox.APIError", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want 500", apiErr.StatusCode)
	}
}
