// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
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
	extractedYAML, err := extractYAMLFromPNG(embeddedPNG)
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

	_, err = extractYAMLFromPNG(pngBytes)
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

	_, err = extractYAMLFromPNG(invalidPNG)
	if err == nil {
		t.Error("expected error extracting from invalid PNG, got nil")
	}
}

// extractYAMLFromPNG reads the embedded YAML back for round-trip testing only.
func extractYAMLFromPNG(pngBytes []byte) ([]byte, error) {
	if len(pngBytes) < 8 {
		return nil, fmt.Errorf("invalid PNG: too short")
	}
	pngSig := []byte{137, 80, 78, 71, 13, 10, 26, 10}
	if !bytes.Equal(pngBytes[:8], pngSig) {
		return nil, fmt.Errorf("invalid PNG signature")
	}

	offset := 8
	for offset < len(pngBytes) {
		if offset+8 > len(pngBytes) {
			break
		}
		length := binary.BigEndian.Uint32(pngBytes[offset : offset+4])
		chunkTypeStr := string(pngBytes[offset+4 : offset+8])
		chunkTotalLength := int(4 + 4 + length + 4)

		if offset+chunkTotalLength > len(pngBytes) {
			return nil, fmt.Errorf("malformed PNG: chunk length exceeds file size")
		}

		if chunkTypeStr == "tEXt" {
			data := pngBytes[offset+8 : offset+8+int(length)]
			nullIdx := bytes.IndexByte(data, 0x00)
			if nullIdx != -1 {
				keyword := string(data[:nullIdx])
				if keyword == pngMetadataKeyword {
					b64YAML := string(data[nullIdx+1:])
					yamlContent, err := base64.StdEncoding.DecodeString(b64YAML)
					if err != nil {
						return nil, fmt.Errorf("failed to decode base64 YAML: %w", err)
					}
					return yamlContent, nil
				}
			}
		}

		offset += chunkTotalLength
	}

	return nil, fmt.Errorf("no awsdac diagram-as-code YAML metadata found in PNG")
}
