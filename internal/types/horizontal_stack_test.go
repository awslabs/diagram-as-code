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
	node := hs.Init()

	// Check if the returned value is a Group
	group, ok := node.(*Group)
	if !ok {
		t.Errorf("Init() did not return a Group")
	}

	// Check the properties of the initialized Group
	if group.bindings != image.Rect(0, 0, 320, 190) {
		t.Errorf("Incorrect bindings: %v", group.bindings)
	}

	if group.iconBounds != image.Rect(0, 0, 0, 0) {
		t.Errorf("Incorrect iconBounds: %v", group.iconBounds)
	}

	if group.borderColor != (color.RGBA{0, 0, 0, 0}) {
		t.Errorf("Incorrect borderColor: %v", group.borderColor)
	}

	if group.fillColor != (color.RGBA{0, 0, 0, 0}) {
		t.Errorf("Incorrect fillColor: %v", group.fillColor)
	}

	if group.label != "" {
		t.Errorf("Incorrect label: %s", group.label)
	}

	if *group.labelColor != (color.RGBA{0, 0, 0, 0}) {
		t.Errorf("Incorrect labelColor: %v", group.labelColor)
	}

	if group.width != 320 {
		t.Errorf("Incorrect width: %d", group.width)
	}

	if group.height != 190 {
		t.Errorf("Incorrect height: %d", group.height)
	}

	if group.margin != (Margin{0, 0, 0, 0}) {
		t.Errorf("Incorrect margin: %v", group.margin)
	}

	if group.padding != (Padding{0, 0, 0, 0}) {
		t.Errorf("Incorrect padding: %v", group.padding)
	}

	if group.direction != "horizontal" {
		t.Errorf("Incorrect direction: %s", group.direction)
	}

	if group.align != "center" {
		t.Errorf("Incorrect align: %s", group.align)
	}
}
