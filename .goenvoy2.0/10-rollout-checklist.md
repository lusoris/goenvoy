# 10 · Rollout checklist — phased, merge-on-green

Six phases. Each phase is one PR, mergeable independently, leaves `main` green. Earlier phases unblock later ones. Whole rollout is ~1 week of focused work or a couple of weeks part-time.

Every phase's PR title is Conventional Commits — since phase 1 ships the PR-title gate, the rollout PRs test it as they land.

---

## Phase 1 — Governance docs + editor config

**PR title**: `chore(governance): adopt golusoris standards (AGENTS, CLAUDE, CoC, SECURITY upgrade)`

**Files landed**:

- [ ] `AGENTS.md` (new) — [03-governance-docs.md §3.1](03-governance-docs.md).
- [ ] `CLAUDE.md` (new) — §3.2.
- [ ] `CONTRIBUTING.md` (rewrite) — §3.3.
- [ ] `SECURITY.md` (rewrite) — §3.4.
- [ ] `CODE_OF_CONDUCT.md` (new) — §3.5.
- [ ] `.editorconfig` (new) — §3.6.
- [ ] `.markdownlintignore` (new) — §3.7.
- [ ] `.gitignore` additions — §3.8.
- [ ] `README.md` — append badges + Standards section — §3.9.

**Gates to hit**: none yet (docs only). No CI changes in this phase.

**Checks before merge**:
- [ ] `AGENTS.md` repository layout tree matches current `find . -name go.mod` output.
- [ ] `SECURITY.md` cosign verify command uses `golusoris/goenvoy` regex.
- [ ] Every new doc cross-links are valid (`.workingdir/PRINCIPLES.md` placeholder — see phase 2).

---

## Phase 2 — Principles + ADRs + working dir

**PR title**: `docs(principles): ship adapted §2 contract + ADR-0001..0008 backfills`

**Files landed**:

- [ ] `.workingdir/PRINCIPLES.md` — copy of [02-principles-adapted.md](02-principles-adapted.md) with links fixed to live locations.
- [ ] `.workingdir/STATE.md` — one-liner session log entry.
- [ ] `docs/adr/README.md` — §8.1.
- [ ] `docs/adr/0000-template.md` — §8.2.
- [ ] `docs/adr/0001-pure-stdlib.md` … `0008-non-rest-wire-protocols.md` — §8.3.

**Gates**: none yet.

**Checks**:
- [ ] Every ADR has the five Nygard sections (Context / Decision / Consequences / Alternatives / Status).
- [ ] `AGENTS.md` from phase 1 links resolve (`.workingdir/PRINCIPLES.md` now exists).

---

## Phase 3 — Lint + tooling baseline

**PR title**: `ci(lint): expand .golangci.yml to 30+ linters (incl. depguard stdlib-only gate)`

**Files landed**:

- [ ] `.golangci.yml` (replace) — [04-lint-and-tooling.md §4.1](04-lint-and-tooling.md).
- [ ] `Makefile` (extend) — §4.3. Adds `vuln-all`, `gosec-all`, `cover-all`, `ci-all`, `tools-install`.
- [ ] `.pre-commit-config.yaml` (optional) — §4.5.

**Gates**: this phase **runs** the new lint config locally + in CI.

**Expected fallout (to land in the same PR)**:

- [ ] Fix any `godot` misses (doc comment period).
- [ ] Fix any `gci` import-group regressions.
- [ ] Fix any `bodyclose` misses.
- [ ] Fix any `contextcheck` / `noctx` misses — rare (every method already takes context).
- [ ] Fix any `errorlint` misses — wrap all errors with `%w`.
- [ ] Refactor or `//nolint:funlen // <reason>` any >120-line functions.
- [ ] Resolve every `depguard` / `gomodguard` hit — should be zero on a clean goenvoy tree. If any module has test-only external deps, delete them.

