// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"image"
	"image/color"
)

type VerticalStack struct {
}

func (h VerticalStack) Init() *Resource {
	sr := Resource{}
	sr.bindings = &image.Rectangle{
		image.Point{0, 0},
		image.Point{320, 190},
	}
	sr.iconImage = image.NewRGBA(*sr.bindings)
	sr.iconBounds = image.Rect(0, 0, 0, 0)
	sr.borderColor = &color.RGBA{0, 0, 0, 0}
	sr.fillColor = color.RGBA{0, 0, 0, 0}
	sr.label = ""
	sr.labelColor = &color.RGBA{0, 0, 0, 0}
	sr.margin = &Margin{0, 0, 0, 0}
	sr.padding = &Padding{0, 0, 0, 0}
	sr.direction = "vertical"
	sr.align = "center"
	return &sr
}
