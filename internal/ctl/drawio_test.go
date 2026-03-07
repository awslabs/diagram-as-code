// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"strings"
	"testing"
)

func TestRGBAToHex(t *testing.T) {
	if got := rgbaToHex("rgba(255,128,0,255)"); got != "#FF8000" {
		t.Fatalf("unexpected hex conversion result: %s", got)
	}

	raw := "not-a-color"
	if got := rgbaToHex(raw); got != raw {
		t.Fatalf("expected fallback to original value, got: %s", got)
	}
}

func TestEscapeXML(t *testing.T) {
	in := "a&b<c>d\"e\nf"
	got := escapeXML(in)

	for _, want := range []string{"&amp;", "&lt;", "&gt;", "&quot;", "&#xa;"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected escaped token %s in %s", want, got)
		}
	}
}

func TestGetDrawioNodeFallback(t *testing.T) {
	got := getDrawioNode("AWS::Unknown::Type")
	if got.isGroup {
		t.Fatal("expected unknown type to fallback to a non-group node")
	}
	if !strings.Contains(got.style, "rounded=1") {
		t.Fatalf("expected fallback style for unknown node, got: %s", got.style)
	}
}
