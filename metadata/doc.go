// Package metadata provides shared types for interacting with media metadata
// provider APIs.
//
// Providers are organized by media category:
//
//   - video: TMDb, TheTVDB, Fanart.tv, OMDb, TVmaze, OpenSubtitles, Letterboxd
//   - anime: AniList, AniDB, Kitsu, MyAnimeList
//   - tracking: Trakt, Simkl
//   - adult: StashBox, TPDB
//   - music: MusicBrainz, AudioDB, Deezer, Discogs, Last.fm, ListenBrainz, Spotify
//   - book: Google Books, Open Library
//   - game: IGDB, RAWG, Steam
//
// Individual provider packages build on these shared types to offer
// provider-specific clients.
//
// # Shared Types
//
// [Rating], [ExternalID], [Image], [Person], and [SearchResult] are common
// across most metadata providers and allow consumers to work with normalized
// data regardless of the source.
package metadata
