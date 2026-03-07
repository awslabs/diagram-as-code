# Development Model

## Branch Strategy
Follow the workflow in [doc/gitflow.md](../gitflow.md).

Default flow:
1. Branch from `develop`:
   - `feature/<scope>`
   - `fix/<scope>`
2. Commit with explicit intent (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`).
3. Run tests (`go test ./...`).
4. Open PR to `develop`.
5. Merge to `main` through release/hotfix flow.

## Coding Patterns
- Keep orchestration in `internal/ctl`.
- Keep data/geometry logic in `internal/types` and `internal/vector`.
- Favor pure helper functions for testability.
- Keep logging informative but not noisy.

## Change Checklist
- Behavior change explained in PR description
- Tests added/updated
- Docs updated (`README` + related docs)
- No secret files introduced

## AI Agent Checklist
- Read `AGENTS.md` first.
- Confirm affected flow (PNG, draw.io, CFn).
- Implement the smallest safe patch.
- Verify with `go test ./...`.
