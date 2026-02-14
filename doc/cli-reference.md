# CLI Reference

Command-line options and environment variables for awsdac.

## Command Syntax

```bash
awsdac <input-file> [flags]
```

## Flags

### `-o, --output <filename>`
Specify output file name.

**Default**: `output.png`

**Example**:
```bash
awsdac diagram.yaml -o my-diagram.png
```

### `-c, --cfn-template`
Create diagram from CloudFormation template. [Beta]

**Example**:
```bash
awsdac template.yaml --cfn-template
```

### `-d, --dac-file`
Generate YAML file in DAC format from CloudFormation template. [Beta]

**Must be used with** `--cfn-template`

**Example**:
```bash
awsdac template.yaml --cfn-template --dac-file
```

### `-t, --template`
Process input file as Go template.

**Example**:
```bash
awsdac diagram-template.yaml --template
```

### `-v, --verbose`
Enable verbose logging.

**Example**:
```bash
awsdac diagram.yaml --verbose
```

### `--version`
Display version information.

**Example**:
```bash
awsdac --version
```

### `-h, --help`
Display help information.

**Example**:
```bash
awsdac --help
```

### `--override-def-file <path>`
Override DefinitionFiles to another URL or local file (for testing).

**Example**:
```bash
awsdac diagram.yaml --override-def-file ./custom-definitions.yaml
```

## Environment Variables

Currently, awsdac does not use environment variables for configuration.

## Exit Codes

- `0`: Success
- `1`: Error (check error message for details)

## Related Documentation

- [Introduction Guide](introduction.md)
- [Templates](template.md)
- [CloudFormation Conversion](cloudformation.md)
