// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package vector

import (
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	v := New(3.0, 4.0)
	if v.X != 3.0 || v.Y != 4.0 {
		t.Errorf("New(3.0, 4.0) = %v, want {3.0, 4.0}", v)
	}
}

func TestAdd(t *testing.T) {
	v1 := New(1.0, 2.0)
	v2 := New(3.0, 4.0)
	result := v1.Add(v2)
	expected := New(4.0, 6.0)
	if result != expected {
		t.Errorf("Add() = %v, want %v", result, expected)
	}
}

func TestSub(t *testing.T) {
	v1 := New(5.0, 7.0)
	v2 := New(2.0, 3.0)
	result := v1.Sub(v2)
	expected := New(3.0, 4.0)
	if result != expected {
		t.Errorf("Sub() = %v, want %v", result, expected)
	}
}

func TestScale(t *testing.T) {
	v := New(2.0, 3.0)
	result := v.Scale(2.5)
	expected := New(5.0, 7.5)
	if result != expected {
		t.Errorf("Scale(2.5) = %v, want %v", result, expected)
	}
}

func TestDot(t *testing.T) {
	v1 := New(2.0, 3.0)
	v2 := New(4.0, 5.0)
	result := v1.Dot(v2)
	expected := 23.0 // 2*4 + 3*5
	if result != expected {
		t.Errorf("Dot() = %v, want %v", result, expected)
	}
}

func TestLength(t *testing.T) {
	v := New(3.0, 4.0)
	result := v.Length()
	expected := 5.0 // sqrt(3^2 + 4^2)
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("Length() = %v, want %v", result, expected)
	}
}

func TestNormalize(t *testing.T) {
	v := New(3.0, 4.0)
	result := v.Normalize()
	expected := New(0.6, 0.8)
	if math.Abs(result.X-expected.X) > 1e-10 || math.Abs(result.Y-expected.Y) > 1e-10 {
		t.Errorf("Normalize() = %v, want %v", result, expected)
	}

	// Test zero vector
	zero := New(0.0, 0.0)
	result = zero.Normalize()
	expected = New(0.0, 0.0)
	if result != expected {
		t.Errorf("Normalize() on zero vector = %v, want %v", result, expected)
	}
}

func TestPerpendicular(t *testing.T) {
	v := New(1.0, 2.0)
	result := v.Perpendicular()
	expected := New(-2.0, 1.0)
	if result != expected {
		t.Errorf("Perpendicular() = %v, want %v", result, expected)
	}

	// Test that perpendicular is actually perpendicular (dot product = 0)
	if math.Abs(v.Dot(result)) > 1e-10 {
		t.Errorf("Perpendicular vector is not perpendicular, dot product = %v", v.Dot(result))
	}
}

func TestIsZero(t *testing.T) {
	zero := New(0.0, 0.0)
	if !zero.IsZero() {
		t.Errorf("IsZero() on zero vector = false, want true")
	}

	nonZero := New(1.0, 0.0)
	if nonZero.IsZero() {
		t.Errorf("IsZero() on non-zero vector = true, want false")
	}
}
