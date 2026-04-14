# 01 · Gap analysis — what goenvoy has today vs golusoris standards

Snapshot: `/home/kilian/dev/goenvoy` as of 2026-04-14, compared against `/home/kilian/dev/golusoris@701d4e4`.

Legend: ✅ present · 🟡 partial · ❌ missing · ➖ not applicable (pure-lib scope).

## Root-level governance docs

| Artefact | goenvoy | golusoris | Gap | Action |
|---|---|---|---|---|
| `README.md` | ✅ solid, module catalog + usage | ✅ | — | Minor tweak: add badges (Go Reference, OSSF Scorecard, OpenSSF Best Practices). |
| `LICENSE` | ✅ MIT | ✅ MIT | — | — |
| `CONTRIBUTING.md` | 🟡 basic | ✅ Conventional-Commits contract + Migration: footer + pre-commit hooks | Missing CC spec, missing Migration: footer, missing pre-commit block | Rewrite — see [03-governance-docs.md](03-governance-docs.md). |
| `CODE_OF_CONDUCT.md` | ❌ | ✅ Contributor Covenant v2.1 | missing | Add — verbatim adapt from golusoris. |
| `SECURITY.md` | ✅ short | ✅ deeper (SLSA verify + cosign + SBOM + Renovate/Dependabot) | Missing cosign-verify block, SBOM+provenance claims, Renovate mention | Rewrite — see [03](03-governance-docs.md). |
| `AGENTS.md` | ❌ | ✅ cross-tool agent guide | missing | Add — adapted per [03](03-governance-docs.md). |
| `CLAUDE.md` | ❌ | ✅ skill/hook catalog + Don't list | missing | Add — adapted per [03](03-governance-docs.md). |
| `CHANGELOG.md` | ✅ Keep-a-Changelog format, manually updated | ✅ release-please-generated | — | Keep manual format; add release-please option in [06](06-release-and-versioning.md). |
| `.editorconfig` | ❌ | ✅ | missing | Add — copy golusoris verbatim, tweak none needed. |
| `.markdownlintignore` | ❌ | ✅ | missing | Add (small). |

## Lint / formatting / tooling

| Artefact | goenvoy | golusoris | Gap | Action |
|---|---|---|---|---|
| `.golangci.yml` at root | ✅ good but minimal set (~25 linters, revive-lite) | ✅ 30+ linters in `tools/golangci.yml` (includes `contextcheck`, `fatcontext`, `containedctx`, `exhaustive`, `wrapcheck`, `depguard`, `gomodguard`, `gomoddirectives`, `forbidigo`, `paralleltest`, `tparallel`, `thelper`, `usetesting`, `testifylint`, `sloglint`, `spancheck`, `loggercheck`) | Missing ~15 linters + `depguard`/`gomodguard` for **pure-stdlib enforcement**, missing `forbidigo` for the few API-client-specific bans | Expand — see [04-lint-and-tooling.md](04-lint-and-tooling.md). |
| `gofumpt` | ❌ (plain `gofmt -s -w`) | ✅ | missing | Add via golangci-lint `formatters:` + hook. |
| `gci` (import grouping) | 🟡 (goimports-local-prefixes) | ✅ gci with 3-group prefix | partial | Upgrade — gci with `prefix(github.com/golusoris/goenvoy)`. |
| `golines` (line-length) | ❌ | ✅ 120-col cap | missing | Optional — add as lint-only warning first. |
| `Makefile` | ✅ `test-all`, `lint-all`, `tidy-all`, `vet-all`, `fmt-all`, `build-all` | ✅ plus `vuln`, `gosec`, `ci`, `cover`, `gen`, `spec-lint` (framework scope) | Missing `vuln-all`, `gosec-all`, `ci-all`, `cover-all` | Extend — see [04](04-lint-and-tooling.md). |
| `tools/` directory | ❌ | ✅ `Makefile.shared`, `golangci.yml`, `mockery.yaml`, … | N/A for a lib | Create only `tools/` iff we want to split Makefile; otherwise keep root Makefile. |

## CI / supply-chain workflows (`.github/workflows/`)

| Workflow | goenvoy | golusoris | Gap | Action |
|---|---|---|---|---|
| `ci.yml` — discover + test matrix | ✅ per-module matrix (discover → test+lint) | ✅ with extra jobs | Missing PR-title CC check, missing gosec, missing govulncheck, missing apidiff-per-module, missing coverage threshold gate, missing build job | Expand — see [05-ci-workflows.md](05-ci-workflows.md). |
| `release.yml` — per-module tag push | ✅ minimal: validate → GH Release | ✅ cosign + syft + SLSA | Missing cosign keyless, missing syft SBOM, missing SLSA provenance, missing goreleaser config | Expand — see [06-release-and-versioning.md](06-release-and-versioning.md). |
| `release-all.yml` — tag all modules at once | ✅ | — | — | Keep as-is; add SLSA hooks. |
| `codeql.yml` | ❌ | ✅ | missing | Add — near-verbatim. |
| `scorecard.yml` | ❌ | ✅ | missing | Add — near-verbatim. |
| `release-please.yml` + config | ❌ | ✅ | missing | Add — **multi-package mode** keyed on each module path (see [06](06-release-and-versioning.md)). |
| `rebuild-on-base.yml` | — | ✅ | ➖ N/A (no Docker image to rebuild on base-image CVE) | skip |
| `auto-assign.yml` | ❌ | ✅ | missing | Add — low-value but free. |
| `dependabot.yml` | ✅ | ✅ | unknown depth | Verify — should cover `gomod` + `github-actions`, grouped by ecosystem. |
| `pull_request_template.md` | ✅ | ✅ (similar) | verify | Verify content aligns with Conventional-Commits + Migration: footer expectation. |
| `ISSUE_TEMPLATE/` | ✅ `.md` | ✅ `.yml` (structured) | format difference | Upgrade to `.yml` forms for better triage. |

