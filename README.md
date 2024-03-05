# Diagram-as-code
This command line interface (CLI) tool enables drawing infrastructure diagrams for Amazon Web Services through YAML code. It facilitates diagram-as-code without relying on image libraries.

The CLI tool promotes code reuse, testing, integration, and automating the diagramming process. It allows managing diagrams with Git by writing human-readable YAML.

![Example diagram](doc/static/example.png)

(generated from [example of NAT Gateway](examples/vpc-natgw.yaml))

## Features
- **Follow AWS architecture guidelines.**  
Easily generate diagrams that follow [AWS diagram guidelines](https://aws.amazon.com/architecture/icons).
- **Scalability**  
Flexible layout allows you to adjust the position and size of groups automatically.
- **Lightweight & Container-friends**  
Experience the comfort of running on a scratch environment without a headless browser and GUI. It has small dependencies, making it suitable for generating large diagrams.
- **Integrate with your Infrastructure as Code**  
Generate diagrams to align with your IaC code without managing diagrams manually.
- **As an AWS diagram engine**  
If you want to integrate with other IaC tools, AI, or drawing GUI tools, you can use this tool as a Golang library.
- **Expandable**  
All resources are defined in files; you can integrate them with other cloud, on-premises diagrams.

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
