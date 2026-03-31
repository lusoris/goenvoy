# Security Policy

## Supported Versions

Only the latest release of each module is supported with security updates.

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly:

1. **Do not** open a public issue.
2. Email **security@lusoris.dev** with:
   - A description of the vulnerability
   - Steps to reproduce
   - Affected module(s)
3. You will receive an acknowledgment within 48 hours.
4. A fix will be developed privately and released as a patch.

## Scope

This library is a collection of HTTP API clients. It does not run servers or handle user authentication directly. Relevant security concerns include:

- Credential leakage (API keys, tokens in logs or error messages)
- Request injection (path traversal, header injection)
- TLS verification bypass
- Improper error handling exposing sensitive data

## Security Practices

- All modules pass [gosec](https://github.com/securego/gosec) static analysis.
- All modules pass [golangci-lint](https://golangci-lint.run/) with `errorlint`, `gosec`, and `noctx` enabled.
- No external dependencies — pure Go stdlib reduces supply chain risk.
- API keys and tokens are never logged or included in error messages.
