// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build darwin
// +build darwin

package font

var Paths = []string{
	"/Library/Fonts/Arial Unicode.ttf",
	"goregular", // As the default font, uses golang.org/x/image/font/gofont/goregular. For more information about this font, go to: https://go.dev/blog/go-fonts
}
