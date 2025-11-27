// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"image"
	"reflect"
	"testing"

	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
)

func TestResizeImage(t *testing.T) {
	// Create a test image
	width, height := 200, 100
	src := image.NewRGBA(image.Rect(0, 0, width, height))

	testCases := []struct {
		name           string
		targetWidth    int
		targetHeight   int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "No resize when both width and height are 0",
			targetWidth:    0,
			targetHeight:   0,
			expectedWidth:  width,
			expectedHeight: height,
		},
		{
			name:           "Resize by width only",
			targetWidth:    100,
			targetHeight:   0,
			expectedWidth:  100,
			expectedHeight: 50, // Maintains aspect ratio: 100/200 * 100 = 50
		},
		{
			name:           "Resize by height only",
			targetWidth:    0,
			targetHeight:   50,
			expectedWidth:  100, // Maintains aspect ratio: 50/100 * 200 = 100
			expectedHeight: 50,
		},
		{
			name:           "Resize by both width and height (width is limiting factor)",
			targetWidth:    50,
			targetHeight:   50,
			expectedWidth:  50,
			expectedHeight: 25, // Maintains aspect ratio: min(50/200, 50/100) * height = 25
		},
		{
			name:           "Resize by both width and height (height is limiting factor)",
			targetWidth:    150,
			targetHeight:   50,
			expectedWidth:  100, // Maintains aspect ratio: min(150/200, 50/100) * width = 100
			expectedHeight: 50,
		},
		{
			name:           "Enlarge image by width",
			targetWidth:    400,
			targetHeight:   0,
			expectedWidth:  400,
			expectedHeight: 200, // Maintains aspect ratio: 400/200 * 100 = 200
		},
		{
			name:           "Enlarge image by height",
			targetWidth:    0,
			targetHeight:   200,
			expectedWidth:  400, // Maintains aspect ratio: 200/100 * 200 = 400
			expectedHeight: 200,
		},
		{
			name:           "Enlarge image by both width and height (width is limiting factor)",
			targetWidth:    300,
			targetHeight:   200,
			expectedWidth:  300,
			expectedHeight: 150, // Maintains aspect ratio: min(300/200, 200/100) * height = 150
		},
		{
			name:           "Enlarge image by both width and height (height is limiting factor)",
			targetWidth:    500,
			targetHeight:   200,
			expectedWidth:  400, // Maintains aspect ratio: min(500/200, 200/100) * width = 400
			expectedHeight: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resized := resizeImage(src, tc.targetWidth, tc.targetHeight)

			actualWidth := resized.Bounds().Dx()
			actualHeight := resized.Bounds().Dy()

			if actualWidth != tc.expectedWidth || actualHeight != tc.expectedHeight {
				t.Errorf("Expected size %dx%d, got %dx%d",
					tc.expectedWidth, tc.expectedHeight, actualWidth, actualHeight)
			}
		})
	}
}

func TestLoadResourcesWithUndefinedTypeFallback(t *testing.T) {
	// Create definition structure with fallback definition
	ds := definition.DefinitionStructure{
		Definitions: map[string]*definition.Definition{
			"AWS::EC2": {
				Type: "Resource",
				Label: &definition.DefinitionLabel{
					Title: "EC2",
				},
			},
		},
	}

	// Create template with undefined type
	template := &TemplateStruct{
		Diagram: Diagram{
			Resources: map[string]Resource{
				"TestResource": {
					Type: "AWS::EC2::UndefinedType", // This type doesn't exist
				},
			},
		},
	}

	actualResources := make(map[string]*types.Resource)

	// Define expected test resource with label set
	expectedTestResource := new(types.Resource).Init()
	labelText := "EC2"
	expectedTestResource.SetLabel(&labelText, nil, nil)

	// This should not panic and should use fallback
	err := loadResources(template, ds, actualResources)
	if err != nil {
		t.Fatalf("loadResources failed: %v", err)
	}

	// Debug: Print full objects
	t.Logf("actualResources: %+v", actualResources)
	t.Logf("expectedTestResource: %+v", expectedTestResource)

	// Verify TestResource was created via fallback
	actualTestResource, exists := actualResources["TestResource"]
	if !exists {
		t.Fatal("TestResource was not created via fallback")
	}

	if actualTestResource == nil {
		t.Fatal("TestResource is nil")
	}

	// Debug: Print actualTestResource details
	t.Logf("actualTestResource: %+v", actualTestResource)

	// Verify it's the same type as expected (both are *types.Resource)
	if reflect.TypeOf(actualTestResource) != reflect.TypeOf(expectedTestResource) {
		t.Errorf("TestResource type mismatch. Expected: %T, Got: %T", expectedTestResource, actualTestResource)
	}

	// Deep check: Compare object contents
	isEqual := reflect.DeepEqual(actualTestResource, expectedTestResource)
	t.Logf("DeepEqual result: %v", isEqual)
	if !isEqual {
		t.Errorf("TestResource deep comparison failed.\nExpected: %+v\nActual: %+v", expectedTestResource, actualTestResource)
	}
}

