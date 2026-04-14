# 08 · ADRs + architecture

Target: `docs/adr/` in goenvoy. No `docs/architecture/` (no C4 — client libs have no deployable architecture to diagram). Nygard format, mirrored from golusoris.

Land in phase 3 alongside CI hardening.

---

## 8.1 `docs/adr/README.md`

```markdown
# Architecture Decision Records

One file per decision. Nygard format. Once Accepted, an ADR is immutable —
a new ADR supersedes it (with a `Supersedes: ADR-NNNN` header in both directions).

## Numbering

- `0001`–`0099` — retroactive backfills of decisions already embodied by the codebase.
- `0100`+ — forward-looking decisions made after this log was instituted.

## Index

| ID | Title | Status |
|---|---|---|
| [0001](0001-pure-stdlib.md) | Pure stdlib — no external dependencies | Accepted |
| [0002](0002-per-module-semver.md) | Per-module independent semver via `<path>/vX.Y.Z` tags | Accepted |
| [0003](0003-new-constructor-shape.md) | `New(baseURL, apiKey) → *Client` constructor shape | Accepted |
| [0004](0004-functional-options.md) | Functional-options configuration (`Option` + `With*`) | Accepted |
| [0005](0005-api-error-type.md) | Per-module `APIError` struct implementing `error` | Accepted |
| [0006](0006-httptest-only-tests.md) | `httptest`-only test policy — no live API calls | Accepted |
| [0007](0007-oauth2-flows.md) | OAuth2 flow helpers — device / auth-code / PKCE / refresh | Accepted |
| [0008](0008-non-rest-wire-protocols.md) | JSON-RPC / XML-RPC / GraphQL clients reuse REST-client patterns | Accepted |

## Template

See [0000-template.md](0000-template.md).
```

---

## 8.2 `docs/adr/0000-template.md`

Verbatim golusoris's template (Nygard):

```markdown
# ADR-NNNN — <title>

* Status: Proposed | Accepted | Superseded by ADR-NNNN | Deprecated
* Date: YYYY-MM-DD
* Deciders: <names / handles>

## Context

What is the issue motivating this decision? What forces are at play
(business, technical, organisational)?

## Decision

The decision itself, stated concisely. Imperative voice.

## Consequences

Positive, negative, and neutral consequences. What becomes easier /
harder? What follow-on decisions does this imply?

## Alternatives considered

Brief treatment of each alternative and why it was rejected.
```

---

## 8.3 Retroactive ADRs

Full drafts follow. Each is short — the decisions are already well-understood in the code; the ADR just makes them citable.

### ADR-0001 — Pure stdlib

