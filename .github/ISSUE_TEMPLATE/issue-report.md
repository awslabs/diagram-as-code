---
name: Issue report
about: Create a report to help us improve
title: ''
labels: ''
assignees: ''

---

## Issue detail
A clear and concise description of what the bug is.

## Reproduce process

### command / error messages
```
$ awsdac ...
```

### Yaml template
```
Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - AWSCloud
        - User
...
```
