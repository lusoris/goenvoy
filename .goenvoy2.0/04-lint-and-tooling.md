# 04 · Lint + tooling baseline

Everything here targets the goenvoy repo root. Linter config is **monorepo-shared** — the matrix in CI runs the same config per module. `depguard` carries the pure-stdlib invariant.

---

## 4.1 `.golangci.yml` (REPLACE current)

Current goenvoy `.golangci.yml` is fine but ~15 linters short and missing the enforcement hooks (`depguard`, `gomodguard`, `forbidigo`, test-hygiene set). Proposed replacement — shared across all modules.

```yaml
# golangci-lint v2 config for goenvoy — applies to every module.
version: "2"

run:
  timeout: 5m
  go: "1.26"
  tests: true
  allow-parallel-runners: true

linters:
  default: none
  enable:
    # correctness
    - govet
    - staticcheck
    - errcheck
    - ineffassign
    - unused
    - gocritic
    - revive
    - bodyclose
    - errorlint
    - wrapcheck
    - nilerr
    - nilnil
    - nilnesserr
    - exhaustive
    - copyloopvar
    - intrange
    - predeclared
    - reassign
    - recvcheck
    - nakedret
    - musttag
    - contextcheck
    - fatcontext
    - containedctx
    - noctx
    # security
    - gosec
    - bidichk
    - asciicheck
    # complexity
    - gocyclo
    - funlen
    - gocognit
    - cyclop
    # style / hygiene
    - misspell
    - godot
    - godoclint
    - whitespace
    - unconvert
    - unparam
    - tagliatelle
    - usestdlibvars
    - dupword
    - canonicalheader
    - importas
    - perfsprint
    - prealloc
    - makezero
    # project discipline (the pure-stdlib gate)
    - forbidigo
    - depguard
    - gochecknoinits
    - gochecknoglobals
    - nolintlint
    # test hygiene
    - paralleltest
    - tparallel
    - thelper
    - usetesting
    - testifylint
    # supply-chain
    - gomoddirectives
    - gomodguard

  settings:
    gocyclo:
      min-complexity: 18
    funlen:
      lines: 120
      statements: 60
    gocognit:
      min-complexity: 30
    cyclop:
      max-complexity: 18
    gosec:
      excludes:
        - G104  # errors handled via wrapcheck/errcheck instead

    govet:
      enable-all: true
      disable:
        - fieldalignment

    revive:
      rules:
        - name: exported
          arguments: ["checkPrivateReceivers"]
        - name: indent-error-flow
        - name: unused-parameter
        - name: unused-receiver
        - name: package-comments
        - name: var-naming

    godot:
      scope: toplevel
      period: true

    goconst:
      min-len: 2
      min-occurrences: 3

    misspell:
      locale: US

    nakedret:
      max-func-lines: 30

    nolintlint:
      allow-unused: false
      require-explanation: true
      require-specific: true

    tagliatelle:
      case:
        rules:
          json: camel

    # The pure-stdlib gate. Every import outside stdlib + golang.org/x
    # is blocked with a pointer to ADR-0001.
    depguard:
      rules:
        stdlib-only:
          list-mode: lax
          files: ["$all"]
          allow:
            - $gostd
            - golang.org/x/
          deny:
            - pkg: "github.com/"
              desc: "goenvoy is pure stdlib (ADR-0001). External deps require a new ADR."
            - pkg: "gopkg.in/"
              desc: "goenvoy is pure stdlib (ADR-0001)."

    gomodguard:
      blocked:
        modules:
          - github.com/stretchr/testify:
              recommendations: ["standard library testing"]
              reason: "No test deps outside stdlib (ADR-0001)."

    forbidigo:
      forbid:
        - pattern: 'fmt\.Print(ln|f)?\('
          msg: "No fmt.Print* in library code. Return errors, don't print."
        - pattern: 'os\.Exit\('
          msg: "No os.Exit in library code."
        - pattern: 'panic\('
          msg: "No panics in library code — return an error."
        - pattern: 'InsecureSkipVerify:\s*true'
          msg: "InsecureSkipVerify: true needs a //nolint:gosec // <reason> justification."
      exclude-godoc-examples: true

    importas:
      no-unaliased: true

  exclusions:
    presets:
      - std-error-handling
      - common-false-positives
    rules:
      - path: (.+)_test\.go
        linters:
          - goconst
          - errcheck
          - gosec
          - funlen
          - gocyclo
          - gocognit
          - cyclop
          - forbidigo    # tests may print for debugging
          - wrapcheck
          - noctx
          - containedctx
      - path: doc\.go
        linters:
          - godoclint
      - path: example_test\.go
        linters:
          - forbidigo
          - wrapcheck
          - gosec

formatters:
  enable:
    - gofumpt
    - gci
    - goimports
  settings:
    gofumpt:
      module-path: github.com/golusoris/goenvoy
      extra-rules: true
    gci:
      sections:
        - standard
        - default
        - prefix(github.com/golusoris/goenvoy)
      custom-order: true
    goimports:
      local-prefixes:
        - github.com/golusoris/goenvoy
```

