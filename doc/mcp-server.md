# MCP Server

Integration guide for using diagram-as-code with AI assistants through the Model Context Protocol (MCP).

## Overview

The awsdac MCP server enables AI assistants and development tools to generate AWS architecture diagrams programmatically. This integration allows seamless diagram creation within your development workflow without leaving your AI assistant interface.

## Prerequisites

- awsdac installed (`brew install awsdac` or `go install`)
- MCP-compatible AI assistant (Claude Desktop, Cline, etc.)
- Basic understanding of diagram-as-code YAML format

## Installation

### For macOS users (Homebrew)

```bash
brew install awsdac
```

The MCP server binary (`awsdac-mcp-server`) is included with the main installation.

### For Go developers

```bash
go install github.com/awslabs/diagram-as-code/cmd/awsdac-mcp-server@latest
```

Verify installation:
```bash
which awsdac-mcp-server
# Output: /opt/homebrew/bin/awsdac-mcp-server (or similar)
```

## Configuration

### Claude Desktop / Cline

Edit your MCP configuration file:

**Location**: `~/Library/Application Support/Claude/claude_desktop_config.json`

**Configuration**:
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

**With custom log file**:
```json
{
  "mcpServers": {
    "awsdac-mcp-server": {
      "command": "/opt/homebrew/bin/awsdac-mcp-server",
      "args": ["--log-file", "/path/to/custom/awsdac-mcp.log"],
      "type": "stdio"
    }
  }
}
```

**Note**: Replace `/opt/homebrew/bin/awsdac-mcp-server` with your actual binary path from `which awsdac-mcp-server`.

### Finding Your Binary Location

```bash
# Check if awsdac-mcp-server is in your PATH
which awsdac-mcp-server

# Or check common Go install locations
ls ~/go/bin/awsdac-mcp-server
ls $GOPATH/bin/awsdac-mcp-server  # if GOPATH is set
```

## Available Tools

The MCP server provides three tools for diagram generation:

### 1. generateDiagram

Generates AWS architecture diagrams from YAML specifications and returns base64-encoded PNG images.

**Use when**: Your AI assistant can display base64-encoded images directly.

**Parameters**:
- `yamlContent` (required): Complete YAML specification following the diagram-as-code format
- `outputFormat` (optional): Output format, default is "png"

### 2. generateDiagramToFile

Same as `generateDiagram` but saves the image directly to a specified file path.

**Use when**: Your AI assistant cannot receive large base64 image data, or you want to save directly to a file.

**Parameters**:
- `yamlContent` (required): Complete YAML specification
- `outputFilePath` (required): Path where the PNG file should be saved

### 3. getDiagramAsCodeFormat

Returns comprehensive format specification, examples, and best practices for creating diagram-as-code YAML files.

**Use when**: You need reference documentation or examples for YAML format.

## Usage Examples

### Example 1: Generate a Simple Diagram

**Prompt to AI assistant**:
```
Generate an AWS architecture diagram showing a VPC with a public subnet containing an EC2 instance.
```

The AI assistant will use the MCP server to generate the diagram automatically.

### Example 2: Save Diagram to File

**Prompt to AI assistant**:
```
Generate an AWS architecture diagram showing an ALB with two EC2 instances, 
and save it to /tmp/my-architecture.png
```

### Example 3: Get Format Documentation

**Prompt to AI assistant**:
```
Show me the diagram-as-code YAML format documentation
```

## Troubleshooting

### MCP Server Not Starting

**Problem**: AI assistant shows "MCP server connection failed"

**Solutions**:
1. Verify binary path in configuration matches `which awsdac-mcp-server` output
2. Check binary has execution permissions: `chmod +x $(which awsdac-mcp-server)`
3. Restart your AI assistant application

### Log File Not Found

**Problem**: Can't find MCP server logs for debugging

**Solution**: Default log location is `/tmp/awsdac-mcp-server.log`

View logs:
```bash
tail -f /tmp/awsdac-mcp-server.log
```

### Connection Issues

**Problem**: AI assistant can't communicate with MCP server

**Solutions**:
1. Verify JSON configuration syntax is valid
2. Check that `type` is set to `"stdio"`
3. Ensure no other MCP server is using the same name
4. Review log files for error messages

## Related Documentation

- **[Introduction Guide](introduction.md)** - Learn diagram-as-code basics
- **[Resource Types](resource-types.md)** - Available AWS resources
- **[Links](links.md)** - Connecting resources
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions
