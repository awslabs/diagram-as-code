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

type Group struct {
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
	direction   string
	align       string
	links       []*Link
	children    []Node
	drawn       bool
}

func (g Group) Init() Node {
	gr := Group{}
	gr.bindings = image.Rect(0, 0, 320, 190)
	gr.iconImage = image.NewRGBA(g.bindings)
	gr.iconBounds = image.Rect(0, 0, 0, 0)
	gr.borderColor = color.RGBA{0, 0, 0, 0}
	gr.fillColor = color.RGBA{0, 0, 0, 0}
	gr.label = ""
	gr.labelFont = ""
	gr.labelColor = nil
	gr.width = 320
	gr.height = 190
	gr.margin = Margin{20, 15, 20, 15}
	gr.padding = Padding{20, 45, 20, 45}
	gr.direction = "horizontal"
	gr.align = "center"
	gr.drawn = false
	return &gr
}

func (g *Group) LoadIcon(imageFilePath string) error {
	imageFile, err := os.Open(imageFilePath)
	if err != nil {
		return err
	}
	defer imageFile.Close()
	iconImage, _, err := image.Decode(imageFile)
	if err != nil {
		return err
	}
	g.iconBounds = image.Rect(0, 0, 64, 64)
	g.iconImage = iconImage
	return nil
}

func (g *Group) SetIconBounds(bounds image.Rectangle) {
	g.iconBounds = bounds
}

func (g *Group) SetBindings(bindings image.Rectangle) {
	g.bindings = bindings
}

func (g Group) GetBindings() image.Rectangle {
	return g.bindings
}

func (g Group) GetMargin() Margin {
	return g.margin
}

func (g Group) GetPadding() Padding {
	return g.padding
}

func (g *Group) SetBorderColor(borderColor color.RGBA) {
	g.borderColor = borderColor
}

func (g *Group) SetFillColor(fillColor color.RGBA) {
	g.fillColor = fillColor
}

func (g *Group) SetLabel(label *string, labelColor *color.RGBA, labelFont *string) {

	if label != nil {
		g.label = *label
	}
	if labelColor != nil {
		g.labelColor = labelColor
	}
	if labelFont != nil {
		g.labelFont = *labelFont
	}
}

func (g *Group) SetAlign(align string) {
	g.align = align
}

func (g *Group) SetDirection(direction string) {
	g.direction = direction
}

func (g *Group) AddLink(link *Link) {
	g.links = append(g.links, link)
}

func (g *Group) AddParent() {
}

func (g *Group) AddChild(child Node) {
	g.children = append(g.children, child)
}

func (g *Group) Scale() {
	log.Infof("Scale %s", g.label)
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
	w := g.padding.Left + g.padding.Right
	h := g.padding.Top + g.padding.Bottom
	for _, subGroup := range g.children {
		subGroup.Scale()
		bindings := subGroup.GetBindings()
		margin := subGroup.GetMargin()
		if g.direction == "horizontal" {
			w += bindings.Dx() + margin.Left + margin.Right
			h = maxInt(h, bindings.Dy()+margin.Top+margin.Bottom)
		} else {
			w = maxInt(w, bindings.Dx()+margin.Left+margin.Right)
			h += bindings.Dy() + margin.Top + margin.Bottom
		}
		if prev != nil {
			prevBindings := prev.GetBindings()
			prevMargin := prev.GetMargin()
			if g.direction == "horizontal" {
				switch g.align {
				case "top":
					subGroup.Translation(
						prevBindings.Max.X+prevMargin.Right+margin.Left-bindings.Min.X,
						prevBindings.Min.Y-prevMargin.Top+margin.Top-bindings.Min.Y,
					)
				case "center":
					subGroup.Translation(
						prevBindings.Max.X+prevMargin.Right+margin.Left-bindings.Min.X,
						prevBindings.Min.Y+(prevBindings.Dy()-bindings.Dy())/2-bindings.Min.Y,
					)
				case "bottom":
					subGroup.Translation(
						prevBindings.Max.X+prevMargin.Right+margin.Left-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom-margin.Bottom-bindings.Max.Y,
					)
				default:
					log.Fatalf("Unknown align %s in the direction(%s) on %s", g.align, g.direction, g.label)
				}
			} else {
				switch g.align {
				case "left":
					subGroup.Translation(
						prevBindings.Min.X-prevMargin.Left+margin.Left-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom+margin.Top-bindings.Min.Y,
					)
				case "center":
					subGroup.Translation(
						prevBindings.Min.X+(prevBindings.Dx()-bindings.Dx())/2-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom+margin.Top-bindings.Min.Y,
					)
				case "right":
					subGroup.Translation(
						prevBindings.Max.X+prevMargin.Right-margin.Right-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom+margin.Top-bindings.Min.Y,
					)
				default:
					log.Fatalf("Unknown align %s in the direction(%s) on %s", g.align, g.direction, g.label)
				}
			}
		}
		bindings = subGroup.GetBindings()
		b.Min.X = minInt(b.Min.X, bindings.Min.X-margin.Left-g.padding.Left)
		b.Min.Y = minInt(b.Min.Y, bindings.Min.Y-margin.Top-g.iconBounds.Dy()-g.padding.Top)
		b.Max.X = maxInt(b.Max.X, bindings.Max.X+margin.Right+g.padding.Right)
		b.Max.Y = maxInt(b.Max.Y, bindings.Max.Y+margin.Bottom+g.padding.Bottom)
		prev = subGroup
	}
	b.Max.X = maxInt(b.Max.X, b.Min.X+w)
	b.Max.Y = maxInt(b.Max.Y, b.Min.Y+h)
	if b.Min.X != math.MaxInt {
		g.SetBindings(b)
	}
}

