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
