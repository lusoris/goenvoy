# Contributing to goenvoy

Thank you for considering a contribution! Here's how to get started.

## Development Setup

1. **Clone the repo**

   ```bash
   git clone https://github.com/lusoris/goenvoy.git
   cd goenvoy
   ```

2. **Use Go 1.26+** — all modules target `go 1.26.1`.

3. **Workspace mode** — the `go.work` file links all 62 modules for local development:

   ```bash
   go work sync
   ```

4. **Run tests**

   ```bash
   make test-all
   ```

5. **Lint**

   ```bash
   make lint-all
   ```

## Project Structure

This is a multi-module monorepo. Each service lives in its own Go module under a category directory:

```
arr/sonarr/         → github.com/lusoris/goenvoy/arr/sonarr
metadata/video/tmdb → github.com/lusoris/goenvoy/metadata/video/tmdb
```

Parent packages (`arr/`, `metadata/`, etc.) contain shared types used by their child modules.

## Adding a New Service

1. Create a directory under the appropriate category (e.g. `downloadclient/aria2/`).
2. Run `go mod init github.com/lusoris/goenvoy/downloadclient/aria2`.
3. Create the standard files:
   - `doc.go` — package documentation
   - `types.go` — request/response types
   - `aria2.go` — client implementation
   - `aria2_test.go` — tests (use `httptest.NewServer`)
   - `example_test.go` — runnable examples for godoc
   - `go.mod` — module definition
4. Add the module to `go.work`.
5. Verify:
   ```bash
   cd downloadclient/aria2
   go test -race ./...
   golangci-lint run ./...
   ```

## Code Style

- **Pure stdlib** — no external dependencies. All HTTP is `net/http`, all JSON is `encoding/json`.
- **Functional options** — use `Option` type and `With*` constructors.
- **Context propagation** — every method that does I/O takes `context.Context` as the first argument.
- **Error types** — each module defines an `APIError` struct implementing `error`.
- **Test with httptest** — mock real API responses, don't hit live servers.
- **Lint clean** — `golangci-lint run ./...` must pass with the repo's `.golangci.yml`.
- **Godoc comments** — all exported types and functions must have doc comments ending with a period.

## Commit Messages

Use clear, descriptive messages. Reference the affected module(s):

```
arr/sonarr: add GetEpisodeFiles method

metadata/video/tvdb: fix token refresh on 401
```

## Pull Requests

- One feature/fix per PR.
- All tests must pass (`make test-all`).
- Lint must be clean (`make lint-all`).
- Include tests for new functionality.
- Update `example_test.go` if the public API changes.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
