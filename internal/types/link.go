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
	Source          *Resource
	SourcePosition  Windrose
	SourceArrowHead ArrowHead
	Target          *Resource
	TargetPosition  Windrose
	TargetArrowHead ArrowHead
	Type            string
	LineWidth       int
	LineStyle       string
	drawn           bool
	lineColor       color.RGBA
}

type ArrowHead struct {
	Type   string  `yaml:"Type"`
	Length float64 `yaml:"Length"`
	Width  string  `yaml:"Width"`
}

func (l Link) Init(source *Resource, sourcePosition Windrose, sourceArrowHead ArrowHead, target *Resource, targetPosition Windrose, targetArrowHead ArrowHead, lineWidth int, lineColor color.RGBA) *Link {
	gl := Link{}
	gl.Source = source
	gl.SourcePosition = sourcePosition
	gl.SourceArrowHead = sourceArrowHead
	gl.Target = target
	gl.TargetPosition = targetPosition
	gl.TargetArrowHead = targetArrowHead
	gl.Type = ""
	gl.LineWidth = lineWidth
	gl.LineStyle = "normal"
	gl.drawn = false
	gl.lineColor = lineColor
	return &gl
}

func (l *Link) SetType(s string) {
	l.Type = s
}

func (l *Link) SetLineStyle (s string) {
	l.LineStyle = s
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

		if l.LineStyle=="dashed" && i%9 > 5 {
			continue
		}
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

	if l.Type == "" || l.Type == "straight" {
		l.drawLine(img, sourcePt, targetPt)
		l.drawArrowHead(img, sourcePt, targetPt, l.SourceArrowHead)
		l.drawArrowHead(img, targetPt, sourcePt, l.TargetArrowHead)
	} else if l.Type == "orthogonal" {
		controlPts := []image.Point{}

		// Convert 4-wind rose
		sourceFourWindrose := ((l.SourcePosition + 2) % 16) / 4
		targetFourWindrose := ((l.TargetPosition + 2) % 16) / 4
		// 0...vertical, 1...horizontal
		sourceDirection := sourceFourWindrose % 2
		targetDirection := targetFourWindrose % 2
		if sourceDirection != targetDirection {
			// orthogonal vector (default control point: 1)
			if sourceDirection == 1 && targetDirection == 0 {
				controlPts = append(controlPts, image.Point{targetPt.X, sourcePt.Y})
			} else {
				controlPts = append(controlPts, image.Point{sourcePt.X, targetPt.Y})
			}
		} else {
			if sourceFourWindrose == targetFourWindrose {
				// same vector (default control point: 4)
				dx := [4]int{0, 1, 0, -1}
				dy := [4]int{-1, 0, 1, 0}
				ptX := (max(sourcePt.X*dx[sourceFourWindrose], targetPt.X*dx[sourceFourWindrose]) + 64) * dx[sourceFourWindrose]
				ptY := (max(sourcePt.Y*dy[sourceFourWindrose], targetPt.Y*dy[sourceFourWindrose]) + 64) * dy[sourceFourWindrose]
				if sourceDirection == 0 && targetDirection == 0 {
					controlPts = append(controlPts, image.Point{sourcePt.X, ptY})
					controlPts = append(controlPts, image.Point{targetPt.X, ptY})
				}
				if sourceDirection == 1 && targetDirection == 1 {
					controlPts = append(controlPts, image.Point{ptX, sourcePt.Y})
					controlPts = append(controlPts, image.Point{ptX, targetPt.Y})
				}
			} else {
				// inverse vector (default control point: 2)
				if sourceDirection == 1 && targetDirection == 1 {
					controlPts = append(controlPts, image.Point{(sourcePt.X + targetPt.X) / 2, sourcePt.Y})
					controlPts = append(controlPts, image.Point{(sourcePt.X + targetPt.X) / 2, targetPt.Y})
				} else {
					controlPts = append(controlPts, image.Point{sourcePt.X, (sourcePt.Y + targetPt.Y) / 2})
					controlPts = append(controlPts, image.Point{targetPt.X, (sourcePt.Y + targetPt.Y) / 2})
				}
			}
		}
		if len(controlPts) >= 1 {
			l.drawLine(img, sourcePt, controlPts[0])
			l.drawArrowHead(img, sourcePt, controlPts[0], l.SourceArrowHead)
			for i := 0; i < len(controlPts)-1; i++ {
				l.drawLine(img, controlPts[i], controlPts[i+1])
			}
			l.drawLine(img, controlPts[len(controlPts)-1], targetPt)
			l.drawArrowHead(img, targetPt, controlPts[len(controlPts)-1], l.TargetArrowHead)

		} else {
			l.drawLine(img, sourcePt, targetPt)
			l.drawArrowHead(img, sourcePt, targetPt, l.SourceArrowHead)
			l.drawArrowHead(img, targetPt, sourcePt, l.TargetArrowHead)
		}
	}
	l.drawn = true
}
