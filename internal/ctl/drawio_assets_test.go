// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"strings"
	"testing"
)

func TestResolveSVGPath(t *testing.T) {
	tests := []struct {
		name    string
		dacType string
		want    bool
	}{
		{name: "exact match", dacType: "AWS::S3::Bucket", want: true},
		{name: "service fallback", dacType: "AWS::ApiGateway::Anything", want: true},
		{name: "unknown type", dacType: "AWS::Unknown::Type", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveSVGPath(tt.dacType)
			if tt.want && got == "" {
				t.Fatalf("expected non-empty SVG path for type %s", tt.dacType)
			}
			if !tt.want && got != "" {
				t.Fatalf("expected empty SVG path for type %s, got %s", tt.dacType, got)
			}
		})
	}
}

func TestSVGToDataURI(t *testing.T) {
	input := `<svg width="10" height="10"><text>100% & #ok</text></svg>`
	got := svgToDataURI(input)

	if !strings.HasPrefix(got, "data:image/svg+xml,") {
		t.Fatalf("unexpected prefix: %s", got)
	}
	if strings.Contains(got, `"`) {
		t.Fatalf("expected no double quotes in encoded SVG: %s", got)
	}
	if !strings.Contains(got, "%25") || !strings.Contains(got, "%23") || !strings.Contains(got, "%3C") || !strings.Contains(got, "%3E") {
		t.Fatalf("expected critical characters to be encoded, got: %s", got)
	}
}

func TestGetAWSIconDataURIUnknownType(t *testing.T) {
	got := GetAWSIconDataURI("AWS::Unknown::Type")
	if got != "" {
		t.Fatalf("expected empty data URI for unknown type, got: %s", got)
	}
}
