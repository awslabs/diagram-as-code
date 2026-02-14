# Troubleshooting

Common issues and solutions for diagram-as-code.

## Table of Contents

- [Installation Issues](#installation-issues)
- [YAML Syntax Errors](#yaml-syntax-errors)
- [Resource Type Errors](#resource-type-errors)
- [Link Positioning Issues](#link-positioning-issues)
- [MCP Server Issues](#mcp-server-issues)
- [CloudFormation Conversion Issues](#cloudformation-conversion-issues)

---

## Installation Issues

### Command not found: awsdac

**Problem**: After installation, `awsdac` command is not recognized.

**Solution for brew users**:
```bash
# Verify installation
brew list awsdac

# If not in PATH, add Homebrew bin to PATH
echo 'export PATH="/opt/homebrew/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**Solution for Go users**:
```bash
# Verify GOPATH/bin is in PATH
echo $PATH | grep go/bin

# If not, add to PATH
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Permission denied

**Problem**: Permission error when running `awsdac`.

**Solution**:
```bash
# Make binary executable
chmod +x $(which awsdac)
```

---

## YAML Syntax Errors

### Invalid YAML format

**Problem**: `Error: yaml: line X: mapping values are not allowed in this context`

**Solution**: Check YAML indentation. Use spaces (not tabs) and ensure consistent indentation.

**Example of incorrect YAML**:
```yaml
Resources:
  VPC:
  Type: AWS::EC2::VPC  # ❌ Wrong indentation
```

**Correct YAML**:
```yaml
Resources:
  VPC:
    Type: AWS::EC2::VPC  # ✅ Correct indentation
```

### Missing required fields

**Problem**: `Error: missing required field 'Type'`

**Solution**: Ensure all resources have a `Type` field.

```yaml
Resources:
  Canvas:
    Type: AWS::Diagram::Canvas  # Required
    Children:
      - AWSCloud
```

---

## Resource Type Errors

### Unknown resource type

**Problem**: `Error: unknown resource type 'AWS::EC2::InvalidType'`

**Solution**: Check the [Resource Types documentation](resource-types.md) for valid types. Common types:
- `AWS::EC2::VPC`
- `AWS::EC2::Subnet`
- `AWS::EC2::Instance`
- `AWS::ElasticLoadBalancingV2::LoadBalancer`
- `AWS::S3::Bucket`

### Resource not found in definition file

**Problem**: `Error: resource type not found in definition file`

**Solution**: Ensure your DefinitionFiles section points to the correct definition file:

```yaml
Diagram:
  DefinitionFiles:
    - Type: URL
      Url: https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml
```

---

## Link Positioning Issues

### Links overlapping resources

**Problem**: Links cross over other resources, making the diagram unclear.

**Solution 1**: Use `UnorderedChildren` feature to automatically reorder resources.

```yaml
Resources:
  VPC:
    Type: AWS::EC2::VPC
    Options:
      UnorderedChildren: true
```

See [UnorderedChildren documentation](advanced/unordered-children.md) for details.

**Solution 2**: Manually adjust link positions using 16-wind rose notation (N, NNE, NE, ENE, E, ESE, SE, SSE, S, SSW, SW, WSW, W, WNW, NW, NNW).

```yaml
Links:
  - Source: ALB
    SourcePosition: SSE  # Adjust position
    Target: Instance
    TargetPosition: N
```

### Auto-positioning not working

**Problem**: Links don't automatically find optimal positions.

**Solution**: Ensure you're using `auto` or omitting position parameters:

```yaml
Links:
  - Source: ALB
    SourcePosition: auto  # or omit entirely
    Target: Instance
    TargetPosition: auto  # or omit entirely
```

See [Auto-positioning documentation](advanced/auto-positioning.md) for details.

---

## MCP Server Issues

### MCP Server not starting

**Problem**: MCP Server fails to start or AI assistant can't connect.

**Solution**: Check MCP server configuration in your AI assistant's config file.

**For Cline/Claude Desktop** (`~/Library/Application Support/Claude/claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "awsdac-mcp-server": {
      "command": "/opt/homebrew/bin/awsdac-mcp-server",
      "args": [],
      "type": "stdio"
    }
  }
}
```

**Verify binary path**:
```bash
which awsdac-mcp-server
# Use the output path in your config
```

### MCP Server log file not found

**Problem**: Can't find MCP Server logs for debugging.

**Solution**: Default log location is `/tmp/awsdac-mcp-server.log`.

To use custom log location:
```json
{
  "mcpServers": {
    "awsdac-mcp-server": {
      "command": "/opt/homebrew/bin/awsdac-mcp-server",
      "args": ["--log-file", "/path/to/custom.log"],
      "type": "stdio"
    }
  }
}
```

See [MCP Server documentation](mcp-server.md) for complete setup guide.

---

## CloudFormation Conversion Issues

### Generated diagram is not optimal

**Problem**: CloudFormation template converts but diagram layout is poor.

**Solution**: Use `--dac-file` option to generate editable YAML, then customize:

```bash
# Generate DAC file
awsdac template.yaml --cfn-template --dac-file -o custom.yaml

# Edit custom.yaml to improve layout
# Then generate diagram
awsdac custom.yaml -o improved.png
```

See [CloudFormation Conversion Guide](cloudformation.md) for customization tips.

### Conversion fails with error

**Problem**: `Error: failed to parse CloudFormation template`

**Solution**: Ensure your CloudFormation template is valid YAML. Check for:
- Correct indentation
- Valid CloudFormation resource types
- No syntax errors

**Known limitations**: Some CloudFormation features are not yet supported. See [known issues](https://github.com/awslabs/diagram-as-code/labels/cfn-template%20feature).

---

## Still Having Issues?

If your issue isn't listed here:

1. **Check existing issues**: [GitHub Issues](https://github.com/awslabs/diagram-as-code/issues)
2. **Search documentation**: Use GitHub's search to find relevant docs
3. **Create new issue**: [Report a bug](https://github.com/awslabs/diagram-as-code/issues/new)

When reporting issues, please include:
- `awsdac` version (`awsdac --version`)
- Operating system
- Complete error message
- Minimal YAML example that reproduces the issue
