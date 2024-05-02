// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"image"
	"image/color"
	"io/ioutil"
	"os"

	fontPath "github.com/awslabs/diagram-as-code/internal/font"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Resource struct {
	bindings    image.Rectangle
	iconImage   image.Image
	iconBounds  image.Rectangle
	borderColor color.RGBA
	fillColor   color.RGBA
	label       string
	labelFont   string
	labelColor  *color.RGBA
	width       int
	height      int
	margin      Margin
	padding     Padding
	links       []*Link
	drawn       bool
}

func (r Resource) Init() Node {
	rr := Resource{}
	rr.bindings = image.Rect(0, 0, 64, 64)
	rr.iconImage = image.NewRGBA(r.bindings)
	rr.iconBounds = image.Rect(0, 0, 64, 64)
	rr.borderColor = color.RGBA{0, 0, 0, 0}
	rr.fillColor = color.RGBA{0, 0, 0, 0}
	rr.label = ""
	rr.labelColor = &color.RGBA{0, 0, 0, 0}
	rr.width = 64
	rr.height = 64
	rr.margin = Margin{30, 100, 30, 100}
	rr.padding = Padding{0, 0, 0, 0}
	rr.drawn = false
	return &rr
}

func (r *Resource) LoadIcon(imageFilePath string) error {
	imageFile, err := os.Open(imageFilePath)
	if err != nil {
		return err
	}
	iconImage, _, err := image.Decode(imageFile)
	if err != nil {
		return err
	}
	r.iconImage = iconImage
	return nil
}

func (r *Resource) SetIconBounds(bounds image.Rectangle) {
	r.iconBounds = bounds
}

func (r *Resource) SetBindings(bindings image.Rectangle) {
	r.bindings = bindings
}

func (r *Resource) GetBindings() image.Rectangle {
	return r.bindings
}

func (r *Resource) GetMargin() Margin {
	return r.margin
}

func (r *Resource) GetPadding() Padding {
	return r.padding
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
}

func (r *Resource) SetDirection(direction string) {
}

func (r *Resource) AddLink(link *Link) {
	r.links = append(r.links, link)
}

func (r *Resource) AddParent() {
}

func (r *Resource) Scale() {
	b := r.bindings
	r.SetBindings(b)
	return
}

func (r *Resource) Translation(dx, dy int) {
	r.bindings = image.Rect(
		r.bindings.Min.X+dx,
		r.bindings.Min.Y+dy,
		r.bindings.Max.X+dx,
		r.bindings.Max.Y+dy,
	)
}

func (r *Resource) ZeroAdjust() {
	r.Translation(-r.bindings.Min.X+r.padding.Left, -r.bindings.Min.Y+r.padding.Top)
}

func (r *Resource) IsDrawn() bool {
	return r.drawn
}

func (r *Resource) Draw(img *image.RGBA, parent *Group) *image.RGBA {
	if img == nil {
		img = image.NewRGBA(r.bindings)
	}
	rctSrc := r.iconImage.Bounds()
	x := image.Rectangle{r.bindings.Min, r.bindings.Min.Add(image.Point{64, 64})}
	draw.CatmullRom.Scale(img, x, r.iconImage, rctSrc, draw.Over, nil)

	r.drawFrame(img)

	r.drawLabel(img, parent)

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

func (r *Resource) drawLabel(img *image.RGBA, parent *Group) {

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

	face := truetype.NewFace(ft, &opt)

	b, _ := font.BoundString(face, r.label)
	w := b.Max.X - b.Min.X + fixed.I(1)
	h := b.Max.Y - b.Min.Y + fixed.I(1)

	p := r.bindings.Min.Add(image.Point{0, r.iconBounds.Max.Y})

	point := fixed.Point26_6{fixed.I(p.X) - (w-fixed.I(r.bindings.Dx()))/2, fixed.I(p.Y+10) + h}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(r.labelColor),
		Face: face,
		Dot:  point,
	}
	d.DrawString(r.label)
}

func (r *Resource) AddChild(child Node) {
	return
}
