// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"

	fontPath "github.com/awslabs/diagram-as-code/internal/font"
	"github.com/golang/freetype/truetype"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Resource struct {
	bindings    *image.Rectangle
	iconImage   image.Image
	iconBounds  image.Rectangle
	borderColor color.RGBA
	fillColor   color.RGBA
	label       string
	labelFont   string
	labelColor  *color.RGBA
	margin      *Margin
	padding     *Padding
	direction   string
	align       string
	links       []*Link
	children    []Node
	drawn       bool
}

var defaultResourceValues = map[bool]Resource{
	false: { // resource has not children and show as Resource
		bindings: &image.Rectangle{
			image.Point{0, 0},
			image.Point{64, 64},
		},
		margin:  &Margin{30, 100, 30, 100},
		padding: &Padding{0, 0, 0, 0},
	},
	true: { // resource has children and show as Group
		bindings: &image.Rectangle{
			image.Point{0, 0},
			image.Point{320, 190},
		},
		margin:  &Margin{20, 15, 20, 15},
		padding: &Padding{20, 45, 20, 45},
	},
}

func (r Resource) Init() Node {
	rr := Resource{}
	rr.bindings = nil
	rr.iconImage = image.NewRGBA(image.Rect(0, 0, 0, 0))
	rr.iconBounds = image.Rect(0, 0, 0, 0)
	rr.borderColor = color.RGBA{0, 0, 0, 0}
	rr.fillColor = color.RGBA{0, 0, 0, 0}
	rr.label = ""
	rr.labelFont = ""
	rr.labelColor = &color.RGBA{0, 0, 0, 0}
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
	r.iconImage = iconImage
	return nil
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
	r.borderColor = borderColor
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

func (r *Resource) AddLink(link *Link) {
	r.links = append(r.links, link)
}

func (r *Resource) AddParent() {
}

func (r *Resource) AddChild(child Node) {
	r.children = append(r.children, child)
}

func (r *Resource) Scale() {
	log.Infof("Scale %s", r.label)
	var prev Node
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
	if r.bindings == nil {
		r.bindings = defaultResourceValues[hasChildren].bindings
	}
	if r.margin == nil {
		r.margin = defaultResourceValues[hasChildren].margin
	}
	if r.padding == nil {
		r.padding = defaultResourceValues[hasChildren].padding
	}

	w := r.padding.Left + r.padding.Right
	h := r.padding.Top + r.padding.Bottom
	for _, subResource := range r.children {
		subResource.Scale()
		bindings := subResource.GetBindings()
		margin := subResource.GetMargin()
		if r.direction == "horizontal" {
			w += bindings.Dx() + margin.Left + margin.Right
			h = maxInt(h, bindings.Dy()+margin.Top+margin.Bottom)
		} else {
			w = maxInt(w, bindings.Dx()+margin.Left+margin.Right)
			h += bindings.Dy() + margin.Top + margin.Bottom
		}
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
		b.Min.Y = minInt(b.Min.Y, bindings.Min.Y-margin.Top-r.iconBounds.Dy()-r.padding.Top)
		b.Max.X = maxInt(b.Max.X, bindings.Max.X+margin.Right+r.padding.Right)
		b.Max.Y = maxInt(b.Max.Y, bindings.Max.Y+margin.Bottom+r.padding.Bottom)
		prev = subResource
	}
	b.Max.X = maxInt(b.Max.X, b.Min.X+w)
	b.Max.Y = maxInt(b.Max.Y, b.Min.Y+h)
	if b.Min.X != math.MaxInt {
		r.SetBindings(b)
	}
}

func (r *Resource) Translation(dx, dy int) {
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
}

func (r *Resource) ZeroAdjust() {
	r.Translation(-r.bindings.Min.X+r.padding.Left, -r.bindings.Min.Y+r.padding.Top)
}

func (r *Resource) IsDrawn() bool {
	return r.drawn
}

func (r *Resource) Draw(img *image.RGBA, parent Node) *image.RGBA {
	if img == nil {
		img = image.NewRGBA(*r.bindings)
	}

	r.drawFrame(img)

	rctSrc := r.iconImage.Bounds()
	x := image.Rectangle{r.bindings.Min, r.bindings.Min.Add(image.Point{64, 64})}
	draw.CatmullRom.Scale(img, x, r.iconImage, rctSrc, draw.Over, nil)

	if parent != nil {
		r.drawLabel(img, parent.(*Resource), len(r.children) > 0)
	}

	for _, subResource := range r.children {
		subResource.Draw(img, r)
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
				img.Set(x, y, _blend_color(c, r.borderColor))
				//img.Set(x, y, border_color)
			} else {
				// Set background
				img.Set(x, y, _blend_color(c, r.fillColor))
				//img.Set(x, y, fill_color)
			}
		}
	}
}

func (r *Resource) drawLabel(img *image.RGBA, parent *Resource, hasChild bool) {

	if r.labelFont == "" {
		if parent != nil && parent.labelFont != "" {
			r.labelFont = parent.labelFont
		} else {
			r.labelFont = fontPath.Arial
		}
	}
	if r.labelColor == nil {
		if parent != nil && parent.labelColor != nil {
			r.labelColor = parent.labelColor
		} else {
			r.labelColor = &color.RGBA{0, 0, 0, 255}
		}
	}
	f, err := os.Open(r.labelFont)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ttfBytes, err := ioutil.ReadAll(f)
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

	face := truetype.NewFace(ft, &opt)

	b, _ := font.BoundString(face, r.label)
	w := b.Max.X - b.Min.X + fixed.I(1)
	h := b.Max.Y - b.Min.Y + fixed.I(1)

	p := r.bindings.Min.Add(image.Point{0, r.iconBounds.Max.Y})

	point := fixed.Point26_6{fixed.I(p.X) - (w-fixed.I(r.bindings.Dx()))/2, fixed.I(p.Y+10) + h}
	if hasChild {
		p = r.bindings.Min.Add(r.iconBounds.Max)
		point = fixed.Point26_6{fixed.I(p.X) + 1000, fixed.I(p.Y) + (64-h)/2}
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(r.labelColor),
		Face: face,
		Dot:  point,
	}
	d.DrawString(r.label)
}
