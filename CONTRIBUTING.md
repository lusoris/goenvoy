# Contributing

## Conventional commits

All commits and PR titles MUST follow [Conventional Commits](https://www.conventionalcommits.org/):

```text
<type>(<scope>): <subject>

[optional body]

[optional footer(s)]
```

**Types**: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `perf`, `build`, `ci`, `revert`.

**Scopes** are module paths: `feat(arr/sonarr):`, `fix(metadata/video/tmdb):`, `chore(tools):`.

## Breaking changes

Append `!` to type and add a `BREAKING CHANGE:` footer plus a `Migration:` block with before/after Go snippets:

```text
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
# Set up workspace (links all 69 modules)
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
