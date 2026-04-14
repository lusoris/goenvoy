# ADR-0008 — JSON-RPC / XML-RPC / GraphQL clients use the same patterns

* Status: Accepted
* Date: 2026-04-14 (retroactive)
* Deciders: @kilian

## Context

Deluge uses JSON-RPC. rTorrent uses XML-RPC. Stash uses GraphQL. NZBGet
uses JSON-RPC. Letterboxd and AniList use REST + GraphQL hybrid.

## Decision

Wire-protocol variety does not unseat the module conventions. Each
client still:

- Exposes `New(...) (*Client, error)` returning a type-safe surface.
- Accepts functional options.
- Returns typed responses and an `APIError` on non-2xx.
- Uses stdlib only (hand-rolled JSON-RPC envelope, `encoding/xml` for
  XML-RPC, raw JSON `{ "query": "..." }` POST for GraphQL).

## Consequences

Positive:

- Consumers learn one shape; protocol is an implementation detail.
- No dependency on `machinebox/graphql` or similar.

Negative:

- The XML-RPC module (`rtorrent`) hand-rolls encoder/decoder — ~150
  lines. Accepted.

## Alternatives considered

- Adopt protocol-specific client libraries (e.g. `machinebox/graphql`,
  `gorilla/rpc`) — rejected: violates ADR-0001 for marginal ergonomic
  gain.
- Generate clients from upstream OpenAPI/GraphQL schemas — rejected
  for now: the upstream schemas are inconsistent in quality, and the
  manual code is small enough to maintain by hand. Re-evaluate per
  service if generated code becomes substantially smaller.
