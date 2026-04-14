# 00 · Overview & principles translation

## What goenvoy is

- A **collection of 63+ HTTP-API client libraries** for media automation, metadata, downloaders, media servers, and anime services.
- A **multi-module Go monorepo** — each service is an independently-versioned module (e.g. `github.com/golusoris/goenvoy/arr/sonarr`).
- **Pure stdlib.** The only imports are `net/http`, `encoding/json`, `context`, `crypto/*`, `net/url`, stdlib test helpers (`net/http/httptest`). No external dependencies, anywhere, ever.
- **Consumer-facing**, not server-facing. Each module exposes `New(baseURL, apiKey, ...Option) (*Client, error)` and a set of typed methods that take `context.Context`.
- Already in use by downstream `lusoris/*` apps (revenge, lurkarr, subdo, arca).

## What golusoris is

- A **single Go module** that wraps a pinned set of best-in-class libraries behind opt-in `fx` modules for building server apps.
- Framework-shaped. Exists to eliminate the "bump one dep in 5 apps" tax.
- Ships Claude/agent tooling (AGENTS.md, CLAUDE.md, hooks, skills, MCP server) as a first-class feature.

These are complementary, not overlapping. goenvoy can (and does) get used *by* golusoris apps via the `integrations/goenvoy/` fx adapter that golusoris ships. That integration direction is fixed — goenvoy does not depend on golusoris.

## Principle translation table

Every standard from golusoris's `.workingdir/PLAN.md §2` maps into goenvoy as **apply / adapt / drop**. Summary first; gap analysis in [01-gap-analysis.md](01-gap-analysis.md) drives the concrete artefacts.

### §2.1 Power of 10 — Go-adapted

| Rule | goenvoy treatment |
|---|---|
| 1 · simple control flow | **Apply.** No `goto`. Recursion only for parser/tree code with a documented bound. `panic/recover` only around user-supplied decoders or inside `http.RoundTripper` recovery shims. |
| 2 · bounded loops | **Apply.** Every retry/backoff loop has an explicit max-attempts bound. Long-running watchers (e.g. SSE/WebSocket tailers, if any) select on `ctx.Done()`. |
| 3 · no dynamic alloc post-init | **Guidance.** Hot paths prealloc; not a gate. |
| 4 · ≤60-line functions | **Apply as CI gate.** `funlen` 120 lines / 60 statements, `gocognit` ≤ 30. |
| 5 · ≥2 assertions | **Apply via tests.** Table-driven tests with ≥2 asserts per case. |
| 6 · smallest scope | **Apply.** No package-level mutable state. |
| 7 · check every return | **Apply.** `errcheck` + `wrapcheck` + `nilerr` on. Every module wraps: `fmt.Errorf("<module>: <op>: %w", err)`. |
| 8 · preprocessor | **N/A.** |
| 9 · pointer restrictions | **Guidance.** Small interfaces, no unsafe. |
| 10 · pedantic warnings | **Apply as CI gate.** 0 lint · 0 gosec · 0 govulncheck · race-green. `//nolint` needs a justification comment. |

### §2.2 SEI CERT for Go — **Apply in full.**

Every module handles API keys; path-injection, TLS bypass, header injection, log leakage are all first-class concerns even though goenvoy doesn't run servers. `gosec` enforces the majority.

### §2.3 Google Go Style — **Apply.**

Canonical. Already largely followed. The diff from today's code is mostly doc-comment-end-with-period (`godot`) and grouped imports (`gci`).

### §2.4 C4 + ADRs — **Apply, adapted.**

- No C4 diagrams needed (goenvoy has no running-system architecture — each module is a one-page client).
- **ADRs in `docs/adr/`** using Nygard format. Retroactive backfill for the decisions already baked in (see [08-adrs-and-architecture.md](08-adrs-and-architecture.md)):
  - ADR-0001: Pure stdlib — no external dependencies.
  - ADR-0002: Per-module independent semver (multi-module monorepo).
  - ADR-0003: `New(baseURL, apiKey) → Client` factory shape.
  - ADR-0004: Functional-options pattern (`Option` type + `With*`).
  - ADR-0005: Per-module `APIError` struct implementing `error`.
  - ADR-0006: `httptest`-only tests — no live API calls.
  - ADR-0007: OAuth2 token acquisition helpers (where applicable) — device code / auth code / PKCE / refresh.

### §2.5 Security + supply-chain — **Apply the client-lib-relevant subset.**

| Standard | Applicability |
|---|---|
| **SLSA Level 3** | **Apply.** Tag-push workflow signs checksums with cosign keyless, emits syft SBOM, attests SLSA-L3 provenance. |
| **OWASP ASVS L2** | **Advisory.** goenvoy runs no servers — not scope. Relevant principles (secret handling, TLS verification, input validation) are enforced via gosec + code review. |
| **NIST SSDF** | **Apply.** OSSF Scorecard covers the majority. |
| **EU CRA / NIS2 / BSI / NCSC / ENISA / GDPR / EU AI Act** | **Advisory only.** goenvoy is a library, not a product with digital elements under CRA. Compliance is a downstream-app concern. Document in `SECURITY.md` that the library does not log secrets, does not disable TLS verification, and emits an SBOM. |

### §2.6 Wire protocols + API standards — **Drop most; apply spot items.**

- RFC 9457 / OpenAPI 3.1 / ogen / spectral / JSON Schema 2020-12 / OTel SemConv — **N/A** (goenvoy consumes, not serves).
- OAuth2 (RFC 6749/6750) + PKCE (RFC 7636) — **Apply** for services that require it (Trakt, Spotify, Letterboxd, TMDb-v4, etc.).
- JWT (RFC 7519) — **Apply where a service uses JWT bearers** (TheTVDB).
- TLS — **Apply.** Never let users disable certificate verification silently; if a module exposes `WithTLSConfig`, document the footgun.

### §2.7 Tooling + formatting — **Apply.**

EditorConfig, gofumpt, gci, golines, Conventional Commits 1.0, SemVer 2.0, Keep-a-Changelog 1.1, Trunk-Based Dev.

One adaptation: `gci` prefix is `github.com/golusoris/goenvoy` (not `github.com/golusoris/golusoris`).

### §2.8 Testing — **Apply.**

Table-driven tests, `go test -race -count=1`, httptest at the HTTP boundary (golusoris prefers integration DBs via testcontainers; goenvoy's analogous choice is httptest for HTTP — already the convention). Coverage gate 70% in CI matrix per module.

### §2.9 Deployment — **Drop.**

goenvoy ships no binaries, no containers, no deploy manifests. Only release artefacts are tagged source + SBOM + provenance attestation.

## Scope boundary rules of thumb

If a proposed artefact…

- …enforces Go code quality or supply-chain hygiene → **apply**.
- …enforces a server-side runtime convention (DI, logging, tracing, DB, jobs, migrations, spec-first HTTP) → **drop**.
- …enforces an AI-assist / agent convention (AGENTS.md, CLAUDE.md, hooks, skills) → **apply, with goenvoy-specific rewiring**.

## What lands in `.workingdir2/` vs in `goenvoy/`

`.workingdir2/` is planning only. The actual artefacts (`.golangci.yml`, `.github/workflows/*.yml`, `.claude/`, `AGENTS.md`, `docs/adr/`, per-module `AGENTS.md`, …) land in the `goenvoy` tree when the rollout in [10-rollout-checklist.md](10-rollout-checklist.md) runs.

The rollout is phased so each step merges green independently — standards adoption should never block functional work.
