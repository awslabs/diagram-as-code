// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package handler implements the Vercel serverless function for diagram generation.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	diagram "github.com/awslabs/diagram-as-code/pkg/diagram"
)

type generateRequest struct {
	YAML   string `json:"yaml"`
	Format string `json:"format"` // "png" or "drawio"
}

// Handler is the Vercel serverless entry point.
// POST /api/generate
// Body: {"yaml": "...", "format": "png"|"drawio"}
// Response: image/png binary OR application/xml text
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.YAML == "" {
		jsonError(w, "field 'yaml' is required", http.StatusBadRequest)
		return
	}

	if req.Format != "drawio" {
		req.Format = "png"
	}

	// Write YAML to temp file
	tmpIn, err := os.CreateTemp("", "dac-in-*.yaml")
	if err != nil {
		jsonError(w, "internal error creating temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpIn.Name())

	if _, err := tmpIn.WriteString(req.YAML); err != nil {
		tmpIn.Close()
		jsonError(w, "internal error writing input", http.StatusInternalServerError)
		return
	}
	tmpIn.Close()

	// Create output temp file with correct extension
	ext := ".png"
	if req.Format == "drawio" {
		ext = ".drawio"
	}
	tmpOut, err := os.CreateTemp("", "dac-out-*"+ext)
	if err != nil {
		jsonError(w, "internal error creating output file", http.StatusInternalServerError)
		return
	}
	tmpOut.Close()
	defer os.Remove(tmpOut.Name())

	opts := &diagram.CreateOptions{OverwriteMode: diagram.Force}
	outputFile := tmpOut.Name()

	if req.Format == "drawio" {
		err = diagram.CreateDrawioFromDacFile(tmpIn.Name(), &outputFile, opts)
	} else {
		err = diagram.CreateDiagramFromDacFile(tmpIn.Name(), &outputFile, opts)
	}

	if err != nil {
		jsonError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		jsonError(w, "failed to read output", http.StatusInternalServerError)
		return
	}

	if req.Format == "drawio" {
		w.Header().Set("Content-Type", "application/xml")
		w.Header().Set("Content-Disposition", `attachment; filename="diagram.drawio"`)
	} else {
		w.Header().Set("Content-Type", "image/png")
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data) //nolint:errcheck
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error":%q}`, msg)
}
