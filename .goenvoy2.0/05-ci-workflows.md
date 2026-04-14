# 05 · CI workflows

Target: `.github/workflows/` in the goenvoy repo. Every new / replaced workflow listed here includes the full file content — ready to commit.

Key design decision: goenvoy is a **multi-module monorepo**, so every check runs inside a per-module matrix (`discover` job emits the module list as JSON, downstream jobs fan out). golusoris runs single-module CI — we re-use golusoris's workflow *structure* but adapt the job bodies to iterate per module.

All GitHub Action refs are **hash-pinned** (same pins as golusoris's workflows, which are audit-approved).

---

## 5.1 `ci.yml` (REPLACE)

Replaces the existing `ci.yml` with one that inherits the golusoris pattern (PR-title, lint, gosec, govulncheck, test+coverage, build, apidiff) — but fans every job across the module matrix.

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

permissions:
  contents: read

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  # ── Conventional-commit PR title ───────────────────────────────────────────
  pr-title:
    name: PR title (Conventional Commits)
    if: github.event_name == 'pull_request'
    runs-on: ubuntu-latest
    steps:
      - name: Check PR title
        env:
          PR_TITLE: ${{ github.event.pull_request.title }}
        run: |
          if ! echo "$PR_TITLE" | grep -qE '^(feat|fix|docs|chore|refactor|test|perf|ci|build|revert)(\(.+\))?(!)?: .+'; then
            echo "::error::PR title does not follow Conventional Commits."
            echo "::error::Got: $PR_TITLE"
            echo "::error::Expected: <type>(<scope>): <description>"
            exit 1
          fi
          echo "OK: $PR_TITLE"

  # ── Discover all modules ───────────────────────────────────────────────────
  discover:
    name: Discover modules
    runs-on: ubuntu-latest
    outputs:
      modules: ${{ steps.find.outputs.modules }}
    steps:
      - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6
      - id: find
        run: |
          modules=$(find . -name 'go.mod' \
            -not -path './.workingdir/*' \
            -not -path './.workingdir2/*' \
            -exec dirname {} \; | sort | jq -R -s -c 'split("\n") | map(select(. != ""))')
          echo "modules=$modules" >> "$GITHUB_OUTPUT"

  # ── Lint ───────────────────────────────────────────────────────────────────
  lint:
    name: Lint ${{ matrix.module }}
    needs: discover
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        module: ${{ fromJSON(needs.discover.outputs.modules) }}
    steps:
      - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6
      - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6
        with:
          go-version: stable
          cache-dependency-path: ${{ matrix.module }}/go.mod
      - uses: golangci/golangci-lint-action@1e7e51e771db61008b38414a730f564565cf7c20 # v9
        with:
          version: latest
          working-directory: ${{ matrix.module }}
          args: --config=${{ github.workspace }}/.golangci.yml --timeout=5m

  # ── Security (gosec) — aggregate across all modules, SARIF to code-scan ───
  gosec:
    name: Security (gosec)
    runs-on: ubuntu-latest
    permissions:
      contents: read
      security-events: write
    steps:
      - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6
      - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6
        with:
          go-version: stable
      - name: Install gosec
        run: go install github.com/securego/gosec/v2/cmd/gosec@v2.25.0
      - name: Run gosec per module
        run: |
          set -e
          modules=$(find . -name 'go.mod' -not -path './.workingdir*/*' -exec dirname {} \;)
          failed=0
          for mod in $modules; do
            echo "::group::gosec $mod"
            (cd "$mod" && gosec -exclude-generated -fmt sarif -out gosec.sarif ./... || echo "--- non-fatal sarif write")
            (cd "$mod" && gosec -exclude-generated ./...) || failed=1
            echo "::endgroup::"
          done
          # collect SARIFs into a single multi-run file
          jq -s '{version:"2.1.0",runs: map(.runs[0])}' $(find . -name gosec.sarif) > gosec-aggregate.sarif
          if [ $failed -ne 0 ]; then exit 1; fi
      - name: Upload SARIF
        if: always()
        uses: github/codeql-action/upload-sarif@c10b8064de6f491fea524254123dbe5e09572f13 # v4
        with:
          sarif_file: gosec-aggregate.sarif
        continue-on-error: true

  # ── Vulnerabilities (govulncheck) per module ──────────────────────────────
  vuln:
    name: Vulnerabilities ${{ matrix.module }}
    needs: discover
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        module: ${{ fromJSON(needs.discover.outputs.modules) }}
    steps:
      - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6
      - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6
        with:
          go-version: stable
          cache-dependency-path: ${{ matrix.module }}/go.mod
      - name: govulncheck
        working-directory: ${{ matrix.module }}
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@v1.1.4
          govulncheck ./...

  # ── Tests + coverage per module ───────────────────────────────────────────
  test:
    name: Test ${{ matrix.module }}
    needs: discover
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        module: ${{ fromJSON(needs.discover.outputs.modules) }}
    steps:
      - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6
      - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6
        with:
          go-version: stable
          cache-dependency-path: ${{ matrix.module }}/go.mod
      - name: go test -race + coverage
        working-directory: ${{ matrix.module }}
        run: go test -race -count=1 -timeout=5m -coverprofile=coverage.out -covermode=atomic ./...
      - name: Coverage threshold (70%)
        working-directory: ${{ matrix.module }}
        run: |
          if [ ! -s coverage.out ]; then echo "no coverage"; exit 0; fi
          COV=$(go tool cover -func=coverage.out | grep '^total' | awk '{print $3}' | tr -d '%')
          echo "Coverage: ${COV}%"
          if awk "BEGIN { exit !($COV < 70) }"; then
            echo "::error::${{ matrix.module }} coverage ${COV}% < 70%"
            exit 1
          fi
      - name: Upload coverage
        if: always()
        uses: actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a # v7
        with:
          name: coverage-${{ hashFiles(format('{0}/go.mod', matrix.module)) }}
          path: ${{ matrix.module }}/coverage.out
          retention-days: 7

  # ── Build per module ──────────────────────────────────────────────────────
  build:
    name: Build ${{ matrix.module }}
    needs: discover
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        module: ${{ fromJSON(needs.discover.outputs.modules) }}
    steps:
      - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6
      - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6
        with:
          go-version: stable
          cache-dependency-path: ${{ matrix.module }}/go.mod
      - working-directory: ${{ matrix.module }}
        run: go build ./...

  # ── apidiff per module vs last <module>/vX.Y.Z tag ────────────────────────
  apidiff:
    name: apidiff ${{ matrix.module }}
    needs: discover
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        module: ${{ fromJSON(needs.discover.outputs.modules) }}
    steps:
      - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6
        with:
          fetch-depth: 0
      - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6
        with:
          go-version: stable
          cache-dependency-path: ${{ matrix.module }}/go.mod
      - name: Install apidiff
        run: go install golang.org/x/exp/cmd/apidiff@v0.0.0-20250218142911-aa4b98e5adaa
      - name: Run apidiff for ${{ matrix.module }}
        env:
          MODULE_DIR: ${{ matrix.module }}
        run: |
          # Strip leading "./"
          MOD_PATH=${MODULE_DIR#./}
          TAG_PREFIX="${MOD_PATH}/v"
          PREV_TAG=$(git tag --list --sort=-v:refname "${TAG_PREFIX}*" | head -1 || true)
          if [ -z "$PREV_TAG" ]; then
            echo "No previous ${TAG_PREFIX}* tag — skipping (first release)."
            exit 0
          fi
          echo "Comparing $MOD_PATH against $PREV_TAG"
          cd "$MOD_PATH"
          MOD_IMPORT=$(go list -m)
          apidiff -m "$MOD_IMPORT" . > /tmp/current.txt 2>/dev/null || true
          SAVED=$(git rev-parse HEAD)
          git stash --include-untracked --quiet || true
          git checkout --quiet "$PREV_TAG"
          apidiff -m "$MOD_IMPORT" . > /tmp/prev.txt 2>/dev/null || true
          git checkout --quiet "$SAVED"
          git stash pop --quiet || true
          if ! apidiff /tmp/prev.txt /tmp/current.txt; then
            echo "::error::Breaking API change in $MOD_PATH vs $PREV_TAG — add a 'BREAKING CHANGE:' + 'Migration:' footer or bump major."
            exit 1
          fi
          echo "API backwards-compatible with $PREV_TAG."

  # ── Aggregate status ──────────────────────────────────────────────────────
  status:
    name: CI Status
    if: always()
    needs: [pr-title, lint, gosec, vuln, test, build, apidiff]
    runs-on: ubuntu-latest
    steps:
      - name: Check results
        run: |
          for j in lint gosec vuln test build apidiff; do
            case "$(eval echo \$${j}_result)" in ""|success|skipped) ;; *)
              echo "job $j: $(eval echo \$${j}_result)"; exit 1 ;;
            esac
          done
          echo "All checks passed."
        env:
          lint_result: ${{ needs.lint.result }}
          gosec_result: ${{ needs.gosec.result }}
          vuln_result: ${{ needs.vuln.result }}
          test_result: ${{ needs.test.result }}
          build_result: ${{ needs.build.result }}
          apidiff_result: ${{ needs.apidiff.result }}
