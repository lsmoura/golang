// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image"
	"image/color"
	"math"
	"reflect"

	"github.com/searKing/golang/go/exp/constraints"
)

type Rectangle[E constraints.Number] struct {
	Min, Max Point[E]
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rectangle[E]) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// Dx returns r's width.
func (r Rectangle[E]) Dx() E {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectangle[E]) Dy() E {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r Rectangle[E]) Size() Point[E] {
	return Point[E]{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Add returns the rectangle r translated by p.
func (r Rectangle[E]) Add(p Point[E]) Rectangle[E] {
	return Rectangle[E]{
		Point[E]{r.Min.X + p.X, r.Min.Y + p.Y},
		Point[E]{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r Rectangle[E]) Sub(p Point[E]) Rectangle[E] {
	return Rectangle[E]{
		Point[E]{r.Min.X - p.X, r.Min.Y - p.Y},
		Point[E]{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Mul returns the rectangle r*k.
func (r Rectangle[E]) Mul(k E) Rectangle[E] {
	return Rectangle[E]{r.Min.Mul(k), r.Max.Mul(k)}
}

// Div returns the rectangle p/k.
func (r Rectangle[E]) Div(k E) Rectangle[E] {
	return Rectangle[E]{r.Min.Div(k), r.Max.Div(k)}
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r Rectangle[E]) Inset(n E) Rectangle[E] {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Rectangle[E]) Intersect(s Rectangle[E]) Rectangle[E] {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	// Letting r0 and s0 be the values of r and s at the time that the method
	// is called, this next line is equivalent to:
	//
	// if max(r0.Min.X, s0.Min.X) >= min(r0.Max.X, s0.Max.X) || likewiseForY { etc }
	if r.Empty() {
		var zero Rectangle[E]
		return zero
	}
	return r
}

// Union returns the smallest rectangle that contains both r and s.
func (r Rectangle[E]) Union(s Rectangle[E]) Rectangle[E] {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

// Empty reports whether the rectangle contains no points.
func (r Rectangle[E]) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r Rectangle[E]) Eq(s Rectangle[E]) bool {
	return r == s || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rectangle[E]) Overlaps(s Rectangle[E]) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r Rectangle[E]) In(s Rectangle[E]) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

// Canon returns the canonical version of r. The returned rectangle has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (r Rectangle[E]) Canon() Rectangle[E] {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// At implements the Image interface.
func (r Rectangle[E]) At(x, y E) color.Color {
	if (Point[E]{x, y}).In(r) {
		return color.Opaque
	}
	return color.Transparent
}

// RGBA64At implements the RGBA64Image interface.
func (r Rectangle[E]) RGBA64At(x, y E) color.RGBA64 {
	if (Point[E]{x, y}).In(r) {
		return color.RGBA64{R: 0xffff, G: 0xffff, B: 0xffff, A: 0xffff}
	}
	return color.RGBA64{}
}

// Bounds implements the Image interface.
func (r Rectangle[E]) Bounds() Rectangle[E] {
	return r
}

// ColorModel implements the Image interface.
func (r Rectangle[E]) ColorModel() color.Model {
	return color.Alpha16Model
}

func (r Rectangle[E]) RoundRectangle() image.Rectangle {
	return image.Rect(round(r.Min.X), round(r.Min.Y), round(r.Max.X), round(r.Max.Y))
}

// UnionPoints returns the smallest rectangle that contains all points.
func (r Rectangle[E]) UnionPoints(pts ...Point[E]) Rectangle[E] {
	if len(pts) == 0 {
		return r
	}
	var pos int
	if r.Empty() { // an empty rectangle is an empty set, Not a point
		r.Min = pts[0]
		r.Max = pts[0]
		pos = 1
	}
	for _, p := range pts[pos:] {
		if p.X < r.Min.X {
			r.Min.X = p.X
		}
		if p.Y < r.Min.Y {
			r.Min.Y = p.Y
		}
		if p.X > r.Max.X {
			r.Max.X = p.X
		}
		if p.Y > r.Max.Y {
			r.Max.Y = p.Y
		}
	}
	return r
}

// ScaleByFactor scale rect to factor*size
func (r Rectangle[E]) ScaleByFactor(factor Point[E]) Rectangle[E] {
	if r.Empty() {
		return r
	}
	factor = factor.Sub(Pt[E](1, 1))
	minOffset := Point[E]{
		X: r.Dx() * factor.X / 2,
		Y: r.Dy() * factor.Y / 2,
	}
	maxOffset := Point[E]{
		X: r.Dx() * factor.X,
		Y: r.Dy() * factor.Y,
	}.Sub(minOffset)

	return Rectangle[E]{
		Min: Point[E]{X: r.Min.X - minOffset.X, Y: r.Min.Y - minOffset.Y},
		Max: Point[E]{X: r.Max.X + maxOffset.X, Y: r.Max.Y + maxOffset.Y},
	}
}

// Rect is shorthand for Rectangle[E]{Pt(x0, y0), Pt(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func Rect[E constraints.Number](x0, y0, x1, y1 E) Rectangle[E] {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle[E]{Point[E]{x0, y0}, Point[E]{x1, y1}}
}

func FromRectInt[E constraints.Number](rect image.Rectangle) Rectangle[E] {
	return Rect(E(rect.Min.X), E(rect.Min.Y), E(rect.Max.X), E(rect.Max.Y))
}

func round[E constraints.Number](x E) int {
	kind := reflect.TypeOf(x).Kind()
	switch kind {
	case reflect.Float32, reflect.Float64:
		return int(math.Round(float64(x)))
	}
	return int(x)
}
