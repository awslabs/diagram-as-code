package types

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/awslabs/diagram-as-code/internal/vector"
)

func TestLinkInit(t *testing.T) {
	source := new(Resource).Init()
	target := new(Resource).Init()
	sourceArrowHead := ArrowHead{Type: "Default", Length: 10, Width: "Default"}
	targetArrowHead := ArrowHead{Type: "Open", Length: 15, Width: "Wide"}

	link := Link{}.Init(source, WINDROSE_N, sourceArrowHead, target, WINDROSE_S, targetArrowHead, 2, color.RGBA{255, 255, 255, 255})

	if link.Source != source {
		t.Errorf("Expected source node to be %v, got %v", source, link.Source)
	}
	if link.SourcePosition != WINDROSE_N {
		t.Errorf("Expected source position to be 'Top', got %d", link.SourcePosition)
	}
	if link.SourceArrowHead != sourceArrowHead {
		t.Errorf("Expected source arrow head to be %v, got %v", sourceArrowHead, link.SourceArrowHead)
	}
	if link.Target != target {
		t.Errorf("Expected target node to be %v, got %v", target, link.Target)
	}
	if link.TargetPosition != WINDROSE_S {
		t.Errorf("Expected target position to be 'Bottom', got %d", link.TargetPosition)
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
	if link.lineColor != (color.RGBA{255, 255, 255, 255}) {
		t.Errorf("Expected line color to be white, got %v", link.lineColor)
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

// Original computeLabelPos function for comparison
func computeLabelPosOriginal(tx, ty, dx, dy, lx, ly float64) (float64, float64) {
	// Compute the dot product of the unit vectors
	dot_product := tx*dx + ty*dy
	// If the angle is 90 degrees or more (dot product <= 0), set a to (0,0)
	if dot_product > 0 {
		// Compute scalar α
		numerator := ly*dx - lx*dy
		denominator := tx*dy - ty*dx
		// Check for division by zero
		if denominator != 0 {
			alpha := numerator / denominator
			// Compute vector a
			return alpha * tx, alpha * ty
		}
	}
	return 0.0, 0.0
}

func TestComputeLabelPosComparison(t *testing.T) {
	link := Link{
		lineColor: color.RGBA{0, 0, 0, 255},
	}

	tests := []struct {
		name string
		tx, ty, dx, dy, lx, ly float64
	}{
		{
			name: "Perpendicular vectors (90 degrees)",
			tx: 1.0, ty: 0.0,  // East
			dx: 0.0, dy: 1.0,  // South
			lx: 10.0, ly: 5.0,
		},
		{
			name: "Obtuse angle (dot product < 0)",
			tx: 1.0, ty: 0.0,   // East
			dx: -1.0, dy: 0.0,  // West (opposite)
			lx: 10.0, ly: 5.0,
		},
		{
			name: "Acute angle with calculation",
			tx: 1.0, ty: 0.0,    // East
			dx: 0.707, dy: 0.707, // Southeast (45 degrees)
			lx: 0.0, ly: 10.0,   // North label
		},
		{
			name: "Same direction (denominator = 0)",
			tx: 1.0, ty: 0.0,  // East
			dx: 1.0, dy: 0.0,  // Same direction
			lx: 0.0, ly: 10.0, // North label
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Original calculation
			origX, origY := computeLabelPosOriginal(tt.tx, tt.ty, tt.dx, tt.dy, tt.lx, tt.ly)
			
			// New vector calculation
			tVec := vector.New(tt.tx, tt.ty)
			dVec := vector.New(tt.dx, tt.dy)
			labelVec := vector.New(tt.lx, tt.ly)
			result := link.computeLabelPos(tVec, dVec, labelVec)
			
			tolerance := 1e-10
			if math.Abs(result.X-origX) > tolerance || math.Abs(result.Y-origY) > tolerance {
				t.Errorf("Vector version differs from original:\n"+
					"  Original: (%v, %v)\n"+
					"  Vector:   (%v, %v)\n"+
					"  Input: tx=%v, ty=%v, dx=%v, dy=%v, lx=%v, ly=%v",
					origX, origY, result.X, result.Y,
					tt.tx, tt.ty, tt.dx, tt.dy, tt.lx, tt.ly)
			}
		})
	}
}
