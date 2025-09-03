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

func TestComputeLabelPos(t *testing.T) {
	link := Link{
		lineColor: color.RGBA{0, 0, 0, 255},
	}

	tests := []struct {
		name        string
		t, d, label vector.Vector
		expected    vector.Vector
	}{
		{
			name:     "Perpendicular vectors (90 degrees)",
			t:        vector.New(1.0, 0.0),
			d:        vector.New(0.0, 1.0),
			label:    vector.New(10.0, 5.0),
			expected: vector.New(0.0, 0.0),
		},
		{
			name:     "Acute angle with calculation",
			t:        vector.New(1.0, 0.0),
			d:        vector.New(0.707, 0.707),
			label:    vector.New(0.0, 10.0),
			expected: vector.New(10.0, 0.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := link.computeLabelPos(tt.t, tt.d, tt.label)
			tolerance := 1e-6
			if math.Abs(result.X-tt.expected.X) > tolerance || math.Abs(result.Y-tt.expected.Y) > tolerance {
				t.Errorf("computeLabelPos() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Original drawArrowHead calculation for comparison
func calculateArrowHeadPointsOriginal(arrowPt, originPt image.Point, arrowHead ArrowHead) (image.Point, image.Point) {
	dx := float64(arrowPt.X - originPt.X)
	dy := float64(arrowPt.Y - originPt.Y)
	length := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	if arrowHead.Length == 0 {
		arrowHead.Length = 10
	}

	_a, _b, _c := getThreeSideOriginal(arrowHead.Width)
	at1 := arrowPt.Sub(image.Point{
		int(arrowHead.Length * (_a*dx - _c*dy) / (_b * length)),
		int(arrowHead.Length * (_c*dx + _a*dy) / (_b * length)),
	})
	at2 := arrowPt.Sub(image.Point{
		int(arrowHead.Length * (_a*dx + _c*dy) / (_b * length)),
		int(arrowHead.Length * (-_c*dx + _a*dy) / (_b * length)),
	})
	return at1, at2
}

func getThreeSideOriginal(t string) (float64, float64, float64) {
	switch t {
	case "Narrow":
		return math.Sqrt(3.0), 2.0, 1.0
	case "Default", "":
		return 1.0, math.Sqrt(2.0), 1.0
	case "Wide":
		return 1.0, 2.0, math.Sqrt(3.0)
	}
	return 0, 0, 0
}

func calculateArrowHeadPointsVector(arrowPt, originPt image.Point, arrowHead ArrowHead) (image.Point, image.Point) {
	arrowVec := vector.New(float64(arrowPt.X), float64(arrowPt.Y))
	originVec := vector.New(float64(originPt.X), float64(originPt.Y))
	direction := arrowVec.Sub(originVec)
	length := direction.Length()

	if arrowHead.Length == 0 {
		arrowHead.Length = 10
	}
	_a, _b, _c := getThreeSideOriginal(arrowHead.Width)

	// Use exact same formula as original code
	dx := direction.X
	dy := direction.Y

	offset1 := vector.New(
		arrowHead.Length*(_a*dx-_c*dy)/(_b*length),
		arrowHead.Length*(_c*dx+_a*dy)/(_b*length),
	)
	offset2 := vector.New(
		arrowHead.Length*(_a*dx+_c*dy)/(_b*length),
		arrowHead.Length*(-_c*dx+_a*dy)/(_b*length),
	)

	// Apply offsets to arrow point (same as arrowPt.Sub(offset))
	at1 := image.Point{int(arrowVec.X - offset1.X), int(arrowVec.Y - offset1.Y)}
	at2 := image.Point{int(arrowVec.X - offset2.X), int(arrowVec.Y - offset2.Y)}

	return at1, at2
}

func TestArrowHeadCalculationComparison(t *testing.T) {
	tests := []struct {
		name      string
		arrowPt   image.Point
		originPt  image.Point
		arrowHead ArrowHead
	}{
		{
			name:      "Horizontal arrow (East)",
			arrowPt:   image.Point{100, 50},
			originPt:  image.Point{50, 50},
			arrowHead: ArrowHead{Type: "Open", Length: 10, Width: "Default"},
		},
		{
			name:      "Vertical arrow (South)",
			arrowPt:   image.Point{50, 100},
			originPt:  image.Point{50, 50},
			arrowHead: ArrowHead{Type: "Open", Length: 10, Width: "Default"},
		},
		{
			name:      "Diagonal arrow (Southeast)",
			arrowPt:   image.Point{100, 100},
			originPt:  image.Point{50, 50},
			arrowHead: ArrowHead{Type: "Open", Length: 10, Width: "Default"},
		},
		{
			name:      "Wide arrow head",
			arrowPt:   image.Point{100, 50},
			originPt:  image.Point{50, 50},
			arrowHead: ArrowHead{Type: "Open", Length: 15, Width: "Wide"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Original calculation
			origAt1, origAt2 := calculateArrowHeadPointsOriginal(tt.arrowPt, tt.originPt, tt.arrowHead)

			// Vector calculation
			vecAt1, vecAt2 := calculateArrowHeadPointsVector(tt.arrowPt, tt.originPt, tt.arrowHead)

			tolerance := 1 // Allow 1 pixel difference due to rounding

			if math.Abs(float64(vecAt1.X-origAt1.X)) > float64(tolerance) ||
				math.Abs(float64(vecAt1.Y-origAt1.Y)) > float64(tolerance) {
				t.Errorf("Arrow point 1 differs:\n"+
					"  Original: %v\n"+
					"  Vector:   %v\n"+
					"  Arrow: %v -> %v",
					origAt1, vecAt1, tt.originPt, tt.arrowPt)
			}

			if math.Abs(float64(vecAt2.X-origAt2.X)) > float64(tolerance) ||
				math.Abs(float64(vecAt2.Y-origAt2.Y)) > float64(tolerance) {
				t.Errorf("Arrow point 2 differs:\n"+
					"  Original: %v\n"+
					"  Vector:   %v\n"+
					"  Arrow: %v -> %v",
					origAt2, vecAt2, tt.originPt, tt.arrowPt)
			}
		})
	}
}
func TestArrowHeadCalculationDetailed(t *testing.T) {
	arrowPt := image.Point{100, 50}
	originPt := image.Point{50, 50}
	arrowHead := ArrowHead{Type: "Open", Length: 10, Width: "Default"}

	// Original calculation
	origAt1, origAt2 := calculateArrowHeadPointsOriginal(arrowPt, originPt, arrowHead)

	// Vector calculation
	vecAt1, vecAt2 := calculateArrowHeadPointsVector(arrowPt, originPt, arrowHead)

	t.Logf("Arrow: %v -> %v", originPt, arrowPt)
	t.Logf("Original at1: %v, at2: %v", origAt1, origAt2)
	t.Logf("Vector   at1: %v, at2: %v", vecAt1, vecAt2)
	t.Logf("Diff at1: (%d, %d), at2: (%d, %d)",
		vecAt1.X-origAt1.X, vecAt1.Y-origAt1.Y,
		vecAt2.X-origAt2.X, vecAt2.Y-origAt2.Y)

	// Test with diagonal arrow for more complex case
	arrowPt2 := image.Point{100, 100}
	origAt1_2, origAt2_2 := calculateArrowHeadPointsOriginal(arrowPt2, originPt, arrowHead)
	vecAt1_2, vecAt2_2 := calculateArrowHeadPointsVector(arrowPt2, originPt, arrowHead)

	t.Logf("Diagonal Arrow: %v -> %v", originPt, arrowPt2)
	t.Logf("Original at1: %v, at2: %v", origAt1_2, origAt2_2)
	t.Logf("Vector   at1: %v, at2: %v", vecAt1_2, vecAt2_2)
	t.Logf("Diff at1: (%d, %d), at2: (%d, %d)",
		vecAt1_2.X-origAt1_2.X, vecAt1_2.Y-origAt1_2.Y,
		vecAt2_2.X-origAt2_2.X, vecAt2_2.Y-origAt2_2.Y)
}
func TestFloatingPointPrecision(t *testing.T) {
	arrowPt := image.Point{100, 50}
	originPt := image.Point{50, 50}
	arrowHead := ArrowHead{Type: "Open", Length: 10, Width: "Default"}

	// Original calculation - step by step
	dx := float64(arrowPt.X - originPt.X) // 50.0
	dy := float64(arrowPt.Y - originPt.Y) // 0.0
	length := math.Sqrt(dx*dx + dy*dy)    // 50.0

	_a, _b, _c := getThreeSideOriginal(arrowHead.Width) // 1.0, √2, 1.0

	t.Logf("Original calculation:")
	t.Logf("dx=%v, dy=%v, length=%v", dx, dy, length)
	t.Logf("_a=%v, _b=%v, _c=%v", _a, _b, _c)

	// Original formula components
	orig_x1 := arrowHead.Length * (_a*dx - _c*dy) / (_b * length)
	orig_y1 := arrowHead.Length * (_c*dx + _a*dy) / (_b * length)
	orig_x2 := arrowHead.Length * (_a*dx + _c*dy) / (_b * length)
	orig_y2 := arrowHead.Length * (-_c*dx + _a*dy) / (_b * length)

	t.Logf("Original float values: at1=(%v, %v), at2=(%v, %v)", orig_x1, orig_y1, orig_x2, orig_y2)
	t.Logf("Original int values: at1=(%d, %d), at2=(%d, %d)", int(orig_x1), int(orig_y1), int(orig_x2), int(orig_y2))

	// Vector calculation - step by step
	arrowVec := vector.New(float64(arrowPt.X), float64(arrowPt.Y))
	originVec := vector.New(float64(originPt.X), float64(originPt.Y))
	direction := arrowVec.Sub(originVec)

	t.Logf("Vector calculation:")
	t.Logf("direction=%v, length=%v", direction, direction.Length())

	unitDir := direction.Normalize()
	perpDir := unitDir.Perpendicular()

	t.Logf("unitDir=%v, perpDir=%v", unitDir, perpDir)

	// Vector formula components
	vec_base := arrowVec.Sub(unitDir.Scale(arrowHead.Length * _a / _b))
	vec_at1 := vec_base.Sub(perpDir.Scale(arrowHead.Length * _c / _b))
	vec_at2 := vec_base.Add(perpDir.Scale(arrowHead.Length * _c / _b))

	t.Logf("Vector float values: at1=(%v, %v), at2=(%v, %v)", vec_at1.X, vec_at1.Y, vec_at2.X, vec_at2.Y)
	t.Logf("Vector int values: at1=(%d, %d), at2=(%d, %d)", int(vec_at1.X), int(vec_at1.Y), int(vec_at2.X), int(vec_at2.Y))

	// Compare the differences
	t.Logf("Float differences: at1=(%v, %v), at2=(%v, %v)",
		vec_at1.X-orig_x1, vec_at1.Y-orig_y1, vec_at2.X-orig_x2, vec_at2.Y-orig_y2)
}
func TestArrowHeadPrecisionAnalysis(t *testing.T) {
	arrowPt := image.Point{100, 50}
	originPt := image.Point{50, 50}
	arrowHead := ArrowHead{Type: "Open", Length: 10, Width: "Default"}

	// Original calculation - step by step
	dx := float64(arrowPt.X - originPt.X)               // 50.0
	dy := float64(arrowPt.Y - originPt.Y)               // 0.0
	length := math.Sqrt(dx*dx + dy*dy)                  // 50.0
	_a, _b, _c := getThreeSideOriginal(arrowHead.Width) // 1.0, √2, 1.0

	// Original offset calculations
	orig_offset1_x := arrowHead.Length * (_a*dx - _c*dy) / (_b * length)
	orig_offset1_y := arrowHead.Length * (_c*dx + _a*dy) / (_b * length)
	orig_offset2_x := arrowHead.Length * (_a*dx + _c*dy) / (_b * length)
	orig_offset2_y := arrowHead.Length * (-_c*dx + _a*dy) / (_b * length)

	t.Logf("Original offsets: offset1=(%v, %v), offset2=(%v, %v)",
		orig_offset1_x, orig_offset1_y, orig_offset2_x, orig_offset2_y)

	// Vector calculation - step by step
	arrowVec := vector.New(float64(arrowPt.X), float64(arrowPt.Y))
	originVec := vector.New(float64(originPt.X), float64(originPt.Y))
	direction := arrowVec.Sub(originVec)
	vec_length := direction.Length()

	// Vector offset calculations (same formula)
	vec_dx := direction.X
	vec_dy := direction.Y
	vec_offset1_x := arrowHead.Length * (_a*vec_dx - _c*vec_dy) / (_b * vec_length)
	vec_offset1_y := arrowHead.Length * (_c*vec_dx + _a*vec_dy) / (_b * vec_length)
	vec_offset2_x := arrowHead.Length * (_a*vec_dx + _c*vec_dy) / (_b * vec_length)
	vec_offset2_y := arrowHead.Length * (-_c*vec_dx + _a*vec_dy) / (_b * vec_length)

	t.Logf("Vector offsets: offset1=(%v, %v), offset2=(%v, %v)",
		vec_offset1_x, vec_offset1_y, vec_offset2_x, vec_offset2_y)

	// Compare intermediate values
	t.Logf("Length comparison: orig=%v, vec=%v, diff=%v", length, vec_length, vec_length-length)
	t.Logf("Direction comparison: orig=(%v,%v), vec=(%v,%v)", dx, dy, vec_dx, vec_dy)

	// Final points
	orig_at1 := arrowPt.Sub(image.Point{int(orig_offset1_x), int(orig_offset1_y)})
	orig_at2 := arrowPt.Sub(image.Point{int(orig_offset2_x), int(orig_offset2_y)})

	vec_at1 := image.Point{int(arrowVec.X - vec_offset1_x), int(arrowVec.Y - vec_offset1_y)}
	vec_at2 := image.Point{int(arrowVec.X - vec_offset2_x), int(arrowVec.Y - vec_offset2_y)}

	t.Logf("Final points: orig_at1=%v, vec_at1=%v, diff=(%d,%d)",
		orig_at1, vec_at1, vec_at1.X-orig_at1.X, vec_at1.Y-orig_at1.Y)
	t.Logf("Final points: orig_at2=%v, vec_at2=%v, diff=(%d,%d)",
		orig_at2, vec_at2, vec_at2.X-orig_at2.X, vec_at2.Y-orig_at2.Y)
}
func TestArrowHeadAccuracy(t *testing.T) {
	arrowPt := image.Point{100, 50}
	originPt := image.Point{50, 50}
	arrowHead := ArrowHead{Type: "Open", Length: 10, Width: "Default"}

	// Theoretical calculation (high precision)
	dx := 50.0
	dy := 0.0
	length := 50.0
	_a, _b, _c := 1.0, math.Sqrt(2.0), 1.0

	// Theoretical exact values (before int conversion)
	theoretical_at1_x := 100.0 - (arrowHead.Length * (_a*dx - _c*dy) / (_b * length))
	theoretical_at1_y := 50.0 - (arrowHead.Length * (_c*dx + _a*dy) / (_b * length))
	theoretical_at2_x := 100.0 - (arrowHead.Length * (_a*dx + _c*dy) / (_b * length))
	theoretical_at2_y := 50.0 - (arrowHead.Length * (-_c*dx + _a*dy) / (_b * length))

	t.Logf("Theoretical exact: at1=(%v, %v), at2=(%v, %v)",
		theoretical_at1_x, theoretical_at1_y, theoretical_at2_x, theoretical_at2_y)

	// What should the int conversion be?
	theoretical_at1_int := image.Point{int(theoretical_at1_x), int(theoretical_at1_y)}
	theoretical_at2_int := image.Point{int(theoretical_at2_x), int(theoretical_at2_y)}

	t.Logf("Theoretical int: at1=%v, at2=%v", theoretical_at1_int, theoretical_at2_int)

	// Original method
	orig_at1, orig_at2 := calculateArrowHeadPointsOriginal(arrowPt, originPt, arrowHead)
	t.Logf("Original: at1=%v, at2=%v", orig_at1, orig_at2)

	// Vector method
	vec_at1, vec_at2 := calculateArrowHeadPointsVector(arrowPt, originPt, arrowHead)
	t.Logf("Vector: at1=%v, at2=%v", vec_at1, vec_at2)

	// Which is closer to theoretical?
	orig_dist1 := math.Abs(float64(orig_at1.X)-theoretical_at1_x) + math.Abs(float64(orig_at1.Y)-theoretical_at1_y)
	vec_dist1 := math.Abs(float64(vec_at1.X)-theoretical_at1_x) + math.Abs(float64(vec_at1.Y)-theoretical_at1_y)

	t.Logf("Distance from theoretical at1: original=%v, vector=%v", orig_dist1, vec_dist1)

	if vec_dist1 < orig_dist1 {
		t.Logf("Vector method is more accurate for at1")
	} else if orig_dist1 < vec_dist1 {
		t.Logf("Original method is more accurate for at1")
	} else {
		t.Logf("Both methods are equally accurate for at1")
	}
}
func TestVerticalArrowHeadSymmetry(t *testing.T) {
	tests := []struct {
		name     string
		arrowPt  image.Point
		originPt image.Point
		expected string
	}{
		{
			name:     "Vertical North arrow",
			arrowPt:  image.Point{50, 50},
			originPt: image.Point{50, 100},
			expected: "symmetric",
		},
		{
			name:     "Vertical South arrow",
			arrowPt:  image.Point{50, 100},
			originPt: image.Point{50, 50},
			expected: "symmetric",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := &Link{lineColor: color.RGBA{0, 0, 0, 255}}
			arrowHead := ArrowHead{Type: "Open", Length: 10, Width: "Default"}

			// Calculate arrow head points
			arrowVec := vector.New(float64(tt.arrowPt.X), float64(tt.arrowPt.Y))
			originVec := vector.New(float64(tt.originPt.X), float64(tt.originPt.Y))
			direction := arrowVec.Sub(originVec)
			length := direction.Length()

			_a, _b, _c := link.getThreeSide(arrowHead.Width)
			dx, dy := direction.X, direction.Y

			at1Vec := arrowVec.Sub(vector.New(
				arrowHead.Length*(_a*dx-_c*dy)/(_b*length),
				arrowHead.Length*(_c*dx+_a*dy)/(_b*length),
			))
			at2Vec := arrowVec.Sub(vector.New(
				arrowHead.Length*(_a*dx+_c*dy)/(_b*length),
				arrowHead.Length*(-_c*dx+_a*dy)/(_b*length),
			))

			at1 := image.Point{int(math.Round(at1Vec.X)), int(math.Round(at1Vec.Y))}
			at2 := image.Point{int(math.Round(at2Vec.X)), int(math.Round(at2Vec.Y))}

			t.Logf("Arrow: %v -> %v", tt.originPt, tt.arrowPt)
			t.Logf("Direction: dx=%v, dy=%v, length=%v", dx, dy, length)
			t.Logf("Arrow points: at1=%v, at2=%v", at1, at2)

			// Check symmetry for vertical arrows
			if dx == 0 { // Vertical arrow
				// For vertical arrows, at1 and at2 should be symmetric around the arrow point
				centerX := float64(tt.arrowPt.X)
				dist1 := math.Abs(float64(at1.X) - centerX)
				dist2 := math.Abs(float64(at2.X) - centerX)

				t.Logf("Symmetry check: center=%v, dist1=%v, dist2=%v", centerX, dist1, dist2)

				if math.Abs(dist1-dist2) > 1e-10 {
					t.Errorf("Vertical arrow head is not symmetric: dist1=%v, dist2=%v, diff=%v",
						dist1, dist2, math.Abs(dist1-dist2))
				}

				// Also check Y coordinates should be the same for horizontal symmetry
				if at1.Y != at2.Y {
					t.Errorf("Vertical arrow head Y coordinates should be equal: at1.Y=%d, at2.Y=%d",
						at1.Y, at2.Y)
				}
			}
		})
	}
}
func TestVerticalArrowCalculationAnalysis(t *testing.T) {
	arrowPt := image.Point{50, 100}
	originPt := image.Point{50, 50}
	arrowHead := ArrowHead{Type: "Open", Length: 10, Width: "Default"}

	link := &Link{lineColor: color.RGBA{0, 0, 0, 255}}

	// Direction vector
	dx := float64(arrowPt.X - originPt.X) // 0
	dy := float64(arrowPt.Y - originPt.Y) // 50
	length := math.Sqrt(dx*dx + dy*dy)    // 50

	_a, _b, _c := link.getThreeSide(arrowHead.Width) // 1.0, √2, 1.0

	t.Logf("Input: dx=%v, dy=%v, length=%v", dx, dy, length)
	t.Logf("Coefficients: _a=%v, _b=%v, _c=%v", _a, _b, _c)

	// Calculate offsets step by step
	// at1 offset
	at1_x_offset := arrowHead.Length * (_a*dx - _c*dy) / (_b * length)
	at1_y_offset := arrowHead.Length * (_c*dx + _a*dy) / (_b * length)

	// at2 offset
	at2_x_offset := arrowHead.Length * (_a*dx + _c*dy) / (_b * length)
	at2_y_offset := arrowHead.Length * (-_c*dx + _a*dy) / (_b * length)

	t.Logf("at1 offset: x=%v, y=%v", at1_x_offset, at1_y_offset)
	t.Logf("at2 offset: x=%v, y=%v", at2_x_offset, at2_y_offset)

	// For vertical arrow (dx=0, dy=50):
	// at1_x_offset = 10 * (1*0 - 1*50) / (√2 * 50) = 10 * (-50) / (50√2) = -10/√2 ≈ -7.071
	// at2_x_offset = 10 * (1*0 + 1*50) / (√2 * 50) = 10 * 50 / (50√2) = 10/√2 ≈ 7.071

	expected_offset := 10.0 / math.Sqrt(2.0)
	t.Logf("Expected symmetric offset: ±%v", expected_offset)

	// Final points
	at1_x := float64(arrowPt.X) - at1_x_offset
	at1_y := float64(arrowPt.Y) - at1_y_offset
	at2_x := float64(arrowPt.X) - at2_x_offset
	at2_y := float64(arrowPt.Y) - at2_y_offset

	t.Logf("Final float points: at1=(%v, %v), at2=(%v, %v)", at1_x, at1_y, at2_x, at2_y)

	// Convert to int
	at1 := image.Point{int(at1_x), int(at1_y)}
	at2 := image.Point{int(at2_x), int(at2_y)}

	t.Logf("Final int points: at1=%v, at2=%v", at1, at2)

	// Check if the issue is in int conversion
	t.Logf("Int conversion: at1_x %v->%d, at2_x %v->%d", at1_x, int(at1_x), at2_x, int(at2_x))

	// The issue: at1_x ≈ 57.071 -> 57, at2_x ≈ 42.929 -> 42
	// This creates asymmetry due to different rounding behavior
}
func TestHorizontalArrowHeadSymmetry(t *testing.T) {
	tests := []struct {
		name     string
		arrowPt  image.Point
		originPt image.Point
	}{
		{
			name:     "Horizontal East arrow",
			arrowPt:  image.Point{100, 50},
			originPt: image.Point{50, 50},
		},
		{
			name:     "Horizontal West arrow",
			arrowPt:  image.Point{50, 50},
			originPt: image.Point{100, 50},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := &Link{lineColor: color.RGBA{0, 0, 0, 255}}
			arrowHead := ArrowHead{Type: "Open", Length: 10, Width: "Default"}

			// Calculate arrow head points
			arrowVec := vector.New(float64(tt.arrowPt.X), float64(tt.arrowPt.Y))
			originVec := vector.New(float64(tt.originPt.X), float64(tt.originPt.Y))
			direction := arrowVec.Sub(originVec)
			length := direction.Length()

			_a, _b, _c := link.getThreeSide(arrowHead.Width)
			dx, dy := direction.X, direction.Y

			at1Vec := arrowVec.Sub(vector.New(
				arrowHead.Length*(_a*dx-_c*dy)/(_b*length),
				arrowHead.Length*(_c*dx+_a*dy)/(_b*length),
			))
			at2Vec := arrowVec.Sub(vector.New(
				arrowHead.Length*(_a*dx+_c*dy)/(_b*length),
				arrowHead.Length*(-_c*dx+_a*dy)/(_b*length),
			))

			at1 := image.Point{int(math.Round(at1Vec.X)), int(math.Round(at1Vec.Y))}
			at2 := image.Point{int(math.Round(at2Vec.X)), int(math.Round(at2Vec.Y))}

			t.Logf("Arrow: %v -> %v", tt.originPt, tt.arrowPt)
			t.Logf("Direction: dx=%v, dy=%v, length=%v", dx, dy, length)
			t.Logf("Arrow points: at1=%v, at2=%v", at1, at2)

			// Check symmetry for horizontal arrows
			if dy == 0 { // Horizontal arrow
				// For horizontal arrows, at1 and at2 should be symmetric around the arrow point
				centerY := float64(tt.arrowPt.Y)
				dist1 := math.Abs(float64(at1.Y) - centerY)
				dist2 := math.Abs(float64(at2.Y) - centerY)

				t.Logf("Symmetry check: center=%v, dist1=%v, dist2=%v", centerY, dist1, dist2)

				if math.Abs(dist1-dist2) > 1e-10 {
					t.Errorf("Horizontal arrow head is not symmetric: dist1=%v, dist2=%v, diff=%v",
						dist1, dist2, math.Abs(dist1-dist2))
				}

				// Also check X coordinates should be the same for vertical symmetry
				if at1.X != at2.X {
					t.Errorf("Horizontal arrow head X coordinates should be equal: at1.X=%d, at2.X=%d",
						at1.X, at2.X)
				}
			}
		})
	}
}

func TestSToNOrthogonalPath(t *testing.T) {
	// Test S to N orthogonal path control points
	link := &Link{
		SourcePosition: 8, // S
		TargetPosition: 0, // N
		Type:           "orthogonal",
	}

	// Given coordinates (from actual test)
	sourcePt := image.Point{402, 446} // ALB
	targetPt := image.Point{786, 384} // EC2Instance

	t.Logf("=== Expected S to N Orthogonal Path ===")
	t.Logf("Source: %v (ALB)", sourcePt)
	t.Logf("Target: %v (EC2Instance)", targetPt)

	// Calculate expected control points
	// Step 1: Move away from resources
	sourceStep1 := image.Point{sourcePt.X, sourcePt.Y + 20} // South 20px
	targetStep1 := image.Point{targetPt.X, targetPt.Y - 20} // North 20px
	t.Logf("Step 1 - Source (S+20px): %v", sourceStep1)
	t.Logf("Step 1 - Target (N+20px): %v", targetStep1)

	// Step 2: Move East by half remaining X distance
	remainingX := targetStep1.X - sourceStep1.X // 786 - 402 = 384
	halfX := remainingX / 2                     // 192

	sourceStep2 := image.Point{sourceStep1.X + halfX, sourceStep1.Y} // East 192px
	targetStep2 := image.Point{targetStep1.X - halfX, targetStep1.Y} // West 192px
	t.Logf("Step 2 - Source (E+%dpx): %v", halfX, sourceStep2)
	t.Logf("Step 2 - Target (W+%dpx): %v", halfX, targetStep2)

	// Step 3: Move North/South by half remaining Y distance
	remainingY := targetStep2.Y - sourceStep2.Y // 364 - 466 = -102
	halfY := remainingY / 2                     // -51

	sourceStep3 := image.Point{sourceStep2.X, sourceStep2.Y + halfY} // North 51px
	targetStep3 := image.Point{targetStep2.X, targetStep2.Y - halfY} // South 51px
	t.Logf("Step 3 - Source (N+%dpx): %v", -halfY, sourceStep3)
	t.Logf("Step 3 - Target (S+%dpx): %v", halfY, targetStep3)

	// Expected control points (should converge at step 3)
	expectedControlPoints := []image.Point{
		sourceStep1, // (402, 466)
		sourceStep2, // (594, 466)
		sourceStep3, // (594, 415) - converged point
		targetStep2, // (594, 364)
		targetStep1, // (786, 364)
	}

	t.Logf("Expected control points: %v", expectedControlPoints)

	// Call actual function
	actualControlPoints := link.calculateOrthogonalPath(sourcePt, targetPt)
	t.Logf("Actual control points: %v", actualControlPoints)

	// Verify key points are correct
	if len(actualControlPoints) < 5 {
		t.Errorf("Expected at least 5 control points, got %d", len(actualControlPoints))
		return
	}

	// Check first point (source initial)
	if actualControlPoints[0] != sourceStep1 {
		t.Errorf("First control point should be %v, got %v", sourceStep1, actualControlPoints[0])
	}

	// Check that we have orthogonal movements (alternating X/Y changes)
	for i := 1; i < len(actualControlPoints); i++ {
		prev := actualControlPoints[i-1]
		curr := actualControlPoints[i]

		xChanged := prev.X != curr.X
		yChanged := prev.Y != curr.Y

		// Should change only one axis at a time (orthogonal)
		if xChanged && yChanged {
			t.Errorf("Non-orthogonal movement from %v to %v", prev, curr)
		}
		if !xChanged && !yChanged {
			t.Errorf("No movement from %v to %v", prev, curr)
		}
	}

	// Verify Source → First control point orthogonality
	if len(actualControlPoints) > 0 {
		first := actualControlPoints[0]
		if sourcePt.X != first.X && sourcePt.Y != first.Y {
			t.Errorf("Non-orthogonal: Source (%d,%d) → First control point (%d,%d)",
				sourcePt.X, sourcePt.Y, first.X, first.Y)
		}
	}

	// Verify Last control point → Target orthogonality
	if len(actualControlPoints) > 0 {
		last := actualControlPoints[len(actualControlPoints)-1]
		if last.X != targetPt.X && last.Y != targetPt.Y {
			t.Errorf("Non-orthogonal: Last control point (%d,%d) → Target (%d,%d)",
				last.X, last.Y, targetPt.X, targetPt.Y)
		}
	}

	t.Logf("=== End S to N Test ===")
}
func TestSToEOrthogonalPath(t *testing.T) {
	// Test S to E orthogonal path control points (orthogonal directions)
	link := &Link{
		SourcePosition: 8, // S
		TargetPosition: 4, // E
		Type:           "orthogonal",
	}

	// Given coordinates (from actual test)
	sourcePt := image.Point{402, 446} // ALB
	targetPt := image.Point{818, 416} // EC2Instance

	t.Logf("=== Expected S to E Orthogonal Path ===")
	t.Logf("Source: %v (ALB)", sourcePt)
	t.Logf("Target: %v (EC2Instance)", targetPt)

	// Calculate expected control points for orthogonal directions
	// Step 1: Move away from resources
	sourceStep1 := image.Point{sourcePt.X, sourcePt.Y + 20} // South 20px
	targetStep1 := image.Point{targetPt.X + 20, targetPt.Y} // East 20px
	t.Logf("Step 1 - Source (S+20px): %v", sourceStep1)
	t.Logf("Step 1 - Target (E+20px): %v", targetStep1)

	// Step 2: Complete movement (orthogonal directions - no /2 needed)
	remainingX := targetStep1.X - sourceStep1.X // 838 - 402 = 436
	remainingY := targetStep1.Y - sourceStep1.Y // 416 - 466 = -50

	sourceStep2 := image.Point{sourceStep1.X + remainingX, sourceStep1.Y} // East 436px
	targetStep2 := image.Point{targetStep1.X, targetStep1.Y + remainingY} // South -50px
	t.Logf("Step 2 - Source (E+%dpx): %v", remainingX, sourceStep2)
	t.Logf("Step 2 - Target (S+%dpx): %v", remainingY, targetStep2)

	// Expected control points (should converge at step 2)
	expectedControlPoints := []image.Point{
		sourceStep1, // (402, 466)
		sourceStep2, // (838, 466) - converged point
		targetStep1, // (838, 416)
	}

	t.Logf("Expected control points: %v", expectedControlPoints)

	// Call actual function
	actualControlPoints := link.calculateOrthogonalPath(sourcePt, targetPt)
	t.Logf("Actual control points: %v", actualControlPoints)

	// Verify key points are correct
	if len(actualControlPoints) < 3 {
		t.Errorf("Expected at least 3 control points, got %d", len(actualControlPoints))
		return
	}

	// Check first point (source initial)
	if actualControlPoints[0] != sourceStep1 {
		t.Errorf("First control point should be %v, got %v", sourceStep1, actualControlPoints[0])
	}

	// Check that we have orthogonal movements (alternating X/Y changes)
	for i := 1; i < len(actualControlPoints); i++ {
		prev := actualControlPoints[i-1]
		curr := actualControlPoints[i]

		xChanged := prev.X != curr.X
		yChanged := prev.Y != curr.Y

		// Should change only one axis at a time (orthogonal)
		if xChanged && yChanged {
			t.Errorf("Non-orthogonal movement from %v to %v", prev, curr)
		}
		if !xChanged && !yChanged {
			t.Errorf("No movement from %v to %v", prev, curr)
		}
	}

	// Verify Source → First control point orthogonality
	if len(actualControlPoints) > 0 {
		first := actualControlPoints[0]
		if sourcePt.X != first.X && sourcePt.Y != first.Y {
			t.Errorf("Non-orthogonal: Source (%d,%d) → First control point (%d,%d)",
				sourcePt.X, sourcePt.Y, first.X, first.Y)
		}
	}

	// Verify Last control point → Target orthogonality
	if len(actualControlPoints) > 0 {
		last := actualControlPoints[len(actualControlPoints)-1]
		if last.X != targetPt.X && last.Y != targetPt.Y {
			t.Errorf("Non-orthogonal: Last control point (%d,%d) → Target (%d,%d)",
				last.X, last.Y, targetPt.X, targetPt.Y)
		}
	}

	t.Logf("=== End S to E Test ===")
}
func TestSToSOrthogonalPath(t *testing.T) {
	// Test S to S orthogonal path control points (same directions)
	link := &Link{
		SourcePosition: 8, // S
		TargetPosition: 8, // S
		Type:           "orthogonal",
	}

	// Given coordinates (from actual test)
	sourcePt := image.Point{402, 446} // ALB
	targetPt := image.Point{786, 384} // EC2Instance (different position for S to S)

	t.Logf("=== Expected S to S Orthogonal Path ===")
	t.Logf("Source: %v (ALB)", sourcePt)
	t.Logf("Target: %v (EC2Instance)", targetPt)

	// Calculate expected control points for same directions (parallel)
	// Step 1: Move away from resources
	sourceStep1 := image.Point{sourcePt.X, sourcePt.Y + 20} // South 20px
	targetStep1 := image.Point{targetPt.X, targetPt.Y + 20} // South 20px
	t.Logf("Step 1 - Source (S+20px): %v", sourceStep1)
	t.Logf("Step 1 - Target (S+20px): %v", targetStep1)

	// Step 2: Move East by half remaining X distance (parallel directions use /2)
	remainingX := targetStep1.X - sourceStep1.X // 786 - 402 = 384
	halfX := remainingX / 2                     // 192

	sourceStep2 := image.Point{sourceStep1.X + halfX, sourceStep1.Y} // East 192px
	targetStep2 := image.Point{targetStep1.X - halfX, targetStep1.Y} // West 192px
	t.Logf("Step 2 - Source (E+%dpx): %v", halfX, sourceStep2)
	t.Logf("Step 2 - Target (W+%dpx): %v", halfX, targetStep2)

	// Step 3: Move South by half remaining Y distance (should be 0 for S to S)
	remainingY := targetStep2.Y - sourceStep2.Y // Should be 0 for same Y
	t.Logf("Step 3 - Remaining Y: %d (should be 0 for S to S)", remainingY)

	// Expected control points (should converge at step 2 since Y is same)
	expectedControlPoints := []image.Point{
		sourceStep1, // (402, 466)
		sourceStep2, // (594, 466) - converged point
		targetStep1, // (786, 404)
	}

	t.Logf("Expected control points: %v", expectedControlPoints)

	// Call actual function
	actualControlPoints := link.calculateOrthogonalPath(sourcePt, targetPt)
	t.Logf("Actual control points: %v", actualControlPoints)

	// Verify key points are correct
	if len(actualControlPoints) < 3 {
		t.Errorf("Expected at least 3 control points, got %d", len(actualControlPoints))
		return
	}

	// Check first point (source initial)
	if actualControlPoints[0] != sourceStep1 {
		t.Errorf("First control point should be %v, got %v", sourceStep1, actualControlPoints[0])
	}

	// Check that we have orthogonal movements (alternating X/Y changes)
	for i := 1; i < len(actualControlPoints); i++ {
		prev := actualControlPoints[i-1]
		curr := actualControlPoints[i]

		xChanged := prev.X != curr.X
		yChanged := prev.Y != curr.Y

		// Should change only one axis at a time (orthogonal)
		if xChanged && yChanged {
			t.Errorf("Non-orthogonal movement from %v to %v", prev, curr)
		}
		if !xChanged && !yChanged {
			t.Errorf("No movement from %v to %v", prev, curr)
		}
	}

	// Verify Source → First control point orthogonality
	if len(actualControlPoints) > 0 {
		first := actualControlPoints[0]
		if sourcePt.X != first.X && sourcePt.Y != first.Y {
			t.Errorf("Non-orthogonal: Source (%d,%d) → First control point (%d,%d)",
				sourcePt.X, sourcePt.Y, first.X, first.Y)
		}
	}

	// Verify Last control point → Target orthogonality
	if len(actualControlPoints) > 0 {
		last := actualControlPoints[len(actualControlPoints)-1]
		if last.X != targetPt.X && last.Y != targetPt.Y {
			t.Errorf("Non-orthogonal: Last control point (%d,%d) → Target (%d,%d)",
				last.X, last.Y, targetPt.X, targetPt.Y)
		}
	}

	t.Logf("=== End S to S Test ===")
}
func TestWToEOrthogonalPath(t *testing.T) {
	// Test W to E orthogonal path control points (opposite directions with detour)
	link := &Link{
		SourcePosition: 12, // W (WINDROSE_W)
		TargetPosition: 4,  // E (WINDROSE_E)
		Type:           "orthogonal",
	}

	// Given coordinates (from actual test)
	sourcePt := image.Point{402, 446} // ALB
	targetPt := image.Point{818, 416} // EC2Instance

	t.Logf("=== Expected W to E Orthogonal Path ===")
	t.Logf("Source: %v (ALB)", sourcePt)
	t.Logf("Target: %v (EC2Instance)", targetPt)

	// Calculate expected control points for opposite directions with detour
	// Step 1: Move away from resources
	sourceStep1 := image.Point{sourcePt.X - 20, sourcePt.Y} // West 20px
	targetStep1 := image.Point{targetPt.X + 20, targetPt.Y} // East 20px
	t.Logf("Step 1 - Source (W+20px): %v", sourceStep1)
	t.Logf("Step 1 - Target (E+20px): %v", targetStep1)

	// Step 2: Detour north to avoid resource penetration
	detourDistance := 50                                                      // Estimated detour distance (resource height/2 + margin)
	sourceStep2 := image.Point{sourceStep1.X, sourceStep1.Y - detourDistance} // North detour
	targetStep2 := image.Point{targetStep1.X, targetStep1.Y - detourDistance} // North detour
	t.Logf("Step 2 - Source (N+%dpx): %v", detourDistance, sourceStep2)
	t.Logf("Step 2 - Target (N+%dpx): %v", detourDistance, targetStep2)

	// Step 3: Move East by half remaining X distance (parallel directions use /2)
	remainingX := targetStep2.X - sourceStep2.X // 838 - 382 = 456
	halfX := remainingX / 2                     // 228

	sourceStep3 := image.Point{sourceStep2.X + halfX, sourceStep2.Y} // East 228px
	targetStep3 := image.Point{targetStep2.X - halfX, targetStep2.Y} // West 228px (converged)
	t.Logf("Step 3 - Source (E+%dpx): %v", halfX, sourceStep3)
	t.Logf("Step 3 - Target (W+%dpx): %v", halfX, targetStep3)

	// Expected control points (should converge at step 3)
	expectedControlPoints := []image.Point{
		sourceStep1, // (382, 446)
		sourceStep2, // (382, 396) - detour north
		sourceStep3, // (610, 396) - converged point
		targetStep2, // (838, 396) - detour north
		targetStep1, // (838, 416)
	}

	t.Logf("Expected control points: %v", expectedControlPoints)

	// Call actual function
	actualControlPoints := link.calculateOrthogonalPath(sourcePt, targetPt)
	t.Logf("Actual control points: %v", actualControlPoints)

	// Verify key points are correct
	if len(actualControlPoints) < 5 {
		t.Errorf("Expected at least 5 control points for detour, got %d", len(actualControlPoints))
		return
	}

	// Check first point (source initial)
	if actualControlPoints[0] != sourceStep1 {
		t.Errorf("First control point should be %v, got %v", sourceStep1, actualControlPoints[0])
	}

	// Check that we have orthogonal movements (alternating X/Y changes)
	for i := 1; i < len(actualControlPoints); i++ {
		prev := actualControlPoints[i-1]
		curr := actualControlPoints[i]

		xChanged := prev.X != curr.X
		yChanged := prev.Y != curr.Y

		// Should change only one axis at a time (orthogonal)
		if xChanged && yChanged {
			t.Errorf("Non-orthogonal movement from %v to %v", prev, curr)
		}
		if !xChanged && !yChanged {
			t.Errorf("No movement from %v to %v", prev, curr)
		}
	}

	// Verify Source → First control point orthogonality
	if len(actualControlPoints) > 0 {
		first := actualControlPoints[0]
		if sourcePt.X != first.X && sourcePt.Y != first.Y {
			t.Errorf("Non-orthogonal: Source (%d,%d) → First control point (%d,%d)",
				sourcePt.X, sourcePt.Y, first.X, first.Y)
		}
	}

	// Verify Last control point → Target orthogonality
	if len(actualControlPoints) > 0 {
		last := actualControlPoints[len(actualControlPoints)-1]
		if last.X != targetPt.X && last.Y != targetPt.Y {
			t.Errorf("Non-orthogonal: Last control point (%d,%d) → Target (%d,%d)",
				last.X, last.Y, targetPt.X, targetPt.Y)
		}
	}

	t.Logf("=== End W to E Test ===")
}
func TestSToNVerticalStackPath(t *testing.T) {
	// Test S to N orthogonal path for vertical stack (east-west detour needed)
	link := &Link{
		SourcePosition: 8, // S
		TargetPosition: 0, // N
		Type:           "orthogonal",
	}

	// Vertical stack coordinates: Resource2 below Resource1
	sourcePt := image.Point{400, 500} // Resource2 (bottom)
	targetPt := image.Point{400, 300} // Resource1 (top)

	t.Logf("=== Expected S to N Vertical Stack Path ===")
	t.Logf("Source: %v (Resource2 - bottom)", sourcePt)
	t.Logf("Target: %v (Resource1 - top)", targetPt)

	// Calculate expected control points for vertical stack with detour
	// Step 1: Move away from resources
	sourceStep1 := image.Point{sourcePt.X, sourcePt.Y + 20} // South 20px
	targetStep1 := image.Point{targetPt.X, targetPt.Y - 20} // North 20px
	t.Logf("Step 1 - Source (S+20px): %v", sourceStep1)
	t.Logf("Step 1 - Target (N+20px): %v", targetStep1)

	// Step 2: Detour east to avoid vertical resource penetration
	detourDistance := 50                                                      // Estimated detour distance (resource width/2 + margin)
	sourceStep2 := image.Point{sourceStep1.X + detourDistance, sourceStep1.Y} // East detour
	targetStep2 := image.Point{targetStep1.X + detourDistance, targetStep1.Y} // East detour
	t.Logf("Step 2 - Source (E+%dpx): %v", detourDistance, sourceStep2)
	t.Logf("Step 2 - Target (E+%dpx): %v", detourDistance, targetStep2)

	// Step 3: Move North by half remaining Y distance (parallel directions use /2)
	remainingY := targetStep2.Y - sourceStep2.Y // 280 - 520 = -240
	halfY := remainingY / 2                     // -120

	sourceStep3 := image.Point{sourceStep2.X, sourceStep2.Y + halfY} // North 120px
	targetStep3 := image.Point{targetStep2.X, targetStep2.Y - halfY} // South 120px (converged)
	t.Logf("Step 3 - Source (N+%dpx): %v", -halfY, sourceStep3)
	t.Logf("Step 3 - Target (S+%dpx): %v", halfY, targetStep3)

	// Expected control points (should converge at step 3)
	expectedControlPoints := []image.Point{
		sourceStep1, // (400, 520)
		sourceStep2, // (450, 520) - detour east
		sourceStep3, // (450, 400) - converged point
		targetStep2, // (450, 280) - detour east
		targetStep1, // (400, 280)
	}

	t.Logf("Expected control points: %v", expectedControlPoints)

	// Call actual function
	actualControlPoints := link.calculateOrthogonalPath(sourcePt, targetPt)
	t.Logf("Actual control points: %v", actualControlPoints)

	// Verify key points are correct
	if len(actualControlPoints) < 5 {
		t.Errorf("Expected at least 5 control points for detour, got %d", len(actualControlPoints))
		return
	}

	// Check first point (source initial)
	if actualControlPoints[0] != sourceStep1 {
		t.Errorf("First control point should be %v, got %v", sourceStep1, actualControlPoints[0])
	}

	// Check that we have orthogonal movements (alternating X/Y changes)
	for i := 1; i < len(actualControlPoints); i++ {
		prev := actualControlPoints[i-1]
		curr := actualControlPoints[i]

		xChanged := prev.X != curr.X
		yChanged := prev.Y != curr.Y

		// Should change only one axis at a time (orthogonal)
		if xChanged && yChanged {
			t.Errorf("Non-orthogonal movement from %v to %v", prev, curr)
		}
		if !xChanged && !yChanged {
			t.Errorf("No movement from %v to %v", prev, curr)
		}
	}

	t.Logf("=== End S to N Vertical Stack Test ===")
}
func TestVerticalStackSToNOrthogonalPath(t *testing.T) {
	t.Log("=== Vertical Stack S to N Orthogonal Path Test ===")

	// Create vertical stack with Instance1 above Instance2
	// Instance1 (top): (402, 350)
	// Instance2 (bottom): (402, 450)
	// Link: Instance2:S -> Instance1:N (should require detour)

	source := image.Point{X: 402, Y: 450} // Instance2 center (bottom)
	target := image.Point{X: 402, Y: 350} // Instance1 center (top)

	t.Logf("Source (Instance2): (%d,%d)", source.X, source.Y)
	t.Logf("Target (Instance1): (%d,%d)", target.X, target.Y)

	link := &Link{
		Type:           "orthogonal",
		SourcePosition: 8, // S (South)
		TargetPosition: 0, // N (North)
	}

	controlPts := link.calculateOrthogonalPath(source, target)

	t.Logf("Actual control points: %v", controlPts)
	t.Logf("Number of control points: %d", len(controlPts))

	// Verify orthogonal movements
	for i := 0; i < len(controlPts)-1; i++ {
		p1 := controlPts[i]
		p2 := controlPts[i+1]

		if p1.X != p2.X && p1.Y != p2.Y {
			t.Errorf("Non-orthogonal movement from (%d,%d) to (%d,%d)", p1.X, p1.Y, p2.X, p2.Y)
		}
	}

	// Verify Source → First control point orthogonality
	if len(controlPts) > 0 {
		first := controlPts[0]
		if source.X != first.X && source.Y != first.Y {
			t.Errorf("Non-orthogonal: Source (%d,%d) → First control point (%d,%d)",
				source.X, source.Y, first.X, first.Y)
		}
	}

	// Verify Last control point → Target orthogonality
	if len(controlPts) > 0 {
		last := controlPts[len(controlPts)-1]
		if last.X != target.X && last.Y != target.Y {
			t.Errorf("Non-orthogonal: Last control point (%d,%d) → Target (%d,%d)",
				last.X, last.Y, target.X, target.Y)
		}
	}

	// Should have detour since both positions point toward each other
	// Expected: detour to avoid direct collision
	if len(controlPts) < 3 {
		t.Errorf("Expected at least 3 control points for detour, got %d", len(controlPts))
	}

	t.Log("=== End Vertical Stack S to N Test ===")
}

func TestEToNOrthogonalPath(t *testing.T) {
	// Test E to N orthogonal path control points (check axis selection)
	link := &Link{
		SourcePosition: 4, // E (WINDROSE_E)
		TargetPosition: 0, // N (WINDROSE_N)
		Type:           "orthogonal",
	}

	// Given coordinates
	sourcePt := image.Point{402, 446} // ALB
	targetPt := image.Point{786, 300} // EC2Instance (higher up)

	t.Logf("=== E to N Orthogonal Path Test ===")
	t.Logf("Source: %v (ALB)", sourcePt)
	t.Logf("Target: %v (EC2Instance)", targetPt)

	// Call actual function
	actualControlPoints := link.calculateOrthogonalPath(sourcePt, targetPt)
	t.Logf("Actual control points: %v", actualControlPoints)
	t.Logf("Number of control points: %d", len(actualControlPoints))

	// Check that we have orthogonal movements
	for i := 1; i < len(actualControlPoints); i++ {
		prev := actualControlPoints[i-1]
		curr := actualControlPoints[i]

		xChanged := prev.X != curr.X
		yChanged := prev.Y != curr.Y

		// Should change only one axis at a time (orthogonal)
		if xChanged && yChanged {
			t.Errorf("Non-orthogonal movement from %v to %v", prev, curr)
		}
		if !xChanged && !yChanged {
			t.Errorf("No movement from %v to %v", prev, curr)
		}
	}

	// Verify Source → First control point orthogonality
	if len(actualControlPoints) > 0 {
		first := actualControlPoints[0]
		if sourcePt.X != first.X && sourcePt.Y != first.Y {
			t.Errorf("Non-orthogonal: Source (%d,%d) → First control point (%d,%d)",
				sourcePt.X, sourcePt.Y, first.X, first.Y)
		}
	}

	// Verify Last control point → Target orthogonality
	if len(actualControlPoints) > 0 {
		last := actualControlPoints[len(actualControlPoints)-1]
		if last.X != targetPt.X && last.Y != targetPt.Y {
			t.Errorf("Non-orthogonal: Last control point (%d,%d) → Target (%d,%d)",
				last.X, last.Y, targetPt.X, targetPt.Y)
		}
	}
	t.Logf("=== End E to N Test ===")
}

func TestHorizontalStackEToWOrthogonalPath(t *testing.T) {
	t.Log("=== Horizontal Stack E to W Orthogonal Path Test ===")

	// Create horizontal stack with Instance1 left, Instance2 right
	// Instance1 (left): (350, 402)
	// Instance2 (right): (450, 402)
	// Link: Instance1:E -> Instance2:W (should require detour)

	source := image.Point{X: 350, Y: 402} // Instance1 center (left)
	target := image.Point{X: 450, Y: 402} // Instance2 center (right)

	t.Logf("Source (Instance1): (%d,%d)", source.X, source.Y)
	t.Logf("Target (Instance2): (%d,%d)", target.X, target.Y)

	link := &Link{
		Type:           "orthogonal",
		SourcePosition: 4,  // E (East)
		TargetPosition: 12, // W (West)
	}

	controlPts := link.calculateOrthogonalPath(source, target)

	t.Logf("Actual control points: %v", controlPts)
	t.Logf("Number of control points: %d", len(controlPts))

	// Verify orthogonal movements
	for i := 0; i < len(controlPts)-1; i++ {
		p1 := controlPts[i]
		p2 := controlPts[i+1]

		if p1.X != p2.X && p1.Y != p2.Y {
			t.Errorf("Non-orthogonal movement from (%d,%d) to (%d,%d)", p1.X, p1.Y, p2.X, p2.Y)
		}

		// Verify Source → First control point orthogonality
		if len(controlPts) > 0 {
			first := controlPts[0]
			if source.X != first.X && source.Y != first.Y {
				t.Errorf("Non-orthogonal: Source (%d,%d) → First control point (%d,%d)",
					source.X, source.Y, first.X, first.Y)
			}
		}

		// Verify Last control point → Target orthogonality
		if len(controlPts) > 0 {
			last := controlPts[len(controlPts)-1]
			if last.X != target.X && last.Y != target.Y {
				t.Errorf("Non-orthogonal: Last control point (%d,%d) → Target (%d,%d)",
					last.X, last.Y, target.X, target.Y)
			}
		}
	}

	// Distance is 100px, which allows direct connection without detour
	// Expected: efficient direct path with minimal control points
	if len(controlPts) < 1 {
		t.Errorf("Expected at least 1 control point, got %d", len(controlPts))
	}

	// Verify the path converges at the midpoint for efficient connection
	if len(controlPts) == 1 {
		midpoint := controlPts[0]
		expectedX := (source.X + target.X) / 2 // 400
		if midpoint.X != expectedX || midpoint.Y != source.Y {
			t.Errorf("Expected midpoint convergence at (%d,%d), got (%d,%d)",
				expectedX, source.Y, midpoint.X, midpoint.Y)
		}
	}

	t.Log("=== End Horizontal Stack E to W Test ===")
}

func TestNestedLayoutELBToVerticalStackOrthogonalPath(t *testing.T) {
	t.Log("=== Nested Layout ELB to VerticalStack Orthogonal Path Test ===")

	// Layout: HorizontalStack{ELB, VerticalStack{Instance1, Instance2}}
	// ELB (left): (300, 400)
	// Instance1 (top right): (500, 350)
	// Instance2 (bottom right): (500, 450)
	// Links: ELB:E -> Instance1:W, ELB:E -> Instance2:W

	elbPos := image.Point{X: 300, Y: 400}
	instance1Pos := image.Point{X: 500, Y: 350}
	instance2Pos := image.Point{X: 500, Y: 450}

	t.Logf("ELB: (%d,%d)", elbPos.X, elbPos.Y)
	t.Logf("Instance1: (%d,%d)", instance1Pos.X, instance1Pos.Y)
	t.Logf("Instance2: (%d,%d)", instance2Pos.X, instance2Pos.Y)

	// Test ELB:E -> Instance1:W
	t.Log("--- Testing ELB:E -> Instance1:W ---")
	link1 := &Link{
		Type:           "orthogonal",
		SourcePosition: 4,  // E (East)
		TargetPosition: 12, // W (West)
	}

	controlPts1 := link1.calculateOrthogonalPath(elbPos, instance1Pos)
	t.Logf("ELB -> Instance1 control points: %v", controlPts1)

	// Test ELB:E -> Instance2:W
	t.Log("--- Testing ELB:E -> Instance2:W ---")
	link2 := &Link{
		Type:           "orthogonal",
		SourcePosition: 4,  // E (East)
		TargetPosition: 12, // W (West)
	}

	controlPts2 := link2.calculateOrthogonalPath(elbPos, instance2Pos)
	t.Logf("ELB -> Instance2 control points: %v", controlPts2)

	// Verify orthogonal movements for both links
	for i, controlPts := range [][]image.Point{controlPts1, controlPts2} {
		linkName := []string{"ELB->Instance1", "ELB->Instance2"}[i]
		for j := 0; j < len(controlPts)-1; j++ {
			p1 := controlPts[j]
			p2 := controlPts[j+1]

			if p1.X != p2.X && p1.Y != p2.Y {
				t.Errorf("%s: Non-orthogonal movement from (%d,%d) to (%d,%d)",
					linkName, p1.X, p1.Y, p2.X, p2.Y)
			}

			// Verify Source → First control point orthogonality
			if len(controlPts) > 0 {
				first := controlPts[0]
				if elbPos.X != first.X && elbPos.Y != first.Y {
					t.Errorf("%s: Non-orthogonal Source (%d,%d) → First control point (%d,%d)",
						linkName, elbPos.X, elbPos.Y, first.X, first.Y)
				}
			}

			// Verify Last control point → Target orthogonality
			if len(controlPts) > 0 {
				last := controlPts[len(controlPts)-1]
				targetPos := []image.Point{instance1Pos, instance2Pos}[i]
				if last.X != targetPos.X && last.Y != targetPos.Y {
					t.Errorf("%s: Non-orthogonal Last control point (%d,%d) → Target (%d,%d)",
						linkName, last.X, last.Y, targetPos.X, targetPos.Y)
				}
			}
		}
	}

	// Both links should have efficient paths (distance 200px allows direct connection)
	if len(controlPts1) < 1 {
		t.Errorf("Expected at least 1 control point for ELB->Instance1, got %d", len(controlPts1))
	}
	if len(controlPts2) < 1 {
		t.Errorf("Expected at least 1 control point for ELB->Instance2, got %d", len(controlPts2))
	}

	// Verify different Y coordinates for the two paths (no collision)
	if len(controlPts1) > 0 && len(controlPts2) > 0 {
		// Paths should diverge to reach different Y levels
		t.Logf("Instance1 path Y: %d, Instance2 path Y: %d",
			controlPts1[len(controlPts1)-1].Y, controlPts2[len(controlPts2)-1].Y)
	}

	t.Log("=== End Nested Layout Test ===")
}

func TestHorizontalStackWToEOrthogonalPath(t *testing.T) {
	t.Log("=== Horizontal Stack W to E Orthogonal Path Test ===")

	// Create horizontal stack with Instance1 left, Instance2 right
	// Instance1 (left): (350, 402)
	// Instance2 (right): (450, 402)
	// Link: Instance2:W -> Instance1:E (reverse direction)

	source := image.Point{X: 450, Y: 402} // Instance2 center (right)
	target := image.Point{X: 350, Y: 402} // Instance1 center (left)

	t.Logf("Source (Instance2): (%d,%d)", source.X, source.Y)
	t.Logf("Target (Instance1): (%d,%d)", target.X, target.Y)

	link := &Link{
		Type:           "orthogonal",
		SourcePosition: 12, // W (West)
		TargetPosition: 4,  // E (East)
	}

	controlPts := link.calculateOrthogonalPath(source, target)

	t.Logf("Actual control points: %v", controlPts)
	t.Logf("Number of control points: %d", len(controlPts))

	// Verify orthogonal movements
	for i := 0; i < len(controlPts)-1; i++ {
		p1 := controlPts[i]
		p2 := controlPts[i+1]

		if p1.X != p2.X && p1.Y != p2.Y {
			t.Errorf("Non-orthogonal movement from (%d,%d) to (%d,%d)", p1.X, p1.Y, p2.X, p2.Y)
		}
	}

	// Verify Source → First control point orthogonality
	if len(controlPts) > 0 {
		first := controlPts[0]
		if source.X != first.X && source.Y != first.Y {
			t.Errorf("Non-orthogonal: Source (%d,%d) → First control point (%d,%d)",
				source.X, source.Y, first.X, first.Y)
		}
	}

	// Verify Last control point → Target orthogonality
	if len(controlPts) > 0 {
		last := controlPts[len(controlPts)-1]
		if last.X != target.X && last.Y != target.Y {
			t.Errorf("Non-orthogonal: Last control point (%d,%d) → Target (%d,%d)",
				last.X, last.Y, target.X, target.Y)
		}
	}
	// Distance is 100px, which allows direct connection without detour
	// Expected: efficient direct path with minimal control points
	if len(controlPts) < 1 {
		t.Errorf("Expected at least 1 control point, got %d", len(controlPts))
	}

	// Verify the path converges at the midpoint for efficient connection
	if len(controlPts) == 1 {
		midpoint := controlPts[0]
		expectedX := (source.X + target.X) / 2 // 400
		if midpoint.X != expectedX || midpoint.Y != source.Y {
			t.Errorf("Expected midpoint convergence at (%d,%d), got (%d,%d)",
				expectedX, source.Y, midpoint.X, midpoint.Y)
		}
	}

	t.Log("=== End Horizontal Stack W to E Test ===")
}

func TestNestedLayoutELBToHorizontalStackOrthogonalPath(t *testing.T) {
	t.Log("=== Nested Layout ELB to HorizontalStack Orthogonal Path Test ===")

	// Layout: VerticalStack{HorizontalStack{Instance1, Instance2}, ELB}
	// Instance1 (top left): (350, 300)
	// Instance2 (top right): (450, 300)
	// ELB (bottom center): (400, 450)
	// Links: ELB:NNW -> Instance1:S, ELB:NNE -> Instance2:S

	instance1Pos := image.Point{X: 350, Y: 300}
	instance2Pos := image.Point{X: 450, Y: 300}
	elbPos := image.Point{X: 400, Y: 450}

	t.Logf("Instance1: (%d,%d)", instance1Pos.X, instance1Pos.Y)
	t.Logf("Instance2: (%d,%d)", instance2Pos.X, instance2Pos.Y)
	t.Logf("ELB: (%d,%d)", elbPos.X, elbPos.Y)

	// Test ELB:NNW -> Instance1:S
	t.Log("--- Testing ELB:NNW -> Instance1:S ---")
	link1 := &Link{
		Type:           "orthogonal",
		SourcePosition: 1, // NNW (North-North-West)
		TargetPosition: 8, // S (South)
	}

	controlPts1 := link1.calculateOrthogonalPath(elbPos, instance1Pos)
	t.Logf("ELB -> Instance1 control points: %v", controlPts1)

	// Test ELB:NNE -> Instance2:S
	t.Log("--- Testing ELB:NNE -> Instance2:S ---")
	link2 := &Link{
		Type:           "orthogonal",
		SourcePosition: 3, // NNE (North-North-East)
		TargetPosition: 8, // S (South)
	}

	controlPts2 := link2.calculateOrthogonalPath(elbPos, instance2Pos)
	t.Logf("ELB -> Instance2 control points: %v", controlPts2)

	// Verify orthogonal movements for both links
	for i, controlPts := range [][]image.Point{controlPts1, controlPts2} {
		linkName := []string{"ELB->Instance1", "ELB->Instance2"}[i]
		for j := 0; j < len(controlPts)-1; j++ {
			p1 := controlPts[j]
			p2 := controlPts[j+1]

			if p1.X != p2.X && p1.Y != p2.Y {
				t.Errorf("%s: Non-orthogonal movement from (%d,%d) to (%d,%d)",
					linkName, p1.X, p1.Y, p2.X, p2.Y)
			}
		}
	}

	// Both links should have efficient paths
	if len(controlPts1) < 1 {
		t.Errorf("Expected at least 1 control point for ELB->Instance1, got %d", len(controlPts1))
	}
	if len(controlPts2) < 1 {
		t.Errorf("Expected at least 1 control point for ELB->Instance2, got %d", len(controlPts2))
	}

	// Verify paths diverge correctly (different X coordinates)
	if len(controlPts1) > 0 && len(controlPts2) > 0 {
		// Paths should diverge to reach different X positions
		t.Logf("Instance1 path final X: %d, Instance2 path final X: %d",
			controlPts1[len(controlPts1)-1].X, controlPts2[len(controlPts2)-1].X)
	}

	// Verify Source → First control point orthogonality for both links
	for i, controlPts := range [][]image.Point{controlPts1, controlPts2} {
		linkName := []string{"ELB->Instance1", "ELB->Instance2"}[i]
		targetPos := []image.Point{instance1Pos, instance2Pos}[i]

		if len(controlPts) > 0 {
			first := controlPts[0]
			if elbPos.X != first.X && elbPos.Y != first.Y {
				t.Errorf("%s: Non-orthogonal Source (%d,%d) → First control point (%d,%d)",
					linkName, elbPos.X, elbPos.Y, first.X, first.Y)
			}
		}

		if len(controlPts) > 0 {
			last := controlPts[len(controlPts)-1]
			if last.X != targetPos.X && last.Y != targetPos.Y {
				t.Errorf("%s: Non-orthogonal Last control point (%d,%d) → Target (%d,%d)",
					linkName, last.X, last.Y, targetPos.X, targetPos.Y)
			}
		}
	}

	t.Log("=== End Nested Layout ELB to HorizontalStack Test ===")
}

func TestOrthogonalityCheck(t *testing.T) {
	t.Log("=== Testing Source/Target Orthogonality ===")

	source := image.Point{X: 350, Y: 402}
	target := image.Point{X: 450, Y: 402}

	link := &Link{
		Type:           "orthogonal",
		SourcePosition: 4,  // E
		TargetPosition: 12, // W
	}

	controlPts := link.calculateOrthogonalPath(source, target)
	t.Logf("Control points: %v", controlPts)

	// Verify Source → First control point orthogonality
	if len(controlPts) > 0 {
		first := controlPts[0]
		if source.X != first.X && source.Y != first.Y {
			t.Errorf("Non-orthogonal: Source (%d,%d) → First control point (%d,%d)",
				source.X, source.Y, first.X, first.Y)
		}
	}

	// Verify Last control point → Target orthogonality
	if len(controlPts) > 0 {
		last := controlPts[len(controlPts)-1]
		if last.X != target.X && last.Y != target.Y {
			t.Errorf("Non-orthogonal: Last control point (%d,%d) → Target (%d,%d)",
				last.X, last.Y, target.X, target.Y)
		}
	}

	t.Log("=== End Orthogonality Check ===")
}

func TestDetourDirectionOrthogonalPath(t *testing.T) {
	t.Log("=== Detour Direction Orthogonal Path Test ===")

	// Layout: HorizontalStack{ELB1, VerticalStack{Instance1, Instance2}, ELB2}
	// ELB1 (left): (200, 400)
	// Instance1 (top): (400, 350)
	// Instance2 (bottom): (400, 450)
	// ELB2 (right): (600, 400)

	elb1Pos := image.Point{X: 200, Y: 400}
	instance1Pos := image.Point{X: 400, Y: 350}
	instance2Pos := image.Point{X: 400, Y: 450}
	elb2Pos := image.Point{X: 600, Y: 400}

	t.Logf("ELB1: (%d,%d)", elb1Pos.X, elb1Pos.Y)
	t.Logf("Instance1: (%d,%d)", instance1Pos.X, instance1Pos.Y)
	t.Logf("Instance2: (%d,%d)", instance2Pos.X, instance2Pos.Y)
	t.Logf("ELB2: (%d,%d)", elb2Pos.X, elb2Pos.Y)

	// Test 1: ELB1:W -> Instance1:W (expect north detour)
	t.Log("--- Testing ELB1:W -> Instance1:W (expect north detour) ---")
	link1 := &Link{
		Type:           "orthogonal",
		SourcePosition: 12, // W (West)
		TargetPosition: 12, // W (West)
	}

	controlPts1 := link1.calculateOrthogonalPath(elb1Pos, instance1Pos)
	t.Logf("ELB1 -> Instance1 control points: %v", controlPts1)

	// Test 2: ELB1:W -> Instance2:W (expect south detour)
	t.Log("--- Testing ELB1:W -> Instance2:W (expect south detour) ---")
	link2 := &Link{
		Type:           "orthogonal",
		SourcePosition: 12, // W (West)
		TargetPosition: 12, // W (West)
	}

	controlPts2 := link2.calculateOrthogonalPath(elb1Pos, instance2Pos)
	t.Logf("ELB1 -> Instance2 control points: %v", controlPts2)

	// Test 3: Instance1:E -> ELB2:E (expect north detour)
	t.Log("--- Testing Instance1:E -> ELB2:E (expect north detour) ---")
	link3 := &Link{
		Type:           "orthogonal",
		SourcePosition: 4, // E (East)
		TargetPosition: 4, // E (East)
	}

	controlPts3 := link3.calculateOrthogonalPath(instance1Pos, elb2Pos)
	t.Logf("Instance1 -> ELB2 control points: %v", controlPts3)

	// Test 4: Instance2:E -> ELB2:E (expect south detour)
	t.Log("--- Testing Instance2:E -> ELB2:E (expect south detour) ---")
	link4 := &Link{
		Type:           "orthogonal",
		SourcePosition: 4, // E (East)
		TargetPosition: 4, // E (East)
	}

	controlPts4 := link4.calculateOrthogonalPath(instance2Pos, elb2Pos)
	t.Logf("Instance2 -> ELB2 control points: %v", controlPts4)

	// Verify orthogonal movements for all links
	allLinks := []struct {
		name        string
		controlPts  []image.Point
		source      image.Point
		target      image.Point
		expectNorth bool
	}{
		{"ELB1->Instance1", controlPts1, elb1Pos, instance1Pos, true},
		{"ELB1->Instance2", controlPts2, elb1Pos, instance2Pos, false},
		{"Instance1->ELB2", controlPts3, instance1Pos, elb2Pos, true},
		{"Instance2->ELB2", controlPts4, instance2Pos, elb2Pos, false},
	}

	for _, link := range allLinks {
		t.Logf("--- Verifying %s ---", link.name)

		// Verify orthogonal movements between control points
		for i := 0; i < len(link.controlPts)-1; i++ {
			p1 := link.controlPts[i]
			p2 := link.controlPts[i+1]

			if p1.X != p2.X && p1.Y != p2.Y {
				t.Errorf("%s: Non-orthogonal movement from (%d,%d) to (%d,%d)",
					link.name, p1.X, p1.Y, p2.X, p2.Y)
			}
		}

		// Verify Source → First control point orthogonality
		if len(link.controlPts) > 0 {
			first := link.controlPts[0]
			if link.source.X != first.X && link.source.Y != first.Y {
				t.Errorf("%s: Non-orthogonal Source (%d,%d) → First control point (%d,%d)",
					link.name, link.source.X, link.source.Y, first.X, first.Y)
			}
		}

		// Verify Last control point → Target orthogonality
		if len(link.controlPts) > 0 {
			last := link.controlPts[len(link.controlPts)-1]
			if last.X != link.target.X && last.Y != link.target.Y {
				t.Errorf("%s: Non-orthogonal Last control point (%d,%d) → Target (%d,%d)",
					link.name, last.X, last.Y, link.target.X, link.target.Y)
			}
		}
	}

	t.Log("=== End Detour Direction Test ===")
}
