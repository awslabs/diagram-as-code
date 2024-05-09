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

### command
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


Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Desktop (please complete the following information):**
 - OS: [e.g. iOS]
 - Browser [e.g. chrome, safari]
 - Version [e.g. 22]

**Smartphone (please complete the following information):**
 - Device: [e.g. iPhone6]
 - OS: [e.g. iOS8.1]
 - Browser [e.g. stock browser, safari]
 - Version [e.g. 22]

**Additional context**
Add any other context about the problem here.
