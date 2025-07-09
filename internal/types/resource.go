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

	fontPath "github.com/awslabs/diagram-as-code/internal/font"
	"github.com/golang/freetype/truetype"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const DEBUG_LAYOUT = false

type BORDER_TYPE int

const (
	BORDER_TYPE_STRAIGHT BORDER_TYPE = iota
	BORDER_TYPE_DASHED
)

type ICON_FILL_TYPE int

const (
	ICON_FILL_TYPE_NONE ICON_FILL_TYPE = iota
	ICON_FILL_TYPE_RECT
)

type Resource struct {
	bindings       *image.Rectangle
	iconImage      image.Image
	iconBounds     image.Rectangle
	borderColor    *color.RGBA
	borderType     BORDER_TYPE
	fillColor      color.RGBA
	label          string
	labelFont      string
	labelColor     *color.RGBA
	headerAlign    string // left(default) / center / right
	margin         *Margin
	padding        *Padding
	direction      string
	align          string
	links          []*Link
	children       []*Resource
	borderChildren []*BorderChild
	iconfill       ResourceIconFill
	drawn          bool
}

type ResourceIconFill struct {
	Type  ICON_FILL_TYPE // none(default) / rect
	Color color.RGBA
}

type BorderChild struct {
	Position Windrose
	Resource *Resource
}

func defaultResourceValues(hasChild bool, setIcon bool) Resource {
	if hasChild {
		return Resource{ // resource has children and show as Group
			bindings: &image.Rectangle{
				image.Point{0, 0},
				image.Point{320, 190},
			},
			margin:      &Margin{20, 15, 20, 15},
			padding:     &Padding{20, 45, 20, 45},
			borderColor: &color.RGBA{0, 0, 0, 255},
		}
	} else {
		if setIcon {
			return Resource{ // resource has not children and show as Resource
				bindings: &image.Rectangle{
					image.Point{0, 0},
					image.Point{64, 64},
				},
				margin:      &Margin{30, 100, 30, 100},
				padding:     &Padding{0, 0, 0, 0},
				borderColor: &color.RGBA{0, 0, 0, 0},
			}
		} else {
			return Resource{ // resource has not children and icon, show as TextBox
				bindings: &image.Rectangle{
					image.Point{0, 0},
					image.Point{0, 0},
				},
				margin:      &Margin{0, 0, 0, 0},
				padding:     &Padding{0, 0, 0, 0},
				borderColor: &color.RGBA{0, 0, 0, 0},
			}
		}
	}
}

func (r *Resource) Init() *Resource {
	rr := Resource{}
	rr.bindings = nil
	rr.iconImage = image.NewRGBA(image.Rect(0, 0, 0, 0))
	rr.iconBounds = image.Rect(0, 0, 0, 0)
	rr.iconfill = ResourceIconFill{
		Type:  ICON_FILL_TYPE_NONE,
		Color: color.RGBA{255, 255, 255, 255},
	}
	rr.borderColor = nil
	rr.borderType = BORDER_TYPE_STRAIGHT
	rr.fillColor = color.RGBA{0, 0, 0, 0}
	rr.label = ""
	rr.labelFont = ""
	rr.labelColor = &color.RGBA{0, 0, 0, 255}
	rr.headerAlign = "left"
	rr.margin = nil
	rr.padding = nil
	rr.direction = "horizontal"
	rr.align = "center"
	rr.drawn = false
	return &rr
}

func (r *Resource) LoadIcon(imageFilePath string) error {
	imageFile, err := os.Open(imageFilePath)
	if err != nil {
		return err
	}
	defer imageFile.Close()
	iconImage, _, err := image.Decode(imageFile)
	if err != nil {
		return err
	}
	r.iconBounds = image.Rect(0, 0, 64, 64)
	_b := image.Rect(0, 0, 64, 64)
	r.bindings = &_b
	r.iconImage = iconImage
	return nil
}

func (r *Resource) SetHeaderAlign(align string) {
	r.headerAlign = align
}

func (r *Resource) SetIconBounds(bounds image.Rectangle) {
	r.iconBounds = bounds
}

func (r *Resource) SetBindings(bindings image.Rectangle) {
	r.bindings = &bindings
}

func (r *Resource) GetBindings() image.Rectangle {
	return *r.bindings
}

func (r *Resource) GetMargin() Margin {
	return *r.margin
}

func (r *Resource) GetPadding() Padding {
	return *r.padding
}

