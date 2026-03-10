// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package diagram is a public wrapper around internal/ctl, exposing diagram
// generation functions that can be imported by code outside the module tree
// (e.g. Vercel serverless handlers).
package diagram

import "github.com/awslabs/diagram-as-code/internal/ctl"

// OverwriteMode mirrors ctl.OverwriteMode.
type OverwriteMode = ctl.OverwriteMode

const (
	Ask         OverwriteMode = ctl.Ask
	Force       OverwriteMode = ctl.Force
	NoOverwrite OverwriteMode = ctl.NoOverwrite
)

// CreateOptions mirrors ctl.CreateOptions.
type CreateOptions = ctl.CreateOptions

// CreateDiagramFromDacFile generates a PNG diagram from a DAC YAML file.
func CreateDiagramFromDacFile(inputfile string, outputfile *string, opts *CreateOptions) error {
	return ctl.CreateDiagramFromDacFile(inputfile, outputfile, opts)
}

// CreateDrawioFromDacFile generates a draw.io XML file from a DAC YAML file.
func CreateDrawioFromDacFile(inputfile string, outputfile *string, opts *CreateOptions) error {
	return ctl.CreateDrawioFromDacFile(inputfile, outputfile, opts)
}

// CreatePDFFromDacFile generates a PDF document from a DAC YAML file.
func CreatePDFFromDacFile(inputfile string, outputfile *string, opts *CreateOptions) error {
	return ctl.CreatePDFFromDacFile(inputfile, outputfile, opts)
}
