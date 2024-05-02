package types

import (
	"image"
	"image/color"
	"testing"

	_ "image/png"
)

func TestGroupInit(t *testing.T) {
	g := Group{}
	node := g.Init()
	group, ok := node.(*Group)
	if !ok {
		t.Errorf("Init() did not return a *Group")
	}

	expectedBindings := image.Rect(0, 0, 320, 190)
	if group.bindings != expectedBindings {
		t.Errorf("Unexpected bindings: got %v, want %v", group.bindings, expectedBindings)
	}

	expectedMargin := Margin{20, 15, 20, 15}
	if group.margin != expectedMargin {
		t.Errorf("Unexpected margin: got %v, want %v", group.margin, expectedMargin)
	}

	expectedPadding := Padding{20, 45, 20, 45}
	if group.padding != expectedPadding {
		t.Errorf("Unexpected padding: got %v, want %v", group.padding, expectedPadding)
	}

	expectedDirection := "horizontal"
	if group.direction != expectedDirection {
		t.Errorf("Unexpected direction: got %s, want %s", group.direction, expectedDirection)
	}

	expectedAlign := "center"
	if group.align != expectedAlign {
		t.Errorf("Unexpected align: got %s, want %s", group.align, expectedAlign)
	}
}

func TestGroupLoadIcon(t *testing.T) {
	g := Group{}
	g.Init()

	// Test with a valid image file
	validFilePath := "testdata/valid_icon.png"
	err := g.LoadIcon(validFilePath)
	if err != nil {
		t.Errorf("Failed to load valid icon image:%v", err)
	}
	if g.iconImage == nil {
		t.Error("Failed to load valid icon image")
	}
	expectedBounds := image.Rect(0, 0, 64, 64)
	if g.iconBounds != expectedBounds {
		t.Errorf("Unexpected icon bounds: got %v, want %v", g.iconBounds, expectedBounds)
	}

	// Test with an invalid file path
	invalidFilePath := "testdata/invalid_path.png"
	err = g.LoadIcon(invalidFilePath)
	if err == nil {
		t.Errorf("LoadIcon with invalid path unexpected successfully")
	}

}

func TestGroupSetters(t *testing.T) {
	g := Group{}
	g.Init()

	// Test SetIconBounds
	expectedBounds := image.Rect(10, 10, 100, 100)
	g.SetIconBounds(expectedBounds)
	if g.iconBounds != expectedBounds {
		t.Errorf("Unexpected icon bounds: got %v, want %v", g.iconBounds, expectedBounds)
	}

	// Test SetBindings
	expectedBindings := image.Rect(50, 50, 200, 300)
	g.SetBindings(expectedBindings)
	if g.bindings != expectedBindings {
		t.Errorf("Unexpected bindings: got %v, want %v", g.bindings, expectedBindings)
	}

	// Test SetBorderColor
	expectedBorderColor := color.RGBA{255, 0, 0, 255}
	g.SetBorderColor(expectedBorderColor)
	if g.borderColor != expectedBorderColor {
		t.Errorf("Unexpected border color: got %v, want %v", g.borderColor, expectedBorderColor)
	}

	// Test SetFillColor
	expectedFillColor := color.RGBA{0, 255, 0, 255}
	g.SetFillColor(expectedFillColor)
	if g.fillColor != expectedFillColor {
		t.Errorf("Unexpected fill color: got %v, want %v", g.fillColor, expectedFillColor)
	}

	// Test SetLabel
	expectedLabel := "Test Label"
	expectedLabelColor := color.RGBA{0, 0, 255, 255}
	expectedLabelFont := "testdata/font.ttf"
	g.SetLabel(&expectedLabel, &expectedLabelColor, &expectedLabelFont)
	if g.label != expectedLabel {
		t.Errorf("Unexpected label: got %s, want %s", g.label, expectedLabel)
	}
	if g.labelColor == nil || *g.labelColor != expectedLabelColor {
		t.Errorf("Unexpected label color: got %v, want %v", g.labelColor, &expectedLabelColor)
	}
	if g.labelFont != expectedLabelFont {
		t.Errorf("Unexpected label font: got %s, want %s", g.labelFont, expectedLabelFont)
	}

	// Test SetAlign
	expectedAlign := "left"
	g.SetAlign(expectedAlign)
	if g.align != expectedAlign {
		t.Errorf("Unexpected align: got %s, want %s", g.align, expectedAlign)
	}

	// Test SetDirection
	expectedDirection := "vertical"
	g.SetDirection(expectedDirection)
	if g.direction != expectedDirection {
		t.Errorf("Unexpected direction: got %s, want %s", g.direction, expectedDirection)
	}
}

func TestGroupAddLink(t *testing.T) {
	// Create a sample Link
	source := Group{}
	src := source.Init()
	target := Group{}
	tg := target.Init()
	link := &Link{Source: &src, Target: &tg}

	// Add the Link
	source.AddLink(link)

	// Check if the Link is added
	if len(source.links) != 1 || source.links[0] != link {
		t.Errorf("Failed to add Link correctly")
	}
}

func TestGroupAddChild(t *testing.T) {
	g := Group{}
	g.Init()

	// Create a sample child Node
	child := &Group{}
	child.Init()

	// Add the child
	g.AddChild(child)

	// Check if the child is added
	if len(g.children) != 1 || g.children[0] != child {
		t.Errorf("Failed to add child correctly")
	}
}

func TestGroupScale(t *testing.T) {
	// Set up a test group with children
	g := new(Group).Init()
	g.LoadIcon("testdata/valid_icon.png")

	child1 := new(Resource).Init()
	child1.LoadIcon("testdata/valid_icon.png")

	child2 := new(Resource).Init()
	child2.LoadIcon("testdata/valid_icon.png")

	g.AddChild(child1)
	g.AddChild(child2)

	// Scale the group
	g.Scale()
	g.ZeroAdjust()

	// Check the scaled bindings
	expectedBindings := image.Rect(45, 20, 663, 248)
	if g.GetBindings() != expectedBindings {
		t.Errorf("Unexpected scaled bindings: got %v, want %v", g.GetBindings(), expectedBindings)
	}
}