func (r *Resource) SetBorderColor(borderColor color.RGBA) {
	r.borderColor = &borderColor
}

func (r *Resource) SetBorderType(borderType BORDER_TYPE) {
	r.borderType = borderType
}

func (r *Resource) SetFillColor(fillColor color.RGBA) {
	r.fillColor = fillColor
}

func (r *Resource) SetLabel(label *string, labelColor *color.RGBA, labelFont *string) {
	if label != nil {
		r.label = *label
	}
	if labelColor != nil {
		r.labelColor = labelColor
	}
	if labelFont != nil {
		r.labelFont = *labelFont
	}
}

func (r *Resource) SetAlign(align string) {
	r.align = align
}

func (r *Resource) SetDirection(direction string) {
	r.direction = direction
}

func (r *Resource) SetIconFill(t ICON_FILL_TYPE, color *color.RGBA) {
	r.iconfill.Type = t
	if color != nil {
		r.iconfill.Color = *color
	}
}

func (r *Resource) AddLink(link *Link) {
	r.links = append(r.links, link)
}

func (r *Resource) AddParent() {
}

func (r *Resource) AddChild(child *Resource) {
	// [TODO] check whether the parent is border children
	if child == nil {
		log.Fatalf("Unknown child. Please see debug logs with -v flag.")
	}
	r.children = append(r.children, child)
}

func (r *Resource) AddBorderChild(borderChild *BorderChild) {
	hasChild := len(borderChild.Resource.children) != 0
	if hasChild {
		panic("Couldn't add Group to Border Childlen")
	}
	r.borderChildren = append(r.borderChildren, borderChild)
}

func (r *Resource) prepareFontFace(hasChild bool, parent *Resource) font.Face {
	if r.labelFont == "" {
		if parent != nil && parent.labelFont != "" {
			r.labelFont = parent.labelFont
		} else {
			for _, x := range fontPath.Paths {
				if _, err := os.Stat(x); !errors.Is(err, os.ErrNotExist) {
					r.labelFont = x
					break
				}
			}
		}
	}
	if r.labelColor == nil {
		if parent != nil && parent.labelColor != nil {
			r.labelColor = parent.labelColor
		} else {
			r.labelColor = &color.RGBA{0, 0, 0, 255}
		}
	}
	if r.labelFont == "" {
		panic("Specified fonts are not installed.")
	}
	f, err := os.Open(r.labelFont)
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
	if hasChild {
		opt.Size = 30
	}

	return truetype.NewFace(ft, &opt)
}

