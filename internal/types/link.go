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
	Source          *Node
	SourcePosition  string
	SourceArrowHead ArrowHead
	Target          *Node
	TargetPosition  string
	TargetArrowHead ArrowHead
	LineWidth       int
	drawn           bool
	lineColor       color.RGBA
}

type ArrowHead struct {
	Type   string  `yaml:"Type"`
	Length float64 `yaml:"Length"`
	Width  string  `yaml:"Width"`
}

func (l Link) Init(source *Node, sourcePosition string, sourceArrowHead ArrowHead, target *Node, targetPosition string, targetArrowHead ArrowHead, lineWidth int) *Link {
	gl := Link{}
	gl.Source = source
	gl.SourcePosition = sourcePosition
	gl.SourceArrowHead = sourceArrowHead
	gl.Target = target
	gl.TargetPosition = targetPosition
	gl.TargetArrowHead = targetArrowHead
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

func (l *Link) drawLine(img *image.RGBA, sourcePt image.Point, targetPt image.Point) {
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
}

func (l *Link) getThreeSide(t string) (float64, float64, float64) {
	switch t {
	case "Narrow":
		return math.Sqrt(3.0), 2.0, 1.0
	case "Default", "":
		return 1.0, math.Sqrt(2.0), 1.0
	case "Wide":
		return 1.0, 2.0, math.Sqrt(3.0)
	}
	return 0, 0, 0
}

func (l *Link) drawArrowHead(img *image.RGBA, arrowPt image.Point, originPt image.Point, arrowHead ArrowHead) {
	dx := float64(arrowPt.X - originPt.X)
	dy := float64(arrowPt.Y - originPt.Y)
	length := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	if arrowHead.Length == 0 {
		arrowHead.Length = 10
	}
	log.Infof("arrowHead.Length:\"%v\", arrowHead.Width:\"%v\"", arrowHead.Length, arrowHead.Width)
	_a, _b, _c := l.getThreeSide(arrowHead.Width)
	at1 := arrowPt.Sub(image.Point{
		int(arrowHead.Length * (_a*dx - _c*dy) / (_b * length)),
		int(arrowHead.Length * (_c*dx + _a*dy) / (_b * length)),
	})
	at2 := arrowPt.Sub(image.Point{
		int(arrowHead.Length * (_a*dx + _c*dy) / (_b * length)),
		int(arrowHead.Length * (-_c*dx + _a*dy) / (_b * length)),
	})

	switch arrowHead.Type {
	case "Default":
		log.Info("Default Arrow Head drawing")
		al := int(arrowHead.Length)
		for i := 0; i < al; i++ {
			a := arrowPt.Mul(i)
			b := a.Add(at1.Mul(al - i)).Div(al)
			c := a.Add(at2.Mul(al - i)).Div(al)
			l.drawLine(img, b, c)
		}
	case "Open":
		log.Info("Open Arrow Head drawing")
		l.drawLine(img, arrowPt, at1)
		l.drawLine(img, arrowPt, at2)
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

	l.drawLine(img, sourcePt, targetPt)

	l.drawArrowHead(img, sourcePt, targetPt, l.SourceArrowHead)
	l.drawArrowHead(img, targetPt, sourcePt, l.TargetArrowHead)

	l.drawn = true
}
