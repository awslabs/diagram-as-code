package main

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	"github.com/awslabs/diagram-as-code/internal/ctl"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
)

//go:embed prompts/*
var promptsFS embed.FS

// Variables to allow mocking file operations in tests
var (
	writeFileFunc = os.WriteFile
	readFileFunc  = os.ReadFile
)

// ToolName constants for the MCP server tools
type ToolName string

const (
	GENERATE_DIAGRAM           ToolName = "generateDiagram"
	GET_DIAGRAM_AS_CODE_FORMAT ToolName = "getDiagramAsCodeFormat"
)

// Default prompt template file paths
const (
	USER_REQUIREMENTS_TEMPLATE_FILE = "prompts/generate_dac_from_user_requirements.txt"
)

// Tool descriptions for better maintainability and detailed MCP host guidance
const (
	GENERATE_DIAGRAM_DESC = `Generate AWS architecture diagrams from YAML-based Diagram-as-code specifications.

DESCRIPTION:
This tool creates professional PNG images of AWS architecture diagrams. The input YAML must follow the Diagram-as-code format with three main sections: DefinitionFiles, Resources, and optional Links.

REQUIREMENTS:
- YAML must include 'Diagram:' as root element
- DefinitionFiles section must reference AWS icon definitions
- Resources section must define hierarchical AWS resource structure  
- Canvas resource is required as the root container

FEATURES:
- Supports all major AWS service icons and types
- Handles complex layouts with grouping (VerticalStack/HorizontalStack)
- Creates network connections and relationships via Links
- Generates high-quality PNG output suitable for documentation

USE CASE:
Perfect for creating technical documentation, architecture reviews, and system design presentations.

PREREQUISITE:
Call 'getDiagramAsCodeFormat' first if you need format specification and examples.`

	GET_FORMAT_DESC = `Get comprehensive Diagram-as-code format specification, examples, and best practices.

PURPOSE:
Returns the complete documentation for creating YAML-based AWS architecture diagrams. This includes format specification, resource types, layout techniques, and practical examples.

WHAT YOU GET:
- Complete YAML schema and syntax rules
- Available AWS resource types and their properties
- Layout strategies (VerticalStack, HorizontalStack, grouping)
- Link configuration for showing relationships
- Best practices for creating beautiful, professional diagrams
- Multiple working examples from simple to complex architectures

WHEN TO USE:
- Before creating your first diagram with generateDiagram
- When you need reference for specific resource types or layouts
- For understanding advanced features like orthogonal links and positioning
- When troubleshooting YAML format issues

OUTPUT:
Extensive documentation including format rules, examples, and architectural guidance for creating effective AWS diagrams.

RECOMMENDATION:
Always call this tool first to understand the format before using generateDiagram.`
)

