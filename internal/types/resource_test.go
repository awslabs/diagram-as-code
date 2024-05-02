package types

import (
	"image"
	"image/color"
	"os"
	"testing"
)

func TestResource(t *testing.T) {
	// Test Init
	ri := new(Resource).Init()
	r, ok := ri.(*Resource)
	if !ok {
		t.Errorf("Cannot convert Node to Resource")
	}
	if r.GetBindings() != image.Rect(0, 0, 64, 64) {
		t.Errorf("Init: expected bindings to be (0, 0, 64, 64), got %v", r.GetBindings())
	}
	if r.GetMargin() != (Margin{30, 100, 30, 100}) {
		t.Errorf("Init: expected margin to be (30, 100, 30, 100), got %v", r.GetMargin())
	}
	if r.GetPadding() != (Padding{0, 0, 0, 0}) {
		t.Errorf("Init: expected padding to be (0, 0, 0, 0), got %v", r.GetPadding())
	}
	if r.IsDrawn() {
		t.Error("Init: expected drawn to be false")
	}

	// Test LoadIcon
	iconFile, err := os.Open("testdata/valid_icon.png")
	if err != nil {
		t.Fatalf("LoadIcon: failed to open test icon file: %v", err)
	}
	defer iconFile.Close()
	r.LoadIcon("testdata/icon.png")
	if r.iconImage == nil {
		t.Error("LoadIcon: iconImage is nil")
	}

	// Test SetIconBounds
	bounds := image.Rect(10, 10, 20, 20)
	r.SetIconBounds(bounds)
	if r.iconBounds != bounds {
		t.Errorf("SetIconBounds: expected bounds to be %v, got %v", bounds, r.iconBounds)
	}

	// Test SetBindings
	bindings := image.Rect(100, 100, 200, 200)
	r.SetBindings(bindings)
	if r.GetBindings() != bindings {
		t.Errorf("SetBindings: expected bindings to be %v, got %v", bindings, r.GetBindings())
	}

	// Test SetBorderColor
	borderColor := color.RGBA{255, 0, 0, 255}
	r.SetBorderColor(borderColor)
	if r.borderColor != borderColor {
		t.Errorf("SetBorderColor: expected borderColor to be %v, got %v", borderColor, r.borderColor)
	}

	// Test SetFillColor
	fillColor := color.RGBA{0, 255, 0, 255}
	r.SetFillColor(fillColor)
	if r.fillColor != fillColor {
		t.Errorf("SetFillColor: expected fillColor to be %v, got %v", fillColor, r.fillColor)
	}

	// Test SetLabel
	/*label := "Test Label"
	labelColor := color.RGBA{0, 0, 255, 255}
	labelFont := "testdata/font.ttf"
	r.SetLabel(&label, &labelColor, &labelFont)
	if r.label != label {
		t.Errorf("SetLabel: expected label to be %q, got %q", label, r.label)
	}
	if *r.labelColor != labelColor {
		t.Errorf("SetLabel: expected labelColor to be %v, got %v", labelColor, *r.labelColor)
	}
	if r.labelFont != labelFont {
		t.Errorf("SetLabel: expected labelFont to be %q, got %q", labelFont, r.labelFont)
	}*/

	// Test Translation
	dx, dy := 10, 20
	origBindings := r.GetBindings()
	r.Translation(dx, dy)
	expected := image.Rect(origBindings.Min.X+dx, origBindings.Min.Y+dy, origBindings.Max.X+dx, origBindings.Max.Y+dy)
	if r.GetBindings() != expected {
		t.Errorf("Translation: expected bindings to be %v, got %v", expected, r.GetBindings())
	}

	// Test ZeroAdjust
	r.ZeroAdjust()
	expected = image.Rect(r.padding.Left, r.padding.Top, r.GetBindings().Max.X, r.GetBindings().Max.Y)
	if r.GetBindings() != expected {
		t.Errorf("ZeroAdjust: expected bindings to be %v, got %v", expected, r.GetBindings())
	}

	// Test Draw (basic)
	img := r.Draw(nil, nil)
	if img == nil {
		t.Error("Draw: returned nil image")
	}
	if !r.IsDrawn() {
		t.Error("Draw: drawn should be true after drawing")
	}
}
