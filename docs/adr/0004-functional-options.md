# ADR-0004 — Functional-options configuration

* Status: Accepted
* Date: 2026-04-14 (retroactive)
* Deciders: @kilian

## Context

Optional client configuration: custom `*http.Client`, timeout,
additional headers, base-URL overrides for proxy, user-agent, ...
Need a pattern that composes without introducing a public config
struct that becomes frozen on first release.

## Decision

Use Rob Pike's functional-options pattern:

```go
type Option func(*Client)
func WithHTTPClient(c *http.Client) Option { ... }
func WithTimeout(d time.Duration) Option  { ... }
```

`New` takes `opts ...Option` and applies each. Adding a new option is
additive — never breaking.

## Consequences

Positive:

- API-stable: adding a new `With*` never breaks the API surface.
- Options self-document.

Negative:

- Slight allocation cost on each call — negligible for a client
  constructor called once.

## Alternatives considered

- Public config struct — rejected: any new field is a breaking change
  if callers use struct-literal initialisation.
- Builder pattern — rejected: heavier API surface; functional options
  cover the same ergonomics with less code.
