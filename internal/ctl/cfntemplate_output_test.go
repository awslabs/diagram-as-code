// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCloudFormationDACOutputDoesNotOverwriteImage(t *testing.T) {
	original := template
	t.Cleanup(func() { template = original })
	for _, ext := range []string{".yaml", ".YAML", ".png", ".yml", ""} {
		t.Run("extension="+ext, func(t *testing.T) {
			dir := t.TempDir()
			input := filepath.Join(dir, "input.yaml")
			if err := os.WriteFile(input, []byte("Resources: {}\n"), 0600); err != nil {
				t.Fatal(err)
			}
			definitions := filepath.Join(dir, "definitions.yaml")
			if err := os.WriteFile(definitions, []byte(`Definitions:
  AWS::Diagram::Canvas:
    Type: Group
    CFn:
      HasChildren: true
  AWS::Diagram::Cloud:
    Type: Group
    CFn:
      HasChildren: true
`), 0600); err != nil {
				t.Fatal(err)
			}
			template = TemplateStruct{Diagram: Diagram{Resources: map[string]Resource{
				"Canvas":   {Type: "AWS::Diagram::Canvas", Children: []string{"AWSCloud"}},
				"AWSCloud": {Type: "AWS::Diagram::Cloud"},
			}}}
			output := filepath.Join(dir, "custom"+ext)
			if err := CreateDiagramFromCFnTemplate(input, &output, true, &CreateOptions{
				OverrideDefFile: definitions, OverrideFont: "goregular", OverwriteMode: Force,
			}); err != nil {
				t.Fatal(err)
			}
			wantImage := filepath.Join(dir, "custom"+ext)
			if ext == ".yaml" || ext == ".YAML" {
				wantImage = filepath.Join(dir, "custom.png")
			}
			if output != wantImage {
				t.Fatalf("image output = %q, want %q", output, wantImage)
			}
			data, err := os.ReadFile(wantImage)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := png.Decode(bytes.NewReader(data)); err != nil {
				t.Fatalf("image is not PNG: %v", err)
			}
			data, err = os.ReadFile(filepath.Join(dir, "custom.yaml"))
			if err != nil {
				t.Fatal(err)
			}
			var dac TemplateStruct
			if err := yaml.Unmarshal(data, &dac); err != nil {
				t.Fatalf("DAC output is not YAML: %v", err)
			}
			if _, ok := dac.Resources["Canvas"]; !ok {
				t.Fatal("DAC output has no Canvas resource")
			}
		})
	}
}
