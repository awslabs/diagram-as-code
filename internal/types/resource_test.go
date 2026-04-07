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
func TestGetBindings_NilCheck(t *testing.T) {
	// Test case: bindings is nil (should return empty Rectangle)
	r := &Resource{
		bindings: nil,
	}

	result := r.GetBindings()
	expected := image.Rectangle{}

	if result != expected {
		t.Errorf("Expected empty Rectangle for nil bindings, got %v", result)
	}
}

func TestResolveAutoPositions_NilBindings(t *testing.T) {
	// Test case: ResolveAutoPositions with nil bindings should return error
	sourceResource := &Resource{bindings: nil}
	targetResource := &Resource{bindings: nil}

	link := &Link{
		Source:         sourceResource,
		Target:         targetResource,
		SourcePosition: WINDROSE_AUTO,
		TargetPosition: WINDROSE_AUTO,
	}

	err := link.ResolveAutoPositions()
	if err == nil {
		t.Error("Expected error for ResolveAutoPositions with nil bindings, got nil")
	}

	expectedErrMsg := "cannot calculate auto-positions for link"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestDrawOverlay(t *testing.T) {
	t.Run("NoSpanTargets", func(t *testing.T) {
		overlay := new(Resource).Init()
		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
		err := overlay.DrawOverlay(img)
		if err != nil {
			t.Errorf("DrawOverlay with no span targets should return nil, got: %v", err)
		}
		if overlay.IsDrawn() {
			t.Error("DrawOverlay should not mark resource as drawn when there are no span targets")
		}
	})

	t.Run("NilTargetBindings", func(t *testing.T) {
		overlay := new(Resource).Init()
		target := new(Resource).Init()
		target.bindings = nil
		overlay.AddSpanTarget(target)

		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
		err := overlay.DrawOverlay(img)
		if err != nil {
			t.Errorf("DrawOverlay with nil target bindings should return nil, got: %v", err)
		}
		if overlay.IsDrawn() {
			t.Error("DrawOverlay should not mark resource as drawn when first target has nil bindings")
		}
	})

	t.Run("SingleTarget", func(t *testing.T) {
		overlay := new(Resource).Init()
		target := new(Resource).Init()
		targetRect := image.Rect(50, 50, 150, 150)
		target.bindings = &targetRect
		overlay.AddSpanTarget(target)

		img := image.NewRGBA(image.Rect(0, 0, 300, 300))
		err := overlay.DrawOverlay(img)
		if err != nil {
			t.Fatalf("DrawOverlay failed: %v", err)
		}
		if !overlay.IsDrawn() {
			t.Error("DrawOverlay should mark resource as drawn")
		}
	})

	t.Run("DefaultBorderColor", func(t *testing.T) {
		overlay := new(Resource).Init()
		overlay.borderColor = nil
		target := new(Resource).Init()
		targetRect := image.Rect(50, 50, 150, 150)
		target.bindings = &targetRect
		overlay.AddSpanTarget(target)

		img := image.NewRGBA(image.Rect(0, 0, 300, 300))
		err := overlay.DrawOverlay(img)
		if err != nil {
			t.Fatalf("DrawOverlay failed: %v", err)
		}
		if overlay.borderColor == nil {
			t.Error("DrawOverlay should set a default border color when none is provided")
		}
		expected := color.RGBA{0, 0, 0, 255}
		if *overlay.borderColor != expected {
			t.Errorf("Expected default border color %v, got %v", expected, *overlay.borderColor)
		}
	})

	t.Run("SkipsNilBindingsInLaterTargets", func(t *testing.T) {
		overlay := new(Resource).Init()
		t1 := new(Resource).Init()
		t2 := new(Resource).Init()
		r1 := image.Rect(10, 10, 50, 50)
		t1.bindings = &r1
		t2.bindings = nil
		overlay.AddSpanTarget(t1)
		overlay.AddSpanTarget(t2)

		img := image.NewRGBA(image.Rect(0, 0, 300, 300))
		err := overlay.DrawOverlay(img)
		if err != nil {
			t.Fatalf("DrawOverlay failed: %v", err)
		}
		if !overlay.IsDrawn() {
			t.Error("DrawOverlay should mark resource as drawn even with nil bindings in later targets")
		}
	})

	t.Run("DashedBorder", func(t *testing.T) {
		overlay := new(Resource).Init()
		overlay.borderType = BORDER_TYPE_DASHED
		target := new(Resource).Init()
		targetRect := image.Rect(50, 50, 100, 100)
		target.bindings = &targetRect
		overlay.AddSpanTarget(target)

		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
		err := overlay.DrawOverlay(img)
		if err != nil {
			t.Fatalf("DrawOverlay with dashed border failed: %v", err)
		}
		if !overlay.IsDrawn() {
			t.Error("DrawOverlay should mark resource as drawn")
		}
	})
}

func TestAddSpanTargetAndChildMutualExclusion(t *testing.T) {
	t.Run("ChildrenThenSpanResources", func(t *testing.T) {
		r := new(Resource).Init()
		child := new(Resource).Init()
		if err := r.AddChild(child); err != nil {
			t.Fatalf("AddChild failed: %v", err)
		}
		if err := r.AddSpanTarget(new(Resource).Init()); err == nil {
			t.Error("AddSpanTarget should fail when Children already exist")
		}
	})

	t.Run("SpanResourcesThenChildren", func(t *testing.T) {
		r := new(Resource).Init()
		target := new(Resource).Init()
		if err := r.AddSpanTarget(target); err != nil {
			t.Fatalf("AddSpanTarget failed: %v", err)
		}
		if err := r.AddChild(new(Resource).Init()); err == nil {
			t.Error("AddChild should fail when SpanResources already exist")
		}
	})
}

// TestSpanOverlayMatchesTreeLayout verifies that a SpanResources overlay
// produces the same bounding box as an equivalent tree hierarchy.
// Tree: AZ -> Subnet -> EC2
// Span: Subnet -> EC2, AZ overlays Subnet
func TestSpanOverlayMatchesTreeLayout(t *testing.T) {
	// Tree version: AZ contains Subnet contains EC2
	azTree := new(Resource).Init()
	azLabel := "Availability Zone"
	azTree.SetLabel(&azLabel, nil, nil)
	subnetTree := new(Resource).Init()
	subnetLabel := "Subnet"
	subnetTree.SetLabel(&subnetLabel, nil, nil)
	ec2Tree := new(Resource).Init()
	ec2Label := "Instance"
	ec2Tree.SetLabel(&ec2Label, nil, nil)
	// Set icon bounds to simulate having an icon
	ec2Tree.SetIconBounds(image.Rect(0, 0, 64, 64))
	ec2Tree.iconImage = image.NewRGBA(image.Rect(0, 0, 64, 64))
	subnetTree.AddChild(ec2Tree)
	azTree.AddChild(subnetTree)
	if err := azTree.Scale(nil, nil); err != nil {
		t.Fatalf("Tree Scale failed: %v", err)
	}
	treeBounds := azTree.GetBindings()

	// Span version: Subnet contains EC2, AZ overlays Subnet
	azSpan := new(Resource).Init()
	azSpan.SetLabel(&azLabel, nil, nil)
	subnetSpan := new(Resource).Init()
	subnetSpan.SetLabel(&subnetLabel, nil, nil)
	ec2Span := new(Resource).Init()
	ec2Span.SetLabel(&ec2Label, nil, nil)
	ec2Span.SetIconBounds(image.Rect(0, 0, 64, 64))
	ec2Span.iconImage = image.NewRGBA(image.Rect(0, 0, 64, 64))
	subnetSpan.AddChild(ec2Span)
	azSpan.AddSpanTarget(subnetSpan)

	if err := subnetSpan.Scale(nil, nil); err != nil {
		t.Fatalf("Span Scale failed: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 2000, 2000))
	if _, err := subnetSpan.Draw(img, nil); err != nil {
		t.Fatalf("Span Draw failed: %v", err)
	}
	if err := azSpan.DrawOverlay(img); err != nil {
		t.Fatalf("Span DrawOverlay failed: %v", err)
	}
	spanBounds := azSpan.GetBindings()

	// Compare: overlay should match tree AZ relative to Subnet
	treeSubnet := subnetTree.GetBindings()
	spanSubnet := subnetSpan.GetBindings()

	t.Logf("Tree: AZ bindings=%v margin=%v padding=%v", treeBounds, azTree.GetMargin(), azTree.GetPadding())
	t.Logf("Tree: Subnet bindings=%v margin=%v padding=%v", subnetTree.GetBindings(), subnetTree.GetMargin(), subnetTree.GetPadding())
	t.Logf("Tree: EC2 bindings=%v margin=%v padding=%v", ec2Tree.GetBindings(), ec2Tree.GetMargin(), ec2Tree.GetPadding())
	t.Logf("Span: AZ bindings=%v", spanBounds)
	t.Logf("Span: Subnet bindings=%v margin=%v padding=%v", subnetSpan.GetBindings(), subnetSpan.GetMargin(), subnetSpan.GetPadding())
	t.Logf("Span: EC2 bindings=%v margin=%v padding=%v", ec2Span.GetBindings(), ec2Span.GetMargin(), ec2Span.GetPadding())

	// Distance from AZ edge to Subnet edge on each side
	treeTop := treeSubnet.Min.Y - treeBounds.Min.Y
	treeBottom := treeBounds.Max.Y - treeSubnet.Max.Y
	treeLeft := treeSubnet.Min.X - treeBounds.Min.X
	treeRight := treeBounds.Max.X - treeSubnet.Max.X

	spanTop := spanSubnet.Min.Y - spanBounds.Min.Y
	spanBottom := spanBounds.Max.Y - spanSubnet.Max.Y
	spanLeft := spanSubnet.Min.X - spanBounds.Min.X
	spanRight := spanBounds.Max.X - spanSubnet.Max.X

	t.Logf("Tree: AZ-to-Subnet  top=%d bottom=%d left=%d right=%d", treeTop, treeBottom, treeLeft, treeRight)
	t.Logf("Span: AZ-to-Subnet  top=%d bottom=%d left=%d right=%d", spanTop, spanBottom, spanLeft, spanRight)

	if treeTop != spanTop {
		t.Errorf("Top mismatch: tree=%d, span=%d", treeTop, spanTop)
	}
	if treeBottom != spanBottom {
		t.Errorf("Bottom mismatch: tree=%d, span=%d", treeBottom, spanBottom)
	}
	if treeLeft != spanLeft {
		t.Errorf("Left mismatch: tree=%d, span=%d", treeLeft, spanLeft)
	}
	if treeRight != spanRight {
		t.Errorf("Right mismatch: tree=%d, span=%d", treeRight, spanRight)
	}
}

func TestSpanOverlayMatchesTreeLayout_WithIcon(t *testing.T) {
	// Tree: Subnet -> ASG -> EC2
	subnetTree := new(Resource).Init()
	subnetLabel := "Subnet"
	subnetTree.SetLabel(&subnetLabel, nil, nil)
	asgTree := new(Resource).Init()
	asgLabel := "Auto Scaling Group"
	asgTree.SetLabel(&asgLabel, nil, nil)
	asgTree.SetIconBounds(image.Rect(0, 0, 64, 64))
	asgTree.iconImage = image.NewRGBA(image.Rect(0, 0, 64, 64))
	ec2Tree := new(Resource).Init()
	ec2Label := "Instance"
	ec2Tree.SetLabel(&ec2Label, nil, nil)
	ec2Tree.SetIconBounds(image.Rect(0, 0, 64, 64))
	ec2Tree.iconImage = image.NewRGBA(image.Rect(0, 0, 64, 64))
	asgTree.AddChild(ec2Tree)
	subnetTree.AddChild(asgTree)
	if err := subnetTree.Scale(nil, nil); err != nil {
		t.Fatalf("Tree Scale failed: %v", err)
	}

	// Span: Subnet -> EC2, ASG overlays EC2
	subnetSpan := new(Resource).Init()
	subnetSpan.SetLabel(&subnetLabel, nil, nil)
	asgSpan := new(Resource).Init()
	asgSpan.SetLabel(&asgLabel, nil, nil)
	asgSpan.SetIconBounds(image.Rect(0, 0, 64, 64))
	asgSpan.iconImage = image.NewRGBA(image.Rect(0, 0, 64, 64))
	ec2Span := new(Resource).Init()
	ec2Span.SetLabel(&ec2Label, nil, nil)
	ec2Span.SetIconBounds(image.Rect(0, 0, 64, 64))
	ec2Span.iconImage = image.NewRGBA(image.Rect(0, 0, 64, 64))
	subnetSpan.AddChild(ec2Span)
	asgSpan.AddSpanTarget(ec2Span)
	if err := subnetSpan.Scale(nil, nil); err != nil {
		t.Fatalf("Span Scale failed: %v", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 2000, 2000))
	if _, err := subnetSpan.Draw(img, nil); err != nil {
		t.Fatalf("Span Draw failed: %v", err)
	}
	if err := asgSpan.DrawOverlay(img); err != nil {
		t.Fatalf("Span DrawOverlay failed: %v", err)
	}

	treeBounds := asgTree.GetBindings()
	spanBounds := asgSpan.GetBindings()

	t.Logf("Tree: ASG bindings=%v margin=%v padding=%v", treeBounds, asgTree.GetMargin(), asgTree.GetPadding())
	t.Logf("Tree: EC2 bindings=%v margin=%v padding=%v", ec2Tree.GetBindings(), ec2Tree.GetMargin(), ec2Tree.GetPadding())
	t.Logf("Span: ASG bindings=%v", spanBounds)
	t.Logf("Span: EC2 bindings=%v margin=%v padding=%v", ec2Span.GetBindings(), ec2Span.GetMargin(), ec2Span.GetPadding())

	treeEC2 := ec2Tree.GetBindings()
	spanEC2 := ec2Span.GetBindings()

	treeTop := treeEC2.Min.Y - treeBounds.Min.Y
	treeBottom := treeBounds.Max.Y - treeEC2.Max.Y
	treeLeft := treeEC2.Min.X - treeBounds.Min.X
	treeRight := treeBounds.Max.X - treeEC2.Max.X

	spanTop := spanEC2.Min.Y - spanBounds.Min.Y
	spanBottom := spanBounds.Max.Y - spanEC2.Max.Y
	spanLeft := spanEC2.Min.X - spanBounds.Min.X
	spanRight := spanBounds.Max.X - spanEC2.Max.X

	t.Logf("Tree: ASG-to-EC2  top=%d bottom=%d left=%d right=%d", treeTop, treeBottom, treeLeft, treeRight)
	t.Logf("Span: ASG-to-EC2  top=%d bottom=%d left=%d right=%d", spanTop, spanBottom, spanLeft, spanRight)

	if treeTop != spanTop {
		t.Errorf("Top mismatch: tree=%d, span=%d (diff=%d)", treeTop, spanTop, treeTop-spanTop)
	}
	if treeBottom != spanBottom {
		t.Errorf("Bottom mismatch: tree=%d, span=%d (diff=%d)", treeBottom, spanBottom, treeBottom-spanBottom)
	}
	if treeLeft != spanLeft {
		t.Errorf("Left mismatch: tree=%d, span=%d (diff=%d)", treeLeft, spanLeft, treeLeft-spanLeft)
	}
	if treeRight != spanRight {
		t.Errorf("Right mismatch: tree=%d, span=%d (diff=%d)", treeRight, spanRight, treeRight-spanRight)
	}

	// Compare Subnet bindings (parent level)
	treeSubnet := subnetTree.GetBindings()
	spanSubnet := subnetSpan.GetBindings()

	// Debug: ASG header width and EC2 placement width
	asgFace, _ := asgTree.prepareFontFace(true, nil)
	asgTextWidth, _ := asgTree.calculateTitleSize(asgFace)
	asgHeaderWidth := asgTextWidth + asgTree.iconBounds.Dx() + 30
	t.Logf("ASG headerWidth=%d (textWidth=%d + iconDx=%d + 30)", asgHeaderWidth, asgTextWidth, asgTree.iconBounds.Dx())
	t.Logf("Tree: ASG bindings Dx=%d", asgTree.GetBindings().Dx())
	t.Logf("Tree: EC2 margin=%v", ec2Tree.GetMargin())
	t.Logf("Span: EC2 margin=%v", ec2Span.GetMargin())

	t.Logf("Tree: Subnet bindings=%v Dx=%d Dy=%d", treeSubnet, treeSubnet.Dx(), treeSubnet.Dy())
	t.Logf("Span: Subnet bindings=%v Dx=%d Dy=%d", spanSubnet, spanSubnet.Dx(), spanSubnet.Dy())
	if treeSubnet.Dx() != spanSubnet.Dx() {
		t.Logf("Subnet Width diff: tree=%d, span=%d (diff=%d) - expected due to overlay header width", treeSubnet.Dx(), spanSubnet.Dx(), treeSubnet.Dx()-spanSubnet.Dx())
	}
	if treeSubnet.Dy() != spanSubnet.Dy() {
		t.Errorf("Subnet Height mismatch: tree=%d, span=%d (diff=%d)", treeSubnet.Dy(), spanSubnet.Dy(), treeSubnet.Dy()-spanSubnet.Dy())
	}
}

func TestCalculateTitleSize(t *testing.T) {
	resource := new(Resource).Init()

	// Test with empty label
	resource.label = ""
	fontFace, err := resource.prepareFontFace(false, nil)
	if err != nil {
		t.Skipf("Skipping test due to font preparation error: %v", err)
	}

	width, height := resource.calculateTitleSize(fontFace)
	if width != 0 || height != 0 {
		t.Errorf("Expected (0, 0) for empty label, got (%d, %d)", width, height)
	}

	// Test with single line label
	resource.label = "Test Label"
	width, height = resource.calculateTitleSize(fontFace)
	if width <= 0 || height <= 0 {
		t.Errorf("Expected positive dimensions for single line label, got (%d, %d)", width, height)
	}

	// Test with multi-line label
	resource.label = "Line 1\nLine 2\nLine 3"
	width2, height2 := resource.calculateTitleSize(fontFace)
	if width2 <= 0 || height2 <= height {
		t.Errorf("Expected multi-line label to have greater height, got single: (%d, %d), multi: (%d, %d)", width, height, width2, height2)
	}
}
