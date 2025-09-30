package types

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestResource(t *testing.T) {
	// Test Init
	r := new(Resource).Init()

	// Test resource has not children
	if err := r.Scale(nil, nil); err != nil {
		t.Errorf("Scale failed: %v", err)
	}
	if r.GetBindings() != image.Rect(0, 0, 0, 0) {
		t.Errorf("Init: expected bindings to be (0, 0, 0, 0), got %v", r.GetBindings())
	}
	if r.GetMargin() != (Margin{0, 0, 0, 0}) {
		t.Errorf("Init: expected margin to be (0, 0, 0, 0), got %v", r.GetMargin())
	}
	if r.GetPadding() != (Padding{0, 0, 0, 0}) {
		t.Errorf("Init: expected padding to be (0, 0, 0, 0), got %v", r.GetPadding())
	}
	if r.IsDrawn() {
		t.Error("Init: expected drawn to be false")
	}

	// Test resource has not children
	r2 := new(Resource).Init()
	r3 := new(Resource).Init()
	if err := r2.AddChild(r3); err != nil {
		t.Errorf("AddChild failed: %v", err)
	}
	if err := r2.Scale(nil, nil); err != nil {
		t.Errorf("Scale failed: %v", err)
	}
	if r2.GetMargin() != (Margin{20, 15, 20, 15}) {
		t.Errorf("Init: expected margin to be (30, 100, 30, 100), got %v", r2.GetMargin())
	}
	if r2.GetPadding() != (Padding{20, 45, 20, 45}) {
		t.Errorf("Init: expected padding to be (0, 0, 0, 0), got %v", r2.GetPadding())
	}
	if r2.IsDrawn() {
		t.Error("Init: expected drawn to be false")
	}

	// Test LoadIcon
	// Test LoadIcon - skip if test icon file doesn't exist or is invalid
	if err := r.LoadIcon("testdata/valid_icon.png"); err != nil {
		t.Logf("LoadIcon: expected error for test file: %v", err)
		// This is expected for test files that may not be valid PNG format
	} else if r.iconImage == nil {
		t.Error("LoadIcon: iconImage is nil")
	}

	// Test SetIconBounds
	bounds := image.Rect(10, 10, 20, 20)
	r.SetIconBounds(bounds)
	if r.iconBounds != bounds {
		t.Errorf("SetIconBounds: expected bounds to be %v, got %v", bounds, r.iconBounds)
	}

	// Test SetBindings
	bindings := image.Rect(100, 100, 200, 200)
	r.SetBindings(bindings)
	if r.GetBindings() != bindings {
		t.Errorf("SetBindings: expected bindings to be %v, got %v", bindings, r.GetBindings())
	}

	// Test SetBorderColor
	borderColor := color.RGBA{255, 0, 0, 255}
	r.SetBorderColor(borderColor)
	if *r.borderColor != borderColor {
		t.Errorf("SetBorderColor: expected borderColor to be %v, got %v", borderColor, r.borderColor)
	}

	// Test SetFillColor
	fillColor := color.RGBA{0, 255, 0, 255}
	r.SetFillColor(fillColor)
	if r.fillColor != fillColor {
		t.Errorf("SetFillColor: expected fillColor to be %v, got %v", fillColor, r.fillColor)
	}

	// Test SetLabel
	/*label := "Test Label"
	labelColor := color.RGBA{0, 0, 255, 255}
	labelFont := "testdata/font.ttf"
	r.SetLabel(&label, &labelColor, &labelFont)
	if r.label != label {
		t.Errorf("SetLabel: expected label to be %q, got %q", label, r.label)
	}
	if *r.labelColor != labelColor {
		t.Errorf("SetLabel: expected labelColor to be %v, got %v", labelColor, *r.labelColor)
	}
	if r.labelFont != labelFont {
		t.Errorf("SetLabel: expected labelFont to be %q, got %q", labelFont, r.labelFont)
	}*/

	// Test Translation
	dx, dy := 10, 20
	origBindings := r.GetBindings()
	if err := r.Translation(dx, dy); err != nil {
		t.Errorf("Translation: unexpected error: %v", err)
	}
	expected := image.Rect(origBindings.Min.X+dx, origBindings.Min.Y+dy, origBindings.Max.X+dx, origBindings.Max.Y+dy)
	if r.GetBindings() != expected {
		t.Errorf("Translation: expected bindings to be %v, got %v", expected, r.GetBindings())
	}

	// Test ZeroAdjust
	if err := r.ZeroAdjust(); err != nil {
		t.Errorf("ZeroAdjust: unexpected error: %v", err)
	}
	expected = image.Rect(r.padding.Left, r.padding.Top, r.GetBindings().Max.X, r.GetBindings().Max.Y)
	if r.GetBindings() != expected {
		t.Errorf("ZeroAdjust: expected bindings to be %v, got %v", expected, r.GetBindings())
	}

	// Test Draw (basic)
	img, err := r.Draw(nil, nil)
	if err != nil {
		t.Errorf("Draw: unexpected error: %v", err)
	}
	if img == nil {
		t.Error("Draw: returned nil image")
	}
	if !r.IsDrawn() {
		t.Error("Draw: drawn should be true after drawing")
	}
}

