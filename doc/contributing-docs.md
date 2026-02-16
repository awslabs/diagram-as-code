# Contributing to Documentation

Guidelines for contributing to diagram-as-code documentation.

## Documentation Standards

### Writing Style
- Use clear, concise language
- Write for non-technical stakeholders when possible
- Include code examples for technical concepts
- Add version badges where applicable: `[Beta]`, `[v1.2+]`

### Structure
Every documentation file should include:
1. Title and brief overview
2. Table of contents (for long documents)
3. Main content with examples
4. Related documentation links

## Review Process

All documentation changes must be reviewed by the documentation review board before merge.

### Submitting Changes

1. Create a feature branch
2. Make your documentation changes
3. Test all code examples
4. Validate all internal links
5. Submit pull request
6. Address review feedback

### Review Checklist

Reviewers will check:
- [ ] Accuracy: Technical content is correct
- [ ] Clarity: Easy to understand for target audience
- [ ] Completeness: Covers all aspects of the topic
- [ ] Consistency: Follows style guide and conventions
- [ ] Cross-references: Related docs are linked

## Style Guide

### Headings
- Use sentence case: "Getting started" not "Getting Started"
- Be descriptive: "Install on macOS" not "Installation"

### Code Blocks
Always specify language for syntax highlighting:

````markdown
```yaml
Resources:
  VPC:
    Type: AWS::EC2::VPC
```
````

### Links
Use relative links for internal documentation:
```markdown
See [Resource Types](resource-types.md) for details.
```

### Version Badges
Indicate version-specific features:
```markdown
## CloudFormation Conversion [Beta]
```

## How to Contribute

1. **Find an issue**: Check [documentation issues](https://github.com/awslabs/diagram-as-code/labels/documentation)
2. **Propose changes**: Open an issue to discuss major changes
3. **Submit PR**: Follow the review process above

## Questions?

- Open an issue on GitHub
- Check existing documentation
- Ask in pull request comments