func TestLoadResourcesWithDefinedType(t *testing.T) {
	// Create definition structure with defined type
	ds := definition.DefinitionStructure{
		Definitions: map[string]*definition.Definition{
			"AWS::EC2::Instance": {
				Type: "Resource",
				Label: &definition.DefinitionLabel{
					Title: "EC2 Instance",
				},
			},
		},
	}

	// Create template with defined type (no fallback needed)
	template := &TemplateStruct{
		Diagram: Diagram{
			Resources: map[string]Resource{
				"TestResource": {
					Type: "AWS::EC2::Instance", // This type exists
				},
			},
		},
	}

	actualResources := make(map[string]*types.Resource)

	// Define expected test resource with label set
	expectedTestResource := new(types.Resource).Init()
	labelText := "EC2 Instance"
	expectedTestResource.SetLabel(&labelText, nil, nil)

	// This should not panic and should use direct definition
	err := loadResources(template, ds, actualResources)
	if err != nil {
		t.Fatalf("loadResources failed: %v", err)
	}

	// Verify TestResource was created without fallback
	actualTestResource, exists := actualResources["TestResource"]
	if !exists {
		t.Fatal("TestResource was not created")
	}

	if actualTestResource == nil {
		t.Fatal("TestResource is nil")
	}

	// Verify it's the same type as expected
	if reflect.TypeOf(actualTestResource) != reflect.TypeOf(expectedTestResource) {
		t.Errorf("TestResource type mismatch. Expected: %T, Got: %T", expectedTestResource, actualTestResource)
	}

	// Deep check: Compare object contents
	if !reflect.DeepEqual(actualTestResource, expectedTestResource) {
		t.Errorf("TestResource deep comparison failed.\nExpected: %+v\nActual: %+v", expectedTestResource, actualTestResource)
	}
}

func TestLoadResourcesWithNoFallbackPossible(t *testing.T) {
	// Create definition structure without fallback definition
	ds := definition.DefinitionStructure{
		Definitions: map[string]*definition.Definition{
			"AWS::S3::Bucket": {
				Type: "Resource",
				Label: &definition.DefinitionLabel{
					Title: "S3 Bucket",
				},
			},
		},
	}

	// Create template with undefined type that cannot fallback
	template := &TemplateStruct{
		Diagram: Diagram{
			Resources: map[string]Resource{
				"TestResource": {
					Type: "AWS::NonExistent::Service", // Neither this nor AWS::NonExistent exists
				},
			},
		},
	}

	actualResources := make(map[string]*types.Resource)

	// This should not panic but should skip the resource
	err := loadResources(template, ds, actualResources)
	if err != nil {
		t.Fatalf("loadResources failed: %v", err)
	}

	// Verify TestResource was NOT created (skipped due to no fallback)
	_, exists := actualResources["TestResource"]
	if exists {
		t.Error("TestResource should not have been created when no fallback is possible")
	}

	// Verify only Canvas was created
	if len(actualResources) != 1 {
		t.Errorf("Expected 1 resource (Canvas only), got %d", len(actualResources))
	}

	canvas, exists := actualResources["Canvas"]
	if !exists || canvas == nil {
		t.Error("Canvas resource should still be created")
	}
}

func TestIsAllowedDefinitionURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "official raw.githubusercontent.com URL",
			url:     "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml",
			wantErr: false,
		},
		{
			name:    "official github.com URL",
			url:     "https://github.com/awslabs/diagram-as-code/releases/download/v1.0.0/definitions.yaml",
			wantErr: false,
		},
		{
			name:    "untrusted URL",
			url:     "https://example.com/malicious-definitions.yaml",
			wantErr: true,
		},
		{
			name:    "localhost URL",
			url:     "http://localhost:8080/definitions.yaml",
			wantErr: true,
		},
		{
			name:    "private IP URL",
			url:     "http://192.168.1.1/definitions.yaml",
			wantErr: true,
		},
		{
			name:    "different github org",
			url:     "https://raw.githubusercontent.com/other-org/diagram-as-code/main/definitions.yaml",
			wantErr: true,
		},
		{
			name:    "different github repo",
			url:     "https://raw.githubusercontent.com/awslabs/other-repo/main/definitions.yaml",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isAllowedDefinitionURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("isAllowedDefinitionURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
