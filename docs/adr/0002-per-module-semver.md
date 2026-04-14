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

`release.yml` fires on any `**/v[0-9]+.*` push. release-please operates
in multi-package mode with one package per module.

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