**Diff vs current**:

- Adds: `wrapcheck`, `contextcheck`, `fatcontext`, `containedctx`, `exhaustive`, `nilnil`, `nilnesserr`, `musttag`, `recvcheck`, `godoclint`, `tagliatelle`, `dupword`, `canonicalheader`, `importas`, `makezero`, `forbidigo`, `depguard`, `gomodguard`, `gomoddirectives`, `gochecknoglobals`, `paralleltest`, `tparallel`, `thelper`, `usetesting`, `testifylint`.
- Drops nothing critical.
- Adds `gofumpt` + `gci` to the `formatters:` section.

**Behavioural impact**:

- `depguard` + `gomodguard` will fire on **any existing test file importing `testify`**. First-pass rollout: expect zero hits (current tree verified pure-stdlib — `go.mod` files have no `require` blocks beyond `go 1.26.1`). If a module has a `testify` import, choose: (a) rewrite to `testing` + `if ... t.Fatalf(...)` (preferred; fast), or (b) file an ADR to exempt that module.
- `forbidigo` on `fmt.Print*` — check each module's `*.go` (excluding `*_test.go`, `example_test.go`, `doc.go`). Replace any with `return fmt.Errorf(...)`.
- `funlen`/`gocognit` may flag a few long decode functions. Refactor or `//nolint:funlen // API mirrors upstream schema, splitting adds no clarity`.

---

## 4.2 `.editorconfig` (NEW)

Verbatim golusoris:

```editorconfig
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.go]
indent_style = tab
indent_size = 4

[*.{md,yaml,yml,json,toml}]
indent_style = space
indent_size = 2

[Makefile]
indent_style = tab

[*.sh]
indent_style = space
indent_size = 2
```

---

## 4.3 Makefile — extend existing

Replace current `Makefile` with:

