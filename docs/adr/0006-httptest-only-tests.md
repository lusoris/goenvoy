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
No module ever opens a connection to a real upstream host during `go test`.
The Claude `guard-go-edit.sh` hook blocks live-host URLs in
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

## Alternatives considered

- Recorded VCR-style tests — rejected: fixture rot is identical, plus
  extra tooling for record/replay; `httptest` covers the same ground
  with stdlib only (ADR-0001).
- Optional live-mode behind an env flag — rejected: tempts CI use,
  inevitably leaks tokens or rate-limits a third party.
