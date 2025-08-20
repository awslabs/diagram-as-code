// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build windows
// +build windows

package font

var Paths = []string{
	"C:\\Windows\\Fonts\\arial.ttf",
	"goregular", // As the default font, uses golang.org/x/image/font/gofont/goregular. For more information about this font, go to: https://go.dev/blog/go-fonts
}
