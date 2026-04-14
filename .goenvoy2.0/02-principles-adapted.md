# 02 · Adapted principles — the §2 contract for goenvoy

This is goenvoy's **foundational contract**, the equivalent of golusoris's `PLAN.md §2` but re-cut for a pure-stdlib multi-module client-lib collection. Source this file into `AGENTS.md` as an authoritative reference. Any deviation requires a PR comment with justification.

Structure mirrors golusoris §2 so reviewers can move between them without re-orienting.

---

## 2.1 Coding rules — Power of 10, Go-adapted

| # | Original rule | goenvoy adaptation |
|---|---|---|
| 1 | Simple control flow; no `goto`, recursion | No `goto`. Hand-written recursion only where it mirrors a naturally-recursive API payload (e.g. folder trees in Plex/Jellyfin) — document the depth bound. `panic/recover` allowed only inside `http.RoundTripper` shims or user-supplied JSON decoders. |
| 2 | Bounded loops | Every retry/backoff loop has a max-attempt counter visible in the loop head. Long-running streaming readers (rare — SSE/WebSocket) select on `ctx.Done()`. |
| 3 | No dynamic alloc post-init | **Guidance.** Preallocate slices when the result size is known (`make([]T, 0, n)`). Not a hard gate. |
| 4 | ≤60-line functions | **Hard gate.** `funlen` 120 lines / 60 stmts, `gocognit` 30, `cyclop` 18. Refactor; don't `//nolint`. |
| 5 | ≥2 assertions per function | **Applies to tests.** Table-driven cases with ≥2 asserts per case. Non-test code asserts at boundaries (URL validation in `New`, status-code check in every method). |
| 6 | Smallest scope | No package-level mutable state. Unexport any field not part of the documented API. Sentinel errors (`ErrNotFound`, `ErrUnauthorized`) are the only permitted package-level `var`s. |
| 7 | Check every return | **Hard gate.** `errcheck` + `wrapcheck` + `nilerr`. Every error crossing a module boundary wraps with `fmt.Errorf("<module>: <op>: %w", err)`. |
| 8 | Preprocessor | N/A in Go. |
| 9 | Pointer restrictions | **Guidance.** No multi-hop deref chains; small interfaces (≤5 methods) defined at the consumer. No `unsafe`. |
| 10 | Pedantic warnings | **Hard gate.** 0 lint · 0 gosec · 0 govulncheck · race-green — **per module**, across the CI matrix. `//nolint` requires a same-line justification comment. |

**Hard gates**: 1, 2, 4, 7, 10. **Guidance** (cite in review): 3, 5, 6, 9. Rule 8 is N/A.

---

## 2.2 Secure coding — SEI CERT for Go

Adopted. Particular vigilance for a client-lib collection:

- **IDS00-Go / IDS01-Go** — Sanitize inputs used to build URLs. Every module validates `baseURL` in `New`, never concatenates user input into paths without URL-escaping.
- **FIO06-Go** — Never log secrets. API keys must not appear in error messages, wrapped errors, or test fixtures.
- **ENV33-Go** — Secrets come from parameters, not env. goenvoy never reads `os.Getenv` inside a client.
- **MSC03-Go** — No hard-coded credentials, even in tests (use `httptest` with throwaway tokens).
- **MEM30-Go / MEM31-Go** — Close every `*http.Response.Body`. `bodyclose` enforces.
- **STR00-Go** — Validate URL schemes; reject anything that isn't `http`/`https`.

---

## 2.3 Go style — Google Go Style Guide (canonical)

<https://google.github.io/styleguide/go/>. Secondary: Effective Go + Go Code Review Comments. Commitments specific to goenvoy:

- **Package naming** — module name and package name align (`arr/sonarr` → `package sonarr`).
- **Client factory** — every service exposes `New(baseURL, apiKey string, opts ...Option) (*Client, error)`. Exceptions (OAuth2 device-code flow, no-auth public APIs) go in the module's `AGENTS.md`.
- **Method naming** — verbs for actions (`GetSeries`, `CreateTag`), `List*` or `GetAll*` for pagination-free listings, `Search*` for query-backed listings. Match the upstream API's terminology where it doesn't clash with Go convention.
- **Doc comments** — full sentences ending with a period (`godot`). Exported types + functions + vars + constants all documented.
- **No stuttering** — `sonarr.Client`, not `sonarr.SonarrClient`.

---

## 2.4 Architecture decisions — ADRs (no C4)

C4 diagrams are omitted — each module is a single-file HTTP client with no sub-system architecture to diagram.

**ADRs** live in `docs/adr/`, Nygard format, one decision per file. Retroactively backfilled records capture every convention that's currently implicit:

| ADR | Title | Status |
|---|---|---|
| 0001 | Pure stdlib — no external dependencies | Accepted (retroactive) |
| 0002 | Per-module independent semver via `<path>/vX.Y.Z` tags | Accepted (retroactive) |
| 0003 | `New(baseURL, apiKey) → *Client` constructor shape | Accepted (retroactive) |
| 0004 | Functional-options configuration (`Option` + `With*`) | Accepted (retroactive) |
| 0005 | Per-module `APIError` struct implementing `error` | Accepted (retroactive) |
| 0006 | `httptest`-only test policy — no live API calls | Accepted (retroactive) |
| 0007 | OAuth2 flow helpers — device / auth-code / PKCE / refresh, per-module | Accepted (retroactive) |
| 0008 | XML-RPC / JSON-RPC / GraphQL clients reuse the same stdlib patterns as REST | Accepted (retroactive) |
| 0100 | *(next forward-looking ADR starts here)* | — |

