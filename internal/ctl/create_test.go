// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"image"
	"testing"
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
