// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
)

// EmbedYAMLInPNG embeds the YAML content in the PNG bytes.
// It parses the PNG structure, and inserts a tEXt chunk containing base64-encoded YAML
// immediately after the IHDR chunk.
func EmbedYAMLInPNG(pngBytes []byte, yamlContent []byte) ([]byte, error) {
	if len(pngBytes) < 8 {
		return nil, fmt.Errorf("invalid PNG: too short")
	}
	pngSig := []byte{137, 80, 78, 71, 13, 10, 26, 10}
	if !bytes.Equal(pngBytes[:8], pngSig) {
		return nil, fmt.Errorf("invalid PNG signature")
	}

	b64YAML := base64.StdEncoding.EncodeToString(yamlContent)

	keyword := []byte("awsdac")
	chunkData := make([]byte, len(keyword)+1+len(b64YAML))
	copy(chunkData[0:len(keyword)], keyword)
	chunkData[len(keyword)] = 0x00
	copy(chunkData[len(keyword)+1:], []byte(b64YAML))

	chunkType := []byte("tEXt")
	chunkLength := uint32(len(chunkData))

	var chunkBuf bytes.Buffer
	if err := binary.Write(&chunkBuf, binary.BigEndian, chunkLength); err != nil {
		return nil, err
	}
	chunkBuf.Write(chunkType)
	chunkBuf.Write(chunkData)

	h := crc32.NewIEEE()
	h.Write(chunkType)
	h.Write(chunkData)
	chunkCRC := h.Sum32()
	if err := binary.Write(&chunkBuf, binary.BigEndian, chunkCRC); err != nil {
		return nil, err
	}
	tEXtChunk := chunkBuf.Bytes()

	var out bytes.Buffer
	out.Write(pngSig)

	offset := 8
	ihdrFound := false

	for offset < len(pngBytes) {
		if offset+8 > len(pngBytes) {
			return nil, fmt.Errorf("malformed PNG: unexpected end of file")
		}
		length := binary.BigEndian.Uint32(pngBytes[offset : offset+4])
		chunkTypeStr := string(pngBytes[offset+4 : offset+8])
		chunkTotalLength := int(4 + 4 + length + 4)

		if offset+chunkTotalLength > len(pngBytes) {
			return nil, fmt.Errorf("malformed PNG: chunk length exceeds file size")
		}

		out.Write(pngBytes[offset : offset+chunkTotalLength])
		offset += chunkTotalLength

		if chunkTypeStr == "IHDR" {
			ihdrFound = true
			out.Write(tEXtChunk)
		}
	}

	if !ihdrFound {
		return nil, fmt.Errorf("IHDR chunk not found in PNG")
	}

	return out.Bytes(), nil
}

// ExtractYAMLFromPNG extracts the base64-encoded YAML content from a PNG file/bytes.
func ExtractYAMLFromPNG(pngBytes []byte) ([]byte, error) {
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
				if keyword == "awsdac" {
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

// ExtractYAMLFromPNGFile reads a PNG file and extracts the embedded YAML.
func ExtractYAMLFromPNGFile(filename string) ([]byte, error) {
	pngBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read PNG file: %w", err)
	}
	return ExtractYAMLFromPNG(pngBytes)
}
