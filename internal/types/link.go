// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	fontPath "github.com/awslabs/diagram-as-code/internal/font"
	"github.com/awslabs/diagram-as-code/internal/vector"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

type LINK_LABEL_TYPE int

const (
	LINK_LABEL_TYPE_HORIZONTAL LINK_LABEL_TYPE = iota
	// [TODO] LINK_LABEL_TYPE_ALONG_PATH
)

// Segment represents a line segment for overlap detection
type Segment struct {
	X1, Y1, X2, Y2 int   // Start and end coordinates
	Link           *Link // Link that created this segment
	AppliedOffset  int   // Offset that was applied to this segment
}

// Global slice to track all convergence segments
var allConvergenceSegments []Segment

// ResetConvergencePointSegments clears the global segment tracking
func ResetConvergencePointSegments() {
	allConvergenceSegments = []Segment{}
}

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
	AutoRight   *LinkLabel
	AutoLeft    *LinkLabel
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

// ResolveAutoPositions converts WINDROSE_AUTO to actual positions after layout is complete
func (l *Link) ResolveAutoPositions() error {
	if l.SourcePosition == WINDROSE_AUTO || l.TargetPosition == WINDROSE_AUTO {
		// Check for uninitialized resources
		if l.Source.bindings == nil || l.Target.bindings == nil {
			return fmt.Errorf("cannot calculate auto-positions for link: source or target resource is not properly initialized (possibly orphaned from Canvas hierarchy)")
		}

		log.Info("Resolving auto-positions after layout")
		autoSourcePos, autoTargetPos := AutoCalculatePositions(l.Source, l.Target)

		if l.SourcePosition == WINDROSE_AUTO {
			l.SourcePosition = autoSourcePos
		}
		if l.TargetPosition == WINDROSE_AUTO {
			l.TargetPosition = autoTargetPos
		}
	}
	return nil
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
	sourceVec := vector.New(float64(sourcePt.X), float64(sourcePt.Y))
	targetVec := vector.New(float64(targetPt.X), float64(targetPt.Y))
	direction := targetVec.Sub(sourceVec)
	length := direction.Length()

	if length == 0 {
		return
	}

	unitDir := direction.Normalize()
	perpDir := unitDir.Perpendicular()

	for i := 0; i < int(length); i++ {
		pos := sourceVec.Add(unitDir.Scale(float64(i)))

		if l.LineStyle == "dashed" && i%9 > 5 {
			continue
		}
		for j := 0; j < l.LineWidth; j++ {
			offset := float64(j) - float64(l.LineWidth-1)/2
			finalPos := pos.Add(perpDir.Scale(offset))
			l.drawNeighborsDot(img, finalPos.X, finalPos.Y)
		}
	}
}

func (l *Link) prepareFontFace(label *LinkLabel, parent1, parent2 *Resource) (font.Face, error) {
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
	var ttfBytes []byte
	if label.Font == "goregular" || label.Font == "" {
		// Use Go-fonts instead system fonts
		ttfBytes = goregular.TTF
	} else {
		f, err := os.Open(label.Font)
		if err != nil {
			return nil, fmt.Errorf("failed to open font file: %w", err)
		}
		defer func() {
			if closeErr := f.Close(); closeErr != nil {
				log.Warnf("Failed to close font file: %v", closeErr)
			}
		}()

		ttfBytes, err = io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("failed to read font file: %w", err)
		}
	}

	ft, err := truetype.Parse(ttfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	opt := truetype.Options{
		Size:              24,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	return truetype.NewFace(ft, &opt), nil
}

func (l *Link) computeLabelPos(t, d, label vector.Vector) vector.Vector {
	// Compute the dot product of the unit vectors
	dotProduct := t.Dot(d)
	// If the angle is 90 degrees or more (dot product <= 0), set a to (0,0)
	if dotProduct > 0 {
		// Compute scalar α using cross product for 2D
		numerator := label.Y*d.X - label.X*d.Y
		denominator := t.X*d.Y - t.Y*d.X // Fixed order to match original
		// Check for division by zero
		if denominator != 0 {
			alpha := numerator / denominator
			// Compute vector a
			return t.Scale(alpha)
		}
	}
	return vector.New(0.0, 0.0)
}

func (l *Link) drawLabel(img *image.RGBA, pos Windrose, source, target *Resource, sourcePt, targetPt image.Point, side string, label *LinkLabel) error {
	if label == nil {
		return nil
	}
	sourceVec := vector.New(float64(sourcePt.X), float64(sourcePt.Y))
	targetVec := vector.New(float64(targetPt.X), float64(targetPt.Y))
	direction := targetVec.Sub(sourceVec).Normalize()

	fourWindrose := ((pos + 2) % 16) / 4
	isCorner := ((pos+2)%16)%4 == 0

	_tx := [4]float64{1.0, 0.0, -1.0, 0.0}
	_ty := [4]float64{0.0, 1.0, 0.0, -1.0}

	if isCorner && side == "Left" {
		fourWindrose = (fourWindrose + 3) % 4
	}

	t := vector.New(_tx[fourWindrose], _ty[fourWindrose])
	if side == "Left" {
		t = t.Scale(-1)
	}

	// calculate text box size
	textWidth := 0
	textHeight := 0
	fontFace, err := l.prepareFontFace(label, source, target)
	if err != nil {
		return fmt.Errorf("failed to prepare font face for link label: %w", err)
	}
	texts := strings.Split(label.Title, "\n")
	for _, line := range texts {
		textBindings, _ := font.BoundString(fontFace, line)
		textWidth = max(textWidth, textBindings.Max.X.Ceil()-textBindings.Min.X.Ceil())
		textHeight += textBindings.Max.Y.Ceil() - textBindings.Min.Y.Ceil()
	}

	// label vector
	labelVec := vector.New(_tx[(fourWindrose+3)%4]*float64(textWidth), _ty[(fourWindrose+3)%4]*float64(textHeight))

	// calculate the base of a right triangle
	p := l.computeLabelPos(t, direction, labelVec)

	// Unit vector for sliding textbox
	ltx := [4]float64{0.0, 0.0, -1.0, -1.0}
	lty := [4]float64{0.0, 1.0, 1.0, 0.0}

	// calculate buffer
	m := t.Add(direction)
	if !m.IsZero() {
		b := m.Normalize().Scale(5)

		l := sourceVec.Add(p).Add(b).Add(vector.New(float64(textWidth)*ltx[fourWindrose], float64(textHeight)*lty[fourWindrose]))

		if side == "Left" {
			l = sourceVec.Add(p).Add(b).Add(vector.New(float64(textWidth)*ltx[(fourWindrose+3)%4], float64(textHeight)*lty[(fourWindrose+3)%4]))
		}

		lineOffset := fixed.I(0)
		for _, line := range texts {
			textBindings, _ := font.BoundString(fontFace, line)
			point := fixed.Point26_6{fixed.I(int(l.X)), fixed.I(int(l.Y)) + lineOffset}
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
	return nil
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
	arrowVec := vector.New(float64(arrowPt.X), float64(arrowPt.Y))
	originVec := vector.New(float64(originPt.X), float64(originPt.Y))
	direction := arrowVec.Sub(originVec)
	length := direction.Length()

	if arrowHead.Length == 0 {
		arrowHead.Length = 10
	}
	log.Infof("arrowHead.Length:\"%v\", arrowHead.Width:\"%v\"", arrowHead.Length, arrowHead.Width)
	_a, _b, _c := l.getThreeSide(arrowHead.Width)

	// Calculate final positions in floating point for better accuracy
	dx := direction.X
	dy := direction.Y

	// Calculate final arrow head positions (not offsets)
	at1Vec := arrowVec.Sub(vector.New(
		arrowHead.Length*(_a*dx-_c*dy)/(_b*length),
		arrowHead.Length*(_c*dx+_a*dy)/(_b*length),
	))
	at2Vec := arrowVec.Sub(vector.New(
		arrowHead.Length*(_a*dx+_c*dy)/(_b*length),
		arrowHead.Length*(-_c*dx+_a*dy)/(_b*length),
	))

	// Convert to int with rounding for better symmetry
	at1 := image.Point{int(math.Round(at1Vec.X)), int(math.Round(at1Vec.Y))}
	at2 := image.Point{int(math.Round(at2Vec.X)), int(math.Round(at2Vec.Y))}

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

func (l *Link) Draw(img *image.RGBA) error {
	source := *l.Source
	target := *l.Target
	if l.drawn {
		log.Info("Link already drawn")
		return nil
	}

	log.Info("Link Drawing")
	sourcePt := l.calcPositionWithOffset(source.GetBindings(), l.SourcePosition, l.Source, true)
	targetPt := l.calcPositionWithOffset(target.GetBindings(), l.TargetPosition, l.Target, false)

	if l.Type == "" || l.Type == "straight" {
		l.drawLine(img, sourcePt, targetPt)
		l.drawArrowHead(img, sourcePt, targetPt, l.SourceArrowHead)
		l.drawArrowHead(img, targetPt, sourcePt, l.TargetArrowHead)
		if err := l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, targetPt, "Right", l.Labels.SourceRight); err != nil {
			return fmt.Errorf("failed to draw source right label: %w", err)
		}
		if err := l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, targetPt, "Left", l.Labels.SourceLeft); err != nil {
			return fmt.Errorf("failed to draw source left label: %w", err)
		}
		if err := l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, sourcePt, "Left", l.Labels.TargetRight); err != nil {
			return fmt.Errorf("failed to draw target right label: %w", err)
		}
		if err := l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, sourcePt, "Right", l.Labels.TargetLeft); err != nil {
			return fmt.Errorf("failed to draw target left label: %w", err)
		}
	} else if l.Type == "orthogonal" {
		controlPts := l.calculateOrthogonalPath(sourcePt, targetPt)

		// Draw the path
		if len(controlPts) >= 1 {
			l.drawLine(img, sourcePt, controlPts[0])
			l.drawArrowHead(img, sourcePt, controlPts[0], l.SourceArrowHead)
			if err := l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, controlPts[0], "Right", l.Labels.SourceRight); err != nil {
				return fmt.Errorf("failed to draw source right label: %w", err)
			}
			if err := l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, controlPts[0], "Left", l.Labels.SourceLeft); err != nil {
				return fmt.Errorf("failed to draw source left label: %w", err)
			}
			for i := 0; i < len(controlPts)-1; i++ {
				l.drawLine(img, controlPts[i], controlPts[i+1])
			}
			l.drawLine(img, controlPts[len(controlPts)-1], targetPt)
			l.drawArrowHead(img, targetPt, controlPts[len(controlPts)-1], l.TargetArrowHead)
			if err := l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, controlPts[len(controlPts)-1], "Left", l.Labels.TargetRight); err != nil {
				return fmt.Errorf("failed to draw target right label: %w", err)
			}
			if err := l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, controlPts[len(controlPts)-1], "Right", l.Labels.TargetLeft); err != nil {
				return fmt.Errorf("failed to draw target left label: %w", err)
			}
		} else {
			l.drawLine(img, sourcePt, targetPt)
			l.drawArrowHead(img, sourcePt, targetPt, l.SourceArrowHead)
			l.drawArrowHead(img, targetPt, sourcePt, l.TargetArrowHead)
			if err := l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, targetPt, "Right", l.Labels.SourceRight); err != nil {
				return fmt.Errorf("failed to draw source right label: %w", err)
			}
			if err := l.drawLabel(img, l.SourcePosition, l.Source, l.Target, sourcePt, targetPt, "Left", l.Labels.SourceLeft); err != nil {
				return fmt.Errorf("failed to draw source left label: %w", err)
			}
			if err := l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, sourcePt, "Left", l.Labels.TargetRight); err != nil {
				return fmt.Errorf("failed to draw target right label: %w", err)
			}
			if err := l.drawLabel(img, l.TargetPosition, l.Target, l.Source, targetPt, sourcePt, "Right", l.Labels.TargetLeft); err != nil {
				return fmt.Errorf("failed to draw target left label: %w", err)
			}
		}

		/* Original orthogonal implementation - commented out for reference
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
		*/
	} else {
		return fmt.Errorf("unknown link type: %s", l.Type)
	}

	// Draw auto-positioned labels
	var controlPts []image.Point
	if l.Type == "orthogonal" {
		controlPts = l.calculateOrthogonalPath(sourcePt, targetPt)
	}

	autoPt1, autoPt2 := l.calculateAutoLabelPoints(sourcePt, targetPt, controlPts)

	// Determine appropriate pos based on line direction
	autoPos := Windrose(4) // Default: East
	if autoPt1.Y == autoPt2.Y {
		// Horizontal line: use East/West for vertical label placement
		if autoPt1.X < autoPt2.X {
			autoPos = 4 // East for left-to-right horizontal line
		} else {
			autoPos = 12 // West for right-to-left horizontal line
		}
	} else if autoPt1.X == autoPt2.X {
		// Vertical line: use North/South for horizontal label placement
		if autoPt1.Y > autoPt2.Y {
			autoPos = 0 // North for top-to-bottom vertical line
		} else {
			autoPos = 8 // South for bottom-to-top vertical line
		}
	}

	// For orthogonal links, check acute angle sides and adjust label placement
	if l.Type == "orthogonal" && len(controlPts) > 0 {
		// Create complete path including source and target points
		fullPath := make([]image.Point, 0, len(controlPts)+2)
		fullPath = append(fullPath, sourcePt)
		fullPath = append(fullPath, controlPts...)
		fullPath = append(fullPath, targetPt)

		// AutoLeft label: find best segment for Left side
		if l.Labels.AutoLeft != nil {
			leftStart, leftEnd, acutePos := l.findBestSegmentForSideWithPosition(fullPath, "Left")
			if leftStart.X != 0 || leftStart.Y != 0 { // valid segment found
				// Calculate autoPos for this specific segment
				leftAutoPos := Windrose(4) // Default: East
				if leftStart.Y == leftEnd.Y {
					// Horizontal segment
					if leftStart.X < leftEnd.X {
						leftAutoPos = 4 // East
					} else {
						leftAutoPos = 12 // West
					}
				} else if leftStart.X == leftEnd.X {
					// Vertical segment
					if leftStart.Y < leftEnd.Y {
						leftAutoPos = 8 // South
					} else {
						leftAutoPos = 0 // North
					}
				}

				// Always use normal placement (no reversal)
				if err := l.drawLabel(img, leftAutoPos, l.Source, l.Target, leftStart, leftEnd, "Left", l.Labels.AutoLeft); err != nil {
					return fmt.Errorf("failed to draw auto left label: %w (acutePos=%s)", err, acutePos)
				}
			}
		}

		// AutoRight label: find best segment for Right side
		if l.Labels.AutoRight != nil {
			rightStart, rightEnd, acutePos := l.findBestSegmentForSideWithPosition(fullPath, "Right")
			if rightStart.X != 0 || rightStart.Y != 0 { // valid segment found
				// Calculate autoPos for this specific segment
				rightAutoPos := Windrose(4) // Default: East
				if rightStart.Y == rightEnd.Y {
					// Horizontal segment
					if rightStart.X < rightEnd.X {
						rightAutoPos = 4 // East
					} else {
						rightAutoPos = 12 // West
					}
				} else if rightStart.X == rightEnd.X {
					// Vertical segment
					if rightStart.Y < rightEnd.Y {
						rightAutoPos = 8 // South
					} else {
						rightAutoPos = 0 // North
					}
				}

				// Always use normal placement (no reversal)
				if err := l.drawLabel(img, rightAutoPos, l.Source, l.Target, rightStart, rightEnd, "Right", l.Labels.AutoRight); err != nil {
					return fmt.Errorf("failed to draw auto right label: %w (acutePos=%s)", err, acutePos)
				}
			}
		}
	} else {
		// Non-orthogonal or no control points, use default placement
		if err := l.drawLabel(img, autoPos, l.Source, l.Target, autoPt1, autoPt2, "Right", l.Labels.AutoRight); err != nil {
			return fmt.Errorf("failed to draw auto right label: %w", err)
		}
		if err := l.drawLabel(img, autoPos, l.Source, l.Target, autoPt1, autoPt2, "Left", l.Labels.AutoLeft); err != nil {
			return fmt.Errorf("failed to draw auto left label: %w", err)
		}
	}

	l.drawn = true
	return nil
}

