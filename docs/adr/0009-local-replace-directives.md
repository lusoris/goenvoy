# ADR-0009 — Drop committed `replace ../` directives in favour of `go.work`

* Status: Accepted
* Date: 2026-04-14
* Deciders: @lusoris

## Context

Eight modules in the `arr/` family commit `replace github.com/golusoris/goenvoy/arr => ../`
plus `require github.com/golusoris/goenvoy/arr v0.0.0` in their `go.mod` files
(`arr/sonarr`, `arr/radarr`, `arr/lidarr`, `arr/readarr`, `arr/prowlarr`,
`arr/bazarr`, `arr/whisparr`, `arr/seerr`).

The pattern was introduced when `arr/` had no published tag — `v0.0.0`
was the placeholder. The shared `arr/` base has since been tagged
independently (`arr/v1.0.0`, `v1.1.0`, `v1.2.0`, `v1.2.1`), so the
`v0.0.0` placeholder is now stale. The published state is:

- `arr/sonarr/v1.2.0` declares `require arr v0.0.0` and `replace arr => ../`.
- Downstream consumers (revenge, subdo, arca) explicitly require the real
  `arr v1.2.x` themselves. Go's Minimum Version Selection (MVS) picks
  `v1.2.x` over the unpublished `v0.0.0`, and the build succeeds — only
  because some other dep brings in the real version.
- The `replace ../` directive in the published module is **ignored** by
  consumers — Go only honours `replace` from the main module.

Failure mode: a downstream that imports only `arr/sonarr` (no other
`arr/*` child, no direct `arr` require) hits `arr v0.0.0` and fails to
build. Today this is invisible because every existing downstream uses
multiple `arr/*` modules. It will not stay invisible.

The `gomoddirectives` linter (phase 3 of the goenvoy-2.0 rollout) flags
the local `replace`. Two responses: (a) configure the linter to allow
local `replace`, or (b) remove the directives and update the requires
to a real version. Path (a) papers over a real architectural defect;
path (b) fixes it.

`go.work` already exists at the repo root and lists every module, so
local cross-module development continues to work without the committed
`replace`.

## Decision

Remove the `replace github.com/golusoris/goenvoy/arr => ../` line from
every `arr/*` child `go.mod`, and replace `require arr v0.0.0` with
`require arr v1.2.1` (the latest tag). `go.work` carries the local
linking from now on.

The lint config does **not** add `gomoddirectives.replace-local: true`.
Local replaces remain a violation, so the same defect cannot reappear
silently.

The eight affected child modules need a semver patch bump (`v1.2.x →
v1.2.(x+1)`) when next released — the cleanup fixes a declared-but-
broken require, which is functionally a bug fix for downstream
consumers who hadn't independently required `arr`. The bump is not
landed in this ADR's PR; it falls out of the next regular release of
each module.

## Consequences

Positive:

- Downstream consumers can `go get` any `arr/*` child standalone and
  the build resolves correctly without an explicit `arr` require.
- The lint baseline stays strict — no exception knob to explain.
- One less goenvoy-specific quirk for future contributors to learn.
- `go.work`'s role as the single source of truth for local linking is
  reinforced (matches what `CONTRIBUTING.md` already documents).

Negative:

- Anyone who clones the repo and runs `go test ./...` inside a single
  `arr/*` child without first running `go work init` (or the documented
  `go work use ...` recipe in `CONTRIBUTING.md`) will fetch `arr v1.2.1`
  from the proxy instead of using the local copy. Acceptable: the docs
  cover the `go.work` workflow, and CI uses `make ci-all` from the repo
  root which has `go.work` active.

Neutral:

- The next release of each affected child carries this fix as a patch
  semver bump in its CHANGELOG entry. No coordinated bump required —
  each child ships on its own cadence.

## Alternatives considered

- **Ratify the status quo with `gomoddirectives.replace-local: true`.**
  Rejected: silences a linter rule to mask a published-state defect
  that already has a fix path. Adds a goenvoy-specific exception that
  future contributors must understand.
- **Pin the placeholder require to a pseudo-version of the latest
  commit instead of `v1.2.1`.** Rejected: pseudo-versions in committed
  go.mod files are a code smell of their own — they bind the module to
  an unstable identifier instead of a release tag.
- **Collapse `arr/*` children into the single `arr` module.** Rejected:
  violates ADR-0002 (per-module independent semver) — consumers would
  inherit churn from twelve siblings.
