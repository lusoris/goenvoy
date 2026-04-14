# 03 · Governance docs — root-level markdown

This doc carries **final drafts** of every root-level governance file. Land order is in [10-rollout-checklist.md](10-rollout-checklist.md); drafts below are ready to copy-paste.

All drafts keep golusoris's tone: terse, first-person-plural where appropriate, imperative voice for rules, full sentences, no marketing.

---

## 3.1 `AGENTS.md` (NEW)

Purpose: cross-tool agent guide (Claude Code, Cursor, Aider, Codex, Continue). Primary source-of-truth for conventions; `CLAUDE.md` extends this with Claude-specific tooling.

```markdown
# Agent guide — goenvoy

> Cross-tool context for Claude Code, Cursor, Aider, Codex, Continue, and other coding assistants.
> **Read this before suggesting changes.** Then read the per-category `AGENTS.md` for the area you're touching.

## What this repo is

`goenvoy` is a **multi-module monorepo** of pure-stdlib Go HTTP-API clients — 63+ services across arr-stack, metadata, download clients, media servers, and anime. Each service is its own Go module at `github.com/lusoris/goenvoy/<category>/<service>`, independently versioned via `<category>/<service>/vX.Y.Z` tags.

## Hard rules

1. **Never add a non-stdlib import.** Pure stdlib is load-bearing (ADR-0001). `net/http`, `encoding/json`, `encoding/xml`, `crypto/*`, `context`, `net/url`, `net/http/httptest` only. `depguard` in `.golangci.yml` enforces. If you *think* you need a dependency, write an ADR first.
2. **Every module exports `New(baseURL, apiKey string, opts ...Option) (*Client, error)`.** Exceptions go in the module's `AGENTS.md` (OAuth2 device-code, public-key-only APIs).
3. **Every I/O method takes `context.Context` as the first argument.** No exceptions.
4. **Every module defines an `APIError` struct implementing `error`.** Status code, message, response body.
5. **Every HTTP call closes the response body.** `bodyclose` enforces.
6. **All errors wrap** — `fmt.Errorf("<module>: <op>: %w", err)`. `wrapcheck` enforces at boundaries.
7. **Never log secrets.** API keys / bearer tokens / refresh tokens must not appear in error messages, wrapped errors, or test fixtures.
8. **Every merged commit: 0 lint · 0 gosec · 0 govulncheck · race-green** — **per module**. `//nolint` requires a same-line justification comment.
9. **Breaking a public API requires a `Migration:` footer** in the commit body with before/after Go snippets. CI `apidiff` runs against the last tag.

See [.workingdir/PRINCIPLES.md](.workingdir/PRINCIPLES.md) (the adapted §2 contract) for the full coding / security / testing rules.

## Repository layout

(copy current README tree here, unchanged)

## Common tasks

| Task | Claude Code skill / command |
|---|---|
| Add a new service client | `/add-service-client <category> <service>` |
| Add a method to an existing client | `/add-service-method <module> <MethodName>` |
| Bump a module's version | `/bump-module <module> <level>` |
| Release a module | `/release-module <module> <version>` |
| Audit + refresh pinned API docs | `/audit-service-docs` |

## CI gates (per module, in matrix)

Every PR runs, per module:

- `golangci-lint run` — 30+ linters, config at `.golangci.yml`.
- `gosec ./...` — SARIF uploaded to code-scanning.
- `govulncheck ./...`.
- `go test -race -count=1 -coverprofile=...` — coverage ≥ 70%.
- `apidiff` vs the previous `<module>/vX.Y.Z` tag — breaking changes require a `BREAKING CHANGE:` / `Migration:` footer.
- `go vet` + `go build`.

And repo-wide:

- PR title matches Conventional Commits 1.0.
- CodeQL (Go, `security-extended,security-and-quality`).
- OSSF Scorecard.

## Pinned upstream-API docs

`docs/upstream/<service>.md` records the canonical API-doc URL, the version this module targets, and the last-verified date. Check these before suggesting a new method — public docs may be ahead or behind.

## When in doubt

Read [.workingdir/PRINCIPLES.md](.workingdir/PRINCIPLES.md) for the full contract, then the per-category `AGENTS.md` for the area you're touching.
```

---

## 3.2 `CLAUDE.md` (NEW)

Claude-specific — pointers into `.claude/` tooling, a short "Don't" list, and a reference to `AGENTS.md`.

```markdown
# Claude Code guide — goenvoy

> Claude Code-specific guide. For cross-tool conventions read [AGENTS.md](AGENTS.md) first; this file extends it.

## Skills available

Located in `.claude/skills/`:

| Skill | When to use |
|---|---|
| `add-service-client` | Scaffold a new API-client module from a one-line prompt. |
| `add-service-method` | Add a new typed method + test case to an existing client. |
| `bump-module` | Bump one module's semver (major/minor/patch) and open a PR. |
| `release-module` | Tag `<module>/vX.Y.Z`, trigger release.yml. |
| `audit-service-docs` | Refresh `docs/upstream/<service>.md` with today's date + current URL. |

Invoke via `/<skill-name>` in Claude Code.

## Hooks active

Located in `.claude/hooks/`:

- **PreToolUse / Bash** — `guard-bash.sh` blocks `--no-verify`, `--no-gpg-sign`, force-push to main/master, `rm -rf .git`, `rm -rf .workingdir*`.
- **PreToolUse / Edit|Write** — `guard-go-edit.sh` blocks: non-stdlib imports (pure-stdlib invariant), `InsecureSkipVerify: true` without a justified `//nolint:gosec`, unjustified `//nolint` directives, live-API URLs in `*_test.go`.
- **PostToolUse / Edit|Write** — `format-go-write.sh` runs `gofumpt -w` + `gci write -s standard -s default -s 'prefix(github.com/lusoris/goenvoy)'`.

## Tone

- Be terse. No preamble.
- When changing a public API: write the `Migration:` footer in the commit body with before/after Go snippets.
- When adding a method: always add a table-driven test case + runnable godoc example.
- Never suggest adding a dependency — goenvoy is pure stdlib by ADR-0001.

## Project principles

Read [.workingdir/PRINCIPLES.md](.workingdir/PRINCIPLES.md) first. Quick hitlist:

- Pure stdlib. No imports outside `net/http`, `encoding/json`, `encoding/xml`, `crypto/*`, `context`, `net/url`, `net/http/httptest`.
- `New(baseURL, apiKey) → *Client` constructor shape.
- Functional options (`Option` + `With*`).
- Every method takes `context.Context` first.
- Every module has an `APIError`.
- Every response body is `defer resp.Body.Close()`-ed.
- `//nolint` needs a same-line justification.

## Don't

- Don't add external dependencies. Ever.
- Don't skip response-body close.
- Don't concatenate user input into URL paths without `url.PathEscape`.
- Don't silence a linter without a justification comment.
- Don't write multi-paragraph comments. One-line `// Why:` comments only.
- Don't create new markdown docs unless explicitly asked.
- Don't run tests against live APIs — `httptest` only.

## Per-commit doc-sync

When touching a module:

- Update `<module>/AGENTS.md` if conventions changed.
- Update `CHANGELOG.md` under the module's unreleased section.
- Update `docs/upstream/<service>.md` if the upstream API surface moved.
```

---

## 3.3 `CONTRIBUTING.md` (REWRITE)

Replaces the current "basic" file. Adds: Conventional Commits, Migration: footer, pre-commit hooks, link to principles.

```markdown
# Contributing

## Conventional commits

All commits and PR titles MUST follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

[optional body]

[optional footer(s)]
```

**Types**: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `perf`, `build`, `ci`, `revert`.

**Scopes** are module paths: `feat(arr/sonarr):`, `fix(metadata/video/tmdb):`, `chore(tools):`.

## Breaking changes

Append `!` to type and add a `BREAKING CHANGE:` footer plus a `Migration:` block with before/after Go snippets:

```
feat(arr/sonarr)!: rename GetSeries to ListSeries

BREAKING CHANGE: GetSeries is now ListSeries to match upstream terminology.

Migration:
  // before
  series, err := client.GetSeries(ctx)
  // after
  series, err := client.ListSeries(ctx)
```

The `Migration:` block is **required** for breaking changes. CI (`apidiff`) fails without it.

## CI gates (per module)

Every PR runs per module in a matrix:

- `golangci-lint run` — 30+ linters (see `.golangci.yml`)
- `gosec` — SARIF to code-scanning
- `govulncheck`
- `go test -race -count=1` — ≥70% coverage per module
- `apidiff` vs the previous `<module>/vX.Y.Z` tag

And repo-wide:

- Conventional-Commits PR title
- CodeQL (`security-extended,security-and-quality`)
- OSSF Scorecard

## Local development

```bash
# Set up workspace (links all 63 modules)
go work init && find . -name 'go.mod' -not -path './.workingdir/*' -exec dirname {} \; | xargs go work use

# Full local CI (all modules)
make ci-all

# One module
cd arr/sonarr
make -f ../../Makefile ci
```

See `Makefile` for `test-all`, `lint-all`, `vuln-all`, `gosec-all`, `ci-all`, `fmt-all`, `tidy-all`.

## Adding a new service

Use the Claude Code skill **`/add-service-client`** or follow the manual steps:

1. Create `<category>/<service>/` (e.g. `downloadclient/aria2/`).
2. `go mod init github.com/lusoris/goenvoy/downloadclient/aria2`.
3. Standard files:
   - `doc.go` — package doc comment.
   - `types.go` — request/response types.
   - `<service>.go` — `New` + methods.
   - `<service>_test.go` — `httptest`-based table-driven tests.
   - `example_test.go` — runnable godoc example.
   - `AGENTS.md` — auth model, pagination style, pinned doc URL.
4. Add the module to `go.work`.
5. `cd <category>/<service> && make -f ../../Makefile ci`.

## Code style

- **Pure stdlib.** `depguard` in `.golangci.yml` blocks external imports. See ADR-0001.
- **Functional options.** `Option` type + `With*` constructors.
- **Context propagation.** Every I/O method takes `context.Context` first.
- **Error types.** Each module defines an `APIError` implementing `error`.
- **Tests.** `httptest.NewServer` only — never hit live APIs.
- **Godoc.** All exported identifiers carry a doc comment ending in a period.

## Pre-commit hooks

```bash
pre-commit install
```

Hooks: `gofumpt`, `golangci-lint`, `gitleaks`, conventional-commit check.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
```

---

## 3.4 `SECURITY.md` (REWRITE)

Replaces the current file. Adds: cosign-verify block, SBOM + provenance claims, Dependabot/Renovate note.

```markdown
# Security policy

## Reporting a vulnerability

Please **do not** open public GitHub issues for security vulnerabilities.

Email: `security@lusoris.dev` (or open a private security advisory on GitHub).

We aim to acknowledge within 72 hours and provide a remediation plan within 7 days.

## Supported versions

Only the latest release of each module is patched. Downstream apps should bump promptly when an advisory is published for a module they use.

## Supply chain

Every tagged module release ships:

- Source tarball + `checksums.txt`.
- `checksums.txt.sig` + `.pem` — [cosign](https://docs.sigstore.dev/cosign/) keyless signature (GitHub OIDC).
- SPDX SBOM via [syft](https://github.com/anchore/syft).
- SLSA-L3 provenance via [actions/attest-build-provenance](https://github.com/actions/attest-build-provenance).

Verify a release archive:

```bash
cosign verify-blob \
  --certificate checksums.txt.pem \
  --signature checksums.txt.sig \
  --certificate-identity-regexp '^https://github.com/lusoris/goenvoy/' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  checksums.txt
```

## Security practices

- **Pure stdlib** — zero external dependencies reduces supply-chain risk to the Go toolchain and `golang.org/x/*` (none currently used).
- **No telemetry** — clients open connections only to the service the caller points them at.
- **No secret persistence** — API keys, bearer tokens, and refresh tokens live in the `*Client` struct only. They are never cached to disk, never written to logs, never included in error messages.
- **TLS on by default** — no module disables certificate verification. `WithHTTPClient` lets callers override the `*http.Client` for specialised cases; doing so is their responsibility.
- **URL validation** — every `New` validates the `baseURL` scheme (`http`/`https` only) and parseability.
- **Static analysis** — every module passes [gosec](https://github.com/securego/gosec) + [golangci-lint](https://golangci-lint.run/) (incl. `errorlint`, `noctx`, `bodyclose`, `contextcheck`, `containedctx`) + [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) + CodeQL (`security-extended,security-and-quality`).
- **OSSF Scorecard** — published and monitored; see the badge on `README.md`.

## Dependencies

Tracked by [Dependabot](https://docs.github.com/en/code-security/dependabot):

- `gomod` — pure-stdlib means `gomod` is quiet; updates only touch test-only tooling pulled via `go install`, if any.
- `github-actions` — all workflow actions are hash-pinned; Dependabot proposes version bumps with updated hashes.

## Scope

This library is a collection of HTTP API clients. Relevant security concerns include:

- Credential leakage (API keys / tokens in logs or error messages).
- Request injection (path traversal, header injection).
- TLS verification bypass.
- Improper error handling exposing sensitive data.

We actively test against these in code review + `gosec` rules G101, G107, G306, G402, G505.
```

---

## 3.5 `CODE_OF_CONDUCT.md` (NEW)

Verbatim golusoris, swap email (already `security@lusoris.dev` — no swap needed).

---

## 3.6 `.editorconfig` (NEW)

Verbatim golusoris.

---

## 3.7 `.markdownlintignore` (NEW)

```text
CHANGELOG.md
node_modules/
.workingdir/
.workingdir2/
```

---

## 3.8 `.gitignore` additions

Append to existing:

```gitignore
# Claude Code runtime state (personal, machine-local).
.claude/settings.local.json
.claude/scheduled_tasks.lock

# Coverage artefacts per module
**/coverage.out
**/coverage.html