**Checks before merge**:
- [ ] `make ci-all` green locally.
- [ ] `make tools-install` installs gofumpt / gci / gosec / govulncheck / apidiff.
- [ ] `gofumpt -l $(find . -name '*.go' -not -path ./.workingdir'*')` empty.
- [ ] `gci diff --skip-generated -s standard -s default -s 'prefix(github.com/gogolusoris/goenvoy)' $(...)` empty.

---

## Phase 4 — CI hardening

**PR title**: `ci: add CodeQL + Scorecard; harden ci.yml (PR-title, gosec, govulncheck, apidiff, coverage gate)`

**Files landed**:

- [ ] `.github/workflows/ci.yml` (replace) — [05-ci-workflows.md §5.1](05-ci-workflows.md).
- [ ] `.github/workflows/codeql.yml` (new) — §5.2.
- [ ] `.github/workflows/scorecard.yml` (new) — §5.3.
- [ ] `.github/workflows/auto-assign.yml` (new) — §5.5.
- [ ] `.github/dependabot.yml` (verify/regenerate) — §5.6.
- [ ] `.github/pull_request_template.md` (verify) — includes Conventional-Commits + Migration: pointer.
- [ ] `.github/ISSUE_TEMPLATE/*.yml` (convert from .md) — §5.7.

**Gates**: every subsequent PR must pass all CI jobs.

**Expected first-run issues**:
- [ ] apidiff: none of the existing module tags exist in standard form yet? If any module has never been tagged, `PREV_TAG` is empty — the job no-ops (safe).
- [ ] Coverage gate: each module should already exceed 70%. If any module falls short, fix the test gap or temporarily adjust `coverage-threshold` for that matrix leg.
- [ ] CodeQL: first run establishes the baseline. Any finding with severity ≥ High needs addressing before merge.
- [ ] Scorecard: first score printed. Fix the low-hanging items (branch protection, pinned actions, signed releases — the last two land in phase 5).

**Checks**:
- [ ] Dry-run locally: `act` or a draft PR shows every matrix leg running.
- [ ] Badges in `README.md` render (phase 1 already added them; first run wires the badge target).

---

## Phase 5 — Release + supply-chain hardening

**PR title**: `ci(release): cosign keyless + syft SBOM + SLSA-L3 provenance on every module tag`

**Files landed**:

- [ ] `.github/workflows/release.yml` (replace) — [06-release-and-versioning.md §6.2](06-release-and-versioning.md).
- [ ] `.github/workflows/release-all.yml` (extend to push tags only, leave release.yml as the fan-out) — §6.5.
- [ ] `tools/release-check.sh` (new) — §6.6.
- [ ] (optional, phase 5b) `.github/workflows/release-please.yml` + `release-please-config.json` + `.release-please-manifest.json` — §6.3.

**Gates**: every future module tag produces signed artefacts.

**Validation protocol**:
- [ ] Push a no-op patch tag to a low-traffic module (e.g. `anime/shoko/v0.0.1-test1`) — confirm:
  - GH Release is created with tarball + sbom + sig + pem.
  - `cosign verify-blob` of `checksums.txt` against the pem succeeds.
  - `gh attestation verify` (or web UI) shows SLSA provenance.
- [ ] Delete the test tag + release afterwards.

**Checks**:
- [ ] `SECURITY.md` verify command runs against the test release successfully.
- [ ] `actions/attest-build-provenance` permissions (`id-token: write`, `attestations: write`) are on `release.yml`.

---

## Phase 6 — Claude hooks + skills + per-module AGENTS

**PR title**: `feat(ai): .claude/ hooks + 5 skills + category/service AGENTS.md bootstrap`

**Files landed**:

