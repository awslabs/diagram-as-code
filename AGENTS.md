# AGENTS.md

This file is the primary entrypoint for AI/code agents working in this repository.

## Mission
Maintain and evolve `diagram-as-code` as a reliable CLI and library to generate AWS architecture diagrams from YAML, CloudFormation, and templates.

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

## Fast Workflow
1. Implement minimal change.
2. Run validation:
   - `go test ./...`
3. Update docs when behavior changes.

## Output Quality Bar
- No silent behavior changes.
- No hidden breaking changes in CLI flags or output format.
- New logic should include tests in the closest package.
