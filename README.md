# Diagram-as-code
This command line interface (CLI) tool enables drawing infrastructure diagrams for Amazon Web Services through YAML code. It facilitates diagram-as-code without relying on image libraries.

The CLI tool promotes code reuse, testing, integration, and automating the diagramming process. It allows managing diagrams with Git by writing human-readable YAML.

Example templates are [here](examples).
Check out the [Introduction Guide](doc/introduction.md) as well for additional information.

<img src="doc/static/introduction2.png" width="800">

![CLI Usage animation](doc/static/command_demo.gif)

## Features
- **Compliant with AWS architecture guidelines**  
Easily generate diagrams that follow [AWS diagram guidelines](https://aws.amazon.com/architecture/icons).
- **Flexible**  
Automatically adjust the position and size of groups.
- **Lightweight & CI/CD-friendly**  
Start quickly on a container; no dependency on headless browser or GUI.
- **Integrate with your Infrastructure as Code**  
Generate diagrams to align with your IaC code without managing diagrams manually.
- **As a drawing library**  
Use as Golang Library and integrate with other IaC tools, AI, or drawing GUI tools.
- **Extensible**  
Add definition files to create non-AWS diagrams as well.

## About This Fork

This is a fork of [awslabs/diagram-as-code](https://github.com/awslabs/diagram-as-code) maintained by **Fernando Azevedo** — [fernando.moretes.com](https://fernando.moretes.com/).

Fork repository: **[github.com/fernandofatech/diagram-as-code](https://github.com/fernandofatech/diagram-as-code)**

This fork extends the original CLI tool with new features and a fully hosted web frontend:

### Web Frontend (New)

A browser-based editor that lets you write YAML and generate diagrams in PNG or draw.io format — no CLI or local install required. Hosted 100% on Vercel.

<img src="doc/static/web-frontend.png" width="900" alt="Web frontend screenshot — YAML editor with live diagram preview">

**Try it live:** [dac.moretes.com](https://dac.moretes.com)

Key features:
- Monaco-powered YAML editor with syntax highlighting
- PNG and draw.io format toggle
- Built-in example templates (ALB+EC2, VPC+NAT, ALB+AutoScaling, Multi-Region)
- One-click download for generated diagrams
- Go serverless backend via Vercel Functions

### CLI Extensions (vs upstream)

- Native `.drawio` export via `--drawio` flag or `-o output.drawio`
- Draw.io export pipeline in Go (`internal/ctl/drawio.go`, `drawio_assets.go`)
- Public wrapper package `pkg/diagram` for embedding in external tools
- `cmd/api-dev` local dev server for testing the API handler without Vercel CLI

## Getting started

### macOS
```bash
brew install awsdac
```

### Linux
```bash
go install github.com/fernandofatech/diagram-as-code/cmd/awsdac@latest
# or download the binary from GitHub Releases and move to /usr/local/bin/
```

### Windows
```powershell
go install github.com/fernandofatech/diagram-as-code/cmd/awsdac@latest
# or download awsdac.exe from GitHub Releases and add to your PATH
```

### From source (all platforms)
```bash
git clone https://github.com/fernandofatech/diagram-as-code.git
cd diagram-as-code
make build   # produces ./awsdac (or awsdac.exe on Windows)
```

**Prerequisites:** Go 1.21+ · Node.js 18+ (only for the web frontend)

## Usage

```
Usage:
  awsdac <input filename> [flags]

Flags:
      --allow-untrusted-definitions  Allow loading definition files from untrusted URLs (not from official repository)
  -c, --cfn-template               [beta] Create diagram from CloudFormation template
  -d, --dac-file                   [beta] Generate YAML file in dac (diagram-as-code) format from CloudFormation template
      --drawio                     Generate draw.io (.drawio) file instead of PNG
  -f, --force                      Overwrite output file without confirmation
  -h, --help                       help for awsdac
      --height int                 Resize output image height (0 means no resizing)
  -o, --output string              Output file name (default "output.png")
      --override-def-file string   For testing purpose, override DefinitionFiles to another url/local file
  -t, --template                   Processes the input file as a template according to text/template.
  -v, --verbose                    Enable verbose logging
      --version                    version for awsdac
      --width int                  Resize output image width (0 means no resizing)
```

### Example

```
$ awsdac examples/alb-ec2.yaml
```

```
$ awsdac privatelink.yaml -o custom-output.png
```

```
$ awsdac examples/alb-ec2.yaml --drawio -o output.drawio
```

```
$ awsdac examples/alb-ec2.yaml -o output.drawio
```

## How Draw.io Export Works
1. The YAML file is parsed with the same resource/link model used by PNG rendering.
2. The same layout engine is executed (`Scale` and `ZeroAdjust`) to keep geometry consistent.
3. Children are reordered from link topology to keep ordering aligned with PNG output.
4. Resources are exported as draw.io cells and links are exported as draw.io edges in `mxGraphModel`.
5. Official AWS SVG icons are loaded from the AWS Asset Package and embedded as data URIs for leaf resources.

## Documentation

### Getting Started
- **[Introduction Guide](doc/introduction.md)** - Quick start (10 minutes) and core concepts
- **[Troubleshooting](doc/troubleshooting.md)** - Common issues and solutions

### Core Features
- **[Resource Types](doc/resource-types.md)** - Available AWS resources and diagram elements
- **[Links](doc/links.md)** - Connecting resources with arrows and lines

### Tools & Integration
- **[MCP Server](doc/mcp-server.md)** - AI assistant integration
- **[CloudFormation Conversion](doc/cloudformation.md)** [Beta] - Convert CloudFormation templates to diagrams

### Advanced Features
- **[Templates](doc/template.md)** - Using Go templates for dynamic diagrams
- **[UnorderedChildren](doc/advanced/unordered-children.md)** - Automatic child reordering for optimal layouts
- **[Auto-positioning](doc/advanced/auto-positioning.md)** - Smart link positioning
- **[Link Grouping Offset](doc/advanced/link-grouping.md)** - Prevent link overlap
- **[BorderChildren](doc/advanced/border-children.md)** - Place resources on borders

### Guides
- **[Best Practices](doc/best-practices.md)** - Design patterns and diagram standards
- **[Gitflow Workflow](doc/gitflow.md)** - Branching model for PRs, tests, review, and release flow
- **[AI Documentation](doc/ai/README.md)** - Context, architecture, guardrails, and development model for AI agents
- **[Agent Entry Point](AGENTS.md)** - Primary project instructions for coding agents

### Contributing
- **[Documentation Guidelines](doc/contributing-docs.md)** - How to contribute to documentation

---

## Development Guide

For contributing guidelines, please see [CONTRIBUTING.md](CONTRIBUTING.md).

### Project Structure
- `cmd/` - CLI tools (awsdac, awsdac-mcp-server, api-dev)
- `internal/` - Core implementation
  - `cache/` - Caching logic
  - `ctl/` - Core control logic
  - `definition/` - Definition file handling
  - `font/` - Font management
  - `types/` - Core types and structures
  - `vector/` - Vector operations
- `pkg/diagram/` - Public wrapper around `internal/ctl` (importable by external packages)
- `api/` - Vercel serverless handler (`POST /api/generate`)
- `web/` - Next.js 14 frontend (Monaco editor + diagram preview)
- `test/` - Integration tests
- `tools/` - Development tools
- `examples/` - Example YAML files
- `doc/` - Documentation

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This project is licensed under the Apache-2.0 License.
