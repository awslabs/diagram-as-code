// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"image"
	"image/color"
	"math"

	log "github.com/sirupsen/logrus"
)

type Link struct {
	Source         *Node
	SourcePosition string
	Target         *Node
	TargetPosition string
	LineWidth      int
	drawn          bool
	lineColor      color.RGBA
}

func (l Link) Init(source *Node, sourcePosition string, target *Node, targetPosition string, lineWidth int) *Link {
	gl := Link{}
	gl.Source = source
	gl.SourcePosition = sourcePosition
	gl.Target = target
	gl.TargetPosition = targetPosition
	gl.LineWidth = lineWidth
	gl.drawn = false
	gl.lineColor = color.RGBA{0, 0, 0, 255}
	return &gl
}

func (l *Link) drawNeighborsDot(img *image.RGBA, x, y float64) {
	lowerPt := image.Point{int(x), int(y)}

	neighborPts := []image.Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	for _, neighborPt := range neighborPts {
		targetPt := lowerPt.Add(neighborPt)
		sourceColor := img.At(targetPt.X, targetPt.Y)
		dx := math.Abs(float64(targetPt.X) - x)
		dy := math.Abs(float64(targetPt.Y) - y)
		c := (1.0 - dx) * (1.0 - dy) * float64(l.lineColor.A)
		targetColor := l.lineColor
		targetColor.A = uint8(c)

		img.Set(targetPt.X, targetPt.Y, _blend_color(sourceColor, targetColor))
	}
}

func (l *Link) Draw(img *image.RGBA) {
	source := *l.Source
	target := *l.Target
	if l.drawn {
		log.Info("Link already drawn")
		return
	}
	log.Info("Link Drawing")
	sourcePt, _ := calcPosition(source.GetBindings(), l.SourcePosition)
	targetPt, _ := calcPosition(target.GetBindings(), l.TargetPosition)

	dx := float64(targetPt.X - sourcePt.X)
	dy := float64(targetPt.Y - sourcePt.Y)
	length := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

	for i := 0; i < int(length); i++ {
		x := float64(sourcePt.X) + dx/length*float64(i)
		y := float64(sourcePt.Y) + dy/length*float64(i)

		for j := 0; j < l.LineWidth; j++ {
			u := float64(j) - float64(l.LineWidth-1)/2
			wx := dy / length * u
			wy := -dx / length * u
			l.drawNeighborsDot(img, x+wx, y+wy)
		}
	}
	l.drawn = true
}
