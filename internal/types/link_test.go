package types

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestLinkInit(t *testing.T) {
	source := new(Resource).Init()
	target := new(Resource).Init()
	sourceArrowHead := ArrowHead{Type: "Default", Length: 10, Width: "Default"}
	targetArrowHead := ArrowHead{Type: "Open", Length: 15, Width: "Wide"}

	link := Link{}.Init(source, "Top", sourceArrowHead, target, "Bottom", targetArrowHead, 2)

	if link.Source != source {
		t.Errorf("Expected source node to be %v, got %v", source, link.Source)
	}
	if link.SourcePosition != "Top" {
		t.Errorf("Expected source position to be 'Top', got %s", link.SourcePosition)
	}
	if link.SourceArrowHead != sourceArrowHead {
		t.Errorf("Expected source arrow head to be %v, got %v", sourceArrowHead, link.SourceArrowHead)
	}
	if link.Target != target {
		t.Errorf("Expected target node to be %v, got %v", target, link.Target)
	}
	if link.TargetPosition != "Bottom" {
		t.Errorf("Expected target position to be 'Bottom', got %s", link.TargetPosition)
	}
	if link.TargetArrowHead != targetArrowHead {
		t.Errorf("Expected target arrow head to be %v, got %v", targetArrowHead, link.TargetArrowHead)
	}
	if link.LineWidth != 2 {
		t.Errorf("Expected line width to be 2, got %d", link.LineWidth)
	}
	if link.drawn {
		t.Error("Expected link to be not drawn initially")
	}
	if link.lineColor != (color.RGBA{0, 0, 0, 255}) {
		t.Errorf("Expected line color to be black, got %v", link.lineColor)
	}
}

func TestGetThreeSide(t *testing.T) {
	link := &Link{}

	a, b, c := link.getThreeSide("Narrow")
	if a != math.Sqrt(3.0) || b != 2.0 || c != 1.0 {
		t.Errorf("Expected (sqrt(3), 2, 1) for 'Narrow', got (%f, %f, %f)", a, b, c)
	}

	a, b, c = link.getThreeSide("Default")
	if a != 1.0 || b != math.Sqrt(2.0) || c != 1.0 {
		t.Errorf("Expected (1, sqrt(2), 1) for 'Default', got (%f, %f, %f)", a, b, c)
	}

	a, b, c = link.getThreeSide("Wide")
	if a != 1.0 || b != 2.0 || c != math.Sqrt(3.0) {
		t.Errorf("Expected (1, 2, sqrt(3)) for 'Wide', got (%f, %f, %f)", a, b, c)
	}

	a, b, c = link.getThreeSide("Invalid")
	if a != 0 || b != 0 || c != 0 {
		t.Errorf("Expected (0, 0, 0) for invalid input, got (%f, %f, %f)", a, b, c)
	}
}

func TestDrawNeighborsDots1(t *testing.T) {
	// Create a test image
	width, height := 10, 10
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	// Set up a test link
	link := &Link{lineColor: color.RGBA{255, 0, 0, 255}}

	// Draw a dot at the center
	x, y := float64(width/2), float64(height/2)
	link.drawNeighborsDot(img, x, y)

	// Check if the center pixel and its neighbors are set correctly
	centerColor := img.At(width/2, height/2).(color.RGBA)
	if centerColor != (color.RGBA{255, 0, 0, 255}) {
		t.Errorf("Expected center pixel to be red, got %v", centerColor)
	}
}

func TestDrawNeighborDots2(t *testing.T) {
	// Create a test image
	width, height := 9, 9
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	// Set up a test link
	link := &Link{lineColor: color.RGBA{255, 0, 0, 255}}

	// Draw a dot at the center
	x, y := float64(width)/2.0, float64(height)/2.0
	link.drawNeighborsDot(img, x, y)

	// Check neighbors
	neighbors := []image.Point{
		{width / 2, height / 2},
		{width/2 + 1, height / 2},
		{width / 2, height/2 + 1},
		{width/2 + 1, height/2 + 1},
	}
	for _, neighbor := range neighbors {
		neighborColor := img.At(neighbor.X, neighbor.Y).(color.RGBA)
		if neighborColor != (color.RGBA{255, 192, 192, 255}) {
			t.Errorf("Expected neighbor pixel at (%d, %d) to be semi-transparent red, got %v", neighbor.X, neighbor.Y, neighborColor)
		}
	}
}
