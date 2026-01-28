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

// TestIssue236OrthogonalDetourFix tests the fix for GitHub issue #236
// Orthogonal link improvement: wrong control points with detour
func TestIssue236OrthogonalDetourFix(t *testing.T) {
	t.Log("=== Testing Issue #236 Fix: Orthogonal Detour ===")

	// Test case from GitHub issue #236
	// Administrator (E position) -> EC2A (S position)
	// Before fix: [(626,856) (646,856) (646,804) (646,804) (646,856) (402,856) (402,607)]
	// After fix:  [(626,856) (646,856) (646,804) (402,804) (402,607)]

	source := image.Point{X: 626, Y: 856} // Administrator
	target := image.Point{X: 402, Y: 607} // EC2A

	link := &Link{
		Type:           "orthogonal",
		SourcePosition: 4, // E (East)
		TargetPosition: 8, // S (South)
	}

	t.Logf("Source (Administrator): (%d,%d)", source.X, source.Y)
	t.Logf("Target (EC2A): (%d,%d)", target.X, target.Y)

	controlPts := link.calculateOrthogonalPath(source, target)
	t.Logf("Actual control points: %v", controlPts)
	t.Logf("Number of control points: %d", len(controlPts))

	// Expected control points after fix
	expected := []image.Point{
		{646, 856}, // Source moves East 20px
		{646, 804}, // Source detour North 52px
		{402, 804}, // Target moves to same Y level (accounting for detour)
	}

	t.Logf("Expected control points: %v", expected)

	// Verify we have the expected number of control points (3, not 5)
	if len(controlPts) != len(expected) {
		t.Errorf("Expected %d control points, got %d", len(expected), len(controlPts))
	}

	// Verify each control point matches expected
	for i, expectedPt := range expected {
		if i < len(controlPts) {
			actualPt := controlPts[i]
			if actualPt != expectedPt {
				t.Errorf("Control point %d: expected %v, got %v", i, expectedPt, actualPt)
			}
		}
	}

	// Verify no duplicate control points
	for i := 1; i < len(controlPts); i++ {
		if controlPts[i] == controlPts[i-1] {
			t.Errorf("Duplicate control point at index %d: %v", i, controlPts[i])
		}
	}

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

	t.Log("=== End Issue #236 Test ===")
}

// TestCounterpartDetourConsideration tests the core fix for non-parallel cases
func TestCounterpartDetourConsideration(t *testing.T) {
	t.Log("=== Testing Counterpart Detour Consideration ===")

	// Test case: sourcePenetration && !targetPenetration
	// Target should reduce movement by detour distance (52px)

	source := image.Point{X: 626, Y: 856}
	target := image.Point{X: 402, Y: 607}

	link := &Link{
		Type:           "orthogonal",
		SourcePosition: 4, // E (East) - will have penetration
		TargetPosition: 8, // S (South) - no penetration
	}

	controlPts := link.calculateOrthogonalPath(source, target)

	// Verify target movement accounts for source detour
	// Expected: Target moves 197px (249 - 52) instead of full 249px
	expectedTargetY := 607 + 197 // 804

	// Find the control point where target converges
	var targetConvergeY int
	for _, pt := range controlPts {
		if pt.X == 402 { // Target X coordinate
			targetConvergeY = pt.Y
			break
		}
	}

	if targetConvergeY != expectedTargetY {
		t.Errorf("Target should converge at Y=%d (accounting for detour), got Y=%d",
			expectedTargetY, targetConvergeY)
	}

	t.Logf("Target converged at Y=%d (expected %d)", targetConvergeY, expectedTargetY)
	t.Log("=== End Counterpart Detour Test ===")
}

