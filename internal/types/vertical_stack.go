// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"image"
	"image/color"
)

type VerticalStack struct {
	bindings    image.Rectangle
	iconImage   image.Image
	iconBounds  image.Rectangle
	borderColor color.RGBA
	fillColor   color.RGBA
	label       string
	fontColor   color.RGBA
	width       int
	height      int
	margin      Margin
	padding     Padding
	//parents *Stack
	children []Node
}

func (h VerticalStack) Init() Node {
	sr := Group{}
	sr.bindings = image.Rect(0, 0, 320, 190)
	sr.iconImage = image.NewRGBA(h.bindings)
	sr.iconBounds = image.Rect(0, 0, 0, 0)
	sr.borderColor = color.RGBA{0, 0, 0, 0}
	sr.fillColor = color.RGBA{0, 0, 0, 0}
	sr.label = ""
	sr.fontColor = color.RGBA{0, 0, 0, 0}
	sr.width = 320
	sr.height = 190
	sr.margin = Margin{0, 0, 0, 0}
	sr.padding = Padding{0, 0, 0, 0}
	sr.direction = "vertical"
	sr.align = "center"
	return &sr
}
