// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"errors"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	fontPath "github.com/awslabs/diagram-as-code/internal/font"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type LINK_LABEL_TYPE int

const (
	LINK_LABEL_TYPE_HORIZONTAL LINK_LABEL_TYPE = iota
	// [TODO] LINK_LABEL_TYPE_ALONG_PATH
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
	Labels          LinkLabels
	drawn           bool
	lineColor       color.RGBA
}

type LinkLabels struct {
	SourceRight *LinkLabel
	SourceLeft  *LinkLabel
	TargetRight *LinkLabel
	TargetLeft  *LinkLabel
}

type LinkLabel struct {
	Type  LINK_LABEL_TYPE
	Title string
	Color *color.RGBA
	Font  string
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

func (l *Link) SetLineStyle(s string) {
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

		if l.LineStyle == "dashed" && i%9 > 5 {
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

func (l *Link) prepareFontFace(label *LinkLabel, parent1, parent2 *Resource) font.Face {
	if label.Font == "" {
		if parent1 != nil && parent1.labelFont != "" {
			label.Font = parent1.labelFont
		} else if parent2 != nil && parent2.labelFont != "" {
			label.Font = parent2.labelFont
		} else {
			for _, x := range fontPath.Paths {
				if _, err := os.Stat(x); !errors.Is(err, os.ErrNotExist) {
					label.Font = x
					break
				}
			}
		}
	}
	if label.Color == nil {
		if parent1 != nil && parent1.labelColor != nil {
			label.Color = parent1.labelColor
		} else if parent2 != nil && parent2.labelFont != "" {
			label.Color = parent2.labelColor
		} else {
			label.Color = &color.RGBA{0, 0, 0, 255}
		}
	}
	if label.Font == "" {
		panic("Specified fonts are not installed.")
	}
	f, err := os.Open(label.Font)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ttfBytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	ft, err := truetype.Parse(ttfBytes)
	if err != nil {
		panic(err)
	}

	opt := truetype.Options{
		Size:              24,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	return truetype.NewFace(ft, &opt)
}

func (l *Link) computeLabelPos(tx, ty, dx, dy, lx, ly float64) (float64, float64) {
	// Compute the dot product of the unit vectors
	dot_product := tx*dx + ty*dy
	// If the angle is 90 degrees or more (dot product <= 0), set a to (0,0)
	if dot_product > 0 {
		// Compute scalar α
		numerator := ly*dx - lx*dy
		denominator := tx*dy - ty*dx
		// Check for division by zero
		if denominator != 0 {
			alpha := numerator / denominator
			// Compute vector a
			return alpha * tx, alpha * ty
		}
	}
	return 0.0, 0.0
}

func (l *Link) drawLabel(img *image.RGBA, pos Windrose, source, target *Resource, sourcePt, targetPt image.Point, side string, label *LinkLabel) {
	if label == nil {
		return
	}
	_dx := float64(targetPt.X - sourcePt.X)
	_dy := float64(targetPt.Y - sourcePt.Y)
	length := math.Sqrt(math.Pow(_dx, 2) + math.Pow(_dy, 2))
	dx := _dx / length
	dy := _dy / length
	fourWindrose := ((pos + 2) % 16) / 4
	isCorner := ((pos+2)%16)%4 == 0

	_tx := [4]float64{1.0, 0.0, -1.0, 0.0}
	_ty := [4]float64{0.0, 1.0, 0.0, -1.0}

	if isCorner && side == "Left" {
		fourWindrose = (fourWindrose + 3) % 4
	}

	tx := _tx[fourWindrose]
	ty := _ty[fourWindrose]
	if side == "Left" {
		tx = _tx[fourWindrose] * -1
		ty = _ty[fourWindrose] * -1
	}
	// calculate text box size
	textWidth := 0
	textHeight := 0
	fontFace := l.prepareFontFace(label, source, target)
	texts := strings.Split(label.Title, "\n")
	for _, line := range texts {
		textBindings, _ := font.BoundString(fontFace, line)
		textWidth = max(textWidth, textBindings.Max.X.Ceil()-textBindings.Min.X.Ceil())
		textHeight += textBindings.Max.Y.Ceil() - textBindings.Min.Y.Ceil()
	}

	// label vector
	ldx := _tx[(fourWindrose+3)%4] * float64(textWidth)
	ldy := _ty[(fourWindrose+3)%4] * float64(textHeight)

	// calculate the base of a right triangle
	px, py := l.computeLabelPos(tx, ty, dx, dy, ldx, ldy)

	// Unit vector for sliding textbox
	ltx := [4]float64{0.0, 0.0, -1.0, -1.0}
	lty := [4]float64{0.0, 1.0, 1.0, 0.0}

	// calculate buffer
	mx := float64(tx) + dx
	my := float64(ty) + dy
	mag := math.Sqrt(mx*mx + my*my)
	if mag == 0 {
		panic("Error: zero length")
	}
	bx := float64(mx) / mag * 5
	by := float64(my) / mag * 5

	lx := float64(sourcePt.X) + px + bx + float64(textWidth)*ltx[fourWindrose]
	ly := float64(sourcePt.Y) + py + by + float64(textHeight)*lty[fourWindrose]

	if side == "Left" {
		lx = float64(sourcePt.X) + px + bx + float64(textWidth)*ltx[(fourWindrose+3)%4]
		ly = float64(sourcePt.Y) + py + by + float64(textHeight)*lty[(fourWindrose+3)%4]
	}

	lineOffset := fixed.I(0)
	for _, line := range texts {
		textBindings, _ := font.BoundString(fontFace, line)
		point := fixed.Point26_6{fixed.I(int(lx)), fixed.I(int(ly)) + lineOffset}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(label.Color),
			Face: fontFace,
			Dot:  point,
		}
		d.DrawString(line)
		lineOffset += lineOffset + textBindings.Max.Y - textBindings.Min.Y
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
		l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, targetPt, "Right", l.Labels.SourceRight)
		l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, targetPt, "Left", l.Labels.SourceLeft)
		l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, sourcePt, "Left", l.Labels.TargetRight)
		l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, sourcePt, "Right", l.Labels.TargetLeft)
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
			l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, controlPts[0], "Right", l.Labels.SourceRight)
			l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, controlPts[0], "Left", l.Labels.SourceLeft)
			for i := 0; i < len(controlPts)-1; i++ {
				l.drawLine(img, controlPts[i], controlPts[i+1])
			}
			l.drawLine(img, controlPts[len(controlPts)-1], targetPt)
			l.drawArrowHead(img, targetPt, controlPts[len(controlPts)-1], l.TargetArrowHead)
			l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, controlPts[len(controlPts)-1], "Left", l.Labels.TargetRight)
			l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, controlPts[len(controlPts)-1], "Right", l.Labels.TargetLeft)
		} else {
			l.drawLine(img, sourcePt, targetPt)
			l.drawArrowHead(img, sourcePt, targetPt, l.SourceArrowHead)
			l.drawArrowHead(img, targetPt, sourcePt, l.TargetArrowHead)
			l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, targetPt, "Right", l.Labels.SourceRight)
			l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, targetPt, "Left", l.Labels.SourceLeft)
			l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, sourcePt, "Left", l.Labels.TargetRight)
			l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, sourcePt, "Right", l.Labels.TargetLeft)
		}
	}
	l.drawn = true
}
