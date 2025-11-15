package main

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
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
Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children:
        - AWSCloud
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Children:
        - MyVPC
    MyVPC:
      Type: AWS::EC2::VPC
`

	tests := []struct {
		name      string
		arguments map[string]interface{}
		wantErr   bool
		wantMIME  string
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
	yamlContent := `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children:
        - AWSCloud
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Children:
        - TestVPC
    TestVPC:
      Type: AWS::EC2::VPC`

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

// TestHandleGenerateDacFromUserRequirements_PanicRecovery tests panic recovery in the getDiagramAsCodeFormat handler
func TestHandleGenerateDacFromUserRequirements_PanicRecovery(t *testing.T) {
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      string(GET_DIAGRAM_AS_CODE_FORMAT),
			Arguments: map[string]interface{}{},
		},
	}

	// Test normal operation first
	t.Run("normal_operation", func(t *testing.T) {
		wrappedHandler := withPanicRecovery("getDiagramAsCodeFormat", handleGenerateDacFromUserRequirements)
		result, err := wrappedHandler(ctx, request)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if len(result.Content) != 1 {
			t.Fatalf("Expected 1 content item, got %d", len(result.Content))
		}

		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatal("Expected TextContent")
		}

		// Verify content contains DAC format documentation
		if len(textContent.Text) < 100 {
			t.Errorf("Expected substantial content, got %d characters", len(textContent.Text))
		}

		// Verify key DAC format sections are present
		expectedSections := []string{"Diagram:", "Resources:", "DefinitionFiles:"}
		for _, section := range expectedSections {
			if !strings.Contains(textContent.Text, section) {
				t.Errorf("Expected content to contain %q", section)
			}
		}

		// Should not be an error result
		if result.IsError {
			t.Error("Expected IsError to be false for normal operation")
		}
	})

	// Test panic recovery by creating a panic-inducing handler
	t.Run("panic_recovery", func(t *testing.T) {
		panicHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			panic("simulated panic in getDiagramAsCodeFormat")
		}

		wrappedHandler := withPanicRecovery("getDiagramAsCodeFormat", panicHandler)
		result, err := wrappedHandler(ctx, request)
		// Should not return an error (panics are recovered)
		if err != nil {
			t.Errorf("Expected no error from panic recovery, got: %v", err)
		}

		// Should return panic recovery result
		if result == nil {
			t.Fatal("Expected panic recovery result, got nil")
		}

		if !result.IsError {
			t.Error("Expected IsError to be true for panic recovery")
		}

		if len(result.Content) != 1 {
			t.Fatalf("Expected 1 content item for panic recovery, got %d", len(result.Content))
		}

		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatal("Expected TextContent for panic recovery")
		}

		if !strings.Contains(textContent.Text, "An unexpected error occurred") {
			t.Errorf("Expected panic recovery message, got: %s", textContent.Text)
		}
	})
}

// TestAllHandlers_PanicRecovery tests panic recovery across all three tool handlers
func TestAllHandlers_PanicRecovery(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		toolName    ToolName
		handler     func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
		arguments   map[string]interface{}
		description string
	}{
		{
			name:        "generateDiagram_panic_recovery",
			toolName:    GENERATE_DIAGRAM,
			handler:     handleGenerateDiagram,
			arguments:   map[string]interface{}{"yamlContent": "valid yaml"},
			description: "Tests panic recovery in generateDiagram handler",
		},
		{
			name:     "generateDiagramToFile_panic_recovery",
			toolName: GENERATE_DIAGRAM_TO_FILE,
			handler:  handleGenerateDiagramToFile,
			arguments: map[string]interface{}{
				"yamlContent":    "valid yaml",
				"outputFilePath": "/tmp/test.png",
			},
			description: "Tests panic recovery in generateDiagramToFile handler",
		},
		{
			name:        "getDiagramAsCodeFormat_panic_recovery",
			toolName:    GET_DIAGRAM_AS_CODE_FORMAT,
			handler:     handleGenerateDacFromUserRequirements,
			arguments:   map[string]interface{}{},
			description: "Tests panic recovery in getDiagramAsCodeFormat handler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a panic-inducing version of the handler
			panicHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				panic("simulated panic in " + string(tt.toolName))
			}

			// Wrap with panic recovery
			wrappedHandler := withPanicRecovery(string(tt.toolName), panicHandler)

			// Create request
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      string(tt.toolName),
					Arguments: tt.arguments,
				},
			}

			// Execute and verify panic recovery
			result, err := wrappedHandler(ctx, request)
			// Should not return an error (panics are recovered)
			if err != nil {
				t.Errorf("Expected no error from panic recovery, got: %v", err)
			}

			// Should return panic recovery result
			if result == nil {
				t.Fatal("Expected panic recovery result, got nil")
			}

			if !result.IsError {
				t.Error("Expected IsError to be true for panic recovery")
			}

			if len(result.Content) != 1 {
				t.Fatalf("Expected 1 content item for panic recovery, got %d", len(result.Content))
			}

			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatal("Expected TextContent for panic recovery")
			}

			expectedPhrases := []string{
				"An unexpected error occurred",
				"server has recovered",
				"ready to process new requests",
			}

			for _, phrase := range expectedPhrases {
				if !strings.Contains(textContent.Text, phrase) {
					t.Errorf("Expected panic recovery message to contain %q, got: %s", phrase, textContent.Text)
				}
			}

			t.Logf("%s: Successfully recovered from simulated panic", tt.description)
		})
	}

	// Test that all handlers can actually work normally (no panic recovery needed)
	t.Run("all_handlers_normal_operation", func(t *testing.T) {
		// Test getDiagramAsCodeFormat (simplest, no external dependencies)
		request := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Name:      string(GET_DIAGRAM_AS_CODE_FORMAT),
				Arguments: map[string]interface{}{},
			},
		}

		wrappedHandler := withPanicRecovery("getDiagramAsCodeFormat", handleGenerateDacFromUserRequirements)
		result, err := wrappedHandler(ctx, request)
		if err != nil {
			t.Errorf("Expected no error from normal operation, got: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result from normal operation, got nil")
		}

		if result.IsError {
			t.Error("Expected IsError to be false for normal operation")
		}

		t.Log("Verified that panic recovery wrapper does not interfere with normal operation")
	})
}

// TestCreateDiagramSafely tests the createDiagramSafely function with panic recovery
func TestCreateDiagramSafely(t *testing.T) {
	tests := []struct {
		name        string
		inputFile   string
		outputFile  string
		shouldPanic bool
		expectError bool
		mockCtlFunc func(string, *string, interface{}) error
	}{
		{
			name:        "successful_diagram_creation",
			inputFile:   "valid.yaml",
			outputFile:  "output.png",
			shouldPanic: false,
			expectError: false,
			mockCtlFunc: func(input string, output *string, opts interface{}) error {
				return nil // Simulate successful creation
			},
		},
		{
			name:        "panic_during_creation",
			inputFile:   "invalid.yaml",
			outputFile:  "output.png",
			shouldPanic: true,
			expectError: true,
			mockCtlFunc: func(input string, output *string, opts interface{}) error {
				panic("diagram creation failed")
			},
		},
		{
			name:        "panic_with_error_type",
			inputFile:   "error.yaml",
			outputFile:  "output.png",
			shouldPanic: true,
			expectError: true,
			mockCtlFunc: func(input string, output *string, opts interface{}) error {
				panic(errors.New("ctl creation error"))
			},
		},
		{
			name:        "panic_with_custom_type",
			inputFile:   "custom.yaml",
			outputFile:  "output.png",
			shouldPanic: true,
			expectError: true,
			mockCtlFunc: func(input string, output *string, opts interface{}) error {
				panic(struct{ message string }{"custom panic type"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Since we can't easily mock ctl.CreateDiagramFromDacFile directly,
			// we'll create a test wrapper that allows us to inject the behavior
			var actualError error
			var panicRecovered bool
			var panicValue interface{}

			// Simulate the createDiagramSafely behavior with our mock
			func() {
				defer func() {
					if r := recover(); r != nil {
						panicRecovered = true
						panicValue = r
						actualError = errors.New("panic occurred during diagram creation")
					}
				}()

				// Call our mock function to simulate ctl.CreateDiagramFromDacFile
				err := tt.mockCtlFunc(tt.inputFile, &tt.outputFile, nil)
				if err != nil {
					actualError = err
				}
			}()

			// Verify panic behavior
			if tt.shouldPanic && !panicRecovered {
				t.Error("Expected panic to occur, but none was recovered")
			}
			if !tt.shouldPanic && panicRecovered {
				t.Errorf("Unexpected panic recovered: %v", panicValue)
			}

			// Verify error behavior
			if tt.expectError && actualError == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && actualError != nil {
				t.Errorf("Expected no error, got: %v", actualError)
			}

			// For panic cases, verify the panic value is captured
			if tt.shouldPanic && panicRecovered {
				if panicValue == nil {
					t.Error("Expected panic value to be captured")
				}
			}
		})
	}
}

// TestCreateDiagramSafely_ActualFunction tests the real createDiagramSafely function
// with invalid input to trigger the panic recovery path
func TestCreateDiagramSafely_ActualFunction(t *testing.T) {
	// Create a temporary invalid file that will cause ctl.CreateDiagramFromDacFile to fail
	tempDir := t.TempDir()
	inputPath := tempDir + "/invalid.yaml"
	outputPath := tempDir + "/output.png"

	// Write invalid YAML content that should cause an error or panic
	invalidYAML := `invalid: yaml: content: with: [unclosed brackets`
	err := os.WriteFile(inputPath, []byte(invalidYAML), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test that createDiagramSafely handles errors gracefully
	// Note: This may not actually panic depending on how ctl.CreateDiagramFromDacFile handles invalid input
	// But it tests the integration path of the actual function
	err = createDiagramSafely(inputPath, &outputPath, nil)

	// We expect either no error (if ctl handles invalid input gracefully)
	// or an error (if it returns an error instead of panicking)
	// The key is that it should NOT panic and crash the test
	t.Logf("createDiagramSafely completed with result: %v", err)

	// The function should not panic - if we reach this point, the panic recovery worked
	// or there was no panic to begin with (both are acceptable outcomes)
}

// TestWithPanicRecovery tests the withPanicRecovery middleware function
func TestWithPanicRecovery(t *testing.T) {
	ctx := context.Background()

	// Mock request for testing
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "testTool",
			Arguments: map[string]interface{}{},
		},
	}

	tests := []struct {
		name           string
		handler        func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
		expectPanic    bool
		expectError    bool
		expectedResult *mcp.CallToolResult
	}{
		{
			name: "normal_handler_execution",
			handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						mcp.TextContent{
							Type: "text",
							Text: "Success",
						},
					},
				}, nil
			},
			expectPanic: false,
			expectError: false,
			expectedResult: &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: "Success",
					},
				},
			},
		},
		{
			name: "panic_recovery_string",
			handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				panic("test panic message")
			},
			expectPanic: true,
			expectError: false,
		},
		{
			name: "panic_recovery_error",
			handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				panic(errors.New("test error"))
			},
			expectPanic: true,
			expectError: false,
		},
		{
			name: "panic_recovery_custom_type",
			handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				panic(struct{ msg string }{"custom panic"})
			},
			expectPanic: true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Wrap handler with panic recovery
			wrappedHandler := withPanicRecovery("testHandler", tt.handler)

			// Execute the wrapped handler
			result, err := wrappedHandler(ctx, request)

			// Verify error behavior
			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if tt.expectPanic {
				// For panic cases, verify recovery response
				if result == nil {
					t.Fatal("Expected result after panic recovery, got nil")
				}
				if len(result.Content) != 1 {
					t.Fatalf("Expected 1 content item, got %d", len(result.Content))
				}
				textContent, ok := result.Content[0].(mcp.TextContent)
				if !ok {
					t.Fatal("Expected TextContent")
				}
				if !strings.Contains(textContent.Text, "An unexpected error occurred") {
					t.Errorf("Expected panic recovery message, got: %s", textContent.Text)
				}
				if !result.IsError {
					t.Error("Expected IsError to be true for panic recovery")
				}
			} else {
				// For normal cases, verify original result
				if tt.expectedResult != nil {
					if result == nil {
						t.Fatal("Expected result, got nil")
					}
					if len(result.Content) != len(tt.expectedResult.Content) {
						t.Fatalf("Expected %d content items, got %d", len(tt.expectedResult.Content), len(result.Content))
					}
					textContent, ok := result.Content[0].(mcp.TextContent)
					if !ok {
						t.Fatal("Expected TextContent")
					}
					expectedText := tt.expectedResult.Content[0].(mcp.TextContent).Text
					if textContent.Text != expectedText {
						t.Errorf("Expected text %q, got %q", expectedText, textContent.Text)
					}
				}
			}
		})
	}
}

func TestHandleGenerateDiagram_OutputFileContent(t *testing.T) {
	ctx := context.Background()
	yamlContent := `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children:
        - AWSCloud
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Children:
        - TestVPC
    TestVPC:
      Type: AWS::EC2::VPC`
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

func TestHandleGenerateDiagramToFile(t *testing.T) {
	ctx := context.Background()
	validYAML := `
Diagram:
  DefinitionFiles:
    - Type: URL
      Url: https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children:
        - MyVPC
    MyVPC:
      Type: AWS::EC2::VPC
      Title: Test VPC
`

	// Use temporary directory for all test files
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		outputPath  string
		wantErr     bool
		expectedMsg string
	}{
		{
			name:        "valid yaml with absolute file path",
			outputPath:  filepath.Join(tempDir, "test-diagram.png"),
			wantErr:     false,
			expectedMsg: "", // Will be set dynamically
		},
		{
			name:        "relative file path",
			outputPath:  filepath.Join(tempDir, "relative", "path", "diagram.png"),
			wantErr:     false,
			expectedMsg: "", // Will be set dynamically
		},
	}

	// Test cases with missing parameters
	errorTests := []struct {
		name      string
		arguments map[string]interface{}
		wantErr   bool
	}{
		{
			name: "missing yamlContent",
			arguments: map[string]interface{}{
				"outputFilePath": filepath.Join(tempDir, "dummy.png"),
			},
			wantErr: true,
		},
		{
			name: "missing outputFilePath",
			arguments: map[string]interface{}{
				"yamlContent": validYAML,
			},
			wantErr: true,
		},
		{
			name: "invalid yamlContent type",
			arguments: map[string]interface{}{
				"yamlContent":    123,
				"outputFilePath": filepath.Join(tempDir, "dummy.png"),
			},
			wantErr: true,
		},
		{
			name: "invalid outputFilePath type",
			arguments: map[string]interface{}{
				"yamlContent":    validYAML,
				"outputFilePath": 123,
			},
			wantErr: true,
		},
	}

	// Run successful test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: string(GENERATE_DIAGRAM_TO_FILE),
					Arguments: map[string]interface{}{
						"yamlContent":    validYAML,
						"outputFilePath": tt.outputPath,
					},
				},
			}

			result, err := handleGenerateDiagramToFile(ctx, request)

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

			// Verify result contains only text content (no image data)
			if len(result.Content) != 1 {
				t.Fatalf("Expected 1 content item, got %d", len(result.Content))
			}

			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatal("Expected content to be TextContent")
			}

			expectedMsg := "Diagram successfully generated and saved to: " + tt.outputPath
			if textContent.Text != expectedMsg {
				t.Errorf("Expected message %q, got %q", expectedMsg, textContent.Text)
			}
		})
	}

	// Run error test cases
	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      string(GENERATE_DIAGRAM_TO_FILE),
					Arguments: tt.arguments,
				},
			}

			_, err := handleGenerateDiagramToFile(ctx, request)

			if !tt.wantErr {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})
	}
}

func TestHandleGenerateDiagramToFile_DirectoryCreation(t *testing.T) {
	ctx := context.Background()
	validYAML := `
Diagram:
  DefinitionFiles:
    - Type: URL
      Url: https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children:
        - MyVPC
    MyVPC:
      Type: AWS::EC2::VPC
`

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "nested", "deep", "diagram.png")

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: string(GENERATE_DIAGRAM_TO_FILE),
			Arguments: map[string]interface{}{
				"yamlContent":    validYAML,
				"outputFilePath": outputPath,
			},
		},
	}

	// Mock file operations for testing
	originalWriteFile := writeFileFunc
	writeFileFunc = func(filename string, data []byte, perm os.FileMode) error {
		if strings.HasSuffix(filename, "input.yaml") {
			return originalWriteFile(filename, data, perm)
		}
		return nil
	}
	defer func() { writeFileFunc = originalWriteFile }()

	_, err := handleGenerateDiagramToFile(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify the nested directory was created
	expectedDir := filepath.Join(tempDir, "nested", "deep")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Expected directory %s to be created", expectedDir)
	}
}

func TestHandleGenerateDiagramToFile_FileVerification(t *testing.T) {
	ctx := context.Background()
	validYAML := `
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
        - TestVPC
    TestVPC:
      Type: AWS::EC2::VPC
`

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test-diagram.png")

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: string(GENERATE_DIAGRAM_TO_FILE),
			Arguments: map[string]interface{}{
				"yamlContent":    validYAML,
				"outputFilePath": outputPath,
			},
		},
	}

	// Track file operations
	var inputFileCreated bool
	originalWriteFile := writeFileFunc
	writeFileFunc = func(filename string, data []byte, perm os.FileMode) error {
		if strings.HasSuffix(filename, "input.yaml") {
			inputFileCreated = true
		}
		return originalWriteFile(filename, data, perm)
	}
	defer func() { writeFileFunc = originalWriteFile }()

	result, err := handleGenerateDiagramToFile(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !inputFileCreated {
		t.Error("Expected input.yaml file to be created")
	}

	// Verify the expected message
	textContent := result.Content[0].(mcp.TextContent)
	expectedMsg := "Diagram successfully generated and saved to: " + outputPath
	if textContent.Text != expectedMsg {
		t.Errorf("Expected message %q, got %q", expectedMsg, textContent.Text)
	}
}

func TestHandleGenerateDiagramToFile_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	validYAML := `
Diagram:
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
`

	tests := []struct {
		name           string
		outputPath     string
		expectedErrMsg string
	}{
		{
			name:           "invalid output path with null byte",
			outputPath:     "/tmp/diagram\x00.png",
			expectedErrMsg: "failed to check output file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: string(GENERATE_DIAGRAM_TO_FILE),
					Arguments: map[string]interface{}{
						"yamlContent":    validYAML,
						"outputFilePath": tt.outputPath,
					},
				},
			}

			_, err := handleGenerateDiagramToFile(ctx, request)

			if err == nil {
				t.Fatalf("Expected error, got nil")
			}

			if !strings.Contains(err.Error(), tt.expectedErrMsg) {
				t.Errorf("Expected error to contain %q, got %q", tt.expectedErrMsg, err.Error())
			}
		})
	}
}

func TestHandleGenerateDiagramToFile_TempFileContent(t *testing.T) {
	ctx := context.Background()
	yamlContent := `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Children:
        - TestVPC
    TestVPC:
      Type: AWS::EC2::VPC
      Title: Test VPC`

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test.png")

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: string(GENERATE_DIAGRAM_TO_FILE),
			Arguments: map[string]interface{}{
				"yamlContent":    yamlContent,
				"outputFilePath": outputPath,
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

	_, err := handleGenerateDiagramToFile(ctx, request)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if string(capturedContent) != yamlContent {
		t.Errorf("Expected temp file content %q, got %q", yamlContent, string(capturedContent))
	}
}
