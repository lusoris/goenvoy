# ADR-0007 — OAuth2 flow helpers — device / auth-code / PKCE / refresh

* Status: Accepted
* Date: 2026-04-14 (retroactive)
* Deciders: @kilian

## Context

Several services (Trakt, Spotify, Letterboxd, TMDb v4, MAL, AniList,
Discogs, Deezer, ListenBrainz) require OAuth2. The std lib
`golang.org/x/oauth2` is the obvious choice, but pulling it in
violates ADR-0001.

## Decision

Implement the minimal OAuth2 flows each service needs by hand, per
module, using `net/http` + `net/url` + `encoding/json` only. The
flows are: Device Authorization Grant (RFC 8628), Authorization Code
+ PKCE (RFC 7636), Refresh Token (RFC 6749 §6), and Client
Credentials where applicable.

Each OAuth-using module exposes:

- `NewWithToken(baseURL, accessToken string, opts ...Option)` — the
  common construction path for already-authenticated callers.
- Helper functions for flow init: `StartDeviceAuth(ctx, clientID, ...)`,
  `ExchangeCode(ctx, ...)`, `Refresh(ctx, refreshToken)`.

## Consequences

Positive:

- Preserves ADR-0001.
- Per-service helpers can match the service's idiosyncrasies
  (Spotify's `show_dialog`, Trakt's 6-digit user code).

Negative:

- Modest code duplication across OAuth2-using modules (~120 lines
  each). Accepted trade.

## Alternatives considered

- Adopt `golang.org/x/oauth2` — rejected: violates ADR-0001 and pulls
  in transitive deps that defeat the pure-stdlib promise.
- Single shared `goenvoy/oauth2` helper module — rejected: every
  consumer would pull a second module just to use one client; the
  per-module-semver decision (ADR-0002) makes the shared module a
  cross-cutting versioning headache.
