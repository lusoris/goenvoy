# goenvoy v1.2.0 — Improvement Plan

> Based on a full audit of [lusoris/revenge](https://github.com/lusoris/revenge) against goenvoy v1.1.0.
> Every item below describes a concrete gap that forces revenge to maintain local workarounds.
> Implementing these changes would eliminate **~26 local types** and **~540 LOC** of wrapper/glue code.

---

## Table of Contents

1. [P1 — TMDb: `append_to_response` Support](#p1--tmdb-append_to_response-support)
2. [P2 — TMDb: Missing Type Fields](#p2--tmdb-missing-type-fields)
3. [P3 — Kitsu: Anime Mappings & Per-Anime Categories](#p3--kitsu-anime-mappings--per-anime-categories)
4. [P4 — AniList: Extended Media Query (Connections)](#p4--anilist-extended-media-query-connections)
5. [P5 — AniDB: Title Dump Search](#p5--anidb-title-dump-search)
6. [P6 — Letterboxd: OAuth2 Client Credentials](#p6--letterboxd-oauth2-client-credentials)
7. [Summary](#summary)

---

## P1 — TMDb: `append_to_response` Support

### Problem

TMDb's detail endpoints (`/movie/{id}`, `/tv/{id}`, `/tv/{id}/season/{n}`, `/person/{id}`) accept an `append_to_response` query parameter that bundles multiple sub-resources into a single API call. This is critical for performance — without it, fetching a movie's full metadata requires 6+ separate HTTP requests instead of 1.

goenvoy's current signatures **do not accept query parameters**:

```go
// Current signatures (goenvoy v1.1.0)
func (c *Client) GetMovie(ctx context.Context, id int, language string) (*MovieDetails, error)
func (c *Client) GetTV(ctx context.Context, id int, language string) (*TVDetails, error)
func (c *Client) GetTVSeason(ctx context.Context, tvID, seasonNumber int, language string) (*SeasonDetails, error)
func (c *Client) GetPerson(ctx context.Context, id int, language string) (*PersonDetails, error)
// GetTVEpisode exists, but also lacks append_to_response support
func (c *Client) GetTVEpisode(ctx context.Context, tvID, seasonNumber, episodeNumber int, language string) (*EpisodeDetails, error)
```

### Current Workaround in revenge

revenge defines **5 wrapper types** that embed goenvoy detail types and add fields for the appended sub-resources:

```go
// revenge/internal/service/metadata/providers/tmdb/types.go

type MovieResponse struct {
    tmdbapi.MovieDetails
    Credits      *tmdbapi.Credits              `json:"credits,omitempty"`
    Images       *ImagesWithStills             `json:"images,omitempty"`
    ReleaseDates *tmdbapi.ReleaseDatesResponse `json:"release_dates,omitempty"`
    Translations *tmdbapi.TranslationsResponse `json:"translations,omitempty"`
    ExternalIDs  *ExternalIDsResponse          `json:"external_ids,omitempty"`
    Videos       *tmdbapi.VideosResponse       `json:"videos,omitempty"`
}

type TVResponse struct {
    tmdbapi.TVDetails
    Credits        *tmdbapi.Credits                `json:"credits,omitempty"`
    Images         *ImagesWithStills               `json:"images,omitempty"`
    ContentRatings *tmdbapi.ContentRatingsResponse `json:"content_ratings,omitempty"`
    Translations   *tmdbapi.TranslationsResponse   `json:"translations,omitempty"`
    ExternalIDs    *ExternalIDsResponse            `json:"external_ids,omitempty"`
    Videos         *tmdbapi.VideosResponse         `json:"videos,omitempty"`
}

type SeasonResponse struct {
    tmdbapi.SeasonDetails
    Credits      *tmdbapi.Credits              `json:"credits,omitempty"`
    Images       *ImagesWithStills             `json:"images,omitempty"`
    Translations *tmdbapi.TranslationsResponse `json:"translations,omitempty"`
}

type EpisodeResponse struct {
    tmdbapi.EpisodeDetails
    Credits *tmdbapi.Credits  `json:"credits,omitempty"`
    Images  *ImagesWithStills `json:"images,omitempty"`
}

type PersonResponse struct {
    tmdbapi.PersonDetails
    ExternalIDs *ExternalIDsResponse  `json:"external_ids,omitempty"`
    Images      *tmdbapi.PersonImages `json:"images,omitempty"`
}
```

revenge then manually builds the URL with query parameters:

```go
func (c *Client) GetMovie(ctx context.Context, id int, language string, appendToResponse string) (*MovieResponse, error) {
    params := make(map[string]string)
    if appendToResponse != "" {
        params["append_to_response"] = appendToResponse
    }
    // ... manual HTTP call with params
}
```

### Proposed Change

Add a `RequestOption` functional option type and enriched "Full" response types:

```go
// metadata/video/tmdb/tmdb.go

// RequestOption configures a single API request.
type RequestOption func(*requestConfig)

type requestConfig struct {
    appendToResponse string
}

// WithAppendToResponse bundles additional sub-resources into the response.
// fields is a comma-separated list, e.g. "credits,images,videos".
func WithAppendToResponse(fields string) RequestOption {
    return func(cfg *requestConfig) { cfg.appendToResponse = fields }
}

// GetMovieFull returns movie details with optional appended sub-resources.
func (c *Client) GetMovieFull(ctx context.Context, id int, language string, opts ...RequestOption) (*MovieDetailsFull, error)

// GetTVFull returns TV show details with optional appended sub-resources.
func (c *Client) GetTVFull(ctx context.Context, id int, language string, opts ...RequestOption) (*TVDetailsFull, error)

// GetTVSeasonFull returns season details with optional appended sub-resources.
func (c *Client) GetTVSeasonFull(ctx context.Context, tvID, seasonNumber int, language string, opts ...RequestOption) (*SeasonDetailsFull, error)

// GetTVEpisodeFull returns episode details with optional appended sub-resources.
func (c *Client) GetTVEpisodeFull(ctx context.Context, tvID, seasonNumber, episodeNumber int, language string, opts ...RequestOption) (*EpisodeDetailsFull, error)

// GetPersonFull returns person details with optional appended sub-resources.
func (c *Client) GetPersonFull(ctx context.Context, id int, language string, opts ...RequestOption) (*PersonDetailsFull, error)
```

New types (in `types.go`):

```go
// MovieDetailsFull embeds MovieDetails with optional appended sub-resources.
type MovieDetailsFull struct {
    MovieDetails
    Credits      *Credits              `json:"credits,omitempty"`
    Images       *Images               `json:"images,omitempty"`
    ReleaseDates *ReleaseDatesResponse `json:"release_dates,omitempty"`
    Translations *TranslationsResponse `json:"translations,omitempty"`
    ExternalIDs  *ExternalIDs          `json:"external_ids,omitempty"`
    Videos       *VideosResponse       `json:"videos,omitempty"`
}

// TVDetailsFull embeds TVDetails with optional appended sub-resources.
type TVDetailsFull struct {
    TVDetails
    Credits        *Credits                `json:"credits,omitempty"`
    Images         *Images                 `json:"images,omitempty"`
    ContentRatings *ContentRatingsResponse `json:"content_ratings,omitempty"`
    Translations   *TranslationsResponse   `json:"translations,omitempty"`
    ExternalIDs    *ExternalIDs            `json:"external_ids,omitempty"`
    Videos         *VideosResponse         `json:"videos,omitempty"`
}

// SeasonDetailsFull embeds SeasonDetails with optional appended sub-resources.
type SeasonDetailsFull struct {
    SeasonDetails
    Credits      *Credits              `json:"credits,omitempty"`
    Images       *Images               `json:"images,omitempty"`
    Translations *TranslationsResponse `json:"translations,omitempty"`
    ExternalIDs  *ExternalIDs          `json:"external_ids,omitempty"`
    Videos       *VideosResponse       `json:"videos,omitempty"`
}

// EpisodeDetailsFull embeds EpisodeDetails with optional appended sub-resources.
type EpisodeDetailsFull struct {
    EpisodeDetails
    Credits      *Credits              `json:"credits,omitempty"`
    Images       *Images               `json:"images,omitempty"`
    Translations *TranslationsResponse `json:"translations,omitempty"`
    ExternalIDs  *ExternalIDs          `json:"external_ids,omitempty"`
    Videos       *VideosResponse       `json:"videos,omitempty"`
}

// PersonDetailsFull embeds PersonDetails with optional appended sub-resources.
type PersonDetailsFull struct {
    PersonDetails
    ExternalIDs *ExternalIDs  `json:"external_ids,omitempty"`
    Images      *PersonImages `json:"images,omitempty"`
}
```

### Impact

- **Eliminates**: `MovieResponse`, `TVResponse`, `SeasonResponse`, `EpisodeResponse`, `PersonResponse` in revenge
- **LOC saved**: ~120

---

## P2 — TMDb: Missing Type Fields

### P2a — `Images`: Add `Stills` Field

**Problem**: TMDb returns a `stills` array for episode images, but goenvoy's `Images` struct omits it.

```go
// Current goenvoy (v1.1.0)
type Images struct {
    ID        int         `json:"id"`
    Backdrops []ImageItem `json:"backdrops,omitempty"`
    Posters   []ImageItem `json:"posters,omitempty"`
    Logos     []ImageItem `json:"logos,omitempty"`
}
```

**Workaround in revenge**:

```go
type ImagesWithStills struct {
    tmdbapi.Images
    Stills []tmdbapi.ImageItem `json:"stills,omitempty"`
}
```

**Fix**: Add `Stills` to the existing `Images` type:

```go
type Images struct {
    ID        int         `json:"id"`
    Backdrops []ImageItem `json:"backdrops,omitempty"`
    Posters   []ImageItem `json:"posters,omitempty"`
    Logos     []ImageItem `json:"logos,omitempty"`
    Stills    []ImageItem `json:"stills,omitempty"`    // NEW
}
```

### P2b — `ExternalIDs`: Add 4 Missing Fields

**Problem**: TMDb returns additional social/external ID fields that goenvoy's `ExternalIDs` type omits.

```go
// Current goenvoy (v1.1.0)
type ExternalIDs struct {
    ID          int    `json:"id"`
    IMDbID      string `json:"imdb_id,omitempty"`
    FacebookID  string `json:"facebook_id,omitempty"`
    InstagramID string `json:"instagram_id,omitempty"`
    TwitterID   string `json:"twitter_id,omitempty"`
    WikidataID  string `json:"wikidata_id,omitempty"`
    TVDbID      int    `json:"tvdb_id,omitempty"`
    TVRageID    int    `json:"tvrage_id,omitempty"`
}
```

**Workaround in revenge**:

```go
type ExternalIDsResponse struct {
    tmdbapi.ExternalIDs
    TikTokID    string `json:"tiktok_id,omitempty"`
    YouTubeID   string `json:"youtube_id,omitempty"`
    FreebaseID  string `json:"freebase_id,omitempty"`
    FreebaseMID string `json:"freebase_mid,omitempty"`
}
```

**Fix**: Add the 4 fields to the existing type:

```go
type ExternalIDs struct {
    ID          int    `json:"id"`
    IMDbID      string `json:"imdb_id,omitempty"`
    FacebookID  string `json:"facebook_id,omitempty"`
    InstagramID string `json:"instagram_id,omitempty"`
    TwitterID   string `json:"twitter_id,omitempty"`
    WikidataID  string `json:"wikidata_id,omitempty"`
    TVDbID      int    `json:"tvdb_id,omitempty"`
    TVRageID    int    `json:"tvrage_id,omitempty"`
    TikTokID    string `json:"tiktok_id,omitempty"`     // NEW
    YouTubeID   string `json:"youtube_id,omitempty"`    // NEW
    FreebaseID  string `json:"freebase_id,omitempty"`   // NEW
    FreebaseMID string `json:"freebase_mid,omitempty"`  // NEW
}
```

### P2c — `PersonResult`: Add `KnownFor` Field

**Problem**: TMDb's person search returns a `known_for` array with the person's most notable works. goenvoy's `PersonResult` omits this entirely.

```go
// Current goenvoy (v1.1.0)
type PersonResult struct {
    ID                 int     `json:"id"`
    Name               string  `json:"name,omitempty"`
    ProfilePath        string  `json:"profile_path,omitempty"`
    Adult              bool    `json:"adult,omitempty"`
    Popularity         float64 `json:"popularity,omitempty"`
    KnownForDepartment string  `json:"known_for_department,omitempty"`
    Gender             int     `json:"gender,omitempty"`
}
```

**Workaround in revenge** (3 custom types):

```go
type PersonSearchResultsResponse struct {
    Page         int                    `json:"page"`
    Results      []PersonSearchResponse `json:"results"`
    TotalPages   int                    `json:"total_pages"`
    TotalResults int                    `json:"total_results"`
}

type PersonSearchResponse struct {
    ID          int                `json:"id"`
    Name        string             `json:"name"`
    ProfilePath string             `json:"profile_path,omitempty"`
    Popularity  float64            `json:"popularity"`
    Adult       bool               `json:"adult"`
    KnownFor    []KnownForResponse `json:"known_for"`
}

type KnownForResponse struct {
    MediaType    string `json:"media_type"`
    ID           int    `json:"id"`
    Title        string `json:"title"`
    Name         string `json:"name"`
    PosterPath   string `json:"poster_path,omitempty"`
    ReleaseDate  string `json:"release_date"`
    FirstAirDate string `json:"first_air_date"`
}
```

**Fix**: Add a `KnownForItem` type and a `KnownFor` field to `PersonResult`:

```go
// KnownForItem represents a media entry a person is known for.
type KnownForItem struct {
    MediaType        string  `json:"media_type"`
    ID               int     `json:"id"`
    Title            string  `json:"title,omitempty"`
    Name             string  `json:"name,omitempty"`
    OriginalTitle    string  `json:"original_title,omitempty"`
    OriginalName     string  `json:"original_name,omitempty"`
    Overview         string  `json:"overview,omitempty"`
    PosterPath       string  `json:"poster_path,omitempty"`
    BackdropPath     string  `json:"backdrop_path,omitempty"`
    ReleaseDate      string  `json:"release_date,omitempty"`
    FirstAirDate     string  `json:"first_air_date,omitempty"`
    VoteAverage      float64 `json:"vote_average,omitempty"`
    VoteCount        int     `json:"vote_count,omitempty"`
    GenreIDs         []int   `json:"genre_ids,omitempty"`
    OriginalLanguage string  `json:"original_language,omitempty"`
    Adult            bool    `json:"adult,omitempty"`
}

type PersonResult struct {
    ID                 int            `json:"id"`
    Name               string         `json:"name,omitempty"`
    ProfilePath        string         `json:"profile_path,omitempty"`
    Adult              bool           `json:"adult,omitempty"`
    Popularity         float64        `json:"popularity,omitempty"`
    KnownForDepartment string         `json:"known_for_department,omitempty"`
    Gender             int            `json:"gender,omitempty"`
    KnownFor           []KnownForItem `json:"known_for,omitempty"` // NEW
}
```

### P2 Impact

- **Eliminates**: `ImagesWithStills`, `ExternalIDsResponse`, `PersonSearchResultsResponse`, `PersonSearchResponse`, `KnownForResponse` (5 types)
- **P1 + P2 combined**: eliminates ALL 12 local TMDb types in revenge — `types.go` becomes empty/deleted
- **LOC saved**: ~80

---

## P3 — Kitsu: Anime Mappings & Per-Anime Categories

### Problem

revenge needs two Kitsu endpoints that goenvoy doesn't expose:

1. **Anime Mappings** — external ID mappings (MAL, TVDb, AniDB, AniList) for a specific anime
   - Endpoint: `GET /api/edge/anime/{id}/mappings?filter[externalSite]=myanimelist/anime,thetvdb/series,thetvdb,anidb,anilist/anime`
   - Returns JSON:API collection of `{externalSite, externalId}` pairs

2. **Per-Anime Categories** — categories for a specific anime (not the global category list)
   - Endpoint: `GET /api/edge/anime/{id}/categories?page[limit]=20`
   - goenvoy has `GetCategories(ctx, limit, offset)` (global) and `GetCategory(ctx, id)`, but no per-anime variant

### Current Workaround in revenge

revenge defines 3 local types and implements raw JSON:API HTTP calls:

```go
// revenge/internal/service/metadata/providers/kitsu/types.go

type jsonAPIList[T any] struct {
    Data []struct {
        ID         string `json:"id"`
        Type       string `json:"type"`
        Attributes T      `json:"attributes"`
    } `json:"data"`
}

type MappingEntry struct {
    ExternalSite string `json:"externalSite"`
    ExternalID   string `json:"externalId"`
}

type CategoryEntry struct {
    Title string `json:"title"`
}
```

revenge's client then does manual HTTP calls using a `getJSONAPI` helper method, parses the JSON:API envelope, and extracts the `Attributes` from each `Data` entry.

### Proposed Change

Add a `Mapping` type and two new client methods to the Kitsu module:

```go
// metadata/anime/kitsu/types.go

// Mapping holds an external ID mapping for an anime.
type Mapping struct {
    ID           string `json:"id"`
    ExternalSite string `json:"externalSite"`
    ExternalID   string `json:"externalId"`
}
```

```go
// metadata/anime/kitsu/kitsu.go

// GetAnimeMappings returns external ID mappings for an anime.
// sites filters by external site names (e.g. "myanimelist/anime", "thetvdb/series", "anidb", "anilist/anime").
// If sites is empty, all mappings are returned.
func (c *Client) GetAnimeMappings(ctx context.Context, animeID int64, sites ...string) ([]Mapping, error)

// GetAnimeCategories returns categories for a specific anime.
func (c *Client) GetAnimeCategories(ctx context.Context, animeID int64, limit int) ([]Category, error)
```

The internal implementation should handle the JSON:API envelope parsing internally, so consumers get clean typed results.

### Impact

- **Eliminates**: `jsonAPIList[T]`, `MappingEntry`, `CategoryEntry` + raw HTTP wrapper code (~60 LOC)
- **LOC saved**: ~60

---

## P4 — AniList: Extended Media Query (Connections)

### Problem

goenvoy's `GetMedia()` uses a GraphQL query that only fetches "basic" media fields (defined in the `mediaFields` fragment). It does **not** request connection fields that are essential for a media server:

- **Studios** (with `isMain` flag to identify the primary studio)
- **Characters** (with roles and voice actors)
- **Staff** (with roles)
- **Relations** (related media with relationship type)
- **External Links** (links to streaming sites, official pages)
- **Trailer** (YouTube/Dailymotion trailer reference)
- **Streaming Episodes** (episode streaming sources)

### Current Workaround in revenge

revenge defines **13 local types** and a custom 80-line GraphQL query:

```go
// revenge/internal/service/metadata/providers/anilist/types.go

type DetailedMedia struct {
    anilistapi.Media                                // embedded goenvoy type
    Studios           StudioConnection              `json:"studios"`
    ExternalLinks     []ExternalLink                `json:"externalLinks"`
    StreamingEpisodes []StreamingEpisode            `json:"streamingEpisodes"`
    Trailer           *Trailer                      `json:"trailer"`
    Characters        CharacterConnection           `json:"characters"`
    Staff             StaffConnection               `json:"staff"`
    Relations         MediaConnection               `json:"relations"`
}

type StudioConnection struct { Edges []StudioEdge }
type StudioEdge struct {
    Node   Studio `json:"node"`
    IsMain bool   `json:"isMain"`
}
type Studio struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
type ExternalLink struct {
    ID       int     `json:"id"`
    URL      *string `json:"url"`
    Site     string  `json:"site"`
    SiteID   *int    `json:"siteId"`
    Type     string  `json:"type"`
    Language *string `json:"language"`
}
type StreamingEpisode struct {
    Title     *string `json:"title"`
    Thumbnail *string `json:"thumbnail"`
    URL       *string `json:"url"`
    Site      *string `json:"site"`
}
type Trailer struct {
    ID        *string `json:"id"`
    Site      *string `json:"site"`
    Thumbnail *string `json:"thumbnail"`
}
type CharacterConnection struct { Edges []CharacterEdge }
type CharacterEdge struct {
    Node        anilistapi.Character `json:"node"`
    Role        string               `json:"role"`
    VoiceActors []anilistapi.Staff   `json:"voiceActors"`
}
type StaffConnection struct { Edges []StaffEdge }
type StaffEdge struct {
    Node anilistapi.Staff `json:"node"`
    Role string           `json:"role"`
}
type MediaConnection struct { Edges []MediaEdge }
type MediaEdge struct {
    Node         anilistapi.Media `json:"node"`
    RelationType string           `json:"relationType"`
}
```

revenge uses goenvoy's `Query()` raw method with a custom GraphQL query (the `mediaQuery` const) to fetch all of this in one call.

### Proposed Change

Add the connection types to goenvoy's AniList module and provide a `GetMediaDetailed()` method with an expanded GraphQL query.

New types in `types.go`:

```go
// Studio represents an animation studio.
type Studio struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// StudioEdge connects a studio to a media entry.
type StudioEdge struct {
    Node   Studio `json:"node"`
    IsMain bool   `json:"isMain"`
}

// StudioConnection contains studio edges.
type StudioConnection struct {
    Edges []StudioEdge `json:"edges"`
}

// CharacterEdge connects a character to a media entry with role and voice actors.
type CharacterEdge struct {
    Node        Character `json:"node"`
    Role        string    `json:"role"`
    VoiceActors []Staff   `json:"voiceActors"`
}

// CharacterConnectionDetailed contains character edges with roles and voice actors.
type CharacterConnectionDetailed struct {
    Edges []CharacterEdge `json:"edges"`
}

// StaffEdge connects a staff member to a media entry.
type StaffEdge struct {
    Node Staff  `json:"node"`
    Role string `json:"role"`
}

// StaffConnectionDetailed contains staff edges with roles.
type StaffConnectionDetailed struct {
    Edges []StaffEdge `json:"edges"`
}

// MediaEdge connects related media entries.
type MediaEdge struct {
    Node         Media  `json:"node"`
    RelationType string `json:"relationType"`
}

// MediaConnectionDetailed contains related media edges.
type MediaConnectionDetailed struct {
    Edges []MediaEdge `json:"edges"`
}

// ExternalLink is a link to an external site for a media entry.
type ExternalLink struct {
    ID       int     `json:"id"`
    URL      *string `json:"url"`
    Site     string  `json:"site"`
    SiteID   *int    `json:"siteId"`
    Type     string  `json:"type"`
    Language *string `json:"language"`
}

// StreamingEpisode is a streaming source for an episode.
type StreamingEpisode struct {
    Title     *string `json:"title"`
    Thumbnail *string `json:"thumbnail"`
    URL       *string `json:"url"`
    Site      *string `json:"site"`
}

// Trailer holds a media trailer reference.
type Trailer struct {
    ID        *string `json:"id"`
    Site      *string `json:"site"`
    Thumbnail *string `json:"thumbnail"`
}

// MediaDetailed extends Media with connection fields (studios, characters, staff, relations, etc.).
type MediaDetailed struct {
    Media
    Studios           StudioConnection             `json:"studios"`
    Characters        CharacterConnectionDetailed  `json:"characters"`
    Staff             StaffConnectionDetailed      `json:"staff"`
    Relations         MediaConnectionDetailed      `json:"relations"`
    ExternalLinks     []ExternalLink               `json:"externalLinks"`
    StreamingEpisodes []StreamingEpisode           `json:"streamingEpisodes"`
    Trailer           *Trailer                     `json:"trailer"`
}
```

New GraphQL fragment and method in `anilist.go`:

```go
const mediaDetailedFields = mediaFields + `
    studios { edges { node { id name } isMain } }
    characters(sort: [ROLE, RELEVANCE], perPage: 25) {
        edges {
            node { ` + characterFields + ` }
            role
            voiceActors(language: JAPANESE, sort: RELEVANCE) {
                ` + staffFields + `
            }
        }
    }
    staff(sort: RELEVANCE, perPage: 25) {
        edges {
            node { ` + staffFields + ` }
            role
        }
    }
    relations {
        edges {
            node { ` + mediaSearchFields + ` }
            relationType
        }
    }
    externalLinks { id url site siteId type language }
    streamingEpisodes { title thumbnail url site }
    trailer { id site thumbnail }
`

const queryGetMediaDetailed = `query ($id: Int) { Media(id: $id) {` + mediaDetailedFields + `} }`
const queryGetMediaDetailedByMalID = `query ($idMal: Int, $type: MediaType) { Media(idMal: $idMal, type: $type) {` + mediaDetailedFields + `} }`

// GetMediaDetailed returns a media entry with all connection fields (studios, characters, staff, relations, etc.).
func (c *Client) GetMediaDetailed(ctx context.Context, id int) (*MediaDetailed, error)

// GetMediaDetailedByMalID returns a detailed media entry by its MyAnimeList ID and type.
func (c *Client) GetMediaDetailedByMalID(ctx context.Context, malID int, mediaType MediaType) (*MediaDetailed, error)
```

### Impact

- **Eliminates**: `DetailedMedia`, `StudioConnection`, `StudioEdge`, `Studio`, `ExternalLink`, `StreamingEpisode`, `Trailer`, `CharacterConnection`, `CharacterEdge`, `StaffConnection`, `StaffEdge`, `MediaConnection`, `MediaEdge` (13 types) + custom `mediaQuery` GraphQL const (~80 lines)
- **LOC saved**: ~150

---

## P5 — AniDB: Title Dump Search

### Problem

AniDB does not have a search API endpoint. The only way to search anime by title is to download their daily title dump file (`http://anidb.net/api/anime-titles.dat.gz`), a gzip-compressed pipe-delimited text file:

```
# Format: aid|type|lang|title
4|1|en|Cowboy Bebop
4|2|ja|カウボーイビバップ
13|1|en|Angel Beats!
```

goenvoy's AniDB module provides `GetAnime(id)` (XML API), `HotAnime()`, `RandomSimilar()`, and `MainPage()` — but **no title search capability**.

### Current Workaround in revenge

```go
// revenge/internal/service/metadata/providers/anidb/types.go
type TitleDumpEntry struct {
    AID   int
    Type  string
    Lang  string
    Title string
}
```

revenge implements:
- `loadTitleDump(ctx)` — downloads + decompresses `anime-titles.dat.gz`, parses into `[]TitleDumpEntry`
- `parseTitleDump(r io.Reader)` — line-by-line pipe-delimited parser with comment/blank line handling
- `SearchAnime(ctx, query, limit)` — in-memory search with ranking: exact match → prefix match → contains match
- 24-hour TTL caching of the parsed dump

### Proposed Change

Add title dump support to the goenvoy AniDB module:

```go
// metadata/anime/anidb/types.go

// TitleEntry represents a single entry from the AniDB anime titles dump.
type TitleEntry struct {
    AID   int    // Anime ID
    Type  string // Title type: "1" (primary), "2" (synonym), "3" (short), "4" (official)
    Lang  string // Language code (e.g. "en", "ja", "x-jat")
    Title string // The title text
}

// TitleMatch is a search result from the title dump with a match quality score.
type TitleMatch struct {
    AID   int     // Anime ID
    Type  string  // Title type
    Lang  string  // Language code
    Title string  // Matched title
    Score float64 // Match quality: 1.0 = exact, 0.75 = prefix, 0.5 = contains
}
```

```go
// metadata/anime/anidb/anidb.go (or a new titles.go file)

// LoadTitleDump downloads and parses the AniDB anime titles dump.
// The dump is cached in memory. Call this before using SearchByTitle.
// The dump URL is http://anidb.net/api/anime-titles.dat.gz.
func (c *Client) LoadTitleDump(ctx context.Context) error

// SearchByTitle searches the loaded title dump for anime matching the query.
// Results are ranked: exact match (1.0) > prefix match (0.75) > contains match (0.5).
// Returns up to limit results sorted by score descending.
// LoadTitleDump must be called first; returns an error if the dump is not loaded.
func (c *Client) SearchByTitle(query string, limit int) ([]TitleMatch, error)

// ParseTitleDump parses an AniDB anime titles dump from the given reader.
// This is exported for consumers who want to provide their own dump source.
func ParseTitleDump(r io.Reader) ([]TitleEntry, error)
```

### Impact

- **Eliminates**: `TitleDumpEntry` + download/parse/search logic (~80 LOC)
- **LOC saved**: ~80

---

## P6 — Letterboxd: OAuth2 Client Credentials

### Problem

The Letterboxd API requires OAuth2 `client_credentials` grant for all API access. goenvoy's Letterboxd client expects a pre-supplied access token:

```go
// Current goenvoy (v1.1.0)
func WithAccessToken(token string) Option {
    return func(c *Client) { c.accessToken = token }
}
```

This forces consumers to implement their own OAuth2 token lifecycle management.

### Current Workaround in revenge

revenge implements a full token manager:

```go
type Client struct {
    api       *letterboxdapi.Client
    apiKey    string
    apiSecret string
    mu        sync.Mutex
    tokenExp  time.Time
    // ...
}

func (c *Client) ensureAPI(ctx context.Context) (*letterboxdapi.Client, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if time.Now().Before(c.tokenExp) {
        return c.api, nil
    }
    // Re-authenticate
    token, expiresIn, err := c.authenticate(ctx)
    // Update c.api with new token
    // Set c.tokenExp = time.Now().Add(expiresIn - 30s buffer)
}

func (c *Client) authenticate(ctx context.Context) (string, int, error) {
    form := url.Values{
        "grant_type":    {"client_credentials"},
        "client_id":     {c.apiKey},
        "client_secret": {c.apiSecret},
    }
    // POST to https://api.letterboxd.com/api/v0/auth/token
    // Parse tokenResponse{AccessToken, ExpiresIn, TokenType}
}
```

Every API method must call `ensureAPI()` first:

```go
func (c *Client) SearchFilms(ctx context.Context, query string) (*letterboxdapi.SearchResponse, error) {
    api, err := c.ensureAPI(ctx)
    if err != nil { return nil, err }
    return api.SearchFilms(ctx, query, 1, 20)
}
```

### Proposed Change

Add a `WithClientCredentials` option to the goenvoy Letterboxd module:

```go
// metadata/video/letterboxd/letterboxd.go

// WithClientCredentials configures the client to automatically acquire and
// refresh an OAuth2 access token using the client_credentials grant.
// The token is refreshed 30 seconds before expiry with mutex-protected state.
func WithClientCredentials(clientID, clientSecret string) Option

// TokenCallback is called whenever a new token is acquired or refreshed.
// This allows consumers to persist or log token events.
type TokenCallback func(accessToken string, expiresIn int)

// WithTokenCallback sets a callback for token lifecycle events.
func WithTokenCallback(cb TokenCallback) Option
```

The client internally manages auth token lifecycle:
- Acquires token on first API call
- Stores expiry with 30-second safety buffer
- Re-acquires before expiry
- Thread-safe via mutex

### Impact

- **Eliminates**: `authenticate()`, `ensureAPI()`, mutex-guarded token state + auth response types (~50 LOC)
- **LOC saved**: ~50

---

## Summary

| Priority | Module | Change | Types Eliminated | LOC Saved |
|:--------:|--------|--------|:----------------:|:---------:|
| **P1** | TMDb | `append_to_response` via `RequestOption` + `*Full` types | 5 | ~120 |
| **P2a** | TMDb | Add `Stills` field to `Images` | 1 | ~10 |
| **P2b** | TMDb | Add 4 fields to `ExternalIDs` | 1 | ~10 |
| **P2c** | TMDb | Add `KnownFor` to `PersonResult` + `KnownForItem` type | 3 | ~60 |
| **P3** | Kitsu | `GetAnimeMappings()` + `GetAnimeCategories()` + `Mapping` type | 3 | ~60 |
| **P4** | AniList | `GetMediaDetailed()` + connection types | 13 | ~150 |
| **P5** | AniDB | `LoadTitleDump()` + `SearchByTitle()` + title types | 1 | ~80 |
| **P6** | Letterboxd | `WithClientCredentials()` OAuth2 auto-management | 0 | ~50 |
| | | **Total** | **~27** | **~540** |

### Suggested Milestone Grouping

**v1.2.0 — Type completeness** (non-breaking, additive):
- P2a (Images.Stills)
- P2b (ExternalIDs fields)
- P2c (PersonResult.KnownFor)

**v1.3.0 — append_to_response + new endpoints**:
- P1 (TMDb append_to_response)
- P3 (Kitsu mappings + per-anime categories)

**v1.4.0 — Extended queries + search**:
- P4 (AniList MediaDetailed)
- P5 (AniDB title dump search)

**v1.5.0 — Auth improvements**:
- P6 (Letterboxd OAuth2 client_credentials)
