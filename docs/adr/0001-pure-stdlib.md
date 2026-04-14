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
