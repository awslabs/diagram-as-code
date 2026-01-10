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
	"sort"
	"strconv"
	"strings"

	fontPath "github.com/awslabs/diagram-as-code/internal/font"
	"github.com/awslabs/diagram-as-code/internal/vector"
	"github.com/golang/freetype/truetype"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
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
	groupingOffset bool // Flag: if true, enable grouping offset for links
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
	defer func() {
		if closeErr := imageFile.Close(); closeErr != nil {
			log.Warnf("Failed to close image file: %v", closeErr)
		}
	}()
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
	if r.bindings == nil {
		return image.Rectangle{}
	}
	return *r.bindings
}

func (r *Resource) GetMargin() Margin {
	return *r.margin
}

func (r *Resource) GetPadding() Padding {
	return *r.padding
}

func (r *Resource) SetMargin(margin Margin) {
	r.margin = &margin
}

func (r *Resource) SetPadding(padding Padding) {
	r.padding = &padding
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

func (r *Resource) SetGroupingOffset(enable bool) {
	r.groupingOffset = enable
}

func (r *Resource) AddLink(link *Link) {
	r.links = append(r.links, link)
}

func (r *Resource) GetLinks() []*Link {
	return r.links
}

func (r *Resource) AddParent() {
}

func (r *Resource) AddChild(child *Resource) error {
	// [TODO] check whether the parent is border children
	if child == nil {
		return fmt.Errorf("unknown child resource - please see debug logs with -v flag")
	}
	r.children = append(r.children, child)
	return nil
}

func (r *Resource) AddBorderChild(borderChild *BorderChild) error {
	hasChild := len(borderChild.Resource.children) != 0
	if hasChild {
		return fmt.Errorf("couldn't add group to border children")
	}
	r.borderChildren = append(r.borderChildren, borderChild)
	return nil
}

func (r *Resource) prepareFontFace(hasChild bool, parent *Resource) (font.Face, error) {
	if r.labelFont == "" {
		if parent != nil {
			log.Infof("parent labelFont: %s", r.labelFont)
		}
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
	log.Infof("labelFont: %s", r.labelFont)
	if r.labelColor == nil {
		if parent != nil && parent.labelColor != nil {
			r.labelColor = parent.labelColor
		} else {
			r.labelColor = &color.RGBA{0, 0, 0, 255}
		}
	}
	var ttfBytes []byte
	if r.labelFont == "goregular" || r.labelFont == "" {
		// Use Go-fonts instead system fonts
		ttfBytes = goregular.TTF
	} else {
		f, err := os.Open(r.labelFont)
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
	if hasChild {
		opt.Size = 30
	}

	return truetype.NewFace(ft, &opt), nil
}

func (r *Resource) Scale(parent *Resource, visited map[*Resource]bool) error {
	log.Infof("Scale %s", r.label)

	if visited == nil {
		visited = make(map[*Resource]bool)
	}
	// Check if resource has been visited to detect cycles
	// Using comma-ok idiom for safe map access
	if isVisited, _ := visited[r]; isVisited {
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
	fontFace, err := r.prepareFontFace(hasChildren, parent)
	if err != nil {
		return fmt.Errorf("failed to prepare font face: %w", err)
	}
	if r.label != "" {
		textHeight = 10
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
		err := subResource.Scale(r, visited)
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
					if err := subResource.Translation(
						prevBindings.Max.X+prevMargin.Right+margin.Left-bindings.Min.X,
						prevBindings.Min.Y-prevMargin.Top+margin.Top-bindings.Min.Y,
					); err != nil {
						return fmt.Errorf("failed to translate subresource: %w", err)
					}
				case "center":
					if err := subResource.Translation(
						prevBindings.Max.X+prevMargin.Right+margin.Left-bindings.Min.X,
						prevBindings.Min.Y+(prevBindings.Dy()-bindings.Dy())/2-bindings.Min.Y,
					); err != nil {
						return fmt.Errorf("failed to translate subresource: %w", err)
					}
				case "bottom":
					if err := subResource.Translation(
						prevBindings.Max.X+prevMargin.Right+margin.Left-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom-margin.Bottom-bindings.Max.Y,
					); err != nil {
						return fmt.Errorf("failed to translate subresource: %w", err)
					}
				default:
					return fmt.Errorf("unknown align %s in the direction(%s) on %s", r.align, r.direction, r.label)
				}
			} else {
				switch r.align {
				case "left":
					if err := subResource.Translation(
						prevBindings.Min.X-prevMargin.Left+margin.Left-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom+margin.Top-bindings.Min.Y,
					); err != nil {
						return fmt.Errorf("failed to translate subresource: %w", err)
					}
				case "center":
					if err := subResource.Translation(
						prevBindings.Min.X+(prevBindings.Dx()-bindings.Dx())/2-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom+margin.Top-bindings.Min.Y,
					); err != nil {
						return fmt.Errorf("failed to translate subresource: %w", err)
					}
				case "right":
					if err := subResource.Translation(
						prevBindings.Max.X+prevMargin.Right-margin.Right-bindings.Min.X,
						prevBindings.Max.Y+prevMargin.Bottom+margin.Top-bindings.Min.Y,
					); err != nil {
						return fmt.Errorf("failed to translate subresource: %w", err)
					}
				default:
					return fmt.Errorf("unknown align %s in the direction(%s) on %s", r.align, r.direction, r.label)
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
			return fmt.Errorf("failed to calculate position for border child: %w", err)
		}
		err = borderChild.Resource.Scale(r, visited) // to initialize default values
		if err != nil {
			return err
		}
		bindings := borderChild.Resource.GetBindings()
		if err := borderChild.Resource.Translation(
			pt.X-(bindings.Min.X+bindings.Max.X)/2,
			pt.Y-(bindings.Min.Y+bindings.Max.Y)/2,
		); err != nil {
			return fmt.Errorf("failed to translate border child resource: %w", err)
		}
	}
	return nil
}

func (r *Resource) Translation(dx, dy int) error {
	if r.bindings == nil {
		return fmt.Errorf("the resource has no binding")
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
		if err := subResource.Translation(dx, dy); err != nil {
			return fmt.Errorf("failed to translate child resource: %w", err)
		}
	}
	for _, borderChild := range r.borderChildren {
		if err := borderChild.Resource.Translation(dx, dy); err != nil {
			return fmt.Errorf("failed to translate border child resource: %w", err)
		}
	}
	return nil
}

func (r *Resource) ZeroAdjust() error {
	return r.Translation(-r.bindings.Min.X+r.padding.Left, -r.bindings.Min.Y+r.padding.Top)
}

func (r *Resource) IsDrawn() bool {
	return r.drawn
}

func (r *Resource) Draw(img *image.RGBA, parent *Resource) (*image.RGBA, error) {
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
		if err := r.drawLabel(img, parent, len(r.children) > 0, hasIcon); err != nil {
			return nil, fmt.Errorf("failed to draw label: %w", err)
		}
	}

	for _, subResource := range r.children {
		if _, err := subResource.Draw(img, r); err != nil {
			return nil, fmt.Errorf("failed to draw child resource: %w", err)
		}
	}
	for _, borderResource := range r.borderChildren {
		if _, err := borderResource.Resource.Draw(img, r); err != nil {
			return nil, fmt.Errorf("failed to draw border child resource: %w", err)
		}
	}
	r.drawn = true

	// Pre-sort links before drawing
	r.sortAllLinks()

	for _, v := range r.links {
		source := *v.Source
		target := *v.Target
		if source.IsDrawn() && target.IsDrawn() {
			if err := v.Draw(img); err != nil {
				return nil, fmt.Errorf("failed to draw link: %w", err)
			}
		}
	}
	return img, nil
}

func (r *Resource) sortAllLinks() {
	log.Infof("=== Sorting links for resource %p ===", r)

	// Group links by same position
	linkGroups := make(map[string][]*Link)

	for _, link := range r.links {
		var key string
		if link.Source == r {
			key = fmt.Sprintf("%d", link.SourcePosition)
		} else if link.Target == r {
			key = fmt.Sprintf("%d", link.TargetPosition)
		}
		if key != "" {
			links, ok := linkGroups[key]
			if !ok {
				links = make([]*Link, 0)
			}
			linkGroups[key] = append(links, link)
		}
	}

	log.Infof("Found %d link groups", len(linkGroups))

	// Sort each group and update original array
	for key, links := range linkGroups {
		if len(links) <= 1 {
			log.Infof("Group %s: only %d link, no sorting needed", key, len(links))
			continue
		}

		log.Infof("Group %s: sorting %d links", key, len(links))

		// Convert key to position
		position, _ := strconv.Atoi(key)

		// Log order before sorting
		for i, link := range links {
			var otherResource *Resource
			var otherPos Windrose
			if link.Source == r {
				otherResource = link.Target
				otherPos = link.TargetPosition
			} else {
				otherResource = link.Source
				otherPos = link.SourcePosition
			}
			pt, _ := calcPosition(otherResource.GetBindings(), otherPos)
			log.Infof("  Before sort [%d]: %s->%s, other pos: (%d, %d)",
				i, getResourceName(link.Source), getResourceName(link.Target), pt.X, pt.Y)
		}

		sort.Slice(links, func(i, j int) bool {
			var pt1, pt2 image.Point
			if links[i].Source == r {
				pt1, _ = calcPosition(links[i].Target.GetBindings(), links[i].TargetPosition)
			} else {
				pt1, _ = calcPosition(links[i].Source.GetBindings(), links[i].SourcePosition)
			}
			if links[j].Source == r {
				pt2, _ = calcPosition(links[j].Target.GetBindings(), links[j].TargetPosition)
			} else {
				pt2, _ = calcPosition(links[j].Source.GetBindings(), links[j].SourcePosition)
			}

			// Sort by perpendicular direction of direction vector
			direction := getDirectionVectorStatic(int(position))
			perpendicular := direction.Perpendicular()

			proj1 := float64(pt1.X)*perpendicular.X + float64(pt1.Y)*perpendicular.Y
			proj2 := float64(pt2.X)*perpendicular.X + float64(pt2.Y)*perpendicular.Y

			log.Infof("    Compare: pt1=(%d,%d) proj1=%.1f vs pt2=(%d,%d) proj2=%.1f -> %v",
				pt1.X, pt1.Y, proj1, pt2.X, pt2.Y, proj2, proj1 < proj2)

			return proj1 < proj2
		})

		// Log order after sorting
		for i, link := range links {
			var otherResource *Resource
			var otherPos Windrose
			if link.Source == r {
				otherResource = link.Target
				otherPos = link.TargetPosition
			} else {
				otherResource = link.Source
				otherPos = link.SourcePosition
			}
			pt, _ := calcPosition(otherResource.GetBindings(), otherPos)
			log.Infof("  After sort [%d]: %s->%s, other pos: (%d, %d)",
				i, getResourceName(link.Source), getResourceName(link.Target), pt.X, pt.Y)
		}

		// Apply sort results to original array
		r.updateLinksOrder(links, key)
	}
	log.Infof("=== End sorting ===")
}

func (r *Resource) updateLinksOrder(sortedLinks []*Link, groupKey string) {
	// Replace links in original array with new order for sort targets
	newLinks := make([]*Link, 0, len(r.links))

	// Add non-sort target links first
	for _, link := range r.links {
		var key string
		if link.Source == r {
			key = fmt.Sprintf("%d", link.SourcePosition)
		} else if link.Target == r {
			key = fmt.Sprintf("%d", link.TargetPosition)
		}

		if key != groupKey {
			newLinks = append(newLinks, link)
		}
	}

	// Add sorted links
	newLinks = append(newLinks, sortedLinks...)

	r.links = newLinks
	log.Infof("Updated links order for group %s", groupKey)
}

func getResourceName(r *Resource) string {
	if r.label != "" {
		return r.label
	}
	return fmt.Sprintf("Resource_%p", r)
}

func getDirectionVectorStatic(position int) vector.Vector {
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
		return vector.New(0, -1)
	}
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

func (r *Resource) drawLabel(img *image.RGBA, parent *Resource, hasChild, hasIcon bool) error {
	face, err := r.prepareFontFace(hasChild, parent)
	if err != nil {
		return fmt.Errorf("failed to prepare font face for drawing label: %w", err)
	}

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
	return nil
}
