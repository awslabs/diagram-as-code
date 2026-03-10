# AGENTS.md

This file is the primary entrypoint for AI/code agents working in this repository.

## Mission
Maintain and evolve `diagram-as-code` as a reliable CLI and library to generate AWS architecture diagrams from YAML, CloudFormation, and templates — distributed via npm, Homebrew, and GitHub Releases.

## Read First
1. [README.md](README.md)
2. [doc/ai/context.md](doc/ai/context.md)
3. [doc/ai/architecture.md](doc/ai/architecture.md)
4. [doc/ai/guardrails.md](doc/ai/guardrails.md)
5. [doc/ai/development-model.md](doc/ai/development-model.md)

## Non-Negotiable Rules
- Keep docs and comments in English.
- Preserve compatibility of PNG generation unless a change is explicitly intended and validated.
- Do not add network dependence in tests unless mocked.
- Keep changes small, scoped, and test-backed.
- When changing CLI flags or output format, update `README.md` and `doc/introduction.md`.

## Distribution
The CLI is published to npm as [`awsdac`](https://www.npmjs.com/package/awsdac).
- `package.json` at root contains the npm package definition
- `bin/awsdac` is the Node.js wrapper script
- `scripts/npm-install.js` downloads the correct binary on `npm install`
- The GitHub Actions release workflow (`.github/workflows/release.yml`) publishes to npm automatically on every `v*` tag

## Fast Workflow
1. Implement minimal change.
2. Run validation:
   - `go test ./...`
3. Update docs when behavior changes.
4. If releasing: update `package.json` version to match the tag before merging to `main`.

## Output Quality Bar
- No silent behavior changes.
- No hidden breaking changes in CLI flags or output format.
- New logic should include tests in the closest package.
