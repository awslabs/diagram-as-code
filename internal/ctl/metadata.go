// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// pngMetadataKeyword is the tEXt chunk keyword used for the embedded YAML.
const pngMetadataKeyword = "awsdac.yaml"

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

	keyword := []byte(pngMetadataKeyword)
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