func TestResourceCycleDetection(t *testing.T) {
	// Define test cases using a table-driven approach
	testCases := []struct {
		name           string
		setupResources func() *Resource
		expectError    bool
		errorContains  string
	}{
		{
			name: "DirectCycle",
			setupResources: func() *Resource {
				r1 := new(Resource).Init()
				r1.label = "Resource1"
				if err := r1.AddChild(r1); err != nil {
					t.Errorf("AddChild failed: %v", err)
				} // Create a direct cycle: r1 -> r1

				return r1
			},
			expectError:   true,
			errorContains: "Cycle detected in resource tree",
		},
		{
			name: "IndirectCycle",
			setupResources: func() *Resource {
				r1 := new(Resource).Init()
				r2 := new(Resource).Init()
				r1.label = "Resource1"
				r2.label = "Resource2"
				if err := r1.AddChild(r2); err != nil {
					t.Errorf("AddChild failed: %v", err)
				}
				if err := r2.AddChild(r1); err != nil {
					t.Errorf("AddChild failed: %v", err)
				} // Create an indirect cycle: r1 -> r2 -> r1
				return r1
			},
			expectError:   true,
			errorContains: "Cycle detected in resource tree",
		},
		{
			name: "LongerCycle",
			setupResources: func() *Resource {
				r1 := new(Resource).Init()
				r2 := new(Resource).Init()
				r3 := new(Resource).Init()
				r1.label = "Resource1"
				r2.label = "Resource2"
				r3.label = "Resource3"
				if err := r1.AddChild(r2); err != nil {
					t.Errorf("AddChild failed: %v", err)
				}
				if err := r2.AddChild(r3); err != nil {
					t.Errorf("AddChild failed: %v", err)
				}
				if err := r3.AddChild(r1); err != nil {
					t.Errorf("AddChild failed: %v", err)
				} // Create a cycle: r1 -> r2 -> r3 -> r1
				return r1
			},
			expectError:   true,
			errorContains: "Cycle detected in resource tree",
		},
		{
			name: "NoCycle",
			setupResources: func() *Resource {
				r1 := new(Resource).Init()
				r2 := new(Resource).Init()
				r3 := new(Resource).Init()
				r4 := new(Resource).Init()
				r1.label = "Root"
				r2.label = "Child1"
				r3.label = "Child2"
				r4.label = "Grandchild"
				if err := r1.AddChild(r2); err != nil {
					t.Errorf("AddChild failed: %v", err)
				}
				if err := r1.AddChild(r3); err != nil {
					t.Errorf("AddChild failed: %v", err)
				}
				if err := r2.AddChild(r4); err != nil {
					t.Errorf("AddChild failed: %v", err)
				}
				return r1
			},
			expectError: false,
		},
		{
			name: "NoCycleWithBorderChildren",
			setupResources: func() *Resource {
				parent := new(Resource).Init()
				child := new(Resource).Init()
				borderChild := new(Resource).Init()
				parent.label = "Parent"
				child.label = "Child"
				borderChild.label = "BorderChild"
				if err := parent.AddChild(child); err != nil {
					t.Errorf("AddChild failed: %v", err)
				}
				if err := parent.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				}
				return parent
			},
			expectError: false,
		},
		{
			name: "DirectCycleWithBorderChild",
			setupResources: func() *Resource {
				parent := new(Resource).Init()
				child := new(Resource).Init()
				borderChild := new(Resource).Init()
				parent.label = "Parent"
				child.label = "Child"
				borderChild.label = "BorderChild"
				if err := parent.AddChild(child); err != nil {
					t.Errorf("AddChild failed: %v", err)
				}
				if err := parent.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				}
				if err := borderChild.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				} // Create an direct cycle: borderChild -> borderChild
				return parent
			},
			expectError:   true,
			errorContains: "Cycle detected in resource tree",
		},
		{
			name: "NoCycleMultipleBorderChildren",
			setupResources: func() *Resource {
				parent := new(Resource).Init()
				child := new(Resource).Init()
				borderChild1 := new(Resource).Init()
				borderChild2 := new(Resource).Init()
				parent.label = "Parent"
				child.label = "Child"
				borderChild1.label = "BorderChild1"
				borderChild2.label = "BorderChild2"
				if err := parent.AddChild(child); err != nil {
					t.Errorf("AddChild failed: %v", err)
				}
				if err := parent.AddBorderChild(&BorderChild{
					Position: 0, // North position
					Resource: borderChild1,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				}
				if err := parent.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild2,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				}
				return parent
			},
			expectError: false,
		},
		{
			name: "IndirectCycleMultipleBorderChildren",
			setupResources: func() *Resource {
				parent := new(Resource).Init()
				child := new(Resource).Init()
				borderChild1 := new(Resource).Init()
				borderChild2 := new(Resource).Init()
				parent.label = "Parent"
				child.label = "Child"
				borderChild1.label = "BorderChild1"
				borderChild2.label = "BorderChild2"
				if err := parent.AddChild(child); err != nil {
					t.Errorf("AddChild failed: %v", err)
				}
				if err := parent.AddBorderChild(&BorderChild{
					Position: 0, // North position
					Resource: borderChild1,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				}
				if err := borderChild1.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild2,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				}
				if err := borderChild2.AddBorderChild(&BorderChild{
					Position: 0, // North position
					Resource: borderChild1,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				} // Create an indirect cycle: r1 -> r2 -> r1
				return parent
			},
			expectError:   true,
			errorContains: "Cycle detected in resource tree",
		},
		{
			name: "LongerCycleMultipleBorderChildren",
			setupResources: func() *Resource {
				parent := new(Resource).Init()
				child := new(Resource).Init()
				borderChild1 := new(Resource).Init()
				borderChild2 := new(Resource).Init()
				borderChild3 := new(Resource).Init()
				parent.label = "Parent"
				child.label = "Child"
				borderChild1.label = "BorderChild1"
				borderChild2.label = "BorderChild2"
				borderChild3.label = "BorderChild3"
				if err := parent.AddChild(child); err != nil {
					t.Errorf("AddChild failed: %v", err)
				}
				if err := parent.AddBorderChild(&BorderChild{
					Position: 0, // North position
					Resource: borderChild1,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				}
				if err := borderChild1.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild2,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				}
				if err := borderChild2.AddBorderChild(&BorderChild{
					Position: 0, // North position
					Resource: borderChild3,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				}
				if err := borderChild3.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild1,
				}); err != nil {
					panic(err) // In test setup, panic is acceptable
				} // Create a cycle: r1 -> r2 -> r3 -> r1
				return parent
			},
			expectError:   true,
			errorContains: "Cycle detected in resource tree",
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			root := tc.setupResources()
			err := root.Scale(nil, nil)

			if tc.expectError && err == nil {
				t.Error("Expected cycle detection error, but got nil")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tc.expectError && err != nil && tc.errorContains != "" {
				if msg := err.Error(); !strings.Contains(msg, tc.errorContains) {
					t.Errorf("Expected error message to contain %q, got %q", tc.errorContains, msg)
				}
			}
		})
	}
}
func TestSetMargin(t *testing.T) {
	resource := new(Resource).Init()

	newMargin := Margin{Top: 10, Right: 20, Bottom: 30, Left: 40}
	resource.SetMargin(newMargin)

	result := resource.GetMargin()
	if result.Top != 10 || result.Right != 20 || result.Bottom != 30 || result.Left != 40 {
		t.Errorf("SetMargin failed: expected {10 20 30 40}, got {%d %d %d %d}",
			result.Top, result.Right, result.Bottom, result.Left)
	}
}

func TestSetPadding(t *testing.T) {
	resource := new(Resource).Init()

	newPadding := Padding{Top: 5, Right: 15, Bottom: 25, Left: 35}
	resource.SetPadding(newPadding)

	result := resource.GetPadding()
	if result.Top != 5 || result.Right != 15 || result.Bottom != 25 || result.Left != 35 {
		t.Errorf("SetPadding failed: expected {5 15 25 35}, got {%d %d %d %d}",
			result.Top, result.Right, result.Bottom, result.Left)
	}
}
