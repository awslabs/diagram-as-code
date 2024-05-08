package types

import (
	"image"
	"image/color"
	"testing"
)

func TestVerticalStackInit(t *testing.T) {
	// Create a new VerticalStack instance
	hs := VerticalStack{}

	// Initialize the VerticalStack
	node := hs.Init()

	// Check if the returned value is a Group
	resource, ok := node.(*Resource)
	if !ok {
		t.Errorf("Init() did not return a Group")
	}

	// Check the properties of the initialized Group
	if *resource.bindings != image.Rect(0, 0, 320, 190) {
		t.Errorf("Incorrect bindings: %v", resource.bindings)
	}

	if resource.iconBounds != image.Rect(0, 0, 0, 0) {
		t.Errorf("Incorrect iconBounds: %v", resource.iconBounds)
	}

	if resource.borderColor != (color.RGBA{0, 0, 0, 0}) {
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

	if resource.direction != "vertical" {
		t.Errorf("Incorrect direction: %s", resource.direction)
	}

	if resource.align != "center" {
		t.Errorf("Incorrect align: %s", resource.align)
	}
}