## AI / agent layer

| Artefact | goenvoy | golusoris | Gap | Action |
|---|---|---|---|---|
| `.claude/settings.json` | ❌ | ✅ PreToolUse (Bash + Edit/Write) + PostToolUse (Edit/Write) | missing | Add — see [07-claude-hooks-and-skills.md](07-claude-hooks-and-skills.md). |
| `.claude/hooks/guard-bash.sh` | ❌ | ✅ blocks `--no-verify`, force-push-to-main, `rm -rf .git`/`.workingdir` | missing | Copy verbatim; add `.workingdir2` path; add `--no-verify` block. |
| `.claude/hooks/guard-go-edit.sh` | ❌ | ✅ bans `time.Now()` outside `clock/`, `fmt.Print*`, unjustified `//nolint` | missing | **Adapt** — drop `time.Now()`/`fmt.Print*` bans (not relevant); add: no non-stdlib imports (pure-stdlib invariant), no TLS `InsecureSkipVerify` without `//nolint:gosec // justification`, keep `//nolint` rule. |
| `.claude/hooks/format-go-write.sh` | ❌ | ✅ gofumpt + gci on save | missing | Copy, swap gci prefix to `github.com/golusoris/goenvoy`. |
| `.claude/skills/*.md` | ❌ | ✅ 5 skills | missing | Write **new** goenvoy-specific skills: `/add-service-client`, `/add-service-method`, `/bump-module`, `/release-module`, `/audit-service-docs`. |
| Pinned upstream docs (`docs/upstream/`) | ❌ | ✅ | ➖ different model — goenvoy services have public OpenAPI specs of varying quality | Add lightweight `docs/upstream/<service>.md` with pinned API-docs URL + version + last-checked date. Not auto-loaded in hooks (for now). |

## Docs + architecture

| Artefact | goenvoy | golusoris | Gap | Action |
|---|---|---|---|---|
| `docs/adr/` | ❌ | ✅ Nygard + index + backfills | missing | Add — see [08-adrs-and-architecture.md](08-adrs-and-architecture.md). |
| `docs/adr/0000-template.md` | ❌ | ✅ | missing | Copy verbatim. |
| `docs/adr/README.md` | ❌ | ✅ | missing | Adapt — record what's backfilled. |
| `docs/architecture/` (C4) | ➖ | ✅ | ➖ N/A | skip — client libs don't have a runtime architecture. |
| `docs/principles.md` | ❌ | ✅ pointer doc | — | Not needed — principles live in [02-principles-adapted.md](02-principles-adapted.md) → materialised into `AGENTS.md`. |
| `docs/compliance/` | ➖ | ✅ | ➖ N/A for a lib | skip. |
| Per-category `AGENTS.md` (arr/, metadata/, downloadclient/, mediaserver/, anime/) | ❌ | ✅ (per subpackage) | missing | Add — see [09-per-module-conventions.md](09-per-module-conventions.md). |
| Per-service `AGENTS.md` (arr/sonarr/, metadata/video/tmdb/, …) | ❌ | ✅ where relevant | missing (high volume — 55+ modules) | Template + bulk-create (skeleton). Fill-in happens organically as services get touched. |

## Working directory

| Artefact | goenvoy | golusoris | Gap | Action |
|---|---|---|---|---|
| `.workingdir/` checked-in session state | 🟡 (mentioned in `.gitignore`-exclude patterns in Makefile/CI but no folder yet) | ✅ has PLAN.md, STATE.md | missing in goenvoy tree | Add `.workingdir/PLAN.md` (goenvoy-scoped) + `.workingdir/STATE.md` (session log) once phase 1 of rollout lands. |

## Dependency + vulnerability management

| Artefact | goenvoy | golusoris | Gap | Action |
|---|---|---|---|---|
| Dependabot | ✅ | ✅ | verify depth | Ensure groups: `gomod-prod`, `gomod-dev`, `github-actions`. Pure-stdlib means `gomod` is quiet — Dependabot mostly patches `.github/workflows/` action versions. |
| Renovate | ❌ | (Dependabot only) | — | Not required. |
| `govulncheck` in CI | ❌ | ✅ | missing | Add — per-module in the matrix. |
| `gosec` in CI | 🟡 (via golangci-lint) | ✅ standalone SARIF upload | upgrade | Add dedicated job; SARIF → code-scanning. |
| `apidiff` vs previous tag | ❌ | ✅ | missing, **high value** for a lib that exports a stable surface | Add — run per-module, compare HEAD vs last `<module>/vX.Y.Z`. |

## Summary (counts)

- **18 missing** artefacts.
- **8 partial** / needs-upgrade artefacts.
- **4 verify** items.
- **~7 N/A** framework-only artefacts.

Rollout order in [10-rollout-checklist.md](10-rollout-checklist.md) groups these into 6 merge-able phases.