- [ ] `.claude/settings.json` — [07-claude-hooks-and-skills.md §7.1](07-claude-hooks-and-skills.md).
- [ ] `.claude/hooks/guard-bash.sh` — §7.2.
- [ ] `.claude/hooks/guard-go-edit.sh` — §7.3.
- [ ] `.claude/hooks/format-go-write.sh` — §7.4.
- [ ] `.claude/hooks/README.md` — §7.5.
- [ ] `.claude/skills/add-service-client.md` — §7.6.1.
- [ ] `.claude/skills/add-service-method.md` — §7.6.2.
- [ ] `.claude/skills/bump-module.md` — §7.6.3.
- [ ] `.claude/skills/release-module.md` — §7.6.4.
- [ ] `.claude/skills/audit-service-docs.md` — §7.6.5.
- [ ] `arr/AGENTS.md`, `metadata/AGENTS.md`, `downloadclient/AGENTS.md`, `mediaserver/AGENTS.md`, `anime/AGENTS.md` — [09-per-module-conventions.md §9.3](09-per-module-conventions.md).
- [ ] `_meta/service-agents-template.md` — §9.5.
- [ ] `tools/bootstrap-service-agents.sh` — §9.4.
- [ ] Bootstrap script run; ~55 service-level `AGENTS.md` files land as skeletons.
- [ ] `docs/upstream/*.md` — one per service with pinned URL + last-verified date.

**Gates**: none directly. The hooks start guiding future edits.

**Smoke tests (before merge)**:
- [ ] Each hook dry-runs per the `hooks/README.md` commands.
- [ ] Invoking `/add-service-client downloadclient aria2 https://github.com/aria2/aria2/wiki/aria2%E2%80%99s-methods` scaffolds a valid module that passes `make ci-all`.
- [ ] Invoking `/audit-service-docs` updates the `Last verified:` date in every `docs/upstream/*.md`.

---

## Exit criteria — whole rollout done when

1. ✅ Every `AGENTS.md`, `CLAUDE.md`, `SECURITY.md`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md` exists and matches the phase-1 drafts.
2. ✅ `.workingdir/PRINCIPLES.md` + `docs/adr/0001…0008` exist.
3. ✅ `.golangci.yml` has 30+ linters active; `make ci-all` is green on a clean checkout.
4. ✅ `ci.yml` + `codeql.yml` + `scorecard.yml` + `release.yml` green on `main`.
5. ✅ A test-tag round-trip verifies cosign + SBOM + SLSA on a real release.
6. ✅ `.claude/` directory live; each hook dry-runs; each skill produces the expected scaffolding.
7. ✅ Every category has `AGENTS.md`; every service has a `AGENTS.md` skeleton.
8. ✅ OSSF Scorecard ≥ 7 (adjust target as we see the baseline run).
9. ✅ The next functional PR touching any module passes every new gate without a `//nolint` ever needing to be added.

---

## Risks + mitigations

| Risk | Mitigation |
|---|---|
| Phase 3 lint config hits 100+ violations, blocking the PR. | Fix them inside the same PR — they're mechanical (`godot`, `gci`, `wrapcheck`). If one class is big (e.g. `gocognit` on a few decode paths), accept `//nolint:gocognit // <reason>` at the few sites rather than refactoring under time pressure. Track leftover refactors as follow-up tasks. |
| Phase 4 `apidiff` fires on an unintentional recent API drift. | Either fix the drift or bump the next release's major version. Don't disable the check. |
| Phase 5 test-tag leaves stray releases visible. | Use `-rc.1` pre-release tag, delete tag + release after validation. |
| Phase 6 bootstrap creates empty `AGENTS.md` skeletons that look abandoned. | Ship the skeleton with a clear `<TODO>` placeholder and a `Last verified:` date. Any touching PR is expected to fill in what it learns. |
| Per-module matrix CI is slow (60+ × 5 jobs). | Phase 4 lands the straight matrix; after it's stable, add a `dorny/paths-filter`-driven "changed modules only" optimisation in phase 7 (out of scope for this plan). |

---

## After exit — ongoing discipline

- Every new module goes through `/add-service-client` (skill ensures conventions land right the first time).
- Every new public API change includes the `Migration:` footer when it's breaking.
- Every ADR-worthy decision lands as a `docs/adr/0100+` file before the code changes.
- `docs/upstream/*.md` get a refresh sweep (`/audit-service-docs`) before each coordinated release-all.
