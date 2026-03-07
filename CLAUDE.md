# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./...

# Run all tests (unit + functional pixel comparison)
go test ./...

# Run a single package's tests
go test ./internal/ctl/...
go test ./internal/types/...

# Run a specific test
go test ./internal/ctl/... -run TestRGBAToHex

# Run the functional integration tests (generates PNGs to /tmp/results/)
go test ./test/...

# Run the CLI directly
go run ./cmd/awsdac examples/alb-ec2.yaml
go run ./cmd/awsdac examples/alb-ec2.yaml --drawio -o output.drawio -f
go run ./cmd/awsdac examples/tgw-nwfw-tmpl.yaml -t -f          # Go template
go run ./cmd/awsdac examples/vpc-subnet-ec2-cfn.yaml -c -f     # CloudFormation
```

## Architecture

### Request-to-Output Pipeline (shared by PNG and draw.io)

All flows go through `internal/ctl/`:

1. **Parse** YAML input (`dacfile.go`: `getTemplate` + `processTemplate`)
2. **Load definitions** from URL/local/embedded (`create.go`: `loadDefinitionFiles`)
3. **Build runtime resources** — each YAML resource becomes a `*types.Resource` with its icon, label, colors, and geometry (`create.go`: `loadResources`)
4. **Associate children & links** (`create.go`: `associateChildren`, `loadLinks`)
5. **Layout** — `canvas.Scale()` recursively computes bounds bottom-up; `ZeroAdjust()` normalizes to origin
6. **Export** — PNG via `createDiagram()`, draw.io via `exportToDrawio()`

The public entrypoints are:

- `ctl.CreateDiagramFromDacFile()` — DAC YAML → PNG
- `ctl.CreateDrawioFromDacFile()` — DAC YAML → draw.io
- `ctl.CreateDiagramFromCFnTemplate()` — CloudFormation → PNG or DAC YAML

### Key Type: `types.Resource`

Every node in the diagram is a `*types.Resource` (`internal/types/resource.go`). It holds:

- `label` (string) — display name, set via `SetLabel()` / read via `GetLabel()`
- `bindings` (`*image.Rectangle`) — pixel bounds computed by the layout engine
- `children` / `borderChildren` — tree structure
- `links` — connections to other resources

Layout is driven by `Scale()` (recursive, child-first) and the concrete subtypes `VerticalStack` / `HorizontalStack`.

### Draw.io Export (`internal/ctl/drawio.go`)

- Runs the same layout engine as PNG to get identical coordinates, then emits `mxGraphModel` XML
- **Label resolution**: YAML `Title` → runtime label from definition (`GetLabel()`) → YAML key name
- **Edge labels**: collected from `Links[].Labels` (SourceLeft, SourceRight, TargetLeft, TargetRight, AutoLeft, AutoRight), joined with ` | `
- **Type → style mapping**: `dacTypeStyles` map covers icons and group containers. Use `groupStyle()` for solid borders and `groupStyleDashed()` for dashed borders (AutoScalingGroup, Region)
- **Icons**: `drawio_assets.go` downloads the official AWS Asset Package ZIP once (cached), extracts SVGs, and inlines them as base64 data URIs

### Definition Files

`definitions/definition-for-aws-icons-light.yaml` is the local copy of the AWS icon definition file. It maps AWS resource types to icon paths, label defaults, fill/border colors, and border styles (`straight`/`dashed`). The tool also fetches the same file from GitHub at runtime unless `--override-def-file` is used.

### Functional Tests (Golden Files)

`test/func_test.go` generates PNGs from all `examples/*.yaml` into `/tmp/results/` and does a pixel-by-pixel comparison against the committed `examples/*.png` golden files.

**If layout changes cause pixel differences**, update the golden files:

```bash
go test ./test/...          # generates fresh PNGs to /tmp/results/
cp /tmp/results/*.png examples/
go test ./test/...          # must pass now
```

## YAML Diagram Structure

```yaml
Diagram:
  DefinitionFiles:
    - Type: URL        # or LocalFile / Embed
      Url: "https://..."
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical   # or horizontal
      Children: [AWSCloud]
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Preset: AWSCloudNoLogo
      Children: [VPC]
    VPC:
      Type: AWS::EC2::VPC
      Children: [Subnet]
      BorderChildren:
        - Position: S        # 16-wind rose: N, NNE, NE, ENE, E, ...
          Resource: IGW
    Subnet:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet   # or PrivateSubnet
      Children: [Instance]
    Instance:
      Type: AWS::EC2::Instance
  Links:
    - Source: ALB
      SourcePosition: NNW    # or "auto"
      Target: Instance
      TargetPosition: S
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "HTTP:80"
```

## Branch Strategy

- `develop` — integration; all feature/fix work targets here
- `main` — production-ready; merge from `develop` only
- Commit prefixes: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`
