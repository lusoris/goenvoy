# goenvoy ← golusoris standards alignment

> **Scope of this working directory.** Plans only. No code changes live here.
> Target repo: [`github.com/golusoris/goenvoy`](https://github.com/golusoris/goenvoy) (local path `/home/kilian/dev/goenvoy`).
> Authored: 2026-04-14. Source standards: golusoris `main` at 701d4e4.

## Problem

`golusoris` (framework) just cut v0.1.0 and hardened its governance, CI, release, and AI-assist layers. `goenvoy` (API-client library collection) is functionally fine but hasn't inherited the same baseline — no golangci-lint, no CodeQL, no Scorecard, no SLSA, no apidiff, no ADRs, no `AGENTS.md`/`CLAUDE.md`, no Claude hooks/skills, no SBOM, no conventional-commits gate, ad-hoc contributing docs.

The ask: **raise goenvoy to golusoris standards without importing the framework itself**. goenvoy must stay a pure-stdlib multi-module monorepo. It is a **consumer** of HTTP APIs, not a server framework — fx, ogen, sqlc, river, migrations, clock, log/slog conventions, etc. do **not** apply. Governance, CI rigor, supply-chain hardening, agent tooling, and document conventions **do** apply.

## Reading order

| # | Doc | What it covers |
|---|---|---|
| 00 | [Overview & principles translation](00-overview.md) | Why we're doing this, what transfers, what doesn't. Boundary map between framework and library discipline. |
| 01 | [Gap analysis](01-gap-analysis.md) | Side-by-side of every standards artefact: present / partial / missing. Drives the rest. |
| 02 | [Adapted principles (the §2 contract for a pure-lib monorepo)](02-principles-adapted.md) | Power-of-10, CERT, Google Style, C4+ADR, SLSA/OSSF, SemVer — re-cut for a client-lib collection. |
| 03 | [Governance docs — AGENTS.md / CLAUDE.md / SECURITY / CONTRIBUTING / COC](03-governance-docs.md) | Full drafts + diffs for each root-level markdown file. |
| 04 | [Lint + tooling baseline](04-lint-and-tooling.md) | `.golangci.yml` (30+ linters), `.editorconfig`, `tools/Makefile.shared`, gofumpt/gci/golines. |
| 05 | [CI workflows](05-ci-workflows.md) | `ci.yml` (multi-module matrix with Conventional-Commits, gosec, govulncheck, coverage, apidiff), `codeql.yml`, `scorecard.yml`, `auto-assign.yml`. |
| 06 | [Release + versioning](06-release-and-versioning.md) | Per-module tags (goenvoy already does this) + release-please multi-package config + SLSA L3 provenance + cosign + syft SBOM. |
| 07 | [Claude hooks + skills](07-claude-hooks-and-skills.md) | `.claude/settings.json`, adapted `guard-bash` / `guard-go-edit` / `format-go-write` hooks, goenvoy-specific skills (`/add-service-client`, `/add-service-method`, `/bump-module`, `/release-module`). |
| 08 | [ADRs + architecture](08-adrs-and-architecture.md) | `docs/adr/` Nygard template + retroactive ADRs for every non-obvious decision already in the repo (pure-stdlib, per-module tags, OAuth flow choices, etc.). |
| 09 | [Per-module conventions](09-per-module-conventions.md) | Per-category `AGENTS.md` (arr, metadata, downloadclient, mediaserver, anime) + per-service `AGENTS.md` template. |
| 10 | [Rollout checklist](10-rollout-checklist.md) | Phased, safe-to-land order. Each phase is one PR. Merge-on-green gates. |

## Acceptance criteria (done when)

1. Every merged commit on `main` is 0 lint · 0 gosec · 0 govulncheck · race-green across **all** modules in the matrix.
2. Every PR title passes Conventional-Commits.
3. CodeQL + OSSF Scorecard workflows run weekly and on every push to `main`.
4. Every module tag (`<path>/vX.Y.Z`) ships with cosign-signed checksums + syft SBOM + SLSA-L3 provenance.
5. `AGENTS.md`, `CLAUDE.md`, `SECURITY.md`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md` exist at repo root and match the golusoris tone + depth (adapted).
6. Every category has an `AGENTS.md`; every service has a minimal `AGENTS.md`.
7. `.claude/settings.json` + hooks + skills land; `/add-service-client` can scaffold a new client from a one-line prompt.
8. `docs/adr/` contains the Nygard template + at least ADR-0001…ADR-0006 retroactive records.
9. Breaking changes to any module's public API fail CI unless the commit body carries a `Migration:` block with before/after Go snippets.

## Non-goals

- Do **not** import or depend on `github.com/golusoris/golusoris`. goenvoy's pure-stdlib promise is load-bearing.
- Do **not** add fx / sqlc / ogen / river / slog-bridge / clock tooling. goenvoy is a client lib, not a server framework.
- Do **not** unify module versions. Per-module semver is a feature, not a bug (ADR-0002 in this plan).
- Do **not** change public APIs as part of the governance rollout. Standards adoption is orthogonal to functional work.