// calculateOrthogonalPath generates control points using convergent approach
func (l *Link) calculateOrthogonalPath(sourcePt, targetPt image.Point) []image.Point {
	log.Infof("=== Convergent Orthogonal Path Calculation ===")
	log.Infof("Source: %v (Position: %v)", sourcePt, l.SourcePosition)
	log.Infof("Target: %v (Position: %v)", targetPt, l.TargetPosition)

	// Calculate LCA-based midpoint for convergence
	lcaMidpoint := l.calculateLCABasedMidpoint(sourcePt, targetPt)
	log.Infof("LCA-based convergence point: %v", lcaMidpoint)

	// 1. Get direction vectors from positions
	sourceDir := l.getDirectionVector(int(l.SourcePosition))
	targetDir := l.getDirectionVector(int(l.TargetPosition))
	log.Infof("Source direction: %v", sourceDir)
	log.Infof("Target direction: %v", targetDir)

	// 2. Start from resource positions
	sourceVec := vector.New(float64(sourcePt.X), float64(sourcePt.Y))
	targetVec := vector.New(float64(targetPt.X), float64(targetPt.Y))

	sourceCurrent := sourceVec
	targetCurrent := targetVec
	log.Infof("Source start: %v", sourceCurrent)
	log.Infof("Target start: %v", targetCurrent)

	// 3. Check for resource penetration
	remaining := targetVec.Sub(sourceVec)

	// Source penetration: moving opposite to source direction
	sourcePenetration := sourceDir.Dot(remaining) < -0.5

	// Target penetration: moving same as target direction (overshooting)
	targetPenetration := targetDir.Dot(remaining) > 0.5

	log.Infof("Penetration check - Source: %v, Target: %v", sourcePenetration, targetPenetration)

	// 4. Check direction relationship between source and target
	isParallel := math.Abs(sourceDir.Dot(targetDir)) > 0.5 // Parallel or opposite directions
	log.Infof("Directions parallel: %v (dot product: %v)", isParallel, sourceDir.Dot(targetDir))

	// 4. Generate convergent path
	sourcePoints := []vector.Vector{}
	targetPoints := []vector.Vector{}

	maxSteps := 4
	if isParallel {
		maxSteps = 5 // Need more steps for parallel directions
	}

	for step := 0; step < maxSteps; step++ {
		log.Infof("Step %d:", step)

		// Calculate remaining distance at step start
		remaining := targetCurrent.Sub(sourceCurrent)
		log.Infof("  Remaining: %v", remaining)

		// Check convergence
		if math.Abs(remaining.X) <= 1.0 && math.Abs(remaining.Y) <= 1.0 {
			log.Infof("  Converged!")
			break
		}

		// Determine movement direction based on position directions
		// Source/Target positions determine initial axis preference
		sourceStartsWithX := math.Abs(sourceDir.X) > 0.5 // Horizontal positions (W/E) start with X-axis
		targetStartsWithX := math.Abs(targetDir.X) > 0.5 // Horizontal positions (W/E) start with X-axis

		// Alternating pattern for each source/target
		sourceUseX := (step%2 == 0) == sourceStartsWithX
		targetUseX := (step%2 == 0) == targetStartsWithX

		log.Infof("  Source use X-axis: %v (starts with X: %v)", sourceUseX, sourceStartsWithX)
		log.Infof("  Target use X-axis: %v (starts with X: %v)", targetUseX, targetStartsWithX)

		// Source movement: detour or normal convergence
		if step == 0 {
			// Step 0: Position direction movement
			// - For detour cases: Move 20px minimum to clear resource boundary
			// - For direct cases: Move efficiently toward target (with 20px minimum)
			// - Always moves in the resource's position direction
			var moveDistance float64
			if sourcePenetration {
				// Source penetration: fixed 20px to clear resource
				moveDistance = 20.0
			} else {
				// Source no penetration: efficient distance
				if sourceUseX {
					moveDistance = math.Abs(remaining.X)
					if isParallel {
						// Check if counterpart (target) will have detour
						counterpartDetour := targetPenetration
						if counterpartDetour {
							// Counterpart has detour: use full distance + 20 for efficiency
							moveDistance = math.Abs(remaining.X) + 20.0
						} else {
							// LCA-based convergence: use LCA midpoint instead of 50%
							moveDistance = math.Abs(float64(lcaMidpoint.X) - sourceCurrent.X)
						}
					} else {
						// Non-parallel case: account for counterpart detour
						if targetPenetration {
							detourDistance := 64.0/2 + 20 // 52px
							moveDistance = math.Abs(remaining.X) - detourDistance
							if moveDistance < 0 {
								moveDistance = 20.0
							}
						}
					}
				} else {
					moveDistance = math.Abs(remaining.Y)
					if isParallel {
						// Check if counterpart (target) will have detour
						counterpartDetour := targetPenetration
						if counterpartDetour {
							// Counterpart has detour: use full distance + 20 for efficiency
							moveDistance = math.Abs(remaining.Y) + 20.0
						} else {
							// LCA-based convergence: use LCA midpoint instead of 50%
							moveDistance = math.Abs(float64(lcaMidpoint.Y) - sourceCurrent.Y)
						}
					} else {
						// Non-parallel case: account for counterpart detour
						if targetPenetration {
							detourDistance := 64.0/2 + 20 // 52px
							moveDistance = math.Abs(remaining.Y) - detourDistance
							if moveDistance < 0 {
								moveDistance = 20.0
							}
						}
					}
				}
				if moveDistance < 20.0 {
					moveDistance = 20.0 // Minimum guarantee
				}
			}
			sourceCurrent = sourceCurrent.Add(sourceDir.Scale(moveDistance))
			log.Infof("  Source position move: %v (distance: %v)", sourceCurrent, moveDistance)
		} else if step == 1 && sourcePenetration {
			// Source detour movement
			detourDistance := 64.0/2 + 20 // Minimum 52px
			if math.Abs(sourceDir.X) > 0.5 {
				// Horizontal position: vertical detour
				// Calculate adaptive distance: max(52px, LCA-based distance)
				adaptiveDistance := math.Abs(remaining.Y)
				if isParallel {
					// LCA-based convergence: use LCA midpoint instead of 50%
					adaptiveDistance = math.Abs(float64(lcaMidpoint.Y) - sourceCurrent.Y)
				}
				if adaptiveDistance > detourDistance {
					detourDistance = adaptiveDistance
				}

				detourOffset := -detourDistance // Default north
				if remaining.Y > 0 {
					detourOffset = detourDistance // South if target is below
				}
				sourceCurrent = sourceCurrent.Add(vector.New(0, detourOffset))
				log.Infof("  Source detour Y-move: %v (distance: %v)", sourceCurrent, detourOffset)
			} else {
				// Vertical position: horizontal detour
				// Calculate adaptive distance: max(52px, LCA-based distance)
				adaptiveDistance := math.Abs(remaining.X)
				if isParallel {
					// LCA-based convergence: use LCA midpoint instead of 50%
					adaptiveDistance = math.Abs(float64(lcaMidpoint.X) - sourceCurrent.X)
				}
				if adaptiveDistance > detourDistance {
					detourDistance = adaptiveDistance
				}

				detourOffset := detourDistance // Default east
				if remaining.X < 0 {
					detourOffset = -detourDistance // West if target is left
				}
				sourceCurrent = sourceCurrent.Add(vector.New(detourOffset, 0))
				log.Infof("  Source detour X-move: %v (distance: %v)", sourceCurrent, detourOffset)
			}
		} else {
			// Normal convergence movement
			if sourceUseX {
				moveDistance := remaining.X
				if isParallel {
					// LCA-based convergence: use LCA midpoint instead of 50%
					moveDistance = float64(lcaMidpoint.X) - sourceCurrent.X
				}
				// Apply minimum distance for step 0
				if step == 0 && math.Abs(moveDistance) < 20.0 {
					moveDistance = math.Copysign(20.0, moveDistance)
				}
				sourceCurrent = vector.New(sourceCurrent.X+moveDistance, sourceCurrent.Y)
				log.Infof("  Source X-move: %v (distance: %v)", sourceCurrent, moveDistance)
			} else {
				moveDistance := remaining.Y
				if isParallel {
					// LCA-based convergence: use LCA midpoint instead of 50%
					moveDistance = float64(lcaMidpoint.Y) - sourceCurrent.Y
				}
				// Apply minimum distance for step 0
				if step == 0 && math.Abs(moveDistance) < 20.0 {
					moveDistance = math.Copysign(20.0, moveDistance)
				}
				sourceCurrent = vector.New(sourceCurrent.X, sourceCurrent.Y+moveDistance)
				log.Infof("  Source Y-move: %v (distance: %v)", sourceCurrent, moveDistance)
			}
		}
		sourcePoints = append(sourcePoints, sourceCurrent)

		// Target movement: detour or normal convergence
		if step == 0 {
			// Step 0: Position direction movement
			// - For detour cases: Move 20px minimum to clear resource boundary
			// - For direct cases: Move efficiently toward source (with 20px minimum)
			// - Always moves in the resource's position direction
			var moveDistance float64
			if targetPenetration {
				// Target penetration: fixed 20px to clear resource
				moveDistance = 20.0
			} else {
				// Target no penetration: efficient distance
				if targetUseX {
					moveDistance = math.Abs(remaining.X)
					if isParallel {
						// Check if counterpart (source) will have detour
						counterpartDetour := sourcePenetration
						if counterpartDetour {
							// Counterpart has detour: use full distance + 20 for efficiency
							moveDistance = math.Abs(remaining.X) + 20.0
						} else {
							// LCA-based convergence: use LCA midpoint instead of 50%
							moveDistance = math.Abs(float64(lcaMidpoint.X) - targetCurrent.X)
						}
					} else {
						// Non-parallel case: account for counterpart detour
						if sourcePenetration {
							detourDistance := 64.0/2 + 20 // 52px
							moveDistance = math.Abs(remaining.X) - detourDistance
							if moveDistance < 0 {
								moveDistance = 20.0
							}
						}
					}
				} else {
					moveDistance = math.Abs(remaining.Y)
					if isParallel {
						// Check if counterpart (source) will have detour
						counterpartDetour := sourcePenetration
						if counterpartDetour {
							// Counterpart has detour: use full distance + 20 for efficiency
							moveDistance = math.Abs(remaining.Y) + 20.0
						} else {
							// LCA-based convergence: use LCA midpoint instead of 50%
							moveDistance = math.Abs(float64(lcaMidpoint.Y) - targetCurrent.Y)
						}
					} else {
						// Non-parallel case: account for counterpart detour
						if sourcePenetration {
							detourDistance := 64.0/2 + 20 // 52px
							moveDistance = math.Abs(remaining.Y) - detourDistance
							if moveDistance < 0 {
								moveDistance = 20.0
							}
						}
					}
				}
				if moveDistance < 20.0 {
					moveDistance = 20.0 // Minimum guarantee
				}
			}
			targetCurrent = targetCurrent.Add(targetDir.Scale(moveDistance))
			log.Infof("  Target position move: %v (distance: %v)", targetCurrent, moveDistance)
		} else if step == 1 && targetPenetration {
			// Target detour movement
			detourDistance := 64.0/2 + 20 // Minimum 52px
			if math.Abs(targetDir.X) > 0.5 {
				// Horizontal position: vertical detour
				// Calculate adaptive distance: max(52px, LCA-based distance)
				adaptiveDistance := math.Abs(remaining.Y)
				if isParallel {
					// LCA-based convergence: use LCA midpoint instead of 50%
					adaptiveDistance = math.Abs(float64(lcaMidpoint.Y) - targetCurrent.Y)
				}
				if adaptiveDistance > detourDistance {
					detourDistance = adaptiveDistance
				}

				detourOffset := -detourDistance // Default north
				if remaining.Y < 0 {            // Inverted: remaining.Y < 0 means Source is above Target
					detourOffset = detourDistance // South if source is above
				}
				targetCurrent = targetCurrent.Add(vector.New(0, detourOffset))
				log.Infof("  Target detour Y-move: %v (distance: %v)", targetCurrent, detourOffset)
			} else {
				// Vertical position: horizontal detour
				// Calculate adaptive distance: max(52px, LCA-based distance)
				adaptiveDistance := math.Abs(remaining.X)
				if isParallel {
					// LCA-based convergence: use LCA midpoint instead of 50%
					adaptiveDistance = math.Abs(float64(lcaMidpoint.X) - targetCurrent.X)
				}
				if adaptiveDistance > detourDistance {
					detourDistance = adaptiveDistance
				}

				detourOffset := detourDistance // Default east
				if remaining.X > 0 {           // Inverted: remaining.X > 0 means Source is right of Target
					detourOffset = -detourDistance // West if source is right
				}
				targetCurrent = targetCurrent.Add(vector.New(detourOffset, 0))
				log.Infof("  Target detour X-move: %v (distance: %v)", targetCurrent, detourOffset)
			}
		} else {
			// Normal convergence movement
			if targetUseX {
				moveDistance := remaining.X
				if isParallel {
					// LCA-based convergence: use LCA midpoint instead of 50%
					moveDistance = targetCurrent.X - float64(lcaMidpoint.X)
				}
				// Apply minimum distance for step 0
				if step == 0 && math.Abs(moveDistance) < 20.0 {
					moveDistance = math.Copysign(20.0, moveDistance)
				}
				targetCurrent = vector.New(targetCurrent.X-moveDistance, targetCurrent.Y)
				log.Infof("  Target X-move: %v (distance: %v)", targetCurrent, moveDistance)
			} else {
				moveDistance := remaining.Y
				if isParallel {
					// LCA-based convergence: use LCA midpoint instead of 50%
					moveDistance = targetCurrent.Y - float64(lcaMidpoint.Y)
				}
				// Apply minimum distance for step 0
				if step == 0 && math.Abs(moveDistance) < 20.0 {
					moveDistance = math.Copysign(20.0, moveDistance)
				}
				targetCurrent = vector.New(targetCurrent.X, targetCurrent.Y-moveDistance)
				log.Infof("  Target Y-move: %v (distance: %v)", targetCurrent, moveDistance)
			}
		}
		targetPoints = append(targetPoints, targetCurrent)
	}

	// 5. Build final control points
	controlPts := []image.Point{}

	// Add source points
	for i, pt := range sourcePoints {
		controlPts = append(controlPts, image.Point{int(math.Round(pt.X)), int(math.Round(pt.Y))})
		log.Infof("Added source point %d: (%d, %d)", i, int(math.Round(pt.X)), int(math.Round(pt.Y)))
	}

	// Add target points in reverse order (excluding duplicates)
	for i := len(targetPoints) - 1; i >= 0; i-- {
		targetPoint := image.Point{int(math.Round(targetPoints[i].X)), int(math.Round(targetPoints[i].Y))}
		log.Infof("Processing target point %d: (%d, %d)", i, targetPoint.X, targetPoint.Y)
		// Skip if duplicate of last control point
		if len(controlPts) > 0 && controlPts[len(controlPts)-1] == targetPoint {
			log.Infof("Skipped duplicate target point %d: (%d, %d)", i, targetPoint.X, targetPoint.Y)
			continue
		}
		controlPts = append(controlPts, targetPoint)
		log.Infof("Added target point %d: (%d, %d)", i, targetPoint.X, targetPoint.Y)
	}

	log.Infof("Final control points: %v", controlPts)
	log.Infof("=== End Convergent Calculation ===")
	return controlPts
}