func (r *Resource) Scale(parent *Resource, visited map[*Resource]bool) error {
	log.Infof("Scale %s", r.label)

	if visited == nil {
		visited = make(map[*Resource]bool)
	}
	if visited[r] {
		return fmt.Errorf("Cycle detected in resource tree at %s", r.label)
	}
	visited[r] = true

	var prev *Resource
	b := image.Rectangle{
		image.Point{
			math.MaxInt,
			math.MaxInt,
		},
		image.Point{
			math.MinInt,
			math.MinInt,
		},
	}
	hasChildren := len(r.children) != 0
	hasBorderChildren := len(r.borderChildren) != 0
	hasIcon := r.iconImage.Bounds().Max.X != 0
	log.Infof("hasIcon: %t\n", hasIcon)
	textWidth := 0
	textHeight := 0
	if r.label != "" {
		textHeight = 10
		fontFace := r.prepareFontFace(hasChildren, parent)
		texts := strings.Split(r.label, "\n")
		for _, line := range texts {
			textBindings, _ := font.BoundString(fontFace, line)
			textWidth = max(textWidth, textBindings.Max.X.Floor()-textBindings.Min.X.Ceil()+20)
			textHeight += textBindings.Max.Y.Floor() - textBindings.Min.Y.Ceil() + 10
		}
	}
	if r.bindings == nil {
		r.bindings = defaultResourceValues(hasChildren, hasIcon).bindings
	}
	if r.margin == nil {
		r.margin = defaultResourceValues(hasChildren, hasIcon).margin
		// Expand bindings to fit text size
		if !hasChildren {
			// Resource (no child)
			log.Infof("textHeight: %d\n", textHeight)
			r.margin.Bottom += textHeight
			_m := (textWidth - r.iconBounds.Dx()) / 2
			r.margin.Left = maxInt(r.margin.Left, _m)
			r.margin.Right = maxInt(r.margin.Right, _m)
		}
		if hasChildren && hasBorderChildren {
			addMargin := Margin{}
			_m := defaultResourceValues(false, hasIcon)
			for _, x := range r.borderChildren {
				switch x.Position / 4 {
				case 0:
					addMargin.Top = _m.margin.Top * 2
				case 1:
					addMargin.Right = _m.margin.Right
				case 2:
					addMargin.Bottom = _m.margin.Bottom * 2
				case 3:
					addMargin.Left = _m.margin.Left
				}
			}
			r.margin.Top += addMargin.Top
			r.margin.Right += addMargin.Right
			r.margin.Bottom += addMargin.Bottom
			r.margin.Left += addMargin.Left
		}
	}
	if r.padding == nil {
		r.padding = defaultResourceValues(hasChildren, hasIcon).padding
		if hasChildren && hasBorderChildren {
			addPadding := Padding{}
			_m := defaultResourceValues(true, hasIcon)
			for _, x := range r.borderChildren {
				switch x.Position / 4 {
				case 0:
					addPadding.Top = _m.padding.Top * 2
				case 1:
					addPadding.Right = _m.padding.Right * 2
				case 2:
					addPadding.Bottom = _m.padding.Bottom * 2
				case 3:
					addPadding.Left = _m.padding.Left * 2
				}
			}
			r.padding.Top += addPadding.Top
			r.padding.Right += addPadding.Right
			r.padding.Bottom += addPadding.Bottom
			r.padding.Left += addPadding.Left
		}
	}
	if r.borderColor == nil {
		r.borderColor = defaultResourceValues(hasChildren, hasIcon).borderColor
	}

	// Expand bindings to fit text size
	if hasChildren && r.direction == "vertical" {
		// Group (has child)
		prev = &Resource{
			margin: &Margin{},
			bindings: &image.Rectangle{
				Min: image.Point{
					0,
					0,
				},
				Max: image.Point{
					textWidth + r.iconBounds.Dx() + 30,
					0,
				},
			},
		}
		b = *prev.bindings
	}

	for _, subResource := range r.children {
		err := subResource.Scale(parent, visited)
		if err != nil {
			return err
		}

		bindings := subResource.GetBindings()
		margin := subResource.GetMargin()
		if prev != nil {
			prevBindings := prev.GetBindings()
			prevMargin := prev.GetMargin()
			if r.direction == "horizontal" {
				switch r.align {
				case "top":
					subResource.Translation(
						prevBindings.Max.X+prevMargin.Right+margin.Left-bindings.Min.X,
						prevBindings.Min.Y-prevMargin.Top+margin.Top-bindings.Min.Y,
					)
				case "center":
					subResource.Translation(
						prevBindings.Max.X+prevMargin.Right+margin.Left-bindings.Min.X,
						prevBindings.Min.Y+(prevBindings.Dy()-bindings.Dy())/2-bindings.Min.Y,
					)
				case "bottom":
					subResource.Translation(
						prevBindings.Max.X+prevMargin.Right+margin.Left-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom-margin.Bottom-bindings.Max.Y,
					)
				default:
					log.Fatalf("Unknown align %s in the direction(%s) on %s", r.align, r.direction, r.label)
				}
			} else {
				switch r.align {
				case "left":
					subResource.Translation(
						prevBindings.Min.X-prevMargin.Left+margin.Left-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom+margin.Top-bindings.Min.Y,
					)
				case "center":
					subResource.Translation(
						prevBindings.Min.X+(prevBindings.Dx()-bindings.Dx())/2-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom+margin.Top-bindings.Min.Y,
					)
				case "right":
					subResource.Translation(
						prevBindings.Max.X+prevMargin.Right-margin.Right-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom+margin.Top-bindings.Min.Y,
					)
				default:
					log.Fatalf("Unknown align %s in the direction(%s) on %s", r.align, r.direction, r.label)
				}
			}
		}
		bindings = subResource.GetBindings()
		b.Min.X = minInt(b.Min.X, bindings.Min.X-margin.Left-r.padding.Left)
		headerHeight := maxInt(r.iconBounds.Dy(), textHeight)
		if r.headerAlign == "center" {
			headerHeight = r.iconBounds.Dy() + textHeight
		}
		b.Min.Y = minInt(b.Min.Y, bindings.Min.Y-margin.Top-headerHeight-r.padding.Top)
		b.Max.X = maxInt(b.Max.X, bindings.Max.X+margin.Right+r.padding.Right)
		b.Max.Y = maxInt(b.Max.Y, bindings.Max.Y+margin.Bottom+r.padding.Bottom)
		prev = subResource
	}
	// Expand bindings to fit text size
	if hasChildren && r.direction == "horizontal" {
		// Group (has child)
		if textWidth+r.iconBounds.Dx()+30 > b.Dx() {
			_dx := b.Dx()
			b.Min.X -= (textWidth + r.iconBounds.Dx() + 30 - _dx) / 2
			b.Max.X += (textWidth + r.iconBounds.Dx() + 30 - _dx) / 2
		}
	}
	if b.Min.X != math.MaxInt {
		r.SetBindings(b)
	}
	for _, borderChild := range r.borderChildren {
		pt, err := calcPosition(r.GetBindings(), borderChild.Position)
		if err != nil {
			panic(err)
		}
		err = borderChild.Resource.Scale(parent, visited) // to initialize default values
		if err != nil {
			return err
		}
		bindings := borderChild.Resource.GetBindings()
		borderChild.Resource.Translation(
			pt.X-(bindings.Min.X+bindings.Max.X)/2,
			pt.Y-(bindings.Min.Y+bindings.Max.Y)/2,
		)
	}
	return nil
}