func (g *Group) Translation(dx, dy int) {
	g.bindings = image.Rect(
		g.bindings.Min.X+dx,
		g.bindings.Min.Y+dy,
		g.bindings.Max.X+dx,
		g.bindings.Max.Y+dy,
	)
	for _, subGroup := range g.children {
		subGroup.Translation(dx, dy)
	}
}

func (g *Group) ZeroAdjust() {
	g.Translation(-g.bindings.Min.X+g.padding.Left, -g.bindings.Min.Y+g.padding.Top)
}

func (g *Group) IsDrawn() bool {
	return g.drawn
}

func (g *Group) Draw(img *image.RGBA, parent *Group) *image.RGBA {
	if img == nil {
		img = image.NewRGBA(g.bindings)
	}
	g.drawFrame(img)

	x := image.Rectangle{g.bindings.Min, g.bindings.Min.Add(image.Point{64, 64})}
	rctSrc := g.iconImage.Bounds()
	draw.CatmullRom.Scale(img, x, g.iconImage, rctSrc, draw.Over, nil)

	g.drawLabel(img, parent)

	for _, subGroup := range g.children {
		subGroup.Draw(img, g)
	}
	g.drawn = true
	for _, v := range g.links {
		target := *v.Target
		source := *v.Source
		if target.IsDrawn() && source.IsDrawn() {
			v.Draw(img)
		}
	}
	return img
}

func (g *Group) drawFrame(img *image.RGBA) {
	x1 := g.bindings.Min.X
	x2 := g.bindings.Max.X
	y1 := g.bindings.Min.Y
	y2 := g.bindings.Max.Y
	for x := x1 - WIDTH + 1; x < x2+WIDTH-1; x++ {
		for y := y1 - WIDTH + 1; y < y2+WIDTH-1; y++ {
			c := img.At(x, y)
			if x <= x1 || x >= x2-1 || y <= y1 || y >= y2-1 {
				// Set border
				img.Set(x, y, _blend_color(c, g.borderColor))
				//img.Set(x, y, border_color)
			} else {
				// Set background
				img.Set(x, y, _blend_color(c, g.fillColor))
				//img.Set(x, y, fill_color)
			}
		}
	}
}

func (g *Group) drawLabel(img *image.RGBA, parent *Group) {

	p := g.bindings.Min.Add(g.iconBounds.Max)

	if g.labelFont == "" {
		if parent != nil && parent.labelFont != "" {
			g.labelFont = parent.labelFont
		} else {
			g.labelFont = fontPath.Arial
		}
	}
	if g.labelColor == nil {
		if parent != nil && parent.labelColor != nil {
			g.labelColor = parent.labelColor
		} else {
			g.labelColor = &color.RGBA{0, 0, 0, 255}
		}
	}
	f, err := os.Open(g.labelFont)
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
		Size:              30,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	face := truetype.NewFace(ft, &opt)

	b, _ := font.BoundString(face, g.label)
	//w := b.Max.X - b.Min.X + fixed.I(1)
	h := b.Max.Y - b.Min.Y + fixed.I(1)

	point := fixed.Point26_6{fixed.I(p.X) + 1000, fixed.I(p.Y) + (64-h)/2}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(g.labelColor),
		Face: face,
		Dot:  point,
	}
	d.DrawString(g.label)
}