```

---

## 5.2 `codeql.yml` (NEW)

Near-verbatim golusoris, minus the `workflow_call` (no downstream apps call this one):

```yaml
name: CodeQL

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    - cron: '30 7 * * 1'

permissions:
  contents: read

jobs:
  analyze:
    name: CodeQL Analysis
    runs-on: ubuntu-latest
    permissions:
      contents: read
      security-events: write
      actions: read
    strategy:
      fail-fast: false
      matrix:
        module:
          - anime
          - arr
          - downloadclient
          - mediaserver
          - metadata
    steps:
      - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6
      - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6
        with:
          go-version: stable
      - name: Initialize CodeQL
        uses: github/codeql-action/init@c10b8064de6f491fea524254123dbe5e09572f13 # v4
        with:
          languages: go
          queries: security-extended,security-and-quality
      - name: Build submodules under ${{ matrix.module }}/
        run: |
          for d in $(find ./${{ matrix.module }} -name 'go.mod' -exec dirname {} \;); do
            echo "==> $d"
            (cd "$d" && go build ./...)
          done
      - uses: github/codeql-action/analyze@c10b8064de6f491fea524254123dbe5e09572f13 # v4
        with:
          category: "/language:go/category:${{ matrix.module }}"
```

---

## 5.3 `scorecard.yml` (NEW)

Verbatim golusoris (drop `workflow_call` body — not needed downstream):

```yaml
name: Scorecard