func (r *Resource) Translation(dx, dy int) {
	if r.bindings == nil {
		panic("The resource has no binding.")
	}
	r.bindings = &image.Rectangle{
		image.Point{
			r.bindings.Min.X + dx,
			r.bindings.Min.Y + dy,
		},
		image.Point{
			r.bindings.Max.X + dx,
			r.bindings.Max.Y + dy,
		},
	}
	for _, subResource := range r.children {
		subResource.Translation(dx, dy)
	}
	for _, borderChild := range r.borderChildren {
		borderChild.Resource.Translation(dx, dy)
	}

}

func (r *Resource) ZeroAdjust() {
	r.Translation(-r.bindings.Min.X+r.padding.Left, -r.bindings.Min.Y+r.padding.Top)
}

func (r *Resource) IsDrawn() bool {
	return r.drawn
}

func (r *Resource) Draw(img *image.RGBA, parent *Resource) *image.RGBA {
	if img == nil {
		img = image.NewRGBA(*r.bindings)
	}

	if DEBUG_LAYOUT {
		r.drawMargin(img)
	}
	r.drawFrame(img)
	if DEBUG_LAYOUT {
		r.drawPadding(img)
	}

	rctSrc := r.iconImage.Bounds()
	x := image.Rectangle{r.bindings.Min, r.bindings.Min.Add(image.Point{64, 64})}
	switch r.headerAlign {
	case "left":
	case "center":
		x.Min = x.Min.Add(image.Point{(r.bindings.Dx() - 64) / 2, 0})
		x.Max = x.Max.Add(image.Point{(r.bindings.Dx() - 64) / 2, 0})
	case "right":
		x.Min = x.Min.Add(image.Point{r.bindings.Dx() - 64, 0})
		x.Max = x.Max.Add(image.Point{r.bindings.Dx() - 64, 0})
	}
	if r.iconfill.Type == ICON_FILL_TYPE_RECT {
		for _x := x.Min.X; _x < x.Max.X; _x++ {
			for _y := x.Min.Y; _y < x.Max.Y; _y++ {
				c := img.At(_x, _y)
				img.Set(_x, _y, _blend_color(c, r.iconfill.Color))
			}
		}
	}
	draw.CatmullRom.Scale(img, x, r.iconImage, rctSrc, draw.Over, nil)

	hasIcon := r.iconImage.Bounds().Max.X != 0
	if parent != nil {
		r.drawLabel(img, parent, len(r.children) > 0, hasIcon)
	}

	for _, subResource := range r.children {
		subResource.Draw(img, r)
	}
	for _, borderResource := range r.borderChildren {
		borderResource.Resource.Draw(img, r)
	}
	r.drawn = true
	for _, v := range r.links {
		source := *v.Source
		target := *v.Target
		if source.IsDrawn() && target.IsDrawn() {
			v.Draw(img)
		}
	}
	return img
}