// AutoCalculatePositions determines optimal source and target positions for a link
func AutoCalculatePositions(source, target *Resource) (sourcePos, targetPos Windrose) {
	sourceBounds := source.GetBindings()
	targetBounds := target.GetBindings()

	// Calculate centers
	sourceCenter := image.Point{
		X: sourceBounds.Min.X + sourceBounds.Dx()/2,
		Y: sourceBounds.Min.Y + sourceBounds.Dy()/2,
	}
	targetCenter := image.Point{
		X: targetBounds.Min.X + targetBounds.Dx()/2,
		Y: targetBounds.Min.Y + targetBounds.Dy()/2,
	}

	log.Infof("Auto-positioning: Source center (%d, %d), Target center (%d, %d)",
		sourceCenter.X, sourceCenter.Y, targetCenter.X, targetCenter.Y)

	// Calculate differences
	dx := targetCenter.X - sourceCenter.X
	dy := targetCenter.Y - sourceCenter.Y

	log.Infof("Auto-positioning: dx=%d, dy=%d", dx, dy)

	// Check for common ancestor layout
	if commonAncestor := findLowestCommonAncestor(source, target); commonAncestor != nil {
		direction := commonAncestor.direction
		log.Infof("Auto-positioning: Found common ancestor with direction=%s", direction)

		// Find LCA direct children that contain source and target
		sourceChild := findChildAncestorInLCA(commonAncestor, source)
		targetChild := findChildAncestorInLCA(commonAncestor, target)

		if sourceChild == nil || targetChild == nil {
			log.Warnf("Could not find LCA child ancestors, falling back to distance-based")
			sourcePos, targetPos = calculateByDistance(dx, dy)
			return sourcePos, targetPos
		}

		// Count resources in each direction (from resource to LCA child)
		sourceCounts := countResourcesInDirections(source, sourceChild)
		targetCounts := countResourcesInDirections(target, targetChild)

		log.Infof("Auto-positioning: Source counts (before adjustment): N=%d, E=%d, W=%d, S=%d",
			sourceCounts.North, sourceCounts.East, sourceCounts.West, sourceCounts.South)
		log.Infof("Auto-positioning: Target counts (before adjustment): N=%d, E=%d, W=%d, S=%d",
			targetCounts.North, targetCounts.East, targetCounts.West, targetCounts.South)

		// Adjust counts based on LCA children relationship
		adjustCountsForLCAChildren(commonAncestor, sourceChild, targetChild, &sourceCounts, &targetCounts)

		log.Infof("Auto-positioning: Source counts (after adjustment): N=%d, E=%d, W=%d, S=%d",
			sourceCounts.North, sourceCounts.East, sourceCounts.West, sourceCounts.South)
		log.Infof("Auto-positioning: Target counts (after adjustment): N=%d, E=%d, W=%d, S=%d",
			targetCounts.North, targetCounts.East, targetCounts.West, targetCounts.South)

		// Select optimal positions based on counts and direction
		sourcePos = selectOptimalPosition(sourceCounts, direction, dx, dy, true)
		targetPos = selectOptimalPosition(targetCounts, direction, -dx, -dy, false)
	} else {
		// No common ancestor: use distance-based logic
		sourcePos, targetPos = calculateByDistance(dx, dy)
	}

	log.Infof("Auto-positioning: Source=%v, Target=%v", sourcePos, targetPos)

	return sourcePos, targetPos
}