# Tooling caches
.golangci-lint-cache/
tmp/
.tmp/
```

---

## 3.9 `README.md` additions (no rewrite)

Append to the top, below the existing title:

- OSSF Scorecard badge: `[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/lusoris/goenvoy/badge)](https://scorecard.dev/viewer/?uri=github.com/lusoris/goenvoy)`
- Go Reference badge: `[![Go Reference](https://pkg.go.dev/badge/github.com/lusoris/goenvoy.svg)](https://pkg.go.dev/github.com/lusoris/goenvoy)`
- CodeQL badge: `[![CodeQL](https://github.com/lusoris/goenvoy/actions/workflows/codeql.yml/badge.svg)](https://github.com/lusoris/goenvoy/actions/workflows/codeql.yml)`

And a short "Standards" section near the end:

```markdown
## Standards

- Pure stdlib (ADR-0001). `depguard` in CI enforces.
- Per-module independent semver (ADR-0002). Tags: `<path>/vX.Y.Z`.
- Every release: cosign-signed checksums + SPDX SBOM + SLSA-L3 provenance.
- Every merged commit: 0 lint / 0 gosec / 0 govulncheck / race-green across all modules.
- AGENTS.md / CLAUDE.md for agent-assisted development conventions.
- docs/adr/ for all architectural decisions (Nygard format).
```
