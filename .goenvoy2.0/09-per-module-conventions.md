# 09 · Per-module conventions — `AGENTS.md` files

goenvoy is a monorepo of ~63 modules grouped under five top-level categories. Each category + each service gets its own `AGENTS.md`. golusoris does the same per-subpackage.

Rule of thumb: **category-level AGENTS.md** captures conventions shared across all services in that category (auth model, pagination style, shared types). **Service-level AGENTS.md** captures per-service oddities (rate-limit hint, upstream-API caveats, deprecated endpoints).

Land in phase 5 — after the root governance + CI + hooks are green. Creating 60+ files is a one-time batch; content fills in organically as services get touched.

---

## 9.1 Category-level template

`<category>/AGENTS.md` carries ~40-80 lines. Template:

```markdown
# AGENTS — <category>

> Per-category conventions for **<category>/** modules. Read repo-root [AGENTS.md](../AGENTS.md) + [.workingdir/PRINCIPLES.md](../.workingdir/PRINCIPLES.md) first.

## What this category is

<one-sentence description — what services live here, what they have in common>

## Shared types (in `<category>/`)

The root of the category (no `/service` segment) is itself a small module containing types shared across services — e.g. `arr/types.go` defines `SystemStatus`, `RootFolder`, `Tag` that every *arr service returns in the same shape.

Service sub-modules **may** import from the category root (e.g. `arr/sonarr` imports `arr`). This is the only in-repo import edge allowed — no service imports another service.

## Shared conventions in this category

- **Auth model**: <API-key in header / cookie-session / OAuth2 / …>
- **Base-URL shape**: <e.g. `/api/v3` for arr-stack; no trailing slash>
- **Pagination**: <cursor / offset+limit / page+size / none>
- **Error body**: <JSON with `error` field / problem+json / plain text>
- **Date format**: <RFC 3339 / Unix seconds / …>

## Modules in this category

| Module | Purpose | Pinned upstream version |
|---|---|---|
| `arr/sonarr` | Sonarr v3 client | see [docs/upstream/arr-sonarr.md](../docs/upstream/arr-sonarr.md) |
| … | … | … |

## When adding a new service here

Use the `/add-service-client <category> <service> <docs-url>` skill.

Must-haves at file creation:
- `doc.go` — one-sentence package comment.
- `types.go` — request/response types + `APIError`.
- `<service>.go` — `New` + `Option`s + `do` helper.
- `<service>_test.go` — table-driven tests, `httptest` only.
- `example_test.go` — runnable godoc example.
- `AGENTS.md` — per-service conventions (copy [_template_](../_meta/service-agents-template.md)).
```

---

## 9.2 Service-level template

`<category>/<service>/AGENTS.md` — minimal, ~20-40 lines:

```markdown
# AGENTS — <category>/<service>

> Per-service notes. Read [../AGENTS.md](../AGENTS.md) first.

## Upstream API

- Canonical docs: <URL>
- Pinned version / commit: <semver / date>
- Last verified: <YYYY-MM-DD>

## Auth model

<e.g. "X-Api-Key header"; "Bearer JWT, rotated via /refresh"; "OAuth2 device code">

## Pagination

<e.g. "cursor-based — nextCursor in JSON envelope"; "Link header with RFC 5988">

## Rate limits

<e.g. "5 req/s per API key, HTTP 429 on exceed">

## Known quirks

- <concrete odd behaviour, e.g. "API returns empty array with 200 for unauthenticated calls to /tag">
- <...>

## Testing notes

<e.g. "JSON responses are camelCase except /status which is snake_case; fixture files under testdata/ reflect both">
```

---

## 9.3 Category-by-category plan

### `arr/AGENTS.md`

- Shared types: `SystemStatus`, `RootFolder`, `Tag`, `QualityProfile`, `LanguageProfile`.
- Auth model: `X-Api-Key` header (all `*arr`).
- Pagination: none — most endpoints return complete lists.
- Base URL: `http(s)://<host>/api/v<major>` — callers pass the `<host>` portion; client adds the API prefix.
- Error body: JSON with `{error: string}` on 4xx; plain text/HTML on 5xx.
- 13 modules listed.

### `metadata/AGENTS.md`

- Sub-categorised: `video/`, `anime/`, `music/`, `tracking/`, `adult/`, `book/`, `game/`.
- Shared types (root `metadata/`): `Rating`, `Image`, `Person`, `Episode`, `Season`.
- Auth diversity: API key, bearer, OAuth2. Each sub-category documents its own.
- 34 modules.

### `downloadclient/AGENTS.md`

- Shared `Downloader` interface (in `downloadclient/types.go`): `ListTorrents`, `AddTorrent`, `GetInfo`, etc. Every service implements the full or a documented subset.
- Wire variety: JSON-RPC (Deluge, NZBGet), XML-RPC (rTorrent), custom (qbit), RPC (Transmission).
- 6 modules.

### `mediaserver/AGENTS.md`

- Auth diversity: Plex (`X-Plex-Token`), Jellyfin/Emby (`X-Emby-Token`), Audiobookshelf (Bearer), etc.
- No shared `MediaServer` interface — differences too large. Shared types at root for cross-references (`Library`, `Item`).
- 10 modules.

### `anime/AGENTS.md`

- Currently one module (`shoko`). Future anime-specific services land here.

---

## 9.4 Per-service AGENTS.md — batch bootstrap

There are ~55 service modules. Bootstrap all of them in one phase-5 PR by running a script that emits a skeleton into each. Each skeleton uses the template from §9.2 with placeholders filled from:

- The module's existing `doc.go` (to infer purpose).
- The module's upstream URL (grepped from comments / README hints).

Expected yield after bootstrap: the auto-filled fields get the upstream URL + last-verified date, the rest is `<TODO: fill in>`. Humans / agents fill these in when they next touch the service.

Script outline (`tools/bootstrap-service-agents.sh`):

```bash
#!/usr/bin/env bash
set -euo pipefail
find . -mindepth 3 -name 'go.mod' -not -path './.workingdir*/*' | while read -r modfile; do
  dir=$(dirname "$modfile")
  [ -f "$dir/AGENTS.md" ] && continue
  svc=$(basename "$dir")
  cat_path=$(dirname "$dir" | sed 's|^\./||')
  cat > "$dir/AGENTS.md" <<EOF
# AGENTS — ${cat_path}/${svc}

> Per-service notes. Read [../AGENTS.md](../AGENTS.md) first.

## Upstream API

- Canonical docs: <TODO>
- Pinned version / commit: <TODO>
- Last verified: 2026-04-14

## Auth model

<TODO>

## Pagination

<TODO>

## Rate limits

<TODO>

## Known quirks

- <TODO>

## Testing notes

<TODO>
EOF
done
```

---

## 9.5 Meta — `_meta/service-agents-template.md`

Keep the template itself in `_meta/service-agents-template.md` so the `/add-service-client` skill and the bootstrap script both pull from one place. Golusoris has a similar templating pattern for its fx-module scaffolding.

---

## 9.6 Relationship with `docs/upstream/<service>.md`

Each `docs/upstream/<category>-<service>.md` is a one-pager pinning the upstream API docs — URL, version, last verified. It's the **source of truth** for "what does the upstream say?". Service `AGENTS.md` files cite this file rather than duplicate the URL.

Split of concerns:
- `docs/upstream/<x>.md` — **external** API docs pin.
- `<category>/<service>/AGENTS.md` — **internal** client conventions / quirks.

`/audit-service-docs` skill (see [07](07-claude-hooks-and-skills.md)) refreshes the `docs/upstream/*.md` files' Last-verified dates without touching the service `AGENTS.md` files.