```makefile
.PHONY: help test-all lint-all vet-all vuln-all gosec-all cover-all tidy-all build-all fmt-all ci-all list-modules clean-all

MODULES := $(shell find . -name 'go.mod' -not -path './.workingdir/*' -not -path './.workingdir2/*' -exec dirname {} \;)

help:
	@awk 'BEGIN{FS=":.*?## "} /^[a-zA-Z_-]+:.*?## /{printf "  \033[36m%-14s\033[0m %s\n",$$1,$$2}' $(MAKEFILE_LIST)

list-modules: ## Print every discovered go.mod directory
	@printf '%s\n' $(MODULES)

test-all: ## go test -race + coverage, all modules
	@for mod in $(MODULES); do \
	  echo "==> Testing $$mod"; \
	  (cd $$mod && go test -race -count=1 -coverprofile=coverage.out -covermode=atomic ./...) || exit 1; \
	done

cover-all: test-all ## Generate coverage.html for each module
	@for mod in $(MODULES); do \
	  (cd $$mod && go tool cover -html=coverage.out -o coverage.html) || true; \
	done

lint-all: ## golangci-lint run, all modules
	@for mod in $(MODULES); do \
	  echo "==> Linting $$mod"; \
	  (cd $$mod && golangci-lint run ./...) || exit 1; \
	done

vet-all: ## go vet, all modules
	@for mod in $(MODULES); do \
	  (cd $$mod && go vet ./...) || exit 1; \
	done

vuln-all: ## govulncheck, all modules
	@for mod in $(MODULES); do \
	  echo "==> govulncheck $$mod"; \
	  (cd $$mod && govulncheck ./...) || exit 1; \
	done

gosec-all: ## gosec, all modules
	@for mod in $(MODULES); do \
	  echo "==> gosec $$mod"; \
	  (cd $$mod && gosec -quiet ./...) || exit 1; \
	done

tidy-all: ## go mod tidy, all modules
	@for mod in $(MODULES); do \
	  (cd $$mod && go mod tidy) || exit 1; \
	done

fmt-all: ## gofumpt + gci, all modules
	@for mod in $(MODULES); do \
	  (cd $$mod && gofumpt -w . && gci write --skip-generated -s standard -s default -s 'prefix(github.com/golusoris/goenvoy)' .) || exit 1; \
	done

build-all: ## go build, all modules
	@for mod in $(MODULES); do \
	  (cd $$mod && go build ./...) || exit 1; \
	done

clean-all: ## remove coverage artefacts
	@for mod in $(MODULES); do \
	  rm -f $$mod/coverage.out $$mod/coverage.html; \
	done

ci-all: lint-all vet-all gosec-all vuln-all test-all ## Full local CI

.PHONY: tools-install
tools-install: ## Install gofumpt, gci, golines, gosec, govulncheck locally
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest
	go install github.com/segmentio/golines@latest
	go install github.com/securego/gosec/v2/cmd/gosec@v2.25.0
	go install golang.org/x/vuln/cmd/govulncheck@v1.1.4
	go install golang.org/x/exp/cmd/apidiff@latest
```

---

## 4.4 `tools/` directory — optional, defer

golusoris uses `tools/golangci.yml` + `tools/Makefile.shared` because apps include it via fx bumps. goenvoy has no downstream-sharing requirement for its own lint config (each app has its own). **Decision: keep lint config at repo root (`/.golangci.yml`), no `tools/` split.**

Reasoning: the multi-module matrix already works with root-relative `.golangci.yml` (goenvoy's CI does this today). Adding `tools/` buys nothing in this scope.

---

## 4.5 Pre-commit config (optional, nice-to-have)

Proposed `.pre-commit-config.yaml` (opt-in; CI still the gate):

```yaml
repos:
  - repo: https://github.com/compilerla/conventional-pre-commit
    rev: v3.4.0
    hooks:
      - id: conventional-pre-commit
        stages: [commit-msg]

  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
      - id: go-fmt
      - id: golangci-lint
        args: [--config=.golangci.yml]

  - repo: https://github.com/gitleaks/gitleaks
    rev: v8.20.0
    hooks:
      - id: gitleaks
```

Land in phase 2; ship docs in `CONTRIBUTING.md` telling contributors to `pre-commit install`.

---

## 4.6 golines (line-length cap 120)

Not enabled in golangci-lint (`golines` isn't a linter, it's a formatter). Add to the Makefile's `fmt-all` target:

```
$(MAKE) -f ... fmt-all   # runs gofumpt + gci + optionally golines
```

Include `golines` as a separate `golines-all` target, gated by opt-in for phase-2 rollout:

```makefile
golines-all:
	@for mod in $(MODULES); do \
	  (cd $$mod && golines -m 120 -w .) || exit 1; \
	done
```

Do **not** run in CI as a gate yet — long lines in existing client code (typed struct fields with long JSON tags) would flood with noise. Gate later once formatter has been applied cleanly once.
