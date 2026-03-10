# Contributing Guidelines

Thank you for your interest in contributing to diagram-as-code! Whether it's a bug report, new feature, documentation improvement, or code contribution — all feedback and contributions are welcome.

Please read through this document before submitting issues or pull requests.

## Reporting Bugs / Feature Requests

Use the [GitHub issue tracker](https://github.com/fernandofatech/diagram-as-code/issues) to report bugs or suggest features.

Before filing, check existing open or recently closed issues to avoid duplicates. Include as much detail as possible:

- A reproducible test case or series of steps
- The version (`awsdac --version`)
- Your OS and architecture
- Any modifications you've made relevant to the bug
- Anything unusual about your environment

## Development Setup

```bash
git clone https://github.com/fernandofatech/diagram-as-code.git
cd diagram-as-code
make build       # compile binary
make test        # run all tests
make dev         # start web frontend (port 3001)
```

**Prerequisites**: Go 1.21+ · Node.js 18+ (web frontend only)

## Code Quality Tools

### mapcheck Static Analyzer

This project includes a custom static analyzer called `mapcheck` that enforces safe map access patterns by requiring the comma-ok idiom.

```bash
# Run mapcheck on the entire codebase
./tools/mapcheck/mapcheck ./...

# Run on specific directories
./tools/mapcheck/mapcheck internal/ctl/
```

**What it detects** — unsafe map access patterns:
```go
// ❌ Unsafe - may panic if key doesn't exist
value := myMap[key]

// ✅ Safe - recommended patterns
value, ok := myMap[key]  // check existence
value, _ := myMap[key]   // ignore check
_, ok := myMap[key]      // only check existence
```

mapcheck runs automatically in CI. All PRs must pass before merging.

**Build mapcheck**:
```bash
cd tools/mapcheck
go build -o mapcheck main.go
```

## Branching Model

This project follows a lightweight Gitflow model. See [doc/gitflow.md](doc/gitflow.md) for the full workflow.

**Short version**:
- Branch from `develop` for features/fixes
- Open PRs against `develop`
- `main` is production-ready and released via tags

## PR Checklist

- All tests pass locally (`go test ./...`)
- mapcheck passes (`./tools/mapcheck/mapcheck ./...`)
- Documentation updated when behavior changes
- Commit messages use conventional format (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`)
- If layout changes: update golden PNG files (`cp /tmp/results/*.png examples/`)

## npm Package

The CLI is distributed on npm as [`awsdac`](https://www.npmjs.com/package/awsdac). The npm package is a thin wrapper that downloads the correct pre-built binary for the user's platform during `npm install`.

If you change CLI behavior or flags, update:
- `README.md`
- `doc/introduction.md`
- Relevant docs in `doc/`

## Finding Issues to Work On

Look at [open issues](https://github.com/fernandofatech/diagram-as-code/issues) — especially those labeled `help wanted` or `good first issue`.

## Code of Conduct

This project follows the [Amazon Open Source Code of Conduct](https://aws.github.io/code-of-conduct).
See the [Code of Conduct FAQ](https://aws.github.io/code-of-conduct-faq) for more information.

## Security

If you discover a potential security issue, please **do not** create a public GitHub issue. Report it via [GitHub's private vulnerability reporting](https://github.com/fernandofatech/diagram-as-code/security/advisories/new).

## Licensing

This project is licensed under Apache-2.0. See the [LICENSE](LICENSE) file. By contributing, you agree that your contributions will be licensed under the same license.
