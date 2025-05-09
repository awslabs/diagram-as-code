package types

import (
	"image"
	"image/color"
	"os"
	"strings"
	"testing"
)

func TestResource(t *testing.T) {
	// Test Init
	r := new(Resource).Init()

	// Test resource has not children
	r.Scale(nil, nil)
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
	r2.AddChild(r3)
	r2.Scale(nil, nil)
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
	iconFile, err := os.Open("testdata/valid_icon.png")
	if err != nil {
		t.Fatalf("LoadIcon: failed to open test icon file: %v", err)
	}
	defer iconFile.Close()
	r.LoadIcon("testdata/icon.png")
	if r.iconImage == nil {
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
	r.Translation(dx, dy)
	expected := image.Rect(origBindings.Min.X+dx, origBindings.Min.Y+dy, origBindings.Max.X+dx, origBindings.Max.Y+dy)
	if r.GetBindings() != expected {
		t.Errorf("Translation: expected bindings to be %v, got %v", expected, r.GetBindings())
	}

	// Test ZeroAdjust
	r.ZeroAdjust()
	expected = image.Rect(r.padding.Left, r.padding.Top, r.GetBindings().Max.X, r.GetBindings().Max.Y)
	if r.GetBindings() != expected {
		t.Errorf("ZeroAdjust: expected bindings to be %v, got %v", expected, r.GetBindings())
	}

	// Test Draw (basic)
	img := r.Draw(nil, nil)
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
				r1.AddChild(r1) // Create a direct cycle: r1 -> r1

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
				r1.AddChild(r2)
				r2.AddChild(r1) // Create an indirect cycle: r1 -> r2 -> r1
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
				r1.AddChild(r2)
				r2.AddChild(r3)
				r3.AddChild(r1) // Create a cycle: r1 -> r2 -> r3 -> r1
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
				r1.AddChild(r2)
				r1.AddChild(r3)
				r2.AddChild(r4)
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
				parent.AddChild(child)
				parent.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild,
				})
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
				parent.AddChild(child)
				parent.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild,
				})
				borderChild.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild,
				}) // Create an direct cycle: borderChild -> borderChild
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
				parent.AddChild(child)
				parent.AddBorderChild(&BorderChild{
					Position: 0, // North position
					Resource: borderChild1,
				})
				parent.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild2,
				})
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
				parent.AddChild(child)
				parent.AddBorderChild(&BorderChild{
					Position: 0, // North position
					Resource: borderChild1,
				})
				borderChild1.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild2,
				})
				borderChild2.AddBorderChild(&BorderChild{
					Position: 0, // North position
					Resource: borderChild1,
				}) // Create an indirect cycle: r1 -> r2 -> r1
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
				parent.AddChild(child)
				parent.AddBorderChild(&BorderChild{
					Position: 0, // North position
					Resource: borderChild1,
				})
				borderChild1.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild2,
				})
				borderChild2.AddBorderChild(&BorderChild{
					Position: 0, // North position
					Resource: borderChild3,
				})
				borderChild3.AddBorderChild(&BorderChild{
					Position: 8, // South position
					Resource: borderChild1,
				}) // Create a cycle: r1 -> r2 -> r3 -> r1
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
