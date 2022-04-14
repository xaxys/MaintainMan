// The MIT License (MIT)

// Copyright (c) 2014 Milan Misak

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package imagehost

import (
	"image"
	"testing"
)

func TestCalculateTopLeftPointFromGravity(t *testing.T) {
	exp := image.Point{200, 0}
	act := calculateTopLeftPointFromGravity(GravityNorth, 400, 300, 800, 600)
	if act != exp {
		t.Error("N failed", act, exp)
	}

	exp = image.Point{400, 0}
	act = calculateTopLeftPointFromGravity(GravityNorthEast, 400, 300, 800, 600)
	if act != exp {
		t.Error("NE failed", act, exp)
	}

	exp = image.Point{400, 150}
	act = calculateTopLeftPointFromGravity(GravityEast, 400, 300, 800, 600)
	if act != exp {
		t.Error("E failed", act, exp)
	}

	exp = image.Point{400, 300}
	act = calculateTopLeftPointFromGravity(GravitySouthEast, 400, 300, 800, 600)
	if act != exp {
		t.Error("SE failed", act, exp)
	}
	exp = image.Point{200, 300}
	act = calculateTopLeftPointFromGravity(GravitySouth, 400, 300, 800, 600)
	if act != exp {
		t.Error("S failed", act, exp)
	}
	exp = image.Point{0, 300}
	act = calculateTopLeftPointFromGravity(GravitySouthWest, 400, 300, 800, 600)
	if act != exp {
		t.Error("SW failed", act, exp)
	}
	exp = image.Point{0, 150}
	act = calculateTopLeftPointFromGravity(GravityWest, 400, 300, 800, 600)
	if act != exp {
		t.Error("W failed", act, exp)
	}

	exp = image.Point{0, 0}
	act = calculateTopLeftPointFromGravity(GravityNorthWest, 400, 300, 800, 600)
	if act != exp {
		t.Error("NW failed", act, exp)
	}

	exp = image.Point{200, 150}
	act = calculateTopLeftPointFromGravity(GravityCenter, 400, 300, 800, 600)
	if act != exp {
		t.Error("C failed", act, exp)
	}
}

func TestParseParameters(t *testing.T) {
	act, _ := parseParameters("w_400,h_300")
	exp := Params{400, 300, DefaultScale, DefaultCroppingMode, DefaultGravity, DefaultFilter}
	if act != exp {
		t.Errorf("Expected: %v, actual: %v", exp, act)
	}

	act, _ = parseParameters("w_200,h_300,c_k,g_c")
	exp = Params{200, 300, DefaultScale, CroppingModeKeepScale, GravityCenter, DefaultFilter}
	if act != exp {
		t.Errorf("Expected: %v, actual: %v", exp, act)
	}
}
