// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package vector

import (
	"math"
)

// Vector represents a 2D vector
type Vector struct {
	X, Y float64
}

// New creates a new vector
func New(x, y float64) Vector {
	return Vector{X: x, Y: y}
}

// Add returns the sum of two vectors
func (v Vector) Add(other Vector) Vector {
	return Vector{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub returns the difference of two vectors
func (v Vector) Sub(other Vector) Vector {
	return Vector{X: v.X - other.X, Y: v.Y - other.Y}
}

// Scale returns the vector scaled by a factor
func (v Vector) Scale(factor float64) Vector {
	return Vector{X: v.X * factor, Y: v.Y * factor}
}

// Dot returns the dot product of two vectors
func (v Vector) Dot(other Vector) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Length returns the magnitude of the vector
func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize returns the unit vector in the same direction
func (v Vector) Normalize() Vector {
	length := v.Length()
	if length == 0 {
		return Vector{X: 0, Y: 0}
	}
	return Vector{X: v.X / length, Y: v.Y / length}
}

// Perpendicular returns the perpendicular vector (90 degrees counterclockwise)
func (v Vector) Perpendicular() Vector {
	return Vector{X: -v.Y, Y: v.X}
}

// IsZero checks if the vector is zero
func (v Vector) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

// DecomposeXY splits a vector into its X and Y components
func (v Vector) DecomposeXY() (Vector, Vector) {
	xComponent := New(v.X, 0)
	yComponent := New(0, v.Y)
	return xComponent, yComponent
}
