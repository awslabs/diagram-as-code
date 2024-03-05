# Diagram-as-code
This command line interface (CLI) tool enables drawing infrastructure diagrams for Amazon Web Services through YAML code. It facilitates diagram-as-code without relying on image libraries.

The CLI tool promotes code reuse, testing, integration, and automating the diagramming process. It allows managing diagrams with Git by writing human-readable YAML.

![Example diagram](doc/static/example.png)

(generated from [example of NAT Gateway](examples/vpc-natgw.yaml))

## Getting started
### for macOS user
In preparing.

### for Gopher
```
# go install github.com/awslabs/diagram-as-code/cmd/awsdac@latest
```

## Resource Types
See [doc/resource-types.md](doc/resource-types.md).

## Example Usage
See [examples templates](examples).

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This project is licensed under the Apache-2.0 License.
