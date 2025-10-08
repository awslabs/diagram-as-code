package types

import (
	"image"
	"image/color"
	"testing"
)

func TestMaxInt(t *testing.T) {
	testCases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 5, 3, 5},
		{"negative numbers", -10, -20, -10},
		{"equal numbers", 42, 42, 42},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := maxInt(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("maxInt(%d, %d) = %d; expected %d", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

func TestMinInt(t *testing.T) {
	testCases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 5, 3, 3},
		{"negative numbers", -10, -20, -20},
		{"equal numbers", 42, 42, 42},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := minInt(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("minInt(%d, %d) = %d; expected %d", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

func TestFetchColor(t *testing.T) {
	testCases := []struct {
		name      string
		color     color.Color
		expectedR uint32
		expectedG uint32
		expectedB uint32
		expectedA uint32
	}{
		{"opaque color", color.RGBA{255, 128, 64, 255}, 255, 128, 64, 255},
		{"transparent color", color.RGBA{100, 200, 150, 128}, 100, 200, 150, 128},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, g, b, a := _fetch_color(tc.color)
			if r != tc.expectedR || g != tc.expectedG || b != tc.expectedB || a != tc.expectedA {
				t.Errorf("_fetch_color(%v) = (%d, %d, %d, %d); expected (%d, %d, %d, %d)", tc.color, r, g, b, a, tc.expectedR, tc.expectedG, tc.expectedB, tc.expectedA)
			}
		})
	}
}

func TestBlendColor(t *testing.T) {
	testCases := []struct {
		name           string
		color1, color2 color.Color
		expected       color.Color
	}{
		{"opaque colors", color.RGBA{255, 0, 0, 125}, color.RGBA{0, 255, 0, 125}, color.RGBA{130, 125, 0, 125}},
		{"transparent colors", color.RGBA{255, 0, 0, 128}, color.RGBA{0, 255, 0, 64}, color.RGBA{191, 64, 0, 128}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := _blend_color(tc.color1, tc.color2)
			if result != tc.expected {
				t.Errorf("_blend_color(%v, %v) = %v; expected %v", tc.color1, tc.color2, result, tc.expected)
			}
		})
	}
}

func TestCalcPosition(t *testing.T) {
	testCases := []struct {
		name     string
		bindings image.Rectangle
		position Windrose
		expected image.Point
		err      bool
	}{
		{"center", image.Rect(0, 0, 100, 100), WINDROSE_N, image.Point{50, 0}, false},
		{"invalid position", image.Rect(0, 0, 100, 100), -1, image.Point{0, 0}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calcPosition(tc.bindings, tc.position)
			if err != nil && !tc.err {
				t.Errorf("calcPosition(%v, %d) returned unexpected error: %v", tc.bindings, tc.position, err)
			} else if err == nil && tc.err {
				t.Errorf("calcPosition(%v, %d) did not return expected error", tc.bindings, tc.position)
			} else if result != tc.expected {
				t.Errorf("calcPosition(%v, %d) = %v; expected %v", tc.bindings, tc.position, result, tc.expected)
			}
		})
	}
}

func TestConvertWindrose(t *testing.T) {
	// Test normal directions
	pos, err := ConvertWindrose("N")
	if err != nil || pos != WINDROSE_N {
		t.Errorf("Expected WINDROSE_N, got %v, err: %v", pos, err)
	}

	// Test empty string (default to auto)
	pos, err = ConvertWindrose("")
	if err != nil || pos != WINDROSE_AUTO {
		t.Errorf("Expected WINDROSE_AUTO for empty string, got %v, err: %v", pos, err)
	}

	// Test explicit "auto"
	pos, err = ConvertWindrose("auto")
	if err != nil || pos != WINDROSE_AUTO {
		t.Errorf("Expected WINDROSE_AUTO for 'auto', got %v, err: %v", pos, err)
	}

	// Test invalid position
	_, err = ConvertWindrose("invalid")
	if err == nil {
		t.Error("Expected error for invalid position")
	}
}
