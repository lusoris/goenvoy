# goenvoy

A collection of Go client libraries for media automation, metadata, and download service APIs.

Designed as a **multi-module monorepo** — import only what you need, from a single library to an entire category.

## Categories

| Category | Module | Services |
|---|---|---|
| **Arr Stack** | `github.com/lusoris/goenvoy/arr` | Sonarr, Radarr, Lidarr, Readarr, Whisparr (v2+v3), Prowlarr, Bazarr, Seerr, Autobrr, Mylar3, FlareSolverr |
| **Metadata** | `github.com/lusoris/goenvoy/metadata` | TMDb, TheTVDB, Fanart.tv, OMDb, TVmaze, AniList, Kitsu, AniDB, MAL, Trakt, Simkl, MusicBrainz, StashBox, TPDB, OpenSubtitles, Last.fm, Discogs, TheAudioDB, Open Library, Google Books |
| **Download Clients** | `github.com/lusoris/goenvoy/downloadclient` | qBittorrent, Transmission, Deluge, rTorrent, SABnzbd, NZBGet |
| **Media Servers** | `github.com/lusoris/goenvoy/mediaserver` | Plex, Jellyfin, Emby, Tautulli, Audiobookshelf, Komga, Navidrome, Kavita, Tdarr |
| **Notifications** | `github.com/lusoris/goenvoy/notification` | Gotify |
| **Anime** | `github.com/lusoris/goenvoy/anime` | Shoko Server |
| **Adult Media** | `github.com/lusoris/goenvoy/stash` | StashApp |

## Install

Import a specific service library:

```go
go get github.com/lusoris/goenvoy/arr/sonarr
```

Or import shared category types:

```go
go get github.com/lusoris/goenvoy/arr
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/lusoris/goenvoy/arr"
)

func main() {
    c, err := arr.NewBaseClient("http://localhost:8989", "your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    var status arr.StatusResponse
    if err := c.Get(context.Background(), "/api/v3/system/status", &status); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s %s\n", status.AppName, status.Version)
}
```

## Structure

Each category has a base module with shared types, plus sub-modules for individual services:

```
goenvoy/
├── arr/              # shared *arr types + base client
│   ├── sonarr/       # Sonarr v3 client
│   ├── radarr/       # Radarr v3 client
│   └── ...
├── metadata/         # shared metadata types (Rating, Image, Person, ...)
│   ├── tmdb/         # TMDb client
│   ├── tvdb/         # TheTVDB client
│   └── ...
├── downloadclient/   # shared download types + Downloader interface
│   ├── qbittorrent/  # qBittorrent WebUI client
│   └── ...
├── mediaserver/      # shared media server types
│   ├── plex/         # Plex Media Server client
│   ├── jellyfin/     # Jellyfin client
│   └── ...
├── anime/            # shared anime types
│   └── shoko/        # Shoko Server client
└── stash/            # StashApp GraphQL client
```

## Development

Requires Go 1.26+. See [CONTRIBUTING.md](CONTRIBUTING.md) for full details.

```bash
# Set up workspace (local dev, links all 56 modules)
go work init && find . -name 'go.mod' -not -path './.workingdir/*' -exec dirname {} \; | xargs go work use

# Run all tests
make test-all

# Lint all modules
make lint-all

# Tidy all modules
make tidy-all

# Format all modules
make fmt-all
```

## License

[MIT](LICENSE)