on:
  push:
    branches: [main]
  schedule:
    - cron: '15 6 * * 1'

permissions:
  contents: read

jobs:
  analysis:
    name: Scorecard Analysis
    runs-on: ubuntu-latest
    permissions:
      contents: read
      security-events: write
      id-token: write
    steps:
      - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6
        with:
          persist-credentials: false
      - uses: ossf/scorecard-action@4eaacf0543bb3f2c246792bd56e8cdeffafb205a # v2.4.3
        with:
          results_file: results.sarif
          results_format: sarif
          publish_results: true
          repo_token: ${{ secrets.SCORECARD_TOKEN || github.token }}
      - uses: actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a # v7
        with:
          name: scorecard-sarif
          path: results.sarif
          retention-days: 5
      - uses: github/codeql-action/upload-sarif@c10b8064de6f491fea524254123dbe5e09572f13 # v4
        with:
          sarif_file: results.sarif
```

`SCORECARD_TOKEN` is optional — see golusoris `scorecard.yml` comment for the fine-grained-PAT scope.

---

## 5.4 `release-please.yml` + `release-please-config.json` (NEW)

See [06-release-and-versioning.md](06-release-and-versioning.md) §6.3.

---

## 5.5 `auto-assign.yml` (NEW)

Low-value but free — auto-assign the PR author as reviewer. Verbatim golusoris.

---

## 5.6 Dependabot — verify

Audit goenvoy's existing `.github/dependabot.yml` against golusoris's (not shown here). Ensure:

- `gomod` ecosystem coverage across **every module directory** (Dependabot is shallow — configure per-directory).
- `github-actions` ecosystem for `.github/workflows/`.
- Grouped updates where practical (e.g. all minor bumps together).

Full example (abbreviated — one entry per module directory, generated by a script at config-time):

```yaml
version: 2
updates:
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
    groups:
      actions:
        patterns: ["*"]

  # gomod — per-module (script-generated; entries omitted here)
  - package-ecosystem: "gomod"
    directory: "/arr"
    schedule: { interval: "weekly" }
  # ... repeat for every module path ...
```

Consider a small `tools/gen-dependabot.sh` that regenerates this on each new module.

---

## 5.7 Jobs to leave alone

- `release-all.yml` — already serviceable; extend in [06](06-release-and-versioning.md) to add cosign + SBOM + SLSA.
- `ISSUE_TEMPLATE/` — upgrade `.md` → `.yml` forms (golusoris has `.yml`). Content diff is minor.
- `pull_request_template.md` — verify contents align with Migration: footer expectation.

---

## 5.8 CI cost / runtime considerations

Per-module matrix at 63+ modules × 5 jobs (lint, vuln, test, build, apidiff) = ~315 jobs per PR. That's within GitHub-hosted-runner free-tier on public repos (unlimited) but will be slow — ~15-25 minute wall-clock.

Optimisations to land in phase 2 (after baseline is green):

- **Changed-files filter** — use [dorny/paths-filter](https://github.com/dorny/paths-filter) to only run the matrix for modules whose files changed. Repo-wide jobs (codeql, scorecard, gosec aggregate) still always run.
- **Action cache keyed on module go.sum** — already done via `cache-dependency-path`.
- **Share golangci-lint install across matrix** — use `golangci-lint-action` with `install-mode: goinstall` + cache.
