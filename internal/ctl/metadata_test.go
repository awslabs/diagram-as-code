// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func generateMinimalPNG() ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func TestEmbedAndExtractYAML(t *testing.T) {
	pngBytes, err := generateMinimalPNG()
	if err != nil {
		t.Fatalf("failed to generate test PNG: %v", err)
	}

	testYAML := []byte("Diagram:\n  Resources:\n    Test:\n      Type: AWS::Diagram::Resource")

	embeddedPNG, err := EmbedYAMLInPNG(pngBytes, testYAML)
	if err != nil {
		t.Fatalf("failed to embed YAML in PNG: %v", err)
	}

	// Verify PNG signature is still intact
	pngSig := []byte{137, 80, 78, 71, 13, 10, 26, 10}
	if !bytes.HasPrefix(embeddedPNG, pngSig) {
		t.Errorf("embedded PNG does not have valid PNG signature")
	}

	// Extract and verify
	extractedYAML, err := ExtractYAMLFromPNG(embeddedPNG)
	if err != nil {
		t.Fatalf("failed to extract YAML from PNG: %v", err)
	}

	if !bytes.Equal(extractedYAML, testYAML) {
		t.Errorf("extracted YAML does not match original. Got %q, want %q", string(extractedYAML), string(testYAML))
	}
}

func TestExtractMissingYAML(t *testing.T) {
	pngBytes, err := generateMinimalPNG()
	if err != nil {
		t.Fatalf("failed to generate test PNG: %v", err)
	}

	_, err = ExtractYAMLFromPNG(pngBytes)
	if err == nil {
		t.Error("expected error extracting missing metadata, got nil")
	}
}

func TestEmbedInvalidPNG(t *testing.T) {
	invalidPNG := []byte("not a PNG image")
	testYAML := []byte("Diagram:")

	_, err := EmbedYAMLInPNG(invalidPNG, testYAML)
	if err == nil {
		t.Error("expected error embedding in invalid PNG, got nil")
	}

	_, err = ExtractYAMLFromPNG(invalidPNG)
	if err == nil {
		t.Error("expected error extracting from invalid PNG, got nil")
	}
}
