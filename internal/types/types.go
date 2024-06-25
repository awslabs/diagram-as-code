// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"image"
	"image/color"
)

const WIDTH = 2

type Margin struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

func _max(a, b uint32) uint32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func _fetch_color(c color.Color) (r, g, b, a uint32) {
	r, g, b, a = c.RGBA()
	r = r >> 8
	g = g >> 8
	b = b >> 8
	a = a >> 8
	return
}

func _blend_color(c1 color.Color, c2 color.Color) color.Color {
	r1, g1, b1, a1 := _fetch_color(c1)
	r2, g2, b2, a2 := _fetch_color(c2)
	r := uint8(((r1 * (255 - a2)) + (r2 * a2)) / 255)
	g := uint8(((g1 * (255 - a2)) + (g2 * a2)) / 255)
	b := uint8(((b1 * (255 - a2)) + (b2 * a2)) / 255)
	a := uint8(_max(a1, a2))
	return color.RGBA{r, g, b, a}
}

func calcPosition(bindings image.Rectangle, position string) (image.Point, error) {
	x := bindings.Min.X
	y := bindings.Min.Y
	dx := bindings.Dx()
	dy := bindings.Dy()

	tx := [5]int{x, x + dx/4, x + dx/2, x + dx/2 + dx/4, x + dx}
	ty := [5]int{y, y + dy/4, y + dy/2, y + dy/2 + dy/4, y + dy}

	switch position {
	case "N":
		return image.Point{tx[2], ty[0]}, nil
	case "NNE":
		return image.Point{tx[3], ty[0]}, nil
	case "NE":
		return image.Point{tx[4], ty[0]}, nil
	case "ENE":
		return image.Point{tx[4], ty[1]}, nil
	case "E":
		return image.Point{tx[4], ty[2]}, nil
	case "ESE":
		return image.Point{tx[4], ty[3]}, nil
	case "SE":
		return image.Point{tx[4], ty[4]}, nil
	case "SSE":
		return image.Point{tx[3], ty[4]}, nil
	case "S":
		return image.Point{tx[2], ty[4]}, nil
	case "SSW":
		return image.Point{tx[1], ty[4]}, nil
	case "SW":
		return image.Point{tx[0], ty[4]}, nil
	case "WSW":
		return image.Point{tx[0], ty[3]}, nil
	case "W":
		return image.Point{tx[0], ty[2]}, nil
	case "WNW":
		return image.Point{tx[0], ty[1]}, nil
	case "NW":
		return image.Point{tx[0], ty[0]}, nil
	case "NNW":
		return image.Point{tx[1], ty[0]}, nil
	}
	return image.Point{tx[0], ty[0]}, fmt.Errorf("Unknown position: %s", position)
}
