// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Require ${response:status} HEALTHY

func processRequire(params, options string) {
	if offline() {
		debug("REQUIRE has no effect in offline mode.")
		return
	}
	comment(echoRequireCommand, "REQUIRE: %s", params)

	left, right := split(params)
	eleft := expand(left)
	eright := expand(right)

	// handle special case here, when mere existence was required
	if len(right) == 0 {
		if len(left) != 0 && len(eleft) == 0 {
			quit("failed required condition: [%s] is not empty", left)
		}
		comment(echoProgress, "Require passed: [%s] is not empty", left)
		return
	}

	if eleft != eright {
		if lower(eleft) != lower(eright) {
			quit("failed required condition: [%s] != [%s]", eleft, eright)
		} else {
			debug("Require command succeeded only in case-insensitive comparison. [%s] and [%s]", left, right)
		}
	}
	comment(echoProgress, "Require passed: [%s] == [%s]", left, right)
}