func (r *Resource) drawFrame(img *image.RGBA) {
	x1 := r.bindings.Min.X
	x2 := r.bindings.Max.X
	y1 := r.bindings.Min.Y
	y2 := r.bindings.Max.Y
	for x := x1 - WIDTH + 1; x < x2+WIDTH-1; x++ {
		for y := y1 - WIDTH + 1; y < y2+WIDTH-1; y++ {
			c := img.At(x, y)
			if x <= x1 || x >= x2-1 || y <= y1 || y >= y2-1 {
				// Set border
				switch r.borderType {
				case BORDER_TYPE_STRAIGHT:
					img.Set(x, y, _blend_color(c, r.borderColor))
				case BORDER_TYPE_DASHED:
					if (x+y)%9 <= 5 {
						img.Set(x, y, _blend_color(c, r.borderColor))
					}
				}
			} else {
				// Set background
				img.Set(x, y, _blend_color(c, r.fillColor))
				if DEBUG_LAYOUT {
					img.Set(x, y, _blend_color(c, color.RGBA{255, 255, 255, 255}))
				}
				//img.Set(x, y, fill_color)
			}
		}
	}
}

func (r *Resource) drawPadding(img *image.RGBA) {
	x1 := r.bindings.Min.X
	x2 := r.bindings.Max.X
	y1 := r.bindings.Min.Y
	y2 := r.bindings.Max.Y
	for x := x1 - WIDTH + 1; x < x2+WIDTH-1; x++ {
		for y := y1 - WIDTH + 1; y < y2+WIDTH-1; y++ {
			c := img.At(x, y)
			img.Set(x, y, _blend_color(c, color.RGBA{0, 255, 0, 127}))
		}
	}
	x1 = r.bindings.Min.X + r.padding.Left
	x2 = r.bindings.Max.X - r.padding.Right
	y1 = r.bindings.Min.Y + r.padding.Top
	y2 = r.bindings.Max.Y - r.padding.Bottom
	for x := x1 - WIDTH + 1; x < x2+WIDTH-1; x++ {
		for y := y1 - WIDTH + 1; y < y2+WIDTH-1; y++ {
			c := img.At(x, y)
			img.Set(x, y, _blend_color(c, color.RGBA{255, 255, 255, 255}))
		}
	}

}

func (r *Resource) drawMargin(img *image.RGBA) {
	x1 := r.bindings.Min.X - r.margin.Left
	x2 := r.bindings.Max.X + r.margin.Right
	y1 := r.bindings.Min.Y - r.margin.Top
	y2 := r.bindings.Max.Y + r.margin.Bottom
	for x := x1 - WIDTH + 1; x < x2+WIDTH-1; x++ {
		for y := y1 - WIDTH + 1; y < y2+WIDTH-1; y++ {
			c := img.At(x, y)
			img.Set(x, y, _blend_color(c, color.RGBA{255, 255, 0, 255}))
		}
	}
}

func (r *Resource) drawLabel(img *image.RGBA, parent *Resource, hasChild, hasIcon bool) {
	face := r.prepareFontFace(hasChild, parent)

	texts := strings.Split(r.label, "\n")
	lineOffset := 0

	for _, line := range texts {
		textBindings, _ := font.BoundString(face, line)

		textWidth := textBindings.Max.X.Floor() - textBindings.Min.X.Ceil()
		textHeight := textBindings.Max.Y.Floor() - textBindings.Min.Y.Ceil()

		w := textBindings.Max.X - textBindings.Min.X
		h := textBindings.Max.Y - textBindings.Min.Y + fixed.I(lineOffset)

		p := r.bindings.Min.Add(image.Point{0, r.iconBounds.Max.Y})

		point := fixed.Point26_6{fixed.I(p.X) - (w-fixed.I(r.bindings.Dx()))/2, fixed.I(p.Y+10) + h}
		if hasChild {
			iconHeight := r.iconBounds.Max.Y
			if iconHeight == 0 {
				iconHeight = 64
			}
			padding := maxInt((iconHeight-textHeight)/2, 0)
			switch r.headerAlign {
			case "left":
				p = r.bindings.Min.Add(image.Point{
					r.iconBounds.Max.X + padding,
					iconHeight - padding + lineOffset,
				})
			case "center":
				p = r.bindings.Min.Add(image.Point{
					(r.bindings.Dx() - textWidth) / 2,
					r.iconBounds.Dy() + iconHeight - padding + lineOffset,
				})
			case "right":
				p = r.bindings.Min.Add(image.Point{
					r.iconBounds.Dx() - r.bindings.Dx() - r.iconBounds.Dx() - padding,
					iconHeight - padding + lineOffset,
				})
			}
			point = fixed.Point26_6{fixed.I(p.X), fixed.I(p.Y)}
		}

		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(r.labelColor),
			Face: face,
			Dot:  point,
		}
		d.DrawString(line)
		lineOffset += textHeight + 10
	}
}
