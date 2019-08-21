// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Require ${response:status} HEALTHY

func processRequire(params string) {
	comment(echoRequireCommand, "REQUIRE: %s", params)

	left, right := split(params)
	eleft := expand(left)
	eright := expand(right)

	if eleft != eright {
		if lower(eleft) != lower(eright) {
			quit("failed required condition: [%s] != [%s]", eleft, eright)
		} else {
			debug("Require command succeeded only in case-insensitive comparison. [%s] and [%s]", left, right)
		}
	}
	comment(echoProgress, "Require passed: [%s] == [%s]", left, right)
}