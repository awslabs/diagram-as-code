# Guardrails

## Security
- Never commit secrets, credentials, API tokens, or key files.
- Prefer environment variables for sensitive runtime configuration.
- Treat remote definition/icon downloads as untrusted input unless explicitly trusted.

## Stability
- Do not change default CLI behavior without updating docs and tests.
- Avoid changing PNG output behavior unintentionally.
- If output semantics change, call it out as breaking/non-breaking in PR notes.

## Testing
- Minimum gate for code changes: `go test ./...`.
- Add unit tests for new logic in the same package.
- Avoid tests that require live internet unless explicitly isolated and optional.

## Documentation
- Keep docs and comments in English.
- Update user-facing docs when adding/changing flags, formats, or workflows.

## Scope Control
- Keep PRs focused and atomic.
- Separate refactor from behavior change when possible.
