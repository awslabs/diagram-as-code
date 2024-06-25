// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"image"
	"image/color"
)

type HorizontalStack struct {
	bindings    image.Rectangle
	iconImage   image.Image
	iconBounds  image.Rectangle
	borderColor color.RGBA
	fillColor   color.RGBA
	label       string
	labelColor  color.RGBA
	margin      Margin
	padding     Padding
	//parents *Stack
	children []*Resource
}

func (v HorizontalStack) Init() *Resource {
	sr := Resource{}
	sr.bindings = &image.Rectangle{
		image.Point{0, 0},
		image.Point{320, 190},
	}
	sr.iconImage = image.NewRGBA(v.bindings)
	sr.iconBounds = image.Rect(0, 0, 0, 0)
	sr.borderColor = &color.RGBA{0, 0, 0, 0}
	sr.fillColor = color.RGBA{0, 0, 0, 0}
	sr.label = ""
	sr.labelColor = &color.RGBA{0, 0, 0, 0}
	sr.margin = &Margin{0, 0, 0, 0}
	sr.padding = &Padding{0, 0, 0, 0}
	sr.direction = "horizontal"
	sr.align = "center"
	return &sr
}