See [08-adrs-and-architecture.md](08-adrs-and-architecture.md) for full drafts.

---

## 2.5 Security + supply-chain standards

Adopted in tiers. goenvoy is a **library**, so compliance claims belong to downstream apps. What the library *does* commit to:

| Standard | Commitment |
|---|---|
| **SLSA Level 3** | Every module tag produces cosign-signed checksums, an SPDX SBOM (syft), and actions/attest-build-provenance attestation. Workflow runs on hash-pinned GH Actions. |
| **OSSF Scorecard** | Weekly + on push-to-main. Aim ≥ 7/10; publish the badge on `README.md`. |
| **NIST SSDF (SP 800-218)** | Covered via Scorecard + CodeQL + Dependabot. |
| **OWASP ASVS L2** | Advisory only — not directly applicable to a consumer library. Sections that apply (V6 stored cryptography, V9 communication): documented in `SECURITY.md`. |
| **EU CRA / NIS2 / GDPR / BSI** | Advisory only. The library does not process personal data; TLS verification is always on by default; API keys are never logged. |

Source-tree hygiene commitments:

- **No telemetry.** Zero outbound connections except to the service the caller points us at.
- **No secret persistence.** API keys live in the `*Client` struct only; never cached to disk, never written to logs.
- **No weak TLS defaults.** Users can override `*http.Client` via `WithHTTPClient`, but the default uses `http.DefaultTransport`'s TLS config.

---

## 2.6 Wire protocols + API standards

Most of golusoris's §2.6 is server-side. goenvoy's subset:

| Standard | Where it applies |
|---|---|
| **RFC 7230/7231 HTTP/1.1** + **RFC 9110 HTTP Semantics** | Every module respects status-code semantics (4xx caller-error, 5xx retry-candidate). 429 triggers backoff. 401/403 do not retry. |
| **OAuth 2.0** (RFC 6749/6750) + **PKCE** (RFC 7636) + **Device Auth Grant** (RFC 8628) | Applied per-service: Trakt, Spotify, Letterboxd, TMDb v4, MAL, AniList, Discogs, Deezer, ListenBrainz. |
| **JWT** (RFC 7519) | Applied where service requires (TheTVDB). |
| **RFC 6265 Cookies** | Applied for session-based clients (some media servers). Cookie jar isolation is required (each `*Client` gets its own `cookiejar.Jar`). |
| **ETag / Cache-Control** (RFC 7232/7234) | Applied as opt-in headers where the upstream supports it (Trakt, TheTVDB) — document in module `AGENTS.md`. |

---

## 2.7 Tooling + formatting

| Standard | Enforcement |
|---|---|
| **EditorConfig** | `.editorconfig` at repo root. |
| **gofumpt** | Via golangci-lint `formatters:` + Claude PostToolUse hook. |
| **gci** | Grouped imports: `standard`, `default`, `prefix(github.com/golusoris/goenvoy)`. |
| **golines** | 120-col line cap. Enforced as warning first, then gate after grace period. |
| **Conventional Commits 1.0** | PR-title check in `ci.yml`. release-please reads history. |
| **Semantic Versioning 2.0** | Per-module tags `<path>/vX.Y.Z`. |
| **Keep a Changelog 1.1** | Root `CHANGELOG.md` — hand-maintained per-module subsections today; release-please can take over per-module in phase 4 (see [06](06-release-and-versioning.md)). |
| **Trunk-Based Development** | Single `main` branch. |

---

## 2.8 Testing standards

| Practice | When it applies |
|---|---|
| **Table-driven tests** | Every function with ≥2 distinct input/output pairs. |
| **`go test -race -count=1`** | Every CI run, per module. |
| **`httptest.NewServer`** | Any HTTP interaction. **No live API calls in tests, ever.** |
| **Example tests** (`example_test.go`) | Every client exposes ≥1 runnable godoc example. |
| **Fuzz tests** | Encouraged on JSON decoders of services with untrusted responses (aggregators like FlareSolverr, Hasheous). |
| **Coverage target** | 70% per module in CI. OAuth flow helpers and error-decoding paths: 85%. |

---

## 2.9 Deployment + configuration

**N/A.** goenvoy ships no binaries, no containers, no manifests. The only release artefact is the tagged module source plus supply-chain attestations (see §2.5 + [06](06-release-and-versioning.md)).

---

## 2.10 Reference links (machine-resolvable)

- Power of 10: <https://web.eecs.umich.edu/~imarkov/10rules.pdf>
- SEI CERT for Go: <https://wiki.sei.cmu.edu/confluence/display/go/>
- Google Go Style Guide: <https://google.github.io/styleguide/go/>
- Nygard ADR format: <https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions>
- Conventional Commits: <https://www.conventionalcommits.org/>
- Semantic Versioning: <https://semver.org/>
- Keep a Changelog: <https://keepachangelog.com/en/1.1.0/>
- SLSA: <https://slsa.dev/>
- OSSF Scorecard: <https://github.com/ossf/scorecard>
- gofumpt: <https://github.com/mvdan/gofumpt>
- gci: <https://github.com/daixiang0/gci>
- golines: <https://github.com/segmentio/golines>
- cosign: <https://docs.sigstore.dev/cosign/>
- syft: <https://github.com/anchore/syft>
- apidiff: <https://pkg.go.dev/golang.org/x/exp/cmd/apidiff>
