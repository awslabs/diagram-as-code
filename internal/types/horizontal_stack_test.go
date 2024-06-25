package types

import (
	"image"
	"image/color"
	"testing"
)

func TestHorizontalStackInit(t *testing.T) {
	// Create a new HorizontalStack instance
	hs := HorizontalStack{}

	// Initialize the HorizontalStack
	resource := hs.Init()

	// Check the properties of the initialized Group
	if *resource.bindings != image.Rect(0, 0, 320, 190) {
		t.Errorf("Incorrect bindings: %v", resource.bindings)
	}

	if resource.iconBounds != image.Rect(0, 0, 0, 0) {
		t.Errorf("Incorrect iconBounds: %v", resource.iconBounds)
	}

	if *resource.borderColor != (color.RGBA{0, 0, 0, 0}) {
		t.Errorf("Incorrect borderColor: %v", resource.borderColor)
	}

	if resource.fillColor != (color.RGBA{0, 0, 0, 0}) {
		t.Errorf("Incorrect fillColor: %v", resource.fillColor)
	}

	if resource.label != "" {
		t.Errorf("Incorrect label: %s", resource.label)
	}

	if *resource.labelColor != (color.RGBA{0, 0, 0, 0}) {
		t.Errorf("Incorrect labelColor: %v", resource.labelColor)
	}

	if *resource.margin != (Margin{0, 0, 0, 0}) {
		t.Errorf("Incorrect margin: %v", resource.margin)
	}

	if *resource.padding != (Padding{0, 0, 0, 0}) {
		t.Errorf("Incorrect padding: %v", resource.padding)
	}

	if resource.direction != "horizontal" {
		t.Errorf("Incorrect direction: %s", resource.direction)
	}

	if resource.align != "center" {
		t.Errorf("Incorrect align: %s", resource.align)
	}
}