// selectOptimalPosition selects the best position based on resource counts and LCA direction
func selectOptimalPosition(counts DirectionCounts, lcaDirection string, dx, dy int, isSource bool) Windrose {
	// Find minimum count
	minCount := counts.North
	if counts.East < minCount {
		minCount = counts.East
	}
	if counts.West < minCount {
		minCount = counts.West
	}
	if counts.South < minCount {
		minCount = counts.South
	}

	// Collect directions with minimum count
	candidates := []Windrose{}
	if counts.North == minCount {
		candidates = append(candidates, WINDROSE_N)
	}
	if counts.East == minCount {
		candidates = append(candidates, WINDROSE_E)
	}
	if counts.West == minCount {
		candidates = append(candidates, WINDROSE_W)
	}
	if counts.South == minCount {
		candidates = append(candidates, WINDROSE_S)
	}

	// If only one candidate, return it
	if len(candidates) == 1 {
		return candidates[0]
	}

	// Multiple candidates: prioritize based on LCA direction
	// Priority: LCA direction (preferred) -> perpendicular (based on dx/dy) -> opposite (avoid)

	if lcaDirection == "vertical" {
		// Preferred: N or S (based on dy)
		// Perpendicular: E or W (based on dx)
		// Opposite: opposite of preferred
		var preferred Windrose
		var opposite Windrose
		if dy > 0 {
			preferred = WINDROSE_S
			opposite = WINDROSE_N
		} else {
			preferred = WINDROSE_N
			opposite = WINDROSE_S
		}

		// Check preferred direction
		if containsWindrose(candidates, preferred) {
			return preferred
		}

		// Check perpendicular (E/W) - choose based on dx
		hasE := containsWindrose(candidates, WINDROSE_E)
		hasW := containsWindrose(candidates, WINDROSE_W)

		if hasE && hasW {
			// Both perpendicular directions available, choose based on dx
			if dx > 0 {
				return WINDROSE_E
			}
			return WINDROSE_W
		} else if hasE {
			return WINDROSE_E
		} else if hasW {
			return WINDROSE_W
		}

		// Last resort: opposite
		if containsWindrose(candidates, opposite) {
			return opposite
		}

	} else if lcaDirection == "horizontal" {
		// Preferred: E or W (based on dx)
		// Perpendicular: N or S (based on dy)
		// Opposite: opposite of preferred
		var preferred Windrose
		var opposite Windrose
		if dx > 0 {
			preferred = WINDROSE_E
			opposite = WINDROSE_W
		} else {
			preferred = WINDROSE_W
			opposite = WINDROSE_E
		}

		// Check preferred direction
		if containsWindrose(candidates, preferred) {
			return preferred
		}

		// Check perpendicular (N/S) - choose based on dy
		hasN := containsWindrose(candidates, WINDROSE_N)
		hasS := containsWindrose(candidates, WINDROSE_S)

		if hasN && hasS {
			// Both perpendicular directions available, choose based on dy
			if dy > 0 {
				return WINDROSE_S
			}
			return WINDROSE_N
		} else if hasN {
			return WINDROSE_N
		} else if hasS {
			return WINDROSE_S
		}

		// Last resort: opposite
		if containsWindrose(candidates, opposite) {
			return opposite
		}

	} else {
		// Unknown direction: fall back to distance-based
		if abs(dx) > abs(dy) {
			if dx > 0 {
				if containsWindrose(candidates, WINDROSE_E) {
					return WINDROSE_E
				}
			} else {
				if containsWindrose(candidates, WINDROSE_W) {
					return WINDROSE_W
				}
			}
		} else {
			if dy > 0 {
				if containsWindrose(candidates, WINDROSE_S) {
					return WINDROSE_S
				}
			} else {
				if containsWindrose(candidates, WINDROSE_N) {
					return WINDROSE_N
				}
			}
		}
	}

	// Fallback: return first candidate
	return candidates[0]
}

