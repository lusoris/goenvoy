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
