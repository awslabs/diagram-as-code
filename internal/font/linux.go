// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build linux
// +build linux

package font

var Paths = []string{
	"/usr/share/fonts/truetype/msttcorefonts/Arial.ttf",                 // For Ubuntu linux ttf-mscorefonts-installer package.
	"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",   // For Ubuntu/Debian Linux fonts-liberation package.
	"/usr/share/fonts/truetype/liberation2/LiberationSans-Regular.ttf",  // For Ubuntu/Debian Linux fonts-liberation2 package.
	"/usr/share/fonts/liberation-sans/LiberationSans-Regular.ttf",       // For Fedora/AL2023 Linux liberation-sans-fonts package.
	"/usr/share/fonts/liberation/LiberationSans-Regular.ttf",            // For Alpine/Arch Linux ttf-liberation package.
	"/run/current-system/sw/share/X11/fonts/LiberationSans-Regular.ttf", // For NixOS Linux liberation_ttf package (enable fontDir in fonts options).
}
