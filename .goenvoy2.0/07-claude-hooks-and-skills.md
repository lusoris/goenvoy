# 07 · Claude hooks + skills

Target: `.claude/` directory at the goenvoy repo root. Structure mirrors golusoris.

```
.claude/
├── settings.json              # team-shared, checked in
├── hooks/
│   ├── README.md
│   ├── guard-bash.sh
│   ├── guard-go-edit.sh
│   └── format-go-write.sh
└── skills/
    ├── add-service-client.md
    ├── add-service-method.md
    ├── bump-module.md
    ├── release-module.md
    └── audit-service-docs.md
```

`.gitignore`-ed: `.claude/settings.local.json`, `.claude/scheduled_tasks.lock`.

---

## 7.1 `.claude/settings.json`

Identical structure to golusoris's:

```json
{
  "$schema": "https://json.schemastore.org/claude-code-settings.json",
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [{ "type": "command", "command": ".claude/hooks/guard-bash.sh" }]
      },
      {
        "matcher": "Edit|Write",
        "hooks": [{ "type": "command", "command": ".claude/hooks/guard-go-edit.sh" }]
      }
    ],
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "hooks": [{ "type": "command", "command": ".claude/hooks/format-go-write.sh" }]
      }
    ]
  }
}
```

---

## 7.2 `hooks/guard-bash.sh` (ADAPT from golusoris)

Mostly verbatim golusoris. Additions:

- Block `rm -rf .workingdir2` as well.
- Keep `--no-verify`, `--no-gpg-sign`, force-push-to-main, `reset --hard main/master`, `rm -rf .git`.

```bash
#!/usr/bin/env bash
# PreToolUse hook for Bash. Blocks trivially-wrong commands for this repo.
# Exit 2 = block with rejection reason on stderr.
set -euo pipefail
payload=$(cat)
if command -v jq >/dev/null 2>&1; then
  cmd=$(printf '%s' "$payload" | jq -r '.tool_input.command // ""')
else
  cmd=$(printf '%s' "$payload" | grep -oE '"command":"[^"]*"' | head -n1 | sed 's/^"command":"//; s/"$//')
fi
[ -z "$cmd" ] && exit 0

deny() { printf 'blocked by .claude/hooks/guard-bash.sh: %s\n' "$1" >&2; exit 2; }

case "$cmd" in
  *--no-verify*)   deny "'--no-verify' disables pre-commit/push hooks. Fix the failure instead." ;;
  *--no-gpg-sign*) deny "'--no-gpg-sign' bypasses commit signing. Not allowed." ;;
esac

if printf '%s' "$cmd" | grep -Eq 'git[[:space:]]+push[[:space:]].*--force(-with-lease)?\b.*\b(origin[[:space:]]+)?(main|master)\b'; then
  deny "force-push to main/master blocked. Create a PR or revert-commit instead."
fi

if printf '%s' "$cmd" | grep -Eq 'git[[:space:]]+reset[[:space:]].*--hard[[:space:]]+(origin/)?(main|master)\b'; then
  deny "'git reset --hard main' drops commits. Confirm with the user first."
fi

if printf '%s' "$cmd" | grep -Eq '\brm[[:space:]]+-[A-Za-z]*r[A-Za-z]*f?[A-Za-z]*[[:space:]]+(\./)?\.git(/[^[:space:]]*)?(\s|$)'; then
  deny "'rm -rf .git*' nukes the repo."
fi
if printf '%s' "$cmd" | grep -Eq '\brm[[:space:]]+-[A-Za-z]*r[A-Za-z]*f?[A-Za-z]*[[:space:]]+(\./)?\.workingdir2?(/|\s|$)'; then
  deny "'rm -rf .workingdir*' destroys plan + state."
fi
exit 0
```

---

## 7.3 `hooks/guard-go-edit.sh` (REBUILD for goenvoy)

Replaces golusoris's rules with goenvoy-specific ones:

**Rules**:

1. **No non-stdlib imports** (pure-stdlib invariant, ADR-0001). Bans any `import "github.com/…"` or `"gopkg.in/…"` landing in a non-test file. `golang.org/x/…` allowed. Testing packages (`*_test.go`, `example_test.go`) exempt (tests can still import stdlib only, but content check is light — let `depguard` catch it in CI).
2. **No `InsecureSkipVerify: true`** without a same-line `//nolint:gosec // <reason>`.
3. **`//nolint` needs justification** — same rule as golusoris.
4. **No live-API URLs in tests** — crude regex over `_test.go` content for `https?://(?!localhost|127\.0\.0\.1)` *inside* a `Test*` function body. Low-false-positive heuristic: match string literals with API hosts from a denylist (`api.tmdb.org`, `api.trakt.tv`, etc.). Start with just the four most-commonly-reached live hosts; expand as burn occurs.

```bash
#!/usr/bin/env bash
# PreToolUse hook for Edit / Write on *.go files.
# Enforces the pure-stdlib invariant and common client-lib footguns.
set -euo pipefail
payload=$(cat)
command -v jq >/dev/null 2>&1 || exit 0

file_path=$(printf '%s' "$payload" | jq -r '.tool_input.file_path // ""')
tool_name=$(printf '%s' "$payload" | jq -r '.tool_name // ""')
[ -z "$file_path" ] && exit 0
[[ "$file_path" != *.go ]] && exit 0

case "$tool_name" in
  Write) content=$(printf '%s' "$payload" | jq -r '.tool_input.content // ""') ;;
  Edit)  content=$(printf '%s' "$payload" | jq -r '.tool_input.new_string // ""') ;;
  *)     exit 0 ;;
esac
[ -z "$content" ] && exit 0

deny() { printf 'blocked by .claude/hooks/guard-go-edit.sh (%s):\n  %s\n' "$file_path" "$1" >&2; exit 2; }

# ── rule 1: no non-stdlib imports (ADR-0001) ─────────────────────────────────
# Fires on Write + Edit where an import line lands. Test files NOT exempt —
# tests must also be stdlib-only. `golang.org/x/...` explicitly allowed.
if printf '%s' "$content" | grep -E '^[[:space:]]*"(github\.com|gopkg\.in|gitlab\.com|bitbucket\.org)/' >/dev/null; then
  deny "non-stdlib import blocked — goenvoy is pure stdlib (ADR-0001). If you truly need a dep, write an ADR first."
fi

# ── rule 2: InsecureSkipVerify without justification ────────────────────────
if printf '%s' "$content" | grep -E 'InsecureSkipVerify:[[:space:]]*true' >/dev/null; then
  if ! printf '%s' "$content" | grep -E 'InsecureSkipVerify:[[:space:]]*true.*//[[:space:]]*nolint:gosec.*//' >/dev/null; then
    deny "InsecureSkipVerify: true requires a same-line '//nolint:gosec // <reason>' justification."
  fi
fi

# ── rule 3: //nolint without justification ──────────────────────────────────
if printf '%s' "$content" | grep -E '//[[:space:]]*nolint(:[[:alnum:],_-]+)?[[:space:]]*$' >/dev/null; then
  deny "//nolint needs a same-line justification, e.g. '//nolint:errcheck // defer Close, error surfaced'."
fi

# ── rule 4: live-API host in test files ─────────────────────────────────────
case "$file_path" in
  *_test.go)
    if printf '%s' "$content" | grep -E '"https?://(api\.tmdb\.org|api\.trakt\.tv|api\.themoviedb\.org|anilist\.co|graphql\.anilist\.co|kitsu\.io|api\.thetvdb\.com)' >/dev/null; then
      deny "live-API URL in test. Use httptest.NewServer — goenvoy tests MUST NOT hit real APIs."
    fi
  ;;
esac

exit 0
```

---

## 7.4 `hooks/format-go-write.sh` (ADAPT)

Swap the gci prefix:

```bash
#!/usr/bin/env bash
set -uo pipefail
payload=$(cat)
command -v jq >/dev/null 2>&1 || exit 0
file_path=$(printf '%s' "$payload" | jq -r '.tool_input.file_path // ""')
[ -z "$file_path" ] && exit 0
[[ "$file_path" != *.go ]] && exit 0
[ ! -f "$file_path" ] && exit 0

if command -v gofumpt >/dev/null 2>&1; then gofumpt -w "$file_path" 2>/dev/null || true; fi
if command -v gci >/dev/null 2>&1; then
  gci write --skip-generated -s standard -s default -s 'prefix(github.com/golusoris/goenvoy)' "$file_path" 2>/dev/null || true
fi
exit 0
```

---

## 7.5 `hooks/README.md`

Adapt golusoris's hooks README to the three goenvoy rules + exemption list.

---

## 7.6 Skills

Each skill is a Markdown file with frontmatter + instructions. Claude Code reads them when the user types `/<skill-name>`. Keep the same structure golusoris uses.

### 7.6.1 `skills/add-service-client.md` — NEW

```markdown
---
name: add-service-client
description: Scaffold a new pure-stdlib Go API-client module under a given category.
---

# Skill — `/add-service-client`

Scaffold a new service-client module from a one-line prompt.

## When to use

The user says "add a client for <service>" and you need to create an entire new module (directory, `go.mod`, doc/types/impl/test/example files, `AGENTS.md`, `docs/upstream/<service>.md`).

## Expected arguments

- `$1` — category, one of `arr | metadata/video | metadata/anime | metadata/music | metadata/tracking | metadata/book | metadata/game | metadata/adult | downloadclient | mediaserver | anime`.
- `$2` — service slug (kebab-case → used as both directory name and Go package name, e.g. `aria2`, `jellyseerr`).
- `$3` — upstream-API docs URL (pinned).
- (optional) `$4` — auth model: one of `apikey | basic | oauth-device | oauth-auth-code-pkce | jwt | none`. Default: `apikey`.

## Steps

1. Verify the target directory does not exist.
2. Create `<category>/<service>/` and inside it:
   - `go.mod` — `module github.com/golusoris/goenvoy/<category>/<service>` · `go 1.26.1`. No `require` block (pure stdlib).
   - `doc.go` — package-level comment: one sentence, ends with a period.
   - `types.go` — placeholder struct for `<Service>Response` + `APIError`.
   - `<service>.go` — `New(baseURL, apiKey string, opts ...Option) (*Client, error)` + `Option` type + `WithHTTPClient` + `WithTimeout` + `WithHeader` + helper `do(ctx, method, path, body, out) error`.
   - `<service>_test.go` — table-driven `TestNew` + one HTTP method test using `httptest.NewServer`.
   - `example_test.go` — `func ExampleNew()` that shows the idiomatic construction.
   - `AGENTS.md` — auth model, pagination style, error shape, pinned upstream URL, last-verified date.
3. Append the new module to `go.work` (`use ./<category>/<service>`).
4. Add `docs/upstream/<category>-<service>.md` with URL + version + today's date + a one-paragraph "what this API does".
5. Add a `## Unreleased` stub under the module in the root `CHANGELOG.md`.
6. From the new module directory run: `go build ./... && go vet ./... && go test -race ./... && golangci-lint run --config ../../.golangci.yml ./...`.
7. Report to the user: module path, files created, next steps (typically: wire the upstream API's methods).

## Template — `<service>.go`

```go
// Package <service> is a pure-stdlib Go client for the <Service> API.
package <service>

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"
)

// Client talks to a <Service> server.
type Client struct {
    baseURL    *url.URL
    apiKey     string
    httpClient *http.Client
    headers    http.Header
}

// Option configures a Client in New.
type Option func(*Client)

// WithHTTPClient replaces the default *http.Client.
func WithHTTPClient(c *http.Client) Option { return func(x *Client) { if c != nil { x.httpClient = c } } }

// WithTimeout sets the client request timeout. Ignored if WithHTTPClient is also set.
func WithTimeout(d time.Duration) Option { return func(x *Client) { x.httpClient.Timeout = d } }

// WithHeader adds a request header sent on every call.
func WithHeader(k, v string) Option { return func(x *Client) { x.headers.Set(k, v) } }

// New constructs a Client. baseURL must be absolute http or https.
func New(baseURL, apiKey string, opts ...Option) (*Client, error) {
    u, err := url.Parse(strings.TrimRight(baseURL, "/"))
    if err != nil {
        return nil, fmt.Errorf("<service>: parse baseURL: %w", err)
    }
    if u.Scheme != "http" && u.Scheme != "https" {
        return nil, fmt.Errorf("<service>: baseURL scheme must be http or https, got %q", u.Scheme)
    }
    c := &Client{
        baseURL:    u,
        apiKey:     apiKey,
        httpClient: &http.Client{Timeout: 30 * time.Second},
        headers:    http.Header{},
    }
    for _, o := range opts {
        o(c)
    }
    return c, nil
}

// APIError is returned when the <Service> API responds with a non-2xx status.
type APIError struct {
    StatusCode int
    Message    string
    Body       string
}

func (e *APIError) Error() string {
    if e.Message != "" {
        return fmt.Sprintf("<service>: HTTP %d: %s", e.StatusCode, e.Message)
    }
    return fmt.Sprintf("<service>: HTTP %d", e.StatusCode)
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
    u := *c.baseURL
    u.Path = strings.TrimRight(u.Path, "/") + "/" + strings.TrimLeft(path, "/")
    var rdr io.Reader
    if body != nil {
        b, err := json.Marshal(body)
        if err != nil {
            return fmt.Errorf("<service>: marshal body: %w", err)
        }
        rdr = strings.NewReader(string(b))
    }
    req, err := http.NewRequestWithContext(ctx, method, u.String(), rdr)
    if err != nil {
        return fmt.Errorf("<service>: new request: %w", err)
    }
    for k, vs := range c.headers {
        for _, v := range vs {
            req.Header.Add(k, v)
        }
    }
    // TODO: authentication scheme (api-key query, bearer, basic, etc.).
    req.Header.Set("Accept", "application/json")
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("<service>: %s %s: %w", method, u.Path, err)
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        raw, _ := io.ReadAll(resp.Body)
        return &APIError{StatusCode: resp.StatusCode, Body: string(raw)}
    }
    if out == nil {
        return nil
    }
    if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
        return fmt.Errorf("<service>: decode response: %w", err)
    }
    return nil
}

var _ = errors.New // placeholder until sentinels added
```
```

### 7.6.2 `skills/add-service-method.md` — NEW

Steps to add a single typed method to an existing client: pick endpoint, add request/response types to `types.go`, implement method calling `c.do`, add table-driven test case, add an `example_test.go` example if the surface is new.

### 7.6.3 `skills/bump-module.md` — NEW

Input: module path + level (`major | minor | patch`). Bumps the module's version in the root `CHANGELOG.md` header, runs `tools/release-check.sh`, opens a PR with title `release(<module>): v<new>`.

### 7.6.4 `skills/release-module.md` — NEW

Given a module + version: verify `CHANGELOG.md` has the stanza, run `tools/release-check.sh`, then create the tag `<module>/v<version>` and push. The tag push triggers `release.yml` (cosign + SBOM + provenance + GH Release).

### 7.6.5 `skills/audit-service-docs.md` — NEW

Walk every `docs/upstream/<service>.md`. For each: HEAD-check the pinned URL, update `Last verified: YYYY-MM-DD`, flag 404s for user follow-up.

---

## 7.7 Smoke tests (for the hook scripts themselves)

Include the same "smoke-testing a hook" section in `hooks/README.md` as golusoris — it's the fastest way to confirm a hook change didn't break parsing.
