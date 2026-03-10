# Contributing Guidelines

Thank you for your interest in contributing to our project. Whether it's a bug report, new feature, correction, or additional
documentation, we greatly value feedback and contributions from our community.

Please read through this document before submitting any issues or pull requests to ensure we have all the necessary
information to effectively respond to your bug report or contribution.


## Reporting Bugs/Feature Requests

We welcome you to use the GitHub issue tracker to report bugs or suggest features.

When filing an issue, please check existing open, or recently closed, issues to make sure somebody else hasn't already
reported the issue. Please try to include as much information as you can. Details like these are incredibly useful:

* A reproducible test case or series of steps
* The version of our code being used
* Any modifications you've made relevant to the bug
* Anything unusual about your environment or deployment



## Code Quality Tools

### mapcheck Static Analyzer

This project includes a custom static analyzer called `mapcheck` that enforces safe map access patterns by requiring the comma-ok idiom.

#### Running mapcheck

```bash
# Run mapcheck on the entire codebase
./tools/mapcheck/mapcheck ./...

# Run mapcheck on specific files or directories
./tools/mapcheck/mapcheck internal/ctl/
```

#### What mapcheck detects

The tool identifies unsafe map access patterns and suggests using the comma-ok idiom:

```go
// ❌ Unsafe - may panic if key doesn't exist
value := myMap[key]

// ✅ Safe - recommended patterns
value, ok := myMap[key]  // Check existence
value, _ := myMap[key]   // Ignore existence check
_, ok := myMap[key]      // Only check existence
```

#### CI Integration

mapcheck is automatically run in our CI pipeline. All pull requests must pass mapcheck validation before merging.

#### Building mapcheck

```bash
cd tools/mapcheck
go build -o mapcheck main.go
```

## Finding contributions to work on
Looking at the existing issues is a great way to find something to contribute on. As our projects, by default, use the default GitHub issue labels (enhancement/bug/duplicate/help wanted/invalid/question/wontfix), looking at any 'help wanted' issues is a great place to start.


## Code of Conduct
This project has adopted the [Amazon Open Source Code of Conduct](https://aws.github.io/code-of-conduct).
For more information see the [Code of Conduct FAQ](https://aws.github.io/code-of-conduct-faq) or contact
opensource-codeofconduct@amazon.com with any additional questions or comments.


## Security issue notifications
If you discover a potential security issue in this project we ask that you notify AWS/Amazon Security via our [vulnerability reporting page](http://aws.amazon.com/security/vulnerability-reporting/). Please do **not** create a public github issue.


## Licensing

See the [LICENSE](LICENSE) file for our project's licensing. We will ask you to confirm the licensing of your contribution.
