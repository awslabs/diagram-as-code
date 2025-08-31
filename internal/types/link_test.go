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