// NewMCPServer creates a new MCP server with the necessary tools and configurations
func NewMCPServer() *server.MCPServer {
	hooks := &server.Hooks{}

	// Add hooks for logging and debugging
	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		log.Infof("beforeAny: %s, %v", method, id)
	})
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		log.Infof("onSuccess: %s, %v", method, id)
	})
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		log.Errorf("onError: %s, %v, %v", method, id, err)
	})

	// Create the MCP server
	mcpServer := server.NewMCPServer(
		"awsdac-mcp-server",
		"0.0.1",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true), // Enable resource capabilities
		server.WithLogging(),
		server.WithHooks(hooks),
		server.WithInstructions(`AWS Diagram-as-Code MCP Server

PURPOSE:
Generate professional AWS architecture diagrams from YAML-based specifications.

RECOMMENDED WORKFLOW:
1. Call 'getDiagramAsCodeFormat' first to understand the format and get examples
2. Use the format guide to create proper YAML content
3. Call 'generateDiagram' with the complete YAML specification
4. Receive a base64-encoded PNG diagram

CAPABILITIES:
- Generate PNG diagrams with AWS resource icons and relationships
- Support hierarchical layouts with Canvas → Cloud → Region → VPC → Subnets → Resources
- Create network connections with Links (straight or orthogonal lines)
- Handle complex layouts using VerticalStack and HorizontalStack groupings

OUTPUT: Base64-encoded PNG images suitable for embedding in responses`),
	)

	// Add the diagram generation tool
	mcpServer.AddTool(mcp.NewTool(string(GENERATE_DIAGRAM),
		mcp.WithDescription(GENERATE_DIAGRAM_DESC),
		mcp.WithString("yamlContent",
			mcp.Description(`Complete YAML specification for the AWS architecture diagram.

REQUIRED STRUCTURE:
Diagram:
  DefinitionFiles:
    - Type: URL
      Url: https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children:
        - AWSCloud
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Children:
        - [your AWS resources]
  Links: # Optional
    - Source: ResourceA
      Target: ResourceB

VALIDATION RULES:
- Must be valid YAML syntax
- Must contain Diagram, DefinitionFiles, and Resources sections
- Canvas must be the root resource with Children
- All resources must be reachable from Canvas
- Link sources and targets must reference existing resources

COMMON RESOURCE TYPES:
- AWS::EC2::VPC, AWS::EC2::Subnet, AWS::EC2::Instance
- AWS::ElasticLoadBalancingV2::LoadBalancer
- AWS::RDS::DBInstance, AWS::S3::Bucket
- AWS::Diagram::VerticalStack, AWS::Diagram::HorizontalStack for grouping`),
			mcp.Required(),
		),
		mcp.WithString("outputFormat",
			mcp.Description(`Output format for the generated diagram.

SUPPORTED FORMATS:
- "png" (default): High-quality PNG image with transparency support

TECHNICAL DETAILS:
- Output is base64-encoded for easy embedding
- Typical resolution: 1200x800 pixels or larger depending on complexity
- Uses official AWS icon set for professional appearance
- Optimized for documentation and presentation use`),
			mcp.DefaultString("png"),
			mcp.Enum("png"), // Future support for other formats
		),
	), handleGenerateDiagram)

	// Add the tool to generate DAC YAML from user requirements
	mcpServer.AddTool(mcp.NewTool(string(GET_DIAGRAM_AS_CODE_FORMAT),
		mcp.WithDescription(GET_FORMAT_DESC),
	), handleGenerateDacFromUserRequirements)

	return mcpServer
}

// handleGenerateDiagram handles the diagram generation tool calls
func handleGenerateDiagram(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	yamlContent, ok := arguments["yamlContent"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid yamlContent argument")
	}

	outputFormat, _ := arguments["outputFormat"].(string)
	if outputFormat == "" {
		outputFormat = "png"
	}

	// Create a temporary directory for processing
	tempDir, err := os.MkdirTemp("", "awsdac-mcp")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create temporary input file
	inputFile := filepath.Join(tempDir, "input.yaml")
	if err := writeFileFunc(inputFile, []byte(yamlContent), 0o644); err != nil {
		return nil, fmt.Errorf("failed to write input file: %v", err)
	}

	// Create output file path
	outputFile := filepath.Join(tempDir, "output.png")

	// Generate diagram directly in the main thread
	// This ensures logs are properly captured
	opts := &ctl.CreateOptions{}
	ctl.CreateDiagramFromDacFile(inputFile, &outputFile, opts)

	// Read the generated diagram
	diagramData, err := readFileFunc(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read generated diagram: %v", err)
	}

	// Encode the diagram as base64
	base64Diagram := base64.StdEncoding.EncodeToString(diagramData)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "Diagram generated successfully",
			},
			mcp.ImageContent{
				Type:     "image",
				Data:     base64Diagram,
				MIMEType: "image/png",
			},
		},
	}, nil
}

// handleGenerateDacFromUserRequirements handles the generation of DAC YAML from user requirements
func handleGenerateDacFromUserRequirements(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	templateContent, err := readPromptFile(USER_REQUIREMENTS_TEMPLATE_FILE)
	if err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(templateContent),
			},
		},
	}, nil
}

// readPromptFile reads a prompt file from the embedded filesystem
func readPromptFile(filePath string) ([]byte, error) {
	content, err := promptsFS.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded prompt file %s: %v", filePath, err)
	}
	return content, nil
}

func main() {
	logFilePath := pflag.String("log-file", "", "Path to log file")
	pflag.Parse()

	// Determine log file path
	var actualLogPath string
	if *logFilePath != "" {
		actualLogPath = *logFilePath
	} else {
		actualLogPath = filepath.Join(os.TempDir(), "awsdac-mcp-server.log")
	}

	// Setup logrus to write to file
	logFile, err := os.OpenFile(actualLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	log.Info("Starting MCP server with stdio transport")
	mcpServer := NewMCPServer()

	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
