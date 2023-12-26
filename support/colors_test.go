// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package support

import (
	"testing"
)

func TestColors(t *testing.T) {

	// "darkolivegreen2":      "#bcee68",
	if value, err := GetColor("darkolivegreen2"); err != nil {
		t.Fatal("should've worked")
	} else if value != "#bcee68" {
		t.Fatal("got wrong color value")
	}
}

func TestReverseColors(t *testing.T) {

	name := GetColorName("#bcee68")
	if name != "darkolivegreen2" {
		t.Fatalf("got wrong color name (%s)", name)
	}
}