// containsWindrose checks if a slice contains a Windrose value
func containsWindrose(slice []Windrose, value Windrose) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// adjustCountsForLCAChildren adjusts counts based on LCA children relationship
// Example: LCA{V1{R1}, V2{R2}, V3{R3}}, R3->R1 link
// LCA children: V3->V1, with V2 in between
// If LCA.direction = "horizontal": V3.West += 1 (V2), V1.East += 1 (V2)
// If LCA.direction = "vertical": V3.North += 1 (V2), V1.South += 1 (V2)
func adjustCountsForLCAChildren(lca, sourceChild, targetChild *Resource, sourceCounts, targetCounts *DirectionCounts) {
	if lca == nil || sourceChild == nil || targetChild == nil {
		return
	}

	// Find indices of source and target children in LCA
	sourceIndex := -1
	targetIndex := -1
	for i, child := range lca.children {
		if child == sourceChild {
			sourceIndex = i
		}
		if child == targetChild {
			targetIndex = i
		}
	}

	if sourceIndex == -1 || targetIndex == -1 || sourceIndex == targetIndex {
		return
	}

	// Count resources between source and target children
	var betweenCount int
	if sourceIndex < targetIndex {
		betweenCount = targetIndex - sourceIndex - 1
	} else {
		betweenCount = sourceIndex - targetIndex - 1
	}

	if betweenCount <= 0 {
		return
	}

	log.Infof("Adjusting counts: %d resources between LCA children (sourceIndex=%d, targetIndex=%d)",
		betweenCount, sourceIndex, targetIndex)

	// Adjust counts based on LCA direction
	if lca.direction == "horizontal" {
		// Horizontal: West->East order
		if sourceIndex < targetIndex {
			// Source is West of Target
			sourceCounts.East += betweenCount
			targetCounts.West += betweenCount
		} else {
			// Source is East of Target
			sourceCounts.West += betweenCount
			targetCounts.East += betweenCount
		}
	} else if lca.direction == "vertical" {
		// Vertical: North->South order
		if sourceIndex < targetIndex {
			// Source is North of Target
			sourceCounts.South += betweenCount
			targetCounts.North += betweenCount
		} else {
			// Source is South of Target
			sourceCounts.North += betweenCount
			targetCounts.South += betweenCount
		}
	}
}

// findLowestCommonAncestor finds the lowest common ancestor using set method
func findLowestCommonAncestor(source, target *Resource) *Resource {
	if source == target {
		return source
	}

	// Step 1: Store all ancestors of source in a set
	ancestors := make(map[*Resource]bool)
	current := source
	for current != nil {
		ancestors[current] = true
		current = current.GetParent()
	}

	// Step 2: Traverse target's ancestors and find first match
	current = target
	for current != nil {
		if _, ok := ancestors[current]; ok {
			return current // First match is LCA
		}
		current = current.GetParent()
	}

	return nil // No common ancestor
}

