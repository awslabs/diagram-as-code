package main

import (
	"context"
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewMCPServer(t *testing.T) {
	server := NewMCPServer()
	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}
}

func TestHandleGenerateDiagram(t *testing.T) {
	ctx := context.Background()
	validYAML := `
Resources:
  - Type: AWS::EC2::VPC
    Name: MyVPC
    Properties:
      CidrBlock: 10.0.0.0/16
`

	tests := []struct {
		name        string
		arguments   map[string]interface{}
		wantErr     bool
		wantMIME    string
		// wantContent: Expected number of elements in result.Content array
		// Success case returns 2 elements:
		//   1. mcp.TextContent - "Diagram generated successfully" message
		//   2. mcp.ImageContent - base64 encoded diagram image data
		wantContent int
	}{
		{
			name: "valid yaml with png format",
			arguments: map[string]interface{}{
				"yamlContent":  validYAML,
				"outputFormat": "png",
			},
			wantErr:     false,
			wantMIME:    "image/png",
			wantContent: 2,
		},
		{
			name: "valid yaml with default format",
			arguments: map[string]interface{}{
				"yamlContent": validYAML,
			},
			wantErr:     false,
			wantMIME:    "image/png",
			wantContent: 2,
		},
		{
			name: "missing yamlContent",
			arguments: map[string]interface{}{
				"outputFormat": "png",
			},
			wantErr: true,
		},
		{
			name: "invalid yamlContent type",
			arguments: map[string]interface{}{
				"yamlContent": 123,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      string(GENERATE_DIAGRAM),
					Arguments: tt.arguments,
				},
			}

			result, err := handleGenerateDiagram(ctx, request)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if result == nil {
				t.Fatal("Expected result, got nil")
			}

			if len(result.Content) != tt.wantContent {
				t.Fatalf("Expected %d content items, got %d", tt.wantContent, len(result.Content))
			}

			// Check text content
			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatal("Expected first content to be TextContent")
			}
			if textContent.Text != "Diagram generated successfully" {
				t.Errorf("Expected success message, got: %s", textContent.Text)
			}

			// Check image content
			imageContent, ok := result.Content[1].(mcp.ImageContent)
			if !ok {
				t.Fatal("Expected second content to be ImageContent")
			}
			if imageContent.MIMEType != tt.wantMIME {
				t.Errorf("Expected %s MIME type, got: %s", tt.wantMIME, imageContent.MIMEType)
			}

			// Verify base64 data is valid
			_, err = base64.StdEncoding.DecodeString(imageContent.Data)
			if err != nil {
				t.Errorf("Expected valid base64 data, got decode error: %v", err)
			}
		})
	}
}

func TestHandleGenerateDiagram_TempFileContent(t *testing.T) {
	ctx := context.Background()
	yamlContent := `Resources:
  - Type: AWS::EC2::VPC
    Name: TestVPC`

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: string(GENERATE_DIAGRAM),
			Arguments: map[string]interface{}{
				"yamlContent": yamlContent,
			},
		},
	}

	// Mock the file operations to capture the temp file content
	originalWriteFile := writeFileFunc
	var capturedContent []byte
	writeFileFunc = func(filename string, data []byte, perm os.FileMode) error {
		if strings.HasSuffix(filename, "input.yaml") {
			capturedContent = data
		}
		return originalWriteFile(filename, data, perm)
	}
	defer func() { writeFileFunc = originalWriteFile }()

	_, err := handleGenerateDiagram(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if string(capturedContent) != yamlContent {
		t.Errorf("Expected temp file content %q, got %q", yamlContent, string(capturedContent))
	}
}

func TestHandleGenerateDiagram_OutputFileContent(t *testing.T) {
	ctx := context.Background()
	yamlContent := `Resources:
  - Type: AWS::EC2::VPC
    Name: TestVPC`
	// Expected content generation approach:
	// 1. Create mock data instead of actual PNG image data for test stability
	// 2. Use predictable mock data to ensure consistent test results
	// 3. Bypass actual diagram generation process to focus on file I/O testing
	mockDiagramData := []byte("mock-png-data")

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: string(GENERATE_DIAGRAM),
			Arguments: map[string]interface{}{
				"yamlContent": yamlContent,
			},
		},
	}

	// Mock the file operations
	// This mock function returns our predefined mock data instead of reading actual PNG file
	originalReadFile := readFileFunc
	var readFilePath string
	readFileFunc = func(filename string) ([]byte, error) {
		if strings.HasSuffix(filename, "output.png") {
			readFilePath = filename
			return mockDiagramData, nil // Return mock data for predictable testing
		}
		return originalReadFile(filename)
	}
	defer func() { readFileFunc = originalReadFile }()

	result, err := handleGenerateDiagram(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify output file was read
	if !strings.HasSuffix(readFilePath, "output.png") {
		t.Errorf("Expected to read output.png file, got: %s", readFilePath)
	}

	// Verify the diagram data in the result
	// Expected value is generated by base64 encoding the same mock data
	imageContent := result.Content[1].(mcp.ImageContent)
	expectedBase64 := base64.StdEncoding.EncodeToString(mockDiagramData)
	if imageContent.Data != expectedBase64 {
		t.Errorf("Expected base64 data %q, got %q", expectedBase64, imageContent.Data)
	}
}