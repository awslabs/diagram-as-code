# Architecture Model

## High-Level Layers
- `cmd/`
  - CLI entrypoints (`awsdac`, `awsdac-mcp-server`)
- `internal/ctl/`
  - Orchestration pipeline: parse input, load definitions/resources/links, run layout, render/export
- `internal/types/`
  - Runtime data model for resources, links, geometry, and rendering behavior
- `internal/definition/`
  - Definition resolution and icon metadata loading
- `internal/cache/`
  - Remote/local resource cache
- `internal/vector/`
  - Geometric utilities

## Request-to-Output Flow
1. Parse input and options.
2. Load and validate definitions.
3. Build runtime resources.
4. Associate children and links.
5. Run layout and positioning.
6. Export:
   - PNG flow
   - draw.io flow

## Draw.io Flow
- Implemented mainly in `internal/ctl/drawio.go`.
- Uses AWS icon asset package via `internal/ctl/drawio_assets.go`.
- Emits `mxGraphModel` XML.

## Design Principles
- Shared model, specialized exporters.
- Deterministic layout whenever possible.
- Keep IO/network side effects isolated.