// calculateByDistance uses the original distance-based logic
func calculateByDistance(dx, dy int) (sourcePos, targetPos Windrose) {
	if abs(dx) > abs(dy) {
		// Horizontal connection
		if dx > 0 {
			sourcePos = WINDROSE_E
			targetPos = WINDROSE_W
		} else {
			sourcePos = WINDROSE_W
			targetPos = WINDROSE_E
		}
	} else {
		// Vertical connection
		if dy > 0 {
			sourcePos = WINDROSE_S
			targetPos = WINDROSE_N
		} else {
			sourcePos = WINDROSE_N
			targetPos = WINDROSE_S
		}
	}
	return sourcePos, targetPos
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DirectionCounts holds the count of resources in each direction
type DirectionCounts struct {
	North int
	East  int
	West  int
	South int
}

// countResourcesInDirections counts resources between target and LCA in each direction
// VerticalStack: North->South order, HorizontalStack: West->East order
func countResourcesInDirections(target, lca *Resource) DirectionCounts {
	counts := DirectionCounts{}

	if target == lca {
		return counts
	}

	// Traverse from target to LCA
	current := target
	for current != nil && current != lca {
		parent := current.GetParent()
		if parent == nil {
			break
		}

		// Find current's index in parent's children
		currentIndex := -1
		for i, child := range parent.children {
			if child == current {
				currentIndex = i
				break
			}
		}

		if currentIndex == -1 {
			current = parent
			continue
		}

		// Count siblings based on parent's direction
		if parent.direction == "vertical" {
			// VerticalStack: North->South order
			counts.North += currentIndex
			counts.South += len(parent.children) - currentIndex - 1
		} else if parent.direction == "horizontal" {
			// HorizontalStack: West->East order
			counts.West += currentIndex
			counts.East += len(parent.children) - currentIndex - 1
		}

		current = parent
	}

	return counts
}

// calculateLCABasedMidpoint calculates midpoint using LCA information
func (l *Link) calculateLCABasedMidpoint(sourcePt, targetPt image.Point) image.Point {
	log.Infof("=== LCA-based Midpoint Calculation ===")
	log.Infof("Source: %v, Target: %v", sourcePt, targetPt)

	// Find LCA between source and target
	lca := findLowestCommonAncestor(l.Source, l.Target)
	if lca == nil {
		log.Infof("No LCA found, using default 50%% midpoint")
		return image.Point{
			X: (sourcePt.X + targetPt.X) / 2,
			Y: (sourcePt.Y + targetPt.Y) / 2,
		}
	}

	log.Infof("LCA found: %s (direction: %s)", getResourceName(lca), lca.direction)

	// Get LCA children list
	children := lca.children
	log.Infof("LCA has %d children:", len(children))
	for i, child := range children {
		log.Infof("  Child[%d]: %s", i, getResourceName(child))
	}

	// Find which LCA child each source/target belongs to
	sourceChild := findChildAncestorInLCA(lca, l.Source)
	targetChild := findChildAncestorInLCA(lca, l.Target)

	if sourceChild != nil {
		log.Infof("Source %s belongs to LCA child: %s", getResourceName(l.Source), getResourceName(sourceChild))
	} else {
		log.Infof("Source %s: no LCA child found", getResourceName(l.Source))
	}

	if targetChild != nil {
		log.Infof("Target %s belongs to LCA child: %s", getResourceName(l.Target), getResourceName(targetChild))
	} else {
		log.Infof("Target %s: no LCA child found", getResourceName(l.Target))
	}

	// Find indices of source and target children in LCA children list
	sourceIndex := -1
	targetIndex := -1
	for i, child := range children {
		if child == sourceChild {
			sourceIndex = i
		}
		if child == targetChild {
			targetIndex = i
		}
	}

	log.Infof("Source child index: %d, Target child index: %d", sourceIndex, targetIndex)

	// Find previous child before target (without loop)
	if sourceIndex != -1 && targetIndex != -1 && sourceIndex != targetIndex {
		// Initialize midpoint with 50% default
		midPt := image.Point{
			X: (sourcePt.X + targetPt.X) / 2,
			Y: (sourcePt.Y + targetPt.Y) / 2,
		}

		var previousChild *Resource = nil

		if sourceIndex < targetIndex {
			// Forward iteration: previous is targetIndex-1
			if targetIndex > 0 {
				previousChild = children[targetIndex-1]
			}
		} else {
			// Reverse iteration: previous is targetIndex+1
			if targetIndex < len(children)-1 {
				previousChild = children[targetIndex+1]
			}
		}

		if previousChild != nil {
			log.Infof("Previous child before target: %s", getResourceName(previousChild))

			// Get bindings of previous child and target child
			previousBounds := previousChild.GetBindings()
			targetChildBounds := targetChild.GetBindings()

			log.Infof("Previous child bounds: %v", previousBounds)
			log.Infof("Target child bounds: %v", targetChildBounds)

			// Calculate gap-based midpoint based on LCA direction
			if lca.direction == "horizontal" {
				// Horizontal layout: calculate X coordinate between bounds
				var gapX int
				if sourceIndex < targetIndex {
					// Forward: target.Min.X and previous.Max.X
					gapX = (targetChildBounds.Min.X + previousBounds.Max.X) / 2
				} else {
					// Reverse: target.Max.X and previous.Min.X
					gapX = (targetChildBounds.Max.X + previousBounds.Min.X) / 2
				}
				midPt.X = gapX
				log.Infof("Horizontal gap X: %d", gapX)
			} else {
				// Vertical layout: calculate Y coordinate between bounds
				var gapY int
				if sourceIndex < targetIndex {
					// Forward: target.Min.Y and previous.Max.Y
					gapY = (targetChildBounds.Min.Y + previousBounds.Max.Y) / 2
				} else {
					// Reverse: target.Max.Y and previous.Min.Y
					gapY = (targetChildBounds.Max.Y + previousBounds.Min.Y) / 2
				}
				midPt.Y = gapY
				log.Infof("Vertical gap Y: %d", gapY)
			}
		} else {
			log.Infof("No previous child found (at boundary)")
		}

		log.Infof("Calculated midpoint: %v", midPt)

		// Check for segment overlap and apply offset if needed
		midPt = l.applySegmentBasedOffset(midPt, sourcePt, targetPt, lca)

		return midPt
	} else {
		log.Infof("Source and target are same child or invalid indices")
	}

	// Fallback: return 50% midpoint
	midPt := image.Point{
		X: (sourcePt.X + targetPt.X) / 2,
		Y: (sourcePt.Y + targetPt.Y) / 2,
	}
	log.Infof("Calculated midpoint: %v", midPt)

	return midPt
}

// applySegmentBasedOffset checks for segment overlap and applies offset if needed
func (l *Link) applySegmentBasedOffset(midPt, sourcePt, targetPt image.Point, lca *Resource) image.Point {
	// Create the convergence segment based on layout direction
	var newSegment Segment

	if lca.direction == "horizontal" {
		// Horizontal layout: vertical line at X=midPt.X from min(sourcePt.Y, targetPt.Y) to max(sourcePt.Y, targetPt.Y)
		minY := sourcePt.Y
		maxY := targetPt.Y
		if minY > maxY {
			minY, maxY = maxY, minY
		}
		newSegment = Segment{
			X1:   midPt.X,
			Y1:   minY,
			X2:   midPt.X,
			Y2:   maxY,
			Link: l,
		}
		log.Infof("New segment: X=%d, Y range [%d-%d] (Source: %s -> Target: %s)",
			midPt.X, minY, maxY, getResourceName(l.Source), getResourceName(l.Target))
	} else {
		// Vertical layout: horizontal line at Y=midPt.Y from min(sourcePt.X, targetPt.X) to max(sourcePt.X, targetPt.X)
		minX := sourcePt.X
		maxX := targetPt.X
		if minX > maxX {
			minX, maxX = maxX, minX
		}
		newSegment = Segment{
			X1:   minX,
			Y1:   midPt.Y,
			X2:   maxX,
			Y2:   midPt.Y,
			Link: l,
		}
		log.Infof("New segment: Y=%d, X range [%d-%d] (Source: %s -> Target: %s)",
			midPt.Y, minX, maxX, getResourceName(l.Source), getResourceName(l.Target))
	}

	// Check for overlaps with existing segments
	maxOffset := 0
	overlapCount := 0
	for i, existing := range allConvergenceSegments {
		if segmentsOverlap(newSegment, existing) {
			// Check if this overlap should be counted based on GroupingOffset settings
			shouldCount := false
			reason := ""

			// Check if same Target + TargetPosition
			if l.Target == existing.Link.Target && l.TargetPosition == existing.Link.TargetPosition {
				if l.Target.groupingOffset && !l.Target.groupingOffsetDirection {
					shouldCount = true
					reason = "same Target+Position, Target GroupingOffset enabled, GroupingOffsetDirection disabled"
				} else {
					reason = "same Target+Position, Target GroupingOffset or GroupingOffsetDirection condition not met"
				}
			} else if l.Source == existing.Link.Source && l.SourcePosition == existing.Link.SourcePosition {
				// Check if same Source + SourcePosition
				if l.Source.groupingOffset && !l.Source.groupingOffsetDirection {
					shouldCount = true
					reason = "same Source+Position, Source GroupingOffset enabled, GroupingOffsetDirection disabled"
				} else {
					reason = "same Source+Position, Source GroupingOffset or GroupingOffsetDirection condition not met"
				}
			} else {
				// Different Source/Target, always count
				shouldCount = true
				reason = "different Source/Target"
			}

			if shouldCount {
				overlapCount++
				// Track maximum offset from overlapping segments
				if existing.AppliedOffset > maxOffset {
					maxOffset = existing.AppliedOffset
				}
				log.Infof("  Overlap with segment[%d]: X[%d-%d] Y[%d-%d] (%s, offset: %d) ✓ COUNTED",
					i, existing.X1, existing.X2, existing.Y1, existing.Y2, reason, existing.AppliedOffset)
			} else {
				log.Infof("  Overlap with segment[%d]: X[%d-%d] Y[%d-%d] (%s) ✗ NOT COUNTED",
					i, existing.X1, existing.X2, existing.Y1, existing.Y2, reason)
			}
		}
	}

	// Apply offset if overlap detected
	var appliedOffset int
	if overlapCount > 0 {
		// Use max offset + 15 instead of count * 15
		appliedOffset = maxOffset + 15
		if lca.direction == "horizontal" {
			midPt.X += appliedOffset
			newSegment.X1 += appliedOffset
			newSegment.X2 += appliedOffset
			log.Infof("✓ Applied X offset %d (X: %d -> %d) based on max offset %d + 15",
				appliedOffset, midPt.X-appliedOffset, midPt.X, maxOffset)
		} else {
			midPt.Y += appliedOffset
			newSegment.Y1 += appliedOffset
			newSegment.Y2 += appliedOffset
			log.Infof("✓ Applied Y offset %d (Y: %d -> %d) based on max offset %d + 15",
				appliedOffset, midPt.Y-appliedOffset, midPt.Y, maxOffset)
		}
	} else {
		log.Infof("✗ No overlap detected, no offset applied")
	}

	// Set applied offset in segment
	newSegment.AppliedOffset = appliedOffset

	// Register this segment
	allConvergenceSegments = append(allConvergenceSegments, newSegment)

	return midPt
}

// segmentsOverlap checks if two segments overlap (excluding boundary-only contact)
func segmentsOverlap(s1, s2 Segment) bool {
	// Check if both are vertical lines (same X)
	if s1.X1 == s1.X2 && s2.X1 == s2.X2 && s1.X1 == s2.X1 {
		// Both vertical on same X, check Y range overlap
		// Use < instead of <= to exclude boundary-only contact
		return !(s1.Y2 <= s2.Y1 || s2.Y2 <= s1.Y1)
	}

	// Check if both are horizontal lines (same Y)
	if s1.Y1 == s1.Y2 && s2.Y1 == s2.Y2 && s1.Y1 == s2.Y1 {
		// Both horizontal on same Y, check X range overlap
		// Use < instead of <= to exclude boundary-only contact
		return !(s1.X2 <= s2.X1 || s2.X2 <= s1.X1)
	}

	return false
}

// findChildAncestorInLCA finds which direct child of LCA contains the given resource
func findChildAncestorInLCA(lca *Resource, resource *Resource) *Resource {
	// Traverse up from resource until we find a direct child of LCA
	current := resource
	for current != nil {
		parent := current.GetParent()
		if parent == lca {
			// current is a direct child of LCA
			return current
		}
		if parent == nil {
			// Reached root without finding LCA
			break
		}
		current = parent
	}
	return nil
}

func (l *Link) calcPositionWithOffset(bindings image.Rectangle, position Windrose, resource *Resource, isSource bool) image.Point {
	// Adjust bindings for SSE/S/SSW positions to include title height
	bindingsWithTitle := bindings
	if position == WINDROSE_SSE || position == WINDROSE_S || position == WINDROSE_SSW {
		fontFace, err := resource.prepareFontFace(len(resource.children) > 0, nil)
		if err == nil {
			_, titleHeight := resource.calculateTitleSize(fontFace)
			bindingsWithTitle.Max.Y += titleHeight
		}
	}

	pt, _ := calcPosition(bindingsWithTitle, position)

	// Check if grouping offset is enabled for this resource
	if !resource.groupingOffset {
		log.Infof("Grouping offset disabled for resource %p, using original position: (%d, %d)", resource, pt.X, pt.Y)
		return pt
	}

	// Get link count and index from the same position (unified counting)
	index, count := l.getLinkIndexAndCount(resource, position)
	log.Infof("Link offset calculation - Resource: %p, Position: %v, Index: %d, Count: %d",
		resource, position, index, count)

	if count <= 1 && !resource.groupingOffsetDirection {
		log.Infof("Single link, no offset needed - Position: (%d, %d)", pt.X, pt.Y)
		return pt
	}

	var groupingOffset int
	if resource.groupingOffsetDirection {
		// GroupingOffsetDirection: Apply group-based offset only
		// index 0 = smaller coordinate average (minus side)
		// index 1 = larger coordinate average (plus side)
		if index == 0 {
			groupingOffset = -5 // Minus side group
		} else {
			groupingOffset = 5 // Plus side group
		}
	} else {
		// Original GroupingOffset: distribute links within group
		groupingOffset = int((float64(index) - float64(count-1)/2.0) * 10)
	}
	log.Infof("Calculated grouping offset: %d (index=%d, count=%d)", groupingOffset, index, count)

	// Apply offset in perpendicular direction to direction vector
	direction := l.getDirectionVector(int(position))
	perpendicular := direction.Perpendicular()
	offset := perpendicular.Scale(float64(groupingOffset))

	finalPt := image.Point{
		X: pt.X + int(math.Round(offset.X)),
		Y: pt.Y + int(math.Round(offset.Y)),
	}

	log.Infof("Position offset applied - Original: (%d, %d), Direction: %v, Perpendicular: %v, Final: (%d, %d)",
		pt.X, pt.Y, direction, perpendicular, finalPt.X, finalPt.Y)

	return finalPt
}

func (l *Link) getLinkIndexAndCount(resource *Resource, position Windrose) (int, int) {
	index := 0
	count := 0

	// Determine if current link is incoming or outgoing for this resource
	currentIsIncoming := (l.Target == resource)

	for _, link := range resource.links {
		var linkPosition Windrose
		var isIncoming bool

		if link.Source == resource {
			linkPosition = link.SourcePosition
			isIncoming = false // outgoing
		} else if link.Target == resource {
			linkPosition = link.TargetPosition
			isIncoming = true // incoming
		} else {
			continue
		}

		// Check if this link should be counted
		shouldCount := linkPosition == position
		if resource.groupingOffsetDirection {
			// Directional grouping: only count links with same direction
			shouldCount = shouldCount && (isIncoming == currentIsIncoming)
		}

		if shouldCount {
			count++
		}
	}

	// For directional grouping, determine group order based on opposite coordinate averages
	if resource.groupingOffsetDirection && count > 0 {
		// Calculate average coordinates for both groups
		var incomingAvg, outgoingAvg float64
		var incomingCount, outgoingCount int

		for _, link := range resource.links {
			var linkPosition Windrose
			var isIncoming bool

			if link.Source == resource {
				linkPosition = link.SourcePosition
				isIncoming = false
			} else if link.Target == resource {
				linkPosition = link.TargetPosition
				isIncoming = true
			} else {
				continue
			}

			if linkPosition == position {
				var coord float64
				if isIncoming {
					// Incoming: use source coordinate
					if position == WINDROSE_N || position == WINDROSE_S {
						// Vertical direction: use X coordinate
						coord = float64(link.Source.iconBounds.Min.X+link.Source.iconBounds.Max.X) / 2.0
					} else {
						// Horizontal direction: use Y coordinate
						coord = float64(link.Source.iconBounds.Min.Y+link.Source.iconBounds.Max.Y) / 2.0
					}
					incomingAvg += coord
					incomingCount++
				} else {
					// Outgoing: use target coordinate
					if position == WINDROSE_N || position == WINDROSE_S {
						// Vertical direction: use X coordinate
						coord = float64(link.Target.iconBounds.Min.X+link.Target.iconBounds.Max.X) / 2.0
					} else {
						// Horizontal direction: use Y coordinate
						coord = float64(link.Target.iconBounds.Min.Y+link.Target.iconBounds.Max.Y) / 2.0
					}
					outgoingAvg += coord
					outgoingCount++
				}
			}
		}

		if incomingCount > 0 {
			incomingAvg /= float64(incomingCount)
		}
		if outgoingCount > 0 {
			outgoingAvg /= float64(outgoingCount)
		}

		// Determine index based on group order (left/top group gets index 0)
		if currentIsIncoming {
			index = 0
			if outgoingCount > 0 && incomingAvg >= outgoingAvg {
				index = 1
			}
			log.Infof("Found current link at group index %d for position %v (incoming, avg=%.1f)", index, position, incomingAvg)
		} else {
			index = 0
			if incomingCount > 0 && outgoingAvg > incomingAvg {
				index = 1
			}
			log.Infof("Found current link at group index %d for position %v (outgoing, avg=%.1f)", index, position, outgoingAvg)
		}
	} else {
		// Original behavior: sequential indexing
		for _, link := range resource.links {
			var linkPosition Windrose
			if link.Source == resource {
				linkPosition = link.SourcePosition
			} else if link.Target == resource {
				linkPosition = link.TargetPosition
			} else {
				continue
			}

			if linkPosition == position {
				if link == l {
					log.Infof("Found current link at sorted index %d for position %v (unified count)", index, position)
					break
				}
				index++
			}
		}
	}

	return index, count
}

// getDirectionVector converts windrose position to unit direction vector
func (l *Link) getDirectionVector(position int) vector.Vector {
	// Convert windrose position to 4-direction
	fourWindrose := ((position + 2) % 16) / 4

	switch fourWindrose {
	case 0:
		return vector.New(0, -1) // North
	case 1:
		return vector.New(1, 0) // East
	case 2:
		return vector.New(0, 1) // South
	case 3:
		return vector.New(-1, 0) // West
	default:
		return vector.New(0, -1) // Default to North
	}
}

// findLongestHorizontalSegment finds the longest horizontal segment in control points
func (l *Link) findLongestHorizontalSegment(controlPts []image.Point) (start, end image.Point, length int) {
	maxLength := 0
	var bestStart, bestEnd image.Point

	for i := 0; i < len(controlPts)-1; i++ {
		p1, p2 := controlPts[i], controlPts[i+1]
		if p1.Y == p2.Y { // horizontal segment
			segLength := int(math.Abs(float64(p2.X - p1.X)))
			if segLength > maxLength {
				maxLength = segLength
				bestStart, bestEnd = p1, p2
			}
		}
	}
	return bestStart, bestEnd, maxLength
}

// findPerpendicularSegments finds if the segment after the longest horizontal segment forms acute angles
// Returns: hasLeftAcute=true if Left side has acute angle, hasRightAcute=true if Right side has acute angle
func (l *Link) findPerpendicularSegments(controlPts []image.Point) (hasLeftAcute, hasRightAcute bool) {
	if len(controlPts) < 3 {
		return false, false
	}

	// Find the longest horizontal segment index, preferring segments with acute angles
	maxLength := 0
	bestIndex := -1
	bestHasAcute := false

	for i := 0; i < len(controlPts)-1; i++ {
		p1, p2 := controlPts[i], controlPts[i+1]
		if p1.Y == p2.Y { // horizontal segment
			segLength := int(math.Abs(float64(p2.X - p1.X)))

			// Check if this segment has acute angles (before or after)
			hasAcute := false

			// Check segment[i-1] to segment[i] (before the horizontal segment)
			if i > 0 {
				seg1 := vector.New(
					float64(controlPts[i].X-controlPts[i-1].X),
					float64(controlPts[i].Y-controlPts[i-1].Y),
				)
				seg2 := vector.New(
					float64(controlPts[i+1].X-controlPts[i].X),
					float64(controlPts[i+1].Y-controlPts[i].Y),
				)

				seg1Norm := seg1.Normalize()
				seg2Norm := seg2.Normalize()
				crossProduct := seg1Norm.Cross(seg2Norm)
				dotProduct := seg1Norm.Dot(seg2Norm)

				// Check if it's a 90-degree turn (perpendicular) with acute angles
				if math.Abs(dotProduct) < 0.1 && math.Abs(crossProduct) > 0.1 {
					hasAcute = true
				}
			}

			// Check segment[i] to segment[i+1] (after the horizontal segment)
			if !hasAcute && i+1 < len(controlPts)-1 {
				seg1 := vector.New(
					float64(controlPts[i+1].X-controlPts[i].X),
					float64(controlPts[i+1].Y-controlPts[i].Y),
				)
				seg2 := vector.New(
					float64(controlPts[i+2].X-controlPts[i+1].X),
					float64(controlPts[i+2].Y-controlPts[i+1].Y),
				)

				seg1Norm := seg1.Normalize()
				seg2Norm := seg2.Normalize()
				crossProduct := seg1Norm.Cross(seg2Norm)
				dotProduct := seg1Norm.Dot(seg2Norm)

				// Check if it's a 90-degree turn (perpendicular) with acute angles
				if math.Abs(dotProduct) < 0.1 && math.Abs(crossProduct) > 0.1 {
					hasAcute = true
				}
			}

			// Select this segment if:
			// 1. It's longer than current best, OR
			// 2. Same length but this one has acute angles and current best doesn't
			if segLength > maxLength || (segLength == maxLength && hasAcute && !bestHasAcute) {
				maxLength = segLength
				bestIndex = i
				bestHasAcute = hasAcute
			}
		}
	}

	if bestIndex == -1 {
		return false, false
	}

	// Check segment[n] to segment[n+1] for acute angles
	if bestIndex+1 < len(controlPts)-1 {
		// segment[n]: controlPts[bestIndex] -> controlPts[bestIndex+1]
		// segment[n+1]: controlPts[bestIndex+1] -> controlPts[bestIndex+2]
		seg1 := vector.New(
			float64(controlPts[bestIndex+1].X-controlPts[bestIndex].X),
			float64(controlPts[bestIndex+1].Y-controlPts[bestIndex].Y),
		)
		seg2 := vector.New(
			float64(controlPts[bestIndex+2].X-controlPts[bestIndex+1].X),
			float64(controlPts[bestIndex+2].Y-controlPts[bestIndex+1].Y),
		)

		seg1Norm := seg1.Normalize()
		seg2Norm := seg2.Normalize()
		crossProduct := seg1Norm.Cross(seg2Norm)
		dotProduct := seg1Norm.Dot(seg2Norm)

		// Debug output
		log.Infof("Acute check: seg1=%v, seg2=%v, dotProduct=%f, crossProduct=%f", seg1, seg2, dotProduct, crossProduct)

		// Check if it's a 90-degree turn (perpendicular)
		if math.Abs(dotProduct) < 0.1 {
			// Right side acute: crossProduct > 0 (90-degree counterclockwise)
			if crossProduct > 0 {
				hasRightAcute = true
			}
			// Left side acute: crossProduct < 0 (90-degree clockwise, 270-degree counterclockwise)
			if crossProduct < 0 {
				hasLeftAcute = true
			}
		}
	}

	return hasLeftAcute, hasRightAcute
}

// getAcuteAngleSide determines which sides of the longest horizontal segment have acute angles
// Returns: leftIsAcute=true if Left side has acute angle, rightIsAcute=true if Right side has acute angle
func (l *Link) getAcuteAngleSide(controlPts []image.Point) (leftIsAcute, rightIsAcute bool) {
	hasLeftAcute, hasRightAcute := l.findPerpendicularSegments(controlPts)
	return hasLeftAcute, hasRightAcute
}

// findBestSegmentForSide finds the best horizontal segment for a specific side (Left or Right)
// Returns: segment start/end points and whether that side has an acute angle
func (l *Link) findBestSegmentForSide(controlPts []image.Point, side string) (start, end image.Point, hasAcute bool) {
	start, end, acutePos := l.findBestSegmentForSideWithPosition(controlPts, side)
	return start, end, acutePos != ""
}

// findBestSegmentForSideWithPosition finds the best horizontal segment for a specific side
// Returns: segment start/end points and acute angle position ("before", "after", or "")
func (l *Link) findBestSegmentForSideWithPosition(controlPts []image.Point, side string) (start, end image.Point, acutePos string) {
	if len(controlPts) < 3 {
		return image.Point{}, image.Point{}, ""
	}

	maxLength := 0
	bestIndex := -1
	bestAcutePos := ""

	log.Infof("Finding best segment for %s side (with position):", side)

	for i := 0; i < len(controlPts)-1; i++ {
		p1, p2 := controlPts[i], controlPts[i+1]
		if p1.Y != p2.Y { // skip non-horizontal segments
			continue
		}

		segLength := int(math.Abs(float64(p2.X - p1.X)))
		acutePosition := ""

		// Check after: segment[i] -> segment[i+1] (horizontal) -> segment[i+2]
		// If acute angle found, use segment[i+1] (next segment after the turn)
		if i+1 < len(controlPts)-1 {
			seg1 := vector.New(
				float64(controlPts[i+1].X-controlPts[i].X),
				float64(controlPts[i+1].Y-controlPts[i].Y),
			)
			seg2 := vector.New(
				float64(controlPts[i+2].X-controlPts[i+1].X),
				float64(controlPts[i+2].Y-controlPts[i+1].Y),
			)

			seg1Norm := seg1.Normalize()
			seg2Norm := seg2.Normalize()
			crossProduct := seg1Norm.Cross(seg2Norm)
			dotProduct := seg1Norm.Dot(seg2Norm)

			if math.Abs(dotProduct) < 0.1 {
				if side == "Right" && crossProduct > 0 {
					acutePosition = "after"
					log.Infof("  Segment %d: length=%d, %s acute=after (cross=%f)", i, segLength, side, crossProduct)
				} else if side == "Left" && crossProduct < 0 {
					acutePosition = "after"
					log.Infof("  Segment %d: length=%d, %s acute=after (cross=%f)", i, segLength, side, crossProduct)
				}
			}
		}

		// Check before: segment[i-1] -> segment[i] (horizontal) -> segment[i+1]
		// If acute angle found, use segment[i] (current segment)
		if acutePosition == "" && i > 0 {
			seg1 := vector.New(
				float64(controlPts[i].X-controlPts[i-1].X),
				float64(controlPts[i].Y-controlPts[i-1].Y),
			)
			seg2 := vector.New(
				float64(controlPts[i+1].X-controlPts[i].X),
				float64(controlPts[i+1].Y-controlPts[i].Y),
			)

			seg1Norm := seg1.Normalize()
			seg2Norm := seg2.Normalize()
			crossProduct := seg1Norm.Cross(seg2Norm)
			dotProduct := seg1Norm.Dot(seg2Norm)

			if math.Abs(dotProduct) < 0.1 {
				if side == "Right" && crossProduct > 0 {
					acutePosition = "before"
					log.Infof("  Segment %d: length=%d, %s acute=before (cross=%f)", i, segLength, side, crossProduct)
				} else if side == "Left" && crossProduct < 0 {
					acutePosition = "before"
					log.Infof("  Segment %d: length=%d, %s acute=before (cross=%f)", i, segLength, side, crossProduct)
				}
			}
		}

		if acutePosition == "" {
			log.Infof("  Segment %d: length=%d, %s acute=none", i, segLength, side)
		}

		// Select this segment if:
		// 1. It's longer than current best, OR
		// 2. Same length but this one has acute angle and current best doesn't
		hasAcute := acutePosition != ""
		bestHasAcute := bestAcutePos != ""

		if segLength > maxLength || (segLength == maxLength && hasAcute && !bestHasAcute) {
			maxLength = segLength
			bestIndex = i
			bestAcutePos = acutePosition
			log.Infof("  → New best segment: index=%d, length=%d, acutePos=%s", bestIndex, maxLength, bestAcutePos)
		}
	}

	if bestIndex == -1 {
		log.Infof("  No horizontal segment found for %s side", side)
		return image.Point{}, image.Point{}, ""
	}

	// Select segment based on acute angle position:
	// - "before" (n-1 to n has acute): use segment n (current horizontal segment)
	// - "after" (n to n+1 has acute): use segment n+1 (next segment, can be vertical)
	selectedIndex := bestIndex
	if bestAcutePos == "after" && bestIndex+1 < len(controlPts)-1 {
		// Use the next segment (n+1)
		selectedIndex = bestIndex + 1
		log.Infof("  Acute angle is 'after' → selecting next segment (n+1): index=%d, points=(%d,%d)->(%d,%d)",
			selectedIndex, controlPts[selectedIndex].X, controlPts[selectedIndex].Y,
			controlPts[selectedIndex+1].X, controlPts[selectedIndex+1].Y)
	}

	log.Infof("  Final best segment for %s: index=%d, points=(%d,%d)->(%d,%d), acutePos=%s",
		side, selectedIndex, controlPts[selectedIndex].X, controlPts[selectedIndex].Y,
		controlPts[selectedIndex+1].X, controlPts[selectedIndex+1].Y, bestAcutePos)

	return controlPts[selectedIndex], controlPts[selectedIndex+1], bestAcutePos
}

// calculateAutoLabelPoints calculates points for auto label positioning
func (l *Link) calculateAutoLabelPoints(sourcePt, targetPt image.Point, controlPts []image.Point) (image.Point, image.Point) {
	if l.Type == "orthogonal" && len(controlPts) > 0 {
		// Create complete path including source and target points
		fullPath := make([]image.Point, 0, len(controlPts)+2)
		fullPath = append(fullPath, sourcePt)
		fullPath = append(fullPath, controlPts...)
		fullPath = append(fullPath, targetPt)

		start, end, length := l.findLongestHorizontalSegment(fullPath)
		if length > 0 {
			return start, end
		}
	}
	// Fallback to straight line points
	return sourcePt, targetPt
}

// ReorderChildrenByLinks reorders children based on link connections to minimize overlap
func ReorderChildrenByLinks(canvas *Resource, links []*Link) {
	for _, link := range links {
		// Find LCA
		lca := findLowestCommonAncestor(link.Source, link.Target)
		if lca == nil || lca.direction == "" {
			continue
		}

		// Find which LCA children contain source and target
		sourceChild := findChildAncestorInLCA(lca, link.Source)
		targetChild := findChildAncestorInLCA(lca, link.Target)

		if sourceChild == nil || targetChild == nil || sourceChild == targetChild {
			continue
		}

		// Find indices
		sourceIndex := -1
		targetIndex := -1
		for i, child := range lca.children {
			if child == sourceChild {
				sourceIndex = i
			}
			if child == targetChild {
				targetIndex = i
			}
		}

		if sourceIndex == -1 || targetIndex == -1 {
			continue
		}

		// Determine movement direction based on order
		sourceMovesRight := sourceIndex < targetIndex

		// Step 6: Reorder intermediate resources
		// Traverse from source to LCA
		reorderPathToLCA(link.Source, lca, sourceMovesRight)

		// Traverse from target to LCA (opposite direction)
		reorderPathToLCA(link.Target, lca, !sourceMovesRight)

		// Step 7: Reorder at LCA level
		if lca.unorderedChildren {
			// Move targetChild adjacent to sourceChild
			if sourceIndex < targetIndex {
				// Pattern: [*, sourceChild, *, targetChild, *]
				// Move targetChild to right of sourceChild
				moveChildAdjacent(lca, targetChild, sourceChild, true)
			} else {
				// Pattern: [*, targetChild, *, sourceChild, *]
				// Move targetChild to left of sourceChild
				moveChildAdjacent(lca, targetChild, sourceChild, false)
			}
		}
	}
}

// getResourceNames returns a slice of resource labels for logging
func getResourceNames(resources []*Resource) []string {
	names := make([]string, len(resources))
	for i, r := range resources {
		names[i] = r.label
	}
	return names
}

// reorderPathToLCA reorders children along the path from resource to LCA (excluding LCA itself)
func reorderPathToLCA(resource *Resource, lca *Resource, moveToEnd bool) {
	current := resource
	for current != nil && current != lca {
		parent := current.GetParent()
		if parent == nil || parent == lca {
			// Stop before reaching LCA - LCA reordering is handled separately in Step 7
			break
		}

		// Find which child of parent contains current
		var childToMove *Resource
		for _, child := range parent.children {
			if child == current || isAncestorOf(child, current) {
				childToMove = child
				break
			}
		}

		if childToMove != nil && parent.unorderedChildren {
			if parent.direction == lca.direction {
				// Same direction: move to edge based on moveToEnd
				positionStr := "start"
				if moveToEnd {
					positionStr = "end"
				}
				log.Infof("Reordering in %s (same direction): moving %s to %v (moveToEnd=%v)",
					parent.label, childToMove.label, positionStr, moveToEnd)
				if moveToEnd {
					// Move to rightmost/bottom (last position)
					moveChildToPosition(parent, childToMove, len(parent.children)-1)
				} else {
					// Move to leftmost/top (first position)
					moveChildToPosition(parent, childToMove, 0)
				}
			} else {
				// Different direction: move to first position
				log.Infof("Reordering in %s (different direction): moving %s to start",
					parent.label, childToMove.label)
				moveChildToPosition(parent, childToMove, 0)
			}
		}

		current = parent
	}
}

// isAncestorOf checks if ancestor is an ancestor of descendant
func isAncestorOf(ancestor *Resource, descendant *Resource) bool {
	current := descendant
	for current != nil {
		if current == ancestor {
			return true
		}
		current = current.GetParent()
	}
	return false
}

// moveChildToPosition moves a child to the specified position in parent's children list
func moveChildToPosition(parent *Resource, child *Resource, targetPos int) {
	// Find current position
	currentPos := -1
	for i, c := range parent.children {
		if c == child {
			currentPos = i
			break
		}
	}

	if currentPos == -1 {
		log.Warnf("Child %s not found in parent %s", child.label, parent.label)
		return
	}

	if currentPos == targetPos {
		log.Infof("Child %s already at position %d in parent %s", child.label, targetPos, parent.label)
		return
	}

	log.Infof("Moving %s from position %d to %d in %s (before: %v)",
		child.label, currentPos, targetPos, parent.label, getResourceNames(parent.children))

	// Create new slice without the child
	newChildren := make([]*Resource, 0, len(parent.children))
	for i, c := range parent.children {
		if i != currentPos {
			newChildren = append(newChildren, c)
		}
	}

	// Insert at target position
	result := make([]*Resource, 0, len(parent.children))
	if targetPos == 0 {
		result = append(result, child)
		result = append(result, newChildren...)
	} else if targetPos >= len(newChildren) {
		result = append(result, newChildren...)
		result = append(result, child)
	} else {
		result = append(result, newChildren[:targetPos]...)
		result = append(result, child)
		result = append(result, newChildren[targetPos:]...)
	}

	parent.children = result

	log.Infof("After move: %v", getResourceNames(parent.children))
}

// moveChildAdjacent moves child to be adjacent to anchor (either right or left)
func moveChildAdjacent(parent *Resource, child *Resource, anchor *Resource, toRight bool) {
	// Find positions
	childPos := -1
	anchorPos := -1
	for i, c := range parent.children {
		if c == child {
			childPos = i
		}
		if c == anchor {
			anchorPos = i
		}
	}

	if childPos == -1 || anchorPos == -1 {
		return
	}

	// Remove child from current position
	parent.children = append(parent.children[:childPos], parent.children[childPos+1:]...)

	// Recalculate anchor position after removal
	anchorPos = -1
	for i, c := range parent.children {
		if c == anchor {
			anchorPos = i
			break
		}
	}

	// Insert adjacent to anchor
	var targetPos int
	if toRight {
		targetPos = anchorPos + 1
	} else {
		targetPos = anchorPos
	}

	parent.children = append(parent.children[:targetPos], append([]*Resource{child}, parent.children[targetPos:]...)...)
}