```markdown
# ADR-0001 — Pure stdlib, no external dependencies

* Status: Accepted
* Date: 2026-04-14 (retroactive, decision pre-dates this log)
* Deciders: @kilian (Lusoris)

## Context

goenvoy is a fan-out of 60+ client libraries. Downstream apps already
pull in their own server stacks (golusoris, chi, pgx, etc.) — adding
another dependency graph per service would inflate the transitive-dep
surface significantly.

## Decision

Every module in goenvoy depends on the Go standard library only.
`golang.org/x/*` is allowed on case-by-case review (currently: none).
No test deps (no testify, no gomega). HTTP via `net/http`, JSON via
`encoding/json`, XML via `encoding/xml`, tests via stdlib `testing`
and `net/http/httptest`.

`.golangci.yml` `depguard` rule enforces this at CI time. The
`.claude/hooks/guard-go-edit.sh` pre-write hook enforces at edit time.

## Consequences

Positive:
- Zero transitive supply-chain surface outside the Go toolchain.
- No version-pin conflicts with downstream apps.
- goenvoy modules are trivially embeddable.

Negative:
- Test fixtures are more verbose than with testify.
- Some algorithms (e.g. consistent hashing, rate limiters) must be
  hand-rolled or skipped.

Neutral:
- Future functionality that genuinely requires an external dep
  (e.g. a binary protocol parser) requires a new ADR superseding or
  amending this one.

## Alternatives considered

- Allow "tiny, stable, vendored" deps — rejected: a fuzzy gate with no
  clean CI enforcement.
- Use testify for ergonomic test asserts — rejected: the ergonomic
  win is small and `depguard` becomes inconsistent.
```

### ADR-0002 — Per-module semver

```markdown
# ADR-0002 — Per-module independent semver

* Status: Accepted
* Date: 2026-04-14 (retroactive)
* Deciders: @kilian

## Context

A monorepo of independent client libraries can version either
(a) all modules together at one repo-wide version, or
(b) each module independently via `<path>/vX.Y.Z` tags (Go's native
multi-module monorepo convention).

Option (a) forces every consumer to bump to the latest repo version
even if only one unrelated module changed. Option (b) preserves
independence at the cost of a longer tag set.

## Decision

Use option (b). Every module is tagged as `<category>/<service>/vX.Y.Z`
(e.g. `arr/sonarr/v1.3.0`). The repo root has no `vX.Y.Z` tag.

`release.yml` fires on any `**/v[0-9]+.*` push. release-please (phase
4) operates in multi-package mode with one package per module.

## Consequences

Positive:
- Consumers pin one module per service, unaffected by churn in other
  services.
- Breaking changes in one module don't force a major bump in all.
- SemVer remains meaningful per surface area.

Negative:
- Tag list is long (hundreds across the lifetime of the repo).
- CI matrix fans out to every module for apidiff.

## Alternatives considered

- Monorepo-wide version — rejected as above.
- Independent repositories per service — rejected: 60+ repos is hostile
  to shared CI + shared conventions + shared AGENTS.md.
```

### ADR-0003 — Constructor shape

```markdown
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
```

### ADR-0004 — Functional options

```markdown
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
```

### ADR-0005 — APIError shape

```markdown
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
implements `error` and is returned wrapped via `fmt.Errorf("<mod>: ...
%w", apiErr)`.

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
```

### ADR-0006 — httptest-only tests

```markdown
# ADR-0006 — `httptest`-only tests — no live API calls

* Status: Accepted
* Date: 2026-04-14 (retroactive)
* Deciders: @kilian

## Context

Client libraries can be tested against (a) a recorded/mocked HTTP
server or (b) live upstream APIs. Live tests produce false-reds on
flaky networks, require secrets in CI, and leak requests to third
parties.

## Decision

All tests use `httptest.NewServer` with hand-crafted response bodies.
No module ever opens a connection to a real upstream host during `go
test`. The Claude `guard-go-edit.sh` hook blocks live-host URLs in
`*_test.go` files for the most-commonly-reached APIs; CI-side there's
no network egress guarantee — convention is the gate.

## Consequences

Positive:
- Deterministic CI.
- No API-key management in CI secrets.
- No rate-limit entanglement.

Negative:
- Fixture bodies need refreshing when upstreams evolve — caught via
  the `docs/upstream/<service>.md` pin + `/audit-service-docs` skill.
```

### ADR-0007 — OAuth2 flows

```markdown
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
- Helper functions for flow init: `StartDeviceAuth(ctx, clientID,
  ...)`, `ExchangeCode(ctx, ...)`, `Refresh(ctx, refreshToken)`.

## Consequences

Positive:
- Preserves ADR-0001.
- Per-service helpers can match the service's idiosyncrasies
  (Spotify's `show_dialog`, Trakt's 6-digit user code).

Negative:
- Modest code duplication across OAuth2-using modules (~120 lines
  each). Accepted trade.
```

### ADR-0008 — Non-REST wire protocols

```markdown
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
```

---

## 8.4 Forward-looking ADRs (slots reserved)

No forward ADRs drafted in this plan. When a new decision lands (e.g. "adopt OpenTelemetry span emission in every HTTP client" or "add Renovate alongside Dependabot"), it takes ADR-0100 upward.

---

## 8.5 How ADRs interact with CLAUDE.md / AGENTS.md

- `AGENTS.md` and per-module `AGENTS.md` cite ADRs by number when explaining a convention ("pure stdlib — see ADR-0001"). This keeps conventions thin and decisions durable.
- When a PR violates a convention, reviewer comments cite the ADR instead of restating the reasoning.