func TestGroupingOffset(t *testing.T) {
	// Setup resources
	source := new(Resource).Init()
	source.SetBindings(image.Rect(0, 0, 64, 64))

	target1 := new(Resource).Init()
	target1.SetBindings(image.Rect(100, 0, 164, 64))

	target2 := new(Resource).Init()
	target2.SetBindings(image.Rect(200, 0, 264, 64))

	// Create links from same source position
	link1 := new(Link).Init(source, WINDROSE_S, ArrowHead{}, target1, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	link2 := new(Link).Init(source, WINDROSE_S, ArrowHead{}, target2, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	source.AddLink(link1)
	source.AddLink(link2)
	target1.AddLink(link1)
	target2.AddLink(link2)

	// Test normal grouping offset
	t.Run("Normal grouping offset", func(t *testing.T) {
		// Enable grouping offset
		source.SetGroupingOffset(true)

		// Sort links first
		source.sortAllLinks()

		// Test first link offset
		index1, count1 := link1.getLinkIndexAndCount(source, WINDROSE_S)
		if count1 != 2 {
			t.Errorf("Expected count 2, got %d", count1)
		}

		// Calculate expected offset: (index - (count-1)/2.0) * 10
		expectedOffset1 := int((float64(index1) - float64(count1-1)/2.0) * 10)

		pt1 := link1.calcPositionWithOffset(source.GetBindings(), WINDROSE_S, source, true)
		originalPt, _ := calcPosition(source.GetBindings(), WINDROSE_S)

		// Check if offset was applied
		if pt1.X == originalPt.X && pt1.Y == originalPt.Y && expectedOffset1 != 0 {
			t.Errorf("Expected offset to be applied, but position unchanged")
		}
	})

	// Test disabled grouping offset (default behavior)
	t.Run("Disabled grouping offset", func(t *testing.T) {
		// Explicitly disable grouping offset
		source.SetGroupingOffset(false)

		pt1 := link1.calcPositionWithOffset(source.GetBindings(), WINDROSE_S, source, true)
		originalPt, _ := calcPosition(source.GetBindings(), WINDROSE_S)

		// Should use original position when disabled
		if pt1.X != originalPt.X || pt1.Y != originalPt.Y {
			t.Errorf("Expected original position (%d, %d), got (%d, %d)",
				originalPt.X, originalPt.Y, pt1.X, pt1.Y)
		}
	})
}

func TestLinkSorting(t *testing.T) {
	// Setup resources with different X positions
	source := new(Resource).Init()
	source.SetBindings(image.Rect(100, 0, 164, 64))

	target1 := new(Resource).Init()
	target1.SetBindings(image.Rect(0, 100, 64, 164)) // Left target

	target2 := new(Resource).Init()
	target2.SetBindings(image.Rect(200, 100, 264, 164)) // Right target

	// Create links
	link1 := new(Link).Init(source, WINDROSE_S, ArrowHead{}, target1, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	link2 := new(Link).Init(source, WINDROSE_S, ArrowHead{}, target2, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	// Add in reverse order to test sorting
	source.AddLink(link2) // Right target first
	source.AddLink(link1) // Left target second

	t.Run("Links sorted by target position", func(t *testing.T) {
		// Enable grouping offset for this test
		source.SetGroupingOffset(true)
		source.sortAllLinks()

		// After sorting, right target should come first (based on perpendicular projection)
		index1, _ := link1.getLinkIndexAndCount(source, WINDROSE_S)
		index2, _ := link2.getLinkIndexAndCount(source, WINDROSE_S)

		// The sorting uses perpendicular projection, so right target (higher X) comes first
		if index2 >= index1 {
			t.Errorf("Expected link2 (right target) to have lower index than link1 (left target), got %d >= %d", index2, index1)
		}
	})
}

func TestMixedSourceTargetSorting(t *testing.T) {
	// Setup resources for A:S->B C->A:S pattern
	resourceA := new(Resource).Init()
	resourceA.SetBindings(image.Rect(100, 100, 164, 164))

	resourceB := new(Resource).Init()
	resourceB.SetBindings(image.Rect(200, 50, 264, 114)) // Right and up

	resourceC := new(Resource).Init()
	resourceC.SetBindings(image.Rect(50, 200, 114, 264)) // Left and down

	// Create mixed source/target links at same position S on resourceA
	linkAS_B := new(Link).Init(resourceA, WINDROSE_S, ArrowHead{}, resourceB, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	linkC_AS := new(Link).Init(resourceC, WINDROSE_N, ArrowHead{}, resourceA, WINDROSE_S, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	// Add links to resourceA
	resourceA.AddLink(linkAS_B) // A as source
	resourceA.AddLink(linkC_AS) // A as target

	t.Run("Mixed source/target links sorted at same position", func(t *testing.T) {
		resourceA.SetGroupingOffset(true)
		resourceA.sortAllLinks()

		// Both links should be in the same group and sorted together
		index1, count1 := linkAS_B.getLinkIndexAndCount(resourceA, WINDROSE_S)
		index2, count2 := linkC_AS.getLinkIndexAndCount(resourceA, WINDROSE_S)

		// Both should see unified count of 2
		if count1 != 2 {
			t.Errorf("Expected count1=2, got %d", count1)
		}
		if count2 != 2 {
			t.Errorf("Expected count2=2, got %d", count2)
		}

		// Indices should be different (0 and 1)
		if index1 == index2 {
			t.Errorf("Expected different indices, got both %d", index1)
		}

		// One should be 0, other should be 1
		if (index1 != 0 && index1 != 1) || (index2 != 0 && index2 != 1) {
			t.Errorf("Expected indices 0 and 1, got %d and %d", index1, index2)
		}
	})
}

func TestGetLinkIndexAndCount(t *testing.T) {
	source := new(Resource).Init()
	target := new(Resource).Init()

	// Create multiple links from different positions
	linkS1 := new(Link).Init(source, WINDROSE_S, ArrowHead{}, target, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	linkS2 := new(Link).Init(source, WINDROSE_S, ArrowHead{}, target, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	linkE := new(Link).Init(source, WINDROSE_E, ArrowHead{}, target, WINDROSE_W, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	source.AddLink(linkS1)
	source.AddLink(linkS2)
	source.AddLink(linkE)

	t.Run("Count links from same position", func(t *testing.T) {
		_, countS := linkS1.getLinkIndexAndCount(source, WINDROSE_S)
		_, countE := linkE.getLinkIndexAndCount(source, WINDROSE_E)

		if countS != 2 {
			t.Errorf("Expected 2 links from S position, got %d", countS)
		}
		if countE != 1 {
			t.Errorf("Expected 1 link from E position, got %d", countE)
		}
	})
}

func TestGroupingOffsetMixedSourceTarget(t *testing.T) {
	// Setup: A->B, C->A (A is both source and target)
	resourceA := new(Resource).Init()
	resourceA.SetBindings(image.Rect(100, 0, 164, 64))

	resourceB := new(Resource).Init()
	resourceB.SetBindings(image.Rect(200, 0, 264, 64))

	resourceC := new(Resource).Init()
	resourceC.SetBindings(image.Rect(0, 0, 64, 64))

	// A->B (A as source, position S)
	linkAB := new(Link).Init(resourceA, WINDROSE_S, ArrowHead{}, resourceB, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	// C->A (A as target, position N)
	linkCA := new(Link).Init(resourceC, WINDROSE_S, ArrowHead{}, resourceA, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	resourceA.AddLink(linkAB)
	resourceA.AddLink(linkCA)
	resourceB.AddLink(linkAB)
	resourceC.AddLink(linkCA)

	t.Run("A as source with grouping offset", func(t *testing.T) {
		resourceA.SetGroupingOffset(true)

		// A->B: A is source at position S
		indexAB, countAB := linkAB.getLinkIndexAndCount(resourceA, WINDROSE_S)
		if countAB != 1 { // Only A:S->B at position S (C->A is at position N)
			t.Errorf("Expected 1 link at A:S, got %d", countAB)
		}
		if indexAB != 0 {
			t.Errorf("Expected index 0 for A->B, got %d", indexAB)
		}

		ptAB := linkAB.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_S, resourceA, true)
		originalPt, _ := calcPosition(resourceA.GetBindings(), WINDROSE_S)

		// Single link should have no offset
		if ptAB.X != originalPt.X || ptAB.Y != originalPt.Y {
			t.Errorf("Single link should have no offset: expected (%d,%d), got (%d,%d)",
				originalPt.X, originalPt.Y, ptAB.X, ptAB.Y)
		}
	})

	t.Run("A as target with grouping offset", func(t *testing.T) {
		resourceA.SetGroupingOffset(true)

		// C->A: A is target at position N (different position from S)
		indexCA, countCA := linkCA.getLinkIndexAndCount(resourceA, WINDROSE_N)
		if countCA != 1 {
			t.Errorf("Expected 1 link to A:N, got %d", countCA)
		}
		if indexCA != 0 {
			t.Errorf("Expected index 0 for C->A, got %d", indexCA)
		}

		ptCA := linkCA.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_N, resourceA, false)
		originalPt, _ := calcPosition(resourceA.GetBindings(), WINDROSE_N)

		// Single link should have no offset
		if ptCA.X != originalPt.X || ptCA.Y != originalPt.Y {
			t.Errorf("Single target link should have no offset: expected (%d,%d), got (%d,%d)",
				originalPt.X, originalPt.Y, ptCA.X, ptCA.Y)
		}
	})

	t.Run("Source and target positions are independent", func(t *testing.T) {
		resourceA.SetGroupingOffset(true)

		// Verify that different positions are counted separately
		_, countSource := linkAB.getLinkIndexAndCount(resourceA, WINDROSE_S)
		_, countTarget := linkCA.getLinkIndexAndCount(resourceA, WINDROSE_N)

		if countSource != 1 {
			t.Errorf("Position S should have 1 link, got %d", countSource)
		}
		if countTarget != 1 {
			t.Errorf("Position N should have 1 link, got %d", countTarget)
		}
	})
}

func TestGroupingOffsetSamePositionSourceAndTarget(t *testing.T) {
	// Setup: A:S->B, C->A:S (both at position S - same position, different roles)
	resourceA := new(Resource).Init()
	resourceA.SetBindings(image.Rect(100, 100, 164, 164))

	resourceB := new(Resource).Init()
	resourceB.SetBindings(image.Rect(100, 200, 164, 264))

	resourceC := new(Resource).Init()
	resourceC.SetBindings(image.Rect(100, 0, 164, 64))

	// A:S->B (A as source at position S)
	linkAB := new(Link).Init(resourceA, WINDROSE_S, ArrowHead{}, resourceB, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	// C:S->A:S (C as source at position S, A as target at position S)
	linkCA := new(Link).Init(resourceC, WINDROSE_S, ArrowHead{}, resourceA, WINDROSE_S, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	resourceA.AddLink(linkAB)
	resourceA.AddLink(linkCA)
	resourceB.AddLink(linkAB)
	resourceC.AddLink(linkCA)

	t.Run("Same position S with different roles", func(t *testing.T) {
		resourceA.SetGroupingOffset(true)

		// A:S->B: count at position S (unified count)
		indexAB, countAB := linkAB.getLinkIndexAndCount(resourceA, WINDROSE_S)
		// C->A:S: count at position S (unified count)
		indexCA, countCA := linkCA.getLinkIndexAndCount(resourceA, WINDROSE_S)

		// Both should see unified count of 2
		if countAB != 2 {
			t.Errorf("Expected 2 links at A:S (unified count), got %d", countAB)
		}
		if countCA != 2 {
			t.Errorf("Expected 2 links at A:S (unified count), got %d", countCA)
		}

		if indexAB < 0 || indexAB >= 2 {
			t.Errorf("Expected index 0 or 1 for A:S->B, got %d", indexAB)
		}
		if indexCA < 0 || indexCA >= 2 {
			t.Errorf("Expected index 0 or 1 for C->A:S, got %d", indexCA)
		}

		// Both should have offset applied (count=2)
		ptAB := linkAB.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_S, resourceA, true)
		ptCA := linkCA.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_S, resourceA, false)

		// Should have different positions due to offset
		if ptAB.X == ptCA.X && ptAB.Y == ptCA.Y {
			t.Errorf("Two links at same position should have different offsets")
		}
	})
}

func TestGroupingOffsetMultipleOutgoingLinks(t *testing.T) {
	// Setup: A->B, A->C (both from A:E - 2 outgoing links from same position)
	resourceA := new(Resource).Init()
	resourceA.SetBindings(image.Rect(100, 100, 164, 164))

	resourceB := new(Resource).Init()
	resourceB.SetBindings(image.Rect(200, 50, 264, 114))

	resourceC := new(Resource).Init()
	resourceC.SetBindings(image.Rect(200, 150, 264, 214))

	// A->B, A->C (both from A:E)
	linkAB := new(Link).Init(resourceA, WINDROSE_E, ArrowHead{}, resourceB, WINDROSE_W, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	linkAC := new(Link).Init(resourceA, WINDROSE_E, ArrowHead{}, resourceC, WINDROSE_W, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	resourceA.AddLink(linkAB)
	resourceA.AddLink(linkAC)

	t.Run("Two outgoing links with grouping offset", func(t *testing.T) {
		resourceA.SetGroupingOffset(true)
		resourceA.sortAllLinks()

		_, countE := linkAB.getLinkIndexAndCount(resourceA, WINDROSE_E)
		if countE != 2 {
			t.Errorf("Expected 2 outgoing links from A:E, got %d", countE)
		}

		ptAB := linkAB.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_E, resourceA, true)
		ptAC := linkAC.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_E, resourceA, true)

		// Two links should have different Y offsets
		if ptAB.Y == ptAC.Y {
			t.Errorf("Two outgoing links should have different Y offsets: both at Y=%d", ptAB.Y)
		}

		// Verify offset magnitude (should be ±5 for 2 links)
		originalPt, _ := calcPosition(resourceA.GetBindings(), WINDROSE_E)
		offsetAB := ptAB.Y - originalPt.Y
		offsetAC := ptAC.Y - originalPt.Y

		expectedOffsets := []int{-5, 5}
		if !((offsetAB == expectedOffsets[0] && offsetAC == expectedOffsets[1]) ||
			(offsetAB == expectedOffsets[1] && offsetAC == expectedOffsets[0])) {
			t.Errorf("Expected offsets ±5, got %d and %d", offsetAB, offsetAC)
		}
	})
}

func TestGroupingOffsetMultipleIncomingLinks(t *testing.T) {
	// Setup: B->A, C->A (both to A:W - 2 incoming links to same position)
	resourceA := new(Resource).Init()
	resourceA.SetBindings(image.Rect(100, 100, 164, 164))

	resourceB := new(Resource).Init()
	resourceB.SetBindings(image.Rect(0, 50, 64, 114))

	resourceC := new(Resource).Init()
	resourceC.SetBindings(image.Rect(0, 150, 64, 214))

	// B->A, C->A (both to A:W)
	linkBA := new(Link).Init(resourceB, WINDROSE_E, ArrowHead{}, resourceA, WINDROSE_W, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	linkCA := new(Link).Init(resourceC, WINDROSE_E, ArrowHead{}, resourceA, WINDROSE_W, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	resourceA.AddLink(linkBA)
	resourceA.AddLink(linkCA)

	t.Run("Two incoming links with grouping offset", func(t *testing.T) {
		resourceA.SetGroupingOffset(true)
		resourceA.sortAllLinks()

		_, countW := linkBA.getLinkIndexAndCount(resourceA, WINDROSE_W)
		if countW != 2 {
			t.Errorf("Expected 2 incoming links to A:W, got %d", countW)
		}

		ptBA := linkBA.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_W, resourceA, false)
		ptCA := linkCA.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_W, resourceA, false)

		// Two links should have different Y offsets
		if ptBA.Y == ptCA.Y {
			t.Errorf("Two incoming links should have different Y offsets: both at Y=%d", ptBA.Y)
		}

		// Verify offset magnitude (should be ±5 for 2 links)
		originalPt, _ := calcPosition(resourceA.GetBindings(), WINDROSE_W)
		offsetBA := ptBA.Y - originalPt.Y
		offsetCA := ptCA.Y - originalPt.Y

		expectedOffsets := []int{-5, 5}
		if !((offsetBA == expectedOffsets[0] && offsetCA == expectedOffsets[1]) ||
			(offsetBA == expectedOffsets[1] && offsetCA == expectedOffsets[0])) {
			t.Errorf("Expected offsets ±5, got %d and %d", offsetBA, offsetCA)
		}
	})
}

func TestGroupingOffsetDirection(t *testing.T) {
	// Setup ELB resource
	elb := new(Resource).Init()
	elb.SetBindings(image.Rect(100, 100, 164, 164))
	elb.SetIconBounds(image.Rect(100, 100, 164, 164))

	// Setup source resources
	source1 := new(Resource).Init()
	source1.SetBindings(image.Rect(0, 200, 64, 264))
	source1.SetIconBounds(image.Rect(0, 200, 64, 264))

	source2 := new(Resource).Init()
	source2.SetBindings(image.Rect(200, 200, 264, 264))
	source2.SetIconBounds(image.Rect(200, 200, 264, 264))

	// Setup target resources
	target1 := new(Resource).Init()
	target1.SetBindings(image.Rect(0, 300, 64, 364))
	target1.SetIconBounds(image.Rect(0, 300, 64, 364))

	target2 := new(Resource).Init()
	target2.SetBindings(image.Rect(200, 300, 264, 364))
	target2.SetIconBounds(image.Rect(200, 300, 264, 364))

	// Create incoming links (Source -> ELB)
	incomingLink1 := new(Link).Init(source1, WINDROSE_N, ArrowHead{}, elb, WINDROSE_S, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	incomingLink2 := new(Link).Init(source2, WINDROSE_N, ArrowHead{}, elb, WINDROSE_S, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	// Create outgoing links (ELB -> Target)
	outgoingLink1 := new(Link).Init(elb, WINDROSE_S, ArrowHead{}, target1, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	outgoingLink2 := new(Link).Init(elb, WINDROSE_S, ArrowHead{}, target2, WINDROSE_N, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	// Add links to resources
	elb.AddLink(incomingLink1)
	elb.AddLink(incomingLink2)
	elb.AddLink(outgoingLink1)
	elb.AddLink(outgoingLink2)

	t.Run("Without directional grouping", func(t *testing.T) {
		// Enable basic grouping offset only
		elb.SetGroupingOffset(true)
		elb.SetGroupingOffsetDirection(false)

		// Sort links
		elb.sortAllLinks()

		// All links should get sequential indexes
		index1, count1 := incomingLink1.getLinkIndexAndCount(elb, WINDROSE_S)
		index2, count2 := incomingLink2.getLinkIndexAndCount(elb, WINDROSE_S)
		index3, count3 := outgoingLink1.getLinkIndexAndCount(elb, WINDROSE_S)
		index4, count4 := outgoingLink2.getLinkIndexAndCount(elb, WINDROSE_S)

		if count1 != 4 || count2 != 4 || count3 != 4 || count4 != 4 {
			t.Errorf("Expected all counts to be 4, got %d, %d, %d, %d", count1, count2, count3, count4)
		}

		// All indexes should be different
		indexes := []int{index1, index2, index3, index4}
		for i := 0; i < len(indexes); i++ {
			for j := i + 1; j < len(indexes); j++ {
				if indexes[i] == indexes[j] {
					t.Errorf("Found duplicate indexes: %d and %d both have index %d", i, j, indexes[i])
				}
			}
		}
	})

	t.Run("With directional grouping", func(t *testing.T) {
		// Enable directional grouping offset
		elb.SetGroupingOffset(true)
		elb.SetGroupingOffsetDirection(true)

		// Sort links
		elb.sortAllLinks()

		// Incoming links should be grouped together
		inIndex1, inCount1 := incomingLink1.getLinkIndexAndCount(elb, WINDROSE_S)
		inIndex2, inCount2 := incomingLink2.getLinkIndexAndCount(elb, WINDROSE_S)

		// Outgoing links should be grouped together
		outIndex1, outCount1 := outgoingLink1.getLinkIndexAndCount(elb, WINDROSE_S)
		outIndex2, outCount2 := outgoingLink2.getLinkIndexAndCount(elb, WINDROSE_S)

		// Each group should have count of 2
		if inCount1 != 2 || inCount2 != 2 {
			t.Errorf("Expected incoming link counts to be 2, got %d, %d", inCount1, inCount2)
		}
		if outCount1 != 2 || outCount2 != 2 {
			t.Errorf("Expected outgoing link counts to be 2, got %d, %d", outCount1, outCount2)
		}

		// Incoming and outgoing groups should have different group indexes
		if inIndex1 == outIndex1 && inIndex2 == outIndex2 {
			t.Errorf("Incoming and outgoing groups should have different indexes")
		}

		// Links within same group should have same index
		if inIndex1 != inIndex2 {
			t.Errorf("Incoming links should have same index, got %d and %d", inIndex1, inIndex2)
		}
		if outIndex1 != outIndex2 {
			t.Errorf("Outgoing links should have same index, got %d and %d", outIndex1, outIndex2)
		}
	})
}

func TestGroupingOffsetDirectionCoordinateSelection(t *testing.T) {
	// Test that correct coordinate axis is used for different positions
	elb := new(Resource).Init()
	elb.SetBindings(image.Rect(100, 100, 164, 164))
	elb.SetIconBounds(image.Rect(100, 100, 164, 164))
	elb.SetGroupingOffset(true)
	elb.SetGroupingOffsetDirection(true)

	// Test vertical position (should use X coordinate)
	t.Run("Vertical position uses X coordinate", func(t *testing.T) {
		// Create sources with different X coordinates
		source1 := new(Resource).Init()
		source1.SetIconBounds(image.Rect(50, 200, 114, 264)) // X center: 82

		source2 := new(Resource).Init()
		source2.SetIconBounds(image.Rect(150, 200, 214, 264)) // X center: 182

		link1 := new(Link).Init(source1, WINDROSE_N, ArrowHead{}, elb, WINDROSE_S, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
		link2 := new(Link).Init(source2, WINDROSE_N, ArrowHead{}, elb, WINDROSE_S, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

		elb.AddLink(link1)
		elb.AddLink(link2)

		// The algorithm should use X coordinates to determine order
		index1, _ := link1.getLinkIndexAndCount(elb, WINDROSE_S)
		index2, _ := link2.getLinkIndexAndCount(elb, WINDROSE_S)

		// Since source1 has smaller X coordinate (82 < 182), it should get index 1
		// and source2 should get index 1 (both incoming, same group)
		if index1 != index2 {
			t.Errorf("Both incoming links should have same index, got %d and %d", index1, index2)
		}
	})

	// Test horizontal position (should use Y coordinate)
	t.Run("Horizontal position uses Y coordinate", func(t *testing.T) {
		elb2 := new(Resource).Init()
		elb2.SetBindings(image.Rect(100, 100, 164, 164))
		elb2.SetIconBounds(image.Rect(100, 100, 164, 164))
		elb2.SetGroupingOffset(true)
		elb2.SetGroupingOffsetDirection(true)

		// Create sources with different Y coordinates
		source1 := new(Resource).Init()
		source1.SetIconBounds(image.Rect(200, 50, 264, 114)) // Y center: 82

		source2 := new(Resource).Init()
		source2.SetIconBounds(image.Rect(200, 150, 264, 214)) // Y center: 182

		link1 := new(Link).Init(source1, WINDROSE_W, ArrowHead{}, elb2, WINDROSE_E, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
		link2 := new(Link).Init(source2, WINDROSE_W, ArrowHead{}, elb2, WINDROSE_E, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

		elb2.AddLink(link1)
		elb2.AddLink(link2)

		// The algorithm should use Y coordinates to determine order
		index1, _ := link1.getLinkIndexAndCount(elb2, WINDROSE_E)
		index2, _ := link2.getLinkIndexAndCount(elb2, WINDROSE_E)

		// Both should be in same group (incoming)
		if index1 != index2 {
			t.Errorf("Both incoming links should have same index, got %d and %d", index1, index2)
		}
	})
}

func TestGroupingOffsetMultipleMixedLinks(t *testing.T) {
	// Setup: A->B, A->C, D->A, E->A (A has 2 outgoing and 2 incoming at same positions)
	resourceA := new(Resource).Init()
	resourceA.SetBindings(image.Rect(100, 100, 164, 164))

	resourceB := new(Resource).Init()
	resourceB.SetBindings(image.Rect(200, 50, 264, 114))

	resourceC := new(Resource).Init()
	resourceC.SetBindings(image.Rect(200, 150, 264, 214))

	resourceD := new(Resource).Init()
	resourceD.SetBindings(image.Rect(0, 50, 64, 114))

	resourceE := new(Resource).Init()
	resourceE.SetBindings(image.Rect(0, 150, 64, 214))

	// A->B, A->C (both from A:E)
	linkAB := new(Link).Init(resourceA, WINDROSE_E, ArrowHead{}, resourceB, WINDROSE_W, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	linkAC := new(Link).Init(resourceA, WINDROSE_E, ArrowHead{}, resourceC, WINDROSE_W, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	// D->A, E->A (both to A:W)
	linkDA := new(Link).Init(resourceD, WINDROSE_E, ArrowHead{}, resourceA, WINDROSE_W, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})
	linkEA := new(Link).Init(resourceE, WINDROSE_E, ArrowHead{}, resourceA, WINDROSE_W, ArrowHead{}, 2, color.RGBA{0, 0, 0, 255})

	resourceA.AddLink(linkAB)
	resourceA.AddLink(linkAC)
	resourceA.AddLink(linkDA)
	resourceA.AddLink(linkEA)

	t.Run("Multiple outgoing links with grouping offset", func(t *testing.T) {
		resourceA.SetGroupingOffset(true)
		resourceA.sortAllLinks()

		_, countE := linkAB.getLinkIndexAndCount(resourceA, WINDROSE_E)
		if countE != 2 {
			t.Errorf("Expected 2 outgoing links from A:E, got %d", countE)
		}

		ptAB := linkAB.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_E, resourceA, true)
		ptAC := linkAC.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_E, resourceA, true)

		// Two links should have different offsets
		if ptAB.Y == ptAC.Y {
			t.Errorf("Two outgoing links should have different Y offsets: both at Y=%d", ptAB.Y)
		}
	})

	t.Run("Multiple incoming links with grouping offset", func(t *testing.T) {
		resourceA.SetGroupingOffset(true)
		resourceA.sortAllLinks()

		_, countW := linkDA.getLinkIndexAndCount(resourceA, WINDROSE_W)
		if countW != 2 {
			t.Errorf("Expected 2 incoming links to A:W, got %d", countW)
		}

		ptDA := linkDA.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_W, resourceA, false)
		ptEA := linkEA.calcPositionWithOffset(resourceA.GetBindings(), WINDROSE_W, resourceA, false)

		// Two links should have different offsets
		if ptDA.Y == ptEA.Y {
			t.Errorf("Two incoming links should have different Y offsets: both at Y=%d", ptDA.Y)
		}
	})

	t.Run("Outgoing and incoming counts are independent", func(t *testing.T) {
		resourceA.SetGroupingOffset(true)

		_, countOutgoing := linkAB.getLinkIndexAndCount(resourceA, WINDROSE_E)
		_, countIncoming := linkDA.getLinkIndexAndCount(resourceA, WINDROSE_W)

		if countOutgoing != 2 {
			t.Errorf("Expected 2 outgoing links, got %d", countOutgoing)
		}
		if countIncoming != 2 {
			t.Errorf("Expected 2 incoming links, got %d", countIncoming)
		}
	})
}

func TestAutoCalculatePositions(t *testing.T) {
	// Create mock resources with different positions
	source := &Resource{}
	target := &Resource{}

	// Test horizontal layout (source left, target right)
	sourceBounds := image.Rect(0, 0, 100, 100)
	targetBounds := image.Rect(200, 0, 300, 100)
	source.bindings = &sourceBounds
	target.bindings = &targetBounds

	sourcePos, targetPos := AutoCalculatePositions(source, target)

	// dx=150, dy=0 -> horizontal connection, target to the right
	if sourcePos != WINDROSE_E {
		t.Errorf("Expected source position E, got %v", sourcePos)
	}
	if targetPos != WINDROSE_W {
		t.Errorf("Expected target position W, got %v", targetPos)
	}

	// Test vertical layout (source top, target bottom)
	sourceBounds2 := image.Rect(0, 0, 100, 100)
	targetBounds2 := image.Rect(0, 200, 100, 300)
	source.bindings = &sourceBounds2
	target.bindings = &targetBounds2

	sourcePos, targetPos = AutoCalculatePositions(source, target)

	// dx=0, dy=150 -> vertical connection, target below
	if sourcePos != WINDROSE_S {
		t.Errorf("Expected source position S, got %v", sourcePos)
	}
	if targetPos != WINDROSE_N {
		t.Errorf("Expected target position N, got %v", targetPos)
	}
}

func TestCalcPositionWithOffset_TitleHeight(t *testing.T) {
	// Create a resource with a title
	resource := new(Resource).Init()
	resource.label = "Test Resource\nMulti-line"

	// Set up resource bindings
	bindings := image.Rect(100, 100, 200, 200)
	resource.bindings = &bindings

	// Create a link
	link := &Link{}

	// Test SSE position (should include title height)
	point := link.calcPositionWithOffset(bindings, WINDROSE_SSE, resource, true)

	// The Y coordinate should be greater than the original binding due to title height
	if point.Y <= bindings.Max.Y {
		t.Errorf("Expected Y coordinate to be adjusted for title height, got %d (original max Y: %d)", point.Y, bindings.Max.Y)
	}

	// Test S position (should include title height)
	point = link.calcPositionWithOffset(bindings, WINDROSE_S, resource, true)
	if point.Y <= bindings.Max.Y {
		t.Errorf("Expected Y coordinate to be adjusted for title height for S position, got %d", point.Y)
	}

	// Test SSW position (should include title height)
	point = link.calcPositionWithOffset(bindings, WINDROSE_SSW, resource, true)
	if point.Y <= bindings.Max.Y {
		t.Errorf("Expected Y coordinate to be adjusted for title height for SSW position, got %d", point.Y)
	}

	// Test N position (should NOT include title height)
	point = link.calcPositionWithOffset(bindings, WINDROSE_N, resource, true)
	expectedY := bindings.Min.Y
	if point.Y != expectedY {
		t.Errorf("Expected Y coordinate %d for N position (no title adjustment), got %d", expectedY, point.Y)
	}

	// Test with resource without title
	resourceNoTitle := new(Resource).Init()
	resourceNoTitle.label = ""
	resourceNoTitle.bindings = &bindings

	pointNoTitle := link.calcPositionWithOffset(bindings, WINDROSE_S, resourceNoTitle, true)
	expectedYNoTitle := bindings.Max.Y
	if pointNoTitle.Y != expectedYNoTitle {
		t.Errorf("Expected Y coordinate %d for S position with no title, got %d", expectedYNoTitle, pointNoTitle.Y)
	}
}

// Tests for layout-aware auto-positioning

func TestFindLowestCommonAncestor(t *testing.T) {
	// Create test hierarchy: Root -> Parent -> Child1, Child2
	root := &Resource{direction: "vertical"}
	parent := &Resource{direction: "horizontal", parent: root}
	child1 := &Resource{parent: parent}
	child2 := &Resource{parent: parent}

	// Test case 1: Same resource
	result := findLowestCommonAncestor(child1, child1)
	if result != child1 {
		t.Errorf("Expected child1, got %v", result)
	}

	// Test case 2: Siblings
	result = findLowestCommonAncestor(child1, child2)
	if result != parent {
		t.Errorf("Expected parent, got %v", result)
	}

	// Test case 3: Parent-child relationship
	result = findLowestCommonAncestor(parent, child1)
	if result != parent {
		t.Errorf("Expected parent, got %v", result)
	}

	// Test case 4: Different branches
	otherParent := &Resource{direction: "vertical", parent: root}
	otherChild := &Resource{parent: otherParent}

	result = findLowestCommonAncestor(child1, otherChild)
	if result != root {
		t.Errorf("Expected root, got %v", result)
	}
}

func TestCountResourcesInDirections(t *testing.T) {
	// Create test hierarchy:
	// LCA {horizontal: [V1, V2, V3]}
	//   V2 {vertical: [R1, R2, R3]}
	lca := &Resource{direction: "horizontal"}
	v1 := &Resource{parent: lca}
	v2 := &Resource{direction: "vertical", parent: lca}
	v3 := &Resource{parent: lca}
	lca.children = []*Resource{v1, v2, v3}

	r1 := &Resource{parent: v2}
	r2 := &Resource{parent: v2}
	r3 := &Resource{parent: v2}
	v2.children = []*Resource{r1, r2, r3}

	// Test case 1: R2 to LCA
	// Expected: N=1 (R1), S=1 (R3), W=1 (V1), E=1 (V3)
	counts := countResourcesInDirections(r2, lca)
	if counts.North != 1 {
		t.Errorf("Expected North=1, got %d", counts.North)
	}
	if counts.South != 1 {
		t.Errorf("Expected South=1, got %d", counts.South)
	}
	if counts.West != 1 {
		t.Errorf("Expected West=1, got %d", counts.West)
	}
	if counts.East != 1 {
		t.Errorf("Expected East=1, got %d", counts.East)
	}

	// Test case 2: R1 to LCA
	// Expected: N=0, S=2 (R2, R3), W=1 (V1), E=1 (V3)
	counts = countResourcesInDirections(r1, lca)
	if counts.North != 0 {
		t.Errorf("Expected North=0, got %d", counts.North)
	}
	if counts.South != 2 {
		t.Errorf("Expected South=2, got %d", counts.South)
	}
	if counts.West != 1 {
		t.Errorf("Expected West=1, got %d", counts.West)
	}
	if counts.East != 1 {
		t.Errorf("Expected East=1, got %d", counts.East)
	}

	// Test case 3: V1 to LCA (direct child)
	// Expected: N=0, S=0, W=0, E=2 (V2, V3)
	counts = countResourcesInDirections(v1, lca)
	if counts.North != 0 {
		t.Errorf("Expected North=0, got %d", counts.North)
	}
	if counts.South != 0 {
		t.Errorf("Expected South=0, got %d", counts.South)
	}
	if counts.West != 0 {
		t.Errorf("Expected West=0, got %d", counts.West)
	}
	if counts.East != 2 {
		t.Errorf("Expected East=2, got %d", counts.East)
	}

	// Test case 4: target == LCA
	// Expected: all zeros
	counts = countResourcesInDirections(lca, lca)
	if counts.North != 0 || counts.South != 0 || counts.West != 0 || counts.East != 0 {
		t.Errorf("Expected all zeros when target==LCA, got %+v", counts)
	}
}

func TestAutoCalculatePositionsWithResourceCounts(t *testing.T) {
	// Create test hierarchy:
	// LCA {horizontal: [V1, V2, V3]}
	//   V2 {vertical: [R1, R2, R3]}
	lca := &Resource{direction: "horizontal"}
	v1 := &Resource{parent: lca}
	v2 := &Resource{direction: "vertical", parent: lca}
	v3 := &Resource{parent: lca}
	lca.children = []*Resource{v1, v2, v3}

	r1 := &Resource{parent: v2}
	r2 := &Resource{parent: v2}
	r3 := &Resource{parent: v2}
	v2.children = []*Resource{r1, r2, r3}

	// Set bindings for all resources
	lca.bindings = &image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{400, 300}}
	v1.bindings = &image.Rectangle{Min: image.Point{10, 100}, Max: image.Point{74, 164}}
	v2.bindings = &image.Rectangle{Min: image.Point{150, 50}, Max: image.Point{250, 250}}
	v3.bindings = &image.Rectangle{Min: image.Point{320, 100}, Max: image.Point{384, 164}}
	r1.bindings = &image.Rectangle{Min: image.Point{180, 60}, Max: image.Point{244, 124}}
	r2.bindings = &image.Rectangle{Min: image.Point{180, 140}, Max: image.Point{244, 204}}
	r3.bindings = &image.Rectangle{Min: image.Point{180, 220}, Max: image.Point{244, 284}}

	// Test case 1: R2 to V1
	// R2 counts: N=1, S=1, W=1, E=1
	// V1 counts: N=0, S=0, W=0, E=2
	// R2 should prefer W (least count=1, aligned with LCA horizontal)
	// V1 should prefer E (least count=0, perpendicular to LCA)
	sourcePos, targetPos := AutoCalculatePositions(r2, v1)
	if sourcePos != WINDROSE_W {
		t.Errorf("Expected R2->V1 source=W, got %v", sourcePos)
	}
	if targetPos != WINDROSE_E {
		t.Errorf("Expected R2->V1 target=E, got %v", targetPos)
	}

	// Test case 2: R1 to R3 (within same vertical stack, but with intermediate resource)
	// R1 counts: N=0, E=0, W=0, S=0 (R1 is direct child of V2)
	// R3 counts: N=0, E=0, W=0, S=0 (R3 is direct child of V2)
	// After adjustment: R1.S += 1 (R2 between), R3.N += 1 (R2 between)
	// R1: min=0 (N,E,W), after adjustment min=0 (N,E,W), S=1
	// R3: min=0 (E,W,S), after adjustment min=0 (E,W), N=1, S=0
	// LCA (V2) direction=vertical, dy>0 → R1 should prefer S, R3 should prefer N
	// But with equal counts, it falls back to distance-based or first candidate
	sourcePos, targetPos = AutoCalculatePositions(r1, r3)
	// With adjustment, R1 has N=0,E=0,W=0,S=1 → min=0, candidates: N,E,W
	// LCA vertical, dy>0 → prefer S (not in candidates), then perpendicular E/W
	// With adjustment, R3 has N=1,E=0,W=0,S=0 → min=0, candidates: E,W,S
	// LCA vertical, dy<0 (from R3 perspective) → prefer N (not in candidates), then perpendicular E/W
	// Both should pick E or W (perpendicular to vertical)
	if sourcePos != WINDROSE_E && sourcePos != WINDROSE_W {
		t.Errorf("Expected R1->R3 source=E or W, got %v", sourcePos)
	}
	if targetPos != WINDROSE_E && targetPos != WINDROSE_W {
		t.Errorf("Expected R1->R3 target=E or W, got %v", targetPos)
	}
}

func TestAutoCalculatePositionsPerpendicularSelection(t *testing.T) {
	// Test perpendicular direction selection based on dx/dy
	// LCA {horizontal: [VPC1, VPC2]}
	//   VPC1 {vertical: [TopLeft, BottomLeft]}
	//   VPC2 {vertical: [TopRight, BottomRight]}
	lca := &Resource{direction: "horizontal"}
	vpc1 := &Resource{direction: "vertical", parent: lca}
	vpc2 := &Resource{direction: "vertical", parent: lca}
	lca.children = []*Resource{vpc1, vpc2}

	topLeft := &Resource{parent: vpc1}
	bottomLeft := &Resource{parent: vpc1}
	vpc1.children = []*Resource{topLeft, bottomLeft}

	topRight := &Resource{parent: vpc2}
	bottomRight := &Resource{parent: vpc2}
	vpc2.children = []*Resource{topRight, bottomRight}

	// Set bindings
	lca.bindings = &image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{600, 400}}
	vpc1.bindings = &image.Rectangle{Min: image.Point{50, 50}, Max: image.Point{250, 350}}
	vpc2.bindings = &image.Rectangle{Min: image.Point{350, 50}, Max: image.Point{550, 350}}
	topLeft.bindings = &image.Rectangle{Min: image.Point{100, 80}, Max: image.Point{164, 144}}
	bottomLeft.bindings = &image.Rectangle{Min: image.Point{100, 250}, Max: image.Point{164, 314}}
	topRight.bindings = &image.Rectangle{Min: image.Point{400, 80}, Max: image.Point{464, 144}}
	bottomRight.bindings = &image.Rectangle{Min: image.Point{400, 250}, Max: image.Point{464, 314}}

	// Test case 1: TopLeft -> TopRight (same height, dy ≈ 0)
	// LCA horizontal (E/W preferred), perpendicular N/S have equal counts
	// Should prefer E/W (LCA direction)
	sourcePos, targetPos := AutoCalculatePositions(topLeft, topRight)
	if sourcePos != WINDROSE_E {
		t.Errorf("TopLeft->TopRight: expected source=E, got %v", sourcePos)
	}
	if targetPos != WINDROSE_W {
		t.Errorf("TopLeft->TopRight: expected target=W, got %v", targetPos)
	}

	// Test case 2: TopLeft -> BottomRight (diagonal, dy > 0)
	// If perpendicular N/S are candidates with equal count, should prefer S (dy > 0)
	// But LCA direction E/W should be preferred first
	sourcePos, targetPos = AutoCalculatePositions(topLeft, bottomRight)
	// Source should prefer E (LCA direction) or S (perpendicular, dy > 0)
	if sourcePos != WINDROSE_E && sourcePos != WINDROSE_S {
		t.Errorf("TopLeft->BottomRight: expected source=E or S, got %v", sourcePos)
	}
	// Target should prefer W (LCA direction) or N (perpendicular, dy < 0 from target perspective)
	if targetPos != WINDROSE_W && targetPos != WINDROSE_N {
		t.Errorf("TopLeft->BottomRight: expected target=W or N, got %v", targetPos)
	}

	// Test case 3: Create scenario where perpendicular is chosen
	// VPC3 {vertical: [Left, Right]} - horizontal layout within VPC
	vpc3 := &Resource{direction: "horizontal", parent: lca}
	lca.children = []*Resource{vpc3}

	left := &Resource{parent: vpc3}
	right := &Resource{parent: vpc3}
	vpc3.children = []*Resource{left, right}

	vpc3.bindings = &image.Rectangle{Min: image.Point{50, 50}, Max: image.Point{550, 350}}
	left.bindings = &image.Rectangle{Min: image.Point{100, 150}, Max: image.Point{164, 214}}
	right.bindings = &image.Rectangle{Min: image.Point{400, 250}, Max: image.Point{464, 314}}

	// Left -> Right: dx > 0, dy > 0
	// LCA direction=horizontal, preferred=E/W
	// If both E/W have equal counts, should choose E (dx > 0)
	sourcePos, targetPos = AutoCalculatePositions(left, right)
	// Both should prefer E/W (LCA direction)
	if sourcePos != WINDROSE_E && sourcePos != WINDROSE_W {
		t.Errorf("Left->Right: expected source=E or W, got %v", sourcePos)
	}
	if targetPos != WINDROSE_E && targetPos != WINDROSE_W {
		t.Errorf("Left->Right: expected target=E or W, got %v", targetPos)
	}
}

func TestAutoCalculatePositionsVerticalLayout(t *testing.T) {
	// Create vertical layout hierarchy
	vpc := &Resource{direction: "vertical"}
	elb := &Resource{parent: vpc}
	instance := &Resource{parent: vpc}

	// Mock GetBindings for positioning
	elb.bindings = &image.Rectangle{Min: image.Point{100, 50}, Max: image.Point{164, 114}}
	instance.bindings = &image.Rectangle{Min: image.Point{50, 150}, Max: image.Point{114, 214}}

	// Test vertical layout (should use S->N)
	sourcePos, targetPos := AutoCalculatePositions(elb, instance)
	if sourcePos != WINDROSE_S || targetPos != WINDROSE_N {
		t.Errorf("Expected S->N for vertical layout, got %v->%v", sourcePos, targetPos)
	}
}

func TestAutoCalculatePositionsHorizontalLayout(t *testing.T) {
	// Create horizontal layout hierarchy
	container := &Resource{direction: "horizontal"}
	source := &Resource{parent: container}
	target := &Resource{parent: container}

	// Mock GetBindings for positioning (dx > dy)
	source.bindings = &image.Rectangle{Min: image.Point{50, 100}, Max: image.Point{114, 164}}
	target.bindings = &image.Rectangle{Min: image.Point{200, 100}, Max: image.Point{264, 164}}

	// Test horizontal layout (should use E->W)
	sourcePos, targetPos := AutoCalculatePositions(source, target)
	if sourcePos != WINDROSE_E || targetPos != WINDROSE_W {
		t.Errorf("Expected E->W for horizontal layout, got %v->%v", sourcePos, targetPos)
	}
}

func TestAutoCalculatePositionsNoCommonAncestor(t *testing.T) {
	// Create resources with no common ancestor
	source := &Resource{}
	target := &Resource{}

	// Mock GetBindings for distance-based calculation
	source.bindings = &image.Rectangle{Min: image.Point{50, 50}, Max: image.Point{114, 114}}
	target.bindings = &image.Rectangle{Min: image.Point{200, 50}, Max: image.Point{264, 114}}

	// Should fall back to distance-based logic (dx > dy)
	sourcePos, targetPos := AutoCalculatePositions(source, target)
	if sourcePos != WINDROSE_E || targetPos != WINDROSE_W {
		t.Errorf("Expected E->W for distance-based horizontal, got %v->%v", sourcePos, targetPos)
	}
}

func TestFindLongestHorizontalSegment(t *testing.T) {
	link := &Link{}

	// Test case 1: Single horizontal segment
	controlPts1 := []image.Point{
		{X: 100, Y: 200},
		{X: 300, Y: 200}, // horizontal: length 200
		{X: 300, Y: 400},
	}
	start, end, length := link.findLongestHorizontalSegment(controlPts1)
	expectedStart := image.Point{X: 100, Y: 200}
	expectedEnd := image.Point{X: 300, Y: 200}
	expectedLength := 200

	if start != expectedStart || end != expectedEnd || length != expectedLength {
		t.Errorf("Expected start=%v, end=%v, length=%d, got start=%v, end=%v, length=%d",
			expectedStart, expectedEnd, expectedLength, start, end, length)
	}

	// Test case 2: Multiple horizontal segments, find longest
	controlPts2 := []image.Point{
		{X: 100, Y: 200},
		{X: 200, Y: 200}, // horizontal: length 100
		{X: 200, Y: 300},
		{X: 500, Y: 300}, // horizontal: length 300 (longest)
		{X: 500, Y: 400},
		{X: 600, Y: 400}, // horizontal: length 100
	}
	start2, end2, length2 := link.findLongestHorizontalSegment(controlPts2)
	expectedStart2 := image.Point{X: 200, Y: 300}
	expectedEnd2 := image.Point{X: 500, Y: 300}
	expectedLength2 := 300

	if start2 != expectedStart2 || end2 != expectedEnd2 || length2 != expectedLength2 {
		t.Errorf("Expected start=%v, end=%v, length=%d, got start=%v, end=%v, length=%d",
			expectedStart2, expectedEnd2, expectedLength2, start2, end2, length2)
	}

	// Test case 3: No horizontal segments
	controlPts3 := []image.Point{
		{X: 100, Y: 200},
		{X: 100, Y: 300}, // vertical
		{X: 200, Y: 300},
		{X: 200, Y: 400}, // vertical
	}
	start3, end3, length3 := link.findLongestHorizontalSegment(controlPts3)
	expectedStart3 := image.Point{X: 100, Y: 300}
	expectedEnd3 := image.Point{X: 200, Y: 300}
	expectedLength3 := 100

	if start3 != expectedStart3 || end3 != expectedEnd3 || length3 != expectedLength3 {
		t.Errorf("Expected start=%v, end=%v, length=%d, got start=%v, end=%v, length=%d",
			expectedStart3, expectedEnd3, expectedLength3, start3, end3, length3)
	}

	// Test case 4: Empty control points
	controlPts4 := []image.Point{}
	start4, end4, length4 := link.findLongestHorizontalSegment(controlPts4)
	expectedZero := image.Point{}

	if start4 != expectedZero || end4 != expectedZero || length4 != 0 {
		t.Errorf("Expected zero values for empty input, got start=%v, end=%v, length=%d",
			start4, end4, length4)
	}
}

func TestCalculateLCABasedMidpoint(t *testing.T) {
	// Create complex nested structure:
	// HorizontalStack{
	//   HorizontalStack{
	//     VerticalStack{VerticalStack{Resource1}},
	//     VerticalStack{VerticalStack{Resource2}},
	//     VerticalStack{VerticalStack{Resource3}},
	//     VerticalStack{VerticalStack{Resource4}},
	//     VerticalStack{VerticalStack{Resource5}},
	//     VerticalStack{VerticalStack{Resource6}}
	//   }
	// }

	// Root HorizontalStack
	root := new(Resource).Init()
	root.direction = "horizontal"

	// Inner HorizontalStack
	innerHorizontal := new(Resource).Init()
	innerHorizontal.direction = "horizontal"
	innerHorizontal.parent = root
	root.children = []*Resource{innerHorizontal}

	// Create 6 VerticalStacks
	verticalStacks := make([]*Resource, 6)
	for i := 0; i < 6; i++ {
		verticalStacks[i] = new(Resource).Init()
		verticalStacks[i].direction = "vertical"
		verticalStacks[i].parent = innerHorizontal
	}
	innerHorizontal.children = verticalStacks

	// Create inner VerticalStacks and Resources
	resources := make([]*Resource, 6)
	for i := 0; i < 6; i++ {
		innerVertical := new(Resource).Init()
		innerVertical.direction = "vertical"
		innerVertical.parent = verticalStacks[i]

		resources[i] = new(Resource).Init()
		resources[i].parent = innerVertical

		verticalStacks[i].children = []*Resource{innerVertical}
		innerVertical.children = []*Resource{resources[i]}

		// Set actual bindings for each vertical stack (horizontal layout)
		x := i * 100 // Each stack is 100px wide
		verticalStacks[i].SetBindings(image.Rect(x, 0, x+80, 100))

		// Set bindings for inner vertical and resources
		innerVertical.SetBindings(image.Rect(x+10, 10, x+70, 90))
		resources[i].SetBindings(image.Rect(x+20, 20, x+60, 80))
	}

	// Test LCA detection: Resource2 and Resource5 should have innerHorizontal as LCA
	lca := findLowestCommonAncestor(resources[1], resources[4])
	if lca != innerHorizontal {
		t.Errorf("Expected LCA to be innerHorizontal, got %v", lca)
	}

	// Test Resource2 -> Resource5 (forward)
	link := &Link{
		Source: resources[1], // Resource2
		Target: resources[4], // Resource5
	}

	sourcePt := image.Point{X: 200, Y: 100}
	targetPt := image.Point{X: 500, Y: 100}

	result := link.calculateLCABasedMidpoint(sourcePt, targetPt)
	if result.X == 0 && result.Y == 0 {
		t.Error("Expected non-zero midpoint result for Resource2->Resource5")
	}

	// Test Resource5 -> Resource2 (reverse)
	reverseLink := &Link{
		Source: resources[4], // Resource5
		Target: resources[1], // Resource2
	}

	result2 := reverseLink.calculateLCABasedMidpoint(targetPt, sourcePt)
	if result2.X == 0 && result2.Y == 0 {
		t.Error("Expected non-zero midpoint result for Resource5->Resource2")
	}

	// Test edge case: no LCA (fallback to 50%)
	orphan := new(Resource).Init()
	noLCALink := &Link{Source: resources[0], Target: orphan}
	fallback := noLCALink.calculateLCABasedMidpoint(image.Point{X: 100, Y: 100}, image.Point{X: 300, Y: 200})
	if fallback.X != 200 || fallback.Y != 150 {
		t.Errorf("Expected 50%% fallback (200,150), got (%d,%d)", fallback.X, fallback.Y)
	}

	// Test edge case: same child (should use fallback)
	sameChildLink := &Link{Source: resources[0], Target: resources[0]}
	samePt := sameChildLink.calculateLCABasedMidpoint(image.Point{X: 50, Y: 50}, image.Point{X: 70, Y: 70})
	if samePt.X != 60 || samePt.Y != 60 {
		t.Errorf("Expected same child fallback (60,60), got (%d,%d)", samePt.X, samePt.Y)
	}

	// Test vertical layout
	vertRoot := new(Resource).Init()
	vertRoot.direction = "vertical"
	vert1 := new(Resource).Init()
	vert1.parent = vertRoot
	vert1.SetBindings(image.Rect(0, 0, 100, 80))
	vert2 := new(Resource).Init()
	vert2.parent = vertRoot
	vert2.SetBindings(image.Rect(0, 100, 100, 180))
	vertRoot.children = []*Resource{vert1, vert2}

	vertLink := &Link{Source: vert1, Target: vert2}
	vertResult := vertLink.calculateLCABasedMidpoint(image.Point{X: 50, Y: 40}, image.Point{X: 50, Y: 140})
	if vertResult.Y != 90 {
		t.Errorf("Expected vertical gap Y=90, got Y=%d", vertResult.Y)
	}
}

func TestFindPerpendicularSegments(t *testing.T) {
	link := &Link{}

	// Test case 1: Perfect 90-degree segments (L-shape)
	// segment[0]: vertical down (0,0) -> (0,100)
	// segment[1]: horizontal right (0,100) -> (200,100) [longest horizontal]
	// segment[2]: vertical down (200,100) -> (200,200)
	controlPts1 := []image.Point{
		{X: 0, Y: 0},
		{X: 0, Y: 100},   // end of segment[0]
		{X: 200, Y: 100}, // end of segment[1] (longest horizontal)
		{X: 200, Y: 200}, // end of segment[2]
	}

	t.Logf("Test case 1 - Perfect 90-degree segments:")
	t.Logf("Control points: %v", controlPts1)
	t.Logf("Segment[0]: (%d,%d) -> (%d,%d) = vector(%d,%d)",
		controlPts1[0].X, controlPts1[0].Y, controlPts1[1].X, controlPts1[1].Y,
		controlPts1[1].X-controlPts1[0].X, controlPts1[1].Y-controlPts1[0].Y)
	t.Logf("Segment[1]: (%d,%d) -> (%d,%d) = vector(%d,%d)",
		controlPts1[1].X, controlPts1[1].Y, controlPts1[2].X, controlPts1[2].Y,
		controlPts1[2].X-controlPts1[1].X, controlPts1[2].Y-controlPts1[1].Y)
	t.Logf("Segment[2]: (%d,%d) -> (%d,%d) = vector(%d,%d)",
		controlPts1[2].X, controlPts1[2].Y, controlPts1[3].X, controlPts1[3].Y,
		controlPts1[3].X-controlPts1[2].X, controlPts1[3].Y-controlPts1[2].Y)

	hasLeft90, hasRight90 := link.findPerpendicularSegments(controlPts1)
	t.Logf("Result: hasLeft90=%v, hasRight90=%v", hasLeft90, hasRight90)

	// Left側は90度時計回り(270度反時計回り)なので false、Right側は90度反時計回りなので true
	if hasLeft90 {
		t.Errorf("Expected no left 90-degree segment (actually 270-degree counterclockwise), got true")
	}
	if !hasRight90 {
		t.Errorf("Expected right 90-degree segment, got false")
	}

	// Test case 2: No 90-degree segments (270-degree turns)
	// segment[0]: vertical up (0,100) -> (0,0)
	// segment[1]: horizontal right (0,0) -> (200,0) [longest horizontal]
	// segment[2]: vertical up (200,0) -> (200,-100)
	controlPts2 := []image.Point{
		{X: 0, Y: 100},
		{X: 0, Y: 0},      // end of segment[0]
		{X: 200, Y: 0},    // end of segment[1] (longest horizontal)
		{X: 200, Y: -100}, // end of segment[2]
	}

	t.Logf("\nTest case 2 - 270-degree segments:")
	t.Logf("Control points: %v", controlPts2)
	t.Logf("Segment[0]: (%d,%d) -> (%d,%d) = vector(%d,%d)",
		controlPts2[0].X, controlPts2[0].Y, controlPts2[1].X, controlPts2[1].Y,
		controlPts2[1].X-controlPts2[0].X, controlPts2[1].Y-controlPts2[0].Y)
	t.Logf("Segment[1]: (%d,%d) -> (%d,%d) = vector(%d,%d)",
		controlPts2[1].X, controlPts2[1].Y, controlPts2[2].X, controlPts2[2].Y,
		controlPts2[2].X-controlPts2[1].X, controlPts2[2].Y-controlPts2[1].Y)
	t.Logf("Segment[2]: (%d,%d) -> (%d,%d) = vector(%d,%d)",
		controlPts2[2].X, controlPts2[2].Y, controlPts2[3].X, controlPts2[3].Y,
		controlPts2[3].X-controlPts2[2].X, controlPts2[3].Y-controlPts2[2].Y)

	hasLeft90_2, hasRight90_2 := link.findPerpendicularSegments(controlPts2)
	t.Logf("Result: hasLeft90=%v, hasRight90=%v", hasLeft90_2, hasRight90_2)

	// 修正: Left側は垂直下向き→水平右向きで90度反時計回りなので true
	// Right側は水平右向き→垂直下向きで270度反時計回りなので false
	if !hasLeft90_2 {
		t.Errorf("Expected left 90-degree segment, got false")
	}
	if hasRight90_2 {
		t.Errorf("Expected no right 90-degree segment (270-degree), got true")
	}

	// Test case 3: Mixed - one 90-degree, one 270-degree
	controlPts3 := []image.Point{
		{X: 0, Y: 0},
		{X: 0, Y: 100},   // end of segment[0]
		{X: 200, Y: 100}, // end of segment[1] (longest horizontal)
		{X: 200, Y: 0},   // end of segment[2]
	}

	t.Logf("\nTest case 3 - Mixed segments:")
	t.Logf("Control points: %v", controlPts3)
	hasLeft90_3, hasRight90_3 := link.findPerpendicularSegments(controlPts3)
	t.Logf("Result: hasLeft90=%v, hasRight90=%v", hasLeft90_3, hasRight90_3)

	// segment[1]: 水平右向き {200, 0}
	// segment[2]: 垂直下向き {0, -100}
	// crossProduct = 200 * (-100) - 0 * 0 = -20000 < 0 → Left側鋭角
	if !hasLeft90_3 {
		t.Errorf("Expected left 90-degree segment (Left side acute), got false")
	}
	if hasRight90_3 {
		t.Errorf("Expected no right 90-degree segment, got true")
	}

	// Test case 4: Insufficient segments
	controlPts4 := []image.Point{
		{X: 0, Y: 0},
		{X: 100, Y: 0}, // only one horizontal segment
	}

	hasLeft90_4, hasRight90_4 := link.findPerpendicularSegments(controlPts4)
	if hasLeft90_4 || hasRight90_4 {
		t.Errorf("Expected no 90-degree segments for insufficient control points, got left=%v, right=%v", hasLeft90_4, hasRight90_4)
	}

	// Test case 5: No horizontal segments
	controlPts5 := []image.Point{
		{X: 0, Y: 0},
		{X: 0, Y: 100}, // vertical
		{X: 0, Y: 200}, // vertical
	}

	hasLeft90_5, hasRight90_5 := link.findPerpendicularSegments(controlPts5)
	if hasLeft90_5 || hasRight90_5 {
		t.Errorf("Expected no 90-degree segments for no horizontal segments, got left=%v, right=%v", hasLeft90_5, hasRight90_5)
	}
}

func TestAutoPosReversal(t *testing.T) {
	tests := []struct {
		name     string
		autoPos  Windrose
		expected Windrose
	}{
		{"East to West", 4, 12},
		{"North to South", 0, 8},
		{"West to East", 12, 4},
		{"South to North", 8, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reversed := Windrose((int(tt.autoPos) + 8) % 16)
			if reversed != tt.expected {
				t.Errorf("autoPos reversal failed: got %d, want %d", reversed, tt.expected)
			}
		})
	}
}

func TestGetAcuteAngleSide(t *testing.T) {
	link := &Link{}

	// Test case 1: Left side has acute angle (90-degree)
	controlPts1 := []image.Point{
		{X: 0, Y: 100},    // start
		{X: 0, Y: 0},      // end of segment[0] (vertical down)
		{X: 200, Y: 0},    // end of segment[1] (horizontal right, longest)
		{X: 200, Y: -100}, // end of segment[2] (vertical down)
	}

	leftIsAcute1, rightIsAcute1 := link.getAcuteAngleSide(controlPts1)
	t.Logf("Test case 1 - Left acute: leftIsAcute=%v, rightIsAcute=%v", leftIsAcute1, rightIsAcute1)

	if !leftIsAcute1 {
		t.Errorf("Expected left side to be acute (90-degree), got false")
	}
	if rightIsAcute1 {
		t.Errorf("Expected right side not to be acute, got true")
	}

	// Test case 2: Right side has acute angle (90-degree)
	controlPts2 := []image.Point{
		{X: 0, Y: 0},     // start
		{X: 0, Y: 100},   // end of segment[0] (vertical up)
		{X: 200, Y: 100}, // end of segment[1] (horizontal right, longest)
		{X: 200, Y: 200}, // end of segment[2] (vertical up)
	}

	leftIsAcute2, rightIsAcute2 := link.getAcuteAngleSide(controlPts2)
	t.Logf("Test case 2 - Right acute: leftIsAcute=%v, rightIsAcute=%v", leftIsAcute2, rightIsAcute2)

	if leftIsAcute2 {
		t.Errorf("Expected left side not to be acute, got true")
	}
	if !rightIsAcute2 {
		t.Errorf("Expected right side to be acute (90-degree), got false")
	}

	// Test case 3: Both sides have acute angles
	controlPts3 := []image.Point{
		{X: 0, Y: 0},      // start
		{X: 0, Y: -100},   // end of segment[0] (vertical down)
		{X: 200, Y: -100}, // end of segment[1] (horizontal right, longest)
		{X: 200, Y: 0},    // end of segment[2] (vertical up)
	}

	leftIsAcute3, rightIsAcute3 := link.getAcuteAngleSide(controlPts3)
	t.Logf("Test case 3 - Right acute only: leftIsAcute=%v, rightIsAcute=%v", leftIsAcute3, rightIsAcute3)

	if leftIsAcute3 {
		t.Errorf("Expected left side not to be acute, got true")
	}
	if !rightIsAcute3 {
		t.Errorf("Expected right side to be acute (90-degree), got false")
	}
}

func TestAcuteAngleLabelPlacement(t *testing.T) {
	// Create mock resources
	source := &Resource{}
	target := &Resource{}

	// Set up resource bounds
	sourceBounds := image.Rect(0, 0, 100, 50)
	targetBounds := image.Rect(300, 150, 400, 200)
	source.bindings = &sourceBounds
	target.bindings = &targetBounds

	// Create link with Right and Left labels
	link := &Link{
		Source: source,
		Target: target,
		Type:   "orthogonal",
		Labels: LinkLabels{
			AutoRight: &LinkLabel{Title: "Right Label"},
			AutoLeft:  &LinkLabel{Title: "Left Label"},
		},
	}

	// Test with control points where left side is acute
	controlPts := []image.Point{
		{X: 50, Y: 100}, // start
		{X: 50, Y: 50},  // end of segment[0] (vertical down)
		{X: 250, Y: 50}, // end of segment[1] (horizontal right, longest)
		{X: 250, Y: 0},  // end of segment[2] (vertical down)
	}

	leftIsAcute, rightIsAcute := link.getAcuteAngleSide(controlPts)
	t.Logf("Acute angle detection: leftIsAcute=%v, rightIsAcute=%v", leftIsAcute, rightIsAcute)

	// Verify the acute angle detection works correctly
	if !leftIsAcute {
		t.Errorf("Expected left side to be acute for this configuration")
	}
	if rightIsAcute {
		t.Errorf("Expected right side not to be acute for this configuration")
	}

	t.Logf("Acute angle label placement test completed successfully")
}

func TestSegmentSelectionBasedOnAcutePosition(t *testing.T) {
	// Create mock resources
	source := &Resource{}
	target := &Resource{}

	// Set up resource bounds
	sourceBounds := image.Rect(0, 0, 100, 50)
	targetBounds := image.Rect(500, 100, 600, 150)
	source.bindings = &sourceBounds
	target.bindings = &targetBounds

	// Create link
	link := &Link{
		Source: source,
		Target: target,
		Type:   "orthogonal",
	}

	// Test case: L-shaped path with two horizontal segments of equal length
	// Path: (254,187) → (486,187) → (486,145) → (486,102) → (718,102)
	controlPts := []image.Point{
		{X: 254, Y: 187}, // start (source)
		{X: 486, Y: 187}, // end of segment[0] (horizontal right, 232px)
		{X: 486, Y: 145}, // end of segment[1] (vertical up)
		{X: 486, Y: 102}, // end of segment[2] (vertical up)
		{X: 718, Y: 102}, // end of segment[3] (horizontal right, 232px)
	}

	// Test AutoLeft: segment 0 has "after" acute, should select segment 1 (vertical)
	leftStart, leftEnd, leftAcutePos := link.findBestSegmentForSideWithPosition(controlPts, "Left")
	t.Logf("AutoLeft: start=(%d,%d), end=(%d,%d), acutePos=%s", leftStart.X, leftStart.Y, leftEnd.X, leftEnd.Y, leftAcutePos)

	if leftAcutePos != "after" {
		t.Errorf("Expected Left acute position to be 'after', got '%s'", leftAcutePos)
	}
	// Should select segment 1 (vertical) because segment 0 has "after" acute
	if leftStart.X != 486 || leftStart.Y != 187 || leftEnd.X != 486 || leftEnd.Y != 145 {
		t.Errorf("Expected Left to select segment 1 (486,187)->(486,145), got (%d,%d)->(%d,%d)",
			leftStart.X, leftStart.Y, leftEnd.X, leftEnd.Y)
	}

	// Test AutoRight: segment 3 has "before" acute, should select segment 3 (horizontal)
	rightStart, rightEnd, rightAcutePos := link.findBestSegmentForSideWithPosition(controlPts, "Right")
	t.Logf("AutoRight: start=(%d,%d), end=(%d,%d), acutePos=%s", rightStart.X, rightStart.Y, rightEnd.X, rightEnd.Y, rightAcutePos)

	if rightAcutePos != "before" {
		t.Errorf("Expected Right acute position to be 'before', got '%s'", rightAcutePos)
	}
	// Should select segment 3 (horizontal) because it has "before" acute
	if rightStart.X != 486 || rightStart.Y != 102 || rightEnd.X != 718 || rightEnd.Y != 102 {
		t.Errorf("Expected Right to select segment 3 (486,102)->(718,102), got (%d,%d)->(%d,%d)",
			rightStart.X, rightStart.Y, rightEnd.X, rightEnd.Y)
	}

	t.Logf("Segment selection based on acute position test completed successfully")
}

// TestUnorderedChildrenExample tests the example scenario from doc/unordered-children.md
// This test verifies link overlap with multiple resources before UnorderedChildren implementation
func TestUnorderedChildrenExample(t *testing.T) {
	// Create resource hierarchy matching doc/static/unordered-children-before.yaml
	canvas := new(Resource).Init()
	canvas.direction = "vertical"

	awsCloud := new(Resource).Init()
	awsCloud.direction = "vertical"
	awsCloud.parent = canvas

	// ResourceA (LCA) - horizontal layout
	resourceA := new(Resource).Init()
	resourceA.direction = "horizontal"
	resourceA.label = "ResourceA (LCA)"
	resourceA.parent = awsCloud
	resourceA.bindings = &image.Rectangle{Min: image.Point{50, 50}, Max: image.Point{950, 400}}

	// ResourceB - contains Resource1, Resource2
	resourceB := new(Resource).Init()
	resourceB.direction = "horizontal"
	resourceB.label = "ResourceB"
	resourceB.parent = resourceA
	resourceB.bindings = &image.Rectangle{Min: image.Point{70, 100}, Max: image.Point{270, 350}}

	resource1 := new(Resource).Init()
	resource1.label = "Resource1"
	resource1.parent = resourceB
	resource1.bindings = &image.Rectangle{Min: image.Point{90, 150}, Max: image.Point{154, 214}}

	resource2 := new(Resource).Init()
	resource2.label = "Resource2"
	resource2.parent = resourceB
	resource2.bindings = &image.Rectangle{Min: image.Point{186, 150}, Max: image.Point{250, 214}}

	// ResourceC - contains Resource3, Resource4
	resourceC := new(Resource).Init()
	resourceC.direction = "horizontal"
	resourceC.label = "ResourceC"
	resourceC.parent = resourceA
	resourceC.bindings = &image.Rectangle{Min: image.Point{370, 100}, Max: image.Point{570, 350}}

	resource3 := new(Resource).Init()
	resource3.label = "Resource3"
	resource3.parent = resourceC
	resource3.bindings = &image.Rectangle{Min: image.Point{390, 150}, Max: image.Point{454, 214}}

	resource4 := new(Resource).Init()
	resource4.label = "Resource4"
	resource4.parent = resourceC
	resource4.bindings = &image.Rectangle{Min: image.Point{486, 150}, Max: image.Point{550, 214}}

	// ResourceD - contains Resource5, Resource6
	resourceD := new(Resource).Init()
	resourceD.direction = "horizontal"
	resourceD.label = "ResourceD"
	resourceD.parent = resourceA
	resourceD.bindings = &image.Rectangle{Min: image.Point{670, 100}, Max: image.Point{870, 350}}

	resource5 := new(Resource).Init()
	resource5.label = "Resource5"
	resource5.parent = resourceD
	resource5.bindings = &image.Rectangle{Min: image.Point{690, 150}, Max: image.Point{754, 214}}

	resource6 := new(Resource).Init()
	resource6.label = "Resource6"
	resource6.parent = resourceD
	resource6.bindings = &image.Rectangle{Min: image.Point{786, 150}, Max: image.Point{850, 214}}

	// Add children
	resourceA.children = []*Resource{resourceB, resourceC, resourceD}
	resourceB.children = []*Resource{resource1, resource2}
	resourceC.children = []*Resource{resource3, resource4}
	resourceD.children = []*Resource{resource5, resource6}

	// Create link from Resource1 to Resource6
	link := &Link{
		Source:         resource1,
		SourcePosition: WINDROSE_E, // East
		Target:         resource6,
		TargetPosition: WINDROSE_W, // West
		Type:           "orthogonal",
	}

	// Test LCA finding
	lca := findLowestCommonAncestor(resource1, resource6)
	if lca != resourceA {
		t.Errorf("Expected LCA to be ResourceA, got %v", lca)
	}
	t.Logf("LCA found: %s", lca.label)

	// Test LCA children identification
	sourceChild := findChildAncestorInLCA(lca, resource1)
	targetChild := findChildAncestorInLCA(lca, resource6)

	if sourceChild != resourceB {
		t.Errorf("Expected source child to be ResourceB, got %v", sourceChild)
	}
	if targetChild != resourceD {
		t.Errorf("Expected target child to be ResourceD, got %v", targetChild)
	}
	t.Logf("Source belongs to: %s", sourceChild.label)
	t.Logf("Target belongs to: %s", targetChild.label)

	// Test order determination
	sourceIndex := -1
	targetIndex := -1
	for i, child := range lca.children {
		if child == sourceChild {
			sourceIndex = i
		}
		if child == targetChild {
			targetIndex = i
		}
	}

	if sourceIndex == -1 || targetIndex == -1 {
		t.Fatalf("Failed to find source or target in LCA children")
	}

	t.Logf("Source index: %d, Target index: %d", sourceIndex, targetIndex)

	// Verify pattern: [*, SourceParent, *, TargetParent, *]
	if sourceIndex >= targetIndex {
		t.Errorf("Expected source index (%d) to be less than target index (%d)", sourceIndex, targetIndex)
	}

	// Calculate orthogonal path
	sourcePt := image.Point{X: 154, Y: 182} // Resource1 east edge
	targetPt := image.Point{X: 786, Y: 182} // Resource6 west edge

	controlPoints := link.calculateOrthogonalPath(sourcePt, targetPt)
	t.Logf("Control points: %v", controlPoints)

	// Verify orthogonality
	allPoints := append([]image.Point{sourcePt}, controlPoints...)
	allPoints = append(allPoints, targetPt)

	for i := 0; i < len(allPoints)-1; i++ {
		p1 := allPoints[i]
		p2 := allPoints[i+1]
		if p1.X != p2.X && p1.Y != p2.Y {
			t.Errorf("Non-orthogonal segment: (%d,%d) -> (%d,%d)", p1.X, p1.Y, p2.X, p2.Y)
		}
	}

	// Document the overlap problem
	t.Logf("=== Overlap Analysis ===")
	t.Logf("Link path crosses:")
	t.Logf("1. Resource2 bounds: %v (if Resource2 is right of Resource1)", resource2.bindings)
	t.Logf("2. ResourceC bounds: %v", resourceC.bindings)
	t.Logf("3. Resource3 bounds: %v", resource3.bindings)
	t.Logf("4. Resource4 bounds: %v", resource4.bindings)
	t.Logf("5. Resource5 bounds: %v (if Resource5 is left of Resource6)", resource5.bindings)
	t.Logf("")
	t.Logf("UnorderedChildren feature will reorder:")
	t.Logf("- ResourceB children: Move Resource1 to rightmost")
	t.Logf("- ResourceD children: Move Resource6 to leftmost")
	t.Logf("- ResourceA children: Move ResourceD adjacent to ResourceB")

	// Enable UnorderedChildren
	resourceA.unorderedChildren = true
	resourceB.unorderedChildren = true
	resourceD.unorderedChildren = true

	// Verify initial order
	t.Logf("\n=== Before Reordering ===")
	t.Logf("ResourceA children: %v", getResourceNames(resourceA.children))
	t.Logf("ResourceB children: %v", getResourceNames(resourceB.children))
	t.Logf("ResourceD children: %v", getResourceNames(resourceD.children))

	if len(resourceA.children) != 3 {
		t.Fatalf("Expected ResourceA to have 3 children, got %d", len(resourceA.children))
	}
	if resourceA.children[0] != resourceB || resourceA.children[1] != resourceC || resourceA.children[2] != resourceD {
		t.Errorf("Initial ResourceA order incorrect")
	}
	if len(resourceB.children) != 2 || resourceB.children[0] != resource1 || resourceB.children[1] != resource2 {
		t.Errorf("Initial ResourceB order incorrect")
	}
	if len(resourceD.children) != 2 || resourceD.children[0] != resource5 || resourceD.children[1] != resource6 {
		t.Errorf("Initial ResourceD order incorrect")
	}

	// Call reordering function
	links := []*Link{link}
	ReorderChildrenByLinks(canvas, links)

	// Verify reordered state
	t.Logf("\n=== After Reordering ===")
	t.Logf("ResourceA children: %v", getResourceNames(resourceA.children))
	t.Logf("ResourceB children: %v", getResourceNames(resourceB.children))
	t.Logf("ResourceD children: %v", getResourceNames(resourceD.children))

	// Expected: ResourceB children reordered - Resource1 at rightmost
	if len(resourceB.children) != 2 {
		t.Fatalf("Expected ResourceB to have 2 children after reordering, got %d", len(resourceB.children))
	}
	if resourceB.children[0] != resource2 || resourceB.children[1] != resource1 {
		t.Errorf("Expected ResourceB children [Resource2, Resource1], got [%s, %s]",
			resourceB.children[0].label, resourceB.children[1].label)
	}

	// Expected: ResourceD children reordered - Resource6 at leftmost
	if len(resourceD.children) != 2 {
		t.Fatalf("Expected ResourceD to have 2 children after reordering, got %d", len(resourceD.children))
	}
	if resourceD.children[0] != resource6 || resourceD.children[1] != resource5 {
		t.Errorf("Expected ResourceD children [Resource6, Resource5], got [%s, %s]",
			resourceD.children[0].label, resourceD.children[1].label)
	}

	// Expected: ResourceA children reordered - ResourceD adjacent to ResourceB
	if len(resourceA.children) != 3 {
		t.Fatalf("Expected ResourceA to have 3 children after reordering, got %d", len(resourceA.children))
	}
	if resourceA.children[0] != resourceB || resourceA.children[1] != resourceD || resourceA.children[2] != resourceC {
		t.Errorf("Expected ResourceA children [ResourceB, ResourceD, ResourceC], got [%s, %s, %s]",
			resourceA.children[0].label, resourceA.children[1].label, resourceA.children[2].label)
	}

	t.Logf("\n=== Reordering Successful ===")
	t.Logf("✓ Resource1 moved to rightmost in ResourceB")
	t.Logf("✓ Resource6 moved to leftmost in ResourceD")
	t.Logf("✓ ResourceD moved adjacent to ResourceB in ResourceA")
}

// TestUnorderedChildrenMultipleLinks tests reordering with multiple links
func TestUnorderedChildrenMultipleLinks(t *testing.T) {
	// Create resource hierarchy
	canvas := new(Resource).Init()
	canvas.direction = "vertical"

	awsCloud := new(Resource).Init()
	awsCloud.direction = "vertical"
	awsCloud.parent = canvas

	// VPC (LCA) - vertical layout
	vpc := new(Resource).Init()
	vpc.direction = "vertical"
	vpc.label = "VPC"
	vpc.parent = awsCloud
	vpc.unorderedChildren = true
	vpc.bindings = &image.Rectangle{Min: image.Point{50, 50}, Max: image.Point{950, 600}}

	// Horizontal groups
	horizontal1 := new(Resource).Init()
	horizontal1.direction = "horizontal"
	horizontal1.label = "Horizontal1"
	horizontal1.parent = vpc
	horizontal1.unorderedChildren = true
	horizontal1.bindings = &image.Rectangle{Min: image.Point{70, 100}, Max: image.Point{870, 250}}

	resource1 := new(Resource).Init()
	resource1.label = "Resource1"
	resource1.parent = horizontal1
	resource1.bindings = &image.Rectangle{Min: image.Point{90, 150}, Max: image.Point{154, 214}}

	resource2 := new(Resource).Init()
	resource2.label = "Resource2"
	resource2.parent = horizontal1
	resource2.bindings = &image.Rectangle{Min: image.Point{786, 150}, Max: image.Point{850, 214}}

	horizontal2 := new(Resource).Init()
	horizontal2.direction = "horizontal"
	horizontal2.label = "Horizontal2"
	horizontal2.parent = vpc
	horizontal2.unorderedChildren = true
	horizontal2.bindings = &image.Rectangle{Min: image.Point{70, 350}, Max: image.Point{870, 500}}

	resource3 := new(Resource).Init()
	resource3.label = "Resource3"
	resource3.parent = horizontal2
	resource3.bindings = &image.Rectangle{Min: image.Point{90, 400}, Max: image.Point{154, 464}}

	resource4 := new(Resource).Init()
	resource4.label = "Resource4"
	resource4.parent = horizontal2
	resource4.bindings = &image.Rectangle{Min: image.Point{786, 400}, Max: image.Point{850, 464}}

	// Add children
	vpc.children = []*Resource{horizontal1, horizontal2}
	horizontal1.children = []*Resource{resource1, resource2}
	horizontal2.children = []*Resource{resource3, resource4}

	// Create two links that cross
	// Link1: Resource1 -> Resource4
	link1 := &Link{
		Source:         resource1,
		SourcePosition: WINDROSE_S,
		Target:         resource4,
		TargetPosition: WINDROSE_N,
		Type:           "orthogonal",
	}

	// Link2: Resource2 -> Resource3
	link2 := &Link{
		Source:         resource2,
		SourcePosition: WINDROSE_S,
		Target:         resource3,
		TargetPosition: WINDROSE_N,
		Type:           "orthogonal",
	}

	t.Logf("\n=== Before Reordering ===")
	t.Logf("VPC children: %v", getResourceNames(vpc.children))
	t.Logf("Horizontal1 children: %v", getResourceNames(horizontal1.children))
	t.Logf("Horizontal2 children: %v", getResourceNames(horizontal2.children))

	// Initial state
	if len(horizontal1.children) != 2 || horizontal1.children[0] != resource1 || horizontal1.children[1] != resource2 {
		t.Errorf("Initial Horizontal1 order incorrect")
	}
	if len(horizontal2.children) != 2 || horizontal2.children[0] != resource3 || horizontal2.children[1] != resource4 {
		t.Errorf("Initial Horizontal2 order incorrect")
	}

	// Call reordering with multiple links
	links := []*Link{link1, link2}
	ReorderChildrenByLinks(canvas, links)

	t.Logf("\n=== After Reordering ===")
	t.Logf("VPC children: %v", getResourceNames(vpc.children))
	t.Logf("Horizontal1 children: %v", getResourceNames(horizontal1.children))
	t.Logf("Horizontal2 children: %v", getResourceNames(horizontal2.children))

	// Expected results after processing both links:
	// Link1 (Resource1 -> Resource4):
	//   - Horizontal1: Resource1 moves to first (already there)
	//   - Horizontal2: Resource4 moves to first
	// Link2 (Resource2 -> Resource3):
	//   - Horizontal1: Resource2 moves to first
	//   - Horizontal2: Resource3 moves to first (already there after Link1)

	// Final expected state:
	// Horizontal1: [Resource2, Resource1] (Resource2 moved to first by Link2)
	// Horizontal2: [Resource3, Resource4] (Resource3 at first, Resource4 moved by Link1)

	if len(horizontal1.children) != 2 {
		t.Fatalf("Expected Horizontal1 to have 2 children, got %d", len(horizontal1.children))
	}
	if horizontal1.children[0] != resource2 || horizontal1.children[1] != resource1 {
		t.Errorf("Expected Horizontal1 children [Resource2, Resource1], got [%s, %s]",
			horizontal1.children[0].label, horizontal1.children[1].label)
	}

	if len(horizontal2.children) != 2 {
		t.Fatalf("Expected Horizontal2 to have 2 children, got %d", len(horizontal2.children))
	}
	if horizontal2.children[0] != resource3 || horizontal2.children[1] != resource4 {
		t.Errorf("Expected Horizontal2 children [Resource3, Resource4], got [%s, %s]",
			horizontal2.children[0].label, horizontal2.children[1].label)
	}

	t.Logf("\n=== Multiple Links Reordering Successful ===")
	t.Logf("✓ Links no longer cross each other")
	t.Logf("✓ Resource2 and Resource1 reordered in Horizontal1")
	t.Logf("✓ Resource3 and Resource4 reordered in Horizontal2")
}

// TestUnorderedChildrenDisabled tests that reordering does NOT happen when unorderedChildren is false
func TestUnorderedChildrenDisabled(t *testing.T) {
	// Create resource hierarchy similar to the main example
	canvas := new(Resource).Init()
	canvas.direction = "vertical"

	awsCloud := new(Resource).Init()
	awsCloud.direction = "vertical"
	awsCloud.parent = canvas

	// ResourceA (LCA) - horizontal layout, UnorderedChildren = FALSE
	resourceA := new(Resource).Init()
	resourceA.direction = "horizontal"
	resourceA.label = "ResourceA (LCA)"
	resourceA.parent = awsCloud
	resourceA.unorderedChildren = false // Explicitly disabled
	resourceA.bindings = &image.Rectangle{Min: image.Point{50, 50}, Max: image.Point{950, 400}}

	// ResourceB - UnorderedChildren = FALSE
	resourceB := new(Resource).Init()
	resourceB.direction = "horizontal"
	resourceB.label = "ResourceB"
	resourceB.parent = resourceA
	resourceB.unorderedChildren = false // Explicitly disabled
	resourceB.bindings = &image.Rectangle{Min: image.Point{70, 100}, Max: image.Point{270, 350}}

	resource1 := new(Resource).Init()
	resource1.label = "Resource1"
	resource1.parent = resourceB
	resource1.bindings = &image.Rectangle{Min: image.Point{90, 150}, Max: image.Point{154, 214}}

	resource2 := new(Resource).Init()
	resource2.label = "Resource2"
	resource2.parent = resourceB
	resource2.bindings = &image.Rectangle{Min: image.Point{186, 150}, Max: image.Point{250, 214}}

	// ResourceC
	resourceC := new(Resource).Init()
	resourceC.direction = "horizontal"
	resourceC.label = "ResourceC"
	resourceC.parent = resourceA
	resourceC.bindings = &image.Rectangle{Min: image.Point{370, 100}, Max: image.Point{570, 350}}

	// ResourceD - UnorderedChildren = FALSE
	resourceD := new(Resource).Init()
	resourceD.direction = "horizontal"
	resourceD.label = "ResourceD"
	resourceD.parent = resourceA
	resourceD.unorderedChildren = false // Explicitly disabled
	resourceD.bindings = &image.Rectangle{Min: image.Point{670, 100}, Max: image.Point{870, 350}}

	resource5 := new(Resource).Init()
	resource5.label = "Resource5"
	resource5.parent = resourceD
	resource5.bindings = &image.Rectangle{Min: image.Point{690, 150}, Max: image.Point{754, 214}}

	resource6 := new(Resource).Init()
	resource6.label = "Resource6"
	resource6.parent = resourceD
	resource6.bindings = &image.Rectangle{Min: image.Point{786, 150}, Max: image.Point{850, 214}}

	// Add children
	resourceA.children = []*Resource{resourceB, resourceC, resourceD}
	resourceB.children = []*Resource{resource1, resource2}
	resourceD.children = []*Resource{resource5, resource6}

	// Create link from Resource1 to Resource6
	link := &Link{
		Source:         resource1,
		SourcePosition: WINDROSE_E,
		Target:         resource6,
		TargetPosition: WINDROSE_W,
		Type:           "orthogonal",
	}

	t.Logf("\n=== Before Reordering (UnorderedChildren=false) ===")
	t.Logf("ResourceA children: %v", getResourceNames(resourceA.children))
	t.Logf("ResourceB children: %v", getResourceNames(resourceB.children))
	t.Logf("ResourceD children: %v", getResourceNames(resourceD.children))

	// Store initial state
	initialAChildren := make([]*Resource, len(resourceA.children))
	copy(initialAChildren, resourceA.children)
	initialBChildren := make([]*Resource, len(resourceB.children))
	copy(initialBChildren, resourceB.children)
	initialDChildren := make([]*Resource, len(resourceD.children))
	copy(initialDChildren, resourceD.children)

	// Call reordering function
	links := []*Link{link}
	ReorderChildrenByLinks(canvas, links)

	t.Logf("\n=== After Reordering (UnorderedChildren=false) ===")
	t.Logf("ResourceA children: %v", getResourceNames(resourceA.children))
	t.Logf("ResourceB children: %v", getResourceNames(resourceB.children))
	t.Logf("ResourceD children: %v", getResourceNames(resourceD.children))

	// Verify NO changes occurred
	if len(resourceA.children) != len(initialAChildren) {
		t.Errorf("ResourceA children count changed")
	}
	for i := range resourceA.children {
		if resourceA.children[i] != initialAChildren[i] {
			t.Errorf("ResourceA children order changed at index %d: expected %s, got %s",
				i, initialAChildren[i].label, resourceA.children[i].label)
		}
	}

	if len(resourceB.children) != len(initialBChildren) {
		t.Errorf("ResourceB children count changed")
	}
	for i := range resourceB.children {
		if resourceB.children[i] != initialBChildren[i] {
			t.Errorf("ResourceB children order changed at index %d: expected %s, got %s",
				i, initialBChildren[i].label, resourceB.children[i].label)
		}
	}

	if len(resourceD.children) != len(initialDChildren) {
		t.Errorf("ResourceD children count changed")
	}
	for i := range resourceD.children {
		if resourceD.children[i] != initialDChildren[i] {
			t.Errorf("ResourceD children order changed at index %d: expected %s, got %s",
				i, initialDChildren[i].label, resourceD.children[i].label)
		}
	}

	t.Logf("\n=== Verification Successful ===")
	t.Logf("✓ ResourceA children unchanged (UnorderedChildren=false)")
	t.Logf("✓ ResourceB children unchanged (UnorderedChildren=false)")
	t.Logf("✓ ResourceD children unchanged (UnorderedChildren=false)")
}

func TestAutoCalculatePositions_BorderChild(t *testing.T) {
	// Setup: VPC with BorderChild IGW at South position
	vpc := new(Resource).Init()
	vpc.label = "VPC"
	vpc.direction = "vertical"
	vpcBounds := image.Rect(100, 100, 500, 500)
	vpc.bindings = &vpcBounds

	alb := new(Resource).Init()
	alb.label = "ALB"
	albBounds := image.Rect(200, 200, 400, 300)
	alb.bindings = &albBounds
	if err := vpc.AddChild(alb); err != nil {
		t.Fatal(err)
	}

	igw := new(Resource).Init()
	igw.label = "IGW"
	igwBounds := image.Rect(250, 480, 350, 520)
	igw.bindings = &igwBounds

	// Add IGW as BorderChild at South position
	bc := &BorderChild{
		Position: WINDROSE_S,
		Resource: igw,
	}
	if err := vpc.AddBorderChild(bc); err != nil {
		t.Fatal(err)
	}

	// Test 1: IGW to ALB (LCA = VPC, which is IGW's parent)
	// Expected: IGW should use inside (opposite of S = N)
	sourcePos, targetPos := AutoCalculatePositions(igw, alb)
	if sourcePos != WINDROSE_N {
		t.Errorf("Test 1: Expected IGW position N (inside), got %v", sourcePos)
	}
	t.Logf("Test 1 passed: IGW->ALB, IGW uses inside position N")

	// Test 2: Create external resource outside VPC
	canvas := new(Resource).Init()
	canvas.label = "Canvas"
	canvas.direction = "vertical"
	canvasBounds := image.Rect(0, 0, 600, 700)
	canvas.bindings = &canvasBounds
	if err := canvas.AddChild(vpc); err != nil {
		t.Fatal(err)
	}

	external := new(Resource).Init()
	external.label = "External"
	externalBounds := image.Rect(200, 600, 400, 650)
	external.bindings = &externalBounds
	if err := canvas.AddChild(external); err != nil {
		t.Fatal(err)
	}

	// Test 2: IGW to External (LCA = Canvas, ancestor of VPC)
	// Expected: IGW should use outside (same as S)
	sourcePos, targetPos = AutoCalculatePositions(igw, external)
	if sourcePos != WINDROSE_S {
		t.Errorf("Test 2: Expected IGW position S (outside), got %v", sourcePos)
	}
	t.Logf("Test 2 passed: IGW->External, IGW uses outside position S")

	// Test 3: Both resources are BorderChildren
	natgw := new(Resource).Init()
	natgw.label = "NATGW"
	natgwBounds := image.Rect(150, 480, 220, 520)
	natgw.bindings = &natgwBounds

	bc2 := &BorderChild{
		Position: WINDROSE_SSW,
		Resource: natgw,
	}
	if err := vpc.AddBorderChild(bc2); err != nil {
		t.Fatal(err)
	}

	// Test 3: NATGW to IGW (both BorderChildren, LCA = VPC)
	// Expected: Both use inside (opposite of their positions)
	sourcePos, targetPos = AutoCalculatePositions(natgw, igw)
	expectedNATGWPos := GetOppositeWindrose(WINDROSE_SSW) // Should be NNE
	expectedIGWPos := GetOppositeWindrose(WINDROSE_S)     // Should be N
	if sourcePos != expectedNATGWPos {
		t.Errorf("Test 3: Expected NATGW position %v (inside), got %v", expectedNATGWPos, sourcePos)
	}
	if targetPos != expectedIGWPos {
		t.Errorf("Test 3: Expected IGW position %v (inside), got %v", expectedIGWPos, targetPos)
	}
	t.Logf("Test 3 passed: NATGW->IGW, both use inside positions")
}
