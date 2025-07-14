package main

import (
	"context"
	"encoding/base64"
	"github.com/spf13/pflag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/awslabs/diagram-as-code/internal/ctl"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
)

// Variables to allow mocking file operations in tests
var writeFileFunc = os.WriteFile
var readFileFunc = os.ReadFile

// ToolName constants for the MCP server tools
type ToolName string

const (
	GENERATE_DIAGRAM ToolName = "generateDiagram"
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
		server.WithInstructions("AWS Diagram-as-Code MCP Server provides tools to generate AWS architecture diagrams from YAML files."),
	)

	// Add the diagram generation tool
	mcpServer.AddTool(mcp.NewTool(string(GENERATE_DIAGRAM),
		mcp.WithDescription("Generates an AWS architecture diagram from a YAML file"),
		mcp.WithString("yamlContent",
			mcp.Description("The YAML content to generate a diagram from"),
			mcp.Required(),
		),
		mcp.WithString("outputFormat",
			mcp.Description("The output format of the diagram (png, svg)"),
			mcp.DefaultString("png"),
		),
	), handleGenerateDiagram)

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

func main() {
	logFilePath := pflag.String("log-file", "awsdac-mcp-server.log", "Path to log file")
	pflag.Parse()

	// Setup logrus to write to file
	logFile, err := os.OpenFile(*logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
