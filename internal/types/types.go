// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"image"
	"image/color"
)

const WIDTH = 2

type Windrose int

const (
	WINDROSE_N Windrose = iota
	WINDROSE_NNE
	WINDROSE_NE
	WINDROSE_ENE
	WINDROSE_E
	WINDROSE_ESE
	WINDROSE_SE
	WINDROSE_SSE
	WINDROSE_S
	WINDROSE_SSW
	WINDROSE_SW
	WINDROSE_WSW
	WINDROSE_W
	WINDROSE_WNW
	WINDROSE_NW
	WINDROSE_NNW
	WINDROSE_AUTO = -1 // Special value for auto-positioning
)

func ConvertWindrose(position string) (Windrose, error) {
	if position == "" || position == "auto" {
		return WINDROSE_AUTO, nil
	}

	switch position {
	case "N":
		return WINDROSE_N, nil
	case "NNE":
		return WINDROSE_NNE, nil
	case "NE":
		return WINDROSE_NE, nil
	case "ENE":
		return WINDROSE_ENE, nil
	case "E":
		return WINDROSE_E, nil
	case "ESE":
		return WINDROSE_ESE, nil
	case "SE":
		return WINDROSE_SE, nil
	case "SSE":
		return WINDROSE_SSE, nil
	case "S":
		return WINDROSE_S, nil
	case "SSW":
		return WINDROSE_SSW, nil
	case "SW":
		return WINDROSE_SW, nil
	case "WSW":
		return WINDROSE_WSW, nil
	case "W":
		return WINDROSE_W, nil
	case "WNW":
		return WINDROSE_WNW, nil
	case "NW":
		return WINDROSE_NW, nil
	case "NNW":
		return WINDROSE_NNW, nil
	}
	return 0, fmt.Errorf("unknown position: %s, supported positions are N, NNE, NE, ENE, E, ESE, SE, SSE, S, SSW, SW, WSW, W, WNW, NW, NNW, auto", position)
}

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

func calcPosition(bindings image.Rectangle, position Windrose) (image.Point, error) {
	x := bindings.Min.X
	y := bindings.Min.Y
	dx := bindings.Dx()
	dy := bindings.Dy()

	tx := [5]int{x, x + dx/4, x + dx/2, x + dx/2 + dx/4, x + dx}
	ty := [5]int{y, y + dy/4, y + dy/2, y + dy/2 + dy/4, y + dy}

	switch position {
	case WINDROSE_N:
		return image.Point{tx[2], ty[0]}, nil
	case WINDROSE_NNE:
		return image.Point{tx[3], ty[0]}, nil
	case WINDROSE_NE:
		return image.Point{tx[4], ty[0]}, nil
	case WINDROSE_ENE:
		return image.Point{tx[4], ty[1]}, nil
	case WINDROSE_E:
		return image.Point{tx[4], ty[2]}, nil
	case WINDROSE_ESE:
		return image.Point{tx[4], ty[3]}, nil
	case WINDROSE_SE:
		return image.Point{tx[4], ty[4]}, nil
	case WINDROSE_SSE:
		return image.Point{tx[3], ty[4]}, nil
	case WINDROSE_S:
		return image.Point{tx[2], ty[4]}, nil
	case WINDROSE_SSW:
		return image.Point{tx[1], ty[4]}, nil
	case WINDROSE_SW:
		return image.Point{tx[0], ty[4]}, nil
	case WINDROSE_WSW:
		return image.Point{tx[0], ty[3]}, nil
	case WINDROSE_W:
		return image.Point{tx[0], ty[2]}, nil
	case WINDROSE_WNW:
		return image.Point{tx[0], ty[1]}, nil
	case WINDROSE_NW:
		return image.Point{tx[0], ty[0]}, nil
	case WINDROSE_NNW:
		return image.Point{tx[1], ty[0]}, nil
	case WINDROSE_AUTO:
		return image.Point{tx[0], ty[0]}, fmt.Errorf("WINDROSE_AUTO should be resolved before calling calcPosition")
	}
	return image.Point{tx[0], ty[0]}, fmt.Errorf("unknown position: %d, supported positions are N, NNE, NE, ENE, E, ESE, SE, SSE, S, SSW, SW, WSW, W, WNW, NW, NNW", position)
}
