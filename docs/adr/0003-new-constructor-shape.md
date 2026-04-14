# ADR-0003 — `New(baseURL, apiKey string, opts ...Option) (*Client, error)`

* Status: Accepted
* Date: 2026-04-14 (retroactive)
* Deciders: @kilian

## Context

Every service module needs a constructor. Options:

1. Struct-literal: `&sonarr.Client{BaseURL: "...", APIKey: "..."}`.
2. Variadic options only: `sonarr.New(opts ...Option)`.
3. Mandatory positional + options: `sonarr.New(baseURL, apiKey string, opts ...Option)`.

## Decision

Option 3. `baseURL` and `apiKey` (or equivalent: token, username+password
pair for basic auth, OAuth token for OAuth-only services) are the only
two positional arguments; everything else is a `With*` option.

Services that don't have a single API key (OAuth2-only, public-no-key,
basic-auth) document their shape in the module's `AGENTS.md` and expose
a parallel `NewWithToken(...)` / `NewPublic(...)` / `NewBasic(...)` as
needed.

## Consequences

Positive:

- One recognisable shape across the whole library.
- Compile-time fail if `baseURL` / `apiKey` omitted.

Negative:

- Services with fundamentally different auth (e.g. OAuth-only) need
  explicit alternates.

## Alternatives considered

- Struct literal — rejected: silently-valid zero values are dangerous.
- Variadic-only — rejected: defers "did you set baseURL" to runtime.
