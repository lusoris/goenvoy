# ADR-0005 — Per-module `APIError` struct

* Status: Accepted
* Date: 2026-04-14 (retroactive)
* Deciders: @kilian

## Context

Services return non-2xx with varying error shapes: some return a
problem-like body, some plain text, some HTML. Callers need a
type-switchable way to distinguish "the service said no" from "the
network said no".

## Decision

Each module defines its own `APIError` struct with at least
`StatusCode int`, `Message string`, and a raw `Body string`. It
implements `error` and is returned wrapped via `fmt.Errorf("<mod>: ... %w", apiErr)`.

Callers use `errors.As(err, &sonarr.APIError{})` to distinguish.

## Consequences

Positive:

- Callers can branch on HTTP semantics (429 → backoff, 401 → reauth)
  without string-matching.
- Per-module shape lets each service expose its own structured
  fields where the upstream has them.

Negative:

- Each module carries an `APIError` definition — some duplication,
  acceptable for the separation.

## Alternatives considered

- Single shared `APIError` in a `goenvoy/errors` package — rejected:
  introduces a coupling between modules that the per-module-semver
  decision (ADR-0002) deliberately avoids.
- Sentinel errors only (`ErrNotFound`, `ErrUnauthorized`) — rejected:
  loses status code + body context callers need for retries.
